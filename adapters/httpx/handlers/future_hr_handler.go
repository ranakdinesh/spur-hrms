package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListBoundedAIAgents(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list bounded ai agents", err, "tenant context is required")
		return
	}
	h.listBoundedAIAgentsForTenant(w, r, tenantID, "list bounded ai agents")
}

func (h *Handler) RunBoundedAIAgents(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "run bounded ai agents", err, "tenant context is required")
		return
	}
	h.runBoundedAIAgentsForTenant(w, r, tenantID, "run bounded ai agents")
}

func (h *Handler) GetPeopleAnalyticsWorkspace(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get people analytics", err, "tenant context is required")
		return
	}
	h.getPeopleAnalyticsForTenant(w, r, tenantID, "get people analytics")
}

func (h *Handler) ListPrivacyEcosystemWorkspace(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list privacy ecosystem", err, "tenant context is required")
		return
	}
	h.listPrivacyEcosystemForTenant(w, r, tenantID, "list privacy ecosystem")
}

func (h *Handler) UpsertPrivacyConsent(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert privacy consent", err, "tenant context is required")
		return
	}
	h.upsertPrivacyConsentForTenant(w, r, tenantID, "upsert privacy consent")
}

func (h *Handler) CreateDataErasureRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create erasure request", err, "tenant context is required")
		return
	}
	h.createDataErasureRequestForTenant(w, r, tenantID, "create erasure request")
}

func (h *Handler) UpdateDataErasureRequestStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "requestID", "update erasure request status")
	if !ok {
		return
	}
	h.updateDataErasureRequestStatusForTenant(w, r, tenantID, id, "update erasure request status")
}

func (h *Handler) UpsertEcosystemIntegrationHook(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert integration hook", err, "tenant context is required")
		return
	}
	h.upsertEcosystemIntegrationHookForTenant(w, r, tenantID, "upsert integration hook")
}

func (h *Handler) UpsertMobileAPIConstraint(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert mobile constraint", err, "tenant context is required")
		return
	}
	h.upsertMobileAPIConstraintForTenant(w, r, tenantID, "upsert mobile constraint")
}

func (h *Handler) ListTenantBoundedAIAgents(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant bounded ai agents"); ok {
		h.listBoundedAIAgentsForTenant(w, r, tenantID, "list tenant bounded ai agents")
	}
}

func (h *Handler) RunTenantBoundedAIAgents(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "run tenant bounded ai agents"); ok {
		h.runBoundedAIAgentsForTenant(w, r, tenantID, "run tenant bounded ai agents")
	}
}

func (h *Handler) GetTenantPeopleAnalyticsWorkspace(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant people analytics"); ok {
		h.getPeopleAnalyticsForTenant(w, r, tenantID, "get tenant people analytics")
	}
}

func (h *Handler) ListTenantPrivacyEcosystemWorkspace(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant privacy ecosystem"); ok {
		h.listPrivacyEcosystemForTenant(w, r, tenantID, "list tenant privacy ecosystem")
	}
}

func (h *Handler) UpsertTenantPrivacyConsent(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant privacy consent"); ok {
		h.upsertPrivacyConsentForTenant(w, r, tenantID, "upsert tenant privacy consent")
	}
}

func (h *Handler) CreateTenantDataErasureRequest(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant erasure request"); ok {
		h.createDataErasureRequestForTenant(w, r, tenantID, "create tenant erasure request")
	}
}

func (h *Handler) UpdateTenantDataErasureRequestStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "update tenant erasure request status")
	if !ok {
		return
	}
	id, err := uuid.Parse(chi.URLParam(r, "requestID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant erasure request status", err, "invalid request id")
		return
	}
	h.updateDataErasureRequestStatusForTenant(w, r, tenantID, id, "update tenant erasure request status")
}

func (h *Handler) UpsertTenantEcosystemIntegrationHook(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant integration hook"); ok {
		h.upsertEcosystemIntegrationHookForTenant(w, r, tenantID, "upsert tenant integration hook")
	}
}

func (h *Handler) UpsertTenantMobileAPIConstraint(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant mobile constraint"); ok {
		h.upsertMobileAPIConstraintForTenant(w, r, tenantID, "upsert tenant mobile constraint")
	}
}

func (h *Handler) listBoundedAIAgentsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListBoundedAIAgents(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) runBoundedAIAgentsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.BoundedAIAgentRunCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	result, err := h.svc.RunBoundedAIAgents(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, result)
}

func (h *Handler) getPeopleAnalyticsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workspace, err := h.svc.GetPeopleAnalyticsWorkspace(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to load people analytics")
		return
	}
	respondJSON(w, http.StatusOK, workspace)
}

func (h *Handler) listPrivacyEcosystemForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workspace, err := h.svc.ListPrivacyEcosystemWorkspace(r.Context(), privacyFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to load privacy ecosystem")
		return
	}
	respondJSON(w, http.StatusOK, workspace)
}

func (h *Handler) upsertPrivacyConsentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PrivacyConsentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertPrivacyConsent(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) createDataErasureRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.DataErasureRequestCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateDataErasureRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateDataErasureRequestStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.DataErasureStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateDataErasureRequestStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) upsertEcosystemIntegrationHookForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EcosystemIntegrationHookCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertEcosystemIntegrationHook(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) upsertMobileAPIConstraintForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.MobileAPIConstraintCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertMobileAPIConstraint(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func privacyFilterFromRequest(r *http.Request, tenantID uuid.UUID) domain.PrivacyEcosystemFilter {
	limit := int32(100)
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}
	offset := int32(0)
	if raw := r.URL.Query().Get("offset"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			offset = int32(parsed)
		}
	}
	return domain.PrivacyEcosystemFilter{TenantID: tenantID, Status: reportOptionalQuery(r, "status"), ConsentArea: reportOptionalQuery(r, "consent_area"), Priority: reportOptionalQuery(r, "priority"), Channel: reportOptionalQuery(r, "channel"), Workflow: reportOptionalQuery(r, "workflow"), Limit: limit, Offset: offset}
}
