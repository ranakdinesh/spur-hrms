package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	TalentOpportunityTypeProjectAssignment = "project_assignment"
	TalentOpportunityTypeGig               = "gig"
	TalentOpportunityTypeRole              = "role"
	TalentOpportunityTypeMentorship        = "mentorship"
	TalentOpportunityTypeStretch           = "stretch"
	TalentOpportunityTypeBackfill          = "backfill"

	TalentOpportunityStatusDraft     = "draft"
	TalentOpportunityStatusOpen      = "open"
	TalentOpportunityStatusPaused    = "paused"
	TalentOpportunityStatusFilled    = "filled"
	TalentOpportunityStatusClosed    = "closed"
	TalentOpportunityStatusCancelled = "cancelled"

	TalentOpportunityVisibilityAllWorkers        = "all_workers"
	TalentOpportunityVisibilityInvited           = "invited"
	TalentOpportunityVisibilityManagerNomination = "manager_nomination"

	TalentOpportunityPriorityLow      = "low"
	TalentOpportunityPriorityNormal   = "normal"
	TalentOpportunityPriorityHigh     = "high"
	TalentOpportunityPriorityCritical = "critical"

	TalentOpportunityLocationOnsite   = "onsite"
	TalentOpportunityLocationRemote   = "remote"
	TalentOpportunityLocationHybrid   = "hybrid"
	TalentOpportunityLocationFlexible = "flexible"

	TalentFallbackNotNeeded   = "not_needed"
	TalentFallbackMonitoring  = "monitoring"
	TalentFallbackRecommended = "recommended"
	TalentFallbackOpened      = "opened"

	TalentApplicationStatusRecommended = "recommended"
	TalentApplicationStatusInvited     = "invited"
	TalentApplicationStatusInterested  = "interested"
	TalentApplicationStatusApplied     = "applied"
	TalentApplicationStatusAccepted    = "accepted"
	TalentApplicationStatusDeclined    = "declined"
	TalentApplicationStatusWithdrawn   = "withdrawn"
	TalentApplicationStatusRejected    = "rejected"
	TalentApplicationStatusAssigned    = "assigned"
)

var (
	ErrInvalidTalentOpportunity      = errors.New("talent marketplace opportunity is invalid")
	ErrTalentOpportunityNotFound     = errors.New("talent marketplace opportunity not found")
	ErrInvalidTalentApplication      = errors.New("talent marketplace application is invalid")
	ErrTalentApplicationNotFound     = errors.New("talent marketplace application not found")
	ErrInvalidTalentMarketplaceEvent = errors.New("talent marketplace event is invalid")
)

type TalentMarketplaceOpportunity struct {
	ID                       uuid.UUID       `json:"id"`
	TenantID                 uuid.UUID       `json:"tenant_id"`
	ProjectID                *uuid.UUID      `json:"project_id,omitempty"`
	EngagementID             *uuid.UUID      `json:"engagement_id,omitempty"`
	SourceRequirementID      *uuid.UUID      `json:"source_requirement_id,omitempty"`
	JobPostingID             *uuid.UUID      `json:"job_posting_id,omitempty"`
	Title                    string          `json:"title"`
	Description              *string         `json:"description,omitempty"`
	OpportunityType          string          `json:"opportunity_type"`
	Status                   string          `json:"status"`
	Visibility               string          `json:"visibility"`
	Priority                 string          `json:"priority"`
	Seats                    int32           `json:"seats"`
	LocationMode             string          `json:"location_mode"`
	MinAllocationPercent     *int32          `json:"min_allocation_percent,omitempty"`
	DurationLabel            *string         `json:"duration_label,omitempty"`
	StartDate                *time.Time      `json:"start_date,omitempty"`
	DueDate                  *time.Time      `json:"due_date,omitempty"`
	CandidateFallbackEnabled bool            `json:"candidate_fallback_enabled"`
	CandidateFallbackStatus  string          `json:"candidate_fallback_status"`
	Metadata                 json.RawMessage `json:"metadata,omitempty"`
	Inactive                 bool            `json:"inactive"`
	CreatedAt                time.Time       `json:"created_at"`
	CreatedBy                *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                time.Time       `json:"updated_at"`
	UpdatedBy                *uuid.UUID      `json:"updated_by,omitempty"`
	ProjectName              *string         `json:"project_name,omitempty"`
	ProjectCode              *string         `json:"project_code,omitempty"`
	EngagementTitle          *string         `json:"engagement_title,omitempty"`
	EngagementCode           *string         `json:"engagement_code,omitempty"`
	JobPostingTitle          *string         `json:"job_posting_title,omitempty"`
	JobPostingCode           *string         `json:"job_posting_code,omitempty"`
	ApplicationCount         int32           `json:"application_count"`
	RecommendedCount         int32           `json:"recommended_count"`
	SelectedCount            int32           `json:"selected_count"`
}

