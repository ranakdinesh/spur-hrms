ALTER TABLE hrms.employee_statutory
    ADD COLUMN IF NOT EXISTS lwf_applicable BOOLEAN NOT NULL DEFAULT FALSE;
