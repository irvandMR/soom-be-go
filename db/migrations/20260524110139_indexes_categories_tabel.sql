-- +goose Up
CREATE INDEX idx_categories_code ON categories(code);
CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_categories_type ON categories(type);

-- +goose Down
DROP INDEX IF EXISTS idx_categories_type;
DROP INDEX IF EXISTS idx_categories_name;
DROP INDEX IF EXISTS idx_categories_code;

