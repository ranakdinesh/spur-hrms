package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateLeaveType(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create leave type", err, "tenant context is required")
		return
	}
	h.createLeaveTypeForTenant(w, r, tenantID, "create leave type")
}

func (h *Handler) ListLeaveTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leave types", err, "tenant context is required")
		return
	}
	h.listLeaveTypesForTenant(w, r, tenantID, "list leave types")
}

func (h *Handler) GetLeaveType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.leaveTypeRequestIDs(w, r, "get leave type")
	if !ok {
		return
	}
	item, err := h.svc.GetLeaveType(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get leave type", err, "leave type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateLeaveType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.leaveTypeRequestIDs(w, r, "update leave type")
	if !ok {
		return
	}
	h.updateLeaveTypeForTenant(w, r, tenantID, id, "update leave type")
}

func (h *Handler) DeleteLeaveType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.leaveTypeRequestIDs(w, r, "delete leave type")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeaveType(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete leave type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantLeaveType(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant leave type")
	if !ok {
		return
	}
	h.createLeaveTypeForTenant(w, r, tenantID, "create tenant leave type")
}

func (h *Handler) ListTenantLeaveTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leave types")
	if !ok {
		return
	}
	h.listLeaveTypesForTenant(w, r, tenantID, "list tenant leave types")
}

func (h *Handler) GetTenantLeaveType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLeaveTypeRequestIDs(w, r, "get tenant leave type")
	if !ok {
		return
	}
	item, err := h.svc.GetLeaveType(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant leave type", err, "leave type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantLeaveType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLeaveTypeRequestIDs(w, r, "update tenant leave type")
	if !ok {
		return
	}
	h.updateLeaveTypeForTenant(w, r, tenantID, id, "update tenant leave type")
}

func (h *Handler) DeleteTenantLeaveType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLeaveTypeRequestIDs(w, r, "delete tenant leave type")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeaveType(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant leave type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createLeaveTypeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LeaveTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLeaveType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listLeaveTypesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListLeaveTypes(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateLeaveTypeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.LeaveTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLeaveType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) leaveTypeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "leaveTypeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse leave type id", err, "invalid leave type id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminLeaveTypeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "leaveTypeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant leave type id", err, "invalid leave type id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
