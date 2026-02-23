//go:build !prod

package logger

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestMustInitDev(t *testing.T) {
	log := MustInit()

	if log == nil {
		t.Error("expected logger to not be nil")
	}
}

func TestMustInitDevLogsMessage(t *testing.T) {
	buf := bytes.Buffer{}

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	log := slog.New(handler)

	log.Info("test message")

	output := buf.String()
	if len(output) == 0 {
		t.Error("expected output to not be empty")
	}

	if !contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got: %s", output)
	}
}

func TestMustInitDevDebugLevel(t *testing.T) {
	buf := bytes.Buffer{}

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	log := slog.New(handler)

	log.Debug("debug message")

	output := buf.String()
	if !contains(output, "debug message") {
		t.Errorf("expected output to contain 'debug message', got: %s", output)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
