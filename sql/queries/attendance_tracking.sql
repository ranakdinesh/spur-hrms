-- name: ListAttendancesByUser :many
SELECT * FROM hrms.attendances
WHERE tenant_id = $1 AND user_id = $2 AND date BETWEEN $3 AND $4 AND NOT inactive
ORDER BY date ASC, time ASC NULLS LAST;

-- name: ListAttendancesByUserDate :many
SELECT * FROM hrms.attendances
WHERE tenant_id = $1 AND user_id = $2 AND date = $3 AND NOT inactive
ORDER BY time ASC NULLS LAST, created_at ASC;

-- name: ListAttendancesByDate :many
SELECT * FROM hrms.attendances
WHERE tenant_id = $1 AND date = $2 AND NOT inactive
ORDER BY time ASC NULLS LAST;

-- name: GetAttendance :one
SELECT * FROM hrms.attendances
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateAttendance :one
INSERT INTO hrms.attendances (
    tenant_id,
    user_id,
    date,
    time,
    type,
    status,
    source,
    latitude,
    longitude,
    work_mode,
    remarks,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: CreateAttendanceWorkdaySegment :one
INSERT INTO hrms.attendance_workday_segments (
    tenant_id,
    user_id,
    date,
    event_time,
    segment_type,
    action,
    work_mode,
    source,
    attendance_location_id,
    reference_type,
    reference_id,
    reference_label,
    latitude,
    longitude,
    location_accuracy_meters,
    location_verification_status,
    remarks,
    metadata,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $19
)
RETURNING *;

-- name: ListAttendanceWorkdaySegmentsByUser :many
SELECT * FROM hrms.attendance_workday_segments
WHERE tenant_id = $1
  AND user_id = $2
  AND date BETWEEN $3 AND $4
  AND NOT inactive
ORDER BY date ASC, event_time ASC, created_at ASC;

-- name: ListAttendanceWorkdaySegmentsByUserDate :many
SELECT * FROM hrms.attendance_workday_segments
WHERE tenant_id = $1
  AND user_id = $2
  AND date = $3
  AND NOT inactive
ORDER BY event_time ASC, created_at ASC;

-- name: UpdateAttendance :one
UPDATE hrms.attendances
SET
    time = $3,
    type = $4,
    status = $5,
    source = $6,
    latitude = $7,
    longitude = $8,
    work_mode = $9,
    remarks = $10,
    updated_by = $11,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteAttendance :exec
UPDATE hrms.attendances
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListAttendanceRequestsByUser :many
SELECT * FROM hrms.attendance_requests
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: ListAttendanceRequestsByStatus :many
SELECT * FROM hrms.attendance_requests
WHERE tenant_id = $1 AND status = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetAttendanceRequest :one
SELECT * FROM hrms.attendance_requests
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteAttendanceRequest :exec
UPDATE hrms.attendance_requests
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListDeviceLogsByUser :many
SELECT * FROM hrms.device_logs
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY logged_at DESC;

-- name: CreateDeviceLog :one
INSERT INTO hrms.device_logs (
    tenant_id,
    user_id,
    device_id,
    device_type,
    ip_address,
    action,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $7
)
RETURNING *;

-- name: GetDeviceLog :one
SELECT * FROM hrms.device_logs
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteDeviceLog :exec
UPDATE hrms.device_logs
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: CreateAttendancePolicy :one
INSERT INTO hrms.attendance_policies (
    tenant_id, name, code, description, branch_id, department_id, user_id, schedule_type, is_default,
    standard_work_minutes, min_half_day_minutes, min_full_day_minutes, grace_late_minutes, grace_early_minutes,
    half_day_late_after_minutes, absent_late_after_minutes, half_day_early_before_minutes, absent_early_before_minutes,
    allow_flexi_hours, core_start_time, core_end_time, allow_wfh, wfh_days_per_week, allow_permanent_remote,
    require_geo, require_device, regularization_window_days, max_regularizations_per_month, approval_mode,
    effective_from, effective_to, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9,
    $10, $11, $12, $13, $14,
    $15, $16, $17, $18,
    $19, $20, $21, $22, $23, $24,
    $25, $26, $27, $28, $29,
    $30, $31, $32, $32
)
RETURNING *;

-- name: ListAttendancePolicies :many
SELECT * FROM hrms.attendance_policies
WHERE tenant_id = $1 AND NOT inactive
ORDER BY is_default DESC, name ASC;

-- name: GetAttendancePolicy :one
SELECT * FROM hrms.attendance_policies
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ResolveAttendancePolicy :one
SELECT * FROM hrms.attendance_policies
WHERE tenant_id = $1
  AND NOT inactive
  AND (effective_from IS NULL OR effective_from <= $5)
  AND (effective_to IS NULL OR effective_to >= $5)
  AND (
    user_id = $2
    OR (user_id IS NULL AND department_id = $3)
    OR (user_id IS NULL AND department_id IS NULL AND branch_id = $4)
    OR (user_id IS NULL AND department_id IS NULL AND branch_id IS NULL AND is_default)
  )
ORDER BY
  CASE
    WHEN user_id = $2 THEN 1
    WHEN department_id = $3 THEN 2
    WHEN branch_id = $4 THEN 3
    WHEN is_default THEN 4
    ELSE 5
  END,
  updated_at DESC
LIMIT 1;

-- name: UpdateAttendancePolicy :one
UPDATE hrms.attendance_policies
SET
    name = $3,
    code = $4,
    description = $5,
    branch_id = $6,
    department_id = $7,
    user_id = $8,
    schedule_type = $9,
    is_default = $10,
    standard_work_minutes = $11,
    min_half_day_minutes = $12,
    min_full_day_minutes = $13,
    grace_late_minutes = $14,
    grace_early_minutes = $15,
    half_day_late_after_minutes = $16,
    absent_late_after_minutes = $17,
    half_day_early_before_minutes = $18,
    absent_early_before_minutes = $19,
    allow_flexi_hours = $20,
    core_start_time = $21,
    core_end_time = $22,
    allow_wfh = $23,
    wfh_days_per_week = $24,
    allow_permanent_remote = $25,
    require_geo = $26,
    require_device = $27,
    regularization_window_days = $28,
    max_regularizations_per_month = $29,
    approval_mode = $30,
    effective_from = $31,
    effective_to = $32,
    updated_by = $33,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteAttendancePolicy :exec
UPDATE hrms.attendance_policies
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateAttendanceRoster :one
INSERT INTO hrms.attendance_rosters (
    tenant_id, user_id, policy_id, date, start_time, end_time, break_minutes, work_mode, location_type, remarks, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11
)
RETURNING *;

-- name: ListAttendanceRostersByDateRange :many
SELECT * FROM hrms.attendance_rosters
WHERE tenant_id = $1 AND date BETWEEN $2 AND $3 AND NOT inactive
ORDER BY date ASC, user_id ASC;

-- name: ListAttendanceRostersByUser :many
SELECT * FROM hrms.attendance_rosters
WHERE tenant_id = $1 AND user_id = $2 AND date BETWEEN $3 AND $4 AND NOT inactive
ORDER BY date ASC;

-- name: GetAttendanceRosterByUserDate :one
SELECT * FROM hrms.attendance_rosters
WHERE tenant_id = $1 AND user_id = $2 AND date = $3 AND NOT inactive;

-- name: UpdateAttendanceRoster :one
UPDATE hrms.attendance_rosters
SET policy_id = $3,
    date = $4,
    start_time = $5,
    end_time = $6,
    break_minutes = $7,
    work_mode = $8,
    location_type = $9,
    remarks = $10,
    updated_by = $11,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteAttendanceRoster :exec
UPDATE hrms.attendance_rosters
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateAttendanceRequest :one
INSERT INTO hrms.attendance_requests (
    tenant_id, user_id, date, requested_type, request_type, requested_checkin_at, requested_checkout_at,
    requested_work_mode, policy_id, roster_id, reason, status, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'pending', $12, $12
)
RETURNING *;

-- name: UpdateAttendanceRequestReview :one
UPDATE hrms.attendance_requests
SET status = $3,
    reviewed_by = $4,
    reviewed_at = NOW(),
    remarks = $5,
    updated_by = $4,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateAttendanceLocation :one
INSERT INTO hrms.attendance_locations (
    tenant_id, branch_id, code, name, location_type, latitude, longitude, radius_meters,
    address, city, state, country, pincode, effective_from, effective_to, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14, $15, $16, $16
)
RETURNING *;

-- name: ListAttendanceLocations :many
SELECT * FROM hrms.attendance_locations
WHERE tenant_id = $1 AND NOT inactive
ORDER BY name ASC;

-- name: GetAttendanceLocation :one
SELECT * FROM hrms.attendance_locations
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateAttendanceLocation :one
UPDATE hrms.attendance_locations
SET branch_id = $3,
    code = $4,
    name = $5,
    location_type = $6,
    latitude = $7,
    longitude = $8,
    radius_meters = $9,
    address = $10,
    city = $11,
    state = $12,
    country = $13,
    pincode = $14,
    effective_from = $15,
    effective_to = $16,
    updated_by = $17,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteAttendanceLocation :exec
UPDATE hrms.attendance_locations
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateAttendanceLocationAssignment :one
INSERT INTO hrms.attendance_location_assignments (
    tenant_id, location_id, branch_id, department_id, user_id, effective_from, effective_to, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $8
)
RETURNING *;

-- name: ListAttendanceLocationAssignments :many
SELECT * FROM hrms.attendance_location_assignments
WHERE tenant_id = $1 AND NOT inactive
ORDER BY created_at DESC;

-- name: ListAttendanceLocationAssignmentsByLocation :many
SELECT * FROM hrms.attendance_location_assignments
WHERE tenant_id = $1 AND location_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: UpdateAttendanceLocationAssignment :one
UPDATE hrms.attendance_location_assignments
SET location_id = $3,
    branch_id = $4,
    department_id = $5,
    user_id = $6,
    effective_from = $7,
    effective_to = $8,
    updated_by = $9,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteAttendanceLocationAssignment :exec
UPDATE hrms.attendance_location_assignments
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateAttendanceDevice :one
INSERT INTO hrms.attendance_devices (
    tenant_id, attendance_location_id, branch_id, code, name, vendor, model, serial_number,
    device_identifier, integration_type, direction_mode, timezone, status, config, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14, $15, $15
)
RETURNING *;

-- name: ListAttendanceDevices :many
SELECT * FROM hrms.attendance_devices
WHERE tenant_id = $1 AND NOT inactive
ORDER BY name ASC;

-- name: GetAttendanceDevice :one
SELECT * FROM hrms.attendance_devices
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateAttendanceDevice :one
UPDATE hrms.attendance_devices
SET attendance_location_id = $3,
    branch_id = $4,
    code = $5,
    name = $6,
    vendor = $7,
    model = $8,
    serial_number = $9,
    device_identifier = $10,
    integration_type = $11,
    direction_mode = $12,
    timezone = $13,
    status = $14,
    config = $15,
    updated_by = $16,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteAttendanceDevice :exec
UPDATE hrms.attendance_devices
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateEmployeeAttendanceDevice :one
INSERT INTO hrms.employee_attendance_devices (
    tenant_id, user_id, device_id, device_user_id, credential_type, card_number, effective_from, effective_to, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $9
)
RETURNING *;

-- name: ListEmployeeAttendanceDevices :many
SELECT * FROM hrms.employee_attendance_devices
WHERE tenant_id = $1 AND NOT inactive
ORDER BY created_at DESC;

-- name: ListEmployeeAttendanceDevicesByUser :many
SELECT * FROM hrms.employee_attendance_devices
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetEmployeeAttendanceDevice :one
SELECT * FROM hrms.employee_attendance_devices
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateEmployeeAttendanceDevice :one
UPDATE hrms.employee_attendance_devices
SET user_id = $3,
    device_id = $4,
    device_user_id = $5,
    credential_type = $6,
    card_number = $7,
    effective_from = $8,
    effective_to = $9,
    updated_by = $10,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteEmployeeAttendanceDevice :exec
UPDATE hrms.employee_attendance_devices
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: CreateRawAttendanceEvent :one
INSERT INTO hrms.raw_attendance_events (
    tenant_id, device_id, employee_device_mapping_id, external_event_id, device_user_id, event_time,
    event_type, attendance_type, import_batch_id, raw_payload, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $11
)
RETURNING *;

-- name: ListRawAttendanceEvents :many
SELECT * FROM hrms.raw_attendance_events
WHERE tenant_id = $1 AND NOT inactive
ORDER BY event_time DESC
LIMIT $2;

-- name: UpdateRawAttendanceEventProcessed :one
UPDATE hrms.raw_attendance_events
SET attendance_id = $3,
    processing_status = $4,
    processing_error = $5,
    processed_at = NOW(),
    updated_by = $6,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;
