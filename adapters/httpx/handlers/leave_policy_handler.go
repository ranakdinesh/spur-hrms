package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateLeavePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create leave policy", err, "tenant context is required")
		return
	}
	h.createLeavePolicyForTenant(w, r, tenantID, "create leave policy")
}

func (h *Handler) ListLeavePolicies(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leave policies", err, "tenant context is required")
		return
	}
	h.listLeavePoliciesForTenant(w, r, tenantID, "list leave policies")
}

func (h *Handler) GetLeavePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.leavePolicyRequestIDs(w, r, "get leave policy")
	if !ok {
		return
	}
	item, err := h.svc.GetLeavePolicy(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get leave policy", err, "leave policy not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateLeavePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.leavePolicyRequestIDs(w, r, "update leave policy")
	if !ok {
		return
	}
	h.updateLeavePolicyForTenant(w, r, tenantID, id, "update leave policy")
}

func (h *Handler) DeleteLeavePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.leavePolicyRequestIDs(w, r, "delete leave policy")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeavePolicy(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete leave policy", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantLeavePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant leave policy")
	if !ok {
		return
	}
	h.createLeavePolicyForTenant(w, r, tenantID, "create tenant leave policy")
}

func (h *Handler) ListTenantLeavePolicies(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leave policies")
	if !ok {
		return
	}
	h.listLeavePoliciesForTenant(w, r, tenantID, "list tenant leave policies")
}

func (h *Handler) GetTenantLeavePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLeavePolicyRequestIDs(w, r, "get tenant leave policy")
	if !ok {
		return
	}
	item, err := h.svc.GetLeavePolicy(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant leave policy", err, "leave policy not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantLeavePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLeavePolicyRequestIDs(w, r, "update tenant leave policy")
	if !ok {
		return
	}
	h.updateLeavePolicyForTenant(w, r, tenantID, id, "update tenant leave policy")
}

func (h *Handler) DeleteTenantLeavePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLeavePolicyRequestIDs(w, r, "delete tenant leave policy")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeavePolicy(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant leave policy", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createLeavePolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LeavePolicyCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLeavePolicy(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listLeavePoliciesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListLeavePolicies(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave policies")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateLeavePolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.LeavePolicyCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLeavePolicy(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) leavePolicyRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "leavePolicyID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse leave policy id", err, "invalid leave policy id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminLeavePolicyRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "leavePolicyID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant leave policy id", err, "invalid leave policy id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
