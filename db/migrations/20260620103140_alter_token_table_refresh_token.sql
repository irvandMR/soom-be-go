-- +goose Up
ALTER TABLE refresh_tokens
    ALTER COLUMN token TYPE text;

-- +goose Down
ALTER TABLE refresh_tokens
    ALTER COLUMN token TYPE VARCHAR(255);
