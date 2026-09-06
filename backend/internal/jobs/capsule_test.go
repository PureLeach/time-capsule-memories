package jobs

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"time_capsule_memories/internal/models"

	"github.com/stretchr/testify/require"
)

type stubClaimer struct {
	capsules []*models.CapsuleResponse
	err      error
	limit    int
}

func (s *stubClaimer) ClaimDue(_ context.Context, limit int) ([]*models.CapsuleResponse, error) {
	s.limit = limit
	return s.capsules, s.err
}

type recordingProcessor struct {
	mu        sync.Mutex
	processed []int
	failOn    map[int]bool

	inFlight atomic.Int32
	peak     atomic.Int32
	block    chan struct{}
}

func (r *recordingProcessor) Process(_ context.Context, capsule *models.CapsuleResponse) error {
	current := r.inFlight.Add(1)
	for {
		peak := r.peak.Load()
		if current <= peak || r.peak.CompareAndSwap(peak, current) {
			break
		}
	}
	if r.block != nil {
		<-r.block
	}
	defer r.inFlight.Add(-1)

	r.mu.Lock()
	r.processed = append(r.processed, capsule.ID)
	r.mu.Unlock()

	if r.failOn[capsule.ID] {
		return errors.New("boom")
	}
	return nil
}

func capsules(ids ...int) []*models.CapsuleResponse {
	out := make([]*models.CapsuleResponse, 0, len(ids))
	for _, id := range ids {
		out = append(out, &models.CapsuleResponse{ID: id, Status: models.StatusInProgress})
	}
	return out
}

func TestRun_ProcessesEveryClaimedCapsule(t *testing.T) {
	claimer := &stubClaimer{capsules: capsules(1, 2, 3)}
	processor := &recordingProcessor{}

	NewDispatcher(claimer, processor).Run()

	require.ElementsMatch(t, []int{1, 2, 3}, processor.processed)
	require.Equal(t, capsuleClaimLimit, claimer.limit)
}

func TestRun_OneFailureDoesNotStopTheBatch(t *testing.T) {
	claimer := &stubClaimer{capsules: capsules(1, 2, 3)}
	processor := &recordingProcessor{failOn: map[int]bool{2: true}}

	NewDispatcher(claimer, processor).Run()

	require.ElementsMatch(t, []int{1, 2, 3}, processor.processed)
}

func TestRun_SkipsWorkWhenNothingIsDue(t *testing.T) {
	processor := &recordingProcessor{}
	NewDispatcher(&stubClaimer{}, processor).Run()
	require.Empty(t, processor.processed)
}

func TestRun_DoesNotProcessOnClaimFailure(t *testing.T) {
	processor := &recordingProcessor{}
	NewDispatcher(&stubClaimer{err: errors.New("db down")}, processor).Run()
	require.Empty(t, processor.processed)
}

func TestRun_RespectsConcurrencyLimit(t *testing.T) {
	ids := make([]int, 0, capsuleSendConcurrency*4)
	for i := 1; i <= capsuleSendConcurrency*4; i++ {
		ids = append(ids, i)
	}

	processor := &recordingProcessor{block: make(chan struct{})}
	claimer := &stubClaimer{capsules: capsules(ids...)}

	done := make(chan struct{})
	go func() {
		NewDispatcher(claimer, processor).Run()
		close(done)
	}()

	require.Eventually(t, func() bool {
		return processor.peak.Load() == int32(capsuleSendConcurrency)
	}, time.Second, time.Millisecond, "workers should saturate the limit")

	close(processor.block)
	<-done

	require.Len(t, processor.processed, len(ids))
	require.LessOrEqual(t, processor.peak.Load(), int32(capsuleSendConcurrency))
}
