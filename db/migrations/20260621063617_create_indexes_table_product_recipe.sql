-- +goose Up
CREATE INDEX idx_product_recipes_product_id ON product_recipes (product_id);
CREATE INDEX idx_product_recipes_unit_id ON product_recipes (unit_id);
CREATE INDEX idx_product_recipes_product_version ON product_recipes (product_id, version_number);

-- +goose Down
DROP INDEX IF EXISTS idx_product_recipes_product_version;
DROP INDEX IF EXISTS idx_product_recipes_unit_id;
DROP INDEX IF EXISTS idx_product_recipes_product_id;