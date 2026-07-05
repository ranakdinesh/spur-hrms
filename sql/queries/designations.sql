-- name: CreateDesignation :one
INSERT INTO hrms.designations (
    tenant_id,
    name,
    level_code,
    seniority_rank,
    description,
    attendance_required,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $7
)
RETURNING *;

-- name: ListDesignations :many
SELECT * FROM hrms.designations
WHERE tenant_id = $1 AND NOT inactive
ORDER BY seniority_rank DESC, name ASC;

-- name: GetDesignation :one
SELECT * FROM hrms.designations
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateDesignation :one
UPDATE hrms.designations
SET name = $3,
    level_code = $4,
    seniority_rank = $5,
    description = $6,
    attendance_required = $7,
    updated_by = $8
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteDesignation :exec
UPDATE hrms.designations
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
