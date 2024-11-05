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

	"github.com/labstack/echo/v4"

	"time_capsule_memories/internal/config" // Импортируйте пакет config
	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/routes"

	_ "time_capsule_memories/docs"

	echoSwagger "github.com/swaggo/echo-swagger"

	"time_capsule_memories/internal/middleware" // Импортируйте пакет middleware
	// Импортируйте пакет middleware
	// "github.com/labstack/echo/v4"
	// "log"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig() // Используйте новый пакет
	if err != nil {
		log.Fatalf("Не удалось загрузить значения из переменных окружения: %v", err)
	}

	// Подключаемся к базе данных
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.Close()

	e := echo.New()
	e.Logger.SetLevel(0) // Установите уровень логирования на Debug

	// Регистрируем middleware
	e.Use(middleware.DBConnectionCheckMiddleware)

	routes.RegisterRoutes(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Infof("Запуск сервера на порту :8000")
	e.Logger.Fatal(e.Start(":8000"))

}
