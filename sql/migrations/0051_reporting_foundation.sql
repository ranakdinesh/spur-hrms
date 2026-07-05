CREATE TABLE IF NOT EXISTS hrms.report_catalog (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    report_code         VARCHAR(120) NOT NULL,
    module              VARCHAR(80) NOT NULL,
    name                VARCHAR(180) NOT NULL,
    description         TEXT,
    category            VARCHAR(100) NOT NULL,
    scope               VARCHAR(30) NOT NULL DEFAULT 'tenant',
    permission_key      VARCHAR(160) NOT NULL,
    default_filters     JSONB NOT NULL DEFAULT '{}'::jsonb,
    supported_filters   JSONB NOT NULL DEFAULT '[]'::jsonb,
    output_columns      JSONB NOT NULL DEFAULT '[]'::jsonb,
    drilldown_contract  JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_system           BOOLEAN NOT NULL DEFAULT TRUE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order          INTEGER NOT NULL DEFAULT 100,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT report_catalog_code_check CHECK (length(trim(report_code)) > 0),
    CONSTRAINT report_catalog_module_check CHECK (length(trim(module)) > 0),
    CONSTRAINT report_catalog_name_check CHECK (length(trim(name)) > 0),
    CONSTRAINT report_catalog_scope_check CHECK (scope IN ('self','team','tenant','system'))
);

CREATE UNIQUE INDEX IF NOT EXISTS report_catalog_tenant_code_idx
    ON hrms.report_catalog(tenant_id, report_code)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS report_catalog_tenant_module_idx
    ON hrms.report_catalog(tenant_id, module, category, sort_order)
    WHERE NOT inactive AND is_active;

CREATE TABLE IF NOT EXISTS hrms.report_saved_views (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    report_id           UUID NOT NULL REFERENCES hrms.report_catalog(id) ON DELETE CASCADE,
    name                VARCHAR(180) NOT NULL,
    description         TEXT,
    visibility          VARCHAR(30) NOT NULL DEFAULT 'private',
    filters             JSONB NOT NULL DEFAULT '{}'::jsonb,
    columns             JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_favorite         BOOLEAN NOT NULL DEFAULT FALSE,
    owner_user_id       UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT report_saved_views_name_check CHECK (length(trim(name)) > 0),
    CONSTRAINT report_saved_views_visibility_check CHECK (visibility IN ('private','team','tenant'))
);

CREATE INDEX IF NOT EXISTS report_saved_views_report_idx
    ON hrms.report_saved_views(tenant_id, report_id, visibility, created_at DESC)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.report_export_jobs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    report_id           UUID NOT NULL REFERENCES hrms.report_catalog(id) ON DELETE CASCADE,
    saved_view_id       UUID REFERENCES hrms.report_saved_views(id) ON DELETE SET NULL,
    export_format       VARCHAR(20) NOT NULL,
    status              VARCHAR(30) NOT NULL DEFAULT 'queued',
    filters             JSONB NOT NULL DEFAULT '{}'::jsonb,
    file_object_key     TEXT,
    error_message       TEXT,
    requested_by        UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    started_at          TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT report_export_jobs_format_check CHECK (export_format IN ('csv','pdf','xlsx')),
    CONSTRAINT report_export_jobs_status_check CHECK (status IN ('queued','running','completed','failed'))
);

CREATE INDEX IF NOT EXISTS report_export_jobs_tenant_idx
    ON hrms.report_export_jobs(tenant_id, created_at DESC)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.report_schedules (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    report_id           UUID NOT NULL REFERENCES hrms.report_catalog(id) ON DELETE CASCADE,
    saved_view_id       UUID REFERENCES hrms.report_saved_views(id) ON DELETE SET NULL,
    name                VARCHAR(180) NOT NULL,
    frequency           VARCHAR(30) NOT NULL,
    timezone            VARCHAR(80) NOT NULL DEFAULT 'Asia/Kolkata',
    delivery_channels   JSONB NOT NULL DEFAULT '["email"]'::jsonb,
    recipient_user_ids  JSONB NOT NULL DEFAULT '[]'::jsonb,
    recipient_emails    JSONB NOT NULL DEFAULT '[]'::jsonb,
    next_run_at         TIMESTAMPTZ,
    last_run_at         TIMESTAMPTZ,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT report_schedules_name_check CHECK (length(trim(name)) > 0),
    CONSTRAINT report_schedules_frequency_check CHECK (frequency IN ('daily','weekly','monthly'))
);

CREATE INDEX IF NOT EXISTS report_schedules_tenant_idx
    ON hrms.report_schedules(tenant_id, is_active, next_run_at)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.report_snapshots (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    report_id           UUID NOT NULL REFERENCES hrms.report_catalog(id) ON DELETE CASCADE,
    saved_view_id       UUID REFERENCES hrms.report_saved_views(id) ON DELETE SET NULL,
    snapshot_key        VARCHAR(180) NOT NULL,
    period_start        DATE NOT NULL,
    period_end          DATE NOT NULL,
    filters             JSONB NOT NULL DEFAULT '{}'::jsonb,
    summary             JSONB NOT NULL DEFAULT '{}'::jsonb,
    row_count           INTEGER NOT NULL DEFAULT 0,
    generated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    generated_by        UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    inactive            BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT report_snapshots_period_check CHECK (period_end >= period_start),
    CONSTRAINT report_snapshots_count_check CHECK (row_count >= 0),
    CONSTRAINT report_snapshots_key_check CHECK (length(trim(snapshot_key)) > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS report_snapshots_tenant_key_idx
    ON hrms.report_snapshots(tenant_id, snapshot_key)
    WHERE NOT inactive;

ALTER TABLE hrms.report_catalog ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.report_saved_views ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.report_export_jobs ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.report_schedules ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.report_snapshots ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE
    table_name text;
BEGIN
    FOREACH table_name IN ARRAY ARRAY[
        'report_catalog',
        'report_saved_views',
        'report_export_jobs',
        'report_schedules',
        'report_snapshots'
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
