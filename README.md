# Website Uptime Monitor

A Go command-line tool that checks website availability, measures response time, and stores check history in SQLite.

Built to practise Go HTTP clients, error handling, package organization, database access, and automated testing.

## Features

- Check a website using a configurable URL and timeout.
- Report HTTP status, response time, and check timestamp.
- Save successful and failed checks in SQLite.
- Display the 10 most recently saved checks.
- Test HTTP checking behavior and database persistence.

## Requirements

- Go 1.23.1 or later

The SQLite driver runs within the application. No separate database server or SQLite CLI installation is required.

## Getting Started

Clone the repository:

```bash
git clone https://github.com/vitamaxoduol/Website-Uptime-Monitor.git
cd Website-Uptime-Monitor
```

Check a website:

```bash
go run ./cmd/uptime-monitor -url https://go.dev -timeout 5s
```

Example output:

```text
URL: https://go.dev | Status: 200 | Response time: 749.738538ms | Checked at: 2026-10-18T17:14:34Z
```

Response times and timestamps vary between runs.

## Command-Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-url` | `https://www.example.com` | Website URL to check |
| `-timeout` | `10s` | Request timeout; must be greater than zero |
| `-history` | `false` | Display the last 10 saved checks without running a new check |
| `-add-monitor` | `""` | Save a website URL as a monitor |
| `-interval` | `60` | Monitor check interval in seconds; must be positive |
| `-list-monitors` | `false` | List saved monitors |
| `-delete-monitor` | `0` | Delete a monitor by its positive ID |

Display help:

```bash
go run ./cmd/uptime-monitor -help
```

## Check History

Display saved results, most recently inserted first:

```bash
go run ./cmd/uptime-monitor -history
```

Failed checks include their error message. Status `0` means no HTTP response was received.

The application creates `uptime.db` in the current working directory and initializes its table automatically. Run commands from the repository root to consistently use the same database.

Local database files are excluded from Git.

## Managing Monitors

Create a monitor:

```bash
go run ./cmd/uptime-monitor -add-monitor https://go.dev -interval 60
```

List monitors:

```bash
go run ./cmd/uptime-monitor -list-monitors
```

Delete a monitor using its ID from the list:

```bash
go run ./cmd/uptime-monitor -delete-monitor 1
```

Use only one action per run: `-add-monitor`, `-list-monitors`,
`-delete-monitor`, or `-history`.

Monitor intervals are stored for future scheduling; automatic checks are
not implemented yet. Deleting a monitor preserves existing check history.

## Running Tests

Run all tests:

```bash
go test -v ./...
```

Tests cover:

- Successful HTTP responses.
- Unexpected HTTP statuses.
- Invalid URLs.
- Check timestamps.
- Saving and reading failure details.
- History ordering and result limits.

HTTP tests use temporary local servers. Storage tests use temporary databases, leaving the application's history untouched.

## Project Structure

```text
cmd/
└── uptime-monitor/
    └── main.go
internal/
├── checker/
│   ├── checker.go
│   └── checker_test.go
└── storage/
    ├── schema.sql
    ├── storage.go
    └── storage_test.go
go.mod
go.sum
README.md
```

- `cmd/uptime-monitor`: Command-line flags and application flow.
- `internal/checker`: HTTP checks and response timing.
- `internal/storage`: SQLite initialization, persistence, and history queries.

## Current Behavior

- Each normal run performs one website check.
- Only HTTP `200 OK` is considered healthy.
- HTTP redirects follow Go's default client behavior.
- Response time measures the request until response headers arrive, including connection setup. It does not measure downloading the full response body.
- Check timestamps are stored in UTC.
- Checks are saved when database writes succeed, including checks that fail.
- Monitoring currently runs on demand; there is no background scheduler or HTTP API yet.

## Planned Improvements

- HTTP API for monitor management and check history.
- Scheduled background checks.
- Bounded concurrency and graceful shutdown.
- Docker support.
- GitLab CI.