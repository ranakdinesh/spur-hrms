package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	PolicyKindAttendance = "attendance"
	PolicyKindLeave      = "leave"

	PolicyScopeTenant      = "tenant"
	PolicyScopeBranch      = "branch"
	PolicyScopeDepartment  = "department"
	PolicyScopeDesignation = "designation"
	PolicyScopeWorkforce   = "workforce_type"
	PolicyScopeRoleGroup   = "role_group"
	PolicyScopeEmployee    = "employee"
	PolicyScopeDefault     = "default"
)

var (
	ErrInvalidPolicySetID        = errors.New("policy_set_id is required")
	ErrInvalidPolicyKind         = errors.New("policy kind is invalid")
	ErrInvalidPolicyCode         = errors.New("policy code is required")
	ErrInvalidPolicyName         = errors.New("policy name is required")
	ErrInvalidPolicyAssignmentID = errors.New("policy_assignment_id is required")
	ErrInvalidPolicyScope        = errors.New("policy assignment scope is invalid")
	ErrInvalidLeavePolicyRuleID  = errors.New("leave_policy_rule_id is required")
	ErrPolicySetNotFound         = errors.New("policy set not found")
	ErrPolicyAssignmentNotFound  = errors.New("policy assignment not found")
	ErrLeavePolicyRuleNotFound   = errors.New("leave policy rule not found")
)

type PolicySet struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	PolicyKind    string          `json:"policy_kind"`
	Code          string          `json:"code"`
	Name          string          `json:"name"`
	Description   *string         `json:"description,omitempty"`
	Config        json.RawMessage `json:"config,omitempty"`
	IsDefault     bool            `json:"is_default"`
	IsActive      bool            `json:"is_active"`
	EffectiveFrom *time.Time      `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time      `json:"effective_to,omitempty"`
	Inactive      bool            `json:"inactive"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     *uuid.UUID      `json:"updated_by,omitempty"`
}

type PolicyAssignment struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	PolicySetID   uuid.UUID  `json:"policy_set_id"`
	PolicyKind    string     `json:"policy_kind"`
	ScopeType     string     `json:"scope_type"`
	ScopeID       *uuid.UUID `json:"scope_id,omitempty"`
	RoleCode      *string    `json:"role_code,omitempty"`
	Priority      int32      `json:"priority"`
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	IsActive      bool       `json:"is_active"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type LeavePolicyRule struct {
	ID                           uuid.UUID       `json:"id"`
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
	Inactive                     bool            `json:"inactive"`
	CreatedAt                    time.Time       `json:"created_at"`
	CreatedBy                    *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                    time.Time       `json:"updated_at"`
	UpdatedBy                    *uuid.UUID      `json:"updated_by,omitempty"`
}

type PolicyResolutionSubject struct {
	TenantID        uuid.UUID
	EmployeeUserID  uuid.UUID
	DesignationID   *uuid.UUID
	WorkforceTypeID *uuid.UUID
	DepartmentID    *uuid.UUID
	BranchID        *uuid.UUID
	RoleCodes       []string
	Date            time.Time
}

type PolicyResolutionResult struct {
	Policy     *PolicySet                  `json:"policy"`
	LeaveRules []*LeavePolicyRule          `json:"leave_rules,omitempty"`
	Candidates []PolicyResolutionCandidate `json:"candidates,omitempty"`
}

type PolicyResolutionCandidate struct {
	PolicySetID uuid.UUID  `json:"policy_set_id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	PolicyKind  string     `json:"policy_kind"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	ScopeType   string     `json:"scope_type"`
	ScopeID     *uuid.UUID `json:"scope_id,omitempty"`
	RoleCode    *string    `json:"role_code,omitempty"`
	Precedence  int32      `json:"precedence"`
}

func ValidatePolicyKind(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	switch clean {
	case PolicyKindAttendance, PolicyKindLeave:
		return clean, nil
	default:
		return "", ErrInvalidPolicyKind
	}
}

func ValidatePolicyScope(scopeType string, scopeID *uuid.UUID, roleCode *string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(scopeType))
	switch clean {
	case PolicyScopeTenant:
		if scopeID != nil || roleCode != nil {
			return "", ErrInvalidPolicyScope
		}
	case PolicyScopeRoleGroup:
		if scopeID != nil || roleCode == nil || strings.TrimSpace(*roleCode) == "" {
			return "", ErrInvalidPolicyScope
		}
	case PolicyScopeBranch, PolicyScopeDepartment, PolicyScopeDesignation, PolicyScopeWorkforce, PolicyScopeEmployee:
		if scopeID == nil || *scopeID == uuid.Nil || roleCode != nil {
			return "", ErrInvalidPolicyScope
		}
	default:
		return "", ErrInvalidPolicyScope
	}
	return clean, nil
}

func NormalizeRoleCodes(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		clean := strings.ToLower(strings.TrimSpace(value))
		if clean == "" {
			continue
		}
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		result = append(result, clean)
	}
	return result
}
