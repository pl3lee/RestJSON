-- name: StoreUserSession :one
INSERT INTO user_sessions (id, user_id, expires_at)
VALUES($1, $2, $3)
RETURNING *;

-- name: GetSession :one
SELECT *
FROM user_sessions
WHERE id=$1;

-- name: InvalidateSession :exec
DELETE FROM user_sessions
WHERE id=$1;

-- name: InvalidateAllSessions :exec
DELETE FROM user_sessions
where user_id=$1;

