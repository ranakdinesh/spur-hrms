ALTER TABLE hrms.financial_years
    ADD COLUMN IF NOT EXISTS payroll_year BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS leave_year BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS holiday_year BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS reporting_year BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS is_locked BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS locked_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS locked_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS close_note TEXT;

UPDATE hrms.financial_years
SET name = CONCAT('FY ', EXTRACT(YEAR FROM start_date)::INT, '-', EXTRACT(YEAR FROM end_date)::INT)
WHERE name IS NULL OR BTRIM(name) = '';

ALTER TABLE hrms.financial_years
    ALTER COLUMN name SET NOT NULL;

CREATE INDEX IF NOT EXISTS financial_years_tenant_period_idx
    ON hrms.financial_years(tenant_id, start_date, end_date)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS financial_years_tenant_locked_idx
    ON hrms.financial_years(tenant_id, is_locked)
    WHERE NOT inactive;
