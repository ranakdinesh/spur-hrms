ALTER TABLE hrms.employees
    ADD COLUMN IF NOT EXISTS probation_status VARCHAR(30) NOT NULL DEFAULT 'confirmed',
    ADD COLUMN IF NOT EXISTS probation_start_date DATE,
    ADD COLUMN IF NOT EXISTS probation_end_date DATE,
    ADD COLUMN IF NOT EXISTS probation_duration_days INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS probation_confirmed_at DATE,
    ADD COLUMN IF NOT EXISTS is_payroll_staff BOOLEAN NOT NULL DEFAULT FALSE;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'employees_probation_status_check'
          AND conrelid = 'hrms.employees'::regclass
    ) THEN
        ALTER TABLE hrms.employees
            ADD CONSTRAINT employees_probation_status_check
            CHECK (probation_status IN ('not_applicable', 'probation', 'confirmed', 'extended'));
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'employees_probation_dates_check'
          AND conrelid = 'hrms.employees'::regclass
    ) THEN
        ALTER TABLE hrms.employees
            ADD CONSTRAINT employees_probation_dates_check
            CHECK (probation_end_date IS NULL OR probation_start_date IS NULL OR probation_end_date >= probation_start_date);
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'employees_probation_duration_check'
          AND conrelid = 'hrms.employees'::regclass
    ) THEN
        ALTER TABLE hrms.employees
            ADD CONSTRAINT employees_probation_duration_check
            CHECK (probation_duration_days >= 0);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS employees_probation_tenant_status_idx
    ON hrms.employees(tenant_id, probation_status, probation_end_date)
    WHERE NOT inactive;
