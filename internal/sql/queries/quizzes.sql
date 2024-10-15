-- name: CreateQuiz :one
INSERT INTO quizzes (
	created_at,
	type,
	category,
	user_id
)
VALUES ($1, $2, $3, $4)
RETURNING quizzes.id;

-- name: GetQuizzesByUserID :many
SELECT * FROM quizzes
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateScoreID :one
UPDATE quizzes
SET score_id = $1
WHERE id = $2
RETURNING *;

-- name: GetLastQuizByUserID :one
SELECT * FROM quizzes
JOIN users ON quizzes.user_id = users.id
WHERE quizzes.user_id = $1
ORDER BY quizzes.id DESC
LIMIT 1;

-- name: PaginateQuizzes :many
SELECT quizzes.created_at, quizzes.type, quizzes.category, quizzes.score_id as score_id, quizzes.id as id FROM quizzes
JOIN users ON quizzes.user_id = users.id
WHERE quizzes.user_id = $1 AND quizzes.id < $2
ORDER BY quizzes.id DESC
LIMIT 5;

-- name: GetUserIDByQuizID :one
SELECT users.id as user_id FROM quizzes
JOIN users ON quizzes.user_id = users.id
WHERE quizzes.id = $1; 
