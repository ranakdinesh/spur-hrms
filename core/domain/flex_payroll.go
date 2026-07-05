package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidFlexPayRun         = errors.New("flexible pay run is invalid")
	ErrFlexPayRunNotFound        = errors.New("flexible pay run not found")
	ErrFlexPayRunLocked          = errors.New("flexible pay run cannot be changed in its current status")
	ErrInvalidContractorInvoice  = errors.New("contractor invoice is invalid")
	ErrContractorInvoiceNotFound = errors.New("contractor invoice not found")
	ErrInvalidFlexPayRunItem     = errors.New("flexible pay run item is invalid")
)

const (
	FlexPayRunHourly    = "hourly"
	FlexPayRunMilestone = "milestone"
	FlexPayRunRetainer  = "retainer"
	FlexPayRunStipend   = "stipend"
	FlexPayRunInvoice   = "invoice"
	FlexPayRunMixed     = "mixed"

	FlexPayStatusDraft          = "draft"
	FlexPayStatusGenerated      = "generated"
	FlexPayStatusSubmitted      = "submitted"
	FlexPayStatusApproved       = "approved"
	FlexPayStatusRejected       = "rejected"
	FlexPayStatusPaymentPending = "payment_pending"
	FlexPayStatusPaid           = "paid"
	FlexPayStatusCancelled      = "cancelled"

	FlexPaySourceWorkLog       = "work_log"
	FlexPaySourceMilestone     = "milestone"
	FlexPaySourceRetainer      = "retainer"
	FlexPaySourceStipend       = "stipend"
	FlexPaySourceManualInvoice = "manual_invoice"
	FlexPaySourceAdjustment    = "adjustment"

	FlexPayTDS194C = "194C"
	FlexPayTDS194J = "194J"
	FlexPayTDSNone = "none"
)

type FlexPayRun struct {
	ID               uuid.UUID            `json:"id"`
	TenantID         uuid.UUID            `json:"tenant_id"`
	RunCode          string               `json:"run_code"`
	Title            string               `json:"title"`
	RunType          string               `json:"run_type"`
	Status           string               `json:"status"`
	PeriodStart      time.Time            `json:"period_start"`
	PeriodEnd        time.Time            `json:"period_end"`
	PayoutDate       *time.Time           `json:"payout_date,omitempty"`
	CurrencyCode     string               `json:"currency_code"`
	SourcePolicy     string               `json:"source_policy"`
	InvoiceCount     int32                `json:"invoice_count"`
	ItemCount        int32                `json:"item_count"`
	GrossAmount      float64              `json:"gross_amount"`
	TDSAmount        float64              `json:"tds_amount"`
	GSTAmount        float64              `json:"gst_amount"`
	NetAmount        float64              `json:"net_amount"`
	GeneratedAt      *time.Time           `json:"generated_at,omitempty"`
	SubmittedAt      *time.Time           `json:"submitted_at,omitempty"`
	SubmittedBy      *uuid.UUID           `json:"submitted_by,omitempty"`
	ApprovedAt       *time.Time           `json:"approved_at,omitempty"`
	ApprovedBy       *uuid.UUID           `json:"approved_by,omitempty"`
	RejectedAt       *time.Time           `json:"rejected_at,omitempty"`
	RejectedBy       *uuid.UUID           `json:"rejected_by,omitempty"`
	PaidAt           *time.Time           `json:"paid_at,omitempty"`
	PaidBy           *uuid.UUID           `json:"paid_by,omitempty"`
	PaymentReference *string              `json:"payment_reference,omitempty"`
	ExportBatchRef   *string              `json:"export_batch_ref,omitempty"`
	Notes            *string              `json:"notes,omitempty"`
	Metadata         json.RawMessage      `json:"metadata,omitempty"`
	Inactive         bool                 `json:"inactive"`
	CreatedAt        time.Time            `json:"created_at"`
	CreatedBy        *uuid.UUID           `json:"created_by,omitempty"`
	UpdatedAt        time.Time            `json:"updated_at"`
	UpdatedBy        *uuid.UUID           `json:"updated_by,omitempty"`
	Items            []*FlexPayRunItem    `json:"items,omitempty"`
	Invoices         []*ContractorInvoice `json:"invoices,omitempty"`
	Events           []*FlexPayRunEvent   `json:"events,omitempty"`
}

