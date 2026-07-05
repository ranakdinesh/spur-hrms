-- name: CreateBenefitPlan :one
INSERT INTO hrms.benefit_plans (
  tenant_id, code, name, plan_type, description, provider_name, policy_number, coverage_amount,
  employer_contribution, employee_contribution, currency_code, eligibility_rule, insurance_metadata,
  effective_from, effective_to, is_active, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17
) RETURNING *;

-- name: UpdateBenefitPlan :one
UPDATE hrms.benefit_plans
SET code=$3, name=$4, plan_type=$5, description=$6, provider_name=$7, policy_number=$8,
    coverage_amount=$9, employer_contribution=$10, employee_contribution=$11, currency_code=$12,
    eligibility_rule=$13, insurance_metadata=$14, effective_from=$15, effective_to=$16,
    is_active=$17, updated_by=$18
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListBenefitPlans :many
SELECT * FROM hrms.benefit_plans
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('active_only')::boolean IS NULL OR is_active = sqlc.narg('active_only')::boolean)
  AND (sqlc.narg('plan_type')::text IS NULL OR plan_type = sqlc.narg('plan_type')::text)
  AND (sqlc.narg('search')::text IS NULL OR code ILIKE '%' || sqlc.narg('search')::text || '%' OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY is_active DESC, name ASC
LIMIT $2 OFFSET $3;

-- name: GetBenefitPlan :one
SELECT * FROM hrms.benefit_plans WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: SoftDeleteBenefitPlan :exec
UPDATE hrms.benefit_plans
SET inactive=TRUE, is_active=FALSE, updated_by=$3
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateBenefitEnrollmentWindow :one
INSERT INTO hrms.benefit_enrollment_windows (
  tenant_id, plan_id, name, opens_on, closes_on, effective_from, effective_to, status, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
) RETURNING *;

-- name: UpdateBenefitEnrollmentWindow :one
UPDATE hrms.benefit_enrollment_windows
SET plan_id=$3, name=$4, opens_on=$5, closes_on=$6, effective_from=$7, effective_to=$8,
    status=$9, metadata=$10, updated_by=$11
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListBenefitEnrollmentWindows :many
SELECT * FROM hrms.benefit_enrollment_windows
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('plan_id')::uuid IS NULL OR plan_id = sqlc.narg('plan_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
ORDER BY opens_on DESC, closes_on DESC
LIMIT $2 OFFSET $3;

-- name: GetBenefitEnrollmentWindow :one
SELECT * FROM hrms.benefit_enrollment_windows WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: SoftDeleteBenefitEnrollmentWindow :exec
UPDATE hrms.benefit_enrollment_windows
SET inactive=TRUE, status='archived', updated_by=$3
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateBenefitDependent :one
INSERT INTO hrms.benefit_dependents (
  tenant_id, employee_user_id, full_name, relationship, date_of_birth, gender, nominee_percentage,
  is_nominee, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
) RETURNING *;

-- name: UpdateBenefitDependent :one
UPDATE hrms.benefit_dependents
SET employee_user_id=$3, full_name=$4, relationship=$5, date_of_birth=$6, gender=$7,
    nominee_percentage=$8, is_nominee=$9, metadata=$10, updated_by=$11
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListBenefitDependents :many
SELECT * FROM hrms.benefit_dependents
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('employee_user_id')::uuid IS NULL OR employee_user_id = sqlc.narg('employee_user_id')::uuid)
  AND (sqlc.narg('nominees_only')::boolean IS NULL OR is_nominee = sqlc.narg('nominees_only')::boolean)
ORDER BY full_name ASC
LIMIT $2 OFFSET $3;

-- name: GetBenefitDependent :one
SELECT * FROM hrms.benefit_dependents WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: SoftDeleteBenefitDependent :exec
UPDATE hrms.benefit_dependents
SET inactive=TRUE, updated_by=$3
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateBenefitEnrollment :one
INSERT INTO hrms.benefit_enrollments (
  tenant_id, plan_id, window_id, employee_user_id, status, coverage_level, selected_amount,
  employee_contribution, employer_contribution, effective_from, effective_to, submitted_at,
  review_remarks, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15
) RETURNING *;

-- name: UpdateBenefitEnrollment :one
UPDATE hrms.benefit_enrollments
SET plan_id=$3, window_id=$4, employee_user_id=$5, status=$6, coverage_level=$7, selected_amount=$8,
    employee_contribution=$9, employer_contribution=$10, effective_from=$11, effective_to=$12,
    submitted_at=$13, review_remarks=$14, metadata=$15, updated_by=$16
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: UpdateBenefitEnrollmentStatus :one
UPDATE hrms.benefit_enrollments
SET status=$3, reviewed_by=$4, reviewed_at=now(), review_remarks=$5, updated_by=$4
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListBenefitEnrollments :many
SELECT * FROM hrms.benefit_enrollments
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('employee_user_id')::uuid IS NULL OR employee_user_id = sqlc.narg('employee_user_id')::uuid)
  AND (sqlc.narg('plan_id')::uuid IS NULL OR plan_id = sqlc.narg('plan_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetBenefitEnrollment :one
SELECT * FROM hrms.benefit_enrollments WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: SoftDeleteBenefitEnrollment :exec
UPDATE hrms.benefit_enrollments
SET inactive=TRUE, status='cancelled', updated_by=$3
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateBenefitClaimType :one
INSERT INTO hrms.benefit_claim_types (
  tenant_id, plan_id, code, name, description, annual_limit, per_claim_limit, requires_attachment,
  taxable, payroll_component_code, eligibility_rule, is_active, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
) RETURNING *;

-- name: UpdateBenefitClaimType :one
UPDATE hrms.benefit_claim_types
SET plan_id=$3, code=$4, name=$5, description=$6, annual_limit=$7, per_claim_limit=$8,
    requires_attachment=$9, taxable=$10, payroll_component_code=$11, eligibility_rule=$12,
    is_active=$13, updated_by=$14
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListBenefitClaimTypes :many
SELECT * FROM hrms.benefit_claim_types
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('active_only')::boolean IS NULL OR is_active = sqlc.narg('active_only')::boolean)
  AND (sqlc.narg('plan_id')::uuid IS NULL OR plan_id = sqlc.narg('plan_id')::uuid)
  AND (sqlc.narg('search')::text IS NULL OR code ILIKE '%' || sqlc.narg('search')::text || '%' OR name ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY is_active DESC, name ASC
LIMIT $2 OFFSET $3;

-- name: GetBenefitClaimType :one
SELECT * FROM hrms.benefit_claim_types WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: SoftDeleteBenefitClaimType :exec
UPDATE hrms.benefit_claim_types
SET inactive=TRUE, is_active=FALSE, updated_by=$3
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateBenefitClaim :one
INSERT INTO hrms.benefit_claims (
  tenant_id, claim_number, claim_type_id, plan_id, employee_user_id, dependent_id, expense_date,
  submitted_at, claim_amount, approved_amount, currency_code, status, payment_status, payment_reference,
  paid_at, reviewed_by, reviewed_at, review_remarks, payroll_export_status, payroll_exported_at,
  payroll_export_reference, notes, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24
) RETURNING *;

-- name: UpdateBenefitClaim :one
UPDATE hrms.benefit_claims
SET claim_type_id=$3, plan_id=$4, employee_user_id=$5, dependent_id=$6, expense_date=$7,
    submitted_at=$8, claim_amount=$9, approved_amount=$10, currency_code=$11, status=$12,
    payment_status=$13, payment_reference=$14, paid_at=$15, reviewed_by=$16, reviewed_at=$17,
    review_remarks=$18, payroll_export_status=$19, payroll_exported_at=$20,
    payroll_export_reference=$21, notes=$22, metadata=$23, updated_by=$24
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: UpdateBenefitClaimStatus :one
UPDATE hrms.benefit_claims
SET status=$3, approved_amount=$4, payment_status=$5, payment_reference=$6, paid_at=$7,
    reviewed_by=$8, reviewed_at=CASE WHEN $8::uuid IS NULL THEN reviewed_at ELSE now() END,
    review_remarks=$9, payroll_export_status=$10, payroll_exported_at=$11,
    payroll_export_reference=$12, updated_by=$8
WHERE tenant_id=$1 AND id=$2 AND NOT inactive
RETURNING *;

-- name: ListBenefitClaims :many
SELECT * FROM hrms.benefit_claims
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('employee_user_id')::uuid IS NULL OR employee_user_id = sqlc.narg('employee_user_id')::uuid)
  AND (sqlc.narg('claim_type_id')::uuid IS NULL OR claim_type_id = sqlc.narg('claim_type_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status')::text)
  AND (sqlc.narg('payment_status')::text IS NULL OR payment_status = sqlc.narg('payment_status')::text)
  AND (sqlc.narg('payroll_export_status')::text IS NULL OR payroll_export_status = sqlc.narg('payroll_export_status')::text)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetBenefitClaim :one
SELECT * FROM hrms.benefit_claims WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: SoftDeleteBenefitClaim :exec
UPDATE hrms.benefit_claims
SET inactive=TRUE, status='cancelled', payroll_export_status='blocked', updated_by=$3
WHERE tenant_id=$1 AND id=$2 AND NOT inactive;

-- name: CreateBenefitClaimAttachment :one
INSERT INTO hrms.benefit_claim_attachments (
  tenant_id, claim_id, file_name, content_type, storage_path, checksum_sha256, size_bytes,
  uploaded_by, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
) RETURNING *;

-- name: ListBenefitClaimAttachments :many
SELECT * FROM hrms.benefit_claim_attachments
WHERE tenant_id=$1 AND claim_id=$2 AND NOT inactive
ORDER BY created_at DESC;

-- name: CreateBenefitEvent :one
INSERT INTO hrms.benefit_events (
  tenant_id, source_type, source_id, action, from_status, to_status, actor_user_id, remarks, metadata, created_by
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
) RETURNING *;

-- name: ListBenefitEvents :many
SELECT * FROM hrms.benefit_events
WHERE tenant_id=$1
  AND NOT inactive
  AND (sqlc.narg('source_type')::text IS NULL OR source_type = sqlc.narg('source_type')::text)
  AND (sqlc.narg('source_id')::uuid IS NULL OR source_id = sqlc.narg('source_id')::uuid)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetBenefitClaimLimitUsage :one
SELECT COALESCE(SUM(COALESCE(approved_amount, claim_amount)), 0)::numeric AS used_amount
FROM hrms.benefit_claims
WHERE tenant_id=$1
  AND employee_user_id=$2
  AND claim_type_id=$3
  AND NOT inactive
  AND status IN ('submitted','under_review','approved','paid')
  AND expense_date >= $4
  AND expense_date <= $5;

-- name: GetBenefitsSummary :many
SELECT 'active_plans'::text AS metric, COUNT(*)::int AS metric_count, 0::numeric AS amount
FROM hrms.benefit_plans p
WHERE p.tenant_id=sqlc.arg('summary_tenant_id')::uuid AND NOT p.inactive AND p.is_active
UNION ALL
SELECT 'open_windows', COUNT(*)::int, 0::numeric
FROM hrms.benefit_enrollment_windows w
WHERE w.tenant_id=sqlc.arg('summary_tenant_id')::uuid AND NOT w.inactive AND w.status='open' AND w.closes_on >= CURRENT_DATE
UNION ALL
SELECT 'pending_enrollments', COUNT(*)::int, 0::numeric
FROM hrms.benefit_enrollments e
WHERE e.tenant_id=sqlc.arg('summary_tenant_id')::uuid AND NOT e.inactive AND e.status='submitted'
UNION ALL
SELECT 'pending_claims', COUNT(*)::int, COALESCE(SUM(claim_amount), 0)::numeric
FROM hrms.benefit_claims c
WHERE c.tenant_id=sqlc.arg('summary_tenant_id')::uuid AND NOT c.inactive AND c.status IN ('submitted','under_review')
UNION ALL
SELECT 'payable_claims', COUNT(*)::int, COALESCE(SUM(COALESCE(approved_amount, claim_amount)), 0)::numeric
FROM hrms.benefit_claims c
WHERE c.tenant_id=sqlc.arg('summary_tenant_id')::uuid AND NOT c.inactive AND c.payment_status='pending'
UNION ALL
SELECT 'payroll_ready', COUNT(*)::int, COALESCE(SUM(COALESCE(approved_amount, claim_amount)), 0)::numeric
FROM hrms.benefit_claims c
WHERE c.tenant_id=sqlc.arg('summary_tenant_id')::uuid AND NOT c.inactive AND c.payroll_export_status='ready'
UNION ALL
SELECT 'payroll_blocked', COUNT(*)::int, 0::numeric
FROM hrms.benefit_claims c
WHERE c.tenant_id=sqlc.arg('summary_tenant_id')::uuid AND NOT c.inactive AND c.payroll_export_status='blocked';
