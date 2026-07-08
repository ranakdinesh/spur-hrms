package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type PolicyEngineRepo interface {
	CreatePolicySet(ctx context.Context, item *domain.PolicySet, actorID *uuid.UUID) (*domain.PolicySet, error)
	UpdatePolicySet(ctx context.Context, item *domain.PolicySet, actorID *uuid.UUID) (*domain.PolicySet, error)
	ListPolicySets(ctx context.Context, tenantID uuid.UUID, policyKind string) ([]*domain.PolicySet, error)
	GetPolicySet(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PolicySet, error)
	DeletePolicySet(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreatePolicyAssignment(ctx context.Context, item *domain.PolicyAssignment, actorID *uuid.UUID) (*domain.PolicyAssignment, error)
	UpdatePolicyAssignment(ctx context.Context, item *domain.PolicyAssignment, actorID *uuid.UUID) (*domain.PolicyAssignment, error)
	ListPolicyAssignments(ctx context.Context, tenantID uuid.UUID, policyKind string) ([]*domain.PolicyAssignment, error)
	DeletePolicyAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateLeavePolicyRule(ctx context.Context, item *domain.LeavePolicyRule, actorID *uuid.UUID) (*domain.LeavePolicyRule, error)
	UpdateLeavePolicyRule(ctx context.Context, item *domain.LeavePolicyRule, actorID *uuid.UUID) (*domain.LeavePolicyRule, error)
	ListLeavePolicyRules(ctx context.Context, tenantID uuid.UUID, policySetID uuid.UUID) ([]*domain.LeavePolicyRule, error)
	DeleteLeavePolicyRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	ResolvePolicySet(ctx context.Context, subject domain.PolicyResolutionSubject, policyKind string) (*domain.PolicySet, error)
	ListPolicyResolutionCandidates(ctx context.Context, subject domain.PolicyResolutionSubject, policyKind string) ([]domain.PolicyResolutionCandidate, error)
}

type PolicySetCommand struct {
	ID            uuid.UUID       `json:"id,omitempty"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	PolicyKind    string          `json:"policy_kind"`
	Code          string          `json:"code"`
	Name          string          `json:"name"`
	Description   *string         `json:"description,omitempty"`
	Config        json.RawMessage `json:"config,omitempty"`
	IsDefault     bool            `json:"is_default"`
	IsActive      bool            `json:"is_active"`
	EffectiveFrom string          `json:"effective_from,omitempty"`
	EffectiveTo   string          `json:"effective_to,omitempty"`
	ActorID       *uuid.UUID      `json:"-"`
}

type PolicyAssignmentCommand struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	PolicySetID   uuid.UUID  `json:"policy_set_id"`
	PolicyKind    string     `json:"policy_kind"`
	ScopeType     string     `json:"scope_type"`
	ScopeID       *uuid.UUID `json:"scope_id,omitempty"`
	RoleCode      *string    `json:"role_code,omitempty"`
	Priority      int32      `json:"priority"`
	EffectiveFrom string     `json:"effective_from,omitempty"`
	EffectiveTo   string     `json:"effective_to,omitempty"`
	IsActive      bool       `json:"is_active"`
	ActorID       *uuid.UUID `json:"-"`
}

type LeavePolicyRuleCommand struct {
	ID                           uuid.UUID       `json:"id,omitempty"`
	TenantID                     uuid.UUID       `json:"tenant_id"`
	PolicySetID                  uuid.UUID       `json:"policy_set_id"`
	LeaveTypeID                  uuid.UUID       `json:"leave_type_id"`
	GrantMode                    string          `json:"grant_mode"`
	AccrualFrequency             *string         `json:"accrual_frequency,omitempty"`
	EntitlementDays              float64         `json:"entitlement_days"`
	AccrualAmountPerPeriod       float64         `json:"accrual_amount_per_period"`
	ProrateJoiners               bool            `json:"prorate_joiners"`
	ProbationHandling            string          `json:"probation_handling"`
	RoundingRule                 string          `json:"rounding_rule"`
	MaxBalanceCap                *float64        `json:"max_balance_cap,omitempty"`
	CarryForwardCap              *float64        `json:"carry_forward_cap,omitempty"`
	EncashmentEligible           bool            `json:"encashment_eligible"`
	NegativeBalanceAllowed       bool            `json:"negative_balance_allowed"`
	InsufficientBalanceAction    string          `json:"insufficient_balance_action"`
	ExpiryDays                   *int32          `json:"expiry_days,omitempty"`
	AllowHalfDay                 bool            `json:"allow_half_day"`
	AttachmentRequiredAfterDays  *float64        `json:"attachment_required_after_days,omitempty"`
	ApprovalWorkflow             json.RawMessage `json:"approval_workflow,omitempty"`
	SandwichEnabled              bool            `json:"sandwich_enabled"`
	SandwichIncludeWeeklyOff     bool            `json:"sandwich_include_weekly_off"`
	SandwichIncludePublicHoliday bool            `json:"sandwich_include_public_holiday"`
	SandwichSameLeaveTypeOnly    bool            `json:"sandwich_same_leave_type_only"`
	SandwichAcrossLeaveTypes     bool            `json:"sandwich_across_leave_types"`
	NoticeRequiredAfterDays      *float64        `json:"notice_required_after_days,omitempty"`
	NoticeDays                   int32           `json:"notice_days"`
	PayrollImpact                string          `json:"payroll_impact"`
	RuleConfig                   json.RawMessage `json:"rule_config,omitempty"`
	ActorID                      *uuid.UUID      `json:"-"`
}
