package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	HRCaseConfidentialityNormal     = "normal"
	HRCaseConfidentialityRestricted = "restricted"
	HRCaseConfidentialitySensitive  = "sensitive"
	HRCaseConfidentialityGrievance  = "grievance"

	HRCasePriorityLow    = "low"
	HRCasePriorityNormal = "normal"
	HRCasePriorityHigh   = "high"
	HRCasePriorityUrgent = "urgent"

	HRCaseStatusNew               = "new"
	HRCaseStatusOpen              = "open"
	HRCaseStatusInProgress        = "in_progress"
	HRCaseStatusWaitingOnEmployee = "waiting_on_employee"
	HRCaseStatusWaitingOnHR       = "waiting_on_hr"
	HRCaseStatusEscalated         = "escalated"
	HRCaseStatusResolved          = "resolved"
	HRCaseStatusClosed            = "closed"
	HRCaseStatusCancelled         = "cancelled"

	HRCaseVisibilityPublic   = "public"
	HRCaseVisibilityInternal = "internal"
)

var (
	ErrInvalidHRCase           = errors.New("hr case is invalid")
	ErrHRCaseNotFound          = errors.New("hr case not found")
	ErrInvalidHRCaseCategory   = errors.New("hr case category is invalid")
	ErrHRCaseCategoryNotFound  = errors.New("hr case category not found")
	ErrInvalidHRCaseSLA        = errors.New("hr case sla policy is invalid")
	ErrHRCaseSLANotFound       = errors.New("hr case sla policy not found")
	ErrInvalidHRCaseComment    = errors.New("hr case comment is invalid")
	ErrInvalidHRCaseAttachment = errors.New("hr case attachment is invalid")
)

type HRCaseCategory struct {
	ID                     uuid.UUID  `json:"id"`
	TenantID               uuid.UUID  `json:"tenant_id"`
	Code                   string     `json:"code"`
	Name                   string     `json:"name"`
	Description            *string    `json:"description,omitempty"`
	ConfidentialityDefault string     `json:"confidentiality_default"`
	DefaultOwnerRole       *string    `json:"default_owner_role,omitempty"`
	IsActive               bool       `json:"is_active"`
	Inactive               bool       `json:"inactive"`
	CreatedAt              time.Time  `json:"created_at"`
	CreatedBy              *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt              time.Time  `json:"updated_at"`
	UpdatedBy              *uuid.UUID `json:"updated_by,omitempty"`
}

type HRCaseSLAPolicy struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty"`
	CategoryName    *string    `json:"category_name,omitempty"`
	Priority        string     `json:"priority"`
	ResponseHours   int32      `json:"response_hours"`
	ResolutionHours int32      `json:"resolution_hours"`
	EscalationHours int32      `json:"escalation_hours"`
	IsActive        bool       `json:"is_active"`
	Inactive        bool       `json:"inactive"`
	CreatedAt       time.Time  `json:"created_at"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at"`
	UpdatedBy       *uuid.UUID `json:"updated_by,omitempty"`
}

