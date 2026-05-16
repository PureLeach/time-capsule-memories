package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterCapsuleRoutes(e *echo.Echo, h *handlers.Handler) {
	e.POST("/capsules", h.CreateCapsule)
}
