CREATE TABLE IF NOT EXISTS hrms.salary_slip_formats (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                       UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    title                           VARCHAR(160) NOT NULL DEFAULT 'Salary Slip',
    subtitle                        VARCHAR(255),
    logo_path                       TEXT,
    primary_color                   VARCHAR(20) NOT NULL DEFAULT '#111827',
    accent_color                    VARCHAR(20) NOT NULL DEFAULT '#588368',
    show_leave_balance              BOOLEAN NOT NULL DEFAULT TRUE,
    show_ytd_summary                BOOLEAN NOT NULL DEFAULT FALSE,
    show_employee_bank              BOOLEAN NOT NULL DEFAULT TRUE,
    show_employer_contributions     BOOLEAN NOT NULL DEFAULT FALSE,
    footer_text                     TEXT,
    custom_fields                   JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive                        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by                      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by                      UUID REFERENCES auth.users(id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS salary_slip_formats_tenant_active_idx
ON hrms.salary_slip_formats(tenant_id)
WHERE NOT inactive;
