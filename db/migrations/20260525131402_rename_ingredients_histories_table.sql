-- +goose Up
ALTER TABLE tbl_ingredients_stock_histories RENAME TO ingredients_stock_histories;

-- +goose Down
ALTER TABLE ingredients_stock_histories RENAME TO tbl_ingredients_stock_histories;