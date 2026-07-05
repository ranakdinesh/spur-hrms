-- name: CreateDepartment :one
INSERT INTO hrms.departments (
    tenant_id,
    name,
    short_code,
    description,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $5
)
RETURNING *;

-- name: ListDepartments :many
SELECT * FROM hrms.departments
WHERE tenant_id = $1 AND NOT inactive
ORDER BY name ASC;

-- name: GetDepartment :one
SELECT * FROM hrms.departments
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateDepartment :one
UPDATE hrms.departments
SET name = $3,
    short_code = $4,
    description = $5,
    updated_by = $6
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteDepartment :exec
UPDATE hrms.departments
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
