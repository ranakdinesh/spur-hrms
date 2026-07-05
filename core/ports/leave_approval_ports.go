package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type LeaveApprovalWorkflowRepo interface {
	CreateLeaveApprovalWorkflow(ctx context.Context, item *domain.LeaveApprovalWorkflow, actorID *uuid.UUID) (*domain.LeaveApprovalWorkflow, error)
	ListLeaveApprovalWorkflows(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeaveApprovalWorkflow, error)
	GetLeaveApprovalWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeaveApprovalWorkflow, error)
	GetDefaultLeaveApprovalWorkflow(ctx context.Context, tenantID uuid.UUID) (*domain.LeaveApprovalWorkflow, error)
	UpdateLeaveApprovalWorkflow(ctx context.Context, item *domain.LeaveApprovalWorkflow, actorID *uuid.UUID) (*domain.LeaveApprovalWorkflow, error)
	DeleteLeaveApprovalWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateLeaveApprovalWorkflowStep(ctx context.Context, item *domain.LeaveApprovalWorkflowStep, actorID *uuid.UUID) (*domain.LeaveApprovalWorkflowStep, error)
	ListLeaveApprovalWorkflowSteps(ctx context.Context, tenantID uuid.UUID, workflowID uuid.UUID) ([]*domain.LeaveApprovalWorkflowStep, error)
	GetLeaveApprovalWorkflowStep(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeaveApprovalWorkflowStep, error)
	UpdateLeaveApprovalWorkflowStep(ctx context.Context, item *domain.LeaveApprovalWorkflowStep, actorID *uuid.UUID) (*domain.LeaveApprovalWorkflowStep, error)
	DeleteLeaveApprovalWorkflowStep(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateWorkflowLeaveApproval(ctx context.Context, item *domain.LeaveApproval, actorID *uuid.UUID) (*domain.LeaveApproval, error)
	GetLeaveApproval(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeaveApproval, error)
	ListLeaveApprovalsByLeave(ctx context.Context, tenantID uuid.UUID, leaveID uuid.UUID) ([]*domain.LeaveApproval, error)
	UpdateLeaveApprovalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.LeaveApproval, error)
	UpdateLeaveStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.Leave, error)
	ListPendingApprovalsByApprover(ctx context.Context, tenantID uuid.UUID, approverID uuid.UUID) ([]*domain.LeaveApproval, error)
}

type LeaveApprovalWorkflowCommand struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description *string    `json:"description,omitempty"`
	IsDefault   bool       `json:"is_default"`
	ActorID     *uuid.UUID `json:"-"`
}

type LeaveApprovalWorkflowStepCommand struct {
	ID                uuid.UUID  `json:"id,omitempty"`
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
	ActorID           *uuid.UUID `json:"-"`
}

type ApproveLeaveCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	ApprovalID uuid.UUID  `json:"approval_id"`
	ApproverID uuid.UUID  `json:"approver_id"`
	Remarks    *string    `json:"remarks,omitempty"`
	ActorID    *uuid.UUID `json:"-"`
}

type RejectLeaveCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	ApprovalID uuid.UUID  `json:"approval_id"`
	ApproverID uuid.UUID  `json:"approver_id"`
	Remarks    *string    `json:"remarks,omitempty"`
	ActorID    *uuid.UUID `json:"-"`
}

type CancelLeaveCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	LeaveID  uuid.UUID  `json:"leave_id"`
	UserID   uuid.UUID  `json:"user_id"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}
