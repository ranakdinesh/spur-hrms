package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetEmailProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get email provider settings", err, "tenant context is required")
		return
	}
	h.getEmailProviderSettings(w, r, tenantID, "get email provider settings")
}

func (h *Handler) UpsertEmailProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert email provider settings", err, "tenant context is required")
		return
	}
	h.upsertEmailProviderSettings(w, r, tenantID, "upsert email provider settings")
}

func (h *Handler) TestEmailProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "test email provider settings", err, "tenant context is required")
		return
	}
	h.testEmailProviderSettings(w, r, tenantID, "test email provider settings")
}

func (h *Handler) DeleteEmailProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete email provider settings", err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "emailProviderSettingsID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, "delete email provider settings", err, "invalid email provider settings id")
		return
	}
	if err := h.svc.DeleteEmailProviderSettings(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete email provider settings", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTenantEmailProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant email provider settings")
	if !ok {
		return
	}
	h.getEmailProviderSettings(w, r, tenantID, "get tenant email provider settings")
}

func (h *Handler) UpsertTenantEmailProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant email provider settings")
	if !ok {
		return
	}
	h.upsertEmailProviderSettings(w, r, tenantID, "upsert tenant email provider settings")
}

func (h *Handler) TestTenantEmailProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "test tenant email provider settings")
	if !ok {
		return
	}
	h.testEmailProviderSettings(w, r, tenantID, "test tenant email provider settings")
}

func (h *Handler) DeleteTenantEmailProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "emailProviderSettingsID", "delete tenant email provider settings")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmailProviderSettings(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant email provider settings", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getEmailProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	item, err := h.svc.GetEmailProviderSettings(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) upsertEmailProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EmailProviderSettingsCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertEmailProviderSettings(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) testEmailProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EmailProviderTestCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.TestEmailProviderSettings(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
