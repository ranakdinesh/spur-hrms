package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateJobPosting(ctx context.Context, cmd ports.JobPostingCommand) (*domain.JobPosting, error) {
	if err := s.applyJobPostingDefaults(ctx, &cmd); err != nil {
		return nil, err
	}
	item, err := domain.NewJobPosting(jobPostingInput(cmd))
	if err != nil {
		s.logError("validate job posting create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if err := s.validateJobPostingReferences(ctx, item); err != nil {
		return nil, err
	}
	result, err := s.jobPostings.CreateJobPosting(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create job posting", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListJobPostings(ctx context.Context, filter domain.JobPostingFilter) (*domain.JobPostingPage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate job posting list tenant", err)
		return nil, err
	}
	filter.Search = cleanStringPtr(filter.Search)
	filter.JobStatus = cleanStringPtr(filter.JobStatus)
	if filter.JobStatus != nil {
		if _, err := domain.ValidateJobPostingStatus(filter.JobStatus); err != nil {
			s.logError("validate job posting list status", err, serviceTenantIDField(filter.TenantID), serviceStringField("job_status", *filter.JobStatus))
			return nil, err
		}
	}
	limit, offset := normalizeListWindow(filter.Limit, filter.Offset)
	filter.Limit = limit
	filter.Offset = offset
	items, err := s.jobPostings.ListJobPostings(ctx, filter)
	if err != nil {
		s.logError("list job postings", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.jobPostings.CountJobPostings(ctx, filter)
	if err != nil {
		s.logError("count job postings", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	page := &domain.JobPostingPage{Items: items, Total: total, Limit: limit, Offset: offset}
	if int64(offset)+int64(len(items)) < total {
		next := offset + limit
		page.NextOffset = &next
	}
	return page, nil
}

func (s *TenantService) ListPublishedJobPostings(ctx context.Context, tenantID uuid.UUID) ([]*domain.JobPosting, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate published job postings tenant", err)
		return nil, err
	}
	return s.jobPostings.ListPublishedJobPostings(ctx, tenantID)
}

func (s *TenantService) GetJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobPosting, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate job posting get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidJobPostingID
		s.logError("validate job posting get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.jobPostings.GetJobPosting(ctx, tenantID, id)
	if err != nil {
		s.logError("get job posting", err, serviceTenantIDField(tenantID), serviceStringField("job_posting_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateJobPosting(ctx context.Context, cmd ports.JobPostingCommand) (*domain.JobPosting, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidJobPostingID
		s.logError("validate job posting update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetJobPosting(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if cmd.JobRequisitionID == nil {
		cmd.JobRequisitionID = existing.JobRequisitionID
	}
	if err := s.applyJobPostingDefaults(ctx, &cmd); err != nil {
		return nil, err
	}
	item, err := domain.NewJobPosting(jobPostingInput(cmd))
	if err != nil {
		s.logError("validate job posting update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_posting_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	if err := s.validateJobPostingReferences(ctx, item); err != nil {
		return nil, err
	}
	result, err := s.jobPostings.UpdateJobPosting(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update job posting", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_posting_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) PublishJobPosting(ctx context.Context, cmd ports.JobPostingPublishCommand) (*domain.JobPosting, error) {
	item, err := s.GetJobPosting(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if item.JobRequisitionID != nil {
		req, err := s.GetJobRequisition(ctx, cmd.TenantID, *item.JobRequisitionID)
		if err != nil {
			return nil, err
		}
		if req.Status != domain.ReqStatusApproved {
			err := domain.ErrInvalidJobPostingPublish
			s.logError("validate job posting publish requisition", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_requisition_id", req.ID.String()), serviceStringField("status", req.Status))
			return nil, err
		}
	}
	result, err := s.jobPostings.PublishJobPosting(ctx, cmd.TenantID, cmd.ID, cmd.ExpiryDate, cmd.ActorID)
	if err != nil {
		s.logError("publish job posting", err, serviceTenantIDField(cmd.TenantID), serviceStringField("job_posting_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ExpireJobPostings(ctx context.Context, tenantID uuid.UUID) ([]*domain.JobPosting, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate expire job postings tenant", err)
		return nil, err
	}
	return s.jobPostings.ExpireJobPostings(ctx, tenantID)
}

func (s *TenantService) CloseJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.JobPosting, error) {
	if _, err := s.GetJobPosting(ctx, tenantID, id); err != nil {
		return nil, err
	}
	result, err := s.jobPostings.CloseJobPosting(ctx, tenantID, id, actorID)
	if err != nil {
		s.logError("close job posting", err, serviceTenantIDField(tenantID), serviceStringField("job_posting_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetJobPosting(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.jobPostings.DeleteJobPosting(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete job posting", err, serviceTenantIDField(tenantID), serviceStringField("job_posting_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) applyJobPostingDefaults(ctx context.Context, cmd *ports.JobPostingCommand) error {
	if cmd.JobRequisitionID == nil || *cmd.JobRequisitionID == uuid.Nil {
		return nil
	}
	req, err := s.GetJobRequisition(ctx, cmd.TenantID, *cmd.JobRequisitionID)
	if err != nil {
		return err
	}
	if req.Status != domain.ReqStatusApproved && cmd.IsPublished {
		return domain.ErrInvalidJobPostingPublish
	}
	if cmd.Title == nil {
		cmd.Title = &req.Title
	}
	if cmd.Description == nil {
		cmd.Description = req.Description
	}
	if cmd.JobCategory == nil {
		cmd.JobCategory = req.Category
	}
	if cmd.DepartmentID == nil {
		cmd.DepartmentID = req.DepartmentID
	}
	if cmd.EmploymentTypeID == nil {
		cmd.EmploymentTypeID = req.EmploymentTypeID
	}
	if cmd.WorkMode == nil {
		cmd.WorkMode = req.WorkMode
	}
	if cmd.MinSalary == nil {
		cmd.MinSalary = req.MinSalary
	}
	if cmd.MaxSalary == nil {
		cmd.MaxSalary = req.MaxSalary
	}
	if cmd.SalaryCurrency == nil && req.Currency != "" {
		cmd.SalaryCurrency = &req.Currency
	}
	if cmd.JobStatus == nil {
		status := domain.JobPostingStatusDraft
		cmd.JobStatus = &status
	}
	return nil
}

func (s *TenantService) validateJobPostingReferences(ctx context.Context, item *domain.JobPosting) error {
	if item.JobRequisitionID != nil {
		if _, err := s.GetJobRequisition(ctx, item.TenantID, *item.JobRequisitionID); err != nil {
			return err
		}
	}
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

func jobPostingInput(cmd ports.JobPostingCommand) domain.JobPostingInput {
	return domain.JobPostingInput{TenantID: cmd.TenantID, JobRequisitionID: cmd.JobRequisitionID, Code: cmd.Code, Title: cmd.Title, JobSummary: cmd.JobSummary, Description: cmd.Description, JobCategory: cmd.JobCategory, DepartmentID: cmd.DepartmentID, Industry: cmd.Industry, EmploymentTypeID: cmd.EmploymentTypeID, WorkMode: cmd.WorkMode, RoleType: cmd.RoleType, MinExperience: cmd.MinExperience, MaxExperience: cmd.MaxExperience, MinSalary: cmd.MinSalary, MaxSalary: cmd.MaxSalary, SalaryCurrency: cmd.SalaryCurrency, SalaryPeriod: cmd.SalaryPeriod, IsSalaryVisible: cmd.IsSalaryVisible, JobStatus: cmd.JobStatus, PublishDate: cmd.PublishDate, ExpiryDate: cmd.ExpiryDate, IsPublished: cmd.IsPublished}
}
