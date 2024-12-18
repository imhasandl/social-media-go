-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
   $1,
   NOW(),
   NOW(),
   $2,
   $3
)
RETURNING *;

-- name: ListAllUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;