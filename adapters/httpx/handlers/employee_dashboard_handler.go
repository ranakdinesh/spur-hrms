package handlers

import (
	"net/http"

	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (h *Handler) GetEmployeeDashboard(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get employee dashboard", err, "tenant context is required")
		return
	}
	actorID := h.actorIDFromRequest(r)
	if actorID == nil {
		h.respondError(w, r, http.StatusUnauthorized, "get employee dashboard", domain.ErrInvalidEmployeeUserID, "user context is required")
		return
	}
	dashboard, err := h.svc.GetEmployeeDashboard(r.Context(), tenantID, *actorID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "get employee dashboard", err, "failed to load employee dashboard")
		return
	}
	respondJSON(w, http.StatusOK, dashboard)
}
