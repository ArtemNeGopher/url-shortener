package grpc

import (
	"context"
	"time"

	"github.com/ArtemNeGopher/url-shortener/pkg/genproto/analytics"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/config"
	"github.com/ArtemNeGopher/url-shortener/services/api-gateway/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalyticsClient struct {
	client analytics.AnalyticsServiceClient
	conn   *grpc.ClientConn
}

func NewAnalyticsClient(cfg *config.AnalyticsClientConfig) *AnalyticsClient {
	conn, err := grpc.NewClient(
		cfg.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err.Error())
	}

	client := analytics.NewAnalyticsServiceClient(conn)
	return &AnalyticsClient{
		client: client,
		conn:   conn,
	}
}

func (c *AnalyticsClient) Close() {
	c.conn.Close()
}

func (c *AnalyticsClient) RegisterClick(ctx context.Context, click *models.Click) {
	in := &analytics.ClickEvent{
		ShortCode: click.ShortCode,
		ClickedAt: timestamppb.New(click.ClickedAt),
		IpAddress: click.IPAdress,
		UserAgent: click.UserAgent,
		Referer:   click.Referer,
	}

	_, _ = c.client.RecordClick(ctx, in)
}

func (c *AnalyticsClient) GetStats(ctx context.Context, shortCode string) (*models.Stats, error) {
	in := &analytics.StatsRequest{
		ShortCode: shortCode,
	}

	resp, err := c.client.GetStatistics(ctx, in)
	if err != nil {
		return nil, err
	}

	return &models.Stats{
		ShortCode:      shortCode,
		TotalClicks:    resp.TotalClicks,
		UniqueVisitors: resp.UniqueVisitors,
		LastClickedAt:  resp.LastClickedAt.AsTime(),
		Referers:       resp.Referers,
	}, nil
}

func (c *AnalyticsClient) GetDayStats(ctx context.Context, shortCode string, date string) (*models.DayStats, error) {
	in := &analytics.DayStatsRequest{
		ShortCode: shortCode,
		Date:      timestamppb.New(time.Now()),
	}

	resp, err := c.client.GetDayStatistics(ctx, in)
	if err != nil {
		return nil, err
	}

	return &models.DayStats{
		ShortCode:      shortCode,
		Date:           date,
		TotalClicks:    resp.TotalClicks,
		UniqueVisitors: resp.UniqueVisitors,
		Referers:       resp.Referers,
	}, nil
}
