package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListSuccessionReviewCycles(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listSuccessionReviewCyclesForTenant(w, r, tenantID, "list succession review cycles")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list succession review cycles", err, "tenant context is required")
	}
}

func (h *Handler) CreateSuccessionReviewCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createSuccessionReviewCycleForTenant(w, r, tenantID, "create succession review cycle")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create succession review cycle", err, "tenant context is required")
	}
}

func (h *Handler) UpdateSuccessionReviewCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.successionTenantAndID(w, r, "cycleID", "update succession review cycle"); ok {
		h.updateSuccessionReviewCycleForTenant(w, r, tenantID, id, "update succession review cycle")
	}
}

func (h *Handler) UpdateSuccessionReviewCycleStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.successionTenantAndID(w, r, "cycleID", "update succession review cycle status"); ok {
		h.updateSuccessionStatus(w, r, tenantID, id, "cycle", "update succession review cycle status")
	}
}

func (h *Handler) ListSuccessionCriticalRoles(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listSuccessionCriticalRolesForTenant(w, r, tenantID, "list succession critical roles")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list succession critical roles", err, "tenant context is required")
	}
}

func (h *Handler) CreateSuccessionCriticalRole(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createSuccessionCriticalRoleForTenant(w, r, tenantID, "create succession critical role")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create succession critical role", err, "tenant context is required")
	}
}

func (h *Handler) UpdateSuccessionCriticalRole(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.successionTenantAndID(w, r, "roleID", "update succession critical role"); ok {
		h.updateSuccessionCriticalRoleForTenant(w, r, tenantID, id, "update succession critical role")
	}
}

func (h *Handler) UpdateSuccessionCriticalRoleStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.successionTenantAndID(w, r, "roleID", "update succession critical role status"); ok {
		h.updateSuccessionStatus(w, r, tenantID, id, "critical_role", "update succession critical role status")
	}
}

func (h *Handler) DeleteSuccessionCriticalRole(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.successionTenantAndID(w, r, "roleID", "delete succession critical role"); ok {
		if err := h.svc.DeleteSuccessionCriticalRole(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete succession critical role", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListSuccessionSuccessorNominations(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listSuccessionSuccessorNominationsForTenant(w, r, tenantID, "list succession successors")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list succession successors", err, "tenant context is required")
	}
}

func (h *Handler) CreateSuccessionSuccessorNomination(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createSuccessionSuccessorNominationForTenant(w, r, tenantID, "create succession successor")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create succession successor", err, "tenant context is required")
	}
}

func (h *Handler) UpdateSuccessionSuccessorNomination(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.successionTenantAndID(w, r, "nominationID", "update succession successor"); ok {
		h.updateSuccessionSuccessorNominationForTenant(w, r, tenantID, id, "update succession successor")
	}
}

func (h *Handler) UpdateSuccessionSuccessorNominationStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.successionTenantAndID(w, r, "nominationID", "update succession successor status"); ok {
		h.updateSuccessionStatus(w, r, tenantID, id, "nomination", "update succession successor status")
	}
}

func (h *Handler) ListSuccessionDevelopmentActions(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listSuccessionDevelopmentActionsForTenant(w, r, tenantID, "list succession development actions")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list succession development actions", err, "tenant context is required")
	}
}

func (h *Handler) CreateSuccessionDevelopmentAction(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createSuccessionDevelopmentActionForTenant(w, r, tenantID, "create succession development action")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create succession development action", err, "tenant context is required")
	}
}

func (h *Handler) UpdateSuccessionDevelopmentAction(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.successionTenantAndID(w, r, "actionID", "update succession development action"); ok {
		h.updateSuccessionDevelopmentActionForTenant(w, r, tenantID, id, "update succession development action")
	}
}

func (h *Handler) UpdateSuccessionDevelopmentActionStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.successionTenantAndID(w, r, "actionID", "update succession development action status"); ok {
		h.updateSuccessionStatus(w, r, tenantID, id, "development_action", "update succession development action status")
	}
}

func (h *Handler) ListSuccessionEvents(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listSuccessionEventsForTenant(w, r, tenantID, "list succession events")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list succession events", err, "tenant context is required")
	}
}

func (h *Handler) GetSuccessionSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.getSuccessionSummaryForTenant(w, r, tenantID, "get succession summary")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "get succession summary", err, "tenant context is required")
	}
}

