//go:build !prod

// Package logger
package logger

import (
	"log/slog"
	"os"
)

func MustInit() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	log := slog.New(handler)

	return log
}
