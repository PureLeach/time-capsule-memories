package middleware

import (
	"time_capsule_memories/internal/logging"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Recover wraps echo's recovery middleware and routes panic logs through slog
// so they pick up the request_id attached by RequestID.
func Recover() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         4 << 10,
		DisableStackAll:   false,
		DisablePrintStack: true,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logging.FromContext(c.Request().Context()).Error("panic recovered",
				"error", err,
				"stack", string(stack),
			)
			return err
		},
	})
}
