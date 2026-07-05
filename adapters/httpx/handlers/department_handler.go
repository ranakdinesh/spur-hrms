package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create department", err, "tenant context is required")
		return
	}
	var cmd ports.DepartmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode department create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	department, err := h.svc.CreateDepartment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create department", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, department)
}

func (h *Handler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list departments", err, "tenant context is required")
		return
	}
	departments, err := h.svc.ListDepartments(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list departments", err, "failed to list departments")
		return
	}
	respondJSON(w, http.StatusOK, departments)
}

func (h *Handler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	tenantID, departmentID, ok := h.departmentRequestIDs(w, r, "get department")
	if !ok {
		return
	}
	department, err := h.svc.GetDepartment(r.Context(), tenantID, departmentID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get department", err, "department not found")
		return
	}
	respondJSON(w, http.StatusOK, department)
}

func (h *Handler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	tenantID, departmentID, ok := h.departmentRequestIDs(w, r, "update department")
	if !ok {
		return
	}
	var cmd ports.DepartmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode department update request", err, "invalid request body")
		return
	}
	cmd.ID = departmentID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	department, err := h.svc.UpdateDepartment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update department", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, department)
}

func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	tenantID, departmentID, ok := h.departmentRequestIDs(w, r, "delete department")
	if !ok {
		return
	}
	if err := h.svc.DeleteDepartment(r.Context(), tenantID, departmentID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete department", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantDepartment(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant department")
	if !ok {
		return
	}
	var cmd ports.DepartmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant department create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	department, err := h.svc.CreateDepartment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant department", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, department)
}

func (h *Handler) ListTenantDepartments(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant departments")
	if !ok {
		return
	}
	departments, err := h.svc.ListDepartments(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant departments", err, "failed to list departments")
		return
	}
	respondJSON(w, http.StatusOK, departments)
}

func (h *Handler) GetTenantDepartment(w http.ResponseWriter, r *http.Request) {
	tenantID, departmentID, ok := h.superAdminDepartmentRequestIDs(w, r, "get tenant department")
	if !ok {
		return
	}
	department, err := h.svc.GetDepartment(r.Context(), tenantID, departmentID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant department", err, "department not found")
		return
	}
	respondJSON(w, http.StatusOK, department)
}

func (h *Handler) UpdateTenantDepartment(w http.ResponseWriter, r *http.Request) {
	tenantID, departmentID, ok := h.superAdminDepartmentRequestIDs(w, r, "update tenant department")
	if !ok {
		return
	}
	var cmd ports.DepartmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant department update request", err, "invalid request body")
		return
	}
	cmd.ID = departmentID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	department, err := h.svc.UpdateDepartment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant department", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, department)
}

func (h *Handler) DeleteTenantDepartment(w http.ResponseWriter, r *http.Request) {
	tenantID, departmentID, ok := h.superAdminDepartmentRequestIDs(w, r, "delete tenant department")
	if !ok {
		return
	}
	if err := h.svc.DeleteDepartment(r.Context(), tenantID, departmentID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant department", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) departmentRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	departmentID, err := uuid.Parse(chi.URLParam(r, "departmentID"))
	if err != nil || departmentID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid department id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, departmentID, true
}

func (h *Handler) superAdminDepartmentRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	departmentID, err := uuid.Parse(chi.URLParam(r, "departmentID"))
	if err != nil || departmentID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid department id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, departmentID, true
}
