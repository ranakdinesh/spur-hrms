package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create employee exit", err, "tenant context is required")
		return
	}
	h.createEmployeeExitForTenant(w, r, tenantID, "create employee exit")
}

func (h *Handler) ListEmployeeExits(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list employee exits", err, "tenant context is required")
		return
	}
	h.listEmployeeExitsForTenant(w, r, tenantID, "list employee exits")
}

func (h *Handler) GetEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get employee exit", err, "tenant context is required")
		return
	}
	h.getEmployeeExitForTenant(w, r, tenantID, "get employee exit")
}

func (h *Handler) ApproveEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "approve employee exit", err, "tenant context is required")
		return
	}
	h.employeeExitActionForTenant(w, r, tenantID, "approve employee exit", h.svc.ApproveEmployeeExit)
}

func (h *Handler) RejectEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "reject employee exit", err, "tenant context is required")
		return
	}
	h.employeeExitActionForTenant(w, r, tenantID, "reject employee exit", h.svc.RejectEmployeeExit)
}

func (h *Handler) CancelEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "cancel employee exit", err, "tenant context is required")
		return
	}
	h.employeeExitActionForTenant(w, r, tenantID, "cancel employee exit", h.svc.CancelEmployeeExit)
}

func (h *Handler) CompleteEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "complete employee exit", err, "tenant context is required")
		return
	}
	h.employeeExitActionForTenant(w, r, tenantID, "complete employee exit", h.svc.CompleteEmployeeExit)
}

func (h *Handler) UpdateEmployeeExitTaskStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "update employee exit task", err, "tenant context is required")
		return
	}
	h.updateEmployeeExitTaskForTenant(w, r, tenantID, "update employee exit task")
}

func (h *Handler) CreateTenantEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant employee exit")
	if !ok {
		return
	}
	h.createEmployeeExitForTenant(w, r, tenantID, "create tenant employee exit")
}

func (h *Handler) ListTenantEmployeeExits(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant employee exits")
	if !ok {
		return
	}
	h.listEmployeeExitsForTenant(w, r, tenantID, "list tenant employee exits")
}

func (h *Handler) GetTenantEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant employee exit")
	if !ok {
		return
	}
	h.getEmployeeExitForTenant(w, r, tenantID, "get tenant employee exit")
}

func (h *Handler) ApproveTenantEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "approve tenant employee exit")
	if !ok {
		return
	}
	h.employeeExitActionForTenant(w, r, tenantID, "approve tenant employee exit", h.svc.ApproveEmployeeExit)
}

func (h *Handler) RejectTenantEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "reject tenant employee exit")
	if !ok {
		return
	}
	h.employeeExitActionForTenant(w, r, tenantID, "reject tenant employee exit", h.svc.RejectEmployeeExit)
}

func (h *Handler) CancelTenantEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "cancel tenant employee exit")
	if !ok {
		return
	}
	h.employeeExitActionForTenant(w, r, tenantID, "cancel tenant employee exit", h.svc.CancelEmployeeExit)
}

func (h *Handler) CompleteTenantEmployeeExit(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "complete tenant employee exit")
	if !ok {
		return
	}
	h.employeeExitActionForTenant(w, r, tenantID, "complete tenant employee exit", h.svc.CompleteEmployeeExit)
}

func (h *Handler) UpdateTenantEmployeeExitTaskStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "update tenant employee exit task")
	if !ok {
		return
	}
	h.updateEmployeeExitTaskForTenant(w, r, tenantID, "update tenant employee exit task")
}

func (h *Handler) createEmployeeExitForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CreateEmployeeExitCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateEmployeeExit(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listEmployeeExitsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter := domain.EmployeeExitFilter{TenantID: tenantID, Limit: 25}
	q := r.URL.Query()
	if status := q.Get("status"); status != "" {
		filter.Status = &status
	}
	if search := q.Get("search"); search != "" {
		filter.Search = &search
	}
	if userID := q.Get("employee_user_id"); userID != "" {
		parsed, err := uuid.Parse(userID)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "parse employee user id", err, "invalid employee user id")
			return
		}
		filter.EmployeeUserID = &parsed
	}
	if limit, err := strconv.Atoi(q.Get("limit")); err == nil && limit > 0 {
		filter.Limit = int32(limit)
	}
	if offset, err := strconv.Atoi(q.Get("offset")); err == nil && offset >= 0 {
		filter.Offset = int32(offset)
	}
	page, err := h.svc.ListEmployeeExits(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list employee exits")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) getEmployeeExitForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	exitID, err := uuid.Parse(chi.URLParam(r, "exitID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse employee exit id", err, "invalid employee exit id")
		return
	}
	item, err := h.svc.GetEmployeeExit(r.Context(), tenantID, exitID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, "employee exit not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) employeeExitActionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string, action func(context.Context, ports.EmployeeExitActionCommand) (*domain.EmployeeExitRequest, error)) {
	exitID, err := uuid.Parse(chi.URLParam(r, "exitID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse employee exit id", err, "invalid employee exit id")
		return
	}
	var cmd ports.EmployeeExitActionCommand
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&cmd)
	}
	cmd.TenantID = tenantID
	cmd.ExitID = exitID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := action(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateEmployeeExitTaskForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse employee exit task id", err, "invalid employee exit task id")
		return
	}
	var cmd ports.EmployeeExitTaskStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.TaskID = taskID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEmployeeExitTaskStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
