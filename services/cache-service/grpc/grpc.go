// Package grpc implements grpc interface
package grpc

import (
	"context"
	"time"

	pb "github.com/ArtemNeGopher/url-shortener/pkg/genproto/cache"
)

type Cache interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, bool, error)
	Delete(ctx context.Context, key string) error
}

type server struct {
	pb.UnimplementedCacheServiceServer
	cache Cache
}

func NewServer(cache Cache) *server {
	return &server{
		cache: cache,
	}
}
