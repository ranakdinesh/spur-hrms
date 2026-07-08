package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
)

func (h *Handler) GetHRDashboard(w http.ResponseWriter, r *http.Request) {
	if !h.requirePermission(w, r, "get hr dashboard", permissions.DashboardHRView) {
		return
	}
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get hr dashboard", err, "tenant context is required")
		return
	}
	h.getHRDashboardForTenant(w, r, tenantID, "get hr dashboard")
}

func (h *Handler) GetTenantHRDashboard(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant hr dashboard")
	if !ok {
		return
	}
	h.getHRDashboardForTenant(w, r, tenantID, "get tenant hr dashboard")
}

func (h *Handler) getHRDashboardForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	now := time.Now().UTC()
	month := int32(now.Month())
	year := int32(now.Year())
	if value := r.URL.Query().Get("month"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation+" parse month", err, "invalid month")
			return
		}
		month = int32(parsed)
	}
	if value := r.URL.Query().Get("year"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation+" parse year", err, "invalid year")
			return
		}
		year = int32(parsed)
	}
	dashboard, err := h.svc.GetHRDashboard(r.Context(), domain.HRDashboardQuery{TenantID: tenantID, Month: month, Year: year})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to load HR dashboard")
		return
	}
	respondJSON(w, http.StatusOK, dashboard)
}
