package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) UpsertPayrollPeriodLock(ctx context.Context, item *domain.PayrollPeriodLock, actorID *uuid.UUID) (*domain.PayrollPeriodLock, error) {
	row, err := s.getQueries(ctx).UpsertPayrollPeriodLock(ctx, sqlc.UpsertPayrollPeriodLockParams{TenantID: item.TenantID, Month: item.Month, Year: item.Year, Status: item.Status, CreatedBy: uuidFromPtr(actorID), UnlockReason: textFromPtr(item.UnlockReason), Notes: textFromPtr(item.Notes)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert payroll period lock", err, tenantIDField(item.TenantID))
	}
	return mapPayrollPeriodLock(row), nil
}

func (s *Store) GetPayrollPeriodLock(ctx context.Context, tenantID uuid.UUID, month int32, year int32) (*domain.PayrollPeriodLock, error) {
	row, err := s.getQueries(ctx).GetPayrollPeriodLock(ctx, sqlc.GetPayrollPeriodLockParams{TenantID: tenantID, Month: month, Year: year})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInvalidPayrollPeriod
		}
		return nil, s.logDBError(ctx, "get payroll period lock", err, tenantIDField(tenantID))
	}
	return mapPayrollPeriodLock(row), nil
}

func (s *Store) ListPayrollPeriodLocks(ctx context.Context, tenantID uuid.UUID) ([]*domain.PayrollPeriodLock, error) {
	rows, err := s.getQueries(ctx).ListPayrollPeriodLocks(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list payroll period locks", err, tenantIDField(tenantID))
	}
	return mapPayrollPeriodLocks(rows), nil
}

func (s *Store) CreatePayrollPeriodLockEvent(ctx context.Context, item *domain.PayrollPeriodLockEvent, actorID *uuid.UUID) (*domain.PayrollPeriodLockEvent, error) {
	row, err := s.getQueries(ctx).CreatePayrollPeriodLockEvent(ctx, sqlc.CreatePayrollPeriodLockEventParams{TenantID: item.TenantID, PayrollLockID: item.PayrollLockID, Action: item.Action, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), Remarks: textFromPtr(item.Remarks), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create payroll lock event", err, tenantIDField(item.TenantID), stringField("payroll_lock_id", item.PayrollLockID.String()))
	}
	return mapPayrollPeriodLockEvent(row), nil
}

func (s *Store) ListPayrollPeriodLockEvents(ctx context.Context, tenantID uuid.UUID, lockID uuid.UUID) ([]*domain.PayrollPeriodLockEvent, error) {
	rows, err := s.getQueries(ctx).ListPayrollPeriodLockEvents(ctx, sqlc.ListPayrollPeriodLockEventsParams{TenantID: tenantID, PayrollLockID: lockID})
	if err != nil {
		return nil, s.logDBError(ctx, "list payroll lock events", err, tenantIDField(tenantID), stringField("payroll_lock_id", lockID.String()))
	}
	return mapPayrollPeriodLockEvents(rows), nil
}

func (s *Store) ListPayrollStatutoryRules(ctx context.Context, tenantID uuid.UUID, ruleType *string) ([]*domain.PayrollStatutoryRule, error) {
	rows, err := s.getQueries(ctx).ListPayrollStatutoryRules(ctx, sqlc.ListPayrollStatutoryRulesParams{TenantID: tenantID, RuleType: textFromPtr(ruleType)})
	if err != nil {
		return nil, s.logDBError(ctx, "list payroll statutory rules", err, tenantIDField(tenantID))
	}
	return mapPayrollStatutoryRules(rows), nil
}

func (s *Store) GetPayrollStatutoryRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayrollStatutoryRule, error) {
	row, err := s.getQueries(ctx).GetPayrollStatutoryRule(ctx, sqlc.GetPayrollStatutoryRuleParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPayrollStatutoryRuleNotFound
		}
		return nil, s.logDBError(ctx, "get payroll statutory rule", err, tenantIDField(tenantID), stringField("rule_id", id.String()))
	}
	return mapPayrollStatutoryRule(row), nil
}

func (s *Store) CreatePayrollStatutoryRule(ctx context.Context, item *domain.PayrollStatutoryRule, actorID *uuid.UUID) (*domain.PayrollStatutoryRule, error) {
	row, err := s.getQueries(ctx).CreatePayrollStatutoryRule(ctx, payrollStatutoryRuleCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create payroll statutory rule", err, tenantIDField(item.TenantID), stringField("rule_type", item.RuleType))
	}
	return mapPayrollStatutoryRule(row), nil
}

