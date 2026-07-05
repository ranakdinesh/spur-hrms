package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	LearningCourseTechnical   = "technical"
	LearningCourseFunctional  = "functional"
	LearningCourseCompliance  = "compliance"
	LearningCourseLeadership  = "leadership"
	LearningCourseBehavioral  = "behavioral"
	LearningCourseAIReadiness = "ai_readiness"
	LearningCourseCustom      = "custom"

	LearningDeliverySelfPaced = "self_paced"
	LearningDeliveryClassroom = "classroom"
	LearningDeliveryVirtual   = "virtual"
	LearningDeliveryBlended   = "blended"
	LearningDeliveryExternal  = "external"

	LearningPathOnboarding  = "onboarding"
	LearningPathCompliance  = "compliance"
	LearningPathUpskilling  = "upskilling"
	LearningPathLeadership  = "leadership"
	LearningPathAIReadiness = "ai_readiness"
	LearningPathCustom      = "custom"

	LearningAssignmentSelf       = "self"
	LearningAssignmentManager    = "manager"
	LearningAssignmentHR         = "hr"
	LearningAssignmentCompliance = "compliance"
	LearningAssignmentSkillGap   = "skill_gap"
	LearningAssignmentAI         = "ai"
	LearningAssignmentManual     = "manual"

	LearningEnrollmentNominated  = "nominated"
	LearningEnrollmentAssigned   = "assigned"
	LearningEnrollmentApproved   = "approved"
	LearningEnrollmentInProgress = "in_progress"
	LearningEnrollmentCompleted  = "completed"
	LearningEnrollmentOverdue    = "overdue"
	LearningEnrollmentWaived     = "waived"
	LearningEnrollmentCancelled  = "cancelled"

	LearningRecommendationSkillGap    = "skill_gap"
	LearningRecommendationCompliance  = "compliance"
	LearningRecommendationPerformance = "performance"
	LearningRecommendationAI          = "ai"
	LearningRecommendationManager     = "manager"
	LearningRecommendationManual      = "manual"

	LearningPriorityLow    = "low"
	LearningPriorityMedium = "medium"
	LearningPriorityHigh   = "high"
	LearningPriorityUrgent = "urgent"

	LearningRecommendationOpen      = "open"
	LearningRecommendationAccepted  = "accepted"
	LearningRecommendationAssigned  = "assigned"
	LearningRecommendationDismissed = "dismissed"
	LearningRecommendationCompleted = "completed"
)

var (
	ErrInvalidLearningCourse          = errors.New("learning course is invalid")
	ErrLearningCourseNotFound         = errors.New("learning course not found")
	ErrInvalidLearningPath            = errors.New("learning path is invalid")
	ErrLearningPathNotFound           = errors.New("learning path not found")
	ErrInvalidLearningPathCourse      = errors.New("learning path course is invalid")
	ErrInvalidLearningEnrollment      = errors.New("learning enrollment is invalid")
	ErrLearningEnrollmentNotFound     = errors.New("learning enrollment not found")
	ErrInvalidLearningRecommendation  = errors.New("learning recommendation is invalid")
	ErrLearningRecommendationNotFound = errors.New("learning recommendation not found")
)

