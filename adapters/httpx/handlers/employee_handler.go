package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
)

func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create employee", err, "tenant context is required")
		return
	}
	h.createEmployeeForTenant(w, r, tenantID, "create employee")
}

func (h *Handler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	if !h.requirePermission(w, r, "list employees", permissions.EmployeesList) {
		return
	}
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list employees", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListEmployees(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list employees", err, "failed to list employees")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetEmployeeProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get employee profile", err, "tenant context is required")
		return
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse employee id", err, "invalid employee id")
		return
	}
	item, err := h.svc.GetEmployeeProfile(r.Context(), tenantID, employeeID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get employee profile", err, "employee profile not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) GetMyEmployeeProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get my employee profile", err, "tenant context is required")
		return
	}
	actorID := h.actorIDFromRequest(r)
	if actorID == nil {
		h.respondError(w, r, http.StatusUnauthorized, "get my employee profile", domain.ErrInvalidEmployeeUserID, "user context is required")
		return
	}
	item, err := h.svc.GetEmployeeSelfProfile(r.Context(), tenantID, *actorID, actorID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get my employee profile", err, "employee profile not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "update employee", err, "tenant context is required")
		return
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse employee id", err, "invalid employee id")
		return
	}
	h.updateEmployeeForTenant(w, r, tenantID, employeeID, "update employee")
}

func (h *Handler) DeactivateEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "deactivate employee", err, "tenant context is required")
		return
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse employee id", err, "invalid employee id")
		return
	}
	if err := h.svc.DeactivateEmployee(r.Context(), tenantID, employeeID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "deactivate employee", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ResendEmployeeCredentials(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.employeeRequestIDs(w, r, "resend employee credentials")
	if !ok {
		return
	}
	item, err := h.svc.ResendEmployeeCredentials(r.Context(), ports.EmployeeCredentialActionCommand{TenantID: tenantID, EmployeeID: employeeID, ActorID: h.actorIDFromRequest(r)})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "resend employee credentials", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ResetEmployeeTemporaryPassword(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.employeeRequestIDs(w, r, "reset employee temporary password")
	if !ok {
		return
	}
	h.resetEmployeeTemporaryPasswordForTenant(w, r, tenantID, employeeID, "reset employee temporary password")
}

func (h *Handler) ListEmployeeCredentialEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.employeeRequestIDs(w, r, "list employee credential events")
	if !ok {
		return
	}
	h.listEmployeeCredentialEventsForTenant(w, r, tenantID, employeeID, "list employee credential events")
}

func (h *Handler) CreateTenantEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant employee")
	if !ok {
		return
	}
	h.createEmployeeForTenant(w, r, tenantID, "create tenant employee")
}

func (h *Handler) ListTenantEmployees(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant employees")
	if !ok {
		return
	}
	items, err := h.svc.ListEmployees(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant employees", err, "failed to list employees")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetTenantEmployeeProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant employee profile")
	if !ok {
		return
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant employee id", err, "invalid employee id")
		return
	}
	item, err := h.svc.GetEmployeeProfile(r.Context(), tenantID, employeeID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant employee profile", err, "employee profile not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "update tenant employee")
	if !ok {
		return
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant employee id", err, "invalid employee id")
		return
	}
	h.updateEmployeeForTenant(w, r, tenantID, employeeID, "update tenant employee")
}

func (h *Handler) DeactivateTenantEmployee(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "deactivate tenant employee")
	if !ok {
		return
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant employee id", err, "invalid employee id")
		return
	}
	if err := h.svc.DeactivateEmployee(r.Context(), tenantID, employeeID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "deactivate tenant employee", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ResendTenantEmployeeCredentials(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.tenantEmployeeRequestIDs(w, r, "resend tenant employee credentials")
	if !ok {
		return
	}
	item, err := h.svc.ResendEmployeeCredentials(r.Context(), ports.EmployeeCredentialActionCommand{TenantID: tenantID, EmployeeID: employeeID, ActorID: h.actorIDFromRequest(r)})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "resend tenant employee credentials", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ResetTenantEmployeeTemporaryPassword(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.tenantEmployeeRequestIDs(w, r, "reset tenant employee temporary password")
	if !ok {
		return
	}
	h.resetEmployeeTemporaryPasswordForTenant(w, r, tenantID, employeeID, "reset tenant employee temporary password")
}

func (h *Handler) ListTenantEmployeeCredentialEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.tenantEmployeeRequestIDs(w, r, "list tenant employee credential events")
	if !ok {
		return
	}
	h.listEmployeeCredentialEventsForTenant(w, r, tenantID, employeeID, "list tenant employee credential events")
}

func (h *Handler) createEmployeeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CreateEmployeeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateEmployee(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateEmployeeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, employeeID uuid.UUID, operation string) {
	var cmd ports.UpdateEmployeeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = employeeID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEmployee(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) resetEmployeeTemporaryPasswordForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, employeeID uuid.UUID, operation string) {
	var cmd ports.EmployeeCredentialActionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.EmployeeID = employeeID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ResetEmployeeTemporaryPassword(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listEmployeeCredentialEventsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, employeeID uuid.UUID, operation string) {
	limit := int32(20)
	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid limit")
			return
		}
		limit = int32(parsed)
	}
	items, err := h.svc.ListEmployeeCredentialEvents(r.Context(), tenantID, employeeID, limit)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list credential events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) employeeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid employee id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, employeeID, true
}

func (h *Handler) tenantEmployeeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid employee id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, employeeID, true
}
