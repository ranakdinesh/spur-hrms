package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) UpsertReportCatalogItem(ctx context.Context, item *domain.ReportCatalogItem, actorID *uuid.UUID) (*domain.ReportCatalogItem, error) {
	row, err := s.getQueries(ctx).UpsertReportCatalogItem(ctx, sqlc.UpsertReportCatalogItemParams{ID: item.ID, TenantID: item.TenantID, ReportCode: item.ReportCode, Module: item.Module, Name: item.Name, Description: textFromPtr(item.Description), Category: item.Category, Scope: item.Scope, PermissionKey: item.PermissionKey, DefaultFilters: []byte(item.DefaultFilters), SupportedFilters: []byte(item.SupportedFilters), OutputColumns: []byte(item.OutputColumns), DrilldownContract: []byte(item.DrilldownContract), IsSystem: item.IsSystem, IsActive: item.IsActive, SortOrder: item.SortOrder, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert report catalog item", err, tenantIDField(item.TenantID), stringField("report_code", item.ReportCode))
	}
	return mapReportCatalogItem(row), nil
}

func (s *Store) ListReportCatalog(ctx context.Context, tenantID uuid.UUID, module *string, scope *string) ([]*domain.ReportCatalogItem, error) {
	rows, err := s.getQueries(ctx).ListReportCatalog(ctx, sqlc.ListReportCatalogParams{TenantID: tenantID, Module: textFromPtr(module), Scope: textFromPtr(scope)})
	if err != nil {
		return nil, s.logDBError(ctx, "list report catalog", err, tenantIDField(tenantID))
	}
	return mapReportCatalogItems(rows), nil
}

func (s *Store) GetReportCatalogItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ReportCatalogItem, error) {
	row, err := s.getQueries(ctx).GetReportCatalogItem(ctx, sqlc.GetReportCatalogItemParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrReportNotFound
		}
		return nil, s.logDBError(ctx, "get report catalog item", err, tenantIDField(tenantID), stringField("report_id", id.String()))
	}
	return mapReportCatalogItem(row), nil
}

