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
	cache    Cache
	cacheTTL time.Duration
}

func NewServer(cache Cache, cacheTTL time.Duration) *server {
	return &server{
		cache:    cache,
		cacheTTL: cacheTTL,
	}
}

func (s *server) Set(ctx context.Context, req *pb.CacheSetRequest) (*pb.CacheSetResponse, error) {
	err := s.cache.Set(ctx, req.Key, req.Value, s.cacheTTL)
	if err != nil {
		return &pb.CacheSetResponse{Success: false}, err
	}
	return &pb.CacheSetResponse{Success: true}, nil
}

func (s *server) Get(ctx context.Context, req *pb.CacheGetRequest) (*pb.CacheGetResponse, error) {
	value, found, err := s.cache.Get(ctx, req.Key)
	return &pb.CacheGetResponse{Value: value, Found: found}, err
}

func (s *server) Delete(ctx context.Context, req *pb.CacheDeleteRequest) (*pb.CacheDeleteResponse, error) {
	err := s.cache.Delete(ctx, req.Key)
	if err != nil {
		return &pb.CacheDeleteResponse{Success: false}, err
	}
	return &pb.CacheDeleteResponse{Success: true}, err
}
