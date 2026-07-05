package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapCandidate(row sqlc.HrmsCandidate) *domain.Candidate {
	return &domain.Candidate{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		Firstname:          ptrFromText(row.Firstname),
		Lastname:           ptrFromText(row.Lastname),
		Email:              ptrFromText(row.Email),
		Phone:              ptrFromText(row.Phone),
		DOB:                ptrFromDate(row.Dob),
		Gender:             ptrFromText(row.Gender),
		TotalExperience:    floatPtrFromNumeric(row.TotalExperience),
		CurrentCompany:     ptrFromText(row.CurrentCompany),
		CurrentDesignation: ptrFromText(row.CurrentDesignation),
		CurrentSalary:      floatPtrFromNumeric(row.CurrentSalary),
		ExpectedSalary:     floatPtrFromNumeric(row.ExpectedSalary),
		NoticePeriod:       ptrFromInt4(row.NoticePeriod),
		CurrentLocation:    ptrFromText(row.CurrentLocation),
		PreferredLocation:  ptrFromText(row.PreferredLocation),
		Source:             ptrFromText(row.Source),
		ResumeURL:          ptrFromText(row.ResumeUrl),
		Inactive:           row.Inactive,
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		CreatedBy:          ptrFromUUID(row.CreatedBy),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:          ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidateListRow(row sqlc.ListCandidatesRow) *domain.Candidate {
	return &domain.Candidate{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		Firstname:          ptrFromText(row.Firstname),
		Lastname:           ptrFromText(row.Lastname),
		Email:              ptrFromText(row.Email),
		Phone:              ptrFromText(row.Phone),
		DOB:                ptrFromDate(row.Dob),
		Gender:             ptrFromText(row.Gender),
		TotalExperience:    floatPtrFromNumeric(row.TotalExperience),
		CurrentCompany:     ptrFromText(row.CurrentCompany),
		CurrentDesignation: ptrFromText(row.CurrentDesignation),
		CurrentSalary:      floatPtrFromNumeric(row.CurrentSalary),
		ExpectedSalary:     floatPtrFromNumeric(row.ExpectedSalary),
		NoticePeriod:       ptrFromInt4(row.NoticePeriod),
		CurrentLocation:    ptrFromText(row.CurrentLocation),
		PreferredLocation:  ptrFromText(row.PreferredLocation),
		Source:             ptrFromText(row.Source),
		ResumeURL:          ptrFromText(row.ResumeUrl),
		Inactive:           row.Inactive,
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		CreatedBy:          ptrFromUUID(row.CreatedBy),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:          ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidates(rows []sqlc.ListCandidatesRow) []*domain.Candidate {
	items := make([]*domain.Candidate, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCandidateListRow(row))
	}
	return items
}
