package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateLeave(ctx context.Context, item *domain.Leave, actorID *uuid.UUID) (*domain.Leave, error) {
	row, err := s.getQueries(ctx).CreateLeave(ctx, sqlc.CreateLeaveParams{
		TenantID:      item.TenantID,
		UserID:        item.UserID,
		LeaveTypeID:   item.LeaveTypeID,
		FyID:          item.FYID,
		StartDate:     dateFromTime(item.StartDate),
		EndDate:       dateFromTime(item.EndDate),
		StartDayType:  item.StartDayType,
		EndDayType:    item.EndDayType,
		Days:          numericFromFloat(item.Days),
		Reason:        textFromPtr(item.Reason),
		Status:        item.Status,
		FromLeaveType: uuidFromPtr(item.FromLeaveType),
		ToLeaveType:   uuidFromPtr(item.ToLeaveType),
		IsSandwich:    item.IsSandwich,
		CreatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create leave", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapLeave(row), nil
}

func (s *Store) CreateLeaveApproval(ctx context.Context, item *domain.LeaveApproval, actorID *uuid.UUID) (*domain.LeaveApproval, error) {
	row, err := s.getQueries(ctx).CreateLeaveApproval(ctx, sqlc.CreateLeaveApprovalParams{TenantID: item.TenantID, LeaveID: item.LeaveID, ApproverID: item.ApproverID, Status: item.Status, Remarks: textFromPtr(item.Remarks), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create leave approval", err, tenantIDField(item.TenantID), stringField("leave_id", item.LeaveID.String()))
	}
	return mapLeaveApproval(row), nil
}

func (s *Store) ListLeavesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.Leave, error) {
	rows, err := s.getQueries(ctx).ListLeavesByUser(ctx, sqlc.ListLeavesByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leaves by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapLeaves(rows), nil
}

func (s *Store) ListLeavesByFY(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) ([]*domain.Leave, error) {
	rows, err := s.getQueries(ctx).ListLeavesByFY(ctx, sqlc.ListLeavesByFYParams{TenantID: tenantID, FyID: fyID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leaves by fy", err, tenantIDField(tenantID), stringField("financial_year_id", fyID.String()))
	}
	return mapLeaves(rows), nil
}

func (s *Store) ListLeaveReportRows(ctx context.Context, filter domain.LeaveReportFilter) ([]*domain.LeaveReportRow, error) {
	if filter.ManagerID != nil && *filter.ManagerID != uuid.Nil {
		rows, err := s.getQueries(ctx).ListManagerLeaveReportRows(ctx, sqlc.ListManagerLeaveReportRowsParams{
			TenantID:     filter.TenantID,
			ManagerID:    uuidFromPtr(filter.ManagerID),
			FyID:         uuidFromPtr(filter.FYID),
			UserID:       uuidFromPtr(filter.UserID),
			DepartmentID: uuidFromPtr(filter.DepartmentID),
			LeaveTypeID:  uuidFromPtr(filter.LeaveTypeID),
			Status:       textFromPtr(filter.Status),
			StartDate:    dateFromPtr(filter.StartDate),
			EndDate:      dateFromPtr(filter.EndDate),
		})
		if err != nil {
			return nil, s.logDBError(ctx, "list manager leave report rows", err, tenantIDField(filter.TenantID), stringField("manager_id", filter.ManagerID.String()))
		}
		return mapManagerLeaveReportRows(rows), nil
	}
	rows, err := s.getQueries(ctx).ListLeaveReportRows(ctx, sqlc.ListLeaveReportRowsParams{
		TenantID:     filter.TenantID,
		FyID:         uuidFromPtr(filter.FYID),
		UserID:       uuidFromPtr(filter.UserID),
		DepartmentID: uuidFromPtr(filter.DepartmentID),
		LeaveTypeID:  uuidFromPtr(filter.LeaveTypeID),
		Status:       textFromPtr(filter.Status),
		StartDate:    dateFromPtr(filter.StartDate),
		EndDate:      dateFromPtr(filter.EndDate),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave report rows", err, tenantIDField(filter.TenantID))
	}
	return mapLeaveReportRows(rows), nil
}

func (s *Store) GetLeaveReportSummary(ctx context.Context, filter domain.LeaveReportFilter) (*domain.LeaveReportSummary, error) {
	row, err := s.getQueries(ctx).GetLeaveReportSummary(ctx, sqlc.GetLeaveReportSummaryParams{
		TenantID:     filter.TenantID,
		ManagerID:    uuidFromPtr(filter.ManagerID),
		FyID:         uuidFromPtr(filter.FYID),
		UserID:       uuidFromPtr(filter.UserID),
		DepartmentID: uuidFromPtr(filter.DepartmentID),
		LeaveTypeID:  uuidFromPtr(filter.LeaveTypeID),
		Status:       textFromPtr(filter.Status),
		StartDate:    dateFromPtr(filter.StartDate),
		EndDate:      dateFromPtr(filter.EndDate),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "get leave report summary", err, tenantIDField(filter.TenantID))
	}
	return mapLeaveReportSummary(row), nil
}

func (s *Store) ListOverlappingLeaves(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.Leave, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListOverlappingLeaves(ctx, sqlc.ListOverlappingLeavesParams{TenantID: tenantID, UserID: userID, EndDate: dateFromTime(start), StartDate: dateFromTime(end)})
	if err != nil {
		return nil, s.logDBError(ctx, "list overlapping leaves", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapLeaves(rows), nil
}

func (s *Store) GetLeave(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Leave, error) {
	row, err := s.getQueries(ctx).GetLeave(ctx, sqlc.GetLeaveParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveNotFound
		}
		return nil, s.logDBError(ctx, "get leave", err, tenantIDField(tenantID), stringField("leave_id", id.String()))
	}
	return mapLeave(row), nil
}
