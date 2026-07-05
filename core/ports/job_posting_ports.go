package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type JobPostingRepo interface {
	CreateJobPosting(ctx context.Context, item *domain.JobPosting, actorID *uuid.UUID) (*domain.JobPosting, error)
	ListJobPostings(ctx context.Context, filter domain.JobPostingFilter) ([]*domain.JobPosting, error)
	CountJobPostings(ctx context.Context, filter domain.JobPostingFilter) (int64, error)
	ListPublishedJobPostings(ctx context.Context, tenantID uuid.UUID) ([]*domain.JobPosting, error)
	GetJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.JobPosting, error)
	GetJobPostingByRequisition(ctx context.Context, tenantID uuid.UUID, requisitionID uuid.UUID) (*domain.JobPosting, error)
	UpdateJobPosting(ctx context.Context, item *domain.JobPosting, actorID *uuid.UUID) (*domain.JobPosting, error)
	PublishJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, expiryDate *time.Time, actorID *uuid.UUID) (*domain.JobPosting, error)
	ExpireJobPostings(ctx context.Context, tenantID uuid.UUID) ([]*domain.JobPosting, error)
	CloseJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.JobPosting, error)
	DeleteJobPosting(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type JobPostingCommand struct {
	ID               uuid.UUID  `json:"id,omitempty"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	JobRequisitionID *uuid.UUID `json:"job_requisition_id,omitempty"`
	Code             *string    `json:"code,omitempty"`
	Title            *string    `json:"title,omitempty"`
	JobSummary       *string    `json:"job_summary,omitempty"`
	Description      *string    `json:"description,omitempty"`
	JobCategory      *string    `json:"job_category,omitempty"`
	DepartmentID     *uuid.UUID `json:"department_id,omitempty"`
	Industry         *string    `json:"industry,omitempty"`
	EmploymentTypeID *uuid.UUID `json:"employment_type_id,omitempty"`
	WorkMode         *string    `json:"work_mode,omitempty"`
	RoleType         *string    `json:"role_type,omitempty"`
	MinExperience    *float64   `json:"min_experience,omitempty"`
	MaxExperience    *float64   `json:"max_experience,omitempty"`
	MinSalary        *float64   `json:"min_salary,omitempty"`
	MaxSalary        *float64   `json:"max_salary,omitempty"`
	SalaryCurrency   *string    `json:"salary_currency,omitempty"`
	SalaryPeriod     *string    `json:"salary_period,omitempty"`
	IsSalaryVisible  bool       `json:"is_salary_visible"`
	JobStatus        *string    `json:"job_status,omitempty"`
	PublishDate      *time.Time `json:"publish_date,omitempty"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	IsPublished      bool       `json:"is_published"`
	ActorID          *uuid.UUID `json:"-"`
}

type JobPostingPublishCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	ID         uuid.UUID  `json:"id"`
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`
	ActorID    *uuid.UUID `json:"-"`
}
