// @title Time Capsule Memories API
// @version 1.0
// @description This is a sample server for Time Capsule Memories.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @license.name MIT
// @host localhost:8000
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

	_ "time_capsule_memories/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Инициализация базы данных
	if err := database.Connect(); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.Close()

	// Инициализация MinIO
	minio_client.MinioInit()

	// Запуск фоновых задач
	jobs.StartScheduler()

	// Создаем экземпляр Echo
	e := echo.New()
	e.Logger.SetLevel(0)

	// Применяем CORS middleware
	e.Use(middleware.CORSConfig())

	// Регистрируем обработчики
	routes.FileRegisterRoutes(e)
	routes.CapsuleRegisterRoutes(e)
	routes.FeedbackRegisterRoutes(e)
	routes.EmailRegisterRoutes(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Infof("Запуск сервера на порту :8000")
	e.Logger.Fatal(e.Start(":8000"))
}
