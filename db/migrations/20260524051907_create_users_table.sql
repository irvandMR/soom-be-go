-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    tenant_role VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role VARCHAR(50) NOT NULL,
    must_change_password BOOLEAN DEFAULT TRUE NOT NULL,
    temp_password_expires_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(50)
);

-- +goose Down
DROP TABLE IF EXISTS users;