type ContractorInvoice struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	FlexPayRunID      *uuid.UUID      `json:"flex_pay_run_id,omitempty"`
	WorkerProfileID   uuid.UUID       `json:"worker_profile_id"`
	EngagementID      *uuid.UUID      `json:"engagement_id,omitempty"`
	InvoiceNumber     string          `json:"invoice_number"`
	InvoiceDate       time.Time       `json:"invoice_date"`
	DueDate           *time.Time      `json:"due_date,omitempty"`
	Status            string          `json:"status"`
	VendorName        string          `json:"vendor_name"`
	VendorGSTIN       *string         `json:"vendor_gstin,omitempty"`
	PlaceOfSupply     *string         `json:"place_of_supply,omitempty"`
	ReverseCharge     bool            `json:"reverse_charge"`
	CurrencyCode      string          `json:"currency_code"`
	GrossAmount       float64         `json:"gross_amount"`
	TDSSection        *string         `json:"tds_section,omitempty"`
	TDSRate           float64         `json:"tds_rate"`
	TDSAmount         float64         `json:"tds_amount"`
	GSTRate           float64         `json:"gst_rate"`
	GSTAmount         float64         `json:"gst_amount"`
	NetAmount         float64         `json:"net_amount"`
	SubmittedAt       *time.Time      `json:"submitted_at,omitempty"`
	SubmittedBy       *uuid.UUID      `json:"submitted_by,omitempty"`
	ApprovedAt        *time.Time      `json:"approved_at,omitempty"`
	ApprovedBy        *uuid.UUID      `json:"approved_by,omitempty"`
	RejectedAt        *time.Time      `json:"rejected_at,omitempty"`
	RejectedBy        *uuid.UUID      `json:"rejected_by,omitempty"`
	RejectionReason   *string         `json:"rejection_reason,omitempty"`
	PaidAt            *time.Time      `json:"paid_at,omitempty"`
	PaidBy            *uuid.UUID      `json:"paid_by,omitempty"`
	PaymentReference  *string         `json:"payment_reference,omitempty"`
	AttachmentPath    *string         `json:"attachment_path,omitempty"`
	Notes             *string         `json:"notes,omitempty"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedBy         *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName *string         `json:"worker_display_name,omitempty"`
	WorkerCode        *string         `json:"worker_code,omitempty"`
	EngagementTitle   *string         `json:"engagement_title,omitempty"`
	EngagementCode    *string         `json:"engagement_code,omitempty"`
}

type FlexPayRunItem struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	FlexPayRunID        uuid.UUID       `json:"flex_pay_run_id"`
	ContractorInvoiceID *uuid.UUID      `json:"contractor_invoice_id,omitempty"`
	WorkerProfileID     uuid.UUID       `json:"worker_profile_id"`
	EngagementID        *uuid.UUID      `json:"engagement_id,omitempty"`
	SourceType          string          `json:"source_type"`
	SourceID            *uuid.UUID      `json:"source_id,omitempty"`
	Description         string          `json:"description"`
	Quantity            float64         `json:"quantity"`
	RateAmount          float64         `json:"rate_amount"`
	GrossAmount         float64         `json:"gross_amount"`
	TDSSection          *string         `json:"tds_section,omitempty"`
	TDSRate             float64         `json:"tds_rate"`
	TDSAmount           float64         `json:"tds_amount"`
	GSTRate             float64         `json:"gst_rate"`
	GSTAmount           float64         `json:"gst_amount"`
	NetAmount           float64         `json:"net_amount"`
	Status              string          `json:"status"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerDisplayName   *string         `json:"worker_display_name,omitempty"`
	WorkerCode          *string         `json:"worker_code,omitempty"`
	EngagementTitle     *string         `json:"engagement_title,omitempty"`
	EngagementCode      *string         `json:"engagement_code,omitempty"`
	InvoiceNumber       *string         `json:"invoice_number,omitempty"`
}

