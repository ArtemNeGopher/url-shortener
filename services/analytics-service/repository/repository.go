package repository

import (
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/models"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/worker"
)

type eventRepository struct{}

func NewEventRepository() *eventRepository {
	return &eventRepository{}
}

var _ worker.Repository = (*eventRepository)(nil)

func (eventrepository *eventRepository) BatchInsertClicks(events []models.ClickEvent) error {
	panic("not implemented") // TODO: Implement
}

func (eventrepository *eventRepository) UpdateStats(shortCode string) error {
	panic("not implemented") // TODO: Implement
}

func (eventrepository *eventRepository) GetStats(shortCode string) (*models.Stats, error) {
	panic("not implemented") // TODO: Implement
}
