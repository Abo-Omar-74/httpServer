-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, body , user_id)
VALUES (gen_random_uuid() , NOW() , NOW() , $1 , $2)
RETURNING *;

-- name: GetAllPosts :many
SELECT * 
FROM posts
ORDER BY created_at ASC;

-- name: GetPost :one
SELECT * 
FROM posts
where posts.id = $1;

-- name: DeletePost :one
DELETE FROM posts
WHERE id = $1
RETURNING *;

-- name: GetPostsByAuthorID :many
SELECT *
FROM posts
WHERE posts.user_id = $1
ORDER BY created_at ASC;