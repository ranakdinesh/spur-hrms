package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListRegisteredScheduledJobs(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.ListRegisteredScheduledJobs(r.Context())
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list registered scheduled jobs", err, "failed to list scheduled jobs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) RunScheduledJob(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "run scheduled job", err, "tenant context is required")
		return
	}
	jobKey := chi.URLParam(r, "jobKey")
	cmd, ok := h.decodeScheduledJobRequest(w, r, "run scheduled job")
	if !ok {
		return
	}
	cmd.JobKey = jobKey
	cmd.TenantID = &tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	result, err := h.svc.RunScheduledJob(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "run scheduled job", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) ListScheduledJobRuns(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list scheduled job runs", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListScheduledJobRuns(r.Context(), tenantID, chi.URLParam(r, "jobKey"), queryInt32(r, "limit", 25), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list scheduled job runs", err, "failed to list scheduled job runs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) RunTenantScheduledJob(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "run tenant scheduled job")
	if !ok {
		return
	}
	jobKey := chi.URLParam(r, "jobKey")
	cmd, ok := h.decodeScheduledJobRequest(w, r, "run tenant scheduled job")
	if !ok {
		return
	}
	cmd.JobKey = jobKey
	cmd.TenantID = &tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	result, err := h.svc.RunScheduledJob(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "run tenant scheduled job", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) ListTenantScheduledJobRuns(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant scheduled job runs")
	if !ok {
		return
	}
	items, err := h.svc.ListScheduledJobRuns(r.Context(), tenantID, chi.URLParam(r, "jobKey"), queryInt32(r, "limit", 25), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant scheduled job runs", err, "failed to list scheduled job runs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) RunCelebrationDailyJob(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "run celebration daily job", err, "tenant context is required")
		return
	}
	cmd, ok := h.decodeCelebrationDailyJobRequest(w, r, "run celebration daily job")
	if !ok {
		return
	}
	cmd.TenantID = &tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	result, err := h.svc.RunCelebrationDailyJob(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "run celebration daily job", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) ListCelebrationDailyJobRuns(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list celebration daily job runs", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListCelebrationDailyJobRuns(r.Context(), tenantID, queryInt32(r, "limit", 25), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list celebration daily job runs", err, "failed to list celebration job runs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) RunTenantCelebrationDailyJob(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "run tenant celebration daily job")
	if !ok {
		return
	}
	cmd, ok := h.decodeCelebrationDailyJobRequest(w, r, "run tenant celebration daily job")
	if !ok {
		return
	}
	cmd.TenantID = &tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	result, err := h.svc.RunCelebrationDailyJob(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "run tenant celebration daily job", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) ListTenantCelebrationDailyJobRuns(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant celebration daily job runs")
	if !ok {
		return
	}
	items, err := h.svc.ListCelebrationDailyJobRuns(r.Context(), tenantID, queryInt32(r, "limit", 25), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant celebration daily job runs", err, "failed to list celebration job runs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) decodeCelebrationDailyJobRequest(w http.ResponseWriter, r *http.Request, operation string) (ports.RunCelebrationDailyJobCommand, bool) {
	var cmd ports.RunCelebrationDailyJobCommand
	if r.Body == nil {
		return cmd, true
	}
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, http.ErrBodyNotAllowed) {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return cmd, false
	}
	return cmd, true
}

func (h *Handler) decodeScheduledJobRequest(w http.ResponseWriter, r *http.Request, operation string) (ports.RunScheduledJobCommand, bool) {
	var cmd ports.RunScheduledJobCommand
	if r.Body == nil {
		return cmd, true
	}
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, http.ErrBodyNotAllowed) {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return cmd, false
	}
	return cmd, true
}

func queryInt32(r *http.Request, key string, fallback int32) int32 {
	if r == nil {
		return fallback
	}
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return fallback
	}
	return int32(parsed)
}
