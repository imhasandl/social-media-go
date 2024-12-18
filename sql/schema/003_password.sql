-- +goose Up
ALTER TABLE users
ADD COLUMN IF NOT EXISTS password TEXT NOT NULL
DEFAULT 'unset'

-- +goose Down
ALTER TABLE users
DROP COLUMN IF EXISTS password;