func (s *Store) UpsertReportSavedView(ctx context.Context, item *domain.ReportSavedView, actorID *uuid.UUID) (*domain.ReportSavedView, error) {
	row, err := s.getQueries(ctx).UpsertReportSavedView(ctx, sqlc.UpsertReportSavedViewParams{ID: item.ID, TenantID: item.TenantID, ReportID: item.ReportID, Name: item.Name, Description: textFromPtr(item.Description), Visibility: item.Visibility, Filters: []byte(item.Filters), Columns: []byte(item.Columns), IsFavorite: item.IsFavorite, OwnerUserID: uuidFromPtr(item.OwnerUserID), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert report saved view", err, tenantIDField(item.TenantID), stringField("report_id", item.ReportID.String()))
	}
	return mapReportSavedView(row), nil
}

func (s *Store) ListReportSavedViews(ctx context.Context, tenantID uuid.UUID, reportID *uuid.UUID) ([]*domain.ReportSavedView, error) {
	rows, err := s.getQueries(ctx).ListReportSavedViews(ctx, sqlc.ListReportSavedViewsParams{TenantID: tenantID, ReportID: uuidFromPtr(reportID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list report saved views", err, tenantIDField(tenantID))
	}
	return mapReportSavedViews(rows), nil
}

func (s *Store) DeleteReportSavedView(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteReportSavedView(ctx, sqlc.SoftDeleteReportSavedViewParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete report saved view", err, tenantIDField(tenantID), stringField("saved_view_id", id.String()))
	}
	return nil
}

func (s *Store) CreateReportExportJob(ctx context.Context, item *domain.ReportExportJob, actorID *uuid.UUID) (*domain.ReportExportJob, error) {
	row, err := s.getQueries(ctx).CreateReportExportJob(ctx, sqlc.CreateReportExportJobParams{ID: item.ID, TenantID: item.TenantID, ReportID: item.ReportID, SavedViewID: uuidFromPtr(item.SavedViewID), ExportFormat: item.ExportFormat, Status: item.Status, Filters: []byte(item.Filters), FileObjectKey: textFromPtr(item.FileObjectKey), ErrorMessage: textFromPtr(item.ErrorMessage), RequestedBy: uuidFromPtr(item.RequestedBy), StartedAt: timestamptzFromPtr(item.StartedAt), CompletedAt: timestamptzFromPtr(item.CompletedAt)})
	if err != nil {
		return nil, s.logDBError(ctx, "create report export job", err, tenantIDField(item.TenantID), stringField("report_id", item.ReportID.String()))
	}
	return mapReportExportJob(row), nil
}

func (s *Store) ListReportExportJobs(ctx context.Context, tenantID uuid.UUID, reportID *uuid.UUID, status *string, limit int32, offset int32) ([]*domain.ReportExportJob, error) {
	rows, err := s.getQueries(ctx).ListReportExportJobs(ctx, sqlc.ListReportExportJobsParams{TenantID: tenantID, Limit: limit, Offset: offset, ReportID: uuidFromPtr(reportID), Status: textFromPtr(status)})
	if err != nil {
		return nil, s.logDBError(ctx, "list report export jobs", err, tenantIDField(tenantID))
	}
	return mapReportExportJobs(rows), nil
}

func (s *Store) UpdateReportExportJobStatus(ctx context.Context, item *domain.ReportExportJob, actorID *uuid.UUID) (*domain.ReportExportJob, error) {
	row, err := s.getQueries(ctx).UpdateReportExportJobStatus(ctx, sqlc.UpdateReportExportJobStatusParams{TenantID: item.TenantID, ID: item.ID, Status: item.Status, FileObjectKey: textFromPtr(item.FileObjectKey), ErrorMessage: textFromPtr(item.ErrorMessage), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update report export job status", err, tenantIDField(item.TenantID), stringField("export_job_id", item.ID.String()))
	}
	return mapReportExportJob(row), nil
}

func (s *Store) UpsertReportSchedule(ctx context.Context, item *domain.ReportSchedule, actorID *uuid.UUID) (*domain.ReportSchedule, error) {
	row, err := s.getQueries(ctx).UpsertReportSchedule(ctx, sqlc.UpsertReportScheduleParams{ID: item.ID, TenantID: item.TenantID, ReportID: item.ReportID, SavedViewID: uuidFromPtr(item.SavedViewID), Name: item.Name, Frequency: item.Frequency, Timezone: item.Timezone, DeliveryChannels: []byte(item.DeliveryChannels), RecipientUserIds: []byte(item.RecipientUserIDs), RecipientEmails: []byte(item.RecipientEmails), NextRunAt: timestamptzFromPtr(item.NextRunAt), LastRunAt: timestamptzFromPtr(item.LastRunAt), IsActive: item.IsActive, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert report schedule", err, tenantIDField(item.TenantID), stringField("report_id", item.ReportID.String()))
	}
	return mapReportSchedule(row), nil
}

func (s *Store) ListReportSchedules(ctx context.Context, tenantID uuid.UUID, reportID *uuid.UUID) ([]*domain.ReportSchedule, error) {
	rows, err := s.getQueries(ctx).ListReportSchedules(ctx, sqlc.ListReportSchedulesParams{TenantID: tenantID, ReportID: uuidFromPtr(reportID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list report schedules", err, tenantIDField(tenantID))
	}
	return mapReportSchedules(rows), nil
}

func (s *Store) DeleteReportSchedule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteReportSchedule(ctx, sqlc.SoftDeleteReportScheduleParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete report schedule", err, tenantIDField(tenantID), stringField("schedule_id", id.String()))
	}
	return nil
}

func (s *Store) CreateReportSnapshot(ctx context.Context, item *domain.ReportSnapshot, actorID *uuid.UUID) (*domain.ReportSnapshot, error) {
	row, err := s.getQueries(ctx).CreateReportSnapshot(ctx, sqlc.CreateReportSnapshotParams{ID: item.ID, TenantID: item.TenantID, ReportID: item.ReportID, SavedViewID: uuidFromPtr(item.SavedViewID), SnapshotKey: item.SnapshotKey, PeriodStart: dateFromTime(item.PeriodStart), PeriodEnd: dateFromTime(item.PeriodEnd), Filters: []byte(item.Filters), Summary: []byte(item.Summary), RowCount: item.RowCount, Column11: timestamptzFromPtr(&item.GeneratedAt), GeneratedBy: uuidFromPtr(item.GeneratedBy)})
	if err != nil {
		return nil, s.logDBError(ctx, "create report snapshot", err, tenantIDField(item.TenantID), stringField("snapshot_key", item.SnapshotKey))
	}
	return mapReportSnapshot(row), nil
}

func (s *Store) ListReportSnapshots(ctx context.Context, tenantID uuid.UUID, reportID *uuid.UUID, limit int32, offset int32) ([]*domain.ReportSnapshot, error) {
	rows, err := s.getQueries(ctx).ListReportSnapshots(ctx, sqlc.ListReportSnapshotsParams{TenantID: tenantID, Limit: limit, Offset: offset, ReportID: uuidFromPtr(reportID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list report snapshots", err, tenantIDField(tenantID))
	}
	return mapReportSnapshots(rows), nil
}
