package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetPayCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get pay cycle", err, "tenant context is required")
		return
	}
	h.getPayCycleForTenant(w, r, tenantID, "get pay cycle")
}

func (h *Handler) UpsertPayCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert pay cycle", err, "tenant context is required")
		return
	}
	h.upsertPayCycleForTenant(w, r, tenantID, "upsert pay cycle")
}

func (h *Handler) DeletePayCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete pay cycle", err, "tenant context is required")
		return
	}
	if err := h.svc.DeletePayCycle(r.Context(), tenantID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete pay cycle", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ResolvePayCyclePeriod(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "resolve pay cycle period", err, "tenant context is required")
		return
	}
	h.resolvePayCyclePeriodForTenant(w, r, tenantID, "resolve pay cycle period")
}

func (h *Handler) GetTenantPayCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant pay cycle")
	if !ok {
		return
	}
	h.getPayCycleForTenant(w, r, tenantID, "get tenant pay cycle")
}

func (h *Handler) UpsertTenantPayCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant pay cycle")
	if !ok {
		return
	}
	h.upsertPayCycleForTenant(w, r, tenantID, "upsert tenant pay cycle")
}

func (h *Handler) DeleteTenantPayCycle(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "delete tenant pay cycle")
	if !ok {
		return
	}
	if err := h.svc.DeletePayCycle(r.Context(), tenantID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant pay cycle", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ResolveTenantPayCyclePeriod(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "resolve tenant pay cycle period")
	if !ok {
		return
	}
	h.resolvePayCyclePeriodForTenant(w, r, tenantID, "resolve tenant pay cycle period")
}

func (h *Handler) getPayCycleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	item, err := h.svc.GetPayCycle(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, "pay cycle configuration not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) upsertPayCycleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PayCycleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode pay cycle request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertPayCycle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) resolvePayCyclePeriodForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	month, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid month")
		return
	}
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid year")
		return
	}
	period, err := h.svc.ResolvePayCyclePeriod(r.Context(), ports.PayCyclePeriodQuery{TenantID: tenantID, Month: month, Year: year})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, period)
}
