-- +goose Up
CREATE INDEX idx_tenants_code ON tenants(code);
CREATE INDEX idx_tenants_email ON tenants(email);
CREATE INDEX idx_tenants_is_active ON tenants(is_active);
CREATE INDEX idx_tenants_deleted_at ON tenants(deleted_at);

-- +goose Down
DROP INDEX IF EXISTS idx_tenants_deleted_at;
DROP INDEX IF EXISTS idx_tenants_is_active;
DROP INDEX IF EXISTS idx_tenants_email;
DROP INDEX IF EXISTS idx_tenants_code;