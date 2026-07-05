CREATE TABLE IF NOT EXISTS hrms.tenant_operation_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operation_number TEXT NOT NULL,
    operation_type TEXT NOT NULL,
    title TEXT NOT NULL,
    target_tenant_id UUID REFERENCES auth.tenants(id) ON DELETE SET NULL,
    target_tenant_name TEXT,
    target_tenant_code TEXT,
    status TEXT NOT NULL DEFAULT 'pending_validation',
    risk_level TEXT NOT NULL DEFAULT 'medium',
    reason TEXT NOT NULL,
    requested_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approved_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,
    completed_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    completed_at TIMESTAMPTZ,
    approval_required BOOLEAN NOT NULL DEFAULT TRUE,
    backup_required BOOLEAN NOT NULL DEFAULT FALSE,
    backup_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    retention_until TIMESTAMPTZ,
    request_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    validation_results JSONB NOT NULL DEFAULT '{}'::jsonb,
    rollback_metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT tenant_operation_requests_number_check CHECK (length(trim(operation_number)) > 0),
    CONSTRAINT tenant_operation_requests_title_check CHECK (length(trim(title)) > 0),
    CONSTRAINT tenant_operation_requests_reason_check CHECK (length(trim(reason)) > 0),
    CONSTRAINT tenant_operation_requests_type_check CHECK (operation_type IN ('create_tenant','suspend_tenant','restore_tenant','schedule_delete_tenant','cancel_delete_tenant','module_enable','module_disable','storage_change','domain_branding_change','admin_reassignment','data_export')),
    CONSTRAINT tenant_operation_requests_status_check CHECK (status IN ('pending_validation','pending_approval','approved','in_progress','completed','rejected','cancelled','failed')),
    CONSTRAINT tenant_operation_requests_risk_check CHECK (risk_level IN ('low','medium','high','critical')),
    CONSTRAINT tenant_operation_requests_payload_check CHECK (jsonb_typeof(request_payload) = 'object'),
    CONSTRAINT tenant_operation_requests_validation_check CHECK (jsonb_typeof(validation_results) = 'object'),
    CONSTRAINT tenant_operation_requests_rollback_check CHECK (jsonb_typeof(rollback_metadata) = 'object'),
    CONSTRAINT tenant_operation_requests_metadata_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS tenant_operation_requests_number_idx ON hrms.tenant_operation_requests(operation_number) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS tenant_operation_requests_status_idx ON hrms.tenant_operation_requests(status, risk_level, created_at DESC) WHERE NOT inactive;
CREATE INDEX IF NOT EXISTS tenant_operation_requests_target_idx ON hrms.tenant_operation_requests(target_tenant_id, operation_type, status) WHERE target_tenant_id IS NOT NULL AND NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.tenant_operation_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL REFERENCES hrms.tenant_operation_requests(id) ON DELETE CASCADE,
    action TEXT NOT NULL,
    from_status TEXT,
    to_status TEXT,
    actor_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    remarks TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT tenant_operation_events_action_check CHECK (length(trim(action)) > 0),
    CONSTRAINT tenant_operation_events_metadata_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS tenant_operation_events_request_idx ON hrms.tenant_operation_events(request_id, created_at DESC) WHERE NOT inactive;

DROP TRIGGER IF EXISTS tenant_operation_requests_set_updated_at ON hrms.tenant_operation_requests;
CREATE TRIGGER tenant_operation_requests_set_updated_at BEFORE UPDATE ON hrms.tenant_operation_requests FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

DROP TRIGGER IF EXISTS tenant_operation_events_set_updated_at ON hrms.tenant_operation_events;
CREATE TRIGGER tenant_operation_events_set_updated_at BEFORE UPDATE ON hrms.tenant_operation_events FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.tenant_operation_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.tenant_operation_events ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE
    table_name TEXT;
BEGIN
    FOREACH table_name IN ARRAY ARRAY['tenant_operation_requests','tenant_operation_events']
    LOOP
        EXECUTE format('DROP POLICY IF EXISTS %I_super_admin_all ON hrms.%I', table_name, table_name);
        EXECUTE format('CREATE POLICY %I_super_admin_all ON hrms.%I FOR ALL USING (current_setting(''app.is_super_admin'', true) = ''true'') WITH CHECK (current_setting(''app.is_super_admin'', true) = ''true'')', table_name, table_name);
    END LOOP;
END $$;
