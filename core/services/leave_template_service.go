package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateLeavePolicyTemplate(ctx context.Context, cmd ports.LeavePolicyTemplateCommand) (*domain.LeavePolicyTemplate, error) {
	item := &domain.LeavePolicyTemplate{TenantID: &cmd.TenantID, Name: cmd.Name, Code: cmd.Code, Description: cmd.Description, EffectiveFrom: cmd.EffectiveFrom, EffectiveTo: cmd.EffectiveTo}
	if err := domain.ValidateLeaveTemplate(item); err != nil {
		s.logError("validate leave policy template", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.leaveTemplates.CreateLeavePolicyTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create leave policy template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("template_code", cmd.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListLeavePolicyTemplates(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeavePolicyTemplate, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy template list tenant", err)
		return nil, err
	}
	items, err := s.leaveTemplates.ListLeavePolicyTemplates(ctx, tenantID)
	if err != nil {
		s.logError("list leave policy templates", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetLeavePolicyTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeavePolicyTemplate, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy template get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidLeaveTemplateID
		s.logError("validate leave policy template get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.leaveTemplates.GetLeavePolicyTemplate(ctx, tenantID, id)
	if err != nil {
		s.logError("get leave policy template", err, serviceTenantIDField(tenantID), serviceStringField("template_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateLeavePolicyTemplate(ctx context.Context, cmd ports.LeavePolicyTemplateCommand) (*domain.LeavePolicyTemplate, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidLeaveTemplateID
		s.logError("validate leave policy template update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item := &domain.LeavePolicyTemplate{ID: cmd.ID, TenantID: &cmd.TenantID, Name: cmd.Name, Code: cmd.Code, Description: cmd.Description, EffectiveFrom: cmd.EffectiveFrom, EffectiveTo: cmd.EffectiveTo}
	if err := domain.ValidateLeaveTemplate(item); err != nil {
		s.logError("validate leave policy template update", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.leaveTemplates.UpdateLeavePolicyTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update leave policy template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("template_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteLeavePolicyTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetLeavePolicyTemplate(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.leaveTemplates.DeleteLeavePolicyTemplate(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete leave policy template", err, serviceTenantIDField(tenantID), serviceStringField("template_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateLeavePolicyTemplateRule(ctx context.Context, cmd ports.LeavePolicyTemplateRuleCommand) (*domain.LeavePolicyTemplateRule, error) {
	item, err := s.buildLeavePolicyTemplateRule(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.leaveTemplates.CreateLeavePolicyTemplateRule(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create leave policy template rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("template_id", cmd.TemplateID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListLeavePolicyTemplateRules(ctx context.Context, tenantID uuid.UUID, templateID uuid.UUID) ([]*domain.LeavePolicyTemplateRule, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy template rule list tenant", err)
		return nil, err
	}
	if templateID == uuid.Nil {
		err := domain.ErrInvalidLeaveTemplateID
		s.logError("validate leave policy template rule list template", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.leaveTemplates.ListLeavePolicyTemplateRules(ctx, tenantID, templateID)
	if err != nil {
		s.logError("list leave policy template rules", err, serviceTenantIDField(tenantID), serviceStringField("template_id", templateID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateLeavePolicyTemplateRule(ctx context.Context, cmd ports.LeavePolicyTemplateRuleCommand) (*domain.LeavePolicyTemplateRule, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidLeaveTemplateRuleID
		s.logError("validate leave policy template rule update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := s.buildLeavePolicyTemplateRule(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.leaveTemplates.UpdateLeavePolicyTemplateRule(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update leave policy template rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("rule_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteLeavePolicyTemplateRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy template rule delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidLeaveTemplateRuleID
		s.logError("validate leave policy template rule delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.leaveTemplates.DeleteLeavePolicyTemplateRule(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete leave policy template rule", err, serviceTenantIDField(tenantID), serviceStringField("rule_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) UpsertEmployeeLeavePolicyAssignment(ctx context.Context, cmd ports.EmployeeLeavePolicyAssignmentCommand) (*domain.EmployeeLeavePolicyAssignment, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy assignment tenant", err)
		return nil, err
	}
	if cmd.UserID == uuid.Nil {
		err := domain.ErrInvalidLeaveBalanceUser
		s.logError("validate leave policy assignment user", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.TemplateID == uuid.Nil {
		err := domain.ErrInvalidLeaveTemplateID
		s.logError("validate leave policy assignment template", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.EffectiveFrom.IsZero() {
		err := domain.ErrInvalidDateRange
		s.logError("validate leave policy assignment effective date", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID); err != nil {
		s.logError("validate leave policy assignment employee", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if _, err := s.GetLeavePolicyTemplate(ctx, cmd.TenantID, cmd.TemplateID); err != nil {
		return nil, err
	}
	if cmd.FYID != nil {
		if _, err := s.financialYears.GetFinancialYear(ctx, cmd.TenantID, *cmd.FYID); err != nil {
			s.logError("validate leave policy assignment fy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
			return nil, err
		}
	}
	item := &domain.EmployeeLeavePolicyAssignment{ID: cmd.ID, TenantID: cmd.TenantID, UserID: cmd.UserID, TemplateID: cmd.TemplateID, FYID: cmd.FYID, EffectiveFrom: cmd.EffectiveFrom, EffectiveTo: cmd.EffectiveTo}
	result, err := s.leaveTemplates.UpsertEmployeeLeavePolicyAssignment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert leave policy assignment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListEmployeeLeavePolicyAssignments(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.EmployeeLeavePolicyAssignment, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy assignments tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidLeaveBalanceUser
		s.logError("validate leave policy assignments user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.leaveTemplates.ListEmployeeLeavePolicyAssignments(ctx, tenantID, userID)
	if err != nil {
		s.logError("list leave policy assignments", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) DeleteEmployeeLeavePolicyAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy assignment delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrLeavePolicyAssignmentNotFound
		s.logError("validate leave policy assignment delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.leaveTemplates.DeleteEmployeeLeavePolicyAssignment(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete leave policy assignment", err, serviceTenantIDField(tenantID), serviceStringField("assignment_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) buildLeavePolicyTemplateRule(ctx context.Context, cmd ports.LeavePolicyTemplateRuleCommand) (*domain.LeavePolicyTemplateRule, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy template rule tenant", err)
		return nil, err
	}
	if _, err := s.GetLeavePolicyTemplate(ctx, cmd.TenantID, cmd.TemplateID); err != nil {
		return nil, err
	}
	if _, err := s.leaveTypes.GetLeaveType(ctx, cmd.TenantID, cmd.LeaveTypeID); err != nil {
		s.logError("validate leave policy template rule leave type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	if cmd.FYID != nil {
		if _, err := s.financialYears.GetFinancialYear(ctx, cmd.TenantID, *cmd.FYID); err != nil {
			s.logError("validate leave policy template rule fy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
			return nil, err
		}
	}
	item := &domain.LeavePolicyTemplateRule{
		ID:                        cmd.ID,
		TenantID:                  cmd.TenantID,
		TemplateID:                cmd.TemplateID,
		LeaveTypeID:               cmd.LeaveTypeID,
		FYID:                      cmd.FYID,
		EmploymentTypeID:          cmd.EmploymentTypeID,
		DepartmentID:              cmd.DepartmentID,
		DesignationID:             cmd.DesignationID,
		ProbationStatus:           cmd.ProbationStatus,
		AccrualMethod:             cmd.AccrualMethod,
		AccrualFrequency:          cmd.AccrualFrequency,
		CreditDays:                cmd.CreditDays,
		CreditHours:               cmd.CreditHours,
		AnnualEntitlement:         cmd.AnnualEntitlement,
		MinWorkedDays:             cmd.MinWorkedDays,
		MaxBalance:                cmd.MaxBalance,
		CarryForwardEnabled:       cmd.CarryForwardEnabled,
		MaxCarryForward:           cmd.MaxCarryForward,
		CarryForwardExpiryMonths:  cmd.CarryForwardExpiryMonths,
		EncashmentEnabled:         cmd.EncashmentEnabled,
		EncashmentLimit:           cmd.EncashmentLimit,
		EncashmentPayablePercent:  cmd.EncashmentPayablePercent,
		NegativeBalanceAllowed:    cmd.NegativeBalanceAllowed,
		MaxNegativeBalance:        cmd.MaxNegativeBalance,
		SandwichApplicable:        cmd.SandwichApplicable,
		IncludeHolidays:           cmd.IncludeHolidays,
		IncludeWeekoffs:           cmd.IncludeWeekoffs,
		RequiresDocumentAfterDays: cmd.RequiresDocumentAfterDays,
		MinRequestDays:            cmd.MinRequestDays,
		MaxRequestDays:            cmd.MaxRequestDays,
		MaxRequestsPerYear:        cmd.MaxRequestsPerYear,
		AccrualDay:                cmd.AccrualDay,
		LapseUnutilized:           cmd.LapseUnutilized,
		AllowHalfDay:              cmd.AllowHalfDay,
		RequiresApproval:          cmd.RequiresApproval,
		CalculationConfig:         cmd.CalculationConfig,
		Priority:                  cmd.Priority,
	}
	if err := domain.ValidateLeaveTemplateRule(item); err != nil {
		s.logError("validate leave policy template rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("template_id", cmd.TemplateID.String()))
		return nil, err
	}
	return item, nil
}
