-- +goose Up
CREATE TABLE sessions (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	quiz_type TEXT NOT NULL,
	quiz_category TEXT NOT NULL,
	user_id UUID REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE sessions;
