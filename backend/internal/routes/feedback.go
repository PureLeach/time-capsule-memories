package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func FeedbackRegisterRoutes(e *echo.Echo) {
	e.POST("/feedback", handlers.CreateFeedback)
}
