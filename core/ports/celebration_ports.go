package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type CelebrationRepo interface {
	CreateCelebrationType(ctx context.Context, item *domain.CelebrationType, actorID *uuid.UUID) (*domain.CelebrationType, error)
	ListCelebrationTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.CelebrationType, error)
	GetCelebrationType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CelebrationType, error)
	UpdateCelebrationType(ctx context.Context, item *domain.CelebrationType, actorID *uuid.UUID) (*domain.CelebrationType, error)
	DeleteCelebrationType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateCelebration(ctx context.Context, item *domain.Celebration, actorID *uuid.UUID) (*domain.Celebration, error)
	ListCelebrations(ctx context.Context, tenantID uuid.UUID) ([]*domain.Celebration, error)
	ListCelebrationsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.Celebration, error)
	GetCelebration(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Celebration, error)
	UpdateCelebration(ctx context.Context, item *domain.Celebration, actorID *uuid.UUID) (*domain.Celebration, error)
	DeleteCelebration(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type CelebrationTypeCommand struct {
	ID                uuid.UUID  `json:"id,omitempty"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	Name              string     `json:"name"`
	IsYearly          bool       `json:"is_yearly"`
	IsUserCelebration bool       `json:"is_user_celebration"`
	ActorID           *uuid.UUID `json:"-"`
}

type CelebrationCommand struct {
	ID                uuid.UUID  `json:"id,omitempty"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	BranchID          *uuid.UUID `json:"branch_id,omitempty"`
	UserID            *uuid.UUID `json:"user_id,omitempty"`
	CelebrationTypeID uuid.UUID  `json:"celebration_type_id"`
	CelebrationDate   string     `json:"celebration_date"`
	CustomTitle       *string    `json:"custom_title,omitempty"`
	Description       *string    `json:"description,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}
