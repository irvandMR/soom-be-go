-- +goose Up
CREATE TABLE ingredients(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    category_id UUID NOT NULL,
    unit_id UUID NOT NULL,
    name VARCHAR(225) NOT NULL,
    stock_qty NUMERIC NOT NULL,
    min_stock_qty NUMERIC NOT NULL,
    price_per_unit NUMERIC NULL,
    purchase_price NUMERIC NOT NULL,
    average_price NUMERIC NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

-- +goose Down
DROP TABLE IF EXISTS categories;
