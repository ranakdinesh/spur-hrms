CREATE TABLE IF NOT EXISTS hrms.leave_request_messages (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES auth.tenants(id) ON DELETE CASCADE,
    leave_id          UUID NOT NULL REFERENCES hrms.leaves(id) ON DELETE CASCADE,
    sender_user_id    UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    recipient_user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    message_type      VARCHAR(40) NOT NULL,
    body              TEXT NOT NULL,
    inactive          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by        UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by        UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT leave_request_messages_type_check CHECK (message_type IN ('clarification_request','employee_reply','comment')),
    CONSTRAINT leave_request_messages_body_check CHECK (length(trim(body)) > 0)
);

CREATE INDEX IF NOT EXISTS leave_request_messages_leave_idx
    ON hrms.leave_request_messages(tenant_id, leave_id, created_at);

DROP TRIGGER IF EXISTS leave_request_messages_updated_at ON hrms.leave_request_messages;
CREATE TRIGGER leave_request_messages_updated_at
    BEFORE UPDATE ON hrms.leave_request_messages
    FOR EACH ROW EXECUTE FUNCTION hrms.update_updated_at();

ALTER TABLE hrms.leave_request_messages ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS leave_request_messages_tenant_isolation ON hrms.leave_request_messages;
CREATE POLICY leave_request_messages_tenant_isolation ON hrms.leave_request_messages
    USING (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
        OR current_setting('app.is_super_admin', true) = 'true'
    )
    WITH CHECK (
        tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
        OR current_setting('app.is_super_admin', true) = 'true'
    );
