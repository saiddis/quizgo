-- +goose Up
CREATE TABLE trivias (
	id UUID PRIMARY KEY NOT NULL,
	type VARCHAR(8) NOT NULL,
	category VARCHAR(64) NOT NULL,
	difficulty VARCHAR(8) NOT NULL,
	question TEXT NOT NULL
);

-- +goose Down
DROP TABLE trivias;
