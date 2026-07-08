package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreatePolicySet(ctx context.Context, item *domain.PolicySet, actorID *uuid.UUID) (*domain.PolicySet, error) {
	row, err := s.getQueries(ctx).CreatePolicySet(ctx, sqlc.CreatePolicySetParams{
		TenantID:      item.TenantID,
		PolicyKind:    item.PolicyKind,
		Code:          item.Code,
		Name:          item.Name,
		Description:   textFromPtr(item.Description),
		Config:        jsonRawDefault(item.Config, "{}"),
		IsDefault:     item.IsDefault,
		IsActive:      item.IsActive,
		EffectiveFrom: dateFromPtr(item.EffectiveFrom),
		EffectiveTo:   dateFromPtr(item.EffectiveTo),
		CreatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create policy set", err, tenantIDField(item.TenantID), stringField("policy_kind", item.PolicyKind), stringField("code", item.Code))
	}
	return mapPolicySet(row), nil
}

func (s *Store) UpdatePolicySet(ctx context.Context, item *domain.PolicySet, actorID *uuid.UUID) (*domain.PolicySet, error) {
	row, err := s.getQueries(ctx).UpdatePolicySet(ctx, sqlc.UpdatePolicySetParams{
		TenantID:      item.TenantID,
		ID:            item.ID,
		Code:          item.Code,
		Name:          item.Name,
		Description:   textFromPtr(item.Description),
		Config:        jsonRawDefault(item.Config, "{}"),
		IsDefault:     item.IsDefault,
		IsActive:      item.IsActive,
		EffectiveFrom: dateFromPtr(item.EffectiveFrom),
		EffectiveTo:   dateFromPtr(item.EffectiveTo),
		UpdatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPolicySetNotFound
		}
		return nil, s.logDBError(ctx, "update policy set", err, tenantIDField(item.TenantID), stringField("policy_set_id", item.ID.String()))
	}
	return mapPolicySet(row), nil
}

func (s *Store) ListPolicySets(ctx context.Context, tenantID uuid.UUID, policyKind string) ([]*domain.PolicySet, error) {
	rows, err := s.getQueries(ctx).ListPolicySets(ctx, sqlc.ListPolicySetsParams{TenantID: tenantID, PolicyKind: policyKind})
	if err != nil {
		return nil, s.logDBError(ctx, "list policy sets", err, tenantIDField(tenantID), stringField("policy_kind", policyKind))
	}
	return mapPolicySets(rows), nil
}

func (s *Store) GetPolicySet(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PolicySet, error) {
	row, err := s.getQueries(ctx).GetPolicySet(ctx, sqlc.GetPolicySetParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPolicySetNotFound
		}
		return nil, s.logDBError(ctx, "get policy set", err, tenantIDField(tenantID), stringField("policy_set_id", id.String()))
	}
	return mapPolicySet(row), nil
}

func (s *Store) DeletePolicySet(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePolicySet(ctx, sqlc.SoftDeletePolicySetParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete policy set", err, tenantIDField(tenantID), stringField("policy_set_id", id.String()))
	}
	return nil
}

func (s *Store) CreatePolicyAssignment(ctx context.Context, item *domain.PolicyAssignment, actorID *uuid.UUID) (*domain.PolicyAssignment, error) {
	row, err := s.getQueries(ctx).CreatePolicyAssignment(ctx, sqlc.CreatePolicyAssignmentParams{
		TenantID:      item.TenantID,
		PolicySetID:   item.PolicySetID,
		PolicyKind:    item.PolicyKind,
		ScopeType:     item.ScopeType,
		ScopeID:       uuidFromPtr(item.ScopeID),
		RoleCode:      textFromPtr(item.RoleCode),
		Priority:      item.Priority,
		EffectiveFrom: dateFromPtr(item.EffectiveFrom),
		EffectiveTo:   dateFromPtr(item.EffectiveTo),
		IsActive:      item.IsActive,
		CreatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create policy assignment", err, tenantIDField(item.TenantID), stringField("policy_kind", item.PolicyKind), stringField("scope_type", item.ScopeType))
	}
	return mapPolicyAssignment(row), nil
}

func (s *Store) ListPolicyAssignments(ctx context.Context, tenantID uuid.UUID, policyKind string) ([]*domain.PolicyAssignment, error) {
	rows, err := s.getQueries(ctx).ListPolicyAssignments(ctx, sqlc.ListPolicyAssignmentsParams{TenantID: tenantID, PolicyKind: policyKind})
	if err != nil {
		return nil, s.logDBError(ctx, "list policy assignments", err, tenantIDField(tenantID), stringField("policy_kind", policyKind))
	}
	return mapPolicyAssignments(rows), nil
}

func (s *Store) DeletePolicyAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePolicyAssignment(ctx, sqlc.SoftDeletePolicyAssignmentParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete policy assignment", err, tenantIDField(tenantID), stringField("policy_assignment_id", id.String()))
	}
	return nil
}

