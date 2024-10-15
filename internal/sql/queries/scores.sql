-- name: CreateScore :one
INSERT INTO scores (
	completion_time,
	hard_quizzes_done,
	medium_quizzes_done,
	easy_quizzes_done, 
	total_score, 
	user_id
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING scores.id;

-- name: GetScoreByID :one
SELECT * FROM scores
WHERE id = $1;

-- name: GetScoresByUserID :many
SELECT * FROM scores
WHERE user_id = $1;
