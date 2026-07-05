package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCandidateID                = errors.New("candidate_id is required")
	ErrInvalidApplicantUserID            = errors.New("applicant user_id is required")
	ErrInvalidCandidateName              = errors.New("candidate first name or last name is required")
	ErrInvalidCandidateEmail             = errors.New("candidate email is invalid")
	ErrInvalidCandidateExperience        = errors.New("candidate total experience is invalid")
	ErrInvalidCandidateSalary            = errors.New("candidate salary is invalid")
	ErrInvalidCandidateNotice            = errors.New("candidate notice period is invalid")
	ErrInvalidCandidateApplicationID     = errors.New("candidate_application_id is required")
	ErrInvalidCandidateApplicationLink   = errors.New("candidate application requires a candidate or job posting")
	ErrInvalidCandidateApplicationStatus = errors.New("candidate application status is invalid")
)

var candidateEmailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type Candidate struct {
	ID                 uuid.UUID  `json:"id"`
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
	Inactive           bool       `json:"inactive"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatedBy          *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt          time.Time  `json:"updated_at"`
	UpdatedBy          *uuid.UUID `json:"updated_by,omitempty"`
}

type CandidateApplicantAccount struct {
	ID          uuid.UUID      `json:"id"`
	TenantID    uuid.UUID      `json:"tenant_id"`
	CandidateID uuid.UUID      `json:"candidate_id"`
	UserID      uuid.UUID      `json:"user_id"`
	Email       string         `json:"email"`
	Status      string         `json:"status"`
	ConsentAt   *time.Time     `json:"consent_at,omitempty"`
	ConsentIP   *string        `json:"consent_ip,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Inactive    bool           `json:"inactive"`
	CreatedAt   time.Time      `json:"created_at"`
	CreatedBy   *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UpdatedBy   *uuid.UUID     `json:"updated_by,omitempty"`
}

type CandidateInput struct {
	TenantID           uuid.UUID
	Firstname          *string
	Lastname           *string
	Email              *string
	Phone              *string
	DOB                *time.Time
	Gender             *string
	TotalExperience    *float64
	CurrentCompany     *string
	CurrentDesignation *string
	CurrentSalary      *float64
	ExpectedSalary     *float64
	NoticePeriod       *int32
	CurrentLocation    *string
	PreferredLocation  *string
	Source             *string
	ResumeURL          *string
}

type CandidateApplicantAccountInput struct {
	TenantID    uuid.UUID
	CandidateID uuid.UUID
	UserID      uuid.UUID
	Email       string
	Status      string
	ConsentAt   *time.Time
	ConsentIP   *string
	Metadata    map[string]any
}

type CandidateFilter struct {
	TenantID uuid.UUID
	Search   *string
	Source   *string
	Gender   *string
	Limit    int32
	Offset   int32
}

type CandidatePage struct {
	Items      []*Candidate `json:"items"`
	Total      int64        `json:"total"`
	Limit      int32        `json:"limit"`
	Offset     int32        `json:"offset"`
	NextOffset *int32       `json:"next_offset,omitempty"`
}

const (
	CandidateApplicationStatusNew       = "New"
	CandidateApplicationStatusScreening = "Screening"
	CandidateApplicationStatusInterview = "Interview"
	CandidateApplicationStatusOffered   = "Offered"
	CandidateApplicationStatusHired     = "Hired"
	CandidateApplicationStatusRejected  = "Rejected"
	CandidateApplicationStatusWithdrawn = "Withdrawn"
)

