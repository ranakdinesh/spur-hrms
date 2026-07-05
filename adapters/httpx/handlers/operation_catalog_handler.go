package handlers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListOperationCatalog(w http.ResponseWriter, r *http.Request) {
	tenantID := uuid.Nil
	if !h.isSuperAdminRequest(r) {
		var err error
		tenantID, err = h.tenantIDFromRequest(r)
		if err != nil {
			h.respondError(w, r, http.StatusUnauthorized, "list operation catalog", err, "tenant context is required")
			return
		}
	} else if parsed, err := h.tenantIDFromRequest(r); err == nil {
		tenantID = parsed
	}
	h.listOperationCatalogForTenant(w, r, tenantID, "list operation catalog")
}

func (h *Handler) ListTenantOperationCatalog(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant operation catalog"); ok {
		h.listOperationCatalogForTenant(w, r, tenantID, "list tenant operation catalog")
	}
}

func (h *Handler) listOperationCatalogForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	includeAll := false
	if flag := optionalBoolQuery(r, "include_all"); flag != nil {
		includeAll = *flag
	}
	result, err := h.svc.ListOperationCatalog(r.Context(), ports.OperationCatalogQuery{
		TenantID:    tenantID,
		IncludeAll:  includeAll,
		Permissions: splitOperationCatalogPermissions(r.URL.Query()["permission"]),
	})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func splitOperationCatalogPermissions(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			if clean := strings.TrimSpace(part); clean != "" {
				out = append(out, clean)
			}
		}
	}
	return out
}
