-- +goose Up
INSERT INTO users (
    id,
    tenant_id,
    tenant_role,
    username,
    email,
    password,
    role,
    must_change_password,
    temp_password_expires_at,
    is_active,
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by
) VALUES (
    gen_random_uuid(),
    NULL,
    'superadmin',
    'superadmin',
    'superadmin@example.com',
    '$2a$12$a9Qq8gjPOPhNKoogRNyvxOQvUspXUtVlvmWMLcT0lseInHGPQY1f6',
    'superadmin',
    FALSE,
    NULL,
    TRUE,
    NOW(),
    'system',
    NULL,
    NULL,
    NULL,
    NULL
);

-- +goose Down
DELETE FROM users WHERE username = 'superadmin' AND role = 'superadmin';