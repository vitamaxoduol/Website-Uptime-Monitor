package storage

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"
)

type Monitor struct {
	ID              int64
	URL             string
	IntervalSeconds int
	CreatedAt       string
}

// CreateMonitor save a website and its check interval
func CreateMonitor(db *sql.DB, websiteURL string, intervalSeconds int) (int64, error) {
	parseURL, err := url.Parse(websiteURL)
	if err != nil {
		return 0, fmt.Errorf("parse monitor URL: %w", err)
	}

	if (parseURL.Scheme != "http" && parseURL.Scheme != "https") ||
		parseURL.Hostname() == "" {
		return 0, fmt.Errorf("monitor URL must use http or https and include a hostname")
	}

	if intervalSeconds <= 0 {
		return 0, fmt.Errorf("monitor interval seconds must be greater than zero")
	}

	result, err := db.Exec(
		`INSERT INTO monitors (url, interval_seconds, created_at)
		VALUES (?, ?, ?)`,
		websiteURL,
		intervalSeconds,
		time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return 0, fmt.Errorf("create: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get monitor ID: %w", err)
	}

	return id, nil
}

// ListMonitors returns all monitors, newest inserted first
func ListMonitors(db *sql.DB) ([]Monitor, error) {
	rows, err := db.Query(
		`SELECT id, url, interval_seconds, created_at
		FROM monitors
		ORDER BY id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list monitors: %w", err)
	}
	defer rows.Close()

	monitors := make([]Monitor, 0)

	for rows.Next() {
		var monitor Monitor

		if err := rows.Scan(
			&monitor.ID,
			&monitor.URL,
			&monitor.IntervalSeconds,
			&monitor.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan monitor: %w", err)
		}

		monitors = append(monitors, monitor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read monitors: %w", err)
	}

	return monitors, nil
}

// DeleteMonitor does delete a monitor by its ID
func DeleteMonitor(db *sql.DB, id int64) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("monitor ID must greater than zero")
	}

	result, err := db.Exec(
		`DELETE FROM monitors WHERE id = ?`,
		id,
	)
	if err != nil {
		return false, fmt.Errorf("delete monitor: %w", err) // Validation or the database operation failed
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("get deleted row count: %w", err)
	}

	return affected > 0, nil
}
