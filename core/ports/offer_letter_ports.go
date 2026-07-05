package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type OfferLetterRepo interface {
	CreateOfferLetterTemplate(ctx context.Context, item *domain.OfferLetterTemplate, actorID *uuid.UUID) (*domain.OfferLetterTemplate, error)
	ListOfferLetterTemplates(ctx context.Context, tenantID uuid.UUID) ([]*domain.OfferLetterTemplate, error)
	GetDefaultOfferLetterTemplate(ctx context.Context, tenantID uuid.UUID) (*domain.OfferLetterTemplate, error)
	GetOfferLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OfferLetterTemplate, error)
	UpdateOfferLetterTemplate(ctx context.Context, item *domain.OfferLetterTemplate, actorID *uuid.UUID) (*domain.OfferLetterTemplate, error)
	DeleteOfferLetterTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateOfferLetter(ctx context.Context, item *domain.OfferLetter, actorID *uuid.UUID) (*domain.OfferLetter, error)
	ListOfferLetters(ctx context.Context, filter domain.OfferLetterFilter) ([]*domain.OfferLetter, error)
	CountOfferLetters(ctx context.Context, filter domain.OfferLetterFilter) (int64, error)
	ListOfferLettersByApplication(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID) ([]*domain.OfferLetter, error)
	GetLatestOfferLetterByApplication(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID) (*domain.OfferLetter, error)
	GetOfferLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OfferLetter, error)
	GetOfferLetterBySignatureToken(ctx context.Context, token string) (*domain.OfferLetter, error)
	UpdateOfferLetter(ctx context.Context, item *domain.OfferLetter, actorID *uuid.UUID) (*domain.OfferLetter, error)
	UpdateOfferLetterStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, reason *string, actorID *uuid.UUID) (*domain.OfferLetter, error)
	SignOfferLetter(ctx context.Context, token string, signerName string, signerEmail string, signerIP *string, userAgent *string, signatureHash string) (*domain.OfferLetter, error)
	DeleteOfferLetter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateOfferLetterEvent(ctx context.Context, event *domain.OfferLetterEvent, actorID *uuid.UUID) (*domain.OfferLetterEvent, error)
	ListOfferLetterEvents(ctx context.Context, tenantID uuid.UUID, offerLetterID uuid.UUID) ([]*domain.OfferLetterEvent, error)
}

type OfferLetterTemplateCommand struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	TenantID    uuid.UUID  `json:"tenant_id"`
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

type OfferLetterCommand struct {
	ID              uuid.UUID      `json:"id,omitempty"`
	TenantID        uuid.UUID      `json:"tenant_id"`
	ApplicationID   uuid.UUID      `json:"application_id"`
	TemplateID      *uuid.UUID     `json:"template_id,omitempty"`
	OfferedCTC      *float64       `json:"offered_ctc,omitempty"`
	Currency        string         `json:"currency"`
	SalaryBreakdown map[string]any `json:"salary_breakdown,omitempty"`
	JoiningDate     *time.Time     `json:"joining_date,omitempty"`
	ValidUntilDate  *time.Time     `json:"valid_until_date,omitempty"`
	Status          *string        `json:"status,omitempty"`
	OfferLetterURL  *string        `json:"offer_letter_url,omitempty"`
	Subject         *string        `json:"subject,omitempty"`
	RenderedHTML    *string        `json:"rendered_html,omitempty"`
	SignerEmail     *string        `json:"signer_email,omitempty"`
	ActorID         *uuid.UUID     `json:"-"`
}

type OfferLetterStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Reason   *string    `json:"reason,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type OfferLetterSignatureCommand struct {
	Token       string  `json:"token"`
	SignerName  string  `json:"signer_name"`
	SignerEmail string  `json:"signer_email"`
	IPAddress   *string `json:"-"`
	UserAgent   *string `json:"-"`
}
