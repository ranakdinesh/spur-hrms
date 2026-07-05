package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateTenantSubscription(ctx context.Context, cmd ports.TenantSubscriptionCommand) (*domain.TenantSubscription, error) {
	item, err := tenantSubscriptionFromCommand(cmd)
	if err != nil {
		s.logError("validate tenant subscription create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("status", cmd.Status))
		return nil, err
	}
	if err := s.applySubscriptionPlan(ctx, item); err != nil {
		s.logError("validate tenant subscription plan", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.subscriptions.CreateTenantSubscription(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create tenant subscription", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("subscription_id", result.ID.String()).Msg("hrms: tenant subscription created")
	return result, nil
}

func (s *TenantService) ListTenantSubscriptions(ctx context.Context, tenantID uuid.UUID) ([]*domain.TenantSubscription, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate tenant subscription list tenant", err)
		return nil, err
	}
	result, err := s.subscriptions.ListTenantSubscriptions(ctx, tenantID)
	if err != nil {
		s.logError("list tenant subscriptions", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetTenantSubscription(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.TenantSubscription, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate tenant subscription get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidSubscriptionID
		s.logError("validate tenant subscription get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.subscriptions.GetTenantSubscription(ctx, tenantID, id)
	if err != nil {
		s.logError("get tenant subscription", err, serviceTenantIDField(tenantID), serviceStringField("subscription_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetCurrentTenantSubscription(ctx context.Context, tenantID uuid.UUID) (*domain.TenantSubscription, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate current tenant subscription tenant", err)
		return nil, err
	}
	result, err := s.subscriptions.GetCurrentTenantSubscription(ctx, tenantID)
	if err != nil {
		s.logError("get current tenant subscription", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateTenantSubscription(ctx context.Context, cmd ports.TenantSubscriptionCommand) (*domain.TenantSubscription, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidSubscriptionID
		s.logError("validate tenant subscription update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := tenantSubscriptionFromCommand(cmd)
	if err != nil {
		s.logError("validate tenant subscription update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("subscription_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	if err := s.applySubscriptionPlan(ctx, item); err != nil {
		s.logError("validate tenant subscription plan", err, serviceTenantIDField(cmd.TenantID), serviceStringField("subscription_id", cmd.ID.String()))
		return nil, err
	}
	result, err := s.subscriptions.UpdateTenantSubscription(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update tenant subscription", err, serviceTenantIDField(cmd.TenantID), serviceStringField("subscription_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("subscription_id", result.ID.String()).Msg("hrms: tenant subscription updated")
	return result, nil
}

func (s *TenantService) DeleteTenantSubscription(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate tenant subscription delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidSubscriptionID
		s.logError("validate tenant subscription delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.subscriptions.DeleteTenantSubscription(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete tenant subscription", err, serviceTenantIDField(tenantID), serviceStringField("subscription_id", id.String()))
		return err
	}
	s.log.Info().Str("tenant_id", tenantID.String()).Str("subscription_id", id.String()).Msg("hrms: tenant subscription deactivated")
	return nil
}

func (s *TenantService) applySubscriptionPlan(ctx context.Context, item *domain.TenantSubscription) error {
	if item == nil || item.PlanID == nil {
		return nil
	}
	plan, err := s.subscriptionPlans.GetSubscriptionPlan(ctx, *item.PlanID)
	if err != nil {
		return err
	}
	if !plan.IsActive {
		return domain.ErrInactiveSubscriptionPlan
	}
	if item.MaxEmployees == 0 && plan.EmployeeLimit > 0 {
		item.MaxEmployees = plan.EmployeeLimit
	}
	return nil
}

func tenantSubscriptionFromCommand(cmd ports.TenantSubscriptionCommand) (*domain.TenantSubscription, error) {
	startDate, err := parseOptionalDate(cmd.StartDate)
	if err != nil {
		return nil, domain.ErrInvalidSubscriptionPeriod
	}
	endDate, err := parseOptionalDate(cmd.EndDate)
	if err != nil {
		return nil, domain.ErrInvalidSubscriptionPeriod
	}
	return domain.NewTenantSubscription(domain.TenantSubscriptionInput{
		TenantID:     cmd.TenantID,
		PlanID:       cmd.PlanID,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       cmd.Status,
		MaxEmployees: cmd.MaxEmployees,
	})
}

func parseOptionalDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
