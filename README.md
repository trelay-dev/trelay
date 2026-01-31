# Trelay

[![Go Version](https://img.shields.io/github/go-mod/go-version/trelay-dev/trelay)](https://golang.org)
[![License](https://img.shields.io/github/license/trelay-dev/trelay)](LICENSE)
[![Status](https://img.shields.io/badge/status-active_development-orange)](#)

A developer-first, privacy-respecting URL manager with a robust API and CLI.

## Current Status
This project is currently under active development.

## Features
- URL Shortening with custom slugs and expiration.
- One-time links (consume on first access).
- Folder management for organization.
- Custom domain routing validation.
- Link analytics and statistics (CSV/JSON export).
- Open Graph metadata fetching for link previews.
- Cross-platform CLI with clipboard support.
- Security-focused: JWT/API Key auth, secure headers, rate limiting.

## Prerequisites
- Go 1.21 or higher
- SQLite 3
- Make (optional)

## Quick Start

### Server
1. Clone the repository.
2. Copy `env.example` to `.env` and configure your secrets.
3. Build and run the server:
```bash
make build-server
./bin/trelay-server
```

### CLI
1. Build the CLI:
```bash
make build-cli
```
2. Configure the CLI:
```bash
./bin/trelay config set api-url http://localhost:8080
./bin/trelay config set api-key your-api-key
```

## Documentation
Full documentation for the CLI is available via:
```bash
trelay --help
```

The API structure follows standard REST conventions. Detailed documentation is planned.

## License
Distributed under the MIT License. See `LICENSE` for more information.
