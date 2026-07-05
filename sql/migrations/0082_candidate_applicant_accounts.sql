CREATE TABLE IF NOT EXISTS hrms.candidate_applicant_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    candidate_id UUID NOT NULL,
    user_id UUID NOT NULL,
    email TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    consent_at TIMESTAMPTZ,
    consent_ip TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    inactive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID,
    CONSTRAINT candidate_applicant_accounts_candidate_fk FOREIGN KEY (candidate_id)
        REFERENCES hrms.candidates (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS candidate_applicant_accounts_user_uidx
    ON hrms.candidate_applicant_accounts (tenant_id, user_id)
    WHERE inactive = FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS candidate_applicant_accounts_candidate_uidx
    ON hrms.candidate_applicant_accounts (tenant_id, candidate_id)
    WHERE inactive = FALSE;

CREATE INDEX IF NOT EXISTS candidate_applicant_accounts_email_idx
    ON hrms.candidate_applicant_accounts (tenant_id, lower(email))
    WHERE inactive = FALSE;

ALTER TABLE hrms.candidate_applicant_accounts ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS candidate_applicant_accounts_tenant_isolation ON hrms.candidate_applicant_accounts;
CREATE POLICY candidate_applicant_accounts_tenant_isolation ON hrms.candidate_applicant_accounts
    USING (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    )
    WITH CHECK (
        current_setting('app.is_super_admin', true) = 'true'
        OR tenant_id::text = current_setting('app.tenant_id', true)
    );
