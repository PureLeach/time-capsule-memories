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
	"time_capsule_memories/internal/middleware"
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/routes"

	"github.com/labstack/echo/v4"

	_ "time_capsule_memories/docs" // Swagger docs

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	logging.Init(config.GetConfig().LogLevel)

	if err := database.Connect(); err != nil {
		logging.Fatal("database connection failed", "error", err)
	}
	defer database.Close()

	minio_client.MinioInit()

	jobs.StartScheduler()

	e := echo.New()
	e.Logger.SetLevel(0)

	e.Use(middleware.CORSConfig())

	routes.RegisterFileRoutes(e)
	routes.RegisterCapsuleRoutes(e)
	routes.RegisterFeedbackRoutes(e)
	routes.RegisterEmailRoutes(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	slog.Info("starting HTTP server", "addr", ":8000")
	e.Logger.Fatal(e.Start(":8000"))
}
