// Package grpc implements grpc interface
package grpc

import (
	"context"
	"fmt"
	"log/slog"
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
	cache    Cache
	cacheTTL time.Duration
	log      *slog.Logger
}

func NewServer(cache Cache, cacheTTL time.Duration, log *slog.Logger) *server {
	return &server{
		cache:    cache,
		cacheTTL: cacheTTL,
		log:      log.With(slog.String("context", "grpc")),
	}
}

func (s *server) Set(ctx context.Context, req *pb.CacheSetRequest) (*pb.CacheSetResponse, error) {
	err := s.cache.Set(ctx, req.Key, req.Value, s.cacheTTL)
	success := true
	if err != nil {
		success = false
	}

	s.log.Info(
		"Cahce set",
		slog.String("key", req.Key),
		slog.String("value", req.Value),
		slog.String("success", fmt.Sprintf("%v", success)),
	)

	return &pb.CacheSetResponse{Success: success}, nil
}

func (s *server) Get(ctx context.Context, req *pb.CacheGetRequest) (*pb.CacheGetResponse, error) {
	value, found, err := s.cache.Get(ctx, req.Key)

	if err != nil {
		s.log.Error(
			"Cache get error",
			slog.String("key", req.Key),
			slog.String("error", err.Error()),
		)
	} else {
		s.log.Info(
			"Cache get",
			slog.String("key", req.Key),
			slog.String("value", value),
		)
	}

	return &pb.CacheGetResponse{Value: value, Found: found}, err
}

func (s *server) Delete(ctx context.Context, req *pb.CacheDeleteRequest) (*pb.CacheDeleteResponse, error) {
	err := s.cache.Delete(ctx, req.Key)
	success := true
	if err != nil {
		success = false
	}

	if err != nil {
		s.log.Error(
			"Cache get error",
			slog.String("key", req.Key),
			slog.String("error", err.Error()),
		)
	} else {
		s.log.Info(
			"Cache get",
			slog.String("key", req.Key),
			slog.String("success", fmt.Sprintf("%v", success)),
		)
	}

	return &pb.CacheDeleteResponse{Success: success}, err
}
