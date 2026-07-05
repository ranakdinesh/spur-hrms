package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateAttendancePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create attendance policy", err, "tenant context is required")
		return
	}
	h.createAttendancePolicyForTenant(w, r, tenantID, "create attendance policy")
}
func (h *Handler) ListAttendancePolicies(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list attendance policies", err, "tenant context is required")
		return
	}
	h.listAttendancePoliciesForTenant(w, r, tenantID, "list attendance policies")
}
func (h *Handler) UpdateAttendancePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "update attendance policy", err, "tenant context is required")
		return
	}
	h.updateAttendancePolicyForTenant(w, r, tenantID, "update attendance policy")
}
func (h *Handler) DeleteAttendancePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete attendance policy", err, "tenant context is required")
		return
	}
	h.deleteAttendancePolicyForTenant(w, r, tenantID, "delete attendance policy")
}

func (h *Handler) CreateTenantAttendancePolicy(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant attendance policy"); ok {
		h.createAttendancePolicyForTenant(w, r, tenantID, "create tenant attendance policy")
	}
}
func (h *Handler) ListTenantAttendancePolicies(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendance policies"); ok {
		h.listAttendancePoliciesForTenant(w, r, tenantID, "list tenant attendance policies")
	}
}
func (h *Handler) UpdateTenantAttendancePolicy(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant attendance policy"); ok {
		h.updateAttendancePolicyForTenant(w, r, tenantID, "update tenant attendance policy")
	}
}
func (h *Handler) DeleteTenantAttendancePolicy(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "delete tenant attendance policy"); ok {
		h.deleteAttendancePolicyForTenant(w, r, tenantID, "delete tenant attendance policy")
	}
}

func (h *Handler) CreateAttendanceRoster(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create attendance roster", err, "tenant context is required")
		return
	}
	h.createAttendanceRosterForTenant(w, r, tenantID, "create attendance roster")
}
func (h *Handler) ListAttendanceRosters(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list attendance rosters", err, "tenant context is required")
		return
	}
	h.listAttendanceRostersForTenant(w, r, tenantID, "list attendance rosters")
}
func (h *Handler) UpdateAttendanceRoster(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "update attendance roster", err, "tenant context is required")
		return
	}
	h.updateAttendanceRosterForTenant(w, r, tenantID, "update attendance roster")
}
func (h *Handler) DeleteAttendanceRoster(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete attendance roster", err, "tenant context is required")
		return
	}
	h.deleteAttendanceRosterForTenant(w, r, tenantID, "delete attendance roster")
}

func (h *Handler) CreateTenantAttendanceRoster(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant attendance roster"); ok {
		h.createAttendanceRosterForTenant(w, r, tenantID, "create tenant attendance roster")
	}
}
func (h *Handler) ListTenantAttendanceRosters(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendance rosters"); ok {
		h.listAttendanceRostersForTenant(w, r, tenantID, "list tenant attendance rosters")
	}
}
func (h *Handler) UpdateTenantAttendanceRoster(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant attendance roster"); ok {
		h.updateAttendanceRosterForTenant(w, r, tenantID, "update tenant attendance roster")
	}
}
func (h *Handler) DeleteTenantAttendanceRoster(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "delete tenant attendance roster"); ok {
		h.deleteAttendanceRosterForTenant(w, r, tenantID, "delete tenant attendance roster")
	}
}

func (h *Handler) CreateAttendanceRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create attendance request", err, "tenant context is required")
		return
	}
	h.createAttendanceRequestForTenant(w, r, tenantID, "create attendance request")
}
func (h *Handler) ListAttendanceRequests(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list attendance requests", err, "tenant context is required")
		return
	}
	h.listAttendanceRequestsForTenant(w, r, tenantID, "list attendance requests")
}
func (h *Handler) ReviewAttendanceRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "review attendance request", err, "tenant context is required")
		return
	}
	h.reviewAttendanceRequestForTenant(w, r, tenantID, "review attendance request")
}

func (h *Handler) ListAttendanceExceptionEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, requestID, ok := h.tenantAndURLUUID(w, r, "requestID", "list attendance exception events")
	if !ok {
		return
	}
	h.listAttendanceExceptionEventsForTenant(w, r, tenantID, requestID, "list attendance exception events")
}

func (h *Handler) ListPayrollBlockingAttendanceRequests(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list payroll blocking attendance requests", err, "tenant context is required")
		return
	}
	h.listPayrollBlockingAttendanceRequestsForTenant(w, r, tenantID, "list payroll blocking attendance requests")
}

func (h *Handler) ListAttendanceExceptionWorkflows(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list attendance exception workflows", err, "tenant context is required")
		return
	}
	h.listAttendanceExceptionWorkflowsForTenant(w, r, tenantID, "list attendance exception workflows")
}

func (h *Handler) CreateAttendanceExceptionWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create attendance exception workflow", err, "tenant context is required")
		return
	}
	h.createAttendanceExceptionWorkflowForTenant(w, r, tenantID, "create attendance exception workflow")
}

func (h *Handler) UpdateAttendanceExceptionWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "update attendance exception workflow", err, "tenant context is required")
		return
	}
	h.updateAttendanceExceptionWorkflowForTenant(w, r, tenantID, "update attendance exception workflow")
}

func (h *Handler) DeleteAttendanceExceptionWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "workflowID", "delete attendance exception workflow")
	if !ok {
		return
	}
	if err := h.svc.DeleteAttendanceExceptionWorkflow(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete attendance exception workflow", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) CreateTenantAttendanceRequest(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant attendance request"); ok {
		h.createAttendanceRequestForTenant(w, r, tenantID, "create tenant attendance request")
	}
}
func (h *Handler) ListTenantAttendanceRequests(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendance requests"); ok {
		h.listAttendanceRequestsForTenant(w, r, tenantID, "list tenant attendance requests")
	}
}
func (h *Handler) ReviewTenantAttendanceRequest(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "review tenant attendance request"); ok {
		h.reviewAttendanceRequestForTenant(w, r, tenantID, "review tenant attendance request")
	}
}

