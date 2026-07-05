-- name: CreateERCaseCategory :one
INSERT INTO hrms.er_case_categories (
    tenant_id, code, name, case_family, description, default_severity, default_owner_role, is_active, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, sqlc.narg('description')::text, $5, sqlc.narg('default_owner_role')::text, $6, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: UpdateERCaseCategory :one
UPDATE hrms.er_case_categories
SET code = $3,
    name = $4,
    case_family = $5,
    description = sqlc.narg('description')::text,
    default_severity = $6,
    default_owner_role = sqlc.narg('default_owner_role')::text,
    is_active = $7,
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListERCaseCategories :many
SELECT *
FROM hrms.er_case_categories
WHERE tenant_id = $1
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND NOT inactive
ORDER BY is_active DESC, name ASC;

-- name: GetERCaseCategory :one
SELECT *
FROM hrms.er_case_categories
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteERCaseCategory :exec
UPDATE hrms.er_case_categories
SET inactive = TRUE,
    is_active = FALSE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateERCase :one
INSERT INTO hrms.er_cases (
    tenant_id, case_number, source_hr_case_id, category_id, title, intake_summary, case_family, severity, status,
    confidentiality_level, complainant_user_id, subject_employee_user_id, owner_user_id, owner_role,
    investigation_lead_user_id, legal_hold, legal_hold_reason, legal_hold_at, legal_hold_by, due_at,
    privacy_notes, created_by, updated_by
) VALUES (
    $1, $2, sqlc.narg('source_hr_case_id')::uuid, sqlc.narg('category_id')::uuid, $3, $4, $5, $6, $7,
    $8, sqlc.narg('complainant_user_id')::uuid, sqlc.narg('subject_employee_user_id')::uuid,
    sqlc.narg('owner_user_id')::uuid, sqlc.narg('owner_role')::text, sqlc.narg('investigation_lead_user_id')::uuid,
    $9, sqlc.narg('legal_hold_reason')::text, sqlc.narg('legal_hold_at')::timestamptz, sqlc.narg('legal_hold_by')::uuid,
    sqlc.narg('due_at')::timestamptz, sqlc.narg('privacy_notes')::text, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: UpdateERCase :one
UPDATE hrms.er_cases
SET source_hr_case_id = sqlc.narg('source_hr_case_id')::uuid,
    category_id = sqlc.narg('category_id')::uuid,
    title = $3,
    intake_summary = $4,
    case_family = $5,
    severity = $6,
    confidentiality_level = $7,
    complainant_user_id = sqlc.narg('complainant_user_id')::uuid,
    subject_employee_user_id = sqlc.narg('subject_employee_user_id')::uuid,
    owner_user_id = sqlc.narg('owner_user_id')::uuid,
    owner_role = sqlc.narg('owner_role')::text,
    investigation_lead_user_id = sqlc.narg('investigation_lead_user_id')::uuid,
    due_at = sqlc.narg('due_at')::timestamptz,
    privacy_notes = sqlc.narg('privacy_notes')::text,
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateERCaseStatus :one
UPDATE hrms.er_cases
SET status = $3,
    resolution_summary = CASE WHEN $3 IN ('closed','cancelled') THEN sqlc.narg('resolution_summary')::text ELSE resolution_summary END,
    closed_at = CASE WHEN $3 IN ('closed','cancelled') THEN NOW() ELSE closed_at END,
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateERCaseLegalHold :one
UPDATE hrms.er_cases
SET legal_hold = $3,
    legal_hold_reason = CASE WHEN $3 THEN sqlc.narg('legal_hold_reason')::text ELSE NULL END,
    legal_hold_at = CASE WHEN $3 THEN NOW() ELSE NULL END,
    legal_hold_by = CASE WHEN $3 THEN sqlc.narg('actor_id')::uuid ELSE NULL END,
    confidentiality_level = CASE WHEN $3 THEN 'legal_hold' ELSE 'restricted' END,
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListERCases :many
SELECT
    c.*,
    cat.name AS category_name,
    cat.code AS category_code,
    COALESCE(allegation_counts.allegation_count, 0)::int AS allegation_count,
    COALESCE(evidence_counts.evidence_count, 0)::int AS evidence_count,
    COALESCE(action_counts.open_action_count, 0)::int AS open_action_count
FROM hrms.er_cases c
LEFT JOIN hrms.er_case_categories cat ON cat.tenant_id = c.tenant_id AND cat.id = c.category_id AND NOT cat.inactive
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS allegation_count
    FROM hrms.er_allegations a
    WHERE a.tenant_id = c.tenant_id AND a.er_case_id = c.id AND NOT a.inactive
) allegation_counts ON TRUE
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS evidence_count
    FROM hrms.er_evidence_attachments e
    WHERE e.tenant_id = c.tenant_id AND e.er_case_id = c.id AND NOT e.inactive
) evidence_counts ON TRUE
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS open_action_count
    FROM hrms.er_action_plans ap
    WHERE ap.tenant_id = c.tenant_id AND ap.er_case_id = c.id AND ap.status IN ('pending','in_progress','overdue') AND NOT ap.inactive
) action_counts ON TRUE
WHERE c.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('status')::text IS NULL OR c.status = sqlc.narg('status')::text)
  AND (sqlc.narg('severity')::text IS NULL OR c.severity = sqlc.narg('severity')::text)
  AND (sqlc.narg('case_family')::text IS NULL OR c.case_family = sqlc.narg('case_family')::text)
  AND (sqlc.narg('category_id')::uuid IS NULL OR c.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('owner_user_id')::uuid IS NULL OR c.owner_user_id = sqlc.narg('owner_user_id')::uuid)
  AND (sqlc.narg('subject_employee_user_id')::uuid IS NULL OR c.subject_employee_user_id = sqlc.narg('subject_employee_user_id')::uuid)
  AND (sqlc.narg('complainant_user_id')::uuid IS NULL OR c.complainant_user_id = sqlc.narg('complainant_user_id')::uuid)
  AND (sqlc.narg('legal_hold')::boolean IS NULL OR c.legal_hold = sqlc.narg('legal_hold')::boolean)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.case_number ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.intake_summary ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT c.inactive
