package storage

import (
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveCheckFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	db, err := Open(path)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	checkedAt := time.Now().UTC()
	checkErr := errors.New("connection failed")

	err = SaveCheck(
		db,
		"https://example.com",
		0,
		25*time.Millisecond,
		checkedAt,
		checkErr,
	)
	if err != nil {
		t.Fatalf("save check: %v", err)
	}

	var url string
	var statusCode int
	var responseTimeNS int64
	var storedTime string
	var errorMessage string

	err = db.QueryRow(
		`SELECT url, status_code, response_time_ns, checked_at, error_message
		FROM checks`,
	).Scan(
		&url,
		&statusCode,
		&responseTimeNS,
		&storedTime,
		&errorMessage,
	)
	if err != nil {
		t.Fatalf("read saved check: %v", err)
	}

	if url != "https://example.com" {
		t.Errorf("unexpected URL: %q", url)
	}
	if statusCode != 0 {
		t.Errorf("expected status 0, got %d", statusCode)
	}
	if responseTimeNS != int64(25*time.Millisecond) {
		t.Errorf("unexpected response time: %d", responseTimeNS)
	}
	if storedTime != checkedAt.Format(time.RFC3339Nano) {
		t.Errorf("unexpected timestamp: %q", storedTime)
	}
	if errorMessage != checkErr.Error() {
		t.Errorf("unexpected error message: %q", errorMessage)
	}
}