type FlexPayRunEvent struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	FlexPayRunID        *uuid.UUID      `json:"flex_pay_run_id,omitempty"`
	ContractorInvoiceID *uuid.UUID      `json:"contractor_invoice_id,omitempty"`
	EventType           string          `json:"event_type"`
	FromStatus          *string         `json:"from_status,omitempty"`
	ToStatus            *string         `json:"to_status,omitempty"`
	Comment             *string         `json:"comment,omitempty"`
	ActorID             *uuid.UUID      `json:"actor_id,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
}

type FlexPayRunFilter struct {
	TenantID uuid.UUID
	Status   *string
	RunType  *string
	DateFrom *time.Time
	DateTo   *time.Time
	Search   *string
}

type ContractorInvoiceFilter struct {
	TenantID        uuid.UUID
	FlexPayRunID    *uuid.UUID
	WorkerProfileID *uuid.UUID
	Status          *string
	Search          *string
}

type FlexPayRunTotals struct {
	InvoiceCount int32   `json:"invoice_count"`
	ItemCount    int32   `json:"item_count"`
	GrossAmount  float64 `json:"gross_amount"`
	TDSAmount    float64 `json:"tds_amount"`
	GSTAmount    float64 `json:"gst_amount"`
	NetAmount    float64 `json:"net_amount"`
}

type WorkLogPaymentCandidate struct {
	WorkLogID         uuid.UUID `json:"work_log_id"`
	TenantID          uuid.UUID `json:"tenant_id"`
	EngagementID      uuid.UUID `json:"engagement_id"`
	WorkerProfileID   uuid.UUID `json:"worker_profile_id"`
	LogDate           time.Time `json:"log_date"`
	BillableHours     float64   `json:"billable_hours"`
	RateAmount        float64   `json:"rate_amount"`
	RateUnit          string    `json:"rate_unit"`
	CurrencyCode      string    `json:"currency_code"`
	EngagementTitle   string    `json:"engagement_title"`
	EngagementCode    *string   `json:"engagement_code,omitempty"`
	WorkerDisplayName string    `json:"worker_display_name"`
	WorkerCode        *string   `json:"worker_code,omitempty"`
}

type MilestonePaymentCandidate struct {
	MilestoneID       uuid.UUID  `json:"milestone_id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	ProjectID         uuid.UUID  `json:"project_id"`
	EngagementID      *uuid.UUID `json:"engagement_id,omitempty"`
	WorkerProfileID   uuid.UUID  `json:"worker_profile_id"`
	Title             string     `json:"title"`
	MilestoneCode     *string    `json:"milestone_code,omitempty"`
	Amount            float64    `json:"amount"`
	CurrencyCode      string     `json:"currency_code"`
	AcceptedAt        *time.Time `json:"accepted_at,omitempty"`
	ProjectName       string     `json:"project_name"`
	ProjectCode       *string    `json:"project_code,omitempty"`
	WorkerDisplayName string     `json:"worker_display_name"`
	WorkerCode        *string    `json:"worker_code,omitempty"`
}

