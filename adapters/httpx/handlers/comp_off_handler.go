package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
)

func (h *Handler) CreateCompOffRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create comp-off request", err, "tenant context is required")
		return
	}
	h.createCompOffRequestForTenant(w, r, tenantID, "create comp-off request")
}

func (h *Handler) ListCompOffRequests(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list comp-off requests", err, "tenant context is required")
		return
	}
	h.listCompOffRequestsForTenant(w, r, tenantID, "list comp-off requests")
}

func (h *Handler) ApproveCompOffRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, requestID, ok := h.compOffRequestIDs(w, r, "approve comp-off request")
	if !ok {
		return
	}
	h.reviewCompOffRequestForTenant(w, r, tenantID, requestID, "approve comp-off request", "approved")
}

func (h *Handler) RejectCompOffRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, requestID, ok := h.compOffRequestIDs(w, r, "reject comp-off request")
	if !ok {
		return
	}
	h.reviewCompOffRequestForTenant(w, r, tenantID, requestID, "reject comp-off request", "rejected")
}

func (h *Handler) CreateTenantCompOffRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant comp-off request")
	if ok {
		h.createCompOffRequestForTenant(w, r, tenantID, "create tenant comp-off request")
	}
}

func (h *Handler) ListTenantCompOffRequests(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant comp-off requests")
	if ok {
		h.listCompOffRequestsForTenant(w, r, tenantID, "list tenant comp-off requests")
	}
}

func (h *Handler) ApproveTenantCompOffRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, requestID, ok := h.superAdminCompOffRequestIDs(w, r, "approve tenant comp-off request")
	if !ok {
		return
	}
	h.reviewCompOffRequestForTenant(w, r, tenantID, requestID, "approve tenant comp-off request", "approved")
}

func (h *Handler) RejectTenantCompOffRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, requestID, ok := h.superAdminCompOffRequestIDs(w, r, "reject tenant comp-off request")
	if !ok {
		return
	}
	h.reviewCompOffRequestForTenant(w, r, tenantID, requestID, "reject tenant comp-off request", "rejected")
}

func (h *Handler) createCompOffRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CompOffRequestCommand
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
		[]string{permissions.LeaveSelfCompOffRequest},
		[]string{permissions.LeaveOperationsManage, permissions.CompOffApprove},
	) {
		return
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateCompOffRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listCompOffRequestsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	query := r.URL.Query()
	if status := query.Get("status"); status != "" {
		if !h.isSuperAdminRequest(r) && !h.hasAnyPermission(r, permissions.LeaveOperationsView, permissions.CompOffApprove) {
			h.respondError(w, r, http.StatusForbidden, operation, nil, "permission required")
			return
		}
		items, err := h.svc.ListCompOffRequestsByStatus(r.Context(), tenantID, status)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list comp-off requests")
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
		[]string{permissions.LeaveSelfView, permissions.LeaveSelfCompOffRequest},
		[]string{permissions.LeaveOperationsView, permissions.CompOffApprove},
	) {
		return
	}
	items, err := h.svc.ListCompOffRequestsByUser(r.Context(), tenantID, userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list comp-off requests")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) reviewCompOffRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, requestID uuid.UUID, operation string, status string) {
	if !h.hasAnyPermission(r, permissions.CompOffApprove, permissions.LeaveOperationsManage) && !h.isSuperAdminRequest(r) {
		h.respondError(w, r, http.StatusForbidden, operation, nil, "permission required")
		return
	}
	var cmd ports.CompOffReviewCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.RequestID = requestID
	cmd.Status = status
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ReviewCompOffRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) compOffRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	requestID, err := uuid.Parse(chi.URLParam(r, "compOffRequestID"))
	if err != nil || requestID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid comp-off request id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, requestID, true
}

func (h *Handler) superAdminCompOffRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	requestID, err := uuid.Parse(chi.URLParam(r, "compOffRequestID"))
	if err != nil || requestID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid comp-off request id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, requestID, true
}
