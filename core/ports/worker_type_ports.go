package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type WorkerTypeRepo interface {
	CreateWorkerType(ctx context.Context, item *domain.WorkerType, actorID *uuid.UUID) (*domain.WorkerType, error)
	ListWorkerTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.WorkerType, error)
	GetWorkerType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerType, error)
	GetWorkerTypeByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.WorkerType, error)
	UpdateWorkerType(ctx context.Context, item *domain.WorkerType, actorID *uuid.UUID) (*domain.WorkerType, error)
	DeleteWorkerType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CountWorkerClassificationRules(ctx context.Context, tenantID uuid.UUID, workerTypeID uuid.UUID) (int64, error)
	CreateWorkerClassificationRule(ctx context.Context, item *domain.WorkerClassificationRule, actorID *uuid.UUID) (*domain.WorkerClassificationRule, error)
	ListWorkerClassificationRules(ctx context.Context, tenantID uuid.UUID, workerTypeID *uuid.UUID) ([]*domain.WorkerClassificationRule, error)
	GetWorkerClassificationRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerClassificationRule, error)
	UpdateWorkerClassificationRule(ctx context.Context, item *domain.WorkerClassificationRule, actorID *uuid.UUID) (*domain.WorkerClassificationRule, error)
	DeleteWorkerClassificationRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type WorkerTypeCommand struct {
	ID                  uuid.UUID       `json:"id,omitempty"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	Code                string          `json:"code"`
	Name                string          `json:"name"`
	ClassificationGroup string          `json:"classification_group"`
	Description         *string         `json:"description,omitempty"`
	AttendanceMode      string          `json:"attendance_mode"`
	PayMode             string          `json:"pay_mode"`
	TDSSection          string          `json:"tds_section"`
	PFApplicable        bool            `json:"pf_applicable"`
	ESICApplicable      bool            `json:"esic_applicable"`
	PTApplicable        bool            `json:"pt_applicable"`
	LWFApplicable       bool            `json:"lwf_applicable"`
	CLRAApplicable      bool            `json:"clra_applicable"`
	LeaveApplicable     bool            `json:"leave_applicable"`
	OvertimeApplicable  bool            `json:"overtime_applicable"`
	RequiresAgreement   bool            `json:"requires_agreement"`
	RequiresInvoice     bool            `json:"requires_invoice"`
	RequiresAttendance  bool            `json:"requires_attendance"`
	StatutoryDefaults   json.RawMessage `json:"statutory_defaults,omitempty"`
	ComplianceNotes     *string         `json:"compliance_notes,omitempty"`
	SortOrder           int32           `json:"sort_order"`
	ActorID             *uuid.UUID      `json:"-"`
}

type WorkerClassificationRuleCommand struct {
	ID           uuid.UUID       `json:"id,omitempty"`
	TenantID     uuid.UUID       `json:"tenant_id"`
	WorkerTypeID uuid.UUID       `json:"worker_type_id"`
	RuleName     string          `json:"rule_name"`
	RuleType     string          `json:"rule_type"`
	Priority     int32           `json:"priority"`
	Conditions   json.RawMessage `json:"conditions,omitempty"`
	Outcome      json.RawMessage `json:"outcome,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	ActorID      *uuid.UUID      `json:"-"`
}
