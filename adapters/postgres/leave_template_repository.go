package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateLeavePolicyTemplate(ctx context.Context, item *domain.LeavePolicyTemplate, actorID *uuid.UUID) (*domain.LeavePolicyTemplate, error) {
	row, err := s.getQueries(ctx).CreateLeavePolicyTemplate(ctx, sqlc.CreateLeavePolicyTemplateParams{
		TenantID:      uuidFromPtr(item.TenantID),
		Name:          item.Name,
		Code:          item.Code,
		Description:   textFromPtr(item.Description),
		EffectiveFrom: dateFromPtr(item.EffectiveFrom),
		EffectiveTo:   dateFromPtr(item.EffectiveTo),
		CreatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create leave policy template", err, optionalTenantIDField(item.TenantID), stringField("template_code", item.Code))
	}
	return mapLeavePolicyTemplate(row), nil
}

func (s *Store) ListLeavePolicyTemplates(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeavePolicyTemplate, error) {
	rows, err := s.getQueries(ctx).ListLeavePolicyTemplates(ctx, uuidFromPtr(&tenantID))
	if err != nil {
		return nil, s.logDBError(ctx, "list leave policy templates", err, tenantIDField(tenantID))
	}
	return mapLeavePolicyTemplates(rows), nil
}

func (s *Store) GetLeavePolicyTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeavePolicyTemplate, error) {
	row, err := s.getQueries(ctx).GetLeavePolicyTemplate(ctx, sqlc.GetLeavePolicyTemplateParams{TenantID: uuidFromPtr(&tenantID), ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveTemplateNotFound
		}
		return nil, s.logDBError(ctx, "get leave policy template", err, tenantIDField(tenantID), stringField("template_id", id.String()))
	}
	return mapLeavePolicyTemplate(row), nil
}

func (s *Store) UpdateLeavePolicyTemplate(ctx context.Context, item *domain.LeavePolicyTemplate, actorID *uuid.UUID) (*domain.LeavePolicyTemplate, error) {
	if item.TenantID == nil {
		return nil, domain.ErrInvalidTenantID
	}
	row, err := s.getQueries(ctx).UpdateLeavePolicyTemplate(ctx, sqlc.UpdateLeavePolicyTemplateParams{
		TenantID:      uuidFromPtr(item.TenantID),
		ID:            item.ID,
		Name:          item.Name,
		Code:          item.Code,
		Description:   textFromPtr(item.Description),
		EffectiveFrom: dateFromPtr(item.EffectiveFrom),
		EffectiveTo:   dateFromPtr(item.EffectiveTo),
		UpdatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveTemplateNotFound
		}
		return nil, s.logDBError(ctx, "update leave policy template", err, optionalTenantIDField(item.TenantID), stringField("template_id", item.ID.String()))
	}
	return mapLeavePolicyTemplate(row), nil
}

func (s *Store) DeleteLeavePolicyTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLeavePolicyTemplate(ctx, sqlc.SoftDeleteLeavePolicyTemplateParams{TenantID: uuidFromPtr(&tenantID), ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete leave policy template", err, tenantIDField(tenantID), stringField("template_id", id.String()))
	}
	return nil
}

func (s *Store) CreateLeavePolicyTemplateRule(ctx context.Context, item *domain.LeavePolicyTemplateRule, actorID *uuid.UUID) (*domain.LeavePolicyTemplateRule, error) {
	row, err := s.getQueries(ctx).CreateLeavePolicyTemplateRule(ctx, leaveTemplateRuleCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create leave policy template rule", err, tenantIDField(item.TenantID), stringField("template_id", item.TemplateID.String()))
	}
	mapped, err := mapLeavePolicyTemplateRule(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map leave policy template rule", fmt.Errorf("hrms: map leave policy template rule: %w", err), tenantIDField(item.TenantID))
	}
	return mapped, nil
}

func (s *Store) ListLeavePolicyTemplateRules(ctx context.Context, tenantID uuid.UUID, templateID uuid.UUID) ([]*domain.LeavePolicyTemplateRule, error) {
	rows, err := s.getQueries(ctx).ListLeavePolicyTemplateRules(ctx, sqlc.ListLeavePolicyTemplateRulesParams{TenantID: tenantID, TemplateID: templateID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave policy template rules", err, tenantIDField(tenantID), stringField("template_id", templateID.String()))
	}
	items, err := mapLeavePolicyTemplateRules(rows)
	if err != nil {
		return nil, s.logDBError(ctx, "map leave policy template rules", fmt.Errorf("hrms: map leave policy template rules: %w", err), tenantIDField(tenantID))
	}
	return items, nil
}

func (s *Store) ListLeavePolicyTemplateRulesByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeavePolicyTemplateRule, error) {
	rows, err := s.getQueries(ctx).ListLeavePolicyTemplateRulesByTenant(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list leave policy template rules by tenant", err, tenantIDField(tenantID))
	}
	items, err := mapLeavePolicyTemplateRules(rows)
	if err != nil {
		return nil, s.logDBError(ctx, "map tenant leave policy template rules", fmt.Errorf("hrms: map tenant leave policy template rules: %w", err), tenantIDField(tenantID))
	}
	return items, nil
}

func (s *Store) GetLeavePolicyTemplateRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeavePolicyTemplateRule, error) {
	row, err := s.getQueries(ctx).GetLeavePolicyTemplateRule(ctx, sqlc.GetLeavePolicyTemplateRuleParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveTemplateRuleNotFound
		}
		return nil, s.logDBError(ctx, "get leave policy template rule", err, tenantIDField(tenantID), stringField("rule_id", id.String()))
	}
	item, err := mapLeavePolicyTemplateRule(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map leave policy template rule", fmt.Errorf("hrms: map leave policy template rule: %w", err), tenantIDField(tenantID), stringField("rule_id", id.String()))
	}
	return item, nil
}

func (s *Store) UpdateLeavePolicyTemplateRule(ctx context.Context, item *domain.LeavePolicyTemplateRule, actorID *uuid.UUID) (*domain.LeavePolicyTemplateRule, error) {
	params := leaveTemplateRuleUpdateParams(item, actorID)
	row, err := s.getQueries(ctx).UpdateLeavePolicyTemplateRule(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLeaveTemplateRuleNotFound
		}
		return nil, s.logDBError(ctx, "update leave policy template rule", err, tenantIDField(item.TenantID), stringField("rule_id", item.ID.String()))
	}
	mapped, err := mapLeavePolicyTemplateRule(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map updated leave policy template rule", fmt.Errorf("hrms: map updated leave policy template rule: %w", err), tenantIDField(item.TenantID))
	}
	return mapped, nil
}

func (s *Store) DeleteLeavePolicyTemplateRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteLeavePolicyTemplateRule(ctx, sqlc.SoftDeleteLeavePolicyTemplateRuleParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete leave policy template rule", err, tenantIDField(tenantID), stringField("rule_id", id.String()))
	}
	return nil
}

