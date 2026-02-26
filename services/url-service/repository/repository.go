package repository

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ArtemNeGopher/url-shortener/services/url-service/config"
	"github.com/ArtemNeGopher/url-shortener/services/url-service/grpc"
	"github.com/ArtemNeGopher/url-shortener/services/url-service/models"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type urlRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewURLRepository(config *config.DatabaseConfig, log *slog.Logger) *urlRepository {
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

	return &urlRepository{
		db:  db,
		log: log.With(slog.String("context", "repository")),
	}
}

var _ grpc.URLRepository = (*urlRepository)(nil)

func (repo *urlRepository) TryRegisterURL(url *models.URL) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем существует ли запись с таким short_code
	// Если существует, то возвращаем ошибку
	query, err := sq.Select("id").
		From("urls").
		Where(sq.Eq{"short_code": url.ShortCode}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		return err
	}
	c, err := query.RowsAffected()
	if err != nil {
		return err
	}
	if c > 0 {
		return errors.New("short code does exist")
	}

	// Вставляем новый url
	q := sq.Insert("urls")

	if url.ExpiresAt != nil {
		q = q.Columns("short_code", "original_url", "expires_at").
			Values(url.ShortCode, url.URL, url.ExpiresAt)
	} else {
		q = q.Columns("short_code", "original_url").
			Values(url.ShortCode, url.URL)
	}
	_, err = q.PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *urlRepository) GetOriginalURL(shortCode string) (*models.URL, error) {
	query := sq.Select("original_url", "expires_at", "is_active").
		From("urls").
		Where(sq.Eq{"short_code": shortCode}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		QueryRow()

	ret := &models.URL{
		ShortCode: shortCode,
	}
	err := query.Scan(&ret.URL, &ret.ExpiresAt, &ret.IsActive)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (repo *urlRepository) DeleteURL(shortCode string) error {
	query, err := sq.Delete("urls").
		Where(sq.Eq{"short_code": shortCode}).
		PlaceholderFormat(sq.Dollar).
		RunWith(repo.db).
		Exec()
	if err != nil {
		return err
	}

	c, err := query.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New("url not found")
	}

	return nil
}

func (repo *urlRepository) Close() {
	repo.db.Close()
}
