-- name: ListEmployees :many
SELECT
    e.id,
    e.tenant_id,
    e.user_id,
    e.employee_code,
    e.firstname,
    e.middle_name,
    e.lastname,
    e.email,
    e.mobile,
    e.dob,
    e.gender,
    e.marital_status,
    e.blood_group,
    e.profile_photo_path,
    e.address,
    e.city,
    e.state,
    e.country,
    e.pincode,
    e.emergency_contact,
    e.joining_date,
    e.resignation_date,
    e.department_id,
    d.name AS department_name,
    e.branch_id,
    b.branch_name AS branch_name,
    e.designation_id,
    dg.name AS designation_name,
    COALESCE(dg.attendance_required, TRUE)::boolean AS attendance_required,
    e.reporting_manager_id,
    e.employment_type_id,
    et.name AS employment_type_name,
    e.role,
    e.grade,
    e.experience_year,
    e.experience_month,
    e.probation_status,
    e.probation_start_date,
    e.probation_end_date,
    e.probation_duration_days,
    e.probation_confirmed_at,
    e.is_payroll_staff,
    e.inactive,
    e.created_at,
    e.created_by,
    e.updated_at,
    e.updated_by
FROM hrms.employees e
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
LEFT JOIN hrms.designations dg ON dg.tenant_id = e.tenant_id AND dg.id = e.designation_id AND NOT dg.inactive
LEFT JOIN hrms.employment_types et ON et.tenant_id = e.tenant_id AND et.id = e.employment_type_id AND NOT et.inactive
WHERE e.tenant_id = $1 AND NOT e.inactive
ORDER BY e.firstname ASC, e.lastname ASC NULLS LAST;

-- name: GetEmployeeProfileItem :one
SELECT
    e.id,
    e.tenant_id,
    e.user_id,
    e.employee_code,
    e.firstname,
    e.middle_name,
    e.lastname,
    e.email,
    e.mobile,
    e.dob,
    e.gender,
    e.marital_status,
    e.blood_group,
    e.profile_photo_path,
    e.address,
    e.city,
    e.state,
    e.country,
    e.pincode,
    e.emergency_contact,
    e.joining_date,
    e.resignation_date,
    e.department_id,
    d.name AS department_name,
    e.branch_id,
    b.branch_name AS branch_name,
    e.designation_id,
    dg.name AS designation_name,
    COALESCE(dg.attendance_required, TRUE)::boolean AS attendance_required,
    e.reporting_manager_id,
    e.employment_type_id,
    et.name AS employment_type_name,
    e.role,
    e.grade,
    e.experience_year,
    e.experience_month,
    e.probation_status,
    e.probation_start_date,
    e.probation_end_date,
    e.probation_duration_days,
    e.probation_confirmed_at,
    e.is_payroll_staff,
    e.inactive,
    e.created_at,
    e.created_by,
    e.updated_at,
    e.updated_by
FROM hrms.employees e
LEFT JOIN hrms.departments d ON d.tenant_id = e.tenant_id AND d.id = e.department_id AND NOT d.inactive
LEFT JOIN hrms.branches b ON b.tenant_id = e.tenant_id AND b.id = e.branch_id AND NOT b.inactive
LEFT JOIN hrms.designations dg ON dg.tenant_id = e.tenant_id AND dg.id = e.designation_id AND NOT dg.inactive
LEFT JOIN hrms.employment_types et ON et.tenant_id = e.tenant_id AND et.id = e.employment_type_id AND NOT et.inactive
WHERE e.tenant_id = $1 AND e.id = $2 AND NOT e.inactive;

-- name: GetEmployeeByID :one
SELECT * FROM hrms.employees
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetEmployeeByUserID :one
SELECT * FROM hrms.employees
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive;

