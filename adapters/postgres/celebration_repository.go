package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateCelebrationType(ctx context.Context, item *domain.CelebrationType, actorID *uuid.UUID) (*domain.CelebrationType, error) {
	row, err := s.getQueries(ctx).CreateCelebrationType(ctx, sqlc.CreateCelebrationTypeParams{TenantID: item.TenantID, Name: item.Name, IsYearly: item.IsYearly, IsUserCelebration: item.IsUserCelebration, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create celebration type", fmt.Errorf("hrms: create celebration type: %w", err), tenantIDField(item.TenantID), stringField("celebration_type_name", item.Name))
	}
	return mapCelebrationType(row), nil
}

func (s *Store) ListCelebrationTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.CelebrationType, error) {
	rows, err := s.getQueries(ctx).ListCelebrationTypes(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list celebration types", err, tenantIDField(tenantID))
	}
	return mapCelebrationTypes(rows), nil
}

func (s *Store) GetCelebrationType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CelebrationType, error) {
	row, err := s.getQueries(ctx).GetCelebrationType(ctx, sqlc.GetCelebrationTypeParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get celebration type", fmt.Errorf("hrms: get celebration type: %w", err), tenantIDField(tenantID), stringField("celebration_type_id", id.String()))
	}
	return mapCelebrationType(row), nil
}

func (s *Store) UpdateCelebrationType(ctx context.Context, item *domain.CelebrationType, actorID *uuid.UUID) (*domain.CelebrationType, error) {
	row, err := s.getQueries(ctx).UpdateCelebrationType(ctx, sqlc.UpdateCelebrationTypeParams{TenantID: item.TenantID, ID: item.ID, Name: item.Name, IsYearly: item.IsYearly, IsUserCelebration: item.IsUserCelebration, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update celebration type", fmt.Errorf("hrms: update celebration type: %w", err), tenantIDField(item.TenantID), stringField("celebration_type_id", item.ID.String()))
	}
	return mapCelebrationType(row), nil
}

func (s *Store) DeleteCelebrationType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCelebrationType(ctx, sqlc.SoftDeleteCelebrationTypeParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete celebration type", fmt.Errorf("hrms: delete celebration type: %w", err), tenantIDField(tenantID), stringField("celebration_type_id", id.String()))
	}
	return nil
}

func (s *Store) CreateCelebration(ctx context.Context, item *domain.Celebration, actorID *uuid.UUID) (*domain.Celebration, error) {
	row, err := s.getQueries(ctx).CreateCelebration(ctx, sqlc.CreateCelebrationParams{TenantID: item.TenantID, BranchID: uuidFromPtr(item.BranchID), UserID: uuidFromPtr(item.UserID), CelebrationTypeID: item.CelebrationTypeID, CelebrationDate: dateFromPtr(item.CelebrationDate), CustomTitle: textFromPtr(item.CustomTitle), Description: textFromPtr(item.Description), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create celebration", fmt.Errorf("hrms: create celebration: %w", err), tenantIDField(item.TenantID), stringField("celebration_type_id", item.CelebrationTypeID.String()))
	}
	return mapCelebration(row), nil
}

func (s *Store) ListCelebrations(ctx context.Context, tenantID uuid.UUID) ([]*domain.Celebration, error) {
	rows, err := s.getQueries(ctx).ListCelebrations(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list celebrations", err, tenantIDField(tenantID))
	}
	return mapCelebrations(rows), nil
}

func (s *Store) GetCelebration(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Celebration, error) {
	row, err := s.getQueries(ctx).GetCelebration(ctx, sqlc.GetCelebrationParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get celebration", fmt.Errorf("hrms: get celebration: %w", err), tenantIDField(tenantID), stringField("celebration_id", id.String()))
	}
	return mapCelebration(row), nil
}

func (s *Store) UpdateCelebration(ctx context.Context, item *domain.Celebration, actorID *uuid.UUID) (*domain.Celebration, error) {
	row, err := s.getQueries(ctx).UpdateCelebration(ctx, sqlc.UpdateCelebrationParams{TenantID: item.TenantID, ID: item.ID, BranchID: uuidFromPtr(item.BranchID), UserID: uuidFromPtr(item.UserID), CelebrationTypeID: item.CelebrationTypeID, CelebrationDate: dateFromPtr(item.CelebrationDate), CustomTitle: textFromPtr(item.CustomTitle), Description: textFromPtr(item.Description), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update celebration", fmt.Errorf("hrms: update celebration: %w", err), tenantIDField(item.TenantID), stringField("celebration_id", item.ID.String()))
	}
	return mapCelebration(row), nil
}

func (s *Store) DeleteCelebration(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCelebration(ctx, sqlc.SoftDeleteCelebrationParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete celebration", fmt.Errorf("hrms: delete celebration: %w", err), tenantIDField(tenantID), stringField("celebration_id", id.String()))
	}
	return nil
}

func (s *Store) ListCelebrationsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.Celebration, error) {
	rows, err := s.getQueries(ctx).ListCelebrationsByUser(ctx, sqlc.ListCelebrationsByUserParams{TenantID: tenantID, UserID: uuidFromPtr(&userID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list celebrations by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapCelebrations(rows), nil
}
