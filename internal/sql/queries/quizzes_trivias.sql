-- name: CreateQuizTrivia :one
INSERT INTO quizzes_trivias (
	quiz_id,
	trivia_id
)
VALUES ($1, $2)
RETURNING quizzes_trivias.id;

-- name: CreateQuizzesTrivias :copyfrom
INSERT INTO quizzes_trivias (
	quiz_id,
	trivia_id
)
VALUES ($1, $2);

-- name: GetQuizTriviaByID :one
SELECT * FROM quizzes_trivias
WHERE id = $1;

-- name: GetTriviasByQuizID :many
SELECT trivias.* FROM quizzes_trivias
JOIN trivias ON quizzes_trivias.trivia_id = trivias.id
WHERE quizzes_trivias.quiz_id = $1;

-- name: GetQuizzesByTriviaID :many
SELECT quizzes.* FROM quizzes_trivias
JOIN quizzes ON quizzes_trivias.quiz_id = quizzes.id
WHERE quizzes_trivias.trivia_id = $1;
