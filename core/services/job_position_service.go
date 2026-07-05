package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateJobPosition(ctx context.Context, cmd ports.JobPositionCommand) (*domain.JobPosition, error) {
	item, err := domain.NewJobPosition(domain.JobPositionInput{TenantID: cmd.TenantID, Code: cmd.Code, Title: cmd.Title, Level: cmd.Level, Category: cmd.Category, Description: cmd.Description, DepartmentID: cmd.DepartmentID, EmploymentTypeID: cmd.EmploymentTypeID, WorkMode: cmd.WorkMode, TotalPosition: cmd.TotalPosition, BudgetedCost: cmd.BudgetedCost})
	if err != nil {
		s.logError("validate job position create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("title", cmd.Title))
		return nil, err
	}
	if err := s.validateJobPositionReferences(ctx, item); err != nil {
		return nil, err
	}
	result, err := s.jobPositions.CreateJobPosition(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create job position", err, serviceTenantIDField(cmd.TenantID), serviceStringField("title", item.Title))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListJobPositions(ctx context.Context, filter domain.JobPositionFilter) (*domain.JobPositionPage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate job position list tenant", err)
		return nil, err
	}
	filter.Search = cleanStringPtr(filter.Search)
	filter.WorkMode = cleanStringPtr(filter.WorkMode)
	if filter.WorkMode != nil {
		mode, err := domain.ValidateJobWorkMode(filter.WorkMode)
		if err != nil {
			s.logError("validate job position list work mode", err, serviceTenantIDField(filter.TenantID))
			return nil, err
		}
		filter.WorkMode = mode
	}
	limit, offset := normalizeListWindow(filter.Limit, filter.Offset)
	filter.Limit = limit
	filter.Offset = offset
	items, err := s.jobPositions.ListJobPositions(ctx, filter)
	if err != nil {
		s.logError("list job positions", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.jobPositions.CountJobPositions(ctx, filter)
	if err != nil {
		s.logError("count job positions", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	page := &domain.JobPositionPage{Items: items, Total: total, Limit: limit, Offset: offset}
	if int64(offset)+int64(len(items)) < total {
		next := offset + limit
		page.NextOffset = &next
	}
	return page, nil
}

func (s *TenantService) GetJobPosition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobPosition, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate job position get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidJobPositionID
		s.logError("validate job position get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.jobPositions.GetJobPosition(ctx, tenantID, id)
	if err != nil {
		s.logError("get job position", err, serviceTenantIDField(tenantID), serviceStringField("job_position_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateJobPosition(ctx context.Context, cmd ports.JobPositionCommand) (*domain.JobPosition, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidJobPositionID
		s.logError("validate job position update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewJobPosition(domain.JobPositionInput{TenantID: cmd.TenantID, Code: cmd.Code, Title: cmd.Title, Level: cmd.Level, Category: cmd.Category, Description: cmd.Description, DepartmentID: cmd.DepartmentID, EmploymentTypeID: cmd.EmploymentTypeID, WorkMode: cmd.WorkMode, TotalPosition: cmd.TotalPosition, BudgetedCost: cmd.BudgetedCost})
	if err != nil {
		s.logError("validate job position update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_position_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	if err := s.validateJobPositionReferences(ctx, item); err != nil {
		return nil, err
	}
	result, err := s.jobPositions.UpdateJobPosition(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update job position", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_position_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteJobPosition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate job position delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidJobPositionID
		s.logError("validate job position delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.jobPositions.DeleteJobPosition(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete job position", err, serviceTenantIDField(tenantID), serviceStringField("job_position_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateJobPositionLocation(ctx context.Context, cmd ports.JobPositionLocationCommand) (*domain.JobPositionLocation, error) {
	item, err := domain.NewJobPositionLocation(domain.JobPositionLocationInput{TenantID: cmd.TenantID, JobPositionID: cmd.JobPositionID, Location: cmd.Location, City: cmd.City, State: cmd.State, Country: cmd.Country, IsRemote: cmd.IsRemote})
	if err != nil {
		s.logError("validate job position location create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_position_id", cmd.JobPositionID.String()))
		return nil, err
	}
	if _, err := s.GetJobPosition(ctx, cmd.TenantID, cmd.JobPositionID); err != nil {
		return nil, err
	}
	result, err := s.jobPositions.CreateJobPositionLocation(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create job position location", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_position_id", cmd.JobPositionID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListJobPositionLocations(ctx context.Context, tenantID uuid.UUID, jobPositionID uuid.UUID) ([]*domain.JobPositionLocation, error) {
	if _, err := s.GetJobPosition(ctx, tenantID, jobPositionID); err != nil {
		return nil, err
	}
	items, err := s.jobPositions.ListJobPositionLocations(ctx, tenantID, jobPositionID)
	if err != nil {
		s.logError("list job position locations", err, serviceTenantIDField(tenantID), serviceStringField("job_position_id", jobPositionID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetJobPositionLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobPositionLocation, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate job position location get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidJobLocationID
		s.logError("validate job position location get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.jobPositions.GetJobPositionLocation(ctx, tenantID, id)
	if err != nil {
		s.logError("get job position location", err, serviceTenantIDField(tenantID), serviceStringField("job_position_location_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateJobPositionLocation(ctx context.Context, cmd ports.JobPositionLocationCommand) (*domain.JobPositionLocation, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidJobLocationID
		s.logError("validate job position location update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetJobPositionLocation(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewJobPositionLocation(domain.JobPositionLocationInput{TenantID: cmd.TenantID, JobPositionID: existing.JobPositionID, Location: cmd.Location, City: cmd.City, State: cmd.State, Country: cmd.Country, IsRemote: cmd.IsRemote})
	if err != nil {
		s.logError("validate job position location update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_position_location_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.jobPositions.UpdateJobPositionLocation(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update job position location", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_position_location_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteJobPositionLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetJobPositionLocation(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.jobPositions.DeleteJobPositionLocation(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete job position location", err, serviceTenantIDField(tenantID), serviceStringField("job_position_location_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) validateJobPositionReferences(ctx context.Context, item *domain.JobPosition) error {
	if item.DepartmentID != nil {
		if _, err := s.GetDepartment(ctx, item.TenantID, *item.DepartmentID); err != nil {
			return err
		}
	}
	if item.EmploymentTypeID != nil {
		if _, err := s.GetEmploymentType(ctx, item.TenantID, *item.EmploymentTypeID); err != nil {
			return err
		}
	}
	return nil
}
