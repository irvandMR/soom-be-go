-- +goose Up
-- +goose StatementBegin
WITH new_tenant AS (
    INSERT INTO tenants (
        id,
        code,
        business_name,
        address,
        phone,
        email,
        invoice_footer,
        is_active,
        created_at,
        created_by
    ) VALUES (
        gen_random_uuid(),
        'TEN-STD',                  -- ganti dengan code tenant baru
        'Tentic Studio',  -- ganti dengan nama bisnis
        'Tangerang Selatan',       -- ganti dengan alamat
        '082386092523',           -- ganti dengan nomor telepon
        'irvandirizky41@gmail.com',         -- ganti dengan email
        NULL,
        TRUE,
        NOW(),
        'system'
    )
    RETURNING id
)
UPDATE users
SET tenant_id = (SELECT id FROM new_tenant),
    updated_at = NOW(),
    updated_by = 'system'
WHERE username = 'superadmin';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Kembalikan tenant_id user ke NULL, lalu hapus tenant yang baru dibuat
UPDATE users
SET tenant_id = NULL,
    updated_at = NOW(),
    updated_by = 'system'
WHERE username = 'superadmin';
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM tenants
WHERE code = 'TEN-STD';
-- +goose StatementEnd