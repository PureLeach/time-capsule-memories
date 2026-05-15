// Package logging configures the application's structured logger.
package logging

import (
	"log/slog"
	"os"
	"strings"
)

// Init installs a JSON-formatted slog logger as the default. The level argument
// accepts "debug", "info", "warn" or "error" (case-insensitive); unrecognised
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
