package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListComplianceRules(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list compliance rules", err, "tenant context is required")
		return
	}
	h.listComplianceRulesForTenant(w, r, tenantID, "list compliance rules")
}

func (h *Handler) CreateComplianceRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create compliance rule", err, "tenant context is required")
		return
	}
	h.createComplianceRuleForTenant(w, r, tenantID, "create compliance rule")
}

func (h *Handler) GetComplianceRule(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.complianceRuleRequestIDs(w, r, "get compliance rule")
	if !ok {
		return
	}
	item, err := h.svc.GetComplianceRule(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get compliance rule", err, "compliance rule not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateComplianceRule(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.complianceRuleRequestIDs(w, r, "update compliance rule")
	if !ok {
		return
	}
	h.updateComplianceRuleForTenant(w, r, tenantID, id, "update compliance rule")
}

func (h *Handler) DeleteComplianceRule(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.complianceRuleRequestIDs(w, r, "delete compliance rule")
	if !ok {
		return
	}
	if err := h.svc.DeleteComplianceRule(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete compliance rule", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SeedComplianceRules(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "seed compliance rules", err, "tenant context is required")
		return
	}
	items, err := h.svc.SeedDefaultComplianceRules(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "seed compliance rules", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListComplianceChecklistItems(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list compliance checklist", err, "tenant context is required")
		return
	}
	h.listComplianceChecklistForTenant(w, r, tenantID, "list compliance checklist")
}

func (h *Handler) GenerateComplianceChecklist(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "generate compliance checklist", err, "tenant context is required")
		return
	}
	h.generateComplianceChecklistForTenant(w, r, tenantID, "generate compliance checklist")
}

func (h *Handler) UpdateComplianceChecklistStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.complianceChecklistRequestIDs(w, r, "update compliance checklist status")
	if !ok {
		return
	}
	h.updateComplianceChecklistStatusForTenant(w, r, tenantID, id, "update compliance checklist status")
}

func (h *Handler) UpdateComplianceChecklistEvidence(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.complianceChecklistRequestIDs(w, r, "update compliance evidence")
	if !ok {
		return
	}
	h.updateComplianceEvidenceForTenant(w, r, tenantID, id, "update compliance evidence")
}

func (h *Handler) WaiveComplianceChecklistItem(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.complianceChecklistRequestIDs(w, r, "waive compliance checklist item")
	if !ok {
		return
	}
	h.waiveComplianceChecklistForTenant(w, r, tenantID, id, "waive compliance checklist item")
}

func (h *Handler) DeleteComplianceChecklistItem(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.complianceChecklistRequestIDs(w, r, "delete compliance checklist item")
	if !ok {
		return
	}
	if err := h.svc.DeleteComplianceChecklistItem(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete compliance checklist item", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetComplianceSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get compliance summary", err, "tenant context is required")
		return
	}
	items, err := h.svc.GetComplianceSummary(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "get compliance summary", err, "failed to get compliance summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListComplianceEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list compliance events", err, "tenant context is required")
		return
	}
	h.listComplianceEventsForTenant(w, r, tenantID, "list compliance events")
}

func (h *Handler) ListTenantComplianceRules(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant compliance rules")
	if !ok {
		return
	}
	h.listComplianceRulesForTenant(w, r, tenantID, "list tenant compliance rules")
}

func (h *Handler) CreateTenantComplianceRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant compliance rule")
	if !ok {
		return
	}
	h.createComplianceRuleForTenant(w, r, tenantID, "create tenant compliance rule")
}

