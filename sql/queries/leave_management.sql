-- name: ListLeaveTypes :many
SELECT * FROM hrms.leave_types
WHERE tenant_id = $1 AND NOT inactive
ORDER BY name ASC;

-- name: GetLeaveType :one
SELECT * FROM hrms.leave_types
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetLeaveTypeByShortcode :one
SELECT * FROM hrms.leave_types
WHERE tenant_id = $1 AND lower(shortcode) = lower($2) AND NOT inactive;

-- name: CreateLeaveType :one
INSERT INTO hrms.leave_types (
    tenant_id,
    name,
    shortcode,
    description,
    is_paid,
    is_carry_forward,
    max_carry_forward,
    is_consecutive_limit,
    consecutive_days_limit,
    is_enabled,
    is_system,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: UpsertSystemLeaveType :one
INSERT INTO hrms.leave_types (
    tenant_id,
    name,
    shortcode,
    description,
    is_paid,
    is_carry_forward,
    max_carry_forward,
    is_consecutive_limit,
    consecutive_days_limit,
    is_enabled,
    is_system,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, TRUE, $11, $11
)
ON CONFLICT (tenant_id, lower(shortcode)) WHERE shortcode IS NOT NULL AND NOT inactive
DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    is_paid = EXCLUDED.is_paid,
    is_carry_forward = EXCLUDED.is_carry_forward,
    max_carry_forward = EXCLUDED.max_carry_forward,
    is_consecutive_limit = EXCLUDED.is_consecutive_limit,
    consecutive_days_limit = EXCLUDED.consecutive_days_limit,
    is_system = TRUE,
    is_enabled = hrms.leave_types.is_enabled,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: UpdateLeaveType :one
UPDATE hrms.leave_types
SET
    name = $3,
    shortcode = $4,
    description = $5,
    is_paid = $6,
    is_carry_forward = $7,
    max_carry_forward = $8,
    is_consecutive_limit = $9,
    consecutive_days_limit = $10,
    is_enabled = $11,
    updated_by = $12,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteLeaveType :exec
UPDATE hrms.leave_types
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListLeavePolicies :many
SELECT * FROM hrms.leave_policies
WHERE tenant_id = $1 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetLeavePolicy :one
SELECT * FROM hrms.leave_policies
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetLeavePolicyByTypeAndFY :one
SELECT * FROM hrms.leave_policies
WHERE tenant_id = $1 AND leave_type_id = $2 AND fy_id = $3 AND NOT inactive;

-- name: CreateLeavePolicy :one
INSERT INTO hrms.leave_policies (
    tenant_id,
    leave_type_id,
    fy_id,
    total_days,
    allocation_type,
    jan,
    feb,
    mar,
    apr,
    may,
    jun,
    jul,
    aug,
    sep,
    oct,
    nov,
    dec,
    is_sandwich_applicable,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $19
)
RETURNING *;

-- name: UpdateLeavePolicy :one
UPDATE hrms.leave_policies
SET
    leave_type_id = $3,
    fy_id = $4,
    total_days = $5,
    allocation_type = $6,
    jan = $7,
    feb = $8,
    mar = $9,
    apr = $10,
    may = $11,
    jun = $12,
    jul = $13,
    aug = $14,
    sep = $15,
    oct = $16,
    nov = $17,
    dec = $18,
    is_sandwich_applicable = $19,
    updated_by = $20,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteLeavePolicy :exec
UPDATE hrms.leave_policies
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListLeaveBalancesByUser :many
SELECT * FROM hrms.leave_balances
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetLeaveBalance :one
SELECT * FROM hrms.leave_balances
WHERE tenant_id = $1 AND user_id = $2 AND leave_type_id = $3 AND fy_id = $4 AND NOT inactive;

-- name: SoftDeleteLeaveBalance :exec
UPDATE hrms.leave_balances
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListLeavesByUser :many
SELECT * FROM hrms.leaves
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY applied_date DESC;

-- name: ListLeavesByFY :many
SELECT * FROM hrms.leaves
WHERE tenant_id = $1 AND fy_id = $2 AND NOT inactive
ORDER BY applied_date DESC;

-- name: ListLeaveReportRows :many
SELECT
    l.id,
    l.tenant_id,
    l.user_id,
    e.employee_code,
    e.firstname,
    e.lastname,
    e.reporting_manager_id,
    e.department_id,
    d.name AS department_name,
    e.designation_id,
    dg.name AS designation_name,
    l.leave_type_id,
    lt.name AS leave_type_name,
    lt.shortcode AS leave_type_shortcode,
    l.fy_id,
    fy.name AS financial_year_name,
    l.start_date,
    l.end_date,
    l.start_day_type,
    l.end_day_type,
    l.days,
    l.reason,
    l.status,
    l.is_sandwich,
    l.applied_date,
    l.created_at,
    l.updated_at
FROM hrms.leaves l
JOIN hrms.employees e ON e.tenant_id = l.tenant_id AND e.user_id = l.user_id AND NOT e.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.designations dg ON dg.tenant_id = e.tenant_id AND dg.id = e.designation_id AND NOT dg.inactive
LEFT JOIN hrms.leave_types lt ON lt.tenant_id = l.tenant_id AND lt.id = l.leave_type_id AND NOT lt.inactive
LEFT JOIN hrms.financial_years fy ON fy.tenant_id = l.tenant_id AND fy.id = l.fy_id AND NOT fy.inactive
WHERE l.tenant_id = sqlc.arg(tenant_id)
  AND NOT l.inactive
  AND (sqlc.narg(fy_id)::uuid IS NULL OR l.fy_id = sqlc.narg(fy_id)::uuid)
  AND (sqlc.narg(user_id)::uuid IS NULL OR l.user_id = sqlc.narg(user_id)::uuid)
  AND (sqlc.narg(department_id)::uuid IS NULL OR e.department_id = sqlc.narg(department_id)::uuid)
  AND (sqlc.narg(leave_type_id)::uuid IS NULL OR l.leave_type_id = sqlc.narg(leave_type_id)::uuid)
  AND (sqlc.narg(status)::text IS NULL OR l.status = sqlc.narg(status)::text)
  AND (sqlc.narg(start_date)::date IS NULL OR l.end_date >= sqlc.narg(start_date)::date)
  AND (sqlc.narg(end_date)::date IS NULL OR l.start_date <= sqlc.narg(end_date)::date)
ORDER BY l.start_date DESC, l.applied_date DESC;

-- name: ListManagerLeaveReportRows :many
SELECT
    l.id,
    l.tenant_id,
    l.user_id,
    e.employee_code,
    e.firstname,
    e.lastname,
    e.reporting_manager_id,
    e.department_id,
    d.name AS department_name,
    e.designation_id,
    dg.name AS designation_name,
    l.leave_type_id,
    lt.name AS leave_type_name,
    lt.shortcode AS leave_type_shortcode,
    l.fy_id,
    fy.name AS financial_year_name,
    l.start_date,
    l.end_date,
    l.start_day_type,
    l.end_day_type,
    l.days,
    l.reason,
    l.status,
    l.is_sandwich,
    l.applied_date,
    l.created_at,
    l.updated_at
FROM hrms.leaves l
JOIN hrms.employees e ON e.tenant_id = l.tenant_id AND e.user_id = l.user_id AND NOT e.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.designations dg ON dg.tenant_id = e.tenant_id AND dg.id = e.designation_id AND NOT dg.inactive
LEFT JOIN hrms.leave_types lt ON lt.tenant_id = l.tenant_id AND lt.id = l.leave_type_id AND NOT lt.inactive
LEFT JOIN hrms.financial_years fy ON fy.tenant_id = l.tenant_id AND fy.id = l.fy_id AND NOT fy.inactive
WHERE l.tenant_id = sqlc.arg(tenant_id)
  AND e.reporting_manager_id = sqlc.arg(manager_id)
  AND NOT l.inactive
  AND (sqlc.narg(fy_id)::uuid IS NULL OR l.fy_id = sqlc.narg(fy_id)::uuid)
  AND (sqlc.narg(user_id)::uuid IS NULL OR l.user_id = sqlc.narg(user_id)::uuid)
  AND (sqlc.narg(department_id)::uuid IS NULL OR e.department_id = sqlc.narg(department_id)::uuid)
  AND (sqlc.narg(leave_type_id)::uuid IS NULL OR l.leave_type_id = sqlc.narg(leave_type_id)::uuid)
  AND (sqlc.narg(status)::text IS NULL OR l.status = sqlc.narg(status)::text)
  AND (sqlc.narg(start_date)::date IS NULL OR l.end_date >= sqlc.narg(start_date)::date)
  AND (sqlc.narg(end_date)::date IS NULL OR l.start_date <= sqlc.narg(end_date)::date)
ORDER BY l.start_date DESC, l.applied_date DESC;

-- name: GetLeaveReportSummary :one
SELECT
    COUNT(*)::int AS total_requests,
    COALESCE(SUM(l.days), 0)::numeric AS total_days,
    COUNT(DISTINCT l.user_id)::int AS employee_count,
    COUNT(*) FILTER (WHERE l.status = 'pending')::int AS pending_count,
    COUNT(*) FILTER (WHERE l.status = 'approved')::int AS approved_count,
    COUNT(*) FILTER (WHERE l.status = 'rejected')::int AS rejected_count,
    COUNT(*) FILTER (WHERE l.status = 'canceled')::int AS canceled_count,
    COALESCE(SUM(l.days) FILTER (WHERE l.status = 'pending'), 0)::numeric AS pending_days,
    COALESCE(SUM(l.days) FILTER (WHERE l.status = 'approved'), 0)::numeric AS approved_days,
    COALESCE(SUM(l.days) FILTER (WHERE l.status = 'rejected'), 0)::numeric AS rejected_days,
    COALESCE(SUM(l.days) FILTER (WHERE l.status = 'canceled'), 0)::numeric AS canceled_days
FROM hrms.leaves l
JOIN hrms.employees e ON e.tenant_id = l.tenant_id AND e.user_id = l.user_id AND NOT e.inactive
WHERE l.tenant_id = sqlc.arg(tenant_id)
  AND NOT l.inactive
  AND (sqlc.narg(manager_id)::uuid IS NULL OR e.reporting_manager_id = sqlc.narg(manager_id)::uuid)
  AND (sqlc.narg(fy_id)::uuid IS NULL OR l.fy_id = sqlc.narg(fy_id)::uuid)
  AND (sqlc.narg(user_id)::uuid IS NULL OR l.user_id = sqlc.narg(user_id)::uuid)
  AND (sqlc.narg(department_id)::uuid IS NULL OR e.department_id = sqlc.narg(department_id)::uuid)
  AND (sqlc.narg(leave_type_id)::uuid IS NULL OR l.leave_type_id = sqlc.narg(leave_type_id)::uuid)
  AND (sqlc.narg(status)::text IS NULL OR l.status = sqlc.narg(status)::text)
  AND (sqlc.narg(start_date)::date IS NULL OR l.end_date >= sqlc.narg(start_date)::date)
  AND (sqlc.narg(end_date)::date IS NULL OR l.start_date <= sqlc.narg(end_date)::date);

-- name: GetLeave :one
SELECT * FROM hrms.leaves
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteLeave :exec
UPDATE hrms.leaves
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListLeaveApprovalsByApprover :many
SELECT * FROM hrms.leave_approvals
WHERE tenant_id = $1 AND approver_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: ListLeaveApprovalsByLeave :many
SELECT * FROM hrms.leave_approvals
WHERE tenant_id = $1 AND leave_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetLeaveApproval :one
SELECT * FROM hrms.leave_approvals
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteLeaveApproval :exec
UPDATE hrms.leave_approvals
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListLeaveLedgerByUser :many
SELECT * FROM hrms.leave_ledger
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: ListLeaveLedgerByLeave :many
SELECT * FROM hrms.leave_ledger
WHERE tenant_id = $1 AND leave_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetLeaveLedgerEntry :one
SELECT * FROM hrms.leave_ledger
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteLeaveLedgerEntry :exec
UPDATE hrms.leave_ledger
SET inactive = TRUE
WHERE tenant_id = $1 AND id = $2;

-- name: ListLeavePolicyTemplates :many
SELECT * FROM hrms.leave_policy_templates
WHERE NOT inactive AND (is_system OR tenant_id = $1)
ORDER BY is_system DESC, name ASC;

-- name: GetLeavePolicyTemplate :one
SELECT * FROM hrms.leave_policy_templates
WHERE id = $2 AND NOT inactive AND (is_system OR tenant_id = $1);

-- name: CreateLeavePolicyTemplate :one
INSERT INTO hrms.leave_policy_templates (
    tenant_id, name, code, description, is_system, effective_from, effective_to, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, FALSE, $5, $6, $7, $7
)
RETURNING *;

-- name: UpdateLeavePolicyTemplate :one
UPDATE hrms.leave_policy_templates
SET name = $3,
    code = $4,
    description = $5,
    effective_from = $6,
    effective_to = $7,
    updated_by = $8,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive AND NOT is_system
RETURNING *;

-- name: SoftDeleteLeavePolicyTemplate :exec
UPDATE hrms.leave_policy_templates
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive AND NOT is_system;

-- name: ListLeavePolicyTemplateRules :many
SELECT * FROM hrms.leave_policy_template_rules
WHERE tenant_id = $1 AND template_id = $2 AND NOT inactive
ORDER BY priority ASC, created_at ASC;

-- name: ListLeavePolicyTemplateRulesByTenant :many
SELECT * FROM hrms.leave_policy_template_rules
WHERE tenant_id = $1 AND NOT inactive
ORDER BY priority ASC, created_at ASC;

-- name: GetLeavePolicyTemplateRule :one
SELECT * FROM hrms.leave_policy_template_rules
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateLeavePolicyTemplateRule :one
INSERT INTO hrms.leave_policy_template_rules (
    tenant_id, template_id, leave_type_id, fy_id, employment_type_id, department_id, designation_id, probation_status,
    accrual_method, accrual_frequency, credit_days, credit_hours, annual_entitlement, min_worked_days, max_balance,
    carry_forward_enabled, max_carry_forward, carry_forward_expiry_months, encashment_enabled,
    encashment_limit, encashment_payable_percent, negative_balance_allowed, max_negative_balance,
    sandwich_applicable, include_holidays, include_weekoffs, requires_document_after_days,
    min_request_days, max_request_days, max_requests_per_year, accrual_day, lapse_unutilized, allow_half_day, requires_approval,
    calculation_config, priority, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14, $15,
    $16, $17, $18, $19,
    $20, $21, $22, $23,
    $24, $25, $26, $27,
    $28, $29, $30, $31, $32, $33, $34,
    $35, $36, $37, $37
)
RETURNING *;

-- name: UpdateLeavePolicyTemplateRule :one
UPDATE hrms.leave_policy_template_rules
SET leave_type_id = $3,
    fy_id = $4,
    employment_type_id = $5,
    department_id = $6,
    designation_id = $7,
    probation_status = $8,
    accrual_method = $9,
    accrual_frequency = $10,
    credit_days = $11,
    credit_hours = $12,
    annual_entitlement = $13,
    min_worked_days = $14,
    max_balance = $15,
    carry_forward_enabled = $16,
    max_carry_forward = $17,
    carry_forward_expiry_months = $18,
    encashment_enabled = $19,
    encashment_limit = $20,
    encashment_payable_percent = $21,
    negative_balance_allowed = $22,
    max_negative_balance = $23,
    sandwich_applicable = $24,
    include_holidays = $25,
    include_weekoffs = $26,
    requires_document_after_days = $27,
    min_request_days = $28,
    max_request_days = $29,
    max_requests_per_year = $30,
    accrual_day = $31,
    lapse_unutilized = $32,
    allow_half_day = $33,
    requires_approval = $34,
    calculation_config = $35,
    priority = $36,
    updated_by = $37,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteLeavePolicyTemplateRule :exec
UPDATE hrms.leave_policy_template_rules
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpsertEmployeeLeavePolicyAssignment :one
INSERT INTO hrms.employee_leave_policy_assignments (
    tenant_id, user_id, template_id, fy_id, effective_from, effective_to, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $7
)
ON CONFLICT (tenant_id, user_id, template_id, fy_id) WHERE NOT inactive
DO UPDATE SET
    effective_from = EXCLUDED.effective_from,
    effective_to = EXCLUDED.effective_to,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: ListEmployeeLeavePolicyAssignments :many
SELECT * FROM hrms.employee_leave_policy_assignments
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY effective_from DESC;

-- name: ListLeavePolicyAssignmentsByTemplate :many
SELECT * FROM hrms.employee_leave_policy_assignments
WHERE tenant_id = $1 AND template_id = $2 AND NOT inactive
ORDER BY effective_from DESC;

-- name: SoftDeleteEmployeeLeavePolicyAssignment :exec
UPDATE hrms.employee_leave_policy_assignments
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpsertLeaveBalance :one
INSERT INTO hrms.leave_balances (
    tenant_id, user_id, leave_type_id, fy_id, total_days, used_days, pending_days, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $8
)
ON CONFLICT (tenant_id, user_id, leave_type_id, fy_id) WHERE NOT inactive
DO UPDATE SET
    total_days = EXCLUDED.total_days,
    used_days = EXCLUDED.used_days,
    pending_days = EXCLUDED.pending_days,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: AddLeaveBalanceCredit :one
UPDATE hrms.leave_balances
SET total_days = total_days + $5,
    updated_by = $6,
    updated_at = NOW()
WHERE tenant_id = $1 AND user_id = $2 AND leave_type_id = $3 AND fy_id = $4 AND NOT inactive
RETURNING *;

-- name: UpdateLeaveBalancePending :one
UPDATE hrms.leave_balances
SET pending_days = pending_days + $5,
    updated_by = $6,
    updated_at = NOW()
WHERE tenant_id = $1 AND user_id = $2 AND leave_type_id = $3 AND fy_id = $4 AND NOT inactive
RETURNING *;

-- name: ReverseLeaveBalancePending :one
UPDATE hrms.leave_balances
SET pending_days = pending_days - $5,
    updated_by = $6,
    updated_at = NOW()
WHERE tenant_id = $1 AND user_id = $2 AND leave_type_id = $3 AND fy_id = $4 AND NOT inactive AND pending_days >= $5
RETURNING *;

-- name: MoveLeaveBalancePendingToUsed :one
UPDATE hrms.leave_balances
SET pending_days = pending_days - $5,
    used_days = used_days + $5,
    updated_by = $6,
    updated_at = NOW()
WHERE tenant_id = $1 AND user_id = $2 AND leave_type_id = $3 AND fy_id = $4 AND NOT inactive AND pending_days >= $5
RETURNING *;

-- name: CreateLeaveLedgerEntry :one
INSERT INTO hrms.leave_ledger (
    tenant_id, user_id, leave_type_id, fy_id, leave_id, transaction_type, days, remarks, source_type, source_id,
    balance_before, balance_after, pending_before, pending_after, used_before, used_after, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18
)
RETURNING *;

-- name: ListLeaveBalancesByTenantFY :many
SELECT * FROM hrms.leave_balances
WHERE tenant_id = $1 AND fy_id = $2 AND NOT inactive
ORDER BY user_id ASC, created_at DESC;

-- name: GetLeaveLedgerBySource :one
SELECT * FROM hrms.leave_ledger
WHERE tenant_id = $1 AND user_id = $2 AND leave_type_id = $3 AND fy_id = $4 AND source_type = $5 AND source_id = $6 AND NOT inactive
LIMIT 1;

-- name: CreateLeave :one
INSERT INTO hrms.leaves (
    tenant_id, user_id, leave_type_id, fy_id, start_date, end_date, start_day_type, end_day_type,
    days, reason, status, from_leave_type, to_leave_type, is_sandwich, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14, $15, $15
)
RETURNING *;

-- name: CreateLeaveApproval :one
INSERT INTO hrms.leave_approvals (
    tenant_id, leave_id, approver_id, status, remarks, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $6
)
RETURNING *;

-- name: CreateLeaveRequestMessage :one
INSERT INTO hrms.leave_request_messages (
    tenant_id, leave_id, sender_user_id, recipient_user_id, message_type, body, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $7
)
RETURNING *;

-- name: ListLeaveRequestMessages :many
SELECT * FROM hrms.leave_request_messages
WHERE tenant_id = $1
  AND leave_id = $2
  AND NOT inactive
ORDER BY created_at ASC;

-- name: ListOverlappingLeaves :many
SELECT * FROM hrms.leaves
WHERE tenant_id = $1
  AND user_id = $2
  AND status IN ('pending','approved')
  AND NOT inactive
  AND start_date <= $4
  AND end_date >= $3
ORDER BY start_date ASC;

-- name: CreateLeaveApprovalWorkflow :one
INSERT INTO hrms.leave_approval_workflows (
    tenant_id, name, code, description, is_default, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $6
)
RETURNING *;

-- name: ListLeaveApprovalWorkflows :many
SELECT * FROM hrms.leave_approval_workflows
WHERE tenant_id = $1 AND NOT inactive
ORDER BY is_default DESC, name ASC;

-- name: GetLeaveApprovalWorkflow :one
SELECT * FROM hrms.leave_approval_workflows
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetDefaultLeaveApprovalWorkflow :one
SELECT * FROM hrms.leave_approval_workflows
WHERE tenant_id = $1 AND is_default AND NOT inactive
LIMIT 1;

-- name: UpdateLeaveApprovalWorkflow :one
UPDATE hrms.leave_approval_workflows
SET name = $3,
    code = $4,
    description = $5,
    is_default = $6,
    updated_by = $7,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ClearDefaultLeaveApprovalWorkflows :exec
UPDATE hrms.leave_approval_workflows
SET is_default = FALSE, updated_by = $2, updated_at = NOW()
WHERE tenant_id = $1 AND is_default AND NOT inactive;

-- name: SoftDeleteLeaveApprovalWorkflow :exec
UPDATE hrms.leave_approval_workflows
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateLeaveApprovalWorkflowStep :one
INSERT INTO hrms.leave_approval_workflow_steps (
    tenant_id, workflow_id, step_order, name, approver_type, approver_user_id, approver_role,
    decision_rule, required_approvals, auto_approve, sla_hours, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: ListLeaveApprovalWorkflowSteps :many
SELECT * FROM hrms.leave_approval_workflow_steps
WHERE tenant_id = $1 AND workflow_id = $2 AND NOT inactive
ORDER BY step_order ASC, created_at ASC;

-- name: GetLeaveApprovalWorkflowStep :one
SELECT * FROM hrms.leave_approval_workflow_steps
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateLeaveApprovalWorkflowStep :one
UPDATE hrms.leave_approval_workflow_steps
SET step_order = $3,
    name = $4,
    approver_type = $5,
    approver_user_id = $6,
    approver_role = $7,
    decision_rule = $8,
    required_approvals = $9,
    auto_approve = $10,
    sla_hours = $11,
    updated_by = $12,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteLeaveApprovalWorkflowStep :exec
UPDATE hrms.leave_approval_workflow_steps
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateWorkflowLeaveApproval :one
INSERT INTO hrms.leave_approvals (
    tenant_id, leave_id, approver_id, status, remarks, workflow_id, workflow_step_id, step_order, decision_rule, required_approvals, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11
)
RETURNING *;

-- name: UpdateLeaveApprovalStatus :one
UPDATE hrms.leave_approvals
SET status = $3,
    remarks = $4,
    action_date = NOW(),
    updated_by = $5,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateLeaveStatus :one
UPDATE hrms.leaves
SET status = $3,
    updated_by = $4,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListPendingApprovalsByApprover :many
SELECT * FROM hrms.leave_approvals
WHERE tenant_id = $1 AND approver_id = $2 AND status = 'pending' AND NOT inactive
ORDER BY created_at ASC;
