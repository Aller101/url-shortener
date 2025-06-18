package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"url_shortener/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);

	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(ctx context.Context, url, alias string) (int64, error) {
	const op = "sqlite.SaveURL"

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO url(url, alias) VALUES(?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, url, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) SaveURLWithTx(ctx context.Context, url, alias string) (int64, error) {
	const op = "sqlite.SaveURLWithTx"

	var id int64

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := stmt.QueryRowContext(ctx, stmt, url, alias).Scan(&id); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	tx.Commit()
	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string
	if err := stmt.QueryRow(alias).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil

}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	const op = "sqlite.DeleteURL"

	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM url WHERE alias=?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.ExecContext(ctx, alias); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrURLNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
