-- +goose Up
ALTER TABLE links ADD COLUMN og_title TEXT DEFAULT '';
ALTER TABLE links ADD COLUMN og_description TEXT DEFAULT '';
ALTER TABLE links ADD COLUMN og_image_url TEXT DEFAULT '';

-- +goose Down
-- SQLite cannot DROP COLUMN in older versions; recreate is heavy — leave columns on rollback in dev.
