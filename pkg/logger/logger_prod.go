//go:build prod

// Package logger
// Логирование в консоль не в проде
// Логирование в файл в проде
package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

func MustInit() *slog.Logger {
	err := os.MkdirAll("/var/log/app", 0o755)
	filePath := fmt.Sprintf("/var/log/app/%s.log", time.Now().Format("2006-01-02_15-04"))
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	log := slog.New(handler)

	return log
}
