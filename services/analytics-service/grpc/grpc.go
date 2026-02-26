package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/ArtemNeGopher/url-shortener/pkg/genproto/analytics"
	pb "github.com/ArtemNeGopher/url-shortener/pkg/genproto/analytics"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/models"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/worker"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StatsRepository interface {
	BatchInsertClicks(events []models.ClickEvent) error
	UpdateStats(shortCode string) error
	GetStats(shortCode string) (*models.Stats, error)
	GetDayStats(shortCode string, date string) (*models.DayStats, error)
}

type server struct {
	repo         StatsRepository
	log          *slog.Logger
	recordWorker *worker.WorkerPool
	pb.UnimplementedAnalyticsServiceServer
}

func NewServer(repo StatsRepository, log *slog.Logger) *server {
	recordWorker := worker.New(10, 100, repo, log)
	recordWorker.Start()

	return &server{
		repo:         repo,
		log:          log,
		recordWorker: recordWorker,
	}
}

var _ pb.AnalyticsServiceServer = (*server)(nil)

func (s *server) Stop() {
	s.recordWorker.Stop()
}

func (s *server) RecordClick(ctx context.Context, req *analytics.ClickEvent) (*analytics.ClickResponse, error) {
	s.log.Info(
		"RecordClick",
		slog.String("short_code", req.ShortCode),
		slog.String("ip", req.IpAddress),
		slog.String("referer", req.Referer),
		slog.Time("stamp", req.ClickedAt.AsTime()),
	)

	// Клик произошёл в будущем, отклоняем
	if req.ClickedAt.AsTime().After(time.Now()) {
		return &analytics.ClickResponse{Success: false}, nil
	}

	event := models.ClickEvent{
		ShortCode: req.ShortCode,
		IPAddress: req.IpAddress,
		UserAgent: req.UserAgent,
		Referer:   req.Referer,
		Timestamp: req.ClickedAt.AsTime(),
	}
	s.recordWorker.Submit(event)

	return &analytics.ClickResponse{Success: true}, nil
}

func (s *server) GetStatistics(ctx context.Context, req *analytics.StatsRequest) (*analytics.StatsResponse, error) {
	s.log.Info("GetStatistics", slog.String("short_code", req.ShortCode))
	done := make(chan struct{})
	var stats *models.Stats
	var err error

	go func() {
		stats, err = s.repo.GetStats(req.ShortCode)
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		s.log.Debug("context canceled")
		return nil, context.Canceled
	}

	if err != nil {
		return nil, err
	}

	resp := &analytics.StatsResponse{
		TotalClicks:    stats.TotalClicks,
		UniqueVisitors: stats.UniqueVisitors,
		Referers:       stats.Referers,
		LastClickedAt:  timestamppb.New(*stats.LastClickedAt),
	}

	return resp, nil
}

func (s *server) GetDayStatistics(ctx context.Context, req *analytics.DayStatsRequest) (*analytics.DayStatsResponse, error) {
	date := req.Date.AsTime().Format("2006-01-02")
	s.log.Info("GetDayStatistics", slog.String("short_code", req.ShortCode), slog.String("date", date))

	// Так как в будущем, возвращаем пустую структуру
	if req.Date.AsTime().After(time.Now()) {
		resp := &analytics.DayStatsResponse{
			TotalClicks:    0,
			UniqueVisitors: 0,
			Referers:       []string{},
		}

		return resp, nil
	}

	done := make(chan struct{})
	var stats *models.DayStats
	var err error

	go func() {
		stats, err = s.repo.GetDayStats(req.ShortCode, date)
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		s.log.Debug("context canceled")
		return nil, context.Canceled
	}

	if err != nil {
		return nil, err
	}

	resp := &analytics.DayStatsResponse{
		TotalClicks:    stats.TotalClicks,
		UniqueVisitors: stats.UniqueVisitors,
		Referers:       stats.Referers,
	}

	return resp, nil
}
