-- name: CreateCompOffRequest :one
INSERT INTO hrms.comp_off_requests (
    tenant_id,
    user_id,
    leave_type_id,
    fy_id,
    work_date,
    worked_minutes,
    requested_days,
    expiry_date,
    reason,
    payroll_impact,
    source_attendance_id,
    source_segment_id,
    metadata,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $14
)
RETURNING *;

-- name: GetCompOffRequest :one
SELECT * FROM hrms.comp_off_requests
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListCompOffRequestsByUser :many
SELECT * FROM hrms.comp_off_requests
WHERE tenant_id = $1
  AND user_id = $2
  AND NOT inactive
ORDER BY created_at DESC;

-- name: ListCompOffRequestsByStatus :many
SELECT * FROM hrms.comp_off_requests
WHERE tenant_id = $1
  AND status = $2
  AND NOT inactive
ORDER BY created_at ASC;

-- name: ReviewCompOffRequest :one
UPDATE hrms.comp_off_requests
SET
    status = $3,
    approved_days = $4,
    expiry_date = $5,
    payroll_impact = $6,
    reviewed_by = $7,
    reviewed_at = NOW(),
    review_remarks = $8,
    metadata = $9,
    updated_by = $7,
    updated_at = NOW()
WHERE tenant_id = $1
  AND id = $2
  AND status = 'pending'
  AND NOT inactive
RETURNING *;
