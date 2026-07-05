package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidJobPositionID        = errors.New("job_position_id is required")
	ErrInvalidJobPositionTitle     = errors.New("job position title is required")
	ErrInvalidJobPositionCode      = errors.New("job position code must use only letters, numbers, underscore, or hyphen")
	ErrInvalidJobPositionHeadcount = errors.New("job position total_position must be zero or greater")
	ErrInvalidJobPositionBudget    = errors.New("job position budgeted_cost must be zero or greater")
	ErrInvalidJobPositionWorkMode  = errors.New("job position work_mode is invalid")
	ErrInvalidJobLocationID        = errors.New("job_position_location_id is required")
	ErrInvalidJobLocation          = errors.New("job position location must include a city/location or be remote")
)

const (
	WorkModeOffice = "Office"
	WorkModeRemote = "Remote"
	WorkModeHybrid = "Hybrid"
)

type JobPosition struct {
	ID                   uuid.UUID  `json:"id"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	Code                 *string    `json:"code,omitempty"`
	Title                string     `json:"title"`
	Level                *string    `json:"level,omitempty"`
	Category             *string    `json:"category,omitempty"`
	Description          *string    `json:"description,omitempty"`
	DepartmentID         *uuid.UUID `json:"department_id,omitempty"`
	DepartmentName       *string    `json:"department_name,omitempty"`
	EmploymentTypeID     *uuid.UUID `json:"employment_type_id,omitempty"`
	EmploymentTypeName   *string    `json:"employment_type_name,omitempty"`
	WorkMode             *string    `json:"work_mode,omitempty"`
	TotalPosition        int32      `json:"total_position"`
	BudgetedCost         *float64   `json:"budgeted_cost,omitempty"`
	LocationCount        int32      `json:"location_count"`
	OpenRequisitionCount int32      `json:"open_requisition_count"`
	Inactive             bool       `json:"inactive"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type JobPositionInput struct {
	TenantID         uuid.UUID
	Code             *string
	Title            string
	Level            *string
	Category         *string
	Description      *string
	DepartmentID     *uuid.UUID
	EmploymentTypeID *uuid.UUID
	WorkMode         *string
	TotalPosition    int32
	BudgetedCost     *float64
}

func NewJobPosition(input JobPositionInput) (*JobPosition, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidJobPositionTitle
	}
	code := normalizeJobPositionCode(input.Code)
	if code != nil && !validJobPositionCode(*code) {
		return nil, ErrInvalidJobPositionCode
	}
	if input.TotalPosition < 0 {
		return nil, ErrInvalidJobPositionHeadcount
	}
	if input.BudgetedCost != nil && *input.BudgetedCost < 0 {
		return nil, ErrInvalidJobPositionBudget
	}
	workMode, err := ValidateJobWorkMode(input.WorkMode)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &JobPosition{
		TenantID:         input.TenantID,
		Code:             code,
		Title:            title,
		Level:            cleanOptional(input.Level),
		Category:         cleanOptional(input.Category),
		Description:      cleanOptional(input.Description),
		DepartmentID:     cleanUUIDOptional(input.DepartmentID),
		EmploymentTypeID: cleanUUIDOptional(input.EmploymentTypeID),
		WorkMode:         workMode,
		TotalPosition:    input.TotalPosition,
		BudgetedCost:     input.BudgetedCost,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

type JobPositionLocation struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	JobPositionID uuid.UUID  `json:"job_position_id"`
	Location      *string    `json:"location,omitempty"`
	City          *string    `json:"city,omitempty"`
	State         *string    `json:"state,omitempty"`
	Country       *string    `json:"country,omitempty"`
	IsRemote      bool       `json:"is_remote"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type JobPositionLocationInput struct {
	TenantID      uuid.UUID
	JobPositionID uuid.UUID
	Location      *string
	City          *string
	State         *string
	Country       *string
	IsRemote      bool
}

func NewJobPositionLocation(input JobPositionLocationInput) (*JobPositionLocation, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.JobPositionID == uuid.Nil {
		return nil, ErrInvalidJobPositionID
	}
	location := cleanOptional(input.Location)
	city := cleanOptional(input.City)
	if !input.IsRemote && location == nil && city == nil {
		return nil, ErrInvalidJobLocation
	}
	now := time.Now().UTC()
	return &JobPositionLocation{
		TenantID:      input.TenantID,
		JobPositionID: input.JobPositionID,
		Location:      location,
		City:          city,
		State:         cleanOptional(input.State),
		Country:       cleanOptional(input.Country),
		IsRemote:      input.IsRemote,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

type JobPositionFilter struct {
	TenantID         uuid.UUID
	DepartmentID     *uuid.UUID
	EmploymentTypeID *uuid.UUID
	WorkMode         *string
	Search           *string
	Limit            int32
	Offset           int32
}

type JobPositionPage struct {
	Items      []*JobPosition `json:"items"`
	Total      int64          `json:"total"`
	Limit      int32          `json:"limit"`
	Offset     int32          `json:"offset"`
	NextOffset *int32         `json:"next_offset,omitempty"`
}

func ValidateJobWorkMode(value *string) (*string, error) {
	clean := cleanOptional(value)
	if clean == nil {
		return nil, nil
	}
	normalized := strings.ToLower(*clean)
	switch normalized {
	case "office", "onsite", "on-site":
		mode := WorkModeOffice
		return &mode, nil
	case "remote":
		mode := WorkModeRemote
		return &mode, nil
	case "hybrid":
		mode := WorkModeHybrid
		return &mode, nil
	default:
		return nil, ErrInvalidJobPositionWorkMode
	}
}

func normalizeJobPositionCode(value *string) *string {
	clean := cleanOptional(value)
	if clean == nil {
		return nil
	}
	normalized := strings.ToUpper(*clean)
	return &normalized
}

func validJobPositionCode(value string) bool {
	if value == "" || len(value) > 50 {
		return false
	}
	for _, char := range value {
		if char >= 'A' && char <= 'Z' {
			continue
		}
		if char >= '0' && char <= '9' {
			continue
		}
		if char == '_' || char == '-' {
			continue
		}
		return false
	}
	return true
}
