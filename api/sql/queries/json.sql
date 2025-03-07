-- name: CreateNewJson :one
INSERT INTO json_files (id, user_id, file_name, url)
VALUES($1, $2, $3, $4)
RETURNING *;
