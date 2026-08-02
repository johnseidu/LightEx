package logger

import (
	"log/slog"
	"os"
)

// New creates the application's default logger.
//
// Development:
//   - Human-readable text output
//
// Production (future):
//   - JSON output
func New() *slog.Logger {

	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	return slog.New(handler)
}
