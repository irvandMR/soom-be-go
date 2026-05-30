-- +goose Up
ALTER TABLE ingredients
    ALTER COLUMN purchase_price DROP NOT NULL,
    ALTER COLUMN average_price DROP NOT NULL;

-- +goose Down
ALTER TABLE ingredients
    ALTER COLUMN purchase_price SET NOT NULL,
    ALTER COLUMN average_price SET NOT NULL;