package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateLeaveApprovalWorkflow(ctx context.Context, cmd ports.LeaveApprovalWorkflowCommand) (*domain.LeaveApprovalWorkflow, error) {
	item := &domain.LeaveApprovalWorkflow{TenantID: cmd.TenantID, Name: cmd.Name, Code: cmd.Code, Description: cmd.Description, IsDefault: cmd.IsDefault}
	if err := domain.ValidateLeaveApprovalWorkflow(item); err != nil {
		s.logError("validate create leave approval workflow", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	var created *domain.LeaveApprovalWorkflow
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		created, err = s.approvalWorkflows.CreateLeaveApprovalWorkflow(txCtx, item, cmd.ActorID)
		return err
	})
	if err != nil {
		s.logError("create leave approval workflow", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", cmd.Code))
		return nil, err
	}
	return created, nil
}

func (s *TenantService) ListLeaveApprovalWorkflows(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeaveApprovalWorkflow, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave approval workflow list tenant", err)
		return nil, err
	}
	items, err := s.approvalWorkflows.ListLeaveApprovalWorkflows(ctx, tenantID)
	if err != nil {
		s.logError("list leave approval workflows", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateLeaveApprovalWorkflow(ctx context.Context, cmd ports.LeaveApprovalWorkflowCommand) (*domain.LeaveApprovalWorkflow, error) {
	item := &domain.LeaveApprovalWorkflow{ID: cmd.ID, TenantID: cmd.TenantID, Name: cmd.Name, Code: cmd.Code, Description: cmd.Description, IsDefault: cmd.IsDefault}
	if item.ID == uuid.Nil {
		err := domain.ErrInvalidLeaveApprovalWorkflowID
		s.logError("validate update leave approval workflow id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if err := domain.ValidateLeaveApprovalWorkflow(item); err != nil {
		s.logError("validate update leave approval workflow", err, serviceTenantIDField(cmd.TenantID), serviceStringField("workflow_id", cmd.ID.String()))
		return nil, err
	}
	var updated *domain.LeaveApprovalWorkflow
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		updated, err = s.approvalWorkflows.UpdateLeaveApprovalWorkflow(txCtx, item, cmd.ActorID)
		return err
	})
	if err != nil {
		s.logError("update leave approval workflow", err, serviceTenantIDField(cmd.TenantID), serviceStringField("workflow_id", cmd.ID.String()))
		return nil, err
	}
	return updated, nil
}

func (s *TenantService) DeleteLeaveApprovalWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate delete leave approval workflow tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidLeaveApprovalWorkflowID
		s.logError("validate delete leave approval workflow id", err, serviceTenantIDField(tenantID))
		return err
	}
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		return s.approvalWorkflows.DeleteLeaveApprovalWorkflow(txCtx, tenantID, id, actorID)
	})
	if err != nil {
		s.logError("delete leave approval workflow", err, serviceTenantIDField(tenantID), serviceStringField("workflow_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateLeaveApprovalWorkflowStep(ctx context.Context, cmd ports.LeaveApprovalWorkflowStepCommand) (*domain.LeaveApprovalWorkflowStep, error) {
	item := leaveApprovalStepFromCommand(cmd)
	if err := domain.ValidateLeaveApprovalWorkflowStep(item); err != nil {
		s.logError("validate create leave approval workflow step", err, serviceTenantIDField(cmd.TenantID), serviceStringField("workflow_id", cmd.WorkflowID.String()))
		return nil, err
	}
	var created *domain.LeaveApprovalWorkflowStep
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		if _, err := s.approvalWorkflows.GetLeaveApprovalWorkflow(txCtx, cmd.TenantID, cmd.WorkflowID); err != nil {
			return err
		}
		var err error
		created, err = s.approvalWorkflows.CreateLeaveApprovalWorkflowStep(txCtx, item, cmd.ActorID)
		return err
	})
	if err != nil {
		s.logError("create leave approval workflow step", err, serviceTenantIDField(cmd.TenantID), serviceStringField("workflow_id", cmd.WorkflowID.String()))
		return nil, err
	}
	return created, nil
}

