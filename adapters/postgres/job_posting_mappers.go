package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapJobPosting(row sqlc.HrmsJobPosting) *domain.JobPosting {
	return &domain.JobPosting{
		ID:               row.ID,
		TenantID:         row.TenantID,
		JobRequisitionID: ptrFromUUID(row.JobRequisitionID),
		Code:             ptrFromText(row.Code),
		Title:            ptrFromText(row.Title),
		JobSummary:       ptrFromText(row.JobSummary),
		Description:      ptrFromText(row.Description),
		JobCategory:      ptrFromText(row.JobCategory),
		DepartmentID:     ptrFromUUID(row.DepartmentID),
		Industry:         ptrFromText(row.Industry),
		EmploymentTypeID: ptrFromUUID(row.EmploymentTypeID),
		WorkMode:         ptrFromText(row.WorkMode),
		RoleType:         ptrFromText(row.RoleType),
		MinExperience:    floatPtrFromNumeric(row.MinExperience),
		MaxExperience:    floatPtrFromNumeric(row.MaxExperience),
		MinSalary:        floatPtrFromNumeric(row.MinSalary),
		MaxSalary:        floatPtrFromNumeric(row.MaxSalary),
		SalaryCurrency:   ptrFromText(row.SalaryCurrency),
		SalaryPeriod:     ptrFromText(row.SalaryPeriod),
		IsSalaryVisible:  row.IsSalaryVisible,
		JobStatus:        ptrFromText(row.JobStatus),
		PublishDate:      ptrFromDate(row.PublishDate),
		ExpiryDate:       ptrFromDate(row.ExpiryDate),
		IsPublished:      row.IsPublished,
		Inactive:         row.Inactive,
		CreatedAt:        timeFromTimestamptz(row.CreatedAt),
		CreatedBy:        ptrFromUUID(row.CreatedBy),
		UpdatedAt:        timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:        ptrFromUUID(row.UpdatedBy),
	}
}

func mapJobPostingListRow(row sqlc.ListJobPostingsRow) *domain.JobPosting {
	status := row.EffectiveStatus
	return &domain.JobPosting{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		JobRequisitionID:     ptrFromUUID(row.JobRequisitionID),
		JobRequisitionCode:   ptrFromText(row.JobRequisitionCode),
		JobRequisitionStatus: ptrFromText(row.JobRequisitionStatus),
		Code:                 ptrFromText(row.Code),
		Title:                ptrFromText(row.Title),
		JobSummary:           ptrFromText(row.JobSummary),
		Description:          ptrFromText(row.Description),
		JobCategory:          ptrFromText(row.JobCategory),
		DepartmentID:         ptrFromUUID(row.DepartmentID),
		DepartmentName:       ptrFromText(row.DepartmentName),
		Industry:             ptrFromText(row.Industry),
		EmploymentTypeID:     ptrFromUUID(row.EmploymentTypeID),
		EmploymentTypeName:   ptrFromText(row.EmploymentTypeName),
		WorkMode:             ptrFromText(row.WorkMode),
		RoleType:             ptrFromText(row.RoleType),
		MinExperience:        floatPtrFromNumeric(row.MinExperience),
		MaxExperience:        floatPtrFromNumeric(row.MaxExperience),
		MinSalary:            floatPtrFromNumeric(row.MinSalary),
		MaxSalary:            floatPtrFromNumeric(row.MaxSalary),
		SalaryCurrency:       ptrFromText(row.SalaryCurrency),
		SalaryPeriod:         ptrFromText(row.SalaryPeriod),
		IsSalaryVisible:      row.IsSalaryVisible,
		EffectiveStatus:      &status,
		JobStatus:            ptrFromText(row.JobStatus),
		PublishDate:          ptrFromDate(row.PublishDate),
		ExpiryDate:           ptrFromDate(row.ExpiryDate),
		IsPublished:          row.IsPublished,
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapJobPostings(rows []sqlc.ListJobPostingsRow) []*domain.JobPosting {
	items := make([]*domain.JobPosting, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapJobPostingListRow(row))
	}
	return items
}

func mapPublishedJobPostings(rows []sqlc.HrmsJobPosting) []*domain.JobPosting {
	items := make([]*domain.JobPosting, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapJobPosting(row))
	}
	return items
}
