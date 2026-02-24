package main

import (
	"fmt"
	"os"

	cfg "github.com/ArtemNeGopher/url-shortener/pkg/config"
	"github.com/ArtemNeGopher/url-shortener/pkg/logger"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log := logger.MustInit()

	config := &config.Config{}
	cfg.MustInit("config/config.yaml", config)
	log.Debug("config loaded")

	files, _ := os.ReadDir("./migrations")
	fmt.Printf("Found migration files: %v\n", len(files))
	for _, f := range files {
		fmt.Printf("  - %s\n", f.Name())
	}

	m, err := migrate.New(
		"file://migrations",
		config.DatabaseURL,
	)
	if err != nil {
		panic(err)
	}

	m.GracefulStop = make(chan bool)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	v, _, err := m.Version()
	log.Info(fmt.Sprintf("version %v", v))

	close(m.GracefulStop)

	if sErr, dbErr := m.Close(); sErr != nil || dbErr != nil {
		log.Error(fmt.Sprintf("Error closing migrate: %v %v", sErr, dbErr))
	}
	log.Debug("migrations completed")
}
