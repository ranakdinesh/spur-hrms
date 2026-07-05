package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDepartmentID        = errors.New("department_id is required")
	ErrInvalidDepartmentName      = errors.New("department name is required")
	ErrInvalidDepartmentShortCode = errors.New("department short_code is required and must use only letters, numbers, underscore, or hyphen")
)

type Department struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Name        string     `json:"name"`
	ShortCode   string     `json:"short_code"`
	Description *string    `json:"description,omitempty"`
	Inactive    bool       `json:"inactive"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
}

type DepartmentInput struct {
	TenantID    uuid.UUID
	Name        string
	ShortCode   string
	Description *string
}

func NewDepartment(input DepartmentInput) (*Department, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidDepartmentName
	}
	shortCode := normalizeDepartmentShortCode(input.ShortCode)
	if !validDepartmentShortCode(shortCode) {
		return nil, ErrInvalidDepartmentShortCode
	}
	now := time.Now().UTC()
	return &Department{
		TenantID:    input.TenantID,
		Name:        name,
		ShortCode:   shortCode,
		Description: cleanOptional(input.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func normalizeDepartmentShortCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func validDepartmentShortCode(value string) bool {
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
