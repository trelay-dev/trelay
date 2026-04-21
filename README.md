<p align="center">
  <img src="frontend/static/assets/logo.png" alt="Trelay Logo" width="80" height="80" />
</p>

<h1 align="center">Trelay</h1>

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/trelay-dev/trelay" alt="Go Version" />
  <img src="https://img.shields.io/github/license/trelay-dev/trelay" alt="License" />
  <img src="https://img.shields.io/badge/version-v2.0-blue" alt="Version" />
  <img src="https://img.shields.io/badge/frontend-SvelteKit-ff3e00" alt="Frontend" />
  <img src="https://img.shields.io/badge/runtime-Bun-fbf0df" alt="Runtime" />
</p>

<p align="center">
  A developer-first, privacy-respecting URL manager for self-hosting.<br/>
  Modern web dashboard, powerful CLI, and automation-friendly API.
</p>

<p align="center">
  <sub>Published open source releases are <strong>v2.0+</strong>. There is no public v1; numbering starts here.</sub>
</p>

## Features

- URL shortening with custom slugs, expiration, and password protection
- Password-protected links: HTML password page (or `?p=` as before) for visitors
- One-time links that self-destruct after first access
- Folder management for organizing links
- Custom domain routing
- Click analytics with CSV/JSON export
- Open Graph metadata fetching for link previews, plus optional per-link OG overrides in the dashboard
- QR code generation with download and clipboard support
- Links list: search and filters (tags, domain, created dates, expiry), bulk move/tag/delete, trash bulk restore
- Click the short link slug or the copy control to copy the full short URL; expiry countdown on rows
- Keyboard shortcuts on the links page (`/` search, `s` selection mode, `m` bulk move, `t` bulk tags)
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
| PATCH | `/api/v1/links/bulk` | Bulk update folder and/or tags |
| POST | `/api/v1/links/bulk/restore` | Bulk restore from trash |
| DELETE | `/api/v1/links/{slug}` | Delete link |
| POST | `/api/v1/links/{slug}/restore` | Restore deleted link |
| GET | `/api/v1/stats/{slug}` | Get link stats |
| GET | `/api/v1/folders` | List folders |
| POST | `/api/v1/folders` | Create folder |
| GET | `/api/v1/preview?url=` | Fetch Open Graph metadata |
| GET | `/healthz` | Health check |

Authentication: Include `X-API-Key` header with your API key.

## Roadmap

Planned features and improvements for releases after **v2.0** are tracked in [`ROADMAP.md`](ROADMAP.md) (UX, analytics, core features, security, and platform work).

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for how we take patches and issues.

## Security

Found something sensitive? Do not use the public issue tracker. Read [`SECURITY.md`](SECURITY.md) and use GitHub’s private reporting under the **Security** tab.

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
