package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateJobPosition(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create job position", err, "tenant context is required")
		return
	}
	h.createJobPosition(w, r, tenantID, "create job position")
}

func (h *Handler) ListJobPositions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list job positions", err, "tenant context is required")
		return
	}
	h.listJobPositions(w, r, tenantID, "list job positions")
}

func (h *Handler) GetJobPosition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobPositionRequestIDs(w, r, "get job position")
	if !ok {
		return
	}
	item, err := h.svc.GetJobPosition(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get job position", err, "job position not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateJobPosition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobPositionRequestIDs(w, r, "update job position")
	if !ok {
		return
	}
	h.updateJobPosition(w, r, tenantID, id, "update job position")
}

func (h *Handler) DeleteJobPosition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobPositionRequestIDs(w, r, "delete job position")
	if !ok {
		return
	}
	if err := h.svc.DeleteJobPosition(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete job position", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateJobPositionLocation(w http.ResponseWriter, r *http.Request) {
	tenantID, jobPositionID, ok := h.jobPositionRequestIDs(w, r, "create job position location")
	if !ok {
		return
	}
	h.createJobPositionLocation(w, r, tenantID, jobPositionID, "create job position location")
}

func (h *Handler) ListJobPositionLocations(w http.ResponseWriter, r *http.Request) {
	tenantID, jobPositionID, ok := h.jobPositionRequestIDs(w, r, "list job position locations")
	if !ok {
		return
	}
	items, err := h.svc.ListJobPositionLocations(r.Context(), tenantID, jobPositionID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list job position locations", err, "failed to list job position locations")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpdateJobPositionLocation(w http.ResponseWriter, r *http.Request) {
	tenantID, locationID, ok := h.jobPositionLocationRequestIDs(w, r, "update job position location")
	if !ok {
		return
	}
	h.updateJobPositionLocation(w, r, tenantID, locationID, "update job position location")
}

func (h *Handler) DeleteJobPositionLocation(w http.ResponseWriter, r *http.Request) {
	tenantID, locationID, ok := h.jobPositionLocationRequestIDs(w, r, "delete job position location")
	if !ok {
		return
	}
	if err := h.svc.DeleteJobPositionLocation(r.Context(), tenantID, locationID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete job position location", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantJobPosition(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant job position")
	if !ok {
		return
	}
	h.createJobPosition(w, r, tenantID, "create tenant job position")
}

func (h *Handler) ListTenantJobPositions(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant job positions")
	if !ok {
		return
	}
	h.listJobPositions(w, r, tenantID, "list tenant job positions")
}

func (h *Handler) GetTenantJobPosition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobPositionID", "get tenant job position")
	if !ok {
		return
	}
	item, err := h.svc.GetJobPosition(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant job position", err, "job position not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantJobPosition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobPositionID", "update tenant job position")
	if !ok {
		return
	}
	h.updateJobPosition(w, r, tenantID, id, "update tenant job position")
}

func (h *Handler) DeleteTenantJobPosition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobPositionID", "delete tenant job position")
	if !ok {
		return
	}
	if err := h.svc.DeleteJobPosition(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant job position", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantJobPositionLocation(w http.ResponseWriter, r *http.Request) {
	tenantID, jobPositionID, ok := h.superAdminLookupRequestIDs(w, r, "jobPositionID", "create tenant job position location")
	if !ok {
		return
	}
	h.createJobPositionLocation(w, r, tenantID, jobPositionID, "create tenant job position location")
}

func (h *Handler) ListTenantJobPositionLocations(w http.ResponseWriter, r *http.Request) {
	tenantID, jobPositionID, ok := h.superAdminLookupRequestIDs(w, r, "jobPositionID", "list tenant job position locations")
	if !ok {
		return
	}
	items, err := h.svc.ListJobPositionLocations(r.Context(), tenantID, jobPositionID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant job position locations", err, "failed to list job position locations")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpdateTenantJobPositionLocation(w http.ResponseWriter, r *http.Request) {
	tenantID, locationID, ok := h.superAdminLookupRequestIDs(w, r, "jobPositionLocationID", "update tenant job position location")
	if !ok {
		return
	}
	h.updateJobPositionLocation(w, r, tenantID, locationID, "update tenant job position location")
}

func (h *Handler) DeleteTenantJobPositionLocation(w http.ResponseWriter, r *http.Request) {
	tenantID, locationID, ok := h.superAdminLookupRequestIDs(w, r, "jobPositionLocationID", "delete tenant job position location")
	if !ok {
		return
	}
	if err := h.svc.DeleteJobPositionLocation(r.Context(), tenantID, locationID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant job position location", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createJobPosition(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.JobPositionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateJobPosition(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateJobPosition(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.JobPositionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateJobPosition(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listJobPositions(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	page, err := h.svc.ListJobPositions(r.Context(), domain.JobPositionFilter{TenantID: tenantID, DepartmentID: optionalUUIDQuery(r, "department_id"), EmploymentTypeID: optionalUUIDQuery(r, "employment_type_id"), WorkMode: optionalStringQuery(r, "work_mode"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 25), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list job positions")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) createJobPositionLocation(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, jobPositionID uuid.UUID, operation string) {
	var cmd ports.JobPositionLocationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.JobPositionID = jobPositionID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateJobPositionLocation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateJobPositionLocation(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.JobPositionLocationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateJobPositionLocation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) jobPositionRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "jobPositionID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid job position id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) jobPositionLocationRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "jobPositionLocationID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid job position location id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
