package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapEmployee(row sqlc.HrmsEmployee) *domain.Employee {
	return &domain.Employee{
		ID:                    row.ID,
		TenantID:              row.TenantID,
		UserID:                row.UserID,
		EmployeeCode:          ptrFromText(row.EmployeeCode),
		Firstname:             row.Firstname,
		MiddleName:            ptrFromText(row.MiddleName),
		Lastname:              ptrFromText(row.Lastname),
		Email:                 ptrFromText(row.Email),
		Mobile:                ptrFromText(row.Mobile),
		DOB:                   ptrFromDate(row.Dob),
		Gender:                ptrFromText(row.Gender),
		MaritalStatus:         ptrFromText(row.MaritalStatus),
		BloodGroup:            ptrFromText(row.BloodGroup),
		ProfilePhotoPath:      ptrFromText(row.ProfilePhotoPath),
		Address:               ptrFromText(row.Address),
		City:                  ptrFromText(row.City),
		State:                 ptrFromText(row.State),
		Country:               ptrFromText(row.Country),
		Pincode:               ptrFromText(row.Pincode),
		EmergencyContact:      ptrFromText(row.EmergencyContact),
		JoiningDate:           ptrFromDate(row.JoiningDate),
		ResignationDate:       ptrFromDate(row.ResignationDate),
		DepartmentID:          ptrFromUUID(row.DepartmentID),
		BranchID:              ptrFromUUID(row.BranchID),
		DesignationID:         ptrFromUUID(row.DesignationID),
		ReportingManagerID:    ptrFromUUID(row.ReportingManagerID),
		EmploymentTypeID:      ptrFromUUID(row.EmploymentTypeID),
		Role:                  ptrFromText(row.Role),
		Grade:                 ptrFromText(row.Grade),
		ExperienceYear:        row.ExperienceYear,
		ExperienceMonth:       row.ExperienceMonth,
		ProbationStatus:       row.ProbationStatus,
		ProbationStartDate:    ptrFromDate(row.ProbationStartDate),
		ProbationEndDate:      ptrFromDate(row.ProbationEndDate),
		ProbationDurationDays: row.ProbationDurationDays,
		ProbationConfirmedAt:  ptrFromDate(row.ProbationConfirmedAt),
		IsPayrollStaff:        row.IsPayrollStaff,
		Inactive:              row.Inactive,
		CreatedAt:             timeFromTimestamptz(row.CreatedAt),
		CreatedBy:             ptrFromUUID(row.CreatedBy),
		UpdatedAt:             timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:             ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmployeeListItem(row sqlc.ListEmployeesRow) *domain.EmployeeListItem {
	employee := mapEmployee(sqlc.HrmsEmployee{
		ID:                    row.ID,
		TenantID:              row.TenantID,
		UserID:                row.UserID,
		EmployeeCode:          row.EmployeeCode,
		Firstname:             row.Firstname,
		MiddleName:            row.MiddleName,
		Lastname:              row.Lastname,
		Email:                 row.Email,
		Mobile:                row.Mobile,
		Dob:                   row.Dob,
		Gender:                row.Gender,
		MaritalStatus:         row.MaritalStatus,
		BloodGroup:            row.BloodGroup,
		ProfilePhotoPath:      row.ProfilePhotoPath,
		Address:               row.Address,
		City:                  row.City,
		State:                 row.State,
		Country:               row.Country,
		Pincode:               row.Pincode,
		EmergencyContact:      row.EmergencyContact,
		JoiningDate:           row.JoiningDate,
		ResignationDate:       row.ResignationDate,
		DepartmentID:          row.DepartmentID,
		BranchID:              row.BranchID,
		DesignationID:         row.DesignationID,
		ReportingManagerID:    row.ReportingManagerID,
		EmploymentTypeID:      row.EmploymentTypeID,
		Role:                  row.Role,
		Grade:                 row.Grade,
		ExperienceYear:        row.ExperienceYear,
		ExperienceMonth:       row.ExperienceMonth,
		ProbationStatus:       row.ProbationStatus,
		ProbationStartDate:    row.ProbationStartDate,
		ProbationEndDate:      row.ProbationEndDate,
		ProbationDurationDays: row.ProbationDurationDays,
		ProbationConfirmedAt:  row.ProbationConfirmedAt,
		IsPayrollStaff:        row.IsPayrollStaff,
		Inactive:              row.Inactive,
		CreatedAt:             row.CreatedAt,
		CreatedBy:             row.CreatedBy,
		UpdatedAt:             row.UpdatedAt,
		UpdatedBy:             row.UpdatedBy,
	})
	return &domain.EmployeeListItem{
		Employee:           *employee,
		DepartmentName:     ptrFromText(row.DepartmentName),
		BranchName:         ptrFromText(row.BranchName),
		DesignationName:    ptrFromText(row.DesignationName),
		AttendanceRequired: row.AttendanceRequired,
		EmploymentTypeName: ptrFromText(row.EmploymentTypeName),
	}
}

func mapEmployeeProfileItem(row sqlc.GetEmployeeProfileItemRow) *domain.EmployeeListItem {
	employee := mapEmployee(sqlc.HrmsEmployee{
		ID:                    row.ID,
		TenantID:              row.TenantID,
		UserID:                row.UserID,
		EmployeeCode:          row.EmployeeCode,
		Firstname:             row.Firstname,
		MiddleName:            row.MiddleName,
		Lastname:              row.Lastname,
		Email:                 row.Email,
		Mobile:                row.Mobile,
		Dob:                   row.Dob,
		Gender:                row.Gender,
		MaritalStatus:         row.MaritalStatus,
		BloodGroup:            row.BloodGroup,
		ProfilePhotoPath:      row.ProfilePhotoPath,
		Address:               row.Address,
		City:                  row.City,
		State:                 row.State,
		Country:               row.Country,
		Pincode:               row.Pincode,
		EmergencyContact:      row.EmergencyContact,
		JoiningDate:           row.JoiningDate,
		ResignationDate:       row.ResignationDate,
		DepartmentID:          row.DepartmentID,
		BranchID:              row.BranchID,
		DesignationID:         row.DesignationID,
		ReportingManagerID:    row.ReportingManagerID,
		EmploymentTypeID:      row.EmploymentTypeID,
		Role:                  row.Role,
		Grade:                 row.Grade,
		ExperienceYear:        row.ExperienceYear,
		ExperienceMonth:       row.ExperienceMonth,
		ProbationStatus:       row.ProbationStatus,
		ProbationStartDate:    row.ProbationStartDate,
		ProbationEndDate:      row.ProbationEndDate,
		ProbationDurationDays: row.ProbationDurationDays,
		ProbationConfirmedAt:  row.ProbationConfirmedAt,
		IsPayrollStaff:        row.IsPayrollStaff,
		Inactive:              row.Inactive,
		CreatedAt:             row.CreatedAt,
		CreatedBy:             row.CreatedBy,
		UpdatedAt:             row.UpdatedAt,
		UpdatedBy:             row.UpdatedBy,
	})
	return &domain.EmployeeListItem{
		Employee:           *employee,
		DepartmentName:     ptrFromText(row.DepartmentName),
		BranchName:         ptrFromText(row.BranchName),
		DesignationName:    ptrFromText(row.DesignationName),
		AttendanceRequired: row.AttendanceRequired,
		EmploymentTypeName: ptrFromText(row.EmploymentTypeName),
	}
}

func mapEmployeeStatutory(row sqlc.HrmsEmployeeStatutory) *domain.EmployeeStatutory {
	return &domain.EmployeeStatutory{
		ID:             row.ID,
		TenantID:       row.TenantID,
		UserID:         row.UserID,
		PFNo:           ptrFromText(row.PfNo),
		UANNo:          ptrFromText(row.UanNo),
		ESICNo:         ptrFromText(row.EsicNo),
		PAN:            ptrFromText(row.Pan),
		Aadhaar:        ptrFromText(row.Aadhaar),
		PTApplicable:   row.PtApplicable,
		PFApplicable:   row.PfApplicable,
		ESICApplicable: row.EsicApplicable,
		LWFApplicable:  row.LwfApplicable,
		Inactive:       row.Inactive,
		CreatedAt:      timeFromTimestamptz(row.CreatedAt),
		CreatedBy:      ptrFromUUID(row.CreatedBy),
		UpdatedAt:      timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:      ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmployeeBank(row sqlc.HrmsEmployeeBank) *domain.EmployeeBank {
	return &domain.EmployeeBank{
		ID:            row.ID,
		TenantID:      row.TenantID,
		UserID:        row.UserID,
		BankName:      ptrFromText(row.BankName),
		AccountNumber: ptrFromText(row.AccountNumber),
		IFSCCode:      ptrFromText(row.IfscCode),
		AccountType:   ptrFromText(row.AccountType),
		BranchName:    ptrFromText(row.BranchName),
		IsPrimary:     row.IsPrimary,
		Inactive:      row.Inactive,
		CreatedAt:     timeFromTimestamptz(row.CreatedAt),
		CreatedBy:     ptrFromUUID(row.CreatedBy),
		UpdatedAt:     timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:     ptrFromUUID(row.UpdatedBy),
	}
}

func mapDocumentType(row sqlc.HrmsDocumentType) *domain.DocumentType {
	return &domain.DocumentType{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		Name:                row.Name,
		Description:         ptrFromText(row.Description),
		IsRequired:          row.IsRequired,
		Instructions:        ptrFromText(row.Instructions),
		AllowedContentTypes: row.AllowedContentTypes,
		MaxFileSizeBytes:    row.MaxFileSizeBytes,
		DisplayOrder:        row.DisplayOrder,
		Inactive:            row.Inactive,
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
		CreatedBy:           ptrFromUUID(row.CreatedBy),
		UpdatedAt:           timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:           ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmployeeDocument(row sqlc.ListEmployeeDocumentsByUserIDRow) *domain.EmployeeDocument {
	return &domain.EmployeeDocument{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		UserID:              row.UserID,
		DocumentTypeID:      ptrFromUUID(row.DocumentTypeID),
		DocumentTypeName:    ptrFromText(row.DocumentTypeName),
		Title:               ptrFromText(row.Title),
		FilePath:            ptrFromText(row.FilePath),
		Status:              row.Status,
		ReviewRemarks:       ptrFromText(row.ReviewRemarks),
		ReviewedBy:          ptrFromUUID(row.ReviewedBy),
		ReviewedAt:          ptrFromTimestamptz(row.ReviewedAt),
		OriginalFileName:    ptrFromText(row.OriginalFileName),
		ContentType:         ptrFromText(row.ContentType),
		FileSizeBytes:       ptrFromInt8(row.FileSizeBytes),
		Encrypted:           row.Encrypted,
		EncryptionAlgorithm: row.EncryptionAlgorithm,
		Inactive:            row.Inactive,
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
		CreatedBy:           ptrFromUUID(row.CreatedBy),
		UpdatedAt:           timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:           ptrFromUUID(row.UpdatedBy),
	}
}

func mapEmployeeDocumentRecord(row sqlc.HrmsEmployeeDocument) *domain.EmployeeDocument {
	return &domain.EmployeeDocument{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		UserID:              row.UserID,
		DocumentTypeID:      ptrFromUUID(row.DocumentTypeID),
		Title:               ptrFromText(row.Title),
		FilePath:            ptrFromText(row.FilePath),
		Status:              row.Status,
		ReviewRemarks:       ptrFromText(row.ReviewRemarks),
		ReviewedBy:          ptrFromUUID(row.ReviewedBy),
		ReviewedAt:          ptrFromTimestamptz(row.ReviewedAt),
		OriginalFileName:    ptrFromText(row.OriginalFileName),
		ContentType:         ptrFromText(row.ContentType),
		FileSizeBytes:       ptrFromInt8(row.FileSizeBytes),
		Encrypted:           row.Encrypted,
		EncryptionAlgorithm: row.EncryptionAlgorithm,
		Inactive:            row.Inactive,
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
		CreatedBy:           ptrFromUUID(row.CreatedBy),
		UpdatedAt:           timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:           ptrFromUUID(row.UpdatedBy),
	}
}

func ptrFromInt8(value pgtype.Int8) *int64 {
	if !value.Valid {
		return nil
	}
	return &value.Int64
}
