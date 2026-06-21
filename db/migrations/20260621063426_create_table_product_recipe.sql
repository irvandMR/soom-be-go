-- +goose Up
CREATE TABLE product_recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL,
    version_number INT4 NOT NULL,
    notes TEXT,
    estimated_yield NUMERIC,
    unit_id UUID,
    is_active BOOL NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(50)
);

-- +goose Down
DROP TABLE IF EXISTS product_recipes;
