-- +goose Up
CREATE TABLE product_recipe_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID NOT NULL,
    ingredient_id UUID NOT NULL,
    quantity NUMERIC,
    unit_id UUID,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(50),
    ingredient_cost NUMERIC NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS product_recipe_items;
