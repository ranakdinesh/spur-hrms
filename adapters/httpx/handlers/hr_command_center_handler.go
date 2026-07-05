package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetHRCommandCenter(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get hr command center", err, "tenant context is required")
		return
	}
	h.getHRCommandCenterForTenant(w, r, tenantID, "get hr command center")
}

func (h *Handler) GetTenantHRCommandCenter(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant hr command center"); ok {
		h.getHRCommandCenterForTenant(w, r, tenantID, "get tenant hr command center")
	}
}

func (h *Handler) getHRCommandCenterForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	result, err := h.svc.GetHRCommandCenter(r.Context(), ports.HRCommandCenterQuery{
		TenantID: tenantID,
		Limit:    queryInt32(r, "limit", 350),
	})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}
