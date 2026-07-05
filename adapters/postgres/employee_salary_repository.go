package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateEmployeeSalary(ctx context.Context, item *domain.EmployeeSalary, actorID *uuid.UUID) (*domain.EmployeeSalary, error) {
	row, err := s.getQueries(ctx).CreateEmployeeSalary(ctx, sqlc.CreateEmployeeSalaryParams{TenantID: item.TenantID, UserID: item.UserID, FyID: item.FYID, TemplateID: item.TemplateID, GrossSalary: numericFromFloat(item.GrossSalary), EffectiveFrom: dateFromPtr(item.EffectiveFrom), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee salary", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()), stringField("financial_year_id", item.FYID.String()))
	}
	return mapEmployeeSalary(row), nil
}

func (s *Store) UpdateEmployeeSalary(ctx context.Context, item *domain.EmployeeSalary, actorID *uuid.UUID) (*domain.EmployeeSalary, error) {
	row, err := s.getQueries(ctx).UpdateEmployeeSalary(ctx, sqlc.UpdateEmployeeSalaryParams{TenantID: item.TenantID, UserID: item.UserID, FyID: item.FYID, TemplateID: item.TemplateID, GrossSalary: numericFromFloat(item.GrossSalary), EffectiveFrom: dateFromPtr(item.EffectiveFrom), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update employee salary", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()), stringField("financial_year_id", item.FYID.String()))
	}
	return mapEmployeeSalary(row), nil
}

func (s *Store) ListEmployeeSalariesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.EmployeeSalary, error) {
	rows, err := s.getQueries(ctx).ListEmployeeSalariesByUser(ctx, sqlc.ListEmployeeSalariesByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee salaries", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapEmployeeSalaries(rows), nil
}

func (s *Store) GetEmployeeSalary(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeSalary, error) {
	row, err := s.getQueries(ctx).GetEmployeeSalary(ctx, sqlc.GetEmployeeSalaryParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get employee salary", err, tenantIDField(tenantID), stringField("employee_salary_id", id.String()))
	}
	return mapEmployeeSalary(row), nil
}

func (s *Store) DeleteEmployeeSalary(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmployeeSalary(ctx, sqlc.SoftDeleteEmployeeSalaryParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete employee salary", err, tenantIDField(tenantID), stringField("employee_salary_id", id.String()))
	}
	return nil
}

func (s *Store) DeleteEmployeeSalariesByUserFY(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, fyID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmployeeSalariesByUserFY(ctx, sqlc.SoftDeleteEmployeeSalariesByUserFYParams{TenantID: tenantID, UserID: userID, FyID: fyID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete employee salaries by user fy", err, tenantIDField(tenantID), stringField("user_id", userID.String()), stringField("financial_year_id", fyID.String()))
	}
	return nil
}

func (s *Store) CreateEmployeeSalaryStructure(ctx context.Context, item *domain.EmployeeSalaryStructure, actorID *uuid.UUID) (*domain.EmployeeSalaryStructure, error) {
	row, err := s.getQueries(ctx).CreateEmployeeSalaryStructure(ctx, sqlc.CreateEmployeeSalaryStructureParams{TenantID: item.TenantID, UserID: item.UserID, TemplateID: item.TemplateID, FyID: item.FYID, ItemType: item.ItemType, Code: item.Code, Name: item.Name, Amount: numericFromFloat(item.Amount), SortOrder: item.SortOrder, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee salary structure", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()), stringField("salary_template_id", item.TemplateID.String()))
	}
	return mapEmployeeSalaryStructure(row), nil
}

func (s *Store) ListEmployeeSalaryStructures(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, fyID uuid.UUID) ([]*domain.EmployeeSalaryStructure, error) {
	rows, err := s.getQueries(ctx).ListEmployeeSalaryStructures(ctx, sqlc.ListEmployeeSalaryStructuresParams{TenantID: tenantID, UserID: userID, FyID: fyID})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee salary structures", err, tenantIDField(tenantID), stringField("user_id", userID.String()), stringField("financial_year_id", fyID.String()))
	}
	return mapEmployeeSalaryStructures(rows), nil
}

func (s *Store) DeleteEmployeeSalaryStructuresByUserFY(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, fyID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmployeeSalaryStructuresByUserFY(ctx, sqlc.SoftDeleteEmployeeSalaryStructuresByUserFYParams{TenantID: tenantID, UserID: userID, FyID: fyID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete employee salary structures by user fy", err, tenantIDField(tenantID), stringField("user_id", userID.String()), stringField("financial_year_id", fyID.String()))
	}
	return nil
}
