package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapSalaryTemplate(row sqlc.HrmsSalaryTemplate) *domain.SalaryTemplate {
	return &domain.SalaryTemplate{ID: row.ID, TenantID: row.TenantID, FYID: row.FyID, Code: row.Code, Name: row.Name, Description: ptrFromText(row.Description), TemplateType: row.TemplateType, AppliesTo: row.AppliesTo, CurrencyCode: row.CurrencyCode, EffectiveFrom: ptrFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), Notes: ptrFromText(row.Notes), IsActive: row.IsActive, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapSalaryTemplates(rows []sqlc.HrmsSalaryTemplate) []*domain.SalaryTemplate {
	items := make([]*domain.SalaryTemplate, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapSalaryTemplate(row))
	}
	return items
}

func mapSalaryTemplateItem(row sqlc.HrmsSalaryTemplateItem) *domain.SalaryTemplateItem {
	return &domain.SalaryTemplateItem{ID: row.ID, TenantID: row.TenantID, TemplateID: row.TemplateID, ItemType: row.ItemType, Code: row.Code, Name: row.Name, Percentage: floatPtrFromNumeric(row.Percentage), Amount: floatPtrFromNumeric(row.Amount), CalculationMode: row.CalculationMode, CalculationBase: row.CalculationBase, Formula: ptrFromText(row.Formula), ContributionSide: row.ContributionSide, IsTaxExempt: row.IsTaxExempt, IsStatutory: row.IsStatutory, IsVariable: row.IsVariable, AffectsGross: row.AffectsGross, AffectsNet: row.AffectsNet, CapAmount: floatPtrFromNumeric(row.CapAmount), MinAmount: floatPtrFromNumeric(row.MinAmount), MaxAmount: floatPtrFromNumeric(row.MaxAmount), SortOrder: row.SortOrder, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapSalaryTemplateItems(rows []sqlc.HrmsSalaryTemplateItem) []*domain.SalaryTemplateItem {
	items := make([]*domain.SalaryTemplateItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapSalaryTemplateItem(row))
	}
	return items
}
