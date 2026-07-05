package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidBranchID   = errors.New("branch_id is required")
	ErrInvalidBranchName = errors.New("branch name is required")
)

type Branch struct {
	ID                  uuid.UUID  `json:"id"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	Name                string     `json:"name"`
	Address             *string    `json:"address,omitempty"`
	City                *string    `json:"city,omitempty"`
	State               *string    `json:"state,omitempty"`
	Country             *string    `json:"country,omitempty"`
	Pincode             *string    `json:"pincode,omitempty"`
	Phone               *string    `json:"phone,omitempty"`
	BranchManagerUserID *uuid.UUID `json:"branch_manager_user_id,omitempty"`
	HRUserID            *uuid.UUID `json:"hr_user_id,omitempty"`
	AccountsUserID      *uuid.UUID `json:"accounts_user_id,omitempty"`
	Inactive            bool       `json:"inactive"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt           time.Time  `json:"updated_at"`
	UpdatedBy           *uuid.UUID `json:"updated_by,omitempty"`
}

type BranchInput struct {
	TenantID            uuid.UUID
	Name                string
	Address             *string
	City                *string
	State               *string
	Country             *string
	Pincode             *string
	Phone               *string
	BranchManagerUserID *uuid.UUID
	HRUserID            *uuid.UUID
	AccountsUserID      *uuid.UUID
}

func NewBranch(input BranchInput) (*Branch, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidBranchName
	}
	now := time.Now().UTC()
	return &Branch{
		TenantID:            input.TenantID,
		Name:                name,
		Address:             cleanOptional(input.Address),
		City:                cleanOptional(input.City),
		State:               cleanOptional(input.State),
		Country:             cleanOptional(input.Country),
		Pincode:             cleanOptional(input.Pincode),
		Phone:               cleanOptional(input.Phone),
		BranchManagerUserID: cleanUUIDOptional(input.BranchManagerUserID),
		HRUserID:            cleanUUIDOptional(input.HRUserID),
		AccountsUserID:      cleanUUIDOptional(input.AccountsUserID),
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}