func (s *Store) UpsertEmployeeLeavePolicyAssignment(ctx context.Context, item *domain.EmployeeLeavePolicyAssignment, actorID *uuid.UUID) (*domain.EmployeeLeavePolicyAssignment, error) {
	row, err := s.getQueries(ctx).UpsertEmployeeLeavePolicyAssignment(ctx, sqlc.UpsertEmployeeLeavePolicyAssignmentParams{
		TenantID:      item.TenantID,
		UserID:        item.UserID,
		TemplateID:    item.TemplateID,
		FyID:          uuidFromPtr(item.FYID),
		EffectiveFrom: dateFromTime(item.EffectiveFrom),
		EffectiveTo:   dateFromPtr(item.EffectiveTo),
		CreatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert employee leave policy assignment", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapEmployeeLeavePolicyAssignment(row), nil
}

func (s *Store) ListEmployeeLeavePolicyAssignments(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.EmployeeLeavePolicyAssignment, error) {
	rows, err := s.getQueries(ctx).ListEmployeeLeavePolicyAssignments(ctx, sqlc.ListEmployeeLeavePolicyAssignmentsParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee leave policy assignments", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapEmployeeLeavePolicyAssignments(rows), nil
}

func (s *Store) ListLeavePolicyAssignmentsByTemplate(ctx context.Context, tenantID uuid.UUID, templateID uuid.UUID) ([]*domain.EmployeeLeavePolicyAssignment, error) {
	rows, err := s.getQueries(ctx).ListLeavePolicyAssignmentsByTemplate(ctx, sqlc.ListLeavePolicyAssignmentsByTemplateParams{TenantID: tenantID, TemplateID: templateID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave policy assignments by template", err, tenantIDField(tenantID), stringField("template_id", templateID.String()))
	}
	return mapEmployeeLeavePolicyAssignments(rows), nil
}

func (s *Store) DeleteEmployeeLeavePolicyAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteEmployeeLeavePolicyAssignment(ctx, sqlc.SoftDeleteEmployeeLeavePolicyAssignmentParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete employee leave policy assignment", err, tenantIDField(tenantID), stringField("assignment_id", id.String()))
	}
	return nil
}

func leaveTemplateRuleCreateParams(item *domain.LeavePolicyTemplateRule, actorID *uuid.UUID) sqlc.CreateLeavePolicyTemplateRuleParams {
	return sqlc.CreateLeavePolicyTemplateRuleParams{
		TenantID:                  item.TenantID,
		TemplateID:                item.TemplateID,
		LeaveTypeID:               item.LeaveTypeID,
		FyID:                      uuidFromPtr(item.FYID),
		EmploymentTypeID:          uuidFromPtr(item.EmploymentTypeID),
		DepartmentID:              uuidFromPtr(item.DepartmentID),
		DesignationID:             uuidFromPtr(item.DesignationID),
		ProbationStatus:           textFromPtr(item.ProbationStatus),
		AccrualMethod:             item.AccrualMethod,
		AccrualFrequency:          item.AccrualFrequency,
		CreditDays:                numericFromFloat(item.CreditDays),
		CreditHours:               numericFromFloat(item.CreditHours),
		AnnualEntitlement:         numericFromFloat(item.AnnualEntitlement),
		MinWorkedDays:             item.MinWorkedDays,
		MaxBalance:                numericFromFloatPtr(item.MaxBalance),
		CarryForwardEnabled:       item.CarryForwardEnabled,
		MaxCarryForward:           numericFromFloat(item.MaxCarryForward),
		CarryForwardExpiryMonths:  item.CarryForwardExpiryMonths,
		EncashmentEnabled:         item.EncashmentEnabled,
		EncashmentLimit:           numericFromFloat(item.EncashmentLimit),
		EncashmentPayablePercent:  numericFromFloat(item.EncashmentPayablePercent),
		NegativeBalanceAllowed:    item.NegativeBalanceAllowed,
		MaxNegativeBalance:        numericFromFloat(item.MaxNegativeBalance),
		SandwichApplicable:        item.SandwichApplicable,
		IncludeHolidays:           item.IncludeHolidays,
		IncludeWeekoffs:           item.IncludeWeekoffs,
		RequiresDocumentAfterDays: numericFromFloatPtr(item.RequiresDocumentAfterDays),
		MinRequestDays:            numericFromFloat(item.MinRequestDays),
		MaxRequestDays:            numericFromFloatPtr(item.MaxRequestDays),
		MaxRequestsPerYear:        item.MaxRequestsPerYear,
		AccrualDay:                item.AccrualDay,
		LapseUnutilized:           item.LapseUnutilized,
		AllowHalfDay:              item.AllowHalfDay,
		RequiresApproval:          item.RequiresApproval,
		CalculationConfig:         jsonBytesFromMap(item.CalculationConfig),
		Priority:                  item.Priority,
		CreatedBy:                 uuidFromPtr(actorID),
	}
}

func leaveTemplateRuleUpdateParams(item *domain.LeavePolicyTemplateRule, actorID *uuid.UUID) sqlc.UpdateLeavePolicyTemplateRuleParams {
	return sqlc.UpdateLeavePolicyTemplateRuleParams{
		TenantID:                  item.TenantID,
		ID:                        item.ID,
		LeaveTypeID:               item.LeaveTypeID,
		FyID:                      uuidFromPtr(item.FYID),
		EmploymentTypeID:          uuidFromPtr(item.EmploymentTypeID),
		DepartmentID:              uuidFromPtr(item.DepartmentID),
		DesignationID:             uuidFromPtr(item.DesignationID),
		ProbationStatus:           textFromPtr(item.ProbationStatus),
		AccrualMethod:             item.AccrualMethod,
		AccrualFrequency:          item.AccrualFrequency,
		CreditDays:                numericFromFloat(item.CreditDays),
		CreditHours:               numericFromFloat(item.CreditHours),
		AnnualEntitlement:         numericFromFloat(item.AnnualEntitlement),
		MinWorkedDays:             item.MinWorkedDays,
		MaxBalance:                numericFromFloatPtr(item.MaxBalance),
		CarryForwardEnabled:       item.CarryForwardEnabled,
		MaxCarryForward:           numericFromFloat(item.MaxCarryForward),
		CarryForwardExpiryMonths:  item.CarryForwardExpiryMonths,
		EncashmentEnabled:         item.EncashmentEnabled,
		EncashmentLimit:           numericFromFloat(item.EncashmentLimit),
		EncashmentPayablePercent:  numericFromFloat(item.EncashmentPayablePercent),
		NegativeBalanceAllowed:    item.NegativeBalanceAllowed,
		MaxNegativeBalance:        numericFromFloat(item.MaxNegativeBalance),
		SandwichApplicable:        item.SandwichApplicable,
		IncludeHolidays:           item.IncludeHolidays,
		IncludeWeekoffs:           item.IncludeWeekoffs,
		RequiresDocumentAfterDays: numericFromFloatPtr(item.RequiresDocumentAfterDays),
		MinRequestDays:            numericFromFloat(item.MinRequestDays),
		MaxRequestDays:            numericFromFloatPtr(item.MaxRequestDays),
		MaxRequestsPerYear:        item.MaxRequestsPerYear,
		AccrualDay:                item.AccrualDay,
		LapseUnutilized:           item.LapseUnutilized,
		AllowHalfDay:              item.AllowHalfDay,
		RequiresApproval:          item.RequiresApproval,
		CalculationConfig:         jsonBytesFromMap(item.CalculationConfig),
		Priority:                  item.Priority,
		UpdatedBy:                 uuidFromPtr(actorID),
	}
}
