package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateCompensationPayBand(ctx context.Context, item *domain.CompensationPayBand, actorID *uuid.UUID) (*domain.CompensationPayBand, error) {
	row, err := s.getQueries(ctx).CreateCompensationPayBand(ctx, sqlc.CreateCompensationPayBandParams{TenantID: item.TenantID, Code: item.Code, Name: item.Name, JobFamily: textFromPtr(item.JobFamily), LevelCode: textFromPtr(item.LevelCode), LocationLabel: textFromPtr(item.LocationLabel), CurrencyCode: item.CurrencyCode, MinPay: compNumeric(item.MinPay), MidpointPay: compNumeric(item.MidpointPay), MaxPay: compNumeric(item.MaxPay), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), IsActive: item.IsActive, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create compensation pay band", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapCompPayBand(row), nil
}

func (s *Store) UpdateCompensationPayBand(ctx context.Context, item *domain.CompensationPayBand, actorID *uuid.UUID) (*domain.CompensationPayBand, error) {
	row, err := s.getQueries(ctx).UpdateCompensationPayBand(ctx, sqlc.UpdateCompensationPayBandParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Name: item.Name, JobFamily: textFromPtr(item.JobFamily), LevelCode: textFromPtr(item.LevelCode), LocationLabel: textFromPtr(item.LocationLabel), CurrencyCode: item.CurrencyCode, MinPay: compNumeric(item.MinPay), MidpointPay: compNumeric(item.MidpointPay), MaxPay: compNumeric(item.MaxPay), EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), IsActive: item.IsActive, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationPayBandNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update compensation pay band", err, tenantIDField(item.TenantID), stringField("pay_band_id", item.ID.String()))
	}
	return mapCompPayBand(row), nil
}

func (s *Store) GetCompensationPayBand(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompensationPayBand, error) {
	row, err := s.getQueries(ctx).GetCompensationPayBand(ctx, sqlc.GetCompensationPayBandParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationPayBandNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get compensation pay band", err, tenantIDField(tenantID), stringField("pay_band_id", id.String()))
	}
	return mapCompPayBand(row), nil
}

func (s *Store) ListCompensationPayBands(ctx context.Context, filter domain.CompensationPayBandFilter) ([]*domain.CompensationPayBand, error) {
	rows, err := s.getQueries(ctx).ListCompensationPayBands(ctx, sqlc.ListCompensationPayBandsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, IsActive: boolFromSkillPtr(filter.IsActive), CurrencyCode: textFromPtr(filter.CurrencyCode), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list compensation pay bands", err, tenantIDField(filter.TenantID))
	}
	return mapCompPayBands(rows), nil
}

func (s *Store) DeleteCompensationPayBand(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCompensationPayBand(ctx, sqlc.SoftDeleteCompensationPayBandParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete compensation pay band", err, tenantIDField(tenantID), stringField("pay_band_id", id.String()))
	}
	return nil
}

