package psql

import (
	"database/sql"
	"fmt"
	"log"

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

	defer db.Close()

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		
		`)
	// CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);

	if err != nil {
		return nil, fmt.Errorf("err1")
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("err2")
	}

	return &Storage{db: db}, nil
}
