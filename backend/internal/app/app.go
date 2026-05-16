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
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/repository"
	"time_capsule_memories/internal/routes"
	"time_capsule_memories/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"

	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	httpAddr         = ":8000"
	minioInitTimeout = 10 * time.Second
	shutdownTimeout  = 30 * time.Second
)

type App struct {
	cfg  *config.Config
	pool *pgxpool.Pool
	cron *cron.Cron
	echo *echo.Echo
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	logging.Init(cfg.LogLevel)

	pool, err := database.New(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	minioCtx, cancel := context.WithTimeout(ctx, minioInitTimeout)
	defer cancel()
	store, err := minio_client.New(minioCtx, cfg)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("init minio: %w", err)
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

	routes.RegisterHealthRoutes(e, handler)
	routes.RegisterFileRoutes(e, handler)
	routes.RegisterCapsuleRoutes(e, handler)
	routes.RegisterFeedbackRoutes(e, handler)
	routes.RegisterEmailRoutes(e, handler)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return &App{
		cfg:  cfg,
		pool: pool,
		cron: cronScheduler,
		echo: e,
	}, nil
}

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
