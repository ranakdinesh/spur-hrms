package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type SalaryTemplateRepo interface {
	CreateSalaryTemplate(ctx context.Context, item *domain.SalaryTemplate, actorID *uuid.UUID) (*domain.SalaryTemplate, error)
	ListSalaryTemplates(ctx context.Context, tenantID uuid.UUID, fyID *uuid.UUID) ([]*domain.SalaryTemplate, error)
	GetSalaryTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SalaryTemplate, error)
	GetActiveSalaryTemplate(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) (*domain.SalaryTemplate, error)
	UpdateSalaryTemplate(ctx context.Context, item *domain.SalaryTemplate, actorID *uuid.UUID) (*domain.SalaryTemplate, error)
	DeleteSalaryTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	ActivateSalaryTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.SalaryTemplate, error)
	CreateSalaryTemplateItem(ctx context.Context, item *domain.SalaryTemplateItem, actorID *uuid.UUID) (*domain.SalaryTemplateItem, error)
	ListSalaryTemplateItems(ctx context.Context, tenantID uuid.UUID, templateID uuid.UUID) ([]*domain.SalaryTemplateItem, error)
	GetSalaryTemplateItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SalaryTemplateItem, error)
	UpdateSalaryTemplateItem(ctx context.Context, item *domain.SalaryTemplateItem, actorID *uuid.UUID) (*domain.SalaryTemplateItem, error)
	DeleteSalaryTemplateItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type SalaryTemplateCommand struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	FYID          uuid.UUID  `json:"fy_id"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	TemplateType  string     `json:"template_type"`
	AppliesTo     string     `json:"applies_to"`
	CurrencyCode  string     `json:"currency_code"`
	EffectiveFrom *string    `json:"effective_from,omitempty"`
	EffectiveTo   *string    `json:"effective_to,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	IsActive      bool       `json:"is_active"`
	ActorID       *uuid.UUID `json:"-"`
}

type SalaryTemplateItemCommand struct {
	ID               uuid.UUID  `json:"id,omitempty"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	TemplateID       uuid.UUID  `json:"template_id"`
	ItemType         string     `json:"item_type"`
	Code             string     `json:"code"`
	Name             string     `json:"name"`
	Percentage       *float64   `json:"percentage,omitempty"`
	Amount           *float64   `json:"amount,omitempty"`
	CalculationMode  string     `json:"calculation_mode"`
	CalculationBase  string     `json:"calculation_base"`
	Formula          *string    `json:"formula,omitempty"`
	ContributionSide string     `json:"contribution_side"`
	IsTaxExempt      bool       `json:"is_tax_exempt"`
	IsStatutory      bool       `json:"is_statutory"`
	IsVariable       bool       `json:"is_variable"`
	AffectsGross     bool       `json:"affects_gross"`
	AffectsNet       bool       `json:"affects_net"`
	CapAmount        *float64   `json:"cap_amount,omitempty"`
	MinAmount        *float64   `json:"min_amount,omitempty"`
	MaxAmount        *float64   `json:"max_amount,omitempty"`
	SortOrder        int32      `json:"sort_order"`
	ActorID          *uuid.UUID `json:"-"`
}
