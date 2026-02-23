//go:build prod

package logger

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestMustInitProdCreatesLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	err = os.MkdirAll("logs", 0755)
	if err != nil {
		t.Fatalf("failed to create logs directory: %v", err)
	}

	log := MustInit()
	log.Info("test message")

	time.Sleep(100 * time.Millisecond)

	files, err := os.ReadDir("logs")
	if err != nil {
		t.Fatalf("failed to read logs directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("expected log file to be created")
	}

	found := false
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".log") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected a .log file to exist in logs directory")
	}
}

func TestMustInitProdJSONHandler(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	err = os.MkdirAll("logs", 0755)
	if err != nil {
		t.Fatalf("failed to create logs directory: %v", err)
	}

	log := MustInit()

	handler := log.Handler()
	if handler == nil {
		t.Error("expected handler to not be nil")
	}

	log.Info("json test")
}
