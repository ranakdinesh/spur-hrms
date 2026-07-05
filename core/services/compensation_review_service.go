package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateCompensationPayBand(ctx context.Context, cmd ports.CompensationPayBandCommand) (*domain.CompensationPayBand, error) {
	item := &domain.CompensationPayBand{TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, JobFamily: cmd.JobFamily, LevelCode: cmd.LevelCode, LocationLabel: cmd.LocationLabel, CurrencyCode: cmd.CurrencyCode, MinPay: cmd.MinPay, MidpointPay: cmd.MidpointPay, MaxPay: cmd.MaxPay, EffectiveFrom: cmd.EffectiveFrom, EffectiveTo: cmd.EffectiveTo, IsActive: cmd.IsActive, Notes: cmd.Notes, Metadata: cmd.Metadata}
	if err := domain.ValidateCompensationPayBand(item); err != nil {
		s.logError("validate compensation pay band", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.compensationReview.CreateCompensationPayBand(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateCompensationPayBand(ctx context.Context, cmd ports.CompensationPayBandCommand) (*domain.CompensationPayBand, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidCompensationPayBand
	}
	item := &domain.CompensationPayBand{ID: cmd.ID, TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, JobFamily: cmd.JobFamily, LevelCode: cmd.LevelCode, LocationLabel: cmd.LocationLabel, CurrencyCode: cmd.CurrencyCode, MinPay: cmd.MinPay, MidpointPay: cmd.MidpointPay, MaxPay: cmd.MaxPay, EffectiveFrom: cmd.EffectiveFrom, EffectiveTo: cmd.EffectiveTo, IsActive: cmd.IsActive, Notes: cmd.Notes, Metadata: cmd.Metadata}
	if err := domain.ValidateCompensationPayBand(item); err != nil {
		return nil, err
	}
	return s.compensationReview.UpdateCompensationPayBand(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListCompensationPayBands(ctx context.Context, filter domain.CompensationPayBandFilter) ([]*domain.CompensationPayBand, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeCompPage(&filter.Limit, &filter.Offset)
	return s.compensationReview.ListCompensationPayBands(ctx, filter)
}

func (s *TenantService) DeleteCompensationPayBand(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidCompensationPayBand
	}
	return s.compensationReview.DeleteCompensationPayBand(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateCompensationCycle(ctx context.Context, cmd ports.CompensationCycleCommand) (*domain.CompensationCycle, error) {
	item := &domain.CompensationCycle{TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, FiscalYearID: cmd.FiscalYearID, Status: cmd.Status, CycleType: cmd.CycleType, StartsOn: cmd.StartsOn, EndsOn: cmd.EndsOn, EffectiveDate: cmd.EffectiveDate, CurrencyCode: cmd.CurrencyCode, BudgetAmount: cmd.BudgetAmount, PlanningGuidance: cmd.PlanningGuidance, ApprovalPolicy: cmd.ApprovalPolicy, Metadata: cmd.Metadata}
	if err := domain.ValidateCompensationCycle(item); err != nil {
		s.logError("validate compensation cycle", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.compensationReview.CreateCompensationCycle(ctx, item, cmd.ActorID)
	if err == nil {
		_, _ = s.compEvent(ctx, result.TenantID, &result.ID, "cycle", &result.ID, "created", nil, &result.Status, nil, cmd.ActorID)
	}
	return result, err
}

func (s *TenantService) UpdateCompensationCycle(ctx context.Context, cmd ports.CompensationCycleCommand) (*domain.CompensationCycle, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidCompensationCycle
	}
	existing, err := s.compensationReview.GetCompensationCycle(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	item := &domain.CompensationCycle{ID: cmd.ID, TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, FiscalYearID: cmd.FiscalYearID, Status: existing.Status, CycleType: cmd.CycleType, StartsOn: cmd.StartsOn, EndsOn: cmd.EndsOn, EffectiveDate: cmd.EffectiveDate, CurrencyCode: cmd.CurrencyCode, BudgetAmount: cmd.BudgetAmount, PlanningGuidance: cmd.PlanningGuidance, ApprovalPolicy: cmd.ApprovalPolicy, Metadata: cmd.Metadata}
	if err := domain.ValidateCompensationCycle(item); err != nil {
		return nil, err
	}
	return s.compensationReview.UpdateCompensationCycle(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateCompensationCycleStatus(ctx context.Context, cmd ports.CompensationStatusCommand) (*domain.CompensationCycle, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidCompensationCycle
	}
	status := domain.NormalizeCompensationCycleStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidCompensationCycle
	}
	before, _ := s.compensationReview.GetCompensationCycle(ctx, cmd.TenantID, cmd.ID)
	result, err := s.compensationReview.UpdateCompensationCycleStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	_, _ = s.compEvent(ctx, cmd.TenantID, &cmd.ID, "cycle", &cmd.ID, "status_changed", from, &result.Status, cmd.Remarks, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListCompensationCycles(ctx context.Context, filter domain.CompensationFilter) ([]*domain.CompensationCycle, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeCompPage(&filter.Limit, &filter.Offset)
	return s.compensationReview.ListCompensationCycles(ctx, filter)
}

func (s *TenantService) CreateCompensationBudgetPool(ctx context.Context, cmd ports.CompensationBudgetPoolCommand) (*domain.CompensationBudgetPool, error) {
	if _, err := s.compensationReview.GetCompensationCycle(ctx, cmd.TenantID, cmd.CycleID); err != nil {
		return nil, err
	}
	item := &domain.CompensationBudgetPool{TenantID: cmd.TenantID, CycleID: cmd.CycleID, Name: cmd.Name, PoolType: cmd.PoolType, OwnerUserID: cmd.OwnerUserID, DepartmentID: cmd.DepartmentID, BranchID: cmd.BranchID, BudgetAmount: cmd.BudgetAmount, AllocatedAmount: cmd.AllocatedAmount, Notes: cmd.Notes}
	if err := domain.ValidateCompensationBudgetPool(item); err != nil {
		return nil, err
	}
	return s.compensationReview.CreateCompensationBudgetPool(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateCompensationBudgetPool(ctx context.Context, cmd ports.CompensationBudgetPoolCommand) (*domain.CompensationBudgetPool, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidCompensationBudgetPool
	}
	item := &domain.CompensationBudgetPool{ID: cmd.ID, TenantID: cmd.TenantID, CycleID: cmd.CycleID, Name: cmd.Name, PoolType: cmd.PoolType, OwnerUserID: cmd.OwnerUserID, DepartmentID: cmd.DepartmentID, BranchID: cmd.BranchID, BudgetAmount: cmd.BudgetAmount, AllocatedAmount: cmd.AllocatedAmount, Notes: cmd.Notes}
	if err := domain.ValidateCompensationBudgetPool(item); err != nil {
		return nil, err
	}
	return s.compensationReview.UpdateCompensationBudgetPool(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListCompensationBudgetPools(ctx context.Context, tenantID uuid.UUID, cycleID uuid.UUID) ([]*domain.CompensationBudgetPool, error) {
	if tenantID == uuid.Nil || cycleID == uuid.Nil {
		return nil, domain.ErrInvalidCompensationBudgetPool
	}
	return s.compensationReview.ListCompensationBudgetPools(ctx, tenantID, cycleID)
}

func (s *TenantService) DeleteCompensationBudgetPool(ctx context.Context, tenantID uuid.UUID, cycleID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || cycleID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidCompensationBudgetPool
	}
	return s.compensationReview.DeleteCompensationBudgetPool(ctx, tenantID, cycleID, id, actorID)
}

func (s *TenantService) CreateCompensationRecommendation(ctx context.Context, cmd ports.CompensationRecommendationCommand) (*domain.CompensationRecommendation, error) {
	if _, err := s.compensationReview.GetCompensationCycle(ctx, cmd.TenantID, cmd.CycleID); err != nil {
		return nil, err
	}
	if _, err := s.workerProfiles.GetWorkerProfile(ctx, cmd.TenantID, cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	item := compRecommendationFromCommand(cmd)
	if err := domain.ValidateCompensationRecommendation(item); err != nil {
		return nil, err
	}
	result, err := s.compensationReview.CreateCompensationRecommendation(ctx, item, cmd.ActorID)
	if err == nil {
		_, _ = s.compEvent(ctx, result.TenantID, &result.CycleID, "recommendation", &result.ID, "created", nil, &result.Status, nil, cmd.ActorID)
	}
	return result, err
}

func (s *TenantService) UpdateCompensationRecommendation(ctx context.Context, cmd ports.CompensationRecommendationCommand) (*domain.CompensationRecommendation, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidCompensationRecommendation
	}
	item := compRecommendationFromCommand(cmd)
	item.ID = cmd.ID
	if err := domain.ValidateCompensationRecommendation(item); err != nil {
		return nil, err
	}
	return s.compensationReview.UpdateCompensationRecommendation(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateCompensationRecommendationStatus(ctx context.Context, cmd ports.CompensationStatusCommand) (*domain.CompensationRecommendation, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidCompensationRecommendation
	}
	status := domain.NormalizeCompensationRecommendationStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidCompensationRecommendation
	}
	before, _ := s.compensationReview.GetCompensationRecommendation(ctx, cmd.TenantID, cmd.ID)
	result, err := s.compensationReview.UpdateCompensationRecommendationStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	_, _ = s.compEvent(ctx, result.TenantID, &result.CycleID, "recommendation", &result.ID, "status_changed", from, &result.Status, cmd.Remarks, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListCompensationRecommendations(ctx context.Context, filter domain.CompensationFilter) ([]*domain.CompensationRecommendation, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeCompPage(&filter.Limit, &filter.Offset)
	return s.compensationReview.ListCompensationRecommendations(ctx, filter)
}

func (s *TenantService) DeleteCompensationRecommendation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidCompensationRecommendation
	}
	return s.compensationReview.DeleteCompensationRecommendation(ctx, tenantID, id, actorID)
}

func (s *TenantService) GenerateCompensationEquityChecks(ctx context.Context, tenantID uuid.UUID, cycleID uuid.UUID, actorID *uuid.UUID) ([]*domain.CompensationEquityCheck, error) {
	if tenantID == uuid.Nil || cycleID == uuid.Nil {
		return nil, domain.ErrInvalidCompensationEquityCheck
	}
	items, err := s.compensationReview.GenerateCompensationEquityChecks(ctx, tenantID, cycleID, actorID)
	if err == nil {
		action := "equity_generated"
		_, _ = s.compEvent(ctx, tenantID, &cycleID, "cycle", &cycleID, action, nil, nil, nil, actorID)
	}
	return items, err
}

func (s *TenantService) ListCompensationEquityChecks(ctx context.Context, filter domain.CompensationFilter) ([]*domain.CompensationEquityCheck, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeCompPage(&filter.Limit, &filter.Offset)
	return s.compensationReview.ListCompensationEquityChecks(ctx, filter)
}

func (s *TenantService) UpdateCompensationEquityCheckStatus(ctx context.Context, cmd ports.CompensationStatusCommand) (*domain.CompensationEquityCheck, error) {
	status := domain.NormalizeCompensationEquityStatus(cmd.Status)
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || status == "" {
		return nil, domain.ErrInvalidCompensationEquityCheck
	}
	return s.compensationReview.UpdateCompensationEquityCheckStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
}

func (s *TenantService) ListCompensationEvents(ctx context.Context, filter domain.CompensationFilter, sourceType *string, sourceID *uuid.UUID) ([]*domain.CompensationEvent, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeCompPage(&filter.Limit, &filter.Offset)
	return s.compensationReview.ListCompensationEvents(ctx, filter, sourceType, sourceID)
}

func (s *TenantService) GetCompensationSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.CompensationSummaryRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.compensationReview.GetCompensationSummary(ctx, tenantID)
}

func compRecommendationFromCommand(cmd ports.CompensationRecommendationCommand) *domain.CompensationRecommendation {
	return &domain.CompensationRecommendation{TenantID: cmd.TenantID, CycleID: cmd.CycleID, WorkerProfileID: cmd.WorkerProfileID, PayBandID: cmd.PayBandID, BudgetPoolID: cmd.BudgetPoolID, CurrentSalary: cmd.CurrentSalary, CurrentCompaRatio: cmd.CurrentCompaRatio, RecommendedSalary: cmd.RecommendedSalary, RecommendedIncrementAmount: cmd.RecommendedIncrementAmount, RecommendedIncrementPercent: cmd.RecommendedIncrementPercent, PromotionRecommended: cmd.PromotionRecommended, RecommendedDesignationID: cmd.RecommendedDesignationID, Reason: cmd.Reason, PerformanceRating: cmd.PerformanceRating, EquityFlag: cmd.EquityFlag, EquityNotes: cmd.EquityNotes, Status: cmd.Status, EffectiveDate: cmd.EffectiveDate, Metadata: cmd.Metadata}
}

func (s *TenantService) compEvent(ctx context.Context, tenantID uuid.UUID, cycleID *uuid.UUID, sourceType string, sourceID *uuid.UUID, action string, fromStatus *string, toStatus *string, remarks *string, actorID *uuid.UUID) (*domain.CompensationEvent, error) {
	return s.compensationReview.CreateCompensationEvent(ctx, &domain.CompensationEvent{TenantID: tenantID, CycleID: cycleID, SourceType: sourceType, SourceID: sourceID, Action: action, FromStatus: fromStatus, ToStatus: toStatus, Remarks: remarks}, actorID)
}

func normalizeCompPage(limit *int32, offset *int32) {
	if *limit <= 0 || *limit > 200 {
		*limit = 50
	}
	if *offset < 0 {
		*offset = 0
	}
}
