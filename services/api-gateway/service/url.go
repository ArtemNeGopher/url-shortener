// Package service
// Сервисы для обработки данных
package service

import (
	"context"
	"time"

	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/models"
)

type URLCacheClient interface {
	GetURL(ctx context.Context, shortCode string) (*models.URL, error)
	SetURL(ctx context.Context, url *models.URL) error
	Delete(ctx context.Context, shortCode string) error
}

type URLServiceClient interface {
	CreateShortURL(ctx context.Context, urlFull string, expiresInDays *uint32) (*models.URL, error)
	GetURL(ctx context.Context, shortCode string) (*models.URL, error)
	Delete(ctx context.Context, shortCode string) error
}

type urlService struct {
	cache URLCacheClient
	url   URLServiceClient
}

func NewURLService(cacheClient URLCacheClient, urlClient URLServiceClient) *urlService {
	return &urlService{
		cache: cacheClient,
		url:   urlClient,
	}
}

func (s *urlService) GetURL(shortCode string) (*models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	u, err := s.cache.GetURL(ctx, shortCode)
	// Нашли в кэше, выходим
	if err == nil {
		return u, nil
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	u, err = s.url.GetURL(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэше ответ
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.cache.SetURL(ctx, u)

	return u, nil
}

func (s *urlService) CreateURL(url string, expiresInDays *uint32) (*models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := s.url.CreateShortURL(ctx, url, expiresInDays)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.cache.SetURL(ctx, ret)

	return ret, nil
}

func (s *urlService) Delete(shortCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = s.cache.Delete(ctx, shortCode)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.url.Delete(ctx, shortCode)

	return err
}
