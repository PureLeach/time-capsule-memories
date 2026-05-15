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
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"time_capsule_memories/internal/config"
	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/jobs"
	"time_capsule_memories/internal/logging"
	appmw "time_capsule_memories/internal/middleware"
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "time_capsule_memories/docs" // Swagger docs

	echoSwagger "github.com/swaggo/echo-swagger"
)

// shutdownTimeout bounds the whole graceful-shutdown phase: HTTP drain, cron
// drain and pool close share this budget. A capsule send can in theory run up
// to capsuleDispatchTimeout (~2 min), but we'd rather force-exit and let the
// next tick retry from the `waiting` row than block the orchestrator.
const shutdownTimeout = 30 * time.Second

func main() {
	cfg := config.GetConfig()
	logging.Init(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := database.Connect(); err != nil {
		logging.Fatal("database connection failed", "error", err)
	}
	defer database.Close()

	minioCtx, cancelMinio := context.WithTimeout(ctx, 10*time.Second)
	if err := minio_client.MinioInit(minioCtx); err != nil {
		cancelMinio()
		logging.Fatal("minio initialisation failed", "error", err)
	}
	cancelMinio()

	cronScheduler, err := jobs.StartScheduler()
	if err != nil {
		logging.Fatal("scheduler start failed", "error", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Order matters: request id first so every later layer can log it; recover
	// next so panics in handlers don't bypass the access log; access log wraps
	// the handler; body limit and CORS sit closest to the handler. Probes are
	// excluded from the access log so they don't drown out real traffic.
	e.Use(appmw.RequestID())
	e.Use(appmw.Recover())
	e.Use(appmw.AccessLog("/healthz", "/readyz"))
	e.Use(middleware.BodyLimit("1M"))
	e.Use(appmw.CORSConfig(cfg.CORSAllowedOrigins))

	routes.RegisterHealthRoutes(e)
	routes.RegisterFileRoutes(e)
	routes.RegisterCapsuleRoutes(e)
	routes.RegisterFeedbackRoutes(e)
	routes.RegisterEmailRoutes(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting HTTP server", "addr", ":8000")
		if err := e.Start(":8000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			logging.Fatal("http server stopped", "error", err)
		}
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("http server shutdown error", "error", err)
	} else {
		slog.Info("http server stopped")
	}

	cronDone := cronScheduler.Stop()
	select {
	case <-cronDone.Done():
		slog.Info("scheduler drained")
	case <-shutdownCtx.Done():
		slog.Warn("scheduler drain timed out; capsules in flight may resume on next tick")
	}
}
