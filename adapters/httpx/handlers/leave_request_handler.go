package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
)

func (h *Handler) ApplyLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "apply leave", err, "tenant context is required")
		return
	}
	h.applyLeaveForTenant(w, r, tenantID, "apply leave")
}

func (h *Handler) ListLeaves(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leaves", err, "tenant context is required")
		return
	}
	h.listLeavesForTenant(w, r, tenantID, "list leaves")
}

func (h *Handler) PreviewLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "preview leave", err, "tenant context is required")
		return
	}
	h.previewLeaveForTenant(w, r, tenantID, "preview leave")
}

func (h *Handler) ListLeaveReport(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leave report", err, "tenant context is required")
		return
	}
	h.listLeaveReportForTenant(w, r, tenantID, "list leave report")
}

func (h *Handler) GetLeaveReportSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get leave report summary", err, "tenant context is required")
		return
	}
	h.getLeaveReportSummaryForTenant(w, r, tenantID, "get leave report summary")
}

func (h *Handler) ApplyTenantLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "apply tenant leave")
	if ok {
		h.applyLeaveForTenant(w, r, tenantID, "apply tenant leave")
	}
}

func (h *Handler) ListTenantLeaves(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leaves")
	if ok {
		h.listLeavesForTenant(w, r, tenantID, "list tenant leaves")
	}
}

func (h *Handler) ListTenantLeaveReport(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leave report")
	if ok {
		h.listLeaveReportForTenant(w, r, tenantID, "list tenant leave report")
	}
}

func (h *Handler) GetTenantLeaveReportSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant leave report summary")
	if ok {
		h.getLeaveReportSummaryForTenant(w, r, tenantID, "get tenant leave report summary")
	}
}

func (h *Handler) PreviewTenantLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "preview tenant leave")
	if ok {
		h.previewLeaveForTenant(w, r, tenantID, "preview tenant leave")
	}
}

func (h *Handler) applyLeaveForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ApplyLeaveCommand
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
		[]string{permissions.LeaveSelfApply, permissions.LeavesApply},
		[]string{permissions.LeaveOperationsManage},
	) {
		return
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ApplyLeave(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) previewLeaveForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ApplyLeaveCommand
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
		[]string{permissions.LeaveSelfView, permissions.LeaveSelfApply, permissions.LeavesView, permissions.LeavesApply},
		[]string{permissions.LeaveOperationsView},
	) {
		return
	}
	item, err := h.svc.PreviewLeave(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listLeavesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if fyRaw := r.URL.Query().Get("fy_id"); fyRaw != "" {
		if !h.requirePermission(w, r, operation, permissions.LeaveOperationsView) {
			return
		}
		fyID, err := uuid.Parse(fyRaw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "parse fy id", err, "invalid fy_id")
			return
		}
		items, err := h.svc.ListLeavesByFY(r.Context(), tenantID, fyID)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leaves")
			return
		}
		respondJSON(w, http.StatusOK, items)
		return
	}
	userRaw := r.URL.Query().Get("user_id")
	var userID uuid.UUID
	var err error
	if userRaw != "" {
		userID, err = uuid.Parse(userRaw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "parse user id", err, "invalid user_id")
			return
		}
	} else if actorID := h.actorIDFromRequest(r); actorID != nil {
		userID = *actorID
	}
	if !h.requireOwnUserOrPermission(w, r, operation, userID,
		[]string{permissions.LeaveSelfView, permissions.LeavesList, permissions.LeavesView},
		[]string{permissions.LeaveOperationsView},
	) {
		return
	}
	items, err := h.svc.ListLeavesByUser(r.Context(), tenantID, userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leaves")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listLeaveReportForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter, ok := h.leaveReportFilterFromRequest(w, r, tenantID, operation)
	if !ok {
		return
	}
	items, err := h.svc.ListLeaveReportRows(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave report")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getLeaveReportSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter, ok := h.leaveReportFilterFromRequest(w, r, tenantID, operation)
	if !ok {
		return
	}
	item, err := h.svc.GetLeaveReportSummary(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to get leave report summary")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) leaveReportFilterFromRequest(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (domain.LeaveReportFilter, bool) {
	query := r.URL.Query()
	filter := domain.LeaveReportFilter{TenantID: tenantID}
	if query.Get("scope") == "manager" {
		filter.ManagerID = h.actorIDFromRequest(r)
	}
	if !h.parseOptionalUUID(w, r, operation, query.Get("manager_id"), "manager_id", &filter.ManagerID) {
		return filter, false
	}
	if !h.parseOptionalUUID(w, r, operation, query.Get("fy_id"), "fy_id", &filter.FYID) {
		return filter, false
	}
	if !h.parseOptionalUUID(w, r, operation, query.Get("user_id"), "user_id", &filter.UserID) {
		return filter, false
	}
	if !h.parseOptionalUUID(w, r, operation, query.Get("department_id"), "department_id", &filter.DepartmentID) {
		return filter, false
	}
	if !h.parseOptionalUUID(w, r, operation, query.Get("leave_type_id"), "leave_type_id", &filter.LeaveTypeID) {
		return filter, false
	}
	if status := query.Get("status"); status != "" {
		filter.Status = &status
	}
	if !h.parseOptionalDate(w, r, operation, query.Get("start_date"), "start_date", &filter.StartDate) {
		return filter, false
	}
	if !h.parseOptionalDate(w, r, operation, query.Get("end_date"), "end_date", &filter.EndDate) {
		return filter, false
	}
	return filter, true
}

func (h *Handler) parseOptionalUUID(w http.ResponseWriter, r *http.Request, operation string, raw string, field string, target **uuid.UUID) bool {
	if raw == "" {
		return true
	}
	id, err := uuid.Parse(raw)
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid "+field)
		return false
	}
	*target = &id
	return true
}

func (h *Handler) parseOptionalDate(w http.ResponseWriter, r *http.Request, operation string, raw string, field string, target **time.Time) bool {
	if raw == "" {
		return true
	}
	date, err := time.Parse("2006-01-02", raw)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid "+field)
		return false
	}
	*target = &date
	return true
}
