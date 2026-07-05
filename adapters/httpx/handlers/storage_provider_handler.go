package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetStorageProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get storage provider settings", err, "tenant context is required")
		return
	}
	h.getStorageProviderSettings(w, r, tenantID, "get storage provider settings")
}

func (h *Handler) UpsertStorageProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert storage provider settings", err, "tenant context is required")
		return
	}
	h.upsertStorageProviderSettings(w, r, tenantID, "upsert storage provider settings")
}

func (h *Handler) TestStorageProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "test storage provider settings", err, "tenant context is required")
		return
	}
	h.testStorageProviderSettings(w, r, tenantID, "test storage provider settings")
}

func (h *Handler) DeleteStorageProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete storage provider settings", err, "tenant context is required")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "storageProviderSettingsID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, "delete storage provider settings", err, "invalid storage provider settings id")
		return
	}
	if err := h.svc.DeleteStorageProviderSettings(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete storage provider settings", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTenantStorageProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant storage provider settings")
	if !ok {
		return
	}
	h.getStorageProviderSettings(w, r, tenantID, "get tenant storage provider settings")
}

func (h *Handler) UpsertTenantStorageProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant storage provider settings")
	if !ok {
		return
	}
	h.upsertStorageProviderSettings(w, r, tenantID, "upsert tenant storage provider settings")
}

func (h *Handler) TestTenantStorageProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "test tenant storage provider settings")
	if !ok {
		return
	}
	h.testStorageProviderSettings(w, r, tenantID, "test tenant storage provider settings")
}

func (h *Handler) DeleteTenantStorageProviderSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "storageProviderSettingsID", "delete tenant storage provider settings")
	if !ok {
		return
	}
	if err := h.svc.DeleteStorageProviderSettings(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant storage provider settings", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getStorageProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	item, err := h.svc.GetStorageProviderSettings(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) upsertStorageProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.StorageProviderSettingsCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertStorageProviderSettings(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) testStorageProviderSettings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	cmd := ports.StorageProviderTestCommand{TenantID: tenantID, ActorID: h.actorIDFromRequest(r)}
	item, err := h.svc.TestStorageProviderSettings(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