type CandidateApplication struct {
	ID                       uuid.UUID  `json:"id"`
	TenantID                 uuid.UUID  `json:"tenant_id"`
	CandidateID              *uuid.UUID `json:"candidate_id,omitempty"`
	CandidateFirstname       *string    `json:"candidate_firstname,omitempty"`
	CandidateLastname        *string    `json:"candidate_lastname,omitempty"`
	CandidateEmail           *string    `json:"candidate_email,omitempty"`
	CandidatePhone           *string    `json:"candidate_phone,omitempty"`
	JobPostingID             *uuid.UUID `json:"job_posting_id,omitempty"`
	JobPostingTitle          *string    `json:"job_posting_title,omitempty"`
	JobPostingCode           *string    `json:"job_posting_code,omitempty"`
	ResumeURL                *string    `json:"resume_url,omitempty"`
	CoverLetter              *string    `json:"cover_letter,omitempty"`
	CurrentCTC               *float64   `json:"current_ctc,omitempty"`
	ExpectedCTC              *float64   `json:"expected_ctc,omitempty"`
	NoticePeriod             *int32     `json:"notice_period,omitempty"`
	ReferredBy               *string    `json:"referred_by,omitempty"`
	Source                   *string    `json:"source,omitempty"`
	SourceDetail             *string    `json:"source_detail,omitempty"`
	Status                   string     `json:"status"`
	Comments                 *string    `json:"comments,omitempty"`
	AppliedAt                time.Time  `json:"applied_at"`
	StatusChangedAt          time.Time  `json:"status_changed_at"`
	StatusChangedBy          *uuid.UUID `json:"status_changed_by,omitempty"`
	RejectionReason          *string    `json:"rejection_reason,omitempty"`
	WithdrawalReason         *string    `json:"withdrawal_reason,omitempty"`
	DuplicateOfApplicationID *uuid.UUID `json:"duplicate_of_application_id,omitempty"`
	DaysInStage              int        `json:"days_in_stage"`
	Inactive                 bool       `json:"inactive"`
	CreatedAt                time.Time  `json:"created_at"`
	CreatedBy                *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt                time.Time  `json:"updated_at"`
	UpdatedBy                *uuid.UUID `json:"updated_by,omitempty"`
}

type CandidateApplicationInput struct {
	TenantID                 uuid.UUID
	CandidateID              *uuid.UUID
	JobPostingID             *uuid.UUID
	ResumeURL                *string
	CoverLetter              *string
	CurrentCTC               *float64
	ExpectedCTC              *float64
	NoticePeriod             *int32
	ReferredBy               *string
	Source                   *string
	SourceDetail             *string
	Status                   *string
	Comments                 *string
	AppliedAt                *time.Time
	RejectionReason          *string
	WithdrawalReason         *string
	DuplicateOfApplicationID *uuid.UUID
}

type CandidateApplicationFilter struct {
	TenantID     uuid.UUID
	Status       *string
	CandidateID  *uuid.UUID
	JobPostingID *uuid.UUID
	Search       *string
	Limit        int32
	Offset       int32
}

type CandidateApplicationPage struct {
	Items      []*CandidateApplication `json:"items"`
	Total      int64                   `json:"total"`
	Limit      int32                   `json:"limit"`
	Offset     int32                   `json:"offset"`
	NextOffset *int32                  `json:"next_offset,omitempty"`
}

type ApplicantPortal struct {
	Account      *CandidateApplicantAccount `json:"account"`
	Candidate    *Candidate                 `json:"candidate,omitempty"`
	Applications []*CandidateApplication    `json:"applications"`
}

