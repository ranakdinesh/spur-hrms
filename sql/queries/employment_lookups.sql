-- name: CreateEmploymentType :one
INSERT INTO hrms.employment_types (
    tenant_id,
    name,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $3
)
RETURNING *;

-- name: ListEmploymentTypes :many
SELECT * FROM hrms.employment_types
WHERE tenant_id = $1 AND NOT inactive
ORDER BY name ASC;

-- name: GetEmploymentType :one
SELECT * FROM hrms.employment_types
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateEmploymentType :one
UPDATE hrms.employment_types
SET name = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteEmploymentType :exec
UPDATE hrms.employment_types
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateMaritalStatus :one
INSERT INTO hrms.marital_statuses (
    tenant_id,
    name,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $3
)
RETURNING *;

-- name: ListMaritalStatuses :many
SELECT * FROM hrms.marital_statuses
WHERE tenant_id = $1 AND NOT inactive
ORDER BY name ASC;

-- name: GetMaritalStatus :one
SELECT * FROM hrms.marital_statuses
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateMaritalStatus :one
UPDATE hrms.marital_statuses
SET name = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteMaritalStatus :exec
UPDATE hrms.marital_statuses
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
