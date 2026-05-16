package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterEmailRoutes(e *echo.Echo, h *handlers.Handler) {
	e.POST("/send-test-email", h.SendTestEmail)
}
