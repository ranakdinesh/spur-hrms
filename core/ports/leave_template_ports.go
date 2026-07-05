package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type LeaveTemplateRepo interface {
	CreateLeavePolicyTemplate(ctx context.Context, item *domain.LeavePolicyTemplate, actorID *uuid.UUID) (*domain.LeavePolicyTemplate, error)
	ListLeavePolicyTemplates(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeavePolicyTemplate, error)
	GetLeavePolicyTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeavePolicyTemplate, error)
	UpdateLeavePolicyTemplate(ctx context.Context, item *domain.LeavePolicyTemplate, actorID *uuid.UUID) (*domain.LeavePolicyTemplate, error)
	DeleteLeavePolicyTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateLeavePolicyTemplateRule(ctx context.Context, item *domain.LeavePolicyTemplateRule, actorID *uuid.UUID) (*domain.LeavePolicyTemplateRule, error)
	ListLeavePolicyTemplateRules(ctx context.Context, tenantID uuid.UUID, templateID uuid.UUID) ([]*domain.LeavePolicyTemplateRule, error)
	ListLeavePolicyTemplateRulesByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeavePolicyTemplateRule, error)
	GetLeavePolicyTemplateRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeavePolicyTemplateRule, error)
	UpdateLeavePolicyTemplateRule(ctx context.Context, item *domain.LeavePolicyTemplateRule, actorID *uuid.UUID) (*domain.LeavePolicyTemplateRule, error)
	DeleteLeavePolicyTemplateRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	UpsertEmployeeLeavePolicyAssignment(ctx context.Context, item *domain.EmployeeLeavePolicyAssignment, actorID *uuid.UUID) (*domain.EmployeeLeavePolicyAssignment, error)
	ListEmployeeLeavePolicyAssignments(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.EmployeeLeavePolicyAssignment, error)
	ListLeavePolicyAssignmentsByTemplate(ctx context.Context, tenantID uuid.UUID, templateID uuid.UUID) ([]*domain.EmployeeLeavePolicyAssignment, error)
	DeleteEmployeeLeavePolicyAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type LeavePolicyTemplateCommand struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	Description   *string    `json:"description,omitempty"`
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}

type LeavePolicyTemplateRuleCommand struct {
	ID                        uuid.UUID      `json:"id,omitempty"`
	TenantID                  uuid.UUID      `json:"tenant_id"`
	TemplateID                uuid.UUID      `json:"template_id"`
	LeaveTypeID               uuid.UUID      `json:"leave_type_id"`
	FYID                      *uuid.UUID     `json:"fy_id,omitempty"`
	EmploymentTypeID          *uuid.UUID     `json:"employment_type_id,omitempty"`
	DepartmentID              *uuid.UUID     `json:"department_id,omitempty"`
	DesignationID             *uuid.UUID     `json:"designation_id,omitempty"`
	ProbationStatus           *string        `json:"probation_status,omitempty"`
	AccrualMethod             string         `json:"accrual_method"`
	AccrualFrequency          string         `json:"accrual_frequency"`
	CreditDays                float64        `json:"credit_days"`
	CreditHours               float64        `json:"credit_hours"`
	AnnualEntitlement         float64        `json:"annual_entitlement"`
	MinWorkedDays             int32          `json:"min_worked_days"`
	MaxBalance                *float64       `json:"max_balance,omitempty"`
	CarryForwardEnabled       bool           `json:"carry_forward_enabled"`
	MaxCarryForward           float64        `json:"max_carry_forward"`
	CarryForwardExpiryMonths  int32          `json:"carry_forward_expiry_months"`
	EncashmentEnabled         bool           `json:"encashment_enabled"`
	EncashmentLimit           float64        `json:"encashment_limit"`
	EncashmentPayablePercent  float64        `json:"encashment_payable_percent"`
	NegativeBalanceAllowed    bool           `json:"negative_balance_allowed"`
	MaxNegativeBalance        float64        `json:"max_negative_balance"`
	SandwichApplicable        bool           `json:"sandwich_applicable"`
	IncludeHolidays           bool           `json:"include_holidays"`
	IncludeWeekoffs           bool           `json:"include_weekoffs"`
	RequiresDocumentAfterDays *float64       `json:"requires_document_after_days,omitempty"`
	MinRequestDays            float64        `json:"min_request_days"`
	MaxRequestDays            *float64       `json:"max_request_days,omitempty"`
	MaxRequestsPerYear        int32          `json:"max_requests_per_year"`
	AccrualDay                int32          `json:"accrual_day"`
	LapseUnutilized           bool           `json:"lapse_unutilized"`
	AllowHalfDay              bool           `json:"allow_half_day"`
	RequiresApproval          bool           `json:"requires_approval"`
	CalculationConfig         map[string]any `json:"calculation_config"`
	Priority                  int32          `json:"priority"`
	ActorID                   *uuid.UUID     `json:"-"`
}

type EmployeeLeavePolicyAssignmentCommand struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	UserID        uuid.UUID  `json:"user_id"`
	TemplateID    uuid.UUID  `json:"template_id"`
	FYID          *uuid.UUID `json:"fy_id,omitempty"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}