func (h *Handler) ListTenantSuccessionReviewCycles(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant succession review cycles"); ok {
		h.listSuccessionReviewCyclesForTenant(w, r, tenantID, "list tenant succession review cycles")
	}
}
func (h *Handler) CreateTenantSuccessionReviewCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant succession review cycle"); ok {
		h.createSuccessionReviewCycleForTenant(w, r, tenantID, "create tenant succession review cycle")
	}
}
func (h *Handler) UpdateTenantSuccessionReviewCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.tenantSuccessionID(w, r, "cycleID", "update tenant succession review cycle"); ok {
		h.updateSuccessionReviewCycleForTenant(w, r, tenantID, id, "update tenant succession review cycle")
	}
}
func (h *Handler) UpdateTenantSuccessionReviewCycleStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateSuccessionReviewCycleStatus(w, r)
}
func (h *Handler) ListTenantSuccessionCriticalRoles(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant succession critical roles"); ok {
		h.listSuccessionCriticalRolesForTenant(w, r, tenantID, "list tenant succession critical roles")
	}
}
func (h *Handler) CreateTenantSuccessionCriticalRole(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant succession critical role"); ok {
		h.createSuccessionCriticalRoleForTenant(w, r, tenantID, "create tenant succession critical role")
	}
}
func (h *Handler) UpdateTenantSuccessionCriticalRole(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.tenantSuccessionID(w, r, "roleID", "update tenant succession critical role"); ok {
		h.updateSuccessionCriticalRoleForTenant(w, r, tenantID, id, "update tenant succession critical role")
	}
}
func (h *Handler) UpdateTenantSuccessionCriticalRoleStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateSuccessionCriticalRoleStatus(w, r)
}
func (h *Handler) DeleteTenantSuccessionCriticalRole(w http.ResponseWriter, r *http.Request) {
	h.DeleteSuccessionCriticalRole(w, r)
}
func (h *Handler) ListTenantSuccessionSuccessorNominations(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant succession successors"); ok {
		h.listSuccessionSuccessorNominationsForTenant(w, r, tenantID, "list tenant succession successors")
	}
}
func (h *Handler) CreateTenantSuccessionSuccessorNomination(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant succession successor"); ok {
		h.createSuccessionSuccessorNominationForTenant(w, r, tenantID, "create tenant succession successor")
	}
}
func (h *Handler) UpdateTenantSuccessionSuccessorNomination(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.tenantSuccessionID(w, r, "nominationID", "update tenant succession successor"); ok {
		h.updateSuccessionSuccessorNominationForTenant(w, r, tenantID, id, "update tenant succession successor")
	}
}
func (h *Handler) UpdateTenantSuccessionSuccessorNominationStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateSuccessionSuccessorNominationStatus(w, r)
}
func (h *Handler) ListTenantSuccessionDevelopmentActions(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant succession development actions"); ok {
		h.listSuccessionDevelopmentActionsForTenant(w, r, tenantID, "list tenant succession development actions")
	}
}
func (h *Handler) CreateTenantSuccessionDevelopmentAction(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant succession development action"); ok {
		h.createSuccessionDevelopmentActionForTenant(w, r, tenantID, "create tenant succession development action")
	}
}
func (h *Handler) UpdateTenantSuccessionDevelopmentAction(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.tenantSuccessionID(w, r, "actionID", "update tenant succession development action"); ok {
		h.updateSuccessionDevelopmentActionForTenant(w, r, tenantID, id, "update tenant succession development action")
	}
}
func (h *Handler) UpdateTenantSuccessionDevelopmentActionStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateSuccessionDevelopmentActionStatus(w, r)
}
func (h *Handler) ListTenantSuccessionEvents(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant succession events"); ok {
		h.listSuccessionEventsForTenant(w, r, tenantID, "list tenant succession events")
	}
}
func (h *Handler) GetTenantSuccessionSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant succession summary"); ok {
		h.getSuccessionSummaryForTenant(w, r, tenantID, "get tenant succession summary")
	}
}

func (h *Handler) listSuccessionReviewCyclesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListSuccessionReviewCycles(r.Context(), successionFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list review cycles")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createSuccessionReviewCycleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.SuccessionReviewCycleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateSuccessionReviewCycle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateSuccessionReviewCycleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.SuccessionReviewCycleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateSuccessionReviewCycle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listSuccessionCriticalRolesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListSuccessionCriticalRoles(r.Context(), successionFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list critical roles")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createSuccessionCriticalRoleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.SuccessionCriticalRoleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateSuccessionCriticalRole(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateSuccessionCriticalRoleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.SuccessionCriticalRoleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateSuccessionCriticalRole(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listSuccessionSuccessorNominationsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListSuccessionSuccessorNominations(r.Context(), successionFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list successors")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createSuccessionSuccessorNominationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.SuccessionSuccessorNominationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateSuccessionSuccessorNomination(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateSuccessionSuccessorNominationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.SuccessionSuccessorNominationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateSuccessionSuccessorNomination(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listSuccessionDevelopmentActionsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListSuccessionDevelopmentActions(r.Context(), successionFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list development actions")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createSuccessionDevelopmentActionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.SuccessionDevelopmentActionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateSuccessionDevelopmentAction(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateSuccessionDevelopmentActionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.SuccessionDevelopmentActionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateSuccessionDevelopmentAction(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateSuccessionStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, kind string, operation string) {
	var cmd ports.SuccessionStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	var item any
	var err error
	switch kind {
	case "cycle":
		item, err = h.svc.UpdateSuccessionReviewCycleStatus(r.Context(), cmd)
	case "critical_role":
		item, err = h.svc.UpdateSuccessionCriticalRoleStatus(r.Context(), cmd)
	case "nomination":
		item, err = h.svc.UpdateSuccessionSuccessorNominationStatus(r.Context(), cmd)
	default:
		item, err = h.svc.UpdateSuccessionDevelopmentActionStatus(r.Context(), cmd)
	}
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listSuccessionEventsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListSuccessionEvents(r.Context(), successionFilterFromRequest(r, tenantID), optionalStringQuery(r, "source_type"), optionalUUIDQuery(r, "source_id"))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list succession events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getSuccessionSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.GetSuccessionSummary(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to get succession summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func successionFilterFromRequest(r *http.Request, tenantID uuid.UUID) domain.SuccessionFilter {
	return domain.SuccessionFilter{TenantID: tenantID, CycleID: optionalUUIDQuery(r, "cycle_id"), CriticalRoleID: optionalUUIDQuery(r, "critical_role_id"), Status: optionalStringQuery(r, "status"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 50), Offset: queryInt32(r, "offset", 0)}
}

func (h *Handler) successionTenantAndID(w http.ResponseWriter, r *http.Request, idParam string, operation string) (uuid.UUID, uuid.UUID, bool) {
	if chi.URLParam(r, "tenantID") != "" {
		return h.tenantSuccessionID(w, r, idParam, operation)
	}
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, idParam))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) tenantSuccessionID(w http.ResponseWriter, r *http.Request, idParam string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, idParam))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
