package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreatePolicySet(ctx context.Context, cmd ports.PolicySetCommand) (*domain.PolicySet, error) {
	item, err := s.buildPolicySet(cmd)
	if err != nil {
		s.logError("validate policy set", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_kind", cmd.PolicyKind), serviceStringField("code", cmd.Code))
		return nil, err
	}
	result, err := s.policyEngine.CreatePolicySet(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create policy set", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_kind", item.PolicyKind), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdatePolicySet(ctx context.Context, cmd ports.PolicySetCommand) (*domain.PolicySet, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidPolicySetID
	}
	item, err := s.buildPolicySet(cmd)
	if err != nil {
		s.logError("validate policy set update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_kind", cmd.PolicyKind), serviceStringField("policy_set_id", cmd.ID.String()))
		return nil, err
	}
	result, err := s.policyEngine.UpdatePolicySet(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update policy set", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListPolicySets(ctx context.Context, tenantID uuid.UUID, policyKind string) ([]*domain.PolicySet, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	kind, err := domain.ValidatePolicyKind(policyKind)
	if err != nil {
		s.logError("validate policy set list kind", err, serviceTenantIDField(tenantID), serviceStringField("policy_kind", policyKind))
		return nil, err
	}
	items, err := s.policyEngine.ListPolicySets(ctx, tenantID, kind)
	if err != nil {
		s.logError("list policy sets", err, serviceTenantIDField(tenantID), serviceStringField("policy_kind", kind))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) DeletePolicySet(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if id == uuid.Nil {
		return domain.ErrInvalidPolicySetID
	}
	if err := s.policyEngine.DeletePolicySet(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete policy set", err, serviceTenantIDField(tenantID), serviceStringField("policy_set_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreatePolicyAssignment(ctx context.Context, cmd ports.PolicyAssignmentCommand) (*domain.PolicyAssignment, error) {
	item, err := s.buildPolicyAssignment(cmd)
	if err != nil {
		s.logError("validate policy assignment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_kind", cmd.PolicyKind), serviceStringField("scope_type", cmd.ScopeType))
		return nil, err
	}
	policySet, err := s.policyEngine.GetPolicySet(ctx, cmd.TenantID, cmd.PolicySetID)
	if err != nil {
		s.logError("validate policy assignment policy set", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	if policySet.PolicyKind != item.PolicyKind {
		err := domain.ErrInvalidPolicyKind
		s.logError("validate policy assignment policy kind", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	result, err := s.policyEngine.CreatePolicyAssignment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create policy assignment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_kind", item.PolicyKind), serviceStringField("scope_type", item.ScopeType))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdatePolicyAssignment(ctx context.Context, cmd ports.PolicyAssignmentCommand) (*domain.PolicyAssignment, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidPolicyAssignmentID
	}
	item, err := s.buildPolicyAssignment(cmd)
	if err != nil {
		s.logError("validate policy assignment update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_kind", cmd.PolicyKind), serviceStringField("policy_assignment_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	policySet, err := s.policyEngine.GetPolicySet(ctx, cmd.TenantID, cmd.PolicySetID)
	if err != nil {
		s.logError("validate policy assignment update policy set", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	if policySet.PolicyKind != item.PolicyKind {
		err := domain.ErrInvalidPolicyKind
		s.logError("validate policy assignment update policy kind", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	result, err := s.policyEngine.UpdatePolicyAssignment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update policy assignment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_assignment_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListPolicyAssignments(ctx context.Context, tenantID uuid.UUID, policyKind string) ([]*domain.PolicyAssignment, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	kind, err := domain.ValidatePolicyKind(policyKind)
	if err != nil {
		s.logError("validate policy assignment list kind", err, serviceTenantIDField(tenantID), serviceStringField("policy_kind", policyKind))
		return nil, err
	}
	items, err := s.policyEngine.ListPolicyAssignments(ctx, tenantID, kind)
	if err != nil {
		s.logError("list policy assignments", err, serviceTenantIDField(tenantID), serviceStringField("policy_kind", kind))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) DeletePolicyAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if id == uuid.Nil {
		return domain.ErrInvalidPolicyAssignmentID
	}
	if err := s.policyEngine.DeletePolicyAssignment(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete policy assignment", err, serviceTenantIDField(tenantID), serviceStringField("policy_assignment_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateLeavePolicyRule(ctx context.Context, cmd ports.LeavePolicyRuleCommand) (*domain.LeavePolicyRule, error) {
	item, err := s.buildLeavePolicyRule(cmd)
	if err != nil {
		s.logError("validate leave policy rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	policySet, err := s.policyEngine.GetPolicySet(ctx, cmd.TenantID, cmd.PolicySetID)
	if err != nil {
		s.logError("validate leave policy rule policy set", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	if policySet.PolicyKind != domain.PolicyKindLeave {
		err := domain.ErrInvalidPolicyKind
		s.logError("validate leave policy rule policy kind", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	if _, err := s.leaveTypes.GetLeaveType(ctx, cmd.TenantID, cmd.LeaveTypeID); err != nil {
		s.logError("validate leave policy rule leave type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	result, err := s.policyEngine.CreateLeavePolicyRule(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create leave policy rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateLeavePolicyRule(ctx context.Context, cmd ports.LeavePolicyRuleCommand) (*domain.LeavePolicyRule, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidLeavePolicyRuleID
	}
	item, err := s.buildLeavePolicyRule(cmd)
	if err != nil {
		s.logError("validate leave policy rule update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_policy_rule_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	policySet, err := s.policyEngine.GetPolicySet(ctx, cmd.TenantID, cmd.PolicySetID)
	if err != nil {
		s.logError("validate leave policy rule update policy set", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	if policySet.PolicyKind != domain.PolicyKindLeave {
		err := domain.ErrInvalidPolicyKind
		s.logError("validate leave policy rule update policy kind", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_set_id", cmd.PolicySetID.String()))
		return nil, err
	}
	if _, err := s.leaveTypes.GetLeaveType(ctx, cmd.TenantID, cmd.LeaveTypeID); err != nil {
		s.logError("validate leave policy rule update leave type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	result, err := s.policyEngine.UpdateLeavePolicyRule(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update leave policy rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_policy_rule_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListLeavePolicyRules(ctx context.Context, tenantID uuid.UUID, policySetID uuid.UUID) ([]*domain.LeavePolicyRule, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	if policySetID == uuid.Nil {
		return nil, domain.ErrInvalidPolicySetID
	}
	items, err := s.policyEngine.ListLeavePolicyRules(ctx, tenantID, policySetID)
	if err != nil {
		s.logError("list leave policy rules", err, serviceTenantIDField(tenantID), serviceStringField("policy_set_id", policySetID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) DeleteLeavePolicyRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if id == uuid.Nil {
		return domain.ErrInvalidLeavePolicyRuleID
	}
	if err := s.policyEngine.DeleteLeavePolicyRule(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete leave policy rule", err, serviceTenantIDField(tenantID), serviceStringField("leave_policy_rule_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ResolveEffectivePolicySet(ctx context.Context, subject domain.PolicyResolutionSubject, policyKind string) (*domain.PolicyResolutionResult, error) {
	if subject.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	kind, err := domain.ValidatePolicyKind(policyKind)
	if err != nil {
		s.logError("validate policy resolution kind", err, serviceTenantIDField(subject.TenantID), serviceStringField("policy_kind", policyKind))
		return nil, err
	}
	if subject.Date.IsZero() {
		subject.Date = time.Now().UTC()
	}
	subject.RoleCodes = domain.NormalizeRoleCodes(subject.RoleCodes)

	policySet, err := s.policyEngine.ResolvePolicySet(ctx, subject, kind)
	if err != nil {
		s.logError("resolve policy set", err, serviceTenantIDField(subject.TenantID), serviceStringField("policy_kind", kind), serviceStringField("employee_user_id", subject.EmployeeUserID.String()))
		return nil, err
	}
	candidates, err := s.policyEngine.ListPolicyResolutionCandidates(ctx, subject, kind)
	if err != nil {
		s.logError("list policy resolution candidates", err, serviceTenantIDField(subject.TenantID), serviceStringField("policy_kind", kind), serviceStringField("employee_user_id", subject.EmployeeUserID.String()))
		return nil, err
	}
	result := &domain.PolicyResolutionResult{Policy: policySet, Candidates: candidates}
	if kind == domain.PolicyKindLeave {
		rules, err := s.policyEngine.ListLeavePolicyRules(ctx, subject.TenantID, policySet.ID)
		if err != nil {
			s.logError("list effective leave policy rules", err, serviceTenantIDField(subject.TenantID), serviceStringField("policy_set_id", policySet.ID.String()))
			return nil, err
		}
		result.LeaveRules = rules
	}
	return result, nil
}

func (s *TenantService) ResolveEmployeeAttendancePolicySet(ctx context.Context, tenantID uuid.UUID, employeeUserID uuid.UUID, date string, roleCodes []string) (*domain.PolicyResolutionResult, error) {
	subject, err := s.employeePolicySubject(ctx, tenantID, employeeUserID, date, roleCodes)
	if err != nil {
		return nil, err
	}
	return s.ResolveEffectivePolicySet(ctx, subject, domain.PolicyKindAttendance)
}

func (s *TenantService) ResolveEmployeeLeavePolicySet(ctx context.Context, tenantID uuid.UUID, employeeUserID uuid.UUID, date string, roleCodes []string) (*domain.PolicyResolutionResult, error) {
	subject, err := s.employeePolicySubject(ctx, tenantID, employeeUserID, date, roleCodes)
	if err != nil {
		return nil, err
	}
	return s.ResolveEffectivePolicySet(ctx, subject, domain.PolicyKindLeave)
}

func (s *TenantService) employeePolicySubject(ctx context.Context, tenantID uuid.UUID, employeeUserID uuid.UUID, date string, roleCodes []string) (domain.PolicyResolutionSubject, error) {
	if tenantID == uuid.Nil {
		return domain.PolicyResolutionSubject{}, domain.ErrInvalidTenantID
	}
	if employeeUserID == uuid.Nil {
		return domain.PolicyResolutionSubject{}, domain.ErrInvalidEmployeeUserID
	}
	resolutionDate := time.Now().UTC()
	if strings.TrimSpace(date) != "" {
		parsed, err := parseOptionalDate(date)
		if err != nil {
			s.logError("parse policy resolution date", err, serviceTenantIDField(tenantID), serviceStringField("employee_user_id", employeeUserID.String()))
			return domain.PolicyResolutionSubject{}, err
		}
		if parsed != nil {
			resolutionDate = *parsed
		}
	}
	employee, err := s.employees.GetEmployeeByUserID(ctx, tenantID, employeeUserID)
	if err != nil {
		s.logError("load employee for policy resolution", err, serviceTenantIDField(tenantID), serviceStringField("employee_user_id", employeeUserID.String()))
		return domain.PolicyResolutionSubject{}, err
	}
	return domain.PolicyResolutionSubject{
		TenantID:       tenantID,
		EmployeeUserID: employeeUserID,
		DesignationID:  employee.DesignationID,
		DepartmentID:   employee.DepartmentID,
		BranchID:       employee.BranchID,
		RoleCodes:      roleCodes,
		Date:           resolutionDate,
	}, nil
}

func (s *TenantService) buildPolicySet(cmd ports.PolicySetCommand) (*domain.PolicySet, error) {
	if cmd.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	kind, err := domain.ValidatePolicyKind(cmd.PolicyKind)
	if err != nil {
		return nil, err
	}
	code := strings.ToLower(strings.TrimSpace(cmd.Code))
	if code == "" {
		return nil, domain.ErrInvalidPolicyCode
	}
	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		return nil, domain.ErrInvalidPolicyName
	}
	effectiveFrom, err := parseOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	effectiveTo, err := parseOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	config := cmd.Config
	if len(config) == 0 {
		config = json.RawMessage("{}")
	}
	if !json.Valid(config) {
		config = json.RawMessage("{}")
	}
	return &domain.PolicySet{ID: cmd.ID, TenantID: cmd.TenantID, PolicyKind: kind, Code: code, Name: name, Description: cleanStringPtr(cmd.Description), Config: config, IsDefault: cmd.IsDefault, IsActive: cmd.IsActive, EffectiveFrom: effectiveFrom, EffectiveTo: effectiveTo}, nil
}

func (s *TenantService) buildPolicyAssignment(cmd ports.PolicyAssignmentCommand) (*domain.PolicyAssignment, error) {
	if cmd.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	if cmd.PolicySetID == uuid.Nil {
		return nil, domain.ErrInvalidPolicySetID
	}
	kind, err := domain.ValidatePolicyKind(cmd.PolicyKind)
	if err != nil {
		return nil, err
	}
	roleCode := cleanStringPtr(cmd.RoleCode)
	scopeType, err := domain.ValidatePolicyScope(cmd.ScopeType, cmd.ScopeID, roleCode)
	if err != nil {
		return nil, err
	}
	effectiveFrom, err := parseOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	effectiveTo, err := parseOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	return &domain.PolicyAssignment{TenantID: cmd.TenantID, PolicySetID: cmd.PolicySetID, PolicyKind: kind, ScopeType: scopeType, ScopeID: cmd.ScopeID, RoleCode: roleCode, Priority: cmd.Priority, EffectiveFrom: effectiveFrom, EffectiveTo: effectiveTo, IsActive: cmd.IsActive}, nil
}

func (s *TenantService) buildLeavePolicyRule(cmd ports.LeavePolicyRuleCommand) (*domain.LeavePolicyRule, error) {
	if cmd.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	if cmd.PolicySetID == uuid.Nil {
		return nil, domain.ErrInvalidPolicySetID
	}
	if cmd.LeaveTypeID == uuid.Nil {
		return nil, domain.ErrInvalidLeavePolicyType
	}
	grantMode := strings.TrimSpace(cmd.GrantMode)
	if grantMode == "" {
		grantMode = "annual_calendar"
	}
	probationHandling := strings.TrimSpace(cmd.ProbationHandling)
	if probationHandling == "" {
		probationHandling = "eligible"
	}
	roundingRule := strings.TrimSpace(cmd.RoundingRule)
	if roundingRule == "" {
		roundingRule = "nearest_half"
	}
	insufficientAction := strings.TrimSpace(cmd.InsufficientBalanceAction)
	if insufficientAction == "" {
		insufficientAction = "block"
	}
	payrollImpact := strings.TrimSpace(cmd.PayrollImpact)
	if payrollImpact == "" {
		payrollImpact = "none"
	}
	approvalWorkflow := cmd.ApprovalWorkflow
	if len(approvalWorkflow) == 0 || !json.Valid(approvalWorkflow) {
		approvalWorkflow = json.RawMessage("{}")
	}
	ruleConfig := cmd.RuleConfig
	if len(ruleConfig) == 0 || !json.Valid(ruleConfig) {
		ruleConfig = json.RawMessage("{}")
	}
	return &domain.LeavePolicyRule{
		TenantID:                     cmd.TenantID,
		PolicySetID:                  cmd.PolicySetID,
		LeaveTypeID:                  cmd.LeaveTypeID,
		GrantMode:                    grantMode,
		AccrualFrequency:             cleanStringPtr(cmd.AccrualFrequency),
		EntitlementDays:              cmd.EntitlementDays,
		AccrualAmountPerPeriod:       cmd.AccrualAmountPerPeriod,
		ProrateJoiners:               cmd.ProrateJoiners,
		ProbationHandling:            probationHandling,
		RoundingRule:                 roundingRule,
		MaxBalanceCap:                cmd.MaxBalanceCap,
		CarryForwardCap:              cmd.CarryForwardCap,
		EncashmentEligible:           cmd.EncashmentEligible,
		NegativeBalanceAllowed:       cmd.NegativeBalanceAllowed,
		InsufficientBalanceAction:    insufficientAction,
		ExpiryDays:                   cmd.ExpiryDays,
		AllowHalfDay:                 cmd.AllowHalfDay,
		AttachmentRequiredAfterDays:  cmd.AttachmentRequiredAfterDays,
		ApprovalWorkflow:             approvalWorkflow,
		SandwichEnabled:              cmd.SandwichEnabled,
		SandwichIncludeWeeklyOff:     cmd.SandwichIncludeWeeklyOff,
		SandwichIncludePublicHoliday: cmd.SandwichIncludePublicHoliday,
		SandwichSameLeaveTypeOnly:    cmd.SandwichSameLeaveTypeOnly,
		SandwichAcrossLeaveTypes:     cmd.SandwichAcrossLeaveTypes,
		NoticeRequiredAfterDays:      cmd.NoticeRequiredAfterDays,
		NoticeDays:                   cmd.NoticeDays,
		PayrollImpact:                payrollImpact,
		RuleConfig:                   ruleConfig,
	}, nil
}
