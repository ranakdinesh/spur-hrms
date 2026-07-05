package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListPulseSurveys(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listPulseSurveysForTenant(w, r, tenantID, "list pulse surveys")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list pulse surveys", err, "tenant context is required")
	}
}

func (h *Handler) CreatePulseSurvey(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createPulseSurveyForTenant(w, r, tenantID, "create pulse survey")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create pulse survey", err, "tenant context is required")
	}
}

func (h *Handler) GetPulseSurvey(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.pulseSurveyRequestIDs(w, r, "get pulse survey")
	if ok {
		h.getPulseSurveyForTenant(w, r, tenantID, id, "get pulse survey")
	}
}

func (h *Handler) UpdatePulseSurvey(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.pulseSurveyRequestIDs(w, r, "update pulse survey")
	if ok {
		h.updatePulseSurveyForTenant(w, r, tenantID, id, "update pulse survey")
	}
}

func (h *Handler) DeletePulseSurvey(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.pulseSurveyRequestIDs(w, r, "delete pulse survey")
	if !ok {
		return
	}
	if err := h.svc.DeletePulseSurvey(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete pulse survey", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdatePulseSurveyStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.pulseSurveyRequestIDs(w, r, "update pulse survey status")
	if !ok {
		return
	}
	var cmd ports.PulseSurveyStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode pulse survey status request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePulseSurveyStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update pulse survey status", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListPulseQuestions(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listPulseQuestionsForTenant(w, r, tenantID, "list pulse questions")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list pulse questions", err, "tenant context is required")
	}
}

func (h *Handler) CreatePulseQuestion(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createPulseQuestionForTenant(w, r, tenantID, "create pulse question")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create pulse question", err, "tenant context is required")
	}
}

func (h *Handler) UpdatePulseQuestion(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.pulseQuestionRequestIDs(w, r, "update pulse question")
	if ok {
		h.updatePulseQuestionForTenant(w, r, tenantID, id, "update pulse question")
	}
}

func (h *Handler) DeletePulseQuestion(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.pulseQuestionRequestIDs(w, r, "delete pulse question")
	if !ok {
		return
	}
	if err := h.svc.DeletePulseQuestion(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete pulse question", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListPulseResponses(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listPulseResponsesForTenant(w, r, tenantID, "list pulse responses")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list pulse responses", err, "tenant context is required")
	}
}

func (h *Handler) CreatePulseResponse(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createPulseResponseForTenant(w, r, tenantID, "create pulse response")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create pulse response", err, "tenant context is required")
	}
}

func (h *Handler) ListWellbeingScores(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listWellbeingScoresForTenant(w, r, tenantID, "list wellbeing scores")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list wellbeing scores", err, "tenant context is required")
	}
}

func (h *Handler) UpsertWellbeingScore(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.upsertWellbeingScoreForTenant(w, r, tenantID, "upsert wellbeing score")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "upsert wellbeing score", err, "tenant context is required")
	}
}

func (h *Handler) ListWellbeingAlerts(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listWellbeingAlertsForTenant(w, r, tenantID, "list wellbeing alerts")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list wellbeing alerts", err, "tenant context is required")
	}
}

func (h *Handler) UpdateWellbeingAlertStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.wellbeingAlertRequestIDs(w, r, "update wellbeing alert")
	if ok {
		h.updateWellbeingAlertStatusForTenant(w, r, tenantID, id, "update wellbeing alert")
	}
}

func (h *Handler) ListWellbeingAggregateRows(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listWellbeingAggregateForTenant(w, r, tenantID, "list wellbeing aggregates")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list wellbeing aggregates", err, "tenant context is required")
	}
}

func (h *Handler) ListTenantPulseSurveys(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant pulse surveys"); ok {
		h.listPulseSurveysForTenant(w, r, tenantID, "list tenant pulse surveys")
	}
}

func (h *Handler) CreateTenantPulseSurvey(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant pulse survey"); ok {
		h.createPulseSurveyForTenant(w, r, tenantID, "create tenant pulse survey")
	}
}

func (h *Handler) GetTenantPulseSurvey(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPulseSurveyIDs(w, r, "get tenant pulse survey")
	if ok {
		h.getPulseSurveyForTenant(w, r, tenantID, id, "get tenant pulse survey")
	}
}

func (h *Handler) UpdateTenantPulseSurvey(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPulseSurveyIDs(w, r, "update tenant pulse survey")
	if ok {
		h.updatePulseSurveyForTenant(w, r, tenantID, id, "update tenant pulse survey")
	}
}

func (h *Handler) DeleteTenantPulseSurvey(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPulseSurveyIDs(w, r, "delete tenant pulse survey")
	if !ok {
		return
	}
	if err := h.svc.DeletePulseSurvey(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant pulse survey", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTenantPulseSurveyStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPulseSurveyIDs(w, r, "update tenant pulse survey status")
	if !ok {
		return
	}
	var cmd ports.PulseSurveyStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant pulse survey status request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePulseSurveyStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant pulse survey status", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListTenantPulseQuestions(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant pulse questions"); ok {
		h.listPulseQuestionsForTenant(w, r, tenantID, "list tenant pulse questions")
	}
}

