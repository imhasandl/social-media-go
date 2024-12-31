-- name: ReportPost :one
INSERT INTO reports (report_id, created_at, updated_at, post_id, user_id, reason)
VALUES(
   $1,
   NOW(),
   NOW(),
   $2,
   $3,
   $4
)
RETURNING *;

-- name: DeleteReportByID :exec
DELETE FROM reports 
WHERE report_id = $1;

-- name: ListAllReports :many
SELECT * FROM reports;
