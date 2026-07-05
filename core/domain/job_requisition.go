package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidJobRequisitionID       = errors.New("job_requisition_id is required")
	ErrInvalidJobRequisitionTitle    = errors.New("job requisition title is required")
	ErrInvalidJobRequisitionCode     = errors.New("job requisition code must use only letters, numbers, underscore, or hyphen")
	ErrInvalidJobRequisitionOpenings = errors.New("job requisition total_openings must be greater than zero")
	ErrInvalidJobRequisitionUser     = errors.New("job requisition requested_by is required")
	ErrInvalidJobRequisitionSalary   = errors.New("job requisition salary range is invalid")
	ErrInvalidJobRequisitionStatus   = errors.New("job requisition status is invalid")
	ErrInvalidJobRequisitionAction   = errors.New("job requisition action is invalid")
)

type JobRequisition struct {
	ID                     uuid.UUID            `json:"id"`
	TenantID               uuid.UUID            `json:"tenant_id"`
	JobPositionID          uuid.UUID            `json:"job_position_id"`
	JobPositionCode        *string              `json:"job_position_code,omitempty"`
	PositionTotalHeadcount int32                `json:"position_total_headcount"`
	PositionBudgetedCost   *float64             `json:"position_budgeted_cost,omitempty"`
	Code                   *string              `json:"code,omitempty"`
	Title                  string               `json:"title"`
	Level                  *string              `json:"level,omitempty"`
	Category               *string              `json:"category,omitempty"`
	DepartmentID           *uuid.UUID           `json:"department_id,omitempty"`
	DepartmentName         *string              `json:"department_name,omitempty"`
	EmploymentTypeID       *uuid.UUID           `json:"employment_type_id,omitempty"`
	EmploymentTypeName     *string              `json:"employment_type_name,omitempty"`
	Description            *string              `json:"description,omitempty"`
	WorkMode               *string              `json:"work_mode,omitempty"`
	TotalOpenings          int32                `json:"total_openings"`
	ReasonForHire          *string              `json:"reason_for_hire,omitempty"`
	MinSalary              *float64             `json:"min_salary,omitempty"`
	MaxSalary              *float64             `json:"max_salary,omitempty"`
	Currency               string               `json:"currency"`
	TargetHireDate         *time.Time           `json:"target_hire_date,omitempty"`
	ExpectedClosureDate    *time.Time           `json:"expected_closure_date,omitempty"`
	RequestedBy            uuid.UUID            `json:"requested_by"`
	RequestedDate          *time.Time           `json:"requested_date,omitempty"`
	IsApproved             bool                 `json:"is_approved"`
	ApprovedBy             *uuid.UUID           `json:"approved_by,omitempty"`
	ApprovedDate           *time.Time           `json:"approved_date,omitempty"`
	Priority               *string              `json:"priority,omitempty"`
	Status                 string               `json:"status"`
	Notes                  *string              `json:"notes,omitempty"`
	LogCount               int32                `json:"log_count"`
	Logs                   []*JobRequisitionLog `json:"logs,omitempty"`
	Inactive               bool                 `json:"inactive"`
	CreatedAt              time.Time            `json:"created_at"`
	CreatedBy              *uuid.UUID           `json:"created_by,omitempty"`
	UpdatedAt              time.Time            `json:"updated_at"`
	UpdatedBy              *uuid.UUID           `json:"updated_by,omitempty"`
}

type JobRequisitionInput struct {
	TenantID            uuid.UUID
	JobPositionID       uuid.UUID
	Code                *string
	Title               string
	Level               *string
	Category            *string
	DepartmentID        *uuid.UUID
	EmploymentTypeID    *uuid.UUID
	Description         *string
	WorkMode            *string
	TotalOpenings       int32
	ReasonForHire       *string
	MinSalary           *float64
	MaxSalary           *float64
	Currency            string
	TargetHireDate      *time.Time
	ExpectedClosureDate *time.Time
	RequestedBy         uuid.UUID
	RequestedDate       *time.Time
	Priority            *string
	Status              string
	Notes               *string
}

