package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateEmploymentType(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create employment type", err, "tenant context is required")
		return
	}
	var cmd ports.EmploymentTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode employment type create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateEmploymentType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create employment type", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListEmploymentTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list employment types", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListEmploymentTypes(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list employment types", err, "failed to list employment types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetEmploymentType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "employmentTypeID", "get employment type")
	if !ok {
		return
	}
	item, err := h.svc.GetEmploymentType(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get employment type", err, "employment type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateEmploymentType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "employmentTypeID", "update employment type")
	if !ok {
		return
	}
	var cmd ports.EmploymentTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode employment type update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEmploymentType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update employment type", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteEmploymentType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "employmentTypeID", "delete employment type")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmploymentType(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete employment type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantEmploymentType(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant employment type")
	if !ok {
		return
	}
	var cmd ports.EmploymentTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant employment type create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateEmploymentType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant employment type", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListTenantEmploymentTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant employment types")
	if !ok {
		return
	}
	items, err := h.svc.ListEmploymentTypes(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant employment types", err, "failed to list employment types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetTenantEmploymentType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employmentTypeID", "get tenant employment type")
	if !ok {
		return
	}
	item, err := h.svc.GetEmploymentType(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant employment type", err, "employment type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantEmploymentType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employmentTypeID", "update tenant employment type")
	if !ok {
		return
	}
	var cmd ports.EmploymentTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant employment type update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEmploymentType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant employment type", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantEmploymentType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employmentTypeID", "delete tenant employment type")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmploymentType(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant employment type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateMaritalStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create marital status", err, "tenant context is required")
		return
	}
	var cmd ports.MaritalStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode marital status create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateMaritalStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create marital status", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListMaritalStatuses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list marital statuses", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListMaritalStatuses(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list marital statuses", err, "failed to list marital statuses")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetMaritalStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "maritalStatusID", "get marital status")
	if !ok {
		return
	}
	item, err := h.svc.GetMaritalStatus(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get marital status", err, "marital status not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateMaritalStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "maritalStatusID", "update marital status")
	if !ok {
		return
	}
	var cmd ports.MaritalStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode marital status update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateMaritalStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update marital status", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteMaritalStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.lookupRequestIDs(w, r, "maritalStatusID", "delete marital status")
	if !ok {
		return
	}
	if err := h.svc.DeleteMaritalStatus(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete marital status", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantMaritalStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant marital status")
	if !ok {
		return
	}
	var cmd ports.MaritalStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant marital status create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateMaritalStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant marital status", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListTenantMaritalStatuses(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant marital statuses")
	if !ok {
		return
	}
	items, err := h.svc.ListMaritalStatuses(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant marital statuses", err, "failed to list marital statuses")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetTenantMaritalStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "maritalStatusID", "get tenant marital status")
	if !ok {
		return
	}
	item, err := h.svc.GetMaritalStatus(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant marital status", err, "marital status not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantMaritalStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "maritalStatusID", "update tenant marital status")
	if !ok {
		return
	}
	var cmd ports.MaritalStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant marital status update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateMaritalStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant marital status", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantMaritalStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "maritalStatusID", "delete tenant marital status")
	if !ok {
		return
	}
	if err := h.svc.DeleteMaritalStatus(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant marital status", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) lookupRequestIDs(w http.ResponseWriter, r *http.Request, param string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid lookup id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminLookupRequestIDs(w http.ResponseWriter, r *http.Request, param string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid lookup id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