type TalentMarketplaceOpportunityInput struct {
	TenantID                 uuid.UUID
	ProjectID                *uuid.UUID
	EngagementID             *uuid.UUID
	SourceRequirementID      *uuid.UUID
	JobPostingID             *uuid.UUID
	Title                    string
	Description              *string
	OpportunityType          string
	Status                   string
	Visibility               string
	Priority                 string
	Seats                    int32
	LocationMode             string
	MinAllocationPercent     *int32
	DurationLabel            *string
	StartDate                *time.Time
	DueDate                  *time.Time
	CandidateFallbackEnabled bool
	CandidateFallbackStatus  string
	Metadata                 json.RawMessage
}

type TalentMarketplaceOpportunityFilter struct {
	TenantID        uuid.UUID
	ProjectID       *uuid.UUID
	EngagementID    *uuid.UUID
	Status          *string
	OpportunityType *string
	Search          *string
}

type TalentMarketplaceApplication struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	OpportunityID     uuid.UUID       `json:"opportunity_id"`
	WorkerProfileID   uuid.UUID       `json:"worker_profile_id"`
	Status            string          `json:"status"`
	MatchScore        *float64        `json:"match_score,omitempty"`
	MatchReasons      json.RawMessage `json:"match_reasons,omitempty"`
	WorkerNote        *string         `json:"worker_note,omitempty"`
	ManagerNote       *string         `json:"manager_note,omitempty"`
	DecidedAt         *time.Time      `json:"decided_at,omitempty"`
	DecidedBy         *uuid.UUID      `json:"decided_by,omitempty"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedBy         *uuid.UUID      `json:"updated_by,omitempty"`
	OpportunityTitle  *string         `json:"opportunity_title,omitempty"`
	OpportunityStatus *string         `json:"opportunity_status,omitempty"`
	WorkerDisplayName *string         `json:"worker_display_name,omitempty"`
	WorkerCode        *string         `json:"worker_code,omitempty"`
	ProjectName       *string         `json:"project_name,omitempty"`
	EngagementTitle   *string         `json:"engagement_title,omitempty"`
}

type TalentMarketplaceApplicationInput struct {
	TenantID        uuid.UUID
	OpportunityID   uuid.UUID
	WorkerProfileID uuid.UUID
	Status          string
	MatchScore      *float64
	MatchReasons    json.RawMessage
	WorkerNote      *string
	ManagerNote     *string
}

type TalentMarketplaceApplicationFilter struct {
	TenantID        uuid.UUID
	OpportunityID   *uuid.UUID
	WorkerProfileID *uuid.UUID
	Status          *string
}

type TalentMarketplaceRecommendation struct {
	WorkerProfileID    uuid.UUID       `json:"worker_profile_id"`
	WorkerDisplayName  string          `json:"worker_display_name"`
	WorkerCode         *string         `json:"worker_code,omitempty"`
	RequiredSkillCount int32           `json:"required_skill_count"`
	MatchedSkillCount  int32           `json:"matched_skill_count"`
	MissingSkillCount  int32           `json:"missing_skill_count"`
	MatchScore         float64         `json:"match_score"`
	MatchReasons       json.RawMessage `json:"match_reasons,omitempty"`
	ApplicationID      *uuid.UUID      `json:"application_id,omitempty"`
	ApplicationStatus  *string         `json:"application_status,omitempty"`
}

type TalentMarketplaceEvent struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	OpportunityID     *uuid.UUID      `json:"opportunity_id,omitempty"`
	ApplicationID     *uuid.UUID      `json:"application_id,omitempty"`
	ActorUserID       *uuid.UUID      `json:"actor_user_id,omitempty"`
	Action            string          `json:"action"`
	FromStatus        *string         `json:"from_status,omitempty"`
	ToStatus          *string         `json:"to_status,omitempty"`
	Notes             *string         `json:"notes,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	OpportunityTitle  *string         `json:"opportunity_title,omitempty"`
	WorkerProfileID   *uuid.UUID      `json:"worker_profile_id,omitempty"`
	WorkerDisplayName *string         `json:"worker_display_name,omitempty"`
}

