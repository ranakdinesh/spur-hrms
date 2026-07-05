package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapProject(row sqlc.HrmsProject) *domain.Project {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.Project{
		ID:               row.ID,
		TenantID:         row.TenantID,
		ProjectCode:      ptrFromText(row.ProjectCode),
		Name:             row.Name,
		Description:      ptrFromText(row.Description),
		Status:           row.Status,
		DepartmentID:     ptrFromUUID(row.DepartmentID),
		BranchID:         ptrFromUUID(row.BranchID),
		ProjectManagerID: ptrFromUUID(row.ProjectManagerID),
		StartDate:        ptrFromDate(row.StartDate),
		DueDate:          ptrFromDate(row.DueDate),
		CompletedAt:      ptrFromTimestamptz(row.CompletedAt),
		BudgetAmount:     ptrFromNumeric(row.BudgetAmount),
		CurrencyCode:     row.CurrencyCode,
		BillingType:      row.BillingType,
		ClientLabel:      ptrFromText(row.ClientLabel),
		Priority:         row.Priority,
		Notes:            ptrFromText(row.Notes),
		Metadata:         metadata,
		Inactive:         row.Inactive,
		CreatedAt:        timeFromTimestamptz(row.CreatedAt),
		CreatedBy:        ptrFromUUID(row.CreatedBy),
		UpdatedAt:        timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:        ptrFromUUID(row.UpdatedBy),
	}
}

func mapProjectListItems(rows []sqlc.ListProjectsRow) []*domain.ProjectListItem {
	items := make([]*domain.ProjectListItem, 0, len(rows))
	for _, row := range rows {
		metadata := json.RawMessage(row.Metadata)
		if len(metadata) == 0 {
			metadata = json.RawMessage(`{}`)
		}
		accepted := floatFromNumeric(row.AcceptedAmount)
		budget := ptrFromNumeric(row.BudgetAmount)
		items = append(items, &domain.ProjectListItem{
			Project: domain.Project{
				ID:               row.ID,
				TenantID:         row.TenantID,
				ProjectCode:      ptrFromText(row.ProjectCode),
				Name:             row.Name,
				Description:      ptrFromText(row.Description),
				Status:           row.Status,
				DepartmentID:     ptrFromUUID(row.DepartmentID),
				BranchID:         ptrFromUUID(row.BranchID),
				ProjectManagerID: ptrFromUUID(row.ProjectManagerID),
				StartDate:        ptrFromDate(row.StartDate),
				DueDate:          ptrFromDate(row.DueDate),
				CompletedAt:      ptrFromTimestamptz(row.CompletedAt),
				BudgetAmount:     budget,
				CurrencyCode:     row.CurrencyCode,
				BillingType:      row.BillingType,
				ClientLabel:      ptrFromText(row.ClientLabel),
				Priority:         row.Priority,
				Notes:            ptrFromText(row.Notes),
				Metadata:         metadata,
				Inactive:         row.Inactive,
				CreatedAt:        timeFromTimestamptz(row.CreatedAt),
				CreatedBy:        ptrFromUUID(row.CreatedBy),
				UpdatedAt:        timeFromTimestamptz(row.UpdatedAt),
				UpdatedBy:        ptrFromUUID(row.UpdatedBy),
			},
			DepartmentName:          ptrFromText(row.DepartmentName),
			BranchName:              ptrFromText(row.BranchName),
			MilestoneCount:          row.MilestoneCount,
			SubmittedMilestoneCount: row.SubmittedMilestoneCount,
			AcceptedMilestoneCount:  row.AcceptedMilestoneCount,
			RejectedMilestoneCount:  row.RejectedMilestoneCount,
			MilestoneAmount:         floatFromNumeric(row.MilestoneAmount),
			AcceptedAmount:          accepted,
			RemainingBudgetAmount:   domain.ProjectBudgetRemaining(budget, accepted),
		})
	}
	return items
}

