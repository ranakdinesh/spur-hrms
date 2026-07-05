-- name: CreateDesignationLevelCode :one
INSERT INTO hrms.designation_level_codes (
    tenant_id,
    code,
    label,
    description,
    sort_order,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $6
)
RETURNING *;

-- name: ListDesignationLevelCodes :many
SELECT * FROM hrms.designation_level_codes
WHERE tenant_id = $1 AND NOT inactive
ORDER BY sort_order ASC, code ASC;

-- name: UpdateDesignationLevelCode :one
UPDATE hrms.designation_level_codes
SET code = $3,
    label = $4,
    description = $5,
    sort_order = $6,
    updated_by = $7
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteDesignationLevelCode :exec
UPDATE hrms.designation_level_codes
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateDesignationSeniorityRank :one
INSERT INTO hrms.designation_seniority_ranks (
    tenant_id,
    rank_value,
    label,
    description,
    sort_order,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $6
)
RETURNING *;

-- name: ListDesignationSeniorityRanks :many
SELECT * FROM hrms.designation_seniority_ranks
WHERE tenant_id = $1 AND NOT inactive
ORDER BY sort_order ASC, rank_value ASC;

-- name: UpdateDesignationSeniorityRank :one
UPDATE hrms.designation_seniority_ranks
SET rank_value = $3,
    label = $4,
    description = $5,
    sort_order = $6,
    updated_by = $7
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteDesignationSeniorityRank :exec
UPDATE hrms.designation_seniority_ranks
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
