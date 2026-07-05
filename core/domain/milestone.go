package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MilestoneStatusDraft     = "draft"
	MilestoneStatusOpen      = "open"
	MilestoneStatusSubmitted = "submitted"
	MilestoneStatusAccepted  = "accepted"
	MilestoneStatusRejected  = "rejected"
	MilestoneStatusCancelled = "cancelled"

	MilestoneEventCreated   = "created"
	MilestoneEventUpdated   = "updated"
	MilestoneEventSubmitted = "submitted"
	MilestoneEventAccepted  = "accepted"
	MilestoneEventRejected  = "rejected"
	MilestoneEventCancelled = "cancelled"
)

var (
	ErrInvalidMilestoneID             = errors.New("milestone_id is required")
	ErrInvalidMilestone               = errors.New("milestone is invalid")
	ErrInvalidMilestoneStatus         = errors.New("milestone status is invalid")
	ErrInvalidMilestoneAmount         = errors.New("milestone amount must be non-negative")
	ErrInvalidMilestoneCurrency       = errors.New("milestone currency_code must be a 3-letter code")
	ErrInvalidMilestonePaymentTrigger = errors.New("milestone payment_trigger must be a valid JSON object")
	ErrInvalidMilestoneMetadata       = errors.New("milestone metadata must be a valid JSON object")
	ErrInvalidMilestoneWorkflow       = errors.New("milestone workflow transition is invalid")
)

type ProjectMilestone struct {
	ID                 uuid.UUID       `json:"id"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	ProjectID          uuid.UUID       `json:"project_id"`
	EngagementID       *uuid.UUID      `json:"engagement_id,omitempty"`
	MilestoneCode      *string         `json:"milestone_code,omitempty"`
	Title              string          `json:"title"`
	Description        *string         `json:"description,omitempty"`
	AcceptanceCriteria *string         `json:"acceptance_criteria,omitempty"`
	DueDate            *time.Time      `json:"due_date,omitempty"`
	Status             string          `json:"status"`
	Amount             *float64        `json:"amount,omitempty"`
	CurrencyCode       string          `json:"currency_code"`
	PaymentTrigger     json.RawMessage `json:"payment_trigger,omitempty"`
	SubmittedAt        *time.Time      `json:"submitted_at,omitempty"`
	SubmittedBy        *uuid.UUID      `json:"submitted_by,omitempty"`
	AcceptedAt         *time.Time      `json:"accepted_at,omitempty"`
	AcceptedBy         *uuid.UUID      `json:"accepted_by,omitempty"`
	RejectedAt         *time.Time      `json:"rejected_at,omitempty"`
	RejectedBy         *uuid.UUID      `json:"rejected_by,omitempty"`
	ReviewComment      *string         `json:"review_comment,omitempty"`
	Notes              *string         `json:"notes,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	Inactive           bool            `json:"inactive"`
	CreatedAt          time.Time       `json:"created_at"`
	CreatedBy          *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt          time.Time       `json:"updated_at"`
	UpdatedBy          *uuid.UUID      `json:"updated_by,omitempty"`
}

type ProjectMilestoneListItem struct {
	ProjectMilestone
	ProjectName       string     `json:"project_name"`
	ProjectCode       *string    `json:"project_code,omitempty"`
	ProjectManagerID  *uuid.UUID `json:"project_manager_id,omitempty"`
	DepartmentID      *uuid.UUID `json:"department_id,omitempty"`
	DepartmentName    *string    `json:"department_name,omitempty"`
	EngagementTitle   *string    `json:"engagement_title,omitempty"`
	EngagementCode    *string    `json:"engagement_code,omitempty"`
	WorkerProfileID   *uuid.UUID `json:"worker_profile_id,omitempty"`
	WorkerDisplayName *string    `json:"worker_display_name,omitempty"`
	WorkerCode        *string    `json:"worker_code,omitempty"`
}

