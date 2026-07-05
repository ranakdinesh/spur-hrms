package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapSalarySlip(row sqlc.HrmsSalarySlip) *domain.SalarySlip {
	return &domain.SalarySlip{ID: row.ID, TenantID: row.TenantID, UserID: row.UserID, FYID: row.FyID, TemplateID: row.TemplateID, Month: row.Month, Year: row.Year, GrossSalary: floatFromNumeric(row.GrossSalary), TotalEarnings: floatFromNumeric(row.TotalEarnings), TotalDeductions: floatFromNumeric(row.TotalDeductions), AbsentDeduction: floatFromNumeric(row.AbsentDeduction), NetSalary: floatFromNumeric(row.NetSalary), AbsentDays: row.AbsentDays, PresentDays: row.PresentDays, TotalDays: row.TotalDays, LWPDays: floatFromNumeric(row.LwpDays), NoOfPHWEO: row.NoOfPhWeo, IsSpecial: row.IsSpecial, IsRegenerated: row.IsRegenerated, PDFPath: ptrFromText(row.PdfPath), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapSalarySlips(rows []sqlc.HrmsSalarySlip) []*domain.SalarySlip {
	items := make([]*domain.SalarySlip, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapSalarySlip(row))
	}
	return items
}

func mapSalarySlipItem(row sqlc.HrmsSalarySlipItem) *domain.SalarySlipItem {
	return &domain.SalarySlipItem{ID: row.ID, TenantID: row.TenantID, SlipID: row.SlipID, ItemType: row.ItemType, Code: row.Code, Name: row.Name, Amount: floatFromNumeric(row.Amount), SortOrder: row.SortOrder, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapSalarySlipItems(rows []sqlc.HrmsSalarySlipItem) []*domain.SalarySlipItem {
	items := make([]*domain.SalarySlipItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapSalarySlipItem(row))
	}
	return items
}

func mapSalarySlipLeave(row sqlc.HrmsSalarySlipLeafe) *domain.SalarySlipLeave {
	return &domain.SalarySlipLeave{ID: row.ID, TenantID: row.TenantID, SlipID: row.SlipID, LeaveTypeID: row.LeaveTypeID, LeaveTypeName: ptrFromText(row.LeaveTypeName), TotalDays: floatFromNumeric(row.TotalDays), UsedDays: floatFromNumeric(row.UsedDays), BalanceDays: floatFromNumeric(row.BalanceDays), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapSalarySlipLeaves(rows []sqlc.HrmsSalarySlipLeafe) []*domain.SalarySlipLeave {
	items := make([]*domain.SalarySlipLeave, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapSalarySlipLeave(row))
	}
	return items
}

func mapSalarySlipFormat(row sqlc.HrmsSalarySlipFormat) (*domain.SalarySlipFormat, error) {
	fields := map[string]any{}
	if len(row.CustomFields) > 0 {
		if err := json.Unmarshal(row.CustomFields, &fields); err != nil {
			return nil, err
		}
	}
	return &domain.SalarySlipFormat{ID: row.ID, TenantID: row.TenantID, Title: row.Title, Subtitle: ptrFromText(row.Subtitle), LogoPath: ptrFromText(row.LogoPath), PrimaryColor: row.PrimaryColor, AccentColor: row.AccentColor, ShowLeaveBalance: row.ShowLeaveBalance, ShowYTDSummary: row.ShowYtdSummary, ShowEmployeeBank: row.ShowEmployeeBank, ShowEmployerContributions: row.ShowEmployerContributions, FooterText: ptrFromText(row.FooterText), CustomFields: fields, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}, nil
}