type LearningCourse struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	Code                string          `json:"code"`
	Title               string          `json:"title"`
	Description         *string         `json:"description,omitempty"`
	CourseType          string          `json:"course_type"`
	DeliveryMode        string          `json:"delivery_mode"`
	Provider            *string         `json:"provider,omitempty"`
	DurationMinutes     int32           `json:"duration_minutes"`
	SkillID             *uuid.UUID      `json:"skill_id,omitempty"`
	ComplianceRuleID    *uuid.UUID      `json:"compliance_rule_id,omitempty"`
	Mandatory           bool            `json:"mandatory"`
	AIReadiness         bool            `json:"ai_readiness"`
	CertificateRequired bool            `json:"certificate_required"`
	BudgetAmount        *float64        `json:"budget_amount,omitempty"`
	CurrencyCode        string          `json:"currency_code"`
	IsActive            bool            `json:"is_active"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
	SkillName           *string         `json:"skill_name,omitempty"`
	SkillCode           *string         `json:"skill_code,omitempty"`
	ComplianceRuleTitle *string         `json:"compliance_rule_title,omitempty"`
}

type LearningPath struct {
	ID           uuid.UUID       `json:"id"`
	TenantID     uuid.UUID       `json:"tenant_id"`
	Code         string          `json:"code"`
	Title        string          `json:"title"`
	Description  *string         `json:"description,omitempty"`
	PathType     string          `json:"path_type"`
	TargetRole   *string         `json:"target_role,omitempty"`
	SkillID      *uuid.UUID      `json:"skill_id,omitempty"`
	AIReadiness  bool            `json:"ai_readiness"`
	IsActive     bool            `json:"is_active"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	Inactive     bool            `json:"inactive"`
	CreatedAt    time.Time       `json:"created_at"`
	CreatedBy    *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt    time.Time       `json:"updated_at"`
	UpdatedBy    *uuid.UUID      `json:"updated_by,omitempty"`
	SkillName    *string         `json:"skill_name,omitempty"`
	CourseCount  int32           `json:"course_count,omitempty"`
	TotalMinutes int32           `json:"total_minutes,omitempty"`
}

type LearningPathCourse struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	PathID          uuid.UUID  `json:"path_id"`
	CourseID        uuid.UUID  `json:"course_id"`
	SortOrder       int32      `json:"sort_order"`
	Required        bool       `json:"required"`
	Inactive        bool       `json:"inactive"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"`
	UpdatedBy       *uuid.UUID `json:"updated_by,omitempty"`
	CourseCode      *string    `json:"course_code,omitempty"`
	CourseTitle     *string    `json:"course_title,omitempty"`
	CourseType      *string    `json:"course_type,omitempty"`
	DeliveryMode    *string    `json:"delivery_mode,omitempty"`
	DurationMinutes int32      `json:"duration_minutes,omitempty"`
}

type LearningEnrollment struct {
	ID                     uuid.UUID       `json:"id"`
	TenantID               uuid.UUID       `json:"tenant_id"`
	CourseID               uuid.UUID       `json:"course_id"`
	PathID                 *uuid.UUID      `json:"path_id,omitempty"`
	WorkerProfileID        uuid.UUID       `json:"worker_profile_id"`
	AssignmentSource       string          `json:"assignment_source"`
	Status                 string          `json:"status"`
	NominatedBy            *uuid.UUID      `json:"nominated_by,omitempty"`
	AssignedBy             *uuid.UUID      `json:"assigned_by,omitempty"`
	DueDate                *time.Time      `json:"due_date,omitempty"`
	StartedAt              *time.Time      `json:"started_at,omitempty"`
	CompletedAt            *time.Time      `json:"completed_at,omitempty"`
	Score                  *float64        `json:"score,omitempty"`
	CertificateURL         *string         `json:"certificate_url,omitempty"`
	CertificateFileName    *string         `json:"certificate_file_name,omitempty"`
	CertificateContentType *string         `json:"certificate_content_type,omitempty"`
	Notes                  *string         `json:"notes,omitempty"`
	Metadata               json.RawMessage `json:"metadata,omitempty"`
	Inactive               bool            `json:"inactive"`
	CreatedAt              time.Time       `json:"created_at"`
	CreatedBy              *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt              time.Time       `json:"updated_at"`
	UpdatedBy              *uuid.UUID      `json:"updated_by,omitempty"`
	CourseTitle            *string         `json:"course_title,omitempty"`
	CourseCode             *string         `json:"course_code,omitempty"`
	CourseType             *string         `json:"course_type,omitempty"`
	PathTitle              *string         `json:"path_title,omitempty"`
	WorkerDisplayName      *string         `json:"worker_display_name,omitempty"`
	WorkerCode             *string         `json:"worker_code,omitempty"`
	AIReadiness            bool            `json:"ai_readiness"`
	Mandatory              bool            `json:"mandatory"`
}

