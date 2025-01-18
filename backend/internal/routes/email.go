package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func EmailRegisterRoutes(e *echo.Echo) {
	e.POST("/send-test-email", handlers.SendTestEmail)
}
