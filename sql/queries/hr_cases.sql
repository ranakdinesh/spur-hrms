-- name: CreateHRCaseCategory :one
INSERT INTO hrms.hr_case_categories (
    tenant_id, code, name, description, confidentiality_default, default_owner_role, is_active, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $8
)
RETURNING *;

-- name: UpdateHRCaseCategory :one
UPDATE hrms.hr_case_categories
SET code = $3,
    name = $4,
    description = $5,
    confidentiality_default = $6,
    default_owner_role = $7,
    is_active = $8,
    updated_by = $9,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListHRCaseCategories :many
SELECT *
FROM hrms.hr_case_categories
WHERE tenant_id = $1
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')::boolean)
  AND NOT inactive
ORDER BY is_active DESC, name ASC;

-- name: GetHRCaseCategory :one
SELECT *
FROM hrms.hr_case_categories
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteHRCaseCategory :exec
UPDATE hrms.hr_case_categories
SET inactive = TRUE,
    is_active = FALSE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateHRCaseSLAPolicy :one
INSERT INTO hrms.hr_case_sla_policies (
    tenant_id, category_id, priority, response_hours, resolution_hours, escalation_hours, is_active, created_by, updated_by
) VALUES (
    $1, sqlc.narg('category_id')::uuid, $2, $3, $4, $5, $6, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: UpdateHRCaseSLAPolicy :one
UPDATE hrms.hr_case_sla_policies
SET category_id = sqlc.narg('category_id')::uuid,
    priority = $3,
    response_hours = $4,
    resolution_hours = $5,
    escalation_hours = $6,
    is_active = $7,
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ListHRCaseSLAPolicies :many
SELECT
    sp.*,
    cc.name AS category_name
FROM hrms.hr_case_sla_policies sp
LEFT JOIN hrms.hr_case_categories cc ON cc.tenant_id = sp.tenant_id AND cc.id = sp.category_id AND NOT cc.inactive
WHERE sp.tenant_id = $1
  AND (sqlc.narg('category_id')::uuid IS NULL OR sp.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('priority')::text IS NULL OR sp.priority = sqlc.narg('priority')::text)
  AND NOT sp.inactive
ORDER BY sp.is_active DESC, cc.name ASC NULLS LAST, sp.priority ASC;

-- name: ResolveHRCaseSLAPolicy :one
SELECT *
FROM hrms.hr_case_sla_policies
WHERE tenant_id = $1
  AND is_active
  AND NOT inactive
  AND (category_id = sqlc.narg('category_id')::uuid OR category_id IS NULL)
  AND priority = $2
ORDER BY CASE WHEN category_id = sqlc.narg('category_id')::uuid THEN 0 ELSE 1 END, updated_at DESC
LIMIT 1;

-- name: SoftDeleteHRCaseSLAPolicy :exec
UPDATE hrms.hr_case_sla_policies
SET inactive = TRUE,
    is_active = FALSE,
    updated_by = $3,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateHRCase :one
INSERT INTO hrms.hr_cases (
    tenant_id, case_number, category_id, case_type, title, description, confidentiality_level,
    requester_user_id, subject_employee_user_id, owner_user_id, owner_role, status, priority,
    source_channel, first_response_due_at, due_at, created_by, updated_by
) VALUES (
    $1, $2, sqlc.narg('category_id')::uuid, $3, $4, $5, $6,
    sqlc.narg('requester_user_id')::uuid, sqlc.narg('subject_employee_user_id')::uuid,
    sqlc.narg('owner_user_id')::uuid, sqlc.narg('owner_role')::text, $7, $8,
    $9, sqlc.narg('first_response_due_at')::timestamptz, sqlc.narg('due_at')::timestamptz,
    sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListHRCases :many
SELECT
    c.*,
    cc.name AS category_name,
    cc.code AS category_code,
    NULL::text AS requester_email,
    NULL::text AS subject_email,
    NULL::text AS owner_email,
    COALESCE(comment_counts.comment_count, 0)::int AS comment_count,
    COALESCE(attachment_counts.attachment_count, 0)::int AS attachment_count
FROM hrms.hr_cases c
LEFT JOIN hrms.hr_case_categories cc ON cc.tenant_id = c.tenant_id AND cc.id = c.category_id AND NOT cc.inactive
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS comment_count
    FROM hrms.hr_case_comments hcc
    WHERE hcc.tenant_id = c.tenant_id AND hcc.case_id = c.id AND NOT hcc.inactive
) comment_counts ON TRUE
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS attachment_count
    FROM hrms.hr_case_attachments hca
    WHERE hca.tenant_id = c.tenant_id AND hca.case_id = c.id AND NOT hca.inactive
) attachment_counts ON TRUE
WHERE c.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('status')::text IS NULL OR c.status = sqlc.narg('status')::text)
  AND (sqlc.narg('priority')::text IS NULL OR c.priority = sqlc.narg('priority')::text)
  AND (sqlc.narg('category_id')::uuid IS NULL OR c.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('requester_user_id')::uuid IS NULL OR c.requester_user_id = sqlc.narg('requester_user_id')::uuid)
  AND (sqlc.narg('subject_employee_user_id')::uuid IS NULL OR c.subject_employee_user_id = sqlc.narg('subject_employee_user_id')::uuid)
  AND (sqlc.narg('owner_user_id')::uuid IS NULL OR c.owner_user_id = sqlc.narg('owner_user_id')::uuid)
  AND (sqlc.narg('confidentiality_level')::text IS NULL OR c.confidentiality_level = sqlc.narg('confidentiality_level')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.case_number ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.description ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.requester_user_id::text ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.subject_employee_user_id::text ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT c.inactive
ORDER BY
  CASE c.priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'normal' THEN 3 ELSE 4 END,
  CASE c.status WHEN 'new' THEN 1 WHEN 'open' THEN 2 WHEN 'escalated' THEN 3 WHEN 'in_progress' THEN 4 ELSE 5 END,
  c.due_at ASC NULLS LAST,
  c.last_activity_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountHRCases :one
SELECT COUNT(*)
FROM hrms.hr_cases c
WHERE c.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('status')::text IS NULL OR c.status = sqlc.narg('status')::text)
  AND (sqlc.narg('priority')::text IS NULL OR c.priority = sqlc.narg('priority')::text)
  AND (sqlc.narg('category_id')::uuid IS NULL OR c.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('requester_user_id')::uuid IS NULL OR c.requester_user_id = sqlc.narg('requester_user_id')::uuid)
  AND (sqlc.narg('subject_employee_user_id')::uuid IS NULL OR c.subject_employee_user_id = sqlc.narg('subject_employee_user_id')::uuid)
  AND (sqlc.narg('owner_user_id')::uuid IS NULL OR c.owner_user_id = sqlc.narg('owner_user_id')::uuid)
  AND (sqlc.narg('confidentiality_level')::text IS NULL OR c.confidentiality_level = sqlc.narg('confidentiality_level')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.case_number ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.description ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.requester_user_id::text ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.subject_employee_user_id::text ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT c.inactive;

-- name: GetHRCase :one
SELECT
    c.*,
    cc.name AS category_name,
    cc.code AS category_code,
    NULL::text AS requester_email,
    NULL::text AS subject_email,
    NULL::text AS owner_email
FROM hrms.hr_cases c
LEFT JOIN hrms.hr_case_categories cc ON cc.tenant_id = c.tenant_id AND cc.id = c.category_id AND NOT cc.inactive
WHERE c.tenant_id = $1 AND c.id = $2 AND NOT c.inactive;

-- name: UpdateHRCaseAssignment :one
UPDATE hrms.hr_cases
SET owner_user_id = sqlc.narg('owner_user_id')::uuid,
    owner_role = sqlc.narg('owner_role')::text,
    last_activity_at = NOW(),
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateHRCaseStatus :one
UPDATE hrms.hr_cases
SET status = $3,
    first_responded_at = CASE WHEN first_responded_at IS NULL AND $3 IN ('open','in_progress','waiting_on_employee','waiting_on_hr','escalated','resolved','closed') THEN NOW() ELSE first_responded_at END,
    resolved_at = CASE WHEN $3 = 'resolved' THEN NOW() ELSE resolved_at END,
    closed_at = CASE WHEN $3 = 'closed' THEN NOW() ELSE closed_at END,
    escalated_at = CASE WHEN $3 = 'escalated' THEN NOW() ELSE escalated_at END,
    escalation_level = CASE WHEN $3 = 'escalated' THEN escalation_level + 1 ELSE escalation_level END,
    resolution_summary = CASE WHEN $3 IN ('resolved','closed') THEN sqlc.narg('resolution_summary')::text ELSE resolution_summary END,
    last_activity_at = NOW(),
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateHRCaseDetails :one
UPDATE hrms.hr_cases
SET category_id = sqlc.narg('category_id')::uuid,
    case_type = $3,
    title = $4,
    description = $5,
    confidentiality_level = $6,
    priority = $7,
    subject_employee_user_id = sqlc.narg('subject_employee_user_id')::uuid,
    due_at = sqlc.narg('due_at')::timestamptz,
    last_activity_at = NOW(),
    updated_by = sqlc.narg('actor_id')::uuid,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateHRCaseComment :one
INSERT INTO hrms.hr_case_comments (
    tenant_id, case_id, author_user_id, visibility, body, created_by, updated_by
) VALUES (
    $1, $2, sqlc.narg('author_user_id')::uuid, $3, $4, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListHRCaseComments :many
SELECT *
FROM hrms.hr_case_comments
WHERE tenant_id = $1
  AND case_id = $2
  AND (sqlc.narg('include_internal')::boolean OR visibility = 'public')
  AND NOT inactive
ORDER BY created_at ASC;

-- name: CreateHRCaseAttachment :one
INSERT INTO hrms.hr_case_attachments (
    tenant_id, case_id, comment_id, file_name, content_type, object_key, visibility, uploaded_by, created_by, updated_by
) VALUES (
    $1, $2, sqlc.narg('comment_id')::uuid, $3, $4, $5, $6, sqlc.narg('uploaded_by')::uuid, sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListHRCaseAttachments :many
SELECT *
FROM hrms.hr_case_attachments
WHERE tenant_id = $1
  AND case_id = $2
  AND (sqlc.narg('include_internal')::boolean OR visibility = 'public')
  AND NOT inactive
ORDER BY created_at DESC;

-- name: CreateHRCaseEvent :one
INSERT INTO hrms.hr_case_events (
    tenant_id, case_id, event_type, from_status, to_status, actor_user_id, comment, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, sqlc.narg('from_status')::text, sqlc.narg('to_status')::text,
    sqlc.narg('actor_user_id')::uuid, sqlc.narg('comment')::text, $4,
    sqlc.narg('actor_id')::uuid, sqlc.narg('actor_id')::uuid
)
RETURNING *;

-- name: ListHRCaseEvents :many
SELECT *
FROM hrms.hr_case_events
WHERE tenant_id = $1 AND case_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetHRCaseSummary :many
SELECT
    COALESCE(status, 'all') AS status,
    COALESCE(priority, 'all') AS priority,
    COUNT(*)::int AS case_count,
    COUNT(*) FILTER (WHERE due_at IS NOT NULL AND due_at < NOW() AND status NOT IN ('resolved','closed','cancelled'))::int AS overdue_count,
    COUNT(*) FILTER (WHERE status = 'escalated')::int AS escalated_count,
    COUNT(*) FILTER (WHERE confidentiality_level IN ('sensitive','grievance'))::int AS restricted_count
FROM hrms.hr_cases
WHERE tenant_id = $1 AND NOT inactive
GROUP BY GROUPING SETS ((status, priority), ())
ORDER BY status ASC, priority ASC;
