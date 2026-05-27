-- +goose Up
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
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
