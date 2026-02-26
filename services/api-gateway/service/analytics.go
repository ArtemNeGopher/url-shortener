package service

import (
	"context"
	"time"

	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/models"
)

type AnalyticsClient interface {
	RegisterClick(ctx context.Context, click *models.Click)
	GetStats(ctx context.Context, shortCode string) (*models.Stats, error)
	GetDayStats(ctx context.Context, shortCode string, date string) (*models.DayStats, error)
}

type analyticsService struct {
	client AnalyticsClient
}

func NewAnalyticsService(client AnalyticsClient) *analyticsService {
	return &analyticsService{
		client: client,
	}
}

func (s *analyticsService) RegisterClick(click *models.Click) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.client.RegisterClick(ctx, click)
}

func (s *analyticsService) GetStats(shortCode string) (*models.Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.GetStats(ctx, shortCode)
}

func (s *analyticsService) GetDayStats(shortCode string, date string) (*models.DayStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.GetDayStats(ctx, shortCode, date)
}
