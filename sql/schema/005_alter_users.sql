-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_password VARCHAR(255) DEFAULT 'default_hashed_password';

ALTER TABLE users
ALTER COLUMN hashed_password SET NOT NULL,
ALTER COLUMN hashed_password DROP DEFAULT;

-- +goose Down
ALTER TABLE users
DROP COLUMN hashed_password;
