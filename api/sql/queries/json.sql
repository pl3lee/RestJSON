-- name: CreateNewJson :one
INSERT INTO json_files (id, user_id, file_name, url)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetJsonFile :one
SELECT *
FROM json_files
WHERE id=$1;
