package repository

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/config"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/models"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestDatabaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return "postgres://urlshortener:password@localhost:5432/urlshortener?sslmode=disable"
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", getTestDatabaseURL())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to database: %v", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		t.Skipf("Skipping test: cannot ping database: %v", err)
		return nil
	}

	db.Exec("TRUNCATE TABLE clicks, url_stats, url_day_stats RESTART IDENTITY CASCADE")

	return db
}

func cleanupTestDB(db *sql.DB) {
	if db == nil {
		return
	}
	db.Exec("TRUNCATE TABLE clicks, url_stats, url_day_stats RESTART IDENTITY CASCADE")
	db.Close()
}

func mockLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestBatchInsertClicks(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	events := []models.ClickEvent{
		{
			ShortCode: "test123",
			IPAddress: "192.168.1.1",
			Referer:   "https://example.com",
			UserAgent: "Mozilla/5.0",
			Timestamp: time.Now(),
		},
		{
			ShortCode: "test123",
			IPAddress: "192.168.1.2",
			Referer:   "https://google.com",
			UserAgent: "Chrome",
			Timestamp: time.Now(),
		},
	}

	err := repo.BatchInsertClicks(events)
	require.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM clicks WHERE short_code = 'test123'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestBatchInsertClicksEmpty(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	err := repo.BatchInsertClicks([]models.ClickEvent{})
	assert.NoError(t, err)
}

func TestUpdateStats(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	events := []models.ClickEvent{
		{
			ShortCode: "statst",
			IPAddress: "10.0.0.1",
			UserAgent: "Agent1",
			Timestamp: time.Now(),
		},
		{
			ShortCode: "statst",
			IPAddress: "10.0.0.1",
			UserAgent: "Agent1",
			Timestamp: time.Now(),
		},
		{
			ShortCode: "statst",
			IPAddress: "10.0.0.2",
			UserAgent: "Agent2",
			Timestamp: time.Now(),
		},
	}
	repo.BatchInsertClicks(events)

	err := repo.UpdateStats("statst")
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	stats, err := repo.GetStats("statst")
	require.NoError(t, err)
	assert.Equal(t, int64(3), stats.TotalClicks)
	assert.Equal(t, int64(2), stats.UniqueVisitors)
}

func TestUpdateStatsNotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	err := repo.UpdateStats("nonexist")
	assert.NoError(t, err)

	stats, err := repo.GetStats("nonexist")
	require.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestGetStats(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	events := []models.ClickEvent{
		{
			ShortCode: "getstat",
			IPAddress: "10.0.0.1",
			UserAgent: "Agent1",
			Timestamp: time.Now(),
		},
	}
	repo.BatchInsertClicks(events)
	err := repo.UpdateStats("getstat")
	require.NoError(t, err)

	stats, err := repo.GetStats("getstat")
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(1), stats.TotalClicks)
}

func TestGetStatsNotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	stats, err := repo.GetStats("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, stats)
}

func TestUpdateDayStats(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	date := time.Now().Format("2006-01-02")

	events := []models.ClickEvent{
		{
			ShortCode: "daystat",
			IPAddress: "10.0.0.1",
			Referer:   "https://example.com",
			UserAgent: "Agent1",
			Timestamp: time.Now(),
		},
		{
			ShortCode: "daystat",
			IPAddress: "10.0.0.1",
			Referer:   "https://example.com",
			UserAgent: "Agent1",
			Timestamp: time.Now(),
		},
		{
			ShortCode: "daystat",
			IPAddress: "10.0.0.2",
			Referer:   "https://google.com",
			UserAgent: "Agent2",
			Timestamp: time.Now(),
		},
	}
	repo.BatchInsertClicks(events)

	err := repo.UpdateDayStats("daystat", date)
	require.NoError(t, err)

	stats, err := repo.GetDayStats("daystat", date)
	require.NoError(t, err)
	assert.Equal(t, int64(3), stats.TotalClicks)
	assert.Equal(t, int64(2), stats.UniqueVisitors)
	assert.Equal(t, 2, len(stats.Referers))
}

func TestGetDayStats(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	date := time.Now().Format("2006-01-02")

	events := []models.ClickEvent{
		{
			ShortCode: "dayget",
			IPAddress: "10.0.0.1",
			UserAgent: "Agent1",
			Timestamp: time.Now(),
		},
	}
	repo.BatchInsertClicks(events)
	repo.UpdateDayStats("dayget", date)

	stats, err := repo.GetDayStats("dayget", date)
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, date, stats.Date)
}

func TestGetDayStatsNotFound(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	date := time.Now().Format("2006-01-02")

	stats, err := repo.GetDayStats("nonexist", date)
	require.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestNewEventRepository(t *testing.T) {
	tests := []struct {
		name    string
		dbURL   string
		wantErr bool
	}{
		{
			name:    "valid connection",
			dbURL:   getTestDatabaseURL(),
			wantErr: false,
		},
		{
			name:    "invalid connection",
			dbURL:   "postgres://invalid:invalid@localhost:9999/nonexistent?sslmode=disable",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic for invalid connection")
					}
				}()
			}

			cfg := &config.DatabaseConfig{
				DatabaseURL:        tt.dbURL,
				MaxOpenConnections: 5,
				MaxIdleConnections: 5,
			}

			repo := NewEventRepository(cfg, mockLogger())
			if repo != nil {
				repo.Close()
			}
		})
	}
}

func TestGetReferers(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer cleanupTestDB(db)

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	events := []models.ClickEvent{
		{
			ShortCode: "referer",
			IPAddress: "10.0.0.1",
			Referer:   "https://example.com",
			UserAgent: "Agent1",
			Timestamp: time.Now(),
		},
		{
			ShortCode: "referer",
			IPAddress: "10.0.0.2",
			Referer:   "https://google.com",
			UserAgent: "Agent2",
			Timestamp: time.Now(),
		},
	}
	repo.BatchInsertClicks(events)
	repo.UpdateStats("referer")

	referers, err := repo.GetStats("referer")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(referers.Referers), 1)
}

func BenchmarkBatchInsertClicks(b *testing.B) {
	db, err := sql.Open("postgres", getTestDatabaseURL())
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to database: %v", err)
		return
	}
	defer db.Close()

	log := mockLogger()
	repo := &eventRepository{db: db, log: log}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		events := []models.ClickEvent{
			{
				ShortCode: fmt.Sprintf("bench%d", i),
				IPAddress: "192.168.1.1",
				Referer:   "https://example.com",
				UserAgent: "Mozilla/5.0",
				Timestamp: time.Now(),
			},
		}
		repo.BatchInsertClicks(events)
	}
}
