# Roadmap

Rough plan for what might land after the current **v2.x** line. Nothing here is a promise or a deadline; it is a backlog we use to stay aligned.

## 1. User experience and dashboard

- [x] **Interactive landing page for protected links**: replace `?p=password` with a small password UI.
- [x] **Bulk actions**: multi-select for move, tag, delete, restore.
- [x] **Search and filters**: slugs, URLs, tags, expiry, domain, dates.
- [x] **Click-to-copy** short URL from a row.
- [x] **Manual OG overrides** for title, description, image on previews.
- [x] **Expiry countdown** for links that are about to expire.
- [x] **More keyboard shortcuts** (e.g. `/` search, `m` move, `t` tags).

## 2. Analytics and tracking

- [ ] **GeoIP** (e.g. MaxMind / IP2Location) for country or city.
- [ ] **Finer device breakdown** from User-Agent (browser/OS versions).
- [ ] **Live-ish view** of recent clicks on the dashboard.
- [ ] **UTM helper** in the create flow.
- [ ] **Richer exports** (PDF, finer CSV/JSON).
- [ ] **Referrer buckets** (social, search, direct, etc.).

## 3. Core product

- [ ] **Click-based expiry** (expire after N clicks).
- [ ] **Link rotator / A/B** (weighted destinations).
- [ ] **Aliases** (several slugs, one target).
- [ ] **Custom domains UI** (settings page).
- [ ] **Webhooks** on click or expiry.
- [ ] **Scheduled links** (active from a start time).

## 4. Security and privacy

- [ ] **2FA** (TOTP) for the dashboard.
- [ ] **Audit log** for admin actions.
- [ ] **Scoped API keys** (read-only, etc.).
- [ ] **Per-link privacy toggles** (e.g. turn off referrer storage).

## 5. Platform and CLI

- [ ] **Cache layer** (Redis or in-process) for hot redirects.
- [ ] **PostgreSQL** as an alternative to SQLite.
- [ ] **Interactive CLI** (`trelay create --interactive`).
- [ ] **Bulk import** from CSV in the CLI (beyond what the API already does).
- [ ] **Richer `/healthz`** (DB, disk, etc.).
