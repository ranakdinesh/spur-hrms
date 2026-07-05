package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) UpsertAISignalLog(ctx context.Context, item *domain.AISignalLog, actorID *uuid.UUID) (*domain.AISignalLog, error) {
	row, err := s.getQueries(ctx).UpsertAISignalLog(ctx, sqlc.UpsertAISignalLogParams{
		TenantID:         item.TenantID,
		SignalKey:        item.SignalKey,
		SignalType:       item.SignalType,
		SourceModule:     item.SourceModule,
		SourceEvent:      item.SourceEvent,
		Severity:         item.Severity,
		ProcessingStatus: item.ProcessingStatus,
		VisibilityScope:  item.VisibilityScope,
		Payload:          jsonRawOrDefault(item.Payload, "{}"),
		Explainability:   jsonRawOrDefault(item.Explainability, "{}"),
		OccurredAt:       timestamptzFromPtr(&item.OccurredAt),
		EntityType:       textFromPtr(item.EntityType),
		EntityID:         uuidFromPtr(item.EntityID),
		EmployeeUserID:   uuidFromPtr(item.EmployeeUserID),
		IdempotencyKey:   textFromPtr(item.IdempotencyKey),
		CorrelationID:    textFromPtr(item.CorrelationID),
		ActorID:          uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert ai signal", err, tenantIDField(item.TenantID), stringField("signal_key", item.SignalKey))
	}
	return mapAISignalLog(row), nil
}

func (s *Store) ListAISignalLogs(ctx context.Context, filter domain.AIActionFilter) ([]*domain.AISignalLog, error) {
	rows, err := s.getQueries(ctx).ListAISignalLogs(ctx, sqlc.ListAISignalLogsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, ProcessingStatus: textFromPtr(filter.ProcessingStatus), Severity: textFromPtr(filter.Severity), SourceModule: textFromPtr(filter.SourceModule), VisibilityScope: textFromPtr(filter.VisibilityScope)})
	if err != nil {
		return nil, s.logDBError(ctx, "list ai signals", err, tenantIDField(filter.TenantID))
	}
	return mapAISignalLogs(rows), nil
}

