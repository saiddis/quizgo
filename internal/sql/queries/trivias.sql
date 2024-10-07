-- name: CreateTrivia :one
INSERT INTO trivias (
	id,
	type,
	category,	
	difficulty,
	question
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTriviaByID :one
SELECT * FROM trivias
WHERE id = $1;

-- name: GetTriviaByQuestion :one
SELECT * FROM trivias
WHERE question = $1;
