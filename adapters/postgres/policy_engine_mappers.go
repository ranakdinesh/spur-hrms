package postgres

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapPolicySet(row sqlc.HrmsPolicySet) *domain.PolicySet {
	return policySetFromParts(row.ID, row.TenantID, row.PolicyKind, row.Code, row.Name, row.Description, row.Config, row.IsDefault, row.IsActive, row.EffectiveFrom, row.EffectiveTo, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy)
}

func mapResolvedPolicySet(row sqlc.ResolvePolicySetRow) *domain.PolicySet {
	return policySetFromParts(row.ID, row.TenantID, row.PolicyKind, row.Code, row.Name, row.Description, row.Config, row.IsDefault, row.IsActive, row.EffectiveFrom, row.EffectiveTo, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy)
}

func policySetFromParts(id uuid.UUID, tenantID uuid.UUID, kind string, code string, name string, description pgtype.Text, config []byte, isDefault bool, isActive bool, effectiveFrom pgtype.Date, effectiveTo pgtype.Date, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID) *domain.PolicySet {
	return &domain.PolicySet{
		ID:            id,
		TenantID:      tenantID,
		PolicyKind:    kind,
		Code:          code,
		Name:          name,
		Description:   ptrFromText(description),
		Config:        jsonRawDefault(config, "{}"),
		IsDefault:     isDefault,
		IsActive:      isActive,
		EffectiveFrom: ptrFromDate(effectiveFrom),
		EffectiveTo:   ptrFromDate(effectiveTo),
		Inactive:      inactive,
		CreatedAt:     timeFromTimestamptz(createdAt),
		CreatedBy:     ptrFromUUID(createdBy),
		UpdatedAt:     timeFromTimestamptz(updatedAt),
		UpdatedBy:     ptrFromUUID(updatedBy),
	}
}

func mapPolicySets(rows []sqlc.HrmsPolicySet) []*domain.PolicySet {
	items := make([]*domain.PolicySet, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPolicySet(row))
	}
	return items
}

func mapPolicyAssignment(row sqlc.HrmsPolicyAssignment) *domain.PolicyAssignment {
	return &domain.PolicyAssignment{
		ID:            row.ID,
		TenantID:      row.TenantID,
		PolicySetID:   row.PolicySetID,
		PolicyKind:    row.PolicyKind,
		ScopeType:     row.ScopeType,
		ScopeID:       ptrFromUUID(row.ScopeID),
		RoleCode:      ptrFromText(row.RoleCode),
		Priority:      row.Priority,
		EffectiveFrom: ptrFromDate(row.EffectiveFrom),
		EffectiveTo:   ptrFromDate(row.EffectiveTo),
		IsActive:      row.IsActive,
		Inactive:      row.Inactive,
		CreatedAt:     timeFromTimestamptz(row.CreatedAt),
		CreatedBy:     ptrFromUUID(row.CreatedBy),
		UpdatedAt:     timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:     ptrFromUUID(row.UpdatedBy),
	}
}

func mapPolicyAssignments(rows []sqlc.HrmsPolicyAssignment) []*domain.PolicyAssignment {
	items := make([]*domain.PolicyAssignment, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapPolicyAssignment(row))
	}
	return items
}

func mapLeavePolicyRule(row sqlc.HrmsLeavePolicyRule) *domain.LeavePolicyRule {
	return &domain.LeavePolicyRule{
		ID:                           row.ID,
		TenantID:                     row.TenantID,
		PolicySetID:                  row.PolicySetID,
		LeaveTypeID:                  row.LeaveTypeID,
		GrantMode:                    row.GrantMode,
		AccrualFrequency:             ptrFromText(row.AccrualFrequency),
		EntitlementDays:              floatFromNumeric(row.EntitlementDays),
		AccrualAmountPerPeriod:       floatFromNumeric(row.AccrualAmountPerPeriod),
		ProrateJoiners:               row.ProrateJoiners,
		ProbationHandling:            row.ProbationHandling,
		RoundingRule:                 row.RoundingRule,
		MaxBalanceCap:                floatPtrFromNumeric(row.MaxBalanceCap),
		CarryForwardCap:              floatPtrFromNumeric(row.CarryForwardCap),
		EncashmentEligible:           row.EncashmentEligible,
		NegativeBalanceAllowed:       row.NegativeBalanceAllowed,
		InsufficientBalanceAction:    row.InsufficientBalanceAction,
		ExpiryDays:                   ptrFromInt4(row.ExpiryDays),
		AllowHalfDay:                 row.AllowHalfDay,
		AttachmentRequiredAfterDays:  floatPtrFromNumeric(row.AttachmentRequiredAfterDays),
		ApprovalWorkflow:             jsonRawDefault(row.ApprovalWorkflow, "{}"),
		SandwichEnabled:              row.SandwichEnabled,
		SandwichIncludeWeeklyOff:     row.SandwichIncludeWeeklyOff,
		SandwichIncludePublicHoliday: row.SandwichIncludePublicHoliday,
		SandwichSameLeaveTypeOnly:    row.SandwichSameLeaveTypeOnly,
		SandwichAcrossLeaveTypes:     row.SandwichAcrossLeaveTypes,
		NoticeRequiredAfterDays:      floatPtrFromNumeric(row.NoticeRequiredAfterDays),
		NoticeDays:                   row.NoticeDays,
		PayrollImpact:                row.PayrollImpact,
		RuleConfig:                   jsonRawDefault(row.RuleConfig, "{}"),
		Inactive:                     row.Inactive,
		CreatedAt:                    timeFromTimestamptz(row.CreatedAt),
		CreatedBy:                    ptrFromUUID(row.CreatedBy),
		UpdatedAt:                    timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:                    ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeavePolicyRules(rows []sqlc.HrmsLeavePolicyRule) []*domain.LeavePolicyRule {
	items := make([]*domain.LeavePolicyRule, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeavePolicyRule(row))
	}
	return items
}

func mapPolicyResolutionCandidates(rows []sqlc.ListPolicyResolutionCandidatesRow) []domain.PolicyResolutionCandidate {
	items := make([]domain.PolicyResolutionCandidate, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.PolicyResolutionCandidate{
			PolicySetID: row.ID,
			TenantID:    row.TenantID,
			PolicyKind:  row.PolicyKind,
			Code:        row.Code,
			Name:        row.Name,
			ScopeType:   row.ScopeType,
			ScopeID:     ptrFromUUID(row.ScopeID),
			RoleCode:    ptrFromText(row.RoleCode),
			Precedence:  row.Precedence,
		})
	}
	return items
}
