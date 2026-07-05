package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type BranchRepo interface {
	CreateBranch(ctx context.Context, branch *domain.Branch, actorID *uuid.UUID) (*domain.Branch, error)
	ListBranches(ctx context.Context, tenantID uuid.UUID) ([]*domain.Branch, error)
	GetBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID) (*domain.Branch, error)
	UpdateBranch(ctx context.Context, branch *domain.Branch, actorID *uuid.UUID) (*domain.Branch, error)
	DeleteBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID, actorID *uuid.UUID) error
}

type BranchCommand struct {
	ID                  uuid.UUID  `json:"id,omitempty"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	Name                string     `json:"name"`
	Address             *string    `json:"address,omitempty"`
	City                *string    `json:"city,omitempty"`
	State               *string    `json:"state,omitempty"`
	Country             *string    `json:"country,omitempty"`
	Pincode             *string    `json:"pincode,omitempty"`
	Phone               *string    `json:"phone,omitempty"`
	BranchManagerUserID *uuid.UUID `json:"branch_manager_user_id,omitempty"`
	HRUserID            *uuid.UUID `json:"hr_user_id,omitempty"`
	AccountsUserID      *uuid.UUID `json:"accounts_user_id,omitempty"`
	ActorID             *uuid.UUID `json:"-"`
}
