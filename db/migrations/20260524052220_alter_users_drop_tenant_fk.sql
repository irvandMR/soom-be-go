-- +goose Up
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_tenant_id_fkey;

-- +goose Down
ALTER TABLE users ADD CONSTRAINT users_tenant_id_fkey FOREIGN KEY (tenant_id) REFERENCES tenants(id);