-- name: GetEmployeeAttendanceRequired :one
SELECT COALESCE(dg.attendance_required, TRUE)::boolean AS attendance_required
FROM hrms.employees e
LEFT JOIN hrms.designations dg ON dg.tenant_id = e.tenant_id AND dg.id = e.designation_id AND NOT dg.inactive
WHERE e.tenant_id = $1 AND e.user_id = $2 AND NOT e.inactive;

-- name: SoftDeleteEmployee :exec
UPDATE hrms.employees
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateEmployee :one
UPDATE hrms.employees
SET
    employee_code = $3,
    firstname = $4,
    middle_name = $5,
    lastname = $6,
    email = $7,
    mobile = $8,
    dob = $9,
    gender = $10,
    marital_status = $11,
    blood_group = $12,
    profile_photo_path = $13,
    address = $14,
    city = $15,
    state = $16,
    country = $17,
    pincode = $18,
    emergency_contact = $19,
    joining_date = $20,
    resignation_date = $21,
    department_id = $22,
    branch_id = $23,
    designation_id = $24,
    reporting_manager_id = $25,
    employment_type_id = $26,
    role = $27,
    grade = $28,
    experience_year = $29,
    experience_month = $30,
    probation_status = $31,
    probation_start_date = $32,
    probation_end_date = $33,
    probation_duration_days = $34,
    probation_confirmed_at = $35,
    is_payroll_staff = $36,
    updated_by = $37,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: GetEmployeeStatutoryByUserID :one
SELECT * FROM hrms.employee_statutory
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
LIMIT 1;

-- name: UpsertEmployeeStatutory :one
INSERT INTO hrms.employee_statutory (
    tenant_id,
    user_id,
    pf_no,
    uan_no,
    esic_no,
    pan,
    aadhaar,
    pt_applicable,
    pf_applicable,
    esic_applicable,
    lwf_applicable,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12
)
ON CONFLICT (tenant_id, user_id) WHERE NOT inactive
DO UPDATE SET
    pf_no = EXCLUDED.pf_no,
    uan_no = EXCLUDED.uan_no,
    esic_no = EXCLUDED.esic_no,
    pan = EXCLUDED.pan,
    aadhaar = EXCLUDED.aadhaar,
    pt_applicable = EXCLUDED.pt_applicable,
    pf_applicable = EXCLUDED.pf_applicable,
    esic_applicable = EXCLUDED.esic_applicable,
    lwf_applicable = EXCLUDED.lwf_applicable,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: SoftDeleteEmployeeStatutory :exec
UPDATE hrms.employee_statutory
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListEmployeeBanksByUserID :many
SELECT * FROM hrms.employee_banks
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY is_primary DESC, created_at DESC;

-- name: GetEmployeeBank :one
SELECT * FROM hrms.employee_banks
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpsertPrimaryEmployeeBank :one
INSERT INTO hrms.employee_banks (
    tenant_id,
    user_id,
    bank_name,
    account_number,
    ifsc_code,
    account_type,
    branch_name,
    is_primary,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, TRUE, $8, $8
)
ON CONFLICT (tenant_id, user_id) WHERE is_primary AND NOT inactive
DO UPDATE SET
    bank_name = EXCLUDED.bank_name,
    account_number = EXCLUDED.account_number,
    ifsc_code = EXCLUDED.ifsc_code,
    account_type = EXCLUDED.account_type,
    branch_name = EXCLUDED.branch_name,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: SoftDeleteEmployeeBank :exec
UPDATE hrms.employee_banks
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListDocumentTypes :many
SELECT * FROM hrms.document_types
WHERE tenant_id = $1 AND NOT inactive
ORDER BY display_order ASC, name ASC;

-- name: GetDocumentType :one
SELECT * FROM hrms.document_types
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateDocumentType :one
INSERT INTO hrms.document_types (
    tenant_id,
    name,
    description,
    is_required,
    instructions,
    allowed_content_types,
    max_file_size_bytes,
    display_order,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $9
)
RETURNING *;

