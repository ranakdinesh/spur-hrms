package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDesignationLevelCodeID     = errors.New("designation_level_code_id is required")
	ErrInvalidDesignationSeniorityRankID = errors.New("designation_seniority_rank_id is required")
	ErrInvalidDesignationMasterLabel     = errors.New("designation master label is required")
	ErrInvalidDesignationMasterSortOrder = errors.New("designation master sort_order must be between 0 and 9999")
)

type DesignationLevelCode struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Code        string     `json:"code"`
	Label       string     `json:"label"`
	Description *string    `json:"description,omitempty"`
	SortOrder   int32      `json:"sort_order"`
	Inactive    bool       `json:"inactive"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
}

type DesignationSeniorityRank struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	RankValue   int32      `json:"rank_value"`
	Label       string     `json:"label"`
	Description *string    `json:"description,omitempty"`
	SortOrder   int32      `json:"sort_order"`
	Inactive    bool       `json:"inactive"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
}

type DesignationLevelCodeInput struct {
	TenantID    uuid.UUID
	Code        string
	Label       string
	Description *string
	SortOrder   int32
}

type DesignationSeniorityRankInput struct {
	TenantID    uuid.UUID
	RankValue   int32
	Label       string
	Description *string
	SortOrder   int32
}

func NewDesignationLevelCode(input DesignationLevelCodeInput) (*DesignationLevelCode, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code := normalizeDesignationLevelCode(input.Code)
	if !validDesignationLevelCode(code) {
		return nil, ErrInvalidDesignationLevelCode
	}
	label := strings.TrimSpace(input.Label)
	if label == "" {
		return nil, ErrInvalidDesignationMasterLabel
	}
	if input.SortOrder < 0 || input.SortOrder > 9999 {
		return nil, ErrInvalidDesignationMasterSortOrder
	}
	now := time.Now().UTC()
	return &DesignationLevelCode{
		TenantID:    input.TenantID,
		Code:        code,
		Label:       label,
		Description: cleanOptional(input.Description),
		SortOrder:   input.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func NewDesignationSeniorityRank(input DesignationSeniorityRankInput) (*DesignationSeniorityRank, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.RankValue < 1 || input.RankValue > 9999 {
		return nil, ErrInvalidDesignationSeniorityRank
	}
	label := strings.TrimSpace(input.Label)
	if label == "" {
		return nil, ErrInvalidDesignationMasterLabel
	}
	if input.SortOrder < 0 || input.SortOrder > 9999 {
		return nil, ErrInvalidDesignationMasterSortOrder
	}
	now := time.Now().UTC()
	return &DesignationSeniorityRank{
		TenantID:    input.TenantID,
		RankValue:   input.RankValue,
		Label:       label,
		Description: cleanOptional(input.Description),
		SortOrder:   input.SortOrder,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func DefaultDesignationLevelCodeInputs(tenantID uuid.UUID) []DesignationLevelCodeInput {
	return []DesignationLevelCodeInput{
		{TenantID: tenantID, Code: "L1", Label: "Entry / Trainee", Description: stringPtr("Entry-level or trainee role."), SortOrder: 10},
		{TenantID: tenantID, Code: "L2", Label: "Junior", Description: stringPtr("Early-career role with guided execution."), SortOrder: 20},
		{TenantID: tenantID, Code: "L3", Label: "Associate / Executive", Description: stringPtr("Independent contributor with defined responsibilities."), SortOrder: 30},
		{TenantID: tenantID, Code: "L4", Label: "Senior / Specialist", Description: stringPtr("Experienced contributor or first-line specialist."), SortOrder: 40},
		{TenantID: tenantID, Code: "L5", Label: "Lead / Manager", Description: stringPtr("Team lead or manager with ownership of outcomes."), SortOrder: 50},
		{TenantID: tenantID, Code: "L6", Label: "Senior Manager / Director", Description: stringPtr("Senior leadership for a function or department."), SortOrder: 60},
		{TenantID: tenantID, Code: "L7", Label: "Executive", Description: stringPtr("Executive leadership with organization-wide impact."), SortOrder: 70},
		{TenantID: tenantID, Code: "M1", Label: "Manager I", Description: stringPtr("First manager level."), SortOrder: 80},
		{TenantID: tenantID, Code: "M2", Label: "Manager II", Description: stringPtr("Experienced manager level."), SortOrder: 90},
		{TenantID: tenantID, Code: "M3", Label: "Senior Manager", Description: stringPtr("Senior manager level."), SortOrder: 100},
		{TenantID: tenantID, Code: "M4", Label: "Director", Description: stringPtr("Director level."), SortOrder: 110},
		{TenantID: tenantID, Code: "M5", Label: "Executive Manager", Description: stringPtr("Executive manager level."), SortOrder: 120},
	}
}

func DefaultDesignationSeniorityRankInputs(tenantID uuid.UUID) []DesignationSeniorityRankInput {
	return []DesignationSeniorityRankInput{
		{TenantID: tenantID, RankValue: 100, Label: "Entry / Trainee", Description: stringPtr("Lowest organization seniority."), SortOrder: 10},
		{TenantID: tenantID, RankValue: 200, Label: "Junior", Description: stringPtr("Junior role seniority."), SortOrder: 20},
		{TenantID: tenantID, RankValue: 300, Label: "Associate / Executive", Description: stringPtr("Standard independent contributor seniority."), SortOrder: 30},
		{TenantID: tenantID, RankValue: 400, Label: "Senior / Specialist", Description: stringPtr("Senior contributor seniority."), SortOrder: 40},
		{TenantID: tenantID, RankValue: 500, Label: "Lead / Manager", Description: stringPtr("Lead or manager seniority."), SortOrder: 50},
		{TenantID: tenantID, RankValue: 600, Label: "Senior Manager / Director", Description: stringPtr("Senior leadership seniority."), SortOrder: 60},
		{TenantID: tenantID, RankValue: 700, Label: "Executive", Description: stringPtr("Executive seniority."), SortOrder: 70},
	}
}

func stringPtr(value string) *string {
	return &value
}
