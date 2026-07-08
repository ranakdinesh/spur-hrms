package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateCompOffRequest(ctx context.Context, item *domain.CompOffRequest, actorID *uuid.UUID) (*domain.CompOffRequest, error) {
	row, err := s.getQueries(ctx).CreateCompOffRequest(ctx, sqlc.CreateCompOffRequestParams{
		TenantID:           item.TenantID,
		UserID:             item.UserID,
		LeaveTypeID:        item.LeaveTypeID,
		FyID:               item.FYID,
		WorkDate:           dateFromTime(item.WorkDate),
		WorkedMinutes:      item.WorkedMinutes,
		RequestedDays:      numericFromFloat(item.RequestedDays),
		ExpiryDate:         dateFromPtr(item.ExpiryDate),
		Reason:             textFromPtr(item.Reason),
		PayrollImpact:      item.PayrollImpact,
		SourceAttendanceID: uuidFromPtr(item.SourceAttendanceID),
		SourceSegmentID:    uuidFromPtr(item.SourceSegmentID),
		Metadata:           jsonBytesFromMap(item.Metadata),
		CreatedBy:          uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create comp-off request", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapCompOffRequest(row), nil
}

func (s *Store) GetCompOffRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompOffRequest, error) {
	row, err := s.getQueries(ctx).GetCompOffRequest(ctx, sqlc.GetCompOffRequestParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCompOffRequestNotFound
		}
		return nil, s.logDBError(ctx, "get comp-off request", err, tenantIDField(tenantID), stringField("comp_off_request_id", id.String()))
	}
	return mapCompOffRequest(row), nil
}

func (s *Store) ListCompOffRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.CompOffRequest, error) {
	rows, err := s.getQueries(ctx).ListCompOffRequestsByUser(ctx, sqlc.ListCompOffRequestsByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list comp-off requests by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapCompOffRequests(rows), nil
}

func (s *Store) ListCompOffRequestsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.CompOffRequest, error) {
	rows, err := s.getQueries(ctx).ListCompOffRequestsByStatus(ctx, sqlc.ListCompOffRequestsByStatusParams{TenantID: tenantID, Status: status})
	if err != nil {
		return nil, s.logDBError(ctx, "list comp-off requests by status", err, tenantIDField(tenantID), stringField("status", status))
	}
	return mapCompOffRequests(rows), nil
}

func (s *Store) ReviewCompOffRequest(ctx context.Context, item *domain.CompOffRequest, actorID *uuid.UUID) (*domain.CompOffRequest, error) {
	row, err := s.getQueries(ctx).ReviewCompOffRequest(ctx, sqlc.ReviewCompOffRequestParams{
		TenantID:      item.TenantID,
		ID:            item.ID,
		Status:        item.Status,
		ApprovedDays:  numericFromFloatPtr(item.ApprovedDays),
		ExpiryDate:    dateFromPtr(item.ExpiryDate),
		PayrollImpact: item.PayrollImpact,
		ReviewedBy:    uuidFromPtr(actorID),
		ReviewRemarks: textFromPtr(item.ReviewRemarks),
		Metadata:      jsonBytesFromMap(item.Metadata),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCompOffRequestNotPending
		}
		return nil, s.logDBError(ctx, "review comp-off request", err, tenantIDField(item.TenantID), stringField("comp_off_request_id", item.ID.String()))
	}
	return mapCompOffRequest(row), nil
}
