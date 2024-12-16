-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
   $1,
   NOW(),
   NOW(),
   $2
)
RETURNING id, created_at, updated_at, email;

-- name: ListAllUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;