-- name: UpdateDocumentType :one
UPDATE hrms.document_types
SET
    name = $3,
    description = $4,
    is_required = $5,
    instructions = $6,
    allowed_content_types = $7,
    max_file_size_bytes = $8,
    display_order = $9,
    updated_by = $10,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteDocumentType :exec
UPDATE hrms.document_types
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: ListEmployeeDocumentsByUserID :many
SELECT
    ed.id,
    ed.tenant_id,
    ed.user_id,
    ed.document_type_id,
    dt.name AS document_type_name,
    ed.title,
    ed.file_path,
    ed.status,
    ed.review_remarks,
    ed.reviewed_by,
    ed.reviewed_at,
    ed.original_file_name,
    ed.content_type,
    ed.file_size_bytes,
    ed.encrypted,
    ed.encryption_algorithm,
    ed.inactive,
    ed.created_at,
    ed.created_by,
    ed.updated_at,
    ed.updated_by
FROM hrms.employee_documents ed
LEFT JOIN hrms.document_types dt ON dt.tenant_id = ed.tenant_id AND dt.id = ed.document_type_id AND NOT dt.inactive
WHERE ed.tenant_id = $1 AND ed.user_id = $2 AND NOT ed.inactive
ORDER BY ed.created_at DESC;

-- name: GetEmployeeDocument :one
SELECT * FROM hrms.employee_documents
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateEmployeeDocument :one
INSERT INTO hrms.employee_documents (
    tenant_id,
    user_id,
    document_type_id,
    title,
    file_path,
    status,
    original_file_name,
    content_type,
    file_size_bytes,
    encrypted,
    encryption_algorithm,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: UpdateEmployeeDocument :one
UPDATE hrms.employee_documents
SET
    document_type_id = $3,
    title = $4,
    file_path = $5,
    original_file_name = $6,
    content_type = $7,
    file_size_bytes = $8,
    status = $9,
    review_remarks = NULL,
    reviewed_by = NULL,
    reviewed_at = NULL,
    updated_by = $10,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteEmployeeDocument :exec
UPDATE hrms.employee_documents
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateEmployee :one
INSERT INTO hrms.employees (
    tenant_id,
    user_id,
    employee_code,
    firstname,
    middle_name,
    lastname,
    email,
    mobile,
    dob,
    gender,
    marital_status,
    blood_group,
    profile_photo_path,
    address,
    city,
    state,
    country,
    pincode,
    emergency_contact,
    joining_date,
    department_id,
    branch_id,
    designation_id,
    reporting_manager_id,
    employment_type_id,
    role,
    grade,
    experience_year,
    experience_month,
    probation_status,
    probation_start_date,
    probation_end_date,
    probation_duration_days,
    probation_confirmed_at,
    is_payroll_staff,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
    $31, $32, $33, $34, $35, $36, $36
)
RETURNING *;

-- name: GetEmployeeByCode :one
SELECT * FROM hrms.employees
WHERE tenant_id = $1 AND employee_code = $2 AND NOT inactive;

-- name: EmployeeCodeExists :one
SELECT EXISTS (
    SELECT 1 FROM hrms.employees
    WHERE tenant_id = $1 AND employee_code = $2 AND NOT inactive
)::boolean;

-- name: EmployeeCodeExistsForOtherEmployee :one
SELECT EXISTS (
    SELECT 1 FROM hrms.employees
    WHERE tenant_id = $1 AND employee_code = $2 AND id <> $3 AND NOT inactive
)::boolean;

-- name: CountActiveEmployees :one
SELECT COUNT(*)::int
FROM hrms.employees
WHERE tenant_id = $1 AND NOT inactive;

-- name: ReviewEmployeeDocument :one
UPDATE hrms.employee_documents
SET
    status = $3,
    review_remarks = $4,
    reviewed_by = $5,
    reviewed_at = NOW(),
    updated_by = $5,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;
