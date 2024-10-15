-- +goose Up
CREATE TABLE scores (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	completion_time BIGINT NOT NULL,
	hard_quizzes_done INT NOT NULL,
	medium_quizzes_done INT NOT NULL,
	easy_quizzes_done INT NOT NULL,
	total_score INT NOT NULL NOT NULL,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE scores;
