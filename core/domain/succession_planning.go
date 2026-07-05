package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	SuccessionCycleDraft     = "draft"
	SuccessionCycleActive    = "active"
	SuccessionCycleReview    = "review"
	SuccessionCycleClosed    = "closed"
	SuccessionCycleArchived  = "archived"
	SuccessionCycleCancelled = "cancelled"

	SuccessionCriticalityLow      = "low"
	SuccessionCriticalityMedium   = "medium"
	SuccessionCriticalityHigh     = "high"
	SuccessionCriticalityCritical = "critical"

	SuccessionRoleOpen     = "open"
	SuccessionRoleCovered  = "covered"
	SuccessionRoleAtRisk   = "at_risk"
	SuccessionRoleClosed   = "closed"
	SuccessionRoleArchived = "archived"

	SuccessionReadyNow       = "ready_now"
	SuccessionReadySoon      = "ready_soon"
	SuccessionReadyLater     = "ready_later"
	SuccessionReadyEmergency = "emergency_cover"

	SuccessionNominationDraft     = "draft"
	SuccessionNominationNominated = "nominated"
	SuccessionNominationReviewed  = "reviewed"
	SuccessionNominationApproved  = "approved"
	SuccessionNominationRejected  = "rejected"
	SuccessionNominationWithdrawn = "withdrawn"

	SuccessionActionOpen       = "open"
	SuccessionActionInProgress = "in_progress"
	SuccessionActionCompleted  = "completed"
	SuccessionActionOverdue    = "overdue"
	SuccessionActionCancelled  = "cancelled"
)

var (
	ErrInvalidSuccessionReviewCycle        = errors.New("invalid succession review cycle")
	ErrSuccessionReviewCycleNotFound       = errors.New("succession review cycle not found")
	ErrInvalidSuccessionCriticalRole       = errors.New("invalid succession critical role")
	ErrSuccessionCriticalRoleNotFound      = errors.New("succession critical role not found")
	ErrInvalidSuccessionNomination         = errors.New("invalid succession nomination")
	ErrSuccessionNominationNotFound        = errors.New("succession nomination not found")
	ErrInvalidSuccessionDevelopmentAction  = errors.New("invalid succession development action")
	ErrSuccessionDevelopmentActionNotFound = errors.New("succession development action not found")
)