func (h *Handler) GetTenantComplianceRule(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminComplianceRuleRequestIDs(w, r, "get tenant compliance rule")
	if !ok {
		return
	}
	item, err := h.svc.GetComplianceRule(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant compliance rule", err, "compliance rule not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantComplianceRule(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminComplianceRuleRequestIDs(w, r, "update tenant compliance rule")
	if !ok {
		return
	}
	h.updateComplianceRuleForTenant(w, r, tenantID, id, "update tenant compliance rule")
}

func (h *Handler) DeleteTenantComplianceRule(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminComplianceRuleRequestIDs(w, r, "delete tenant compliance rule")
	if !ok {
		return
	}
	if err := h.svc.DeleteComplianceRule(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant compliance rule", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SeedTenantComplianceRules(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "seed tenant compliance rules")
	if !ok {
		return
	}
	items, err := h.svc.SeedDefaultComplianceRules(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "seed tenant compliance rules", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListTenantComplianceChecklistItems(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant compliance checklist")
	if !ok {
		return
	}
	h.listComplianceChecklistForTenant(w, r, tenantID, "list tenant compliance checklist")
}

func (h *Handler) GenerateTenantComplianceChecklist(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "generate tenant compliance checklist")
	if !ok {
		return
	}
	h.generateComplianceChecklistForTenant(w, r, tenantID, "generate tenant compliance checklist")
}

func (h *Handler) UpdateTenantComplianceChecklistStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminComplianceChecklistRequestIDs(w, r, "update tenant compliance checklist status")
	if !ok {
		return
	}
	h.updateComplianceChecklistStatusForTenant(w, r, tenantID, id, "update tenant compliance checklist status")
}

func (h *Handler) UpdateTenantComplianceChecklistEvidence(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminComplianceChecklistRequestIDs(w, r, "update tenant compliance evidence")
	if !ok {
		return
	}
	h.updateComplianceEvidenceForTenant(w, r, tenantID, id, "update tenant compliance evidence")
}

func (h *Handler) WaiveTenantComplianceChecklistItem(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminComplianceChecklistRequestIDs(w, r, "waive tenant compliance checklist item")
	if !ok {
		return
	}
	h.waiveComplianceChecklistForTenant(w, r, tenantID, id, "waive tenant compliance checklist item")
}

func (h *Handler) DeleteTenantComplianceChecklistItem(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminComplianceChecklistRequestIDs(w, r, "delete tenant compliance checklist item")
	if !ok {
		return
	}
	if err := h.svc.DeleteComplianceChecklistItem(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant compliance checklist item", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTenantComplianceSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant compliance summary")
	if !ok {
		return
	}
	items, err := h.svc.GetComplianceSummary(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "get tenant compliance summary", err, "failed to get compliance summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListTenantComplianceEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant compliance events")
	if !ok {
		return
	}
	h.listComplianceEventsForTenant(w, r, tenantID, "list tenant compliance events")
}

func (h *Handler) createComplianceRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ComplianceRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateComplianceRule(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateComplianceRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ComplianceRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateComplianceRule(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listComplianceRulesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	active := optionalBoolQuery(r, "is_active")
	items, err := h.svc.ListComplianceRules(r.Context(), domain.ComplianceRuleFilter{TenantID: tenantID, Category: optionalStringQuery(r, "category"), Scope: optionalStringQuery(r, "scope"), Severity: optionalStringQuery(r, "severity"), IsActive: active, Search: optionalStringQuery(r, "search")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list compliance rules")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listComplianceChecklistForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	engagementID, ok := h.optionalUUIDQuery(w, r, "engagement_id", operation)
	if !ok {
		return
	}
	ruleID, ok := h.optionalUUIDQuery(w, r, "rule_id", operation)
	if !ok {
		return
	}
	dueBefore, ok := optionalDateQuery(w, r, "due_before", operation, h)
	if !ok {
		return
	}
	items, err := h.svc.ListComplianceChecklistItems(r.Context(), domain.ComplianceChecklistFilter{TenantID: tenantID, WorkerProfileID: workerID, EngagementID: engagementID, RuleID: ruleID, Status: optionalStringQuery(r, "status"), Category: optionalStringQuery(r, "category"), DueBefore: dueBefore, Search: optionalStringQuery(r, "search")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list compliance checklist")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) generateComplianceChecklistForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ComplianceChecklistGenerateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	items, err := h.svc.GenerateComplianceChecklist(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateComplianceChecklistStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ComplianceChecklistStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateComplianceChecklistStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateComplianceEvidenceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ComplianceEvidenceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateComplianceChecklistEvidence(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) waiveComplianceChecklistForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ComplianceWaiverCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.WaiveComplianceChecklistItem(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listComplianceEventsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	itemID, ok := h.optionalUUIDQuery(w, r, "checklist_item_id", operation)
	if !ok {
		return
	}
	ruleID, ok := h.optionalUUIDQuery(w, r, "rule_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListComplianceEvents(r.Context(), tenantID, itemID, ruleID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list compliance events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) complianceRuleRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.complianceRuleIDFromURL(w, r, operation)
	return tenantID, id, ok
}

func (h *Handler) superAdminComplianceRuleRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.complianceRuleIDFromURL(w, r, operation)
	return tenantID, id, ok
}

func (h *Handler) complianceChecklistRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.complianceChecklistItemIDFromURL(w, r, operation)
	return tenantID, id, ok
}

func (h *Handler) superAdminComplianceChecklistRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.complianceChecklistItemIDFromURL(w, r, operation)
	return tenantID, id, ok
}

func (h *Handler) complianceRuleIDFromURL(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "complianceRuleID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid compliance rule id")
		return uuid.Nil, false
	}
	return id, true
}

func (h *Handler) complianceChecklistItemIDFromURL(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "checklistItemID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid checklist item id")
		return uuid.Nil, false
	}
	return id, true
}

func optionalDateQuery(w http.ResponseWriter, r *http.Request, key string, operation string, h *Handler) (*time.Time, bool) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, true
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid "+key)
		return nil, false
	}
	return &parsed, true
}
