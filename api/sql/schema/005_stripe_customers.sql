-- +goose Up
ALTER TABLE users
ADD stripe_customer_id TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE users
DROP COLUMN stripe_customer_id;