ORDER BY
    CASE c.severity WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 ELSE 4 END,
    CASE c.status WHEN 'intake' THEN 1 WHEN 'triage' THEN 2 WHEN 'investigation' THEN 3 WHEN 'findings' THEN 4 ELSE 5 END,
    c.legal_hold DESC,
    c.due_at ASC NULLS LAST,
    c.updated_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountERCases :one
SELECT COUNT(*)
FROM hrms.er_cases c
WHERE c.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('status')::text IS NULL OR c.status = sqlc.narg('status')::text)
  AND (sqlc.narg('severity')::text IS NULL OR c.severity = sqlc.narg('severity')::text)
  AND (sqlc.narg('case_family')::text IS NULL OR c.case_family = sqlc.narg('case_family')::text)
  AND (sqlc.narg('category_id')::uuid IS NULL OR c.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('owner_user_id')::uuid IS NULL OR c.owner_user_id = sqlc.narg('owner_user_id')::uuid)
  AND (sqlc.narg('subject_employee_user_id')::uuid IS NULL OR c.subject_employee_user_id = sqlc.narg('subject_employee_user_id')::uuid)
  AND (sqlc.narg('complainant_user_id')::uuid IS NULL OR c.complainant_user_id = sqlc.narg('complainant_user_id')::uuid)
  AND (sqlc.narg('legal_hold')::boolean IS NULL OR c.legal_hold = sqlc.narg('legal_hold')::boolean)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.case_number ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.intake_summary ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT c.inactive;

-- name: GetERCase :one
SELECT
    c.*,
    cat.name AS category_name,
    cat.code AS category_code
FROM hrms.er_cases c
LEFT JOIN hrms.er_case_categories cat ON cat.tenant_id = c.tenant_id AND cat.id = c.category_id AND NOT cat.inactive
WHERE c.tenant_id = $1 AND c.id = $2 AND NOT c.inactive;

