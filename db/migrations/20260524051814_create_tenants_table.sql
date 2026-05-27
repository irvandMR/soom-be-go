-- +goose Up
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    business_name VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    phone VARCHAR(25) NOT NULL,
    email TEXT NOT NULL UNIQUE,
    invoice_footer TEXT,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

-- +goose Down
DROP TABLE IF EXISTS tenants;