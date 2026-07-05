package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateJobPosition(ctx context.Context, item *domain.JobPosition, actorID *uuid.UUID) (*domain.JobPosition, error) {
	row, err := s.getQueries(ctx).CreateJobPosition(ctx, sqlc.CreateJobPositionParams{TenantID: item.TenantID, Code: textFromPtr(item.Code), Title: item.Title, Level: textFromPtr(item.Level), Category: textFromPtr(item.Category), Description: textFromPtr(item.Description), DepartmentID: uuidFromPtr(item.DepartmentID), EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID), WorkMode: textFromPtr(item.WorkMode), TotalPosition: item.TotalPosition, BudgetedCost: numericFromFloatPtr(item.BudgetedCost), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create job position", fmt.Errorf("hrms: create job position: %w", err), tenantIDField(item.TenantID), stringField("title", item.Title))
	}
	return mapJobPosition(row), nil
}

func (s *Store) ListJobPositions(ctx context.Context, filter domain.JobPositionFilter) ([]*domain.JobPosition, error) {
	rows, err := s.getQueries(ctx).ListJobPositions(ctx, sqlc.ListJobPositionsParams{TenantID: filter.TenantID, DepartmentID: uuidFromPtr(filter.DepartmentID), EmploymentTypeID: uuidFromPtr(filter.EmploymentTypeID), WorkMode: textFromPtr(filter.WorkMode), Search: textFromPtr(filter.Search), Offset: filter.Offset, Limit: filter.Limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list job positions", err, tenantIDField(filter.TenantID))
	}
	return mapJobPositions(rows), nil
}

func (s *Store) CountJobPositions(ctx context.Context, filter domain.JobPositionFilter) (int64, error) {
	count, err := s.getQueries(ctx).CountJobPositions(ctx, sqlc.CountJobPositionsParams{TenantID: filter.TenantID, DepartmentID: uuidFromPtr(filter.DepartmentID), EmploymentTypeID: uuidFromPtr(filter.EmploymentTypeID), WorkMode: textFromPtr(filter.WorkMode), Search: textFromPtr(filter.Search)})
	if err != nil {
		return 0, s.logDBError(ctx, "count job positions", err, tenantIDField(filter.TenantID))
	}
	return count, nil
}

func (s *Store) GetJobPosition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobPosition, error) {
	row, err := s.getQueries(ctx).GetJobPosition(ctx, sqlc.GetJobPositionParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get job position", fmt.Errorf("hrms: get job position: %w", err), tenantIDField(tenantID), stringField("job_position_id", id.String()))
	}
	return mapJobPosition(row), nil
}

func (s *Store) UpdateJobPosition(ctx context.Context, item *domain.JobPosition, actorID *uuid.UUID) (*domain.JobPosition, error) {
	row, err := s.getQueries(ctx).UpdateJobPosition(ctx, sqlc.UpdateJobPositionParams{TenantID: item.TenantID, ID: item.ID, Code: textFromPtr(item.Code), Title: item.Title, Level: textFromPtr(item.Level), Category: textFromPtr(item.Category), Description: textFromPtr(item.Description), DepartmentID: uuidFromPtr(item.DepartmentID), EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID), WorkMode: textFromPtr(item.WorkMode), TotalPosition: item.TotalPosition, BudgetedCost: numericFromFloatPtr(item.BudgetedCost), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update job position", fmt.Errorf("hrms: update job position: %w", err), tenantIDField(item.TenantID), stringField("job_position_id", item.ID.String()))
	}
	return mapJobPosition(row), nil
}

func (s *Store) DeleteJobPosition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteJobPosition(ctx, sqlc.SoftDeleteJobPositionParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete job position", fmt.Errorf("hrms: delete job position: %w", err), tenantIDField(tenantID), stringField("job_position_id", id.String()))
	}
	return nil
}

func (s *Store) CreateJobPositionLocation(ctx context.Context, item *domain.JobPositionLocation, actorID *uuid.UUID) (*domain.JobPositionLocation, error) {
	row, err := s.getQueries(ctx).CreateJobPositionLocation(ctx, sqlc.CreateJobPositionLocationParams{TenantID: item.TenantID, JobPositionID: item.JobPositionID, Location: textFromPtr(item.Location), City: textFromPtr(item.City), State: textFromPtr(item.State), Country: textFromPtr(item.Country), IsRemote: item.IsRemote, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create job position location", fmt.Errorf("hrms: create job position location: %w", err), tenantIDField(item.TenantID), stringField("job_position_id", item.JobPositionID.String()))
	}
	return mapJobPositionLocation(row), nil
}

func (s *Store) ListJobPositionLocations(ctx context.Context, tenantID uuid.UUID, jobPositionID uuid.UUID) ([]*domain.JobPositionLocation, error) {
	rows, err := s.getQueries(ctx).ListJobPositionLocations(ctx, sqlc.ListJobPositionLocationsParams{TenantID: tenantID, JobPositionID: jobPositionID})
	if err != nil {
		return nil, s.logDBError(ctx, "list job position locations", err, tenantIDField(tenantID), stringField("job_position_id", jobPositionID.String()))
	}
	return mapJobPositionLocations(rows), nil
}

func (s *Store) GetJobPositionLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobPositionLocation, error) {
	row, err := s.getQueries(ctx).GetJobPositionLocation(ctx, sqlc.GetJobPositionLocationParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get job position location", fmt.Errorf("hrms: get job position location: %w", err), tenantIDField(tenantID), stringField("job_position_location_id", id.String()))
	}
	return mapJobPositionLocation(row), nil
}

func (s *Store) UpdateJobPositionLocation(ctx context.Context, item *domain.JobPositionLocation, actorID *uuid.UUID) (*domain.JobPositionLocation, error) {
	row, err := s.getQueries(ctx).UpdateJobPositionLocation(ctx, sqlc.UpdateJobPositionLocationParams{TenantID: item.TenantID, ID: item.ID, Location: textFromPtr(item.Location), City: textFromPtr(item.City), State: textFromPtr(item.State), Country: textFromPtr(item.Country), IsRemote: item.IsRemote, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update job position location", fmt.Errorf("hrms: update job position location: %w", err), tenantIDField(item.TenantID), stringField("job_position_location_id", item.ID.String()))
	}
	return mapJobPositionLocation(row), nil
}

func (s *Store) DeleteJobPositionLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteJobPositionLocation(ctx, sqlc.SoftDeleteJobPositionLocationParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete job position location", fmt.Errorf("hrms: delete job position location: %w", err), tenantIDField(tenantID), stringField("job_position_location_id", id.String()))
	}
	return nil
}
