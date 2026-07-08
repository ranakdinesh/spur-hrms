-- name: CreatePolicySet :one
INSERT INTO hrms.policy_sets (
    tenant_id, policy_kind, code, name, description, config,
    is_default, is_active, effective_from, effective_to, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11
)
RETURNING *;

-- name: UpdatePolicySet :one
UPDATE hrms.policy_sets
SET
    code = $3,
    name = $4,
    description = $5,
    config = $6,
    is_default = $7,
    is_active = $8,
    effective_from = $9,
    effective_to = $10,
    updated_by = $11,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListPolicySets :many
SELECT * FROM hrms.policy_sets
WHERE tenant_id = $1 AND policy_kind = $2 AND NOT inactive
ORDER BY is_default DESC, name ASC;

-- name: GetPolicySet :one
SELECT * FROM hrms.policy_sets
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeletePolicySet :exec
UPDATE hrms.policy_sets
SET inactive = TRUE, is_active = FALSE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreatePolicyAssignment :one
INSERT INTO hrms.policy_assignments (
    tenant_id, policy_set_id, policy_kind, scope_type, scope_id, role_code,
    priority, effective_from, effective_to, is_active, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11
)
RETURNING *;

-- name: ListPolicyAssignments :many
SELECT * FROM hrms.policy_assignments
WHERE tenant_id = $1 AND policy_kind = $2 AND NOT inactive
ORDER BY scope_type ASC, priority DESC, created_at DESC;

-- name: SoftDeletePolicyAssignment :exec
UPDATE hrms.policy_assignments
SET inactive = TRUE, is_active = FALSE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateLeavePolicyRule :one
INSERT INTO hrms.leave_policy_rules (
    tenant_id, policy_set_id, leave_type_id, grant_mode, accrual_frequency,
    entitlement_days, accrual_amount_per_period, prorate_joiners, probation_handling,
    rounding_rule, max_balance_cap, carry_forward_cap, encashment_eligible,
    negative_balance_allowed, insufficient_balance_action, expiry_days, allow_half_day,
    attachment_required_after_days, approval_workflow, sandwich_enabled,
    sandwich_include_weekly_off, sandwich_include_public_holiday,
    sandwich_same_leave_type_only, sandwich_across_leave_types,
    notice_required_after_days, notice_days, payroll_impact, rule_config,
    created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    $10, $11, $12, $13,
    $14, $15, $16, $17,
    $18, $19, $20,
    $21, $22,
    $23, $24,
    $25, $26, $27, $28,
    $29, $29
)
RETURNING *;

-- name: ListLeavePolicyRules :many
SELECT * FROM hrms.leave_policy_rules
WHERE tenant_id = $1 AND policy_set_id = $2 AND NOT inactive
ORDER BY created_at ASC;

