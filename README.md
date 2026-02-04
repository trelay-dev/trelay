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

## Features

- URL shortening with custom slugs, expiration, and password protection
- One-time links that self-destruct after first access
- Folder management for organizing links
- Custom domain routing
- Click analytics with CSV/JSON export
- Open Graph metadata fetching for link previews
- QR code generation with download and clipboard support
- Trash with soft-delete and restore functionality
- Cross-platform CLI with shell completion
- Security-focused: JWT/API key auth, secure headers, rate limiting
- Single Docker container deployment

## Quick Start

### Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/trelay-dev/trelay.git
cd trelay

# Configure environment
cp env.example .env
# Edit .env with your API_KEY and JWT_SECRET

# Run with Docker Compose
docker compose up -d
```

Access the dashboard at `http://localhost:8080`

### Manual Setup

#### Prerequisites
- Go 1.21+
- Bun 1.0+ (for frontend)
- SQLite 3
- Make (optional)

#### Server
```bash
# Copy and configure environment
cp env.example .env

# Build and run
make build-server
./bin/trelay-server
```

#### Frontend (Development)
```bash
cd frontend
bun install
bun run dev
```

The dev server runs at `http://localhost:5173` with hot reload.

#### CLI
```bash
# Build
make build-cli

# Configure
./bin/trelay config set api-url http://localhost:8080
./bin/trelay config set api-key your-api-key

# Use
./bin/trelay create https://example.com --slug my-link
./bin/trelay list
./bin/trelay qr my-link --open
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `trelay create <url>` | Create a short link |
| `trelay list` | List all links |
| `trelay get <slug>` | Get link details |
| `trelay delete <slug>` | Delete a link |
| `trelay stats <slug>` | View link statistics |
| `trelay qr <slug>` | Generate QR code |
| `trelay folder create <name>` | Create a folder |
| `trelay folder list` | List folders |
| `trelay config set <key> <value>` | Set CLI configuration |
| `trelay completion <shell>` | Generate shell completion |

Run `trelay --help` for full documentation.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/links` | Create link |
| GET | `/api/v1/links` | List links |
| GET | `/api/v1/links/{slug}` | Get link |
| PATCH | `/api/v1/links/{slug}` | Update link |
| DELETE | `/api/v1/links/{slug}` | Delete link |
| POST | `/api/v1/links/{slug}/restore` | Restore deleted link |
| GET | `/api/v1/stats/{slug}` | Get link stats |
| GET | `/api/v1/folders` | List folders |
| POST | `/api/v1/folders` | Create folder |
| GET | `/api/v1/preview?url=` | Fetch Open Graph metadata |
| GET | `/healthz` | Health check |

Authentication: Include `X-API-Key` header with your API key.

## Configuration

Environment variables (set in `.env`):

| Variable | Description | Default |
|----------|-------------|---------|
| `API_KEY` | API authentication key | Required |
| `JWT_SECRET` | JWT signing secret | Required |
| `SERVER_PORT` | HTTP server port | `8080` |
| `DB_PATH` | SQLite database path | `trelay.db` |
| `BASE_URL` | Public URL for short links | `http://localhost:8080` |
| `STATIC_DIR` | Frontend build directory | (empty) |
| `ANALYTICS_ENABLED` | Enable click tracking | `true` |
| `IP_ANONYMIZATION` | Anonymize IP addresses | `true` |
| `RATE_LIMIT_PER_MIN` | API rate limit | `100` |

## License

Distributed under the MIT License. See `LICENSE` for more information.
