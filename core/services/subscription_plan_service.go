package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateSubscriptionPlan(ctx context.Context, cmd ports.SubscriptionPlanCommand) (*domain.SubscriptionPlan, error) {
	item, err := subscriptionPlanFromCommand(cmd)
	if err != nil {
		s.logError("validate subscription plan create", err, serviceStringField("plan_code", cmd.Code))
		return nil, err
	}
	result, err := s.subscriptionPlans.CreateSubscriptionPlan(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create subscription plan", err, serviceStringField("plan_code", item.Code))
		return nil, err
	}
	s.log.Info().Str("plan_id", result.ID.String()).Str("plan_code", result.Code).Msg("hrms: subscription plan created")
	return result, nil
}

func (s *TenantService) ListSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	result, err := s.subscriptionPlans.ListSubscriptionPlans(ctx)
	if err != nil {
		s.logError("list subscription plans", err)
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListActiveSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	result, err := s.subscriptionPlans.ListActiveSubscriptionPlans(ctx)
	if err != nil {
		s.logError("list active subscription plans", err)
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetSubscriptionPlan(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error) {
	if id == uuid.Nil {
		err := domain.ErrInvalidSubscriptionPlanID
		s.logError("validate subscription plan get id", err)
		return nil, err
	}
	result, err := s.subscriptionPlans.GetSubscriptionPlan(ctx, id)
	if err != nil {
		s.logError("get subscription plan", err, serviceStringField("plan_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateSubscriptionPlan(ctx context.Context, cmd ports.SubscriptionPlanCommand) (*domain.SubscriptionPlan, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidSubscriptionPlanID
		s.logError("validate subscription plan update id", err)
		return nil, err
	}
	item, err := subscriptionPlanFromCommand(cmd)
	if err != nil {
		s.logError("validate subscription plan update", err, serviceStringField("plan_id", cmd.ID.String()), serviceStringField("plan_code", cmd.Code))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.subscriptionPlans.UpdateSubscriptionPlan(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update subscription plan", err, serviceStringField("plan_id", cmd.ID.String()), serviceStringField("plan_code", item.Code))
		return nil, err
	}
	s.log.Info().Str("plan_id", result.ID.String()).Str("plan_code", result.Code).Msg("hrms: subscription plan updated")
	return result, nil
}

func (s *TenantService) DeleteSubscriptionPlan(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error {
	if id == uuid.Nil {
		err := domain.ErrInvalidSubscriptionPlanID
		s.logError("validate subscription plan delete id", err)
		return err
	}
	if err := s.subscriptionPlans.DeleteSubscriptionPlan(ctx, id, actorID); err != nil {
		s.logError("delete subscription plan", err, serviceStringField("plan_id", id.String()))
		return err
	}
	s.log.Info().Str("plan_id", id.String()).Msg("hrms: subscription plan deactivated")
	return nil
}

func subscriptionPlanFromCommand(cmd ports.SubscriptionPlanCommand) (*domain.SubscriptionPlan, error) {
	return domain.NewSubscriptionPlan(domain.SubscriptionPlanInput{
		ID:                cmd.ID,
		Code:              cmd.Code,
		Name:              cmd.Name,
		Description:       cmd.Description,
		PriceAmount:       cmd.PriceAmount,
		PriceBasis:        cmd.PriceBasis,
		MinimumAmount:     cmd.MinimumAmount,
		IncludedEmployees: cmd.IncludedEmployees,
		OverageAmount:     cmd.OverageAmount,
		CurrencyCode:      cmd.CurrencyCode,
		BillingCycle:      cmd.BillingCycle,
		EmployeeLimit:     cmd.EmployeeLimit,
		TrialDays:         cmd.TrialDays,
		Visibility:        cmd.Visibility,
		IsActive:          cmd.IsActive,
	})
}
