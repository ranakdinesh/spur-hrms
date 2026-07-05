package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	AgreementTypeSOW               = "sow"
	AgreementTypeNDA               = "nda"
	AgreementTypeRetainer          = "retainer"
	AgreementTypeFreelanceContract = "freelance_contract"
	AgreementTypeInternshipLetter  = "internship_letter"
	AgreementTypeAmendment         = "amendment"

	AgreementStatusGenerated = "Generated"
	AgreementStatusApproved  = "Approved"
	AgreementStatusSent      = "Sent"
	AgreementStatusSigned    = "Signed"
	AgreementStatusRevoked   = "Revoked"
)

var (
	ErrInvalidAgreementID       = errors.New("agreement_id is required")
	ErrInvalidAgreementType     = errors.New("agreement type is invalid")
	ErrInvalidAgreementStatus   = errors.New("agreement status is invalid")
	ErrInvalidAgreementTemplate = errors.New("agreement template is invalid")
	ErrInvalidAgreementBody     = errors.New("agreement body is required")
	ErrInvalidAgreement         = errors.New("agreement is invalid")
	ErrInvalidAgreementDates    = errors.New("agreement dates are invalid")
	ErrInvalidAgreementMetadata = errors.New("agreement metadata must be a valid JSON object")
	ErrAgreementCannotBeChanged = errors.New("agreement lifecycle transition is invalid")
	ErrAgreementPDFMissing      = errors.New("agreement PDF renderer or storage is not configured")
	ErrAgreementSignature       = errors.New("agreement signature is invalid")
	ErrAgreementTokenMissing    = errors.New("agreement signature token is missing")
)

type AgreementTemplate struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	AgreementType string          `json:"agreement_type"`
	Name          string          `json:"name"`
	Description   *string         `json:"description,omitempty"`
	Subject       *string         `json:"subject,omitempty"`
	BodyHTML      string          `json:"body_html"`
	FooterHTML    *string         `json:"footer_html,omitempty"`
	Locale        string          `json:"locale"`
	IsDefault     bool            `json:"is_default"`
	IsActive      bool            `json:"is_active"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	Inactive      bool            `json:"inactive"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     *uuid.UUID      `json:"updated_by,omitempty"`
}

