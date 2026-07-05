package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateCandidateApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create candidate application", err, "tenant context is required")
		return
	}
	h.createCandidateApplication(w, r, tenantID, "create candidate application")
}

func (h *Handler) ListCandidateApplications(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list candidate applications", err, "tenant context is required")
		return
	}
	h.listCandidateApplications(w, r, tenantID, "list candidate applications")
}

func (h *Handler) GetCandidateApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateApplicationRequestIDs(w, r, "get candidate application")
	if !ok {
		return
	}
	item, err := h.svc.GetCandidateApplication(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get candidate application", err, "candidate application not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateCandidateApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateApplicationRequestIDs(w, r, "update candidate application")
	if !ok {
		return
	}
	h.updateCandidateApplication(w, r, tenantID, id, "update candidate application")
}

func (h *Handler) MoveCandidateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateApplicationRequestIDs(w, r, "move candidate application status")
	if !ok {
		return
	}
	h.moveCandidateApplicationStatus(w, r, tenantID, id, "move candidate application status")
}

func (h *Handler) ListCandidateApplicationEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateApplicationRequestIDs(w, r, "list candidate application events")
	if !ok {
		return
	}
	items, err := h.svc.ListCandidateApplicationEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list candidate application events", err, "failed to list candidate application events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteCandidateApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateApplicationRequestIDs(w, r, "delete candidate application")
	if !ok {
		return
	}
	if err := h.svc.DeleteCandidateApplication(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete candidate application", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantCandidateApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant candidate application")
	if !ok {
		return
	}
	h.createCandidateApplication(w, r, tenantID, "create tenant candidate application")
}

func (h *Handler) ListTenantCandidateApplications(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant candidate applications")
	if !ok {
		return
	}
	h.listCandidateApplications(w, r, tenantID, "list tenant candidate applications")
}

func (h *Handler) GetTenantCandidateApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateApplicationID", "get tenant candidate application")
	if !ok {
		return
	}
	item, err := h.svc.GetCandidateApplication(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant candidate application", err, "candidate application not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantCandidateApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateApplicationID", "update tenant candidate application")
	if !ok {
		return
	}
	h.updateCandidateApplication(w, r, tenantID, id, "update tenant candidate application")
}

func (h *Handler) MoveTenantCandidateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateApplicationID", "move tenant candidate application status")
	if !ok {
		return
	}
	h.moveCandidateApplicationStatus(w, r, tenantID, id, "move tenant candidate application status")
}

func (h *Handler) ListTenantCandidateApplicationEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateApplicationID", "list tenant candidate application events")
	if !ok {
		return
	}
	items, err := h.svc.ListCandidateApplicationEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant candidate application events", err, "failed to list candidate application events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteTenantCandidateApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateApplicationID", "delete tenant candidate application")
	if !ok {
		return
	}
	if err := h.svc.DeleteCandidateApplication(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant candidate application", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createCandidateApplication(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CandidateApplicationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateCandidateApplication(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateCandidateApplication(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.CandidateApplicationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCandidateApplication(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) moveCandidateApplicationStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.CandidateApplicationMoveCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.MoveCandidateApplicationStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listCandidateApplications(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	page, err := h.svc.ListCandidateApplications(r.Context(), domain.CandidateApplicationFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), CandidateID: optionalUUIDQuery(r, "candidate_id"), JobPostingID: optionalUUIDQuery(r, "job_posting_id"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list candidate applications")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) candidateApplicationRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "candidateApplicationID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid candidate application id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
