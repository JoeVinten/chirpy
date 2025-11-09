-- +goose Up
CREATE TABLE chirps (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NUlL,
	updated_at TIMESTAMP NOT NULL,
	body TEXT NOT NULL,
	user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
