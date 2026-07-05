-- name: CreateSubscriptionPlan :one
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
    is_active,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $15
)
RETURNING *;

-- name: ListSubscriptionPlans :many
SELECT * FROM hrms.subscription_plans
WHERE NOT inactive
ORDER BY is_active DESC, name ASC;

-- name: ListActiveSubscriptionPlans :many
SELECT * FROM hrms.subscription_plans
WHERE is_active AND NOT inactive
ORDER BY name ASC;

-- name: GetSubscriptionPlan :one
SELECT * FROM hrms.subscription_plans
WHERE id = $1 AND NOT inactive;

-- name: GetSubscriptionPlanByCode :one
SELECT * FROM hrms.subscription_plans
WHERE lower(code) = lower($1) AND NOT inactive;

-- name: UpdateSubscriptionPlan :one
UPDATE hrms.subscription_plans
SET code = $2,
    name = $3,
    description = $4,
    price_amount = $5,
    price_basis = $6,
    minimum_amount = $7,
    included_employees = $8,
    overage_amount = $9,
    currency_code = $10,
    billing_cycle = $11,
    employee_limit = $12,
    trial_days = $13,
    visibility = $14,
    is_active = $15,
    updated_by = $16,
    updated_at = NOW()
WHERE id = $1 AND NOT inactive
RETURNING *;

-- name: SoftDeleteSubscriptionPlan :exec
UPDATE hrms.subscription_plans
SET inactive = TRUE,
    is_active = FALSE,
    updated_by = $2,
    updated_at = NOW()
WHERE id = $1 AND NOT inactive;
