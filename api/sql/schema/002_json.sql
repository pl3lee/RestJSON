-- +goose Up
CREATE TABLE json_files (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  user_id UUID NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  url TEXT NOT NULL,
  CONSTRAINT fk_user
  FOREIGN KEY (user_id) REFERENCES users(id)
  ON DELETE CASCADE
);

-- +goose Down
DROP TABLE json_files;
