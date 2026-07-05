package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type FlexPayrollRepo interface {
	CreateFlexPayRun(ctx context.Context, item *domain.FlexPayRun, actorID *uuid.UUID) (*domain.FlexPayRun, error)
	GetFlexPayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.FlexPayRun, error)
	ListFlexPayRuns(ctx context.Context, filter domain.FlexPayRunFilter) ([]*domain.FlexPayRun, error)
	UpdateFlexPayRunStatus(ctx context.Context, item *domain.FlexPayRun, actorID *uuid.UUID) (*domain.FlexPayRun, error)
	DeleteFlexPayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateContractorInvoice(ctx context.Context, item *domain.ContractorInvoice, actorID *uuid.UUID) (*domain.ContractorInvoice, error)
	UpdateContractorInvoice(ctx context.Context, item *domain.ContractorInvoice, actorID *uuid.UUID) (*domain.ContractorInvoice, error)
	GetContractorInvoice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ContractorInvoice, error)
	ListContractorInvoices(ctx context.Context, filter domain.ContractorInvoiceFilter) ([]*domain.ContractorInvoice, error)
	UpdateContractorInvoiceStatus(ctx context.Context, item *domain.ContractorInvoice, actorID *uuid.UUID) (*domain.ContractorInvoice, error)
	DeleteContractorInvoice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateFlexPayRunItem(ctx context.Context, item *domain.FlexPayRunItem, actorID *uuid.UUID) (*domain.FlexPayRunItem, error)
	ListFlexPayRunItems(ctx context.Context, tenantID uuid.UUID, flexPayRunID uuid.UUID) ([]*domain.FlexPayRunItem, error)
	UpdateFlexPayRunItemsStatusByRun(ctx context.Context, tenantID uuid.UUID, flexPayRunID uuid.UUID, status string, actorID *uuid.UUID) error
	GetFlexPayRunTotals(ctx context.Context, tenantID uuid.UUID, flexPayRunID uuid.UUID) (*domain.FlexPayRunTotals, error)
	CreateFlexPayRunEvent(ctx context.Context, item *domain.FlexPayRunEvent) (*domain.FlexPayRunEvent, error)
	ListFlexPayRunEvents(ctx context.Context, tenantID uuid.UUID, flexPayRunID *uuid.UUID, contractorInvoiceID *uuid.UUID) ([]*domain.FlexPayRunEvent, error)
	ListApprovedWorkLogPaymentCandidates(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.WorkLogPaymentCandidate, error)
	ListAcceptedMilestonePaymentCandidates(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.MilestonePaymentCandidate, error)
}

type FlexPayRunCommand struct {
	ID           uuid.UUID       `json:"id,omitempty"`
	TenantID     uuid.UUID       `json:"tenant_id"`
	RunCode      string          `json:"run_code"`
	Title        string          `json:"title"`
	RunType      string          `json:"run_type"`
	PeriodStart  string          `json:"period_start"`
	PeriodEnd    string          `json:"period_end"`
	PayoutDate   *string         `json:"payout_date,omitempty"`
	CurrencyCode string          `json:"currency_code"`
	SourcePolicy string          `json:"source_policy"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	ActorID      *uuid.UUID      `json:"-"`
}

type FlexPayRunGenerateCommand struct {
	TenantID     uuid.UUID  `json:"tenant_id"`
	FlexPayRunID uuid.UUID  `json:"flex_pay_run_id"`
	TDSSection   *string    `json:"tds_section,omitempty"`
	TDSRate      *float64   `json:"tds_rate,omitempty"`
	GSTRate      *float64   `json:"gst_rate,omitempty"`
	ActorID      *uuid.UUID `json:"-"`
}

type FlexPayRunActionCommand struct {
	TenantID         uuid.UUID  `json:"tenant_id"`
	FlexPayRunID     uuid.UUID  `json:"flex_pay_run_id"`
	Comment          *string    `json:"comment,omitempty"`
	PaymentReference *string    `json:"payment_reference,omitempty"`
	ExportBatchRef   *string    `json:"export_batch_ref,omitempty"`
	ActorID          *uuid.UUID `json:"-"`
}

type ContractorInvoiceCommand struct {
	ID              uuid.UUID       `json:"id,omitempty"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	FlexPayRunID    *uuid.UUID      `json:"flex_pay_run_id,omitempty"`
	WorkerProfileID uuid.UUID       `json:"worker_profile_id"`
	EngagementID    *uuid.UUID      `json:"engagement_id,omitempty"`
	InvoiceNumber   string          `json:"invoice_number"`
	InvoiceDate     string          `json:"invoice_date"`
	DueDate         *string         `json:"due_date,omitempty"`
	VendorName      string          `json:"vendor_name"`
	VendorGSTIN     *string         `json:"vendor_gstin,omitempty"`
	PlaceOfSupply   *string         `json:"place_of_supply,omitempty"`
	ReverseCharge   bool            `json:"reverse_charge"`
	CurrencyCode    string          `json:"currency_code"`
	GrossAmount     float64         `json:"gross_amount"`
	TDSSection      *string         `json:"tds_section,omitempty"`
	TDSRate         float64         `json:"tds_rate"`
	GSTRate         float64         `json:"gst_rate"`
	AttachmentPath  *string         `json:"attachment_path,omitempty"`
	Notes           *string         `json:"notes,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}

type ContractorInvoiceActionCommand struct {
	TenantID         uuid.UUID  `json:"tenant_id"`
	InvoiceID        uuid.UUID  `json:"invoice_id"`
	Comment          *string    `json:"comment,omitempty"`
	PaymentReference *string    `json:"payment_reference,omitempty"`
	ActorID          *uuid.UUID `json:"-"`
}

type FlexPayRunItemCommand struct {
	TenantID            uuid.UUID       `json:"tenant_id"`
	FlexPayRunID        uuid.UUID       `json:"flex_pay_run_id"`
	ContractorInvoiceID *uuid.UUID      `json:"contractor_invoice_id,omitempty"`
	WorkerProfileID     uuid.UUID       `json:"worker_profile_id"`
	EngagementID        *uuid.UUID      `json:"engagement_id,omitempty"`
	SourceType          string          `json:"source_type"`
	Description         string          `json:"description"`
	Quantity            float64         `json:"quantity"`
	RateAmount          float64         `json:"rate_amount"`
	TDSSection          *string         `json:"tds_section,omitempty"`
	TDSRate             float64         `json:"tds_rate"`
	GSTRate             float64         `json:"gst_rate"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	ActorID             *uuid.UUID      `json:"-"`
}
