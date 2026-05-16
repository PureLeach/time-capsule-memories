package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterFileRoutes(e *echo.Echo, h *handlers.Handler) {
	e.GET("/generate-presigned-url", h.GeneratePresignedURL)
}
