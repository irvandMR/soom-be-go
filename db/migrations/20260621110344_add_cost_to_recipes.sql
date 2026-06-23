-- +goose Up
-- +goose StatementBegin
ALTER TABLE product_recipes
ADD COLUMN total_cost NUMERIC(15, 2) NOT NULL DEFAULT 0,
ADD COLUMN cost_per_unit NUMERIC(15, 2) NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product_recipes
DROP COLUMN total_cost,
DROP COLUMN cost_per_unit;
-- +goose StatementEnd
