<p align="center">
  <img src="frontend/static/assets/logo.png" alt="Trelay Logo" width="80" height="80" />
</p>

<h1 align="center">Trelay</h1>

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/trelay-dev/trelay" alt="Go Version" />
  <img src="https://img.shields.io/github/license/trelay-dev/trelay" alt="License" />
  <img src="https://img.shields.io/badge/status-active_development-orange" alt="Status" />
  <img src="https://img.shields.io/badge/frontend-SvelteKit-ff3e00" alt="Frontend" />
  <img src="https://img.shields.io/badge/runtime-Bun-fbf0df" alt="Runtime" />
</p>

<p align="center">
  A developer-first, privacy-respecting URL manager with a robust API and CLI.
</p>

## Current Status
This project is currently under active development.

## Features
- URL Shortening with custom slugs and expiration.
- Web Dashboard for managing links and folders.
- One-time links (consume on first access).
- Folder management for organization.
- Custom domain routing validation.
- Link analytics and statistics (CSV/JSON export).
- Open Graph metadata fetching for link previews.
- Cross-platform CLI with clipboard support.
- Security-focused: JWT/API Key auth, secure headers, rate limiting.

## Prerequisites
- Go 1.21 or higher
- Bun 1.0 or higher (for frontend)
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

### Web Dashboard
1. Navigate to the frontend directory and install dependencies:
```bash
cd frontend
bun install
```
2. Start the development server:
```bash
bun run dev
```
The dashboard will be available at `http://localhost:5173`.

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
