package jobs

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

func StartScheduler(spec string, d *Dispatcher) (*cron.Cron, error) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, fmt.Errorf("load timezone: %w", err)
	}

	c := cron.New(cron.WithLocation(loc))

	if _, err := c.AddFunc(spec, d.Run); err != nil {
		return nil, fmt.Errorf("schedule capsule dispatch: %w", err)
	}

	c.Start()
	slog.Info("scheduler started", "job", "capsule_dispatch", "spec", spec)

	return c, nil
}
