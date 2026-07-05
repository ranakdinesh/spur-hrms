package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidJobPostingID         = errors.New("job_posting_id is required")
	ErrInvalidJobPostingTitle      = errors.New("job posting title is required")
	ErrInvalidJobPostingCode       = errors.New("job posting code must use only letters, numbers, underscore, or hyphen")
	ErrInvalidJobPostingSalary     = errors.New("job posting salary range is invalid")
	ErrInvalidJobPostingExperience = errors.New("job posting experience range is invalid")
	ErrInvalidJobPostingStatus     = errors.New("job posting status is invalid")
	ErrInvalidJobPostingPublish    = errors.New("job posting can only be published from an approved requisition")
)

const (
	JobPostingStatusDraft   = "Draft"
	JobPostingStatusOpen    = "Open"
	JobPostingStatusClosed  = "Closed"
	JobPostingStatusExpired = "Expired"
)

type JobPosting struct {
	ID                   uuid.UUID  `json:"id"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	JobRequisitionID     *uuid.UUID `json:"job_requisition_id,omitempty"`
	JobRequisitionCode   *string    `json:"job_requisition_code,omitempty"`
	JobRequisitionStatus *string    `json:"job_requisition_status,omitempty"`
	Code                 *string    `json:"code,omitempty"`
	Title                *string    `json:"title,omitempty"`
	JobSummary           *string    `json:"job_summary,omitempty"`
	Description          *string    `json:"description,omitempty"`
	JobCategory          *string    `json:"job_category,omitempty"`
	DepartmentID         *uuid.UUID `json:"department_id,omitempty"`
	DepartmentName       *string    `json:"department_name,omitempty"`
	Industry             *string    `json:"industry,omitempty"`
	EmploymentTypeID     *uuid.UUID `json:"employment_type_id,omitempty"`
	EmploymentTypeName   *string    `json:"employment_type_name,omitempty"`
	WorkMode             *string    `json:"work_mode,omitempty"`
	RoleType             *string    `json:"role_type,omitempty"`
	MinExperience        *float64   `json:"min_experience,omitempty"`
	MaxExperience        *float64   `json:"max_experience,omitempty"`
	MinSalary            *float64   `json:"min_salary,omitempty"`
	MaxSalary            *float64   `json:"max_salary,omitempty"`
	SalaryCurrency       *string    `json:"salary_currency,omitempty"`
	SalaryPeriod         *string    `json:"salary_period,omitempty"`
	IsSalaryVisible      bool       `json:"is_salary_visible"`
	EffectiveStatus      *string    `json:"effective_status,omitempty"`
	JobStatus            *string    `json:"job_status,omitempty"`
	PublishDate          *time.Time `json:"publish_date,omitempty"`
	ExpiryDate           *time.Time `json:"expiry_date,omitempty"`
	IsPublished          bool       `json:"is_published"`
	Inactive             bool       `json:"inactive"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type JobPostingInput struct {
	TenantID         uuid.UUID
	JobRequisitionID *uuid.UUID
	Code             *string
	Title            *string
	JobSummary       *string
	Description      *string
	JobCategory      *string
	DepartmentID     *uuid.UUID
	Industry         *string
	EmploymentTypeID *uuid.UUID
	WorkMode         *string
	RoleType         *string
	MinExperience    *float64
	MaxExperience    *float64
	MinSalary        *float64
	MaxSalary        *float64
	SalaryCurrency   *string
	SalaryPeriod     *string
	IsSalaryVisible  bool
	JobStatus        *string
	PublishDate      *time.Time
	ExpiryDate       *time.Time
	IsPublished      bool
}

func NewJobPosting(input JobPostingInput) (*JobPosting, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	title := cleanOptional(input.Title)
	if title == nil || strings.TrimSpace(*title) == "" {
		return nil, ErrInvalidJobPostingTitle
	}
	code := normalizeJobPositionCode(input.Code)
	if code != nil && !validJobPositionCode(*code) {
		return nil, ErrInvalidJobPostingCode
	}
	if input.MinSalary != nil && *input.MinSalary < 0 {
		return nil, ErrInvalidJobPostingSalary
	}
	if input.MaxSalary != nil && *input.MaxSalary < 0 {
		return nil, ErrInvalidJobPostingSalary
	}
	if input.MinSalary != nil && input.MaxSalary != nil && *input.MinSalary > *input.MaxSalary {
		return nil, ErrInvalidJobPostingSalary
	}
	if input.MinExperience != nil && *input.MinExperience < 0 {
		return nil, ErrInvalidJobPostingExperience
	}
	if input.MaxExperience != nil && *input.MaxExperience < 0 {
		return nil, ErrInvalidJobPostingExperience
	}
	if input.MinExperience != nil && input.MaxExperience != nil && *input.MinExperience > *input.MaxExperience {
		return nil, ErrInvalidJobPostingExperience
	}
	status, err := ValidateJobPostingStatus(input.JobStatus)
	if err != nil {
		return nil, err
	}
	if status == nil {
		value := JobPostingStatusDraft
		status = &value
	}
	workMode, err := ValidateJobWorkMode(input.WorkMode)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &JobPosting{
		TenantID:         input.TenantID,
		JobRequisitionID: cleanUUIDOptional(input.JobRequisitionID),
		Code:             code,
		Title:            title,
		JobSummary:       cleanOptional(input.JobSummary),
		Description:      cleanOptional(input.Description),
		JobCategory:      cleanOptional(input.JobCategory),
		DepartmentID:     cleanUUIDOptional(input.DepartmentID),
		Industry:         cleanOptional(input.Industry),
		EmploymentTypeID: cleanUUIDOptional(input.EmploymentTypeID),
		WorkMode:         workMode,
		RoleType:         cleanOptional(input.RoleType),
		MinExperience:    input.MinExperience,
		MaxExperience:    input.MaxExperience,
		MinSalary:        input.MinSalary,
		MaxSalary:        input.MaxSalary,
		SalaryCurrency:   cleanOptional(input.SalaryCurrency),
		SalaryPeriod:     cleanOptional(input.SalaryPeriod),
		IsSalaryVisible:  input.IsSalaryVisible,
		JobStatus:        status,
		PublishDate:      cleanTimeOptional(input.PublishDate),
		ExpiryDate:       cleanTimeOptional(input.ExpiryDate),
		IsPublished:      input.IsPublished,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

type JobPostingFilter struct {
	TenantID     uuid.UUID
	JobStatus    *string
	IsPublished  *bool
	DepartmentID *uuid.UUID
	Search       *string
	Limit        int32
	Offset       int32
}

type JobPostingPage struct {
	Items      []*JobPosting `json:"items"`
	Total      int64         `json:"total"`
	Limit      int32         `json:"limit"`
	Offset     int32         `json:"offset"`
	NextOffset *int32        `json:"next_offset,omitempty"`
}

func ValidateJobPostingStatus(value *string) (*string, error) {
	status := cleanOptional(value)
	if status == nil {
		return nil, nil
	}
	switch *status {
	case JobPostingStatusDraft, JobPostingStatusOpen, JobPostingStatusClosed, JobPostingStatusExpired:
		return status, nil
	default:
		return nil, ErrInvalidJobPostingStatus
	}
}
