# Website Uptime Monitor

A Go command-line tool that checks whether a website responds with HTTP `200 OK` and measures its response time.

Built as a hands-on project to practise Go HTTP clients, error handling, packages, and automated testing.

## Features

- Check a website using a configurable URL.
- Set a request timeout.
- Report the HTTP status and response time for successful checks.
- Report connection errors and unexpected HTTP statuses.
- Test checker behavior using local HTTP test servers.

## Requirements

- Go 1.23.1 or later

## Getting Started

Clone the repository:

```bash
git clone https://github.com/vitamaxoduol/Website-Uptime-Monitor.git
cd Website-Uptime-Monitor
```

Run a check:

```bash
go run ./cmd/uptime-monitor -url https://go.dev -timeout 5s
```

Example output:

```text
URL: https://go.dev | Status: 200 | Response time: 243.817ms
```

Response times vary between requests.

### Command-Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-url` | `https://www.example.com` | Website URL to check |
| `-timeout` | `10s` | Request timeout; must be greater than zero |

Display help:

```bash
go run ./cmd/uptime-monitor -help
```

## Running Tests

```bash
go test -v ./...
```

Tests cover:

- A successful `200 OK` response.
- A `500 Internal Server Error` response.
- An invalid URL.

## Project Structure

```text
cmd/
└── uptime-monitor/
    └── main.go          # CLI entry point
internal/
├── checker.go           # Website checking logic
└── checker_test.go      # Checker tests
go.mod
README.md
```

## Current Behavior

The tool performs one check per run. Only HTTP `200 OK` is considered healthy.

Response time measures the request until response headers arrive, including any connection setup. It does not measure downloading the entire response body.

## Planned Improvements

- SQLite storage for monitored URLs and check history.
- HTTP API for managing monitors and viewing results.
- Scheduled background checks with bounded concurrency.
- Docker support.
- GitLab CI.