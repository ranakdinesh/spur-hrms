-- name: CreateBranch :one
INSERT INTO hrms.branches (
    tenant_id,
    branch_name,
    address,
    city,
    state,
    country,
    pincode,
    phone,
    branch_manager_user_id,
    hr_user_id,
    accounts_user_id,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: ListBranches :many
SELECT * FROM hrms.branches
WHERE tenant_id = $1 AND NOT inactive
ORDER BY branch_name ASC;

-- name: GetBranch :one
SELECT * FROM hrms.branches
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateBranch :one
UPDATE hrms.branches
SET branch_name = $3,
    address = $4,
    city = $5,
    state = $6,
    country = $7,
    pincode = $8,
    phone = $9,
    branch_manager_user_id = $10,
    hr_user_id = $11,
    accounts_user_id = $12,
    updated_by = $13
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteBranch :exec
UPDATE hrms.branches
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
