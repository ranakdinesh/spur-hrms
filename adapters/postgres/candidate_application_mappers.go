package postgres

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapCandidateApplicantAccount(row sqlc.HrmsCandidateApplicantAccount) *domain.CandidateApplicantAccount {
	metadata := map[string]any{}
	if len(row.Metadata) > 0 {
		_ = json.Unmarshal(row.Metadata, &metadata)
	}
	return &domain.CandidateApplicantAccount{
		ID:          row.ID,
		TenantID:    row.TenantID,
		CandidateID: row.CandidateID,
		UserID:      row.UserID,
		Email:       row.Email,
		Status:      row.Status,
		ConsentAt:   ptrFromTimestamptz(row.ConsentAt),
		ConsentIP:   ptrFromText(row.ConsentIp),
		Metadata:    metadata,
		Inactive:    row.Inactive,
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
		CreatedBy:   ptrFromUUID(row.CreatedBy),
		UpdatedAt:   timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:   ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidateApplication(row sqlc.HrmsCandidateApplication) *domain.CandidateApplication {
	return &domain.CandidateApplication{
		ID:                       row.ID,
		TenantID:                 row.TenantID,
		CandidateID:              ptrFromUUID(row.CandidateID),
		JobPostingID:             ptrFromUUID(row.JobPostingID),
		ResumeURL:                ptrFromText(row.ResumeUrl),
		CoverLetter:              ptrFromText(row.CoverLetter),
		CurrentCTC:               floatPtrFromNumeric(row.CurrentCtc),
		ExpectedCTC:              floatPtrFromNumeric(row.ExpectedCtc),
		NoticePeriod:             ptrFromInt4(row.NoticePeriod),
		ReferredBy:               ptrFromText(row.ReferredBy),
		Source:                   ptrFromText(row.Source),
		SourceDetail:             ptrFromText(row.SourceDetail),
		Status:                   row.Status,
		Comments:                 ptrFromText(row.Comments),
		AppliedAt:                timeFromTimestamptz(row.AppliedAt),
		StatusChangedAt:          timeFromTimestamptz(row.StatusChangedAt),
		StatusChangedBy:          ptrFromUUID(row.StatusChangedBy),
		RejectionReason:          ptrFromText(row.RejectionReason),
		WithdrawalReason:         ptrFromText(row.WithdrawalReason),
		DuplicateOfApplicationID: ptrFromUUID(row.DuplicateOfApplicationID),
		DaysInStage:              daysSince(row.StatusChangedAt),
		Inactive:                 row.Inactive,
		CreatedAt:                timeFromTimestamptz(row.CreatedAt),
		CreatedBy:                ptrFromUUID(row.CreatedBy),
		UpdatedAt:                timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:                ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidateApplicationListRow(row sqlc.ListCandidateApplicationsRow) *domain.CandidateApplication {
	return &domain.CandidateApplication{
		ID:                       row.ID,
		TenantID:                 row.TenantID,
		CandidateID:              ptrFromUUID(row.CandidateID),
		CandidateFirstname:       ptrFromText(row.CandidateFirstname),
		CandidateLastname:        ptrFromText(row.CandidateLastname),
		CandidateEmail:           ptrFromText(row.CandidateEmail),
		CandidatePhone:           ptrFromText(row.CandidatePhone),
		JobPostingID:             ptrFromUUID(row.JobPostingID),
		JobPostingTitle:          ptrFromText(row.JobPostingTitle),
		JobPostingCode:           ptrFromText(row.JobPostingCode),
		ResumeURL:                ptrFromText(row.ResumeUrl),
		CoverLetter:              ptrFromText(row.CoverLetter),
		CurrentCTC:               floatPtrFromNumeric(row.CurrentCtc),
		ExpectedCTC:              floatPtrFromNumeric(row.ExpectedCtc),
		NoticePeriod:             ptrFromInt4(row.NoticePeriod),
		ReferredBy:               ptrFromText(row.ReferredBy),
		Source:                   ptrFromText(row.Source),
		SourceDetail:             ptrFromText(row.SourceDetail),
		Status:                   row.Status,
		Comments:                 ptrFromText(row.Comments),
		AppliedAt:                timeFromTimestamptz(row.AppliedAt),
		StatusChangedAt:          timeFromTimestamptz(row.StatusChangedAt),
		StatusChangedBy:          ptrFromUUID(row.StatusChangedBy),
		RejectionReason:          ptrFromText(row.RejectionReason),
		WithdrawalReason:         ptrFromText(row.WithdrawalReason),
		DuplicateOfApplicationID: ptrFromUUID(row.DuplicateOfApplicationID),
		DaysInStage:              daysSince(row.StatusChangedAt),
		Inactive:                 row.Inactive,
		CreatedAt:                timeFromTimestamptz(row.CreatedAt),
		CreatedBy:                ptrFromUUID(row.CreatedBy),
		UpdatedAt:                timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:                ptrFromUUID(row.UpdatedBy),
	}
}

func mapCandidateApplicationEvent(row sqlc.HrmsCandidateApplicationEvent) *domain.CandidateApplicationEvent {
	return &domain.CandidateApplicationEvent{ID: row.ID, TenantID: row.TenantID, ApplicationID: row.ApplicationID, FromStatus: ptrFromText(row.FromStatus), ToStatus: row.ToStatus, Action: row.Action, Reason: ptrFromText(row.Reason), Remarks: ptrFromText(row.Remarks), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapCandidateApplicationEvents(rows []sqlc.HrmsCandidateApplicationEvent) []*domain.CandidateApplicationEvent {
	items := make([]*domain.CandidateApplicationEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCandidateApplicationEvent(row))
	}
	return items
}

func daysSince(value pgtype.Timestamptz) int {
	if !value.Valid {
		return 0
	}
	days := int(time.Since(value.Time).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

func mapCandidateApplications(rows []sqlc.ListCandidateApplicationsRow) []*domain.CandidateApplication {
	items := make([]*domain.CandidateApplication, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCandidateApplicationListRow(row))
	}
	return items
}
