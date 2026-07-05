package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type DesignationRepo interface {
	CreateDesignation(ctx context.Context, designation *domain.Designation, actorID *uuid.UUID) (*domain.Designation, error)
	ListDesignations(ctx context.Context, tenantID uuid.UUID) ([]*domain.Designation, error)
	GetDesignation(ctx context.Context, tenantID uuid.UUID, designationID uuid.UUID) (*domain.Designation, error)
	UpdateDesignation(ctx context.Context, designation *domain.Designation, actorID *uuid.UUID) (*domain.Designation, error)
	DeleteDesignation(ctx context.Context, tenantID uuid.UUID, designationID uuid.UUID, actorID *uuid.UUID) error
}

type DesignationMasterRepo interface {
	CreateDesignationLevelCode(ctx context.Context, item *domain.DesignationLevelCode, actorID *uuid.UUID) (*domain.DesignationLevelCode, error)
	ListDesignationLevelCodes(ctx context.Context, tenantID uuid.UUID) ([]*domain.DesignationLevelCode, error)
	UpdateDesignationLevelCode(ctx context.Context, item *domain.DesignationLevelCode, actorID *uuid.UUID) (*domain.DesignationLevelCode, error)
	DeleteDesignationLevelCode(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateDesignationSeniorityRank(ctx context.Context, item *domain.DesignationSeniorityRank, actorID *uuid.UUID) (*domain.DesignationSeniorityRank, error)
	ListDesignationSeniorityRanks(ctx context.Context, tenantID uuid.UUID) ([]*domain.DesignationSeniorityRank, error)
	UpdateDesignationSeniorityRank(ctx context.Context, item *domain.DesignationSeniorityRank, actorID *uuid.UUID) (*domain.DesignationSeniorityRank, error)
	DeleteDesignationSeniorityRank(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type DesignationCommand struct {
	ID                 uuid.UUID  `json:"id,omitempty"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	Name               string     `json:"name"`
	LevelCode          string     `json:"level_code"`
	SeniorityRank      int32      `json:"seniority_rank"`
	Description        *string    `json:"description,omitempty"`
	AttendanceRequired *bool      `json:"attendance_required,omitempty"`
	ActorID            *uuid.UUID `json:"-"`
}

type DesignationLevelCodeCommand struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Code        string     `json:"code"`
	Label       string     `json:"label"`
	Description *string    `json:"description,omitempty"`
	SortOrder   int32      `json:"sort_order"`
	ActorID     *uuid.UUID `json:"-"`
}

type DesignationSeniorityRankCommand struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	RankValue   int32      `json:"rank_value"`
	Label       string     `json:"label"`
	Description *string    `json:"description,omitempty"`
	SortOrder   int32      `json:"sort_order"`
	ActorID     *uuid.UUID `json:"-"`
}
