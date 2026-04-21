# Trelay Implementation Roadmap

This document outlines the planned features, improvements, and "good-to-have" additions for Trelay.

## 1. User Experience (UX) & Dashboard
- [ ] **Interactive Landing Page for Protected Links**: Replace `?p=password` with a dedicated UI for password entry.
- [ ] **Bulk Actions**: Multi-select links for mass move, tag, delete, or restore operations.
- [ ] **Advanced Search & Filtering**: Full-text search (slugs, URLs, tags) and filters for expiry, domain, and date ranges.
- [ ] **Click-to-Copy URL**: One-click shortening URL copying from the link row.
- [ ] **Live Preview Overrides**: Allow manual editing of Open Graph metadata (Title, Description, Image).
- [ ] **Expiration Countdown**: Human-readable countdown for links expiring soon.
- [ ] **Enhanced Keyboard Shortcuts**: Add `/` for search, `m` for move, `t` for tags, etc.

## 2. Analytics & Tracking
- [ ] **Geo-location Tracking**: Integrate GeoIP (MaxMind/IP2Location) for country/city breakdown.
- [ ] **Detailed Device Analytics**: Robust User-Agent parsing for specific Browsers and OS versions.
- [ ] **Real-time "Live" View**: Dashboard section showing clicks as they happen.
- [ ] **UTM Parameter Builder**: Utility in "Create Link" modal for building tracking URLs.
- [ ] **Enhanced Exporting**: PDF reports and more granular CSV/JSON export options.
- [ ] **Referrer Categorization**: Group traffic into Social, Search, Direct, etc.

## 3. Core Functional Additions
- [ ] **Click-based Expiration**: Expire links after a specific number of clicks.
- [ ] **Link Rotator (A/B Testing)**: Single slug redirecting to multiple URLs based on weights.
- [ ] **Alias Support**: Support multiple slugs for a single destination link.
- [ ] **Custom Domains UI**: Dedicated settings page for managing custom domains.
- [ ] **Webhooks**: Trigger external URLs on click events or link expiration.
- [ ] **Scheduled Links**: Set a "start time" for links to become active.

## 4. Security & Privacy
- [ ] **Two-Factor Authentication (2FA)**: Support TOTP for dashboard access.
- [ ] **Audit Logs**: Track administrative actions (creation, deletion, updates).
- [ ] **API Key Scoping**: Create keys with restricted permissions (e.g., read-only).
- [ ] **Granular Privacy Toggles**: Option to disable specific tracking (e.g., referrers) per link.

## 5. Technical & CLI
- [ ] **Caching Layer**: Integrate Redis or in-memory cache for high-traffic redirects.
- [ ] **PostgreSQL Support**: Add support for Postgres as an alternative to SQLite.
- [ ] **Interactive CLI Mode**: Walk users through link creation via `trelay create --interactive`.
- [ ] **Bulk Import**: CLI command to import links from CSV files.
- [ ] **Enhanced Monitoring**: Detailed metrics on `/healthz` (DB health, disk usage, etc.).