type HRCase struct {
	ID                    uuid.UUID  `json:"id"`
	TenantID              uuid.UUID  `json:"tenant_id"`
	CaseNumber            string     `json:"case_number"`
	CategoryID            *uuid.UUID `json:"category_id,omitempty"`
	CategoryName          *string    `json:"category_name,omitempty"`
	CategoryCode          *string    `json:"category_code,omitempty"`
	CaseType              string     `json:"case_type"`
	Title                 string     `json:"title"`
	Description           string     `json:"description"`
	ConfidentialityLevel  string     `json:"confidentiality_level"`
	RequesterUserID       *uuid.UUID `json:"requester_user_id,omitempty"`
	SubjectEmployeeUserID *uuid.UUID `json:"subject_employee_user_id,omitempty"`
	OwnerUserID           *uuid.UUID `json:"owner_user_id,omitempty"`
	OwnerRole             *string    `json:"owner_role,omitempty"`
	Status                string     `json:"status"`
	Priority              string     `json:"priority"`
	SourceChannel         string     `json:"source_channel"`
	FirstResponseDueAt    *time.Time `json:"first_response_due_at,omitempty"`
	FirstRespondedAt      *time.Time `json:"first_responded_at,omitempty"`
	DueAt                 *time.Time `json:"due_at,omitempty"`
	ResolvedAt            *time.Time `json:"resolved_at,omitempty"`
	ClosedAt              *time.Time `json:"closed_at,omitempty"`
	EscalatedAt           *time.Time `json:"escalated_at,omitempty"`
	EscalationLevel       int32      `json:"escalation_level"`
	LastActivityAt        time.Time  `json:"last_activity_at"`
	ResolutionSummary     *string    `json:"resolution_summary,omitempty"`
	Inactive              bool       `json:"inactive"`
	CreatedAt             time.Time  `json:"created_at"`
	CreatedBy             *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt             time.Time  `json:"updated_at"`
	UpdatedBy             *uuid.UUID `json:"updated_by,omitempty"`
	RequesterEmail        *string    `json:"requester_email,omitempty"`
	SubjectEmail          *string    `json:"subject_email,omitempty"`
	OwnerEmail            *string    `json:"owner_email,omitempty"`
	CommentCount          int32      `json:"comment_count"`
	AttachmentCount       int32      `json:"attachment_count"`
}

type HRCaseComment struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	CaseID       uuid.UUID  `json:"case_id"`
	AuthorUserID *uuid.UUID `json:"author_user_id,omitempty"`
	Visibility   string     `json:"visibility"`
	Body         string     `json:"body"`
	Inactive     bool       `json:"inactive"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
}

type HRCaseAttachment struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	CaseID      uuid.UUID  `json:"case_id"`
	CommentID   *uuid.UUID `json:"comment_id,omitempty"`
	FileName    string     `json:"file_name"`
	ContentType string     `json:"content_type"`
	ObjectKey   string     `json:"object_key"`
	Visibility  string     `json:"visibility"`
	UploadedBy  *uuid.UUID `json:"uploaded_by,omitempty"`
	Inactive    bool       `json:"inactive"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
}

