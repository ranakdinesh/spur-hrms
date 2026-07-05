package postgres

import (
	"encoding/json"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapReportCatalogItem(row sqlc.HrmsReportCatalog) *domain.ReportCatalogItem {
	return &domain.ReportCatalogItem{ID: row.ID, TenantID: row.TenantID, ReportCode: row.ReportCode, Module: row.Module, Name: row.Name, Description: ptrFromText(row.Description), Category: row.Category, Scope: row.Scope, PermissionKey: row.PermissionKey, DefaultFilters: json.RawMessage(row.DefaultFilters), SupportedFilters: json.RawMessage(row.SupportedFilters), OutputColumns: json.RawMessage(row.OutputColumns), DrilldownContract: json.RawMessage(row.DrilldownContract), IsSystem: row.IsSystem, IsActive: row.IsActive, SortOrder: row.SortOrder, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapReportCatalogItems(rows []sqlc.HrmsReportCatalog) []*domain.ReportCatalogItem {
	items := make([]*domain.ReportCatalogItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapReportCatalogItem(row))
	}
	return items
}

func mapReportSavedView(row sqlc.HrmsReportSavedView) *domain.ReportSavedView {
	return &domain.ReportSavedView{ID: row.ID, TenantID: row.TenantID, ReportID: row.ReportID, Name: row.Name, Description: ptrFromText(row.Description), Visibility: row.Visibility, Filters: json.RawMessage(row.Filters), Columns: json.RawMessage(row.Columns), IsFavorite: row.IsFavorite, OwnerUserID: ptrFromUUID(row.OwnerUserID), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapReportSavedViews(rows []sqlc.HrmsReportSavedView) []*domain.ReportSavedView {
	items := make([]*domain.ReportSavedView, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapReportSavedView(row))
	}
	return items
}

func mapReportExportJob(row sqlc.HrmsReportExportJob) *domain.ReportExportJob {
	return &domain.ReportExportJob{ID: row.ID, TenantID: row.TenantID, ReportID: row.ReportID, SavedViewID: ptrFromUUID(row.SavedViewID), ExportFormat: row.ExportFormat, Status: row.Status, Filters: json.RawMessage(row.Filters), FileObjectKey: ptrFromText(row.FileObjectKey), ErrorMessage: ptrFromText(row.ErrorMessage), RequestedBy: ptrFromUUID(row.RequestedBy), StartedAt: ptrFromTimestamptz(row.StartedAt), CompletedAt: ptrFromTimestamptz(row.CompletedAt), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapReportExportJobs(rows []sqlc.HrmsReportExportJob) []*domain.ReportExportJob {
	items := make([]*domain.ReportExportJob, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapReportExportJob(row))
	}
	return items
}

func mapReportSchedule(row sqlc.HrmsReportSchedule) *domain.ReportSchedule {
	return &domain.ReportSchedule{ID: row.ID, TenantID: row.TenantID, ReportID: row.ReportID, SavedViewID: ptrFromUUID(row.SavedViewID), Name: row.Name, Frequency: row.Frequency, Timezone: row.Timezone, DeliveryChannels: json.RawMessage(row.DeliveryChannels), RecipientUserIDs: json.RawMessage(row.RecipientUserIds), RecipientEmails: json.RawMessage(row.RecipientEmails), NextRunAt: ptrFromTimestamptz(row.NextRunAt), LastRunAt: ptrFromTimestamptz(row.LastRunAt), IsActive: row.IsActive, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapReportSchedules(rows []sqlc.HrmsReportSchedule) []*domain.ReportSchedule {
	items := make([]*domain.ReportSchedule, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapReportSchedule(row))
	}
	return items
}

func mapReportSnapshot(row sqlc.HrmsReportSnapshot) *domain.ReportSnapshot {
	return &domain.ReportSnapshot{ID: row.ID, TenantID: row.TenantID, ReportID: row.ReportID, SavedViewID: ptrFromUUID(row.SavedViewID), SnapshotKey: row.SnapshotKey, PeriodStart: timeFromDate(row.PeriodStart), PeriodEnd: timeFromDate(row.PeriodEnd), Filters: json.RawMessage(row.Filters), Summary: json.RawMessage(row.Summary), RowCount: row.RowCount, GeneratedAt: timeFromTimestamptz(row.GeneratedAt), GeneratedBy: ptrFromUUID(row.GeneratedBy), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapReportSnapshots(rows []sqlc.HrmsReportSnapshot) []*domain.ReportSnapshot {
	items := make([]*domain.ReportSnapshot, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapReportSnapshot(row))
	}
	return items
}
