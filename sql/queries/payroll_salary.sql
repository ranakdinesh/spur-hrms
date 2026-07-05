-- name: ListSalaryTemplates :many
SELECT * FROM hrms.salary_templates
WHERE tenant_id = $1 AND NOT inactive
ORDER BY is_active DESC, name ASC;

-- name: ListSalaryTemplatesByFY :many
SELECT * FROM hrms.salary_templates
WHERE tenant_id = $1 AND fy_id = $2 AND NOT inactive
ORDER BY is_active DESC, name ASC;

-- name: GetSalaryTemplate :one
SELECT * FROM hrms.salary_templates
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetActiveSalaryTemplate :one
SELECT * FROM hrms.salary_templates
WHERE tenant_id = $1 AND fy_id = $2 AND is_active AND NOT inactive;

-- name: CreateSalaryTemplate :one
INSERT INTO hrms.salary_templates (
    tenant_id, fy_id, code, name, description, template_type, applies_to,
    currency_code, effective_from, effective_to, notes, is_active, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, FALSE, $12, $12)
RETURNING *;

-- name: UpdateSalaryTemplate :one
UPDATE hrms.salary_templates
SET fy_id = $3,
    code = $4,
    name = $5,
    description = $6,
    template_type = $7,
    applies_to = $8,
    currency_code = $9,
    effective_from = $10,
    effective_to = $11,
    notes = $12,
    updated_by = $13,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: DeactivateSalaryTemplatesForFY :exec
UPDATE hrms.salary_templates
SET is_active = FALSE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND fy_id = $2 AND NOT inactive;

-- name: ActivateSalaryTemplate :one
UPDATE hrms.salary_templates
SET is_active = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteSalaryTemplate :exec
UPDATE hrms.salary_templates
SET inactive = TRUE, is_active = FALSE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: ListSalaryTemplateItems :many
SELECT * FROM hrms.salary_template_items
WHERE tenant_id = $1 AND template_id = $2 AND NOT inactive
ORDER BY sort_order ASC, name ASC;

-- name: GetSalaryTemplateItem :one
SELECT * FROM hrms.salary_template_items
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: CreateSalaryTemplateItem :one
INSERT INTO hrms.salary_template_items (
    tenant_id, template_id, item_type, code, name, percentage, amount,
    calculation_mode, calculation_base, formula, contribution_side,
    is_tax_exempt, is_statutory, is_variable, affects_gross, affects_net,
    cap_amount, min_amount, max_amount, sort_order, created_by, updated_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $21
)
RETURNING *;

-- name: UpdateSalaryTemplateItem :one
UPDATE hrms.salary_template_items
SET item_type = $4,
    code = $5,
    name = $6,
    percentage = $7,
    amount = $8,
    calculation_mode = $9,
    calculation_base = $10,
    formula = $11,
    contribution_side = $12,
    is_tax_exempt = $13,
    is_statutory = $14,
    is_variable = $15,
    affects_gross = $16,
    affects_net = $17,
    cap_amount = $18,
    min_amount = $19,
    max_amount = $20,
    sort_order = $21,
    updated_by = $22,
    updated_at = NOW()
WHERE tenant_id = $1 AND template_id = $2 AND id = $3 AND NOT inactive
RETURNING *;

-- name: SoftDeleteSalaryTemplateItem :exec
UPDATE hrms.salary_template_items
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: ListEmployeeSalariesByUser :many
SELECT * FROM hrms.employee_salaries
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY effective_from DESC NULLS LAST, created_at DESC;

