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

type StoredCheck struct {
	ID             int64
	URL            string
	StatusCode     int
	ResponseTimeNS int64
	CheckedAt      string
	ErrorMessage   string
}

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

func ListChecks(db *sql.DB, limit int) ([]StoredCheck, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must greater that zero")
	}

	rows, err := db.Query(
		`SELECT id, url, status_code, response_time_ns, checked_at, error_message
		FROM checks
		ORDER BY id DESC
		LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list checks: %w", err)
	}
	defer rows.Close()

	checks := make([]StoredCheck, 0)

	for rows.Next() {
		var check StoredCheck

		err := rows.Scan(
			&check.ID,
			&check.URL,
			&check.StatusCode,
			&check.ResponseTimeNS,
			&check.CheckedAt,
			&check.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("Scan check: %w", err)
		}

		checks = append(checks, check)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read checks: %w", err)
	}
	return checks, nil
}
