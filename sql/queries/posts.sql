-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, user_id, body, likes)
VALUES (
   $1,
   NOW(),
   NOW(),
   $2,
   $3,
   $4
)
RETURNING id, created_at, updated_at, user_id, body, likes;

-- name: GetPosts :many
SELECT * FROM posts;

-- name: GetPostByID :one
SELECT * FROM posts
WHERE id = $1;

-- name: ChangePostByID :exec
UPDATE posts SET
body = $1, updated_at = NOW()
WHERE id = $2;

-- name: DeletePostByID :exec
DELETE FROM posts WHERE id = $1;

-- name: GetMostLikedPosts :many
SELECT * FROM posts
ORDER BY likes ASC LIMIT 10;