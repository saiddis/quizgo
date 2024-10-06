-- +goose Up
CREATE TABLE scores (
	id UUID PRIMARY KEY NOT NULL,
	completion_time BIGINT NOT NULL,
	hard_quizzes_done INT,
	medium_quizzes_done INT,
	easy_quizzes_done INT,
	total_score INT NOT NULL,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE scores;
