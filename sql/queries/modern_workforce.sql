-- name: CreateWorkerType :one
INSERT INTO hrms.worker_types (
    tenant_id,
    code,
    name,
    classification_group,
    description,
    attendance_mode,
    pay_mode,
    tds_section,
    pf_applicable,
    esic_applicable,
    pt_applicable,
    lwf_applicable,
    clra_applicable,
    leave_applicable,
    overtime_applicable,
    requires_agreement,
    requires_invoice,
    requires_attendance,
    statutory_defaults,
    compliance_notes,
    is_system_default,
    sort_order,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
    $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $23
)
RETURNING *;

-- name: ListWorkerTypes :many
SELECT * FROM hrms.worker_types
WHERE tenant_id = $1 AND NOT inactive
ORDER BY sort_order ASC, name ASC;

-- name: GetWorkerType :one
SELECT * FROM hrms.worker_types
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetWorkerTypeByCode :one
SELECT * FROM hrms.worker_types
WHERE tenant_id = $1 AND code = $2 AND NOT inactive;

-- name: UpdateWorkerType :one
UPDATE hrms.worker_types
SET name = $3,
    classification_group = $4,
    description = $5,
    attendance_mode = $6,
    pay_mode = $7,
    tds_section = $8,
    pf_applicable = $9,
    esic_applicable = $10,
    pt_applicable = $11,
    lwf_applicable = $12,
    clra_applicable = $13,
    leave_applicable = $14,
    overtime_applicable = $15,
    requires_agreement = $16,
    requires_invoice = $17,
    requires_attendance = $18,
    statutory_defaults = $19,
    compliance_notes = $20,
    sort_order = $21,
    updated_by = $22
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteWorkerType :exec
UPDATE hrms.worker_types
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CountWorkerClassificationRules :one
SELECT COUNT(*) FROM hrms.worker_classification_rules
WHERE tenant_id = $1 AND worker_type_id = $2 AND NOT inactive;

-- name: CreateWorkerClassificationRule :one
INSERT INTO hrms.worker_classification_rules (
    tenant_id,
    worker_type_id,
    rule_name,
    rule_type,
    priority,
    conditions,
    outcome,
    notes,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $9
)
RETURNING *;

-- name: ListWorkerClassificationRules :many
SELECT * FROM hrms.worker_classification_rules
WHERE tenant_id = $1
  AND (sqlc.narg('worker_type_id')::uuid IS NULL OR worker_type_id = sqlc.narg('worker_type_id')::uuid)
  AND NOT inactive
ORDER BY priority ASC, rule_name ASC;

-- name: GetWorkerClassificationRule :one
SELECT * FROM hrms.worker_classification_rules
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateWorkerClassificationRule :one
UPDATE hrms.worker_classification_rules
SET worker_type_id = $3,
    rule_name = $4,
    rule_type = $5,
    priority = $6,
    conditions = $7,
    outcome = $8,
    notes = $9,
    updated_by = $10
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteWorkerClassificationRule :exec
UPDATE hrms.worker_classification_rules
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateWorkerProfile :one
INSERT INTO hrms.worker_profiles (
    tenant_id,
    worker_type_id,
    employee_id,
    employee_user_id,
    worker_code,
    display_name,
    legal_name,
    email,
    mobile,
    profile_status,
    start_date,
    end_date,
    branch_id,
    department_id,
    reporting_manager_id,
    work_location_label,
    source_partner,
    external_reference,
    compliance_status,
    payroll_status,
    notes,
    metadata,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $23
)
RETURNING *;

-- name: UpdateWorkerProfile :one
UPDATE hrms.worker_profiles
SET worker_type_id = $3,
    employee_id = $4,
    employee_user_id = $5,
    worker_code = $6,
    display_name = $7,
    legal_name = $8,
    email = $9,
    mobile = $10,
    profile_status = $11,
    start_date = $12,
    end_date = $13,
    branch_id = $14,
    department_id = $15,
    reporting_manager_id = $16,
    work_location_label = $17,
    source_partner = $18,
    external_reference = $19,
    compliance_status = $20,
    payroll_status = $21,
    notes = $22,
    metadata = $23,
    updated_by = $24
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetWorkerProfile :one
SELECT * FROM hrms.worker_profiles
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetWorkerProfileByEmployeeID :one
SELECT * FROM hrms.worker_profiles
WHERE tenant_id = $1 AND employee_id = $2 AND NOT inactive;

-- name: ListWorkerProfiles :many
SELECT
    wp.id,
    wp.tenant_id,
    wp.worker_type_id,
    wt.code AS worker_type_code,
    wt.name AS worker_type_name,
    wt.classification_group,
    wt.attendance_mode,
    wt.pay_mode,
    wp.employee_id,
    wp.employee_user_id,
    e.employee_code,
    wp.worker_code,
    wp.display_name,
    wp.legal_name,
    wp.email,
    wp.mobile,
    wp.profile_status,
    wp.start_date,
    wp.end_date,
    wp.branch_id,
    b.branch_name AS branch_name,
    wp.department_id,
    d.name AS department_name,
    wp.reporting_manager_id,
    wp.work_location_label,
    wp.source_partner,
    wp.external_reference,
    wp.compliance_status,
    wp.payroll_status,
    wp.notes,
    wp.metadata,
    wp.inactive,
    wp.created_at,
    wp.created_by,
    wp.updated_at,
    wp.updated_by
FROM hrms.worker_profiles wp
JOIN hrms.worker_types wt ON wt.tenant_id = wp.tenant_id AND wt.id = wp.worker_type_id AND NOT wt.inactive
LEFT JOIN hrms.employees e ON e.tenant_id = wp.tenant_id AND e.id = wp.employee_id AND NOT e.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = wp.tenant_id AND b.id = wp.branch_id AND NOT b.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = wp.tenant_id AND d.id = wp.department_id AND NOT d.inactive
WHERE wp.tenant_id = $1
  AND (sqlc.narg('worker_type_id')::uuid IS NULL OR wp.worker_type_id = sqlc.narg('worker_type_id')::uuid)
  AND (sqlc.narg('classification_group')::text IS NULL OR wt.classification_group = sqlc.narg('classification_group')::text)
  AND (sqlc.narg('profile_status')::text IS NULL OR wp.profile_status = sqlc.narg('profile_status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.worker_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.email ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.mobile ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.employee_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wt.name ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT wp.inactive
ORDER BY wp.updated_at DESC, wp.display_name ASC;

-- name: SoftDeleteWorkerProfile :exec
UPDATE hrms.worker_profiles
SET inactive = TRUE,
    profile_status = 'ended',
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;
