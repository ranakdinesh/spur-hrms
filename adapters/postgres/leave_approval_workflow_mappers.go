package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapLeaveApprovalWorkflow(row sqlc.HrmsLeaveApprovalWorkflow) *domain.LeaveApprovalWorkflow {
	return &domain.LeaveApprovalWorkflow{
		ID:          row.ID,
		TenantID:    row.TenantID,
		Name:        row.Name,
		Code:        row.Code,
		Description: ptrFromText(row.Description),
		IsDefault:   row.IsDefault,
		Inactive:    row.Inactive,
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
		CreatedBy:   ptrFromUUID(row.CreatedBy),
		UpdatedAt:   timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:   ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeaveApprovalWorkflows(rows []sqlc.HrmsLeaveApprovalWorkflow) []*domain.LeaveApprovalWorkflow {
	items := make([]*domain.LeaveApprovalWorkflow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveApprovalWorkflow(row))
	}
	return items
}

func mapLeaveApprovalWorkflowStep(row sqlc.HrmsLeaveApprovalWorkflowStep) *domain.LeaveApprovalWorkflowStep {
	return &domain.LeaveApprovalWorkflowStep{
		ID:                row.ID,
		TenantID:          row.TenantID,
		WorkflowID:        row.WorkflowID,
		StepOrder:         row.StepOrder,
		Name:              row.Name,
		ApproverType:      row.ApproverType,
		ApproverUserID:    ptrFromUUID(row.ApproverUserID),
		ApproverRole:      ptrFromText(row.ApproverRole),
		DecisionRule:      row.DecisionRule,
		RequiredApprovals: row.RequiredApprovals,
		AutoApprove:       row.AutoApprove,
		SLAHours:          row.SlaHours,
		Inactive:          row.Inactive,
		CreatedAt:         timeFromTimestamptz(row.CreatedAt),
		CreatedBy:         ptrFromUUID(row.CreatedBy),
		UpdatedAt:         timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:         ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeaveApprovalWorkflowSteps(rows []sqlc.HrmsLeaveApprovalWorkflowStep) []*domain.LeaveApprovalWorkflowStep {
	items := make([]*domain.LeaveApprovalWorkflowStep, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveApprovalWorkflowStep(row))
	}
	return items
}