func (s *Store) UpdateAISignalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, errorMessage *string, actorID *uuid.UUID) (*domain.AISignalLog, error) {
	row, err := s.getQueries(ctx).UpdateAISignalStatus(ctx, sqlc.UpdateAISignalStatusParams{TenantID: tenantID, ID: id, ProcessingStatus: status, ErrorMessage: textFromPtr(errorMessage), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update ai signal status", err, tenantIDField(tenantID), stringField("signal_id", id.String()))
	}
	return mapAISignalLog(row), nil
}

func (s *Store) UpsertAIAgentActionLog(ctx context.Context, item *domain.AIAgentActionLog, actorID *uuid.UUID) (*domain.AIAgentActionLog, error) {
	row, err := s.getQueries(ctx).UpsertAIAgentActionLog(ctx, sqlc.UpsertAIAgentActionLogParams{
		TenantID:            item.TenantID,
		ActionKey:           item.ActionKey,
		AgentKey:            item.AgentKey,
		AgentName:           item.AgentName,
		ActionType:          item.ActionType,
		Status:              item.Status,
		Severity:            item.Severity,
		Title:               item.Title,
		Summary:             item.Summary,
		VisibilityScope:     item.VisibilityScope,
		ProposedAction:      jsonRawOrDefault(item.ProposedAction, "{}"),
		InputSnapshot:       jsonRawOrDefault(item.InputSnapshot, "{}"),
		OutputSnapshot:      jsonRawOrDefault(item.OutputSnapshot, "{}"),
		Explainability:      jsonRawOrDefault(item.Explainability, "{}"),
		ConfidenceScore:     numericFromFloat(item.ConfidenceScore),
		RequiresHumanReview: item.RequiresHumanReview,
		InsightID:           uuidFromPtr(item.InsightID),
		SignalID:            uuidFromPtr(item.SignalID),
		EntityType:          textFromPtr(item.EntityType),
		EntityID:            uuidFromPtr(item.EntityID),
		EmployeeUserID:      uuidFromPtr(item.EmployeeUserID),
		ModelVersion:        textFromPtr(item.ModelVersion),
		SidecarRunID:        textFromPtr(item.SidecarRunID),
		ActorID:             uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert ai action", err, tenantIDField(item.TenantID), stringField("action_key", item.ActionKey))
	}
	return mapAIAgentActionLog(row), nil
}

func (s *Store) ListAIAgentActionLogs(ctx context.Context, filter domain.AIActionFilter) ([]*domain.AIAgentActionLog, error) {
	rows, err := s.getQueries(ctx).ListAIAgentActionLogs(ctx, sqlc.ListAIAgentActionLogsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, Status: textFromPtr(filter.Status), Severity: textFromPtr(filter.Severity), AgentKey: textFromPtr(filter.AgentKey), VisibilityScope: textFromPtr(filter.VisibilityScope), InsightID: uuidFromPtr(filter.InsightID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list ai actions", err, tenantIDField(filter.TenantID))
	}
	return mapAIAgentActionLogs(rows), nil
}

func (s *Store) GetAIAgentActionLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AIAgentActionLog, error) {
	row, err := s.getQueries(ctx).GetAIAgentActionLog(ctx, sqlc.GetAIAgentActionLogParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAIActionNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get ai action", err, tenantIDField(tenantID), stringField("action_id", id.String()))
	}
	return mapAIAgentActionLog(row), nil
}

func (s *Store) UpdateAIAgentActionStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, failureMessage *string, actorID *uuid.UUID) (*domain.AIAgentActionLog, error) {
	row, err := s.getQueries(ctx).UpdateAIAgentActionStatus(ctx, sqlc.UpdateAIAgentActionStatusParams{TenantID: tenantID, ID: id, Status: status, ReviewedBy: uuidFromPtr(actorID), FailureMessage: textFromPtr(failureMessage), ActorID: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAIActionNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update ai action status", err, tenantIDField(tenantID), stringField("action_id", id.String()))
	}
	return mapAIAgentActionLog(row), nil
}

func (s *Store) CreateAIHumanOverride(ctx context.Context, item *domain.AIHumanOverride, actorID *uuid.UUID) (*domain.AIHumanOverride, error) {
	row, err := s.getQueries(ctx).CreateAIHumanOverride(ctx, sqlc.CreateAIHumanOverrideParams{TenantID: item.TenantID, OverrideType: item.OverrideType, OverrideStatus: item.OverrideStatus, Reason: item.Reason, Decision: item.Decision, Metadata: jsonRawOrDefault(item.Metadata, "{}"), InsightID: uuidFromPtr(item.InsightID), ActionID: uuidFromPtr(item.ActionID), OriginalStatus: textFromPtr(item.OriginalStatus), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create ai override", err, tenantIDField(item.TenantID))
	}
	return mapAIHumanOverride(row), nil
}

func (s *Store) ListAIHumanOverrides(ctx context.Context, filter domain.AIActionFilter) ([]*domain.AIHumanOverride, error) {
	rows, err := s.getQueries(ctx).ListAIHumanOverrides(ctx, sqlc.ListAIHumanOverridesParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, InsightID: uuidFromPtr(filter.InsightID), ActionID: uuidFromPtr(filter.ActionID), Decision: textFromPtr(filter.Decision)})
	if err != nil {
		return nil, s.logDBError(ctx, "list ai overrides", err, tenantIDField(filter.TenantID))
	}
	return mapAIHumanOverrides(rows), nil
}

func (s *Store) UpsertAIEventOutbox(ctx context.Context, item *domain.AIEventOutbox, actorID *uuid.UUID) (*domain.AIEventOutbox, error) {
	row, err := s.getQueries(ctx).UpsertAIEventOutbox(ctx, sqlc.UpsertAIEventOutboxParams{TenantID: item.TenantID, EventKey: item.EventKey, EventType: item.EventType, TargetBus: item.TargetBus, Status: item.Status, Payload: jsonRawOrDefault(item.Payload, "{}"), CorrelationID: textFromPtr(item.CorrelationID), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert ai event outbox", err, tenantIDField(item.TenantID), stringField("event_key", item.EventKey))
	}
	return mapAIEventOutbox(row), nil
}

func (s *Store) ListAIEventOutbox(ctx context.Context, filter domain.AIActionFilter) ([]*domain.AIEventOutbox, error) {
	rows, err := s.getQueries(ctx).ListAIEventOutbox(ctx, sqlc.ListAIEventOutboxParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, Status: textFromPtr(filter.Status), EventType: textFromPtr(filter.EventType)})
	if err != nil {
		return nil, s.logDBError(ctx, "list ai event outbox", err, tenantIDField(filter.TenantID))
	}
	return mapAIEventOutboxRows(rows), nil
}

func (s *Store) UpdateAIEventOutboxStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, lastError *string, actorID *uuid.UUID) (*domain.AIEventOutbox, error) {
	row, err := s.getQueries(ctx).UpdateAIEventOutboxStatus(ctx, sqlc.UpdateAIEventOutboxStatusParams{TenantID: tenantID, ID: id, Status: status, LastError: textFromPtr(lastError), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update ai event outbox status", err, tenantIDField(tenantID), stringField("event_id", id.String()))
	}
	return mapAIEventOutbox(row), nil
}

func limitOrDefault(limit int32) int32 {
	if limit <= 0 {
		return 100
	}
	return limit
}
