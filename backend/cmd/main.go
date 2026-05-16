// @title Time Capsule Memories API
// @version 1.0
// @description REST API backend for the Time Capsule Memories project.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @license.name MIT
// @host backend.localhost
// @BasePath /

package main

import (
	"context"
	"os/signal"
	"syscall"

	"time_capsule_memories/internal/app"
	"time_capsule_memories/internal/logging"

	_ "time_capsule_memories/docs" // Swagger docs
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
