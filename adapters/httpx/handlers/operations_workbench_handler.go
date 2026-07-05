package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetOperationsWorkbench(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get operations workbench", err, "tenant context is required")
		return
	}
	h.getOperationsWorkbenchForTenant(w, r, tenantID, "get operations workbench")
}

func (h *Handler) GetTenantOperationsWorkbench(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant operations workbench"); ok {
		h.getOperationsWorkbenchForTenant(w, r, tenantID, "get tenant operations workbench")
	}
}

func (h *Handler) ActOperationsWorkbenchCard(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "act operations workbench card", err, "tenant context is required")
		return
	}
	h.actOperationsWorkbenchCardForTenant(w, r, tenantID, "act operations workbench card")
}

func (h *Handler) ActTenantOperationsWorkbenchCard(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "act tenant operations workbench card"); ok {
		h.actOperationsWorkbenchCardForTenant(w, r, tenantID, "act tenant operations workbench card")
	}
}

func (h *Handler) getOperationsWorkbenchForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	result, err := h.svc.GetOperationsWorkbench(r.Context(), ports.OperationsWorkbenchQuery{
		TenantID: tenantID,
		Lane:     optionalStringQuery(r, "lane"),
		Category: optionalStringQuery(r, "category"),
		Severity: optionalStringQuery(r, "severity"),
		Search:   optionalStringQuery(r, "search"),
		Limit:    queryInt32(r, "limit", 200),
		Offset:   queryInt32(r, "offset", 0),
	})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) actOperationsWorkbenchCardForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.OperationsWorkbenchActionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	result, err := h.svc.ActOperationsWorkbenchCard(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}
