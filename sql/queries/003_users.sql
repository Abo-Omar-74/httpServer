-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email , hashed_password)
VALUES (gen_random_uuid() , NOW() , NOW() , $1 , $2)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: FindUserByID :one
SELECT * from users
where id = $1;

-- name: FindUserByEmail :one
SELECT * from users
where email = $1;

-- name: EditUserByID :one
UPDATE users 
SET email = $2 , hashed_password = $3
where id = $1
RETURNING *;

-- name: UpgradeUserByID :one
UPDATE users
SET is_premium = TRUE
WHERE id = $1
RETURNING *;