type Agreement struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	AgreementType        string          `json:"agreement_type"`
	Title                string          `json:"title"`
	TemplateID           *uuid.UUID      `json:"template_id,omitempty"`
	TemplateName         *string         `json:"template_name,omitempty"`
	WorkerProfileID      *uuid.UUID      `json:"worker_profile_id,omitempty"`
	WorkerDisplayName    *string         `json:"worker_display_name,omitempty"`
	WorkerCode           *string         `json:"worker_code,omitempty"`
	EngagementID         *uuid.UUID      `json:"engagement_id,omitempty"`
	EngagementTitle      *string         `json:"engagement_title,omitempty"`
	EngagementCode       *string         `json:"engagement_code,omitempty"`
	ProjectID            *uuid.UUID      `json:"project_id,omitempty"`
	ProjectName          *string         `json:"project_name,omitempty"`
	ProjectCode          *string         `json:"project_code,omitempty"`
	Subject              *string         `json:"subject,omitempty"`
	RenderedHTML         *string         `json:"rendered_html,omitempty"`
	Status               string          `json:"status"`
	IssueDate            *time.Time      `json:"issue_date,omitempty"`
	EffectiveDate        *time.Time      `json:"effective_date,omitempty"`
	EndDate              *time.Time      `json:"end_date,omitempty"`
	PDFPath              *string         `json:"pdf_path,omitempty"`
	Version              int32           `json:"version"`
	IsLatest             bool            `json:"is_latest"`
	SentAt               *time.Time      `json:"sent_at,omitempty"`
	RevokedAt            *time.Time      `json:"revoked_at,omitempty"`
	SignatureToken       *string         `json:"signature_token,omitempty"`
	SignatureRequestedAt *time.Time      `json:"signature_requested_at,omitempty"`
	SignatureCompletedAt *time.Time      `json:"signature_completed_at,omitempty"`
	SignerName           *string         `json:"signer_name,omitempty"`
	SignerEmail          *string         `json:"signer_email,omitempty"`
	SignerIP             *string         `json:"signer_ip,omitempty"`
	SignerUserAgent      *string         `json:"signer_user_agent,omitempty"`
	SignatureHash        *string         `json:"signature_hash,omitempty"`
	AuditCertificateURL  *string         `json:"audit_certificate_url,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type AgreementEvent struct {
	ID          uuid.UUID      `json:"id"`
	TenantID    uuid.UUID      `json:"tenant_id"`
	AgreementID uuid.UUID      `json:"agreement_id"`
	FromStatus  *string        `json:"from_status,omitempty"`
	ToStatus    string         `json:"to_status"`
	Action      string         `json:"action"`
	Remarks     *string        `json:"remarks,omitempty"`
	ActorEmail  *string        `json:"actor_email,omitempty"`
	IPAddress   *string        `json:"ip_address,omitempty"`
	UserAgent   *string        `json:"user_agent,omitempty"`
	Metadata    map[string]any `json:"metadata"`
	Inactive    bool           `json:"inactive"`
	CreatedAt   time.Time      `json:"created_at"`
	CreatedBy   *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UpdatedBy   *uuid.UUID     `json:"updated_by,omitempty"`
}

type AgreementTemplateInput struct {
	TenantID      uuid.UUID
	AgreementType string
	Name          string
	Description   *string
	Subject       *string
	BodyHTML      string
	FooterHTML    *string
	Locale        string
	IsDefault     bool
	IsActive      bool
	Metadata      json.RawMessage
}

type AgreementInput struct {
	TenantID          uuid.UUID
	AgreementType     string
	Title             string
	TemplateID        *uuid.UUID
	WorkerProfileID   *uuid.UUID
	EngagementID      *uuid.UUID
	ProjectID         *uuid.UUID
	Subject           *string
	RenderedHTML      *string
	Status            *string
	IssueDate         *time.Time
	EffectiveDate     *time.Time
	EndDate           *time.Time
	PDFPath           *string
	SignatureRequired bool
	SignerEmail       *string
	Metadata          json.RawMessage
}

type AgreementFilter struct {
	TenantID        uuid.UUID
	AgreementType   *string
	Status          *string
	WorkerProfileID *uuid.UUID
	EngagementID    *uuid.UUID
	ProjectID       *uuid.UUID
	Search          *string
}

func NewAgreementTemplate(input AgreementTemplateInput) (*AgreementTemplate, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	agreementType, err := ValidateAgreementType(input.AgreementType)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidAgreementTemplate
	}
	body := strings.TrimSpace(input.BodyHTML)
	if body == "" {
		return nil, ErrInvalidAgreementBody
	}
	locale := strings.TrimSpace(input.Locale)
	if locale == "" {
		locale = "en-IN"
	}
	metadata := normalizeWorkerJSONObject(input.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidAgreementMetadata
	}
	now := time.Now().UTC()
	return &AgreementTemplate{TenantID: input.TenantID, AgreementType: agreementType, Name: name, Description: cleanOptional(input.Description), Subject: cleanOptional(input.Subject), BodyHTML: body, FooterHTML: cleanOptional(input.FooterHTML), Locale: locale, IsDefault: input.IsDefault, IsActive: input.IsActive, Metadata: metadata, CreatedAt: now, UpdatedAt: now}, nil
}

func NewAgreement(input AgreementInput) (*Agreement, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	agreementType, err := ValidateAgreementType(input.AgreementType)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidAgreement
	}
	status, err := validateOptionalAgreementStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if status == nil {
		value := AgreementStatusGenerated
		status = &value
	}
	effectiveDate := datePtrUTC(input.EffectiveDate)
	endDate := datePtrUTC(input.EndDate)
	if effectiveDate != nil && endDate != nil && endDate.Before(*effectiveDate) {
		return nil, ErrInvalidAgreementDates
	}
	metadata := normalizeWorkerJSONObject(input.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidAgreementMetadata
	}
	now := time.Now().UTC()
	return &Agreement{TenantID: input.TenantID, AgreementType: agreementType, Title: title, TemplateID: cleanUUIDOptional(input.TemplateID), WorkerProfileID: cleanUUIDOptional(input.WorkerProfileID), EngagementID: cleanUUIDOptional(input.EngagementID), ProjectID: cleanUUIDOptional(input.ProjectID), Subject: cleanOptional(input.Subject), RenderedHTML: cleanOptional(input.RenderedHTML), Status: *status, IssueDate: datePtrUTC(input.IssueDate), EffectiveDate: effectiveDate, EndDate: endDate, PDFPath: cleanOptional(input.PDFPath), SignerEmail: cleanOptional(input.SignerEmail), Metadata: metadata, CreatedAt: now, UpdatedAt: now}, nil
}

func ValidateAgreementType(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case AgreementTypeSOW, AgreementTypeNDA, AgreementTypeRetainer, AgreementTypeFreelanceContract, AgreementTypeInternshipLetter, AgreementTypeAmendment:
		return normalized, nil
	default:
		return "", ErrInvalidAgreementType
	}
}

func ValidateAgreementStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case AgreementStatusGenerated, AgreementStatusApproved, AgreementStatusSent, AgreementStatusSigned, AgreementStatusRevoked:
		return status, nil
	default:
		return "", ErrInvalidAgreementStatus
	}
}

func NewAgreementSignatureToken() (string, error) {
	return NewOfferSignatureToken()
}

func validateOptionalAgreementStatus(value *string) (*string, error) {
	status := cleanOptional(value)
	if status == nil {
		return nil, nil
	}
	normalized, err := ValidateAgreementStatus(*status)
	if err != nil {
		return nil, err
	}
	return &normalized, nil
}
