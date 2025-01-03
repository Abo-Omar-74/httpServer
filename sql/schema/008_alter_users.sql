-- +goose Up
ALTER TABLE users
ADD COLUMN is_premium BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_premium;