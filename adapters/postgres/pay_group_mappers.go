package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapPayGroup(row sqlc.HrmsPayGroup) *domain.PayGroup {
	return &domain.PayGroup{
		ID:               row.ID,
		TenantID:         row.TenantID,
		Code:             row.Code,
		Name:             row.Name,
		Description:      ptrFromText(row.Description),
		GroupingType:     row.GroupingType,
		BranchID:         ptrFromUUID(row.BranchID),
		DepartmentID:     ptrFromUUID(row.DepartmentID),
		EmploymentTypeID: ptrFromUUID(row.EmploymentTypeID),
		ReportingTag:     ptrFromText(row.ReportingTag),
		Rules:            json.RawMessage(row.Rules),
		IsActive:         row.IsActive,
		Inactive:         row.Inactive,
		CreatedAt:        timeFromTimestamptz(row.CreatedAt),
		CreatedBy:        ptrFromUUID(row.CreatedBy),
		UpdatedAt:        timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:        ptrFromUUID(row.UpdatedBy),
	}
}

func mapPayGroups(rows []sqlc.HrmsPayGroup) []*domain.PayGroup {
	items := make([]*domain.PayGroup, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayGroup(row))
	}
	return items
}

func mapPayGroupMember(row sqlc.HrmsPayGroupMember) *domain.PayGroupMember {
	return &domain.PayGroupMember{
		ID:             row.ID,
		TenantID:       row.TenantID,
		PayGroupID:     row.PayGroupID,
		UserID:         row.UserID,
		MembershipType: row.MembershipType,
		EffectiveFrom:  ptrFromDate(row.EffectiveFrom),
		EffectiveTo:    ptrFromDate(row.EffectiveTo),
		Inactive:       row.Inactive,
		CreatedAt:      timeFromTimestamptz(row.CreatedAt),
		CreatedBy:      ptrFromUUID(row.CreatedBy),
		UpdatedAt:      timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:      ptrFromUUID(row.UpdatedBy),
	}
}

func mapPayGroupMembers(rows []sqlc.HrmsPayGroupMember) []*domain.PayGroupMember {
	items := make([]*domain.PayGroupMember, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayGroupMember(row))
	}
	return items
}

func mapPayGroupEmployee(row sqlc.ListPayGroupEmployeesRow) *domain.PayGroupEmployee {
	return &domain.PayGroupEmployee{
		EmployeeID:         row.EmployeeID,
		UserID:             row.UserID,
		EmployeeCode:       ptrFromText(row.EmployeeCode),
		Firstname:          row.Firstname,
		Lastname:           ptrFromText(row.Lastname),
		BranchID:           ptrFromUUID(row.BranchID),
		BranchName:         ptrFromText(row.BranchName),
		DepartmentID:       ptrFromUUID(row.DepartmentID),
		DepartmentName:     ptrFromText(row.DepartmentName),
		EmploymentTypeID:   ptrFromUUID(row.EmploymentTypeID),
		EmploymentTypeName: ptrFromText(row.EmploymentTypeName),
		MatchSource:        row.MatchSource,
	}
}

func mapPayGroupEmployees(rows []sqlc.ListPayGroupEmployeesRow) []*domain.PayGroupEmployee {
	items := make([]*domain.PayGroupEmployee, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayGroupEmployee(row))
	}
	return items
}

func mapPayRun(row sqlc.HrmsPayRun) *domain.PayRun {
	return &domain.PayRun{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		PayGroupID:          row.PayGroupID,
		FYID:                row.FyID,
		Month:               row.Month,
		Year:                row.Year,
		Status:              row.Status,
		EmployeeCount:       row.EmployeeCount,
		ReadyCount:          row.ReadyCount,
		BlockedCount:        row.BlockedCount,
		GeneratedCount:      row.GeneratedCount,
		AttendanceFrozenAt:  ptrFromTimestamptz(row.AttendanceFrozenAt),
		LOPFrozenAt:         ptrFromTimestamptz(row.LopFrozenAt),
		AdjustmentsFrozenAt: ptrFromTimestamptz(row.AdjustmentsFrozenAt),
		GeneratedAt:         ptrFromTimestamptz(row.GeneratedAt),
		LockedAt:            ptrFromTimestamptz(row.LockedAt),
		LockedBy:            ptrFromUUID(row.LockedBy),
		UnlockedAt:          ptrFromTimestamptz(row.UnlockedAt),
		UnlockedBy:          ptrFromUUID(row.UnlockedBy),
		Readiness:           json.RawMessage(row.Readiness),
		Notes:               ptrFromText(row.Notes),
		Inactive:            row.Inactive,
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
		CreatedBy:           ptrFromUUID(row.CreatedBy),
		UpdatedAt:           timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:           ptrFromUUID(row.UpdatedBy),
	}
}

