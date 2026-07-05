package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	LeaveApproverReportingManager = "reporting_manager"
	LeaveApproverManagerManager   = "manager_manager"
	LeaveApproverHRUser           = "hr_user"
	LeaveApproverSpecificUser     = "specific_user"
	LeaveApproverRole             = "role"
	LeaveApproverApplicant        = "applicant"

	LeaveApprovalDecisionAll = "all"
	LeaveApprovalDecisionAny = "any"
)

var (
	ErrInvalidLeaveApprovalWorkflowID   = errors.New("leave approval workflow id is required")
	ErrInvalidLeaveApprovalWorkflowName = errors.New("leave approval workflow name is required")
	ErrInvalidLeaveApprovalWorkflowCode = errors.New("leave approval workflow code is required")
	ErrInvalidLeaveApprovalStepID       = errors.New("leave approval workflow step id is required")
	ErrInvalidLeaveApproverType         = errors.New("leave approver type is invalid")
	ErrInvalidLeaveApprovalDecisionRule = errors.New("leave approval decision rule is invalid")
	ErrInvalidLeaveApprovalStepOrder    = errors.New("leave approval step order is invalid")
	ErrLeaveApprovalWorkflowNotFound    = errors.New("leave approval workflow not found")
	ErrLeaveApprovalStepNotFound        = errors.New("leave approval workflow step not found")
	ErrLeaveApprovalNotFound            = errors.New("leave approval not found")
	ErrLeaveApprovalNotPending          = errors.New("leave approval is not pending")
	ErrLeaveApprovalUnauthorized        = errors.New("approver is not allowed to act on this leave approval")
	ErrLeaveApproverNotResolved         = errors.New("leave approver could not be resolved")
	ErrLeaveNotPending                  = errors.New("leave is not pending")
	ErrLeaveCancelUnauthorized          = errors.New("user is not allowed to cancel this leave")
)

type LeaveApprovalWorkflow struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description *string    `json:"description,omitempty"`
	IsDefault   bool       `json:"is_default"`
	Inactive    bool       `json:"inactive"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
}

type LeaveApprovalWorkflowStep struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	WorkflowID        uuid.UUID  `json:"workflow_id"`
	StepOrder         int32      `json:"step_order"`
	Name              string     `json:"name"`
	ApproverType      string     `json:"approver_type"`
	ApproverUserID    *uuid.UUID `json:"approver_user_id,omitempty"`
	ApproverRole      *string    `json:"approver_role,omitempty"`
	DecisionRule      string     `json:"decision_rule"`
	RequiredApprovals int32      `json:"required_approvals"`
	AutoApprove       bool       `json:"auto_approve"`
	SLAHours          int32      `json:"sla_hours"`
	Inactive          bool       `json:"inactive"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty"`
}

func ValidateLeaveApprovalWorkflow(item *LeaveApprovalWorkflow) error {
	if item == nil {
		return ErrInvalidLeaveApprovalWorkflowID
	}
	if item.TenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return ErrInvalidLeaveApprovalWorkflowName
	}
	item.Code = strings.TrimSpace(item.Code)
	if item.Code == "" {
		return ErrInvalidLeaveApprovalWorkflowCode
	}
	return nil
}

func ValidateLeaveApprovalWorkflowStep(item *LeaveApprovalWorkflowStep) error {
	if item == nil {
		return ErrInvalidLeaveApprovalStepID
	}
	if item.TenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	if item.WorkflowID == uuid.Nil {
		return ErrInvalidLeaveApprovalWorkflowID
	}
	if item.StepOrder <= 0 {
		return ErrInvalidLeaveApprovalStepOrder
	}
	if item.Name == "" {
		item.Name = "Approval step"
	}
	if !validLeaveApproverType(item.ApproverType) {
		return ErrInvalidLeaveApproverType
	}
	if item.DecisionRule == "" {
		item.DecisionRule = LeaveApprovalDecisionAll
	}
	if !validLeaveApprovalDecisionRule(item.DecisionRule) {
		return ErrInvalidLeaveApprovalDecisionRule
	}
	if item.RequiredApprovals <= 0 {
		item.RequiredApprovals = 1
	}
	return nil
}

func validLeaveApproverType(value string) bool {
	switch value {
	case LeaveApproverReportingManager, LeaveApproverManagerManager, LeaveApproverHRUser, LeaveApproverSpecificUser, LeaveApproverRole, LeaveApproverApplicant:
		return true
	default:
		return false
	}
}

func validLeaveApprovalDecisionRule(value string) bool {
	switch value {
	case LeaveApprovalDecisionAll, LeaveApprovalDecisionAny:
		return true
	default:
		return false
	}
}