func (s *Store) UpdatePayrollStatutoryRule(ctx context.Context, item *domain.PayrollStatutoryRule, actorID *uuid.UUID) (*domain.PayrollStatutoryRule, error) {
	row, err := s.getQueries(ctx).UpdatePayrollStatutoryRule(ctx, sqlc.UpdatePayrollStatutoryRuleParams{TenantID: item.TenantID, ID: item.ID, RuleType: item.RuleType, Name: item.Name, State: textFromPtr(item.State), BranchID: uuidFromPtr(item.BranchID), EffectiveFrom: dateFromTime(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), MinGrossSalary: numericFromFloatPtr(item.MinGrossSalary), MaxGrossSalary: numericFromFloatPtr(item.MaxGrossSalary), EmployeeAmount: numericFromFloat(item.EmployeeAmount), EmployerAmount: numericFromFloat(item.EmployerAmount), Frequency: item.Frequency, DeductionMonth: int4FromPtr(item.DeductionMonth), Notes: textFromPtr(item.Notes), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update payroll statutory rule", err, tenantIDField(item.TenantID), stringField("rule_id", item.ID.String()))
	}
	return mapPayrollStatutoryRule(row), nil
}

func payrollStatutoryRuleCreateParams(item *domain.PayrollStatutoryRule, actorID *uuid.UUID) sqlc.CreatePayrollStatutoryRuleParams {
	return sqlc.CreatePayrollStatutoryRuleParams{TenantID: item.TenantID, RuleType: item.RuleType, Name: item.Name, State: textFromPtr(item.State), BranchID: uuidFromPtr(item.BranchID), EffectiveFrom: dateFromTime(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), MinGrossSalary: numericFromFloatPtr(item.MinGrossSalary), MaxGrossSalary: numericFromFloatPtr(item.MaxGrossSalary), EmployeeAmount: numericFromFloat(item.EmployeeAmount), EmployerAmount: numericFromFloat(item.EmployerAmount), Frequency: item.Frequency, DeductionMonth: int4FromPtr(item.DeductionMonth), Notes: textFromPtr(item.Notes), CreatedBy: uuidFromPtr(actorID)}
}

func (s *Store) DeletePayrollStatutoryRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePayrollStatutoryRule(ctx, sqlc.SoftDeletePayrollStatutoryRuleParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete payroll statutory rule", err, tenantIDField(tenantID), stringField("rule_id", id.String()))
	}
	return nil
}

func (s *Store) ResolvePayrollStatutoryRule(ctx context.Context, tenantID uuid.UUID, ruleType string, state *string, branchID *uuid.UUID, effectiveDate string, grossSalary float64, month int32) (*domain.PayrollStatutoryRule, error) {
	row, err := s.getQueries(ctx).ResolvePayrollStatutoryRule(ctx, sqlc.ResolvePayrollStatutoryRuleParams{TenantID: tenantID, RuleType: ruleType, BranchID: uuidFromPtr(branchID), State: textFromPtr(state), EffectiveFrom: dateFromString(effectiveDate), MinGrossSalary: numericFromFloat(grossSalary), DeductionMonth: int4FromPtr(&month)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPayrollStatutoryRuleNotFound
		}
		return nil, s.logDBError(ctx, "resolve payroll statutory rule", err, tenantIDField(tenantID), stringField("rule_type", ruleType))
	}
	return mapPayrollStatutoryRule(row), nil
}

