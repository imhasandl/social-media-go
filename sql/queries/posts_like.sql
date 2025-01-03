-- name: LikePost :one
INSERT INTO posts_likes (id, post_id, user_id, created_at)
VALUES (
   $1,
   $2,
   $3,
   NOW()
)
RETURNING *;

-- name: DislikePost :exec
DELETE FROM posts_likes WHERE user_id = $1;

-- name: ListLikePost :many
SELECT * FROM posts_likes;