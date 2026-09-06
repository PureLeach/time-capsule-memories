// Command time-capsule-memories serves the HTTP API and runs the capsule dispatcher.
//
// @title Time Capsule Memories API
// @version 1.0
// @description Schedules messages, with optional image attachments, for delivery by email on a future date.
// @contact.name Source and issues
// @contact.url https://github.com/MaxBarannikov/time-capsule-memories
// @host backend.localhost
// @BasePath /
package main

import (
	"context"
	"os/signal"
	"syscall"

	"time_capsule_memories/internal/app"
	"time_capsule_memories/internal/logging"

	_ "time_capsule_memories/docs" // generated spec, served at /swagger
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.New(ctx)
	if err != nil {
		logging.Fatal("application init failed", "error", err)
	}

	if err := a.Run(ctx); err != nil {
		logging.Fatal("application stopped with error", "error", err)
	}
}
