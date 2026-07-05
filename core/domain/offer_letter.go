package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidOfferLetterID       = errors.New("offer_letter_id is required")
	ErrInvalidOfferTemplateID     = errors.New("offer_template_id is required")
	ErrInvalidOfferTemplateName   = errors.New("offer template name is required")
	ErrInvalidOfferTemplateBody   = errors.New("offer template body is required")
	ErrInvalidOfferApplicationID  = errors.New("application_id is required")
	ErrInvalidOfferStatus         = errors.New("offer status is invalid")
	ErrInvalidOfferAmount         = errors.New("offer amount is invalid")
	ErrInvalidOfferSignature      = errors.New("offer signature is invalid")
	ErrOfferCannotBeChanged       = errors.New("only generated offers can be edited")
	ErrOfferSignatureTokenMissing = errors.New("offer signature token is missing")
)

type OfferLetterTemplate struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Subject     *string    `json:"subject,omitempty"`
	BodyHTML    string     `json:"body_html"`
	FooterHTML  *string    `json:"footer_html,omitempty"`
	Locale      string     `json:"locale"`
	IsDefault   bool       `json:"is_default"`
	IsActive    bool       `json:"is_active"`
	Inactive    bool       `json:"inactive"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty"`
}

type OfferLetter struct {
	ID                       uuid.UUID      `json:"id"`
	TenantID                 uuid.UUID      `json:"tenant_id"`
	ApplicationID            uuid.UUID      `json:"application_id"`
	CandidateID              *uuid.UUID     `json:"candidate_id,omitempty"`
	CandidateFirstname       *string        `json:"candidate_firstname,omitempty"`
	CandidateLastname        *string        `json:"candidate_lastname,omitempty"`
	CandidateEmail           *string        `json:"candidate_email,omitempty"`
	JobPostingTitle          *string        `json:"job_posting_title,omitempty"`
	JobPostingCode           *string        `json:"job_posting_code,omitempty"`
	TemplateID               *uuid.UUID     `json:"template_id,omitempty"`
	TemplateName             *string        `json:"template_name,omitempty"`
	OfferedCTC               *float64       `json:"offered_ctc,omitempty"`
	Currency                 string         `json:"currency"`
	SalaryBreakdown          map[string]any `json:"salary_breakdown,omitempty"`
	JoiningDate              *time.Time     `json:"joining_date,omitempty"`
	ValidUntilDate           *time.Time     `json:"valid_until_date,omitempty"`
	Status                   string         `json:"status"`
	OfferLetterURL           *string        `json:"offer_letter_url,omitempty"`
	CandidateReactionDate    *time.Time     `json:"candidate_reaction_date,omitempty"`
	CandidateRejectionReason *string        `json:"candidate_rejection_reason,omitempty"`
	Version                  int32          `json:"version"`
	IsLatest                 bool           `json:"is_latest"`
	Subject                  *string        `json:"subject,omitempty"`
	RenderedHTML             *string        `json:"rendered_html,omitempty"`
	SentAt                   *time.Time     `json:"sent_at,omitempty"`
	RevokedAt                *time.Time     `json:"revoked_at,omitempty"`
	SignatureToken           *string        `json:"signature_token,omitempty"`
	SignatureRequestedAt     *time.Time     `json:"signature_requested_at,omitempty"`
	SignatureCompletedAt     *time.Time     `json:"signature_completed_at,omitempty"`
	SignerName               *string        `json:"signer_name,omitempty"`
	SignerEmail              *string        `json:"signer_email,omitempty"`
	SignerIP                 *string        `json:"signer_ip,omitempty"`
	SignerUserAgent          *string        `json:"signer_user_agent,omitempty"`
	SignatureHash            *string        `json:"signature_hash,omitempty"`
	AuditCertificateURL      *string        `json:"audit_certificate_url,omitempty"`
	Inactive                 bool           `json:"inactive"`
	CreatedAt                time.Time      `json:"created_at"`
	CreatedBy                *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt                time.Time      `json:"updated_at"`
	UpdatedBy                *uuid.UUID     `json:"updated_by,omitempty"`
}

