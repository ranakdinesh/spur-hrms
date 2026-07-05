package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapShiftTemplate(row sqlc.HrmsShiftTemplate) *domain.ShiftTemplate {
	return &domain.ShiftTemplate{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		Code:                 row.Code,
		Name:                 row.Name,
		Description:          ptrFromText(row.Description),
		StartTime:            clockStringFromTime(row.StartTime),
		EndTime:              clockStringFromTime(row.EndTime),
		BreakMinutes:         row.BreakMinutes,
		PaidMinutes:          row.PaidMinutes,
		WorkMode:             row.WorkMode,
		LocationType:         row.LocationType,
		AttendancePolicyID:   ptrFromUUID(row.AttendancePolicyID),
		AttendanceLocationID: ptrFromUUID(row.AttendanceLocationID),
		AllowOvertime:        row.AllowOvertime,
		PayrollCode:          ptrFromText(row.PayrollCode),
		Metadata:             shiftJSONRaw(row.Metadata),
		IsActive:             row.IsActive,
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapShiftTemplates(rows []sqlc.HrmsShiftTemplate) []*domain.ShiftTemplate {
	items := make([]*domain.ShiftTemplate, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapShiftTemplate(row))
	}
	return items
}

func mapStaffingRequirement(row sqlc.HrmsStaffingRequirement) *domain.StaffingRequirement {
	return &domain.StaffingRequirement{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		Name:                 row.Name,
		RequirementDate:      ptrFromDate(row.RequirementDate),
		StartDate:            ptrFromDate(row.StartDate),
		EndDate:              ptrFromDate(row.EndDate),
		DayOfWeek:            ptrFromInt4(row.DayOfWeek),
		BranchID:             ptrFromUUID(row.BranchID),
		DepartmentID:         ptrFromUUID(row.DepartmentID),
		AttendanceLocationID: ptrFromUUID(row.AttendanceLocationID),
		RoleLabel:            ptrFromText(row.RoleLabel),
		TeamLabel:            ptrFromText(row.TeamLabel),
		ShiftTemplateID:      ptrFromUUID(row.ShiftTemplateID),
		RequiredCount:        row.RequiredCount,
		MinCount:             row.MinCount,
		MaxCount:             ptrFromInt4(row.MaxCount),
		Priority:             row.Priority,
		Status:               row.Status,
		PayrollBlocking:      row.PayrollBlocking,
		Notes:                ptrFromText(row.Notes),
		Metadata:             shiftJSONRaw(row.Metadata),
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapStaffingRequirements(rows []sqlc.HrmsStaffingRequirement) []*domain.StaffingRequirement {
	items := make([]*domain.StaffingRequirement, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapStaffingRequirement(row))
	}
	return items
}

func mapShiftScheduleAssignment(row sqlc.HrmsShiftScheduleAssignment) *domain.ShiftScheduleAssignment {
	return &domain.ShiftScheduleAssignment{
		ID:                     row.ID,
		TenantID:               row.TenantID,
		ScheduleDate:           timeFromDate(row.ScheduleDate),
		WorkerProfileID:        ptrFromUUID(row.WorkerProfileID),
		EmployeeUserID:         ptrFromUUID(row.EmployeeUserID),
		ShiftTemplateID:        ptrFromUUID(row.ShiftTemplateID),
		AttendancePolicyID:     ptrFromUUID(row.AttendancePolicyID),
		AttendanceLocationID:   ptrFromUUID(row.AttendanceLocationID),
		AttendanceRosterID:     ptrFromUUID(row.AttendanceRosterID),
		BranchID:               ptrFromUUID(row.BranchID),
		DepartmentID:           ptrFromUUID(row.DepartmentID),
		StartTime:              clockStringFromTime(row.StartTime),
		EndTime:                clockStringFromTime(row.EndTime),
		BreakMinutes:           row.BreakMinutes,
		WorkMode:               row.WorkMode,
		LocationType:           row.LocationType,
		RoleLabel:              ptrFromText(row.RoleLabel),
		TeamLabel:              ptrFromText(row.TeamLabel),
		Status:                 row.Status,
		Source:                 row.Source,
		OvertimePlannedMinutes: row.OvertimePlannedMinutes,
		HasConflict:            row.HasConflict,
		ConflictReason:         ptrFromText(row.ConflictReason),
		PayrollBlocking:        row.PayrollBlocking,
		Notes:                  ptrFromText(row.Notes),
		Metadata:               shiftJSONRaw(row.Metadata),
		Inactive:               row.Inactive,
		CreatedAt:              timeFromTimestamptz(row.CreatedAt),
		CreatedBy:              ptrFromUUID(row.CreatedBy),
		UpdatedAt:              timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:              ptrFromUUID(row.UpdatedBy),
	}
}

func mapShiftScheduleAssignments(rows []sqlc.HrmsShiftScheduleAssignment) []*domain.ShiftScheduleAssignment {
	items := make([]*domain.ShiftScheduleAssignment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapShiftScheduleAssignment(row))
	}
	return items
}

func mapShiftSwapRequest(row sqlc.HrmsShiftSwapRequest) *domain.ShiftSwapRequest {
	return &domain.ShiftSwapRequest{
		ID:                       row.ID,
		TenantID:                 row.TenantID,
		RequesterAssignmentID:    row.RequesterAssignmentID,
		RequesterWorkerProfileID: ptrFromUUID(row.RequesterWorkerProfileID),
		RequesterUserID:          ptrFromUUID(row.RequesterUserID),
		TargetWorkerProfileID:    ptrFromUUID(row.TargetWorkerProfileID),
		TargetUserID:             ptrFromUUID(row.TargetUserID),
		OfferedAssignmentID:      ptrFromUUID(row.OfferedAssignmentID),
		RequestedDate:            ptrFromDate(row.RequestedDate),
		RequestedShiftTemplateID: ptrFromUUID(row.RequestedShiftTemplateID),
		Reason:                   ptrFromText(row.Reason),
		Status:                   row.Status,
		ReviewedBy:               ptrFromUUID(row.ReviewedBy),
		ReviewedAt:               ptrFromTimestamptz(row.ReviewedAt),
		ReviewRemarks:            ptrFromText(row.ReviewRemarks),
		PayrollBlocking:          row.PayrollBlocking,
		Metadata:                 shiftJSONRaw(row.Metadata),
		Inactive:                 row.Inactive,
		CreatedAt:                timeFromTimestamptz(row.CreatedAt),
		CreatedBy:                ptrFromUUID(row.CreatedBy),
		UpdatedAt:                timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:                ptrFromUUID(row.UpdatedBy),
	}
}

func mapShiftSwapRequests(rows []sqlc.HrmsShiftSwapRequest) []*domain.ShiftSwapRequest {
	items := make([]*domain.ShiftSwapRequest, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapShiftSwapRequest(row))
	}
	return items
}

