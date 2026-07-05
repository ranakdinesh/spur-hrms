package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListPerformanceCheckIns(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listPerformanceCheckInsForTenant(w, r, tenantID, "list performance check-ins")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list performance check-ins", err, "tenant context is required")
	}
}

func (h *Handler) CreatePerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createPerformanceCheckInForTenant(w, r, tenantID, "create performance check-in")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create performance check-in", err, "tenant context is required")
	}
}

func (h *Handler) GetPerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.performanceCheckInRequestIDs(w, r, "get performance check-in")
	if ok {
		h.getPerformanceCheckInForTenant(w, r, tenantID, id, "get performance check-in")
	}
}

func (h *Handler) UpdatePerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.performanceCheckInRequestIDs(w, r, "update performance check-in")
	if ok {
		h.updatePerformanceCheckInForTenant(w, r, tenantID, id, "update performance check-in")
	}
}

func (h *Handler) DeletePerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.performanceCheckInRequestIDs(w, r, "delete performance check-in")
	if !ok {
		return
	}
	if err := h.svc.DeletePerformanceCheckIn(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete performance check-in", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SubmitPerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.performanceCheckInRequestIDs(w, r, "submit performance check-in")
	if !ok {
		return
	}
	item, err := h.svc.SubmitPerformanceCheckIn(r.Context(), ports.PerformanceStatusCommand{TenantID: tenantID, ID: id, ActorID: h.actorIDFromRequest(r)})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "submit performance check-in", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ReviewPerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.performanceCheckInRequestIDs(w, r, "review performance check-in")
	if !ok {
		return
	}
	h.reviewPerformanceCheckInForTenant(w, r, tenantID, id, "review performance check-in")
}

func (h *Handler) GetPerformanceCheckInSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.getPerformanceCheckInSummaryForTenant(w, r, tenantID, "get performance check-in summary")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "get performance check-in summary", err, "tenant context is required")
	}
}

func (h *Handler) ListFeedbackRequests(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listFeedbackRequestsForTenant(w, r, tenantID, "list feedback requests")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list feedback requests", err, "tenant context is required")
	}
}

func (h *Handler) CreateFeedbackRequest(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createFeedbackRequestForTenant(w, r, tenantID, "create feedback request")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create feedback request", err, "tenant context is required")
	}
}

func (h *Handler) GetFeedbackRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.feedbackRequestIDs(w, r, "get feedback request")
	if ok {
		h.getFeedbackRequestForTenant(w, r, tenantID, id, "get feedback request")
	}
}

func (h *Handler) UpdateFeedbackRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.feedbackRequestIDs(w, r, "update feedback request")
	if ok {
		h.updateFeedbackRequestForTenant(w, r, tenantID, id, "update feedback request")
	}
}

func (h *Handler) UpdateFeedbackRequestStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.feedbackRequestIDs(w, r, "update feedback request status")
	if !ok {
		return
	}
	var cmd ports.FeedbackStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode feedback status request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateFeedbackRequestStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update feedback request status", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListFeedbackResponses(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listFeedbackResponsesForTenant(w, r, tenantID, "list feedback responses")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list feedback responses", err, "tenant context is required")
	}
}

func (h *Handler) CreateFeedbackResponse(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createFeedbackResponseForTenant(w, r, tenantID, "create feedback response")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create feedback response", err, "tenant context is required")
	}
}

func (h *Handler) ListPerformanceTimelineEvents(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listPerformanceTimelineForTenant(w, r, tenantID, "list performance timeline")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list performance timeline", err, "tenant context is required")
	}
}

func (h *Handler) ListPerformanceCalibrationRows(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listPerformanceCalibrationForTenant(w, r, tenantID, "list performance calibration")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list performance calibration", err, "tenant context is required")
	}
}

func (h *Handler) ListTenantPerformanceCheckIns(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant performance check-ins"); ok {
		h.listPerformanceCheckInsForTenant(w, r, tenantID, "list tenant performance check-ins")
	}
}

func (h *Handler) CreateTenantPerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant performance check-in"); ok {
		h.createPerformanceCheckInForTenant(w, r, tenantID, "create tenant performance check-in")
	}
}

func (h *Handler) GetTenantPerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPerformanceCheckInIDs(w, r, "get tenant performance check-in")
	if ok {
		h.getPerformanceCheckInForTenant(w, r, tenantID, id, "get tenant performance check-in")
	}
}

func (h *Handler) UpdateTenantPerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPerformanceCheckInIDs(w, r, "update tenant performance check-in")
	if ok {
		h.updatePerformanceCheckInForTenant(w, r, tenantID, id, "update tenant performance check-in")
	}
}

func (h *Handler) DeleteTenantPerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPerformanceCheckInIDs(w, r, "delete tenant performance check-in")
	if !ok {
		return
	}
	if err := h.svc.DeletePerformanceCheckIn(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant performance check-in", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SubmitTenantPerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPerformanceCheckInIDs(w, r, "submit tenant performance check-in")
	if !ok {
		return
	}
	item, err := h.svc.SubmitPerformanceCheckIn(r.Context(), ports.PerformanceStatusCommand{TenantID: tenantID, ID: id, ActorID: h.actorIDFromRequest(r)})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "submit tenant performance check-in", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ReviewTenantPerformanceCheckIn(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPerformanceCheckInIDs(w, r, "review tenant performance check-in")
	if ok {
		h.reviewPerformanceCheckInForTenant(w, r, tenantID, id, "review tenant performance check-in")
	}
}

func (h *Handler) GetTenantPerformanceCheckInSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant performance check-in summary"); ok {
		h.getPerformanceCheckInSummaryForTenant(w, r, tenantID, "get tenant performance check-in summary")
	}
}

