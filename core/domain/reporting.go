package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrReportCatalogInvalid   = errors.New("report catalog item is invalid")
	ErrReportNotFound         = errors.New("report not found")
	ErrReportSavedViewInvalid = errors.New("report saved view is invalid")
	ErrReportExportJobInvalid = errors.New("report export job is invalid")
	ErrReportScheduleInvalid  = errors.New("report schedule is invalid")
	ErrReportSnapshotInvalid  = errors.New("report snapshot is invalid")
)

const (
	ReportScopeSelf   = "self"
	ReportScopeTeam   = "team"
	ReportScopeTenant = "tenant"
	ReportScopeSystem = "system"

	ReportViewPrivate = "private"
	ReportViewTeam    = "team"
	ReportViewTenant  = "tenant"

	ReportExportCSV  = "csv"
	ReportExportPDF  = "pdf"
	ReportExportXLSX = "xlsx"

	ReportExportQueued    = "queued"
	ReportExportRunning   = "running"
	ReportExportCompleted = "completed"
	ReportExportFailed    = "failed"

	ReportFrequencyDaily   = "daily"
	ReportFrequencyWeekly  = "weekly"
	ReportFrequencyMonthly = "monthly"
)

type ReportCatalogItem struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	ReportCode        string          `json:"report_code"`
	Module            string          `json:"module"`
	Name              string          `json:"name"`
	Description       *string         `json:"description,omitempty"`
	Category          string          `json:"category"`
	Scope             string          `json:"scope"`
	PermissionKey     string          `json:"permission_key"`
	DefaultFilters    json.RawMessage `json:"default_filters,omitempty"`
	SupportedFilters  json.RawMessage `json:"supported_filters,omitempty"`
	OutputColumns     json.RawMessage `json:"output_columns,omitempty"`
	DrilldownContract json.RawMessage `json:"drilldown_contract,omitempty"`
	IsSystem          bool            `json:"is_system"`
	IsActive          bool            `json:"is_active"`
	SortOrder         int32           `json:"sort_order"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedBy         *uuid.UUID      `json:"updated_by,omitempty"`
}

type ReportSavedView struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	ReportID    uuid.UUID       `json:"report_id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Visibility  string          `json:"visibility"`
	Filters     json.RawMessage `json:"filters,omitempty"`
	Columns     json.RawMessage `json:"columns,omitempty"`
	IsFavorite  bool            `json:"is_favorite"`
	OwnerUserID *uuid.UUID      `json:"owner_user_id,omitempty"`
	Inactive    bool            `json:"inactive"`
	CreatedAt   time.Time       `json:"created_at"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at"`
	UpdatedBy   *uuid.UUID      `json:"updated_by,omitempty"`
}

type ReportExportJob struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	ReportID      uuid.UUID       `json:"report_id"`
	SavedViewID   *uuid.UUID      `json:"saved_view_id,omitempty"`
	ExportFormat  string          `json:"export_format"`
	Status        string          `json:"status"`
	Filters       json.RawMessage `json:"filters,omitempty"`
	FileObjectKey *string         `json:"file_object_key,omitempty"`
	ErrorMessage  *string         `json:"error_message,omitempty"`
	RequestedBy   *uuid.UUID      `json:"requested_by,omitempty"`
	StartedAt     *time.Time      `json:"started_at,omitempty"`
	CompletedAt   *time.Time      `json:"completed_at,omitempty"`
	Inactive      bool            `json:"inactive"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     *uuid.UUID      `json:"updated_by,omitempty"`
}

type ReportSchedule struct {
	ID               uuid.UUID       `json:"id"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	ReportID         uuid.UUID       `json:"report_id"`
	SavedViewID      *uuid.UUID      `json:"saved_view_id,omitempty"`
	Name             string          `json:"name"`
	Frequency        string          `json:"frequency"`
	Timezone         string          `json:"timezone"`
	DeliveryChannels json.RawMessage `json:"delivery_channels,omitempty"`
	RecipientUserIDs json.RawMessage `json:"recipient_user_ids,omitempty"`
	RecipientEmails  json.RawMessage `json:"recipient_emails,omitempty"`
	NextRunAt        *time.Time      `json:"next_run_at,omitempty"`
	LastRunAt        *time.Time      `json:"last_run_at,omitempty"`
	IsActive         bool            `json:"is_active"`
	Inactive         bool            `json:"inactive"`
	CreatedAt        time.Time       `json:"created_at"`
	CreatedBy        *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt        time.Time       `json:"updated_at"`
	UpdatedBy        *uuid.UUID      `json:"updated_by,omitempty"`
}

type ReportSnapshot struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	ReportID    uuid.UUID       `json:"report_id"`
	SavedViewID *uuid.UUID      `json:"saved_view_id,omitempty"`
	SnapshotKey string          `json:"snapshot_key"`
	PeriodStart time.Time       `json:"period_start"`
	PeriodEnd   time.Time       `json:"period_end"`
	Filters     json.RawMessage `json:"filters,omitempty"`
	Summary     json.RawMessage `json:"summary,omitempty"`
	RowCount    int32           `json:"row_count"`
	GeneratedAt time.Time       `json:"generated_at"`
	GeneratedBy *uuid.UUID      `json:"generated_by,omitempty"`
	Inactive    bool            `json:"inactive"`
	CreatedAt   time.Time       `json:"created_at"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at"`
	UpdatedBy   *uuid.UUID      `json:"updated_by,omitempty"`
}

