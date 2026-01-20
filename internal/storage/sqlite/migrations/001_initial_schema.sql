-- +goose Up
-- Initial schema for Trelay

CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    domain TEXT DEFAULT '',
    password_hash TEXT DEFAULT '',
    expires_at DATETIME,
    tags TEXT DEFAULT '[]',
    folder_id INTEGER,
    click_count INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE INDEX idx_links_slug ON links(slug);
CREATE INDEX idx_links_domain ON links(domain);
CREATE INDEX idx_links_folder_id ON links(folder_id);
CREATE INDEX idx_links_created_at ON links(created_at);
CREATE INDEX idx_links_deleted_at ON links(deleted_at);

CREATE TABLE IF NOT EXISTS clicks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    link_id INTEGER NOT NULL,
    timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    referrer TEXT DEFAULT '',
    device_hash TEXT DEFAULT '',
    user_agent TEXT DEFAULT '',
    ip_hash TEXT DEFAULT '',
    FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
);

CREATE INDEX idx_clicks_link_id ON clicks(link_id);
CREATE INDEX idx_clicks_timestamp ON clicks(timestamp);
CREATE INDEX idx_clicks_referrer ON clicks(referrer);

CREATE TABLE IF NOT EXISTS config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS folders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    parent_id INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES folders(id) ON DELETE CASCADE
);

CREATE INDEX idx_folders_parent_id ON folders(parent_id);

-- +goose Down
DROP INDEX IF EXISTS idx_folders_parent_id;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS config;
DROP INDEX IF EXISTS idx_clicks_referrer;
DROP INDEX IF EXISTS idx_clicks_timestamp;
DROP INDEX IF EXISTS idx_clicks_link_id;
DROP TABLE IF EXISTS clicks;
DROP INDEX IF EXISTS idx_links_deleted_at;
DROP INDEX IF EXISTS idx_links_created_at;
DROP INDEX IF EXISTS idx_links_folder_id;
DROP INDEX IF EXISTS idx_links_domain;
DROP INDEX IF EXISTS idx_links_slug;
DROP TABLE IF EXISTS links;
