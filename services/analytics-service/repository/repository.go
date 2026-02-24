package repository

import (
	"database/sql"
	"log/slog"

	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/config"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/models"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/worker"

	_ "github.com/lib/pq"
)

type eventRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewEventRepository(config *config.DatabaseConfig, log *slog.Logger) *eventRepository {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)

	return &eventRepository{
		db:  db,
		log: log.With(slog.String("context", "repository")),
	}
}

var _ worker.Repository = (*eventRepository)(nil)

func (repo *eventRepository) Close() {
	repo.db.Close()
}

func (repo *eventRepository) BatchInsertClicks(events []models.ClickEvent) error {
	panic("not implemented") // TODO: Implement
}

func (repo *eventRepository) UpdateStats(shortCode string) error {
	panic("not implemented") // TODO: Implement
}

func (repo *eventRepository) GetStats(shortCode string) (*models.Stats, error) {
	panic("not implemented") // TODO: Implement
}

func (repo *eventRepository) GetDayStats(shortCode string, date string) (*models.DayStats, error) {
	panic("not implemented") // TODO: Implement
}