func NewReportCatalogItem(item ReportCatalogItem) (*ReportCatalogItem, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.ReportCode) == "" || strings.TrimSpace(item.Module) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrReportCatalogInvalid
	}
	item.ReportCode = strings.TrimSpace(item.ReportCode)
	item.Module = strings.TrimSpace(item.Module)
	item.Name = strings.TrimSpace(item.Name)
	item.Category = firstNonEmpty(item.Category, "General")
	scope, ok := normalizeReportEnum(item.Scope, ReportScopeTenant, ReportScopeSelf, ReportScopeTeam, ReportScopeTenant, ReportScopeSystem)
	if !ok || strings.TrimSpace(item.PermissionKey) == "" {
		return nil, ErrReportCatalogInvalid
	}
	item.Scope = scope
	item.PermissionKey = strings.TrimSpace(item.PermissionKey)
	item.DefaultFilters = normalizeJSONRaw(item.DefaultFilters, "{}")
	item.SupportedFilters = normalizeJSONRaw(item.SupportedFilters, "[]")
	item.OutputColumns = normalizeJSONRaw(item.OutputColumns, "[]")
	item.DrilldownContract = normalizeJSONRaw(item.DrilldownContract, "{}")
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	if item.SortOrder == 0 {
		item.SortOrder = 100
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewReportSavedView(item ReportSavedView) (*ReportSavedView, error) {
	if item.TenantID == uuid.Nil || item.ReportID == uuid.Nil || strings.TrimSpace(item.Name) == "" {
		return nil, ErrReportSavedViewInvalid
	}
	item.Name = strings.TrimSpace(item.Name)
	visibility, ok := normalizeReportEnum(item.Visibility, ReportViewPrivate, ReportViewPrivate, ReportViewTeam, ReportViewTenant)
	if !ok {
		return nil, ErrReportSavedViewInvalid
	}
	item.Visibility = visibility
	item.Filters = normalizeJSONRaw(item.Filters, "{}")
	item.Columns = normalizeJSONRaw(item.Columns, "[]")
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewReportExportJob(item ReportExportJob) (*ReportExportJob, error) {
	if item.TenantID == uuid.Nil || item.ReportID == uuid.Nil {
		return nil, ErrReportExportJobInvalid
	}
	format, ok := normalizeReportEnum(item.ExportFormat, ReportExportCSV, ReportExportCSV, ReportExportPDF, ReportExportXLSX)
	if !ok {
		return nil, ErrReportExportJobInvalid
	}
	status, ok := normalizeReportEnum(item.Status, ReportExportQueued, ReportExportQueued, ReportExportRunning, ReportExportCompleted, ReportExportFailed)
	if !ok {
		return nil, ErrReportExportJobInvalid
	}
	item.ExportFormat = format
	item.Status = status
	item.Filters = normalizeJSONRaw(item.Filters, "{}")
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewReportSchedule(item ReportSchedule) (*ReportSchedule, error) {
	if item.TenantID == uuid.Nil || item.ReportID == uuid.Nil || strings.TrimSpace(item.Name) == "" {
		return nil, ErrReportScheduleInvalid
	}
	frequency, ok := normalizeReportEnum(item.Frequency, ReportFrequencyMonthly, ReportFrequencyDaily, ReportFrequencyWeekly, ReportFrequencyMonthly)
	if !ok {
		return nil, ErrReportScheduleInvalid
	}
	item.Name = strings.TrimSpace(item.Name)
	item.Frequency = frequency
	item.Timezone = firstNonEmpty(item.Timezone, "Asia/Kolkata")
	item.DeliveryChannels = normalizeJSONRaw(item.DeliveryChannels, `["email"]`)
	item.RecipientUserIDs = normalizeJSONRaw(item.RecipientUserIDs, "[]")
	item.RecipientEmails = normalizeJSONRaw(item.RecipientEmails, "[]")
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewReportSnapshot(item ReportSnapshot) (*ReportSnapshot, error) {
	if item.TenantID == uuid.Nil || item.ReportID == uuid.Nil || strings.TrimSpace(item.SnapshotKey) == "" || item.PeriodStart.IsZero() || item.PeriodEnd.IsZero() || item.PeriodEnd.Before(item.PeriodStart) || item.RowCount < 0 {
		return nil, ErrReportSnapshotInvalid
	}
	item.SnapshotKey = strings.TrimSpace(item.SnapshotKey)
	item.Filters = normalizeJSONRaw(item.Filters, "{}")
	item.Summary = normalizeJSONRaw(item.Summary, "{}")
	if item.GeneratedAt.IsZero() {
		item.GeneratedAt = time.Now().UTC()
	}
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func normalizeReportEnum(value string, fallback string, allowed ...string) (string, bool) {
	cleaned := strings.TrimSpace(strings.ToLower(value))
	if cleaned == "" {
		cleaned = fallback
	}
	for _, candidate := range allowed {
		if cleaned == candidate {
			return cleaned, true
		}
	}
	return "", false
}

func normalizeJSONRaw(value json.RawMessage, fallback string) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(fallback)
	}
	return value
}

func firstNonEmpty(value string, fallback string) string {
	if trimmed := strings.TrimSpace(value); trimmed != "" {
		return trimmed
	}
	return fallback
}