-- name: CreateEmployeeSalary :one
INSERT INTO hrms.employee_salaries (
    tenant_id, user_id, fy_id, template_id, gross_salary, effective_from, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
RETURNING *;

-- name: GetEmployeeSalary :one
SELECT * FROM hrms.employee_salaries
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: UpdateEmployeeSalary :one
UPDATE hrms.employee_salaries
SET template_id = $4,
    gross_salary = $5,
    effective_from = $6,
    updated_by = $7,
    updated_at = NOW()
WHERE tenant_id = $1 AND user_id = $2 AND fy_id = $3 AND NOT inactive
RETURNING *;

-- name: SoftDeleteEmployeeSalary :exec
UPDATE hrms.employee_salaries
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: SoftDeleteEmployeeSalariesByUserFY :exec
UPDATE hrms.employee_salaries
SET inactive = TRUE, updated_by = $4, updated_at = NOW()
WHERE tenant_id = $1 AND user_id = $2 AND fy_id = $3 AND NOT inactive;

-- name: ListEmployeeSalaryStructures :many
SELECT * FROM hrms.employee_salary_structures
WHERE tenant_id = $1 AND user_id = $2 AND fy_id = $3 AND NOT inactive
ORDER BY sort_order ASC, name ASC;

-- name: CreateEmployeeSalaryStructure :one
INSERT INTO hrms.employee_salary_structures (
    tenant_id, user_id, template_id, fy_id, item_type, code, name, amount, sort_order, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)
RETURNING *;

-- name: GetEmployeeSalaryStructure :one
SELECT * FROM hrms.employee_salary_structures
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: SoftDeleteEmployeeSalaryStructure :exec
UPDATE hrms.employee_salary_structures
SET inactive = TRUE, updated_by = $3, updated_at = NOW()
WHERE tenant_id = $1 AND id = $2;

-- name: SoftDeleteEmployeeSalaryStructuresByUserFY :exec
UPDATE hrms.employee_salary_structures
SET inactive = TRUE, updated_by = $4, updated_at = NOW()
WHERE tenant_id = $1 AND user_id = $2 AND fy_id = $3 AND NOT inactive;

-- name: ListSalarySlipsByUser :many
SELECT * FROM hrms.salary_slips
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY year DESC, month DESC;

-- name: ListSalarySlipsByTenantPeriod :many
SELECT * FROM hrms.salary_slips
WHERE tenant_id = $1 AND month = $2 AND year = $3 AND NOT inactive
ORDER BY created_at DESC;

-- name: ListRecentSalarySlipsByUser :many
SELECT * FROM hrms.salary_slips
WHERE tenant_id = $1 AND user_id = $2 AND NOT inactive
ORDER BY year DESC, month DESC
LIMIT $3;

-- name: GetSalarySlip :one
SELECT * FROM hrms.salary_slips
WHERE tenant_id = $1 AND id = $2 AND NOT inactive;

-- name: GetSalarySlipByPeriod :one
SELECT * FROM hrms.salary_slips
WHERE tenant_id = $1 AND user_id = $2 AND month = $3 AND year = $4 AND NOT inactive;

-- name: CreateSalarySlip :one
INSERT INTO hrms.salary_slips (
    tenant_id, user_id, fy_id, template_id, month, year, gross_salary,
    total_earnings, total_deductions, absent_deduction, net_salary,
    absent_days, present_days, total_days, lwp_days, no_of_ph_weo,
    is_special, is_regenerated, pdf_path, created_by, updated_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11,
    $12, $13, $14, $15, $16,
    $17, $18, $19, $20, $20
)
RETURNING *;

-- name: UpdateSalarySlipByID :one
UPDATE hrms.salary_slips
SET fy_id = $3,
    template_id = $4,
    gross_salary = $5,
    total_earnings = $6,
    total_deductions = $7,
    absent_deduction = $8,
    net_salary = $9,
    absent_days = $10,
    present_days = $11,
    total_days = $12,
    lwp_days = $13,
    no_of_ph_weo = $14,
    is_special = $15,
    is_regenerated = $16,
    pdf_path = $17,
    updated_by = $18,
    updated_at = NOW()
WHERE tenant_id = $1 AND id = $2 AND NOT inactive
RETURNING *;

-- name: SoftDeleteSalarySlip :exec
UPDATE hrms.salary_slips
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND id = $2;

-- name: ListSalarySlipItems :many
SELECT * FROM hrms.salary_slip_items
WHERE tenant_id = $1 AND slip_id = $2 AND NOT inactive
ORDER BY sort_order ASC, name ASC;

-- name: CreateSalarySlipItem :one
INSERT INTO hrms.salary_slip_items (
    tenant_id, slip_id, item_type, code, name, amount, sort_order, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
RETURNING *;

-- name: SoftDeleteSalarySlipItemsBySlip :exec
UPDATE hrms.salary_slip_items
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND slip_id = $2;

-- name: ListSalarySlipLeaves :many
SELECT * FROM hrms.salary_slip_leaves
WHERE tenant_id = $1 AND slip_id = $2 AND NOT inactive
ORDER BY leave_type_name ASC NULLS LAST;

-- name: CreateSalarySlipLeave :one
INSERT INTO hrms.salary_slip_leaves (
    tenant_id, slip_id, leave_type_id, leave_type_name, total_days, used_days, balance_days, created_by, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
RETURNING *;

-- name: SoftDeleteSalarySlipLeavesBySlip :exec
UPDATE hrms.salary_slip_leaves
SET inactive = TRUE, updated_by = $3
WHERE tenant_id = $1 AND slip_id = $2;

-- name: GetSalarySlipFormat :one
SELECT * FROM hrms.salary_slip_formats
WHERE tenant_id = $1 AND NOT inactive;

-- name: UpsertSalarySlipFormat :one
INSERT INTO hrms.salary_slip_formats (
    tenant_id, title, subtitle, logo_path, primary_color, accent_color,
    show_leave_balance, show_ytd_summary, show_employee_bank, show_employer_contributions,
    footer_text, custom_fields, updated_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (tenant_id) WHERE NOT inactive DO UPDATE
SET title = EXCLUDED.title,
    subtitle = EXCLUDED.subtitle,
    logo_path = EXCLUDED.logo_path,
    primary_color = EXCLUDED.primary_color,
    accent_color = EXCLUDED.accent_color,
    show_leave_balance = EXCLUDED.show_leave_balance,
    show_ytd_summary = EXCLUDED.show_ytd_summary,
    show_employee_bank = EXCLUDED.show_employee_bank,
    show_employer_contributions = EXCLUDED.show_employer_contributions,
    footer_text = EXCLUDED.footer_text,
    custom_fields = EXCLUDED.custom_fields,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;
