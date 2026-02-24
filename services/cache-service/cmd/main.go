package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	configPkg "github.com/ArtemNeGopher/url-shortener/pkg/config"
	pb "github.com/ArtemNeGopher/url-shortener/pkg/genproto/cache"
	"github.com/ArtemNeGopher/url-shortener/pkg/logger"
	"github.com/ArtemNeGopher/url-shortener/services/cache-service/cache"
	cfg "github.com/ArtemNeGopher/url-shortener/services/cache-service/config"

	implGRPC "github.com/ArtemNeGopher/url-shortener/services/cache-service/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load config
	config := &cfg.Config{}
	configPkg.MustInit("config/config.yaml", config)

	// Init logger
	log := logger.MustInit()

	// Redis client
	rdb := InitRedis(config.RedisAddr, config.RedisPoolSize)
	log.Debug("Redis client created",
		slog.String("addr", config.RedisAddr))

	// Create Cache
	cache := cache.New(rdb, config.LocalTTL)
	defer cache.Close() // Закрываем кэш
	log.Debug("Cache created")

	// Create grpc server
	implServer := implGRPC.NewServer(cache, config.LocalTTL, log)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logger.GetLoggingInterceptor(log)),
	)
	log.Debug("gRPC server created")

	// Reflections
	reflection.Register(server)
	log.Debug("gRPC reflection registered")

	// Register grpc
	pb.RegisterCacheServiceServer(server, implServer)
	log.Debug("gRPC server registered")

	// Run grpc
	go func() {
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		lis, err := net.Listen("tcp", addr)
		log.Info(fmt.Sprintf("gRPC server listening on %s", addr))
		if err != nil {
			panic(err)
		}
		log.Debug("gRPC serve")
		server.Serve(lis)
		log.Info("gRPC stopped")
	}()

	// Create quit chan
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Debug("Quit chan created")

	// Wait signal
	<-quit
	log.Info("Shutting down gRPC server...")

	// GracefulStop
	server.GracefulStop()
}
