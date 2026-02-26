package config

import (
	"os"
	"testing"
)

type TestConfig struct {
	Host string `yaml:"host" env:"HOST"`
	Port int    `yaml:"port" env:"PORT"`
}

func TestInit(t *testing.T) {
	content := `host: "localhost"
port: 8080
`
	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	var cfg TestConfig
	err = Init(tmpFile.Name(), &cfg)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected host 'localhost', got '%s'", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
}

func TestInitFileNotFound(t *testing.T) {
	var cfg TestConfig
	err := Init("nonexistent.yaml", &cfg)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestMustInit(t *testing.T) {
	content := `host: "127.0.0.1"
port: 3000
`
	tmpFile, err := os.CreateTemp("", "test_must_init_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	var cfg TestConfig
	MustInit(tmpFile.Name(), &cfg)

	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected host '127.0.0.1', got '%s'", cfg.Host)
	}
	if cfg.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Port)
	}
}