func (h *Handler) ListTenantFeedbackRequests(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant feedback requests"); ok {
		h.listFeedbackRequestsForTenant(w, r, tenantID, "list tenant feedback requests")
	}
}

func (h *Handler) CreateTenantFeedbackRequest(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant feedback request"); ok {
		h.createFeedbackRequestForTenant(w, r, tenantID, "create tenant feedback request")
	}
}

func (h *Handler) GetTenantFeedbackRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFeedbackRequestIDs(w, r, "get tenant feedback request")
	if ok {
		h.getFeedbackRequestForTenant(w, r, tenantID, id, "get tenant feedback request")
	}
}

func (h *Handler) UpdateTenantFeedbackRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFeedbackRequestIDs(w, r, "update tenant feedback request")
	if ok {
		h.updateFeedbackRequestForTenant(w, r, tenantID, id, "update tenant feedback request")
	}
}

func (h *Handler) UpdateTenantFeedbackRequestStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFeedbackRequestIDs(w, r, "update tenant feedback request status")
	if !ok {
		return
	}
	var cmd ports.FeedbackStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant feedback status request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateFeedbackRequestStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant feedback request status", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListTenantFeedbackResponses(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant feedback responses"); ok {
		h.listFeedbackResponsesForTenant(w, r, tenantID, "list tenant feedback responses")
	}
}

func (h *Handler) CreateTenantFeedbackResponse(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant feedback response"); ok {
		h.createFeedbackResponseForTenant(w, r, tenantID, "create tenant feedback response")
	}
}

func (h *Handler) ListTenantPerformanceTimelineEvents(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant performance timeline"); ok {
		h.listPerformanceTimelineForTenant(w, r, tenantID, "list tenant performance timeline")
	}
}

func (h *Handler) ListTenantPerformanceCalibrationRows(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant performance calibration"); ok {
		h.listPerformanceCalibrationForTenant(w, r, tenantID, "list tenant performance calibration")
	}
}

func (h *Handler) listPerformanceCheckInsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	reviewerID, ok := h.optionalUUIDQuery(w, r, "reviewer_worker_profile_id", operation)
	if !ok {
		return
	}
	cycleID, ok := h.optionalUUIDQuery(w, r, "cycle_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListPerformanceCheckIns(r.Context(), domain.PerformanceCheckInFilter{TenantID: tenantID, WorkerProfileID: workerID, ReviewerWorkerProfileID: reviewerID, CycleID: cycleID, Status: optionalStringQuery(r, "status"), Mood: optionalStringQuery(r, "mood")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list performance check-ins")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPerformanceCheckInForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PerformanceCheckInCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePerformanceCheckIn(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getPerformanceCheckInForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	item, err := h.svc.GetPerformanceCheckIn(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updatePerformanceCheckInForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.PerformanceCheckInCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePerformanceCheckIn(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) reviewPerformanceCheckInForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.PerformanceCheckInReviewCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ReviewPerformanceCheckIn(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) getPerformanceCheckInSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	cycleID, ok := h.optionalUUIDQuery(w, r, "cycle_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.GetPerformanceCheckInSummary(r.Context(), tenantID, cycleID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to get performance check-in summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listFeedbackRequestsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	subjectID, ok := h.optionalUUIDQuery(w, r, "subject_worker_profile_id", operation)
	if !ok {
		return
	}
	requesterID, ok := h.optionalUUIDQuery(w, r, "requester_worker_profile_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListFeedbackRequests(r.Context(), domain.FeedbackRequestFilter{TenantID: tenantID, SubjectWorkerProfileID: subjectID, RequesterWorkerProfileID: requesterID, Status: optionalStringQuery(r, "status"), FeedbackType: optionalStringQuery(r, "feedback_type")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list feedback requests")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createFeedbackRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.FeedbackRequestCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateFeedbackRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getFeedbackRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	item, err := h.svc.GetFeedbackRequest(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateFeedbackRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.FeedbackRequestCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateFeedbackRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listFeedbackResponsesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	requestID, ok := h.optionalUUIDQuery(w, r, "request_id", operation)
	if !ok {
		return
	}
	subjectID, ok := h.optionalUUIDQuery(w, r, "subject_worker_profile_id", operation)
	if !ok {
		return
	}
	respondentID, ok := h.optionalUUIDQuery(w, r, "respondent_worker_profile_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListFeedbackResponses(r.Context(), domain.FeedbackResponseFilter{TenantID: tenantID, RequestID: requestID, SubjectWorkerProfileID: subjectID, RespondentWorkerProfileID: respondentID})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list feedback responses")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createFeedbackResponseForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.FeedbackResponseCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateFeedbackResponse(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listPerformanceTimelineForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListPerformanceTimelineEvents(r.Context(), domain.PerformanceTimelineFilter{TenantID: tenantID, WorkerProfileID: workerID, EventType: optionalStringQuery(r, "event_type")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list performance timeline")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listPerformanceCalibrationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	cycleID, ok := h.optionalUUIDQuery(w, r, "cycle_id", operation)
	if !ok {
		return
	}
	workerID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListPerformanceCalibrationRows(r.Context(), domain.PerformanceCalibrationFilter{TenantID: tenantID, CycleID: cycleID, WorkerProfileID: workerID})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list performance calibration")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) performanceCheckInRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "checkInID", operation, "invalid check-in id")
	return tenantID, id, ok
}

func (h *Handler) feedbackRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "feedbackRequestID", operation, "invalid feedback request id")
	return tenantID, id, ok
}

func (h *Handler) superAdminPerformanceCheckInIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "checkInID", operation, "invalid check-in id")
	return tenantID, id, ok
}

func (h *Handler) superAdminFeedbackRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "feedbackRequestID", operation, "invalid feedback request id")
	return tenantID, id, ok
}
