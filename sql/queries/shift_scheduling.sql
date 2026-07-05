-- name: CreateShiftTemplate :one
INSERT INTO hrms.shift_templates (
  tenant_id, code, name, description, start_time, end_time, break_minutes, paid_minutes, work_mode, location_type,
  attendance_policy_id, attendance_location_id, allow_overtime, payroll_code, metadata, is_active, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17
) RETURNING *;

-- name: UpdateShiftTemplate :one
UPDATE hrms.shift_templates
SET code=$3, name=$4, description=$5, start_time=$6, end_time=$7, break_minutes=$8, paid_minutes=$9,
    work_mode=$10, location_type=$11, attendance_policy_id=$12, attendance_location_id=$13,
    allow_overtime=$14, payroll_code=$15, metadata=$16, is_active=$17, updated_by=$18
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListShiftTemplates :many
SELECT * FROM hrms.shift_templates
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('active_only')::boolean IS NULL OR is_active = sqlc.narg('active_only')::boolean)
  AND (sqlc.narg('search')::text IS NULL OR code ILIKE '%' || sqlc.narg('search')::text || '%' OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY name ASC, code ASC
LIMIT $2 OFFSET $3;

-- name: GetShiftTemplate :one
SELECT * FROM hrms.shift_templates
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: SoftDeleteShiftTemplate :exec
UPDATE hrms.shift_templates
SET inactive=TRUE, is_active=FALSE, updated_by=$3
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateStaffingRequirement :one
INSERT INTO hrms.staffing_requirements (
  tenant_id, name, requirement_date, start_date, end_date, day_of_week, branch_id, department_id,
  attendance_location_id, role_label, team_label, shift_template_id, required_count, min_count, max_count,
  priority, status, payroll_blocking, notes, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21
) RETURNING *;

-- name: UpdateStaffingRequirement :one
UPDATE hrms.staffing_requirements
SET name=$3, requirement_date=$4, start_date=$5, end_date=$6, day_of_week=$7, branch_id=$8,
    department_id=$9, attendance_location_id=$10, role_label=$11, team_label=$12, shift_template_id=$13,
    required_count=$14, min_count=$15, max_count=$16, priority=$17, status=$18, payroll_blocking=$19,
    notes=$20, metadata=$21, updated_by=$22
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListStaffingRequirements :many
SELECT * FROM hrms.staffing_requirements
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (
    sqlc.narg('start_date')::date IS NULL
    OR requirement_date IS NULL
    OR requirement_date BETWEEN sqlc.narg('start_date')::date AND COALESCE(sqlc.narg('end_date')::date, sqlc.narg('start_date')::date)
    OR start_date IS NULL
    OR COALESCE(end_date, start_date) >= sqlc.narg('start_date')::date
  )
ORDER BY COALESCE(requirement_date, start_date, created_at::date) DESC, name ASC
LIMIT $2 OFFSET $3;

-- name: GetStaffingRequirement :one
SELECT * FROM hrms.staffing_requirements
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: SoftDeleteStaffingRequirement :exec
UPDATE hrms.staffing_requirements
SET inactive=TRUE, status='archived', updated_by=$3
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateShiftScheduleAssignment :one
INSERT INTO hrms.shift_schedule_assignments (
  tenant_id, schedule_date, worker_profile_id, employee_user_id, shift_template_id, attendance_policy_id,
  attendance_location_id, attendance_roster_id, branch_id, department_id, start_time, end_time, break_minutes,
  work_mode, location_type, role_label, team_label, status, source, overtime_planned_minutes, has_conflict,
  conflict_reason, payroll_blocking, notes, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26
) RETURNING *;

-- name: UpdateShiftScheduleAssignment :one
UPDATE hrms.shift_schedule_assignments
SET schedule_date=$3, worker_profile_id=$4, employee_user_id=$5, shift_template_id=$6, attendance_policy_id=$7,
    attendance_location_id=$8, attendance_roster_id=$9, branch_id=$10, department_id=$11, start_time=$12,
    end_time=$13, break_minutes=$14, work_mode=$15, location_type=$16, role_label=$17, team_label=$18,
    status=$19, source=$20, overtime_planned_minutes=$21, has_conflict=$22, conflict_reason=$23,
    payroll_blocking=$24, notes=$25, metadata=$26, updated_by=$27
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: UpdateShiftScheduleAssignmentStatus :one
UPDATE hrms.shift_schedule_assignments
SET status=$3, has_conflict=$4, conflict_reason=$5, payroll_blocking=$6, updated_by=$7
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListShiftScheduleAssignments :many
SELECT * FROM hrms.shift_schedule_assignments
WHERE tenant_id=$1
  AND NOT inactive
  AND schedule_date BETWEEN $2 AND $3
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('employee_user_id')::uuid IS NULL OR employee_user_id = sqlc.narg('employee_user_id')::uuid)
  AND (sqlc.narg('branch_id')::uuid IS NULL OR branch_id = sqlc.narg('branch_id')::uuid)
  AND (sqlc.narg('department_id')::uuid IS NULL OR department_id = sqlc.narg('department_id')::uuid)
  AND (sqlc.narg('attendance_location_id')::uuid IS NULL OR attendance_location_id = sqlc.narg('attendance_location_id')::uuid)
