package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create work log", err, "tenant context is required")
		return
	}
	h.createWorkLogForTenant(w, r, tenantID, "create work log")
}

func (h *Handler) ListWorkLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list work logs", err, "tenant context is required")
		return
	}
	h.listWorkLogsForTenant(w, r, tenantID, "list work logs")
}

func (h *Handler) ListWorkLogRollups(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list work log rollups", err, "tenant context is required")
		return
	}
	h.listWorkLogRollupsForTenant(w, r, tenantID, "list work log rollups")
}

func (h *Handler) GetWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.workLogRequestIDs(w, r, "get work log")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkLog(r.Context(), tenantID, workLogID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get work log", err, "work log not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.workLogRequestIDs(w, r, "update work log")
	if !ok {
		return
	}
	h.updateWorkLogForTenant(w, r, tenantID, workLogID, "update work log")
}

func (h *Handler) SubmitWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.workLogRequestIDs(w, r, "submit work log")
	if !ok {
		return
	}
	item, err := h.svc.SubmitWorkLog(r.Context(), tenantID, workLogID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "submit work log", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ReviewWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.workLogRequestIDs(w, r, "review work log")
	if !ok {
		return
	}
	h.reviewWorkLogForTenant(w, r, tenantID, workLogID, "review work log")
}

func (h *Handler) DeleteWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.workLogRequestIDs(w, r, "delete work log")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkLog(r.Context(), tenantID, workLogID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete work log", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant work log")
	if !ok {
		return
	}
	h.createWorkLogForTenant(w, r, tenantID, "create tenant work log")
}

func (h *Handler) ListTenantWorkLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant work logs")
	if !ok {
		return
	}
	h.listWorkLogsForTenant(w, r, tenantID, "list tenant work logs")
}

func (h *Handler) ListTenantWorkLogRollups(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant work log rollups")
	if !ok {
		return
	}
	h.listWorkLogRollupsForTenant(w, r, tenantID, "list tenant work log rollups")
}

func (h *Handler) GetTenantWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.superAdminWorkLogRequestIDs(w, r, "get tenant work log")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkLog(r.Context(), tenantID, workLogID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant work log", err, "work log not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.superAdminWorkLogRequestIDs(w, r, "update tenant work log")
	if !ok {
		return
	}
	h.updateWorkLogForTenant(w, r, tenantID, workLogID, "update tenant work log")
}

func (h *Handler) SubmitTenantWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.superAdminWorkLogRequestIDs(w, r, "submit tenant work log")
	if !ok {
		return
	}
	item, err := h.svc.SubmitWorkLog(r.Context(), tenantID, workLogID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "submit tenant work log", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ReviewTenantWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.superAdminWorkLogRequestIDs(w, r, "review tenant work log")
	if !ok {
		return
	}
	h.reviewWorkLogForTenant(w, r, tenantID, workLogID, "review tenant work log")
}

func (h *Handler) DeleteTenantWorkLog(w http.ResponseWriter, r *http.Request) {
	tenantID, workLogID, ok := h.superAdminWorkLogRequestIDs(w, r, "delete tenant work log")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkLog(r.Context(), tenantID, workLogID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant work log", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createWorkLogForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.WorkLogCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkLog(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateWorkLogForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, workLogID uuid.UUID, operation string) {
	var cmd ports.WorkLogCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = workLogID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkLog(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) reviewWorkLogForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, workLogID uuid.UUID, operation string) {
	var cmd ports.WorkLogStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.WorkLogID = workLogID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ReviewWorkLog(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listWorkLogsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter, ok := h.workLogFilterFromRequest(w, r, tenantID, operation)
	if !ok {
		return
	}
	items, err := h.svc.ListWorkLogs(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list work logs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listWorkLogRollupsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter, ok := h.workLogFilterFromRequest(w, r, tenantID, operation)
	if !ok {
		return
	}
	items, err := h.svc.ListWorkLogRollups(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list work log rollups")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) workLogFilterFromRequest(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (domain.WorkLogFilter, bool) {
	engagementID, ok := h.optionalUUIDQuery(w, r, "engagement_id", operation)
	if !ok {
		return domain.WorkLogFilter{}, false
	}
	workerProfileID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return domain.WorkLogFilter{}, false
	}
	dateFrom, ok := h.optionalWorkLogDateQuery(w, r, "date_from", operation)
	if !ok {
		return domain.WorkLogFilter{}, false
	}
	dateTo, ok := h.optionalWorkLogDateQuery(w, r, "date_to", operation)
	if !ok {
		return domain.WorkLogFilter{}, false
	}
	return domain.WorkLogFilter{
		TenantID:        tenantID,
		EngagementID:    engagementID,
		WorkerProfileID: workerProfileID,
		Status:          optionalStringQuery(r, "status"),
		DateFrom:        dateFrom,
		DateTo:          dateTo,
		Search:          optionalStringQuery(r, "search"),
	}, true
}

func (h *Handler) optionalWorkLogDateQuery(w http.ResponseWriter, r *http.Request, key string, operation string) (*time.Time, bool) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, true
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid "+key)
		return nil, false
	}
	return &parsed, true
}

func (h *Handler) workLogRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	workLogID, ok := h.workLogIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, workLogID, true
}

func (h *Handler) superAdminWorkLogRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	workLogID, ok := h.workLogIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, workLogID, true
}

func (h *Handler) workLogIDFromURL(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	workLogID, err := uuid.Parse(chi.URLParam(r, "workLogID"))
	if err != nil || workLogID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid work log id")
		return uuid.Nil, false
	}
	return workLogID, true
}