type LearningRecommendation struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	WorkerProfileID   *uuid.UUID      `json:"worker_profile_id,omitempty"`
	SkillID           *uuid.UUID      `json:"skill_id,omitempty"`
	CourseID          *uuid.UUID      `json:"course_id,omitempty"`
	PathID            *uuid.UUID      `json:"path_id,omitempty"`
	SourceType        string          `json:"source_type"`
	Reason            string          `json:"reason"`
	Priority          string          `json:"priority"`
	ConfidenceScore   *float64        `json:"confidence_score,omitempty"`
	Status            string          `json:"status"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedBy         *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName *string         `json:"worker_display_name,omitempty"`
	WorkerCode        *string         `json:"worker_code,omitempty"`
	SkillName         *string         `json:"skill_name,omitempty"`
	CourseTitle       *string         `json:"course_title,omitempty"`
	PathTitle         *string         `json:"path_title,omitempty"`
}

type LearningCourseInput struct {
	TenantID            uuid.UUID
	Code                string
	Title               string
	Description         *string
	CourseType          string
	DeliveryMode        string
	Provider            *string
	DurationMinutes     int32
	SkillID             *uuid.UUID
	ComplianceRuleID    *uuid.UUID
	Mandatory           bool
	AIReadiness         bool
	CertificateRequired bool
	BudgetAmount        *float64
	CurrencyCode        string
	IsActive            bool
	Metadata            json.RawMessage
}

type LearningPathInput struct {
	TenantID    uuid.UUID
	Code        string
	Title       string
	Description *string
	PathType    string
	TargetRole  *string
	SkillID     *uuid.UUID
	AIReadiness bool
	IsActive    bool
	Metadata    json.RawMessage
}

type LearningEnrollmentInput struct {
	TenantID         uuid.UUID
	CourseID         uuid.UUID
	PathID           *uuid.UUID
	WorkerProfileID  uuid.UUID
	AssignmentSource string
	Status           string
	NominatedBy      *uuid.UUID
	AssignedBy       *uuid.UUID
	DueDate          *time.Time
	Notes            *string
	Metadata         json.RawMessage
}

type LearningRecommendationInput struct {
	TenantID        uuid.UUID
	WorkerProfileID *uuid.UUID
	SkillID         *uuid.UUID
	CourseID        *uuid.UUID
	PathID          *uuid.UUID
	SourceType      string
	Reason          string
	Priority        string
	ConfidenceScore *float64
	Status          string
	Metadata        json.RawMessage
}

type LearningCourseFilter struct {
	TenantID    uuid.UUID
	CourseType  *string
	SkillID     *uuid.UUID
	Mandatory   *bool
	AIReadiness *bool
	IsActive    *bool
	Search      *string
	Limit       int32
	Offset      int32
}

type LearningPathFilter struct {
	TenantID    uuid.UUID
	PathType    *string
	SkillID     *uuid.UUID
	AIReadiness *bool
	IsActive    *bool
	Search      *string
	Limit       int32
	Offset      int32
}

type LearningEnrollmentFilter struct {
	TenantID         uuid.UUID
	WorkerProfileID  *uuid.UUID
	CourseID         *uuid.UUID
	Status           *string
	AssignmentSource *string
	Search           *string
	Limit            int32
	Offset           int32
}

type LearningRecommendationFilter struct {
	TenantID        uuid.UUID
	WorkerProfileID *uuid.UUID
	SkillID         *uuid.UUID
	SourceType      *string
	Status          *string
	Search          *string
	Limit           int32
	Offset          int32
}

type LearningSummaryRow struct {
	Metric      string `json:"metric"`
	MetricCount int32  `json:"metric_count"`
}

func NewLearningCourse(input LearningCourseInput) (*LearningCourse, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code := normalizeSkillCode(input.Code)
	title := strings.TrimSpace(input.Title)
	courseType := normalizeWorkerProfileEnum(input.CourseType, LearningCourseTechnical)
	deliveryMode := normalizeWorkerProfileEnum(input.DeliveryMode, LearningDeliverySelfPaced)
	currency := strings.ToUpper(strings.TrimSpace(input.CurrencyCode))
	if currency == "" {
		currency = "INR"
	}
	if code == "" || title == "" || input.DurationMinutes < 0 || !containsString(learningCourseTypes(), courseType) || !containsString(learningDeliveryModes(), deliveryMode) {
		return nil, ErrInvalidLearningCourse
	}
	if input.BudgetAmount != nil && *input.BudgetAmount < 0 {
		return nil, ErrInvalidLearningCourse
	}
	metadata, err := validLearningMetadata(input.Metadata)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &LearningCourse{TenantID: input.TenantID, Code: code, Title: title, Description: cleanOptional(input.Description), CourseType: courseType, DeliveryMode: deliveryMode, Provider: cleanOptional(input.Provider), DurationMinutes: input.DurationMinutes, SkillID: cleanUUIDOptional(input.SkillID), ComplianceRuleID: cleanUUIDOptional(input.ComplianceRuleID), Mandatory: input.Mandatory, AIReadiness: input.AIReadiness, CertificateRequired: input.CertificateRequired, BudgetAmount: input.BudgetAmount, CurrencyCode: currency, IsActive: input.IsActive, Metadata: metadata, CreatedAt: now, UpdatedAt: now}, nil
}

func NewLearningPath(input LearningPathInput) (*LearningPath, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code := normalizeSkillCode(input.Code)
	title := strings.TrimSpace(input.Title)
	pathType := normalizeWorkerProfileEnum(input.PathType, LearningPathUpskilling)
	if code == "" || title == "" || !containsString(learningPathTypes(), pathType) {
		return nil, ErrInvalidLearningPath
	}
	metadata, err := validLearningMetadata(input.Metadata)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &LearningPath{TenantID: input.TenantID, Code: code, Title: title, Description: cleanOptional(input.Description), PathType: pathType, TargetRole: cleanOptional(input.TargetRole), SkillID: cleanUUIDOptional(input.SkillID), AIReadiness: input.AIReadiness, IsActive: input.IsActive, Metadata: metadata, CreatedAt: now, UpdatedAt: now}, nil
}

func NewLearningEnrollment(input LearningEnrollmentInput) (*LearningEnrollment, error) {
	if input.TenantID == uuid.Nil || input.CourseID == uuid.Nil || input.WorkerProfileID == uuid.Nil {
		return nil, ErrInvalidLearningEnrollment
	}
	source := normalizeWorkerProfileEnum(input.AssignmentSource, LearningAssignmentHR)
	status := normalizeWorkerProfileEnum(input.Status, LearningEnrollmentAssigned)
	if !containsString(learningAssignmentSources(), source) || !containsString(learningEnrollmentStatuses(), status) {
		return nil, ErrInvalidLearningEnrollment
	}
	metadata, err := validLearningMetadata(input.Metadata)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &LearningEnrollment{TenantID: input.TenantID, CourseID: input.CourseID, PathID: cleanUUIDOptional(input.PathID), WorkerProfileID: input.WorkerProfileID, AssignmentSource: source, Status: status, NominatedBy: cleanUUIDOptional(input.NominatedBy), AssignedBy: cleanUUIDOptional(input.AssignedBy), DueDate: datePtrUTC(input.DueDate), Notes: cleanOptional(input.Notes), Metadata: metadata, CreatedAt: now, UpdatedAt: now}, nil
}

func NewLearningRecommendation(input LearningRecommendationInput) (*LearningRecommendation, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidLearningRecommendation
	}
	source := normalizeWorkerProfileEnum(input.SourceType, LearningRecommendationManual)
	priority := normalizeWorkerProfileEnum(input.Priority, LearningPriorityMedium)
	status := normalizeWorkerProfileEnum(input.Status, LearningRecommendationOpen)
	reason := strings.TrimSpace(input.Reason)
	if reason == "" || !containsString(learningRecommendationSources(), source) || !containsString(learningPriorities(), priority) || !containsString(learningRecommendationStatuses(), status) {
		return nil, ErrInvalidLearningRecommendation
	}
	if input.ConfidenceScore != nil && (*input.ConfidenceScore < 0 || *input.ConfidenceScore > 100) {
		return nil, ErrInvalidLearningRecommendation
	}
	metadata, err := validLearningMetadata(input.Metadata)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &LearningRecommendation{TenantID: input.TenantID, WorkerProfileID: cleanUUIDOptional(input.WorkerProfileID), SkillID: cleanUUIDOptional(input.SkillID), CourseID: cleanUUIDOptional(input.CourseID), PathID: cleanUUIDOptional(input.PathID), SourceType: source, Reason: reason, Priority: priority, ConfidenceScore: input.ConfidenceScore, Status: status, Metadata: metadata, CreatedAt: now, UpdatedAt: now}, nil
}

func NormalizeLearningEnrollmentStatus(value string) string {
	status := normalizeWorkerProfileEnum(value, "")
	if !containsString(learningEnrollmentStatuses(), status) {
		return ""
	}
	return status
}

func NormalizeLearningRecommendationStatus(value string) string {
	status := normalizeWorkerProfileEnum(value, "")
	if !containsString(learningRecommendationStatuses(), status) {
		return ""
	}
	return status
}

func validLearningMetadata(value json.RawMessage) (json.RawMessage, error) {
	metadata := normalizeWorkerJSONObject(value, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidLearningCourse
	}
	return metadata, nil
}

func learningCourseTypes() []string {
	return []string{LearningCourseTechnical, LearningCourseFunctional, LearningCourseCompliance, LearningCourseLeadership, LearningCourseBehavioral, LearningCourseAIReadiness, LearningCourseCustom}
}

func learningDeliveryModes() []string {
	return []string{LearningDeliverySelfPaced, LearningDeliveryClassroom, LearningDeliveryVirtual, LearningDeliveryBlended, LearningDeliveryExternal}
}

func learningPathTypes() []string {
	return []string{LearningPathOnboarding, LearningPathCompliance, LearningPathUpskilling, LearningPathLeadership, LearningPathAIReadiness, LearningPathCustom}
}

func learningAssignmentSources() []string {
	return []string{LearningAssignmentSelf, LearningAssignmentManager, LearningAssignmentHR, LearningAssignmentCompliance, LearningAssignmentSkillGap, LearningAssignmentAI, LearningAssignmentManual}
}

func learningEnrollmentStatuses() []string {
	return []string{LearningEnrollmentNominated, LearningEnrollmentAssigned, LearningEnrollmentApproved, LearningEnrollmentInProgress, LearningEnrollmentCompleted, LearningEnrollmentOverdue, LearningEnrollmentWaived, LearningEnrollmentCancelled}
}

func learningRecommendationSources() []string {
	return []string{LearningRecommendationSkillGap, LearningRecommendationCompliance, LearningRecommendationPerformance, LearningRecommendationAI, LearningRecommendationManager, LearningRecommendationManual}
}

func learningPriorities() []string {
	return []string{LearningPriorityLow, LearningPriorityMedium, LearningPriorityHigh, LearningPriorityUrgent}
}

func learningRecommendationStatuses() []string {
	return []string{LearningRecommendationOpen, LearningRecommendationAccepted, LearningRecommendationAssigned, LearningRecommendationDismissed, LearningRecommendationCompleted}
}
