package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapEngagement(row sqlc.HrmsEngagement) *domain.Engagement {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.Engagement{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		WorkerProfileID:    row.WorkerProfileID,
		EngagementCode:     ptrFromText(row.EngagementCode),
		Title:              row.Title,
		Description:        ptrFromText(row.Description),
		EngagementType:     row.EngagementType,
		Status:             row.Status,
		StartDate:          timeFromDate(row.StartDate),
		EndDate:            ptrFromDate(row.EndDate),
		HoursBudget:        ptrFromNumeric(row.HoursBudget),
		RateAmount:         ptrFromNumeric(row.RateAmount),
		CurrencyCode:       row.CurrencyCode,
		RateUnit:           row.RateUnit,
		BranchID:           ptrFromUUID(row.BranchID),
		DepartmentID:       ptrFromUUID(row.DepartmentID),
		ReportingManagerID: ptrFromUUID(row.ReportingManagerID),
		ProjectLabel:       ptrFromText(row.ProjectLabel),
		ProjectCode:        ptrFromText(row.ProjectCode),
		CostCenter:         ptrFromText(row.CostCenter),
		RenewalDueDate:     ptrFromDate(row.RenewalDueDate),
		RenewalStatus:      row.RenewalStatus,
		TerminationReason:  ptrFromText(row.TerminationReason),
		TerminatedAt:       ptrFromTimestamptz(row.TerminatedAt),
		Notes:              ptrFromText(row.Notes),
		Metadata:           metadata,
		Inactive:           row.Inactive,
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		CreatedBy:          ptrFromUUID(row.CreatedBy),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:          ptrFromUUID(row.UpdatedBy),
	}
}

func mapEngagementListItems(rows []sqlc.ListEngagementsRow) []*domain.EngagementListItem {
	items := make([]*domain.EngagementListItem, 0, len(rows))
	for _, row := range rows {
		metadata := json.RawMessage(row.Metadata)
		if len(metadata) == 0 {
			metadata = json.RawMessage(`{}`)
		}
		items = append(items, &domain.EngagementListItem{
			Engagement: domain.Engagement{
				ID:                 row.ID,
				TenantID:           row.TenantID,
				WorkerProfileID:    row.WorkerProfileID,
				EngagementCode:     ptrFromText(row.EngagementCode),
				Title:              row.Title,
				Description:        ptrFromText(row.Description),
				EngagementType:     row.EngagementType,
				Status:             row.Status,
				StartDate:          timeFromDate(row.StartDate),
				EndDate:            ptrFromDate(row.EndDate),
				HoursBudget:        ptrFromNumeric(row.HoursBudget),
				RateAmount:         ptrFromNumeric(row.RateAmount),
				CurrencyCode:       row.CurrencyCode,
				RateUnit:           row.RateUnit,
				BranchID:           ptrFromUUID(row.BranchID),
				DepartmentID:       ptrFromUUID(row.DepartmentID),
				ReportingManagerID: ptrFromUUID(row.ReportingManagerID),
				ProjectLabel:       ptrFromText(row.ProjectLabel),
				ProjectCode:        ptrFromText(row.ProjectCode),
				CostCenter:         ptrFromText(row.CostCenter),
				RenewalDueDate:     ptrFromDate(row.RenewalDueDate),
				RenewalStatus:      row.RenewalStatus,
				TerminationReason:  ptrFromText(row.TerminationReason),
				TerminatedAt:       ptrFromTimestamptz(row.TerminatedAt),
				Notes:              ptrFromText(row.Notes),
				Metadata:           metadata,
				Inactive:           row.Inactive,
				CreatedAt:          timeFromTimestamptz(row.CreatedAt),
				CreatedBy:          ptrFromUUID(row.CreatedBy),
				UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
				UpdatedBy:          ptrFromUUID(row.UpdatedBy),
			},
			WorkerDisplayName:   row.WorkerDisplayName,
			WorkerCode:          ptrFromText(row.WorkerCode),
			EmployeeID:          ptrFromUUID(row.EmployeeID),
			WorkerTypeName:      row.WorkerTypeName,
			ClassificationGroup: row.ClassificationGroup,
			BranchName:          ptrFromText(row.BranchName),
			DepartmentName:      ptrFromText(row.DepartmentName),
		})
	}
	return items
}
