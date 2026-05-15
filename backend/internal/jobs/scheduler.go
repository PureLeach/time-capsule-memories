package jobs

import (
	"fmt"
	"log/slog"
	"time"
	"time_capsule_memories/internal/config"

	"github.com/robfig/cron/v3"
)

// StartScheduler initialises the cron scheduler, registers the capsule
// dispatch job, and starts ticking. The returned *cron.Cron lets the caller
// drain in-flight jobs at shutdown via c.Stop().
func StartScheduler() (*cron.Cron, error) {
	cfg := config.GetConfig()

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, fmt.Errorf("load timezone: %w", err)
	}

	c := cron.New(cron.WithLocation(loc))

	if _, err := c.AddFunc(cfg.CronCapsuleDispatch, JobCapsule); err != nil {
		return nil, fmt.Errorf("schedule capsule dispatch: %w", err)
	}

	c.Start()
	slog.Info("scheduler started", "job", "capsule_dispatch", "spec", cfg.CronCapsuleDispatch)

	return c, nil
}
