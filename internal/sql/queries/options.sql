-- name: CreateOption :one
INSERT INTO options (
	id,
	option,
	correct,
	trivia_id
)
VALUES ($1, $2, $3, $4)
RETURNING options.id;

-- name: CreateOptions :copyfrom
INSERT INTO options (
	id,
	option,
	correct,
	trivia_id
)
VALUES ($1, $2, $3, $4);

-- name: GetOptionByID :one
SELECT * FROM options
WHERE id = $1;

-- name: GetOptionsIDByTriviaID :many
SELECT options.id as id, options.correct as correct FROM options
JOIN trivias ON trivia_id = trivias.id
WHERE trivia_id = $1;
