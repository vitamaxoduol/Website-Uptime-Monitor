package storage

import (
	"path/filepath"
	"testing"
)

func TestCreateMonitor(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	id, err := CreateMonitor(db, "https://go.dev", 60)
	if err != nil {
		t.Fatalf("create monitor: %v", err)
	}

	if id <= 0 {
		t.Fatalf("expected a positive monitor ID, got %d", id)
	}

	var storedURL string
	var storedInterval int

	err = db.QueryRow(
		`SELECT url, interval_seconds FROM monitors WHERE id= ?`,
		id,
	).Scan(&storedURL, &storedInterval)
	if err != nil {
		t.Fatalf("read monitor: %v", err)
	}

	if storedURL != "https://go.dev" {
		t.Errorf("expected URL https://go.dev, got %q", storedURL)
	}

	if storedInterval != 60 {
		t.Errorf("expected interval 60, got %d", storedInterval)
	}

	_, err = CreateMonitor(db, "https://go.dev", 120)
	if err == nil {
		t.Fatalf("expected an error for duplIcate URL, got nil")
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM monitors`).Scan(&count)
	if err != nil {
		t.Fatalf("count monitors: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 monitor after duplicate attempt, got %d", count)
	}
}

// This test verifies:
// An empty table returns no monitors and no error.
// Both saved monitors are returned.
// The newest inserted monitor comes first.
// Its URL and interval are read correctly
func TestListMonitors(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	monitors, err := ListMonitors(db)
	if err != nil {
		t.Fatalf("list empty database: %v", err)
	}

	if len(monitors) != 0 {
		t.Fatalf("expected no monitors, got %d", len(monitors))
	}

	firstID, err := CreateMonitor(db, "https://go.dev", 60)
	if err != nil {
		t.Fatalf("create first monitor: %v", err)
	}

	secondID, err := CreateMonitor(db, "https://example.com", 120)
	if err != nil {
		t.Fatalf("create second monitor: %v", err)
	}

	monitors, err = ListMonitors(db)
	if err != nil {
		t.Fatalf("list monitors: %v", err)
	}

	if len(monitors) != 2 {
		t.Fatalf("expected 2 monitors, got %d", len(monitors))
	}

	if monitors[0].ID != secondID {
		t.Errorf("expected newest ID %d, got %d", secondID, monitors[0].ID)
	}

	if monitors[1].ID != firstID {
		t.Errorf("expected oldest ID %d, got %d", firstID, monitors[1].ID)
	}

	if monitors[0].URL != "https://example.com" {
		t.Errorf("unexpected newest URL: %q", monitors[0].URL)
	}

	if monitors[0].IntervalSeconds != 120 {
		t.Errorf("expected interval 120, got %d", monitors[0].IntervalSeconds)
	}
}

func TestDeleteMonitor(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	deleteID, err := CreateMonitor(db, "https://go.dev", 60)
	if err != nil {
		t.Fatalf("create monitor to delete: %v", err)
	}

	keepID, err := CreateMonitor(db, "https://example.com", 120)
	if err != nil {
		t.Fatalf("create monitor to keep: %v", err)
	}

	deleted, err := DeleteMonitor(db, deleteID)
	if err != nil {
		t.Fatalf("delete monitor: %v", err)
	}
	if !deleted {
		t.Fatal("expected monitor to be deleted")
	}

	monitors, err := ListMonitors(db)
	if err != nil {
		t.Fatalf("list remaining monitors: %v", err)
	}
	if len(monitors) != 1 {
		t.Fatalf("expected 1 remaining monitor, got %d", len(monitors))
	}
	if monitors[0].ID != keepID {
		t.Errorf("expected remaining ID %d, got %d", keepID, monitors[0].ID)
	}

	deleted, err = DeleteMonitor(db, deleteID)
	if err != nil {
		t.Fatalf("delete nonexistent monitor: %v", err)
	}
	if deleted {
		t.Fatal("expected false when deleting a nonexistent monitor")
	}
}
