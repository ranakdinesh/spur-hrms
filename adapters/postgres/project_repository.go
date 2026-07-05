package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateProject(ctx context.Context, item *domain.Project, actorID *uuid.UUID) (*domain.Project, error) {
	row, err := s.getQueries(ctx).CreateProject(ctx, projectParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create project", err, tenantIDField(item.TenantID), stringField("project_name", item.Name))
	}
	return mapProject(row), nil
}

func (s *Store) UpdateProject(ctx context.Context, item *domain.Project, actorID *uuid.UUID) (*domain.Project, error) {
	params := updateProjectParams(item, actorID)
	row, err := s.getQueries(ctx).UpdateProject(ctx, params)
	if err != nil {
		return nil, s.logDBError(ctx, "update project", err, tenantIDField(item.TenantID), stringField("project_id", item.ID.String()))
	}
	return mapProject(row), nil
}

func (s *Store) UpdateProjectStatus(ctx context.Context, item *domain.Project, actorID *uuid.UUID) (*domain.Project, error) {
	row, err := s.getQueries(ctx).UpdateProjectStatus(ctx, sqlc.UpdateProjectStatusParams{
		TenantID:    item.TenantID,
		ID:          item.ID,
		Status:      item.Status,
		CompletedAt: timestamptzFromPtr(item.CompletedAt),
		UpdatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update project status", err, tenantIDField(item.TenantID), stringField("project_id", item.ID.String()))
	}
	return mapProject(row), nil
}

func (s *Store) GetProject(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Project, error) {
	row, err := s.getQueries(ctx).GetProject(ctx, sqlc.GetProjectParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get project", err, tenantIDField(tenantID), stringField("project_id", id.String()))
	}
	return mapProject(row), nil
}

func (s *Store) ListProjects(ctx context.Context, filter domain.ProjectFilter) ([]*domain.ProjectListItem, error) {
	rows, err := s.getQueries(ctx).ListProjects(ctx, sqlc.ListProjectsParams{
		TenantID:         filter.TenantID,
		Status:           textFromPtr(filter.Status),
		DepartmentID:     uuidFromPtr(filter.DepartmentID),
		BranchID:         uuidFromPtr(filter.BranchID),
		ProjectManagerID: uuidFromPtr(filter.ProjectManagerID),
		Search:           textFromPtr(filter.Search),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list projects", err, tenantIDField(filter.TenantID))
	}
	return mapProjectListItems(rows), nil
}

func (s *Store) DeleteProject(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteProject(ctx, sqlc.SoftDeleteProjectParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete project", err, tenantIDField(tenantID), stringField("project_id", id.String()))
	}
	return nil
}

func (s *Store) CreateProjectMilestone(ctx context.Context, item *domain.ProjectMilestone, actorID *uuid.UUID) (*domain.ProjectMilestone, error) {
	row, err := s.getQueries(ctx).CreateProjectMilestone(ctx, createProjectMilestoneParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create project milestone", err, tenantIDField(item.TenantID), stringField("project_id", item.ProjectID.String()))
	}
	return mapProjectMilestone(row), nil
}

func (s *Store) UpdateProjectMilestone(ctx context.Context, item *domain.ProjectMilestone, actorID *uuid.UUID) (*domain.ProjectMilestone, error) {
	row, err := s.getQueries(ctx).UpdateProjectMilestone(ctx, updateProjectMilestoneParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "update project milestone", err, tenantIDField(item.TenantID), stringField("milestone_id", item.ID.String()))
	}
	return mapProjectMilestone(row), nil
}

func (s *Store) UpdateProjectMilestoneStatus(ctx context.Context, item *domain.ProjectMilestone, actorID *uuid.UUID) (*domain.ProjectMilestone, error) {
	row, err := s.getQueries(ctx).UpdateProjectMilestoneStatus(ctx, sqlc.UpdateProjectMilestoneStatusParams{
		TenantID:      item.TenantID,
		ID:            item.ID,
		Status:        item.Status,
		SubmittedAt:   timestamptzFromPtr(item.SubmittedAt),
		SubmittedBy:   uuidFromPtr(item.SubmittedBy),
		AcceptedAt:    timestamptzFromPtr(item.AcceptedAt),
		AcceptedBy:    uuidFromPtr(item.AcceptedBy),
		RejectedAt:    timestamptzFromPtr(item.RejectedAt),
		RejectedBy:    uuidFromPtr(item.RejectedBy),
		ReviewComment: textFromPtr(item.ReviewComment),
		UpdatedBy:     uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update project milestone status", err, tenantIDField(item.TenantID), stringField("milestone_id", item.ID.String()))
	}
	return mapProjectMilestone(row), nil
}

func (s *Store) GetProjectMilestone(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ProjectMilestone, error) {
	row, err := s.getQueries(ctx).GetProjectMilestone(ctx, sqlc.GetProjectMilestoneParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get project milestone", err, tenantIDField(tenantID), stringField("milestone_id", id.String()))
	}
	return mapProjectMilestone(row), nil
}

func (s *Store) ListProjectMilestones(ctx context.Context, filter domain.ProjectMilestoneFilter) ([]*domain.ProjectMilestoneListItem, error) {
	rows, err := s.getQueries(ctx).ListProjectMilestones(ctx, sqlc.ListProjectMilestonesParams{
		TenantID:     filter.TenantID,
		ProjectID:    uuidFromPtr(filter.ProjectID),
		EngagementID: uuidFromPtr(filter.EngagementID),
		Status:       textFromPtr(filter.Status),
		Search:       textFromPtr(filter.Search),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list project milestones", err, tenantIDField(filter.TenantID))
	}
	return mapProjectMilestoneListItems(rows), nil
}

func (s *Store) DeleteProjectMilestone(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteProjectMilestone(ctx, sqlc.SoftDeleteProjectMilestoneParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete project milestone", err, tenantIDField(tenantID), stringField("milestone_id", id.String()))
	}
	return nil
}

func (s *Store) CreateProjectMilestoneEvent(ctx context.Context, event *domain.ProjectMilestoneEvent) (*domain.ProjectMilestoneEvent, error) {
	row, err := s.getQueries(ctx).CreateProjectMilestoneEvent(ctx, sqlc.CreateProjectMilestoneEventParams{
		TenantID:    event.TenantID,
		ProjectID:   event.ProjectID,
		MilestoneID: event.MilestoneID,
		EventType:   event.EventType,
		FromStatus:  textFromPtr(event.FromStatus),
		ToStatus:    textFromPtr(event.ToStatus),
		Comment:     textFromPtr(event.Comment),
		ActorID:     uuidFromPtr(event.ActorID),
		Metadata:    jsonBytesFromRaw(event.Metadata),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create project milestone event", err, tenantIDField(event.TenantID), stringField("milestone_id", event.MilestoneID.String()))
	}
	return mapProjectMilestoneEvent(row), nil
}

func (s *Store) ListProjectMilestoneEvents(ctx context.Context, tenantID uuid.UUID, milestoneID uuid.UUID) ([]*domain.ProjectMilestoneEvent, error) {
	rows, err := s.getQueries(ctx).ListProjectMilestoneEvents(ctx, sqlc.ListProjectMilestoneEventsParams{TenantID: tenantID, MilestoneID: milestoneID})
	if err != nil {
		return nil, s.logDBError(ctx, "list project milestone events", err, tenantIDField(tenantID), stringField("milestone_id", milestoneID.String()))
	}
	return mapProjectMilestoneEvents(rows), nil
}

func projectParams(item *domain.Project, actorID *uuid.UUID) sqlc.CreateProjectParams {
	return sqlc.CreateProjectParams{
		TenantID:         item.TenantID,
		ProjectCode:      textFromPtr(item.ProjectCode),
		Name:             item.Name,
		Description:      textFromPtr(item.Description),
		Status:           item.Status,
		DepartmentID:     uuidFromPtr(item.DepartmentID),
		BranchID:         uuidFromPtr(item.BranchID),
		ProjectManagerID: uuidFromPtr(item.ProjectManagerID),
		StartDate:        dateFromPtr(item.StartDate),
		DueDate:          dateFromPtr(item.DueDate),
		CompletedAt:      timestamptzFromPtr(item.CompletedAt),
		BudgetAmount:     numericFromEngagementDecimalPtr(item.BudgetAmount),
		CurrencyCode:     item.CurrencyCode,
		BillingType:      item.BillingType,
		ClientLabel:      textFromPtr(item.ClientLabel),
		Priority:         item.Priority,
		Notes:            textFromPtr(item.Notes),
		Metadata:         jsonBytesFromRaw(item.Metadata),
		CreatedBy:        uuidFromPtr(actorID),
	}
}

func updateProjectParams(item *domain.Project, actorID *uuid.UUID) sqlc.UpdateProjectParams {
	return sqlc.UpdateProjectParams{
		TenantID:         item.TenantID,
		ID:               item.ID,
		ProjectCode:      textFromPtr(item.ProjectCode),
		Name:             item.Name,
		Description:      textFromPtr(item.Description),
		Status:           item.Status,
		DepartmentID:     uuidFromPtr(item.DepartmentID),
		BranchID:         uuidFromPtr(item.BranchID),
		ProjectManagerID: uuidFromPtr(item.ProjectManagerID),
		StartDate:        dateFromPtr(item.StartDate),
		DueDate:          dateFromPtr(item.DueDate),
		CompletedAt:      timestamptzFromPtr(item.CompletedAt),
		BudgetAmount:     numericFromEngagementDecimalPtr(item.BudgetAmount),
		CurrencyCode:     item.CurrencyCode,
		BillingType:      item.BillingType,
		ClientLabel:      textFromPtr(item.ClientLabel),
		Priority:         item.Priority,
		Notes:            textFromPtr(item.Notes),
		Metadata:         jsonBytesFromRaw(item.Metadata),
		UpdatedBy:        uuidFromPtr(actorID),
	}
}

func createProjectMilestoneParams(item *domain.ProjectMilestone, actorID *uuid.UUID) sqlc.CreateProjectMilestoneParams {
	return sqlc.CreateProjectMilestoneParams{
		TenantID:           item.TenantID,
		ProjectID:          item.ProjectID,
		EngagementID:       uuidFromPtr(item.EngagementID),
		MilestoneCode:      textFromPtr(item.MilestoneCode),
		Title:              item.Title,
		Description:        textFromPtr(item.Description),
		AcceptanceCriteria: textFromPtr(item.AcceptanceCriteria),
		DueDate:            dateFromPtr(item.DueDate),
		Status:             item.Status,
		Amount:             numericFromEngagementDecimalPtr(item.Amount),
		CurrencyCode:       item.CurrencyCode,
		PaymentTrigger:     jsonBytesFromRaw(item.PaymentTrigger),
		SubmittedAt:        timestamptzFromPtr(item.SubmittedAt),
		SubmittedBy:        uuidFromPtr(item.SubmittedBy),
		AcceptedAt:         timestamptzFromPtr(item.AcceptedAt),
		AcceptedBy:         uuidFromPtr(item.AcceptedBy),
		RejectedAt:         timestamptzFromPtr(item.RejectedAt),
		RejectedBy:         uuidFromPtr(item.RejectedBy),
		ReviewComment:      textFromPtr(item.ReviewComment),
		Notes:              textFromPtr(item.Notes),
		Metadata:           jsonBytesFromRaw(item.Metadata),
		CreatedBy:          uuidFromPtr(actorID),
	}
}

func updateProjectMilestoneParams(item *domain.ProjectMilestone, actorID *uuid.UUID) sqlc.UpdateProjectMilestoneParams {
	return sqlc.UpdateProjectMilestoneParams{
		TenantID:           item.TenantID,
		ID:                 item.ID,
		ProjectID:          item.ProjectID,
		EngagementID:       uuidFromPtr(item.EngagementID),
		MilestoneCode:      textFromPtr(item.MilestoneCode),
		Title:              item.Title,
		Description:        textFromPtr(item.Description),
		AcceptanceCriteria: textFromPtr(item.AcceptanceCriteria),
		DueDate:            dateFromPtr(item.DueDate),
		Status:             item.Status,
		Amount:             numericFromEngagementDecimalPtr(item.Amount),
		CurrencyCode:       item.CurrencyCode,
		PaymentTrigger:     jsonBytesFromRaw(item.PaymentTrigger),
		SubmittedAt:        timestamptzFromPtr(item.SubmittedAt),
		SubmittedBy:        uuidFromPtr(item.SubmittedBy),
		AcceptedAt:         timestamptzFromPtr(item.AcceptedAt),
		AcceptedBy:         uuidFromPtr(item.AcceptedBy),
		RejectedAt:         timestamptzFromPtr(item.RejectedAt),
		RejectedBy:         uuidFromPtr(item.RejectedBy),
		ReviewComment:      textFromPtr(item.ReviewComment),
		Notes:              textFromPtr(item.Notes),
		Metadata:           jsonBytesFromRaw(item.Metadata),
		UpdatedBy:          uuidFromPtr(actorID),
	}
}