type OfferLetterEvent struct {
	ID            uuid.UUID      `json:"id"`
	TenantID      uuid.UUID      `json:"tenant_id"`
	OfferLetterID uuid.UUID      `json:"offer_letter_id"`
	FromStatus    *string        `json:"from_status,omitempty"`
	ToStatus      string         `json:"to_status"`
	Action        string         `json:"action"`
	Remarks       *string        `json:"remarks,omitempty"`
	ActorEmail    *string        `json:"actor_email,omitempty"`
	IPAddress     *string        `json:"ip_address,omitempty"`
	UserAgent     *string        `json:"user_agent,omitempty"`
	Metadata      map[string]any `json:"metadata"`
	Inactive      bool           `json:"inactive"`
	CreatedAt     time.Time      `json:"created_at"`
	CreatedBy     *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt     time.Time      `json:"updated_at"`
	UpdatedBy     *uuid.UUID     `json:"updated_by,omitempty"`
}

type OfferLetterTemplateInput struct {
	TenantID    uuid.UUID
	Name        string
	Description *string
	Subject     *string
	BodyHTML    string
	FooterHTML  *string
	Locale      string
	IsDefault   bool
	IsActive    bool
}

type OfferLetterInput struct {
	TenantID        uuid.UUID
	ApplicationID   uuid.UUID
	CandidateID     *uuid.UUID
	TemplateID      *uuid.UUID
	OfferedCTC      *float64
	Currency        string
	SalaryBreakdown map[string]any
	JoiningDate     *time.Time
	ValidUntilDate  *time.Time
	Status          *string
	OfferLetterURL  *string
	Subject         *string
	RenderedHTML    *string
	SignerEmail     *string
}

type OfferLetterFilter struct {
	TenantID      uuid.UUID
	ApplicationID *uuid.UUID
	Status        *string
	Search        *string
	Limit         int32
	Offset        int32
}

type OfferLetterPage struct {
	Items      []*OfferLetter `json:"items"`
	Total      int64          `json:"total"`
	Limit      int32          `json:"limit"`
	Offset     int32          `json:"offset"`
	NextOffset *int32         `json:"next_offset,omitempty"`
}

func NewOfferLetterTemplate(input OfferLetterTemplateInput) (*OfferLetterTemplate, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidOfferTemplateName
	}
	body := strings.TrimSpace(input.BodyHTML)
	if body == "" {
		return nil, ErrInvalidOfferTemplateBody
	}
	locale := strings.TrimSpace(input.Locale)
	if locale == "" {
		locale = "en-IN"
	}
	now := time.Now().UTC()
	return &OfferLetterTemplate{TenantID: input.TenantID, Name: name, Description: cleanOptional(input.Description), Subject: cleanOptional(input.Subject), BodyHTML: body, FooterHTML: cleanOptional(input.FooterHTML), Locale: locale, IsDefault: input.IsDefault, IsActive: input.IsActive, CreatedAt: now, UpdatedAt: now}, nil
}

func NewOfferLetter(input OfferLetterInput) (*OfferLetter, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.ApplicationID == uuid.Nil {
		return nil, ErrInvalidOfferApplicationID
	}
	if input.OfferedCTC != nil && *input.OfferedCTC < 0 {
		return nil, ErrInvalidOfferAmount
	}
	status, err := validateOptionalOfferStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if status == nil {
		value := OfferStatusGenerated
		status = &value
	}
	currency := strings.ToUpper(strings.TrimSpace(input.Currency))
	if currency == "" {
		currency = "INR"
	}
	breakdown := input.SalaryBreakdown
	if breakdown == nil {
		breakdown = map[string]any{}
	}
	now := time.Now().UTC()
	return &OfferLetter{TenantID: input.TenantID, ApplicationID: input.ApplicationID, CandidateID: cleanUUIDOptional(input.CandidateID), TemplateID: cleanUUIDOptional(input.TemplateID), OfferedCTC: input.OfferedCTC, Currency: currency, SalaryBreakdown: breakdown, JoiningDate: cleanTimeOptional(input.JoiningDate), ValidUntilDate: cleanTimeOptional(input.ValidUntilDate), Status: *status, OfferLetterURL: cleanOptional(input.OfferLetterURL), Subject: cleanOptional(input.Subject), RenderedHTML: cleanOptional(input.RenderedHTML), SignerEmail: cleanOptional(input.SignerEmail), CreatedAt: now, UpdatedAt: now}, nil
}

func validateOptionalOfferStatus(value *string) (*string, error) {
	status := cleanOptional(value)
	if status == nil {
		return nil, nil
	}
	normalized, err := ValidateOfferStatus(*status)
	if err != nil {
		return nil, ErrInvalidOfferStatus
	}
	return &normalized, nil
}

func NewOfferSignatureToken() (string, error) {
	var data [24]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(data[:]), nil
}
