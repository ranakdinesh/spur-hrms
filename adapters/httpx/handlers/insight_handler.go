package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListInsights(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list insights", err, "tenant context is required")
		return
	}
	h.listInsightsForTenant(w, r, tenantID, "list insights")
}

func (h *Handler) RefreshInsights(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "refresh insights", err, "tenant context is required")
		return
	}
	h.refreshInsightsForTenant(w, r, tenantID, "refresh insights")
}

func (h *Handler) UpdateInsightStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, insightID, ok := h.tenantAndURLUUID(w, r, "insightID", "update insight status")
	if !ok {
		return
	}
	h.updateInsightStatusForTenant(w, r, tenantID, insightID, "update insight status")
}

func (h *Handler) ListInsightEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, insightID, ok := h.tenantAndURLUUID(w, r, "insightID", "list insight events")
	if !ok {
		return
	}
	h.listInsightEventsForTenant(w, r, tenantID, insightID, "list insight events")
}

func (h *Handler) ListTenantInsights(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant insights")
	if !ok {
		return
	}
	h.listInsightsForTenant(w, r, tenantID, "list tenant insights")
}

func (h *Handler) RefreshTenantInsights(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "refresh tenant insights")
	if !ok {
		return
	}
	h.refreshInsightsForTenant(w, r, tenantID, "refresh tenant insights")
}

func (h *Handler) UpdateTenantInsightStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "update tenant insight status")
	if !ok {
		return
	}
	insightID, err := uuid.Parse(chi.URLParam(r, "insightID"))
	if err != nil || insightID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant insight status", err, "invalid insight id")
		return
	}
	h.updateInsightStatusForTenant(w, r, tenantID, insightID, "update tenant insight status")
}

func (h *Handler) ListTenantInsightEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant insight events")
	if !ok {
		return
	}
	insightID, err := uuid.Parse(chi.URLParam(r, "insightID"))
	if err != nil || insightID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, "list tenant insight events", err, "invalid insight id")
		return
	}
	h.listInsightEventsForTenant(w, r, tenantID, insightID, "list tenant insight events")
}

func (h *Handler) listInsightsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workspace, err := h.svc.ListInsightWorkspace(r.Context(), insightFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list insights")
		return
	}
	respondJSON(w, http.StatusOK, workspace)
}

func (h *Handler) refreshInsightsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workspace, err := h.svc.RefreshInsights(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to refresh insights")
		return
	}
	respondJSON(w, http.StatusOK, workspace)
}

func (h *Handler) updateInsightStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, insightID uuid.UUID, operation string) {
	var cmd ports.InsightStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = insightID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateInsightStatus(r.Context(), cmd)
	if err != nil {
		status := http.StatusBadRequest
		if err == domain.ErrInsightNotFound {
			status = http.StatusNotFound
		}
		h.respondError(w, r, status, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listInsightEventsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, insightID uuid.UUID, operation string) {
	events, err := h.svc.ListInsightEvents(r.Context(), tenantID, insightID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list insight audit events")
		return
	}
	respondJSON(w, http.StatusOK, events)
}

func insightFilterFromRequest(r *http.Request, tenantID uuid.UUID) domain.InsightFilter {
	limit := int32(100)
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}
	offset := int32(0)
	if raw := r.URL.Query().Get("offset"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			offset = int32(parsed)
		}
	}
	return domain.InsightFilter{TenantID: tenantID, Status: reportOptionalQuery(r, "status"), Severity: reportOptionalQuery(r, "severity"), Category: reportOptionalQuery(r, "category"), InsightType: reportOptionalQuery(r, "insight_type"), Limit: limit, Offset: offset}
}
