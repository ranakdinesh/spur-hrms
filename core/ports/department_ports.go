package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type DepartmentRepo interface {
	CreateDepartment(ctx context.Context, department *domain.Department, actorID *uuid.UUID) (*domain.Department, error)
	ListDepartments(ctx context.Context, tenantID uuid.UUID) ([]*domain.Department, error)
	GetDepartment(ctx context.Context, tenantID uuid.UUID, departmentID uuid.UUID) (*domain.Department, error)
	UpdateDepartment(ctx context.Context, department *domain.Department, actorID *uuid.UUID) (*domain.Department, error)
	DeleteDepartment(ctx context.Context, tenantID uuid.UUID, departmentID uuid.UUID, actorID *uuid.UUID) error
}

type DepartmentCommand struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Name        string     `json:"name"`
	ShortCode   string     `json:"short_code"`
	Description *string    `json:"description,omitempty"`
	ActorID     *uuid.UUID `json:"-"`
}
