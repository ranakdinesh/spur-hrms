package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateWorkLog(ctx context.Context, item *domain.WorkLog, actorID *uuid.UUID) (*domain.WorkLog, error) {
	row, err := s.getQueries(ctx).CreateWorkLog(ctx, createWorkLogParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create work log", err, tenantIDField(item.TenantID), stringField("engagement_id", item.EngagementID.String()))
	}
	return mapWorkLog(row), nil
}

func (s *Store) UpdateWorkLog(ctx context.Context, item *domain.WorkLog, actorID *uuid.UUID) (*domain.WorkLog, error) {
	row, err := s.getQueries(ctx).UpdateWorkLog(ctx, sqlc.UpdateWorkLogParams{
		TenantID:             item.TenantID,
		ID:                   item.ID,
		LogDate:              dateFromTime(item.LogDate),
		HoursWorked:          numericFromEngagementDecimalPtr(&item.HoursWorked),
		BillableHours:        numericFromEngagementDecimalPtr(item.BillableHours),
		WorkSummary:          textFromPtr(item.WorkSummary),
		DeliverableReference: textFromPtr(item.DeliverableReference),
		Status:               item.Status,
		SubmittedAt:          timestamptzFromPtr(item.SubmittedAt),
		SubmittedBy:          uuidFromPtr(item.SubmittedBy),
		ReviewedAt:           timestamptzFromPtr(item.ReviewedAt),
		ReviewedBy:           uuidFromPtr(item.ReviewedBy),
		ReviewComment:        textFromPtr(item.ReviewComment),
		Metadata:             jsonBytesFromRaw(item.Metadata),
		UpdatedBy:            uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update work log", err, tenantIDField(item.TenantID), stringField("work_log_id", item.ID.String()))
	}
	return mapWorkLog(row), nil
}

func (s *Store) UpdateWorkLogStatus(ctx context.Context, item *domain.WorkLog, actorID *uuid.UUID) (*domain.WorkLog, error) {
	row, err := s.getQueries(ctx).UpdateWorkLogStatus(ctx, sqlc.UpdateWorkLogStatusParams{
		TenantID:      item.TenantID,
		ID:            item.ID,
		Status:        item.Status,
		SubmittedAt:   timestamptzFromPtr(item.SubmittedAt),
		SubmittedBy:   uuidFromPtr(item.SubmittedBy),
		ReviewedAt:    timestamptzFromPtr(item.ReviewedAt),
		ReviewedBy:    uuidFromPtr(item.ReviewedBy),
		ReviewComment: textFromPtr(item.ReviewComment),
		UpdatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update work log status", err, tenantIDField(item.TenantID), stringField("work_log_id", item.ID.String()))
	}
	return mapWorkLog(row), nil
}

func (s *Store) GetWorkLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkLog, error) {
	row, err := s.getQueries(ctx).GetWorkLog(ctx, sqlc.GetWorkLogParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get work log", err, tenantIDField(tenantID), stringField("work_log_id", id.String()))
	}
	return mapWorkLog(row), nil
}

func (s *Store) ListWorkLogs(ctx context.Context, filter domain.WorkLogFilter) ([]*domain.WorkLogListItem, error) {
	rows, err := s.getQueries(ctx).ListWorkLogs(ctx, sqlc.ListWorkLogsParams{
		TenantID:        filter.TenantID,
		EngagementID:    uuidFromPtr(filter.EngagementID),
		WorkerProfileID: uuidFromPtr(filter.WorkerProfileID),
		Status:          textFromPtr(filter.Status),
		DateFrom:        dateFromPtr(filter.DateFrom),
		DateTo:          dateFromPtr(filter.DateTo),
		Search:          textFromPtr(filter.Search),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list work logs", err, tenantIDField(filter.TenantID))
	}
	return mapWorkLogListItems(rows), nil
}

func (s *Store) ListWorkLogRollups(ctx context.Context, filter domain.WorkLogFilter) ([]*domain.WorkLogRollup, error) {
	rows, err := s.getQueries(ctx).ListWorkLogRollups(ctx, sqlc.ListWorkLogRollupsParams{
		TenantID:        filter.TenantID,
		EngagementID:    uuidFromPtr(filter.EngagementID),
		WorkerProfileID: uuidFromPtr(filter.WorkerProfileID),
		DateFrom:        dateFromPtr(filter.DateFrom),
		DateTo:          dateFromPtr(filter.DateTo),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list work log rollups", err, tenantIDField(filter.TenantID))
	}
	return mapWorkLogRollups(rows), nil
}

func (s *Store) GetWorkLogBudgetUsage(ctx context.Context, tenantID uuid.UUID, engagementID uuid.UUID, excludeWorkLogID *uuid.UUID) (*domain.WorkLogBudgetUsage, error) {
	row, err := s.getQueries(ctx).GetWorkLogBudgetUsage(ctx, sqlc.GetWorkLogBudgetUsageParams{
		TenantID:         tenantID,
		ID:               engagementID,
		ExcludeWorkLogID: uuidFromPtr(excludeWorkLogID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "get work log budget usage", err, tenantIDField(tenantID), stringField("engagement_id", engagementID.String()))
	}
	return &domain.WorkLogBudgetUsage{HoursBudget: ptrFromNumeric(row.HoursBudget), UsedHours: floatFromNumeric(row.UsedHours)}, nil
}

func (s *Store) DeleteWorkLog(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteWorkLog(ctx, sqlc.SoftDeleteWorkLogParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete work log", err, tenantIDField(tenantID), stringField("work_log_id", id.String()))
	}
	return nil
}

func createWorkLogParams(item *domain.WorkLog, actorID *uuid.UUID) sqlc.CreateWorkLogParams {
	return sqlc.CreateWorkLogParams{
		TenantID:             item.TenantID,
		EngagementID:         item.EngagementID,
		WorkerProfileID:      item.WorkerProfileID,
		LogDate:              dateFromTime(item.LogDate),
		HoursWorked:          numericFromEngagementDecimalPtr(&item.HoursWorked),
		BillableHours:        numericFromEngagementDecimalPtr(item.BillableHours),
		WorkSummary:          textFromPtr(item.WorkSummary),
		DeliverableReference: textFromPtr(item.DeliverableReference),
		Status:               item.Status,
		SubmittedAt:          timestamptzFromPtr(item.SubmittedAt),
		SubmittedBy:          uuidFromPtr(item.SubmittedBy),
		ReviewedAt:           timestamptzFromPtr(item.ReviewedAt),
		ReviewedBy:           uuidFromPtr(item.ReviewedBy),
		ReviewComment:        textFromPtr(item.ReviewComment),
		Metadata:             jsonBytesFromRaw(item.Metadata),
		CreatedBy:            uuidFromPtr(actorID),
	}
}
