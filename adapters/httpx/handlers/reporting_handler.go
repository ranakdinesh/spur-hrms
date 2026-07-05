package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListTenantReportCatalog(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list report catalog", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListReportCatalog(r.Context(), ports.ReportCatalogQuery{TenantID: tenantID, Module: reportOptionalQuery(r, "module"), Scope: reportOptionalQuery(r, "scope"), ActorID: h.actorIDFromRequest(r)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list report catalog", err, "failed to list reports")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListTenantReportSavedViews(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list report saved views", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListReportSavedViews(r.Context(), reportListQueryFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list report saved views", err, "failed to list saved views")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpsertTenantReportSavedView(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert report saved view", err, "tenant context is required")
		return
	}
	var cmd ports.ReportSavedViewCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode report saved view", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	if id := chi.URLParam(r, "savedViewID"); id != "" {
		parsed, err := uuid.Parse(id)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "parse report saved view id", err, "invalid saved view id")
			return
		}
		cmd.ID = parsed
	}
	item, err := h.svc.UpsertReportSavedView(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert report saved view", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantReportSavedView(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "savedViewID", "delete report saved view")
	if !ok {
		return
	}
	if err := h.svc.DeleteReportSavedView(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete report saved view", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantReportExportJobs(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list report export jobs", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListReportExportJobs(r.Context(), reportListQueryFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list report export jobs", err, "failed to list export jobs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateTenantReportExportJob(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create report export job", err, "tenant context is required")
		return
	}
	var cmd ports.ReportExportJobCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode report export job", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateReportExportJob(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create report export job", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateTenantReportExportJobStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "exportJobID", "update report export job status")
	if !ok {
		return
	}
	var cmd ports.ReportExportJobStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode report export status", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateReportExportJobStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update report export job status", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListTenantReportSchedules(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list report schedules", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListReportSchedules(r.Context(), reportListQueryFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list report schedules", err, "failed to list schedules")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpsertTenantReportSchedule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert report schedule", err, "tenant context is required")
		return
	}
	var cmd ports.ReportScheduleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode report schedule", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	if id := chi.URLParam(r, "scheduleID"); id != "" {
		parsed, err := uuid.Parse(id)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "parse report schedule id", err, "invalid schedule id")
			return
		}
		cmd.ID = parsed
	}
	item, err := h.svc.UpsertReportSchedule(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert report schedule", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantReportSchedule(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "scheduleID", "delete report schedule")
	if !ok {
		return
	}
	if err := h.svc.DeleteReportSchedule(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete report schedule", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantReportSnapshots(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list report snapshots", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListReportSnapshots(r.Context(), reportListQueryFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list report snapshots", err, "failed to list snapshots")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateTenantReportSnapshot(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create report snapshot", err, "tenant context is required")
		return
	}
	var cmd ports.ReportSnapshotCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode report snapshot", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateReportSnapshot(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create report snapshot", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) GetTenantReportDataset(w http.ResponseWriter, r *http.Request) {
	tenantID, reportID, ok := h.tenantAndURLUUID(w, r, "reportID", "get report dataset")
	if !ok {
		return
	}
	dataset, err := h.svc.BuildReportDataset(r.Context(), reportDatasetQueryFromRequest(r, tenantID, reportID, ""))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "get report dataset", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, dataset)
}

func (h *Handler) DownloadTenantReport(w http.ResponseWriter, r *http.Request) {
	tenantID, reportID, ok := h.tenantAndURLUUID(w, r, "reportID", "download report")
	if !ok {
		return
	}
	h.downloadReportForTenant(w, r, reportDatasetQueryFromRequest(r, tenantID, reportID, ""), "download report")
}

func (h *Handler) GetTenantReportDatasetByCode(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get report dataset by code", err, "tenant context is required")
		return
	}
	dataset, err := h.svc.BuildReportDataset(r.Context(), reportDatasetQueryFromRequest(r, tenantID, uuid.Nil, chi.URLParam(r, "reportCode")))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "get report dataset by code", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, dataset)
}

func (h *Handler) DownloadTenantReportByCode(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "download report by code", err, "tenant context is required")
		return
	}
	h.downloadReportForTenant(w, r, reportDatasetQueryFromRequest(r, tenantID, uuid.Nil, chi.URLParam(r, "reportCode")), "download report by code")
}

func (h *Handler) tenantAndURLUUID(w http.ResponseWriter, r *http.Request, param string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func reportListQueryFromRequest(r *http.Request, tenantID uuid.UUID) ports.ReportListQuery {
	limit := parseInt32Query(r, "limit", 25)
	offset := parseInt32Query(r, "offset", 0)
	return ports.ReportListQuery{TenantID: tenantID, ReportID: reportOptionalUUIDQuery(r, "report_id"), Status: reportOptionalQuery(r, "status"), Limit: limit, Offset: offset}
}

func reportDatasetQueryFromRequest(r *http.Request, tenantID uuid.UUID, reportID uuid.UUID, reportCode string) ports.ReportDatasetQuery {
	month := parseInt32Query(r, "month", 0)
	year := parseInt32Query(r, "year", 0)
	return ports.ReportDatasetQuery{TenantID: tenantID, ReportID: reportID, ReportCode: reportCode, StartDate: r.URL.Query().Get("start_date"), EndDate: r.URL.Query().Get("end_date"), Month: month, Year: year, FYID: reportOptionalUUIDQuery(r, "fy_id")}
}

func (h *Handler) downloadReportForTenant(w http.ResponseWriter, r *http.Request, query ports.ReportDatasetQuery, operation string) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "pdf"
	}
	download, err := h.svc.ExportReportDataset(r.Context(), query, format)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.Header().Set("Content-Type", download.ContentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+download.FileName+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(download.Content)
}

func reportOptionalQuery(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func reportOptionalUUIDQuery(r *http.Request, key string) *uuid.UUID {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	parsed, err := uuid.Parse(value)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseInt32Query(r *http.Request, key string, fallback int32) int32 {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return int32(parsed)
}
