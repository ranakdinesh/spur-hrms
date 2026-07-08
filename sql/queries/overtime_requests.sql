-- name: CreateOvertimeRequest :one
INSERT INTO hrms.overtime_requests (
    tenant_id,
    user_id,
    work_date,
    requested_minutes,
    reason,
    calculation_type,
    rate_multiplier,
    payroll_component_code,
    source_attendance_id,
    source_segment_id,
    metadata,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: GetOvertimeRequest :one
SELECT * FROM hrms.overtime_requests
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListOvertimeRequestsByUser :many
SELECT * FROM hrms.overtime_requests
WHERE tenant_id = $1
  AND user_id = $2
  AND NOT inactive
ORDER BY created_at DESC;

-- name: ListOvertimeRequestsByStatus :many
SELECT * FROM hrms.overtime_requests
WHERE tenant_id = $1
  AND status = $2
  AND NOT inactive
ORDER BY created_at ASC;

-- name: ListOvertimeRequestsByPayrollExportStatus :many
SELECT * FROM hrms.overtime_requests
WHERE tenant_id = $1
  AND payroll_export_status = $2
  AND NOT inactive
ORDER BY work_date ASC, created_at ASC;

-- name: ReviewOvertimeRequest :one
UPDATE hrms.overtime_requests
SET
    status = $3,
    approved_minutes = $4,
    review_remarks = $5,
    calculation_type = $6,
    rate_multiplier = $7,
    payroll_component_code = $8,
    payroll_export_status = $9,
    reviewed_by = $10,
    reviewed_at = NOW(),
    metadata = $11,
    updated_by = $10,
    updated_at = NOW()
WHERE tenant_id = $1
  AND id = $2
  AND status = 'pending'
  AND NOT inactive
RETURNING *;
