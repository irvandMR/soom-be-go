-- +goose Up
ALTER TABLE ingredients
    ALTER COLUMN stock_qty DROP NOT NULL;

-- +goose Down
ALTER TABLE ingredients
    ALTER COLUMN stock_qty SET NOT NULL;