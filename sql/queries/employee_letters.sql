-- name: ListEmployeeLetterTemplates :many
SELECT * FROM hrms.employee_letter_templates
WHERE tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('letter_type')::text IS NULL OR letter_type = sqlc.narg('letter_type')::text)
  AND NOT inactive
ORDER BY letter_type ASC, is_default DESC, name ASC;

-- name: GetDefaultEmployeeLetterTemplate :one
SELECT * FROM hrms.employee_letter_templates
WHERE tenant_id = $1 AND letter_type = $2 AND is_default AND is_active AND NOT inactive
LIMIT 1;

-- name: GetEmployeeLetterTemplate :one
SELECT * FROM hrms.employee_letter_templates
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateEmployeeLetterTemplate :one
INSERT INTO hrms.employee_letter_templates (
    tenant_id, letter_type, name, description, subject, body_html, footer_html, locale, is_default, is_active, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateEmployeeLetterTemplate :one
UPDATE hrms.employee_letter_templates
SET letter_type = $3,
    name = $4,
    description = $5,
    subject = $6,
    body_html = $7,
    footer_html = $8,
    locale = $9,
    is_default = $10,
    is_active = $11,
    updated_at = NOW(),
    updated_by = $12
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ClearDefaultEmployeeLetterTemplates :exec
UPDATE hrms.employee_letter_templates
SET is_default = FALSE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND letter_type = $2 AND is_default AND NOT inactive;

-- name: SoftDeleteEmployeeLetterTemplate :exec
UPDATE hrms.employee_letter_templates
SET inactive = TRUE, is_active = FALSE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListEmployeeLetters :many
SELECT
    el.id,
    el.tenant_id,
    el.employee_id,
    el.user_id,
    e.employee_code,
    e.firstname AS employee_firstname,
    e.lastname AS employee_lastname,
    e.email AS employee_email,
    d.name AS department_name,
    b.branch_name AS branch_name,
    des.name AS designation_name,
    el.template_id,
    elt.name AS template_name,
    el.document_type_id,
    dt.name AS document_type_name,
    el.employee_document_id,
    el.letter_type,
    el.subject,
    el.rendered_html,
    el.status,
    el.issue_date,
    el.effective_date,
    el.end_date,
    el.pdf_path,
    el.version,
    el.is_latest,
    el.approval_requested_at,
    el.approved_at,
    el.approved_by,
    el.sent_at,
    el.revoked_at,
    el.signature_token,
    el.signature_requested_at,
    el.signature_completed_at,
    el.signer_name,
    el.signer_email,
    el.signer_ip,
    el.signer_user_agent,
    el.signature_hash,
    el.audit_certificate_url,
    el.inactive,
    el.created_at,
    el.created_by,
    el.updated_at,
    el.updated_by
FROM hrms.employee_letters el
JOIN hrms.employees e ON e.tenant_id = el.tenant_id AND e.id = el.employee_id AND NOT e.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
LEFT JOIN hrms.designations des ON des.tenant_id = e.tenant_id AND des.id = e.designation_id AND NOT des.inactive
LEFT JOIN hrms.employee_letter_templates elt ON elt.tenant_id = el.tenant_id AND elt.id = el.template_id
LEFT JOIN hrms.document_types dt ON dt.tenant_id = el.tenant_id AND dt.id = el.document_type_id
WHERE el.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('employee_id')::uuid IS NULL OR el.employee_id = sqlc.narg('employee_id')::uuid)
  AND (sqlc.narg('letter_type')::text IS NULL OR el.letter_type = sqlc.narg('letter_type')::text)
  AND (sqlc.narg('status')::text IS NULL OR el.status = sqlc.narg('status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR e.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.email ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.employee_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR el.subject ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT el.inactive
ORDER BY el.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountEmployeeLetters :one
SELECT COUNT(*)
FROM hrms.employee_letters el
JOIN hrms.employees e ON e.tenant_id = el.tenant_id AND e.id = el.employee_id AND NOT e.inactive
WHERE el.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('employee_id')::uuid IS NULL OR el.employee_id = sqlc.narg('employee_id')::uuid)
  AND (sqlc.narg('letter_type')::text IS NULL OR el.letter_type = sqlc.narg('letter_type')::text)
  AND (sqlc.narg('status')::text IS NULL OR el.status = sqlc.narg('status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR e.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.email ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.employee_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR el.subject ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT el.inactive;

-- name: GetEmployeeLetter :one
SELECT * FROM hrms.employee_letters
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetEmployeeLetterBySignatureToken :one
SELECT * FROM hrms.employee_letters
WHERE signature_token = $1 AND NOT inactive;

-- name: ClearLatestEmployeeLetters :exec
UPDATE hrms.employee_letters
SET is_latest = FALSE, updated_at = NOW(), updated_by = $4
WHERE tenant_id = $1 AND employee_id = $2 AND letter_type = $3 AND is_latest AND NOT inactive;

-- name: NextEmployeeLetterVersion :one
SELECT COALESCE(MAX(version), 0)::int + 1 AS next_version
FROM hrms.employee_letters
WHERE tenant_id = $1 AND employee_id = $2 AND letter_type = $3;

-- name: CreateEmployeeLetter :one
INSERT INTO hrms.employee_letters (
    tenant_id, employee_id, user_id, template_id, document_type_id, employee_document_id,
    letter_type, subject, rendered_html, status, issue_date, effective_date, end_date,
    pdf_path, version, is_latest, approval_requested_at, signature_token,
    signature_requested_at, signer_email, created_by
)
VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12, $13,
    $14, $15, $16, $17, $18,
    $19, $20, $21
)
RETURNING *;

-- name: UpdateEmployeeLetterPDF :one
UPDATE hrms.employee_letters
SET pdf_path = $3,
    employee_document_id = COALESCE($4, employee_document_id),
    updated_at = NOW(),
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateEmployeeLetterStatus :one
UPDATE hrms.employee_letters
SET status = $3,
    approval_requested_at = CASE WHEN $3 = 'Generated' THEN COALESCE(approval_requested_at, NOW()) ELSE approval_requested_at END,
    approved_at = CASE WHEN $3 = 'Approved' THEN COALESCE(approved_at, NOW()) ELSE approved_at END,
    approved_by = CASE WHEN $3 = 'Approved' THEN COALESCE($4, approved_by) ELSE approved_by END,
    sent_at = CASE WHEN $3 = 'Sent' THEN COALESCE(sent_at, NOW()) ELSE sent_at END,
    signature_requested_at = CASE WHEN $3 = 'Sent' AND signature_token IS NOT NULL THEN COALESCE(signature_requested_at, NOW()) ELSE signature_requested_at END,
    revoked_at = CASE WHEN $3 = 'Revoked' THEN COALESCE(revoked_at, NOW()) ELSE revoked_at END,
    updated_at = NOW(),
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SignEmployeeLetter :one
UPDATE hrms.employee_letters
SET status = 'Signed',
    signature_completed_at = NOW(),
    signer_name = $2,
    signer_email = $3,
    signer_ip = $4,
    signer_user_agent = $5,
    signature_hash = $6,
    updated_at = NOW()
WHERE signature_token = $1 AND status = 'Sent' AND NOT inactive
RETURNING *;

-- name: CreateEmployeeLetterEvent :one
INSERT INTO hrms.employee_letter_events (
    tenant_id, employee_letter_id, from_status, to_status, action, remarks, actor_email, ip_address, user_agent, metadata, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, COALESCE(sqlc.arg('metadata'), '{}'::jsonb), $10)
RETURNING *;

-- name: ListEmployeeLetterEvents :many
SELECT * FROM hrms.employee_letter_events
WHERE tenant_id = $1 AND employee_letter_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: SoftDeleteEmployeeLetter :exec
UPDATE hrms.employee_letters
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;
