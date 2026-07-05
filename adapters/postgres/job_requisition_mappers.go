package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapJobRequisition(row sqlc.HrmsJobRequisition) *domain.JobRequisition {
	return &domain.JobRequisition{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		JobPositionID:       row.JobPositionID,
		Code:                ptrFromText(row.Code),
		Title:               row.Title,
		Level:               ptrFromText(row.Level),
		Category:            ptrFromText(row.Category),
		DepartmentID:        ptrFromUUID(row.DepartmentID),
		EmploymentTypeID:    ptrFromUUID(row.EmploymentTypeID),
		Description:         ptrFromText(row.Description),
		WorkMode:            ptrFromText(row.WorkMode),
		TotalOpenings:       row.TotalOpenings,
		ReasonForHire:       ptrFromText(row.ReasonForHire),
		MinSalary:           floatPtrFromNumeric(row.MinSalary),
		MaxSalary:           floatPtrFromNumeric(row.MaxSalary),
		Currency:            row.Currency,
		TargetHireDate:      ptrFromDate(row.TargetHireDate),
		ExpectedClosureDate: ptrFromDate(row.ExpectedClosureDate),
		RequestedBy:         row.RequestedBy,
		RequestedDate:       ptrFromDate(row.RequestedDate),
		IsApproved:          row.IsApproved,
		ApprovedBy:          ptrFromUUID(row.ApprovedBy),
		ApprovedDate:        ptrFromTimestamptz(row.ApprovedDate),
		Priority:            ptrFromText(row.Priority),
		Status:              row.Status,
		Notes:               ptrFromText(row.Notes),
		Inactive:            row.Inactive,
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
		CreatedBy:           ptrFromUUID(row.CreatedBy),
		UpdatedAt:           timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:           ptrFromUUID(row.UpdatedBy),
	}
}

func mapJobRequisitionListRow(row sqlc.ListJobRequisitionsRow) *domain.JobRequisition {
	item := &domain.JobRequisition{
		ID:                     row.ID,
		TenantID:               row.TenantID,
		JobPositionID:          row.JobPositionID,
		JobPositionCode:        ptrFromText(row.JobPositionCode),
		PositionTotalHeadcount: row.PositionTotalHeadcount,
		PositionBudgetedCost:   floatPtrFromNumeric(row.PositionBudgetedCost),
		Code:                   ptrFromText(row.Code),
		Title:                  row.Title,
		Level:                  ptrFromText(row.Level),
		Category:               ptrFromText(row.Category),
		DepartmentID:           ptrFromUUID(row.DepartmentID),
		DepartmentName:         ptrFromText(row.DepartmentName),
		EmploymentTypeID:       ptrFromUUID(row.EmploymentTypeID),
		EmploymentTypeName:     ptrFromText(row.EmploymentTypeName),
		Description:            ptrFromText(row.Description),
		WorkMode:               ptrFromText(row.WorkMode),
		TotalOpenings:          row.TotalOpenings,
		ReasonForHire:          ptrFromText(row.ReasonForHire),
		MinSalary:              floatPtrFromNumeric(row.MinSalary),
		MaxSalary:              floatPtrFromNumeric(row.MaxSalary),
		Currency:               row.Currency,
		TargetHireDate:         ptrFromDate(row.TargetHireDate),
		ExpectedClosureDate:    ptrFromDate(row.ExpectedClosureDate),
		RequestedBy:            row.RequestedBy,
		RequestedDate:          ptrFromDate(row.RequestedDate),
		IsApproved:             row.IsApproved,
		ApprovedBy:             ptrFromUUID(row.ApprovedBy),
		ApprovedDate:           ptrFromTimestamptz(row.ApprovedDate),
		Priority:               ptrFromText(row.Priority),
		Status:                 row.Status,
		Notes:                  ptrFromText(row.Notes),
		LogCount:               row.LogCount,
		Inactive:               row.Inactive,
		CreatedAt:              timeFromTimestamptz(row.CreatedAt),
		CreatedBy:              ptrFromUUID(row.CreatedBy),
		UpdatedAt:              timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:              ptrFromUUID(row.UpdatedBy),
	}
	return item
}

func mapJobRequisitions(rows []sqlc.ListJobRequisitionsRow) []*domain.JobRequisition {
	items := make([]*domain.JobRequisition, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapJobRequisitionListRow(row))
	}
	return items
}

func mapJobRequisitionLog(row sqlc.HrmsJobRequisitionLog) *domain.JobRequisitionLog {
	return &domain.JobRequisitionLog{
		ID:               row.ID,
		TenantID:         row.TenantID,
		JobRequisitionID: row.JobRequisitionID,
		FromStatus:       ptrFromText(row.FromStatus),
		ToStatus:         row.ToStatus,
		Action:           row.Action,
		Remarks:          ptrFromText(row.Remarks),
		Inactive:         row.Inactive,
		CreatedAt:        timeFromTimestamptz(row.CreatedAt),
		CreatedBy:        ptrFromUUID(row.CreatedBy),
		UpdatedAt:        timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:        ptrFromUUID(row.UpdatedBy),
	}
}

func mapJobRequisitionLogs(rows []sqlc.HrmsJobRequisitionLog) []*domain.JobRequisitionLog {
	items := make([]*domain.JobRequisitionLog, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapJobRequisitionLog(row))
	}
	return items
}
