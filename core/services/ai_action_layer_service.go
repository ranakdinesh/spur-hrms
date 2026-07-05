package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) ListAIActionWorkspace(ctx context.Context, filter domain.AIActionFilter) (*domain.AIActionWorkspace, error) {
	filter.Limit = limitAI(filter.Limit)
	signals, err := s.aiActions.ListAISignalLogs(ctx, filter)
	if err != nil {
		s.logError("list ai signals", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	actions, err := s.aiActions.ListAIAgentActionLogs(ctx, filter)
	if err != nil {
		s.logError("list ai actions", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	overrides, err := s.aiActions.ListAIHumanOverrides(ctx, filter)
	if err != nil {
		s.logError("list ai overrides", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	events, err := s.aiActions.ListAIEventOutbox(ctx, filter)
	if err != nil {
		s.logError("list ai events", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return &domain.AIActionWorkspace{Signals: signals, Actions: actions, Overrides: overrides, Events: events, Summary: aiActionSummary(signals, actions, overrides, events)}, nil
}

func (s *TenantService) CreateAISignal(ctx context.Context, cmd ports.AISignalCommand) (*domain.AISignalLog, error) {
	occurredAt, err := parseAITime(cmd.OccurredAt)
	if err != nil {
		s.logError("parse ai signal time", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewAISignalLog(domain.AISignalLog{TenantID: cmd.TenantID, SignalKey: cmd.SignalKey, SignalType: cmd.SignalType, SourceModule: cmd.SourceModule, SourceEvent: cmd.SourceEvent, Severity: cmd.Severity, ProcessingStatus: cmd.ProcessingStatus, EntityType: cmd.EntityType, EntityID: cmd.EntityID, EmployeeUserID: cmd.EmployeeUserID, VisibilityScope: cmd.VisibilityScope, IdempotencyKey: cmd.IdempotencyKey, CorrelationID: cmd.CorrelationID, Payload: cmd.Payload, Explainability: cmd.Explainability, OccurredAt: occurredAt})
	if err != nil {
		s.logError("validate ai signal", err, serviceTenantIDField(cmd.TenantID), serviceStringField("signal_key", cmd.SignalKey))
		return nil, err
	}
	return s.aiActions.UpsertAISignalLog(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateAISignalStatus(ctx context.Context, cmd ports.AIStatusCommand) (*domain.AISignalLog, error) {
	status, err := domain.NormalizeAISignalStatus(cmd.Status)
	if err != nil {
		s.logError("validate ai signal status", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.aiActions.UpdateAISignalStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.FailureMessage, cmd.ActorID)
}

func (s *TenantService) CreateAIAgentAction(ctx context.Context, cmd ports.AIActionCommand) (*domain.AIAgentActionLog, error) {
	item, err := domain.NewAIAgentActionLog(domain.AIAgentActionLog{TenantID: cmd.TenantID, ActionKey: cmd.ActionKey, AgentKey: cmd.AgentKey, AgentName: cmd.AgentName, ActionType: cmd.ActionType, Status: cmd.Status, Severity: cmd.Severity, Title: cmd.Title, Summary: cmd.Summary, InsightID: cmd.InsightID, SignalID: cmd.SignalID, EntityType: cmd.EntityType, EntityID: cmd.EntityID, EmployeeUserID: cmd.EmployeeUserID, VisibilityScope: cmd.VisibilityScope, ProposedAction: cmd.ProposedAction, InputSnapshot: cmd.InputSnapshot, OutputSnapshot: cmd.OutputSnapshot, Explainability: cmd.Explainability, ConfidenceScore: cmd.ConfidenceScore, ModelVersion: cmd.ModelVersion, SidecarRunID: cmd.SidecarRunID, RequiresHumanReview: cmd.RequiresHumanReview})
	if err != nil {
		s.logError("validate ai action", err, serviceTenantIDField(cmd.TenantID), serviceStringField("action_key", cmd.ActionKey))
		return nil, err
	}
	return s.aiActions.UpsertAIAgentActionLog(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateAIAgentActionStatus(ctx context.Context, cmd ports.AIStatusCommand) (*domain.AIAgentActionLog, error) {
	status, err := domain.NormalizeAIActionStatus(cmd.Status)
	if err != nil {
		s.logError("validate ai action status", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.aiActions.UpdateAIAgentActionStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.FailureMessage, cmd.ActorID)
}

func (s *TenantService) OverrideAIAction(ctx context.Context, cmd ports.AIOverrideCommand) (*domain.AIHumanOverride, error) {
	item, err := domain.NewAIHumanOverride(domain.AIHumanOverride{TenantID: cmd.TenantID, InsightID: cmd.InsightID, ActionID: cmd.ActionID, OverrideType: cmd.OverrideType, OriginalStatus: cmd.OriginalStatus, OverrideStatus: cmd.OverrideStatus, Reason: cmd.Reason, Decision: cmd.Decision, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate ai override", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	override, err := s.aiActions.CreateAIHumanOverride(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create ai override", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.ActionID != nil {
		if _, err := s.aiActions.UpdateAIAgentActionStatus(ctx, cmd.TenantID, *cmd.ActionID, actionStatusForOverride(item.OverrideStatus), nil, cmd.ActorID); err != nil {
			s.log.Warn().Err(err).Str("tenant_id", cmd.TenantID.String()).Str("action_id", cmd.ActionID.String()).Msg("hrms: ai override recorded but action status update failed")
		}
	}
	return override, nil
}

func (s *TenantService) EmitAIWorkflowEvent(ctx context.Context, cmd ports.AIWorkflowEventCommand) (*domain.AISignalLog, error) {
	if strings.TrimSpace(cmd.SignalType) == "" {
		cmd.SignalType = cmd.EventType
	}
	signal, err := s.CreateAISignal(ctx, ports.AISignalCommand{TenantID: cmd.TenantID, SignalKey: "signal-" + cmd.EventKey, SignalType: cmd.SignalType, SourceModule: cmd.SourceModule, SourceEvent: cmd.SourceEvent, Severity: cmd.Severity, ProcessingStatus: domain.AISignalStatusQueued, EntityType: cmd.EntityType, EntityID: cmd.EntityID, EmployeeUserID: cmd.EmployeeUserID, VisibilityScope: cmd.VisibilityScope, IdempotencyKey: &cmd.EventKey, CorrelationID: cmd.CorrelationID, Payload: cmd.Payload, Explainability: cmd.Explainability, ActorID: cmd.ActorID})
	if err != nil {
		return nil, err
	}
	outbox, err := domain.NewAIEventOutbox(domain.AIEventOutbox{TenantID: cmd.TenantID, EventKey: cmd.EventKey, EventType: cmd.EventType, TargetBus: domain.AIEventTargetRedisStream, Status: domain.AIEventStatusPending, Payload: cmd.Payload, CorrelationID: cmd.CorrelationID})
	if err != nil {
		return nil, err
	}
	stored, err := s.aiActions.UpsertAIEventOutbox(ctx, outbox, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	if s.aiEventPublisher != nil {
		if err := s.aiEventPublisher.PublishAIEvent(ctx, stored); err != nil {
			message := err.Error()
			_, _ = s.aiActions.UpdateAIEventOutboxStatus(ctx, stored.TenantID, stored.ID, domain.AIEventStatusFailed, &message, cmd.ActorID)
			return signal, err
		}
		_, _ = s.aiActions.UpdateAIEventOutboxStatus(ctx, stored.TenantID, stored.ID, domain.AIEventStatusPublished, nil, cmd.ActorID)
	}
	return signal, nil
}

func (s *TenantService) recordInsightAIAction(ctx context.Context, insight *domain.Insight, actorID *uuid.UUID) {
	entityType := insight.EntityType
	signalKey := "insight-" + insight.InsightKey
	signal, err := domain.NewAISignalLog(domain.AISignalLog{TenantID: insight.TenantID, SignalKey: signalKey, SignalType: insight.InsightType, SourceModule: "insights", SourceEvent: "insight_detected", Severity: insight.Severity, ProcessingStatus: domain.AISignalStatusProcessed, EntityType: entityType, EntityID: insight.EntityID, EmployeeUserID: insight.EmployeeUserID, VisibilityScope: domain.AIVisibilityHR, IdempotencyKey: &signalKey, Payload: rawJSON(map[string]any{"insight_id": insight.ID, "title": insight.Title, "summary": insight.Summary, "score": insight.Score}), Explainability: insight.Explainability, OccurredAt: insight.DetectedAt})
	if err != nil {
		s.log.Warn().Err(err).Str("tenant_id", insight.TenantID.String()).Str("insight_key", insight.InsightKey).Msg("hrms: skipped ai signal for insight")
		return
	}
	storedSignal, err := s.aiActions.UpsertAISignalLog(ctx, signal, actorID)
	if err != nil {
		s.log.Warn().Err(err).Str("tenant_id", insight.TenantID.String()).Str("insight_key", insight.InsightKey).Msg("hrms: failed ai signal upsert for insight")
		return
	}
	actionKey := "review-" + insight.InsightKey
	action, err := domain.NewAIAgentActionLog(domain.AIAgentActionLog{TenantID: insight.TenantID, ActionKey: actionKey, AgentKey: "insight_review_agent", AgentName: "Insight Review Agent", ActionType: "human_review", Status: domain.AIActionStatusProposed, Severity: insight.Severity, Title: "Review: " + insight.Title, Summary: insight.Summary, InsightID: &insight.ID, SignalID: &storedSignal.ID, EntityType: entityType, EntityID: insight.EntityID, EmployeeUserID: insight.EmployeeUserID, VisibilityScope: domain.AIVisibilityHR, ProposedAction: insight.Recommendations, InputSnapshot: insight.Context, OutputSnapshot: rawJSON(map[string]any{"status": insight.Status, "confidence_score": insight.ConfidenceScore}), Explainability: insight.Explainability, ConfidenceScore: insight.ConfidenceScore, ModelVersion: insight.ModelVersion, RequiresHumanReview: true})
	if err != nil {
		s.log.Warn().Err(err).Str("tenant_id", insight.TenantID.String()).Str("insight_key", insight.InsightKey).Msg("hrms: skipped ai action proposal for insight")
		return
	}
	if _, err := s.aiActions.UpsertAIAgentActionLog(ctx, action, actorID); err != nil {
		s.log.Warn().Err(err).Str("tenant_id", insight.TenantID.String()).Str("insight_key", insight.InsightKey).Msg("hrms: failed ai action proposal upsert for insight")
	}
	eventKey := "insight-event-" + insight.InsightKey
	outbox, err := domain.NewAIEventOutbox(domain.AIEventOutbox{TenantID: insight.TenantID, EventKey: eventKey, EventType: "insight.detected", TargetBus: domain.AIEventTargetRedisStream, Status: domain.AIEventStatusPending, Payload: rawJSON(map[string]any{"insight_id": insight.ID, "signal_id": storedSignal.ID, "category": insight.Category, "severity": insight.Severity})})
	if err != nil {
		return
	}
	if _, err := s.aiActions.UpsertAIEventOutbox(ctx, outbox, actorID); err != nil {
		s.log.Warn().Err(err).Str("tenant_id", insight.TenantID.String()).Str("event_key", eventKey).Msg("hrms: failed ai outbox upsert")
	}
}

func aiActionSummary(signals []*domain.AISignalLog, actions []*domain.AIAgentActionLog, overrides []*domain.AIHumanOverride, events []*domain.AIEventOutbox) domain.AIActionSummary {
	summary := domain.AIActionSummary{Overrides: int32(len(overrides))}
	for _, signal := range signals {
		if signal.ProcessingStatus == domain.AISignalStatusNew || signal.ProcessingStatus == domain.AISignalStatusQueued {
			summary.SignalsPending++
		}
		if signal.ProcessingStatus == domain.AISignalStatusFailed {
			summary.SignalsFailed++
		}
	}
	for _, action := range actions {
		switch action.Status {
		case domain.AIActionStatusProposed:
			summary.ActionsProposed++
		case domain.AIActionStatusApproved:
			summary.ActionsApproved++
		case domain.AIActionStatusExecuted:
			summary.ActionsExecuted++
		case domain.AIActionStatusFailed:
			summary.ActionsFailed++
		}
	}
	for _, event := range events {
		if event.Status == domain.AIEventStatusPending {
			summary.OutboxPending++
		}
		if event.Status == domain.AIEventStatusFailed {
			summary.OutboxFailed++
		}
	}
	return summary
}

func actionStatusForOverride(status string) string {
	switch status {
	case domain.InsightStatusResolved:
		return domain.AIActionStatusApproved
	case domain.InsightStatusDismissed, domain.AIActionStatusRejected:
		return domain.AIActionStatusRejected
	case domain.AIActionStatusCancelled:
		return domain.AIActionStatusCancelled
	default:
		return domain.AIActionStatusOverridden
	}
}

func parseAITime(value string) (time.Time, error) {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return time.Now().UTC(), nil
	}
	if parsed, err := time.Parse(time.RFC3339, clean); err == nil {
		return parsed.UTC(), nil
	}
	parsed, err := time.Parse("2006-01-02", clean)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid ai timestamp: %w", err)
	}
	return parsed.UTC(), nil
}

func limitAI(limit int32) int32 {
	if limit <= 0 {
		return 100
	}
	if limit > 500 {
		return 500
	}
	return limit
}