func mapShiftScheduleEvent(row sqlc.HrmsShiftScheduleEvent) *domain.ShiftScheduleEvent {
	return &domain.ShiftScheduleEvent{
		ID:          row.ID,
		TenantID:    row.TenantID,
		SourceType:  row.SourceType,
		SourceID:    row.SourceID,
		Action:      row.Action,
		FromStatus:  ptrFromText(row.FromStatus),
		ToStatus:    ptrFromText(row.ToStatus),
		ActorUserID: ptrFromUUID(row.ActorUserID),
		Remarks:     ptrFromText(row.Remarks),
		Metadata:    shiftJSONRaw(row.Metadata),
		Inactive:    row.Inactive,
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
		CreatedBy:   ptrFromUUID(row.CreatedBy),
		UpdatedAt:   timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:   ptrFromUUID(row.UpdatedBy),
	}
}

func mapShiftScheduleEvents(rows []sqlc.HrmsShiftScheduleEvent) []*domain.ShiftScheduleEvent {
	items := make([]*domain.ShiftScheduleEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapShiftScheduleEvent(row))
	}
	return items
}

func mapShiftScheduleSummaryRows(rows []sqlc.GetShiftScheduleSummaryRow) []*domain.ShiftScheduleSummaryRow {
	items := make([]*domain.ShiftScheduleSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.ShiftScheduleSummaryRow{Metric: row.Metric, MetricCount: row.MetricCount})
	}
	return items
}

func mapShiftStaffingGapRows(rows []sqlc.ListShiftStaffingGapsRow) []*domain.ShiftStaffingGapRow {
	items := make([]*domain.ShiftStaffingGapRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.ShiftStaffingGapRow{
			RequirementID:        row.RequirementID,
			RequirementName:      row.RequirementName,
			BranchID:             ptrFromUUID(row.BranchID),
			DepartmentID:         ptrFromUUID(row.DepartmentID),
			AttendanceLocationID: ptrFromUUID(row.AttendanceLocationID),
			ShiftTemplateID:      ptrFromUUID(row.ShiftTemplateID),
			RequiredCount:        row.RequiredCount,
			AssignedCount:        row.AssignedCount,
			GapCount:             row.GapCount,
			Priority:             row.Priority,
			PayrollBlocking:      row.PayrollBlocking,
		})
	}
	return items
}

func shiftJSONRaw(value []byte) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(value)
}
