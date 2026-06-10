package ports

import (
	"context"

	"github.com/google/uuid"
	"y/core/domain"
)

// ─── Repository ───────────────────────────────────────────────────────────────

type HrmsRepo interface {
	Create(ctx context.Context, e *domain.Hrms) (*domain.Hrms, error)
	GetByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*domain.Hrms, error)
	List(ctx context.Context, tenantID uuid.UUID) ([]*domain.Hrms, error)
	Delete(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) error
}

// ─── Service ──────────────────────────────────────────────────────────────────

type HrmsService interface {
	Create(ctx context.Context, cmd CreateHrmsCmd) (*domain.Hrms, error)
	Get(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*domain.Hrms, error)
	List(ctx context.Context, tenantID uuid.UUID) ([]*domain.Hrms, error)
	Delete(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) error
}

// ─── Commands ─────────────────────────────────────────────────────────────────

type CreateHrmsCmd struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	CreatedBy uuid.UUID `json:"created_by"`
	// TODO: add your command fields
}
