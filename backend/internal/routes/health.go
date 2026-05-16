package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterHealthRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET("/healthz", h.Healthz)
	e.GET("/readyz", h.Readyz)
}