func mapProjectMilestone(row sqlc.HrmsProjectMilestone) *domain.ProjectMilestone {
	paymentTrigger := json.RawMessage(row.PaymentTrigger)
	if len(paymentTrigger) == 0 {
		paymentTrigger = json.RawMessage(`{}`)
	}
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.ProjectMilestone{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		ProjectID:          row.ProjectID,
		EngagementID:       ptrFromUUID(row.EngagementID),
		MilestoneCode:      ptrFromText(row.MilestoneCode),
		Title:              row.Title,
		Description:        ptrFromText(row.Description),
		AcceptanceCriteria: ptrFromText(row.AcceptanceCriteria),
		DueDate:            ptrFromDate(row.DueDate),
		Status:             row.Status,
		Amount:             ptrFromNumeric(row.Amount),
		CurrencyCode:       row.CurrencyCode,
		PaymentTrigger:     paymentTrigger,
		SubmittedAt:        ptrFromTimestamptz(row.SubmittedAt),
		SubmittedBy:        ptrFromUUID(row.SubmittedBy),
		AcceptedAt:         ptrFromTimestamptz(row.AcceptedAt),
		AcceptedBy:         ptrFromUUID(row.AcceptedBy),
		RejectedAt:         ptrFromTimestamptz(row.RejectedAt),
		RejectedBy:         ptrFromUUID(row.RejectedBy),
		ReviewComment:      ptrFromText(row.ReviewComment),
		Notes:              ptrFromText(row.Notes),
		Metadata:           metadata,
		Inactive:           row.Inactive,
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		CreatedBy:          ptrFromUUID(row.CreatedBy),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:          ptrFromUUID(row.UpdatedBy),
	}
}

func mapProjectMilestoneListItems(rows []sqlc.ListProjectMilestonesRow) []*domain.ProjectMilestoneListItem {
	items := make([]*domain.ProjectMilestoneListItem, 0, len(rows))
	for _, row := range rows {
		paymentTrigger := json.RawMessage(row.PaymentTrigger)
		if len(paymentTrigger) == 0 {
			paymentTrigger = json.RawMessage(`{}`)
		}
		metadata := json.RawMessage(row.Metadata)
		if len(metadata) == 0 {
			metadata = json.RawMessage(`{}`)
		}
		items = append(items, &domain.ProjectMilestoneListItem{
			ProjectMilestone: domain.ProjectMilestone{
				ID:                 row.ID,
				TenantID:           row.TenantID,
				ProjectID:          row.ProjectID,
				EngagementID:       ptrFromUUID(row.EngagementID),
				MilestoneCode:      ptrFromText(row.MilestoneCode),
				Title:              row.Title,
				Description:        ptrFromText(row.Description),
				AcceptanceCriteria: ptrFromText(row.AcceptanceCriteria),
				DueDate:            ptrFromDate(row.DueDate),
				Status:             row.Status,
				Amount:             ptrFromNumeric(row.Amount),
				CurrencyCode:       row.CurrencyCode,
				PaymentTrigger:     paymentTrigger,
				SubmittedAt:        ptrFromTimestamptz(row.SubmittedAt),
				SubmittedBy:        ptrFromUUID(row.SubmittedBy),
				AcceptedAt:         ptrFromTimestamptz(row.AcceptedAt),
				AcceptedBy:         ptrFromUUID(row.AcceptedBy),
				RejectedAt:         ptrFromTimestamptz(row.RejectedAt),
				RejectedBy:         ptrFromUUID(row.RejectedBy),
				ReviewComment:      ptrFromText(row.ReviewComment),
				Notes:              ptrFromText(row.Notes),
				Metadata:           metadata,
				Inactive:           row.Inactive,
				CreatedAt:          timeFromTimestamptz(row.CreatedAt),
				CreatedBy:          ptrFromUUID(row.CreatedBy),
				UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
				UpdatedBy:          ptrFromUUID(row.UpdatedBy),
			},
			ProjectName:       row.ProjectName,
			ProjectCode:       ptrFromText(row.ProjectCode),
			ProjectManagerID:  ptrFromUUID(row.ProjectManagerID),
			DepartmentID:      ptrFromUUID(row.DepartmentID),
			DepartmentName:    ptrFromText(row.DepartmentName),
			EngagementTitle:   ptrFromText(row.EngagementTitle),
			EngagementCode:    ptrFromText(row.EngagementCode),
			WorkerProfileID:   ptrFromUUID(row.WorkerProfileID),
			WorkerDisplayName: ptrFromText(row.WorkerDisplayName),
			WorkerCode:        ptrFromText(row.WorkerCode),
		})
	}
	return items
}

func mapProjectMilestoneEvent(row sqlc.HrmsProjectMilestoneEvent) *domain.ProjectMilestoneEvent {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.ProjectMilestoneEvent{
		ID:          row.ID,
		TenantID:    row.TenantID,
		ProjectID:   row.ProjectID,
		MilestoneID: row.MilestoneID,
		EventType:   row.EventType,
		FromStatus:  ptrFromText(row.FromStatus),
		ToStatus:    ptrFromText(row.ToStatus),
		Comment:     ptrFromText(row.Comment),
		ActorID:     ptrFromUUID(row.ActorID),
		Metadata:    metadata,
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
	}
}

func mapProjectMilestoneEvents(rows []sqlc.HrmsProjectMilestoneEvent) []*domain.ProjectMilestoneEvent {
	items := make([]*domain.ProjectMilestoneEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapProjectMilestoneEvent(row))
	}
	return items
}
