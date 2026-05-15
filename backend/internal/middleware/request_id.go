package middleware

import (
	"time_capsule_memories/internal/logging"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestID wraps echo's RequestID middleware so the generated id is also
// stashed on the request context. Downstream code reads it through
// logging.FromContext to attach request_id to every log line.
func RequestID() echo.MiddlewareFunc {
	base := middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		TargetHeader: echo.HeaderXRequestID,
	})
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		wrapped := base(func(c echo.Context) error {
			id := c.Response().Header().Get(echo.HeaderXRequestID)
			if id != "" {
				req := c.Request()
				c.SetRequest(req.WithContext(logging.WithRequestID(req.Context(), id)))
			}
			return next(c)
		})
		return wrapped
	}
}
