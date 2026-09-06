// Package app wires every component together and owns the process lifecycle.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"time_capsule_memories/internal/config"
	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/handlers"
	"time_capsule_memories/internal/jobs"
	"time_capsule_memories/internal/logging"
	appmw "time_capsule_memories/internal/middleware"
	"time_capsule_memories/internal/repository"
	"time_capsule_memories/internal/routes"
	"time_capsule_memories/internal/services"
	"time_capsule_memories/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"

	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	httpAddr        = ":8000"
	readHeaderLimit = 10 * time.Second
	storageTimeout  = 10 * time.Second
	shutdownTimeout = 30 * time.Second

	// Generous next to the 4096-character message cap; attachments never pass
	// through this service.
	requestBodyLimit = "1M"
)

type App struct {
	pool *pgxpool.Pool
	cron *cron.Cron
	echo *echo.Echo
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	logging.Init(cfg.LogLevel)

	pool, err := database.New(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	storageCtx, cancel := context.WithTimeout(ctx, storageTimeout)
	defer cancel()
	store, err := storage.New(storageCtx, cfg)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("init object store: %w", err)
	}

	capsuleRepo := repository.NewCapsule(pool)
	feedbackRepo := repository.NewFeedback(pool)
	mailer := services.NewSMTPMailer(cfg)
	capsuleService := services.NewCapsuleService(capsuleRepo, store, mailer)

	dispatcher := jobs.NewDispatcher(capsuleRepo, capsuleService)
	cronScheduler, err := jobs.StartScheduler(cfg.CronCapsuleDispatch, dispatcher)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("start scheduler: %w", err)
	}

	handler := handlers.New(pool, store, capsuleRepo, feedbackRepo, mailer)

	return &App{
		pool: pool,
		cron: cronScheduler,
		echo: newRouter(cfg, handler),
	}, nil
}

func newRouter(cfg *config.Config, handler *handlers.Handler) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Server.ReadHeaderTimeout = readHeaderLimit

	// Order matters: request id first so later layers can log it, recover next
	// so panics still reach the access log.
	e.Use(appmw.RequestID())
	e.Use(appmw.Recover())
	e.Use(appmw.AccessLog("/healthz", "/readyz"))
	e.Use(middleware.BodyLimit(requestBodyLimit))
	e.Use(appmw.CORSConfig(cfg.CORSAllowedOrigins))

	routes.Register(e, handler)
	if cfg.EnableTestEmailEndpoint {
		slog.Warn("test email endpoint enabled; it is unauthenticated and must not be exposed publicly")
		routes.RegisterTestEmail(e, handler)
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}

// Run serves until ctx is canceled, then drains the scheduler so a capsule in
// flight is not abandoned mid-delivery.
func (a *App) Run(ctx context.Context) error {
	defer a.pool.Close()

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting HTTP server", "addr", httpAddr)
		if err := a.echo.Start(httpAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.echo.Shutdown(shutdownCtx); err != nil {
		slog.Error("http server shutdown error", "error", err)
	} else {
		slog.Info("http server stopped")
	}

	cronDone := a.cron.Stop()
	select {
	case <-cronDone.Done():
		slog.Info("scheduler drained")
	case <-shutdownCtx.Done():
		slog.Warn("scheduler drain timed out; capsules in flight may resume on next tick")
	}

	return nil
}
