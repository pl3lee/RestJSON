-- name: CreateNewJson :one
INSERT INTO json_files (id, user_id, file_name, url)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetJsonFile :one
SELECT *
FROM json_files
WHERE id=$1;

-- name: GetJsonFiles :many
SELECT *
FROM json_files
WHERE user_id=$1;

-- name: RenameJsonFile :one
UPDATE json_files
SET file_name=$2, updated_at=NOW()
WHERE id=$1
RETURNING *;

-- name: DeleteJsonFile :exec
DELETE FROM json_files
WHERE id=$1;
