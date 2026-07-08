package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
)

func (h *Handler) CreateOvertimeRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create overtime request", err, "tenant context is required")
		return
	}
	h.createOvertimeRequestForTenant(w, r, tenantID, "create overtime request")
}

func (h *Handler) ListOvertimeRequests(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list overtime requests", err, "tenant context is required")
		return
	}
	h.listOvertimeRequestsForTenant(w, r, tenantID, "list overtime requests")
}

func (h *Handler) ApproveOvertimeRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, requestID, ok := h.overtimeRequestIDs(w, r, "approve overtime request")
	if !ok {
		return
	}
	h.reviewOvertimeRequestForTenant(w, r, tenantID, requestID, "approve overtime request", "approved")
}

func (h *Handler) RejectOvertimeRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, requestID, ok := h.overtimeRequestIDs(w, r, "reject overtime request")
	if !ok {
		return
	}
	h.reviewOvertimeRequestForTenant(w, r, tenantID, requestID, "reject overtime request", "rejected")
}

func (h *Handler) CreateTenantOvertimeRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant overtime request")
	if ok {
		h.createOvertimeRequestForTenant(w, r, tenantID, "create tenant overtime request")
	}
}

func (h *Handler) ListTenantOvertimeRequests(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant overtime requests")
	if ok {
		h.listOvertimeRequestsForTenant(w, r, tenantID, "list tenant overtime requests")
	}
}

func (h *Handler) ApproveTenantOvertimeRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, requestID, ok := h.superAdminOvertimeRequestIDs(w, r, "approve tenant overtime request")
	if !ok {
		return
	}
	h.reviewOvertimeRequestForTenant(w, r, tenantID, requestID, "approve tenant overtime request", "approved")
}

func (h *Handler) RejectTenantOvertimeRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, requestID, ok := h.superAdminOvertimeRequestIDs(w, r, "reject tenant overtime request")
	if !ok {
		return
	}
	h.reviewOvertimeRequestForTenant(w, r, tenantID, requestID, "reject tenant overtime request", "rejected")
}

func (h *Handler) createOvertimeRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.OvertimeRequestCommand
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
		[]string{permissions.AttendanceSelfOvertimeRequest},
		[]string{permissions.AttendanceOperationsManage, permissions.OvertimeApprove},
	) {
		return
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateOvertimeRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listOvertimeRequestsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	query := r.URL.Query()
	if exportStatus := query.Get("payroll_export_status"); exportStatus != "" {
		if !h.isSuperAdminRequest(r) && !h.hasAnyPermission(r, permissions.OvertimeExport, permissions.PayrollSalarySheetExport) {
			h.respondError(w, r, http.StatusForbidden, operation, nil, "permission required")
			return
		}
		items, err := h.svc.ListOvertimeRequestsByPayrollExportStatus(r.Context(), tenantID, exportStatus)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list overtime requests")
			return
		}
		respondJSON(w, http.StatusOK, items)
		return
	}
	if status := query.Get("status"); status != "" {
		if !h.isSuperAdminRequest(r) && !h.hasAnyPermission(r, permissions.AttendanceOperationsView, permissions.OvertimeApprove) {
			h.respondError(w, r, http.StatusForbidden, operation, nil, "permission required")
			return
		}
		items, err := h.svc.ListOvertimeRequestsByStatus(r.Context(), tenantID, status)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list overtime requests")
			return
		}
		respondJSON(w, http.StatusOK, items)
		return
	}
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
		[]string{permissions.AttendanceSelfView, permissions.AttendanceSelfOvertimeRequest},
		[]string{permissions.AttendanceOperationsView, permissions.OvertimeApprove},
	) {
		return
	}
	items, err := h.svc.ListOvertimeRequestsByUser(r.Context(), tenantID, userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list overtime requests")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) reviewOvertimeRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, requestID uuid.UUID, operation string, status string) {
	if !h.isSuperAdminRequest(r) && !h.hasAnyPermission(r, permissions.OvertimeApprove, permissions.AttendanceOperationsManage) {
		h.respondError(w, r, http.StatusForbidden, operation, nil, "permission required")
		return
	}
	var cmd ports.OvertimeReviewCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.RequestID = requestID
	cmd.Status = status
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ReviewOvertimeRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) overtimeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	requestID, err := uuid.Parse(chi.URLParam(r, "overtimeRequestID"))
	if err != nil || requestID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid overtime request id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, requestID, true
}

func (h *Handler) superAdminOvertimeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	requestID, err := uuid.Parse(chi.URLParam(r, "overtimeRequestID"))
	if err != nil || requestID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid overtime request id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, requestID, true
}
