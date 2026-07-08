package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
)

func (h *Handler) ListPolicyEnginePolicySets(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list policy engine policy sets", err, "tenant context is required")
		return
	}
	h.listPolicyEnginePolicySetsForTenant(w, r, tenantID, "list policy engine policy sets")
}

func (h *Handler) CreatePolicyEnginePolicySet(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create policy engine policy set", err, "tenant context is required")
		return
	}
	h.createPolicyEnginePolicySetForTenant(w, r, tenantID, "create policy engine policy set")
}

func (h *Handler) UpdatePolicyEnginePolicySet(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "update policy engine policy set", err, "tenant context is required")
		return
	}
	h.updatePolicyEnginePolicySetForTenant(w, r, tenantID, "update policy engine policy set")
}

func (h *Handler) DeletePolicyEnginePolicySet(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete policy engine policy set", err, "tenant context is required")
		return
	}
	h.deletePolicyEnginePolicySetForTenant(w, r, tenantID, "delete policy engine policy set")
}

func (h *Handler) ListPolicyEngineAssignments(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list policy engine assignments", err, "tenant context is required")
		return
	}
	h.listPolicyEngineAssignmentsForTenant(w, r, tenantID, "list policy engine assignments")
}

func (h *Handler) CreatePolicyEngineAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create policy engine assignment", err, "tenant context is required")
		return
	}
	h.createPolicyEngineAssignmentForTenant(w, r, tenantID, "create policy engine assignment")
}

func (h *Handler) UpdatePolicyEngineAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "update policy engine assignment", err, "tenant context is required")
		return
	}
	h.updatePolicyEngineAssignmentForTenant(w, r, tenantID, "update policy engine assignment")
}

func (h *Handler) DeletePolicyEngineAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete policy engine assignment", err, "tenant context is required")
		return
	}
	h.deletePolicyEngineAssignmentForTenant(w, r, tenantID, "delete policy engine assignment")
}

func (h *Handler) ListPolicyEngineLeaveRules(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list policy engine leave rules", err, "tenant context is required")
		return
	}
	h.listPolicyEngineLeaveRulesForTenant(w, r, tenantID, "list policy engine leave rules")
}

func (h *Handler) CreatePolicyEngineLeaveRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create policy engine leave rule", err, "tenant context is required")
		return
	}
	h.createPolicyEngineLeaveRuleForTenant(w, r, tenantID, "create policy engine leave rule")
}

func (h *Handler) UpdatePolicyEngineLeaveRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "update policy engine leave rule", err, "tenant context is required")
		return
	}
	h.updatePolicyEngineLeaveRuleForTenant(w, r, tenantID, "update policy engine leave rule")
}

func (h *Handler) DeletePolicyEngineLeaveRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete policy engine leave rule", err, "tenant context is required")
		return
	}
	h.deletePolicyEngineLeaveRuleForTenant(w, r, tenantID, "delete policy engine leave rule")
}

func (h *Handler) PreviewPolicyEngineEffectivePolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "preview effective policy", err, "tenant context is required")
		return
	}
	h.previewPolicyEngineEffectivePolicyForTenant(w, r, tenantID, "preview effective policy")
}

func (h *Handler) ListTenantPolicyEnginePolicySets(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant policy engine policy sets"); ok {
		h.listPolicyEnginePolicySetsForTenant(w, r, tenantID, "list tenant policy engine policy sets")
	}
}

func (h *Handler) CreateTenantPolicyEnginePolicySet(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant policy engine policy set"); ok {
		h.createPolicyEnginePolicySetForTenant(w, r, tenantID, "create tenant policy engine policy set")
	}
}

func (h *Handler) UpdateTenantPolicyEnginePolicySet(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant policy engine policy set"); ok {
		h.updatePolicyEnginePolicySetForTenant(w, r, tenantID, "update tenant policy engine policy set")
	}
}

func (h *Handler) DeleteTenantPolicyEnginePolicySet(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "delete tenant policy engine policy set"); ok {
		h.deletePolicyEnginePolicySetForTenant(w, r, tenantID, "delete tenant policy engine policy set")
	}
}

func (h *Handler) ListTenantPolicyEngineAssignments(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant policy engine assignments"); ok {
		h.listPolicyEngineAssignmentsForTenant(w, r, tenantID, "list tenant policy engine assignments")
	}
}

func (h *Handler) CreateTenantPolicyEngineAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant policy engine assignment"); ok {
		h.createPolicyEngineAssignmentForTenant(w, r, tenantID, "create tenant policy engine assignment")
	}
}

func (h *Handler) UpdateTenantPolicyEngineAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant policy engine assignment"); ok {
		h.updatePolicyEngineAssignmentForTenant(w, r, tenantID, "update tenant policy engine assignment")
	}
}

func (h *Handler) DeleteTenantPolicyEngineAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "delete tenant policy engine assignment"); ok {
		h.deletePolicyEngineAssignmentForTenant(w, r, tenantID, "delete tenant policy engine assignment")
	}
}

