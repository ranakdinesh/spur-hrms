package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapPayrollPeriodLock(row sqlc.HrmsPayrollPeriodLock) *domain.PayrollPeriodLock {
	return &domain.PayrollPeriodLock{ID: row.ID, TenantID: row.TenantID, Month: row.Month, Year: row.Year, Status: row.Status, LockedAt: ptrFromTimestamptz(row.LockedAt), LockedBy: ptrFromUUID(row.LockedBy), UnlockedAt: ptrFromTimestamptz(row.UnlockedAt), UnlockedBy: ptrFromUUID(row.UnlockedBy), UnlockReason: ptrFromText(row.UnlockReason), Notes: ptrFromText(row.Notes), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapPayrollPeriodLocks(rows []sqlc.HrmsPayrollPeriodLock) []*domain.PayrollPeriodLock {
	items := make([]*domain.PayrollPeriodLock, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayrollPeriodLock(row))
	}
	return items
}

func mapPayrollPeriodLockEvent(row sqlc.HrmsPayrollPeriodLockEvent) *domain.PayrollPeriodLockEvent {
	return &domain.PayrollPeriodLockEvent{ID: row.ID, TenantID: row.TenantID, PayrollLockID: row.PayrollLockID, Action: row.Action, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), Remarks: ptrFromText(row.Remarks), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapPayrollPeriodLockEvents(rows []sqlc.HrmsPayrollPeriodLockEvent) []*domain.PayrollPeriodLockEvent {
	items := make([]*domain.PayrollPeriodLockEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayrollPeriodLockEvent(row))
	}
	return items
}

func mapPayrollStatutoryRule(row sqlc.HrmsPayrollStatutoryRule) *domain.PayrollStatutoryRule {
	return &domain.PayrollStatutoryRule{ID: row.ID, TenantID: row.TenantID, RuleType: row.RuleType, Name: row.Name, State: ptrFromText(row.State), BranchID: ptrFromUUID(row.BranchID), EffectiveFrom: timeFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), MinGrossSalary: ptrFromNumeric(row.MinGrossSalary), MaxGrossSalary: ptrFromNumeric(row.MaxGrossSalary), EmployeeAmount: floatFromNumeric(row.EmployeeAmount), EmployerAmount: floatFromNumeric(row.EmployerAmount), Frequency: row.Frequency, DeductionMonth: ptrFromInt4(row.DeductionMonth), Notes: ptrFromText(row.Notes), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapPayrollStatutoryRules(rows []sqlc.HrmsPayrollStatutoryRule) []*domain.PayrollStatutoryRule {
	items := make([]*domain.PayrollStatutoryRule, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayrollStatutoryRule(row))
	}
	return items
}

func mapPayrollImportBatch(row sqlc.HrmsPayrollImportBatch) *domain.PayrollImportBatch {
	return &domain.PayrollImportBatch{ID: row.ID, TenantID: row.TenantID, ImportType: row.ImportType, Month: ptrFromInt4(row.Month), Year: ptrFromInt4(row.Year), FYID: ptrFromUUID(row.FyID), TemplateID: ptrFromUUID(row.TemplateID), FileName: ptrFromText(row.FileName), Status: row.Status, TotalRows: row.TotalRows, ValidRows: row.ValidRows, InvalidRows: row.InvalidRows, AppliedRows: row.AppliedRows, ErrorReport: json.RawMessage(row.ErrorReport), Notes: ptrFromText(row.Notes), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapPayrollImportBatches(rows []sqlc.HrmsPayrollImportBatch) []*domain.PayrollImportBatch {
	items := make([]*domain.PayrollImportBatch, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayrollImportBatch(row))
	}
	return items
}

func mapPayrollImportRow(row sqlc.HrmsPayrollImportRow) *domain.PayrollImportRow {
	return &domain.PayrollImportRow{ID: row.ID, TenantID: row.TenantID, BatchID: row.BatchID, RowNumber: row.RowNumber, EmployeeCode: ptrFromText(row.EmployeeCode), EmployeeUserID: ptrFromUUID(row.EmployeeUserID), EmployeeName: ptrFromText(row.EmployeeName), GrossSalary: ptrFromNumeric(row.GrossSalary), PresentDays: ptrFromNumeric(row.PresentDays), AbsentDays: ptrFromNumeric(row.AbsentDays), LOPDays: ptrFromNumeric(row.LopDays), VariableEarnings: ptrFromNumeric(row.VariableEarnings), VariableDeductions: ptrFromNumeric(row.VariableDeductions), Status: row.Status, ErrorMessage: ptrFromText(row.ErrorMessage), RawData: json.RawMessage(row.RawData), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapPayrollImportRows(rows []sqlc.HrmsPayrollImportRow) []*domain.PayrollImportRow {
	items := make([]*domain.PayrollImportRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayrollImportRow(row))
	}
	return items
}

func mapConsolidatedSalarySheetRow(row sqlc.ListConsolidatedSalarySheetRow) *domain.ConsolidatedSalarySheetRow {
	return &domain.ConsolidatedSalarySheetRow{SalarySlipID: row.ID, TenantID: row.TenantID, UserID: row.UserID, EmployeeCode: ptrFromText(row.EmployeeCode), Firstname: row.Firstname, Lastname: ptrFromText(row.Lastname), Email: ptrFromText(row.Email), BranchName: ptrFromText(row.BranchName), DepartmentName: ptrFromText(row.DepartmentName), Month: row.Month, Year: row.Year, GrossSalary: floatFromNumeric(row.GrossSalary), TotalEarnings: floatFromNumeric(row.TotalEarnings), TotalDeductions: floatFromNumeric(row.TotalDeductions), AbsentDeduction: floatFromNumeric(row.AbsentDeduction), NetSalary: floatFromNumeric(row.NetSalary), PresentDays: row.PresentDays, AbsentDays: row.AbsentDays, LWPDays: floatFromNumeric(row.LwpDays), PDFPath: ptrFromText(row.PdfPath), CreatedAt: timeFromTimestamptz(row.CreatedAt)}
}

func mapConsolidatedSalarySheetRows(rows []sqlc.ListConsolidatedSalarySheetRow) []*domain.ConsolidatedSalarySheetRow {
	items := make([]*domain.ConsolidatedSalarySheetRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapConsolidatedSalarySheetRow(row))
	}
	return items
}

func mapPayrollReconciliationRow(row sqlc.ListPayrollReconciliationRowsRow) *domain.PayrollReconciliationRow {
	return &domain.PayrollReconciliationRow{EmployeeID: row.EmployeeID, UserID: row.UserID, EmployeeCode: ptrFromText(row.EmployeeCode), Firstname: row.Firstname, Lastname: ptrFromText(row.Lastname), Email: ptrFromText(row.Email), BranchName: ptrFromText(row.BranchName), DepartmentName: ptrFromText(row.DepartmentName), SalarySlipID: ptrFromUUID(row.SalarySlipID), PresentDays: ptrFromInt4(row.PresentDays), AbsentDays: ptrFromInt4(row.AbsentDays), LWPDays: ptrFromNumeric(row.LwpDays), NetSalary: ptrFromNumeric(row.NetSalary), ReconciliationStatus: row.ReconciliationStatus}
}

func mapPayrollReconciliationRows(rows []sqlc.ListPayrollReconciliationRowsRow) []*domain.PayrollReconciliationRow {
	items := make([]*domain.PayrollReconciliationRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPayrollReconciliationRow(row))
	}
	return items
}