func (h *Handler) ListTenantAttendanceExceptionEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendance exception events")
	if !ok {
		return
	}
	requestID, err := uuid.Parse(chi.URLParam(r, "requestID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "list tenant attendance exception events", err, "invalid request id")
		return
	}
	h.listAttendanceExceptionEventsForTenant(w, r, tenantID, requestID, "list tenant attendance exception events")
}

func (h *Handler) ListTenantPayrollBlockingAttendanceRequests(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant payroll blocking attendance requests"); ok {
		h.listPayrollBlockingAttendanceRequestsForTenant(w, r, tenantID, "list tenant payroll blocking attendance requests")
	}
}

func (h *Handler) ListTenantAttendanceExceptionWorkflows(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendance exception workflows"); ok {
		h.listAttendanceExceptionWorkflowsForTenant(w, r, tenantID, "list tenant attendance exception workflows")
	}
}

func (h *Handler) CreateTenantAttendanceExceptionWorkflow(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant attendance exception workflow"); ok {
		h.createAttendanceExceptionWorkflowForTenant(w, r, tenantID, "create tenant attendance exception workflow")
	}
}

func (h *Handler) UpdateTenantAttendanceExceptionWorkflow(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant attendance exception workflow"); ok {
		h.updateAttendanceExceptionWorkflowForTenant(w, r, tenantID, "update tenant attendance exception workflow")
	}
}

func (h *Handler) DeleteTenantAttendanceExceptionWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "delete tenant attendance exception workflow")
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant attendance exception workflow", err, "invalid workflow id")
		return
	}
	if err := h.svc.DeleteAttendanceExceptionWorkflow(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant attendance exception workflow", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createAttendancePolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendancePolicyCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAttendancePolicy(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) listAttendancePoliciesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAttendancePolicies(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list attendance policies")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) updateAttendancePolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendancePolicyCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "policyID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid policy id")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateAttendancePolicy(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) deleteAttendancePolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "policyID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid policy id")
		return
	}
	if err := h.svc.DeleteAttendancePolicy(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createAttendanceRosterForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceRosterCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAttendanceRoster(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) listAttendanceRostersForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	q := r.URL.Query()
	if raw := q.Get("user_id"); raw != "" {
		userID, err := uuid.Parse(raw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid user_id")
			return
		}
		items, err := h.svc.ListAttendanceRostersByUser(r.Context(), tenantID, userID, q.Get("start_date"), q.Get("end_date"))
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list rosters")
			return
		}
		respondJSON(w, http.StatusOK, items)
		return
	}
	items, err := h.svc.ListAttendanceRostersByDateRange(r.Context(), tenantID, q.Get("start_date"), q.Get("end_date"))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list rosters")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) updateAttendanceRosterForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceRosterCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "rosterID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid roster id")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateAttendanceRoster(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) deleteAttendanceRosterForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "rosterID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid roster id")
		return
	}
	if err := h.svc.DeleteAttendanceRoster(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createAttendanceRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceRequestCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	if cmd.UserID == uuid.Nil {
		if actor := h.actorIDFromRequest(r); actor != nil {
			cmd.UserID = *actor
		}
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAttendanceRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) listAttendanceRequestsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	q := r.URL.Query()
	if raw := q.Get("user_id"); raw != "" {
		userID, err := uuid.Parse(raw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid user_id")
			return
		}
		items, err := h.svc.ListAttendanceRequestsByUser(r.Context(), tenantID, userID)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list requests")
			return
		}
		respondJSON(w, http.StatusOK, items)
		return
	}
	items, err := h.svc.ListAttendanceRequestsByStatus(r.Context(), tenantID, q.Get("status"))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list requests")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) reviewAttendanceRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceReviewCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "requestID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request id")
		return
	}
	cmd.TenantID = tenantID
	cmd.RequestID = id
	if actor := h.actorIDFromRequest(r); actor != nil {
		cmd.ReviewerID = *actor
	}
	item, err := h.svc.ReviewAttendanceRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listAttendanceExceptionWorkflowsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAttendanceExceptionWorkflows(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list workflows")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createAttendanceExceptionWorkflowForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendanceExceptionWorkflowCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAttendanceExceptionWorkflow(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateAttendanceExceptionWorkflowForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid workflow id")
		return
	}
	var cmd ports.AttendanceExceptionWorkflowCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateAttendanceExceptionWorkflow(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listAttendanceExceptionEventsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, requestID uuid.UUID, operation string) {
	items, err := h.svc.ListAttendanceExceptionEvents(r.Context(), tenantID, requestID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listPayrollBlockingAttendanceRequestsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	q := r.URL.Query()
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")
	if startDate == "" && q.Get("month") != "" && q.Get("year") != "" {
		month, monthErr := strconv.Atoi(q.Get("month"))
		year, yearErr := strconv.Atoi(q.Get("year"))
		if monthErr != nil || yearErr != nil || month < 1 || month > 12 {
			h.respondError(w, r, http.StatusBadRequest, operation, domain.ErrInvalidPayrollPeriod, "invalid month or year")
			return
		}
		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		startDate = start.Format("2006-01-02")
		endDate = start.AddDate(0, 1, -1).Format("2006-01-02")
	}
	items, err := h.svc.ListPayrollBlockingAttendanceRequests(r.Context(), tenantID, startDate, endDate)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}
