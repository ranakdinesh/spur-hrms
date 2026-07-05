package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapJobPosition(row sqlc.HrmsJobPosition) *domain.JobPosition {
	return &domain.JobPosition{
		ID:               row.ID,
		TenantID:         row.TenantID,
		Code:             ptrFromText(row.Code),
		Title:            row.Title,
		Level:            ptrFromText(row.Level),
		Category:         ptrFromText(row.Category),
		Description:      ptrFromText(row.Description),
		DepartmentID:     ptrFromUUID(row.DepartmentID),
		EmploymentTypeID: ptrFromUUID(row.EmploymentTypeID),
		WorkMode:         ptrFromText(row.WorkMode),
		TotalPosition:    row.TotalPosition,
		BudgetedCost:     floatPtrFromNumeric(row.BudgetedCost),
		Inactive:         row.Inactive,
		CreatedAt:        timeFromTimestamptz(row.CreatedAt),
		CreatedBy:        ptrFromUUID(row.CreatedBy),
		UpdatedAt:        timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:        ptrFromUUID(row.UpdatedBy),
	}
}

func mapJobPositionListRow(row sqlc.ListJobPositionsRow) *domain.JobPosition {
	return &domain.JobPosition{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		Code:                 ptrFromText(row.Code),
		Title:                row.Title,
		Level:                ptrFromText(row.Level),
		Category:             ptrFromText(row.Category),
		Description:          ptrFromText(row.Description),
		DepartmentID:         ptrFromUUID(row.DepartmentID),
		DepartmentName:       ptrFromText(row.DepartmentName),
		EmploymentTypeID:     ptrFromUUID(row.EmploymentTypeID),
		EmploymentTypeName:   ptrFromText(row.EmploymentTypeName),
		WorkMode:             ptrFromText(row.WorkMode),
		TotalPosition:        row.TotalPosition,
		BudgetedCost:         floatPtrFromNumeric(row.BudgetedCost),
		LocationCount:        row.LocationCount,
		OpenRequisitionCount: row.OpenRequisitionCount,
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapJobPositions(rows []sqlc.ListJobPositionsRow) []*domain.JobPosition {
	items := make([]*domain.JobPosition, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapJobPositionListRow(row))
	}
	return items
}

func mapJobPositionLocation(row sqlc.HrmsJobPositionLocation) *domain.JobPositionLocation {
	return &domain.JobPositionLocation{
		ID:            row.ID,
		TenantID:      row.TenantID,
		JobPositionID: row.JobPositionID,
		Location:      ptrFromText(row.Location),
		City:          ptrFromText(row.City),
		State:         ptrFromText(row.State),
		Country:       ptrFromText(row.Country),
		IsRemote:      row.IsRemote,
		Inactive:      row.Inactive,
		CreatedAt:     timeFromTimestamptz(row.CreatedAt),
		CreatedBy:     ptrFromUUID(row.CreatedBy),
		UpdatedAt:     timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:     ptrFromUUID(row.UpdatedBy),
	}
}

func mapJobPositionLocations(rows []sqlc.HrmsJobPositionLocation) []*domain.JobPositionLocation {
	items := make([]*domain.JobPositionLocation, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapJobPositionLocation(row))
	}
	return items
}
