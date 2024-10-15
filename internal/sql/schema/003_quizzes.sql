-- +goose Up
CREATE TABLE quizzes (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	created_at TIMESTAMP NOT NULL,
	type VARCHAR(8) NOT NULL,
	category VARCHAR(4) NOT NULL,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	score_id BIGINT REFERENCES scores(id)
);

-- +goose Down
DROP TABLE quizzes;
