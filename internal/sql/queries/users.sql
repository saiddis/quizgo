-- name: CreateUser :one
INSERT INTO users (
	id,
	created_at,
	email
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
