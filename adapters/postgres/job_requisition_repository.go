package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateJobRequisition(ctx context.Context, item *domain.JobRequisition, actorID *uuid.UUID) (*domain.JobRequisition, error) {
	row, err := s.getQueries(ctx).CreateJobRequisition(ctx, sqlc.CreateJobRequisitionParams{TenantID: item.TenantID, JobPositionID: item.JobPositionID, Code: textFromPtr(item.Code), Title: item.Title, Level: textFromPtr(item.Level), Category: textFromPtr(item.Category), DepartmentID: uuidFromPtr(item.DepartmentID), EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID), Description: textFromPtr(item.Description), WorkMode: textFromPtr(item.WorkMode), TotalOpenings: item.TotalOpenings, ReasonForHire: textFromPtr(item.ReasonForHire), MinSalary: numericFromFloatPtr(item.MinSalary), MaxSalary: numericFromFloatPtr(item.MaxSalary), Currency: item.Currency, TargetHireDate: dateFromPtr(item.TargetHireDate), ExpectedClosureDate: dateFromPtr(item.ExpectedClosureDate), RequestedBy: item.RequestedBy, RequestedDate: dateFromPtr(item.RequestedDate), Priority: textFromPtr(item.Priority), Status: item.Status, Notes: textFromPtr(item.Notes), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create job requisition", fmt.Errorf("hrms: create job requisition: %w", err), tenantIDField(item.TenantID), stringField("title", item.Title))
	}
	return mapJobRequisition(row), nil
}

func (s *Store) ListJobRequisitions(ctx context.Context, filter domain.JobRequisitionFilter) ([]*domain.JobRequisition, error) {
	rows, err := s.getQueries(ctx).ListJobRequisitions(ctx, sqlc.ListJobRequisitionsParams{TenantID: filter.TenantID, Status: textFromPtr(filter.Status), JobPositionID: uuidFromPtr(filter.JobPositionID), DepartmentID: uuidFromPtr(filter.DepartmentID), Search: textFromPtr(filter.Search), Offset: filter.Offset, Limit: filter.Limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list job requisitions", err, tenantIDField(filter.TenantID))
	}
	return mapJobRequisitions(rows), nil
}

func (s *Store) CountJobRequisitions(ctx context.Context, filter domain.JobRequisitionFilter) (int64, error) {
	count, err := s.getQueries(ctx).CountJobRequisitions(ctx, sqlc.CountJobRequisitionsParams{TenantID: filter.TenantID, Status: textFromPtr(filter.Status), JobPositionID: uuidFromPtr(filter.JobPositionID), DepartmentID: uuidFromPtr(filter.DepartmentID), Search: textFromPtr(filter.Search)})
	if err != nil {
		return 0, s.logDBError(ctx, "count job requisitions", err, tenantIDField(filter.TenantID))
	}
	return count, nil
}

func (s *Store) GetJobRequisition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobRequisition, error) {
	row, err := s.getQueries(ctx).GetJobRequisition(ctx, sqlc.GetJobRequisitionParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get job requisition", fmt.Errorf("hrms: get job requisition: %w", err), tenantIDField(tenantID), stringField("job_requisition_id", id.String()))
	}
	return mapJobRequisition(row), nil
}

func (s *Store) UpdateJobRequisition(ctx context.Context, item *domain.JobRequisition, actorID *uuid.UUID) (*domain.JobRequisition, error) {
	row, err := s.getQueries(ctx).UpdateJobRequisition(ctx, sqlc.UpdateJobRequisitionParams{TenantID: item.TenantID, ID: item.ID, JobPositionID: item.JobPositionID, Code: textFromPtr(item.Code), Title: item.Title, Level: textFromPtr(item.Level), Category: textFromPtr(item.Category), DepartmentID: uuidFromPtr(item.DepartmentID), EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID), Description: textFromPtr(item.Description), WorkMode: textFromPtr(item.WorkMode), TotalOpenings: item.TotalOpenings, ReasonForHire: textFromPtr(item.ReasonForHire), MinSalary: numericFromFloatPtr(item.MinSalary), MaxSalary: numericFromFloatPtr(item.MaxSalary), Currency: item.Currency, TargetHireDate: dateFromPtr(item.TargetHireDate), ExpectedClosureDate: dateFromPtr(item.ExpectedClosureDate), RequestedBy: item.RequestedBy, RequestedDate: dateFromPtr(item.RequestedDate), Priority: textFromPtr(item.Priority), Notes: textFromPtr(item.Notes), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update job requisition", fmt.Errorf("hrms: update job requisition: %w", err), tenantIDField(item.TenantID), stringField("job_requisition_id", item.ID.String()))
	}
	return mapJobRequisition(row), nil
}

func (s *Store) UpdateJobRequisitionStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, notes *string, actorID *uuid.UUID) (*domain.JobRequisition, error) {
	row, err := s.getQueries(ctx).UpdateJobRequisitionStatus(ctx, sqlc.UpdateJobRequisitionStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID), Notes: textFromPtr(notes)})
	if err != nil {
		return nil, s.logDBError(ctx, "update job requisition status", fmt.Errorf("hrms: update job requisition status: %w", err), tenantIDField(tenantID), stringField("job_requisition_id", id.String()), stringField("status", status))
	}
	return mapJobRequisition(row), nil
}

func (s *Store) DeleteJobRequisition(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteJobRequisition(ctx, sqlc.SoftDeleteJobRequisitionParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete job requisition", fmt.Errorf("hrms: delete job requisition: %w", err), tenantIDField(tenantID), stringField("job_requisition_id", id.String()))
	}
	return nil
}

func (s *Store) CreateJobRequisitionLog(ctx context.Context, item *domain.JobRequisitionLog, actorID *uuid.UUID) (*domain.JobRequisitionLog, error) {
	row, err := s.getQueries(ctx).CreateJobRequisitionLog(ctx, sqlc.CreateJobRequisitionLogParams{TenantID: item.TenantID, JobRequisitionID: item.JobRequisitionID, FromStatus: textFromPtr(item.FromStatus), ToStatus: item.ToStatus, Action: item.Action, Remarks: textFromPtr(item.Remarks), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create job requisition log", fmt.Errorf("hrms: create job requisition log: %w", err), tenantIDField(item.TenantID), stringField("job_requisition_id", item.JobRequisitionID.String()), stringField("action", item.Action))
	}
	return mapJobRequisitionLog(row), nil
}

func (s *Store) ListJobRequisitionLogs(ctx context.Context, tenantID uuid.UUID, jobRequisitionID uuid.UUID) ([]*domain.JobRequisitionLog, error) {
	rows, err := s.getQueries(ctx).ListJobRequisitionLogs(ctx, sqlc.ListJobRequisitionLogsParams{TenantID: tenantID, JobRequisitionID: jobRequisitionID})
	if err != nil {
		return nil, s.logDBError(ctx, "list job requisition logs", err, tenantIDField(tenantID), stringField("job_requisition_id", jobRequisitionID.String()))
	}
	return mapJobRequisitionLogs(rows), nil
}
