package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetCommunicationProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get communication provider settings", err, "tenant context is required")
		return
	}
	h.getCommunicationProviderSettings(w, r, tenantID, "get communication provider settings")
}

func (h *Handler) UpsertCommunicationProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert communication provider settings", err, "tenant context is required")
		return
	}
	h.upsertCommunicationProviderSettings(w, r, tenantID, "upsert communication provider settings")
}

func (h *Handler) TestCommunicationProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "test communication provider settings", err, "tenant context is required")
		return
	}
	h.testCommunicationProviderSettings(w, r, tenantID, "test communication provider settings")
}

func (h *Handler) DeleteCommunicationProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete communication provider settings", err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "communicationProviderSettingsID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, "delete communication provider settings", err, "invalid communication provider settings id")
		return
	}
	if err := h.svc.DeleteCommunicationProviderSettings(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete communication provider settings", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTenantCommunicationProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant communication provider settings")
	if !ok {
		return
	}
	h.getCommunicationProviderSettings(w, r, tenantID, "get tenant communication provider settings")
}

func (h *Handler) UpsertTenantCommunicationProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant communication provider settings")
	if !ok {
		return
	}
	h.upsertCommunicationProviderSettings(w, r, tenantID, "upsert tenant communication provider settings")
}

func (h *Handler) TestTenantCommunicationProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "test tenant communication provider settings")
	if !ok {
		return
	}
	h.testCommunicationProviderSettings(w, r, tenantID, "test tenant communication provider settings")
}

func (h *Handler) DeleteTenantCommunicationProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "communicationProviderSettingsID", "delete tenant communication provider settings")
	if !ok {
		return
	}
	if err := h.svc.DeleteCommunicationProviderSettings(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant communication provider settings", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getCommunicationProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	item, err := h.svc.GetCommunicationProviderSettings(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) upsertCommunicationProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CommunicationProviderSettingsCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertCommunicationProviderSettings(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) testCommunicationProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CommunicationProviderTestCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.TestCommunicationProviderSettings(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
