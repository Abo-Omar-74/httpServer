// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: 003_users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email , hashed_password)
VALUES (gen_random_uuid() , NOW() , NOW() , $1 , $2)
RETURNING id, created_at, updated_at, email, hashed_password, is_premium
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsPremium,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const editUserByID = `-- name: EditUserByID :one
UPDATE users 
SET email = $2 , hashed_password = $3
where id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_premium
`

type EditUserByIDParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
}

func (q *Queries) EditUserByID(ctx context.Context, arg EditUserByIDParams) (User, error) {
	row := q.db.QueryRowContext(ctx, editUserByID, arg.ID, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsPremium,
	)
	return i, err
}

const findUserByEmail = `-- name: FindUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_premium from users
where email = $1
`

func (q *Queries) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsPremium,
	)
	return i, err
}

const findUserByID = `-- name: FindUserByID :one
SELECT id, created_at, updated_at, email, hashed_password, is_premium from users
where id = $1
`

func (q *Queries) FindUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsPremium,
	)
	return i, err
}

const upgradeUserByID = `-- name: UpgradeUserByID :one
UPDATE users
SET is_premium = TRUE
WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_premium
`

func (q *Queries) UpgradeUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, upgradeUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsPremium,
	)
	return i, err
}
