-- name: CreateOption :one
INSERT INTO options (
	id,
	option,
	correct,
	trivia_id
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOptionByID :one
SELECT * FROM options
WHERE id = $1;

-- name: GetOptionsByTriviaID :many
SELECT * FROM options
JOIN trivias ON trivia_id = trivias.id
WHERE trivia_id = $1;
