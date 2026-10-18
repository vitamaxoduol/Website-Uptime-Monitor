CREATE TABLE IF NOT EXISTS checks (
    id INTEGER PRIMARY KEY,
    url TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    response_time_ns NOT NULL,
    checked_at TEXT NOT NULL,
    error_message TEXT NOT NULL DEFAULT ''
);