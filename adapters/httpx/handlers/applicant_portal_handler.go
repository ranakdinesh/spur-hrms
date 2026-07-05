package handlers

import (
	"net/http"

	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (h *Handler) GetApplicantPortal(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get applicant portal", err, "tenant context is required")
		return
	}
	actorID := h.actorIDFromRequest(r)
	if actorID == nil {
		h.respondError(w, r, http.StatusUnauthorized, "get applicant portal", domain.ErrInvalidApplicantUserID, "user context is required")
		return
	}
	portal, err := h.svc.GetApplicantPortal(r.Context(), tenantID, *actorID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get applicant portal", err, "applicant portal not found")
		return
	}
	respondJSON(w, http.StatusOK, portal)
}