type SuccessionReviewCycle struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	Code                 string          `json:"code"`
	Name                 string          `json:"name"`
	Status               string          `json:"status"`
	StartsOn             *time.Time      `json:"starts_on,omitempty"`
	EndsOn               *time.Time      `json:"ends_on,omitempty"`
	ConfidentialityLevel string          `json:"confidentiality_level"`
	Notes                *string         `json:"notes,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type SuccessionCriticalRole struct {
	ID                            uuid.UUID       `json:"id"`
	TenantID                      uuid.UUID       `json:"tenant_id"`
	CycleID                       *uuid.UUID      `json:"cycle_id,omitempty"`
	Code                          string          `json:"code"`
	Title                         string          `json:"title"`
	DepartmentID                  *uuid.UUID      `json:"department_id,omitempty"`
	DesignationID                 *uuid.UUID      `json:"designation_id,omitempty"`
	IncumbentWorkerProfileID      *uuid.UUID      `json:"incumbent_worker_profile_id,omitempty"`
	EmergencyCoverWorkerProfileID *uuid.UUID      `json:"emergency_cover_worker_profile_id,omitempty"`
	Criticality                   string          `json:"criticality"`
	ImpactLevel                   string          `json:"impact_level"`
	VacancyRisk                   string          `json:"vacancy_risk"`
	AttritionRisk                 string          `json:"attrition_risk"`
	ReadinessTarget               string          `json:"readiness_target"`
	SuccessorRequiredCount        int32           `json:"successor_required_count"`
	RoleSummary                   *string         `json:"role_summary,omitempty"`
	Status                        string          `json:"status"`
	Metadata                      json.RawMessage `json:"metadata,omitempty"`
	Inactive                      bool            `json:"inactive"`
	CreatedAt                     time.Time       `json:"created_at"`
	CreatedBy                     *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                     time.Time       `json:"updated_at"`
	UpdatedBy                     *uuid.UUID      `json:"updated_by,omitempty"`
	DepartmentName                *string         `json:"department_name,omitempty"`
	DesignationName               *string         `json:"designation_name,omitempty"`
	IncumbentName                 *string         `json:"incumbent_name,omitempty"`
	IncumbentCode                 *string         `json:"incumbent_code,omitempty"`
	EmergencyCoverName            *string         `json:"emergency_cover_name,omitempty"`
	EmergencyCoverCode            *string         `json:"emergency_cover_code,omitempty"`
	SuccessorCount                int64           `json:"successor_count"`
	ReadyNowCount                 int64           `json:"ready_now_count"`
}

type SuccessionSuccessorNomination struct {
	ID                 uuid.UUID       `json:"id"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	CriticalRoleID     uuid.UUID       `json:"critical_role_id"`
	WorkerProfileID    uuid.UUID       `json:"worker_profile_id"`
	NominatedBy        *uuid.UUID      `json:"nominated_by,omitempty"`
	ReadinessLevel     string          `json:"readiness_level"`
	ReadinessMonths    int32           `json:"readiness_months"`
	PotentialRating    *string         `json:"potential_rating,omitempty"`
	PerformanceRating  *string         `json:"performance_rating,omitempty"`
	RetentionRisk      string          `json:"retention_risk"`
	MobilityPreference *string         `json:"mobility_preference,omitempty"`
	NominationStatus   string          `json:"nomination_status"`
	DevelopmentNotes   *string         `json:"development_notes,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	Inactive           bool            `json:"inactive"`
	CreatedAt          time.Time       `json:"created_at"`
	CreatedBy          *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt          time.Time       `json:"updated_at"`
	UpdatedBy          *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName  *string         `json:"worker_display_name,omitempty"`
	WorkerCode         *string         `json:"worker_code,omitempty"`
	CriticalRoleTitle  *string         `json:"critical_role_title,omitempty"`
}

type SuccessionDevelopmentAction struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	NominationID        *uuid.UUID      `json:"nomination_id,omitempty"`
	CriticalRoleID      *uuid.UUID      `json:"critical_role_id,omitempty"`
	WorkerProfileID     uuid.UUID       `json:"worker_profile_id"`
	ActionType          string          `json:"action_type"`
	Title               string          `json:"title"`
	LearningCourseID    *uuid.UUID      `json:"learning_course_id,omitempty"`
	LearningPathID      *uuid.UUID      `json:"learning_path_id,omitempty"`
	OwnerUserID         *uuid.UUID      `json:"owner_user_id,omitempty"`
	DueDate             *time.Time      `json:"due_date,omitempty"`
	Status              string          `json:"status"`
	Notes               *string         `json:"notes,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName   *string         `json:"worker_display_name,omitempty"`
	WorkerCode          *string         `json:"worker_code,omitempty"`
	CriticalRoleTitle   *string         `json:"critical_role_title,omitempty"`
	LearningCourseTitle *string         `json:"learning_course_title,omitempty"`
	LearningPathTitle   *string         `json:"learning_path_title,omitempty"`
}

type SuccessionEvent struct {
	ID         uuid.UUID       `json:"id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	SourceType string          `json:"source_type"`
	SourceID   *uuid.UUID      `json:"source_id,omitempty"`
	Action     string          `json:"action"`
	FromStatus *string         `json:"from_status,omitempty"`
	ToStatus   *string         `json:"to_status,omitempty"`
	Remarks    *string         `json:"remarks,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	CreatedBy  *uuid.UUID      `json:"created_by,omitempty"`
}

type SuccessionSummaryRow struct {
	Metric      string `json:"metric"`
	MetricCount int64  `json:"metric_count"`
}

type SuccessionFilter struct {
	TenantID       uuid.UUID
	CycleID        *uuid.UUID
	CriticalRoleID *uuid.UUID
	Status         *string
	Search         *string
	Limit          int32
	Offset         int32
}

func NormalizeSuccessionCycleStatus(value string) string {
	return normalizeAllowed(value, []string{SuccessionCycleDraft, SuccessionCycleActive, SuccessionCycleReview, SuccessionCycleClosed, SuccessionCycleArchived, SuccessionCycleCancelled})
}

func NormalizeSuccessionRisk(value string) string {
	return normalizeAllowed(value, []string{SuccessionCriticalityLow, SuccessionCriticalityMedium, SuccessionCriticalityHigh, SuccessionCriticalityCritical})
}

func NormalizeSuccessionRoleStatus(value string) string {
	return normalizeAllowed(value, []string{SuccessionRoleOpen, SuccessionRoleCovered, SuccessionRoleAtRisk, SuccessionRoleClosed, SuccessionRoleArchived})
}

func NormalizeSuccessionReadiness(value string) string {
	return normalizeAllowed(value, []string{SuccessionReadyNow, SuccessionReadySoon, SuccessionReadyLater, SuccessionReadyEmergency})
}

func NormalizeSuccessionNominationStatus(value string) string {
	return normalizeAllowed(value, []string{SuccessionNominationDraft, SuccessionNominationNominated, SuccessionNominationReviewed, SuccessionNominationApproved, SuccessionNominationRejected, SuccessionNominationWithdrawn})
}

