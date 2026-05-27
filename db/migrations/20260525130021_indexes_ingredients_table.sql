-- +goose Up
CREATE INDEX idx_ingredients_tenant_id ON ingredients(tenant_id);
CREATE INDEX idx_ingredients_category_id ON ingredients(category_id);
CREATE INDEX idx_ingredients_unit_id ON ingredients(unit_id);
CREATE INDEX idx_ingredients_name ON ingredients(name);
CREATE INDEX idx_ingredients_is_active ON ingredients(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_ingredients_deleted_at ON ingredients(deleted_at) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_ingredients_deleted_at;
DROP INDEX IF EXISTS idx_ingredients_is_active;
DROP INDEX IF EXISTS idx_ingredients_name;
DROP INDEX IF EXISTS idx_ingredients_unit_id;
DROP INDEX IF EXISTS idx_ingredients_category_id;
DROP INDEX IF EXISTS idx_ingredients_tenant_id;
