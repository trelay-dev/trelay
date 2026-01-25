-- +goose Up
ALTER TABLE links ADD COLUMN is_one_time BOOLEAN DEFAULT 0;

-- +goose Down
ALTER TABLE links DROP COLUMN is_one_time;
