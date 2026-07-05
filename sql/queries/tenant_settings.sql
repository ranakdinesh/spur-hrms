-- name: UpsertTenantSetting :one
INSERT INTO hrms.tenant_settings (tenant_id, key, value)
VALUES ($1, $2, $3)
ON CONFLICT (tenant_id, key) DO UPDATE SET
    value = EXCLUDED.value
RETURNING *;

-- name: GetTenantSetting :one
SELECT * FROM hrms.tenant_settings
WHERE tenant_id = $1 AND key = $2;

-- name: ListTenantSettings :many
SELECT * FROM hrms.tenant_settings
WHERE tenant_id = $1
ORDER BY key ASC;

-- name: DeleteTenantSetting :exec
DELETE FROM hrms.tenant_settings
WHERE tenant_id = $1 AND key = $2;
