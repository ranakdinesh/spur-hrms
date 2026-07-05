CREATE TABLE IF NOT EXISTS hrms.workflow_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    workflow_key VARCHAR(100) NOT NULL,
    name VARCHAR(180) NOT NULL,
    module_key VARCHAR(80) NOT NULL DEFAULT 'hrms',
    description TEXT,
    status VARCHAR(40) NOT NULL DEFAULT 'draft',
    visibility_scope VARCHAR(40) NOT NULL DEFAULT 'tenant',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT workflow_definitions_key_check CHECK (BTRIM(workflow_key) <> ''),
    CONSTRAINT workflow_definitions_name_check CHECK (BTRIM(name) <> ''),
    CONSTRAINT workflow_definitions_status_check CHECK (status IN ('draft','active','paused','archived')),
    CONSTRAINT workflow_definitions_visibility_check CHECK (visibility_scope IN ('tenant','restricted','confidential')),
    CONSTRAINT workflow_definitions_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS workflow_definitions_key_idx
    ON hrms.workflow_definitions(tenant_id, lower(workflow_key))
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.workflow_definition_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    workflow_definition_id UUID NOT NULL REFERENCES hrms.workflow_definitions(id) ON DELETE CASCADE,
    step_order INT NOT NULL,
    step_key VARCHAR(100) NOT NULL,
    name VARCHAR(180) NOT NULL,
    step_type VARCHAR(40) NOT NULL DEFAULT 'approval',
    assignment_type VARCHAR(40) NOT NULL DEFAULT 'role',
    assignment_value VARCHAR(160),
    required BOOLEAN NOT NULL DEFAULT TRUE,
    due_offset_hours INT NOT NULL DEFAULT 24,
    allowed_actions JSONB NOT NULL DEFAULT '["approve","reject","request_info","delegate","comment","complete"]'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT workflow_steps_order_check CHECK (step_order > 0),
    CONSTRAINT workflow_steps_key_check CHECK (BTRIM(step_key) <> ''),
    CONSTRAINT workflow_steps_name_check CHECK (BTRIM(name) <> ''),
    CONSTRAINT workflow_steps_type_check CHECK (step_type IN ('approval','review','action','checklist','notification')),
    CONSTRAINT workflow_steps_assignment_check CHECK (assignment_type IN ('user','role','team','requester','manager','source_owner')),
    CONSTRAINT workflow_steps_due_check CHECK (due_offset_hours >= 0),
    CONSTRAINT workflow_steps_actions_array_check CHECK (jsonb_typeof(allowed_actions) = 'array'),
    CONSTRAINT workflow_steps_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS workflow_steps_key_idx
    ON hrms.workflow_definition_steps(tenant_id, workflow_definition_id, lower(step_key))
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS workflow_steps_order_idx
    ON hrms.workflow_definition_steps(tenant_id, workflow_definition_id, step_order)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.operation_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    template_key VARCHAR(100) NOT NULL,
    name VARCHAR(180) NOT NULL,
    category VARCHAR(80) NOT NULL DEFAULT 'hr_operation',
    source_module VARCHAR(80) NOT NULL DEFAULT 'hrms',
    source_type VARCHAR(100) NOT NULL DEFAULT 'manual',
    workflow_definition_id UUID REFERENCES hrms.workflow_definitions(id) ON DELETE SET NULL,
    default_priority INT NOT NULL DEFAULT 50,
    default_severity VARCHAR(40) NOT NULL DEFAULT 'medium',
    allowed_actions JSONB NOT NULL DEFAULT '["approve","reject","request_info","delegate","comment","complete"]'::jsonb,
    launch_schema JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT operation_templates_key_check CHECK (BTRIM(template_key) <> ''),
    CONSTRAINT operation_templates_name_check CHECK (BTRIM(name) <> ''),
    CONSTRAINT operation_templates_priority_check CHECK (default_priority BETWEEN 1 AND 100),
    CONSTRAINT operation_templates_severity_check CHECK (default_severity IN ('low','medium','high','critical')),
    CONSTRAINT operation_templates_actions_array_check CHECK (jsonb_typeof(allowed_actions) = 'array'),
    CONSTRAINT operation_templates_launch_object_check CHECK (jsonb_typeof(launch_schema) = 'object'),
    CONSTRAINT operation_templates_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS operation_templates_key_idx
    ON hrms.operation_templates(tenant_id, lower(template_key))
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.workflow_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    task_number VARCHAR(40) NOT NULL,
    template_id UUID REFERENCES hrms.operation_templates(id) ON DELETE SET NULL,
    workflow_definition_id UUID REFERENCES hrms.workflow_definitions(id) ON DELETE SET NULL,
    workflow_step_id UUID REFERENCES hrms.workflow_definition_steps(id) ON DELETE SET NULL,
    parent_task_id UUID REFERENCES hrms.workflow_tasks(id) ON DELETE SET NULL,
    source_module VARCHAR(80) NOT NULL DEFAULT 'hrms',
    source_type VARCHAR(100) NOT NULL DEFAULT 'manual',
    source_id UUID,
    source_record_label VARCHAR(220),
    title VARCHAR(220) NOT NULL,
    description TEXT,
    requester_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    assignee_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    assignee_role VARCHAR(120),
    assignee_team VARCHAR(120),
    delegated_from_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    status VARCHAR(40) NOT NULL DEFAULT 'pending',
    priority INT NOT NULL DEFAULT 50,
    severity VARCHAR(40) NOT NULL DEFAULT 'medium',
    visibility_scope VARCHAR(40) NOT NULL DEFAULT 'tenant',
    due_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    completed_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    action_schema JSONB NOT NULL DEFAULT '[]'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT workflow_tasks_number_check CHECK (BTRIM(task_number) <> ''),
    CONSTRAINT workflow_tasks_title_check CHECK (BTRIM(title) <> ''),
    CONSTRAINT workflow_tasks_status_check CHECK (status IN ('pending','in_progress','waiting_info','approved','rejected','completed','cancelled','delegated','blocked')),
    CONSTRAINT workflow_tasks_priority_check CHECK (priority BETWEEN 1 AND 100),
    CONSTRAINT workflow_tasks_severity_check CHECK (severity IN ('low','medium','high','critical')),
    CONSTRAINT workflow_tasks_visibility_check CHECK (visibility_scope IN ('tenant','restricted','confidential')),
    CONSTRAINT workflow_tasks_action_schema_array_check CHECK (jsonb_typeof(action_schema) = 'array'),
    CONSTRAINT workflow_tasks_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE UNIQUE INDEX IF NOT EXISTS workflow_tasks_number_idx
    ON hrms.workflow_tasks(tenant_id, task_number)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS workflow_tasks_inbox_idx
    ON hrms.workflow_tasks(tenant_id, status, assignee_user_id, assignee_role, due_at)
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS workflow_tasks_source_idx
    ON hrms.workflow_tasks(tenant_id, source_module, source_type, source_id)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.workflow_task_watchers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    task_id UUID NOT NULL REFERENCES hrms.workflow_tasks(id) ON DELETE CASCADE,
    watcher_user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    watch_reason VARCHAR(120),
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS workflow_task_watchers_unique_idx
    ON hrms.workflow_task_watchers(tenant_id, task_id, watcher_user_id)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.workflow_task_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    task_id UUID NOT NULL REFERENCES hrms.workflow_tasks(id) ON DELETE CASCADE,
    visibility VARCHAR(40) NOT NULL DEFAULT 'tenant',
    body TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT workflow_task_comments_body_check CHECK (BTRIM(body) <> ''),
    CONSTRAINT workflow_task_comments_visibility_check CHECK (visibility IN ('tenant','restricted','confidential')),
    CONSTRAINT workflow_task_comments_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS workflow_task_comments_task_idx
    ON hrms.workflow_task_comments(tenant_id, task_id, created_at)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.workflow_task_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    task_id UUID NOT NULL REFERENCES hrms.workflow_tasks(id) ON DELETE CASCADE,
    comment_id UUID REFERENCES hrms.workflow_task_comments(id) ON DELETE SET NULL,
    file_name VARCHAR(220) NOT NULL,
    content_type VARCHAR(160) NOT NULL,
    storage_path TEXT NOT NULL,
    checksum_sha256 VARCHAR(96),
    size_bytes BIGINT NOT NULL DEFAULT 0,
    visibility VARCHAR(40) NOT NULL DEFAULT 'tenant',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT workflow_task_attachments_file_check CHECK (BTRIM(file_name) <> ''),
    CONSTRAINT workflow_task_attachments_content_type_check CHECK (BTRIM(content_type) <> ''),
    CONSTRAINT workflow_task_attachments_path_check CHECK (BTRIM(storage_path) <> ''),
    CONSTRAINT workflow_task_attachments_size_check CHECK (size_bytes >= 0),
    CONSTRAINT workflow_task_attachments_visibility_check CHECK (visibility IN ('tenant','restricted','confidential')),
    CONSTRAINT workflow_task_attachments_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS workflow_task_attachments_task_idx
    ON hrms.workflow_task_attachments(tenant_id, task_id, created_at DESC)
    WHERE NOT inactive;

