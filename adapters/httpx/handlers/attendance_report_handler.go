package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetAttendanceReport(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get attendance report", err, "tenant context is required")
		return
	}
	h.getAttendanceReportForTenant(w, r, tenantID, "get attendance report")
}

func (h *Handler) GetTenantAttendanceReport(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant attendance report"); ok {
		h.getAttendanceReportForTenant(w, r, tenantID, "get tenant attendance report")
	}
}

func (h *Handler) getAttendanceReportForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	query, ok := h.attendanceReportQuery(w, r, tenantID, operation)
	if !ok {
		return
	}
	report, err := h.svc.GetAttendanceReport(r.Context(), query)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to build attendance report")
		return
	}
	respondJSON(w, http.StatusOK, report)
}

func (h *Handler) attendanceReportQuery(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (ports.AttendanceReportQuery, bool) {
	userID, ok := h.optionalUUIDQuery(w, r, "user_id", operation)
	if !ok {
		return ports.AttendanceReportQuery{}, false
	}
	departmentID, ok := h.optionalUUIDQuery(w, r, "department_id", operation)
	if !ok {
		return ports.AttendanceReportQuery{}, false
	}
	branchID, ok := h.optionalUUIDQuery(w, r, "branch_id", operation)
	if !ok {
		return ports.AttendanceReportQuery{}, false
	}
	values := r.URL.Query()
	return ports.AttendanceReportQuery{
		TenantID:     tenantID,
		UserID:       userID,
		DepartmentID: departmentID,
		BranchID:     branchID,
		StartDate:    values.Get("start_date"),
		EndDate:      values.Get("end_date"),
	}, true
}
