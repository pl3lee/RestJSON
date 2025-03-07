-- name: CreateNewJson :one
INSERT INTO json_files (user_id, file_name, url)
VALUES($1, $2, $3)
RETURNING *;
