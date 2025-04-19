package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterCapsuleRoutes(e *echo.Echo) {
	e.POST("/capsules", handlers.CreateCapsule)
}
