package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPolicyTypeID   = errors.New("policy_type_id is required")
	ErrInvalidPolicyTypeName = errors.New("policy type name is required")
	ErrInvalidPolicyID       = errors.New("policy_id is required")
	ErrInvalidPolicyTitle    = errors.New("policy title is required")
	ErrPolicyStorageMissing  = errors.New("policy file storage is not configured")
	ErrInvalidPolicyFile     = errors.New("policy file content is invalid")
	ErrSystemPolicyReadOnly  = errors.New("system policy types are read-only")
	ErrPolicyTypeNotFound    = errors.New("policy type not found")
	ErrCompanyPolicyNotFound = errors.New("company policy not found")
)

type PolicyType struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  *uuid.UUID `json:"tenant_id,omitempty"`
	Name      string     `json:"name"`
	IsSystem  bool       `json:"is_system"`
	Inactive  bool       `json:"inactive"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
}

type CompanyPolicy struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	PolicyTypeID *uuid.UUID `json:"policy_type_id,omitempty"`
	Title        string     `json:"title"`
	FilePath     *string    `json:"file_path,omitempty"`
	Description  *string    `json:"description,omitempty"`
	Inactive     bool       `json:"inactive"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
}

type PolicyTypeInput struct {
	TenantID uuid.UUID
	Name     string
	IsSystem bool
}

type CompanyPolicyInput struct {
	TenantID     uuid.UUID
	PolicyTypeID *uuid.UUID
	Title        string
	FilePath     *string
	Description  *string
}

func NewPolicyType(input PolicyTypeInput) (*PolicyType, error) {
	if input.TenantID == uuid.Nil && !input.IsSystem {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidPolicyTypeName
	}
	now := time.Now().UTC()
	item := &PolicyType{Name: name, IsSystem: input.IsSystem, CreatedAt: now, UpdatedAt: now}
	if input.TenantID != uuid.Nil && !input.IsSystem {
		item.TenantID = &input.TenantID
	}
	return item, nil
}

func NewCompanyPolicy(input CompanyPolicyInput) (*CompanyPolicy, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidPolicyTitle
	}
	now := time.Now().UTC()
	return &CompanyPolicy{
		TenantID:     input.TenantID,
		PolicyTypeID: cleanUUIDOptional(input.PolicyTypeID),
		Title:        title,
		FilePath:     cleanStringPtr(input.FilePath),
		Description:  cleanStringPtr(input.Description),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func cleanStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}
