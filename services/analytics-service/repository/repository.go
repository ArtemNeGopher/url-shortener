package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/config"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/models"
	"github.com/ArtemNeGopher/url-shortener/services/analytics-service/worker"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type eventRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewEventRepository(config *config.DatabaseConfig, log *slog.Logger) *eventRepository {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)

	return &eventRepository{
		db:  db,
		log: log.With(slog.String("context", "repository")),
	}
}

var _ worker.Repository = (*eventRepository)(nil)

func (repo *eventRepository) Close() {
	repo.db.Close()
}

func (repo *eventRepository) BatchInsertClicks(events []models.ClickEvent) error {
	if len(events) == 0 {
		return nil
	}

	repo.log.Debug("batch inserting clicks", slog.Int("count", len(events)))

	sql := sq.Insert("clicks").
		Columns("short_code", "ip_address", "referer", "user_agent", "clicked_at").
		PlaceholderFormat(sq.Dollar)

	for _, event := range events {
		sql = sql.Values(event.ShortCode, event.IPAddress, event.Referer, event.UserAgent, event.Timestamp)
	}

	_, err := sql.RunWith(repo.db).Exec()
	if err != nil {
		repo.log.Error("failed to batch insert clicks", slog.String("error", err.Error()))
		return err
	}

	repo.log.Debug("successfully inserted clicks", slog.Int("count", len(events)))
	return nil
}

func (repo *eventRepository) UpdateStats(shortCode string) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	lastClickQuery := sq.Select("MAX(clicked_at)").
		From("clicks").
		Where(sq.Eq{"short_code": shortCode}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx)

	lastClick := sql.NullTime{}
	err = lastClickQuery.QueryRow().Scan(&lastClick)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var uniqueVisitors sql.NullInt64
	uniqueQuery := sq.Select("COALESCE(SUM(user_count), 0) AS unique_visitors").
		FromSelect(
			sq.Select("COUNT(DISTINCT user_agent) AS user_count").
				From("clicks").
				Where(sq.Eq{"short_code": shortCode}).
				GroupBy("ip_address"),
			"unique_visitors_ips",
		).
		PlaceholderFormat(sq.Dollar)

	err = uniqueQuery.RunWith(tx).QueryRow().Scan(&uniqueVisitors)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var totalClicks int64
	totalQuery := sq.Select("COUNT(id) AS total_clicks").
		From("clicks").
		Where(sq.Eq{"short_code": shortCode}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx)

	err = totalQuery.QueryRow().Scan(&totalClicks)
	if err != nil {
		fmt.Println(err)
		return err
	}

	referers, err := getReferers(tx, shortCode, "")
	if err != nil {
		return err
	}

	referersJSON, err := json.Marshal(referers)
	if err != nil {
		return err
	}

	uQuery, err := sq.Update("url_stats").
		Set("total_clicks", totalClicks).
		Set("unique_visitors", uniqueVisitors).
		Set("last_clicked_at", lastClick).
		Set("referers", sq.Expr("?::jsonb", string(referersJSON))).
		Where(sq.Eq{"short_code": shortCode}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}

	rofAttached, _ := uQuery.RowsAffected()
	if rofAttached == 0 {
		_, err = sq.Insert("url_stats").
			Columns("short_code", "total_clicks", "unique_visitors", "last_clicked_at", "referers").
			Values(shortCode, totalClicks, uniqueVisitors, lastClick, sq.Expr("?::jsonb", string(referersJSON))).
			PlaceholderFormat(sq.Dollar).
			RunWith(tx).
			Exec()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return tx.Commit()
}

func (repo *eventRepository) GetStats(shortCode string) (*models.Stats, error) {
	statsQuery := sq.Select("total_clicks", "unique_visitors", "last_clicked_at", "referers").
		From("url_stats").
		Where(sq.Eq{"short_code": shortCode}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db)

	stats := &models.Stats{}
	var referersBytes []byte
	err := statsQuery.QueryRow().Scan(&stats.TotalClicks, &stats.UniqueVisitors, &stats.LastClickedAt, &referersBytes)
	if err != nil {
		return nil, err
	}

	if referersBytes != nil {
		if err := json.Unmarshal(referersBytes, &stats.Referers); err != nil {
			stats.Referers = []string{}
		}
	} else {
		stats.Referers = []string{}
	}

	return stats, nil
}