ORDER BY schedule_date ASC, start_time ASC, created_at ASC
LIMIT $4 OFFSET $5;

-- name: GetShiftScheduleAssignment :one
SELECT * FROM hrms.shift_schedule_assignments
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: ListShiftAssignmentsForWorkerDate :many
SELECT * FROM hrms.shift_schedule_assignments
WHERE tenant_id=$1
  AND NOT inactive
  AND schedule_date=$2
  AND status <> 'cancelled'
  AND (worker_profile_id = sqlc.narg('worker_profile_id')::uuid OR employee_user_id = sqlc.narg('employee_user_id')::uuid)
ORDER BY start_time ASC;

-- name: SoftDeleteShiftScheduleAssignment :exec
UPDATE hrms.shift_schedule_assignments
SET inactive=TRUE, updated_by=$3
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateShiftSwapRequest :one
INSERT INTO hrms.shift_swap_requests (
  tenant_id, requester_assignment_id, requester_worker_profile_id, requester_user_id, target_worker_profile_id,
  target_user_id, offered_assignment_id, requested_date, requested_shift_template_id, reason, status,
  payroll_blocking, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14
) RETURNING *;

-- name: UpdateShiftSwapRequestReview :one
UPDATE hrms.shift_swap_requests
SET status=$3, reviewed_by=$4, reviewed_at=now(), review_remarks=$5, payroll_blocking=$6, updated_by=$4
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListShiftSwapRequests :many
SELECT * FROM hrms.shift_swap_requests
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('requester_user_id')::uuid IS NULL OR requester_user_id = sqlc.narg('requester_user_id')::uuid)
  AND (sqlc.narg('target_user_id')::uuid IS NULL OR target_user_id = sqlc.narg('target_user_id')::uuid)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetShiftSwapRequest :one
SELECT * FROM hrms.shift_swap_requests
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateShiftScheduleEvent :one
INSERT INTO hrms.shift_schedule_events (
  tenant_id, source_type, source_id, action, from_status, to_status, actor_user_id, remarks, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
) RETURNING *;

-- name: ListShiftScheduleEvents :many
SELECT * FROM hrms.shift_schedule_events
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('source_type')::text IS NULL OR source_type = sqlc.narg('source_type')::text)
  AND (sqlc.narg('source_id')::uuid IS NULL OR source_id = sqlc.narg('source_id')::uuid)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetShiftScheduleSummary :many
