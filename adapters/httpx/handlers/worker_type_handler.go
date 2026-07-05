package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateWorkerType(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create worker type", err, "tenant context is required")
		return
	}
	h.createWorkerTypeForTenant(w, r, tenantID, "create worker type")
}

func (h *Handler) ListWorkerTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list worker types", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListWorkerTypes(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list worker types", err, "failed to list worker types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetWorkerType(w http.ResponseWriter, r *http.Request) {
	tenantID, workerTypeID, ok := h.workerTypeRequestIDs(w, r, "get worker type")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkerType(r.Context(), tenantID, workerTypeID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get worker type", err, "worker type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateWorkerType(w http.ResponseWriter, r *http.Request) {
	tenantID, workerTypeID, ok := h.workerTypeRequestIDs(w, r, "update worker type")
	if !ok {
		return
	}
	h.updateWorkerTypeForTenant(w, r, tenantID, workerTypeID, "update worker type")
}

func (h *Handler) DeleteWorkerType(w http.ResponseWriter, r *http.Request) {
	tenantID, workerTypeID, ok := h.workerTypeRequestIDs(w, r, "delete worker type")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkerType(r.Context(), tenantID, workerTypeID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete worker type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantWorkerType(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant worker type")
	if !ok {
		return
	}
	h.createWorkerTypeForTenant(w, r, tenantID, "create tenant worker type")
}

func (h *Handler) ListTenantWorkerTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant worker types")
	if !ok {
		return
	}
	items, err := h.svc.ListWorkerTypes(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant worker types", err, "failed to list worker types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetTenantWorkerType(w http.ResponseWriter, r *http.Request) {
	tenantID, workerTypeID, ok := h.superAdminWorkerTypeRequestIDs(w, r, "get tenant worker type")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkerType(r.Context(), tenantID, workerTypeID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant worker type", err, "worker type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantWorkerType(w http.ResponseWriter, r *http.Request) {
	tenantID, workerTypeID, ok := h.superAdminWorkerTypeRequestIDs(w, r, "update tenant worker type")
	if !ok {
		return
	}
	h.updateWorkerTypeForTenant(w, r, tenantID, workerTypeID, "update tenant worker type")
}

func (h *Handler) DeleteTenantWorkerType(w http.ResponseWriter, r *http.Request) {
	tenantID, workerTypeID, ok := h.superAdminWorkerTypeRequestIDs(w, r, "delete tenant worker type")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkerType(r.Context(), tenantID, workerTypeID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant worker type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateWorkerClassificationRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create worker classification rule", err, "tenant context is required")
		return
	}
	h.createWorkerClassificationRuleForTenant(w, r, tenantID, "create worker classification rule")
}

func (h *Handler) ListWorkerClassificationRules(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list worker classification rules", err, "tenant context is required")
		return
	}
	workerTypeID, ok := h.optionalUUIDQuery(w, r, "worker_type_id", "list worker classification rules")
	if !ok {
		return
	}
	items, err := h.svc.ListWorkerClassificationRules(r.Context(), tenantID, workerTypeID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list worker classification rules", err, "failed to list worker classification rules")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetWorkerClassificationRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.workerClassificationRuleRequestIDs(w, r, "get worker classification rule")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkerClassificationRule(r.Context(), tenantID, ruleID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get worker classification rule", err, "worker classification rule not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateWorkerClassificationRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.workerClassificationRuleRequestIDs(w, r, "update worker classification rule")
	if !ok {
		return
	}
	h.updateWorkerClassificationRuleForTenant(w, r, tenantID, ruleID, "update worker classification rule")
}

func (h *Handler) DeleteWorkerClassificationRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.workerClassificationRuleRequestIDs(w, r, "delete worker classification rule")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkerClassificationRule(r.Context(), tenantID, ruleID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete worker classification rule", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantWorkerClassificationRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant worker classification rule")
	if !ok {
		return
	}
	h.createWorkerClassificationRuleForTenant(w, r, tenantID, "create tenant worker classification rule")
}

func (h *Handler) ListTenantWorkerClassificationRules(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant worker classification rules")
	if !ok {
		return
	}
	workerTypeID, ok := h.optionalUUIDQuery(w, r, "worker_type_id", "list tenant worker classification rules")
	if !ok {
		return
	}
	items, err := h.svc.ListWorkerClassificationRules(r.Context(), tenantID, workerTypeID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant worker classification rules", err, "failed to list worker classification rules")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetTenantWorkerClassificationRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.superAdminWorkerClassificationRuleRequestIDs(w, r, "get tenant worker classification rule")
	if !ok {
		return
	}
	item, err := h.svc.GetWorkerClassificationRule(r.Context(), tenantID, ruleID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant worker classification rule", err, "worker classification rule not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantWorkerClassificationRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.superAdminWorkerClassificationRuleRequestIDs(w, r, "update tenant worker classification rule")
	if !ok {
		return
	}
	h.updateWorkerClassificationRuleForTenant(w, r, tenantID, ruleID, "update tenant worker classification rule")
}

func (h *Handler) DeleteTenantWorkerClassificationRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.superAdminWorkerClassificationRuleRequestIDs(w, r, "delete tenant worker classification rule")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkerClassificationRule(r.Context(), tenantID, ruleID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant worker classification rule", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createWorkerTypeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.WorkerTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkerType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateWorkerTypeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, workerTypeID uuid.UUID, operation string) {
	var cmd ports.WorkerTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = workerTypeID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkerType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createWorkerClassificationRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.WorkerClassificationRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkerClassificationRule(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateWorkerClassificationRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, ruleID uuid.UUID, operation string) {
	var cmd ports.WorkerClassificationRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = ruleID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkerClassificationRule(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) workerTypeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	workerTypeID, ok := h.workerTypeIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, workerTypeID, true
}

func (h *Handler) superAdminWorkerTypeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	workerTypeID, ok := h.workerTypeIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, workerTypeID, true
}

func (h *Handler) workerTypeIDFromURL(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	workerTypeID, err := uuid.Parse(chi.URLParam(r, "workerTypeID"))
	if err != nil || workerTypeID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid worker type id")
		return uuid.Nil, false
	}
	return workerTypeID, true
}

func (h *Handler) workerClassificationRuleRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	ruleID, ok := h.workerClassificationRuleIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, ruleID, true
}

func (h *Handler) superAdminWorkerClassificationRuleRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	ruleID, ok := h.workerClassificationRuleIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, ruleID, true
}

func (h *Handler) workerClassificationRuleIDFromURL(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	ruleID, err := uuid.Parse(chi.URLParam(r, "ruleID"))
	if err != nil || ruleID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid worker classification rule id")
		return uuid.Nil, false
	}
	return ruleID, true
}
