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
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/routes"

	"github.com/labstack/echo/v4"

	_ "time_capsule_memories/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.Close()

	minio_client.MinioInit()

	// Создаем экземпляр Echo
	e := echo.New()
	e.Logger.SetLevel(0) // Установите уровень логирования на Debug

	// Передаем конфигурацию и minioClient в контекст Echo
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		c.Set("cfg", &cfg)                // Передаем указатель на cfg
	// 		c.Set("minioClient", minioClient) // Передаем minioClient
	// 		return next(c)
	// 	}
	// })

	// Регистрируем middleware
	// e.Use(middleware.DBConnectionCheckMiddleware)

	// Регистрируем обработчики
	routes.FileRegisterRoutes(e)
	routes.CapsuleRegisterRoutes(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Infof("Запуск сервера на порту :8000")
	e.Logger.Fatal(e.Start(":8000"))
}
