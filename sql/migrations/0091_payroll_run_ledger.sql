CREATE TABLE IF NOT EXISTS hrms.pay_run_inputs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    pay_run_id          UUID NOT NULL REFERENCES hrms.pay_runs(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    input_type          VARCHAR(50) NOT NULL,
    source_type         VARCHAR(50) NOT NULL,
    source_id           UUID,
    description         TEXT NOT NULL,
    quantity            NUMERIC(12, 2),
    amount              NUMERIC(14, 2),
    metadata            JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT pay_run_inputs_type_check CHECK (input_type IN ('attendance','leave','salary','statutory','adjustment','claim','overtime','system')),
    CONSTRAINT pay_run_inputs_source_check CHECK (length(trim(source_type)) > 0),
    CONSTRAINT pay_run_inputs_description_check CHECK (length(trim(description)) > 0)
);

CREATE INDEX IF NOT EXISTS pay_run_inputs_run_idx
    ON hrms.pay_run_inputs(tenant_id, pay_run_id, user_id, input_type)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.pay_run_components (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    pay_run_id          UUID NOT NULL REFERENCES hrms.pay_runs(id) ON DELETE CASCADE,
    user_id             UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    component_type      VARCHAR(40) NOT NULL,
    code                VARCHAR(80) NOT NULL,
    name                VARCHAR(180) NOT NULL,
    amount              NUMERIC(14, 2) NOT NULL DEFAULT 0,
    source_input_id     UUID REFERENCES hrms.pay_run_inputs(id) ON DELETE SET NULL,
    salary_template_id  UUID REFERENCES hrms.salary_templates(id) ON DELETE SET NULL,
    taxable             BOOLEAN NOT NULL DEFAULT FALSE,
    statutory           BOOLEAN NOT NULL DEFAULT FALSE,
    employer_cost       BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order          INTEGER NOT NULL DEFAULT 0,
    metadata            JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT pay_run_components_type_check CHECK (component_type IN ('earning','deduction','employer_contribution','reimbursement')),
    CONSTRAINT pay_run_components_code_check CHECK (length(trim(code)) > 0),
    CONSTRAINT pay_run_components_name_check CHECK (length(trim(name)) > 0)
);

CREATE INDEX IF NOT EXISTS pay_run_components_run_idx
    ON hrms.pay_run_components(tenant_id, pay_run_id, user_id, component_type, sort_order)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS pay_run_components_code_idx
    ON hrms.pay_run_components(tenant_id, pay_run_id, code)
    WHERE NOT inactive;

ALTER TABLE hrms.pay_run_inputs ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.pay_run_components ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE
    table_name text;
BEGIN
    FOREACH table_name IN ARRAY ARRAY[
        'pay_run_inputs',
        'pay_run_components'
    ]
    LOOP
        EXECUTE format('DROP POLICY IF EXISTS tenant_isolation_%I ON hrms.%I', table_name, table_name);
        EXECUTE format(
            'CREATE POLICY tenant_isolation_%I ON hrms.%I USING (current_setting(''app.is_super_admin'', true) = ''true'' OR tenant_id::text = current_setting(''app.tenant_id'', true)) WITH CHECK (current_setting(''app.is_super_admin'', true) = ''true'' OR tenant_id::text = current_setting(''app.tenant_id'', true))',
            table_name,
            table_name
        );
    END LOOP;
END $$;
