-- +goose Up
ALTER TABLE users
ADD COLUMN username TEXT NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN username;
