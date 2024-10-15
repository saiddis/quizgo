-- +goose Up
CREATE TABLE answers (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	quiz_id BIGINT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE, 
	trivia_id UUID NOT NULL REFERENCES trivias(id) ON DELETE CASCADE,
	option_id UUID NOT NULL REFERENCES options(id)
);

-- +goose Down
DROP TABLE answers;
