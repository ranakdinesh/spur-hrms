package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateFinancialYear(ctx context.Context, item *domain.FinancialYear, actorID *uuid.UUID) (*domain.FinancialYear, error) {
	row, err := s.getQueries(ctx).CreateFinancialYear(ctx, sqlc.CreateFinancialYearParams{
		TenantID:      item.TenantID,
		Name:          item.Name,
		StartDate:     dateFromTime(item.StartDate),
		EndDate:       dateFromTime(item.EndDate),
		PayrollYear:   item.PayrollYear,
		LeaveYear:     item.LeaveYear,
		HolidayYear:   item.HolidayYear,
		ReportingYear: item.ReportingYear,
		CloseNote:     textFromPtr(item.CloseNote),
		CreatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create financial year", err, tenantIDField(item.TenantID), stringField("financial_year_name", item.Name))
	}
	return mapFinancialYear(row), nil
}

func (s *Store) ListFinancialYears(ctx context.Context, tenantID uuid.UUID) ([]*domain.FinancialYear, error) {
	rows, err := s.getQueries(ctx).ListFinancialYears(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list financial years", err, tenantIDField(tenantID))
	}
	return mapFinancialYears(rows), nil
}

func (s *Store) GetFinancialYear(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.FinancialYear, error) {
	row, err := s.getQueries(ctx).GetFinancialYear(ctx, sqlc.GetFinancialYearParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get financial year", err, tenantIDField(tenantID), stringField("financial_year_id", id.String()))
	}
	return mapFinancialYear(row), nil
}

func (s *Store) GetActiveFinancialYear(ctx context.Context, tenantID uuid.UUID) (*domain.FinancialYear, error) {
	row, err := s.getQueries(ctx).GetActiveFinancialYear(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get active financial year", err, tenantIDField(tenantID))
	}
	return mapFinancialYear(row), nil
}

func (s *Store) UpdateFinancialYear(ctx context.Context, item *domain.FinancialYear, actorID *uuid.UUID) (*domain.FinancialYear, error) {
	row, err := s.getQueries(ctx).UpdateFinancialYear(ctx, sqlc.UpdateFinancialYearParams{
		TenantID:      item.TenantID,
		ID:            item.ID,
		Name:          item.Name,
		StartDate:     dateFromTime(item.StartDate),
		EndDate:       dateFromTime(item.EndDate),
		PayrollYear:   item.PayrollYear,
		LeaveYear:     item.LeaveYear,
		HolidayYear:   item.HolidayYear,
		ReportingYear: item.ReportingYear,
		CloseNote:     textFromPtr(item.CloseNote),
		UpdatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update financial year", err, tenantIDField(item.TenantID), stringField("financial_year_id", item.ID.String()))
	}
	return mapFinancialYear(row), nil
}

func (s *Store) DeleteFinancialYear(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteFinancialYear(ctx, sqlc.SoftDeleteFinancialYearParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete financial year", err, tenantIDField(tenantID), stringField("financial_year_id", id.String()))
	}
	return nil
}

func (s *Store) SetActiveFinancialYear(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.FinancialYear, error) {
	queries := s.getQueries(ctx)
	if err := queries.ClearActiveFinancialYears(ctx, sqlc.ClearActiveFinancialYearsParams{TenantID: tenantID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return nil, s.logDBError(ctx, "clear active financial years", err, tenantIDField(tenantID), stringField("financial_year_id", id.String()))
	}
	row, err := queries.MarkFinancialYearActive(ctx, sqlc.MarkFinancialYearActiveParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "set active financial year", err, tenantIDField(tenantID), stringField("financial_year_id", id.String()))
	}
	return mapFinancialYear(row), nil
}

func (s *Store) SetFinancialYearLock(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, isLocked bool, closeNote *string, actorID *uuid.UUID) (*domain.FinancialYear, error) {
	row, err := s.getQueries(ctx).SetFinancialYearLock(ctx, sqlc.SetFinancialYearLockParams{
		TenantID:  tenantID,
		ID:        id,
		IsLocked:  isLocked,
		UpdatedBy: uuidFromPtr(actorID),
		CloseNote: textFromPtr(closeNote),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "set financial year lock", err, tenantIDField(tenantID), stringField("financial_year_id", id.String()))
	}
	return mapFinancialYear(row), nil
}
