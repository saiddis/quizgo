-- name: CreateScore :one
INSERT INTO scores (
	id,
	completion_time,
	hard_quizzes_done,
	medium_quizzes_done,
	easy_quizzes_done, 
	total_score, 
	user_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetScoresByUserID :many
SELECT * FROM scores
WHERE user_id = $1;