WITH assignments AS (
  SELECT a.*
  FROM hrms.shift_schedule_assignments a
  WHERE a.tenant_id=$1
    AND NOT a.inactive
    AND a.schedule_date BETWEEN $2 AND $3
),
swaps AS (
  SELECT sw.*
  FROM hrms.shift_swap_requests sw
  WHERE sw.tenant_id=$1
    AND NOT sw.inactive
    AND sw.created_at::date BETWEEN $2 AND $3
),
requirements AS (
  SELECT r.*
  FROM hrms.staffing_requirements r
  WHERE r.tenant_id=$1
    AND NOT r.inactive
    AND r.status='active'
),
gap_rows AS (
  SELECT
    r.id,
    GREATEST(r.required_count - COALESCE(COUNT(a.id), 0)::int, 0) AS gap_count
  FROM requirements r
  LEFT JOIN assignments a ON a.status IN ('published','locked')
    AND (r.requirement_date IS NULL OR a.schedule_date = r.requirement_date)
    AND (r.start_date IS NULL OR a.schedule_date >= r.start_date)
    AND (r.end_date IS NULL OR a.schedule_date <= r.end_date)
    AND (r.day_of_week IS NULL OR EXTRACT(DOW FROM a.schedule_date)::int = r.day_of_week)
    AND (r.branch_id IS NULL OR a.branch_id = r.branch_id)
    AND (r.department_id IS NULL OR a.department_id = r.department_id)
    AND (r.attendance_location_id IS NULL OR a.attendance_location_id = r.attendance_location_id)
    AND (r.shift_template_id IS NULL OR a.shift_template_id = r.shift_template_id)
  GROUP BY r.id, r.required_count
)
SELECT 'assignments_total'::text AS metric, COUNT(*)::int AS metric_count FROM assignments
UNION ALL SELECT 'assignments_published', COUNT(*)::int FROM assignments WHERE status='published'
UNION ALL SELECT 'assignments_locked', COUNT(*)::int FROM assignments WHERE status='locked'
UNION ALL SELECT 'assignments_conflict', COUNT(*)::int FROM assignments WHERE has_conflict
UNION ALL SELECT 'payroll_blockers', COUNT(*)::int FROM assignments WHERE payroll_blocking
UNION ALL SELECT 'swap_requests_pending', COUNT(*)::int FROM swaps WHERE status='pending'
UNION ALL SELECT 'staffing_gap_positions', COUNT(*)::int FROM gap_rows WHERE gap_count > 0
UNION ALL SELECT 'staffing_gap_people', COALESCE(SUM(gap_count), 0)::int FROM gap_rows;

-- name: ListShiftStaffingGaps :many
WITH assignments AS (
  SELECT a.*
  FROM hrms.shift_schedule_assignments a
  WHERE a.tenant_id=$1
    AND NOT a.inactive
    AND a.schedule_date BETWEEN $2 AND $3
    AND a.status IN ('published','locked')
)
SELECT
  r.id AS requirement_id,
  r.name AS requirement_name,
  r.branch_id,
  r.department_id,
  r.attendance_location_id,
  r.shift_template_id,
  r.required_count,
  COALESCE(COUNT(a.id), 0)::int AS assigned_count,
  GREATEST(r.required_count - COALESCE(COUNT(a.id), 0)::int, 0)::int AS gap_count,
  r.priority,
  r.payroll_blocking
FROM hrms.staffing_requirements r
LEFT JOIN assignments a ON (r.requirement_date IS NULL OR a.schedule_date = r.requirement_date)
  AND (r.start_date IS NULL OR a.schedule_date >= r.start_date)
  AND (r.end_date IS NULL OR a.schedule_date <= r.end_date)
  AND (r.day_of_week IS NULL OR EXTRACT(DOW FROM a.schedule_date)::int = r.day_of_week)
  AND (r.branch_id IS NULL OR a.branch_id = r.branch_id)
  AND (r.department_id IS NULL OR a.department_id = r.department_id)
  AND (r.attendance_location_id IS NULL OR a.attendance_location_id = r.attendance_location_id)
  AND (r.shift_template_id IS NULL OR a.shift_template_id = r.shift_template_id)
WHERE r.tenant_id=$1
  AND NOT r.inactive
  AND r.status='active'
GROUP BY r.id, r.name, r.branch_id, r.department_id, r.attendance_location_id, r.shift_template_id, r.required_count, r.priority, r.payroll_blocking
HAVING GREATEST(r.required_count - COALESCE(COUNT(a.id), 0)::int, 0) > 0
ORDER BY r.payroll_blocking DESC, r.priority DESC, gap_count DESC, r.name ASC;
