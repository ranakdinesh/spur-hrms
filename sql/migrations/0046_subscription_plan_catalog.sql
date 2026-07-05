-- HRMS-015A: super-admin managed subscription plan catalog.

CREATE TABLE IF NOT EXISTS hrms.subscription_plans (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            TEXT NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT,
    price_amount    NUMERIC(12,2) NOT NULL DEFAULT 0,
    price_basis     TEXT NOT NULL DEFAULT 'per_employee',
    minimum_amount  NUMERIC(12,2) NOT NULL DEFAULT 0,
    included_employees INT NOT NULL DEFAULT 0,
    overage_amount  NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency_code   CHAR(3) NOT NULL DEFAULT 'INR',
    billing_cycle   TEXT NOT NULL DEFAULT 'monthly',
    employee_limit  INT NOT NULL DEFAULT 0,
    trial_days      INT NOT NULL DEFAULT 14,
    visibility      TEXT NOT NULL DEFAULT 'public',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    inactive        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    CONSTRAINT subscription_plans_code_not_blank CHECK (btrim(code) <> ''),
    CONSTRAINT subscription_plans_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT subscription_plans_price_nonnegative CHECK (price_amount >= 0),
    CONSTRAINT subscription_plans_price_basis_check CHECK (price_basis IN ('per_employee','package_plus_overage','flat','custom_quote')),
    CONSTRAINT subscription_plans_minimum_nonnegative CHECK (minimum_amount >= 0),
    CONSTRAINT subscription_plans_included_employees_nonnegative CHECK (included_employees >= 0),
    CONSTRAINT subscription_plans_overage_nonnegative CHECK (overage_amount >= 0),
    CONSTRAINT subscription_plans_employee_limit_nonnegative CHECK (employee_limit >= 0),
    CONSTRAINT subscription_plans_trial_days_nonnegative CHECK (trial_days >= 0),
    CONSTRAINT subscription_plans_visibility_check CHECK (visibility IN ('public','internal')),
    CONSTRAINT subscription_plans_billing_cycle_check CHECK (billing_cycle IN ('monthly','quarterly','yearly','one_time','custom'))
);

CREATE UNIQUE INDEX IF NOT EXISTS subscription_plans_code_active_idx
    ON hrms.subscription_plans (lower(code))
    WHERE NOT inactive;

CREATE INDEX IF NOT EXISTS subscription_plans_active_idx
    ON hrms.subscription_plans (is_active, name)
    WHERE NOT inactive;