func (s *Store) CreateLeavePolicyRule(ctx context.Context, item *domain.LeavePolicyRule, actorID *uuid.UUID) (*domain.LeavePolicyRule, error) {
	row, err := s.getQueries(ctx).CreateLeavePolicyRule(ctx, sqlc.CreateLeavePolicyRuleParams{
		TenantID:                     item.TenantID,
		PolicySetID:                  item.PolicySetID,
		LeaveTypeID:                  item.LeaveTypeID,
		GrantMode:                    item.GrantMode,
		AccrualFrequency:             textFromPtr(item.AccrualFrequency),
		EntitlementDays:              numericFromFloat(item.EntitlementDays),
		AccrualAmountPerPeriod:       numericFromFloat(item.AccrualAmountPerPeriod),
		ProrateJoiners:               item.ProrateJoiners,
		ProbationHandling:            item.ProbationHandling,
		RoundingRule:                 item.RoundingRule,
		MaxBalanceCap:                numericFromFloatPtr(item.MaxBalanceCap),
		CarryForwardCap:              numericFromFloatPtr(item.CarryForwardCap),
		EncashmentEligible:           item.EncashmentEligible,
		NegativeBalanceAllowed:       item.NegativeBalanceAllowed,
		InsufficientBalanceAction:    item.InsufficientBalanceAction,
		ExpiryDays:                   int4FromPtr(item.ExpiryDays),
		AllowHalfDay:                 item.AllowHalfDay,
		AttachmentRequiredAfterDays:  numericFromFloatPtr(item.AttachmentRequiredAfterDays),
		ApprovalWorkflow:             jsonRawDefault(item.ApprovalWorkflow, "{}"),
		SandwichEnabled:              item.SandwichEnabled,
		SandwichIncludeWeeklyOff:     item.SandwichIncludeWeeklyOff,
		SandwichIncludePublicHoliday: item.SandwichIncludePublicHoliday,
		SandwichSameLeaveTypeOnly:    item.SandwichSameLeaveTypeOnly,
		SandwichAcrossLeaveTypes:     item.SandwichAcrossLeaveTypes,
		NoticeRequiredAfterDays:      numericFromFloatPtr(item.NoticeRequiredAfterDays),
		NoticeDays:                   item.NoticeDays,
		PayrollImpact:                item.PayrollImpact,
		RuleConfig:                   jsonRawDefault(item.RuleConfig, "{}"),
		CreatedBy:                    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create leave policy rule", err, tenantIDField(item.TenantID), stringField("policy_set_id", item.PolicySetID.String()))
	}
	return mapLeavePolicyRule(row), nil
}

func (s *Store) ListLeavePolicyRules(ctx context.Context, tenantID uuid.UUID, policySetID uuid.UUID) ([]*domain.LeavePolicyRule, error) {
	rows, err := s.getQueries(ctx).ListLeavePolicyRules(ctx, sqlc.ListLeavePolicyRulesParams{TenantID: tenantID, PolicySetID: policySetID})
	if err != nil {
		return nil, s.logDBError(ctx, "list leave policy rules", err, tenantIDField(tenantID), stringField("policy_set_id", policySetID.String()))
	}
	return mapLeavePolicyRules(rows), nil
}

func (s *Store) ResolvePolicySet(ctx context.Context, subject domain.PolicyResolutionSubject, policyKind string) (*domain.PolicySet, error) {
	row, err := s.getQueries(ctx).ResolvePolicySet(ctx, policyResolutionParams(subject, policyKind))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPolicySetNotFound
		}
		return nil, s.logDBError(ctx, "resolve policy set", err, tenantIDField(subject.TenantID), stringField("policy_kind", policyKind), stringField("employee_user_id", subject.EmployeeUserID.String()))
	}
	return mapResolvedPolicySet(row), nil
}

func (s *Store) ListPolicyResolutionCandidates(ctx context.Context, subject domain.PolicyResolutionSubject, policyKind string) ([]domain.PolicyResolutionCandidate, error) {
	rows, err := s.getQueries(ctx).ListPolicyResolutionCandidates(ctx, policyResolutionCandidateParams(subject, policyKind))
	if err != nil {
		return nil, s.logDBError(ctx, "list policy resolution candidates", err, tenantIDField(subject.TenantID), stringField("policy_kind", policyKind), stringField("employee_user_id", subject.EmployeeUserID.String()))
	}
	return mapPolicyResolutionCandidates(rows), nil
}

func policyResolutionParams(subject domain.PolicyResolutionSubject, policyKind string) sqlc.ResolvePolicySetParams {
	date := subject.Date
	if date.IsZero() {
		date = time.Now().UTC()
	}
	return sqlc.ResolvePolicySetParams{
		TenantID:      subject.TenantID,
		PolicyKind:    policyKind,
		ScopeID:       uuidFromPtr(&subject.EmployeeUserID),
		ScopeID_2:     uuidFromPtr(subject.DesignationID),
		ScopeID_3:     uuidFromPtr(subject.WorkforceTypeID),
		ScopeID_4:     uuidFromPtr(subject.DepartmentID),
		ScopeID_5:     uuidFromPtr(subject.BranchID),
		Column8:       domain.NormalizeRoleCodes(subject.RoleCodes),
		EffectiveFrom: dateFromTime(date),
	}
}

func policyResolutionCandidateParams(subject domain.PolicyResolutionSubject, policyKind string) sqlc.ListPolicyResolutionCandidatesParams {
	params := policyResolutionParams(subject, policyKind)
	return sqlc.ListPolicyResolutionCandidatesParams{
		TenantID:      params.TenantID,
		PolicyKind:    params.PolicyKind,
		ScopeID:       params.ScopeID,
		ScopeID_2:     params.ScopeID_2,
		ScopeID_3:     params.ScopeID_3,
		ScopeID_4:     params.ScopeID_4,
		ScopeID_5:     params.ScopeID_5,
		Column8:       params.Column8,
		EffectiveFrom: params.EffectiveFrom,
	}
}
