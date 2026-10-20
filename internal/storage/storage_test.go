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

func TestListChecksWebsiteNewestFirst(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	urls := []string{
		"https://first.example",
		"https://second.example",
		"https://third.example",
	}

	for _, url := range urls {
		err := SaveCheck(
			db,
			url,
			200,
			25*time.Millisecond,
			time.Now().UTC(),
			nil,
		)
		if err != nil {
			t.Fatalf("save check for %s: %v", url, err)
		}
	}

	checks, err := ListChecks(db, 2)
	if err != nil {
		t.Fatalf("list checks, got %v", err)
	}

	if len(checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(checks))
	}

	if checks[0].URL != urls[2] {
		t.Errorf("expected newest URL %q, got %q", urls[2], checks[0].URL)
	}

	if checks[1].URL != urls[1] {
		t.Errorf("expected second newest URL %q, got %q", urls[1], checks[1].URL)
	}
}
