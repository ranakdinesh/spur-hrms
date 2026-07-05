package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateTenantSubscription(ctx context.Context, item *domain.TenantSubscription, actorID *uuid.UUID) (*domain.TenantSubscription, error) {
	row, err := s.getQueries(ctx).CreateTenantSubscription(ctx, sqlc.CreateTenantSubscriptionParams{
		TenantID:     item.TenantID,
		PlanID:       uuidFromPtr(item.PlanID),
		StartDate:    dateFromPtr(item.StartDate),
		EndDate:      dateFromPtr(item.EndDate),
		Status:       item.Status,
		MaxEmployees: item.MaxEmployees,
		CreatedBy:    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create tenant subscription", err, tenantIDField(item.TenantID), stringField("subscription_status", item.Status))
	}
	return mapTenantSubscription(row), nil
}

func (s *Store) ListTenantSubscriptions(ctx context.Context, tenantID uuid.UUID) ([]*domain.TenantSubscription, error) {
	rows, err := s.getQueries(ctx).ListTenantSubscriptions(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list tenant subscriptions", err, tenantIDField(tenantID))
	}
	return mapTenantSubscriptions(rows), nil
}

func (s *Store) GetTenantSubscription(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.TenantSubscription, error) {
	row, err := s.getQueries(ctx).GetTenantSubscription(ctx, sqlc.GetTenantSubscriptionParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get tenant subscription", err, tenantIDField(tenantID), stringField("subscription_id", id.String()))
	}
	return mapTenantSubscription(row), nil
}

func (s *Store) GetCurrentTenantSubscription(ctx context.Context, tenantID uuid.UUID) (*domain.TenantSubscription, error) {
	row, err := s.getQueries(ctx).GetCurrentTenantSubscription(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get current tenant subscription", err, tenantIDField(tenantID))
	}
	return mapTenantSubscription(row), nil
}

func (s *Store) UpdateTenantSubscription(ctx context.Context, item *domain.TenantSubscription, actorID *uuid.UUID) (*domain.TenantSubscription, error) {
	row, err := s.getQueries(ctx).UpdateTenantSubscription(ctx, sqlc.UpdateTenantSubscriptionParams{
		TenantID:     item.TenantID,
		ID:           item.ID,
		PlanID:       uuidFromPtr(item.PlanID),
		StartDate:    dateFromPtr(item.StartDate),
		EndDate:      dateFromPtr(item.EndDate),
		Status:       item.Status,
		MaxEmployees: item.MaxEmployees,
		UpdatedBy:    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update tenant subscription", err, tenantIDField(item.TenantID), stringField("subscription_id", item.ID.String()))
	}
	return mapTenantSubscription(row), nil
}

func (s *Store) DeleteTenantSubscription(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteTenantSubscription(ctx, sqlc.SoftDeleteTenantSubscriptionParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete tenant subscription", err, tenantIDField(tenantID), stringField("subscription_id", id.String()))
	}
	return nil
}