type TalentMarketplaceEventFilter struct {
	TenantID      uuid.UUID
	OpportunityID *uuid.UUID
	ApplicationID *uuid.UUID
}

func NewTalentMarketplaceOpportunity(input TalentMarketplaceOpportunityInput) (*TalentMarketplaceOpportunity, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTalentOpportunity
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidTalentOpportunity
	}
	status := normalizeWorkerProfileEnum(input.Status, TalentOpportunityStatusDraft)
	if !containsString(talentOpportunityStatuses(), status) {
		return nil, ErrInvalidTalentOpportunity
	}
	opportunityType := normalizeWorkerProfileEnum(input.OpportunityType, TalentOpportunityTypeProjectAssignment)
	if !containsString(talentOpportunityTypes(), opportunityType) {
		return nil, ErrInvalidTalentOpportunity
	}
	visibility := normalizeWorkerProfileEnum(input.Visibility, TalentOpportunityVisibilityAllWorkers)
	if !containsString(talentOpportunityVisibilities(), visibility) {
		return nil, ErrInvalidTalentOpportunity
	}
	priority := normalizeWorkerProfileEnum(input.Priority, TalentOpportunityPriorityNormal)
	if !containsString(talentOpportunityPriorities(), priority) {
		return nil, ErrInvalidTalentOpportunity
	}
	locationMode := normalizeWorkerProfileEnum(input.LocationMode, TalentOpportunityLocationFlexible)
	if !containsString(talentOpportunityLocations(), locationMode) {
		return nil, ErrInvalidTalentOpportunity
	}
	seats := input.Seats
	if seats <= 0 {
		seats = 1
	}
	if input.MinAllocationPercent != nil && (*input.MinAllocationPercent < 1 || *input.MinAllocationPercent > 100) {
		return nil, ErrInvalidTalentOpportunity
	}
	startDate := datePtrUTC(input.StartDate)
	dueDate := datePtrUTC(input.DueDate)
	if startDate != nil && dueDate != nil && dueDate.Before(*startDate) {
		return nil, ErrInvalidTalentOpportunity
	}
	fallbackStatus := normalizeWorkerProfileEnum(input.CandidateFallbackStatus, TalentFallbackNotNeeded)
	if input.CandidateFallbackEnabled && fallbackStatus == TalentFallbackNotNeeded {
		fallbackStatus = TalentFallbackMonitoring
	}
	if !containsString(talentFallbackStatuses(), fallbackStatus) {
		return nil, ErrInvalidTalentOpportunity
	}
	metadata := normalizeWorkerJSONObject(input.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidTalentOpportunity
	}
	now := time.Now().UTC()
	return &TalentMarketplaceOpportunity{
		TenantID:                 input.TenantID,
		ProjectID:                cleanUUIDOptional(input.ProjectID),
		EngagementID:             cleanUUIDOptional(input.EngagementID),
		SourceRequirementID:      cleanUUIDOptional(input.SourceRequirementID),
		JobPostingID:             cleanUUIDOptional(input.JobPostingID),
		Title:                    title,
		Description:              cleanOptional(input.Description),
		OpportunityType:          opportunityType,
		Status:                   status,
		Visibility:               visibility,
		Priority:                 priority,
		Seats:                    seats,
		LocationMode:             locationMode,
		MinAllocationPercent:     cleanInt32Optional(input.MinAllocationPercent),
		DurationLabel:            cleanOptional(input.DurationLabel),
		StartDate:                startDate,
		DueDate:                  dueDate,
		CandidateFallbackEnabled: input.CandidateFallbackEnabled,
		CandidateFallbackStatus:  fallbackStatus,
		Metadata:                 metadata,
		CreatedAt:                now,
		UpdatedAt:                now,
	}, nil
}

