package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateDesignation(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create designation", err, "tenant context is required")
		return
	}
	var cmd ports.DesignationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode designation create request", err, "invalid request body")
		return
	}
	if !h.canSetDesignationAttendanceRequirement(r, cmd.AttendanceRequired) {
		h.respondError(w, r, http.StatusForbidden, "create designation attendance requirement", nil, "tenant admin permission is required to change attendance requirement")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	designation, err := h.svc.CreateDesignation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create designation", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, designation)
}

func (h *Handler) ListDesignations(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list designations", err, "tenant context is required")
		return
	}
	designations, err := h.svc.ListDesignations(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list designations", err, "failed to list designations")
		return
	}
	respondJSON(w, http.StatusOK, designations)
}

func (h *Handler) GetDesignation(w http.ResponseWriter, r *http.Request) {
	tenantID, designationID, ok := h.designationRequestIDs(w, r, "get designation")
	if !ok {
		return
	}
	designation, err := h.svc.GetDesignation(r.Context(), tenantID, designationID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get designation", err, "designation not found")
		return
	}
	respondJSON(w, http.StatusOK, designation)
}

func (h *Handler) UpdateDesignation(w http.ResponseWriter, r *http.Request) {
	tenantID, designationID, ok := h.designationRequestIDs(w, r, "update designation")
	if !ok {
		return
	}
	var cmd ports.DesignationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode designation update request", err, "invalid request body")
		return
	}
	if !h.canSetDesignationAttendanceRequirement(r, cmd.AttendanceRequired) {
		h.respondError(w, r, http.StatusForbidden, "update designation attendance requirement", nil, "tenant admin permission is required to change attendance requirement")
		return
	}
	cmd.ID = designationID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	designation, err := h.svc.UpdateDesignation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update designation", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, designation)
}

func (h *Handler) DeleteDesignation(w http.ResponseWriter, r *http.Request) {
	tenantID, designationID, ok := h.designationRequestIDs(w, r, "delete designation")
	if !ok {
		return
	}
	if err := h.svc.DeleteDesignation(r.Context(), tenantID, designationID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete designation", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantDesignation(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant designation")
	if !ok {
		return
	}
	var cmd ports.DesignationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant designation create request", err, "invalid request body")
		return
	}
	if !h.canSetDesignationAttendanceRequirement(r, cmd.AttendanceRequired) {
		h.respondError(w, r, http.StatusForbidden, "create tenant designation attendance requirement", nil, "tenant admin permission is required to change attendance requirement")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	designation, err := h.svc.CreateDesignation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant designation", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, designation)
}

func (h *Handler) ListTenantDesignations(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant designations")
	if !ok {
		return
	}
	designations, err := h.svc.ListDesignations(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant designations", err, "failed to list designations")
		return
	}
	respondJSON(w, http.StatusOK, designations)
}

func (h *Handler) GetTenantDesignation(w http.ResponseWriter, r *http.Request) {
	tenantID, designationID, ok := h.superAdminDesignationRequestIDs(w, r, "get tenant designation")
	if !ok {
		return
	}
	designation, err := h.svc.GetDesignation(r.Context(), tenantID, designationID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant designation", err, "designation not found")
		return
	}
	respondJSON(w, http.StatusOK, designation)
}

func (h *Handler) UpdateTenantDesignation(w http.ResponseWriter, r *http.Request) {
	tenantID, designationID, ok := h.superAdminDesignationRequestIDs(w, r, "update tenant designation")
	if !ok {
		return
	}
	var cmd ports.DesignationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant designation update request", err, "invalid request body")
		return
	}
	if !h.canSetDesignationAttendanceRequirement(r, cmd.AttendanceRequired) {
		h.respondError(w, r, http.StatusForbidden, "update tenant designation attendance requirement", nil, "tenant admin permission is required to change attendance requirement")
		return
	}
	cmd.ID = designationID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	designation, err := h.svc.UpdateDesignation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant designation", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, designation)
}

func (h *Handler) DeleteTenantDesignation(w http.ResponseWriter, r *http.Request) {
	tenantID, designationID, ok := h.superAdminDesignationRequestIDs(w, r, "delete tenant designation")
	if !ok {
		return
	}
	if err := h.svc.DeleteDesignation(r.Context(), tenantID, designationID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant designation", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) designationRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	designationID, err := uuid.Parse(chi.URLParam(r, "designationID"))
	if err != nil || designationID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid designation id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, designationID, true
}

func (h *Handler) canSetDesignationAttendanceRequirement(r *http.Request, value *bool) bool {
	if value == nil {
		return true
	}
	if h.isSuperAdmin != nil && h.isSuperAdmin(r.Context()) {
		return true
	}
	if h.rolesFromContext == nil {
		return false
	}
	for _, role := range h.rolesFromContext(r.Context()) {
		if normalizeDesignationRole(role) == "TENANT_ADMIN" {
			return true
		}
	}
	return false
}

func normalizeDesignationRole(role string) string {
	return strings.NewReplacer(" ", "_", "-", "_").Replace(strings.ToUpper(strings.TrimSpace(role)))
}

func (h *Handler) superAdminDesignationRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	designationID, err := uuid.Parse(chi.URLParam(r, "designationID"))
	if err != nil || designationID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid designation id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, designationID, true
}
