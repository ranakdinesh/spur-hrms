package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapInterviewRound(row sqlc.HrmsInterviewRound) *domain.InterviewRound {
	return &domain.InterviewRound{
		ID:                row.ID,
		TenantID:          row.TenantID,
		ApplicationID:     row.ApplicationID,
		RoundName:         ptrFromText(row.RoundName),
		RoundNumber:       ptrFromInt4(row.RoundNumber),
		ScheduledDate:     ptrFromTimestamptz(row.ScheduledDate),
		DurationMinutes:   ptrFromInt4(row.DurationMinutes),
		InterviewerUserID: ptrFromUUID(row.InterviewerUserID),
		Mode:              ptrFromText(row.Mode),
		MeetingLink:       ptrFromText(row.MeetingLink),
		Location:          ptrFromText(row.Location),
		Status:            row.Status,
		Remarks:           ptrFromText(row.Remarks),
		Timezone:          row.Timezone,
		Feedback:          ptrFromText(row.Feedback),
		Score:             floatPtrFromNumeric(row.Score),
		Decision:          ptrFromText(row.Decision),
		CompletedAt:       ptrFromTimestamptz(row.CompletedAt),
		Inactive:          row.Inactive,
		CreatedAt:         timeFromTimestamptz(row.CreatedAt),
		CreatedBy:         ptrFromUUID(row.CreatedBy),
		UpdatedAt:         timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:         ptrFromUUID(row.UpdatedBy),
	}
}

func mapInterviewRoundListRow(row sqlc.ListInterviewRoundsRow) *domain.InterviewRound {
	applicationStatus := row.ApplicationStatus
	return &domain.InterviewRound{
		ID:                 row.ID,
		TenantID:           row.TenantID,
		ApplicationID:      row.ApplicationID,
		ApplicationStatus:  &applicationStatus,
		CandidateFirstname: ptrFromText(row.CandidateFirstname),
		CandidateLastname:  ptrFromText(row.CandidateLastname),
		CandidateEmail:     ptrFromText(row.CandidateEmail),
		JobPostingTitle:    ptrFromText(row.JobPostingTitle),
		JobPostingCode:     ptrFromText(row.JobPostingCode),
		RoundName:          ptrFromText(row.RoundName),
		RoundNumber:        ptrFromInt4(row.RoundNumber),
		ScheduledDate:      ptrFromTimestamptz(row.ScheduledDate),
		DurationMinutes:    ptrFromInt4(row.DurationMinutes),
		InterviewerUserID:  ptrFromUUID(row.InterviewerUserID),
		Mode:               ptrFromText(row.Mode),
		MeetingLink:        ptrFromText(row.MeetingLink),
		Location:           ptrFromText(row.Location),
		Status:             row.Status,
		Remarks:            ptrFromText(row.Remarks),
		Timezone:           row.Timezone,
		Feedback:           ptrFromText(row.Feedback),
		Score:              floatPtrFromNumeric(row.Score),
		Decision:           ptrFromText(row.Decision),
		CompletedAt:        ptrFromTimestamptz(row.CompletedAt),
		Inactive:           row.Inactive,
		CreatedAt:          timeFromTimestamptz(row.CreatedAt),
		CreatedBy:          ptrFromUUID(row.CreatedBy),
		UpdatedAt:          timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:          ptrFromUUID(row.UpdatedBy),
	}
}

func mapInterviewRounds(rows []sqlc.ListInterviewRoundsRow) []*domain.InterviewRound {
	items := make([]*domain.InterviewRound, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapInterviewRoundListRow(row))
	}
	return items
}
