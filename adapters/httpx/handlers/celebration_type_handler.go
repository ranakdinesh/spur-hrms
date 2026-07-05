package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateCelebrationType(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create celebration type", err, "tenant context is required")
		return
	}
	var cmd ports.CelebrationTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode celebration type create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateCelebrationType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create celebration type", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListCelebrationTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list celebration types", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListCelebrationTypes(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list celebration types", err, "failed to list celebration types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetCelebrationType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "celebrationTypeID", "get celebration type")
	if !ok {
		return
	}
	item, err := h.svc.GetCelebrationType(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get celebration type", err, "celebration type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateCelebrationType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "celebrationTypeID", "update celebration type")
	if !ok {
		return
	}
	var cmd ports.CelebrationTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode celebration type update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCelebrationType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update celebration type", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteCelebrationType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "celebrationTypeID", "delete celebration type")
	if !ok {
		return
	}
	if err := h.svc.DeleteCelebrationType(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete celebration type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantCelebrationType(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant celebration type")
	if !ok {
		return
	}
	var cmd ports.CelebrationTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant celebration type create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateCelebrationType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant celebration type", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListTenantCelebrationTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant celebration types")
	if !ok {
		return
	}
	items, err := h.svc.ListCelebrationTypes(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant celebration types", err, "failed to list celebration types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetTenantCelebrationType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "celebrationTypeID", "get tenant celebration type")
	if !ok {
		return
	}
	item, err := h.svc.GetCelebrationType(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant celebration type", err, "celebration type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantCelebrationType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "celebrationTypeID", "update tenant celebration type")
	if !ok {
		return
	}
	var cmd ports.CelebrationTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant celebration type update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCelebrationType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant celebration type", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantCelebrationType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "celebrationTypeID", "delete tenant celebration type")
	if !ok {
		return
	}
	if err := h.svc.DeleteCelebrationType(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant celebration type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