func (s *Store) CreateCompensationCycle(ctx context.Context, item *domain.CompensationCycle, actorID *uuid.UUID) (*domain.CompensationCycle, error) {
	row, err := s.getQueries(ctx).CreateCompensationCycle(ctx, sqlc.CreateCompensationCycleParams{TenantID: item.TenantID, Code: item.Code, Name: item.Name, FiscalYearID: uuidFromPtr(item.FiscalYearID), Status: item.Status, CycleType: item.CycleType, StartsOn: dateFromPtr(item.StartsOn), EndsOn: dateFromPtr(item.EndsOn), EffectiveDate: dateFromPtr(item.EffectiveDate), CurrencyCode: item.CurrencyCode, BudgetAmount: compNumeric(item.BudgetAmount), PlanningGuidance: textFromPtr(item.PlanningGuidance), ApprovalPolicy: textFromPtr(item.ApprovalPolicy), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create compensation cycle", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapCompCycle(row), nil
}

func (s *Store) UpdateCompensationCycle(ctx context.Context, item *domain.CompensationCycle, actorID *uuid.UUID) (*domain.CompensationCycle, error) {
	row, err := s.getQueries(ctx).UpdateCompensationCycle(ctx, sqlc.UpdateCompensationCycleParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Name: item.Name, FiscalYearID: uuidFromPtr(item.FiscalYearID), CycleType: item.CycleType, StartsOn: dateFromPtr(item.StartsOn), EndsOn: dateFromPtr(item.EndsOn), EffectiveDate: dateFromPtr(item.EffectiveDate), CurrencyCode: item.CurrencyCode, BudgetAmount: compNumeric(item.BudgetAmount), PlanningGuidance: textFromPtr(item.PlanningGuidance), ApprovalPolicy: textFromPtr(item.ApprovalPolicy), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationCycleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update compensation cycle", err, tenantIDField(item.TenantID), stringField("cycle_id", item.ID.String()))
	}
	return mapCompCycle(row), nil
}

func (s *Store) UpdateCompensationCycleStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.CompensationCycle, error) {
	row, err := s.getQueries(ctx).UpdateCompensationCycleStatus(ctx, sqlc.UpdateCompensationCycleStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationCycleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update compensation cycle status", err, tenantIDField(tenantID), stringField("cycle_id", id.String()), stringField("status", status))
	}
	return mapCompCycle(row), nil
}

func (s *Store) GetCompensationCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompensationCycle, error) {
	row, err := s.getQueries(ctx).GetCompensationCycle(ctx, sqlc.GetCompensationCycleParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationCycleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get compensation cycle", err, tenantIDField(tenantID), stringField("cycle_id", id.String()))
	}
	return mapCompCycle(row), nil
}

func (s *Store) ListCompensationCycles(ctx context.Context, filter domain.CompensationFilter) ([]*domain.CompensationCycle, error) {
	rows, err := s.getQueries(ctx).ListCompensationCycles(ctx, sqlc.ListCompensationCyclesParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list compensation cycles", err, tenantIDField(filter.TenantID))
	}
	return mapCompCycles(rows), nil
}

