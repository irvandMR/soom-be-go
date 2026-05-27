-- +goose Up
CREATE INDEX idx_stock_histories_ingredient_id ON tbl_ingredients_stock_histories(ingredient_id);
CREATE INDEX idx_stock_histories_ref ON tbl_ingredients_stock_histories(reference_type, reference_id);
CREATE INDEX idx_stock_histories_created_at ON tbl_ingredients_stock_histories(created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_stock_histories_created_at;
DROP INDEX IF EXISTS idx_stock_histories_ref;
DROP INDEX IF EXISTS idx_stock_histories_ingredient_id;