func (h *Handler) CreateTenantPulseQuestion(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant pulse question"); ok {
		h.createPulseQuestionForTenant(w, r, tenantID, "create tenant pulse question")
	}
}

func (h *Handler) UpdateTenantPulseQuestion(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPulseQuestionIDs(w, r, "update tenant pulse question")
	if ok {
		h.updatePulseQuestionForTenant(w, r, tenantID, id, "update tenant pulse question")
	}
}

func (h *Handler) DeleteTenantPulseQuestion(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPulseQuestionIDs(w, r, "delete tenant pulse question")
	if !ok {
		return
	}
	if err := h.svc.DeletePulseQuestion(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant pulse question", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantPulseResponses(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant pulse responses"); ok {
		h.listPulseResponsesForTenant(w, r, tenantID, "list tenant pulse responses")
	}
}

func (h *Handler) CreateTenantPulseResponse(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant pulse response"); ok {
		h.createPulseResponseForTenant(w, r, tenantID, "create tenant pulse response")
	}
}

func (h *Handler) ListTenantWellbeingScores(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant wellbeing scores"); ok {
		h.listWellbeingScoresForTenant(w, r, tenantID, "list tenant wellbeing scores")
	}
}

func (h *Handler) UpsertTenantWellbeingScore(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant wellbeing score"); ok {
		h.upsertWellbeingScoreForTenant(w, r, tenantID, "upsert tenant wellbeing score")
	}
}

func (h *Handler) ListTenantWellbeingAlerts(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant wellbeing alerts"); ok {
		h.listWellbeingAlertsForTenant(w, r, tenantID, "list tenant wellbeing alerts")
	}
}

func (h *Handler) UpdateTenantWellbeingAlertStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminWellbeingAlertIDs(w, r, "update tenant wellbeing alert")
	if ok {
		h.updateWellbeingAlertStatusForTenant(w, r, tenantID, id, "update tenant wellbeing alert")
	}
}

func (h *Handler) ListTenantWellbeingAggregateRows(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant wellbeing aggregates"); ok {
		h.listWellbeingAggregateForTenant(w, r, tenantID, "list tenant wellbeing aggregates")
	}
}

func (h *Handler) listPulseSurveysForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListPulseSurveys(r.Context(), domain.PulseSurveyFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), SurveyType: optionalStringQuery(r, "survey_type"), Search: optionalStringQuery(r, "search")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list pulse surveys")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPulseSurveyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PulseSurveyCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePulseSurvey(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getPulseSurveyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	item, err := h.svc.GetPulseSurvey(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updatePulseSurveyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.PulseSurveyCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePulseSurvey(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listPulseQuestionsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	surveyID, ok := h.optionalUUIDQuery(w, r, "survey_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListPulseQuestions(r.Context(), tenantID, surveyID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list pulse questions")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPulseQuestionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PulseQuestionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePulseQuestion(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updatePulseQuestionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.PulseQuestionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePulseQuestion(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listPulseResponsesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	surveyID, ok := h.optionalUUIDQuery(w, r, "survey_id", operation)
	if !ok {
		return
	}
	workerID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListPulseResponses(r.Context(), domain.PulseResponseFilter{TenantID: tenantID, SurveyID: surveyID, WorkerProfileID: workerID, RiskLevel: optionalStringQuery(r, "risk_level")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list pulse responses")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPulseResponseForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PulseResponseCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePulseResponse(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listWellbeingScoresForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListWellbeingScores(r.Context(), domain.WellbeingScoreFilter{TenantID: tenantID, WorkerProfileID: workerID, RiskLevel: optionalStringQuery(r, "risk_level")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list wellbeing scores")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) upsertWellbeingScoreForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.WellbeingScoreCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertWellbeingScore(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listWellbeingAlertsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListWellbeingAlerts(r.Context(), domain.WellbeingAlertFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), Severity: optionalStringQuery(r, "severity")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list wellbeing alerts")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateWellbeingAlertStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.WellbeingAlertStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWellbeingAlertStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listWellbeingAggregateForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	surveyID, ok := h.optionalUUIDQuery(w, r, "survey_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListWellbeingAggregateRows(r.Context(), tenantID, surveyID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list wellbeing aggregates")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) pulseSurveyRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "surveyID", operation, "invalid survey id")
	return tenantID, id, ok
}

func (h *Handler) pulseQuestionRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "questionID", operation, "invalid question id")
	return tenantID, id, ok
}

func (h *Handler) wellbeingAlertRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "alertID", operation, "invalid alert id")
	return tenantID, id, ok
}

func (h *Handler) superAdminPulseSurveyIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "surveyID", operation, "invalid survey id")
	return tenantID, id, ok
}

func (h *Handler) superAdminPulseQuestionIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "questionID", operation, "invalid question id")
	return tenantID, id, ok
}

func (h *Handler) superAdminWellbeingAlertIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "alertID", operation, "invalid alert id")
	return tenantID, id, ok
}
