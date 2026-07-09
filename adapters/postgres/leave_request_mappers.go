package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapLeave(row sqlc.HrmsLeafe) *domain.Leave {
	return &domain.Leave{
		ID:            row.ID,
		TenantID:      row.TenantID,
		UserID:        row.UserID,
		LeaveTypeID:   row.LeaveTypeID,
		FYID:          row.FyID,
		StartDate:     timeFromDate(row.StartDate),
		EndDate:       timeFromDate(row.EndDate),
		StartDayType:  row.StartDayType,
		EndDayType:    row.EndDayType,
		Days:          floatFromNumeric(row.Days),
		Reason:        ptrFromText(row.Reason),
		Status:        row.Status,
		AppliedDate:   timeFromTimestamptz(row.AppliedDate),
		FromLeaveType: ptrFromUUID(row.FromLeaveType),
		ToLeaveType:   ptrFromUUID(row.ToLeaveType),
		IsSandwich:    row.IsSandwich,
		Inactive:      row.Inactive,
		CreatedAt:     timeFromTimestamptz(row.CreatedAt),
		CreatedBy:     ptrFromUUID(row.CreatedBy),
		UpdatedAt:     timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:     ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeaves(rows []sqlc.HrmsLeafe) []*domain.Leave {
	items := make([]*domain.Leave, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeave(row))
	}
	return items
}

func mapLeaveApproval(row sqlc.HrmsLeaveApproval) *domain.LeaveApproval {
	return &domain.LeaveApproval{
		ID:                row.ID,
		TenantID:          row.TenantID,
		LeaveID:           row.LeaveID,
		ApproverID:        row.ApproverID,
		Status:            row.Status,
		Remarks:           ptrFromText(row.Remarks),
		ActionDate:        ptrFromTimestamptz(row.ActionDate),
		WorkflowID:        ptrFromUUID(row.WorkflowID),
		WorkflowStepID:    ptrFromUUID(row.WorkflowStepID),
		StepOrder:         row.StepOrder,
		DecisionRule:      row.DecisionRule,
		RequiredApprovals: row.RequiredApprovals,
		Inactive:          row.Inactive,
		CreatedAt:         timeFromTimestamptz(row.CreatedAt),
		CreatedBy:         ptrFromUUID(row.CreatedBy),
		UpdatedAt:         timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:         ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeaveApprovals(rows []sqlc.HrmsLeaveApproval) []*domain.LeaveApproval {
	items := make([]*domain.LeaveApproval, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveApproval(row))
	}
	return items
}

func mapLeaveRequestMessage(row sqlc.HrmsLeaveRequestMessage) *domain.LeaveRequestMessage {
	return &domain.LeaveRequestMessage{
		ID:              row.ID,
		TenantID:        row.TenantID,
		LeaveID:         row.LeaveID,
		SenderUserID:    row.SenderUserID,
		RecipientUserID: ptrFromUUID(row.RecipientUserID),
		MessageType:     row.MessageType,
		Body:            row.Body,
		Inactive:        row.Inactive,
		CreatedAt:       timeFromTimestamptz(row.CreatedAt),
		CreatedBy:       ptrFromUUID(row.CreatedBy),
		UpdatedAt:       timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:       ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeaveRequestMessages(rows []sqlc.HrmsLeaveRequestMessage) []*domain.LeaveRequestMessage {
	items := make([]*domain.LeaveRequestMessage, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveRequestMessage(row))
	}
	return items
}

func mapLeaveReportRow(row sqlc.ListLeaveReportRowsRow) *domain.LeaveReportRow {
	return &domain.LeaveReportRow{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		UserID:             row.UserID,
		EmployeeCode:       ptrFromText(row.EmployeeCode),
		Firstname:          row.Firstname,
		Lastname:           ptrFromText(row.Lastname),
		ReportingManagerID: ptrFromUUID(row.ReportingManagerID),
		DepartmentID:       ptrFromUUID(row.DepartmentID),
		DepartmentName:     ptrFromText(row.DepartmentName),
		DesignationID:      ptrFromUUID(row.DesignationID),
		DesignationName:    ptrFromText(row.DesignationName),
		LeaveTypeID:        row.LeaveTypeID,
		LeaveTypeName:      ptrFromText(row.LeaveTypeName),
		LeaveTypeShortcode: ptrFromText(row.LeaveTypeShortcode),
		FYID:               row.FyID,
		FinancialYearName:  ptrFromText(row.FinancialYearName),
		StartDate:          timeFromDate(row.StartDate),
		EndDate:            timeFromDate(row.EndDate),
		StartDayType:       row.StartDayType,
		EndDayType:         row.EndDayType,
		Days:               floatFromNumeric(row.Days),
		Reason:             ptrFromText(row.Reason),
		Status:             row.Status,
		IsSandwich:         row.IsSandwich,
		AppliedDate:        timeFromTimestamptz(row.AppliedDate),
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
	}
}

func mapLeaveReportRows(rows []sqlc.ListLeaveReportRowsRow) []*domain.LeaveReportRow {
	items := make([]*domain.LeaveReportRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveReportRow(row))
	}
	return items
}

func mapManagerLeaveReportRow(row sqlc.ListManagerLeaveReportRowsRow) *domain.LeaveReportRow {
	return &domain.LeaveReportRow{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		UserID:             row.UserID,
		EmployeeCode:       ptrFromText(row.EmployeeCode),
		Firstname:          row.Firstname,
		Lastname:           ptrFromText(row.Lastname),
		ReportingManagerID: ptrFromUUID(row.ReportingManagerID),
		DepartmentID:       ptrFromUUID(row.DepartmentID),
		DepartmentName:     ptrFromText(row.DepartmentName),
		DesignationID:      ptrFromUUID(row.DesignationID),
		DesignationName:    ptrFromText(row.DesignationName),
		LeaveTypeID:        row.LeaveTypeID,
		LeaveTypeName:      ptrFromText(row.LeaveTypeName),
		LeaveTypeShortcode: ptrFromText(row.LeaveTypeShortcode),
		FYID:               row.FyID,
		FinancialYearName:  ptrFromText(row.FinancialYearName),
		StartDate:          timeFromDate(row.StartDate),
		EndDate:            timeFromDate(row.EndDate),
		StartDayType:       row.StartDayType,
		EndDayType:         row.EndDayType,
		Days:               floatFromNumeric(row.Days),
		Reason:             ptrFromText(row.Reason),
		Status:             row.Status,
		IsSandwich:         row.IsSandwich,
		AppliedDate:        timeFromTimestamptz(row.AppliedDate),
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
	}
}

func mapManagerLeaveReportRows(rows []sqlc.ListManagerLeaveReportRowsRow) []*domain.LeaveReportRow {
	items := make([]*domain.LeaveReportRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapManagerLeaveReportRow(row))
	}
	return items
}

func mapLeaveReportSummary(row sqlc.GetLeaveReportSummaryRow) *domain.LeaveReportSummary {
	return &domain.LeaveReportSummary{
		TotalRequests: row.TotalRequests,
		TotalDays:     floatFromNumeric(row.TotalDays),
		EmployeeCount: row.EmployeeCount,
		PendingCount:  row.PendingCount,
		ApprovedCount: row.ApprovedCount,
		RejectedCount: row.RejectedCount,
		CanceledCount: row.CanceledCount,
		PendingDays:   floatFromNumeric(row.PendingDays),
		ApprovedDays:  floatFromNumeric(row.ApprovedDays),
		RejectedDays:  floatFromNumeric(row.RejectedDays),
		CanceledDays:  floatFromNumeric(row.CanceledDays),
	}
}
