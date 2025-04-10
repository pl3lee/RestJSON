// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(provider_id, email, name)
VALUES ($1, $2, $3)
ON CONFLICT(provider_id)
DO UPDATE SET email=$2, name=$3
RETURNING id, provider_id, created_at, updated_at, email, name
`

type CreateUserParams struct {
	ProviderID string
	Email      string
	Name       string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ProviderID, arg.Email, arg.Name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Name,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id=$1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserById = `-- name: GetUserById :one
SELECT id, provider_id, created_at, updated_at, email, name
FROM users
WHERE id=$1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Name,
	)
	return i, err
}

const getUserByProviderId = `-- name: GetUserByProviderId :one
SELECT id, provider_id, created_at, updated_at, email, name
FROM users
WHERE provider_id=$1
`

func (q *Queries) GetUserByProviderId(ctx context.Context, providerID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByProviderId, providerID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.ProviderID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Name,
	)
	return i, err
}
