CREATE TABLE IF NOT EXISTS hrms.employee_credential_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    employee_id UUID NOT NULL REFERENCES hrms.employees(id),
    user_id UUID NOT NULL,
    event_type TEXT NOT NULL CHECK (event_type IN ('resend_credentials', 'reset_temporary_password')),
    delivery_channel TEXT NOT NULL DEFAULT 'email',
    delivery_target TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL CHECK (status IN ('sent', 'failed')),
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID
);

CREATE INDEX IF NOT EXISTS idx_employee_credential_events_tenant_employee
    ON hrms.employee_credential_events(tenant_id, employee_id, created_at DESC);

ALTER TABLE hrms.employee_credential_events ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_policies
        WHERE schemaname = 'hrms' AND tablename = 'employee_credential_events' AND policyname = 'employee_credential_events_tenant_isolation'
    ) THEN
        CREATE POLICY employee_credential_events_tenant_isolation ON hrms.employee_credential_events
            USING (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true) = 'true')
            WITH CHECK (tenant_id = current_setting('app.tenant_id', true)::uuid OR current_setting('app.is_super_admin', true) = 'true');
    END IF;
END $$;
