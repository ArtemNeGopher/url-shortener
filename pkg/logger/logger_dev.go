//go:build !prod

// Package logger
// Логирование в консоль не в проде
// Логирование в файл в проде
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
