package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapTalentMarketplaceOpportunity(row sqlc.HrmsTalentMarketplaceOpportunity) *domain.TalentMarketplaceOpportunity {
	return talentOpportunityFromParts(row.ID, row.TenantID, row.ProjectID, row.EngagementID, row.SourceRequirementID, row.JobPostingID, row.Title, row.Description, row.OpportunityType, row.Status, row.Visibility, row.Priority, row.Seats, row.LocationMode, row.MinAllocationPercent, row.DurationLabel, row.StartDate, row.DueDate, row.CandidateFallbackEnabled, row.CandidateFallbackStatus, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil, nil, 0, 0, 0)
}

func mapTalentMarketplaceOpportunityRows(rows []sqlc.ListTalentMarketplaceOpportunitiesRow) []*domain.TalentMarketplaceOpportunity {
	items := make([]*domain.TalentMarketplaceOpportunity, 0, len(rows))
	for _, row := range rows {
		items = append(items, talentOpportunityFromParts(row.ID, row.TenantID, row.ProjectID, row.EngagementID, row.SourceRequirementID, row.JobPostingID, row.Title, row.Description, row.OpportunityType, row.Status, row.Visibility, row.Priority, row.Seats, row.LocationMode, row.MinAllocationPercent, row.DurationLabel, row.StartDate, row.DueDate, row.CandidateFallbackEnabled, row.CandidateFallbackStatus, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.ProjectName), ptrFromText(row.ProjectCode), ptrFromText(row.EngagementTitle), ptrFromText(row.EngagementCode), ptrFromText(row.JobPostingTitle), ptrFromText(row.JobPostingCode), row.ApplicationCount, row.RecommendedCount, row.SelectedCount))
	}
	return items
}

func talentOpportunityFromParts(id uuid.UUID, tenantID uuid.UUID, projectID pgtype.UUID, engagementID pgtype.UUID, sourceRequirementID pgtype.UUID, jobPostingID pgtype.UUID, title string, description pgtype.Text, opportunityType string, status string, visibility string, priority string, seats int32, locationMode string, minAllocationPercent pgtype.Int4, durationLabel pgtype.Text, startDate pgtype.Date, dueDate pgtype.Date, fallbackEnabled bool, fallbackStatus string, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, projectName *string, projectCode *string, engagementTitle *string, engagementCode *string, jobPostingTitle *string, jobPostingCode *string, applicationCount int32, recommendedCount int32, selectedCount int32) *domain.TalentMarketplaceOpportunity {
	metadata := json.RawMessage(metadataBytes)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.TalentMarketplaceOpportunity{
		ID:                       id,
		TenantID:                 tenantID,
		ProjectID:                ptrFromUUID(projectID),
		EngagementID:             ptrFromUUID(engagementID),
		SourceRequirementID:      ptrFromUUID(sourceRequirementID),
		JobPostingID:             ptrFromUUID(jobPostingID),
		Title:                    title,
		Description:              ptrFromText(description),
		OpportunityType:          opportunityType,
		Status:                   status,
		Visibility:               visibility,
		Priority:                 priority,
		Seats:                    seats,
		LocationMode:             locationMode,
		MinAllocationPercent:     ptrFromInt4(minAllocationPercent),
		DurationLabel:            ptrFromText(durationLabel),
		StartDate:                ptrFromDate(startDate),
		DueDate:                  ptrFromDate(dueDate),
		CandidateFallbackEnabled: fallbackEnabled,
		CandidateFallbackStatus:  fallbackStatus,
		Metadata:                 metadata,
		Inactive:                 inactive,
		CreatedAt:                timeFromTimestamptz(createdAt),
		CreatedBy:                ptrFromUUID(createdBy),
		UpdatedAt:                timeFromTimestamptz(updatedAt),
		UpdatedBy:                ptrFromUUID(updatedBy),
		ProjectName:              projectName,
		ProjectCode:              projectCode,
		EngagementTitle:          engagementTitle,
		EngagementCode:           engagementCode,
		JobPostingTitle:          jobPostingTitle,
		JobPostingCode:           jobPostingCode,
		ApplicationCount:         applicationCount,
		RecommendedCount:         recommendedCount,
		SelectedCount:            selectedCount,
	}
}

func mapTalentMarketplaceApplication(row sqlc.HrmsTalentMarketplaceApplication) *domain.TalentMarketplaceApplication {
	return talentApplicationFromParts(row.ID, row.TenantID, row.OpportunityID, row.WorkerProfileID, row.Status, row.MatchScore, row.MatchReasons, row.WorkerNote, row.ManagerNote, row.DecidedAt, row.DecidedBy, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil, nil)
}

func mapTalentMarketplaceApplicationRows(rows []sqlc.ListTalentMarketplaceApplicationsRow) []*domain.TalentMarketplaceApplication {
	items := make([]*domain.TalentMarketplaceApplication, 0, len(rows))
	for _, row := range rows {
		opportunityTitle := row.OpportunityTitle
		opportunityStatus := row.OpportunityStatus
		workerDisplayName := row.WorkerDisplayName
		items = append(items, talentApplicationFromParts(row.ID, row.TenantID, row.OpportunityID, row.WorkerProfileID, row.Status, row.MatchScore, row.MatchReasons, row.WorkerNote, row.ManagerNote, row.DecidedAt, row.DecidedBy, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &opportunityTitle, &opportunityStatus, &workerDisplayName, ptrFromText(row.WorkerCode), ptrFromText(row.ProjectName), ptrFromText(row.EngagementTitle)))
	}
	return items
}

