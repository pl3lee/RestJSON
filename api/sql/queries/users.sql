-- name: CreateUser :one
INSERT INTO users(provider_id, email, name)
VALUES ($1, $2, $3)
ON CONFLICT(provider_id)
DO UPDATE SET email=$2, name=$3
RETURNING *;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id=$1;

-- name: GetUserByProviderId :one
SELECT *
FROM users
WHERE provider_id=$1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id=$1;

-- name: UpdateCustomerId :one
UPDATE users
SET stripe_customer_id=$2, updated_at=NOW()
WHERE id=$1
RETURNING *;

-- name: GetUserByCustomerId :one
SELECT *
FROM users
WHERE stripe_customer_id=$1;
