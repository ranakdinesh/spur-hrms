package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type ReportingRepo interface {
	UpsertReportCatalogItem(ctx context.Context, item *domain.ReportCatalogItem, actorID *uuid.UUID) (*domain.ReportCatalogItem, error)
	ListReportCatalog(ctx context.Context, tenantID uuid.UUID, module *string, scope *string) ([]*domain.ReportCatalogItem, error)
	GetReportCatalogItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ReportCatalogItem, error)
	UpsertReportSavedView(ctx context.Context, item *domain.ReportSavedView, actorID *uuid.UUID) (*domain.ReportSavedView, error)
	ListReportSavedViews(ctx context.Context, tenantID uuid.UUID, reportID *uuid.UUID) ([]*domain.ReportSavedView, error)
	DeleteReportSavedView(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateReportExportJob(ctx context.Context, item *domain.ReportExportJob, actorID *uuid.UUID) (*domain.ReportExportJob, error)
	ListReportExportJobs(ctx context.Context, tenantID uuid.UUID, reportID *uuid.UUID, status *string, limit int32, offset int32) ([]*domain.ReportExportJob, error)
	UpdateReportExportJobStatus(ctx context.Context, item *domain.ReportExportJob, actorID *uuid.UUID) (*domain.ReportExportJob, error)
	UpsertReportSchedule(ctx context.Context, item *domain.ReportSchedule, actorID *uuid.UUID) (*domain.ReportSchedule, error)
	ListReportSchedules(ctx context.Context, tenantID uuid.UUID, reportID *uuid.UUID) ([]*domain.ReportSchedule, error)
	DeleteReportSchedule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateReportSnapshot(ctx context.Context, item *domain.ReportSnapshot, actorID *uuid.UUID) (*domain.ReportSnapshot, error)
	ListReportSnapshots(ctx context.Context, tenantID uuid.UUID, reportID *uuid.UUID, limit int32, offset int32) ([]*domain.ReportSnapshot, error)
}

type ReportCatalogQuery struct {
	TenantID uuid.UUID
	Module   *string
	Scope    *string
	ActorID  *uuid.UUID
}

type ReportSavedViewCommand struct {
	ID          uuid.UUID       `json:"id,omitempty"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	ReportID    uuid.UUID       `json:"report_id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Visibility  string          `json:"visibility"`
	Filters     json.RawMessage `json:"filters,omitempty"`
	Columns     json.RawMessage `json:"columns,omitempty"`
	IsFavorite  bool            `json:"is_favorite"`
	OwnerUserID *uuid.UUID      `json:"owner_user_id,omitempty"`
	ActorID     *uuid.UUID      `json:"-"`
}

type ReportListQuery struct {
	TenantID uuid.UUID
	ReportID *uuid.UUID
	Status   *string
	Limit    int32
	Offset   int32
	ActorID  *uuid.UUID
}

type ReportExportJobCommand struct {
	TenantID     uuid.UUID       `json:"tenant_id"`
	ReportID     uuid.UUID       `json:"report_id"`
	SavedViewID  *uuid.UUID      `json:"saved_view_id,omitempty"`
	ExportFormat string          `json:"export_format"`
	Filters      json.RawMessage `json:"filters,omitempty"`
	ActorID      *uuid.UUID      `json:"-"`
}

type ReportExportJobStatusCommand struct {
	TenantID      uuid.UUID  `json:"tenant_id"`
	ID            uuid.UUID  `json:"id"`
	Status        string     `json:"status"`
	FileObjectKey *string    `json:"file_object_key,omitempty"`
	ErrorMessage  *string    `json:"error_message,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}

type ReportScheduleCommand struct {
	ID               uuid.UUID       `json:"id,omitempty"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	ReportID         uuid.UUID       `json:"report_id"`
	SavedViewID      *uuid.UUID      `json:"saved_view_id,omitempty"`
	Name             string          `json:"name"`
	Frequency        string          `json:"frequency"`
	Timezone         string          `json:"timezone"`
	DeliveryChannels json.RawMessage `json:"delivery_channels,omitempty"`
	RecipientUserIDs json.RawMessage `json:"recipient_user_ids,omitempty"`
	RecipientEmails  json.RawMessage `json:"recipient_emails,omitempty"`
	NextRunAt        *string         `json:"next_run_at,omitempty"`
	IsActive         bool            `json:"is_active"`
	ActorID          *uuid.UUID      `json:"-"`
}

type ReportSnapshotCommand struct {
	TenantID    uuid.UUID       `json:"tenant_id"`
	ReportID    uuid.UUID       `json:"report_id"`
	SavedViewID *uuid.UUID      `json:"saved_view_id,omitempty"`
	SnapshotKey string          `json:"snapshot_key"`
	PeriodStart string          `json:"period_start"`
	PeriodEnd   string          `json:"period_end"`
	Filters     json.RawMessage `json:"filters,omitempty"`
	Summary     json.RawMessage `json:"summary,omitempty"`
	RowCount    int32           `json:"row_count"`
	ActorID     *uuid.UUID      `json:"-"`
}

type ReportDatasetQuery struct {
	TenantID   uuid.UUID
	ReportID   uuid.UUID
	ReportCode string
	StartDate  string
	EndDate    string
	Month      int32
	Year       int32
	FYID       *uuid.UUID
	ActorID    *uuid.UUID
}

type ReportDownload struct {
	Content     []byte
	FileName    string
	ContentType string
}
