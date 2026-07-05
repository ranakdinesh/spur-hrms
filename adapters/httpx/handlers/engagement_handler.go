package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateEngagement(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create engagement", err, "tenant context is required")
		return
	}
	h.createEngagementForTenant(w, r, tenantID, "create engagement")
}

func (h *Handler) ListEngagements(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list engagements", err, "tenant context is required")
		return
	}
	h.listEngagementsForTenant(w, r, tenantID, "list engagements")
}

func (h *Handler) GetEngagement(w http.ResponseWriter, r *http.Request) {
	tenantID, engagementID, ok := h.engagementRequestIDs(w, r, "get engagement")
	if !ok {
		return
	}
	item, err := h.svc.GetEngagement(r.Context(), tenantID, engagementID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get engagement", err, "engagement not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateEngagement(w http.ResponseWriter, r *http.Request) {
	tenantID, engagementID, ok := h.engagementRequestIDs(w, r, "update engagement")
	if !ok {
		return
	}
	h.updateEngagementForTenant(w, r, tenantID, engagementID, "update engagement")
}

func (h *Handler) UpdateEngagementStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, engagementID, ok := h.engagementRequestIDs(w, r, "update engagement status")
	if !ok {
		return
	}
	h.updateEngagementStatusForTenant(w, r, tenantID, engagementID, "update engagement status")
}

func (h *Handler) DeleteEngagement(w http.ResponseWriter, r *http.Request) {
	tenantID, engagementID, ok := h.engagementRequestIDs(w, r, "delete engagement")
	if !ok {
		return
	}
	if err := h.svc.DeleteEngagement(r.Context(), tenantID, engagementID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete engagement", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantEngagement(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant engagement")
	if !ok {
		return
	}
	h.createEngagementForTenant(w, r, tenantID, "create tenant engagement")
}

func (h *Handler) ListTenantEngagements(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant engagements")
	if !ok {
		return
	}
	h.listEngagementsForTenant(w, r, tenantID, "list tenant engagements")
}

func (h *Handler) GetTenantEngagement(w http.ResponseWriter, r *http.Request) {
	tenantID, engagementID, ok := h.superAdminEngagementRequestIDs(w, r, "get tenant engagement")
	if !ok {
		return
	}
	item, err := h.svc.GetEngagement(r.Context(), tenantID, engagementID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant engagement", err, "engagement not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantEngagement(w http.ResponseWriter, r *http.Request) {
	tenantID, engagementID, ok := h.superAdminEngagementRequestIDs(w, r, "update tenant engagement")
	if !ok {
		return
	}
	h.updateEngagementForTenant(w, r, tenantID, engagementID, "update tenant engagement")
}

func (h *Handler) UpdateTenantEngagementStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, engagementID, ok := h.superAdminEngagementRequestIDs(w, r, "update tenant engagement status")
	if !ok {
		return
	}
	h.updateEngagementStatusForTenant(w, r, tenantID, engagementID, "update tenant engagement status")
}

func (h *Handler) DeleteTenantEngagement(w http.ResponseWriter, r *http.Request) {
	tenantID, engagementID, ok := h.superAdminEngagementRequestIDs(w, r, "delete tenant engagement")
	if !ok {
		return
	}
	if err := h.svc.DeleteEngagement(r.Context(), tenantID, engagementID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant engagement", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createEngagementForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EngagementCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateEngagement(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateEngagementForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, engagementID uuid.UUID, operation string) {
	var cmd ports.EngagementCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = engagementID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEngagement(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateEngagementStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, engagementID uuid.UUID, operation string) {
	var cmd ports.EngagementStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.EngagementID = engagementID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEngagementStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listEngagementsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerProfileID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	departmentID, ok := h.optionalUUIDQuery(w, r, "department_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListEngagements(r.Context(), domain.EngagementFilter{
		TenantID:        tenantID,
		WorkerProfileID: workerProfileID,
		EngagementType:  optionalStringQuery(r, "engagement_type"),
		Status:          optionalStringQuery(r, "status"),
		DepartmentID:    departmentID,
		Search:          optionalStringQuery(r, "search"),
	})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list engagements")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) engagementRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	engagementID, ok := h.engagementIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, engagementID, true
}

func (h *Handler) superAdminEngagementRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	engagementID, ok := h.engagementIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, engagementID, true
}

func (h *Handler) engagementIDFromURL(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	engagementID, err := uuid.Parse(chi.URLParam(r, "engagementID"))
	if err != nil || engagementID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid engagement id")
		return uuid.Nil, false
	}
	return engagementID, true
}
