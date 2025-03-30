package databases

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./link_fizz.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS links (
		id TEXT PRIMARY KEY,
		original_url TEXT NOT NULL,
		short_code TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create table 'links': %v", err)
	}

	createIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_original_url ON links (original_url);`

	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create index for 'original_url': %v", err)
	}

	return db, nil
}
