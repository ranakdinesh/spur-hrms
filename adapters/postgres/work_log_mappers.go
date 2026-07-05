package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapWorkLog(row sqlc.HrmsWorkLog) *domain.WorkLog {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.WorkLog{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		EngagementID:         row.EngagementID,
		WorkerProfileID:      row.WorkerProfileID,
		LogDate:              timeFromDate(row.LogDate),
		HoursWorked:          floatFromNumeric(row.HoursWorked),
		BillableHours:        ptrFromNumeric(row.BillableHours),
		WorkSummary:          ptrFromText(row.WorkSummary),
		DeliverableReference: ptrFromText(row.DeliverableReference),
		Status:               row.Status,
		SubmittedAt:          ptrFromTimestamptz(row.SubmittedAt),
		SubmittedBy:          ptrFromUUID(row.SubmittedBy),
		ReviewedAt:           ptrFromTimestamptz(row.ReviewedAt),
		ReviewedBy:           ptrFromUUID(row.ReviewedBy),
		ReviewComment:        ptrFromText(row.ReviewComment),
		Metadata:             metadata,
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapWorkLogListItems(rows []sqlc.ListWorkLogsRow) []*domain.WorkLogListItem {
	items := make([]*domain.WorkLogListItem, 0, len(rows))
	for _, row := range rows {
		metadata := json.RawMessage(row.Metadata)
		if len(metadata) == 0 {
			metadata = json.RawMessage(`{}`)
		}
		items = append(items, &domain.WorkLogListItem{
			WorkLog: domain.WorkLog{
				ID:                   row.ID,
				TenantID:             row.TenantID,
				EngagementID:         row.EngagementID,
				WorkerProfileID:      row.WorkerProfileID,
				LogDate:              timeFromDate(row.LogDate),
				HoursWorked:          floatFromNumeric(row.HoursWorked),
				BillableHours:        ptrFromNumeric(row.BillableHours),
				WorkSummary:          ptrFromText(row.WorkSummary),
				DeliverableReference: ptrFromText(row.DeliverableReference),
				Status:               row.Status,
				SubmittedAt:          ptrFromTimestamptz(row.SubmittedAt),
				SubmittedBy:          ptrFromUUID(row.SubmittedBy),
				ReviewedAt:           ptrFromTimestamptz(row.ReviewedAt),
				ReviewedBy:           ptrFromUUID(row.ReviewedBy),
				ReviewComment:        ptrFromText(row.ReviewComment),
				Metadata:             metadata,
				Inactive:             row.Inactive,
				CreatedAt:            timeFromTimestamptz(row.CreatedAt),
				CreatedBy:            ptrFromUUID(row.CreatedBy),
				UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
				UpdatedBy:            ptrFromUUID(row.UpdatedBy),
			},
			EngagementTitle:    row.EngagementTitle,
			EngagementCode:     ptrFromText(row.EngagementCode),
			ProjectLabel:       ptrFromText(row.ProjectLabel),
			ProjectCode:        ptrFromText(row.ProjectCode),
			CostCenter:         ptrFromText(row.CostCenter),
			WorkerDisplayName:  row.WorkerDisplayName,
			WorkerCode:         ptrFromText(row.WorkerCode),
			EmployeeID:         ptrFromUUID(row.EmployeeID),
			ReportingManagerID: ptrFromUUID(row.ReportingManagerID),
			DepartmentID:       ptrFromUUID(row.DepartmentID),
			DepartmentName:     ptrFromText(row.DepartmentName),
			BranchID:           ptrFromUUID(row.BranchID),
			BranchName:         ptrFromText(row.BranchName),
		})
	}
	return items
}

func mapWorkLogRollups(rows []sqlc.ListWorkLogRollupsRow) []*domain.WorkLogRollup {
	items := make([]*domain.WorkLogRollup, 0, len(rows))
	for _, row := range rows {
		budget := ptrFromNumeric(row.HoursBudget)
		approved := floatFromNumeric(row.ApprovedHours)
		items = append(items, &domain.WorkLogRollup{
			TenantID:          row.TenantID,
			EngagementID:      row.EngagementID,
			EngagementTitle:   row.EngagementTitle,
			EngagementCode:    ptrFromText(row.EngagementCode),
			WorkerProfileID:   row.WorkerProfileID,
			WorkerDisplayName: row.WorkerDisplayName,
			LogCount:          row.LogCount,
			TotalHours:        floatFromNumeric(row.TotalHours),
			BillableHours:     floatFromNumeric(row.BillableHours),
			ApprovedHours:     approved,
			SubmittedHours:    floatFromNumeric(row.SubmittedHours),
			RejectedHours:     floatFromNumeric(row.RejectedHours),
			HoursBudget:       budget,
			RemainingHours:    domain.WorkLogBudgetRemaining(budget, approved),
		})
	}
	return items
}
