-- name: ListJobPositions :many
SELECT
    jp.id,
    jp.tenant_id,
    jp.code,
    jp.title,
    jp.level,
    jp.category,
    jp.description,
    jp.department_id,
    d.name AS department_name,
    jp.employment_type_id,
    et.name AS employment_type_name,
    jp.work_mode,
    jp.total_position,
    jp.budgeted_cost,
    COALESCE(loc.location_count, 0)::int AS location_count,
    COALESCE(req.open_requisition_count, 0)::int AS open_requisition_count,
    jp.inactive,
    jp.created_at,
    jp.created_by,
    jp.updated_at,
    jp.updated_by
FROM hrms.job_positions jp
LEFT JOIN hrms.departments d ON d.tenant_id = jp.tenant_id AND d.id = jp.department_id AND NOT d.inactive
LEFT JOIN hrms.employment_types et ON et.tenant_id = jp.tenant_id AND et.id = jp.employment_type_id AND NOT et.inactive
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS location_count
    FROM hrms.job_position_locations jpl
    WHERE jpl.tenant_id = jp.tenant_id AND jpl.job_position_id = jp.id AND NOT jpl.inactive
) loc ON TRUE
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS open_requisition_count
    FROM hrms.job_requisitions jr
    WHERE jr.tenant_id = jp.tenant_id
      AND jr.job_position_id = jp.id
      AND jr.status IN ('Pending','Approved')
      AND NOT jr.inactive
) req ON TRUE
WHERE jp.tenant_id = $1
  AND (sqlc.narg('department_id')::uuid IS NULL OR jp.department_id = sqlc.narg('department_id')::uuid)
  AND (sqlc.narg('employment_type_id')::uuid IS NULL OR jp.employment_type_id = sqlc.narg('employment_type_id')::uuid)
  AND (sqlc.narg('work_mode')::text IS NULL OR jp.work_mode = sqlc.narg('work_mode')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR jp.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.category ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT jp.inactive
ORDER BY jp.title ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountJobPositions :one
SELECT COUNT(*) FROM hrms.job_positions jp
WHERE jp.tenant_id = $1
  AND (sqlc.narg('department_id')::uuid IS NULL OR jp.department_id = sqlc.narg('department_id')::uuid)
  AND (sqlc.narg('employment_type_id')::uuid IS NULL OR jp.employment_type_id = sqlc.narg('employment_type_id')::uuid)
  AND (sqlc.narg('work_mode')::text IS NULL OR jp.work_mode = sqlc.narg('work_mode')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR jp.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.category ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT jp.inactive;

-- name: GetJobPosition :one
SELECT * FROM hrms.job_positions
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateJobPosition :one
INSERT INTO hrms.job_positions (
    tenant_id, code, title, level, category, description, department_id, employment_type_id,
    work_mode, total_position, budgeted_cost, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: UpdateJobPosition :one
UPDATE hrms.job_positions
SET code = $3,
    title = $4,
    level = $5,
    category = $6,
    description = $7,
    department_id = $8,
    employment_type_id = $9,
    work_mode = $10,
    total_position = $11,
    budgeted_cost = $12,
    updated_at = NOW(),
    updated_by = $13
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteJobPosition :exec
UPDATE hrms.job_positions
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListJobPositionLocations :many
SELECT * FROM hrms.job_position_locations
WHERE tenant_id = $1 AND job_position_id = $2 AND NOT inactive
ORDER BY is_remote DESC, city ASC NULLS LAST;

-- name: GetJobPositionLocation :one
SELECT * FROM hrms.job_position_locations
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateJobPositionLocation :one
INSERT INTO hrms.job_position_locations (
    tenant_id, job_position_id, location, city, state, country, is_remote, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateJobPositionLocation :one
UPDATE hrms.job_position_locations
SET location = $3,
    city = $4,
    state = $5,
    country = $6,
    is_remote = $7,
    updated_at = NOW(),
    updated_by = $8
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteJobPositionLocation :exec
UPDATE hrms.job_position_locations
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListJobRequisitions :many
SELECT
    jr.id,
    jr.tenant_id,
    jr.job_position_id,
    jp.code AS job_position_code,
    jp.total_position AS position_total_headcount,
    jp.budgeted_cost AS position_budgeted_cost,
    jr.code,
    jr.title,
    jr.level,
    jr.category,
    jr.department_id,
    d.name AS department_name,
    jr.employment_type_id,
    et.name AS employment_type_name,
    jr.description,
    jr.work_mode,
    jr.total_openings,
    jr.reason_for_hire,
    jr.min_salary,
    jr.max_salary,
    jr.currency,
    jr.target_hire_date,
    jr.expected_closure_date,
    jr.requested_by,
    jr.requested_date,
    jr.is_approved,
    jr.approved_by,
    jr.approved_date,
    jr.priority,
    jr.status,
    jr.notes,
    COALESCE(logs.log_count, 0)::int AS log_count,
    jr.inactive,
    jr.created_at,
    jr.created_by,
    jr.updated_at,
    jr.updated_by
FROM hrms.job_requisitions jr
JOIN hrms.job_positions jp ON jp.tenant_id = jr.tenant_id AND jp.id = jr.job_position_id
LEFT JOIN hrms.departments d ON d.tenant_id = jr.tenant_id AND d.id = jr.department_id AND NOT d.inactive
LEFT JOIN hrms.employment_types et ON et.tenant_id = jr.tenant_id AND et.id = jr.employment_type_id AND NOT et.inactive
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS log_count
    FROM hrms.job_requisition_logs jrl
    WHERE jrl.tenant_id = jr.tenant_id AND jrl.job_requisition_id = jr.id AND NOT jrl.inactive
) logs ON TRUE
WHERE jr.tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR jr.status = sqlc.narg('status')::text)
  AND (sqlc.narg('job_position_id')::uuid IS NULL OR jr.job_position_id = sqlc.narg('job_position_id')::uuid)
  AND (sqlc.narg('department_id')::uuid IS NULL OR jr.department_id = sqlc.narg('department_id')::uuid)
  AND (
      sqlc.narg('search')::text IS NULL
      OR jr.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jr.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jr.reason_for_hire ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT jr.inactive
ORDER BY jr.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountJobRequisitions :one
SELECT COUNT(*) FROM hrms.job_requisitions jr
WHERE jr.tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR jr.status = sqlc.narg('status')::text)
  AND (sqlc.narg('job_position_id')::uuid IS NULL OR jr.job_position_id = sqlc.narg('job_position_id')::uuid)
  AND (sqlc.narg('department_id')::uuid IS NULL OR jr.department_id = sqlc.narg('department_id')::uuid)
  AND (
      sqlc.narg('search')::text IS NULL
      OR jr.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jr.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jr.reason_for_hire ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT jr.inactive;

-- name: GetJobRequisition :one
SELECT * FROM hrms.job_requisitions
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateJobRequisition :one
INSERT INTO hrms.job_requisitions (
    tenant_id, job_position_id, code, title, level, category, department_id, employment_type_id,
    description, work_mode, total_openings, reason_for_hire, min_salary, max_salary, currency,
    target_hire_date, expected_closure_date, requested_by, requested_date, priority, status, notes, created_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8,
    $9, $10, $11, $12, $13, $14, $15,
    $16, $17, $18, $19, $20, $21, $22, $23
)
RETURNING *;

-- name: UpdateJobRequisition :one
UPDATE hrms.job_requisitions
SET job_position_id = $3,
    code = $4,
    title = $5,
    level = $6,
    category = $7,
    department_id = $8,
    employment_type_id = $9,
    description = $10,
    work_mode = $11,
    total_openings = $12,
    reason_for_hire = $13,
    min_salary = $14,
    max_salary = $15,
    currency = $16,
    target_hire_date = $17,
    expected_closure_date = $18,
    requested_by = $19,
    requested_date = $20,
    priority = $21,
    notes = $22,
    updated_at = NOW(),
    updated_by = $23
WHERE tenant_id = $1 AND id = $2 AND status IN ('Draft','Rejected') AND NOT inactive
RETURNING *;

-- name: UpdateJobRequisitionStatus :one
UPDATE hrms.job_requisitions
SET status = $3,
    is_approved = CASE WHEN $3 = 'Approved' THEN TRUE WHEN $3 IN ('Rejected','Closed') THEN FALSE ELSE is_approved END,
    approved_by = CASE WHEN $3 = 'Approved' THEN $4 ELSE approved_by END,
    approved_date = CASE WHEN $3 = 'Approved' THEN NOW() ELSE approved_date END,
    notes = COALESCE($5, notes),
    updated_at = NOW(),
    updated_by = $4
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateJobRequisitionLog :one
INSERT INTO hrms.job_requisition_logs (
    tenant_id, job_requisition_id, from_status, to_status, action, remarks, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListJobRequisitionLogs :many
SELECT * FROM hrms.job_requisition_logs
WHERE tenant_id = $1 AND job_requisition_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: SoftDeleteJobRequisition :exec
UPDATE hrms.job_requisitions
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListJobPostings :many
SELECT
    jp.id,
    jp.tenant_id,
    jp.job_requisition_id,
    jr.code AS job_requisition_code,
    jr.status AS job_requisition_status,
    jp.code,
    jp.title,
    jp.job_summary,
    jp.description,
    jp.job_category,
    jp.department_id,
    d.name AS department_name,
    jp.industry,
    jp.employment_type_id,
    et.name AS employment_type_name,
    jp.work_mode,
    jp.role_type,
    jp.min_experience,
    jp.max_experience,
    jp.min_salary,
    jp.max_salary,
    jp.salary_currency,
    jp.salary_period,
    jp.is_salary_visible,
    CASE
        WHEN jp.is_published AND jp.expiry_date IS NOT NULL AND jp.expiry_date < CURRENT_DATE THEN 'Expired'
        ELSE COALESCE(jp.job_status, 'Draft')
    END::text AS effective_status,
    jp.job_status,
    jp.publish_date,
    jp.expiry_date,
    jp.is_published,
    jp.inactive,
    jp.created_at,
    jp.created_by,
    jp.updated_at,
    jp.updated_by
FROM hrms.job_postings jp
LEFT JOIN hrms.job_requisitions jr ON jr.tenant_id = jp.tenant_id AND jr.id = jp.job_requisition_id AND NOT jr.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = jp.tenant_id AND d.id = jp.department_id AND NOT d.inactive
LEFT JOIN hrms.employment_types et ON et.tenant_id = jp.tenant_id AND et.id = jp.employment_type_id AND NOT et.inactive
WHERE jp.tenant_id = $1
  AND (sqlc.narg('job_status')::text IS NULL OR jp.job_status = sqlc.narg('job_status')::text)
  AND (sqlc.narg('is_published')::boolean IS NULL OR jp.is_published = sqlc.narg('is_published')::boolean)
  AND (sqlc.narg('department_id')::uuid IS NULL OR jp.department_id = sqlc.narg('department_id')::uuid)
  AND (sqlc.narg('search')::text IS NULL OR jp.title ILIKE '%' || sqlc.narg('search')::text || '%' OR jp.code ILIKE '%' || sqlc.narg('search')::text || '%' OR jp.job_category ILIKE '%' || sqlc.narg('search')::text || '%')
  AND NOT jp.inactive
ORDER BY jp.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountJobPostings :one
SELECT COUNT(*) FROM hrms.job_postings jp
WHERE jp.tenant_id = $1
  AND (sqlc.narg('job_status')::text IS NULL OR jp.job_status = sqlc.narg('job_status')::text)
  AND (sqlc.narg('is_published')::boolean IS NULL OR jp.is_published = sqlc.narg('is_published')::boolean)
  AND (sqlc.narg('department_id')::uuid IS NULL OR jp.department_id = sqlc.narg('department_id')::uuid)
  AND (sqlc.narg('search')::text IS NULL OR jp.title ILIKE '%' || sqlc.narg('search')::text || '%' OR jp.code ILIKE '%' || sqlc.narg('search')::text || '%' OR jp.job_category ILIKE '%' || sqlc.narg('search')::text || '%')
  AND NOT jp.inactive;

-- name: ListPublishedJobPostings :many
SELECT * FROM hrms.job_postings
WHERE tenant_id = $1
  AND is_published
  AND (expiry_date IS NULL OR expiry_date >= CURRENT_DATE)
  AND NOT inactive
ORDER BY publish_date DESC NULLS LAST;

-- name: GetJobPosting :one
SELECT * FROM hrms.job_postings
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetJobPostingByRequisition :one
SELECT * FROM hrms.job_postings
WHERE tenant_id = $1 AND job_requisition_id = $2 AND NOT inactive
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateJobPosting :one
INSERT INTO hrms.job_postings (
    tenant_id, job_requisition_id, code, title, job_summary, description, job_category,
    department_id, industry, employment_type_id, work_mode, role_type, min_experience,
    max_experience, min_salary, max_salary, salary_currency, salary_period,
    is_salary_visible, job_status, publish_date, expiry_date, is_published, created_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12, $13,
    $14, $15, $16, $17, $18,
    $19, $20, $21, $22, $23, $24
)
RETURNING *;

-- name: UpdateJobPosting :one
UPDATE hrms.job_postings
SET code = $3,
    title = $4,
    job_summary = $5,
    description = $6,
    job_category = $7,
    department_id = $8,
    industry = $9,
    employment_type_id = $10,
    work_mode = $11,
    role_type = $12,
    min_experience = $13,
    max_experience = $14,
    min_salary = $15,
    max_salary = $16,
    salary_currency = $17,
    salary_period = $18,
    is_salary_visible = $19,
    job_status = $20,
    publish_date = $21,
    expiry_date = $22,
    is_published = $23,
    updated_at = NOW(),
    updated_by = $24
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: PublishJobPosting :one
UPDATE hrms.job_postings
SET is_published = TRUE,
    job_status = 'Open',
    publish_date = COALESCE(publish_date, CURRENT_DATE),
    expiry_date = COALESCE(sqlc.narg('expiry_date')::date, expiry_date),
    updated_at = NOW(),
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ExpireJobPostings :many
UPDATE hrms.job_postings
SET is_published = FALSE,
    job_status = 'Expired',
    updated_at = NOW()
WHERE tenant_id = $1
  AND is_published
  AND expiry_date IS NOT NULL
  AND expiry_date < CURRENT_DATE
  AND NOT inactive
RETURNING *;

-- name: CloseJobPosting :one
UPDATE hrms.job_postings
SET is_published = FALSE,
    job_status = 'Closed',
    updated_at = NOW(),
    updated_by = $3
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteJobPosting :exec
UPDATE hrms.job_postings
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListCandidates :many
SELECT *,
       COUNT(*) OVER()::bigint AS total_count
FROM hrms.candidates
WHERE tenant_id = sqlc.arg(tenant_id)
  AND NOT inactive
  AND (
    sqlc.narg(search)::text IS NULL
    OR firstname ILIKE '%' || sqlc.narg(search)::text || '%'
    OR lastname ILIKE '%' || sqlc.narg(search)::text || '%'
    OR email ILIKE '%' || sqlc.narg(search)::text || '%'
    OR phone ILIKE '%' || sqlc.narg(search)::text || '%'
    OR current_company ILIKE '%' || sqlc.narg(search)::text || '%'
    OR current_designation ILIKE '%' || sqlc.narg(search)::text || '%'
    OR current_location ILIKE '%' || sqlc.narg(search)::text || '%'
    OR preferred_location ILIKE '%' || sqlc.narg(search)::text || '%'
    OR source ILIKE '%' || sqlc.narg(search)::text || '%'
  )
  AND (sqlc.narg(source)::text IS NULL OR source = sqlc.narg(source)::text)
  AND (sqlc.narg(gender)::text IS NULL OR gender = sqlc.narg(gender)::text)
ORDER BY created_at DESC
OFFSET sqlc.arg(offset_rows)
LIMIT sqlc.arg(limit_rows);

-- name: CountCandidates :one
SELECT COUNT(*)::bigint
FROM hrms.candidates
WHERE tenant_id = sqlc.arg(tenant_id)
  AND NOT inactive
  AND (
    sqlc.narg(search)::text IS NULL
    OR firstname ILIKE '%' || sqlc.narg(search)::text || '%'
    OR lastname ILIKE '%' || sqlc.narg(search)::text || '%'
    OR email ILIKE '%' || sqlc.narg(search)::text || '%'
    OR phone ILIKE '%' || sqlc.narg(search)::text || '%'
    OR current_company ILIKE '%' || sqlc.narg(search)::text || '%'
    OR current_designation ILIKE '%' || sqlc.narg(search)::text || '%'
    OR current_location ILIKE '%' || sqlc.narg(search)::text || '%'
    OR preferred_location ILIKE '%' || sqlc.narg(search)::text || '%'
    OR source ILIKE '%' || sqlc.narg(search)::text || '%'
  )
  AND (sqlc.narg(source)::text IS NULL OR source = sqlc.narg(source)::text)
  AND (sqlc.narg(gender)::text IS NULL OR gender = sqlc.narg(gender)::text);

-- name: GetCandidate :one
SELECT * FROM hrms.candidates
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateCandidate :one
INSERT INTO hrms.candidates (
  tenant_id, firstname, lastname, email, phone, dob, gender, total_experience,
  current_company, current_designation, current_salary, expected_salary, notice_period,
  current_location, preferred_location, source, resume_url, created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8,
  $9, $10, $11, $12, $13,
  $14, $15, $16, $17, $18
)
RETURNING *;

-- name: UpdateCandidate :one
UPDATE hrms.candidates
SET firstname = $3,
    lastname = $4,
    email = $5,
    phone = $6,
    dob = $7,
    gender = $8,
    total_experience = $9,
    current_company = $10,
    current_designation = $11,
    current_salary = $12,
    expected_salary = $13,
    notice_period = $14,
    current_location = $15,
    preferred_location = $16,
    source = $17,
    resume_url = $18,
    updated_at = NOW(),
    updated_by = $19
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteCandidate :exec
UPDATE hrms.candidates
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: CreateCandidateApplicantAccount :one
INSERT INTO hrms.candidate_applicant_accounts (
    tenant_id, candidate_id, user_id, email, status, consent_at, consent_ip, metadata, created_by, updated_by
)
VALUES (
    sqlc.arg(tenant_id), sqlc.arg(candidate_id), sqlc.arg(user_id), lower(sqlc.arg(email)),
    COALESCE(sqlc.narg(status)::text, 'active'),
    sqlc.narg(consent_at)::timestamptz,
    sqlc.narg(consent_ip)::text,
    COALESCE(sqlc.narg(metadata)::jsonb, '{}'::jsonb),
    sqlc.narg(created_by)::uuid,
    sqlc.narg(created_by)::uuid
)
ON CONFLICT (tenant_id, user_id) WHERE inactive = FALSE DO UPDATE
SET candidate_id = EXCLUDED.candidate_id,
    email = EXCLUDED.email,
    status = EXCLUDED.status,
    consent_at = COALESCE(EXCLUDED.consent_at, hrms.candidate_applicant_accounts.consent_at),
    consent_ip = COALESCE(EXCLUDED.consent_ip, hrms.candidate_applicant_accounts.consent_ip),
    metadata = COALESCE(EXCLUDED.metadata, hrms.candidate_applicant_accounts.metadata),
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by
RETURNING *;

-- name: GetCandidateApplicantAccountByUser :one
SELECT * FROM hrms.candidate_applicant_accounts
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive;

-- name: GetCandidateApplicantAccountByCandidate :one
SELECT * FROM hrms.candidate_applicant_accounts
WHERE tenant_id = $1 AND candidate_id = $2 AND NOT inactive;

-- name: ListCandidateApplications :many
SELECT ca.*,
       c.firstname AS candidate_firstname,
       c.lastname AS candidate_lastname,
       c.email AS candidate_email,
       c.phone AS candidate_phone,
       jp.title AS job_posting_title,
       jp.code AS job_posting_code,
       COUNT(*) OVER()::bigint AS total_count
FROM hrms.candidate_applications ca
LEFT JOIN hrms.candidates c ON c.id = ca.candidate_id AND c.tenant_id = ca.tenant_id
LEFT JOIN hrms.job_postings jp ON jp.id = ca.job_posting_id AND jp.tenant_id = ca.tenant_id
WHERE ca.tenant_id = sqlc.arg(tenant_id)
  AND NOT ca.inactive
  AND (sqlc.narg(status)::text IS NULL OR ca.status = sqlc.narg(status)::text)
  AND (sqlc.narg(candidate_id)::uuid IS NULL OR ca.candidate_id = sqlc.narg(candidate_id)::uuid)
  AND (sqlc.narg(job_posting_id)::uuid IS NULL OR ca.job_posting_id = sqlc.narg(job_posting_id)::uuid)
  AND (
    sqlc.narg(search)::text IS NULL
    OR c.firstname ILIKE '%' || sqlc.narg(search)::text || '%'
    OR c.lastname ILIKE '%' || sqlc.narg(search)::text || '%'
    OR c.email ILIKE '%' || sqlc.narg(search)::text || '%'
    OR ca.source ILIKE '%' || sqlc.narg(search)::text || '%'
    OR ca.referred_by ILIKE '%' || sqlc.narg(search)::text || '%'
    OR ca.comments ILIKE '%' || sqlc.narg(search)::text || '%'
    OR jp.title ILIKE '%' || sqlc.narg(search)::text || '%'
    OR jp.code ILIKE '%' || sqlc.narg(search)::text || '%'
  )
ORDER BY ca.updated_at DESC, ca.created_at DESC
OFFSET sqlc.arg(offset_rows)
LIMIT sqlc.arg(limit_rows);

-- name: CountCandidateApplications :one
SELECT COUNT(*)::bigint
FROM hrms.candidate_applications ca
LEFT JOIN hrms.candidates c ON c.id = ca.candidate_id AND c.tenant_id = ca.tenant_id
LEFT JOIN hrms.job_postings jp ON jp.id = ca.job_posting_id AND jp.tenant_id = ca.tenant_id
WHERE ca.tenant_id = sqlc.arg(tenant_id)
  AND NOT ca.inactive
  AND (sqlc.narg(status)::text IS NULL OR ca.status = sqlc.narg(status)::text)
  AND (sqlc.narg(candidate_id)::uuid IS NULL OR ca.candidate_id = sqlc.narg(candidate_id)::uuid)
  AND (sqlc.narg(job_posting_id)::uuid IS NULL OR ca.job_posting_id = sqlc.narg(job_posting_id)::uuid)
  AND (
    sqlc.narg(search)::text IS NULL
    OR c.firstname ILIKE '%' || sqlc.narg(search)::text || '%'
    OR c.lastname ILIKE '%' || sqlc.narg(search)::text || '%'
    OR c.email ILIKE '%' || sqlc.narg(search)::text || '%'
    OR ca.source ILIKE '%' || sqlc.narg(search)::text || '%'
    OR ca.referred_by ILIKE '%' || sqlc.narg(search)::text || '%'
    OR ca.comments ILIKE '%' || sqlc.narg(search)::text || '%'
    OR jp.title ILIKE '%' || sqlc.narg(search)::text || '%'
    OR jp.code ILIKE '%' || sqlc.narg(search)::text || '%'
  );

-- name: ListCandidateApplicationsByStatus :many
SELECT * FROM hrms.candidate_applications
WHERE tenant_id = $1 AND status = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: ListCandidateApplicationsByCandidate :many
SELECT * FROM hrms.candidate_applications
WHERE tenant_id = $1 AND candidate_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: GetCandidateApplication :one
SELECT * FROM hrms.candidate_applications
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateCandidateApplication :one
INSERT INTO hrms.candidate_applications (
  tenant_id, candidate_id, job_posting_id, resume_url, cover_letter, current_ctc,
  expected_ctc, notice_period, referred_by, source, status, comments, created_by
  , applied_at, status_changed_at, status_changed_by, rejection_reason, withdrawal_reason,
  source_detail, duplicate_of_application_id
) VALUES (
  $1, $2, $3, $4, $5, $6,
  $7, $8, $9, $10, $11, $12, $13,
  COALESCE($14::timestamptz, NOW()), NOW(), $13, $15, $16,
  $17, $18
)
RETURNING *;

-- name: UpdateCandidateApplication :one
UPDATE hrms.candidate_applications
SET candidate_id = $3,
    job_posting_id = $4,
    resume_url = $5,
    cover_letter = $6,
    current_ctc = $7,
    expected_ctc = $8,
    notice_period = $9,
    referred_by = $10,
    source = $11,
    status = $12,
    comments = $13,
    applied_at = $15,
    status_changed_at = CASE WHEN status <> $12 THEN NOW() ELSE status_changed_at END,
    status_changed_by = CASE WHEN status <> $12 THEN $14 ELSE status_changed_by END,
    rejection_reason = $16,
    withdrawal_reason = $17,
    source_detail = $18,
    duplicate_of_application_id = $19,
    updated_at = NOW(),
    updated_by = $14
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: MoveCandidateApplicationStatus :one
UPDATE hrms.candidate_applications
SET status = $3,
    comments = COALESCE($4, comments),
    rejection_reason = CASE WHEN $3 = 'Rejected' THEN $6 ELSE rejection_reason END,
    withdrawal_reason = CASE WHEN $3 = 'Withdrawn' THEN $6 ELSE withdrawal_reason END,
    status_changed_at = NOW(),
    status_changed_by = $5,
    updated_at = NOW(),
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateCandidateApplicationEvent :one
INSERT INTO hrms.candidate_application_events (
  tenant_id, application_id, from_status, to_status, action, reason, remarks, created_by
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListCandidateApplicationEvents :many
SELECT * FROM hrms.candidate_application_events
WHERE tenant_id = $1 AND application_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: SoftDeleteCandidateApplication :exec
UPDATE hrms.candidate_applications
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListCandidateEducation :many
SELECT * FROM hrms.candidate_education
WHERE tenant_id = $1 AND candidate_id = $2 AND NOT inactive
ORDER BY end_date DESC NULLS LAST;

-- name: SoftDeleteCandidateEducation :exec
UPDATE hrms.candidate_education
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListCandidateExperience :many
SELECT * FROM hrms.candidate_experience
WHERE tenant_id = $1 AND candidate_id = $2 AND NOT inactive
ORDER BY end_date DESC NULLS LAST;

-- name: SoftDeleteCandidateExperience :exec
UPDATE hrms.candidate_experience
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListInterviewRoundsByApplication :many
SELECT * FROM hrms.interview_rounds
WHERE tenant_id = $1 AND application_id = $2 AND NOT inactive
ORDER BY round_number ASC NULLS LAST, scheduled_date ASC NULLS LAST;

-- name: ListInterviewRounds :many
SELECT
    ir.id,
    ir.tenant_id,
    ir.application_id,
    ca.status AS application_status,
    c.firstname AS candidate_firstname,
    c.lastname AS candidate_lastname,
    c.email AS candidate_email,
    jp.title AS job_posting_title,
    jp.code AS job_posting_code,
    ir.round_name,
    ir.round_number,
    ir.scheduled_date,
    ir.duration_minutes,
    ir.interviewer_user_id,
    ir.mode,
    ir.meeting_link,
    ir.location,
    ir.status,
    ir.remarks,
    ir.timezone,
    ir.feedback,
    ir.score,
    ir.decision,
    ir.completed_at,
    ir.inactive,
    ir.created_at,
    ir.created_by,
    ir.updated_at,
    ir.updated_by
FROM hrms.interview_rounds ir
JOIN hrms.candidate_applications ca ON ca.tenant_id = ir.tenant_id AND ca.id = ir.application_id AND NOT ca.inactive
LEFT JOIN hrms.candidates c ON c.tenant_id = ca.tenant_id AND c.id = ca.candidate_id AND NOT c.inactive
LEFT JOIN hrms.job_postings jp ON jp.tenant_id = ca.tenant_id AND jp.id = ca.job_posting_id AND NOT jp.inactive
WHERE ir.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('application_id')::uuid IS NULL OR ir.application_id = sqlc.narg('application_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR ir.status = sqlc.narg('status')::text)
  AND (sqlc.narg('interviewer_user_id')::uuid IS NULL OR ir.interviewer_user_id = sqlc.narg('interviewer_user_id')::uuid)
  AND (sqlc.narg('date_from')::timestamptz IS NULL OR ir.scheduled_date >= sqlc.narg('date_from')::timestamptz)
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR ir.scheduled_date <= sqlc.narg('date_to')::timestamptz)
  AND (
      sqlc.narg('search')::text IS NULL
      OR ir.round_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.email ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT ir.inactive
ORDER BY ir.scheduled_date ASC NULLS LAST, ir.round_number ASC NULLS LAST, ir.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountInterviewRounds :one
SELECT COUNT(*)
FROM hrms.interview_rounds ir
JOIN hrms.candidate_applications ca ON ca.tenant_id = ir.tenant_id AND ca.id = ir.application_id AND NOT ca.inactive
LEFT JOIN hrms.candidates c ON c.tenant_id = ca.tenant_id AND c.id = ca.candidate_id AND NOT c.inactive
LEFT JOIN hrms.job_postings jp ON jp.tenant_id = ca.tenant_id AND jp.id = ca.job_posting_id AND NOT jp.inactive
WHERE ir.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('application_id')::uuid IS NULL OR ir.application_id = sqlc.narg('application_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR ir.status = sqlc.narg('status')::text)
  AND (sqlc.narg('interviewer_user_id')::uuid IS NULL OR ir.interviewer_user_id = sqlc.narg('interviewer_user_id')::uuid)
  AND (sqlc.narg('date_from')::timestamptz IS NULL OR ir.scheduled_date >= sqlc.narg('date_from')::timestamptz)
  AND (sqlc.narg('date_to')::timestamptz IS NULL OR ir.scheduled_date <= sqlc.narg('date_to')::timestamptz)
  AND (
      sqlc.narg('search')::text IS NULL
      OR ir.round_name ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.email ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.code ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT ir.inactive;

-- name: GetInterviewRound :one
SELECT * FROM hrms.interview_rounds
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateInterviewRound :one
INSERT INTO hrms.interview_rounds (
    tenant_id, application_id, round_name, round_number, scheduled_date, duration_minutes,
    interviewer_user_id, mode, meeting_link, location, status, remarks, timezone,
    feedback, score, decision, completed_at, created_by
)
VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12, $13,
    $14, $15, $16, $17, $18
)
RETURNING *;

-- name: UpdateInterviewRound :one
UPDATE hrms.interview_rounds
SET application_id = $3,
    round_name = $4,
    round_number = $5,
    scheduled_date = $6,
    duration_minutes = $7,
    interviewer_user_id = $8,
    mode = $9,
    meeting_link = $10,
    location = $11,
    status = $12,
    remarks = $13,
    timezone = $14,
    feedback = $15,
    score = $16,
    decision = $17,
    completed_at = $18,
    updated_at = NOW(),
    updated_by = $19
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: UpdateInterviewRoundStatus :one
UPDATE hrms.interview_rounds
SET status = $3,
    remarks = COALESCE($4, remarks),
    feedback = COALESCE($5, feedback),
    score = COALESCE($6, score),
    decision = COALESCE($7, decision),
    completed_at = CASE WHEN $3 = 'Completed' THEN COALESCE($8, NOW()) ELSE completed_at END,
    updated_at = NOW(),
    updated_by = $9
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteInterviewRound :exec
UPDATE hrms.interview_rounds
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListOfferLetterTemplates :many
SELECT * FROM hrms.offer_letter_templates
WHERE tenant_id = $1 AND NOT inactive
ORDER BY is_default DESC, name ASC;

-- name: GetDefaultOfferLetterTemplate :one
SELECT * FROM hrms.offer_letter_templates
WHERE tenant_id = $1 AND is_default AND is_active AND NOT inactive
LIMIT 1;

-- name: GetOfferLetterTemplate :one
SELECT * FROM hrms.offer_letter_templates
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateOfferLetterTemplate :one
INSERT INTO hrms.offer_letter_templates (
    tenant_id, name, description, subject, body_html, footer_html, locale, is_default, is_active, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateOfferLetterTemplate :one
UPDATE hrms.offer_letter_templates
SET name = $3,
    description = $4,
    subject = $5,
    body_html = $6,
    footer_html = $7,
    locale = $8,
    is_default = $9,
    is_active = $10,
    updated_at = NOW(),
    updated_by = $11
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ClearDefaultOfferLetterTemplates :exec
UPDATE hrms.offer_letter_templates
SET is_default = FALSE, updated_at = NOW(), updated_by = $2
WHERE tenant_id = $1 AND is_default AND NOT inactive;

-- name: SoftDeleteOfferLetterTemplate :exec
UPDATE hrms.offer_letter_templates
SET inactive = TRUE, is_active = FALSE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListOfferLetters :many
SELECT
    ol.id,
    ol.tenant_id,
    ol.application_id,
    ol.candidate_id,
    c.firstname AS candidate_firstname,
    c.lastname AS candidate_lastname,
    c.email AS candidate_email,
    jp.title AS job_posting_title,
    jp.code AS job_posting_code,
    ol.template_id,
    olt.name AS template_name,
    ol.offered_ctc,
    ol.currency,
    ol.salary_breakdown,
    ol.joining_date,
    ol.valid_until_date,
    ol.status,
    ol.offer_letter_url,
    ol.candidate_reaction_date,
    ol.candidate_rejection_reason,
    ol.version,
    ol.is_latest,
    ol.subject,
    ol.rendered_html,
    ol.sent_at,
    ol.revoked_at,
    ol.signature_token,
    ol.signature_requested_at,
    ol.signature_completed_at,
    ol.signer_name,
    ol.signer_email,
    ol.signer_ip,
    ol.signer_user_agent,
    ol.signature_hash,
    ol.audit_certificate_url,
    ol.inactive,
    ol.created_at,
    ol.created_by,
    ol.updated_at,
    ol.updated_by
FROM hrms.offer_letters ol
JOIN hrms.candidate_applications ca ON ca.tenant_id = ol.tenant_id AND ca.id = ol.application_id AND NOT ca.inactive
LEFT JOIN hrms.candidates c ON c.tenant_id = ol.tenant_id AND c.id = ol.candidate_id AND NOT c.inactive
LEFT JOIN hrms.job_postings jp ON jp.tenant_id = ca.tenant_id AND jp.id = ca.job_posting_id AND NOT jp.inactive
LEFT JOIN hrms.offer_letter_templates olt ON olt.tenant_id = ol.tenant_id AND olt.id = ol.template_id
WHERE ol.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('application_id')::uuid IS NULL OR ol.application_id = sqlc.narg('application_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR ol.status = sqlc.narg('status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.email ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ol.subject ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT ol.inactive
ORDER BY ol.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountOfferLetters :one
SELECT COUNT(*)
FROM hrms.offer_letters ol
JOIN hrms.candidate_applications ca ON ca.tenant_id = ol.tenant_id AND ca.id = ol.application_id AND NOT ca.inactive
LEFT JOIN hrms.candidates c ON c.tenant_id = ol.tenant_id AND c.id = ol.candidate_id AND NOT c.inactive
LEFT JOIN hrms.job_postings jp ON jp.tenant_id = ca.tenant_id AND jp.id = ca.job_posting_id AND NOT jp.inactive
WHERE ol.tenant_id = sqlc.arg('tenant_id')
  AND (sqlc.narg('application_id')::uuid IS NULL OR ol.application_id = sqlc.narg('application_id')::uuid)
  AND (sqlc.narg('status')::text IS NULL OR ol.status = sqlc.narg('status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.email ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.title ILIKE '%' || sqlc.narg('search')::text || '%'
      OR jp.code ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ol.subject ILIKE '%' || sqlc.narg('search')::text || '%'
  )
  AND NOT ol.inactive;

-- name: ListOfferLettersByApplication :many
SELECT * FROM hrms.offer_letters
WHERE tenant_id = $1 AND application_id = $2 AND NOT inactive
ORDER BY version DESC;

-- name: GetLatestOfferLetterByApplication :one
SELECT * FROM hrms.offer_letters
WHERE tenant_id = $1 AND application_id = $2 AND is_latest AND NOT inactive;

-- name: GetOfferLetter :one
SELECT * FROM hrms.offer_letters
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetOfferLetterBySignatureToken :one
SELECT * FROM hrms.offer_letters
WHERE signature_token = $1 AND NOT inactive;

-- name: ClearLatestOfferLetters :exec
UPDATE hrms.offer_letters
SET is_latest = FALSE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND application_id = $2 AND is_latest AND NOT inactive;

-- name: NextOfferLetterVersion :one
SELECT COALESCE(MAX(version), 0)::int + 1 AS next_version
FROM hrms.offer_letters
WHERE tenant_id = $1 AND application_id = $2;

-- name: CreateOfferLetter :one
INSERT INTO hrms.offer_letters (
    tenant_id, application_id, candidate_id, template_id, offered_ctc, currency, salary_breakdown,
    joining_date, valid_until_date, status, offer_letter_url, version, is_latest, subject,
    rendered_html, signature_token, signature_requested_at, signer_email, created_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12, $13, $14,
    $15, $16, $17, $18, $19
)
RETURNING *;

-- name: UpdateOfferLetter :one
UPDATE hrms.offer_letters
SET template_id = $3,
    offered_ctc = $4,
    currency = $5,
    salary_breakdown = $6,
    joining_date = $7,
    valid_until_date = $8,
    status = $9,
    offer_letter_url = $10,
    subject = $11,
    rendered_html = $12,
    signer_email = $13,
    updated_at = NOW(),
    updated_by = $14
WHERE tenant_id = $1 AND id = $2 AND status = 'Generated' AND NOT inactive
RETURNING *;

-- name: UpdateOfferLetterStatus :one
UPDATE hrms.offer_letters
SET status = $3,
    sent_at = CASE WHEN $3 = 'Sent' THEN COALESCE(sent_at, NOW()) ELSE sent_at END,
    signature_requested_at = CASE WHEN $3 = 'Sent' THEN COALESCE(signature_requested_at, NOW()) ELSE signature_requested_at END,
    candidate_reaction_date = CASE WHEN $3 IN ('Accepted','Declined') THEN COALESCE(candidate_reaction_date, NOW()) ELSE candidate_reaction_date END,
    revoked_at = CASE WHEN $3 = 'Revoked' THEN COALESCE(revoked_at, NOW()) ELSE revoked_at END,
    candidate_rejection_reason = COALESCE($4, candidate_rejection_reason),
    updated_at = NOW(),
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SignOfferLetter :one
UPDATE hrms.offer_letters
SET status = 'Accepted',
    candidate_reaction_date = NOW(),
    signature_completed_at = NOW(),
    signer_name = $2,
    signer_email = $3,
    signer_ip = $4,
    signer_user_agent = $5,
    signature_hash = $6,
    updated_at = NOW()
WHERE signature_token = $1 AND status = 'Sent' AND NOT inactive
RETURNING *;

-- name: CreateOfferLetterEvent :one
INSERT INTO hrms.offer_letter_events (
    tenant_id, offer_letter_id, from_status, to_status, action, remarks, actor_email, ip_address, user_agent, metadata, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, COALESCE($10, '{}'::jsonb), $11)
RETURNING *;

-- name: ListOfferLetterEvents :many
SELECT * FROM hrms.offer_letter_events
WHERE tenant_id = $1 AND offer_letter_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: SoftDeleteOfferLetter :exec
UPDATE hrms.offer_letters
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListOnboardingWorkflows :many
SELECT * FROM hrms.onboarding_workflows
WHERE tenant_id = $1 AND NOT inactive
ORDER BY name ASC;

-- name: GetOnboardingWorkflow :one
SELECT * FROM hrms.onboarding_workflows
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateOnboardingWorkflow :one
INSERT INTO hrms.onboarding_workflows (
    tenant_id, name, description, is_default, is_active, created_by
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateOnboardingWorkflow :one
UPDATE hrms.onboarding_workflows
SET name = $3,
    description = $4,
    is_default = $5,
    is_active = $6,
    updated_at = NOW(),
    updated_by = $7
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: ClearDefaultOnboardingWorkflows :exec
UPDATE hrms.onboarding_workflows
SET is_default = FALSE, updated_at = NOW(), updated_by = $2
WHERE tenant_id = $1 AND is_default AND NOT inactive;

-- name: SoftDeleteOnboardingWorkflow :exec
UPDATE hrms.onboarding_workflows
SET inactive = TRUE, is_active = FALSE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListOnboardingTasksByWorkflow :many
SELECT * FROM hrms.onboarding_tasks
WHERE tenant_id = $1 AND workflow_id = $2 AND NOT inactive
ORDER BY sort_order ASC, title ASC;

-- name: GetOnboardingTask :one
SELECT * FROM hrms.onboarding_tasks
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateOnboardingTask :one
INSERT INTO hrms.onboarding_tasks (
    tenant_id, workflow_id, title, description, due_days, is_required, sort_order, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateOnboardingTask :one
UPDATE hrms.onboarding_tasks
SET title = $3,
    description = $4,
    due_days = $5,
    is_required = $6,
    sort_order = $7,
    updated_at = NOW(),
    updated_by = $8
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteOnboardingTask :exec
UPDATE hrms.onboarding_tasks
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListOnboardingWorkflowAssignments :many
SELECT
    owa.*,
    ow.name AS workflow_name,
    jp.title AS job_posting_title,
    jpos.title AS job_position_title,
    d.name AS department_name,
    et.name AS employment_type_name
FROM hrms.onboarding_workflow_assignments owa
JOIN hrms.onboarding_workflows ow ON ow.tenant_id = owa.tenant_id AND ow.id = owa.workflow_id AND NOT ow.inactive
LEFT JOIN hrms.job_postings jp ON jp.tenant_id = owa.tenant_id AND jp.id = owa.job_posting_id AND NOT jp.inactive
LEFT JOIN hrms.job_positions jpos ON jpos.tenant_id = owa.tenant_id AND jpos.id = owa.job_position_id AND NOT jpos.inactive
LEFT JOIN hrms.departments d ON d.tenant_id = owa.tenant_id AND d.id = owa.department_id AND NOT d.inactive
LEFT JOIN hrms.employment_types et ON et.tenant_id = owa.tenant_id AND et.id = owa.employment_type_id AND NOT et.inactive
WHERE owa.tenant_id = $1 AND NOT owa.inactive
ORDER BY owa.priority ASC, owa.name ASC;

-- name: GetOnboardingWorkflowAssignment :one
SELECT * FROM hrms.onboarding_workflow_assignments
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateOnboardingWorkflowAssignment :one
INSERT INTO hrms.onboarding_workflow_assignments (
    tenant_id, workflow_id, name, job_posting_id, job_position_id, department_id, employment_type_id, priority, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateOnboardingWorkflowAssignment :one
UPDATE hrms.onboarding_workflow_assignments
SET workflow_id = $3,
    name = $4,
    job_posting_id = $5,
    job_position_id = $6,
    department_id = $7,
    employment_type_id = $8,
    priority = $9,
    updated_at = NOW(),
    updated_by = $10
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteOnboardingWorkflowAssignment :exec
UPDATE hrms.onboarding_workflow_assignments
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListCandidateOnboardings :many
SELECT
    co.*,
    c.firstname AS candidate_firstname,
    c.lastname AS candidate_lastname,
    c.email AS candidate_email,
    ow.name AS workflow_name,
    COUNT(cot.id)::int AS total_tasks,
    COUNT(cot.id) FILTER (WHERE cot.status = 'Completed')::int AS completed_tasks,
    COUNT(cot.id) FILTER (WHERE ot.is_required)::int AS required_tasks,
    COUNT(cot.id) FILTER (WHERE ot.is_required AND cot.status = 'Completed')::int AS completed_required_tasks,
    COUNT(cot.id) FILTER (WHERE cot.due_at IS NOT NULL AND cot.due_at < NOW() AND cot.status <> 'Completed')::int AS overdue_tasks
FROM hrms.candidate_onboardings co
JOIN hrms.candidates c ON c.tenant_id = co.tenant_id AND c.id = co.candidate_id AND NOT c.inactive
JOIN hrms.onboarding_workflows ow ON ow.tenant_id = co.tenant_id AND ow.id = co.workflow_id AND NOT ow.inactive
LEFT JOIN hrms.candidate_onboarding_tasks cot ON cot.tenant_id = co.tenant_id AND cot.candidate_onboarding_id = co.id AND NOT cot.inactive
LEFT JOIN hrms.onboarding_tasks ot ON ot.tenant_id = co.tenant_id AND ot.id = cot.onboarding_task_id AND NOT ot.inactive
WHERE co.tenant_id = $1
  AND NOT co.inactive
  AND (sqlc.narg('status')::text IS NULL OR co.onboarding_status = sqlc.narg('status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.email ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ow.name ILIKE '%' || sqlc.narg('search')::text || '%'
  )
GROUP BY co.id, c.firstname, c.lastname, c.email, ow.name
ORDER BY co.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountCandidateOnboardings :one
SELECT COUNT(*)::bigint
FROM hrms.candidate_onboardings co
JOIN hrms.candidates c ON c.tenant_id = co.tenant_id AND c.id = co.candidate_id AND NOT c.inactive
JOIN hrms.onboarding_workflows ow ON ow.tenant_id = co.tenant_id AND ow.id = co.workflow_id AND NOT ow.inactive
WHERE co.tenant_id = $1
  AND NOT co.inactive
  AND (sqlc.narg('status')::text IS NULL OR co.onboarding_status = sqlc.narg('status')::text)
  AND (
      sqlc.narg('search')::text IS NULL
      OR c.firstname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.lastname ILIKE '%' || sqlc.narg('search')::text || '%'
      OR c.email ILIKE '%' || sqlc.narg('search')::text || '%'
      OR ow.name ILIKE '%' || sqlc.narg('search')::text || '%'
  );

-- name: GetCandidateOnboarding :one
SELECT
    co.*,
    c.firstname AS candidate_firstname,
    c.lastname AS candidate_lastname,
    c.email AS candidate_email,
    ow.name AS workflow_name,
    COUNT(cot.id)::int AS total_tasks,
    COUNT(cot.id) FILTER (WHERE cot.status = 'Completed')::int AS completed_tasks,
    COUNT(cot.id) FILTER (WHERE ot.is_required)::int AS required_tasks,
    COUNT(cot.id) FILTER (WHERE ot.is_required AND cot.status = 'Completed')::int AS completed_required_tasks,
    COUNT(cot.id) FILTER (WHERE cot.due_at IS NOT NULL AND cot.due_at < NOW() AND cot.status <> 'Completed')::int AS overdue_tasks
FROM hrms.candidate_onboardings co
JOIN hrms.candidates c ON c.tenant_id = co.tenant_id AND c.id = co.candidate_id AND NOT c.inactive
JOIN hrms.onboarding_workflows ow ON ow.tenant_id = co.tenant_id AND ow.id = co.workflow_id AND NOT ow.inactive
LEFT JOIN hrms.candidate_onboarding_tasks cot ON cot.tenant_id = co.tenant_id AND cot.candidate_onboarding_id = co.id AND NOT cot.inactive
LEFT JOIN hrms.onboarding_tasks ot ON ot.tenant_id = co.tenant_id AND ot.id = cot.onboarding_task_id AND NOT ot.inactive
WHERE co.tenant_id = $1 AND co.id = $2 AND NOT co.inactive
GROUP BY co.id, c.firstname, c.lastname, c.email, ow.name;

-- name: GetCandidateOnboardingByCandidate :one
SELECT
    co.*,
    c.firstname AS candidate_firstname,
    c.lastname AS candidate_lastname,
    c.email AS candidate_email,
    ow.name AS workflow_name,
    COUNT(cot.id)::int AS total_tasks,
    COUNT(cot.id) FILTER (WHERE cot.status = 'Completed')::int AS completed_tasks,
    COUNT(cot.id) FILTER (WHERE ot.is_required)::int AS required_tasks,
    COUNT(cot.id) FILTER (WHERE ot.is_required AND cot.status = 'Completed')::int AS completed_required_tasks,
    COUNT(cot.id) FILTER (WHERE cot.due_at IS NOT NULL AND cot.due_at < NOW() AND cot.status <> 'Completed')::int AS overdue_tasks
FROM hrms.candidate_onboardings co
JOIN hrms.candidates c ON c.tenant_id = co.tenant_id AND c.id = co.candidate_id AND NOT c.inactive
JOIN hrms.onboarding_workflows ow ON ow.tenant_id = co.tenant_id AND ow.id = co.workflow_id AND NOT ow.inactive
LEFT JOIN hrms.candidate_onboarding_tasks cot ON cot.tenant_id = co.tenant_id AND cot.candidate_onboarding_id = co.id AND NOT cot.inactive
LEFT JOIN hrms.onboarding_tasks ot ON ot.tenant_id = co.tenant_id AND ot.id = cot.onboarding_task_id AND NOT ot.inactive
WHERE co.tenant_id = $1 AND co.candidate_id = $2 AND NOT co.inactive
GROUP BY co.id, c.firstname, c.lastname, c.email, ow.name;

-- name: GetDefaultOnboardingWorkflow :one
SELECT * FROM hrms.onboarding_workflows
WHERE tenant_id = $1 AND is_default AND is_active AND NOT inactive
ORDER BY updated_at DESC
LIMIT 1;

-- name: ResolveOnboardingWorkflowForCandidate :one
SELECT ow.*
FROM hrms.candidate_applications ca
JOIN hrms.job_postings jp ON jp.tenant_id = ca.tenant_id AND jp.id = ca.job_posting_id AND NOT jp.inactive
JOIN hrms.onboarding_workflow_assignments owa ON owa.tenant_id = ca.tenant_id AND NOT owa.inactive
JOIN hrms.onboarding_workflows ow ON ow.tenant_id = owa.tenant_id AND ow.id = owa.workflow_id AND ow.is_active AND NOT ow.inactive
WHERE ca.tenant_id = $1
  AND ca.candidate_id = $2
  AND NOT ca.inactive
  AND ca.status IN ('Offered','Hired')
  AND (owa.job_posting_id IS NULL OR owa.job_posting_id = ca.job_posting_id)
  AND (owa.job_position_id IS NULL OR owa.job_position_id = jp.job_position_id)
  AND (owa.department_id IS NULL OR owa.department_id = jp.department_id)
  AND (owa.employment_type_id IS NULL OR owa.employment_type_id = jp.employment_type_id)
ORDER BY
  CASE WHEN owa.job_posting_id = ca.job_posting_id THEN 0 ELSE 1 END,
  CASE WHEN owa.job_position_id = jp.job_position_id THEN 0 ELSE 1 END,
  CASE WHEN owa.department_id = jp.department_id THEN 0 ELSE 1 END,
  CASE WHEN owa.employment_type_id = jp.employment_type_id THEN 0 ELSE 1 END,
  owa.priority ASC,
  ca.updated_at DESC
LIMIT 1;

-- name: CreateCandidateOnboarding :one
INSERT INTO hrms.candidate_onboardings (
    tenant_id, candidate_id, workflow_id, onboarding_status, progress_percentage, started_at, created_by
)
VALUES ($1, $2, $3, $4, 0, NOW(), $5)
RETURNING *;

-- name: CreateCandidateOnboardingTasksFromWorkflow :many
INSERT INTO hrms.candidate_onboarding_tasks (
    tenant_id, candidate_onboarding_id, onboarding_task_id, status, due_at, created_by
)
SELECT
    $1,
    $2,
    ot.id,
    'Pending',
    NOW() + make_interval(days => ot.due_days),
    $3
FROM hrms.onboarding_tasks ot
WHERE ot.tenant_id = $1
  AND ot.workflow_id = $4
  AND NOT ot.inactive
ORDER BY ot.sort_order ASC, ot.title ASC
ON CONFLICT (tenant_id, candidate_onboarding_id, onboarding_task_id) WHERE NOT inactive DO NOTHING
RETURNING *;

-- name: RecalculateCandidateOnboardingProgress :one
WITH totals AS (
    SELECT
        COUNT(cot.id)::int AS total_tasks,
        COUNT(cot.id) FILTER (WHERE cot.status = 'Completed')::int AS completed_tasks,
        COUNT(cot.id) FILTER (WHERE ot.is_required)::int AS required_tasks,
        COUNT(cot.id) FILTER (WHERE ot.is_required AND cot.status = 'Completed')::int AS completed_required_tasks
    FROM hrms.candidate_onboarding_tasks cot
    JOIN hrms.onboarding_tasks ot ON ot.tenant_id = cot.tenant_id AND ot.id = cot.onboarding_task_id AND NOT ot.inactive
    WHERE cot.tenant_id = $1
      AND cot.candidate_onboarding_id = $2
      AND NOT cot.inactive
), next_values AS (
    SELECT
        CASE WHEN total_tasks = 0 THEN 0 ELSE ROUND((completed_tasks::numeric / total_tasks::numeric) * 100)::int END AS progress_percentage,
        CASE
            WHEN required_tasks > 0 AND completed_required_tasks = required_tasks THEN 'Completed'
            WHEN completed_tasks > 0 THEN 'InProgress'
            ELSE 'InProgress'
        END AS onboarding_status
    FROM totals
)
UPDATE hrms.candidate_onboardings co
SET progress_percentage = nv.progress_percentage,
    onboarding_status = nv.onboarding_status,
    completed_at = CASE WHEN nv.onboarding_status = 'Completed' THEN COALESCE(co.completed_at, NOW()) ELSE NULL END,
    updated_at = NOW(),
    updated_by = $3
FROM next_values nv
WHERE co.tenant_id = $1 AND co.id = $2 AND NOT co.inactive
RETURNING co.*;

-- name: SoftDeleteCandidateOnboarding :exec
UPDATE hrms.candidate_onboardings
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListCandidateOnboardingTasks :many
SELECT
    cot.*,
    ot.title AS task_title,
    ot.description AS task_description,
    ot.due_days AS task_due_days,
    ot.is_required AS task_is_required,
    ot.sort_order AS task_sort_order,
    (cot.due_at IS NOT NULL AND cot.due_at < NOW() AND cot.status <> 'Completed')::boolean AS is_overdue
FROM hrms.candidate_onboarding_tasks cot
JOIN hrms.onboarding_tasks ot ON ot.tenant_id = cot.tenant_id AND ot.id = cot.onboarding_task_id AND NOT ot.inactive
WHERE cot.tenant_id = $1 AND cot.candidate_onboarding_id = $2 AND NOT cot.inactive
ORDER BY ot.sort_order ASC, ot.title ASC;

-- name: GetCandidateOnboardingTask :one
SELECT
    cot.*,
    ot.title AS task_title,
    ot.description AS task_description,
    ot.due_days AS task_due_days,
    ot.is_required AS task_is_required,
    ot.sort_order AS task_sort_order,
    (cot.due_at IS NOT NULL AND cot.due_at < NOW() AND cot.status <> 'Completed')::boolean AS is_overdue
FROM hrms.candidate_onboarding_tasks cot
JOIN hrms.onboarding_tasks ot ON ot.tenant_id = cot.tenant_id AND ot.id = cot.onboarding_task_id AND NOT ot.inactive
WHERE cot.tenant_id = $1 AND cot.id = $2 AND NOT cot.inactive;

-- name: UpdateCandidateOnboardingTaskStatus :one
UPDATE hrms.candidate_onboarding_tasks
SET status = $3,
    started_at = CASE WHEN $3 = 'InProgress' AND started_at IS NULL THEN NOW() ELSE started_at END,
    completed_at = CASE WHEN $3 = 'Completed' THEN NOW() ELSE NULL END,
    completed_by = CASE WHEN $3 = 'Completed' THEN $5 ELSE NULL END,
    remarks = $4,
    updated_at = NOW(),
    updated_by = $5
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: CreateCandidateOnboardingEvent :one
INSERT INTO hrms.candidate_onboarding_events (
    tenant_id, candidate_onboarding_id, candidate_onboarding_task_id, action, from_status, to_status, remarks, metadata, created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, COALESCE($8, '{}'::jsonb), $9)
RETURNING *;

-- name: ListCandidateOnboardingEvents :many
SELECT * FROM hrms.candidate_onboarding_events
WHERE tenant_id = $1 AND candidate_onboarding_id = $2 AND NOT inactive
ORDER BY created_at DESC;

-- name: SoftDeleteCandidateOnboardingTask :exec
UPDATE hrms.candidate_onboarding_tasks
SET inactive = TRUE, updated_at = NOW(), updated_by = $3
WHERE tenant_id = $1 AND id = $2;
