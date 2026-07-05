ALTER TABLE hrms.policy_types
    ALTER COLUMN tenant_id DROP NOT NULL,
    ADD COLUMN IF NOT EXISTS is_system BOOLEAN NOT NULL DEFAULT FALSE;

DROP INDEX IF EXISTS hrms.policy_types_tenant_name_idx;
CREATE UNIQUE INDEX IF NOT EXISTS policy_types_tenant_name_idx
    ON hrms.policy_types(tenant_id, lower(name))
    WHERE NOT inactive AND NOT is_system;
CREATE UNIQUE INDEX IF NOT EXISTS policy_types_system_name_idx
    ON hrms.policy_types(lower(name))
    WHERE NOT inactive AND is_system;
CREATE INDEX IF NOT EXISTS policy_types_system_idx
    ON hrms.policy_types(is_system)
    WHERE NOT inactive;

INSERT INTO hrms.policy_types (tenant_id, name, is_system)
VALUES
    (NULL, 'HR Policy', TRUE),
    (NULL, 'Leave Policy', TRUE),
    (NULL, 'Attendance Policy', TRUE),
    (NULL, 'Payroll Policy', TRUE),
    (NULL, 'Code of Conduct', TRUE),
    (NULL, 'Security Policy', TRUE)
ON CONFLICT DO NOTHING;

DROP POLICY IF EXISTS policy_types_tenant_isolation ON hrms.policy_types;
CREATE POLICY policy_types_tenant_isolation ON hrms.policy_types
USING (
    is_system
    OR tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
    OR current_setting('app.is_super_admin', true) = 'true'
)
WITH CHECK (
    (is_system AND current_setting('app.is_super_admin', true) = 'true')
    OR tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid
    OR current_setting('app.is_super_admin', true) = 'true'
);
