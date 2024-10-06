-- +goose Up
CREATE TABLE quizzes (
	id UUID PRIMARY KEY NOT NULL,
	created_at TIMESTAMP NOT NULL,
	quiz_type VARCHAR(8) NOT NULL,
	quiz_category VARCHAR(4) NOT NULL,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE quizzes;
