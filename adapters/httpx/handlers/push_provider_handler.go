package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetPushProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get push provider settings", err, "tenant context is required")
		return
	}
	h.getPushProviderSettings(w, r, tenantID, "get push provider settings")
}

func (h *Handler) UpsertPushProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert push provider settings", err, "tenant context is required")
		return
	}
	h.upsertPushProviderSettings(w, r, tenantID, "upsert push provider settings")
}

func (h *Handler) TestPushProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "test push provider settings", err, "tenant context is required")
		return
	}
	h.testPushProviderSettings(w, r, tenantID, "test push provider settings")
}

func (h *Handler) DeletePushProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete push provider settings", err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "pushProviderSettingsID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, "delete push provider settings", err, "invalid push provider settings id")
		return
	}
	if err := h.svc.DeletePushProviderSettings(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete push provider settings", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTenantPushProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant push provider settings")
	if !ok {
		return
	}
	h.getPushProviderSettings(w, r, tenantID, "get tenant push provider settings")
}

func (h *Handler) UpsertTenantPushProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant push provider settings")
	if !ok {
		return
	}
	h.upsertPushProviderSettings(w, r, tenantID, "upsert tenant push provider settings")
}

func (h *Handler) TestTenantPushProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "test tenant push provider settings")
	if !ok {
		return
	}
	h.testPushProviderSettings(w, r, tenantID, "test tenant push provider settings")
}

func (h *Handler) DeleteTenantPushProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "pushProviderSettingsID", "delete tenant push provider settings")
	if !ok {
		return
	}
	if err := h.svc.DeletePushProviderSettings(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant push provider settings", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getPushProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	item, err := h.svc.GetPushProviderSettings(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) upsertPushProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PushProviderSettingsCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertPushProviderSettings(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) testPushProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PushProviderTestCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.TestPushProviderSettings(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
