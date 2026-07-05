package postgres

import (
	"encoding/json"
	"math"
	"math/big"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapCompPayBand(row sqlc.HrmsCompensationPayBand) *domain.CompensationPayBand {
	return &domain.CompensationPayBand{ID: row.ID, TenantID: row.TenantID, Code: row.Code, Name: row.Name, JobFamily: ptrFromText(row.JobFamily), LevelCode: ptrFromText(row.LevelCode), LocationLabel: ptrFromText(row.LocationLabel), CurrencyCode: row.CurrencyCode, MinPay: compFloat(row.MinPay), MidpointPay: compFloat(row.MidpointPay), MaxPay: compFloat(row.MaxPay), EffectiveFrom: ptrFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), IsActive: row.IsActive, Notes: ptrFromText(row.Notes), Metadata: compRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapCompPayBands(rows []sqlc.HrmsCompensationPayBand) []*domain.CompensationPayBand {
	items := make([]*domain.CompensationPayBand, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCompPayBand(row))
	}
	return items
}

func mapCompCycle(row sqlc.HrmsCompensationCycle) *domain.CompensationCycle {
	return &domain.CompensationCycle{ID: row.ID, TenantID: row.TenantID, Code: row.Code, Name: row.Name, FiscalYearID: ptrFromUUID(row.FiscalYearID), Status: row.Status, CycleType: row.CycleType, StartsOn: ptrFromDate(row.StartsOn), EndsOn: ptrFromDate(row.EndsOn), EffectiveDate: ptrFromDate(row.EffectiveDate), CurrencyCode: row.CurrencyCode, BudgetAmount: compFloat(row.BudgetAmount), PlanningGuidance: ptrFromText(row.PlanningGuidance), ApprovalPolicy: ptrFromText(row.ApprovalPolicy), FinalizedAt: ptrFromTimestamptz(row.FinalizedAt), Metadata: compRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapCompCycles(rows []sqlc.HrmsCompensationCycle) []*domain.CompensationCycle {
	items := make([]*domain.CompensationCycle, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCompCycle(row))
	}
	return items
}

func mapCompBudgetPool(row sqlc.HrmsCompensationBudgetPool) *domain.CompensationBudgetPool {
	return &domain.CompensationBudgetPool{ID: row.ID, TenantID: row.TenantID, CycleID: row.CycleID, Name: row.Name, PoolType: row.PoolType, OwnerUserID: ptrFromUUID(row.OwnerUserID), DepartmentID: ptrFromUUID(row.DepartmentID), BranchID: ptrFromUUID(row.BranchID), BudgetAmount: compFloat(row.BudgetAmount), AllocatedAmount: compFloat(row.AllocatedAmount), Notes: ptrFromText(row.Notes), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapCompBudgetPools(rows []sqlc.ListCompensationBudgetPoolsRow) []*domain.CompensationBudgetPool {
	items := make([]*domain.CompensationBudgetPool, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.CompensationBudgetPool{ID: row.ID, TenantID: row.TenantID, CycleID: row.CycleID, Name: row.Name, PoolType: row.PoolType, OwnerUserID: ptrFromUUID(row.OwnerUserID), DepartmentID: ptrFromUUID(row.DepartmentID), BranchID: ptrFromUUID(row.BranchID), BudgetAmount: compFloat(row.BudgetAmount), AllocatedAmount: compFloat(row.AllocatedAmount), CommittedAmount: compFloat(row.CommittedAmount), RecommendationCount: row.RecommendationCount, Notes: ptrFromText(row.Notes), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)})
	}
	return items
}

func mapCompRecommendation(row sqlc.HrmsCompensationRecommendation) *domain.CompensationRecommendation {
	return compRecommendationFromParts(row.ID, row.TenantID, row.CycleID, row.WorkerProfileID, row.PayBandID, row.BudgetPoolID, row.CurrentSalary, row.CurrentCompaRatio, row.RecommendedSalary, row.RecommendedIncrementAmount, row.RecommendedIncrementPercent, row.PromotionRecommended, row.RecommendedDesignationID, row.Reason, row.PerformanceRating, row.EquityFlag, row.EquityNotes, row.Status, row.EffectiveDate, row.PayrollHandoffAt, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil)
}

func mapCompRecommendationList(rows []sqlc.ListCompensationRecommendationsRow) []*domain.CompensationRecommendation {
	items := make([]*domain.CompensationRecommendation, 0, len(rows))
	for _, row := range rows {
		items = append(items, compRecommendationFromParts(row.ID, row.TenantID, row.CycleID, row.WorkerProfileID, row.PayBandID, row.BudgetPoolID, row.CurrentSalary, row.CurrentCompaRatio, row.RecommendedSalary, row.RecommendedIncrementAmount, row.RecommendedIncrementPercent, row.PromotionRecommended, row.RecommendedDesignationID, row.Reason, row.PerformanceRating, row.EquityFlag, row.EquityNotes, row.Status, row.EffectiveDate, row.PayrollHandoffAt, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.WorkerDisplayName, ptrFromText(row.WorkerCode), ptrFromText(row.PayBandCode), ptrFromText(row.PayBandName), ptrFromText(row.BudgetPoolName)))
	}
	return items
}

func mapCompRecommendationDetail(row sqlc.GetCompensationRecommendationRow) *domain.CompensationRecommendation {
	return compRecommendationFromParts(row.ID, row.TenantID, row.CycleID, row.WorkerProfileID, row.PayBandID, row.BudgetPoolID, row.CurrentSalary, row.CurrentCompaRatio, row.RecommendedSalary, row.RecommendedIncrementAmount, row.RecommendedIncrementPercent, row.PromotionRecommended, row.RecommendedDesignationID, row.Reason, row.PerformanceRating, row.EquityFlag, row.EquityNotes, row.Status, row.EffectiveDate, row.PayrollHandoffAt, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.WorkerDisplayName, ptrFromText(row.WorkerCode), ptrFromText(row.PayBandCode), ptrFromText(row.PayBandName), ptrFromText(row.BudgetPoolName))
}

func compRecommendationFromParts(id uuid.UUID, tenantID uuid.UUID, cycleID uuid.UUID, workerProfileID uuid.UUID, payBandID pgtype.UUID, budgetPoolID pgtype.UUID, currentSalary pgtype.Numeric, currentCompaRatio pgtype.Numeric, recommendedSalary pgtype.Numeric, recommendedIncrementAmount pgtype.Numeric, recommendedIncrementPercent pgtype.Numeric, promotionRecommended bool, recommendedDesignationID pgtype.UUID, reason pgtype.Text, performanceRating pgtype.Text, equityFlag bool, equityNotes pgtype.Text, status string, effectiveDate pgtype.Date, payrollHandoffAt pgtype.Timestamptz, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, workerDisplayName *string, workerCode *string, payBandCode *string, payBandName *string, budgetPoolName *string) *domain.CompensationRecommendation {
	return &domain.CompensationRecommendation{ID: id, TenantID: tenantID, CycleID: cycleID, WorkerProfileID: workerProfileID, PayBandID: ptrFromUUID(payBandID), BudgetPoolID: ptrFromUUID(budgetPoolID), CurrentSalary: compFloat(currentSalary), CurrentCompaRatio: compFloat(currentCompaRatio), RecommendedSalary: compFloat(recommendedSalary), RecommendedIncrementAmount: compFloat(recommendedIncrementAmount), RecommendedIncrementPercent: compFloat(recommendedIncrementPercent), PromotionRecommended: promotionRecommended, RecommendedDesignationID: ptrFromUUID(recommendedDesignationID), Reason: ptrFromText(reason), PerformanceRating: ptrFromText(performanceRating), EquityFlag: equityFlag, EquityNotes: ptrFromText(equityNotes), Status: status, EffectiveDate: ptrFromDate(effectiveDate), PayrollHandoffAt: ptrFromTimestamptz(payrollHandoffAt), Metadata: compRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), WorkerDisplayName: workerDisplayName, WorkerCode: workerCode, PayBandCode: payBandCode, PayBandName: payBandName, BudgetPoolName: budgetPoolName}
}

func mapCompEquityGenerated(rows []sqlc.GenerateCompensationEquityChecksRow) []*domain.CompensationEquityCheck {
	items := make([]*domain.CompensationEquityCheck, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.CompensationEquityCheck{ID: row.ID, TenantID: row.TenantID, CycleID: row.CycleID, WorkerProfileID: row.WorkerProfileID, PayBandID: ptrFromUUID(row.PayBandID), CheckType: row.CheckType, Severity: row.Severity, CurrentSalary: compFloat(row.CurrentSalary), BandMin: compFloat(row.BandMin), BandMidpoint: compFloat(row.BandMidpoint), BandMax: compFloat(row.BandMax), VariancePercent: compFloat(row.VariancePercent), Finding: row.Finding, Recommendation: ptrFromText(row.Recommendation), Status: row.Status, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)})
	}
	return items
}

