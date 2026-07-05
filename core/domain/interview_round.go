package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInterviewRoundID       = errors.New("interview_round_id is required")
	ErrInvalidInterviewApplicationID = errors.New("application_id is required")
	ErrInvalidInterviewStatus        = errors.New("interview status is invalid")
	ErrInvalidInterviewMode          = errors.New("interview mode is invalid")
	ErrInvalidInterviewDuration      = errors.New("interview duration is invalid")
	ErrInvalidInterviewScore         = errors.New("interview score is invalid")
	ErrInvalidInterviewDecision      = errors.New("interview decision is invalid")
	ErrInvalidInterviewDate          = errors.New("interview date is invalid")
)

const (
	InterviewStatusScheduled   = "Scheduled"
	InterviewStatusRescheduled = "Rescheduled"
	InterviewStatusCompleted   = "Completed"
	InterviewStatusCancelled   = "Cancelled"
	InterviewStatusNoShow      = "NoShow"

	InterviewModePhone      = "Phone"
	InterviewModeVideo      = "Video"
	InterviewModeInPerson   = "InPerson"
	InterviewModePanel      = "Panel"
	InterviewModeAssignment = "Assignment"

	InterviewDecisionStrongHire = "StrongHire"
	InterviewDecisionHire       = "Hire"
	InterviewDecisionHold       = "Hold"
	InterviewDecisionNoHire     = "NoHire"
)

type InterviewRound struct {
	ID                 uuid.UUID  `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	ApplicationID      uuid.UUID  `json:"application_id"`
	ApplicationStatus  *string    `json:"application_status,omitempty"`
	CandidateFirstname *string    `json:"candidate_firstname,omitempty"`
	CandidateLastname  *string    `json:"candidate_lastname,omitempty"`
	CandidateEmail     *string    `json:"candidate_email,omitempty"`
	JobPostingTitle    *string    `json:"job_posting_title,omitempty"`
	JobPostingCode     *string    `json:"job_posting_code,omitempty"`
	RoundName          *string    `json:"round_name,omitempty"`
	RoundNumber        *int32     `json:"round_number,omitempty"`
	ScheduledDate      *time.Time `json:"scheduled_date,omitempty"`
	DurationMinutes    *int32     `json:"duration_minutes,omitempty"`
	InterviewerUserID  *uuid.UUID `json:"interviewer_user_id,omitempty"`
	Mode               *string    `json:"mode,omitempty"`
	MeetingLink        *string    `json:"meeting_link,omitempty"`
	Location           *string    `json:"location,omitempty"`
	Status             string     `json:"status"`
	Remarks            *string    `json:"remarks,omitempty"`
	Timezone           string     `json:"timezone"`
	Feedback           *string    `json:"feedback,omitempty"`
	Score              *float64   `json:"score,omitempty"`
	Decision           *string    `json:"decision,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	Inactive           bool       `json:"inactive"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatedBy          *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt          time.Time  `json:"updated_at"`
	UpdatedBy          *uuid.UUID `json:"updated_by,omitempty"`
}

type InterviewRoundInput struct {
	TenantID          uuid.UUID
	ApplicationID     uuid.UUID
	RoundName         *string
	RoundNumber       *int32
	ScheduledDate     *time.Time
	DurationMinutes   *int32
	InterviewerUserID *uuid.UUID
	Mode              *string
	MeetingLink       *string
	Location          *string
	Status            *string
	Remarks           *string
	Timezone          *string
	Feedback          *string
	Score             *float64
	Decision          *string
	CompletedAt       *time.Time
}

type InterviewRoundFilter struct {
	TenantID          uuid.UUID
	ApplicationID     *uuid.UUID
	Status            *string
	InterviewerUserID *uuid.UUID
	DateFrom          *time.Time
	DateTo            *time.Time
	Search            *string
	Limit             int32
	Offset            int32
}

type InterviewRoundPage struct {
	Items      []*InterviewRound `json:"items"`
	Total      int64             `json:"total"`
	Limit      int32             `json:"limit"`
	Offset     int32             `json:"offset"`
	NextOffset *int32            `json:"next_offset,omitempty"`
}

func NewInterviewRound(input InterviewRoundInput) (*InterviewRound, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.ApplicationID == uuid.Nil {
		return nil, ErrInvalidInterviewApplicationID
	}
	if input.DurationMinutes != nil && *input.DurationMinutes <= 0 {
		return nil, ErrInvalidInterviewDuration
	}
	if input.Score != nil && (*input.Score < 0 || *input.Score > 5) {
		return nil, ErrInvalidInterviewScore
	}
	status, err := ValidateInterviewStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if status == nil {
		value := InterviewStatusScheduled
		status = &value
	}
	mode, err := ValidateInterviewMode(input.Mode)
	if err != nil {
		return nil, err
	}
	decision, err := ValidateInterviewDecision(input.Decision)
	if err != nil {
		return nil, err
	}
	timezone := "UTC"
	if cleaned := cleanOptional(input.Timezone); cleaned != nil {
		timezone = *cleaned
	}
	completedAt := cleanTimeOptional(input.CompletedAt)
	if *status == InterviewStatusCompleted && completedAt == nil {
		now := time.Now().UTC()
		completedAt = &now
	}
	now := time.Now().UTC()
	return &InterviewRound{
		TenantID:          input.TenantID,
		ApplicationID:     input.ApplicationID,
		RoundName:         cleanOptional(input.RoundName),
		RoundNumber:       input.RoundNumber,
		ScheduledDate:     cleanTimeOptional(input.ScheduledDate),
		DurationMinutes:   input.DurationMinutes,
		InterviewerUserID: cleanUUIDOptional(input.InterviewerUserID),
		Mode:              mode,
		MeetingLink:       cleanOptional(input.MeetingLink),
		Location:          cleanOptional(input.Location),
		Status:            *status,
		Remarks:           cleanOptional(input.Remarks),
		Timezone:          timezone,
		Feedback:          cleanOptional(input.Feedback),
		Score:             input.Score,
		Decision:          decision,
		CompletedAt:       completedAt,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func ValidateInterviewStatus(value *string) (*string, error) {
	status := cleanOptional(value)
	if status == nil {
		return nil, nil
	}
	switch *status {
	case InterviewStatusScheduled, InterviewStatusRescheduled, InterviewStatusCompleted, InterviewStatusCancelled, InterviewStatusNoShow:
		return status, nil
	default:
		return nil, ErrInvalidInterviewStatus
	}
}

func ValidateInterviewMode(value *string) (*string, error) {
	mode := cleanOptional(value)
	if mode == nil {
		return nil, nil
	}
	switch *mode {
	case InterviewModePhone, InterviewModeVideo, InterviewModeInPerson, InterviewModePanel, InterviewModeAssignment:
		return mode, nil
	default:
		return nil, ErrInvalidInterviewMode
	}
}

func ValidateInterviewDecision(value *string) (*string, error) {
	decision := cleanOptional(value)
	if decision == nil {
		return nil, nil
	}
	switch *decision {
	case InterviewDecisionStrongHire, InterviewDecisionHire, InterviewDecisionHold, InterviewDecisionNoHire:
		return decision, nil
	default:
		return nil, ErrInvalidInterviewDecision
	}
}
