package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmployeeLetterID       = errors.New("employee_letter_id is required")
	ErrInvalidEmployeeLetterType     = errors.New("employee letter type is invalid")
	ErrInvalidEmployeeLetterStatus   = errors.New("employee letter status is invalid")
	ErrInvalidEmployeeLetterTemplate = errors.New("employee letter template is invalid")
	ErrInvalidEmployeeLetterBody     = errors.New("employee letter body is required")
	ErrEmployeeLetterCannotBeChanged = errors.New("only generated employee letters can be changed")
	ErrEmployeeLetterPDFMissing      = errors.New("employee letter PDF renderer or storage is not configured")
	ErrEmployeeLetterSignature       = errors.New("employee letter signature is invalid")
	ErrEmployeeLetterTokenMissing    = errors.New("employee letter signature token is missing")
)

const (
	EmployeeLetterTypeAppointment = "appointment"
	EmployeeLetterTypeExperience  = "experience"
	EmployeeLetterTypeRelieving   = "relieving"

	EmployeeLetterStatusGenerated = "Generated"
	EmployeeLetterStatusApproved  = "Approved"
	EmployeeLetterStatusSent      = "Sent"
	EmployeeLetterStatusSigned    = "Signed"
	EmployeeLetterStatusRevoked   = "Revoked"
)

type EmployeeLetterTemplate struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	LetterType  string     `json:"letter_type"`
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

