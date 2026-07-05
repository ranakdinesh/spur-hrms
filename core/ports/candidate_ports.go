package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type CandidateRepo interface {
	CreateCandidate(ctx context.Context, item *domain.Candidate, actorID *uuid.UUID) (*domain.Candidate, error)
	ListCandidates(ctx context.Context, filter domain.CandidateFilter) ([]*domain.Candidate, error)
	CountCandidates(ctx context.Context, filter domain.CandidateFilter) (int64, error)
	GetCandidate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Candidate, error)
	UpdateCandidate(ctx context.Context, item *domain.Candidate, actorID *uuid.UUID) (*domain.Candidate, error)
	DeleteCandidate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateCandidateApplicantAccount(ctx context.Context, item *domain.CandidateApplicantAccount, actorID *uuid.UUID) (*domain.CandidateApplicantAccount, error)
	GetCandidateApplicantAccountByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.CandidateApplicantAccount, error)
	GetCandidateApplicantAccountByCandidate(ctx context.Context, tenantID uuid.UUID, candidateID uuid.UUID) (*domain.CandidateApplicantAccount, error)
	CreateCandidateApplication(ctx context.Context, item *domain.CandidateApplication, actorID *uuid.UUID) (*domain.CandidateApplication, error)
	ListCandidateApplications(ctx context.Context, filter domain.CandidateApplicationFilter) ([]*domain.CandidateApplication, error)
	CountCandidateApplications(ctx context.Context, filter domain.CandidateApplicationFilter) (int64, error)
	GetCandidateApplication(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CandidateApplication, error)
	UpdateCandidateApplication(ctx context.Context, item *domain.CandidateApplication, actorID *uuid.UUID) (*domain.CandidateApplication, error)
	MoveCandidateApplicationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, comments *string, reason *string, actorID *uuid.UUID) (*domain.CandidateApplication, error)
	DeleteCandidateApplication(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateCandidateApplicationEvent(ctx context.Context, event *domain.CandidateApplicationEvent, actorID *uuid.UUID) (*domain.CandidateApplicationEvent, error)
	ListCandidateApplicationEvents(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID) ([]*domain.CandidateApplicationEvent, error)
	CreateInterviewRound(ctx context.Context, item *domain.InterviewRound, actorID *uuid.UUID) (*domain.InterviewRound, error)
	ListInterviewRounds(ctx context.Context, filter domain.InterviewRoundFilter) ([]*domain.InterviewRound, error)
	CountInterviewRounds(ctx context.Context, filter domain.InterviewRoundFilter) (int64, error)
	GetInterviewRound(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.InterviewRound, error)
	UpdateInterviewRound(ctx context.Context, item *domain.InterviewRound, actorID *uuid.UUID) (*domain.InterviewRound, error)
	UpdateInterviewRoundStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, feedback *string, score *float64, decision *string, completedAt *time.Time, actorID *uuid.UUID) (*domain.InterviewRound, error)
	DeleteInterviewRound(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type CandidateCommand struct {
	ID                 uuid.UUID  `json:"id,omitempty"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	Firstname          *string    `json:"firstname,omitempty"`
	Lastname           *string    `json:"lastname,omitempty"`
	Email              *string    `json:"email,omitempty"`
	Phone              *string    `json:"phone,omitempty"`
	DOB                *time.Time `json:"dob,omitempty"`
	Gender             *string    `json:"gender,omitempty"`
	TotalExperience    *float64   `json:"total_experience,omitempty"`
	CurrentCompany     *string    `json:"current_company,omitempty"`
	CurrentDesignation *string    `json:"current_designation,omitempty"`
	CurrentSalary      *float64   `json:"current_salary,omitempty"`
	ExpectedSalary     *float64   `json:"expected_salary,omitempty"`
	NoticePeriod       *int32     `json:"notice_period,omitempty"`
	CurrentLocation    *string    `json:"current_location,omitempty"`
	PreferredLocation  *string    `json:"preferred_location,omitempty"`
	Source             *string    `json:"source,omitempty"`
	ResumeURL          *string    `json:"resume_url,omitempty"`
	ActorID            *uuid.UUID `json:"-"`
}

type CandidateApplicantAccountCommand struct {
	TenantID    uuid.UUID      `json:"tenant_id"`
	CandidateID uuid.UUID      `json:"candidate_id"`
	UserID      uuid.UUID      `json:"user_id"`
	Email       string         `json:"email"`
	Status      string         `json:"status,omitempty"`
	ConsentAt   *time.Time     `json:"consent_at,omitempty"`
	ConsentIP   *string        `json:"consent_ip,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	ActorID     *uuid.UUID     `json:"-"`
}

type CandidateApplicationCommand struct {
	ID                       uuid.UUID  `json:"id,omitempty"`
	TenantID                 uuid.UUID  `json:"tenant_id"`
	CandidateID              *uuid.UUID `json:"candidate_id,omitempty"`
	JobPostingID             *uuid.UUID `json:"job_posting_id,omitempty"`
	ResumeURL                *string    `json:"resume_url,omitempty"`
	CoverLetter              *string    `json:"cover_letter,omitempty"`
	CurrentCTC               *float64   `json:"current_ctc,omitempty"`
	ExpectedCTC              *float64   `json:"expected_ctc,omitempty"`
	NoticePeriod             *int32     `json:"notice_period,omitempty"`
	ReferredBy               *string    `json:"referred_by,omitempty"`
	Source                   *string    `json:"source,omitempty"`
	SourceDetail             *string    `json:"source_detail,omitempty"`
	Status                   *string    `json:"status,omitempty"`
	Comments                 *string    `json:"comments,omitempty"`
	AppliedAt                *time.Time `json:"applied_at,omitempty"`
	RejectionReason          *string    `json:"rejection_reason,omitempty"`
	WithdrawalReason         *string    `json:"withdrawal_reason,omitempty"`
	DuplicateOfApplicationID *uuid.UUID `json:"duplicate_of_application_id,omitempty"`
	ActorID                  *uuid.UUID `json:"-"`
}

type CandidateApplicationMoveCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Comments *string    `json:"comments,omitempty"`
	Reason   *string    `json:"reason,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type InterviewRoundCommand struct {
	ID                uuid.UUID  `json:"id,omitempty"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	ApplicationID     uuid.UUID  `json:"application_id"`
	RoundName         *string    `json:"round_name,omitempty"`
	RoundNumber       *int32     `json:"round_number,omitempty"`
	ScheduledDate     *time.Time `json:"scheduled_date,omitempty"`
	DurationMinutes   *int32     `json:"duration_minutes,omitempty"`
	InterviewerUserID *uuid.UUID `json:"interviewer_user_id,omitempty"`
	Mode              *string    `json:"mode,omitempty"`
	MeetingLink       *string    `json:"meeting_link,omitempty"`
	Location          *string    `json:"location,omitempty"`
	Status            *string    `json:"status,omitempty"`
	Remarks           *string    `json:"remarks,omitempty"`
	Timezone          *string    `json:"timezone,omitempty"`
	Feedback          *string    `json:"feedback,omitempty"`
	Score             *float64   `json:"score,omitempty"`
	Decision          *string    `json:"decision,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}

type InterviewRoundStatusCommand struct {
	TenantID    uuid.UUID  `json:"tenant_id"`
	ID          uuid.UUID  `json:"id"`
	Status      string     `json:"status"`
	Remarks     *string    `json:"remarks,omitempty"`
	Feedback    *string    `json:"feedback,omitempty"`
	Score       *float64   `json:"score,omitempty"`
	Decision    *string    `json:"decision,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ActorID     *uuid.UUID `json:"-"`
}
