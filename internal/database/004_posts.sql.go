// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: 004_posts.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, body , user_id)
VALUES (gen_random_uuid() , NOW() , NOW() , $1 , $2)
RETURNING id, created_at, updated_at, body, user_id
`

type CreatePostParams struct {
	Body   string
	UserID uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost, arg.Body, arg.UserID)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :one
DELETE FROM posts
WHERE id = $1
RETURNING id, created_at, updated_at, body, user_id
`

func (q *Queries) DeletePost(ctx context.Context, id uuid.UUID) (Post, error) {
	row := q.db.QueryRowContext(ctx, deletePost, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const getAllPosts = `-- name: GetAllPosts :many
SELECT id, created_at, updated_at, body, user_id 
FROM posts
ORDER BY created_at ASC
`

func (q *Queries) GetAllPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPost = `-- name: GetPost :one
SELECT id, created_at, updated_at, body, user_id 
FROM posts
where posts.id = $1
`

func (q *Queries) GetPost(ctx context.Context, id uuid.UUID) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const getPostsByAuthorID = `-- name: GetPostsByAuthorID :many
SELECT id, created_at, updated_at, body, user_id
FROM posts
WHERE posts.user_id = $1
ORDER BY created_at ASC
`

func (q *Queries) GetPostsByAuthorID(ctx context.Context, userID uuid.UUID) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByAuthorID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