-- name: CreateERCaseParty :one
INSERT INTO hrms.er_case_parties (
    tenant_id, er_case_id, party_user_id, party_name, party_role, representation_notes, contact_notes, created_by, updated_by
) VALUES (
    $1, $2, sqlc.narg('party_user_id')::uuid, sqlc.narg('party_name')::text, $3, sqlc.narg('representation_notes')::text, sqlc.narg('contact_notes')::text, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListERCaseParties :many
SELECT *
FROM hrms.er_case_parties
WHERE tenant_id = $1 AND er_case_id = $2 AND NOT inactive
ORDER BY party_role ASC, created_at ASC;

-- name: CreateERAllegation :one
INSERT INTO hrms.er_allegations (
    tenant_id, er_case_id, allegation_type, incident_date, incident_location, description, policy_reference, status, created_by, updated_by
) VALUES (
    $1, $2, $3, sqlc.narg('incident_date')::date, sqlc.narg('incident_location')::text, $4, sqlc.narg('policy_reference')::text, $5, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListERAllegations :many
SELECT *
FROM hrms.er_allegations
WHERE tenant_id = $1 AND er_case_id = $2 AND NOT inactive
ORDER BY created_at ASC;

-- name: CreateERInvestigationStep :one
INSERT INTO hrms.er_investigation_steps (
    tenant_id, er_case_id, step_type, title, description, owner_user_id, due_at, completed_at, status, outcome_notes, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, sqlc.narg('description')::text, sqlc.narg('owner_user_id')::uuid, sqlc.narg('due_at')::timestamptz,
    sqlc.narg('completed_at')::timestamptz, $5, sqlc.narg('outcome_notes')::text, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: UpdateERInvestigationStepStatus :one
UPDATE hrms.er_investigation_steps
SET status = $3,
    completed_at = CASE WHEN $3 = 'completed' THEN COALESCE(sqlc.narg('completed_at')::timestamptz, NOW()) ELSE completed_at END,
    outcome_notes = COALESCE(sqlc.narg('outcome_notes')::text, outcome_notes),
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListERInvestigationSteps :many
SELECT *
FROM hrms.er_investigation_steps
WHERE tenant_id = $1 AND er_case_id = $2 AND NOT inactive
ORDER BY due_at ASC NULLS LAST, created_at ASC;

-- name: CreateERWitnessNote :one
INSERT INTO hrms.er_witness_notes (
    tenant_id, er_case_id, witness_user_id, witness_name, interview_at, interviewer_user_id, statement_summary, confidentiality_level, created_by, updated_by
) VALUES (
    $1, $2, sqlc.narg('witness_user_id')::uuid, sqlc.narg('witness_name')::text, sqlc.narg('interview_at')::timestamptz,
    sqlc.narg('interviewer_user_id')::uuid, $3, $4, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListERWitnessNotes :many
SELECT *
FROM hrms.er_witness_notes
WHERE tenant_id = $1 AND er_case_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: CreateEREvidenceAttachment :one
INSERT INTO hrms.er_evidence_attachments (
    tenant_id, er_case_id, allegation_id, file_name, content_type, storage_path, checksum_sha256, size_bytes, evidence_type, description, uploaded_by, legal_hold, created_by, updated_by
) VALUES (
    $1, $2, sqlc.narg('allegation_id')::uuid, $3, $4, $5, sqlc.narg('checksum_sha256')::text, $6, $7,
    sqlc.narg('description')::text, sqlc.narg('uploaded_by')::uuid, $8, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListEREvidenceAttachments :many
SELECT *
FROM hrms.er_evidence_attachments
WHERE tenant_id = $1 AND er_case_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: CreateERFinding :one
INSERT INTO hrms.er_findings (
    tenant_id, er_case_id, allegation_id, finding, rationale, recommended_action, decided_by, decided_at, created_by, updated_by
) VALUES (
    $1, $2, sqlc.narg('allegation_id')::uuid, $3, $4, sqlc.narg('recommended_action')::text, sqlc.narg('decided_by')::uuid, sqlc.narg('decided_at')::timestamptz, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListERFindings :many
SELECT *
FROM hrms.er_findings
WHERE tenant_id = $1 AND er_case_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: CreateERActionPlan :one
INSERT INTO hrms.er_action_plans (
    tenant_id, er_case_id, action_type, description, assigned_to_user_id, due_at, completed_at, status, follow_up_notes, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, sqlc.narg('assigned_to_user_id')::uuid, sqlc.narg('due_at')::timestamptz, sqlc.narg('completed_at')::timestamptz, $5, sqlc.narg('follow_up_notes')::text, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: UpdateERActionPlanStatus :one
UPDATE hrms.er_action_plans
SET status = $3,
    completed_at = CASE WHEN $3 = 'completed' THEN COALESCE(sqlc.narg('completed_at')::timestamptz, NOW()) ELSE completed_at END,
    follow_up_notes = COALESCE(sqlc.narg('follow_up_notes')::text, follow_up_notes),
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListERActionPlans :many
SELECT *
FROM hrms.er_action_plans
WHERE tenant_id = $1 AND er_case_id = $2 AND NOT inactive
ORDER BY due_at ASC NULLS LAST, created_at ASC;

-- name: CreateERCaseEvent :one
INSERT INTO hrms.er_case_events (
    tenant_id, er_case_id, event_type, from_status, to_status, actor_user_id, comment, metadata, created_by
) VALUES (
    $1, $2, $3, sqlc.narg('from_status')::text, sqlc.narg('to_status')::text,
    sqlc.narg('actor_user_id')::uuid, sqlc.narg('comment')::text, $4, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListERCaseEvents :many
SELECT *
FROM hrms.er_case_events
WHERE tenant_id = $1 AND er_case_id = $2
ORDER BY created_at DESC;

-- name: GetERCaseSummary :many
SELECT
    status,
    severity,
    COUNT(*)::int AS case_count,
    COUNT(*) FILTER (WHERE legal_hold)::int AS legal_hold_count,
    COUNT(*) FILTER (WHERE due_at IS NOT NULL AND due_at < NOW() AND status NOT IN ('closed','cancelled'))::int AS overdue_count
FROM hrms.er_cases
WHERE tenant_id = $1 AND NOT inactive
GROUP BY status, severity
ORDER BY status, severity;