type HRCaseEvent struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	CaseID      uuid.UUID       `json:"case_id"`
	EventType   string          `json:"event_type"`
	FromStatus  *string         `json:"from_status,omitempty"`
	ToStatus    *string         `json:"to_status,omitempty"`
	ActorUserID *uuid.UUID      `json:"actor_user_id,omitempty"`
	Comment     *string         `json:"comment,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type HRCaseFilter struct {
	TenantID              uuid.UUID
	Status                *string
	Priority              *string
	CategoryID            *uuid.UUID
	RequesterUserID       *uuid.UUID
	SubjectEmployeeUserID *uuid.UUID
	OwnerUserID           *uuid.UUID
	ConfidentialityLevel  *string
	Search                *string
	Limit                 int32
	Offset                int32
}

type HRCasePage struct {
	Items      []*HRCase           `json:"items"`
	Total      int64               `json:"total"`
	Summary    []*HRCaseSummaryRow `json:"summary"`
	Categories []*HRCaseCategory   `json:"categories,omitempty"`
}

type HRCaseSummaryRow struct {
	Status          string `json:"status"`
	Priority        string `json:"priority"`
	CaseCount       int32  `json:"case_count"`
	OverdueCount    int32  `json:"overdue_count"`
	EscalatedCount  int32  `json:"escalated_count"`
	RestrictedCount int32  `json:"restricted_count"`
}

type HRCaseWorkspace struct {
	Case        *HRCase             `json:"case"`
	Comments    []*HRCaseComment    `json:"comments"`
	Attachments []*HRCaseAttachment `json:"attachments"`
	Events      []*HRCaseEvent      `json:"events"`
}

func NewHRCaseCategory(item HRCaseCategory) (*HRCaseCategory, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidHRCaseCategory
	}
	item.Code = strings.ToUpper(strings.TrimSpace(item.Code))
	item.Name = strings.TrimSpace(item.Name)
	item.Description = cleanOptional(item.Description)
	item.DefaultOwnerRole = cleanOptional(item.DefaultOwnerRole)
	item.ConfidentialityDefault = normalizeHRCaseConfidentiality(item.ConfidentialityDefault)
	return &item, nil
}

func NewHRCaseSLAPolicy(item HRCaseSLAPolicy) (*HRCaseSLAPolicy, error) {
	if item.TenantID == uuid.Nil {
		return nil, ErrInvalidHRCaseSLA
	}
	item.Priority = normalizeHRCasePriority(item.Priority)
	if item.ResponseHours < 0 || item.ResolutionHours <= 0 || item.EscalationHours <= 0 {
		return nil, ErrInvalidHRCaseSLA
	}
	return &item, nil
}

func NewHRCase(item HRCase) (*HRCase, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Description) == "" {
		return nil, ErrInvalidHRCase
	}
	item.Title = strings.TrimSpace(item.Title)
	item.Description = strings.TrimSpace(item.Description)
	item.CaseType = defaultString(strings.TrimSpace(item.CaseType), "general")
	item.ConfidentialityLevel = normalizeHRCaseConfidentiality(item.ConfidentialityLevel)
	item.Priority = normalizeHRCasePriority(item.Priority)
	item.Status = normalizeHRCaseStatus(item.Status)
	item.SourceChannel = defaultString(strings.TrimSpace(item.SourceChannel), "web")
	return &item, nil
}

func NewHRCaseComment(item HRCaseComment) (*HRCaseComment, error) {
	if item.TenantID == uuid.Nil || item.CaseID == uuid.Nil || strings.TrimSpace(item.Body) == "" {
		return nil, ErrInvalidHRCaseComment
	}
	item.Body = strings.TrimSpace(item.Body)
	item.Visibility = normalizeHRCaseVisibility(item.Visibility)
	return &item, nil
}

func NewHRCaseAttachment(item HRCaseAttachment) (*HRCaseAttachment, error) {
	if item.TenantID == uuid.Nil || item.CaseID == uuid.Nil || strings.TrimSpace(item.FileName) == "" || strings.TrimSpace(item.ObjectKey) == "" {
		return nil, ErrInvalidHRCaseAttachment
	}
	item.FileName = strings.TrimSpace(item.FileName)
	item.ContentType = defaultString(strings.TrimSpace(item.ContentType), "application/octet-stream")
	item.ObjectKey = strings.TrimSpace(item.ObjectKey)
	item.Visibility = normalizeHRCaseVisibility(item.Visibility)
	return &item, nil
}

func NewHRCaseNumber(now time.Time) string {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return fmt.Sprintf("HR-%s-%06d", now.Format("20060102"), now.UnixNano()%1000000)
}

func normalizeHRCaseConfidentiality(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case HRCaseConfidentialityRestricted, HRCaseConfidentialitySensitive, HRCaseConfidentialityGrievance:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return HRCaseConfidentialityNormal
	}
}

func normalizeHRCasePriority(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case HRCasePriorityLow, HRCasePriorityHigh, HRCasePriorityUrgent:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return HRCasePriorityNormal
	}
}

func normalizeHRCaseStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case HRCaseStatusOpen, HRCaseStatusInProgress, HRCaseStatusWaitingOnEmployee, HRCaseStatusWaitingOnHR, HRCaseStatusEscalated, HRCaseStatusResolved, HRCaseStatusClosed, HRCaseStatusCancelled:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return HRCaseStatusNew
	}
}

func normalizeHRCaseVisibility(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), HRCaseVisibilityInternal) {
		return HRCaseVisibilityInternal
	}
	return HRCaseVisibilityPublic
}
