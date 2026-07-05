package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmploymentTypeID   = errors.New("employment_type_id is required")
	ErrInvalidEmploymentTypeName = errors.New("employment type name is required")
	ErrInvalidMaritalStatusID    = errors.New("marital_status_id is required")
	ErrInvalidMaritalStatusName  = errors.New("marital status name is required")
)

type EmploymentType struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	Name      string     `json:"name"`
	Inactive  bool       `json:"inactive"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
}

type MaritalStatus struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	Name      string     `json:"name"`
	Inactive  bool       `json:"inactive"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
}

type EmploymentTypeInput struct {
	TenantID uuid.UUID
	Name     string
}

type MaritalStatusInput struct {
	TenantID uuid.UUID
	Name     string
}

func NewEmploymentType(input EmploymentTypeInput) (*EmploymentType, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidEmploymentTypeName
	}
	now := time.Now().UTC()
	return &EmploymentType{
		TenantID:  input.TenantID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func NewMaritalStatus(input MaritalStatusInput) (*MaritalStatus, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidMaritalStatusName
	}
	now := time.Now().UTC()
	return &MaritalStatus{
		TenantID:  input.TenantID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func DefaultEmploymentTypeInputs(tenantID uuid.UUID) []EmploymentTypeInput {
	return []EmploymentTypeInput{
		{TenantID: tenantID, Name: "Permanent"},
		{TenantID: tenantID, Name: "Probation"},
		{TenantID: tenantID, Name: "Contract"},
		{TenantID: tenantID, Name: "Consultant"},
		{TenantID: tenantID, Name: "Intern"},
		{TenantID: tenantID, Name: "Part-time"},
	}
}

func DefaultMaritalStatusInputs(tenantID uuid.UUID) []MaritalStatusInput {
	return []MaritalStatusInput{
		{TenantID: tenantID, Name: "Single"},
		{TenantID: tenantID, Name: "Married"},
		{TenantID: tenantID, Name: "Divorced"},
		{TenantID: tenantID, Name: "Widowed"},
		{TenantID: tenantID, Name: "Separated"},
	}
}
