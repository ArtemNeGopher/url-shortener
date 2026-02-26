package grpc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	u "net/url"

	"github.com/ArtemNeGopher/url-shortener/pkg/genproto/url"
	pb "github.com/ArtemNeGopher/url-shortener/pkg/genproto/url"
	"github.com/ArtemNeGopher/url-shortener/pkg/shortcode"
	"github.com/ArtemNeGopher/url-shortener/services/url-service/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type URLRepository interface {
	TryRegisterURL(url *models.URL) error
	GetOriginalURL(shortCode string) (*models.URL, error)
	DeleteURL(shortCode string) error
}

type server struct {
	repo             URLRepository
	urlCreateRetries uint
	pb.UnimplementedURLServiceServer
}

func NewServer(repo URLRepository, urlCreateRetries uint, log *slog.Logger) *server {
	return &server{
		repo:             repo,
		urlCreateRetries: urlCreateRetries,
	}
}

func isURL(str string) bool {
	u, err := u.ParseRequestURI(str)
	// Для абсолютных URL: u.Scheme и u.Host должны быть не пустыми
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (s *server) CreateShortURL(ctx context.Context, req *url.CreateURLRequest) (*url.CreateURLResponse, error) {
	var expiresAt *time.Time = nil
	if req.ExpiresInDays != nil && *req.ExpiresInDays != 0 {
		expiresAt = &time.Time{}
		*expiresAt = time.Now().Add(time.Hour * 24 * time.Duration(*req.ExpiresInDays))
	}
	if !isURL(req.Url) {
		return nil, errors.New("invalid url")
	}

	done := make(chan struct{})

	crUrl := &models.URL{
		URL:       req.Url,
		ExpiresAt: expiresAt,
	}

	var err error
	go func() {
		for range s.urlCreateRetries {
			crUrl.ShortCode, err = shortcode.Generate()
			if err != nil {
				continue
			}
			err = s.repo.TryRegisterURL(crUrl)
		}

		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil, context.Canceled
	}

	if err != nil {
		return nil, err
	}

	resp := &url.CreateURLResponse{
		ShortCode: crUrl.ShortCode,
	}
	if crUrl.ExpiresAt != nil {
		resp.ExpiresAt = timestamppb.New(*crUrl.ExpiresAt)
	}

	return resp, nil
}

func (s *server) GetOriginalURL(ctx context.Context, req *url.GetURLRequest) (*url.GetURLResponse, error) {
	done := make(chan struct{})

	var origURL *models.URL
	var err error
	go func() {
		origURL, err = s.repo.GetOriginalURL(req.ShortCode)
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil, context.Canceled
	}

	if err != nil {
		return nil, err
	}

	resp := &url.GetURLResponse{
		Url:      origURL.URL,
		IsActive: origURL.IsActive,
	}

	if origURL.ExpiresAt != nil {
		resp.ExpiresAt = timestamppb.New(*origURL.ExpiresAt)
	}

	return resp, nil
}

func (s *server) DeleteURL(ctx context.Context, req *url.DeleteURLRequest) (*url.DeleteURLResponse, error) {
	done := make(chan struct{})

	var err error
	go func() {
		err = s.repo.DeleteURL(req.ShortCode)
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil, context.Canceled
	}

	success := err != nil

	resp := &url.DeleteURLResponse{
		Success: success,
	}

	return resp, nil
}
