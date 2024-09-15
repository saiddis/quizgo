-- name: CreateSession :one
INSERT INTO sessions (
	id,
	created_at,
	quiz_type,
	quiz_category,
	user_id
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionsByUserID :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC;

