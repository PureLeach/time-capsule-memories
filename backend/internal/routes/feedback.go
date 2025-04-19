package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterFeedbackRoutes(e *echo.Echo) {
	e.POST("/feedback", handlers.CreateFeedback)
}
