package middleware

import (
	"log"
	"net/http"
	"time_capsule_memories/internal/repository"

	"github.com/labstack/echo/v4"
)

// DBConnectionCheckMiddleware проверяет соединение с базой данных
func DBConnectionCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Проверяем соединение с базой данных
		if !repository.CheckDatabaseConnection() {
			log.Println("Ошибка подключения к базе данных")
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Внутренняя ошибка сервера"})
		}

		// Если соединение успешно, продолжаем
		return next(c)
	}
}
