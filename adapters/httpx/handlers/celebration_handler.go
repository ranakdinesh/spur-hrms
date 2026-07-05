package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateCelebration(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create celebration", err, "tenant context is required")
		return
	}
	var cmd ports.CelebrationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode celebration create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateCelebration(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create celebration", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListCelebrations(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list celebrations", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListCelebrations(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list celebrations", err, "failed to list celebrations")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetCelebration(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "celebrationID", "get celebration")
	if !ok {
		return
	}
	item, err := h.svc.GetCelebration(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get celebration", err, "celebration not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateCelebration(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "celebrationID", "update celebration")
	if !ok {
		return
	}
	var cmd ports.CelebrationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode celebration update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCelebration(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update celebration", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteCelebration(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "celebrationID", "delete celebration")
	if !ok {
		return
	}
	if err := h.svc.DeleteCelebration(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete celebration", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantCelebration(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant celebration")
	if !ok {
		return
	}
	var cmd ports.CelebrationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant celebration create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateCelebration(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant celebration", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListTenantCelebrations(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant celebrations")
	if !ok {
		return
	}
	items, err := h.svc.ListCelebrations(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant celebrations", err, "failed to list celebrations")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetTenantCelebration(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "celebrationID", "get tenant celebration")
	if !ok {
		return
	}
	item, err := h.svc.GetCelebration(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant celebration", err, "celebration not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantCelebration(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "celebrationID", "update tenant celebration")
	if !ok {
		return
	}
	var cmd ports.CelebrationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant celebration update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCelebration(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant celebration", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantCelebration(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "celebrationID", "delete tenant celebration")
	if !ok {
		return
	}
	if err := h.svc.DeleteCelebration(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant celebration", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
