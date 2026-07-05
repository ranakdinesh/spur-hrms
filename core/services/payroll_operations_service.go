package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) UpsertPayrollPeriodLock(ctx context.Context, cmd ports.PayrollPeriodLockCommand) (*domain.PayrollPeriodLock, error) {
	if err := domain.ValidatePayrollPeriod(cmd.Month, cmd.Year); err != nil {
		s.logError("validate payroll period lock", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status, err := domain.ValidatePayrollLockStatus(cmd.Status)
	if err != nil {
		return nil, err
	}
	if status == domain.PayrollLockStatusLocked {
		start := time.Date(int(cmd.Year), time.Month(cmd.Month), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, -1)
		blockers, err := s.attendanceExceptionWorkflows.ListPayrollBlockingAttendanceRequests(ctx, cmd.TenantID, start.Format("2006-01-02"), end.Format("2006-01-02"))
		if err != nil {
			s.logError("check attendance exception payroll blockers", err, serviceTenantIDField(cmd.TenantID))
			return nil, err
		}
		if len(blockers) > 0 {
			return nil, fmt.Errorf("payroll period has %d unresolved attendance exception(s)", len(blockers))
		}
	}
	var beforeStatus *string
	if existing, err := s.payrollOperations.GetPayrollPeriodLock(ctx, cmd.TenantID, cmd.Month, cmd.Year); err == nil {
		beforeStatus = &existing.Status
	}
	item := &domain.PayrollPeriodLock{TenantID: cmd.TenantID, Month: cmd.Month, Year: cmd.Year, Status: status, UnlockReason: cleanCommandString(cmd.UnlockReason), Notes: cleanCommandString(cmd.Notes)}
	var result *domain.PayrollPeriodLock
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		result, err = s.payrollOperations.UpsertPayrollPeriodLock(txCtx, item, cmd.ActorID)
		if err != nil {
			return err
		}
		_, err = s.payrollOperations.CreatePayrollPeriodLockEvent(txCtx, &domain.PayrollPeriodLockEvent{TenantID: cmd.TenantID, PayrollLockID: result.ID, Action: status, FromStatus: beforeStatus, ToStatus: &result.Status, Remarks: item.Notes}, cmd.ActorID)
		return err
	}); err != nil {
		s.logError("upsert payroll period lock", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result.Events, _ = s.payrollOperations.ListPayrollPeriodLockEvents(ctx, cmd.TenantID, result.ID)
	return result, nil
}

func (s *TenantService) ListPayrollPeriodLocks(ctx context.Context, tenantID uuid.UUID) ([]*domain.PayrollPeriodLock, error) {
	items, err := s.payrollOperations.ListPayrollPeriodLocks(ctx, tenantID)
	if err != nil {
		s.logError("list payroll period locks", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	for _, item := range items {
		item.Events, _ = s.payrollOperations.ListPayrollPeriodLockEvents(ctx, tenantID, item.ID)
	}
	return items, nil
}

func (s *TenantService) ensurePayrollPeriodOpen(ctx context.Context, tenantID uuid.UUID, month int32, year int32) error {
	lock, err := s.payrollOperations.GetPayrollPeriodLock(ctx, tenantID, month, year)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPayrollPeriod) {
			return nil
		}
		return err
	}
	if lock.Status == domain.PayrollLockStatusLocked {
		return domain.ErrPayrollPeriodLocked
	}
	return nil
}

func (s *TenantService) CreatePayrollStatutoryRule(ctx context.Context, cmd ports.PayrollStatutoryRuleCommand) (*domain.PayrollStatutoryRule, error) {
	item, err := payrollRuleFromCommand(cmd)
	if err != nil {
		s.logError("validate payroll statutory rule", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.payrollOperations.CreatePayrollStatutoryRule(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdatePayrollStatutoryRule(ctx context.Context, cmd ports.PayrollStatutoryRuleCommand) (*domain.PayrollStatutoryRule, error) {
	item, err := payrollRuleFromCommand(cmd)
	if err != nil {
		s.logError("validate payroll statutory rule update", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item.ID = cmd.ID
	return s.payrollOperations.UpdatePayrollStatutoryRule(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListPayrollStatutoryRules(ctx context.Context, tenantID uuid.UUID, ruleType *string) ([]*domain.PayrollStatutoryRule, error) {
	return s.payrollOperations.ListPayrollStatutoryRules(ctx, tenantID, ruleType)
}

func (s *TenantService) DeletePayrollStatutoryRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	return s.payrollOperations.DeletePayrollStatutoryRule(ctx, tenantID, id, actorID)
}

func payrollRuleFromCommand(cmd ports.PayrollStatutoryRuleCommand) (*domain.PayrollStatutoryRule, error) {
	from, err := parseOptionalDate(cmd.EffectiveFrom)
	if err != nil || from == nil {
		return nil, domain.ErrInvalidPayrollStatutoryRule
	}
	to, err := parseOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, domain.ErrInvalidPayrollStatutoryRule
	}
	return domain.NewPayrollStatutoryRule(domain.PayrollStatutoryRule{TenantID: cmd.TenantID, RuleType: cmd.RuleType, Name: cmd.Name, State: cleanCommandString(cmd.State), BranchID: cmd.BranchID, EffectiveFrom: *from, EffectiveTo: to, MinGrossSalary: cmd.MinGrossSalary, MaxGrossSalary: cmd.MaxGrossSalary, EmployeeAmount: cmd.EmployeeAmount, EmployerAmount: cmd.EmployerAmount, Frequency: cmd.Frequency, DeductionMonth: cmd.DeductionMonth, Notes: cleanCommandString(cmd.Notes)})
}

func (s *TenantService) ImportPayrollData(ctx context.Context, cmd ports.PayrollImportCommand) (*domain.PayrollImportBatch, error) {
	if cmd.TenantID == uuid.Nil || len(cmd.Rows) == 0 {
		return nil, domain.ErrInvalidPayrollPeriod
	}
	if cmd.Month != nil && cmd.Year != nil {
		if err := s.ensurePayrollPeriodOpen(ctx, cmd.TenantID, *cmd.Month, *cmd.Year); err != nil {
			return nil, err
		}
	}
	importRows := make([]*domain.PayrollImportRow, 0, len(cmd.Rows))
	validRows := int32(0)
	appliedRows := int32(0)
	for idx, input := range cmd.Rows {
		row := s.validatePayrollImportRow(ctx, cmd, input, int32(idx+1))
		if row.Status == domain.PayrollImportRowValid {
			validRows++
		}
		importRows = append(importRows, row)
	}
	if cmd.Apply {
		for _, row := range importRows {
			if row.Status != domain.PayrollImportRowValid || row.EmployeeUserID == nil || row.GrossSalary == nil || cmd.FYID == nil || cmd.TemplateID == nil {
				continue
			}
			_, err := s.AssignEmployeeSalary(ctx, ports.EmployeeSalaryCommand{TenantID: cmd.TenantID, UserID: *row.EmployeeUserID, FYID: *cmd.FYID, TemplateID: *cmd.TemplateID, GrossSalary: *row.GrossSalary, ActorID: cmd.ActorID})
			if err != nil {
				message := err.Error()
				row.Status = domain.PayrollImportRowInvalid
				row.ErrorMessage = &message
				validRows--
				continue
			}
			row.Status = domain.PayrollImportRowApplied
			appliedRows++
		}
	}
	status := domain.PayrollImportStatusValidated
	invalidRows := int32(len(importRows)) - validRows
	if cmd.Apply {
		if invalidRows > 0 && appliedRows > 0 {
			status = domain.PayrollImportStatusPartial
		} else if invalidRows > 0 && appliedRows == 0 {
			status = domain.PayrollImportStatusFailed
		} else {
			status = domain.PayrollImportStatusApplied
		}
	}
	batch := &domain.PayrollImportBatch{TenantID: cmd.TenantID, ImportType: payrollFirstClean(cmd.ImportType, domain.PayrollImportSalaryRevision), Month: cmd.Month, Year: cmd.Year, FYID: cmd.FYID, TemplateID: cmd.TemplateID, FileName: cleanCommandString(cmd.FileName), Status: status, TotalRows: int32(len(importRows)), ValidRows: validRows, InvalidRows: invalidRows, AppliedRows: appliedRows, ErrorReport: payrollImportErrorReport(importRows), Notes: cleanCommandString(cmd.Notes)}
	var result *domain.PayrollImportBatch
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		result, err = s.payrollOperations.CreatePayrollImportBatch(txCtx, batch, cmd.ActorID)
		if err != nil {
			return err
		}
		for _, row := range importRows {
			row.BatchID = result.ID
			if _, err := s.payrollOperations.CreatePayrollImportRow(txCtx, row, cmd.ActorID); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		s.logError("import payroll data", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result.Rows, _ = s.payrollOperations.ListPayrollImportRows(ctx, cmd.TenantID, result.ID)
	return result, nil
}

func (s *TenantService) validatePayrollImportRow(ctx context.Context, cmd ports.PayrollImportCommand, input ports.PayrollImportRowCommand, rowNumber int32) *domain.PayrollImportRow {
	raw, _ := json.Marshal(input.RawData)
	row := &domain.PayrollImportRow{TenantID: cmd.TenantID, RowNumber: rowNumber, EmployeeCode: payrollCleanStringPtr(input.EmployeeCode), GrossSalary: input.GrossSalary, PresentDays: input.PresentDays, AbsentDays: input.AbsentDays, LOPDays: input.LOPDays, VariableEarnings: input.VariableEarnings, VariableDeductions: input.VariableDeductions, Status: domain.PayrollImportRowValid, RawData: raw}
	if row.EmployeeCode == nil {
		message := "employee_code is required"
		row.Status = domain.PayrollImportRowInvalid
		row.ErrorMessage = &message
		return row
	}
	employee, err := s.employees.GetEmployeeByCode(ctx, cmd.TenantID, *row.EmployeeCode)
	if err != nil {
		message := "employee_code was not found"
		row.Status = domain.PayrollImportRowInvalid
		row.ErrorMessage = &message
		return row
	}
	row.EmployeeUserID = &employee.UserID
	row.EmployeeName = stringPtr(strings.TrimSpace(employee.Firstname + " " + valueFromPtr(employee.Lastname)))
	if input.GrossSalary != nil && *input.GrossSalary < 0 {
		message := "gross_salary cannot be negative"
		row.Status = domain.PayrollImportRowInvalid
		row.ErrorMessage = &message
	}
	return row
}

func (s *TenantService) ListPayrollImportBatches(ctx context.Context, tenantID uuid.UUID, limit int32, offset int32) ([]*domain.PayrollImportBatch, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	return s.payrollOperations.ListPayrollImportBatches(ctx, tenantID, limit, offset)
}

func (s *TenantService) GetPayrollImportBatch(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayrollImportBatch, error) {
	batch, err := s.payrollOperations.GetPayrollImportBatch(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	batch.Rows, err = s.payrollOperations.ListPayrollImportRows(ctx, tenantID, id)
	return batch, err
}

func (s *TenantService) ListConsolidatedSalarySheet(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]*domain.ConsolidatedSalarySheetRow, error) {
	if err := domain.ValidatePayrollPeriod(month, year); err != nil {
		return nil, err
	}
	return s.payrollOperations.ListConsolidatedSalarySheet(ctx, tenantID, month, year)
}

func (s *TenantService) ExportConsolidatedSalarySheetCSV(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]byte, string, error) {
	rows, err := s.ListConsolidatedSalarySheet(ctx, tenantID, month, year)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	_ = writer.Write([]string{"employee_code", "employee_name", "email", "branch", "department", "month", "year", "gross_salary", "total_earnings", "total_deductions", "absent_deduction", "net_salary", "present_days", "absent_days", "lwp_days", "pdf_path"})
	for _, row := range rows {
		_ = writer.Write([]string{valueFromPtr(row.EmployeeCode), strings.TrimSpace(row.Firstname + " " + valueFromPtr(row.Lastname)), valueFromPtr(row.Email), valueFromPtr(row.BranchName), valueFromPtr(row.DepartmentName), strconv.Itoa(int(row.Month)), strconv.Itoa(int(row.Year)), moneyString(row.GrossSalary), moneyString(row.TotalEarnings), moneyString(row.TotalDeductions), moneyString(row.AbsentDeduction), moneyString(row.NetSalary), strconv.Itoa(int(row.PresentDays)), strconv.Itoa(int(row.AbsentDays)), moneyString(row.LWPDays), valueFromPtr(row.PDFPath)})
	}
	writer.Flush()
	return buf.Bytes(), fmt.Sprintf("consolidated-salary-sheet-%04d-%02d.csv", year, month), writer.Error()
}

func (s *TenantService) ListPayrollReconciliationRows(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]*domain.PayrollReconciliationRow, error) {
	if err := domain.ValidatePayrollPeriod(month, year); err != nil {
		return nil, err
	}
	return s.payrollOperations.ListPayrollReconciliationRows(ctx, tenantID, month, year)
}

func payrollImportErrorReport(rows []*domain.PayrollImportRow) json.RawMessage {
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

func moneyString(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func payrollFirstClean(value string, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func payrollCleanStringPtr(value string) *string {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return nil
	}
	return &cleaned
}
