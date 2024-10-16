-- +goose up
CREATE TABLE refresh_tokens (
    token PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    user_id NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP,
    recoked_at TIMESTAMP
);

-- +goose down
DROP TABLE refresh_tokens;