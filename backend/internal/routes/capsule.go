package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	// e.GET("/capsules", handlers.GetCapsules)
	e.POST("/capsules", handlers.CreateCapsule)
}
