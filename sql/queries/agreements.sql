-- name: CreateAgreementTemplate :one
INSERT INTO hrms.agreement_templates (
    tenant_id, agreement_type, name, description, subject, body_html, footer_html,
    locale, is_default, is_active, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: ListAgreementTemplates :many
SELECT * FROM hrms.agreement_templates
WHERE tenant_id = $1
  AND (sqlc.narg('agreement_type')::text IS NULL OR agreement_type = sqlc.narg('agreement_type')::text)
  AND NOT inactive
ORDER BY agreement_type, is_default DESC, name;

-- name: GetDefaultAgreementTemplate :one
SELECT * FROM hrms.agreement_templates
WHERE tenant_id = $1 AND agreement_type = $2 AND is_default AND is_active AND NOT inactive
ORDER BY updated_at DESC
LIMIT 1;

-- name: GetAgreementTemplate :one
SELECT * FROM hrms.agreement_templates
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateAgreementTemplate :one
UPDATE hrms.agreement_templates
SET agreement_type = $3,
    name = $4,
    description = $5,
    subject = $6,
    body_html = $7,
    footer_html = $8,
    locale = $9,
    is_default = $10,
    is_active = $11,
    metadata = $12,
    updated_by = $13
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteAgreementTemplate :exec
UPDATE hrms.agreement_templates
SET inactive = TRUE,
    is_active = FALSE,
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateAgreement :one
INSERT INTO hrms.agreements (
    tenant_id, agreement_type, title, template_id, worker_profile_id, engagement_id, project_id,
    subject, rendered_html, status, issue_date, effective_date, end_date, pdf_path, signature_token,
    signature_requested_at, signer_email, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $19
)
RETURNING *;

-- name: ListAgreements :many
SELECT
    a.id,
    a.tenant_id,
    a.agreement_type,
    a.title,
    a.template_id,
    t.name AS template_name,
    a.worker_profile_id,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    a.engagement_id,
    e.title AS engagement_title,
    e.engagement_code,
    a.project_id,
    p.name AS project_name,
    p.project_code,
    a.subject,
    a.rendered_html,
    a.status,
    a.issue_date,
    a.effective_date,
    a.end_date,
    a.pdf_path,
    a.version,
    a.is_latest,
    a.sent_at,
    a.revoked_at,
    a.signature_token,
    a.signature_requested_at,
    a.signature_completed_at,
    a.signer_name,
    a.signer_email,
    a.signer_ip,
    a.signer_user_agent,
    a.signature_hash,
    a.audit_certificate_url,
    a.metadata,
    a.inactive,
    a.created_at,
    a.created_by,
    a.updated_at,
    a.updated_by
FROM hrms.agreements a
LEFT JOIN hrms.agreement_templates t ON t.tenant_id = a.tenant_id AND t.id = a.template_id AND NOT t.inactive
LEFT JOIN hrms.worker_profiles wp ON wp.tenant_id = a.tenant_id AND wp.id = a.worker_profile_id AND NOT wp.inactive
LEFT JOIN hrms.engagements e ON e.tenant_id = a.tenant_id AND e.id = a.engagement_id AND NOT e.inactive
LEFT JOIN hrms.projects p ON p.tenant_id = a.tenant_id AND p.id = a.project_id AND NOT p.inactive
WHERE a.tenant_id = $1
  AND (sqlc.narg('agreement_type')::text IS NULL OR a.agreement_type = sqlc.narg('agreement_type')::text)
  AND (sqlc.narg('status')::text IS NULL OR a.status = sqlc.narg('status')::text)
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR a.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('engagement_id')::uuid IS NULL OR a.engagement_id = sqlc.narg('engagement_id')::uuid)
  AND (sqlc.narg('project_id')::uuid IS NULL OR a.project_id = sqlc.narg('project_id')::uuid)
  AND (
      sqlc.narg('search')::text IS NULL
      OR a.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR a.subject ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR e.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR p.name ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT a.inactive
ORDER BY
    CASE a.status WHEN 'Generated' THEN 0 WHEN 'Approved' THEN 1 WHEN 'Sent' THEN 2 WHEN 'Signed' THEN 3 ELSE 4 END,
    a.updated_at DESC;

-- name: GetAgreement :one
SELECT * FROM hrms.agreements
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetAgreementBySignatureToken :one
SELECT * FROM hrms.agreements
WHERE signature_token = $1 AND NOT inactive;

-- name: UpdateAgreementPDF :one
UPDATE hrms.agreements
SET pdf_path = $3,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateAgreementStatus :one
UPDATE hrms.agreements
SET status = $3,
    sent_at = CASE WHEN $3 = 'Sent' THEN now() ELSE sent_at END,
    revoked_at = CASE WHEN $3 = 'Revoked' THEN now() ELSE revoked_at END,
    signature_requested_at = CASE WHEN $3 = 'Sent' AND signature_token IS NOT NULL THEN now() ELSE signature_requested_at END,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SignAgreement :one
UPDATE hrms.agreements
SET status = 'Signed',
    signature_completed_at = now(),
    signer_name = $2,
    signer_email = $3,
    signer_ip = $4,
    signer_user_agent = $5,
    signature_hash = $6
WHERE signature_token = $1 AND NOT inactive
RETURNING *;

-- name: SoftDeleteAgreement :exec
UPDATE hrms.agreements
SET inactive = TRUE,
    status = 'Revoked',
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateAgreementEvent :one
INSERT INTO hrms.agreement_events (
    tenant_id, agreement_id, from_status, to_status, action, remarks, actor_email,
    ip_address, user_agent, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11
)
RETURNING *;

-- name: ListAgreementEvents :many
SELECT * FROM hrms.agreement_events
WHERE tenant_id = $1 AND agreement_id = $2 AND NOT inactive
ORDER BY created_at DESC;