func NormalizeSuccessionActionStatus(value string) string {
	return normalizeAllowed(value, []string{SuccessionActionOpen, SuccessionActionInProgress, SuccessionActionCompleted, SuccessionActionOverdue, SuccessionActionCancelled})
}

func ValidateSuccessionReviewCycle(item *SuccessionReviewCycle) error {
	if item == nil || item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" {
		return ErrInvalidSuccessionReviewCycle
	}
	item.Code = strings.TrimSpace(item.Code)
	item.Name = strings.TrimSpace(item.Name)
	item.Status = defaultAllowed(item.Status, SuccessionCycleDraft, NormalizeSuccessionCycleStatus)
	item.ConfidentialityLevel = defaultAllowed(item.ConfidentialityLevel, "hr_only", func(v string) string {
		return normalizeAllowed(v, []string{"hr_only", "leadership", "restricted"})
	})
	if item.Status == "" || item.ConfidentialityLevel == "" {
		return ErrInvalidSuccessionReviewCycle
	}
	return nil
}

func ValidateSuccessionCriticalRole(item *SuccessionCriticalRole) error {
	if item == nil || item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Title) == "" {
		return ErrInvalidSuccessionCriticalRole
	}
	item.Code = strings.TrimSpace(item.Code)
	item.Title = strings.TrimSpace(item.Title)
	item.Criticality = defaultAllowed(item.Criticality, SuccessionCriticalityHigh, NormalizeSuccessionRisk)
	item.ImpactLevel = defaultAllowed(item.ImpactLevel, SuccessionCriticalityHigh, NormalizeSuccessionRisk)
	item.VacancyRisk = defaultAllowed(item.VacancyRisk, SuccessionCriticalityMedium, NormalizeSuccessionRisk)
	item.AttritionRisk = defaultAllowed(item.AttritionRisk, SuccessionCriticalityMedium, NormalizeSuccessionRisk)
	item.ReadinessTarget = defaultAllowed(item.ReadinessTarget, SuccessionReadySoon, NormalizeSuccessionReadiness)
	item.Status = defaultAllowed(item.Status, SuccessionRoleOpen, NormalizeSuccessionRoleStatus)
	if item.SuccessorRequiredCount <= 0 {
		item.SuccessorRequiredCount = 2
	}
	if item.Criticality == "" || item.ImpactLevel == "" || item.VacancyRisk == "" || item.AttritionRisk == "" || item.ReadinessTarget == "" || item.Status == "" {
		return ErrInvalidSuccessionCriticalRole
	}
	return nil
}

func ValidateSuccessionSuccessorNomination(item *SuccessionSuccessorNomination) error {
	if item == nil || item.TenantID == uuid.Nil || item.CriticalRoleID == uuid.Nil || item.WorkerProfileID == uuid.Nil {
		return ErrInvalidSuccessionNomination
	}
	item.ReadinessLevel = defaultAllowed(item.ReadinessLevel, SuccessionReadySoon, NormalizeSuccessionReadiness)
	item.RetentionRisk = defaultAllowed(item.RetentionRisk, SuccessionCriticalityMedium, NormalizeSuccessionRisk)
	item.NominationStatus = defaultAllowed(item.NominationStatus, SuccessionNominationNominated, NormalizeSuccessionNominationStatus)
	if item.ReadinessMonths < 0 {
		item.ReadinessMonths = 0
	}
	if item.ReadinessLevel == "" || item.RetentionRisk == "" || item.NominationStatus == "" {
		return ErrInvalidSuccessionNomination
	}
	return nil
}

func ValidateSuccessionDevelopmentAction(item *SuccessionDevelopmentAction) error {
	if item == nil || item.TenantID == uuid.Nil || item.WorkerProfileID == uuid.Nil || strings.TrimSpace(item.Title) == "" {
		return ErrInvalidSuccessionDevelopmentAction
	}
	item.Title = strings.TrimSpace(item.Title)
	item.ActionType = defaultSuccessionString(item.ActionType, "development")
	item.Status = defaultAllowed(item.Status, SuccessionActionOpen, NormalizeSuccessionActionStatus)
	if item.Status == "" {
		return ErrInvalidSuccessionDevelopmentAction
	}
	return nil
}

func normalizeAllowed(value string, allowed []string) string {
	normalized := normalizeWorkerProfileEnum(value, "")
	for _, item := range allowed {
		if normalized == item {
			return normalized
		}
	}
	return ""
}

func defaultAllowed(value string, fallback string, normalizer func(string) string) string {
	normalized := normalizer(value)
	if normalized != "" {
		return normalized
	}
	return fallback
}

func defaultSuccessionString(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return normalizeWorkerProfileEnum(value, fallback)
}
