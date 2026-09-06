// Package jobs contains the scheduled background work.
package jobs

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"time_capsule_memories/internal/models"
)

const (
	capsuleDispatchTimeout = 2 * time.Minute
	capsuleClaimLimit      = 100
	// Keeps the fan-out below the connection limit of a typical SMTP relay.
	capsuleSendConcurrency = 8
)

type capsuleClaimer interface {
	ClaimDue(ctx context.Context, limit int) ([]*models.CapsuleResponse, error)
}

type capsuleProcessor interface {
	Process(ctx context.Context, capsule *models.CapsuleResponse) error
}

type Dispatcher struct {
	repo    capsuleClaimer
	service capsuleProcessor
}

func NewDispatcher(repo capsuleClaimer, service capsuleProcessor) *Dispatcher {
	return &Dispatcher{repo: repo, service: service}
}

// Run claims the capsules that are due and delivers them. Failures are logged
// rather than returned: cron has nowhere to put an error, and returning failed
// rows to 'waiting' is the capsule service's job.
func (d *Dispatcher) Run() {
	started := time.Now()

	capsules, err := d.repo.ClaimDue(context.Background(), capsuleClaimLimit)
	if err != nil {
		slog.Error("failed to claim due capsules", "error", err)
		return
	}

	if len(capsules) == 0 {
		slog.Debug("no capsules due")
		return
	}

	slog.Info("capsules claimed", "count", len(capsules))

	var (
		wg   sync.WaitGroup
		slot = make(chan struct{}, capsuleSendConcurrency)
	)
	for _, capsule := range capsules {
		wg.Add(1)
		slot <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-slot }()

			ctx, cancel := context.WithTimeout(context.Background(), capsuleDispatchTimeout)
			defer cancel()

			if err := d.service.Process(ctx, capsule); err != nil {
				slog.Error("failed to process capsule", "capsule_id", capsule.ID, "error", err)
				return
			}
			slog.Info("capsule processed", "capsule_id", capsule.ID)
		}()
	}

	wg.Wait()

	slog.Info("capsule dispatch finished",
		"count", len(capsules),
		"duration_ms", time.Since(started).Milliseconds(),
	)
}
