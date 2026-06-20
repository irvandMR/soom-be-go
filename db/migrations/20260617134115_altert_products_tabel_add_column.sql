-- +goose Up
ALTER TABLE products RENAME COLUMN stock_quantity TO stock_qty;
ALTER TABLE products ADD COLUMN is_active BOOLEAN DEFAULT TRUE NOT NULL;


-- +goose Down
ALTER TABLE products DROP COLUMN is_active;
ALTER TABLE products RENAME COLUMN stock_qty TO stock_quantity;
