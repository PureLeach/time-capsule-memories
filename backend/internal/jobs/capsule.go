package jobs

import (
	"log"
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
	log.Printf("[CRON] Capsule dispatch job started at %s", startTime.Format(time.RFC3339))

	// Fetch the capsules that are scheduled for delivery today
	capsules, err := repository.GetCapsulesByToday()
	if err != nil {
		// Log the error and exit early if fetching capsules failed
		log.Printf("[CRON] Failed to fetch capsules: %v", err)
		return
	}

	// If no capsules are found, log it and exit the function
	if len(capsules) == 0 {
		log.Println("[CRON] No capsules to process")
		return
	}

	// Log the number of capsules to process
	log.Printf("[CRON] Found %d capsules to process", len(capsules))

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Iterate over each capsule to process them concurrently
	for _, capsule := range capsules {
		capsule := capsule // Copy the capsule to avoid race conditions in goroutines
		wg.Add(1)

		go func(capsuleID string) {
			defer wg.Done()

			// Process each capsule using the services layer
			if err := services.ProcessCapsule(capsule); err != nil {
				// Log the failure for each capsule
				log.Printf("[CRON] Failed to process capsule ID %s: %v", capsuleID, err)
			} else {
				// Log successful processing of the capsule
				log.Printf("[CRON] Successfully processed capsule ID %s", capsuleID)
			}
		}(strconv.Itoa(capsule.ID)) // Convert capsule ID to string to match the function argument type
	}

	// Wait for all goroutines to finish before logging job completion
	wg.Wait()

	// Log the job completion time
	log.Printf("[CRON] Capsule dispatch job finished at %s", time.Now().Format(time.RFC3339))
}
