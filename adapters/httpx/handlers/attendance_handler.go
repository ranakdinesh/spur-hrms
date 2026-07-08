package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
)

func (h *Handler) PunchAttendance(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "punch attendance", err, "tenant context is required")
		return
	}
	h.punchAttendanceForTenant(w, r, tenantID, "punch attendance")
}

func (h *Handler) ListAttendances(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list attendances", err, "tenant context is required")
		return
	}
	h.listAttendancesForTenant(w, r, tenantID, "list attendances")
}

func (h *Handler) ListAttendanceDailyStatuses(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list attendance daily statuses", err, "tenant context is required")
		return
	}
	h.listAttendanceDailyStatusesForTenant(w, r, tenantID, "list attendance daily statuses")
}

func (h *Handler) GetAttendanceStatusSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get attendance status summary", err, "tenant context is required")
		return
	}
	h.getAttendanceStatusSummaryForTenant(w, r, tenantID, "get attendance status summary")
}

func (h *Handler) PunchTenantAttendance(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "punch tenant attendance")
	if ok {
		h.punchAttendanceForTenant(w, r, tenantID, "punch tenant attendance")
	}
}

func (h *Handler) ListTenantAttendances(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendances")
	if ok {
		h.listAttendancesForTenant(w, r, tenantID, "list tenant attendances")
	}
}

func (h *Handler) ListTenantAttendanceDailyStatuses(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant attendance daily statuses")
	if ok {
		h.listAttendanceDailyStatusesForTenant(w, r, tenantID, "list tenant attendance daily statuses")
	}
}

func (h *Handler) GetTenantAttendanceStatusSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant attendance status summary")
	if ok {
		h.getAttendanceStatusSummaryForTenant(w, r, tenantID, "get tenant attendance status summary")
	}
}

func (h *Handler) punchAttendanceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AttendancePunchCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	if cmd.UserID == uuid.Nil {
		if actorID := h.actorIDFromRequest(r); actorID != nil {
			cmd.UserID = *actorID
		}
	}
	if !h.requireOwnUserOrPermission(w, r, operation, cmd.UserID,
		[]string{permissions.AttendanceSelfPunch, permissions.AttendanceCheckIn, permissions.AttendanceCheckOut},
		[]string{permissions.AttendanceOperationsManage},
	) {
		return
	}
	if cmd.Action == "" {
		h.respondError(w, r, http.StatusBadRequest, operation, nil, "attendance action is required")
		return
	}
	if cmd.IPAddress == nil {
		cmd.IPAddress = stringPtr(clientIP(r))
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.PunchAttendance(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listAttendancesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	query := r.URL.Query()
	userRaw := query.Get("user_id")
	var userID uuid.UUID
	if userRaw != "" {
		parsed, err := uuid.Parse(userRaw)
		if err != nil || parsed == uuid.Nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid user_id")
			return
		}
		userID = parsed
	} else if actorID := h.actorIDFromRequest(r); actorID != nil {
		userID = *actorID
	}
	if !h.requireOwnUserOrPermission(w, r, operation, userID,
		[]string{permissions.AttendanceSelfView, permissions.AttendanceList, permissions.AttendanceView},
		[]string{permissions.AttendanceOperationsView},
	) {
		return
	}
	if date := query.Get("date"); date != "" {
		items, err := h.svc.ListAttendancesByUserDate(r.Context(), tenantID, userID, date)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list attendances")
			return
		}
		respondJSON(w, http.StatusOK, items)
		return
	}
	items, err := h.svc.ListAttendancesByUser(r.Context(), tenantID, userID, query.Get("start_date"), query.Get("end_date"))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list attendances")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listAttendanceDailyStatusesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	query, ok := h.attendanceStatusQuery(w, r, tenantID, operation)
	if !ok {
		return
	}
	if query.UserID == nil {
		if actorID := h.actorIDFromRequest(r); actorID != nil && !h.hasAnyPermission(r, permissions.AttendanceOperationsView) && !h.isSuperAdminRequest(r) {
			query.UserID = actorID
		}
	}
	if query.UserID != nil && !h.requireOwnUserOrPermission(w, r, operation, *query.UserID,
		[]string{permissions.AttendanceSelfView, permissions.AttendanceList, permissions.AttendanceView},
		[]string{permissions.AttendanceOperationsView},
	) {
		return
	}
	if query.UserID == nil && !h.requirePermission(w, r, operation, permissions.AttendanceOperationsView) {
		return
	}
	items, err := h.svc.ListAttendanceDailyStatuses(r.Context(), query)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to resolve attendance status")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getAttendanceStatusSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	query, ok := h.attendanceStatusQuery(w, r, tenantID, operation)
	if !ok {
		return
	}
	if query.UserID == nil {
		if actorID := h.actorIDFromRequest(r); actorID != nil && !h.hasAnyPermission(r, permissions.AttendanceOperationsView) && !h.isSuperAdminRequest(r) {
			query.UserID = actorID
		}
	}
	if query.UserID != nil && !h.requireOwnUserOrPermission(w, r, operation, *query.UserID,
		[]string{permissions.AttendanceSelfView, permissions.AttendanceList, permissions.AttendanceView},
		[]string{permissions.AttendanceOperationsView},
	) {
		return
	}
	if query.UserID == nil && !h.requirePermission(w, r, operation, permissions.AttendanceOperationsView) {
		return
	}
	summary, err := h.svc.GetAttendanceStatusSummary(r.Context(), query)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to resolve attendance summary")
		return
	}
	respondJSON(w, http.StatusOK, summary)
}

func (h *Handler) attendanceStatusQuery(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (ports.AttendanceStatusQuery, bool) {
	query := r.URL.Query()
	result := ports.AttendanceStatusQuery{TenantID: tenantID, Date: query.Get("date")}
	if userRaw := query.Get("user_id"); userRaw != "" {
		parsed, err := uuid.Parse(userRaw)
		if err != nil || parsed == uuid.Nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid user_id")
			return result, false
		}
		result.UserID = &parsed
	}
	return result, true
}

func clientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		forwarded = strings.TrimSpace(strings.Split(forwarded, ",")[0])
		if host, _, err := net.SplitHostPort(forwarded); err == nil {
			return host
		}
		return forwarded
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func stringPtr(value string) *string {
	return &value
}
