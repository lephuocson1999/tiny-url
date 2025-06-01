package app

import (
	"context"
	"database/sql"
	"time"

	"tiny-url/internal/domain"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(ctx context.Context, url *domain.URL) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO urls (id, short_code, original_url, created_at, expires_at, access_count, last_accessed)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (short_code) DO UPDATE SET original_url = EXCLUDED.original_url
	`, url.ID, url.ShortCode, url.OriginalURL, url.CreatedAt, url.ExpiresAt, url.AccessCount, url.LastAccessed)
	return err
}

func (r *PostgresRepo) FindByCode(ctx context.Context, code string) (*domain.URL, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, short_code, original_url, created_at, expires_at, access_count, last_accessed
		FROM urls WHERE short_code = $1
	`, code)
	var u domain.URL
	var lastAccessed sql.NullTime
	if err := row.Scan(&u.ID, &u.ShortCode, &u.OriginalURL, &u.CreatedAt, &u.ExpiresAt, &u.AccessCount, &lastAccessed); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	if lastAccessed.Valid {
		u.LastAccessed = &lastAccessed.Time
	}
	if u.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrExpired
	}
	return &u, nil
}
