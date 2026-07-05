package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateCandidate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create candidate", err, "tenant context is required")
		return
	}
	h.createCandidate(w, r, tenantID, "create candidate")
}

func (h *Handler) ListCandidates(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list candidates", err, "tenant context is required")
		return
	}
	h.listCandidates(w, r, tenantID, "list candidates")
}

func (h *Handler) GetCandidate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateRequestIDs(w, r, "get candidate")
	if !ok {
		return
	}
	item, err := h.svc.GetCandidate(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get candidate", err, "candidate not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateCandidate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateRequestIDs(w, r, "update candidate")
	if !ok {
		return
	}
	h.updateCandidate(w, r, tenantID, id, "update candidate")
}

func (h *Handler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateRequestIDs(w, r, "delete candidate")
	if !ok {
		return
	}
	if err := h.svc.DeleteCandidate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete candidate", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantCandidate(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant candidate")
	if !ok {
		return
	}
	h.createCandidate(w, r, tenantID, "create tenant candidate")
}

func (h *Handler) ListTenantCandidates(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant candidates")
	if !ok {
		return
	}
	h.listCandidates(w, r, tenantID, "list tenant candidates")
}

func (h *Handler) GetTenantCandidate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateID", "get tenant candidate")
	if !ok {
		return
	}
	item, err := h.svc.GetCandidate(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant candidate", err, "candidate not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantCandidate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateID", "update tenant candidate")
	if !ok {
		return
	}
	h.updateCandidate(w, r, tenantID, id, "update tenant candidate")
}

func (h *Handler) DeleteTenantCandidate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateID", "delete tenant candidate")
	if !ok {
		return
	}
	if err := h.svc.DeleteCandidate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant candidate", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createCandidate(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CandidateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateCandidate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateCandidate(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.CandidateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCandidate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listCandidates(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	page, err := h.svc.ListCandidates(r.Context(), domain.CandidateFilter{TenantID: tenantID, Search: optionalStringQuery(r, "search"), Source: optionalStringQuery(r, "source"), Gender: optionalStringQuery(r, "gender"), Limit: queryInt32(r, "limit", 25), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list candidates")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) candidateRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "candidateID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid candidate id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