func NewJobRequisition(input JobRequisitionInput) (*JobRequisition, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.JobPositionID == uuid.Nil {
		return nil, ErrInvalidJobPositionID
	}
	if input.RequestedBy == uuid.Nil {
		return nil, ErrInvalidJobRequisitionUser
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidJobRequisitionTitle
	}
	code := normalizeJobPositionCode(input.Code)
	if code != nil && !validJobPositionCode(*code) {
		return nil, ErrInvalidJobRequisitionCode
	}
	if input.TotalOpenings <= 0 {
		return nil, ErrInvalidJobRequisitionOpenings
	}
	if input.MinSalary != nil && *input.MinSalary < 0 {
		return nil, ErrInvalidJobRequisitionSalary
	}
	if input.MaxSalary != nil && *input.MaxSalary < 0 {
		return nil, ErrInvalidJobRequisitionSalary
	}
	if input.MinSalary != nil && input.MaxSalary != nil && *input.MinSalary > *input.MaxSalary {
		return nil, ErrInvalidJobRequisitionSalary
	}
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = ReqStatusDraft
	}
	var err error
	status, err = ValidateRequisitionStatus(status)
	if err != nil {
		return nil, ErrInvalidJobRequisitionStatus
	}
	workMode, err := ValidateJobWorkMode(input.WorkMode)
	if err != nil {
		return nil, err
	}
	currency := strings.ToUpper(strings.TrimSpace(input.Currency))
	if currency == "" {
		currency = "INR"
	}
	now := time.Now().UTC()
	requestedDate := input.RequestedDate
	if requestedDate == nil || requestedDate.IsZero() {
		today := dateOnly(now)
		requestedDate = &today
	}
	return &JobRequisition{
		TenantID:            input.TenantID,
		JobPositionID:       input.JobPositionID,
		Code:                code,
		Title:               title,
		Level:               cleanOptional(input.Level),
		Category:            cleanOptional(input.Category),
		DepartmentID:        cleanUUIDOptional(input.DepartmentID),
		EmploymentTypeID:    cleanUUIDOptional(input.EmploymentTypeID),
		Description:         cleanOptional(input.Description),
		WorkMode:            workMode,
		TotalOpenings:       input.TotalOpenings,
		ReasonForHire:       cleanOptional(input.ReasonForHire),
		MinSalary:           input.MinSalary,
		MaxSalary:           input.MaxSalary,
		Currency:            currency,
		TargetHireDate:      cleanTimeOptional(input.TargetHireDate),
		ExpectedClosureDate: cleanTimeOptional(input.ExpectedClosureDate),
		RequestedBy:         input.RequestedBy,
		RequestedDate:       cleanTimeOptional(requestedDate),
		Priority:            cleanOptional(input.Priority),
		Status:              status,
		Notes:               cleanOptional(input.Notes),
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

type JobRequisitionLog struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	JobRequisitionID uuid.UUID  `json:"job_requisition_id"`
	FromStatus       *string    `json:"from_status,omitempty"`
	ToStatus         string     `json:"to_status"`
	Action           string     `json:"action"`
	Remarks          *string    `json:"remarks,omitempty"`
	Inactive         bool       `json:"inactive"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
}

type JobRequisitionFilter struct {
	TenantID      uuid.UUID
	Status        *string
	JobPositionID *uuid.UUID
	DepartmentID  *uuid.UUID
	Search        *string
	Limit         int32
	Offset        int32
}

type JobRequisitionPage struct {
	Items      []*JobRequisition `json:"items"`
	Total      int64             `json:"total"`
	Limit      int32             `json:"limit"`
	Offset     int32             `json:"offset"`
	NextOffset *int32            `json:"next_offset,omitempty"`
}

func cleanTimeOptional(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	clean := dateOnly(*value)
	return &clean
}