func (h *Handler) ListTenantPolicyEngineLeaveRules(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant policy engine leave rules"); ok {
		h.listPolicyEngineLeaveRulesForTenant(w, r, tenantID, "list tenant policy engine leave rules")
	}
}

func (h *Handler) CreateTenantPolicyEngineLeaveRule(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant policy engine leave rule"); ok {
		h.createPolicyEngineLeaveRuleForTenant(w, r, tenantID, "create tenant policy engine leave rule")
	}
}

func (h *Handler) UpdateTenantPolicyEngineLeaveRule(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "update tenant policy engine leave rule"); ok {
		h.updatePolicyEngineLeaveRuleForTenant(w, r, tenantID, "update tenant policy engine leave rule")
	}
}

func (h *Handler) DeleteTenantPolicyEngineLeaveRule(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "delete tenant policy engine leave rule"); ok {
		h.deletePolicyEngineLeaveRuleForTenant(w, r, tenantID, "delete tenant policy engine leave rule")
	}
}

func (h *Handler) PreviewTenantPolicyEngineEffectivePolicy(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "preview tenant effective policy"); ok {
		h.previewPolicyEngineEffectivePolicyForTenant(w, r, tenantID, "preview tenant effective policy")
	}
}

func (h *Handler) listPolicyEnginePolicySetsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	kind, ok := h.policyKindFromQuery(w, r, operation)
	if !ok {
		return
	}
	if !h.requirePolicyEnginePermission(w, r, operation, kind, false) {
		return
	}
	items, err := h.svc.ListPolicySets(r.Context(), tenantID, kind)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPolicyEnginePolicySetForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PolicySetCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	kind, err := domain.ValidatePolicyKind(cmd.PolicyKind)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	if !h.requirePolicyEnginePermission(w, r, operation, kind, true) {
		return
	}
	cmd.PolicyKind = kind
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePolicySet(r.Context(), cmd)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updatePolicyEnginePolicySetForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	policySetID, ok := h.policyEngineURLUUID(w, r, "policySetID", "policy set id", operation)
	if !ok {
		return
	}
	var cmd ports.PolicySetCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	kind, err := domain.ValidatePolicyKind(cmd.PolicyKind)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	if !h.requirePolicyEnginePermission(w, r, operation, kind, true) {
		return
	}
	cmd.ID = policySetID
	cmd.PolicyKind = kind
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePolicySet(r.Context(), cmd)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) deletePolicyEnginePolicySetForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	kind, ok := h.policyKindFromQuery(w, r, operation)
	if !ok {
		return
	}
	if !h.requirePolicyEnginePermission(w, r, operation, kind, true) {
		return
	}
	policySetID, ok := h.policyEngineURLUUID(w, r, "policySetID", "policy set id", operation)
	if !ok {
		return
	}
	if err := h.svc.DeletePolicySet(r.Context(), tenantID, policySetID, h.actorIDFromRequest(r)); err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listPolicyEngineAssignmentsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	kind, ok := h.policyKindFromQuery(w, r, operation)
	if !ok {
		return
	}
	if !h.requirePolicyEnginePermission(w, r, operation, kind, false) {
		return
	}
	items, err := h.svc.ListPolicyAssignments(r.Context(), tenantID, kind)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPolicyEngineAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PolicyAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	kind, err := domain.ValidatePolicyKind(cmd.PolicyKind)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	if !h.requirePolicyEnginePermission(w, r, operation, kind, true) {
		return
	}
	cmd.PolicyKind = kind
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePolicyAssignment(r.Context(), cmd)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updatePolicyEngineAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	assignmentID, ok := h.policyEngineURLUUID(w, r, "assignmentID", "assignment id", operation)
	if !ok {
		return
	}
	var cmd ports.PolicyAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	kind, err := domain.ValidatePolicyKind(cmd.PolicyKind)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	if !h.requirePolicyEnginePermission(w, r, operation, kind, true) {
		return
	}
	cmd.ID = assignmentID
	cmd.PolicyKind = kind
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePolicyAssignment(r.Context(), cmd)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) deletePolicyEngineAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	kind, ok := h.policyKindFromQuery(w, r, operation)
	if !ok {
		return
	}
	if !h.requirePolicyEnginePermission(w, r, operation, kind, true) {
		return
	}
	assignmentID, ok := h.policyEngineURLUUID(w, r, "assignmentID", "assignment id", operation)
	if !ok {
		return
	}
	if err := h.svc.DeletePolicyAssignment(r.Context(), tenantID, assignmentID, h.actorIDFromRequest(r)); err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listPolicyEngineLeaveRulesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if !h.requirePolicyEnginePermission(w, r, operation, domain.PolicyKindLeave, false) {
		return
	}
	policySetID, ok := h.policyEngineURLUUID(w, r, "policySetID", "policy set id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListLeavePolicyRules(r.Context(), tenantID, policySetID)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPolicyEngineLeaveRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if !h.requirePolicyEnginePermission(w, r, operation, domain.PolicyKindLeave, true) {
		return
	}
	policySetID, ok := h.policyEngineURLUUID(w, r, "policySetID", "policy set id", operation)
	if !ok {
		return
	}
	var cmd ports.LeavePolicyRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.PolicySetID = policySetID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLeavePolicyRule(r.Context(), cmd)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updatePolicyEngineLeaveRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if !h.requirePolicyEnginePermission(w, r, operation, domain.PolicyKindLeave, true) {
		return
	}
	ruleID, ok := h.policyEngineURLUUID(w, r, "ruleID", "leave policy rule id", operation)
	if !ok {
		return
	}
	var cmd ports.LeavePolicyRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = ruleID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLeavePolicyRule(r.Context(), cmd)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) deletePolicyEngineLeaveRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if !h.requirePolicyEnginePermission(w, r, operation, domain.PolicyKindLeave, true) {
		return
	}
	ruleID, ok := h.policyEngineURLUUID(w, r, "ruleID", "leave policy rule id", operation)
	if !ok {
		return
	}
	if err := h.svc.DeleteLeavePolicyRule(r.Context(), tenantID, ruleID, h.actorIDFromRequest(r)); err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) previewPolicyEngineEffectivePolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var body struct {
		PolicyKind      string     `json:"policy_kind"`
		EmployeeUserID  *uuid.UUID `json:"employee_user_id,omitempty"`
		DesignationID   *uuid.UUID `json:"designation_id,omitempty"`
		WorkforceTypeID *uuid.UUID `json:"workforce_type_id,omitempty"`
		DepartmentID    *uuid.UUID `json:"department_id,omitempty"`
		BranchID        *uuid.UUID `json:"branch_id,omitempty"`
		RoleCodes       []string   `json:"role_codes,omitempty"`
		Date            string     `json:"date,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	kind, err := domain.ValidatePolicyKind(body.PolicyKind)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	if !h.requirePolicyEnginePermission(w, r, operation, kind, false) {
		return
	}
	resolutionDate := time.Now().UTC()
	if body.Date != "" {
		parsed, err := time.Parse("2006-01-02", body.Date)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid date")
			return
		}
		resolutionDate = parsed
	}
	subject := domain.PolicyResolutionSubject{
		TenantID:        tenantID,
		DesignationID:   body.DesignationID,
		WorkforceTypeID: body.WorkforceTypeID,
		DepartmentID:    body.DepartmentID,
		BranchID:        body.BranchID,
		RoleCodes:       body.RoleCodes,
		Date:            resolutionDate,
	}
	if body.EmployeeUserID != nil {
		subject.EmployeeUserID = *body.EmployeeUserID
	}
	result, err := h.svc.ResolveEffectivePolicySet(r.Context(), subject, kind)
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) policyKindFromQuery(w http.ResponseWriter, r *http.Request, operation string) (string, bool) {
	kind, err := domain.ValidatePolicyKind(r.URL.Query().Get("kind"))
	if err != nil {
		h.respondPolicyEngineError(w, r, operation, err)
		return "", false
	}
	return kind, true
}

func (h *Handler) requirePolicyEnginePermission(w http.ResponseWriter, r *http.Request, operation string, policyKind string, manage bool) bool {
	var permission string
	switch policyKind {
	case domain.PolicyKindAttendance:
		if manage {
			permission = permissions.AttendancePolicyManage
		} else {
			permission = permissions.AttendancePolicyView
		}
	case domain.PolicyKindLeave:
		if manage {
			permission = permissions.LeavePolicyManage
		} else {
			permission = permissions.LeavePolicyView
		}
	default:
		h.respondPolicyEngineError(w, r, operation, domain.ErrInvalidPolicyKind)
		return false
	}
	return h.requirePermission(w, r, operation, permission)
}

func (h *Handler) policyEngineURLUUID(w http.ResponseWriter, r *http.Request, param string, label string, operation string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid "+label)
		return uuid.Nil, false
	}
	return id, true
}

func (h *Handler) respondPolicyEngineError(w http.ResponseWriter, r *http.Request, operation string, err error) {
	switch {
	case errors.Is(err, domain.ErrPolicySetNotFound),
		errors.Is(err, domain.ErrPolicyAssignmentNotFound),
		errors.Is(err, domain.ErrLeavePolicyRuleNotFound):
		h.respondError(w, r, http.StatusNotFound, operation, err, err.Error())
	case errors.Is(err, domain.ErrInvalidTenantID),
		errors.Is(err, domain.ErrInvalidPolicySetID),
		errors.Is(err, domain.ErrInvalidPolicyKind),
		errors.Is(err, domain.ErrInvalidPolicyCode),
		errors.Is(err, domain.ErrInvalidPolicyName),
		errors.Is(err, domain.ErrInvalidPolicyAssignmentID),
		errors.Is(err, domain.ErrInvalidPolicyScope),
		errors.Is(err, domain.ErrInvalidLeavePolicyRuleID),
		errors.Is(err, domain.ErrInvalidLeavePolicyType),
		errors.Is(err, domain.ErrInvalidEmployeeUserID):
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
	default:
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "policy engine operation failed")
	}
}
