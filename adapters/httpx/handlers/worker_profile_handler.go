package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateWorkerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create worker profile", err, "tenant context is required")
		return
	}
	h.createWorkerProfileForTenant(w, r, tenantID, "create worker profile")
}

func (h *Handler) ListWorkerProfiles(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list worker profiles", err, "tenant context is required")
		return
	}
	h.listWorkerProfilesForTenant(w, r, tenantID, "list worker profiles")
}

func (h *Handler) GetWorkerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, workerProfileID, ok := h.workerProfileRequestIDs(w, r, "get worker profile")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkerProfile(r.Context(), tenantID, workerProfileID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get worker profile", err, "worker profile not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateWorkerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, workerProfileID, ok := h.workerProfileRequestIDs(w, r, "update worker profile")
	if !ok {
		return
	}
	h.updateWorkerProfileForTenant(w, r, tenantID, workerProfileID, "update worker profile")
}

func (h *Handler) DeleteWorkerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, workerProfileID, ok := h.workerProfileRequestIDs(w, r, "delete worker profile")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkerProfile(r.Context(), tenantID, workerProfileID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete worker profile", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantWorkerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant worker profile")
	if !ok {
		return
	}
	h.createWorkerProfileForTenant(w, r, tenantID, "create tenant worker profile")
}

func (h *Handler) ListTenantWorkerProfiles(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant worker profiles")
	if !ok {
		return
	}
	h.listWorkerProfilesForTenant(w, r, tenantID, "list tenant worker profiles")
}

func (h *Handler) GetTenantWorkerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, workerProfileID, ok := h.superAdminWorkerProfileRequestIDs(w, r, "get tenant worker profile")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkerProfile(r.Context(), tenantID, workerProfileID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant worker profile", err, "worker profile not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantWorkerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, workerProfileID, ok := h.superAdminWorkerProfileRequestIDs(w, r, "update tenant worker profile")
	if !ok {
		return
	}
	h.updateWorkerProfileForTenant(w, r, tenantID, workerProfileID, "update tenant worker profile")
}

func (h *Handler) DeleteTenantWorkerProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, workerProfileID, ok := h.superAdminWorkerProfileRequestIDs(w, r, "delete tenant worker profile")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkerProfile(r.Context(), tenantID, workerProfileID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant worker profile", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createWorkerProfileForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.WorkerProfileCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkerProfile(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateWorkerProfileForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, workerProfileID uuid.UUID, operation string) {
	var cmd ports.WorkerProfileCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = workerProfileID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkerProfile(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listWorkerProfilesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerTypeID, ok := h.optionalUUIDQuery(w, r, "worker_type_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListWorkerProfiles(r.Context(), domain.WorkerProfileFilter{
		TenantID:            tenantID,
		WorkerTypeID:        workerTypeID,
		ClassificationGroup: optionalStringQuery(r, "classification_group"),
		ProfileStatus:       optionalStringQuery(r, "profile_status"),
		Search:              optionalStringQuery(r, "search"),
	})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list worker profiles")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) workerProfileRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	workerProfileID, ok := h.workerProfileIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, workerProfileID, true
}

func (h *Handler) superAdminWorkerProfileRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	workerProfileID, ok := h.workerProfileIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, workerProfileID, true
}

func (h *Handler) workerProfileIDFromURL(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	workerProfileID, err := uuid.Parse(chi.URLParam(r, "workerProfileID"))
	if err != nil || workerProfileID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid worker profile id")
		return uuid.Nil, false
	}
	return workerProfileID, true
}
