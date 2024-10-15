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

-- name: GetTheHighestTotalScore :one
SELECT max(scores.total_score) as score FROM scores
GROUP BY scores.total_score
ORDER BY scores.total_score DESC
LIMIT 1;

-- name: UsersBestScorePagination :many
WITH max_scores as (
    SELECT 
        user_id, 
        max(total_score) as max_score
    FROM 
        scores
    WHERE 
        scores.total_score > 0 AND scores.total_score < $1
    GROUP BY 
        user_id
)
SELECT 
    users.email, 
    scores.easy_quizzes_done, 
    scores.medium_quizzes_done, 
    scores.hard_quizzes_done, 
    scores.completion_time, 
    scores.total_score as score
FROM 
    max_scores
JOIN 
    scores ON scores.user_id = max_scores.user_id AND scores.total_score = max_scores.max_score
RIGHT JOIN 
    users ON scores.user_id = users.id
ORDER BY 
    scores.total_score DESC
LIMIT 10;
