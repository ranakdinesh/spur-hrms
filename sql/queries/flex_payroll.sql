-- name: CreateFlexPayRun :one
INSERT INTO hrms.flex_pay_runs (
    tenant_id, run_code, title, run_type, status, period_start, period_end, payout_date,
    currency_code, source_policy, notes, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $13
)
RETURNING *;

-- name: GetFlexPayRun :one
SELECT * FROM hrms.flex_pay_runs
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListFlexPayRuns :many
SELECT * FROM hrms.flex_pay_runs
WHERE tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('run_type')::text IS NULL OR run_type = sqlc.narg('run_type')::text)
  AND (sqlc.narg('date_from')::date IS NULL OR period_end >= sqlc.narg('date_from')::date)
  AND (sqlc.narg('date_to')::date IS NULL OR period_start <= sqlc.narg('date_to')::date)
  AND (
      sqlc.narg('search')::text IS NULL
      OR run_code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR title ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT inactive
ORDER BY period_start DESC, updated_at DESC;

-- name: UpdateFlexPayRunStatus :one
UPDATE hrms.flex_pay_runs
SET status = $3,
    invoice_count = $4,
    item_count = $5,
    gross_amount = $6,
    tds_amount = $7,
    gst_amount = $8,
    net_amount = $9,
    generated_at = CASE WHEN $3 = 'generated' THEN COALESCE(generated_at, now()) ELSE generated_at END,
    submitted_at = CASE WHEN $3 = 'submitted' THEN now() ELSE submitted_at END,
    submitted_by = CASE WHEN $3 = 'submitted' THEN $10 ELSE submitted_by END,
    approved_at = CASE WHEN $3 = 'approved' THEN now() ELSE approved_at END,
    approved_by = CASE WHEN $3 = 'approved' THEN $10 ELSE approved_by END,
    rejected_at = CASE WHEN $3 = 'rejected' THEN now() ELSE rejected_at END,
    rejected_by = CASE WHEN $3 = 'rejected' THEN $10 ELSE rejected_by END,
    paid_at = CASE WHEN $3 = 'paid' THEN now() ELSE paid_at END,
    paid_by = CASE WHEN $3 = 'paid' THEN $10 ELSE paid_by END,
    payment_reference = $11,
    export_batch_ref = $12,
    notes = $13,
    metadata = $14,
    updated_by = $10
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteFlexPayRun :exec
UPDATE hrms.flex_pay_runs
SET inactive = TRUE,
    status = 'cancelled',
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateContractorInvoice :one
INSERT INTO hrms.contractor_invoices (
    tenant_id, flex_pay_run_id, worker_profile_id, engagement_id, invoice_number, invoice_date,
    due_date, status, vendor_name, vendor_gstin, place_of_supply, reverse_charge,
    currency_code, gross_amount, tds_section, tds_rate, tds_amount, gst_rate, gst_amount,
    net_amount, attachment_path, notes, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12,
    $13, $14, $15, $16, $17, $18, $19,
    $20, $21, $22, $23, $24, $24
)
RETURNING *;

-- name: UpdateContractorInvoice :one
UPDATE hrms.contractor_invoices
SET flex_pay_run_id = $3,
    worker_profile_id = $4,
    engagement_id = $5,
    invoice_number = $6,
    invoice_date = $7,
    due_date = $8,
    vendor_name = $9,
    vendor_gstin = $10,
    place_of_supply = $11,
    reverse_charge = $12,
    currency_code = $13,
    gross_amount = $14,
    tds_section = $15,
    tds_rate = $16,
    tds_amount = $17,
    gst_rate = $18,
    gst_amount = $19,
    net_amount = $20,
    attachment_path = $21,
    notes = $22,
    metadata = $23,
    updated_by = $24
WHERE tenant_id = $1 AND id = $2 AND status IN ('draft', 'rejected') AND NOT inactive
RETURNING *;

-- name: GetContractorInvoice :one
SELECT * FROM hrms.contractor_invoices
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListContractorInvoices :many
SELECT
    ci.*,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    e.title AS engagement_title,
    e.engagement_code
FROM hrms.contractor_invoices ci
JOIN hrms.worker_profiles wp ON wp.tenant_id = ci.tenant_id AND wp.id = ci.worker_profile_id AND NOT wp.inactive
LEFT JOIN hrms.engagements e ON e.tenant_id = ci.tenant_id AND e.id = ci.engagement_id AND NOT e.inactive
WHERE ci.tenant_id = $1
  AND (sqlc.narg('flex_pay_run_id')::uuid IS NULL OR ci.flex_pay_run_id = sqlc.narg('flex_pay_run_id')::uuid)
  AND (sqlc.narg('worker_profile_id')::uuid IS NULL OR ci.worker_profile_id = sqlc.narg('worker_profile_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR ci.status = sqlc.narg('status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR ci.invoice_number ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ci.vendor_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.display_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR wp.worker_code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT ci.inactive
ORDER BY ci.invoice_date DESC, ci.updated_at DESC;

-- name: UpdateContractorInvoiceStatus :one
UPDATE hrms.contractor_invoices
SET status = $3,
    submitted_at = CASE WHEN $3 = 'submitted' THEN now() ELSE submitted_at END,
    submitted_by = CASE WHEN $3 = 'submitted' THEN $4 ELSE submitted_by END,
    approved_at = CASE WHEN $3 = 'approved' THEN now() ELSE approved_at END,
    approved_by = CASE WHEN $3 = 'approved' THEN $4 ELSE approved_by END,
    rejected_at = CASE WHEN $3 = 'rejected' THEN now() ELSE rejected_at END,
    rejected_by = CASE WHEN $3 = 'rejected' THEN $4 ELSE rejected_by END,
    rejection_reason = CASE WHEN $3 = 'rejected' THEN $5 ELSE rejection_reason END,
    paid_at = CASE WHEN $3 = 'paid' THEN now() ELSE paid_at END,
    paid_by = CASE WHEN $3 = 'paid' THEN $4 ELSE paid_by END,
    payment_reference = CASE WHEN $3 = 'paid' THEN $6 ELSE payment_reference END,
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteContractorInvoice :exec
UPDATE hrms.contractor_invoices
SET inactive = TRUE,
    status = 'cancelled',
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND status IN ('draft', 'rejected') AND NOT inactive;

-- name: CreateFlexPayRunItem :one
INSERT INTO hrms.flex_pay_run_items (
    tenant_id, flex_pay_run_id, contractor_invoice_id, worker_profile_id, engagement_id,
    source_type, source_id, description, quantity, rate_amount, gross_amount,
    tds_section, tds_rate, tds_amount, gst_rate, gst_amount, net_amount,
    status, metadata, created_by, updated_by
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10, $11,
    $12, $13, $14, $15, $16, $17,
    $18, $19, $20, $20
)
RETURNING *;

-- name: ListFlexPayRunItems :many
SELECT
    item.*,
    wp.display_name AS worker_display_name,
    wp.worker_code,
    e.title AS engagement_title,
    e.engagement_code,
    ci.invoice_number
FROM hrms.flex_pay_run_items item
JOIN hrms.worker_profiles wp ON wp.tenant_id = item.tenant_id AND wp.id = item.worker_profile_id AND NOT wp.inactive
LEFT JOIN hrms.engagements e ON e.tenant_id = item.tenant_id AND e.id = item.engagement_id AND NOT e.inactive
LEFT JOIN hrms.contractor_invoices ci ON ci.tenant_id = item.tenant_id AND ci.id = item.contractor_invoice_id AND NOT ci.inactive
WHERE item.tenant_id = $1
  AND item.flex_pay_run_id = $2
  AND NOT item.inactive
ORDER BY wp.display_name ASC, item.created_at ASC;

-- name: UpdateFlexPayRunItemsStatusByRun :exec
UPDATE hrms.flex_pay_run_items
SET status = $3,
    updated_by = $4
WHERE tenant_id = $1 AND flex_pay_run_id = $2 AND NOT inactive;

-- name: GetFlexPayRunTotals :one
SELECT
    COUNT(DISTINCT contractor_invoice_id)::int AS invoice_count,
    COUNT(*)::int AS item_count,
    COALESCE(SUM(gross_amount), 0)::numeric AS gross_amount,
    COALESCE(SUM(tds_amount), 0)::numeric AS tds_amount,
    COALESCE(SUM(gst_amount), 0)::numeric AS gst_amount,
    COALESCE(SUM(net_amount), 0)::numeric AS net_amount
FROM hrms.flex_pay_run_items
WHERE tenant_id = $1 AND flex_pay_run_id = $2 AND NOT inactive;

-- name: CreateFlexPayRunEvent :one
INSERT INTO hrms.flex_pay_run_events (
    tenant_id, flex_pay_run_id, contractor_invoice_id, event_type, from_status, to_status,
    comment, actor_id, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9
)
RETURNING *;

-- name: ListFlexPayRunEvents :many
SELECT * FROM hrms.flex_pay_run_events
WHERE tenant_id = $1
  AND (sqlc.narg('flex_pay_run_id')::uuid IS NULL OR flex_pay_run_id = sqlc.narg('flex_pay_run_id')::uuid)
  AND (sqlc.narg('contractor_invoice_id')::uuid IS NULL OR contractor_invoice_id = sqlc.narg('contractor_invoice_id')::uuid)
ORDER BY created_at DESC;

-- name: ListApprovedWorkLogPaymentCandidates :many
SELECT
    wl.id AS work_log_id,
    wl.tenant_id,
    wl.engagement_id,
    wl.worker_profile_id,
    wl.log_date,
    COALESCE(wl.billable_hours, wl.hours_worked)::numeric AS billable_hours,
    e.rate_amount,
    e.rate_unit,
    e.currency_code,
    e.title AS engagement_title,
    e.engagement_code,
    wp.display_name AS worker_display_name,
    wp.worker_code
FROM hrms.work_logs wl
JOIN hrms.engagements e ON e.tenant_id = wl.tenant_id AND e.id = wl.engagement_id AND NOT e.inactive
JOIN hrms.worker_profiles wp ON wp.tenant_id = wl.tenant_id AND wp.id = wl.worker_profile_id AND NOT wp.inactive
WHERE wl.tenant_id = $1
  AND wl.status = 'approved'
  AND wl.log_date BETWEEN $2 AND $3
  AND e.rate_unit IN ('hour', 'day')
  AND e.rate_amount IS NOT NULL
  AND NOT wl.inactive
  AND NOT EXISTS (
      SELECT 1 FROM hrms.flex_pay_run_items item
      WHERE item.tenant_id = wl.tenant_id
        AND item.source_type = 'work_log'
        AND item.source_id = wl.id
        AND NOT item.inactive
  )
ORDER BY wp.display_name ASC, wl.log_date ASC;

-- name: ListAcceptedMilestonePaymentCandidates :many
SELECT
    pm.id AS milestone_id,
    pm.tenant_id,
    pm.project_id,
    pm.engagement_id,
    e.worker_profile_id,
    pm.title,
    pm.milestone_code,
    pm.amount,
    pm.currency_code,
    pm.accepted_at,
    p.name AS project_name,
    p.project_code,
    wp.display_name AS worker_display_name,
    wp.worker_code
FROM hrms.project_milestones pm
JOIN hrms.projects p ON p.tenant_id = pm.tenant_id AND p.id = pm.project_id AND NOT p.inactive
JOIN hrms.engagements e ON e.tenant_id = pm.tenant_id AND e.id = pm.engagement_id AND NOT e.inactive
JOIN hrms.worker_profiles wp ON wp.tenant_id = pm.tenant_id AND wp.id = e.worker_profile_id AND NOT wp.inactive
WHERE pm.tenant_id = $1
  AND pm.status = 'accepted'
  AND pm.amount IS NOT NULL
  AND pm.accepted_at::date BETWEEN $2 AND $3
  AND NOT pm.inactive
  AND NOT EXISTS (
      SELECT 1 FROM hrms.flex_pay_run_items item
      WHERE item.tenant_id = pm.tenant_id
        AND item.source_type = 'milestone'
        AND item.source_id = pm.id
        AND NOT item.inactive
  )
ORDER BY wp.display_name ASC, pm.accepted_at ASC;
