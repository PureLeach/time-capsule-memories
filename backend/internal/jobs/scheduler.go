package jobs

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

// StartScheduler runs the dispatcher on the given cron spec, in UTC so a host's
// timezone cannot shift delivery dates.
func StartScheduler(spec string, d *Dispatcher) (*cron.Cron, error) {
	c := cron.New(cron.WithLocation(time.UTC))

	if _, err := c.AddFunc(spec, d.Run); err != nil {
		return nil, fmt.Errorf("schedule capsule dispatch %q: %w", spec, err)
	}

	c.Start()
	slog.Info("scheduler started", "job", "capsule_dispatch", "spec", spec)

	return c, nil
}
