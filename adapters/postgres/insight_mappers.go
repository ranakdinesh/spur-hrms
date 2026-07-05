package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapInsight(row sqlc.HrmsInsight) *domain.Insight {
	return &domain.Insight{
		ID:              row.ID,
		TenantID:        row.TenantID,
		InsightKey:      row.InsightKey,
		InsightType:     row.InsightType,
		Category:        row.Category,
		Severity:        row.Severity,
		Status:          row.Status,
		Title:           row.Title,
		Summary:         row.Summary,
		ConfidenceScore: floatFromNumeric(row.ConfidenceScore),
		Score:           floatFromNumeric(row.Score),
		Source:          row.Source,
		ModelVersion:    ptrFromText(row.ModelVersion),
		EntityType:      ptrFromText(row.EntityType),
		EntityID:        ptrFromUUID(row.EntityID),
		EmployeeUserID:  ptrFromUUID(row.EmployeeUserID),
		Reasons:         row.Reasons,
		Recommendations: row.Recommendations,
		Context:         row.Context,
		Explainability:  row.Explainability,
		DetectedAt:      timeFromTimestamptz(row.DetectedAt),
		DueAt:           ptrFromTimestamptz(row.DueAt),
		AssignedTo:      ptrFromUUID(row.AssignedTo),
		ReviewedBy:      ptrFromUUID(row.ReviewedBy),
		ReviewedAt:      ptrFromTimestamptz(row.ReviewedAt),
		ResolvedAt:      ptrFromTimestamptz(row.ResolvedAt),
		ResolutionNote:  ptrFromText(row.ResolutionNote),
		Inactive:        row.Inactive,
		CreatedAt:       timeFromTimestamptz(row.CreatedAt),
		CreatedBy:       ptrFromUUID(row.CreatedBy),
		UpdatedAt:       timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:       ptrFromUUID(row.UpdatedBy),
	}
}

func mapInsights(rows []sqlc.HrmsInsight) []*domain.Insight {
	items := make([]*domain.Insight, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapInsight(row))
	}
	return items
}

func mapInsightEvent(row sqlc.HrmsInsightEvent) *domain.InsightEvent {
	return &domain.InsightEvent{
		ID:         row.ID,
		TenantID:   row.TenantID,
		InsightID:  row.InsightID,
		Action:     row.Action,
		FromStatus: ptrFromText(row.FromStatus),
		ToStatus:   ptrFromText(row.ToStatus),
		Remarks:    ptrFromText(row.Remarks),
		Metadata:   row.Metadata,
		Inactive:   row.Inactive,
		CreatedAt:  timeFromTimestamptz(row.CreatedAt),
		CreatedBy:  ptrFromUUID(row.CreatedBy),
		UpdatedAt:  timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:  ptrFromUUID(row.UpdatedBy),
	}
}

func mapInsightEvents(rows []sqlc.HrmsInsightEvent) []*domain.InsightEvent {
	items := make([]*domain.InsightEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapInsightEvent(row))
	}
	return items
}
