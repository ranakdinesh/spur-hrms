package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateInterviewRound(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create interview round", err, "tenant context is required")
		return
	}
	h.createInterviewRound(w, r, tenantID, "create interview round")
}

func (h *Handler) ListInterviewRounds(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list interview rounds", err, "tenant context is required")
		return
	}
	h.listInterviewRounds(w, r, tenantID, "list interview rounds")
}

func (h *Handler) GetInterviewRound(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.interviewRoundRequestIDs(w, r, "get interview round")
	if !ok {
		return
	}
	item, err := h.svc.GetInterviewRound(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get interview round", err, "interview round not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateInterviewRound(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.interviewRoundRequestIDs(w, r, "update interview round")
	if !ok {
		return
	}
	h.updateInterviewRound(w, r, tenantID, id, "update interview round")
}

func (h *Handler) UpdateInterviewRoundStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.interviewRoundRequestIDs(w, r, "update interview round status")
	if !ok {
		return
	}
	h.updateInterviewRoundStatus(w, r, tenantID, id, "update interview round status")
}

func (h *Handler) DeleteInterviewRound(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.interviewRoundRequestIDs(w, r, "delete interview round")
	if !ok {
		return
	}
	if err := h.svc.DeleteInterviewRound(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete interview round", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantInterviewRound(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant interview round")
	if !ok {
		return
	}
	h.createInterviewRound(w, r, tenantID, "create tenant interview round")
}

func (h *Handler) ListTenantInterviewRounds(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant interview rounds")
	if !ok {
		return
	}
	h.listInterviewRounds(w, r, tenantID, "list tenant interview rounds")
}

func (h *Handler) GetTenantInterviewRound(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "interviewRoundID", "get tenant interview round")
	if !ok {
		return
	}
	item, err := h.svc.GetInterviewRound(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant interview round", err, "interview round not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantInterviewRound(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "interviewRoundID", "update tenant interview round")
	if !ok {
		return
	}
	h.updateInterviewRound(w, r, tenantID, id, "update tenant interview round")
}

func (h *Handler) UpdateTenantInterviewRoundStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "interviewRoundID", "update tenant interview round status")
	if !ok {
		return
	}
	h.updateInterviewRoundStatus(w, r, tenantID, id, "update tenant interview round status")
}

func (h *Handler) DeleteTenantInterviewRound(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "interviewRoundID", "delete tenant interview round")
	if !ok {
		return
	}
	if err := h.svc.DeleteInterviewRound(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant interview round", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createInterviewRound(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.InterviewRoundCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateInterviewRound(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateInterviewRound(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.InterviewRoundCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateInterviewRound(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateInterviewRoundStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.InterviewRoundStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateInterviewRoundStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listInterviewRounds(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	dateFrom, ok := optionalTimeQuery(w, r, "date_from", operation, h)
	if !ok {
		return
	}
	dateTo, ok := optionalTimeQuery(w, r, "date_to", operation, h)
	if !ok {
		return
	}
	page, err := h.svc.ListInterviewRounds(r.Context(), domain.InterviewRoundFilter{TenantID: tenantID, ApplicationID: optionalUUIDQuery(r, "application_id"), Status: optionalStringQuery(r, "status"), InterviewerUserID: optionalUUIDQuery(r, "interviewer_user_id"), DateFrom: dateFrom, DateTo: dateTo, Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list interview rounds")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) interviewRoundRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "interviewRoundID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid interview round id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func optionalTimeQuery(w http.ResponseWriter, r *http.Request, name string, operation string, h *Handler) (*time.Time, bool) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return nil, true
	}
	if value, err := time.Parse(time.RFC3339, raw); err == nil {
		return &value, true
	}
	if value, err := time.Parse("2006-01-02", raw); err == nil {
		return &value, true
	}
	h.respondError(w, r, http.StatusBadRequest, operation, domain.ErrInvalidInterviewDate, "invalid "+name)
	return nil, false
}
