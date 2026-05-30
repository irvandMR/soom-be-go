-- +goose Up
ALTER TABLE ingredients
    ALTER COLUMN purchase_price DROP NOT NULL,
    ALTER COLUMN average_price DROP NOT NULL,
    ALTER COLUMN is_active DROP NOT NULL,
    ALTER COLUMN created_at DROP NOT NULL,
    ALTER COLUMN created_by DROP NOT NULL;

-- +goose Down
ALTER TABLE ingredients
    ALTER COLUMN purchase_price SET NOT NULL,
    ALTER COLUMN average_price SET NOT NULL,
    ALTER COLUMN is_active SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN created_by SET NOT NULL;