package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListLeavePolicyTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leave policy templates", err, "tenant context is required")
		return
	}
	h.listLeavePolicyTemplatesForTenant(w, r, tenantID, "list leave policy templates")
}

func (h *Handler) CreateLeavePolicyTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create leave policy template", err, "tenant context is required")
		return
	}
	h.createLeavePolicyTemplateForTenant(w, r, tenantID, "create leave policy template")
}

func (h *Handler) UpdateLeavePolicyTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.leaveTemplateRequestIDs(w, r, "update leave policy template")
	if !ok {
		return
	}
	h.updateLeavePolicyTemplateForTenant(w, r, tenantID, id, "update leave policy template")
}

func (h *Handler) DeleteLeavePolicyTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.leaveTemplateRequestIDs(w, r, "delete leave policy template")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeavePolicyTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete leave policy template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListLeavePolicyTemplateRules(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.leaveTemplateRuleParentIDs(w, r, "list leave policy template rules")
	if !ok {
		return
	}
	h.listLeavePolicyTemplateRulesForTenant(w, r, tenantID, templateID, "list leave policy template rules")
}

func (h *Handler) CreateLeavePolicyTemplateRule(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.leaveTemplateRuleParentIDs(w, r, "create leave policy template rule")
	if !ok {
		return
	}
	h.createLeavePolicyTemplateRuleForTenant(w, r, tenantID, templateID, "create leave policy template rule")
}

func (h *Handler) UpdateLeavePolicyTemplateRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.leaveTemplateRuleRequestIDs(w, r, "update leave policy template rule")
	if !ok {
		return
	}
	h.updateLeavePolicyTemplateRuleForTenant(w, r, tenantID, ruleID, "update leave policy template rule")
}

func (h *Handler) DeleteLeavePolicyTemplateRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.leaveTemplateRuleRequestIDs(w, r, "delete leave policy template rule")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeavePolicyTemplateRule(r.Context(), tenantID, ruleID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete leave policy template rule", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpsertEmployeeLeavePolicyAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert leave policy assignment", err, "tenant context is required")
		return
	}
	h.upsertEmployeeLeavePolicyAssignmentForTenant(w, r, tenantID, "upsert leave policy assignment")
}

func (h *Handler) ListEmployeeLeavePolicyAssignments(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leave policy assignments", err, "tenant context is required")
		return
	}
	h.listEmployeeLeavePolicyAssignmentsForTenant(w, r, tenantID, "list leave policy assignments")
}

func (h *Handler) DeleteEmployeeLeavePolicyAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.assignmentRequestIDs(w, r, "delete leave policy assignment")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeLeavePolicyAssignment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete leave policy assignment", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpsertLeaveBalance(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert leave balance", err, "tenant context is required")
		return
	}
	h.upsertLeaveBalanceForTenant(w, r, tenantID, "upsert leave balance")
}

func (h *Handler) ListLeaveBalances(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leave balances", err, "tenant context is required")
		return
	}
	h.listLeaveBalancesForTenant(w, r, tenantID, "list leave balances")
}

func (h *Handler) AdjustLeaveBalance(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "adjust leave balance", err, "tenant context is required")
		return
	}
	h.adjustLeaveBalanceForTenant(w, r, tenantID, "adjust leave balance")
}

func (h *Handler) ListLeaveLedger(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leave ledger", err, "tenant context is required")
		return
	}
	h.listLeaveLedgerForTenant(w, r, tenantID, "list leave ledger")
}

func (h *Handler) RunLeaveAccrual(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "run leave accrual", err, "tenant context is required")
		return
	}
	h.runLeaveAccrualForTenant(w, r, tenantID, "run leave accrual")
}

func (h *Handler) ListTenantLeavePolicyTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leave policy templates")
	if ok {
		h.listLeavePolicyTemplatesForTenant(w, r, tenantID, "list tenant leave policy templates")
	}
}
func (h *Handler) CreateTenantLeavePolicyTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant leave policy template")
	if ok {
		h.createLeavePolicyTemplateForTenant(w, r, tenantID, "create tenant leave policy template")
	}
}
func (h *Handler) UpdateTenantLeavePolicyTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLeaveTemplateRequestIDs(w, r, "update tenant leave policy template")
	if ok {
		h.updateLeavePolicyTemplateForTenant(w, r, tenantID, id, "update tenant leave policy template")
	}
}
func (h *Handler) DeleteTenantLeavePolicyTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLeaveTemplateRequestIDs(w, r, "delete tenant leave policy template")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeavePolicyTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant leave policy template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) ListTenantLeavePolicyTemplateRules(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.superAdminLeaveTemplateRuleParentIDs(w, r, "list tenant leave policy template rules")
	if ok {
		h.listLeavePolicyTemplateRulesForTenant(w, r, tenantID, templateID, "list tenant leave policy template rules")
	}
}
func (h *Handler) CreateTenantLeavePolicyTemplateRule(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.superAdminLeaveTemplateRuleParentIDs(w, r, "create tenant leave policy template rule")
	if ok {
		h.createLeavePolicyTemplateRuleForTenant(w, r, tenantID, templateID, "create tenant leave policy template rule")
	}
}
func (h *Handler) UpdateTenantLeavePolicyTemplateRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.superAdminLeaveTemplateRuleRequestIDs(w, r, "update tenant leave policy template rule")
	if ok {
		h.updateLeavePolicyTemplateRuleForTenant(w, r, tenantID, ruleID, "update tenant leave policy template rule")
	}
}
func (h *Handler) DeleteTenantLeavePolicyTemplateRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ruleID, ok := h.superAdminLeaveTemplateRuleRequestIDs(w, r, "delete tenant leave policy template rule")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeavePolicyTemplateRule(r.Context(), tenantID, ruleID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant leave policy template rule", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) UpsertTenantEmployeeLeavePolicyAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant leave policy assignment")
	if ok {
		h.upsertEmployeeLeavePolicyAssignmentForTenant(w, r, tenantID, "upsert tenant leave policy assignment")
	}
}
func (h *Handler) ListTenantEmployeeLeavePolicyAssignments(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leave policy assignments")
	if ok {
		h.listEmployeeLeavePolicyAssignmentsForTenant(w, r, tenantID, "list tenant leave policy assignments")
	}
}
func (h *Handler) DeleteTenantEmployeeLeavePolicyAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminAssignmentRequestIDs(w, r, "delete tenant leave policy assignment")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeLeavePolicyAssignment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant leave policy assignment", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) UpsertTenantLeaveBalance(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant leave balance")
	if ok {
		h.upsertLeaveBalanceForTenant(w, r, tenantID, "upsert tenant leave balance")
	}
}
func (h *Handler) ListTenantLeaveBalances(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leave balances")
	if ok {
		h.listLeaveBalancesForTenant(w, r, tenantID, "list tenant leave balances")
	}
}
func (h *Handler) AdjustTenantLeaveBalance(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "adjust tenant leave balance")
	if ok {
		h.adjustLeaveBalanceForTenant(w, r, tenantID, "adjust tenant leave balance")
	}
}
func (h *Handler) ListTenantLeaveLedger(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leave ledger")
	if ok {
		h.listLeaveLedgerForTenant(w, r, tenantID, "list tenant leave ledger")
	}
}
func (h *Handler) RunTenantLeaveAccrual(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "run tenant leave accrual")
	if ok {
		h.runLeaveAccrualForTenant(w, r, tenantID, "run tenant leave accrual")
	}
}

func (h *Handler) listLeavePolicyTemplatesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListLeavePolicyTemplates(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave policy templates")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) createLeavePolicyTemplateForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LeavePolicyTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLeavePolicyTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) updateLeavePolicyTemplateForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.LeavePolicyTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLeavePolicyTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) listLeavePolicyTemplateRulesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, templateID uuid.UUID, operation string) {
	items, err := h.svc.ListLeavePolicyTemplateRules(r.Context(), tenantID, templateID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave template rules")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) createLeavePolicyTemplateRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, templateID uuid.UUID, operation string) {
	var cmd ports.LeavePolicyTemplateRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.TemplateID = templateID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLeavePolicyTemplateRule(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) updateLeavePolicyTemplateRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, ruleID uuid.UUID, operation string) {
	var cmd ports.LeavePolicyTemplateRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = ruleID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLeavePolicyTemplateRule(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) upsertEmployeeLeavePolicyAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EmployeeLeavePolicyAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertEmployeeLeavePolicyAssignment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) listEmployeeLeavePolicyAssignmentsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse user id", err, "user_id query is required")
		return
	}
	items, err := h.svc.ListEmployeeLeavePolicyAssignments(r.Context(), tenantID, userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave policy assignments")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) upsertLeaveBalanceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LeaveBalanceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertLeaveBalance(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) listLeaveBalancesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if fyRaw := r.URL.Query().Get("fy_id"); fyRaw != "" {
		fyID, err := uuid.Parse(fyRaw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "parse fy id", err, "invalid fy_id")
			return
		}
		items, err := h.svc.ListLeaveBalancesByTenantFY(r.Context(), tenantID, fyID)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave balances")
			return
		}
		respondJSON(w, http.StatusOK, items)
		return
	}
	userRaw := r.URL.Query().Get("user_id")
	var userID uuid.UUID
	var err error
	if userRaw != "" {
		userID, err = uuid.Parse(userRaw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "parse user id", err, "invalid user_id")
			return
		}
	} else if actorID := h.actorIDFromRequest(r); actorID != nil {
		userID = *actorID
	}
	items, err := h.svc.ListLeaveBalancesByUser(r.Context(), tenantID, userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave balances")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) adjustLeaveBalanceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LeaveBalanceAdjustmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.AdjustLeaveBalance(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) listLeaveLedgerForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse user id", err, "user_id query is required")
		return
	}
	items, err := h.svc.ListLeaveLedgerByUser(r.Context(), tenantID, userID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave ledger")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) runLeaveAccrualForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.RunLeaveAccrualCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	items, err := h.svc.RunLeaveAccrual(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) leaveTemplateRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "templateID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse template id", err, "invalid template id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
func (h *Handler) superAdminLeaveTemplateRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "templateID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse template id", err, "invalid template id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
func (h *Handler) leaveTemplateRuleParentIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	templateID, err := uuid.Parse(chi.URLParam(r, "templateID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse template id", err, "invalid template id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, templateID, true
}
func (h *Handler) superAdminLeaveTemplateRuleParentIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	templateID, err := uuid.Parse(chi.URLParam(r, "templateID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse template id", err, "invalid template id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, templateID, true
}
func (h *Handler) leaveTemplateRuleRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	ruleID, err := uuid.Parse(chi.URLParam(r, "ruleID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse rule id", err, "invalid rule id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, ruleID, true
}
func (h *Handler) superAdminLeaveTemplateRuleRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	ruleID, err := uuid.Parse(chi.URLParam(r, "ruleID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse rule id", err, "invalid rule id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, ruleID, true
}
func (h *Handler) assignmentRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "assignmentID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse assignment id", err, "invalid assignment id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
func (h *Handler) superAdminAssignmentRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "assignmentID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse assignment id", err, "invalid assignment id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
