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

func (s *Store) CreateHoliday(ctx context.Context, item *domain.Holiday, actorID *uuid.UUID) (*domain.Holiday, error) {
	row, err := s.getQueries(ctx).CreateHoliday(ctx, sqlc.CreateHolidayParams{
		TenantID:   item.TenantID,
		BranchID:   uuidFromPtr(item.BranchID),
		FyID:       uuidFromPtr(item.FYID),
		Name:       item.Name,
		Date:       dateFromTime(item.Date),
		IsOptional: item.IsOptional,
		CreatedBy:  uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create holiday", err, tenantIDField(item.TenantID), stringField("holiday_name", item.Name))
	}
	return mapHoliday(row), nil
}

func (s *Store) ListHolidays(ctx context.Context, tenantID uuid.UUID) ([]*domain.Holiday, error) {
	rows, err := s.getQueries(ctx).ListHolidays(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list holidays", err, tenantIDField(tenantID))
	}
	return mapHolidays(rows), nil
}

func (s *Store) ListHolidaysByDateRange(ctx context.Context, tenantID uuid.UUID, startDate time.Time, endDate time.Time) ([]*domain.Holiday, error) {
	rows, err := s.getQueries(ctx).ListHolidaysByDateRange(ctx, sqlc.ListHolidaysByDateRangeParams{TenantID: tenantID, Date: dateFromTime(startDate), Date_2: dateFromTime(endDate)})
	if err != nil {
		return nil, s.logDBError(ctx, "list holidays by date range", err, tenantIDField(tenantID), stringField("start_date", startDate.Format("2006-01-02")), stringField("end_date", endDate.Format("2006-01-02")))
	}
	return mapHolidays(rows), nil
}

func (s *Store) ListHolidaysByFinancialYear(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) ([]*domain.Holiday, error) {
	rows, err := s.getQueries(ctx).ListHolidaysByFinancialYear(ctx, sqlc.ListHolidaysByFinancialYearParams{TenantID: tenantID, FyID: uuidFromPtr(&fyID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list holidays by financial year", err, tenantIDField(tenantID), stringField("financial_year_id", fyID.String()))
	}
	return mapHolidays(rows), nil
}

func (s *Store) ListHolidaysByBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID) ([]*domain.Holiday, error) {
	rows, err := s.getQueries(ctx).ListHolidaysByBranch(ctx, sqlc.ListHolidaysByBranchParams{TenantID: tenantID, BranchID: uuidFromPtr(&branchID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list holidays by branch", err, tenantIDField(tenantID), stringField("branch_id", branchID.String()))
	}
	return mapHolidays(rows), nil
}

func (s *Store) ListUpcomingHolidays(ctx context.Context, tenantID uuid.UUID, limit int32) ([]*domain.Holiday, error) {
	rows, err := s.getQueries(ctx).ListUpcomingHolidays(ctx, sqlc.ListUpcomingHolidaysParams{TenantID: tenantID, Limit: limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list upcoming holidays", err, tenantIDField(tenantID))
	}
	return mapHolidays(rows), nil
}

func (s *Store) GetHoliday(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Holiday, error) {
	row, err := s.getQueries(ctx).GetHoliday(ctx, sqlc.GetHolidayParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrHolidayNotFound
		}
		return nil, s.logDBError(ctx, "get holiday", err, tenantIDField(tenantID), stringField("holiday_id", id.String()))
	}
	return mapHoliday(row), nil
}

func (s *Store) UpdateHoliday(ctx context.Context, item *domain.Holiday, actorID *uuid.UUID) (*domain.Holiday, error) {
	row, err := s.getQueries(ctx).UpdateHoliday(ctx, sqlc.UpdateHolidayParams{
		TenantID:   item.TenantID,
		ID:         item.ID,
		BranchID:   uuidFromPtr(item.BranchID),
		FyID:       uuidFromPtr(item.FYID),
		Name:       item.Name,
		Date:       dateFromTime(item.Date),
		IsOptional: item.IsOptional,
		UpdatedBy:  uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update holiday", err, tenantIDField(item.TenantID), stringField("holiday_id", item.ID.String()))
	}
	return mapHoliday(row), nil
}

func (s *Store) DeleteHoliday(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteHoliday(ctx, sqlc.SoftDeleteHolidayParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete holiday", err, tenantIDField(tenantID), stringField("holiday_id", id.String()))
	}
	return nil
}