func (s *Store) CreateCompensationBudgetPool(ctx context.Context, item *domain.CompensationBudgetPool, actorID *uuid.UUID) (*domain.CompensationBudgetPool, error) {
	row, err := s.getQueries(ctx).CreateCompensationBudgetPool(ctx, sqlc.CreateCompensationBudgetPoolParams{TenantID: item.TenantID, CycleID: item.CycleID, Name: item.Name, PoolType: item.PoolType, OwnerUserID: uuidFromPtr(item.OwnerUserID), DepartmentID: uuidFromPtr(item.DepartmentID), BranchID: uuidFromPtr(item.BranchID), BudgetAmount: compNumeric(item.BudgetAmount), AllocatedAmount: compNumeric(item.AllocatedAmount), Notes: textFromPtr(item.Notes), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create compensation budget pool", err, tenantIDField(item.TenantID), stringField("cycle_id", item.CycleID.String()))
	}
	return mapCompBudgetPool(row), nil
}

func (s *Store) UpdateCompensationBudgetPool(ctx context.Context, item *domain.CompensationBudgetPool, actorID *uuid.UUID) (*domain.CompensationBudgetPool, error) {
	row, err := s.getQueries(ctx).UpdateCompensationBudgetPool(ctx, sqlc.UpdateCompensationBudgetPoolParams{TenantID: item.TenantID, ID: item.ID, CycleID: item.CycleID, Name: item.Name, PoolType: item.PoolType, OwnerUserID: uuidFromPtr(item.OwnerUserID), DepartmentID: uuidFromPtr(item.DepartmentID), BranchID: uuidFromPtr(item.BranchID), BudgetAmount: compNumeric(item.BudgetAmount), AllocatedAmount: compNumeric(item.AllocatedAmount), Notes: textFromPtr(item.Notes), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationBudgetPoolNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update compensation budget pool", err, tenantIDField(item.TenantID), stringField("pool_id", item.ID.String()))
	}
	return mapCompBudgetPool(row), nil
}

func (s *Store) ListCompensationBudgetPools(ctx context.Context, tenantID uuid.UUID, cycleID uuid.UUID) ([]*domain.CompensationBudgetPool, error) {
	rows, err := s.getQueries(ctx).ListCompensationBudgetPools(ctx, sqlc.ListCompensationBudgetPoolsParams{TenantID: tenantID, CycleID: cycleID})
	if err != nil {
		return nil, s.logDBError(ctx, "list compensation budget pools", err, tenantIDField(tenantID), stringField("cycle_id", cycleID.String()))
	}
	return mapCompBudgetPools(rows), nil
}

func (s *Store) DeleteCompensationBudgetPool(ctx context.Context, tenantID uuid.UUID, cycleID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCompensationBudgetPool(ctx, sqlc.SoftDeleteCompensationBudgetPoolParams{TenantID: tenantID, CycleID: cycleID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete compensation budget pool", err, tenantIDField(tenantID), stringField("pool_id", id.String()))
	}
	return nil
}

func (s *Store) CreateCompensationRecommendation(ctx context.Context, item *domain.CompensationRecommendation, actorID *uuid.UUID) (*domain.CompensationRecommendation, error) {
	row, err := s.getQueries(ctx).CreateCompensationRecommendation(ctx, sqlc.CreateCompensationRecommendationParams{TenantID: item.TenantID, CycleID: item.CycleID, WorkerProfileID: item.WorkerProfileID, PayBandID: uuidFromPtr(item.PayBandID), BudgetPoolID: uuidFromPtr(item.BudgetPoolID), CurrentSalary: compNumeric(item.CurrentSalary), CurrentCompaRatio: compNumeric(item.CurrentCompaRatio), RecommendedSalary: compNumeric(item.RecommendedSalary), RecommendedIncrementAmount: compNumeric(item.RecommendedIncrementAmount), RecommendedIncrementPercent: compNumeric(item.RecommendedIncrementPercent), PromotionRecommended: item.PromotionRecommended, RecommendedDesignationID: uuidFromPtr(item.RecommendedDesignationID), Reason: textFromPtr(item.Reason), PerformanceRating: textFromPtr(item.PerformanceRating), EquityFlag: item.EquityFlag, EquityNotes: textFromPtr(item.EquityNotes), Status: item.Status, EffectiveDate: dateFromPtr(item.EffectiveDate), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create compensation recommendation", err, tenantIDField(item.TenantID), stringField("cycle_id", item.CycleID.String()))
	}
	return mapCompRecommendation(row), nil
}

func (s *Store) UpdateCompensationRecommendation(ctx context.Context, item *domain.CompensationRecommendation, actorID *uuid.UUID) (*domain.CompensationRecommendation, error) {
	row, err := s.getQueries(ctx).UpdateCompensationRecommendation(ctx, sqlc.UpdateCompensationRecommendationParams{TenantID: item.TenantID, ID: item.ID, CycleID: item.CycleID, PayBandID: uuidFromPtr(item.PayBandID), BudgetPoolID: uuidFromPtr(item.BudgetPoolID), CurrentSalary: compNumeric(item.CurrentSalary), CurrentCompaRatio: compNumeric(item.CurrentCompaRatio), RecommendedSalary: compNumeric(item.RecommendedSalary), RecommendedIncrementAmount: compNumeric(item.RecommendedIncrementAmount), RecommendedIncrementPercent: compNumeric(item.RecommendedIncrementPercent), PromotionRecommended: item.PromotionRecommended, RecommendedDesignationID: uuidFromPtr(item.RecommendedDesignationID), Reason: textFromPtr(item.Reason), PerformanceRating: textFromPtr(item.PerformanceRating), EquityFlag: item.EquityFlag, EquityNotes: textFromPtr(item.EquityNotes), EffectiveDate: dateFromPtr(item.EffectiveDate), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationRecommendationNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update compensation recommendation", err, tenantIDField(item.TenantID), stringField("recommendation_id", item.ID.String()))
	}
	return mapCompRecommendation(row), nil
}

func (s *Store) UpdateCompensationRecommendationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.CompensationRecommendation, error) {
	row, err := s.getQueries(ctx).UpdateCompensationRecommendationStatus(ctx, sqlc.UpdateCompensationRecommendationStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationRecommendationNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update compensation recommendation status", err, tenantIDField(tenantID), stringField("recommendation_id", id.String()), stringField("status", status))
	}
	return mapCompRecommendation(row), nil
}

func (s *Store) GetCompensationRecommendation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompensationRecommendation, error) {
	row, err := s.getQueries(ctx).GetCompensationRecommendation(ctx, sqlc.GetCompensationRecommendationParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationRecommendationNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get compensation recommendation", err, tenantIDField(tenantID), stringField("recommendation_id", id.String()))
	}
	return mapCompRecommendationDetail(row), nil
}

func (s *Store) ListCompensationRecommendations(ctx context.Context, filter domain.CompensationFilter) ([]*domain.CompensationRecommendation, error) {
	rows, err := s.getQueries(ctx).ListCompensationRecommendations(ctx, sqlc.ListCompensationRecommendationsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, CycleID: uuidFromPtr(filter.CycleID), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list compensation recommendations", err, tenantIDField(filter.TenantID))
	}
	return mapCompRecommendationList(rows), nil
}

func (s *Store) DeleteCompensationRecommendation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCompensationRecommendation(ctx, sqlc.SoftDeleteCompensationRecommendationParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete compensation recommendation", err, tenantIDField(tenantID), stringField("recommendation_id", id.String()))
	}
	return nil
}

func (s *Store) GenerateCompensationEquityChecks(ctx context.Context, tenantID uuid.UUID, cycleID uuid.UUID, actorID *uuid.UUID) ([]*domain.CompensationEquityCheck, error) {
	rows, err := s.getQueries(ctx).GenerateCompensationEquityChecks(ctx, sqlc.GenerateCompensationEquityChecksParams{TenantID: tenantID, CycleID: cycleID, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "generate compensation equity checks", err, tenantIDField(tenantID), stringField("cycle_id", cycleID.String()))
	}
	return mapCompEquityGenerated(rows), nil
}

func (s *Store) ListCompensationEquityChecks(ctx context.Context, filter domain.CompensationFilter) ([]*domain.CompensationEquityCheck, error) {
	rows, err := s.getQueries(ctx).ListCompensationEquityChecks(ctx, sqlc.ListCompensationEquityChecksParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, CycleID: uuidFromPtr(filter.CycleID), Status: textFromPtr(filter.Status)})
	if err != nil {
		return nil, s.logDBError(ctx, "list compensation equity checks", err, tenantIDField(filter.TenantID))
	}
	return mapCompEquityList(rows), nil
}

func (s *Store) UpdateCompensationEquityCheckStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.CompensationEquityCheck, error) {
	row, err := s.getQueries(ctx).UpdateCompensationEquityCheckStatus(ctx, sqlc.UpdateCompensationEquityCheckStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCompensationEquityCheckNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update compensation equity check status", err, tenantIDField(tenantID), stringField("equity_check_id", id.String()))
	}
	return mapCompEquity(row), nil
}

