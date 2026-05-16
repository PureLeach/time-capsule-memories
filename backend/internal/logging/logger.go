// Package logging configures the application's structured logger.
package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// ctxKey is unexported so callers cannot accidentally collide with our context values.
type ctxKey int

const (
	requestIDKey ctxKey = iota
)

// Init installs a JSON-formatted slog logger as the default. The level argument
// accepts "debug", "info", "warn" or "error" (case-insensitive); unrecognized
// values fall back to info.
func Init(level string) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     parseLevel(level),
		AddSource: true,
	})
	slog.SetDefault(slog.New(handler))
}

// Fatal logs the message at error level and terminates the process with status 1.
// Intended for unrecoverable startup failures.
func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

// WithRequestID returns a child context carrying the given request id.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestIDFromContext returns the request id stored on ctx, or "" if absent.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// FromContext returns a logger pre-populated with attributes pulled off ctx
// (currently just request_id). Falls back to the default logger when ctx
// carries nothing useful.
func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}
	if id := RequestIDFromContext(ctx); id != "" {
		return slog.Default().With("request_id", id)
	}
	return slog.Default()
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
