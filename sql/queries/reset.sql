-- name: ResetUsers :exec
DELETE FROM users;

-- name: ResetPosts :exec
DELETE FROM posts;

-- name: ResetReports :exec
DELETE FROM reports;

-- name: ResetLikePost :exec
DELETE FROM posts_likes;