func (s *Store) CreateCompensationEvent(ctx context.Context, item *domain.CompensationEvent, actorID *uuid.UUID) (*domain.CompensationEvent, error) {
	row, err := s.getQueries(ctx).CreateCompensationEvent(ctx, sqlc.CreateCompensationEventParams{TenantID: item.TenantID, CycleID: uuidFromPtr(item.CycleID), SourceType: item.SourceType, SourceID: uuidFromPtr(item.SourceID), Action: item.Action, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), Remarks: textFromPtr(item.Remarks), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create compensation event", err, tenantIDField(item.TenantID), stringField("source_type", item.SourceType))
	}
	return mapCompEvent(row), nil
}

func (s *Store) ListCompensationEvents(ctx context.Context, filter domain.CompensationFilter, sourceType *string, sourceID *uuid.UUID) ([]*domain.CompensationEvent, error) {
	rows, err := s.getQueries(ctx).ListCompensationEvents(ctx, sqlc.ListCompensationEventsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, CycleID: uuidFromPtr(filter.CycleID), SourceType: textFromPtr(sourceType), SourceID: uuidFromPtr(sourceID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list compensation events", err, tenantIDField(filter.TenantID))
	}
	return mapCompEvents(rows), nil
}

func (s *Store) GetCompensationSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.CompensationSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetCompensationSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get compensation summary", err, tenantIDField(tenantID))
	}
	return mapCompSummary(rows), nil
}
