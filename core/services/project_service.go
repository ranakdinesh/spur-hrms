package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateProject(ctx context.Context, cmd ports.ProjectCommand) (*domain.Project, error) {
	input, err := s.prepareProjectCommand(ctx, cmd, nil)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewProject(input)
	if err != nil {
		s.logError("validate project create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_name", cmd.Name))
		return nil, err
	}
	result, err := s.projects.CreateProject(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create project", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_name", item.Name))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("project_id", result.ID.String()).Msg("hrms: project created")
	return result, nil
}

func (s *TenantService) UpdateProject(ctx context.Context, cmd ports.ProjectCommand) (*domain.Project, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidProjectID
		s.logError("validate project update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetProject(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existing.Status == domain.ProjectStatusCompleted || existing.Status == domain.ProjectStatusCancelled {
		err := domain.ErrInvalidProjectStatus
		s.logError("validate project update state", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_id", cmd.ID.String()))
		return nil, err
	}
	input, err := s.prepareProjectCommand(ctx, cmd, existing)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewProject(input)
	if err != nil {
		s.logError("validate project update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.projects.UpdateProject(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update project", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateProjectStatus(ctx context.Context, cmd ports.ProjectStatusCommand) (*domain.Project, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project status tenant", err)
		return nil, err
	}
	if cmd.ProjectID == uuid.Nil {
		err := domain.ErrInvalidProjectID
		s.logError("validate project status id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status, err := domain.ValidateProjectStatus(cmd.Status)
	if err != nil {
		s.logError("validate project status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_id", cmd.ProjectID.String()))
		return nil, err
	}
	item, err := s.GetProject(ctx, cmd.TenantID, cmd.ProjectID)
	if err != nil {
		return nil, err
	}
	if item.Status == domain.ProjectStatusCancelled && status != domain.ProjectStatusCancelled {
		err := domain.ErrInvalidProjectStatus
		s.logError("validate cancelled project status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_id", cmd.ProjectID.String()))
		return nil, err
	}
	item.Status = status
	if status == domain.ProjectStatusCompleted {
		now := time.Now().UTC()
		item.CompletedAt = &now
	} else {
		item.CompletedAt = nil
	}
	result, err := s.projects.UpdateProjectStatus(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update project status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("project_id", cmd.ProjectID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetProject(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Project, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidProjectID
		s.logError("validate project get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.projects.GetProject(ctx, tenantID, id)
	if err != nil {
		s.logError("get project", err, serviceTenantIDField(tenantID), serviceStringField("project_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListProjects(ctx context.Context, filter domain.ProjectFilter) ([]*domain.ProjectListItem, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project list tenant", err)
		return nil, err
	}
	result, err := s.projects.ListProjects(ctx, filter)
	if err != nil {
		s.logError("list projects", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteProject(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidProjectID
		s.logError("validate project delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.projects.DeleteProject(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete project", err, serviceTenantIDField(tenantID), serviceStringField("project_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) prepareProjectCommand(ctx context.Context, cmd ports.ProjectCommand, existing *domain.Project) (domain.ProjectInput, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate project tenant", err)
		return domain.ProjectInput{}, err
	}
	startDate, err := parseEngagementDate(cmd.StartDate)
	if err != nil {
		s.logError("parse project start date", err, serviceTenantIDField(cmd.TenantID))
		return domain.ProjectInput{}, err
	}
	dueDate, err := parseEngagementDate(cmd.DueDate)
	if err != nil {
		s.logError("parse project due date", err, serviceTenantIDField(cmd.TenantID))
		return domain.ProjectInput{}, err
	}
	if cmd.DepartmentID != nil {
		if _, err := s.GetDepartment(ctx, cmd.TenantID, *cmd.DepartmentID); err != nil {
			return domain.ProjectInput{}, err
		}
	}
	if cmd.BranchID != nil {
		if _, err := s.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			return domain.ProjectInput{}, err
		}
	}
	completedAt := (*time.Time)(nil)
	status := cmd.Status
	if existing != nil {
		completedAt = existing.CompletedAt
		if status == "" {
			status = existing.Status
		}
		if cmd.CurrencyCode == "" {
			cmd.CurrencyCode = existing.CurrencyCode
		}
		if cmd.BillingType == "" {
			cmd.BillingType = existing.BillingType
		}
		if cmd.Priority == "" {
			cmd.Priority = existing.Priority
		}
	}
	return domain.ProjectInput{
		TenantID:         cmd.TenantID,
		ProjectCode:      cmd.ProjectCode,
		Name:             cmd.Name,
		Description:      cmd.Description,
		Status:           status,
		DepartmentID:     cmd.DepartmentID,
		BranchID:         cmd.BranchID,
		ProjectManagerID: cmd.ProjectManagerID,
		StartDate:        startDate,
		DueDate:          dueDate,
		CompletedAt:      completedAt,
		BudgetAmount:     cmd.BudgetAmount,
		CurrencyCode:     cmd.CurrencyCode,
		BillingType:      cmd.BillingType,
		ClientLabel:      cmd.ClientLabel,
		Priority:         cmd.Priority,
		Notes:            cmd.Notes,
		Metadata:         cmd.Metadata,
	}, nil
}
