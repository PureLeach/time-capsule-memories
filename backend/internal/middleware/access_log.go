package middleware

import (
	"time"
	"time_capsule_memories/internal/logging"

	"github.com/labstack/echo/v4"
)

// AccessLog emits a structured slog line per request once the handler returns.
// Paths listed in skipPaths are not logged — handy for liveness/readiness
// probes that otherwise dominate the log volume.
func AccessLog(skipPaths ...string) echo.MiddlewareFunc {
	skip := make(map[string]struct{}, len(skipPaths))
	for _, p := range skipPaths {
		skip[p] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if _, ok := skip[c.Request().URL.Path]; ok {
				return next(c)
			}

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