func NewFlexPayRun(item FlexPayRun) (*FlexPayRun, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.RunCode) == "" || strings.TrimSpace(item.Title) == "" || item.PeriodStart.IsZero() || item.PeriodEnd.IsZero() || item.PeriodEnd.Before(item.PeriodStart) {
		return nil, ErrInvalidFlexPayRun
	}
	runType := normalizeWorkerProfileEnum(item.RunType, FlexPayRunMixed)
	if !containsString([]string{FlexPayRunHourly, FlexPayRunMilestone, FlexPayRunRetainer, FlexPayRunStipend, FlexPayRunInvoice, FlexPayRunMixed}, runType) {
		return nil, ErrInvalidFlexPayRun
	}
	status, ok := ValidateFlexPayStatus(item.Status)
	if !ok {
		return nil, ErrInvalidFlexPayRun
	}
	currency := normalizeCurrencyCode(item.CurrencyCode)
	if len(currency) != 3 {
		return nil, ErrInvalidFlexPayRun
	}
	if negative(item.GrossAmount, item.TDSAmount, item.GSTAmount, item.NetAmount) {
		return nil, ErrInvalidFlexPayRun
	}
	metadata := normalizeWorkerJSONObject(item.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidFlexPayRun
	}
	item.RunCode = strings.TrimSpace(item.RunCode)
	item.Title = strings.TrimSpace(item.Title)
	item.RunType = runType
	item.Status = status
	item.CurrencyCode = currency
	item.SourcePolicy = strings.TrimSpace(item.SourcePolicy)
	if item.SourcePolicy == "" {
		item.SourcePolicy = "approved_work_logs_and_accepted_milestones"
	}
	item.Notes = cleanOptional(item.Notes)
	item.Metadata = metadata
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewContractorInvoice(item ContractorInvoice) (*ContractorInvoice, error) {
	if item.TenantID == uuid.Nil || item.WorkerProfileID == uuid.Nil || strings.TrimSpace(item.InvoiceNumber) == "" || strings.TrimSpace(item.VendorName) == "" || item.InvoiceDate.IsZero() {
		return nil, ErrInvalidContractorInvoice
	}
	status, ok := ValidateFlexPayStatus(item.Status)
	if !ok {
		return nil, ErrInvalidContractorInvoice
	}
	currency := normalizeCurrencyCode(item.CurrencyCode)
	if len(currency) != 3 || negative(item.GrossAmount, item.TDSRate, item.TDSAmount, item.GSTRate, item.GSTAmount, item.NetAmount) {
		return nil, ErrInvalidContractorInvoice
	}
	if item.TDSSection != nil {
		section := strings.TrimSpace(*item.TDSSection)
		if section == "" {
			item.TDSSection = nil
		} else if !containsString([]string{FlexPayTDS194C, FlexPayTDS194J, FlexPayTDSNone}, section) {
			return nil, ErrInvalidContractorInvoice
		} else {
			item.TDSSection = &section
		}
	}
	metadata := normalizeWorkerJSONObject(item.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidContractorInvoice
	}
	item.FlexPayRunID = cleanUUIDOptional(item.FlexPayRunID)
	item.EngagementID = cleanUUIDOptional(item.EngagementID)
	item.InvoiceNumber = strings.TrimSpace(item.InvoiceNumber)
	item.VendorName = strings.TrimSpace(item.VendorName)
	item.CurrencyCode = currency
	item.Status = status
	item.VendorGSTIN = cleanOptional(item.VendorGSTIN)
	item.PlaceOfSupply = cleanOptional(item.PlaceOfSupply)
	item.AttachmentPath = cleanOptional(item.AttachmentPath)
	item.Notes = cleanOptional(item.Notes)
	item.Metadata = metadata
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewFlexPayRunItem(item FlexPayRunItem) (*FlexPayRunItem, error) {
	if item.TenantID == uuid.Nil || item.FlexPayRunID == uuid.Nil || item.WorkerProfileID == uuid.Nil || strings.TrimSpace(item.Description) == "" {
		return nil, ErrInvalidFlexPayRunItem
	}
	sourceType := normalizeWorkerProfileEnum(item.SourceType, FlexPaySourceManualInvoice)
	if !containsString([]string{FlexPaySourceWorkLog, FlexPaySourceMilestone, FlexPaySourceRetainer, FlexPaySourceStipend, FlexPaySourceManualInvoice, FlexPaySourceAdjustment}, sourceType) {
		return nil, ErrInvalidFlexPayRunItem
	}
	status, ok := ValidateFlexPayStatus(item.Status)
	if !ok || negative(item.Quantity, item.RateAmount, item.GrossAmount, item.TDSRate, item.TDSAmount, item.GSTRate, item.GSTAmount, item.NetAmount) {
		return nil, ErrInvalidFlexPayRunItem
	}
	metadata := normalizeWorkerJSONObject(item.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidFlexPayRunItem
	}
	item.ContractorInvoiceID = cleanUUIDOptional(item.ContractorInvoiceID)
	item.EngagementID = cleanUUIDOptional(item.EngagementID)
	item.SourceID = cleanUUIDOptional(item.SourceID)
	item.SourceType = sourceType
	item.Status = status
	item.Description = strings.TrimSpace(item.Description)
	item.Metadata = metadata
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func ValidateFlexPayStatus(value string) (string, bool) {
	status := normalizeWorkerProfileEnum(value, FlexPayStatusDraft)
	return status, containsString([]string{FlexPayStatusDraft, FlexPayStatusGenerated, FlexPayStatusSubmitted, FlexPayStatusApproved, FlexPayStatusRejected, FlexPayStatusPaymentPending, FlexPayStatusPaid, FlexPayStatusCancelled}, status)
}

func FlexTaxAmounts(gross float64, tdsRate float64, gstRate float64) (float64, float64, float64) {
	tds := roundMoney(gross * tdsRate / 100)
	gst := roundMoney(gross * gstRate / 100)
	return tds, gst, roundMoney(gross - tds + gst)
}

func negative(values ...float64) bool {
	for _, value := range values {
		if value < 0 {
			return true
		}
	}
	return false
}
