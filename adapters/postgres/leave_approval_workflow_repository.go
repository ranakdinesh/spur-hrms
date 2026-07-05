package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateLeaveApprovalWorkflow(ctx context.Context, item *domain.LeaveApprovalWorkflow, actorID *uuid.UUID) (*domain.LeaveApprovalWorkflow, error) {
	q := s.getQueries(ctx)
	if item.IsDefault {
		if err := q.ClearDefaultLeaveApprovalWorkflows(ctx, sqlc.ClearDefaultLeaveApprovalWorkflowsParams{TenantID: item.TenantID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
			return nil, s.logDBError(ctx, "clear default leave approval workflows", err, tenantIDField(item.TenantID))
		}
	}
	row, err := q.CreateLeaveApprovalWorkflow(ctx, sqlc.CreateLeaveApprovalWorkflowParams{TenantID: item.TenantID, Name: item.Name, Code: item.Code, Description: textFromPtr(item.Description), IsDefault: item.IsDefault, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create leave approval workflow", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapLeaveApprovalWorkflow(row), nil
}

func (s *Store) ListLeaveApprovalWorkflows(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeaveApprovalWorkflow, error) {
	rows, err := s.getQueries(ctx).ListLeaveApprovalWorkflows(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list leave approval workflows", err, tenantIDField(tenantID))
	}
	return mapLeaveApprovalWorkflows(rows), nil
}

func (s *Store) GetLeaveApprovalWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeaveApprovalWorkflow, error) {
	row, err := s.getQueries(ctx).GetLeaveApprovalWorkflow(ctx, sqlc.GetLeaveApprovalWorkflowParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveApprovalWorkflowNotFound
		}
		return nil, s.logDBError(ctx, "get leave approval workflow", err, tenantIDField(tenantID), stringField("workflow_id", id.String()))
	}
	return mapLeaveApprovalWorkflow(row), nil
}

func (s *Store) GetDefaultLeaveApprovalWorkflow(ctx context.Context, tenantID uuid.UUID) (*domain.LeaveApprovalWorkflow, error) {
	row, err := s.getQueries(ctx).GetDefaultLeaveApprovalWorkflow(ctx, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveApprovalWorkflowNotFound
		}
		return nil, s.logDBError(ctx, "get default leave approval workflow", err, tenantIDField(tenantID))
	}
	return mapLeaveApprovalWorkflow(row), nil
}

func (s *Store) UpdateLeaveApprovalWorkflow(ctx context.Context, item *domain.LeaveApprovalWorkflow, actorID *uuid.UUID) (*domain.LeaveApprovalWorkflow, error) {
	q := s.getQueries(ctx)
	if item.IsDefault {
		if err := q.ClearDefaultLeaveApprovalWorkflows(ctx, sqlc.ClearDefaultLeaveApprovalWorkflowsParams{TenantID: item.TenantID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
			return nil, s.logDBError(ctx, "clear default leave approval workflows", err, tenantIDField(item.TenantID))
		}
	}
	row, err := q.UpdateLeaveApprovalWorkflow(ctx, sqlc.UpdateLeaveApprovalWorkflowParams{TenantID: item.TenantID, ID: item.ID, Name: item.Name, Code: item.Code, Description: textFromPtr(item.Description), IsDefault: item.IsDefault, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveApprovalWorkflowNotFound
		}
		return nil, s.logDBError(ctx, "update leave approval workflow", err, tenantIDField(item.TenantID), stringField("workflow_id", item.ID.String()))
	}
	return mapLeaveApprovalWorkflow(row), nil
}

func (s *Store) DeleteLeaveApprovalWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLeaveApprovalWorkflow(ctx, sqlc.SoftDeleteLeaveApprovalWorkflowParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete leave approval workflow", err, tenantIDField(tenantID), stringField("workflow_id", id.String()))
	}
	return nil
}

