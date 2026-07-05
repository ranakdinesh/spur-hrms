package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateLeavePolicy(ctx context.Context, cmd ports.LeavePolicyCommand) (*domain.LeavePolicy, error) {
	item, err := s.buildLeavePolicy(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if existing, err := s.leavePolicies.GetLeavePolicyByTypeAndFY(ctx, cmd.TenantID, cmd.LeaveTypeID, cmd.FYID); err == nil && existing != nil {
		err := domain.ErrLeavePolicyAlreadyExists
		s.logError("validate leave policy uniqueness", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	} else if err != nil && !errors.Is(err, domain.ErrLeavePolicyNotFound) {
		s.logError("check leave policy uniqueness", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	result, err := s.leavePolicies.CreateLeavePolicy(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create leave policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListLeavePolicies(ctx context.Context, tenantID uuid.UUID) ([]*domain.LeavePolicy, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy list tenant", err)
		return nil, err
	}
	items, err := s.leavePolicies.ListLeavePolicies(ctx, tenantID)
	if err != nil {
		s.logError("list leave policies", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetLeavePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.LeavePolicy, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidLeavePolicyID
		s.logError("validate leave policy get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.leavePolicies.GetLeavePolicy(ctx, tenantID, id)
	if err != nil {
		s.logError("get leave policy", err, serviceTenantIDField(tenantID), serviceStringField("leave_policy_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateLeavePolicy(ctx context.Context, cmd ports.LeavePolicyCommand) (*domain.LeavePolicy, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidLeavePolicyID
		s.logError("validate leave policy update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetLeavePolicy(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := s.buildLeavePolicy(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	if existing, err := s.leavePolicies.GetLeavePolicyByTypeAndFY(ctx, cmd.TenantID, cmd.LeaveTypeID, cmd.FYID); err == nil && existing != nil && existing.ID != cmd.ID {
		err := domain.ErrLeavePolicyAlreadyExists
		s.logError("validate leave policy update uniqueness", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_policy_id", cmd.ID.String()))
		return nil, err
	} else if err != nil && !errors.Is(err, domain.ErrLeavePolicyNotFound) {
		s.logError("check leave policy update uniqueness", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_policy_id", cmd.ID.String()))
		return nil, err
	}
	result, err := s.leavePolicies.UpdateLeavePolicy(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update leave policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_policy_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteLeavePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidLeavePolicyID
		s.logError("validate leave policy delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if _, err := s.GetLeavePolicy(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.leavePolicies.DeleteLeavePolicy(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete leave policy", err, serviceTenantIDField(tenantID), serviceStringField("leave_policy_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) RunMonthlyLeaveAllocation(ctx context.Context, cmd ports.MonthlyLeaveAllocationCommand) ([]*domain.MonthlyLeaveAllocation, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate monthly leave allocation tenant", err)
		return nil, err
	}
	if cmd.FYID == uuid.Nil {
		err := domain.ErrInvalidLeavePolicyFY
		s.logError("validate monthly leave allocation financial year", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.Month < 1 || cmd.Month > 12 {
		err := domain.ErrInvalidLeaveAllocationMonth
		s.logError("validate monthly leave allocation month", err, serviceTenantIDField(cmd.TenantID), serviceStringField("month", fmt.Sprint(cmd.Month)))
		return nil, err
	}
	if _, err := s.financialYears.GetFinancialYear(ctx, cmd.TenantID, cmd.FYID); err != nil {
		s.logError("validate monthly leave allocation financial year exists", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	policies, err := s.leavePolicies.ListLeavePolicies(ctx, cmd.TenantID)
	if err != nil {
		s.logError("list monthly leave allocation policies", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	allocations := make([]*domain.MonthlyLeaveAllocation, 0)
	for _, policy := range policies {
		if policy.FYID != cmd.FYID || policy.AllocationType != domain.LeaveAllocationMonthly {
			continue
		}
		days, err := policy.AllocationForMonth(cmd.Month)
		if err != nil {
			s.logError("resolve monthly leave allocation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_policy_id", policy.ID.String()), serviceStringField("month", fmt.Sprint(cmd.Month)))
			return nil, err
		}
		if days <= 0 {
			continue
		}
		allocations = append(allocations, &domain.MonthlyLeaveAllocation{
			TenantID:    cmd.TenantID,
			PolicyID:    policy.ID,
			LeaveTypeID: policy.LeaveTypeID,
			FYID:        policy.FYID,
			Month:       cmd.Month,
			Days:        days,
		})
	}
	return allocations, nil
}

func (s *TenantService) buildLeavePolicy(ctx context.Context, cmd ports.LeavePolicyCommand) (*domain.LeavePolicy, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave policy tenant", err)
		return nil, err
	}
	if _, err := s.leaveTypes.GetLeaveType(ctx, cmd.TenantID, cmd.LeaveTypeID); err != nil {
		s.logError("validate leave policy leave type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	if _, err := s.financialYears.GetFinancialYear(ctx, cmd.TenantID, cmd.FYID); err != nil {
		s.logError("validate leave policy financial year", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	item, err := domain.NewLeavePolicy(domain.LeavePolicyInput{
		TenantID:             cmd.TenantID,
		LeaveTypeID:          cmd.LeaveTypeID,
		FYID:                 cmd.FYID,
		TotalDays:            cmd.TotalDays,
		AllocationType:       cmd.AllocationType,
		Monthly:              [12]int{int(cmd.Jan), int(cmd.Feb), int(cmd.Mar), int(cmd.Apr), int(cmd.May), int(cmd.Jun), int(cmd.Jul), int(cmd.Aug), int(cmd.Sep), int(cmd.Oct), int(cmd.Nov), int(cmd.Dec)},
		IsSandwichApplicable: cmd.IsSandwichApplicable,
	})
	if err != nil {
		s.logError("validate leave policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	return item, nil
}