func talentApplicationFromParts(id uuid.UUID, tenantID uuid.UUID, opportunityID uuid.UUID, workerProfileID uuid.UUID, status string, matchScore pgtype.Numeric, matchReasonsBytes []byte, workerNote pgtype.Text, managerNote pgtype.Text, decidedAt pgtype.Timestamptz, decidedBy pgtype.UUID, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, opportunityTitle *string, opportunityStatus *string, workerDisplayName *string, workerCode *string, projectName *string, engagementTitle *string) *domain.TalentMarketplaceApplication {
	matchReasons := json.RawMessage(matchReasonsBytes)
	if len(matchReasons) == 0 {
		matchReasons = json.RawMessage(`{}`)
	}
	return &domain.TalentMarketplaceApplication{
		ID:                id,
		TenantID:          tenantID,
		OpportunityID:     opportunityID,
		WorkerProfileID:   workerProfileID,
		Status:            status,
		MatchScore:        floatPtrFromNumeric(matchScore),
		MatchReasons:      matchReasons,
		WorkerNote:        ptrFromText(workerNote),
		ManagerNote:       ptrFromText(managerNote),
		DecidedAt:         ptrFromTimestamptz(decidedAt),
		DecidedBy:         ptrFromUUID(decidedBy),
		Inactive:          inactive,
		CreatedAt:         timeFromTimestamptz(createdAt),
		CreatedBy:         ptrFromUUID(createdBy),
		UpdatedAt:         timeFromTimestamptz(updatedAt),
		UpdatedBy:         ptrFromUUID(updatedBy),
		OpportunityTitle:  opportunityTitle,
		OpportunityStatus: opportunityStatus,
		WorkerDisplayName: workerDisplayName,
		WorkerCode:        workerCode,
		ProjectName:       projectName,
		EngagementTitle:   engagementTitle,
	}
}

func mapTalentMarketplaceRecommendations(rows []sqlc.ListTalentMarketplaceRecommendationsRow) []*domain.TalentMarketplaceRecommendation {
	items := make([]*domain.TalentMarketplaceRecommendation, 0, len(rows))
	for _, row := range rows {
		reasons := json.RawMessage(row.MatchReasons)
		if len(reasons) == 0 {
			reasons = json.RawMessage(`{}`)
		}
		score := floatFromNumeric(row.MatchScore)
		items = append(items, &domain.TalentMarketplaceRecommendation{
			WorkerProfileID:    row.WorkerProfileID,
			WorkerDisplayName:  row.WorkerDisplayName,
			WorkerCode:         ptrFromText(row.WorkerCode),
			RequiredSkillCount: row.RequiredSkillCount,
			MatchedSkillCount:  row.MatchedSkillCount,
			MissingSkillCount:  row.MissingSkillCount,
			MatchScore:         score,
			MatchReasons:       reasons,
			ApplicationID:      ptrFromUUID(row.ApplicationID),
			ApplicationStatus:  ptrFromText(row.ApplicationStatus),
		})
	}
	return items
}

func mapTalentMarketplaceEvent(row sqlc.HrmsTalentMarketplaceEvent) *domain.TalentMarketplaceEvent {
	return talentEventFromParts(row.ID, row.TenantID, row.OpportunityID, row.ApplicationID, row.ActorUserID, row.Action, row.FromStatus, row.ToStatus, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, nil, pgtype.UUID{}, nil)
}

func mapTalentMarketplaceEventRows(rows []sqlc.ListTalentMarketplaceEventsRow) []*domain.TalentMarketplaceEvent {
	items := make([]*domain.TalentMarketplaceEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, talentEventFromParts(row.ID, row.TenantID, row.OpportunityID, row.ApplicationID, row.ActorUserID, row.Action, row.FromStatus, row.ToStatus, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, ptrFromText(row.OpportunityTitle), row.WorkerProfileID, ptrFromText(row.WorkerDisplayName)))
	}
	return items
}

func talentEventFromParts(id uuid.UUID, tenantID uuid.UUID, opportunityID pgtype.UUID, applicationID pgtype.UUID, actorUserID pgtype.UUID, action string, fromStatus pgtype.Text, toStatus pgtype.Text, notes pgtype.Text, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, opportunityTitle *string, workerProfileID pgtype.UUID, workerDisplayName *string) *domain.TalentMarketplaceEvent {
	metadata := json.RawMessage(metadataBytes)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.TalentMarketplaceEvent{
		ID:                id,
		TenantID:          tenantID,
		OpportunityID:     ptrFromUUID(opportunityID),
		ApplicationID:     ptrFromUUID(applicationID),
		ActorUserID:       ptrFromUUID(actorUserID),
		Action:            action,
		FromStatus:        ptrFromText(fromStatus),
		ToStatus:          ptrFromText(toStatus),
		Notes:             ptrFromText(notes),
		Metadata:          metadata,
		Inactive:          inactive,
		CreatedAt:         timeFromTimestamptz(createdAt),
		OpportunityTitle:  opportunityTitle,
		WorkerProfileID:   ptrFromUUID(workerProfileID),
		WorkerDisplayName: workerDisplayName,
	}
}