-- name: ResolvePolicySet :one
WITH candidates AS (
    SELECT ps.*,
           CASE pa.scope_type
               WHEN 'employee' THEN 60
               WHEN 'designation' THEN 50
               WHEN 'workforce_type' THEN 50
               WHEN 'role_group' THEN 45
               WHEN 'department' THEN 40
               WHEN 'branch' THEN 30
               WHEN 'tenant' THEN 20
               ELSE 0
           END + pa.priority AS precedence,
           pa.created_at AS assignment_created_at
    FROM hrms.policy_assignments pa
    JOIN hrms.policy_sets ps ON ps.tenant_id = pa.tenant_id
        AND ps.id = pa.policy_set_id
        AND ps.policy_kind = pa.policy_kind
        AND ps.is_active
        AND NOT ps.inactive
    WHERE pa.tenant_id = $1
      AND pa.policy_kind = $2
      AND pa.is_active
      AND NOT pa.inactive
      AND (pa.effective_from IS NULL OR pa.effective_from <= $9)
      AND (pa.effective_to IS NULL OR pa.effective_to >= $9)
      AND (ps.effective_from IS NULL OR ps.effective_from <= $9)
      AND (ps.effective_to IS NULL OR ps.effective_to >= $9)
      AND (
        pa.scope_type = 'tenant'
        OR (pa.scope_type = 'employee' AND pa.scope_id = $3)
        OR (pa.scope_type = 'designation' AND pa.scope_id = $4)
        OR (pa.scope_type = 'workforce_type' AND pa.scope_id = $5)
        OR (pa.scope_type = 'department' AND pa.scope_id = $6)
        OR (pa.scope_type = 'branch' AND pa.scope_id = $7)
        OR (pa.scope_type = 'role_group' AND lower(pa.role_code) = ANY($8::text[]))
      )
    UNION ALL
    SELECT ps.*, 10 AS precedence, ps.created_at AS assignment_created_at
    FROM hrms.policy_sets ps
    WHERE ps.tenant_id = $1
      AND ps.policy_kind = $2
      AND ps.is_default
      AND ps.is_active
      AND NOT ps.inactive
      AND (ps.effective_from IS NULL OR ps.effective_from <= $9)
      AND (ps.effective_to IS NULL OR ps.effective_to >= $9)
)
SELECT * FROM candidates
ORDER BY precedence DESC, assignment_created_at DESC
LIMIT 1;

-- name: ListPolicyResolutionCandidates :many
WITH candidates AS (
    SELECT ps.id, ps.tenant_id, ps.policy_kind, ps.code, ps.name, pa.scope_type, pa.scope_id, pa.role_code,
           CASE pa.scope_type
               WHEN 'employee' THEN 60
               WHEN 'designation' THEN 50
               WHEN 'workforce_type' THEN 50
               WHEN 'role_group' THEN 45
               WHEN 'department' THEN 40
               WHEN 'branch' THEN 30
               WHEN 'tenant' THEN 20
               ELSE 0
           END + pa.priority AS precedence
    FROM hrms.policy_assignments pa
    JOIN hrms.policy_sets ps ON ps.tenant_id = pa.tenant_id
        AND ps.id = pa.policy_set_id
        AND ps.policy_kind = pa.policy_kind
        AND ps.is_active
        AND NOT ps.inactive
    WHERE pa.tenant_id = $1
      AND pa.policy_kind = $2
      AND pa.is_active
      AND NOT pa.inactive
      AND (pa.effective_from IS NULL OR pa.effective_from <= $9)
      AND (pa.effective_to IS NULL OR pa.effective_to >= $9)
      AND (ps.effective_from IS NULL OR ps.effective_from <= $9)
      AND (ps.effective_to IS NULL OR ps.effective_to >= $9)
      AND (
        pa.scope_type = 'tenant'
        OR (pa.scope_type = 'employee' AND pa.scope_id = $3)
        OR (pa.scope_type = 'designation' AND pa.scope_id = $4)
        OR (pa.scope_type = 'workforce_type' AND pa.scope_id = $5)
        OR (pa.scope_type = 'department' AND pa.scope_id = $6)
        OR (pa.scope_type = 'branch' AND pa.scope_id = $7)
        OR (pa.scope_type = 'role_group' AND lower(pa.role_code) = ANY($8::text[]))
      )
    UNION ALL
    SELECT ps.id, ps.tenant_id, ps.policy_kind, ps.code, ps.name, 'default'::text AS scope_type, NULL::uuid AS scope_id, NULL::varchar AS role_code, 10 AS precedence
    FROM hrms.policy_sets ps
    WHERE ps.tenant_id = $1
      AND ps.policy_kind = $2
      AND ps.is_default
      AND ps.is_active
      AND NOT ps.inactive
      AND (ps.effective_from IS NULL OR ps.effective_from <= $9)
      AND (ps.effective_to IS NULL OR ps.effective_to >= $9)
)
SELECT * FROM candidates
ORDER BY precedence DESC, name ASC;
