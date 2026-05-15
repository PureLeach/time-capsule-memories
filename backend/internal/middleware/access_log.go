package middleware

import (
	"time"
	"time_capsule_memories/internal/logging"

	"github.com/labstack/echo/v4"
)

// AccessLog emits a structured slog line per request once the handler returns.
// Sits after RequestID so request_id is available on the context.
func AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			req := c.Request()
			res := c.Response()

			logger := logging.FromContext(req.Context())
			logger.Info("http request",
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"duration_ms", time.Since(start).Milliseconds(),
				"remote_ip", c.RealIP(),
			)
			return err
		}
	}
}
