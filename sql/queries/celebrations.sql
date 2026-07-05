-- name: ListCelebrationTypes :many
SELECT * FROM hrms.celebration_types
WHERE tenant_id = $1 AND NOT inactive
ORDER BY name ASC;

-- name: CreateCelebrationType :one
INSERT INTO hrms.celebration_types (tenant_id, name, is_yearly, is_user_celebration, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCelebrationType :one
SELECT * FROM hrms.celebration_types
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateCelebrationType :one
UPDATE hrms.celebration_types
SET name = $3,
    is_yearly = $4,
    is_user_celebration = $5,
    updated_at = NOW(),
    updated_by = $6
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteCelebrationType :exec
UPDATE hrms.celebration_types
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListCelebrations :many
SELECT * FROM hrms.celebrations
WHERE tenant_id = $1 AND NOT inactive
ORDER BY celebration_date ASC NULLS LAST, created_at DESC;

-- name: CreateCelebration :one
INSERT INTO hrms.celebrations (tenant_id, branch_id, user_id, celebration_type_id, celebration_date, custom_title, description, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListCelebrationsByUser :many
SELECT * FROM hrms.celebrations
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY celebration_date ASC NULLS LAST;

-- name: GetCelebration :one
SELECT * FROM hrms.celebrations
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateCelebration :one
UPDATE hrms.celebrations
SET branch_id = $3,
    user_id = $4,
    celebration_type_id = $5,
    celebration_date = $6,
    custom_title = $7,
    description = $8,
    updated_at = NOW(),
    updated_by = $9
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteCelebration :exec
UPDATE hrms.celebrations
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;
