package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type AgreementRepo interface {
	CreateAgreementTemplate(ctx context.Context, item *domain.AgreementTemplate, actorID *uuid.UUID) (*domain.AgreementTemplate, error)
	ListAgreementTemplates(ctx context.Context, tenantID uuid.UUID, agreementType *string) ([]*domain.AgreementTemplate, error)
	GetDefaultAgreementTemplate(ctx context.Context, tenantID uuid.UUID, agreementType string) (*domain.AgreementTemplate, error)
	GetAgreementTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AgreementTemplate, error)
	UpdateAgreementTemplate(ctx context.Context, item *domain.AgreementTemplate, actorID *uuid.UUID) (*domain.AgreementTemplate, error)
	DeleteAgreementTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateAgreement(ctx context.Context, item *domain.Agreement, actorID *uuid.UUID) (*domain.Agreement, error)
	ListAgreements(ctx context.Context, filter domain.AgreementFilter) ([]*domain.Agreement, error)
	GetAgreement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Agreement, error)
	GetAgreementBySignatureToken(ctx context.Context, token string) (*domain.Agreement, error)
	UpdateAgreementPDF(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, pdfPath *string, actorID *uuid.UUID) (*domain.Agreement, error)
	UpdateAgreementStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.Agreement, error)
	SignAgreement(ctx context.Context, token string, signerName string, signerEmail string, signerIP *string, userAgent *string, signatureHash string) (*domain.Agreement, error)
	DeleteAgreement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateAgreementEvent(ctx context.Context, event *domain.AgreementEvent, actorID *uuid.UUID) (*domain.AgreementEvent, error)
	ListAgreementEvents(ctx context.Context, tenantID uuid.UUID, agreementID uuid.UUID) ([]*domain.AgreementEvent, error)
}

type AgreementPDFRenderer interface {
	RenderAgreementPDF(ctx context.Context, doc AgreementDocument) ([]byte, error)
}

type AgreementStorage interface {
	StoreAgreementPDF(ctx context.Context, input StoreAgreementPDFInput) (string, error)
}

type StoreAgreementPDFInput struct {
	TenantID      uuid.UUID
	AgreementID   uuid.UUID
	AgreementType string
	FileName      string
	ContentType   string
	Content       []byte
}

type AgreementDocument struct {
	Agreement *domain.Agreement `json:"agreement"`
}

type AgreementTemplateCommand struct {
	ID            uuid.UUID       `json:"id,omitempty"`
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
	ActorID       *uuid.UUID      `json:"-"`
}

type AgreementCommand struct {
	ID                uuid.UUID       `json:"id,omitempty"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	AgreementType     string          `json:"agreement_type"`
	Title             string          `json:"title"`
	TemplateID        *uuid.UUID      `json:"template_id,omitempty"`
	WorkerProfileID   *uuid.UUID      `json:"worker_profile_id,omitempty"`
	EngagementID      *uuid.UUID      `json:"engagement_id,omitempty"`
	ProjectID         *uuid.UUID      `json:"project_id,omitempty"`
	Subject           *string         `json:"subject,omitempty"`
	RenderedHTML      *string         `json:"rendered_html,omitempty"`
	IssueDate         *time.Time      `json:"issue_date,omitempty"`
	EffectiveDate     *time.Time      `json:"effective_date,omitempty"`
	EndDate           *time.Time      `json:"end_date,omitempty"`
	SignatureRequired bool            `json:"signature_required"`
	SignerEmail       *string         `json:"signer_email,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	ActorID           *uuid.UUID      `json:"-"`
}

type AgreementStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type AgreementSignatureCommand struct {
	Token       string  `json:"token"`
	SignerName  string  `json:"signer_name"`
	SignerEmail string  `json:"signer_email"`
	IPAddress   *string `json:"-"`
	UserAgent   *string `json:"-"`
}
