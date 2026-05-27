-- +goose Up
CREATE TABLE tbl_ingredients_stock_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ingredient_id UUID NOT NULL,
    type VARCHAR(25) NOT NULL, -- 'IN', 'OUT', 'ADJUSTMENT', 'WASTE'
    quantity NUMERIC NOT NULL,
    purchase_price NUMERIC NOT NULL DEFAULT 0,
    notes VARCHAR(255) NULL,
    reference_id VARCHAR(255) NULL,
    reference_type VARCHAR(100) NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ NULL,
    updated_by VARCHAR(50) NULL
);

-- +goose Down
DROP TABLE IF EXISTS tbl_ingredients_stock_histories;