type ProjectMilestoneEvent struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	ProjectID   uuid.UUID       `json:"project_id"`
	MilestoneID uuid.UUID       `json:"milestone_id"`
	EventType   string          `json:"event_type"`
	FromStatus  *string         `json:"from_status,omitempty"`
	ToStatus    *string         `json:"to_status,omitempty"`
	Comment     *string         `json:"comment,omitempty"`
	ActorID     *uuid.UUID      `json:"actor_id,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type ProjectMilestoneInput struct {
	TenantID           uuid.UUID
	ProjectID          uuid.UUID
	EngagementID       *uuid.UUID
	MilestoneCode      *string
	Title              string
	Description        *string
	AcceptanceCriteria *string
	DueDate            *time.Time
	Status             string
	Amount             *float64
	CurrencyCode       string
	PaymentTrigger     json.RawMessage
	SubmittedAt        *time.Time
	SubmittedBy        *uuid.UUID
	AcceptedAt         *time.Time
	AcceptedBy         *uuid.UUID
	RejectedAt         *time.Time
	RejectedBy         *uuid.UUID
	ReviewComment      *string
	Notes              *string
	Metadata           json.RawMessage
}

type ProjectMilestoneFilter struct {
	TenantID     uuid.UUID
	ProjectID    *uuid.UUID
	EngagementID *uuid.UUID
	Status       *string
	Search       *string
}

func NewProjectMilestone(input ProjectMilestoneInput) (*ProjectMilestone, error) {
	if input.TenantID == uuid.Nil || input.ProjectID == uuid.Nil {
		return nil, ErrInvalidMilestone
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidMilestone
	}
	status := normalizeWorkerProfileEnum(input.Status, MilestoneStatusDraft)
	if !containsString(milestoneStatuses(), status) {
		return nil, ErrInvalidMilestoneStatus
	}
	if invalidPositiveFloat(input.Amount) {
		return nil, ErrInvalidMilestoneAmount
	}
	currency := normalizeCurrencyCode(input.CurrencyCode)
	if len(currency) != 3 {
		return nil, ErrInvalidMilestoneCurrency
	}
	paymentTrigger := normalizeWorkerJSONObject(input.PaymentTrigger, "{}")
	if !json.Valid(paymentTrigger) || !jsonObject(paymentTrigger) {
		return nil, ErrInvalidMilestonePaymentTrigger
	}
	metadata := normalizeWorkerJSONObject(input.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidMilestoneMetadata
	}
	now := time.Now().UTC()
	return &ProjectMilestone{
		TenantID:           input.TenantID,
		ProjectID:          input.ProjectID,
		EngagementID:       cleanUUIDOptional(input.EngagementID),
		MilestoneCode:      cleanOptional(input.MilestoneCode),
		Title:              title,
		Description:        cleanOptional(input.Description),
		AcceptanceCriteria: cleanOptional(input.AcceptanceCriteria),
		DueDate:            datePtrUTC(input.DueDate),
		Status:             status,
		Amount:             cleanFloatOptional(input.Amount),
		CurrencyCode:       currency,
		PaymentTrigger:     paymentTrigger,
		SubmittedAt:        input.SubmittedAt,
		SubmittedBy:        cleanUUIDOptional(input.SubmittedBy),
		AcceptedAt:         input.AcceptedAt,
		AcceptedBy:         cleanUUIDOptional(input.AcceptedBy),
		RejectedAt:         input.RejectedAt,
		RejectedBy:         cleanUUIDOptional(input.RejectedBy),
		ReviewComment:      cleanOptional(input.ReviewComment),
		Notes:              cleanOptional(input.Notes),
		Metadata:           metadata,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func ValidateMilestoneStatus(value string) (string, error) {
	status := normalizeWorkerProfileEnum(value, "")
	if !containsString(milestoneStatuses(), status) {
		return "", ErrInvalidMilestoneStatus
	}
	return status, nil
}

func milestoneStatuses() []string {
	return []string{MilestoneStatusDraft, MilestoneStatusOpen, MilestoneStatusSubmitted, MilestoneStatusAccepted, MilestoneStatusRejected, MilestoneStatusCancelled}
}
