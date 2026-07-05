package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapEmployeeSalary(row sqlc.HrmsEmployeeSalary) *domain.EmployeeSalary {
	return &domain.EmployeeSalary{ID: row.ID, TenantID: row.TenantID, UserID: row.UserID, FYID: row.FyID, TemplateID: row.TemplateID, GrossSalary: floatFromNumeric(row.GrossSalary), EffectiveFrom: ptrFromDate(row.EffectiveFrom), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapEmployeeSalaries(rows []sqlc.HrmsEmployeeSalary) []*domain.EmployeeSalary {
	items := make([]*domain.EmployeeSalary, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeSalary(row))
	}
	return items
}

func mapEmployeeSalaryStructure(row sqlc.HrmsEmployeeSalaryStructure) *domain.EmployeeSalaryStructure {
	return &domain.EmployeeSalaryStructure{ID: row.ID, TenantID: row.TenantID, UserID: row.UserID, TemplateID: row.TemplateID, FYID: row.FyID, ItemType: row.ItemType, Code: row.Code, Name: row.Name, Amount: floatFromNumeric(row.Amount), SortOrder: row.SortOrder, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapEmployeeSalaryStructures(rows []sqlc.HrmsEmployeeSalaryStructure) []*domain.EmployeeSalaryStructure {
	items := make([]*domain.EmployeeSalaryStructure, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeSalaryStructure(row))
	}
	return items
}
