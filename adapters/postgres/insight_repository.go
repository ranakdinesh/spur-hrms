package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *Store) UpsertInsight(ctx context.Context, item *domain.Insight, actorID *uuid.UUID) (*domain.Insight, error) {
	row, err := s.getQueries(ctx).UpsertInsight(ctx, sqlc.UpsertInsightParams{
		TenantID:        item.TenantID,
		InsightKey:      item.InsightKey,
		InsightType:     item.InsightType,
		Category:        item.Category,
		Severity:        item.Severity,
		Status:          item.Status,
		Title:           item.Title,
		Summary:         item.Summary,
		ConfidenceScore: numericFromFloat(item.ConfidenceScore),
		Score:           numericFromFloat(item.Score),
		Source:          item.Source,
		Reasons:         jsonRawOrDefault(item.Reasons, "[]"),
		Recommendations: jsonRawOrDefault(item.Recommendations, "[]"),
		Context:         jsonRawOrDefault(item.Context, "{}"),
		Explainability:  jsonRawOrDefault(item.Explainability, "{}"),
		DetectedAt:      timestamptzFromPtr(&item.DetectedAt),
		ID:              uuidFromPtr(nonNilUUID(item.ID)),
		ModelVersion:    textFromPtr(item.ModelVersion),
		EntityType:      textFromPtr(item.EntityType),
		EntityID:        uuidFromPtr(item.EntityID),
		EmployeeUserID:  uuidFromPtr(item.EmployeeUserID),
		DueAt:           timestamptzFromPtr(item.DueAt),
		AssignedTo:      uuidFromPtr(item.AssignedTo),
		ReviewedBy:      uuidFromPtr(item.ReviewedBy),
		ReviewedAt:      timestamptzFromPtr(item.ReviewedAt),
		ResolvedAt:      timestamptzFromPtr(item.ResolvedAt),
		ResolutionNote:  textFromPtr(item.ResolutionNote),
		ActorID:         uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert insight", err, tenantIDField(item.TenantID), stringField("insight_key", item.InsightKey))
	}
	return mapInsight(row), nil
}

func (s *Store) ListInsights(ctx context.Context, filter domain.InsightFilter) ([]*domain.Insight, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.getQueries(ctx).ListInsights(ctx, sqlc.ListInsightsParams{TenantID: filter.TenantID, Limit: limit, Offset: filter.Offset, Status: textFromPtr(filter.Status), Severity: textFromPtr(filter.Severity), Category: textFromPtr(filter.Category), InsightType: textFromPtr(filter.InsightType), AssignedTo: uuidFromPtr(filter.AssignedTo)})
	if err != nil {
		return nil, s.logDBError(ctx, "list insights", err, tenantIDField(filter.TenantID))
	}
	return mapInsights(rows), nil
}

func (s *Store) GetInsight(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Insight, error) {
	row, err := s.getQueries(ctx).GetInsight(ctx, sqlc.GetInsightParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInsightNotFound
		}
		return nil, s.logDBError(ctx, "get insight", err, tenantIDField(tenantID), stringField("insight_id", id.String()))
	}
	return mapInsight(row), nil
}

func (s *Store) UpdateInsightStatus(ctx context.Context, cmd ports.InsightStatusCommand) (*domain.Insight, error) {
	row, err := s.getQueries(ctx).UpdateInsightStatus(ctx, sqlc.UpdateInsightStatusParams{TenantID: cmd.TenantID, ID: cmd.ID, Status: cmd.Status, AssignedTo: uuidFromPtr(cmd.AssignedTo), ReviewedBy: uuidFromPtr(cmd.ActorID), ResolutionNote: textFromPtr(cmd.ResolutionNote), ActorID: uuidFromPtr(cmd.ActorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update insight status", err, tenantIDField(cmd.TenantID), stringField("insight_id", cmd.ID.String()))
	}
	return mapInsight(row), nil
}

func (s *Store) CreateInsightEvent(ctx context.Context, event *domain.InsightEvent, actorID *uuid.UUID) (*domain.InsightEvent, error) {
	row, err := s.getQueries(ctx).CreateInsightEvent(ctx, sqlc.CreateInsightEventParams{TenantID: event.TenantID, InsightID: event.InsightID, Action: event.Action, Metadata: jsonRawOrDefault(event.Metadata, "{}"), FromStatus: textFromPtr(event.FromStatus), ToStatus: textFromPtr(event.ToStatus), Remarks: textFromPtr(event.Remarks), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create insight event", err, tenantIDField(event.TenantID), stringField("insight_id", event.InsightID.String()))
	}
	return mapInsightEvent(row), nil
}

func (s *Store) ListInsightEvents(ctx context.Context, tenantID uuid.UUID, insightID uuid.UUID) ([]*domain.InsightEvent, error) {
	rows, err := s.getQueries(ctx).ListInsightEvents(ctx, sqlc.ListInsightEventsParams{TenantID: tenantID, InsightID: insightID})
	if err != nil {
		return nil, s.logDBError(ctx, "list insight events", err, tenantIDField(tenantID), stringField("insight_id", insightID.String()))
	}
	return mapInsightEvents(rows), nil
}

func nonNilUUID(value uuid.UUID) *uuid.UUID {
	if value == uuid.Nil {
		return nil
	}
	return &value
}

func jsonRawOrDefault(value []byte, fallback string) []byte {
	if len(value) == 0 {
		return []byte(fallback)
	}
	return value
}
