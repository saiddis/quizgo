-- name: CreateAnswer :one
INSERT INTO answers (
	quiz_id,
	trivia_id,
	option_id
)
VALUES ($1, $2, $3)
RETURNING answers.id;

-- name: GetAnswerByID :one
SELECT * FROM answers
WHERE id = $1;

-- name: GetOptionByAnswerID :one
SELECT options.* FROM answers
JOIN options ON answers.option_id = options.id
WHERE answers.id = $1;

-- name: GetAnswersByQuizID :many
SELECT answers.id as id, trivias.id as trivia_id FROM answers
JOIN trivias ON answers.trivia_id = trivias.id
WHERE answers.quiz_id = $1;