func (s *TenantService) ListLeaveApprovalWorkflowSteps(ctx context.Context, tenantID uuid.UUID, workflowID uuid.UUID) ([]*domain.LeaveApprovalWorkflowStep, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave approval workflow step list tenant", err)
		return nil, err
	}
	if workflowID == uuid.Nil {
		err := domain.ErrInvalidLeaveApprovalWorkflowID
		s.logError("validate leave approval workflow step list workflow", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.approvalWorkflows.ListLeaveApprovalWorkflowSteps(ctx, tenantID, workflowID)
	if err != nil {
		s.logError("list leave approval workflow steps", err, serviceTenantIDField(tenantID), serviceStringField("workflow_id", workflowID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateLeaveApprovalWorkflowStep(ctx context.Context, cmd ports.LeaveApprovalWorkflowStepCommand) (*domain.LeaveApprovalWorkflowStep, error) {
	item := leaveApprovalStepFromCommand(cmd)
	if item.ID == uuid.Nil {
		err := domain.ErrInvalidLeaveApprovalStepID
		s.logError("validate update leave approval workflow step id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if err := domain.ValidateLeaveApprovalWorkflowStep(item); err != nil {
		s.logError("validate update leave approval workflow step", err, serviceTenantIDField(cmd.TenantID), serviceStringField("step_id", cmd.ID.String()))
		return nil, err
	}
	var updated *domain.LeaveApprovalWorkflowStep
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		updated, err = s.approvalWorkflows.UpdateLeaveApprovalWorkflowStep(txCtx, item, cmd.ActorID)
		return err
	})
	if err != nil {
		s.logError("update leave approval workflow step", err, serviceTenantIDField(cmd.TenantID), serviceStringField("step_id", cmd.ID.String()))
		return nil, err
	}
	return updated, nil
}

func (s *TenantService) DeleteLeaveApprovalWorkflowStep(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate delete leave approval workflow step tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidLeaveApprovalStepID
		s.logError("validate delete leave approval workflow step id", err, serviceTenantIDField(tenantID))
		return err
	}
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		return s.approvalWorkflows.DeleteLeaveApprovalWorkflowStep(txCtx, tenantID, id, actorID)
	})
	if err != nil {
		s.logError("delete leave approval workflow step", err, serviceTenantIDField(tenantID), serviceStringField("step_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListPendingApprovalsByApprover(ctx context.Context, tenantID uuid.UUID, approverID uuid.UUID) ([]*domain.LeaveApproval, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate pending approvals tenant", err)
		return nil, err
	}
	if approverID == uuid.Nil {
		err := domain.ErrInvalidLeaveApproverType
		s.logError("validate pending approvals approver", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.approvalWorkflows.ListPendingApprovalsByApprover(ctx, tenantID, approverID)
	if err != nil {
		s.logError("list pending approvals by approver", err, serviceTenantIDField(tenantID), serviceStringField("approver_id", approverID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ApproveLeave(ctx context.Context, cmd ports.ApproveLeaveCommand) (*domain.LeaveApplication, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate approve leave tenant", err)
		return nil, err
	}
	if cmd.ApprovalID == uuid.Nil {
		err := domain.ErrLeaveApprovalNotFound
		s.logError("validate approve leave approval id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.ApproverID == uuid.Nil {
		err := domain.ErrLeaveApprovalUnauthorized
		s.logError("validate approve leave approver", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	var application *domain.LeaveApplication
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		approval, err := s.approvalWorkflows.GetLeaveApproval(txCtx, cmd.TenantID, cmd.ApprovalID)
		if err != nil {
			return err
		}
		if approval.Status != domain.LeaveStatusPending {
			return domain.ErrLeaveApprovalNotPending
		}
		if approval.ApproverID != cmd.ApproverID {
			return domain.ErrLeaveApprovalUnauthorized
		}
		updatedApproval, err := s.approvalWorkflows.UpdateLeaveApprovalStatus(txCtx, cmd.TenantID, cmd.ApprovalID, domain.LeaveStatusApproved, cmd.Remarks, cmd.ActorID)
		if err != nil {
			return err
		}
		leave, err := s.leaveRequests.GetLeave(txCtx, cmd.TenantID, approval.LeaveID)
		if err != nil {
			return err
		}
		if leave.Status != domain.LeaveStatusPending {
			application = &domain.LeaveApplication{Leave: leave, Approval: updatedApproval}
			return nil
		}
		approvals, err := s.approvalWorkflows.ListLeaveApprovalsByLeave(txCtx, cmd.TenantID, leave.ID)
		if err != nil {
			return err
		}
		if !leaveApprovalStepComplete(approvals, updatedApproval) {
			application = &domain.LeaveApplication{Leave: leave, Approval: updatedApproval}
			return nil
		}
		nextCreated, err := s.createNextLeaveApprovalStep(txCtx, leave, approvals, updatedApproval, cmd.ActorID)
		if err != nil {
			return err
		}
		if nextCreated != nil {
			application = &domain.LeaveApplication{Leave: leave, Approval: nextCreated}
			return nil
		}
		approvedLeave, balance, err := s.finalizeApprovedLeave(txCtx, leave, cmd.ActorID)
		if err != nil {
			return err
		}
		application = &domain.LeaveApplication{Leave: approvedLeave, Approval: updatedApproval, Balance: balance}
		return nil
	})
	if err != nil {
		s.logError("approve leave", err, serviceTenantIDField(cmd.TenantID), serviceStringField("approval_id", cmd.ApprovalID.String()))
		return nil, err
	}
	if application != nil && application.Leave != nil {
		if application.Leave.Status == domain.LeaveStatusApproved {
			s.notifyLeaveReviewed(ctx, application, domain.NotifLeaveApproved, cmd.ActorID)
		} else if application.Approval != nil && application.Approval.Status == domain.LeaveStatusPending {
			s.notifyLeaveApplied(ctx, application, cmd.ActorID)
		}
	}
	return application, nil
}

func (s *TenantService) RejectLeave(ctx context.Context, cmd ports.RejectLeaveCommand) (*domain.LeaveApplication, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate reject leave tenant", err)
		return nil, err
	}
	if cmd.ApprovalID == uuid.Nil {
		err := domain.ErrLeaveApprovalNotFound
		s.logError("validate reject leave approval id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.ApproverID == uuid.Nil {
		err := domain.ErrLeaveApprovalUnauthorized
		s.logError("validate reject leave approver", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	var application *domain.LeaveApplication
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		approval, err := s.approvalWorkflows.GetLeaveApproval(txCtx, cmd.TenantID, cmd.ApprovalID)
		if err != nil {
			return err
		}
		if approval.Status != domain.LeaveStatusPending {
			return domain.ErrLeaveApprovalNotPending
		}
		if approval.ApproverID != cmd.ApproverID {
			return domain.ErrLeaveApprovalUnauthorized
		}
		updatedApproval, err := s.approvalWorkflows.UpdateLeaveApprovalStatus(txCtx, cmd.TenantID, cmd.ApprovalID, domain.LeaveStatusRejected, cmd.Remarks, cmd.ActorID)
		if err != nil {
			return err
		}
		leave, err := s.leaveRequests.GetLeave(txCtx, cmd.TenantID, approval.LeaveID)
		if err != nil {
			return err
		}
		if leave.Status != domain.LeaveStatusPending {
			return domain.ErrLeaveNotPending
		}
		rejectedLeave, balance, err := s.finalizeRejectedOrCanceledLeave(txCtx, leave, domain.LeaveStatusRejected, domain.LeaveLedgerSourceLeaveReject, "leave rejected", cmd.ActorID)
		if err != nil {
			return err
		}
		application = &domain.LeaveApplication{Leave: rejectedLeave, Approval: updatedApproval, Balance: balance}
		return nil
	})
	if err != nil {
		s.logError("reject leave", err, serviceTenantIDField(cmd.TenantID), serviceStringField("approval_id", cmd.ApprovalID.String()))
		return nil, err
	}
	s.notifyLeaveReviewed(ctx, application, domain.NotifLeaveRejected, cmd.ActorID)
	return application, nil
}

func (s *TenantService) CancelLeave(ctx context.Context, cmd ports.CancelLeaveCommand) (*domain.LeaveApplication, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate cancel leave tenant", err)
		return nil, err
	}
	if cmd.LeaveID == uuid.Nil {
		err := domain.ErrInvalidLeaveID
		s.logError("validate cancel leave id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.UserID == uuid.Nil && cmd.ActorID != nil {
		cmd.UserID = *cmd.ActorID
	}
	if cmd.UserID == uuid.Nil {
		err := domain.ErrInvalidLeaveUser
		s.logError("validate cancel leave user", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	var application *domain.LeaveApplication
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		leave, err := s.leaveRequests.GetLeave(txCtx, cmd.TenantID, cmd.LeaveID)
		if err != nil {
			return err
		}
		if leave.Status != domain.LeaveStatusPending {
			return domain.ErrLeaveNotPending
		}
		if leave.UserID != cmd.UserID {
			return domain.ErrLeaveCancelUnauthorized
		}
		canceledLeave, balance, err := s.finalizeRejectedOrCanceledLeave(txCtx, leave, domain.LeaveStatusCanceled, domain.LeaveLedgerSourceLeaveCancel, "leave canceled", cmd.ActorID)
		if err != nil {
			return err
		}
		approvals, err := s.approvalWorkflows.ListLeaveApprovalsByLeave(txCtx, cmd.TenantID, leave.ID)
		if err != nil {
			return err
		}
		var updatedApproval *domain.LeaveApproval
		remarks := cmd.Remarks
		if remarks == nil {
			defaultRemarks := "leave canceled by applicant"
			remarks = &defaultRemarks
		}
		for _, approval := range approvals {
			if approval.Status != domain.LeaveStatusPending {
				continue
			}
			item, err := s.approvalWorkflows.UpdateLeaveApprovalStatus(txCtx, cmd.TenantID, approval.ID, domain.LeaveStatusCanceled, remarks, cmd.ActorID)
			if err != nil {
				return err
			}
			if updatedApproval == nil {
				updatedApproval = item
			}
		}
		application = &domain.LeaveApplication{Leave: canceledLeave, Approval: updatedApproval, Balance: balance}
		return nil
	})
	if err != nil {
		s.logError("cancel leave", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_id", cmd.LeaveID.String()))
		return nil, err
	}
	return application, nil
}

func leaveApprovalStepFromCommand(cmd ports.LeaveApprovalWorkflowStepCommand) *domain.LeaveApprovalWorkflowStep {
	return &domain.LeaveApprovalWorkflowStep{ID: cmd.ID, TenantID: cmd.TenantID, WorkflowID: cmd.WorkflowID, StepOrder: cmd.StepOrder, Name: cmd.Name, ApproverType: cmd.ApproverType, ApproverUserID: cmd.ApproverUserID, ApproverRole: cmd.ApproverRole, DecisionRule: cmd.DecisionRule, RequiredApprovals: cmd.RequiredApprovals, AutoApprove: cmd.AutoApprove, SLAHours: cmd.SLAHours}
}

func leaveApprovalStepComplete(approvals []*domain.LeaveApproval, current *domain.LeaveApproval) bool {
	if current == nil {
		return false
	}
	total, approved := int32(0), int32(0)
	required := current.RequiredApprovals
	if required <= 0 {
		required = 1
	}
	for _, approval := range approvals {
		if approval.LeaveID != current.LeaveID || approval.StepOrder != current.StepOrder {
			continue
		}
		total++
		if approval.Status == domain.LeaveStatusApproved || approval.ID == current.ID {
			approved++
		}
	}
	if current.DecisionRule == domain.LeaveApprovalDecisionAny {
		return approved >= required
	}
	return approved >= required && approved >= total
}

func (s *TenantService) createNextLeaveApprovalStep(ctx context.Context, leave *domain.Leave, approvals []*domain.LeaveApproval, current *domain.LeaveApproval, actorID *uuid.UUID) (*domain.LeaveApproval, error) {
	if current.WorkflowID == nil || *current.WorkflowID == uuid.Nil {
		return nil, nil
	}
	steps, err := s.approvalWorkflows.ListLeaveApprovalWorkflowSteps(ctx, leave.TenantID, *current.WorkflowID)
	if err != nil {
		return nil, err
	}
	var next *domain.LeaveApprovalWorkflowStep
	for _, step := range steps {
		if step.StepOrder <= current.StepOrder || leaveApprovalStepAlreadyCreated(approvals, step.StepOrder) {
			continue
		}
		if next == nil || step.StepOrder < next.StepOrder {
			next = step
		}
	}
	if next == nil {
		return nil, nil
	}
	return s.createApprovalsForStep(ctx, leave, next, actorID)
}

func leaveApprovalStepAlreadyCreated(approvals []*domain.LeaveApproval, stepOrder int32) bool {
	for _, approval := range approvals {
		if approval.StepOrder == stepOrder {
			return true
		}
	}
	return false
}

func (s *TenantService) createApprovalsForStep(ctx context.Context, leave *domain.Leave, step *domain.LeaveApprovalWorkflowStep, actorID *uuid.UUID) (*domain.LeaveApproval, error) {
	if step.AutoApprove {
		return nil, nil
	}
	approvers, err := s.resolveLeaveApprovers(ctx, leave.TenantID, leave.UserID, step)
	if err != nil {
		return nil, err
	}
	var first *domain.LeaveApproval
	for _, approverID := range approvers {
		workflowID := step.WorkflowID
		stepID := step.ID
		item := &domain.LeaveApproval{TenantID: leave.TenantID, LeaveID: leave.ID, ApproverID: approverID, Status: domain.LeaveStatusPending, WorkflowID: &workflowID, WorkflowStepID: &stepID, StepOrder: step.StepOrder, DecisionRule: step.DecisionRule, RequiredApprovals: step.RequiredApprovals}
		created, err := s.approvalWorkflows.CreateWorkflowLeaveApproval(ctx, item, actorID)
		if err != nil {
			return nil, err
		}
		if first == nil {
			first = created
		}
	}
	return first, nil
}

func (s *TenantService) resolveLeaveApprovers(ctx context.Context, tenantID uuid.UUID, applicantID uuid.UUID, step *domain.LeaveApprovalWorkflowStep) ([]uuid.UUID, error) {
	employee, err := s.employees.GetEmployeeByUserID(ctx, tenantID, applicantID)
	if err != nil {
		return nil, err
	}
	switch step.ApproverType {
	case domain.LeaveApproverReportingManager:
		if employee.ReportingManagerID != nil && *employee.ReportingManagerID != uuid.Nil {
			return []uuid.UUID{*employee.ReportingManagerID}, nil
		}
	case domain.LeaveApproverManagerManager:
		if employee.ReportingManagerID != nil && *employee.ReportingManagerID != uuid.Nil {
			manager, err := s.employees.GetEmployeeByUserID(ctx, tenantID, *employee.ReportingManagerID)
			if err != nil {
				return nil, err
			}
			if manager.ReportingManagerID != nil && *manager.ReportingManagerID != uuid.Nil {
				return []uuid.UUID{*manager.ReportingManagerID}, nil
			}
		}
	case domain.LeaveApproverHRUser, domain.LeaveApproverSpecificUser:
		if step.ApproverUserID != nil && *step.ApproverUserID != uuid.Nil {
			return []uuid.UUID{*step.ApproverUserID}, nil
		}
	case domain.LeaveApproverApplicant:
		return []uuid.UUID{applicantID}, nil
	case domain.LeaveApproverRole:
		if step.ApproverUserID != nil && *step.ApproverUserID != uuid.Nil {
			return []uuid.UUID{*step.ApproverUserID}, nil
		}
		s.log.Warn().Str("tenant_id", tenantID.String()).Str("role", stringValue(step.ApproverRole)).Msg("hrms: role-based leave approver resolution requires identity role lookup; configure a specific approver user for now")
	}
	return nil, domain.ErrLeaveApproverNotResolved
}

func (s *TenantService) createInitialLeaveApproval(ctx context.Context, leave *domain.Leave, fallbackApproverID uuid.UUID, actorID *uuid.UUID) (*domain.LeaveApproval, error) {
	workflow, err := s.approvalWorkflows.GetDefaultLeaveApprovalWorkflow(ctx, leave.TenantID)
	if err != nil && !errors.Is(err, domain.ErrLeaveApprovalWorkflowNotFound) {
		return nil, err
	}
	if workflow != nil {
		steps, err := s.approvalWorkflows.ListLeaveApprovalWorkflowSteps(ctx, leave.TenantID, workflow.ID)
		if err != nil {
			return nil, err
		}
		for _, step := range steps {
			approval, err := s.createApprovalsForStep(ctx, leave, step, actorID)
			if err != nil {
				return nil, err
			}
			if approval != nil {
				return approval, nil
			}
		}
	}
	approval := &domain.LeaveApproval{TenantID: leave.TenantID, LeaveID: leave.ID, ApproverID: fallbackApproverID, Status: domain.LeaveStatusPending, StepOrder: 1, DecisionRule: domain.LeaveApprovalDecisionAll, RequiredApprovals: 1}
	return s.leaveRequests.CreateLeaveApproval(ctx, approval, actorID)
}

func (s *TenantService) finalizeApprovedLeave(ctx context.Context, leave *domain.Leave, actorID *uuid.UUID) (*domain.Leave, *domain.LeaveBalance, error) {
	leaveType, err := s.leaveTypes.GetLeaveType(ctx, leave.TenantID, leave.LeaveTypeID)
	if err != nil {
		return nil, nil, err
	}
	var balance *domain.LeaveBalance
	if leaveType.IsPaid {
		before, err := s.leaveBalances.GetLeaveBalance(ctx, leave.TenantID, leave.UserID, leave.LeaveTypeID, leave.FYID)
		if err != nil {
			return nil, nil, err
		}
		balance, err = s.leaveBalances.MoveLeaveBalancePendingToUsed(ctx, leave.TenantID, leave.UserID, leave.LeaveTypeID, leave.FYID, leave.Days, actorID)
		if err != nil {
			return nil, nil, err
		}
		remarks := "leave approved"
		ledger := &domain.LeaveLedgerEntry{TenantID: leave.TenantID, UserID: leave.UserID, LeaveTypeID: leave.LeaveTypeID, FYID: leave.FYID, LeaveID: &leave.ID, TransactionType: domain.LeaveLedgerDebit, Days: leave.Days, Remarks: &remarks, SourceType: domain.LeaveLedgerSourceLeaveApprove, SourceID: &leave.ID, BalanceBefore: &before.BalanceDays, BalanceAfter: &balance.BalanceDays, PendingBefore: &before.PendingDays, PendingAfter: &balance.PendingDays, UsedBefore: &before.UsedDays, UsedAfter: &balance.UsedDays}
		if _, err := s.leaveBalances.CreateLeaveLedgerEntry(ctx, ledger, actorID); err != nil {
			return nil, nil, err
		}
	}
	approvedLeave, err := s.approvalWorkflows.UpdateLeaveStatus(ctx, leave.TenantID, leave.ID, domain.LeaveStatusApproved, actorID)
	if err != nil {
		return nil, nil, err
	}
	return approvedLeave, balance, nil
}

func (s *TenantService) finalizeRejectedOrCanceledLeave(ctx context.Context, leave *domain.Leave, status string, sourceType string, defaultRemarks string, actorID *uuid.UUID) (*domain.Leave, *domain.LeaveBalance, error) {
	leaveType, err := s.leaveTypes.GetLeaveType(ctx, leave.TenantID, leave.LeaveTypeID)
	if err != nil {
		return nil, nil, err
	}
	var balance *domain.LeaveBalance
	if leaveType.IsPaid {
		before, err := s.leaveBalances.GetLeaveBalance(ctx, leave.TenantID, leave.UserID, leave.LeaveTypeID, leave.FYID)
		if err != nil {
			return nil, nil, err
		}
		balance, err = s.leaveBalances.ReverseLeaveBalancePending(ctx, leave.TenantID, leave.UserID, leave.LeaveTypeID, leave.FYID, leave.Days, actorID)
		if err != nil {
			return nil, nil, err
		}
		remarks := defaultRemarks
		ledger := &domain.LeaveLedgerEntry{TenantID: leave.TenantID, UserID: leave.UserID, LeaveTypeID: leave.LeaveTypeID, FYID: leave.FYID, LeaveID: &leave.ID, TransactionType: domain.LeaveLedgerCredit, Days: leave.Days, Remarks: &remarks, SourceType: sourceType, SourceID: &leave.ID, BalanceBefore: &before.BalanceDays, BalanceAfter: &balance.BalanceDays, PendingBefore: &before.PendingDays, PendingAfter: &balance.PendingDays, UsedBefore: &before.UsedDays, UsedAfter: &balance.UsedDays}
		if _, err := s.leaveBalances.CreateLeaveLedgerEntry(ctx, ledger, actorID); err != nil {
			return nil, nil, err
		}
	}
	updatedLeave, err := s.approvalWorkflows.UpdateLeaveStatus(ctx, leave.TenantID, leave.ID, status, actorID)
	if err != nil {
		return nil, nil, err
	}
	return updatedLeave, balance, nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
