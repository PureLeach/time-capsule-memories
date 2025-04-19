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
	"log"

	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/jobs"
	"time_capsule_memories/internal/middleware"
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/routes"

	"github.com/labstack/echo/v4"

	_ "time_capsule_memories/docs" // Swagger docs

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Initialize the database connection
	if err := database.Connect(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.Close()

	// Initialize the MinIO client
	minio_client.MinioInit()

	// Start background jobs (e.g., scheduled tasks)
	jobs.StartScheduler()

	// Create a new Echo instance
	e := echo.New()
	e.Logger.SetLevel(0)

	// Register global middleware
	e.Use(middleware.CORSConfig())

	// Register API routes
	routes.RegisterFileRoutes(e)
	routes.RegisterCapsuleRoutes(e)
	routes.RegisterFeedbackRoutes(e)
	routes.RegisterEmailRoutes(e)

	// Swagger documentation endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start the HTTP server
	log.Println("Starting server on port :8000")
	e.Logger.Fatal(e.Start(":8000"))
}
