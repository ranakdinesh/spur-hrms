package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type EmploymentLookupRepo interface {
	CreateEmploymentType(ctx context.Context, item *domain.EmploymentType, actorID *uuid.UUID) (*domain.EmploymentType, error)
	ListEmploymentTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.EmploymentType, error)
	GetEmploymentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmploymentType, error)
	UpdateEmploymentType(ctx context.Context, item *domain.EmploymentType, actorID *uuid.UUID) (*domain.EmploymentType, error)
	DeleteEmploymentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateMaritalStatus(ctx context.Context, item *domain.MaritalStatus, actorID *uuid.UUID) (*domain.MaritalStatus, error)
	ListMaritalStatuses(ctx context.Context, tenantID uuid.UUID) ([]*domain.MaritalStatus, error)
	GetMaritalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.MaritalStatus, error)
	UpdateMaritalStatus(ctx context.Context, item *domain.MaritalStatus, actorID *uuid.UUID) (*domain.MaritalStatus, error)
	DeleteMaritalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type EmploymentTypeCommand struct {
	ID       uuid.UUID  `json:"id,omitempty"`
	TenantID uuid.UUID  `json:"tenant_id"`
	Name     string     `json:"name"`
	ActorID  *uuid.UUID `json:"-"`
}

type MaritalStatusCommand struct {
	ID       uuid.UUID  `json:"id,omitempty"`
	TenantID uuid.UUID  `json:"tenant_id"`
	Name     string     `json:"name"`
	ActorID  *uuid.UUID `json:"-"`
}
