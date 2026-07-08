package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateOvertimeRequest(ctx context.Context, item *domain.OvertimeRequest, actorID *uuid.UUID) (*domain.OvertimeRequest, error) {
	row, err := s.getQueries(ctx).CreateOvertimeRequest(ctx, sqlc.CreateOvertimeRequestParams{
		TenantID:             item.TenantID,
		UserID:               item.UserID,
		WorkDate:             dateFromTime(item.WorkDate),
		RequestedMinutes:     item.RequestedMinutes,
		Reason:               textFromPtr(item.Reason),
		CalculationType:      item.CalculationType,
		RateMultiplier:       numericFromFloat(item.RateMultiplier),
		PayrollComponentCode: textFromPtr(item.PayrollComponentCode),
		SourceAttendanceID:   uuidFromPtr(item.SourceAttendanceID),
		SourceSegmentID:      uuidFromPtr(item.SourceSegmentID),
		Metadata:             jsonBytesFromMap(item.Metadata),
		CreatedBy:            uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create overtime request", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapOvertimeRequest(row), nil
}

func (s *Store) GetOvertimeRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OvertimeRequest, error) {
	row, err := s.getQueries(ctx).GetOvertimeRequest(ctx, sqlc.GetOvertimeRequestParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOvertimeRequestNotFound
		}
		return nil, s.logDBError(ctx, "get overtime request", err, tenantIDField(tenantID), stringField("overtime_request_id", id.String()))
	}
	return mapOvertimeRequest(row), nil
}

func (s *Store) ListOvertimeRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.OvertimeRequest, error) {
	rows, err := s.getQueries(ctx).ListOvertimeRequestsByUser(ctx, sqlc.ListOvertimeRequestsByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list overtime requests by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapOvertimeRequests(rows), nil
}

func (s *Store) ListOvertimeRequestsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.OvertimeRequest, error) {
	rows, err := s.getQueries(ctx).ListOvertimeRequestsByStatus(ctx, sqlc.ListOvertimeRequestsByStatusParams{TenantID: tenantID, Status: status})
	if err != nil {
		return nil, s.logDBError(ctx, "list overtime requests by status", err, tenantIDField(tenantID), stringField("status", status))
	}
	return mapOvertimeRequests(rows), nil
}

func (s *Store) ListOvertimeRequestsByPayrollExportStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.OvertimeRequest, error) {
	rows, err := s.getQueries(ctx).ListOvertimeRequestsByPayrollExportStatus(ctx, sqlc.ListOvertimeRequestsByPayrollExportStatusParams{TenantID: tenantID, PayrollExportStatus: status})
	if err != nil {
		return nil, s.logDBError(ctx, "list overtime requests by payroll export status", err, tenantIDField(tenantID), stringField("payroll_export_status", status))
	}
	return mapOvertimeRequests(rows), nil
}

func (s *Store) ReviewOvertimeRequest(ctx context.Context, item *domain.OvertimeRequest, actorID *uuid.UUID) (*domain.OvertimeRequest, error) {
	row, err := s.getQueries(ctx).ReviewOvertimeRequest(ctx, sqlc.ReviewOvertimeRequestParams{
		TenantID:             item.TenantID,
		ID:                   item.ID,
		Status:               item.Status,
		ApprovedMinutes:      int4FromInt32Ptr(item.ApprovedMinutes),
		ReviewRemarks:        textFromPtr(item.ReviewRemarks),
		CalculationType:      item.CalculationType,
		RateMultiplier:       numericFromFloat(item.RateMultiplier),
		PayrollComponentCode: textFromPtr(item.PayrollComponentCode),
		PayrollExportStatus:  item.PayrollExportStatus,
		ReviewedBy:           uuidFromPtr(actorID),
		Metadata:             jsonBytesFromMap(item.Metadata),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOvertimeRequestNotPending
		}
		return nil, s.logDBError(ctx, "review overtime request", err, tenantIDField(item.TenantID), stringField("overtime_request_id", item.ID.String()))
	}
	return mapOvertimeRequest(row), nil
}

func int4FromInt32Ptr(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *value, Valid: true}
}

func int32PtrFromInt4(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	return &value.Int32
}
