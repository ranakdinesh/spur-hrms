package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateProjectSkillRequirement(ctx context.Context, cmd ports.ProjectSkillRequirementCommand) (*domain.ProjectSkillRequirement, error) {
	item, err := s.prepareProjectSkillRequirement(ctx, cmd)
	if err != nil {
		s.logError("validate project skill requirement", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.skillGaps.CreateProjectSkillRequirement(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create project skill requirement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("skill_id", cmd.SkillID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateProjectSkillRequirement(ctx context.Context, cmd ports.ProjectSkillRequirementCommand) (*domain.ProjectSkillRequirement, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidProjectSkillRequirement
	}
	if _, err := s.skillGaps.GetProjectSkillRequirement(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := s.prepareProjectSkillRequirement(ctx, cmd)
	if err != nil {
		s.logError("validate project skill requirement update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("requirement_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.skillGaps.UpdateProjectSkillRequirement(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update project skill requirement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("requirement_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetProjectSkillRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ProjectSkillRequirement, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidProjectSkillRequirement
	}
	return s.skillGaps.GetProjectSkillRequirement(ctx, tenantID, id)
}

func (s *TenantService) ListProjectSkillRequirements(ctx context.Context, filter domain.ProjectSkillRequirementFilter) ([]*domain.ProjectSkillRequirement, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.Search = domain.NormalizeProjectSkillSearch(filter.Search)
	return s.skillGaps.ListProjectSkillRequirements(ctx, filter)
}

func (s *TenantService) DeleteProjectSkillRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidProjectSkillRequirement
	}
	return s.skillGaps.DeleteProjectSkillRequirement(ctx, tenantID, id, actorID)
}

func (s *TenantService) ListProjectSkillGapRows(ctx context.Context, filter domain.ProjectSkillRequirementFilter) ([]*domain.ProjectSkillGapRow, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.skillGaps.ListProjectSkillGapRows(ctx, filter)
}

func (s *TenantService) ListSkillGapSummary(ctx context.Context, tenantID uuid.UUID, projectID *uuid.UUID) ([]*domain.SkillGapSummaryRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.skillGaps.ListSkillGapSummary(ctx, tenantID, projectID)
}

func (s *TenantService) ListSinglePersonSkillDependencies(ctx context.Context, tenantID uuid.UUID, projectID *uuid.UUID) ([]*domain.SinglePersonSkillDependency, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.skillGaps.ListSinglePersonSkillDependencies(ctx, tenantID, projectID)
}

func (s *TenantService) prepareProjectSkillRequirement(ctx context.Context, cmd ports.ProjectSkillRequirementCommand) (*domain.ProjectSkillRequirement, error) {
	if cmd.ProjectID != nil {
		if _, err := s.GetProject(ctx, cmd.TenantID, *cmd.ProjectID); err != nil {
			return nil, err
		}
	}
	if cmd.EngagementID != nil {
		if _, err := s.GetEngagement(ctx, cmd.TenantID, *cmd.EngagementID); err != nil {
			return nil, err
		}
	}
	if _, err := s.GetSkill(ctx, cmd.TenantID, cmd.SkillID); err != nil {
		return nil, err
	}
	return domain.NewProjectSkillRequirement(domain.ProjectSkillRequirementInput{
		TenantID:            cmd.TenantID,
		ProjectID:           cmd.ProjectID,
		EngagementID:        cmd.EngagementID,
		SkillID:             cmd.SkillID,
		RequiredProficiency: cmd.RequiredProficiency,
		MinYearsExperience:  cmd.MinYearsExperience,
		RequiredCount:       cmd.RequiredCount,
		Importance:          cmd.Importance,
		RequirementSource:   cmd.RequirementSource,
		Notes:               cmd.Notes,
		Metadata:            cmd.Metadata,
	})
}
