package postgres

import (
	"context"
	"math"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateEngagement(ctx context.Context, item *domain.Engagement, actorID *uuid.UUID) (*domain.Engagement, error) {
	row, err := s.getQueries(ctx).CreateEngagement(ctx, createEngagementParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create engagement", err, tenantIDField(item.TenantID), stringField("worker_profile_id", item.WorkerProfileID.String()))
	}
	return mapEngagement(row), nil
}

func (s *Store) UpdateEngagement(ctx context.Context, item *domain.Engagement, actorID *uuid.UUID) (*domain.Engagement, error) {
	row, err := s.getQueries(ctx).UpdateEngagement(ctx, sqlc.UpdateEngagementParams{
		TenantID:           item.TenantID,
		ID:                 item.ID,
		WorkerProfileID:    item.WorkerProfileID,
		EngagementCode:     textFromPtr(item.EngagementCode),
		Title:              item.Title,
		Description:        textFromPtr(item.Description),
		EngagementType:     item.EngagementType,
		Status:             item.Status,
		StartDate:          dateFromTime(item.StartDate),
		EndDate:            dateFromPtr(item.EndDate),
		HoursBudget:        numericFromEngagementDecimalPtr(item.HoursBudget),
		RateAmount:         numericFromEngagementDecimalPtr(item.RateAmount),
		CurrencyCode:       item.CurrencyCode,
		RateUnit:           item.RateUnit,
		BranchID:           uuidFromPtr(item.BranchID),
		DepartmentID:       uuidFromPtr(item.DepartmentID),
		ReportingManagerID: uuidFromPtr(item.ReportingManagerID),
		ProjectLabel:       textFromPtr(item.ProjectLabel),
		ProjectCode:        textFromPtr(item.ProjectCode),
		CostCenter:         textFromPtr(item.CostCenter),
		RenewalDueDate:     dateFromPtr(item.RenewalDueDate),
		RenewalStatus:      item.RenewalStatus,
		TerminationReason:  textFromPtr(item.TerminationReason),
		TerminatedAt:       timestamptzFromPtr(item.TerminatedAt),
		Notes:              textFromPtr(item.Notes),
		Metadata:           jsonBytesFromRaw(item.Metadata),
		UpdatedBy:          uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update engagement", err, tenantIDField(item.TenantID), stringField("engagement_id", item.ID.String()))
	}
	return mapEngagement(row), nil
}

func (s *Store) UpdateEngagementStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, terminationReason *string, terminatedAt *time.Time, actorID *uuid.UUID) (*domain.Engagement, error) {
	row, err := s.getQueries(ctx).UpdateEngagementStatus(ctx, sqlc.UpdateEngagementStatusParams{
		TenantID:          tenantID,
		ID:                id,
		Status:            status,
		TerminationReason: textFromPtr(terminationReason),
		TerminatedAt:      timestamptzFromPtr(terminatedAt),
		UpdatedBy:         uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update engagement status", err, tenantIDField(tenantID), stringField("engagement_id", id.String()))
	}
	return mapEngagement(row), nil
}

func (s *Store) GetEngagement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Engagement, error) {
	row, err := s.getQueries(ctx).GetEngagement(ctx, sqlc.GetEngagementParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get engagement", err, tenantIDField(tenantID), stringField("engagement_id", id.String()))
	}
	return mapEngagement(row), nil
}

func (s *Store) ListEngagements(ctx context.Context, filter domain.EngagementFilter) ([]*domain.EngagementListItem, error) {
	rows, err := s.getQueries(ctx).ListEngagements(ctx, sqlc.ListEngagementsParams{
		TenantID:        filter.TenantID,
		WorkerProfileID: uuidFromPtr(filter.WorkerProfileID),
		EngagementType:  textFromPtr(filter.EngagementType),
		Status:          textFromPtr(filter.Status),
		DepartmentID:    uuidFromPtr(filter.DepartmentID),
		Search:          textFromPtr(filter.Search),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list engagements", err, tenantIDField(filter.TenantID))
	}
	return mapEngagementListItems(rows), nil
}

func (s *Store) DeleteEngagement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEngagement(ctx, sqlc.SoftDeleteEngagementParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete engagement", err, tenantIDField(tenantID), stringField("engagement_id", id.String()))
	}
	return nil
}

func createEngagementParams(item *domain.Engagement, actorID *uuid.UUID) sqlc.CreateEngagementParams {
	return sqlc.CreateEngagementParams{
		TenantID:           item.TenantID,
		WorkerProfileID:    item.WorkerProfileID,
		EngagementCode:     textFromPtr(item.EngagementCode),
		Title:              item.Title,
		Description:        textFromPtr(item.Description),
		EngagementType:     item.EngagementType,
		Status:             item.Status,
		StartDate:          dateFromTime(item.StartDate),
		EndDate:            dateFromPtr(item.EndDate),
		HoursBudget:        numericFromEngagementDecimalPtr(item.HoursBudget),
		RateAmount:         numericFromEngagementDecimalPtr(item.RateAmount),
		CurrencyCode:       item.CurrencyCode,
		RateUnit:           item.RateUnit,
		BranchID:           uuidFromPtr(item.BranchID),
		DepartmentID:       uuidFromPtr(item.DepartmentID),
		ReportingManagerID: uuidFromPtr(item.ReportingManagerID),
		ProjectLabel:       textFromPtr(item.ProjectLabel),
		ProjectCode:        textFromPtr(item.ProjectCode),
		CostCenter:         textFromPtr(item.CostCenter),
		RenewalDueDate:     dateFromPtr(item.RenewalDueDate),
		RenewalStatus:      item.RenewalStatus,
		TerminationReason:  textFromPtr(item.TerminationReason),
		TerminatedAt:       timestamptzFromPtr(item.TerminatedAt),
		Notes:              textFromPtr(item.Notes),
		Metadata:           jsonBytesFromRaw(item.Metadata),
		CreatedBy:          uuidFromPtr(actorID),
	}
}

func numericFromEngagementDecimalPtr(value *float64) pgtype.Numeric {
	if value == nil {
		return pgtype.Numeric{Valid: false}
	}
	scaled := int64(math.Round(*value * 100))
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -2, Valid: true}
}
