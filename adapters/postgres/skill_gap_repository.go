package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateProjectSkillRequirement(ctx context.Context, item *domain.ProjectSkillRequirement, actorID *uuid.UUID) (*domain.ProjectSkillRequirement, error) {
	row, err := s.getQueries(ctx).CreateProjectSkillRequirement(ctx, projectSkillRequirementCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create project skill requirement", err, tenantIDField(item.TenantID), stringField("skill_id", item.SkillID.String()))
	}
	return mapProjectSkillRequirement(row), nil
}

func (s *Store) UpdateProjectSkillRequirement(ctx context.Context, item *domain.ProjectSkillRequirement, actorID *uuid.UUID) (*domain.ProjectSkillRequirement, error) {
	row, err := s.getQueries(ctx).UpdateProjectSkillRequirement(ctx, projectSkillRequirementUpdateParams(item, actorID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrProjectSkillRequirementNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update project skill requirement", err, tenantIDField(item.TenantID), stringField("requirement_id", item.ID.String()))
	}
	return mapProjectSkillRequirement(row), nil
}

func (s *Store) GetProjectSkillRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ProjectSkillRequirement, error) {
	row, err := s.getQueries(ctx).GetProjectSkillRequirement(ctx, sqlc.GetProjectSkillRequirementParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrProjectSkillRequirementNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get project skill requirement", err, tenantIDField(tenantID), stringField("requirement_id", id.String()))
	}
	return mapProjectSkillRequirement(row), nil
}

func (s *Store) ListProjectSkillRequirements(ctx context.Context, filter domain.ProjectSkillRequirementFilter) ([]*domain.ProjectSkillRequirement, error) {
	rows, err := s.getQueries(ctx).ListProjectSkillRequirements(ctx, sqlc.ListProjectSkillRequirementsParams{TenantID: filter.TenantID, ProjectID: uuidFromPtr(filter.ProjectID), EngagementID: uuidFromPtr(filter.EngagementID), SkillID: uuidFromPtr(filter.SkillID), Importance: textFromPtr(filter.Importance), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list project skill requirements", err, tenantIDField(filter.TenantID))
	}
	return mapProjectSkillRequirementRows(rows), nil
}

func (s *Store) DeleteProjectSkillRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteProjectSkillRequirement(ctx, sqlc.SoftDeleteProjectSkillRequirementParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete project skill requirement", err, tenantIDField(tenantID), stringField("requirement_id", id.String()))
	}
	return nil
}

func (s *Store) ListProjectSkillGapRows(ctx context.Context, filter domain.ProjectSkillRequirementFilter) ([]*domain.ProjectSkillGapRow, error) {
	rows, err := s.getQueries(ctx).ListProjectSkillGapRows(ctx, sqlc.ListProjectSkillGapRowsParams{TenantID: filter.TenantID, ProjectID: uuidFromPtr(filter.ProjectID), EngagementID: uuidFromPtr(filter.EngagementID), Importance: textFromPtr(filter.Importance)})
	if err != nil {
		return nil, s.logDBError(ctx, "list project skill gap rows", err, tenantIDField(filter.TenantID))
	}
	return mapProjectSkillGapRows(rows), nil
}

func (s *Store) ListSkillGapSummary(ctx context.Context, tenantID uuid.UUID, projectID *uuid.UUID) ([]*domain.SkillGapSummaryRow, error) {
	rows, err := s.getQueries(ctx).ListSkillGapSummary(ctx, sqlc.ListSkillGapSummaryParams{TenantID: tenantID, ProjectID: uuidFromPtr(projectID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list skill gap summary", err, tenantIDField(tenantID))
	}
	return mapSkillGapSummaryRows(rows), nil
}

func (s *Store) ListSinglePersonSkillDependencies(ctx context.Context, tenantID uuid.UUID, projectID *uuid.UUID) ([]*domain.SinglePersonSkillDependency, error) {
	rows, err := s.getQueries(ctx).ListSinglePersonSkillDependencies(ctx, sqlc.ListSinglePersonSkillDependenciesParams{TenantID: tenantID, ProjectID: uuidFromPtr(projectID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list single person skill dependencies", err, tenantIDField(tenantID))
	}
	return mapSinglePersonSkillDependencyRows(rows), nil
}

func projectSkillRequirementCreateParams(item *domain.ProjectSkillRequirement, actorID *uuid.UUID) sqlc.CreateProjectSkillRequirementParams {
	return sqlc.CreateProjectSkillRequirementParams{TenantID: item.TenantID, ProjectID: uuidFromPtr(item.ProjectID), EngagementID: uuidFromPtr(item.EngagementID), SkillID: item.SkillID, RequiredProficiency: item.RequiredProficiency, MinYearsExperience: numericFromSkillFloat(item.MinYearsExperience), RequiredCount: item.RequiredCount, Importance: item.Importance, RequirementSource: item.RequirementSource, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func projectSkillRequirementUpdateParams(item *domain.ProjectSkillRequirement, actorID *uuid.UUID) sqlc.UpdateProjectSkillRequirementParams {
	return sqlc.UpdateProjectSkillRequirementParams{TenantID: item.TenantID, ID: item.ID, ProjectID: uuidFromPtr(item.ProjectID), EngagementID: uuidFromPtr(item.EngagementID), SkillID: item.SkillID, RequiredProficiency: item.RequiredProficiency, MinYearsExperience: numericFromSkillFloat(item.MinYearsExperience), RequiredCount: item.RequiredCount, Importance: item.Importance, RequirementSource: item.RequirementSource, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}
