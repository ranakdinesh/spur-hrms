package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ProjectSkillRequirementImportanceNiceToHave = "nice_to_have"
	ProjectSkillRequirementImportanceRequired   = "required"
	ProjectSkillRequirementImportanceCritical   = "critical"

	ProjectSkillRequirementSourceProject    = "project"
	ProjectSkillRequirementSourceEngagement = "engagement"
	ProjectSkillRequirementSourceRole       = "role"
	ProjectSkillRequirementSourceClient     = "client"
	ProjectSkillRequirementSourceCompliance = "compliance"

	SkillGapActionCovered        = "covered"
	SkillGapActionTrainOrAssign  = "train_or_assign"
	SkillGapActionHireOrContract = "hire_or_contract"
	SkillGapActionTrain          = "train"
)

var (
	ErrInvalidProjectSkillRequirement  = errors.New("project skill requirement is invalid")
	ErrProjectSkillRequirementNotFound = errors.New("project skill requirement not found")
)

type ProjectSkillRequirement struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	ProjectID           *uuid.UUID      `json:"project_id,omitempty"`
	EngagementID        *uuid.UUID      `json:"engagement_id,omitempty"`
	SkillID             uuid.UUID       `json:"skill_id"`
	RequiredProficiency string          `json:"required_proficiency"`
	MinYearsExperience  *float64        `json:"min_years_experience,omitempty"`
	RequiredCount       int32           `json:"required_count"`
	Importance          string          `json:"importance"`
	RequirementSource   string          `json:"requirement_source"`
	Notes               *string         `json:"notes,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
	ProjectName         *string         `json:"project_name,omitempty"`
	ProjectCode         *string         `json:"project_code,omitempty"`
	EngagementTitle     *string         `json:"engagement_title,omitempty"`
	EngagementCode      *string         `json:"engagement_code,omitempty"`
	WorkerProfileID     *uuid.UUID      `json:"worker_profile_id,omitempty"`
	WorkerDisplayName   *string         `json:"worker_display_name,omitempty"`
	WorkerCode          *string         `json:"worker_code,omitempty"`
	SkillName           *string         `json:"skill_name,omitempty"`
	SkillCode           *string         `json:"skill_code,omitempty"`
	SkillType           *string         `json:"skill_type,omitempty"`
	CategoryName        *string         `json:"category_name,omitempty"`
}

type ProjectSkillRequirementInput struct {
	TenantID            uuid.UUID
	ProjectID           *uuid.UUID
	EngagementID        *uuid.UUID
	SkillID             uuid.UUID
	RequiredProficiency string
	MinYearsExperience  *float64
	RequiredCount       int32
	Importance          string
	RequirementSource   string
	Notes               *string
	Metadata            json.RawMessage
}

type ProjectSkillRequirementFilter struct {
	TenantID     uuid.UUID
	ProjectID    *uuid.UUID
	EngagementID *uuid.UUID
	SkillID      *uuid.UUID
	Importance   *string
	Search       *string
}

type ProjectSkillGapRow struct {
	RequirementID          uuid.UUID  `json:"requirement_id"`
	TenantID               uuid.UUID  `json:"tenant_id"`
	ProjectID              *uuid.UUID `json:"project_id,omitempty"`
	ProjectName            *string    `json:"project_name,omitempty"`
	ProjectCode            *string    `json:"project_code,omitempty"`
	EngagementID           *uuid.UUID `json:"engagement_id,omitempty"`
	EngagementTitle        *string    `json:"engagement_title,omitempty"`
	EngagementCode         *string    `json:"engagement_code,omitempty"`
	SkillID                uuid.UUID  `json:"skill_id"`
	SkillName              string     `json:"skill_name"`
	SkillCode              string     `json:"skill_code"`
	SkillType              string     `json:"skill_type"`
	RequiredProficiency    string     `json:"required_proficiency"`
	MinYearsExperience     *float64   `json:"min_years_experience,omitempty"`
	RequiredCount          int32      `json:"required_count"`
	Importance             string     `json:"importance"`
	AssignedMatchCount     int32      `json:"assigned_match_count"`
	AvailableMatchCount    int32      `json:"available_match_count"`
	GapCount               int32      `json:"gap_count"`
	MatchPercent           int32      `json:"match_percent"`
	SinglePersonDependency bool       `json:"single_person_dependency"`
	SuggestedAction        string     `json:"suggested_action"`
}

type SkillGapSummaryRow struct {
	ProjectID                   *uuid.UUID `json:"project_id,omitempty"`
	ProjectName                 string     `json:"project_name"`
	ProjectCode                 *string    `json:"project_code,omitempty"`
	RequirementCount            int32      `json:"requirement_count"`
	MissingSkillCount           int32      `json:"missing_skill_count"`
	MandatoryGapCount           int32      `json:"mandatory_gap_count"`
	AverageMatchPercent         int32      `json:"average_match_percent"`
	SinglePersonDependencyCount int32      `json:"single_person_dependency_count"`
}

type SinglePersonSkillDependency struct {
	RequirementID     uuid.UUID  `json:"requirement_id"`
	ProjectID         *uuid.UUID `json:"project_id,omitempty"`
	ProjectName       *string    `json:"project_name,omitempty"`
	EngagementID      *uuid.UUID `json:"engagement_id,omitempty"`
	EngagementTitle   *string    `json:"engagement_title,omitempty"`
	SkillID           uuid.UUID  `json:"skill_id"`
	SkillName         string     `json:"skill_name"`
	Importance        string     `json:"importance"`
	WorkerProfileID   uuid.UUID  `json:"worker_profile_id"`
	WorkerDisplayName string     `json:"worker_display_name"`
	WorkerCode        *string    `json:"worker_code,omitempty"`
	Proficiency       string     `json:"proficiency"`
	YearsExperience   *float64   `json:"years_experience,omitempty"`
}

func NewProjectSkillRequirement(input ProjectSkillRequirementInput) (*ProjectSkillRequirement, error) {
	if input.TenantID == uuid.Nil || input.SkillID == uuid.Nil {
		return nil, ErrInvalidProjectSkillRequirement
	}
	projectID := cleanUUIDOptional(input.ProjectID)
	engagementID := cleanUUIDOptional(input.EngagementID)
	if projectID == nil && engagementID == nil {
		return nil, ErrInvalidProjectSkillRequirement
	}
	proficiency := normalizeWorkerProfileEnum(input.RequiredProficiency, SkillProficiencyIntermediate)
	if !containsString(skillProficiencies(), proficiency) {
		return nil, ErrInvalidProjectSkillRequirement
	}
	requiredCount := input.RequiredCount
	if requiredCount <= 0 {
		requiredCount = 1
	}
	if input.MinYearsExperience != nil && *input.MinYearsExperience < 0 {
		return nil, ErrInvalidProjectSkillRequirement
	}
	importance := normalizeWorkerProfileEnum(input.Importance, ProjectSkillRequirementImportanceRequired)
	if !containsString(projectSkillRequirementImportances(), importance) {
		return nil, ErrInvalidProjectSkillRequirement
	}
	source := normalizeWorkerProfileEnum(input.RequirementSource, "")
	if source == "" {
		source = ProjectSkillRequirementSourceProject
		if engagementID != nil {
			source = ProjectSkillRequirementSourceEngagement
		}
	}
	if !containsString(projectSkillRequirementSources(), source) {
		return nil, ErrInvalidProjectSkillRequirement
	}
	metadata := normalizeWorkerJSONObject(input.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidProjectSkillRequirement
	}
	now := time.Now().UTC()
	return &ProjectSkillRequirement{
		TenantID:            input.TenantID,
		ProjectID:           projectID,
		EngagementID:        engagementID,
		SkillID:             input.SkillID,
		RequiredProficiency: proficiency,
		MinYearsExperience:  cleanFloatOptional(input.MinYearsExperience),
		RequiredCount:       requiredCount,
		Importance:          importance,
		RequirementSource:   source,
		Notes:               cleanOptional(input.Notes),
		Metadata:            metadata,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

func NormalizeProjectSkillRequirementImportance(value string) string {
	return normalizeWorkerProfileEnum(value, "")
}

func projectSkillRequirementImportances() []string {
	return []string{ProjectSkillRequirementImportanceNiceToHave, ProjectSkillRequirementImportanceRequired, ProjectSkillRequirementImportanceCritical}
}

func projectSkillRequirementSources() []string {
	return []string{ProjectSkillRequirementSourceProject, ProjectSkillRequirementSourceEngagement, ProjectSkillRequirementSourceRole, ProjectSkillRequirementSourceClient, ProjectSkillRequirementSourceCompliance}
}

func skillProficiencies() []string {
	return []string{SkillProficiencyBeginner, SkillProficiencyIntermediate, SkillProficiencyAdvanced, SkillProficiencyExpert}
}

func NormalizeProjectSkillSearch(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}
