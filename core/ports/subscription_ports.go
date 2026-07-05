package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type TenantSubscriptionRepo interface {
	CreateTenantSubscription(ctx context.Context, item *domain.TenantSubscription, actorID *uuid.UUID) (*domain.TenantSubscription, error)
	ListTenantSubscriptions(ctx context.Context, tenantID uuid.UUID) ([]*domain.TenantSubscription, error)
	GetTenantSubscription(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.TenantSubscription, error)
	GetCurrentTenantSubscription(ctx context.Context, tenantID uuid.UUID) (*domain.TenantSubscription, error)
	UpdateTenantSubscription(ctx context.Context, item *domain.TenantSubscription, actorID *uuid.UUID) (*domain.TenantSubscription, error)
	DeleteTenantSubscription(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type TenantSubscriptionCommand struct {
	ID           uuid.UUID  `json:"id,omitempty"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	PlanID       *uuid.UUID `json:"plan_id,omitempty"`
	StartDate    string     `json:"start_date"`
	EndDate      string     `json:"end_date"`
	Status       string     `json:"status"`
	MaxEmployees int32      `json:"max_employees"`
	ActorID      *uuid.UUID `json:"-"`
}
