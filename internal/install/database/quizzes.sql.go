// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: quizzes.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createQuiz = `-- name: CreateQuiz :one
INSERT INTO quizzes (
	created_at,
	type,
	category,
	user_id
)
VALUES ($1, $2, $3, $4)
RETURNING quizzes.id
`

type CreateQuizParams struct {
	CreatedAt pgtype.Timestamp
	Type      string
	Category  string
	UserID    uuid.UUID
}

func (q *Queries) CreateQuiz(ctx context.Context, arg CreateQuizParams) (int64, error) {
	row := q.db.QueryRow(ctx, createQuiz,
		arg.CreatedAt,
		arg.Type,
		arg.Category,
		arg.UserID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getLastQuizByUserID = `-- name: GetLastQuizByUserID :one
SELECT quizzes.id, quizzes.created_at, type, category, user_id, score_id, users.id, users.created_at, email FROM quizzes
JOIN users ON quizzes.user_id = users.id
WHERE quizzes.user_id = $1
ORDER BY quizzes.id DESC
LIMIT 1
`

type GetLastQuizByUserIDRow struct {
	ID          int64
	CreatedAt   pgtype.Timestamp
	Type        string
	Category    string
	UserID      uuid.UUID
	ScoreID     pgtype.Int8
	ID_2        uuid.UUID
	CreatedAt_2 pgtype.Timestamp
	Email       string
}

func (q *Queries) GetLastQuizByUserID(ctx context.Context, userID uuid.UUID) (GetLastQuizByUserIDRow, error) {
	row := q.db.QueryRow(ctx, getLastQuizByUserID, userID)
	var i GetLastQuizByUserIDRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Type,
		&i.Category,
		&i.UserID,
		&i.ScoreID,
		&i.ID_2,
		&i.CreatedAt_2,
		&i.Email,
	)
	return i, err
}

const getQuizzesByUserID = `-- name: GetQuizzesByUserID :many
SELECT id, created_at, type, category, user_id, score_id FROM quizzes
WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetQuizzesByUserID(ctx context.Context, userID uuid.UUID) ([]Quiz, error) {
	rows, err := q.db.Query(ctx, getQuizzesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Quiz
	for rows.Next() {
		var i Quiz
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Type,
			&i.Category,
			&i.UserID,
			&i.ScoreID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserIDByQuizID = `-- name: GetUserIDByQuizID :one
SELECT users.id as user_id FROM quizzes
JOIN users ON quizzes.user_id = users.id
WHERE quizzes.id = $1
`

func (q *Queries) GetUserIDByQuizID(ctx context.Context, id int64) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, getUserIDByQuizID, id)
	var user_id uuid.UUID
	err := row.Scan(&user_id)
	return user_id, err
}

const quizzesPagination = `-- name: QuizzesPagination :many
SELECT quizzes.created_at, quizzes.type, quizzes.category, quizzes.score_id as score_id, quizzes.id as id FROM quizzes
JOIN users ON quizzes.user_id = users.id
WHERE quizzes.user_id = $1 AND quizzes.id < $2
ORDER BY quizzes.id DESC
LIMIT 5
`

type QuizzesPaginationParams struct {
	UserID uuid.UUID
	ID     int64
}

type QuizzesPaginationRow struct {
	CreatedAt pgtype.Timestamp
	Type      string
	Category  string
	ScoreID   pgtype.Int8
	ID        int64
}

func (q *Queries) QuizzesPagination(ctx context.Context, arg QuizzesPaginationParams) ([]QuizzesPaginationRow, error) {
	rows, err := q.db.Query(ctx, quizzesPagination, arg.UserID, arg.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QuizzesPaginationRow
	for rows.Next() {
		var i QuizzesPaginationRow
		if err := rows.Scan(
			&i.CreatedAt,
			&i.Type,
			&i.Category,
			&i.ScoreID,
			&i.ID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateScoreID = `-- name: UpdateScoreID :one
UPDATE quizzes
SET score_id = $1
WHERE id = $2
RETURNING id, created_at, type, category, user_id, score_id
`

type UpdateScoreIDParams struct {
	ScoreID pgtype.Int8
	ID      int64
}

func (q *Queries) UpdateScoreID(ctx context.Context, arg UpdateScoreIDParams) (Quiz, error) {
	row := q.db.QueryRow(ctx, updateScoreID, arg.ScoreID, arg.ID)
	var i Quiz
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Type,
		&i.Category,
		&i.UserID,
		&i.ScoreID,
	)
	return i, err
}