func NewTalentMarketplaceApplication(input TalentMarketplaceApplicationInput) (*TalentMarketplaceApplication, error) {
	if input.TenantID == uuid.Nil || input.OpportunityID == uuid.Nil || input.WorkerProfileID == uuid.Nil {
		return nil, ErrInvalidTalentApplication
	}
	status := normalizeWorkerProfileEnum(input.Status, TalentApplicationStatusApplied)
	if !containsString(talentApplicationStatuses(), status) {
		return nil, ErrInvalidTalentApplication
	}
	if input.MatchScore != nil && (*input.MatchScore < 0 || *input.MatchScore > 100) {
		return nil, ErrInvalidTalentApplication
	}
	reasons := normalizeWorkerJSONObject(input.MatchReasons, "{}")
	if !json.Valid(reasons) || !jsonObject(reasons) {
		return nil, ErrInvalidTalentApplication
	}
	now := time.Now().UTC()
	return &TalentMarketplaceApplication{
		TenantID:        input.TenantID,
		OpportunityID:   input.OpportunityID,
		WorkerProfileID: input.WorkerProfileID,
		Status:          status,
		MatchScore:      cleanFloatOptional(input.MatchScore),
		MatchReasons:    reasons,
		WorkerNote:      cleanOptional(input.WorkerNote),
		ManagerNote:     cleanOptional(input.ManagerNote),
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func NormalizeTalentMarketplaceSearch(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func ValidateTalentApplicationStatus(value string) (string, error) {
	status := normalizeWorkerProfileEnum(value, "")
	if !containsString(talentApplicationStatuses(), status) {
		return "", ErrInvalidTalentApplication
	}
	return status, nil
}

func ValidateTalentFallbackStatus(value string) (string, error) {
	status := normalizeWorkerProfileEnum(value, TalentFallbackRecommended)
	if !containsString(talentFallbackStatuses(), status) || status == TalentFallbackNotNeeded {
		return "", ErrInvalidTalentOpportunity
	}
	return status, nil
}

func talentOpportunityTypes() []string {
	return []string{TalentOpportunityTypeProjectAssignment, TalentOpportunityTypeGig, TalentOpportunityTypeRole, TalentOpportunityTypeMentorship, TalentOpportunityTypeStretch, TalentOpportunityTypeBackfill}
}

func talentOpportunityStatuses() []string {
	return []string{TalentOpportunityStatusDraft, TalentOpportunityStatusOpen, TalentOpportunityStatusPaused, TalentOpportunityStatusFilled, TalentOpportunityStatusClosed, TalentOpportunityStatusCancelled}
}

func talentOpportunityVisibilities() []string {
	return []string{TalentOpportunityVisibilityAllWorkers, TalentOpportunityVisibilityInvited, TalentOpportunityVisibilityManagerNomination}
}

func talentOpportunityPriorities() []string {
	return []string{TalentOpportunityPriorityLow, TalentOpportunityPriorityNormal, TalentOpportunityPriorityHigh, TalentOpportunityPriorityCritical}
}

func talentOpportunityLocations() []string {
	return []string{TalentOpportunityLocationOnsite, TalentOpportunityLocationRemote, TalentOpportunityLocationHybrid, TalentOpportunityLocationFlexible}
}

func talentFallbackStatuses() []string {
	return []string{TalentFallbackNotNeeded, TalentFallbackMonitoring, TalentFallbackRecommended, TalentFallbackOpened}
}

func talentApplicationStatuses() []string {
	return []string{TalentApplicationStatusRecommended, TalentApplicationStatusInvited, TalentApplicationStatusInterested, TalentApplicationStatusApplied, TalentApplicationStatusAccepted, TalentApplicationStatusDeclined, TalentApplicationStatusWithdrawn, TalentApplicationStatusRejected, TalentApplicationStatusAssigned}
}

func cleanInt32Optional(value *int32) *int32 {
	if value == nil {
		return nil
	}
	return value
}
