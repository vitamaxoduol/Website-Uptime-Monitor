-- Stores what happened each time a check ran
CREATE TABLE IF NOT EXISTS checks (
    id INTEGER PRIMARY KEY,
    url TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    response_time_ns NOT NULL,
    checked_at TEXT NOT NULL,
    error_message TEXT NOT NULL DEFAULT ''
);

-- Stores what to check: a website URL and its check intervaL
CREATE TABLE IF NOT EXISTS monitors (
    id INTEGER PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    interval_seconds INTEGER NOT NULL CHECK (interval_seconds > 0),
    created_at TEXT NOT NULL
);