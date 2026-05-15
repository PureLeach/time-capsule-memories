package jobs

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"time_capsule_memories/internal/repository"
	"time_capsule_memories/internal/services"
)

const capsuleDispatchTimeout = 2 * time.Minute

const capsuleClaimLimit = 100

func JobCapsule() {
	startTime := time.Now()
	slog.Info("capsule dispatch started", "started_at", startTime.Format(time.RFC3339))

	capsules, err := repository.ClaimDueCapsules(context.Background(), capsuleClaimLimit)
	if err != nil {
		slog.Error("failed to claim due capsules", "error", err)
		return
	}

	if len(capsules) == 0 {
		slog.Info("no capsules due")
		return
	}

	slog.Info("capsules claimed", "count", len(capsules))

	var wg sync.WaitGroup
	for _, capsule := range capsules {
		capsule := capsule
		wg.Add(1)

		go func(capsuleID string) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), capsuleDispatchTimeout)
			defer cancel()

			if err := services.ProcessCapsule(ctx, capsule); err != nil {
				slog.Error("failed to process capsule", "capsule_id", capsuleID, "error", err)
			} else {
				slog.Info("capsule processed", "capsule_id", capsuleID)
			}
		}(strconv.Itoa(capsule.ID))
	}

	wg.Wait()

	slog.Info("capsule dispatch finished", "finished_at", time.Now().Format(time.RFC3339))
}
