// Package grpc
// Содержит клиенты для обращения к другим сервисам.
package grpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ArtemNeGopher/url-shortener/pkg/genproto/cache"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/config"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/models"
	"google.golang.org/grpc"
)

type CacheClient struct {
	client cache.CacheServiceClient
	conn   *grpc.ClientConn
}

func NewCacheClient(cfg *config.CacheClientConfig) *CacheClient {
	conn, err := grpc.NewClient(
		cfg.Addr,
	)
	if err != nil {
		panic(err.Error())
	}

	client := cache.NewCacheServiceClient(conn)
	return &CacheClient{
		client: client,
		conn:   conn,
	}
}

func (c *CacheClient) Close() {
	c.conn.Close()
}

func (c *CacheClient) SetURL(ctx context.Context, url *models.URL) error {
	json, err := json.Marshal(url)
	if err != nil {
		return err
	}

	in := &cache.CacheSetRequest{
		Key:   url.ShortCode,
		Value: string(json),
	}
	_, err = c.client.Set(ctx, in)
	return err
}

func (c *CacheClient) GetURL(ctx context.Context, shortCode string) (*models.URL, error) {
	// Валидация shortCode
	if len(shortCode) != 7 {
		return nil, errors.New("shortCode length must be 7")
	}

	// Идём в cache-service
	in := &cache.CacheGetRequest{
		Key: shortCode,
	}
	resp, err := c.client.Get(ctx, in)
	if err != nil {
		return nil, err
	}

	// Записи нет, кэш мис
	if !resp.Found {
		return nil, errors.New("cache miss")
	}

	// Парсим выход
	u := &models.URL{}
	err = json.Unmarshal([]byte(resp.Value), u)
	if err != nil {
		// удаляем из кэша плохую запись
		go c.Delete(ctx, shortCode)

		return nil, errors.New("invalid URL in cache")
	}

	return u, nil
}

func (c *CacheClient) Delete(ctx context.Context, shortCode string) error {
	in := &cache.CacheDeleteRequest{
		Key: shortCode,
	}
	_, err := c.client.Delete(ctx, in)
	return err
}
