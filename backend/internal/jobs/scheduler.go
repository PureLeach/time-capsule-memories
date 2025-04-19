package jobs

import (
	"log"
	"time"
	"time_capsule_memories/internal/config"

	"github.com/robfig/cron/v3"
)

// StartScheduler initializes and starts the cron scheduler for recurring tasks.
func StartScheduler() {
	// Fetch the configuration settings
	cfg := config.GetConfig()

	// Load the timezone (e.g., UTC). Log an error and exit if it fails.
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err) // Use log.Fatalf for immediate exit on critical errors
	}

	// Create a new cron scheduler with the specified timezone
	c := cron.New(cron.WithLocation(loc))

	// Add the capsule dispatch job with the cron expression from the configuration
	_, err = c.AddFunc(cfg.CronCapsuleDispatch, func() {
		JobCapsule() // Execute the capsule dispatch job
	})
	if err != nil {
		log.Fatalf("Failed to schedule capsule dispatch job: %v", err)
	}

	// Start the scheduler in a separate goroutine for asynchronous execution
	go func() {
		log.Println("[CRON] Scheduler started successfully.")
		c.Start()
	}()

	// Ensure the scheduler is stopped when the program exits
	defer c.Stop()
}
