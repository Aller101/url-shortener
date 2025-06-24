package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"url-shortener/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(connStr string) (*Storage, error) {

	const op = "storage.psql.New"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка ping: %v", err)
	}

	// defer db.Close()

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
			);
			
			`)
	// CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.psql.SaveURL"
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = stmt.QueryRow(urlToSave, alias).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.psql.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=$1")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.psql.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=$1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	err = stmt.QueryRow(alias).Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrURLNotFound
		}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}
