package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create working hour", err, "tenant context is required")
		return
	}
	var cmd ports.WorkingHourCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode working hour create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkingHour(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create working hour", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListWorkingHours(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list working hours", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListWorkingHours(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list working hours", err, "failed to list working hours")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ResolveWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "resolve working hour", err, "tenant context is required")
		return
	}
	cmd, ok := h.resolveWorkingHourCommand(w, r, tenantID, "resolve working hour")
	if !ok {
		return
	}
	item, err := h.svc.ResolveWorkingHour(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "resolve working hour", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) GetWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.workingHourRequestIDs(w, r, "get working hour")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkingHour(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get working hour", err, "working hour not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.workingHourRequestIDs(w, r, "update working hour")
	if !ok {
		return
	}
	var cmd ports.WorkingHourCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode working hour update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkingHour(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update working hour", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.workingHourRequestIDs(w, r, "delete working hour")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkingHour(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete working hour", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CopyWorkingHoursToBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "copy working hours to branch", err, "tenant context is required")
		return
	}
	branchID, ok := h.copyWorkingHoursBranchID(w, r, "copy working hours to branch")
	if !ok {
		return
	}
	items, err := h.svc.CopyTenantWorkingHoursToBranch(r.Context(), ports.CopyWorkingHoursCommand{TenantID: tenantID, BranchID: branchID, ActorID: h.actorIDFromRequest(r)})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "copy working hours to branch", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateTenantWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant working hour")
	if !ok {
		return
	}
	var cmd ports.WorkingHourCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant working hour create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkingHour(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant working hour", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListTenantWorkingHours(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant working hours")
	if !ok {
		return
	}
	items, err := h.svc.ListWorkingHours(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant working hours", err, "failed to list working hours")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ResolveTenantWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "resolve tenant working hour")
	if !ok {
		return
	}
	cmd, ok := h.resolveWorkingHourCommand(w, r, tenantID, "resolve tenant working hour")
	if !ok {
		return
	}
	item, err := h.svc.ResolveWorkingHour(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "resolve tenant working hour", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) GetTenantWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminWorkingHourRequestIDs(w, r, "get tenant working hour")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkingHour(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant working hour", err, "working hour not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminWorkingHourRequestIDs(w, r, "update tenant working hour")
	if !ok {
		return
	}
	var cmd ports.WorkingHourCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant working hour update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkingHour(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant working hour", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantWorkingHour(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminWorkingHourRequestIDs(w, r, "delete tenant working hour")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkingHour(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant working hour", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CopyTenantWorkingHoursToBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "copy tenant working hours to branch")
	if !ok {
		return
	}
	branchID, ok := h.copyWorkingHoursBranchID(w, r, "copy tenant working hours to branch")
	if !ok {
		return
	}
	items, err := h.svc.CopyTenantWorkingHoursToBranch(r.Context(), ports.CopyWorkingHoursCommand{TenantID: tenantID, BranchID: branchID, ActorID: h.actorIDFromRequest(r)})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "copy tenant working hours to branch", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) workingHourRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "workingHourID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid working hour id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminWorkingHourRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "workingHourID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid working hour id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) resolveWorkingHourCommand(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (ports.ResolveWorkingHourCommand, bool) {
	branchID, ok := h.optionalUUIDQuery(w, r, "branch_id", operation)
	if !ok {
		return ports.ResolveWorkingHourCommand{}, false
	}
	userID, ok := h.optionalUUIDQuery(w, r, "user_id", operation)
	if !ok {
		return ports.ResolveWorkingHourCommand{}, false
	}
	return ports.ResolveWorkingHourCommand{TenantID: tenantID, BranchID: branchID, UserID: userID, DayOfWeek: r.URL.Query().Get("day_of_week")}, true
}

func (h *Handler) copyWorkingHoursBranchID(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	var body struct {
		BranchID uuid.UUID `json:"branch_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return uuid.Nil, false
	}
	if body.BranchID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, nil, "branch_id is required")
		return uuid.Nil, false
	}
	return body.BranchID, true
}

func (h *Handler) optionalUUIDQuery(w http.ResponseWriter, r *http.Request, key string, operation string) (*uuid.UUID, bool) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, true
	}
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid "+key)
		return nil, false
	}
	return &id, true
}
