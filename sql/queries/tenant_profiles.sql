-- name: UpsertTenantProfile :one
INSERT INTO hrms.tenant_profiles (
    tenant_id, subdomain, mobile_activation_code, display_name, logo_object_key
)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (tenant_id) DO UPDATE SET
    subdomain = EXCLUDED.subdomain,
    mobile_activation_code = EXCLUDED.mobile_activation_code,
    display_name = EXCLUDED.display_name,
    logo_object_key = EXCLUDED.logo_object_key
RETURNING *;

-- name: GetTenantProfile :one
SELECT * FROM hrms.tenant_profiles
WHERE tenant_id = $1;

-- name: GetTenantProfileBySubdomain :one
SELECT * FROM hrms.tenant_profiles
WHERE subdomain = $1;

-- name: GetTenantProfileByActivationCode :one
SELECT * FROM hrms.tenant_profiles
WHERE mobile_activation_code = $1;
