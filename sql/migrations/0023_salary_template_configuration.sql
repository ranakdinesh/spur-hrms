-- HRMS-035: richer salary templates for payroll setup.
-- Keeps existing salary_templates and salary_template_items compatible while adding
-- enough configuration for India-style CTC, statutory and variable pay components.

ALTER TABLE hrms.salary_templates
    ADD COLUMN IF NOT EXISTS code VARCHAR(60),
    ADD COLUMN IF NOT EXISTS template_type VARCHAR(40) NOT NULL DEFAULT 'ctc',
    ADD COLUMN IF NOT EXISTS applies_to VARCHAR(40) NOT NULL DEFAULT 'all',
    ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3) NOT NULL DEFAULT 'INR',
    ADD COLUMN IF NOT EXISTS effective_from DATE,
    ADD COLUMN IF NOT EXISTS effective_to DATE,
    ADD COLUMN IF NOT EXISTS notes TEXT;

UPDATE hrms.salary_templates
SET code = LOWER(REGEXP_REPLACE(COALESCE(code, name), '[^a-zA-Z0-9]+', '_', 'g'))
WHERE code IS NULL;

ALTER TABLE hrms.salary_templates
    ALTER COLUMN code SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS salary_templates_tenant_code_idx
    ON hrms.salary_templates(tenant_id, LOWER(code)) WHERE NOT inactive;

ALTER TABLE hrms.salary_templates
    DROP CONSTRAINT IF EXISTS salary_templates_configuration_check;

ALTER TABLE hrms.salary_templates
    ADD CONSTRAINT salary_templates_configuration_check CHECK (
        template_type IN ('ctc','gross','net','allowance','deduction')
        AND applies_to IN ('all','grade','department','designation','employee_type','custom')
        AND currency_code ~ '^[A-Z]{3}$'
        AND (effective_to IS NULL OR effective_from IS NULL OR effective_to >= effective_from)
    );

ALTER TABLE hrms.salary_template_items
    ADD COLUMN IF NOT EXISTS calculation_mode VARCHAR(40) NOT NULL DEFAULT 'fixed',
    ADD COLUMN IF NOT EXISTS calculation_base VARCHAR(60) NOT NULL DEFAULT 'gross',
    ADD COLUMN IF NOT EXISTS formula TEXT,
    ADD COLUMN IF NOT EXISTS contribution_side VARCHAR(40) NOT NULL DEFAULT 'employee',
    ADD COLUMN IF NOT EXISTS is_statutory BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_variable BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS affects_gross BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS affects_net BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS cap_amount NUMERIC(12,2),
    ADD COLUMN IF NOT EXISTS min_amount NUMERIC(12,2),
    ADD COLUMN IF NOT EXISTS max_amount NUMERIC(12,2);

CREATE UNIQUE INDEX IF NOT EXISTS salary_template_items_template_code_idx
    ON hrms.salary_template_items(template_id, LOWER(code)) WHERE NOT inactive;

ALTER TABLE hrms.salary_template_items
    DROP CONSTRAINT IF EXISTS salary_template_items_configuration_check;

ALTER TABLE hrms.salary_template_items
    ADD CONSTRAINT salary_template_items_configuration_check CHECK (
        item_type IN ('earning','deduction','employer_contribution','reimbursement')
        AND calculation_mode IN ('fixed','percentage','formula','manual')
        AND calculation_base IN ('ctc','gross','basic','taxable','net','custom')
        AND contribution_side IN ('employee','employer','none')
        AND (calculation_mode <> 'percentage' OR percentage IS NOT NULL)
        AND (calculation_mode <> 'fixed' OR amount IS NOT NULL)
        AND (calculation_mode <> 'formula' OR formula IS NOT NULL)
        AND (cap_amount IS NULL OR cap_amount >= 0)
        AND (min_amount IS NULL OR min_amount >= 0)
        AND (max_amount IS NULL OR max_amount >= 0)
        AND (max_amount IS NULL OR min_amount IS NULL OR max_amount >= min_amount)
    );