type EmployeeLetter struct {
	ID                   uuid.UUID  `json:"id"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	EmployeeID           uuid.UUID  `json:"employee_id"`
	UserID               uuid.UUID  `json:"user_id"`
	EmployeeCode         *string    `json:"employee_code,omitempty"`
	EmployeeFirstname    *string    `json:"employee_firstname,omitempty"`
	EmployeeLastname     *string    `json:"employee_lastname,omitempty"`
	EmployeeEmail        *string    `json:"employee_email,omitempty"`
	DepartmentName       *string    `json:"department_name,omitempty"`
	BranchName           *string    `json:"branch_name,omitempty"`
	DesignationName      *string    `json:"designation_name,omitempty"`
	TemplateID           *uuid.UUID `json:"template_id,omitempty"`
	TemplateName         *string    `json:"template_name,omitempty"`
	DocumentTypeID       *uuid.UUID `json:"document_type_id,omitempty"`
	DocumentTypeName     *string    `json:"document_type_name,omitempty"`
	EmployeeDocumentID   *uuid.UUID `json:"employee_document_id,omitempty"`
	LetterType           string     `json:"letter_type"`
	Subject              *string    `json:"subject,omitempty"`
	RenderedHTML         *string    `json:"rendered_html,omitempty"`
	Status               string     `json:"status"`
	IssueDate            *time.Time `json:"issue_date,omitempty"`
	EffectiveDate        *time.Time `json:"effective_date,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	PDFPath              *string    `json:"pdf_path,omitempty"`
	Version              int32      `json:"version"`
	IsLatest             bool       `json:"is_latest"`
	ApprovalRequestedAt  *time.Time `json:"approval_requested_at,omitempty"`
	ApprovedAt           *time.Time `json:"approved_at,omitempty"`
	ApprovedBy           *uuid.UUID `json:"approved_by,omitempty"`
	SentAt               *time.Time `json:"sent_at,omitempty"`
	RevokedAt            *time.Time `json:"revoked_at,omitempty"`
	SignatureToken       *string    `json:"signature_token,omitempty"`
	SignatureRequestedAt *time.Time `json:"signature_requested_at,omitempty"`
	SignatureCompletedAt *time.Time `json:"signature_completed_at,omitempty"`
	SignerName           *string    `json:"signer_name,omitempty"`
	SignerEmail          *string    `json:"signer_email,omitempty"`
	SignerIP             *string    `json:"signer_ip,omitempty"`
	SignerUserAgent      *string    `json:"signer_user_agent,omitempty"`
	SignatureHash        *string    `json:"signature_hash,omitempty"`
	AuditCertificateURL  *string    `json:"audit_certificate_url,omitempty"`
	Inactive             bool       `json:"inactive"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type EmployeeLetterEvent struct {
	ID               uuid.UUID      `json:"id"`
	TenantID         uuid.UUID      `json:"tenant_id"`
	EmployeeLetterID uuid.UUID      `json:"employee_letter_id"`
	FromStatus       *string        `json:"from_status,omitempty"`
	ToStatus         string         `json:"to_status"`
	Action           string         `json:"action"`
	Remarks          *string        `json:"remarks,omitempty"`
	ActorEmail       *string        `json:"actor_email,omitempty"`
	IPAddress        *string        `json:"ip_address,omitempty"`
	UserAgent        *string        `json:"user_agent,omitempty"`
	Metadata         map[string]any `json:"metadata"`
	Inactive         bool           `json:"inactive"`
	CreatedAt        time.Time      `json:"created_at"`
	CreatedBy        *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt        time.Time      `json:"updated_at"`
	UpdatedBy        *uuid.UUID     `json:"updated_by,omitempty"`
}

type EmployeeLetterTemplateInput struct {
	TenantID    uuid.UUID
	LetterType  string
	Name        string
	Description *string
	Subject     *string
	BodyHTML    string
	FooterHTML  *string
	Locale      string
	IsDefault   bool
	IsActive    bool
}

type EmployeeLetterInput struct {
	TenantID           uuid.UUID
	EmployeeID         uuid.UUID
	UserID             uuid.UUID
	TemplateID         *uuid.UUID
	DocumentTypeID     *uuid.UUID
	EmployeeDocumentID *uuid.UUID
	LetterType         string
	Subject            *string
	RenderedHTML       *string
	Status             *string
	IssueDate          *time.Time
	EffectiveDate      *time.Time
	EndDate            *time.Time
	PDFPath            *string
	SignatureRequired  bool
	SignerEmail        *string
}

type EmployeeLetterFilter struct {
	TenantID   uuid.UUID
	EmployeeID *uuid.UUID
	LetterType *string
	Status     *string
	Search     *string
	Limit      int32
	Offset     int32
}

type EmployeeLetterPage struct {
	Items      []*EmployeeLetter `json:"items"`
	Total      int64             `json:"total"`
	Limit      int32             `json:"limit"`
	Offset     int32             `json:"offset"`
	NextOffset *int32            `json:"next_offset,omitempty"`
}

func NewEmployeeLetterTemplate(input EmployeeLetterTemplateInput) (*EmployeeLetterTemplate, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	letterType, err := ValidateEmployeeLetterType(input.LetterType)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidEmployeeLetterTemplate
	}
	body := strings.TrimSpace(input.BodyHTML)
	if body == "" {
		return nil, ErrInvalidEmployeeLetterBody
	}
	locale := strings.TrimSpace(input.Locale)
	if locale == "" {
		locale = "en-IN"
	}
	now := time.Now().UTC()
	return &EmployeeLetterTemplate{TenantID: input.TenantID, LetterType: letterType, Name: name, Description: cleanOptional(input.Description), Subject: cleanOptional(input.Subject), BodyHTML: body, FooterHTML: cleanOptional(input.FooterHTML), Locale: locale, IsDefault: input.IsDefault, IsActive: input.IsActive, CreatedAt: now, UpdatedAt: now}, nil
}

func NewEmployeeLetter(input EmployeeLetterInput) (*EmployeeLetter, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.EmployeeID == uuid.Nil {
		return nil, ErrInvalidEmployeeID
	}
	if input.UserID == uuid.Nil {
		return nil, ErrInvalidEmployeeUserID
	}
	letterType, err := ValidateEmployeeLetterType(input.LetterType)
	if err != nil {
		return nil, err
	}
	status, err := validateOptionalEmployeeLetterStatus(input.Status)
	if err != nil {
		return nil, err
	}
	if status == nil {
		value := EmployeeLetterStatusGenerated
		status = &value
	}
	now := time.Now().UTC()
	return &EmployeeLetter{TenantID: input.TenantID, EmployeeID: input.EmployeeID, UserID: input.UserID, TemplateID: cleanUUIDOptional(input.TemplateID), DocumentTypeID: cleanUUIDOptional(input.DocumentTypeID), EmployeeDocumentID: cleanUUIDOptional(input.EmployeeDocumentID), LetterType: letterType, Subject: cleanOptional(input.Subject), RenderedHTML: cleanOptional(input.RenderedHTML), Status: *status, IssueDate: datePtr(input.IssueDate), EffectiveDate: datePtr(input.EffectiveDate), EndDate: datePtr(input.EndDate), PDFPath: cleanOptional(input.PDFPath), SignerEmail: cleanOptional(input.SignerEmail), CreatedAt: now, UpdatedAt: now}, nil
}

func ValidateEmployeeLetterType(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case EmployeeLetterTypeAppointment, EmployeeLetterTypeExperience, EmployeeLetterTypeRelieving:
		return normalized, nil
	default:
		return "", ErrInvalidEmployeeLetterType
	}
}

func ValidateEmployeeLetterStatus(value string) (string, error) {
	status := strings.TrimSpace(value)
	switch status {
	case EmployeeLetterStatusGenerated, EmployeeLetterStatusApproved, EmployeeLetterStatusSent, EmployeeLetterStatusSigned, EmployeeLetterStatusRevoked:
		return status, nil
	default:
		return "", ErrInvalidEmployeeLetterStatus
	}
}

func validateOptionalEmployeeLetterStatus(value *string) (*string, error) {
	status := cleanOptional(value)
	if status == nil {
		return nil, nil
	}
	normalized, err := ValidateEmployeeLetterStatus(*status)
	if err != nil {
		return nil, err
	}
	return &normalized, nil
}

func NewEmployeeLetterSignatureToken() (string, error) {
	return NewOfferSignatureToken()
}
