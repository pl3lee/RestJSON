-- +goose Up
ALTER TABLE users
ADD subscribed BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE users
DROP COLUMN subscribed;
