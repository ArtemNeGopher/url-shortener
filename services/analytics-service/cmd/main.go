package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	cfg "github.com/ArtemNeGopher/url-shortener/pkg/config"
	"github.com/ArtemNeGopher/url-shortener/pkg/logger"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/config"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/repository"

	pb "github.com/ArtemNeGopher/url-shortener/pkg/genproto/analytics"
	implGRPC "github.com/ArtemNeGopher/url-shortener/services/analytics-service/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Logger
	log := logger.MustInit()

	// Config
	config := &config.Config{}
	cfg.MustInit("config/config.yaml", config)
	log.Debug("config loaded")

	// Migrations
	migrations(config.Database.DatabaseURL, log)
	log.Debug("migrations completed")

	// Repository
	repo := repository.NewEventRepository(&config.Database, log)
	log.Debug("repository created")

	// GRPC
	implServer := implGRPC.NewService(repo, log)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logger.GetLoggingInterceptor(log)),
	)

	// implServer должен быть остановлен раньше
	// чем repo. Так как implServer
	// должен сохранить данные в репо
	defer func() {
		implServer.Stop()
		repo.Close()
	}()

	// Reflections
	reflection.Register(server)
	log.Debug("gRPC reflection registered")

	// Register grpc
	pb.RegisterAnalyticsServiceServer(server, implServer)
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

func migrations(databaseURL string, log *slog.Logger) {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	v, _, err := m.Version()
	log.Info(fmt.Sprintf("version %v", v))

	if sErr, dbErr := m.Close(); sErr != nil || dbErr != nil {
		log.Error(fmt.Sprintf("Error closing migrate: %v %v", sErr, dbErr))
	}
}
