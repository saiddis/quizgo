-- +goose Up
CREATE TABLE options (
	id UUID PRIMARY KEY NOT NULL,
	option VARCHAR(256) NOT NULL,
	correct BOOL NOT NULL,
	trivia_id UUID NOT NULL REFERENCES trivias(id) ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE options;
