-- +goose Up
CREATE TABLE refresh_tokens (
token TEXT primary key,
created_at TIMESTAMP not null,
updated_at TIMESTAMP not null,
expires_at TIMESTAMP not null,
revoked_at TIMESTAMP,
user_id UUID not null,
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;