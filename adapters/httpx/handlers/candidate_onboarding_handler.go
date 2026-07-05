package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListCandidateOnboardings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list candidate onboardings", err, "tenant context is required")
		return
	}
	h.listCandidateOnboardings(w, r, tenantID, "list candidate onboardings")
}

func (h *Handler) StartCandidateOnboarding(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "start candidate onboarding", err, "tenant context is required")
		return
	}
	h.startCandidateOnboarding(w, r, tenantID, "start candidate onboarding")
}

func (h *Handler) GetCandidateOnboarding(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateOnboardingRequestIDs(w, r, "get candidate onboarding")
	if !ok {
		return
	}
	item, err := h.svc.GetCandidateOnboarding(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get candidate onboarding", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteCandidateOnboarding(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateOnboardingRequestIDs(w, r, "delete candidate onboarding")
	if !ok {
		return
	}
	if err := h.svc.DeleteCandidateOnboarding(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete candidate onboarding", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateCandidateOnboardingTaskStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, taskID, ok := h.candidateOnboardingTaskRequestIDs(w, r, "update candidate onboarding task status")
	if !ok {
		return
	}
	h.updateCandidateOnboardingTaskStatus(w, r, tenantID, taskID, "update candidate onboarding task status")
}

func (h *Handler) ListCandidateOnboardingEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.candidateOnboardingRequestIDs(w, r, "list candidate onboarding events")
	if !ok {
		return
	}
	events, err := h.svc.ListCandidateOnboardingEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "list candidate onboarding events", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, events)
}

func (h *Handler) ListTenantCandidateOnboardings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant candidate onboardings")
	if !ok {
		return
	}
	h.listCandidateOnboardings(w, r, tenantID, "list tenant candidate onboardings")
}

func (h *Handler) StartTenantCandidateOnboarding(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "start tenant candidate onboarding")
	if !ok {
		return
	}
	h.startCandidateOnboarding(w, r, tenantID, "start tenant candidate onboarding")
}

func (h *Handler) GetTenantCandidateOnboarding(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateOnboardingID", "get tenant candidate onboarding")
	if !ok {
		return
	}
	item, err := h.svc.GetCandidateOnboarding(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant candidate onboarding", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantCandidateOnboarding(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateOnboardingID", "delete tenant candidate onboarding")
	if !ok {
		return
	}
	if err := h.svc.DeleteCandidateOnboarding(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant candidate onboarding", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTenantCandidateOnboardingTaskStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, taskID, ok := h.superAdminLookupRequestIDs(w, r, "candidateOnboardingTaskID", "update tenant candidate onboarding task status")
	if !ok {
		return
	}
	h.updateCandidateOnboardingTaskStatus(w, r, tenantID, taskID, "update tenant candidate onboarding task status")
}

func (h *Handler) ListTenantCandidateOnboardingEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "candidateOnboardingID", "list tenant candidate onboarding events")
	if !ok {
		return
	}
	events, err := h.svc.ListCandidateOnboardingEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "list tenant candidate onboarding events", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, events)
}

func (h *Handler) listCandidateOnboardings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	page, err := h.svc.ListCandidateOnboardings(r.Context(), domain.CandidateOnboardingFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 25), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list candidate onboardings")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) startCandidateOnboarding(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.StartCandidateOnboardingCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.StartCandidateOnboarding(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateCandidateOnboardingTaskStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, taskID uuid.UUID, operation string) {
	var cmd ports.CandidateOnboardingTaskStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.TaskID = taskID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCandidateOnboardingTaskStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) candidateOnboardingRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "candidateOnboardingID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid candidate onboarding id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) candidateOnboardingTaskRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "candidateOnboardingTaskID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid candidate onboarding task id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
