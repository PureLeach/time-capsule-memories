// Package routes maps URL paths onto handler methods.
package routes

import (
	"time_capsule_memories/internal/handlers"

	"github.com/labstack/echo/v4"
)

// Register wires the always-on routes.
func Register(e *echo.Echo, h *handlers.Handler) {
	e.GET("/healthz", h.Healthz)
	e.GET("/readyz", h.Readyz)

	e.GET("/generate-presigned-url", h.GeneratePresignedURL)
	e.POST("/capsules", h.CreateCapsule)
	e.POST("/feedback", h.CreateFeedback)
}

// RegisterTestEmail exposes an unauthenticated SMTP smoke test. Gate it on
// ENABLE_TEST_EMAIL_ENDPOINT.
func RegisterTestEmail(e *echo.Echo, h *handlers.Handler) {
	e.POST("/send-test-email", h.SendTestEmail)
}
