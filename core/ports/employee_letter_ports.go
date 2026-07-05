package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type EmployeeLetterRepo interface {
	CreateEmployeeLetterTemplate(ctx context.Context, item *domain.EmployeeLetterTemplate, actorID *uuid.UUID) (*domain.EmployeeLetterTemplate, error)
	ListEmployeeLetterTemplates(ctx context.Context, tenantID uuid.UUID, letterType *string) ([]*domain.EmployeeLetterTemplate, error)
	GetDefaultEmployeeLetterTemplate(ctx context.Context, tenantID uuid.UUID, letterType string) (*domain.EmployeeLetterTemplate, error)
	GetEmployeeLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeLetterTemplate, error)
	UpdateEmployeeLetterTemplate(ctx context.Context, item *domain.EmployeeLetterTemplate, actorID *uuid.UUID) (*domain.EmployeeLetterTemplate, error)
	DeleteEmployeeLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateEmployeeLetter(ctx context.Context, item *domain.EmployeeLetter, actorID *uuid.UUID) (*domain.EmployeeLetter, error)
	ListEmployeeLetters(ctx context.Context, filter domain.EmployeeLetterFilter) ([]*domain.EmployeeLetter, error)
	CountEmployeeLetters(ctx context.Context, filter domain.EmployeeLetterFilter) (int64, error)
	GetEmployeeLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeLetter, error)
	GetEmployeeLetterBySignatureToken(ctx context.Context, token string) (*domain.EmployeeLetter, error)
	UpdateEmployeeLetterPDF(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, pdfPath *string, employeeDocumentID *uuid.UUID, actorID *uuid.UUID) (*domain.EmployeeLetter, error)
	UpdateEmployeeLetterStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.EmployeeLetter, error)
	SignEmployeeLetter(ctx context.Context, token string, signerName string, signerEmail string, signerIP *string, userAgent *string, signatureHash string) (*domain.EmployeeLetter, error)
	DeleteEmployeeLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateEmployeeLetterEvent(ctx context.Context, event *domain.EmployeeLetterEvent, actorID *uuid.UUID) (*domain.EmployeeLetterEvent, error)
	ListEmployeeLetterEvents(ctx context.Context, tenantID uuid.UUID, employeeLetterID uuid.UUID) ([]*domain.EmployeeLetterEvent, error)
}

type EmployeeLetterPDFRenderer interface {
	RenderEmployeeLetterPDF(ctx context.Context, doc EmployeeLetterDocument) ([]byte, error)
}

type EmployeeLetterStorage interface {
	StoreEmployeeLetterPDF(ctx context.Context, input StoreEmployeeLetterPDFInput) (string, error)
}

type StoreEmployeeLetterPDFInput struct {
	TenantID    uuid.UUID
	EmployeeID  uuid.UUID
	LetterID    uuid.UUID
	LetterType  string
	FileName    string
	ContentType string
	Content     []byte
}

type EmployeeLetterDocument struct {
	Letter   *domain.EmployeeLetter `json:"letter"`
	Employee *domain.Employee       `json:"employee,omitempty"`
}

type EmployeeLetterTemplateCommand struct {
	ID          uuid.UUID  `json:"id,omitempty"`
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
	ActorID     *uuid.UUID `json:"-"`
}

type EmployeeLetterCommand struct {
	ID                uuid.UUID  `json:"id,omitempty"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	EmployeeID        uuid.UUID  `json:"employee_id"`
	TemplateID        *uuid.UUID `json:"template_id,omitempty"`
	DocumentTypeID    *uuid.UUID `json:"document_type_id,omitempty"`
	LetterType        string     `json:"letter_type"`
	Subject           *string    `json:"subject,omitempty"`
	RenderedHTML      *string    `json:"rendered_html,omitempty"`
	IssueDate         *time.Time `json:"issue_date,omitempty"`
	EffectiveDate     *time.Time `json:"effective_date,omitempty"`
	EndDate           *time.Time `json:"end_date,omitempty"`
	SignatureRequired bool       `json:"signature_required"`
	SignerEmail       *string    `json:"signer_email,omitempty"`
	LinkDocument      bool       `json:"link_document"`
	ActorID           *uuid.UUID `json:"-"`
}

type EmployeeLetterStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type EmployeeLetterSignatureCommand struct {
	Token       string  `json:"token"`
	SignerName  string  `json:"signer_name"`
	SignerEmail string  `json:"signer_email"`
	IPAddress   *string `json:"-"`
	UserAgent   *string `json:"-"`
}
