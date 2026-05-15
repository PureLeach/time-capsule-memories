package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterHealthRoutes(e *echo.Echo) {
	e.GET("/healthz", handlers.Healthz)
	e.GET("/readyz", handlers.Readyz)
}