func mapPayRuns(rows []sqlc.HrmsPayRun) []*domain.PayRun {
	items := make([]*domain.PayRun, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayRun(row))
	}
	return items
}

func mapPayRunEmployeeRecord(row sqlc.HrmsPayRunEmployee) *domain.PayRunEmployee {
	return &domain.PayRunEmployee{
		ID:              row.ID,
		TenantID:        row.TenantID,
		PayRunID:        row.PayRunID,
		UserID:          row.UserID,
		ReadinessStatus: row.ReadinessStatus,
		BlockerReason:   ptrFromText(row.BlockerReason),
		SalarySlipID:    ptrFromUUID(row.SalarySlipID),
		GeneratedAt:     ptrFromTimestamptz(row.GeneratedAt),
		Inactive:        row.Inactive,
		CreatedAt:       timeFromTimestamptz(row.CreatedAt),
		CreatedBy:       ptrFromUUID(row.CreatedBy),
		UpdatedAt:       timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:       ptrFromUUID(row.UpdatedBy),
	}
}

func mapPayRunEmployee(row sqlc.ListPayRunEmployeesRow) *domain.PayRunEmployee {
	item := mapPayRunEmployeeRecord(sqlc.HrmsPayRunEmployee{
		ID:              row.ID,
		TenantID:        row.TenantID,
		PayRunID:        row.PayRunID,
		UserID:          row.UserID,
		ReadinessStatus: row.ReadinessStatus,
		BlockerReason:   row.BlockerReason,
		SalarySlipID:    row.SalarySlipID,
		GeneratedAt:     row.GeneratedAt,
		Inactive:        row.Inactive,
		CreatedAt:       row.CreatedAt,
		CreatedBy:       row.CreatedBy,
		UpdatedAt:       row.UpdatedAt,
		UpdatedBy:       row.UpdatedBy,
	})
	item.EmployeeCode = ptrFromText(row.EmployeeCode)
	item.Firstname = row.Firstname
	item.Lastname = ptrFromText(row.Lastname)
	item.BranchName = ptrFromText(row.BranchName)
	item.DepartmentName = ptrFromText(row.DepartmentName)
	return item
}

func mapPayRunEmployees(rows []sqlc.ListPayRunEmployeesRow) []*domain.PayRunEmployee {
	items := make([]*domain.PayRunEmployee, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayRunEmployee(row))
	}
	return items
}

func mapPayRunEvent(row sqlc.HrmsPayRunEvent) *domain.PayRunEvent {
	return &domain.PayRunEvent{
		ID:         row.ID,
		TenantID:   row.TenantID,
		PayRunID:   row.PayRunID,
		Action:     row.Action,
		FromStatus: ptrFromText(row.FromStatus),
		ToStatus:   ptrFromText(row.ToStatus),
		Remarks:    ptrFromText(row.Remarks),
		Metadata:   json.RawMessage(row.Metadata),
		Inactive:   row.Inactive,
		CreatedAt:  timeFromTimestamptz(row.CreatedAt),
		CreatedBy:  ptrFromUUID(row.CreatedBy),
		UpdatedAt:  timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:  ptrFromUUID(row.UpdatedBy),
	}
}

func mapPayRunEvents(rows []sqlc.HrmsPayRunEvent) []*domain.PayRunEvent {
	items := make([]*domain.PayRunEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayRunEvent(row))
	}
	return items
}
