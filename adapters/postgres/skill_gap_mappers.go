package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapProjectSkillRequirement(row sqlc.HrmsProjectSkillRequirement) *domain.ProjectSkillRequirement {
	return projectSkillRequirementFromParts(row.ID, row.TenantID, row.ProjectID, row.EngagementID, row.SkillID, row.RequiredProficiency, row.MinYearsExperience, row.RequiredCount, row.Importance, row.RequirementSource, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, pgtype.UUID{}, nil, nil, nil, nil, nil, nil)
}

func mapProjectSkillRequirementRows(rows []sqlc.ListProjectSkillRequirementsRow) []*domain.ProjectSkillRequirement {
	items := make([]*domain.ProjectSkillRequirement, 0, len(rows))
	for _, row := range rows {
		skillName := row.SkillName
		skillCode := row.SkillCode
		skillType := row.SkillType
		items = append(items, projectSkillRequirementFromParts(row.ID, row.TenantID, row.ProjectID, row.EngagementID, row.SkillID, row.RequiredProficiency, row.MinYearsExperience, row.RequiredCount, row.Importance, row.RequirementSource, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.ProjectName), ptrFromText(row.ProjectCode), ptrFromText(row.EngagementTitle), ptrFromText(row.EngagementCode), row.WorkerProfileID, ptrFromText(row.WorkerDisplayName), ptrFromText(row.WorkerCode), &skillName, &skillCode, &skillType, ptrFromText(row.CategoryName)))
	}
	return items
}

func projectSkillRequirementFromParts(id uuid.UUID, tenantID uuid.UUID, projectID pgtype.UUID, engagementID pgtype.UUID, skillID uuid.UUID, requiredProficiency string, minYearsExperience pgtype.Numeric, requiredCount int32, importance string, requirementSource string, notes pgtype.Text, metadataBytes []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, projectName *string, projectCode *string, engagementTitle *string, engagementCode *string, workerProfileID pgtype.UUID, workerDisplayName *string, workerCode *string, skillName *string, skillCode *string, skillType *string, categoryName *string) *domain.ProjectSkillRequirement {
	metadata := json.RawMessage(metadataBytes)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.ProjectSkillRequirement{
		ID:                  id,
		TenantID:            tenantID,
		ProjectID:           ptrFromUUID(projectID),
		EngagementID:        ptrFromUUID(engagementID),
		SkillID:             skillID,
		RequiredProficiency: requiredProficiency,
		MinYearsExperience:  ptrFromNumeric(minYearsExperience),
		RequiredCount:       requiredCount,
		Importance:          importance,
		RequirementSource:   requirementSource,
		Notes:               ptrFromText(notes),
		Metadata:            metadata,
		Inactive:            inactive,
		CreatedAt:           timeFromTimestamptz(createdAt),
		CreatedBy:           ptrFromUUID(createdBy),
		UpdatedAt:           timeFromTimestamptz(updatedAt),
		UpdatedBy:           ptrFromUUID(updatedBy),
		ProjectName:         projectName,
		ProjectCode:         projectCode,
		EngagementTitle:     engagementTitle,
		EngagementCode:      engagementCode,
		WorkerProfileID:     ptrFromUUID(workerProfileID),
		WorkerDisplayName:   workerDisplayName,
		WorkerCode:          workerCode,
		SkillName:           skillName,
		SkillCode:           skillCode,
		SkillType:           skillType,
		CategoryName:        categoryName,
	}
}

func mapProjectSkillGapRows(rows []sqlc.ListProjectSkillGapRowsRow) []*domain.ProjectSkillGapRow {
	items := make([]*domain.ProjectSkillGapRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.ProjectSkillGapRow{
			RequirementID:          row.RequirementID,
			TenantID:               row.TenantID,
			ProjectID:              ptrFromUUID(row.ProjectID),
			ProjectName:            ptrFromText(row.ProjectName),
			ProjectCode:            ptrFromText(row.ProjectCode),
			EngagementID:           ptrFromUUID(row.EngagementID),
			EngagementTitle:        ptrFromText(row.EngagementTitle),
			EngagementCode:         ptrFromText(row.EngagementCode),
			SkillID:                row.SkillID,
			SkillName:              row.SkillName,
			SkillCode:              row.SkillCode,
			SkillType:              row.SkillType,
			RequiredProficiency:    row.RequiredProficiency,
			MinYearsExperience:     ptrFromNumeric(row.MinYearsExperience),
			RequiredCount:          row.RequiredCount,
			Importance:             row.Importance,
			AssignedMatchCount:     row.AssignedMatchCount,
			AvailableMatchCount:    row.AvailableMatchCount,
			GapCount:               row.GapCount,
			MatchPercent:           row.MatchPercent,
			SinglePersonDependency: row.SinglePersonDependency.Valid && row.SinglePersonDependency.Bool,
			SuggestedAction:        row.SuggestedAction,
		})
	}
	return items
}

func mapSkillGapSummaryRows(rows []sqlc.ListSkillGapSummaryRow) []*domain.SkillGapSummaryRow {
	items := make([]*domain.SkillGapSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.SkillGapSummaryRow{ProjectID: ptrFromUUID(row.ProjectID), ProjectName: row.ProjectName, ProjectCode: ptrFromText(row.ProjectCode), RequirementCount: row.RequirementCount, MissingSkillCount: row.MissingSkillCount, MandatoryGapCount: row.MandatoryGapCount, AverageMatchPercent: row.AverageMatchPercent, SinglePersonDependencyCount: row.SinglePersonDependencyCount})
	}
	return items
}

func mapSinglePersonSkillDependencyRows(rows []sqlc.ListSinglePersonSkillDependenciesRow) []*domain.SinglePersonSkillDependency {
	items := make([]*domain.SinglePersonSkillDependency, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.SinglePersonSkillDependency{
			RequirementID:     row.RequirementID,
			ProjectID:         ptrFromUUID(row.ProjectID),
			ProjectName:       ptrFromText(row.ProjectName),
			EngagementID:      ptrFromUUID(row.EngagementID),
			EngagementTitle:   ptrFromText(row.EngagementTitle),
			SkillID:           row.SkillID,
			SkillName:         row.SkillName,
			Importance:        row.Importance,
			WorkerProfileID:   row.WorkerProfileID,
			WorkerDisplayName: row.WorkerDisplayName,
			WorkerCode:        ptrFromText(row.WorkerCode),
			Proficiency:       row.Proficiency,
			YearsExperience:   ptrFromNumeric(row.YearsExperience),
		})
	}
	return items
}
