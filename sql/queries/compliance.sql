-- name: CreateComplianceRule :one
INSERT INTO hrms.compliance_rules (
    tenant_id, code, title, description, category, scope, severity,
    classification_group, worker_type_id, engagement_type, branch_id, department_id,
    country_code, state_code, trigger_event, default_due_days, recurring_days,
    requires_evidence, evidence_label, auto_detect_key, blocks_payroll,
    is_active, effective_from, effective_to, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12,
    $13, $14, $15, $16, $17,
    $18, $19, $20, $21,
    $22, $23, $24, $25, $26, $26
)
RETURNING *;

-- name: UpdateComplianceRule :one
UPDATE hrms.compliance_rules
SET code = $3,
    title = $4,
    description = $5,
    category = $6,
    scope = $7,
    severity = $8,
    classification_group = $9,
    worker_type_id = $10,
    engagement_type = $11,
    branch_id = $12,
    department_id = $13,
    country_code = $14,
    state_code = $15,
    trigger_event = $16,
    default_due_days = $17,
    recurring_days = $18,
    requires_evidence = $19,
    evidence_label = $20,
    auto_detect_key = $21,
    blocks_payroll = $22,
    is_active = $23,
    effective_from = $24,
    effective_to = $25,
    metadata = $26,
    updated_by = $27
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetComplianceRule :one
SELECT * FROM hrms.compliance_rules
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetComplianceRuleByCode :one
SELECT * FROM hrms.compliance_rules
WHERE tenant_id = $1 AND lower(code) = lower($2) AND NOT inactive;

-- name: ListComplianceRules :many
SELECT
    cr.*,
    wt.name AS worker_type_name,
    b.branch_name AS branch_name,
    d.name AS department_name
