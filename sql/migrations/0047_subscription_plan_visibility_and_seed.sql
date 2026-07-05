-- HRMS-015B: public/internal subscription catalog and launch pricing.

ALTER TABLE hrms.subscription_plans
    ADD COLUMN IF NOT EXISTS price_basis TEXT NOT NULL DEFAULT 'per_employee',
    ADD COLUMN IF NOT EXISTS minimum_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS included_employees INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS overage_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS visibility TEXT NOT NULL DEFAULT 'public';

ALTER TABLE hrms.subscription_plans
    DROP CONSTRAINT IF EXISTS subscription_plans_price_basis_check;

ALTER TABLE hrms.subscription_plans
    ADD CONSTRAINT subscription_plans_price_basis_check
    CHECK (price_basis IN ('per_employee','package_plus_overage','flat','custom_quote'));

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'subscription_plans_visibility_check'
    ) THEN
        ALTER TABLE hrms.subscription_plans
            ADD CONSTRAINT subscription_plans_visibility_check
            CHECK (visibility IN ('public','internal'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'subscription_plans_minimum_nonnegative'
    ) THEN
        ALTER TABLE hrms.subscription_plans
            ADD CONSTRAINT subscription_plans_minimum_nonnegative
            CHECK (minimum_amount >= 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'subscription_plans_included_employees_nonnegative'
    ) THEN
        ALTER TABLE hrms.subscription_plans
            ADD CONSTRAINT subscription_plans_included_employees_nonnegative
            CHECK (included_employees >= 0);
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'subscription_plans_overage_nonnegative'
    ) THEN
        ALTER TABLE hrms.subscription_plans
            ADD CONSTRAINT subscription_plans_overage_nonnegative
            CHECK (overage_amount >= 0);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS subscription_plans_visibility_active_idx
    ON hrms.subscription_plans (visibility, is_active, name)
    WHERE NOT inactive;

INSERT INTO hrms.subscription_plans (
    code,
    name,
    description,
    price_amount,
    price_basis,
    minimum_amount,
    included_employees,
    overage_amount,
    currency_code,
    billing_cycle,
    employee_limit,
    trial_days,
    visibility,
    is_active
)
SELECT *
FROM (
    VALUES
        ('STARTER', 'Starter', 'Public plan for small teams. Minimum monthly commitment with per-employee billing; capped so small-plan customers do not consume enterprise capacity.', 39.00, 'per_employee', 999.00, 0, 39.00, 'INR', 'monthly', 50, 14, 'public', TRUE),
        ('GROWTH_100', 'Growth 100', 'Public package plan: fixed monthly price includes 100 employees, then per-employee overage within the plan cap.', 5000.00, 'package_plus_overage', 5000.00, 100, 45.00, 'INR', 'monthly', 150, 14, 'public', TRUE),
        ('BUSINESS_250', 'Business 250', 'Public package plan for established teams: includes 250 employees with a lower overage rate and stronger operating commitment.', 11000.00, 'package_plus_overage', 11000.00, 250, 40.00, 'INR', 'monthly', 350, 14, 'public', TRUE),
        ('ENTERPRISE_500', 'Enterprise 500', 'Public enterprise package with larger included headcount, lower overage, and room for annual commercial commitments.', 20000.00, 'package_plus_overage', 20000.00, 500, 35.00, 'INR', 'monthly', 750, 14, 'public', TRUE),
        ('NEGOTIATED_STARTUP', 'Negotiated Startup', 'Internal negotiated startup plan. Use only when sales approves a discounted floor and smaller cap.', 29.00, 'per_employee', 699.00, 0, 29.00, 'INR', 'monthly', 35, 14, 'internal', TRUE),
        ('NEGOTIATED_GROWTH_100', 'Negotiated Growth 100', 'Internal negotiated package for growing teams with commercial approval.', 4500.00, 'package_plus_overage', 4500.00, 100, 39.00, 'INR', 'monthly', 150, 14, 'internal', TRUE),
        ('ENTERPRISE_CUSTOM', 'Enterprise Custom', 'Internal quote-based enterprise plan for custom contracts, annual prepay, bundled services, or non-standard limits.', 0.00, 'custom_quote', 0.00, 0, 0.00, 'INR', 'custom', 1000, 0, 'internal', TRUE)
) AS seed(code, name, description, price_amount, price_basis, minimum_amount, included_employees, overage_amount, currency_code, billing_cycle, employee_limit, trial_days, visibility, is_active)
WHERE NOT EXISTS (
    SELECT 1
    FROM hrms.subscription_plans existing
    WHERE lower(existing.code) = lower(seed.code)
      AND NOT existing.inactive
);

UPDATE hrms.subscription_plans
SET name = seed.name,
    description = seed.description,
    price_amount = seed.price_amount,
    price_basis = seed.price_basis,
    minimum_amount = seed.minimum_amount,
    included_employees = seed.included_employees,
    overage_amount = seed.overage_amount,
    currency_code = seed.currency_code,
    billing_cycle = seed.billing_cycle,
    employee_limit = seed.employee_limit,
    trial_days = seed.trial_days,
    visibility = seed.visibility,
    is_active = TRUE,
    updated_at = NOW()
FROM (
    VALUES
        ('STARTER', 'Starter', 'Public plan for small teams. Minimum monthly commitment with per-employee billing; capped so small-plan customers do not consume enterprise capacity.', 39.00, 'per_employee', 999.00, 0, 39.00, 'INR', 'monthly', 50, 14, 'public'),
        ('PROFESSIONAL', 'Professional', 'Legacy public plan replaced by Growth 100. Kept inactive for historical subscriptions.', 0.00, 'custom_quote', 0.00, 0, 0.00, 'INR', 'custom', 0, 0, 'public'),
        ('CUSTOM', 'Custom', 'Legacy public plan replaced by Business 250 and Enterprise 500. Kept inactive for historical subscriptions.', 0.00, 'custom_quote', 0.00, 0, 0.00, 'INR', 'custom', 0, 0, 'public'),
        ('NEGOTIATED_GROWTH', 'Negotiated Growth', 'Legacy negotiated plan replaced by Negotiated Growth 100. Kept inactive for historical subscriptions.', 0.00, 'custom_quote', 0.00, 0, 0.00, 'INR', 'custom', 0, 0, 'internal')
) AS seed(code, name, description, price_amount, price_basis, minimum_amount, included_employees, overage_amount, currency_code, billing_cycle, employee_limit, trial_days, visibility)
WHERE lower(hrms.subscription_plans.code) = lower(seed.code)
  AND NOT hrms.subscription_plans.inactive;

UPDATE hrms.subscription_plans
SET is_active = FALSE,
    updated_at = NOW()
WHERE code IN ('PROFESSIONAL', 'CUSTOM', 'NEGOTIATED_GROWTH')
  AND NOT inactive;
