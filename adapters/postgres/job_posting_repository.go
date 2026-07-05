package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateJobPosting(ctx context.Context, item *domain.JobPosting, actorID *uuid.UUID) (*domain.JobPosting, error) {
	row, err := s.getQueries(ctx).CreateJobPosting(ctx, jobPostingCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create job posting", fmt.Errorf("hrms: create job posting: %w", err), tenantIDField(item.TenantID))
	}
	return mapJobPosting(row), nil
}

func (s *Store) ListJobPostings(ctx context.Context, filter domain.JobPostingFilter) ([]*domain.JobPosting, error) {
	rows, err := s.getQueries(ctx).ListJobPostings(ctx, sqlc.ListJobPostingsParams{TenantID: filter.TenantID, JobStatus: textFromPtr(filter.JobStatus), IsPublished: boolFromPtr(filter.IsPublished), DepartmentID: uuidFromPtr(filter.DepartmentID), Search: textFromPtr(filter.Search), Offset: filter.Offset, Limit: filter.Limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list job postings", err, tenantIDField(filter.TenantID))
	}
	return mapJobPostings(rows), nil
}

func (s *Store) CountJobPostings(ctx context.Context, filter domain.JobPostingFilter) (int64, error) {
	count, err := s.getQueries(ctx).CountJobPostings(ctx, sqlc.CountJobPostingsParams{TenantID: filter.TenantID, JobStatus: textFromPtr(filter.JobStatus), IsPublished: boolFromPtr(filter.IsPublished), DepartmentID: uuidFromPtr(filter.DepartmentID), Search: textFromPtr(filter.Search)})
	if err != nil {
		return 0, s.logDBError(ctx, "count job postings", err, tenantIDField(filter.TenantID))
	}
	return count, nil
}

func (s *Store) ListPublishedJobPostings(ctx context.Context, tenantID uuid.UUID) ([]*domain.JobPosting, error) {
	rows, err := s.getQueries(ctx).ListPublishedJobPostings(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list published job postings", err, tenantIDField(tenantID))
	}
	return mapPublishedJobPostings(rows), nil
}

func (s *Store) GetJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobPosting, error) {
	row, err := s.getQueries(ctx).GetJobPosting(ctx, sqlc.GetJobPostingParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get job posting", fmt.Errorf("hrms: get job posting: %w", err), tenantIDField(tenantID), stringField("job_posting_id", id.String()))
	}
	return mapJobPosting(row), nil
}

func (s *Store) GetJobPostingByRequisition(ctx context.Context, tenantID uuid.UUID, requisitionID uuid.UUID) (*domain.JobPosting, error) {
	row, err := s.getQueries(ctx).GetJobPostingByRequisition(ctx, sqlc.GetJobPostingByRequisitionParams{TenantID: tenantID, JobRequisitionID: uuidFromPtr(&requisitionID)})
	if err != nil {
		return nil, s.logDBError(ctx, "get job posting by requisition", fmt.Errorf("hrms: get job posting by requisition: %w", err), tenantIDField(tenantID), stringField("job_requisition_id", requisitionID.String()))
	}
	return mapJobPosting(row), nil
}

func (s *Store) UpdateJobPosting(ctx context.Context, item *domain.JobPosting, actorID *uuid.UUID) (*domain.JobPosting, error) {
	params := jobPostingUpdateParams(item, actorID)
	row, err := s.getQueries(ctx).UpdateJobPosting(ctx, params)
	if err != nil {
		return nil, s.logDBError(ctx, "update job posting", fmt.Errorf("hrms: update job posting: %w", err), tenantIDField(item.TenantID), stringField("job_posting_id", item.ID.String()))
	}
	return mapJobPosting(row), nil
}

func (s *Store) PublishJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, expiryDate *time.Time, actorID *uuid.UUID) (*domain.JobPosting, error) {
	row, err := s.getQueries(ctx).PublishJobPosting(ctx, sqlc.PublishJobPostingParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID), ExpiryDate: dateFromPtr(expiryDate)})
	if err != nil {
		return nil, s.logDBError(ctx, "publish job posting", fmt.Errorf("hrms: publish job posting: %w", err), tenantIDField(tenantID), stringField("job_posting_id", id.String()))
	}
	return mapJobPosting(row), nil
}

