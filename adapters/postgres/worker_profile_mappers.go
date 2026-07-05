package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapWorkerProfile(row sqlc.HrmsWorkerProfile) *domain.WorkerProfile {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.WorkerProfile{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		WorkerTypeID:       row.WorkerTypeID,
		EmployeeID:         ptrFromUUID(row.EmployeeID),
		EmployeeUserID:     ptrFromUUID(row.EmployeeUserID),
		WorkerCode:         ptrFromText(row.WorkerCode),
		DisplayName:        row.DisplayName,
		LegalName:          ptrFromText(row.LegalName),
		Email:              ptrFromText(row.Email),
		Mobile:             ptrFromText(row.Mobile),
		ProfileStatus:      row.ProfileStatus,
		StartDate:          ptrFromDate(row.StartDate),
		EndDate:            ptrFromDate(row.EndDate),
		BranchID:           ptrFromUUID(row.BranchID),
		DepartmentID:       ptrFromUUID(row.DepartmentID),
		ReportingManagerID: ptrFromUUID(row.ReportingManagerID),
		WorkLocationLabel:  ptrFromText(row.WorkLocationLabel),
		SourcePartner:      ptrFromText(row.SourcePartner),
		ExternalReference:  ptrFromText(row.ExternalReference),
		ComplianceStatus:   row.ComplianceStatus,
		PayrollStatus:      row.PayrollStatus,
		Notes:              ptrFromText(row.Notes),
		Metadata:           metadata,
		Inactive:           row.Inactive,
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		CreatedBy:          ptrFromUUID(row.CreatedBy),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:          ptrFromUUID(row.UpdatedBy),
	}
}

func mapWorkerProfileListItems(rows []sqlc.ListWorkerProfilesRow) []*domain.WorkerProfileListItem {
	items := make([]*domain.WorkerProfileListItem, 0, len(rows))
	for _, row := range rows {
		metadata := json.RawMessage(row.Metadata)
		if len(metadata) == 0 {
			metadata = json.RawMessage(`{}`)
		}
		items = append(items, &domain.WorkerProfileListItem{
			WorkerProfile: domain.WorkerProfile{
				ID:                 row.ID,
				TenantID:           row.TenantID,
				WorkerTypeID:       row.WorkerTypeID,
				EmployeeID:         ptrFromUUID(row.EmployeeID),
				EmployeeUserID:     ptrFromUUID(row.EmployeeUserID),
				WorkerCode:         ptrFromText(row.WorkerCode),
				DisplayName:        row.DisplayName,
				LegalName:          ptrFromText(row.LegalName),
				Email:              ptrFromText(row.Email),
				Mobile:             ptrFromText(row.Mobile),
				ProfileStatus:      row.ProfileStatus,
				StartDate:          ptrFromDate(row.StartDate),
				EndDate:            ptrFromDate(row.EndDate),
				BranchID:           ptrFromUUID(row.BranchID),
				DepartmentID:       ptrFromUUID(row.DepartmentID),
				ReportingManagerID: ptrFromUUID(row.ReportingManagerID),
				WorkLocationLabel:  ptrFromText(row.WorkLocationLabel),
				SourcePartner:      ptrFromText(row.SourcePartner),
				ExternalReference:  ptrFromText(row.ExternalReference),
				ComplianceStatus:   row.ComplianceStatus,
				PayrollStatus:      row.PayrollStatus,
				Notes:              ptrFromText(row.Notes),
				Metadata:           metadata,
				Inactive:           row.Inactive,
				CreatedAt:          timeFromTimestamptz(row.CreatedAt),
				CreatedBy:          ptrFromUUID(row.CreatedBy),
				UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
				UpdatedBy:          ptrFromUUID(row.UpdatedBy),
			},
			WorkerTypeCode:      row.WorkerTypeCode,
			WorkerTypeName:      row.WorkerTypeName,
			ClassificationGroup: row.ClassificationGroup,
			AttendanceMode:      row.AttendanceMode,
			PayMode:             row.PayMode,
			EmployeeCode:        ptrFromText(row.EmployeeCode),
			BranchName:          ptrFromText(row.BranchName),
			DepartmentName:      ptrFromText(row.DepartmentName),
		})
	}
	return items
}
