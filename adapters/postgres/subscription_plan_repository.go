package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateSubscriptionPlan(ctx context.Context, item *domain.SubscriptionPlan, actorID *uuid.UUID) (*domain.SubscriptionPlan, error) {
	row, err := s.getQueries(ctx).CreateSubscriptionPlan(ctx, sqlc.CreateSubscriptionPlanParams{
		Code:              item.Code,
		Name:              item.Name,
		Description:       textFromPtr(item.Description),
		PriceAmount:       numericFromMoney(item.PriceAmount),
		PriceBasis:        item.PriceBasis,
		MinimumAmount:     numericFromMoney(item.MinimumAmount),
		IncludedEmployees: item.IncludedEmployees,
		OverageAmount:     numericFromMoney(item.OverageAmount),
		CurrencyCode:      item.CurrencyCode,
		BillingCycle:      item.BillingCycle,
		EmployeeLimit:     item.EmployeeLimit,
		TrialDays:         item.TrialDays,
		Visibility:        item.Visibility,
		IsActive:          item.IsActive,
		CreatedBy:         uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create subscription plan", err, stringField("plan_code", item.Code))
	}
	return mapSubscriptionPlan(row), nil
}

func (s *Store) ListSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	rows, err := s.getQueries(ctx).ListSubscriptionPlans(ctx)
	if err != nil {
		return nil, s.logDBError(ctx, "list subscription plans", err)
	}
	return mapSubscriptionPlans(rows), nil
}

func (s *Store) ListActiveSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	rows, err := s.getQueries(ctx).ListActiveSubscriptionPlans(ctx)
	if err != nil {
		return nil, s.logDBError(ctx, "list active subscription plans", err)
	}
	return mapSubscriptionPlans(rows), nil
}

func (s *Store) GetSubscriptionPlan(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error) {
	row, err := s.getQueries(ctx).GetSubscriptionPlan(ctx, id)
	if err != nil {
		return nil, s.logDBError(ctx, "get subscription plan", err, stringField("plan_id", id.String()))
	}
	return mapSubscriptionPlan(row), nil
}

func (s *Store) UpdateSubscriptionPlan(ctx context.Context, item *domain.SubscriptionPlan, actorID *uuid.UUID) (*domain.SubscriptionPlan, error) {
	row, err := s.getQueries(ctx).UpdateSubscriptionPlan(ctx, sqlc.UpdateSubscriptionPlanParams{
		ID:                item.ID,
		Code:              item.Code,
		Name:              item.Name,
		Description:       textFromPtr(item.Description),
		PriceAmount:       numericFromMoney(item.PriceAmount),
		PriceBasis:        item.PriceBasis,
		MinimumAmount:     numericFromMoney(item.MinimumAmount),
		IncludedEmployees: item.IncludedEmployees,
		OverageAmount:     numericFromMoney(item.OverageAmount),
		CurrencyCode:      item.CurrencyCode,
		BillingCycle:      item.BillingCycle,
		EmployeeLimit:     item.EmployeeLimit,
		TrialDays:         item.TrialDays,
		Visibility:        item.Visibility,
		IsActive:          item.IsActive,
		UpdatedBy:         uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update subscription plan", err, stringField("plan_id", item.ID.String()), stringField("plan_code", item.Code))
	}
	return mapSubscriptionPlan(row), nil
}

func (s *Store) DeleteSubscriptionPlan(ctx context.Context, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteSubscriptionPlan(ctx, sqlc.SoftDeleteSubscriptionPlanParams{ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete subscription plan", err, stringField("plan_id", id.String()))
	}
	return nil
}