func mapCompEquityList(rows []sqlc.ListCompensationEquityChecksRow) []*domain.CompensationEquityCheck {
	items := make([]*domain.CompensationEquityCheck, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.CompensationEquityCheck{ID: row.ID, TenantID: row.TenantID, CycleID: row.CycleID, WorkerProfileID: row.WorkerProfileID, PayBandID: ptrFromUUID(row.PayBandID), CheckType: row.CheckType, Severity: row.Severity, CurrentSalary: compFloat(row.CurrentSalary), BandMin: compFloat(row.BandMin), BandMidpoint: compFloat(row.BandMidpoint), BandMax: compFloat(row.BandMax), VariancePercent: compFloat(row.VariancePercent), Finding: row.Finding, Recommendation: ptrFromText(row.Recommendation), Status: row.Status, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy), WorkerDisplayName: &row.WorkerDisplayName, WorkerCode: ptrFromText(row.WorkerCode), PayBandCode: ptrFromText(row.PayBandCode), PayBandName: ptrFromText(row.PayBandName)})
	}
	return items
}

func mapCompEquity(row sqlc.HrmsCompensationEquityCheck) *domain.CompensationEquityCheck {
	return &domain.CompensationEquityCheck{ID: row.ID, TenantID: row.TenantID, CycleID: row.CycleID, WorkerProfileID: row.WorkerProfileID, PayBandID: ptrFromUUID(row.PayBandID), CheckType: row.CheckType, Severity: row.Severity, CurrentSalary: compFloat(row.CurrentSalary), BandMin: compFloat(row.BandMin), BandMidpoint: compFloat(row.BandMidpoint), BandMax: compFloat(row.BandMax), VariancePercent: compFloat(row.VariancePercent), Finding: row.Finding, Recommendation: ptrFromText(row.Recommendation), Status: row.Status, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapCompEvent(row sqlc.HrmsCompensationEvent) *domain.CompensationEvent {
	return &domain.CompensationEvent{ID: row.ID, TenantID: row.TenantID, CycleID: ptrFromUUID(row.CycleID), SourceType: row.SourceType, SourceID: ptrFromUUID(row.SourceID), Action: row.Action, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), Remarks: ptrFromText(row.Remarks), Metadata: compRaw(row.Metadata), CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy)}
}

func mapCompEvents(rows []sqlc.HrmsCompensationEvent) []*domain.CompensationEvent {
	items := make([]*domain.CompensationEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCompEvent(row))
	}
	return items
}

func mapCompSummary(rows []sqlc.GetCompensationSummaryRow) []*domain.CompensationSummaryRow {
	items := make([]*domain.CompensationSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.CompensationSummaryRow{Metric: row.Metric, MetricCount: row.MetricCount, MetricAmount: compFloat(row.MetricAmount)})
	}
	return items
}

func compRaw(value []byte) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(value)
}

func compNumeric(value float64) pgtype.Numeric {
	scaled := int64(math.Round(value * 10000))
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -4, Valid: true}
}

func compFloat(value pgtype.Numeric) float64 {
	floatValue := ptrFromNumeric(value)
	if floatValue == nil {
		return 0
	}
	return *floatValue
}
