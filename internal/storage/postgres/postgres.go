package postgres

import (
	"context"

	"fmt"

	"url_shortner/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(context.Background(), `
			CREATE TABLE IF NOT EXISTS url (
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil

}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	var id int64
	err := s.db.QueryRow(
		context.Background(),
		`INSERT INTO url(alias, url) VALUES($1, $2)
		 ON CONFLICT (alias) DO UPDATE SET url = EXCLUDED.url
		 RETURNING id`,
		alias,
		urlToSave,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.SaveURL"

	var url string
	err := s.db.QueryRow(context.Background(), `SELECT url FROM url WHERE alias = $1`, alias).Scan(&url)
	if err != nil {
		if err != pgx.ErrNoRows {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil

}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	_, err := s.db.Exec(context.Background(), `delete from url where alias = $1`, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (s *Storage) UpdateURL(alias, urlToSave string) (int64, error) {
	const op = "storage.postgres.UpdateURL"

	result, err := s.db.Exec(context.Background(), `update url set url = $1 where alias = $2`, urlToSave, alias)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected := result.RowsAffected()
	return rowsAffected, nil

}

func (s *Storage) Close() {
	s.db.Close()
}
