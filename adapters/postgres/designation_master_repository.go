package postgres

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateDesignationLevelCode(ctx context.Context, item *domain.DesignationLevelCode, actorID *uuid.UUID) (*domain.DesignationLevelCode, error) {
	row, err := s.getQueries(ctx).CreateDesignationLevelCode(ctx, sqlc.CreateDesignationLevelCodeParams{
		TenantID:    item.TenantID,
		Code:        item.Code,
		Label:       item.Label,
		Description: textFromPtr(item.Description),
		SortOrder:   item.SortOrder,
		CreatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create designation level code", err, tenantIDField(item.TenantID), stringField("level_code", item.Code))
	}
	return mapDesignationLevelCode(row), nil
}

func (s *Store) ListDesignationLevelCodes(ctx context.Context, tenantID uuid.UUID) ([]*domain.DesignationLevelCode, error) {
	rows, err := s.getQueries(ctx).ListDesignationLevelCodes(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list designation level codes", err, tenantIDField(tenantID))
	}
	return mapDesignationLevelCodes(rows), nil
}

func (s *Store) UpdateDesignationLevelCode(ctx context.Context, item *domain.DesignationLevelCode, actorID *uuid.UUID) (*domain.DesignationLevelCode, error) {
	row, err := s.getQueries(ctx).UpdateDesignationLevelCode(ctx, sqlc.UpdateDesignationLevelCodeParams{
		TenantID:    item.TenantID,
		ID:          item.ID,
		Code:        item.Code,
		Label:       item.Label,
		Description: textFromPtr(item.Description),
		SortOrder:   item.SortOrder,
		UpdatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update designation level code", err, tenantIDField(item.TenantID), stringField("level_code_id", item.ID.String()))
	}
	return mapDesignationLevelCode(row), nil
}

func (s *Store) DeleteDesignationLevelCode(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteDesignationLevelCode(ctx, sqlc.SoftDeleteDesignationLevelCodeParams{
		TenantID:  tenantID,
		ID:        id,
		UpdatedBy: uuidFromPtr(actorID),
	}); err != nil {
		return s.logDBError(ctx, "delete designation level code", err, tenantIDField(tenantID), stringField("level_code_id", id.String()))
	}
	return nil
}

func (s *Store) CreateDesignationSeniorityRank(ctx context.Context, item *domain.DesignationSeniorityRank, actorID *uuid.UUID) (*domain.DesignationSeniorityRank, error) {
	row, err := s.getQueries(ctx).CreateDesignationSeniorityRank(ctx, sqlc.CreateDesignationSeniorityRankParams{
		TenantID:    item.TenantID,
		RankValue:   item.RankValue,
		Label:       item.Label,
		Description: textFromPtr(item.Description),
		SortOrder:   item.SortOrder,
		CreatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create designation seniority rank", err, tenantIDField(item.TenantID), stringField("rank_value", strconv.Itoa(int(item.RankValue))))
	}
	return mapDesignationSeniorityRank(row), nil
}

func (s *Store) ListDesignationSeniorityRanks(ctx context.Context, tenantID uuid.UUID) ([]*domain.DesignationSeniorityRank, error) {
	rows, err := s.getQueries(ctx).ListDesignationSeniorityRanks(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list designation seniority ranks", err, tenantIDField(tenantID))
	}
	return mapDesignationSeniorityRanks(rows), nil
}

func (s *Store) UpdateDesignationSeniorityRank(ctx context.Context, item *domain.DesignationSeniorityRank, actorID *uuid.UUID) (*domain.DesignationSeniorityRank, error) {
	row, err := s.getQueries(ctx).UpdateDesignationSeniorityRank(ctx, sqlc.UpdateDesignationSeniorityRankParams{
		TenantID:    item.TenantID,
		ID:          item.ID,
		RankValue:   item.RankValue,
		Label:       item.Label,
		Description: textFromPtr(item.Description),
		SortOrder:   item.SortOrder,
		UpdatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update designation seniority rank", err, tenantIDField(item.TenantID), stringField("rank_id", item.ID.String()))
	}
	return mapDesignationSeniorityRank(row), nil
}

func (s *Store) DeleteDesignationSeniorityRank(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteDesignationSeniorityRank(ctx, sqlc.SoftDeleteDesignationSeniorityRankParams{
		TenantID:  tenantID,
		ID:        id,
		UpdatedBy: uuidFromPtr(actorID),
	}); err != nil {
		return s.logDBError(ctx, "delete designation seniority rank", err, tenantIDField(tenantID), stringField("rank_id", id.String()))
	}
	return nil
}
