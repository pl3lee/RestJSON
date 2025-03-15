-- +goose Up
ALTER TABLE api_keys
ADD name TEXT NOT NULL;

-- +goose Down
ALTER TABLE api_keys
DROP COLUMN name;

