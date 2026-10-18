package storage

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql

var schema string

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return db, nil

}

func SaveCheck(
	db *sql.DB,
	url string,
	statusCode int,
	responseTime time.Duration,
	checkedAt time.Time,
	checkErr error,
) error {
	errorMessage := ""
	if checkErr != nil {
		errorMessage = checkErr.Error()
	}

	_, err := db.Exec(
		`INSERT INTO checks (
		url,
		status_code,
		response_time_ns,
		checked_at,
		error_message
		) VALUES (?, ?, ?, ?, ?)`,
		url,
		statusCode,
		int64(responseTime),
		checkedAt.UTC().Format(time.RFC3339Nano),
		errorMessage,
	)
	if err != nil {
		return fmt.Errorf("save check: %w", err)
	}

	return nil
}
