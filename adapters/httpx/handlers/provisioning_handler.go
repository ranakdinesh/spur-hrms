package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ProvisionTenant(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "provision tenant")
	if !ok {
		return
	}
	var cmd ports.ProvisionTenantCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode provision tenant request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	result, err := h.svc.ProvisionTenant(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "provision tenant", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

func (h *Handler) GetTenantProvisioningStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant provisioning status")
	if !ok {
		return
	}
	h.getProvisioningStatusForTenant(w, r, tenantID, "get tenant provisioning status")
}

func (h *Handler) GetProvisioningStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get provisioning status", err, "tenant context is required")
		return
	}
	h.getProvisioningStatusForTenant(w, r, tenantID, "get provisioning status")
}

func (h *Handler) getProvisioningStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	result, err := h.svc.GetTenantProvisioningStatus(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to load provisioning status")
		return
	}
	respondJSON(w, http.StatusOK, result)
}
