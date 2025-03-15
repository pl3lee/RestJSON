-- +goose Up
CREATE TABLE api_keys (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  user_id UUID NOT NULL,
  key_hash TEXT NOT NULL,
  last_used_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user
  FOREIGN KEY (user_id) REFERENCES users(id)
  ON DELETE CASCADE
);

CREATE INDEX api_key_hash ON api_keys(key_hash);

-- +goose Down
DROP INDEX api_key_hash ON api_keys;
DROP TABLE api_keys;

