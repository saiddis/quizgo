-- +goose Up
CREATE TABLE quizzes_trivias (
	id UUID PRIMARY KEY NOT NULL,
	quiz_id UUID NOT NULL REFERENCES quizzes(id),
	trivia_id UUID NOT NULL REFERENCES trivias(id)
);

-- +goose Down
DROP TABLE quizzes_trivias;
