package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapAISignalLog(row sqlc.HrmsAiSignalLog) *domain.AISignalLog {
	return &domain.AISignalLog{
		ID:               row.ID,
		TenantID:         row.TenantID,
		SignalKey:        row.SignalKey,
		SignalType:       row.SignalType,
		SourceModule:     row.SourceModule,
		SourceEvent:      row.SourceEvent,
		Severity:         row.Severity,
		ProcessingStatus: row.ProcessingStatus,
		EntityType:       ptrFromText(row.EntityType),
		EntityID:         ptrFromUUID(row.EntityID),
		EmployeeUserID:   ptrFromUUID(row.EmployeeUserID),
		VisibilityScope:  row.VisibilityScope,
		IdempotencyKey:   ptrFromText(row.IdempotencyKey),
		CorrelationID:    ptrFromText(row.CorrelationID),
		Payload:          row.Payload,
		Explainability:   row.Explainability,
		OccurredAt:       timeFromTimestamptz(row.OccurredAt),
		ProcessedAt:      ptrFromTimestamptz(row.ProcessedAt),
		ErrorMessage:     ptrFromText(row.ErrorMessage),
		Inactive:         row.Inactive,
		CreatedAt:        timeFromTimestamptz(row.CreatedAt),
		CreatedBy:        ptrFromUUID(row.CreatedBy),
		UpdatedAt:        timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:        ptrFromUUID(row.UpdatedBy),
	}
}

func mapAISignalLogs(rows []sqlc.HrmsAiSignalLog) []*domain.AISignalLog {
	items := make([]*domain.AISignalLog, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAISignalLog(row))
	}
	return items
}

func mapAIAgentActionLog(row sqlc.HrmsAiAgentActionLog) *domain.AIAgentActionLog {
	return &domain.AIAgentActionLog{
		ID:                  row.ID,
		TenantID:            row.TenantID,
		ActionKey:           row.ActionKey,
		AgentKey:            row.AgentKey,
		AgentName:           row.AgentName,
		ActionType:          row.ActionType,
		Status:              row.Status,
		Severity:            row.Severity,
		Title:               row.Title,
		Summary:             row.Summary,
		InsightID:           ptrFromUUID(row.InsightID),
		SignalID:            ptrFromUUID(row.SignalID),
		EntityType:          ptrFromText(row.EntityType),
		EntityID:            ptrFromUUID(row.EntityID),
		EmployeeUserID:      ptrFromUUID(row.EmployeeUserID),
		VisibilityScope:     row.VisibilityScope,
		ProposedAction:      row.ProposedAction,
		InputSnapshot:       row.InputSnapshot,
		OutputSnapshot:      row.OutputSnapshot,
		Explainability:      row.Explainability,
		ConfidenceScore:     floatFromNumeric(row.ConfidenceScore),
		ModelVersion:        ptrFromText(row.ModelVersion),
		SidecarRunID:        ptrFromText(row.SidecarRunID),
		RequiresHumanReview: row.RequiresHumanReview,
		ReviewedBy:          ptrFromUUID(row.ReviewedBy),
		ReviewedAt:          ptrFromTimestamptz(row.ReviewedAt),
		ExecutedAt:          ptrFromTimestamptz(row.ExecutedAt),
		FailedAt:            ptrFromTimestamptz(row.FailedAt),
		FailureMessage:      ptrFromText(row.FailureMessage),
		Inactive:            row.Inactive,
		CreatedAt:           timeFromTimestamptz(row.CreatedAt),
		CreatedBy:           ptrFromUUID(row.CreatedBy),
		UpdatedAt:           timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:           ptrFromUUID(row.UpdatedBy),
	}
}

func mapAIAgentActionLogs(rows []sqlc.HrmsAiAgentActionLog) []*domain.AIAgentActionLog {
	items := make([]*domain.AIAgentActionLog, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAIAgentActionLog(row))
	}
	return items
}

func mapAIHumanOverride(row sqlc.HrmsAiHumanOverride) *domain.AIHumanOverride {
	return &domain.AIHumanOverride{
		ID:             row.ID,
		TenantID:       row.TenantID,
		InsightID:      ptrFromUUID(row.InsightID),
		ActionID:       ptrFromUUID(row.ActionID),
		OverrideType:   row.OverrideType,
		OriginalStatus: ptrFromText(row.OriginalStatus),
		OverrideStatus: row.OverrideStatus,
		Reason:         row.Reason,
		Decision:       row.Decision,
		Metadata:       row.Metadata,
		Inactive:       row.Inactive,
		CreatedAt:      timeFromTimestamptz(row.CreatedAt),
		CreatedBy:      ptrFromUUID(row.CreatedBy),
		UpdatedAt:      timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:      ptrFromUUID(row.UpdatedBy),
	}
}

func mapAIHumanOverrides(rows []sqlc.HrmsAiHumanOverride) []*domain.AIHumanOverride {
	items := make([]*domain.AIHumanOverride, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAIHumanOverride(row))
	}
	return items
}

func mapAIEventOutbox(row sqlc.HrmsAiEventOutbox) *domain.AIEventOutbox {
	return &domain.AIEventOutbox{
		ID:            row.ID,
		TenantID:      row.TenantID,
		EventKey:      row.EventKey,
		EventType:     row.EventType,
		TargetBus:     row.TargetBus,
		Status:        row.Status,
		Payload:       row.Payload,
		CorrelationID: ptrFromText(row.CorrelationID),
		Attempts:      row.Attempts,
		PublishedAt:   ptrFromTimestamptz(row.PublishedAt),
		LastError:     ptrFromText(row.LastError),
		Inactive:      row.Inactive,
		CreatedAt:     timeFromTimestamptz(row.CreatedAt),
		CreatedBy:     ptrFromUUID(row.CreatedBy),
		UpdatedAt:     timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:     ptrFromUUID(row.UpdatedBy),
	}
}

func mapAIEventOutboxRows(rows []sqlc.HrmsAiEventOutbox) []*domain.AIEventOutbox {
	items := make([]*domain.AIEventOutbox, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAIEventOutbox(row))
	}
	return items
}
