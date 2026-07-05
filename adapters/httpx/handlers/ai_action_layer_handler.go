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

func (h *Handler) ListAIActionWorkspace(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list ai action workspace", err, "tenant context is required")
		return
	}
	h.listAIActionWorkspaceForTenant(w, r, tenantID, "list ai action workspace")
}

func (h *Handler) CreateAISignal(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create ai signal", err, "tenant context is required")
		return
	}
	h.createAISignalForTenant(w, r, tenantID, "create ai signal")
}

func (h *Handler) CreateAIAgentAction(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create ai action", err, "tenant context is required")
		return
	}
	h.createAIAgentActionForTenant(w, r, tenantID, "create ai action")
}

func (h *Handler) UpdateAIAgentActionStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, actionID, ok := h.tenantAndURLUUID(w, r, "actionID", "update ai action status")
	if !ok {
		return
	}
	h.updateAIAgentActionStatusForTenant(w, r, tenantID, actionID, "update ai action status")
}

func (h *Handler) OverrideAIAction(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "override ai action", err, "tenant context is required")
		return
	}
	h.overrideAIActionForTenant(w, r, tenantID, "override ai action")
}

func (h *Handler) EmitAIWorkflowEvent(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "emit ai workflow event", err, "tenant context is required")
		return
	}
	h.emitAIWorkflowEventForTenant(w, r, tenantID, "emit ai workflow event")
}

func (h *Handler) ListTenantAIActionWorkspace(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant ai action workspace"); ok {
		h.listAIActionWorkspaceForTenant(w, r, tenantID, "list tenant ai action workspace")
	}
}

func (h *Handler) CreateTenantAISignal(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant ai signal"); ok {
		h.createAISignalForTenant(w, r, tenantID, "create tenant ai signal")
	}
}

func (h *Handler) CreateTenantAIAgentAction(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant ai action"); ok {
		h.createAIAgentActionForTenant(w, r, tenantID, "create tenant ai action")
	}
}

func (h *Handler) UpdateTenantAIAgentActionStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "update tenant ai action status")
	if !ok {
		return
	}
	actionID, err := uuid.Parse(chi.URLParam(r, "actionID"))
	if err != nil || actionID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant ai action status", err, "invalid action id")
		return
	}
	h.updateAIAgentActionStatusForTenant(w, r, tenantID, actionID, "update tenant ai action status")
}

func (h *Handler) OverrideTenantAIAction(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "override tenant ai action"); ok {
		h.overrideAIActionForTenant(w, r, tenantID, "override tenant ai action")
	}
}

func (h *Handler) EmitTenantAIWorkflowEvent(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "emit tenant ai workflow event"); ok {
		h.emitAIWorkflowEventForTenant(w, r, tenantID, "emit tenant ai workflow event")
	}
}

func (h *Handler) listAIActionWorkspaceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workspace, err := h.svc.ListAIActionWorkspace(r.Context(), aiActionFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list ai action workspace")
		return
	}
	respondJSON(w, http.StatusOK, workspace)
}

func (h *Handler) createAISignalForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AISignalCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAISignal(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) createAIAgentActionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AIActionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAIAgentAction(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateAIAgentActionStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, actionID uuid.UUID, operation string) {
	var cmd ports.AIStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = actionID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateAIAgentActionStatus(r.Context(), cmd)
	if err != nil {
		status := http.StatusBadRequest
		if err == domain.ErrAIActionNotFound {
			status = http.StatusNotFound
		}
		h.respondError(w, r, status, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) overrideAIActionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AIOverrideCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.OverrideAIAction(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) emitAIWorkflowEventForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AIWorkflowEventCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.EmitAIWorkflowEvent(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func aiActionFilterFromRequest(r *http.Request, tenantID uuid.UUID) domain.AIActionFilter {
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
	var insightID *uuid.UUID
	if raw := r.URL.Query().Get("insight_id"); raw != "" {
		if parsed, err := uuid.Parse(raw); err == nil && parsed != uuid.Nil {
			insightID = &parsed
		}
	}
	return domain.AIActionFilter{TenantID: tenantID, Status: reportOptionalQuery(r, "status"), Severity: reportOptionalQuery(r, "severity"), AgentKey: reportOptionalQuery(r, "agent_key"), SourceModule: reportOptionalQuery(r, "source_module"), ProcessingStatus: reportOptionalQuery(r, "processing_status"), VisibilityScope: reportOptionalQuery(r, "visibility_scope"), EventType: reportOptionalQuery(r, "event_type"), Decision: reportOptionalQuery(r, "decision"), InsightID: insightID, Limit: limit, Offset: offset}
}
