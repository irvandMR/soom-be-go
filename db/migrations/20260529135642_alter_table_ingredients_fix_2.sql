-- +goose Up
ALTER TABLE ingredients
    ALTER COLUMN is_active SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN created_by SET NOT NULL;

-- +goose Down
ALTER TABLE ingredients
    ALTER COLUMN is_active DROP NOT NULL,
    ALTER COLUMN created_at DROP NOT NULL,
    ALTER COLUMN created_by DROP NOT NULL;