func (repo *eventRepository) UpdateDayStats(shortCode string, date string) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var uniqueVisitors sql.NullInt64
	uniqueQuery := sq.Select("COALESCE(SUM(user_count), 0) AS unique_visitors").
		FromSelect(
			sq.Select("COUNT(DISTINCT user_agent) AS user_count").
				From("clicks").
				Where(sq.Eq{"short_code": shortCode}).
				Where(sq.Eq{"clicked_at::date": date}).
				GroupBy("ip_address"),
			"unique_visitors_ips",
		).
		PlaceholderFormat(sq.Dollar)

	err = uniqueQuery.RunWith(tx).QueryRow().Scan(&uniqueVisitors)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var totalClicks int64
	totalQuery := sq.Select("COUNT(id) AS total_clicks").
		From("clicks").
		Where(sq.Eq{"short_code": shortCode}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx)

	err = totalQuery.QueryRow().Scan(&totalClicks)
	if err != nil {
		fmt.Println(err)
		return err
	}

	referers, err := getReferers(tx, shortCode, date)
	if err != nil {
		return err
	}

	referersJSON, err := json.Marshal(referers)
	if err != nil {
		return err
	}

	uQuery, err := sq.Update("url_day_stats").
		Set("total_clicks", totalClicks).
		Set("unique_visitors", uniqueVisitors).
		Set("referers", sq.Expr("?::jsonb", string(referersJSON))).
		Where(sq.Eq{"short_code": shortCode}).
		Where(sq.Eq{"date": date}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		fmt.Println(err)
		return err
	}

	rofAttached, _ := uQuery.RowsAffected()
	if rofAttached == 0 {
		_, err = sq.Insert("url_day_stats").
			Columns("short_code", "date", "total_clicks", "unique_visitors", "referers").
			Values(shortCode, date, totalClicks, uniqueVisitors, sq.Expr("?::jsonb", string(referersJSON))).
			PlaceholderFormat(sq.Dollar).
			RunWith(tx).
			Exec()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return tx.Commit()
}

func (repo *eventRepository) GetDayStats(shortCode string, date string) (*models.DayStats, error) {
	statsQuery := sq.Select("total_clicks", "unique_visitors", "referers", "updated_at").
		From("url_day_stats").
		Where(sq.Eq{"short_code": shortCode}).
		Where(sq.Eq{"date": date}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db)

	stats := &models.DayStats{}
	var referersBytes []byte
	var updatedAt sql.NullTime
	err := statsQuery.QueryRow().Scan(&stats.TotalClicks, &stats.UniqueVisitors, &referersBytes, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) || (updatedAt.Valid && time.Now().Sub(updatedAt.Time).Minutes() > 60) {
		err = repo.UpdateDayStats(shortCode, date)
		if err != nil {
			return nil, err
		}
		err = statsQuery.QueryRow().Scan(&stats.TotalClicks, &stats.UniqueVisitors, &referersBytes, &updatedAt)
	}
	if err != nil {
		return nil, err
	}

	stats.Date = date
	stats.ShortCode = shortCode
	if referersBytes != nil {
		if err := json.Unmarshal(referersBytes, &stats.Referers); err != nil {
			stats.Referers = []string{}
		}
	} else {
		stats.Referers = []string{}
	}

	return stats, nil
}

func getReferers(runner sq.BaseRunner, shortCode string, date string) ([]string, error) {
	query := sq.Select("DISTINCT referer").
		From("clicks").
		Where(sq.Eq{"short_code": shortCode})

	if date != "" {
		query = query.Where(sq.Eq{"clicked_at::date": date})
	}

	query = query.PlaceholderFormat(sq.Dollar).
		RunWith(runner)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	referers := make([]string, 0)

	for rows.Next() {
		var referer sql.NullString
		err = rows.Scan(&referer)
		if err != nil {
			return nil, err
		}

		if referer.Valid {
			referers = append(referers, referer.String)
		}
	}

	return referers, nil
}