type CandidateApplicationEvent struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	ApplicationID uuid.UUID  `json:"application_id"`
	FromStatus    *string    `json:"from_status,omitempty"`
	ToStatus      string     `json:"to_status"`
	Action        string     `json:"action"`
	Reason        *string    `json:"reason,omitempty"`
	Remarks       *string    `json:"remarks,omitempty"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

func NewCandidate(input CandidateInput) (*Candidate, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	firstname := cleanOptional(input.Firstname)
	lastname := cleanOptional(input.Lastname)
	if firstname == nil && lastname == nil {
		return nil, ErrInvalidCandidateName
	}
	email := cleanOptional(input.Email)
	if email != nil {
		normalized := strings.ToLower(*email)
		email = &normalized
		if !candidateEmailPattern.MatchString(*email) {
			return nil, ErrInvalidCandidateEmail
		}
	}
	if input.TotalExperience != nil && *input.TotalExperience < 0 {
		return nil, ErrInvalidCandidateExperience
	}
	if input.CurrentSalary != nil && *input.CurrentSalary < 0 {
		return nil, ErrInvalidCandidateSalary
	}
	if input.ExpectedSalary != nil && *input.ExpectedSalary < 0 {
		return nil, ErrInvalidCandidateSalary
	}
	if input.NoticePeriod != nil && *input.NoticePeriod < 0 {
		return nil, ErrInvalidCandidateNotice
	}
	now := time.Now().UTC()
	return &Candidate{
		TenantID:           input.TenantID,
		Firstname:          firstname,
		Lastname:           lastname,
		Email:              email,
		Phone:              cleanOptional(input.Phone),
		DOB:                cleanTimeOptional(input.DOB),
		Gender:             cleanOptional(input.Gender),
		TotalExperience:    input.TotalExperience,
		CurrentCompany:     cleanOptional(input.CurrentCompany),
		CurrentDesignation: cleanOptional(input.CurrentDesignation),
		CurrentSalary:      input.CurrentSalary,
		ExpectedSalary:     input.ExpectedSalary,
		NoticePeriod:       input.NoticePeriod,
		CurrentLocation:    cleanOptional(input.CurrentLocation),
		PreferredLocation:  cleanOptional(input.PreferredLocation),
		Source:             cleanOptional(input.Source),
		ResumeURL:          cleanOptional(input.ResumeURL),
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func NewCandidateApplicantAccount(input CandidateApplicantAccountInput) (*CandidateApplicantAccount, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.CandidateID == uuid.Nil {
		return nil, ErrInvalidCandidateID
	}
	if input.UserID == uuid.Nil {
		return nil, ErrInvalidApplicantUserID
	}
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || !candidateEmailPattern.MatchString(email) {
		return nil, ErrInvalidCandidateEmail
	}
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "active"
	}
	now := time.Now().UTC()
	return &CandidateApplicantAccount{
		TenantID:    input.TenantID,
		CandidateID: input.CandidateID,
		UserID:      input.UserID,
		Email:       email,
		Status:      status,
		ConsentAt:   cleanTimeOptional(input.ConsentAt),
		ConsentIP:   cleanOptional(input.ConsentIP),
		Metadata:    input.Metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func NewCandidateApplication(input CandidateApplicationInput) (*CandidateApplication, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if cleanUUIDOptional(input.CandidateID) == nil && cleanUUIDOptional(input.JobPostingID) == nil {
		return nil, ErrInvalidCandidateApplicationLink
	}
	if input.CurrentCTC != nil && *input.CurrentCTC < 0 {
		return nil, ErrInvalidCandidateSalary
	}
	if input.ExpectedCTC != nil && *input.ExpectedCTC < 0 {
		return nil, ErrInvalidCandidateSalary
	}
	if input.NoticePeriod != nil && *input.NoticePeriod < 0 {
		return nil, ErrInvalidCandidateNotice
	}
	status, err := ValidateCandidateApplicationStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if status == nil {
		value := CandidateApplicationStatusNew
		status = &value
	}
	now := time.Now().UTC()
	return &CandidateApplication{
		TenantID:                 input.TenantID,
		CandidateID:              cleanUUIDOptional(input.CandidateID),
		JobPostingID:             cleanUUIDOptional(input.JobPostingID),
		ResumeURL:                cleanOptional(input.ResumeURL),
		CoverLetter:              cleanOptional(input.CoverLetter),
		CurrentCTC:               input.CurrentCTC,
		ExpectedCTC:              input.ExpectedCTC,
		NoticePeriod:             input.NoticePeriod,
		ReferredBy:               cleanOptional(input.ReferredBy),
		Source:                   cleanOptional(input.Source),
		SourceDetail:             cleanOptional(input.SourceDetail),
		Status:                   *status,
		Comments:                 cleanOptional(input.Comments),
		AppliedAt:                valueOrNow(cleanTimeOptional(input.AppliedAt), now),
		StatusChangedAt:          now,
		RejectionReason:          cleanOptional(input.RejectionReason),
		WithdrawalReason:         cleanOptional(input.WithdrawalReason),
		DuplicateOfApplicationID: cleanUUIDOptional(input.DuplicateOfApplicationID),
		CreatedAt:                now,
		UpdatedAt:                now,
	}, nil
}

func valueOrNow(value *time.Time, fallback time.Time) time.Time {
	if value == nil {
		return fallback
	}
	return *value
}

func ValidateCandidateApplicationStatus(value *string) (*string, error) {
	status := cleanOptional(value)
	if status == nil {
		return nil, nil
	}
	switch *status {
	case CandidateApplicationStatusNew, CandidateApplicationStatusScreening, CandidateApplicationStatusInterview, CandidateApplicationStatusOffered, CandidateApplicationStatusHired, CandidateApplicationStatusRejected, CandidateApplicationStatusWithdrawn:
		return status, nil
	default:
		return nil, ErrInvalidCandidateApplicationStatus
	}
}