func (s *Store) CreateLeaveApprovalWorkflowStep(ctx context.Context, item *domain.LeaveApprovalWorkflowStep, actorID *uuid.UUID) (*domain.LeaveApprovalWorkflowStep, error) {
	row, err := s.getQueries(ctx).CreateLeaveApprovalWorkflowStep(ctx, sqlc.CreateLeaveApprovalWorkflowStepParams{TenantID: item.TenantID, WorkflowID: item.WorkflowID, StepOrder: item.StepOrder, Name: item.Name, ApproverType: item.ApproverType, ApproverUserID: uuidFromPtr(item.ApproverUserID), ApproverRole: textFromPtr(item.ApproverRole), DecisionRule: item.DecisionRule, RequiredApprovals: item.RequiredApprovals, AutoApprove: item.AutoApprove, SlaHours: item.SLAHours, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create leave approval workflow step", err, tenantIDField(item.TenantID), stringField("workflow_id", item.WorkflowID.String()))
	}
	return mapLeaveApprovalWorkflowStep(row), nil
}

func (s *Store) ListLeaveApprovalWorkflowSteps(ctx context.Context, tenantID uuid.UUID, workflowID uuid.UUID) ([]*domain.LeaveApprovalWorkflowStep, error) {
	rows, err := s.getQueries(ctx).ListLeaveApprovalWorkflowSteps(ctx, sqlc.ListLeaveApprovalWorkflowStepsParams{TenantID: tenantID, WorkflowID: workflowID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave approval workflow steps", err, tenantIDField(tenantID), stringField("workflow_id", workflowID.String()))
	}
	return mapLeaveApprovalWorkflowSteps(rows), nil
}

func (s *Store) GetLeaveApprovalWorkflowStep(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeaveApprovalWorkflowStep, error) {
	row, err := s.getQueries(ctx).GetLeaveApprovalWorkflowStep(ctx, sqlc.GetLeaveApprovalWorkflowStepParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveApprovalStepNotFound
		}
		return nil, s.logDBError(ctx, "get leave approval workflow step", err, tenantIDField(tenantID), stringField("step_id", id.String()))
	}
	return mapLeaveApprovalWorkflowStep(row), nil
}

func (s *Store) UpdateLeaveApprovalWorkflowStep(ctx context.Context, item *domain.LeaveApprovalWorkflowStep, actorID *uuid.UUID) (*domain.LeaveApprovalWorkflowStep, error) {
	row, err := s.getQueries(ctx).UpdateLeaveApprovalWorkflowStep(ctx, sqlc.UpdateLeaveApprovalWorkflowStepParams{TenantID: item.TenantID, ID: item.ID, StepOrder: item.StepOrder, Name: item.Name, ApproverType: item.ApproverType, ApproverUserID: uuidFromPtr(item.ApproverUserID), ApproverRole: textFromPtr(item.ApproverRole), DecisionRule: item.DecisionRule, RequiredApprovals: item.RequiredApprovals, AutoApprove: item.AutoApprove, SlaHours: item.SLAHours, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveApprovalStepNotFound
		}
		return nil, s.logDBError(ctx, "update leave approval workflow step", err, tenantIDField(item.TenantID), stringField("step_id", item.ID.String()))
	}
	return mapLeaveApprovalWorkflowStep(row), nil
}

func (s *Store) DeleteLeaveApprovalWorkflowStep(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLeaveApprovalWorkflowStep(ctx, sqlc.SoftDeleteLeaveApprovalWorkflowStepParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete leave approval workflow step", err, tenantIDField(tenantID), stringField("step_id", id.String()))
	}
	return nil
}

func (s *Store) CreateWorkflowLeaveApproval(ctx context.Context, item *domain.LeaveApproval, actorID *uuid.UUID) (*domain.LeaveApproval, error) {
	row, err := s.getQueries(ctx).CreateWorkflowLeaveApproval(ctx, sqlc.CreateWorkflowLeaveApprovalParams{TenantID: item.TenantID, LeaveID: item.LeaveID, ApproverID: item.ApproverID, Status: item.Status, Remarks: textFromPtr(item.Remarks), WorkflowID: uuidFromPtr(item.WorkflowID), WorkflowStepID: uuidFromPtr(item.WorkflowStepID), StepOrder: item.StepOrder, DecisionRule: item.DecisionRule, RequiredApprovals: item.RequiredApprovals, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create workflow leave approval", err, tenantIDField(item.TenantID), stringField("leave_id", item.LeaveID.String()))
	}
	return mapLeaveApproval(row), nil
}

func (s *Store) GetLeaveApproval(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeaveApproval, error) {
	row, err := s.getQueries(ctx).GetLeaveApproval(ctx, sqlc.GetLeaveApprovalParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveApprovalNotFound
		}
		return nil, s.logDBError(ctx, "get leave approval", err, tenantIDField(tenantID), stringField("approval_id", id.String()))
	}
	return mapLeaveApproval(row), nil
}

func (s *Store) ListLeaveApprovalsByLeave(ctx context.Context, tenantID uuid.UUID, leaveID uuid.UUID) ([]*domain.LeaveApproval, error) {
	rows, err := s.getQueries(ctx).ListLeaveApprovalsByLeave(ctx, sqlc.ListLeaveApprovalsByLeaveParams{TenantID: tenantID, LeaveID: leaveID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave approvals by leave", err, tenantIDField(tenantID), stringField("leave_id", leaveID.String()))
	}
	return mapLeaveApprovals(rows), nil
}

func (s *Store) UpdateLeaveApprovalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.LeaveApproval, error) {
	row, err := s.getQueries(ctx).UpdateLeaveApprovalStatus(ctx, sqlc.UpdateLeaveApprovalStatusParams{TenantID: tenantID, ID: id, Status: status, Remarks: textFromPtr(remarks), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveApprovalNotFound
		}
		return nil, s.logDBError(ctx, "update leave approval status", err, tenantIDField(tenantID), stringField("approval_id", id.String()))
	}
	return mapLeaveApproval(row), nil
}

func (s *Store) UpdateLeaveStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.Leave, error) {
	row, err := s.getQueries(ctx).UpdateLeaveStatus(ctx, sqlc.UpdateLeaveStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveNotFound
		}
		return nil, s.logDBError(ctx, "update leave status", err, tenantIDField(tenantID), stringField("leave_id", id.String()))
	}
	return mapLeave(row), nil
}

func (s *Store) ListPendingApprovalsByApprover(ctx context.Context, tenantID uuid.UUID, approverID uuid.UUID) ([]*domain.LeaveApproval, error) {
	rows, err := s.getQueries(ctx).ListPendingApprovalsByApprover(ctx, sqlc.ListPendingApprovalsByApproverParams{TenantID: tenantID, ApproverID: approverID})
	if err != nil {
		return nil, s.logDBError(ctx, "list pending approvals by approver", err, tenantIDField(tenantID), stringField("approver_id", approverID.String()))
	}
	return mapLeaveApprovals(rows), nil
}