func (s *Store) ExpireJobPostings(ctx context.Context, tenantID uuid.UUID) ([]*domain.JobPosting, error) {
	rows, err := s.getQueries(ctx).ExpireJobPostings(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "expire job postings", err, tenantIDField(tenantID))
	}
	return mapPublishedJobPostings(rows), nil
}

func (s *Store) CloseJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.JobPosting, error) {
	row, err := s.getQueries(ctx).CloseJobPosting(ctx, sqlc.CloseJobPostingParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "close job posting", fmt.Errorf("hrms: close job posting: %w", err), tenantIDField(tenantID), stringField("job_posting_id", id.String()))
	}
	return mapJobPosting(row), nil
}

func (s *Store) DeleteJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteJobPosting(ctx, sqlc.SoftDeleteJobPostingParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete job posting", fmt.Errorf("hrms: delete job posting: %w", err), tenantIDField(tenantID), stringField("job_posting_id", id.String()))
	}
	return nil
}

func jobPostingCreateParams(item *domain.JobPosting, actorID *uuid.UUID) sqlc.CreateJobPostingParams {
	return sqlc.CreateJobPostingParams{TenantID: item.TenantID, JobRequisitionID: uuidFromPtr(item.JobRequisitionID), Code: textFromPtr(item.Code), Title: textFromPtr(item.Title), JobSummary: textFromPtr(item.JobSummary), Description: textFromPtr(item.Description), JobCategory: textFromPtr(item.JobCategory), DepartmentID: uuidFromPtr(item.DepartmentID), Industry: textFromPtr(item.Industry), EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID), WorkMode: textFromPtr(item.WorkMode), RoleType: textFromPtr(item.RoleType), MinExperience: numericFromFloatPtr(item.MinExperience), MaxExperience: numericFromFloatPtr(item.MaxExperience), MinSalary: numericFromFloatPtr(item.MinSalary), MaxSalary: numericFromFloatPtr(item.MaxSalary), SalaryCurrency: textFromPtr(item.SalaryCurrency), SalaryPeriod: textFromPtr(item.SalaryPeriod), IsSalaryVisible: item.IsSalaryVisible, JobStatus: textFromPtr(item.JobStatus), PublishDate: dateFromPtr(item.PublishDate), ExpiryDate: dateFromPtr(item.ExpiryDate), IsPublished: item.IsPublished, CreatedBy: uuidFromPtr(actorID)}
}

func jobPostingUpdateParams(item *domain.JobPosting, actorID *uuid.UUID) sqlc.UpdateJobPostingParams {
	return sqlc.UpdateJobPostingParams{TenantID: item.TenantID, ID: item.ID, Code: textFromPtr(item.Code), Title: textFromPtr(item.Title), JobSummary: textFromPtr(item.JobSummary), Description: textFromPtr(item.Description), JobCategory: textFromPtr(item.JobCategory), DepartmentID: uuidFromPtr(item.DepartmentID), Industry: textFromPtr(item.Industry), EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID), WorkMode: textFromPtr(item.WorkMode), RoleType: textFromPtr(item.RoleType), MinExperience: numericFromFloatPtr(item.MinExperience), MaxExperience: numericFromFloatPtr(item.MaxExperience), MinSalary: numericFromFloatPtr(item.MinSalary), MaxSalary: numericFromFloatPtr(item.MaxSalary), SalaryCurrency: textFromPtr(item.SalaryCurrency), SalaryPeriod: textFromPtr(item.SalaryPeriod), IsSalaryVisible: item.IsSalaryVisible, JobStatus: textFromPtr(item.JobStatus), PublishDate: dateFromPtr(item.PublishDate), ExpiryDate: dateFromPtr(item.ExpiryDate), IsPublished: item.IsPublished, UpdatedBy: uuidFromPtr(actorID)}
}
