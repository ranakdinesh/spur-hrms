package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDesignationID            = errors.New("designation_id is required")
	ErrInvalidDesignationName          = errors.New("designation name is required")
	ErrInvalidDesignationLevelCode     = errors.New("designation level_code is required and must use only letters, numbers, underscore, or hyphen")
	ErrInvalidDesignationSeniorityRank = errors.New("designation seniority_rank must be between 1 and 9999")
)

type Designation struct {
	ID                 uuid.UUID  `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	Name               string     `json:"name"`
	LevelCode          string     `json:"level_code"`
	SeniorityRank      int32      `json:"seniority_rank"`
	Description        *string    `json:"description,omitempty"`
	AttendanceRequired bool       `json:"attendance_required"`
	Inactive           bool       `json:"inactive"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatedBy          *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt          time.Time  `json:"updated_at"`
	UpdatedBy          *uuid.UUID `json:"updated_by,omitempty"`
}

type DesignationInput struct {
	TenantID           uuid.UUID
	Name               string
	LevelCode          string
	SeniorityRank      int32
	Description        *string
	AttendanceRequired *bool
}

func NewDesignation(input DesignationInput) (*Designation, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidDesignationName
	}
	levelCode := normalizeDesignationLevelCode(input.LevelCode)
	if !validDesignationLevelCode(levelCode) {
		return nil, ErrInvalidDesignationLevelCode
	}
	if input.SeniorityRank < 1 || input.SeniorityRank > 9999 {
		return nil, ErrInvalidDesignationSeniorityRank
	}
	attendanceRequired := true
	if input.AttendanceRequired != nil {
		attendanceRequired = *input.AttendanceRequired
	}
	now := time.Now().UTC()
	return &Designation{
		TenantID:           input.TenantID,
		Name:               name,
		LevelCode:          levelCode,
		SeniorityRank:      input.SeniorityRank,
		Description:        cleanOptional(input.Description),
		AttendanceRequired: attendanceRequired,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func normalizeDesignationLevelCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func validDesignationLevelCode(value string) bool {
	if value == "" || len(value) > 32 {
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
