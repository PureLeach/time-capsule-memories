package jobs

import (
	"log/slog"
	"strconv"
	"sync"
	"time"

	"time_capsule_memories/internal/repository"
	"time_capsule_memories/internal/services"
)

// JobCapsule is a scheduled task that processes capsules scheduled for today's delivery date.
func JobCapsule() {
	// Record the job start time for monitoring purposes
	startTime := time.Now()
	slog.Info("capsule dispatch started", "started_at", startTime.Format(time.RFC3339))

	// Fetch the capsules that are scheduled for delivery today
	capsules, err := repository.GetCapsulesByToday()
	if err != nil {
		slog.Error("failed to fetch due capsules", "error", err)
		return
	}

	// If no capsules are found, log it and exit the function
	if len(capsules) == 0 {
		slog.Info("no capsules due")
		return
	}

	slog.Info("capsules due", "count", len(capsules))

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Iterate over each capsule to process them concurrently
	for _, capsule := range capsules {
		capsule := capsule // Copy the capsule to avoid race conditions in goroutines
		wg.Add(1)

		go func(capsuleID string) {
			defer wg.Done()

			if err := services.ProcessCapsule(capsule); err != nil {
				slog.Error("failed to process capsule", "capsule_id", capsuleID, "error", err)
			} else {
				slog.Info("capsule processed", "capsule_id", capsuleID)
			}
		}(strconv.Itoa(capsule.ID))
	}

	// Wait for all goroutines to finish before logging job completion
	wg.Wait()

	slog.Info("capsule dispatch finished", "finished_at", time.Now().Format(time.RFC3339))
}
