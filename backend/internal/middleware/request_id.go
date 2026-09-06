package middleware

import (
	"time_capsule_memories/internal/logging"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestID also stashes echo's generated id on the request context, where
// logging.FromContext picks it up.
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