func (s *Store) CreatePayrollImportBatch(ctx context.Context, item *domain.PayrollImportBatch, actorID *uuid.UUID) (*domain.PayrollImportBatch, error) {
	row, err := s.getQueries(ctx).CreatePayrollImportBatch(ctx, sqlc.CreatePayrollImportBatchParams{TenantID: item.TenantID, ImportType: item.ImportType, Month: int4FromPtr(item.Month), Year: int4FromPtr(item.Year), FyID: uuidFromPtr(item.FYID), TemplateID: uuidFromPtr(item.TemplateID), FileName: textFromPtr(item.FileName), Status: item.Status, TotalRows: item.TotalRows, ValidRows: item.ValidRows, InvalidRows: item.InvalidRows, AppliedRows: item.AppliedRows, ErrorReport: []byte(item.ErrorReport), Notes: textFromPtr(item.Notes), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create payroll import batch", err, tenantIDField(item.TenantID))
	}
	return mapPayrollImportBatch(row), nil
}

func (s *Store) ListPayrollImportBatches(ctx context.Context, tenantID uuid.UUID, limit int32, offset int32) ([]*domain.PayrollImportBatch, error) {
	rows, err := s.getQueries(ctx).ListPayrollImportBatches(ctx, sqlc.ListPayrollImportBatchesParams{TenantID: tenantID, Limit: limit, Offset: offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list payroll import batches", err, tenantIDField(tenantID))
	}
	return mapPayrollImportBatches(rows), nil
}

func (s *Store) GetPayrollImportBatch(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayrollImportBatch, error) {
	row, err := s.getQueries(ctx).GetPayrollImportBatch(ctx, sqlc.GetPayrollImportBatchParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPayrollImportNotFound
		}
		return nil, s.logDBError(ctx, "get payroll import batch", err, tenantIDField(tenantID), stringField("batch_id", id.String()))
	}
	return mapPayrollImportBatch(row), nil
}

func (s *Store) CreatePayrollImportRow(ctx context.Context, item *domain.PayrollImportRow, actorID *uuid.UUID) (*domain.PayrollImportRow, error) {
	raw := []byte(item.RawData)
	if len(raw) == 0 {
		raw = []byte(`{}`)
	}
	row, err := s.getQueries(ctx).CreatePayrollImportRow(ctx, sqlc.CreatePayrollImportRowParams{TenantID: item.TenantID, BatchID: item.BatchID, RowNumber: item.RowNumber, EmployeeCode: textFromPtr(item.EmployeeCode), EmployeeUserID: uuidFromPtr(item.EmployeeUserID), EmployeeName: textFromPtr(item.EmployeeName), GrossSalary: numericFromFloatPtr(item.GrossSalary), PresentDays: numericFromFloatPtr(item.PresentDays), AbsentDays: numericFromFloatPtr(item.AbsentDays), LopDays: numericFromFloatPtr(item.LOPDays), VariableEarnings: numericFromFloatPtr(item.VariableEarnings), VariableDeductions: numericFromFloatPtr(item.VariableDeductions), Status: item.Status, ErrorMessage: textFromPtr(item.ErrorMessage), RawData: raw, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create payroll import row", err, tenantIDField(item.TenantID), stringField("batch_id", item.BatchID.String()))
	}
	return mapPayrollImportRow(row), nil
}

func (s *Store) ListPayrollImportRows(ctx context.Context, tenantID uuid.UUID, batchID uuid.UUID) ([]*domain.PayrollImportRow, error) {
	rows, err := s.getQueries(ctx).ListPayrollImportRows(ctx, sqlc.ListPayrollImportRowsParams{TenantID: tenantID, BatchID: batchID})
	if err != nil {
		return nil, s.logDBError(ctx, "list payroll import rows", err, tenantIDField(tenantID), stringField("batch_id", batchID.String()))
	}
	return mapPayrollImportRows(rows), nil
}

func (s *Store) ListConsolidatedSalarySheet(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]*domain.ConsolidatedSalarySheetRow, error) {
	rows, err := s.getQueries(ctx).ListConsolidatedSalarySheet(ctx, sqlc.ListConsolidatedSalarySheetParams{TenantID: tenantID, Month: month, Year: year})
	if err != nil {
		return nil, s.logDBError(ctx, "list consolidated salary sheet", err, tenantIDField(tenantID))
	}
	return mapConsolidatedSalarySheetRows(rows), nil
}

func (s *Store) ListPayrollReconciliationRows(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]*domain.PayrollReconciliationRow, error) {
	rows, err := s.getQueries(ctx).ListPayrollReconciliationRows(ctx, sqlc.ListPayrollReconciliationRowsParams{TenantID: tenantID, Month: month, Year: year})
	if err != nil {
		return nil, s.logDBError(ctx, "list payroll reconciliation rows", err, tenantIDField(tenantID))
	}
	return mapPayrollReconciliationRows(rows), nil
}

func dateFromString(value string) pgtype.Date {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return pgtype.Date{Valid: false}
	}
	return dateFromTime(parsed)
}

func marshalPayrollImportErrors(rows []*domain.PayrollImportRow) json.RawMessage {
	type rowError struct {
		Row     int32   `json:"row"`
		Code    *string `json:"employee_code,omitempty"`
		Message *string `json:"message,omitempty"`
	}
	errors := make([]rowError, 0)
	for _, row := range rows {
		if row.Status == domain.PayrollImportRowInvalid {
			errors = append(errors, rowError{Row: row.RowNumber, Code: row.EmployeeCode, Message: row.ErrorMessage})
		}
	}
	data, _ := json.Marshal(errors)
	return data
}
