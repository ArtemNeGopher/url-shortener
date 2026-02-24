// Package grpc implements grpc interface
package grpc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	pb "github.com/ArtemNeGopher/url-shortener/pkg/genproto/cache"
)

var ErrCanceled = errors.New("canceled")

type Cache interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, bool, error)
	Delete(ctx context.Context, key string) error
}

type server struct {
	cache    Cache
	cacheTTL time.Duration
	log      *slog.Logger
	pb.UnimplementedCacheServiceServer
}

func NewServer(cache Cache, cacheTTL time.Duration, log *slog.Logger) *server {
	return &server{
		cache:    cache,
		cacheTTL: cacheTTL,
		log:      log.With(slog.String("context", "grpc-server")),
	}
}

var _ pb.CacheServiceServer = (*server)(nil)

func (s *server) Set(ctx context.Context, req *pb.CacheSetRequest) (*pb.CacheSetResponse, error) {
	done := make(chan struct{})
	var err error

	go func() {
		err = s.cache.Set(ctx, req.Key, req.Value, s.cacheTTL)
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil, ErrCanceled
	}

	success := err != nil

	s.log.Info(
		"Cahce set",
		slog.String("key", req.Key),
		slog.String("value", req.Value),
		slog.Bool("success", success),
	)

	return &pb.CacheSetResponse{Success: success}, nil
}

func (s *server) Get(ctx context.Context, req *pb.CacheGetRequest) (*pb.CacheGetResponse, error) {
	done := make(chan struct{})
	var value string
	var found bool
	var err error

	go func() {
		value, found, err = s.cache.Get(ctx, req.Key)
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil, ErrCanceled
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
			slog.String("value", value),
			slog.Bool("found", found),
		)
	}

	return &pb.CacheGetResponse{Value: value, Found: found}, err
}

func (s *server) Delete(ctx context.Context, req *pb.CacheDeleteRequest) (*pb.CacheDeleteResponse, error) {
	done := make(chan struct{})
	var err error

	go func() {
		err = s.cache.Delete(ctx, req.Key)
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil, ErrCanceled
	}

	success := err != nil

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
			slog.Bool("success", success),
		)
	}

	return &pb.CacheDeleteResponse{Success: success}, err
}