CREATE TABLE IF NOT EXISTS hrms.workflow_task_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    task_id UUID NOT NULL REFERENCES hrms.workflow_tasks(id) ON DELETE CASCADE,
    action VARCHAR(80) NOT NULL,
    from_status VARCHAR(40),
    to_status VARCHAR(40),
    actor_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    remarks TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT workflow_task_events_action_check CHECK (BTRIM(action) <> ''),
    CONSTRAINT workflow_task_events_metadata_object_check CHECK (jsonb_typeof(metadata) = 'object')
);

CREATE INDEX IF NOT EXISTS workflow_task_events_task_idx
    ON hrms.workflow_task_events(tenant_id, task_id, created_at)
    WHERE NOT inactive;

DROP TRIGGER IF EXISTS workflow_definitions_set_updated_at ON hrms.workflow_definitions;
CREATE TRIGGER workflow_definitions_set_updated_at BEFORE UPDATE ON hrms.workflow_definitions FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS workflow_definition_steps_set_updated_at ON hrms.workflow_definition_steps;
CREATE TRIGGER workflow_definition_steps_set_updated_at BEFORE UPDATE ON hrms.workflow_definition_steps FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS operation_templates_set_updated_at ON hrms.operation_templates;
CREATE TRIGGER operation_templates_set_updated_at BEFORE UPDATE ON hrms.operation_templates FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS workflow_tasks_set_updated_at ON hrms.workflow_tasks;
CREATE TRIGGER workflow_tasks_set_updated_at BEFORE UPDATE ON hrms.workflow_tasks FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS workflow_task_watchers_set_updated_at ON hrms.workflow_task_watchers;
CREATE TRIGGER workflow_task_watchers_set_updated_at BEFORE UPDATE ON hrms.workflow_task_watchers FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS workflow_task_comments_set_updated_at ON hrms.workflow_task_comments;
CREATE TRIGGER workflow_task_comments_set_updated_at BEFORE UPDATE ON hrms.workflow_task_comments FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS workflow_task_attachments_set_updated_at ON hrms.workflow_task_attachments;
CREATE TRIGGER workflow_task_attachments_set_updated_at BEFORE UPDATE ON hrms.workflow_task_attachments FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();
DROP TRIGGER IF EXISTS workflow_task_events_set_updated_at ON hrms.workflow_task_events;
CREATE TRIGGER workflow_task_events_set_updated_at BEFORE UPDATE ON hrms.workflow_task_events FOR EACH ROW EXECUTE FUNCTION hrms.set_updated_at();

ALTER TABLE hrms.workflow_definitions ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.workflow_definition_steps ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.operation_templates ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.workflow_tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.workflow_task_watchers ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.workflow_task_comments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.workflow_task_attachments ENABLE ROW LEVEL SECURITY;
ALTER TABLE hrms.workflow_task_events ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE
    table_name text;
BEGIN
    FOREACH table_name IN ARRAY ARRAY['workflow_definitions','workflow_definition_steps','operation_templates','workflow_tasks','workflow_task_watchers','workflow_task_comments','workflow_task_attachments','workflow_task_events']
    LOOP
        EXECUTE format('DROP POLICY IF EXISTS tenant_isolation_%1$s ON hrms.%1$s', table_name);
        EXECUTE format(
            'CREATE POLICY tenant_isolation_%1$s ON hrms.%1$s USING (current_setting(''app.is_super_admin'', true) = ''true'' OR tenant_id::text = current_setting(''app.tenant_id'', true)) WITH CHECK (current_setting(''app.is_super_admin'', true) = ''true'' OR tenant_id::text = current_setting(''app.tenant_id'', true))',
            table_name
        );
    END LOOP;
END $$;
