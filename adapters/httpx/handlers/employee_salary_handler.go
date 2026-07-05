package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) AssignEmployeeSalary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "assign employee salary", err, "tenant context is required")
		return
	}
	h.assignEmployeeSalaryForTenant(w, r, tenantID, "assign employee salary")
}

func (h *Handler) ListEmployeeSalaries(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list employee salaries", err, "tenant context is required")
		return
	}
	h.listEmployeeSalariesForTenant(w, r, tenantID, "list employee salaries")
}

func (h *Handler) GetEmployeeSalary(w http.ResponseWriter, r *http.Request) {
	tenantID, salaryID, ok := h.employeeSalaryRequestIDs(w, r, "get employee salary")
	if !ok {
		return
	}
	item, err := h.svc.GetEmployeeSalary(r.Context(), tenantID, salaryID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get employee salary", err, "employee salary not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteEmployeeSalary(w http.ResponseWriter, r *http.Request) {
	tenantID, salaryID, ok := h.employeeSalaryRequestIDs(w, r, "delete employee salary")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeSalary(r.Context(), tenantID, salaryID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete employee salary", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CalculateEmployeeSalary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "calculate employee salary", err, "tenant context is required")
		return
	}
	h.calculateEmployeeSalaryForTenant(w, r, tenantID, "calculate employee salary")
}

func (h *Handler) AssignTenantEmployeeSalary(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "assign tenant employee salary")
	if !ok {
		return
	}
	h.assignEmployeeSalaryForTenant(w, r, tenantID, "assign tenant employee salary")
}

func (h *Handler) ListTenantEmployeeSalaries(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant employee salaries")
	if !ok {
		return
	}
	h.listEmployeeSalariesForTenant(w, r, tenantID, "list tenant employee salaries")
}

func (h *Handler) GetTenantEmployeeSalary(w http.ResponseWriter, r *http.Request) {
	tenantID, salaryID, ok := h.tenantEmployeeSalaryRequestIDs(w, r, "get tenant employee salary")
	if !ok {
		return
	}
	item, err := h.svc.GetEmployeeSalary(r.Context(), tenantID, salaryID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant employee salary", err, "employee salary not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantEmployeeSalary(w http.ResponseWriter, r *http.Request) {
	tenantID, salaryID, ok := h.tenantEmployeeSalaryRequestIDs(w, r, "delete tenant employee salary")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeSalary(r.Context(), tenantID, salaryID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant employee salary", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CalculateTenantEmployeeSalary(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "calculate tenant employee salary")
	if !ok {
		return
	}
	h.calculateEmployeeSalaryForTenant(w, r, tenantID, "calculate tenant employee salary")
}

func (h *Handler) assignEmployeeSalaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EmployeeSalaryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode employee salary request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.AssignEmployeeSalary(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listEmployeeSalariesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid user id")
		return
	}
	items, err := h.svc.ListEmployeeSalariesByUser(r.Context(), tenantID, userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list employee salaries")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) calculateEmployeeSalaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EmployeeSalaryCalculationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode employee salary calculation request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	item, err := h.svc.CalculateEmployeeSalary(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) employeeSalaryRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	salaryID, err := uuid.Parse(chi.URLParam(r, "salaryID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid employee salary id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, salaryID, true
}

func (h *Handler) tenantEmployeeSalaryRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	salaryID, err := uuid.Parse(chi.URLParam(r, "salaryID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid employee salary id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, salaryID, true
}
