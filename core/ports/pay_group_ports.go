package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type PayGroupRepo interface {
	CreatePayGroup(ctx context.Context, item *domain.PayGroup, actorID *uuid.UUID) (*domain.PayGroup, error)
	UpdatePayGroup(ctx context.Context, item *domain.PayGroup, actorID *uuid.UUID) (*domain.PayGroup, error)
	GetPayGroup(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayGroup, error)
	ListPayGroups(ctx context.Context, tenantID uuid.UUID) ([]*domain.PayGroup, error)
	DeletePayGroup(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	UpsertPayGroupMember(ctx context.Context, item *domain.PayGroupMember, actorID *uuid.UUID) (*domain.PayGroupMember, error)
	ListPayGroupMembers(ctx context.Context, tenantID uuid.UUID, payGroupID uuid.UUID) ([]*domain.PayGroupMember, error)
	DeletePayGroupMember(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	ListPayGroupEmployees(ctx context.Context, tenantID uuid.UUID, payGroupID uuid.UUID) ([]*domain.PayGroupEmployee, error)
	CreatePayRun(ctx context.Context, item *domain.PayRun, actorID *uuid.UUID) (*domain.PayRun, error)
	GetPayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayRun, error)
	ListPayRuns(ctx context.Context, tenantID uuid.UUID, payGroupID *uuid.UUID, month *int32, year *int32) ([]*domain.PayRun, error)
	UpdatePayRunStatus(ctx context.Context, item *domain.PayRun, actorID *uuid.UUID) (*domain.PayRun, error)
	DeletePayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	UpsertPayRunEmployee(ctx context.Context, item *domain.PayRunEmployee, actorID *uuid.UUID) (*domain.PayRunEmployee, error)
	ListPayRunEmployees(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) ([]*domain.PayRunEmployee, error)
	DeletePayRunLedger(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID, actorID *uuid.UUID) error
	CreatePayRunInput(ctx context.Context, item *domain.PayRunInput, actorID *uuid.UUID) (*domain.PayRunInput, error)
	CreatePayRunComponent(ctx context.Context, item *domain.PayRunComponent, actorID *uuid.UUID) (*domain.PayRunComponent, error)
	ListPayRunInputs(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) ([]*domain.PayRunInput, error)
	ListPayRunComponents(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) ([]*domain.PayRunComponent, error)
	GetPayRunLedgerSummary(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) (*domain.PayRunLedgerSummary, error)
	CreatePayRunEvent(ctx context.Context, item *domain.PayRunEvent, actorID *uuid.UUID) (*domain.PayRunEvent, error)
	ListPayRunEvents(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) ([]*domain.PayRunEvent, error)
}

type PayGroupCommand struct {
	ID               uuid.UUID       `json:"id,omitempty"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	Code             string          `json:"code"`
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	GroupingType     string          `json:"grouping_type"`
	BranchID         *uuid.UUID      `json:"branch_id,omitempty"`
	DepartmentID     *uuid.UUID      `json:"department_id,omitempty"`
	EmploymentTypeID *uuid.UUID      `json:"employment_type_id,omitempty"`
	ReportingTag     *string         `json:"reporting_tag,omitempty"`
	Rules            json.RawMessage `json:"rules,omitempty"`
	IsActive         bool            `json:"is_active"`
	ActorID          *uuid.UUID      `json:"-"`
}

type PayGroupMemberCommand struct {
	TenantID       uuid.UUID  `json:"tenant_id"`
	PayGroupID     uuid.UUID  `json:"pay_group_id"`
	UserID         uuid.UUID  `json:"user_id"`
	MembershipType string     `json:"membership_type"`
	EffectiveFrom  string     `json:"effective_from,omitempty"`
	EffectiveTo    string     `json:"effective_to,omitempty"`
	ActorID        *uuid.UUID `json:"-"`
}

type PayRunCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	PayGroupID uuid.UUID  `json:"pay_group_id"`
	FYID       uuid.UUID  `json:"fy_id"`
	Month      int32      `json:"month"`
	Year       int32      `json:"year"`
	Notes      *string    `json:"notes,omitempty"`
	ActorID    *uuid.UUID `json:"-"`
}

type PayRunListQuery struct {
	TenantID   uuid.UUID
	PayGroupID *uuid.UUID
	Month      *int32
	Year       *int32
}

type PayRunActionCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	PayRunID   uuid.UUID  `json:"pay_run_id"`
	Action     string     `json:"action"`
	Remarks    *string    `json:"remarks,omitempty"`
	Regenerate bool       `json:"regenerate"`
	ActorID    *uuid.UUID `json:"-"`
}
