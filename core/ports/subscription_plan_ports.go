package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type SubscriptionPlanRepo interface {
	CreateSubscriptionPlan(ctx context.Context, item *domain.SubscriptionPlan, actorID *uuid.UUID) (*domain.SubscriptionPlan, error)
	ListSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error)
	ListActiveSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error)
	GetSubscriptionPlan(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error)
	UpdateSubscriptionPlan(ctx context.Context, item *domain.SubscriptionPlan, actorID *uuid.UUID) (*domain.SubscriptionPlan, error)
	DeleteSubscriptionPlan(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error
}

type SubscriptionPlanCommand struct {
	ID                uuid.UUID  `json:"id,omitempty"`
	Code              string     `json:"code"`
	Name              string     `json:"name"`
	Description       *string    `json:"description,omitempty"`
	PriceAmount       float64    `json:"price_amount"`
	PriceBasis        string     `json:"price_basis"`
	MinimumAmount     float64    `json:"minimum_amount"`
	IncludedEmployees int32      `json:"included_employees"`
	OverageAmount     float64    `json:"overage_amount"`
	CurrencyCode      string     `json:"currency_code"`
	BillingCycle      string     `json:"billing_cycle"`
	EmployeeLimit     int32      `json:"employee_limit"`
	TrialDays         int32      `json:"trial_days"`
	Visibility        string     `json:"visibility"`
	IsActive          bool       `json:"is_active"`
	ActorID           *uuid.UUID `json:"-"`
}
