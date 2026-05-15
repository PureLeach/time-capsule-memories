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
	"log/slog"

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

func main() {
	cfg := config.GetConfig()
	logging.Init(cfg.LogLevel)

	if err := database.Connect(); err != nil {
		logging.Fatal("database connection failed", "error", err)
	}
	defer database.Close()

	minio_client.MinioInit()

	jobs.StartScheduler()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Order matters: request id first so every later layer can log it; recover
	// next so panics in handlers don't bypass the access log; access log wraps
	// the handler; body limit and CORS sit closest to the handler.
	e.Use(appmw.RequestID())
	e.Use(appmw.Recover())
	e.Use(appmw.AccessLog())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(appmw.CORSConfig(cfg.CORSAllowedOrigins))

	routes.RegisterFileRoutes(e)
	routes.RegisterCapsuleRoutes(e)
	routes.RegisterFeedbackRoutes(e)
	routes.RegisterEmailRoutes(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	slog.Info("starting HTTP server", "addr", ":8000")
	if err := e.Start(":8000"); err != nil {
		logging.Fatal("http server stopped", "error", err)
	}
}
