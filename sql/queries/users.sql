-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, username, password)
VALUES (
   $1,
   NOW(),
   NOW(),
   $2,
   $3,
   $4
)
RETURNING *;

-- name: ChangeUser :one
UPDATE users SET email = $1, updated_at = NOW(), password = $2
WHERE id = $3
RETURNING *;

-- name: ListAllUsers :many
SELECT id, created_at, updated_at, email, username FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: CheckIfUsernameOrEmailTaken :one
SELECT id from users
WHERE username = $1 or email = $2;

-- name: UpgradeToPremium :one
UPDATE users SET is_premium = true, updated_at = NOW()
WHERE id = $1
RETURNING *;