FROM hrms.compliance_rules cr
LEFT JOIN hrms.worker_types wt ON wt.tenant_id = cr.tenant_id AND wt.id = cr.worker_type_id AND NOT wt.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = cr.tenant_id AND b.id = cr.branch_id AND NOT b.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = cr.tenant_id AND d.id = cr.department_id AND NOT d.inactive
WHERE cr.tenant_id = $1
  AND (sqlc.narg('category')::text IS NULL OR cr.category = sqlc.narg('category')::text)
  AND (sqlc.narg('scope')::text IS NULL OR cr.scope = sqlc.narg('scope')::text)
  AND (sqlc.narg('severity')::text IS NULL OR cr.severity = sqlc.narg('severity')::text)
  AND (sqlc.narg('is_active')::boolean IS NULL OR cr.is_active = sqlc.narg('is_active')::boolean)
  AND (
      sqlc.narg('search')::text IS NULL
      OR cr.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR cr.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR cr.description ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT cr.inactive
ORDER BY cr.category ASC, cr.severity DESC, cr.title ASC;

-- name: ListActiveComplianceRulesForWorker :many
SELECT cr.*
FROM hrms.compliance_rules cr
JOIN hrms.worker_profiles wp ON wp.tenant_id = cr.tenant_id AND wp.id = $2 AND NOT wp.inactive
JOIN hrms.worker_types wt ON wt.tenant_id = wp.tenant_id AND wt.id = wp.worker_type_id AND NOT wt.inactive
WHERE cr.tenant_id = $1
  AND cr.is_active
  AND NOT cr.inactive
  AND cr.scope IN ('worker', 'worker_or_engagement')
  AND (cr.worker_type_id IS NULL OR cr.worker_type_id = wp.worker_type_id)
  AND (cr.classification_group IS NULL OR cr.classification_group = wt.classification_group)
  AND (cr.branch_id IS NULL OR cr.branch_id = wp.branch_id)
  AND (cr.department_id IS NULL OR cr.department_id = wp.department_id)
  AND (cr.effective_from IS NULL OR cr.effective_from <= CURRENT_DATE)
  AND (cr.effective_to IS NULL OR cr.effective_to >= CURRENT_DATE)
ORDER BY cr.category ASC, cr.title ASC;

-- name: ListActiveComplianceRulesForEngagement :many
SELECT cr.*
FROM hrms.compliance_rules cr
JOIN hrms.engagements e ON e.tenant_id = cr.tenant_id AND e.id = $2 AND NOT e.inactive
JOIN hrms.worker_profiles wp ON wp.tenant_id = e.tenant_id AND wp.id = e.worker_profile_id AND NOT wp.inactive
JOIN hrms.worker_types wt ON wt.tenant_id = wp.tenant_id AND wt.id = wp.worker_type_id AND NOT wt.inactive
WHERE cr.tenant_id = $1
  AND cr.is_active
  AND NOT cr.inactive
  AND cr.scope IN ('engagement', 'worker_or_engagement')
  AND (cr.worker_type_id IS NULL OR cr.worker_type_id = wp.worker_type_id)
  AND (cr.classification_group IS NULL OR cr.classification_group = wt.classification_group)
  AND (cr.engagement_type IS NULL OR cr.engagement_type = e.engagement_type)
  AND (cr.branch_id IS NULL OR cr.branch_id = COALESCE(e.branch_id, wp.branch_id))
  AND (cr.department_id IS NULL OR cr.department_id = COALESCE(e.department_id, wp.department_id))
  AND (cr.effective_from IS NULL OR cr.effective_from <= CURRENT_DATE)
  AND (cr.effective_to IS NULL OR cr.effective_to >= CURRENT_DATE)
ORDER BY cr.category ASC, cr.title ASC;

-- name: SoftDeleteComplianceRule :exec
UPDATE hrms.compliance_rules
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateComplianceChecklistItem :one
INSERT INTO hrms.compliance_checklist_items (
    tenant_id, rule_id, worker_profile_id, engagement_id, status, due_date,
    evidence_path, evidence_file_name, evidence_content_type, detected_value,
    notes, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10,
    $11, $12, $13, $13
)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetComplianceChecklistItem :one
SELECT * FROM hrms.compliance_checklist_items
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListComplianceChecklistItems :many
SELECT
    ci.*,
    cr.code AS rule_code,
    cr.title AS rule_title,
    cr.category AS rule_category,
    cr.scope AS rule_scope,
    cr.severity AS rule_severity,
    cr.requires_evidence,
    cr.evidence_label,
    cr.blocks_payroll,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    e.title AS engagement_title,
    e.engagement_code
FROM hrms.compliance_checklist_items ci
JOIN hrms.compliance_rules cr ON cr.tenant_id = ci.tenant_id AND cr.id = ci.rule_id AND NOT cr.inactive
LEFT JOIN hrms.worker_profiles wp ON wp.tenant_id = ci.tenant_id AND wp.id = ci.worker_profile_id AND NOT wp.inactive
LEFT JOIN hrms.engagements e ON e.tenant_id = ci.tenant_id AND e.id = ci.engagement_id AND NOT e.inactive
WHERE ci.tenant_id = $1
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR ci.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('engagement_id')::uuid IS NULL OR ci.engagement_id = sqlc.narg('engagement_id')::uuid)
  AND (sqlc.narg('rule_id')::uuid IS NULL OR ci.rule_id = sqlc.narg('rule_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR ci.status = sqlc.narg('status')::text)
  AND (sqlc.narg('category')::text IS NULL OR cr.category = sqlc.narg('category')::text)
  AND (sqlc.narg('due_before')::date IS NULL OR ci.due_date <= sqlc.narg('due_before')::date)
  AND (
      sqlc.narg('search')::text IS NULL
      OR cr.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR cr.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.worker_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.engagement_code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT ci.inactive
ORDER BY
    CASE ci.status
        WHEN 'non_compliant' THEN 1
        WHEN 'expired' THEN 2
        WHEN 'pending' THEN 3
        WHEN 'in_review' THEN 4
        WHEN 'waived' THEN 5
        WHEN 'compliant' THEN 6
        ELSE 7
    END,
    ci.due_date ASC NULLS LAST,
    cr.severity DESC,
    cr.title ASC;

-- name: UpdateComplianceChecklistStatus :one
UPDATE hrms.compliance_checklist_items
SET status = $3,
    reviewed_at = CASE WHEN $3 IN ('compliant', 'non_compliant', 'not_applicable') THEN now() ELSE reviewed_at END,
    reviewed_by = CASE WHEN $3 IN ('compliant', 'non_compliant', 'not_applicable') THEN $4 ELSE reviewed_by END,
    completed_at = CASE WHEN $3 = 'compliant' THEN now() ELSE completed_at END,
    notes = $5,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateComplianceChecklistEvidence :one
UPDATE hrms.compliance_checklist_items
SET evidence_path = $3,
    evidence_file_name = $4,
    evidence_content_type = $5,
    evidence_uploaded_at = now(),
    evidence_uploaded_by = $6,
    status = CASE WHEN status IN ('pending', 'non_compliant', 'expired') THEN 'in_review' ELSE status END,
    notes = COALESCE($7, notes),
    updated_by = $6
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: WaiveComplianceChecklistItem :one
UPDATE hrms.compliance_checklist_items
SET status = 'waived',
    waiver_reason = $3,
    waiver_until = $4,
    waived_at = now(),
    waived_by = $5,
    notes = COALESCE($6, notes),
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: RefreshComplianceChecklistDueStatus :exec
UPDATE hrms.compliance_checklist_items
SET status = 'expired'
WHERE tenant_id = $1
  AND status IN ('pending', 'in_review', 'non_compliant')
  AND due_date IS NOT NULL
  AND due_date < CURRENT_DATE
  AND NOT inactive;

-- name: SoftDeleteComplianceChecklistItem :exec
UPDATE hrms.compliance_checklist_items
SET inactive = TRUE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetComplianceSummary :many
SELECT
    COALESCE(cr.category, 'all') AS category,
    ci.status,
    COUNT(*)::int AS item_count,
    COUNT(*) FILTER (WHERE cr.blocks_payroll)::int AS payroll_blocker_count,
    COUNT(*) FILTER (WHERE ci.due_date IS NOT NULL AND ci.due_date <= CURRENT_DATE + INTERVAL '7 days' AND ci.status IN ('pending', 'in_review', 'non_compliant'))::int AS due_soon_count
FROM hrms.compliance_checklist_items ci
JOIN hrms.compliance_rules cr ON cr.tenant_id = ci.tenant_id AND cr.id = ci.rule_id AND NOT cr.inactive
WHERE ci.tenant_id = $1 AND NOT ci.inactive
GROUP BY GROUPING SETS ((cr.category, ci.status), (ci.status))
ORDER BY category ASC, status ASC;

-- name: CreateComplianceEvent :one
INSERT INTO hrms.compliance_events (
    tenant_id, checklist_item_id, rule_id, event_type, from_status, to_status,
    comment, actor_id, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9
)
RETURNING *;

-- name: ListComplianceEvents :many
SELECT * FROM hrms.compliance_events
WHERE tenant_id = $1
  AND (sqlc.narg('checklist_item_id')::uuid IS NULL OR checklist_item_id = sqlc.narg('checklist_item_id')::uuid)
  AND (sqlc.narg('rule_id')::uuid IS NULL OR rule_id = sqlc.narg('rule_id')::uuid)
ORDER BY created_at DESC
LIMIT 100;
