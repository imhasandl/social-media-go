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
DELETE FROM posts_likes 
WHERE user_id = $1 AND post_id = $2;

-- name: ListLikePost :many
SELECT * FROM posts_likes;

-- name: CheckIfUserLikeAlready :exec
SELECT id FROM posts_likes
WHERE user_id = $1;

-- name: IncrementPostLike :exec
UPDATE posts SET likes = likes + 1
WHERE id = $1;

-- name: DecrementPostLike :exec
UPDATE posts SET likes = likes - 1
WHERE id = $1;