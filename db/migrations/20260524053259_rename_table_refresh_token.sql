-- +goose Up
ALTER TABLE tbl_refresh_tokens RENAME TO refresh_tokens;

-- +goose Down
ALTER TABLE refresh_tokens RENAME TO tbl_refresh_tokens;