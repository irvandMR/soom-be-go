-- +goose Up
CREATE TABLE uoms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(100) NOT NULL,
    have_conversion BOOLEAN NOT NULL DEFAULT FALSE,
    base_unit VARCHAR(25),
    conversion_factor NUMERIC,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

-- +goose Down
DROP TABLE IF EXISTS uoms;