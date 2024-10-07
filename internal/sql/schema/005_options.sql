-- +goose Up
CREATE TABLE options (
	id BIGSERIAL NOT NULL,
	option VARCHAR(256) NOT NULL,
	correct BOOL NOT NULL,
	trivia_id UUID NOT NULL REFERENCES trivias(id) ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE answers;
