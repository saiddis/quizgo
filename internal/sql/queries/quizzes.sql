-- name: CreateQuiz :one
INSERT INTO quizzes (
	id,
	created_at,
	quiz_type,
	quiz_category,
	user_id
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetQuizByUserID :many
SELECT * FROM quizzes
WHERE user_id = $1
ORDER BY created_at DESC;

