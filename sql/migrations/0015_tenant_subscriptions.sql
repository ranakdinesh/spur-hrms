-- HRMS-015: tenant subscription lifecycle and employee limits.

ALTER TABLE hrms.tenant_subscriptions
    ALTER COLUMN status SET DEFAULT 'active';

UPDATE hrms.tenant_subscriptions
SET status = 'active'
WHERE status IS NULL OR btrim(status) = '';

ALTER TABLE hrms.tenant_subscriptions
    ALTER COLUMN status SET NOT NULL;

ALTER TABLE hrms.tenant_subscriptions
    DROP CONSTRAINT IF EXISTS tenant_subscriptions_status_check;

ALTER TABLE hrms.tenant_subscriptions
    ADD CONSTRAINT tenant_subscriptions_status_check CHECK (status IN ('trialing','active','past_due','cancelled','expired'));

CREATE INDEX IF NOT EXISTS tenant_subscriptions_status_idx ON hrms.tenant_subscriptions(tenant_id, status) WHERE NOT inactive;

CREATE UNIQUE INDEX IF NOT EXISTS tenant_subscriptions_current_idx
    ON hrms.tenant_subscriptions(tenant_id)
    WHERE NOT inactive AND status IN ('trialing','active','past_due');
