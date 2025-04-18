-- name: CreateApiKey :one
INSERT INTO api_keys(user_id, key_hash, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserFromApiKeyHash :one
SELECT *
FROM api_keys
WHERE key_hash=$1;

-- name: GetAllApiKeys :many
SELECT *
FROM api_keys
WHERE user_id=$1;

-- name: UpdateApiKeyLastUsed :exec
UPDATE api_keys
SET updated_at=NOW(), last_used_at=NOW()
WHERE key_hash=$1;

-- name: DeleteApiKey :exec
DELETE FROM api_keys
WHERE key_hash=$1;

