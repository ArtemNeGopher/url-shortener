package grpc

import (
	"context"
	"time"

	"github.com/ArtemNeGopher/url-shortener/pkg/genproto/url"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/config"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type URLClient struct {
	client url.URLServiceClient
	conn   *grpc.ClientConn
}

func NewURLClient(cfg *config.URLClientConfig) *URLClient {
	conn, err := grpc.NewClient(
		cfg.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err.Error())
	}

	client := url.NewURLServiceClient(conn)
	return &URLClient{
		client: client,
		conn:   conn,
	}
}

func (c *URLClient) Close() {
	c.conn.Close()
}

func (c *URLClient) CreateShortURL(ctx context.Context, urlFull string, expiresInDays *uint32) (*models.URL, error) {
	in := &url.CreateURLRequest{
		Url:           urlFull,
		ExpiresInDays: expiresInDays,
	}

	resp, err := c.client.CreateShortURL(ctx, in)
	if err != nil {
		return nil, err
	}

	expiresAt := new(time.Time)
	if resp.ExpiresAt != nil {
		*expiresAt = resp.ExpiresAt.AsTime()
	}
	return &models.URL{
		ShortCode: resp.ShortCode,
		URL:       urlFull,
		ExpiresAt: expiresAt,
		IsActive:  true,
	}, nil
}

func (c *URLClient) GetURL(ctx context.Context, shortCode string) (*models.URL, error) {
	in := &url.GetURLRequest{
		ShortCode: shortCode,
	}

	resp, err := c.client.GetOriginalURL(ctx, in)
	if err != nil {
		return nil, err
	}

	expiresAt := new(time.Time)
	if resp.ExpiresAt != nil {
		*expiresAt = resp.ExpiresAt.AsTime()
	}
	return &models.URL{
		ShortCode: shortCode,
		URL:       resp.Url,
		ExpiresAt: expiresAt,
		IsActive:  resp.IsActive,
	}, nil
}

func (c *URLClient) Delete(ctx context.Context, shortCode string) error {
	in := &url.DeleteURLRequest{
		ShortCode: shortCode,
	}

	_, err := c.client.DeleteURL(ctx, in)
	return err
}
