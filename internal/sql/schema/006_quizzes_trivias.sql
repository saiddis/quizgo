-- +goose Up
CREATE TABLE quizzes_trivias (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	quiz_id BIGINT NOT NULL REFERENCES quizzes(id),
	trivia_id UUID NOT NULL REFERENCES trivias(id),
	UNIQUE(quiz_id, trivia_id)
);

-- +goose Down
DROP TABLE quizzes_trivias;
