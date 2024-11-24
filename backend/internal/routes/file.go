package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func FileRegisterRoutes(e *echo.Echo) {
	e.GET("/generate-presigned-url", handlers.GeneratePresignedURLHandler)
}
