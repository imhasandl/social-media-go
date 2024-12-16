-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, user_id, body)
VALUES (
   $1,
   NOW(),
   NOW(),
   $2,
   $3
)
RETURNING id, created_at, updated_at, user_id, body;

-- name: GetPosts :many
SELECT * FROM posts;