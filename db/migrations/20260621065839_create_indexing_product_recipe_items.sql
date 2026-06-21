-- +goose Up
CREATE INDEX idx_recipe_items_recipe_id ON product_recipe_items (recipe_id);
CREATE INDEX idx_recipe_items_ingredient_id ON product_recipe_items (ingredient_id);
CREATE INDEX idx_recipe_items_unit_id ON product_recipe_items (unit_id);

-- +goose Down
DROP INDEX IF EXISTS idx_recipe_items_unit_id;
DROP INDEX IF EXISTS idx_recipe_items_ingredient_id;
DROP INDEX IF EXISTS idx_recipe_items_recipe_id;
