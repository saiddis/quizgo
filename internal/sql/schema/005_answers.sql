-- +goose Up
CREATE TABLE answers (
	id BIGSERIAL NOT NULL,
	answer VARCHAR(256) NOT NULL,
	correct BOOL NOT NULL,
	trivia_id UUID NOT NULL REFERENCES trivias(id) ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE answers;
