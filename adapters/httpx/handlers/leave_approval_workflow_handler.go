package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListLeaveApprovalWorkflows(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leave approval workflows", err, "tenant context is required")
		return
	}
	h.listLeaveApprovalWorkflowsForTenant(w, r, tenantID, "list leave approval workflows")
}

func (h *Handler) CreateLeaveApprovalWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create leave approval workflow", err, "tenant context is required")
		return
	}
	h.createLeaveApprovalWorkflowForTenant(w, r, tenantID, "create leave approval workflow")
}

func (h *Handler) UpdateLeaveApprovalWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.leaveApprovalWorkflowRequestIDs(w, r, "update leave approval workflow")
	if ok {
		h.updateLeaveApprovalWorkflowForTenant(w, r, tenantID, workflowID, "update leave approval workflow")
	}
}

func (h *Handler) DeleteLeaveApprovalWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.leaveApprovalWorkflowRequestIDs(w, r, "delete leave approval workflow")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeaveApprovalWorkflow(r.Context(), tenantID, workflowID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete leave approval workflow", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListLeaveApprovalWorkflowSteps(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.leaveApprovalWorkflowRequestIDs(w, r, "list leave approval workflow steps")
	if ok {
		h.listLeaveApprovalWorkflowStepsForTenant(w, r, tenantID, workflowID, "list leave approval workflow steps")
	}
}

func (h *Handler) CreateLeaveApprovalWorkflowStep(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.leaveApprovalWorkflowRequestIDs(w, r, "create leave approval workflow step")
	if ok {
		h.createLeaveApprovalWorkflowStepForTenant(w, r, tenantID, workflowID, "create leave approval workflow step")
	}
}

func (h *Handler) UpdateLeaveApprovalWorkflowStep(w http.ResponseWriter, r *http.Request) {
	tenantID, stepID, ok := h.leaveApprovalWorkflowStepRequestIDs(w, r, "update leave approval workflow step")
	if ok {
		h.updateLeaveApprovalWorkflowStepForTenant(w, r, tenantID, stepID, "update leave approval workflow step")
	}
}

func (h *Handler) DeleteLeaveApprovalWorkflowStep(w http.ResponseWriter, r *http.Request) {
	tenantID, stepID, ok := h.leaveApprovalWorkflowStepRequestIDs(w, r, "delete leave approval workflow step")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeaveApprovalWorkflowStep(r.Context(), tenantID, stepID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete leave approval workflow step", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListLeaveApprovals(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list leave approvals", err, "tenant context is required")
		return
	}
	h.listLeaveApprovalsForTenant(w, r, tenantID, "list leave approvals")
}

func (h *Handler) ApproveLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, approvalID, ok := h.leaveApprovalRequestIDs(w, r, "approve leave")
	if ok {
		h.approveLeaveForTenant(w, r, tenantID, approvalID, "approve leave")
	}
}

func (h *Handler) RejectLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, approvalID, ok := h.leaveApprovalRequestIDs(w, r, "reject leave")
	if ok {
		h.rejectLeaveForTenant(w, r, tenantID, approvalID, "reject leave")
	}
}

func (h *Handler) CancelLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, leaveID, ok := h.leaveRequestIDs(w, r, "cancel leave")
	if ok {
		h.cancelLeaveForTenant(w, r, tenantID, leaveID, "cancel leave")
	}
}

func (h *Handler) ListTenantLeaveApprovalWorkflows(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leave approval workflows")
	if ok {
		h.listLeaveApprovalWorkflowsForTenant(w, r, tenantID, "list tenant leave approval workflows")
	}
}
func (h *Handler) CreateTenantLeaveApprovalWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant leave approval workflow")
	if ok {
		h.createLeaveApprovalWorkflowForTenant(w, r, tenantID, "create tenant leave approval workflow")
	}
}
func (h *Handler) UpdateTenantLeaveApprovalWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.superAdminLeaveApprovalWorkflowRequestIDs(w, r, "update tenant leave approval workflow")
	if ok {
		h.updateLeaveApprovalWorkflowForTenant(w, r, tenantID, workflowID, "update tenant leave approval workflow")
	}
}
func (h *Handler) DeleteTenantLeaveApprovalWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.superAdminLeaveApprovalWorkflowRequestIDs(w, r, "delete tenant leave approval workflow")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeaveApprovalWorkflow(r.Context(), tenantID, workflowID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant leave approval workflow", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) ListTenantLeaveApprovalWorkflowSteps(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.superAdminLeaveApprovalWorkflowRequestIDs(w, r, "list tenant leave approval workflow steps")
	if ok {
		h.listLeaveApprovalWorkflowStepsForTenant(w, r, tenantID, workflowID, "list tenant leave approval workflow steps")
	}
}
func (h *Handler) CreateTenantLeaveApprovalWorkflowStep(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.superAdminLeaveApprovalWorkflowRequestIDs(w, r, "create tenant leave approval workflow step")
	if ok {
		h.createLeaveApprovalWorkflowStepForTenant(w, r, tenantID, workflowID, "create tenant leave approval workflow step")
	}
}
func (h *Handler) UpdateTenantLeaveApprovalWorkflowStep(w http.ResponseWriter, r *http.Request) {
	tenantID, stepID, ok := h.superAdminLeaveApprovalWorkflowStepRequestIDs(w, r, "update tenant leave approval workflow step")
	if ok {
		h.updateLeaveApprovalWorkflowStepForTenant(w, r, tenantID, stepID, "update tenant leave approval workflow step")
	}
}
func (h *Handler) DeleteTenantLeaveApprovalWorkflowStep(w http.ResponseWriter, r *http.Request) {
	tenantID, stepID, ok := h.superAdminLeaveApprovalWorkflowStepRequestIDs(w, r, "delete tenant leave approval workflow step")
	if !ok {
		return
	}
	if err := h.svc.DeleteLeaveApprovalWorkflowStep(r.Context(), tenantID, stepID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant leave approval workflow step", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) ListTenantLeaveApprovals(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant leave approvals")
	if ok {
		h.listLeaveApprovalsForTenant(w, r, tenantID, "list tenant leave approvals")
	}
}
func (h *Handler) ApproveTenantLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, approvalID, ok := h.superAdminLeaveApprovalRequestIDs(w, r, "approve tenant leave")
	if ok {
		h.approveLeaveForTenant(w, r, tenantID, approvalID, "approve tenant leave")
	}
}

func (h *Handler) RejectTenantLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, approvalID, ok := h.superAdminLeaveApprovalRequestIDs(w, r, "reject tenant leave")
	if ok {
		h.rejectLeaveForTenant(w, r, tenantID, approvalID, "reject tenant leave")
	}
}

func (h *Handler) CancelTenantLeave(w http.ResponseWriter, r *http.Request) {
	tenantID, leaveID, ok := h.superAdminLeaveRequestIDs(w, r, "cancel tenant leave")
	if ok {
		h.cancelLeaveForTenant(w, r, tenantID, leaveID, "cancel tenant leave")
	}
}

func (h *Handler) listLeaveApprovalWorkflowsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListLeaveApprovalWorkflows(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave approval workflows")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createLeaveApprovalWorkflowForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LeaveApprovalWorkflowCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLeaveApprovalWorkflow(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateLeaveApprovalWorkflowForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, workflowID uuid.UUID, operation string) {
	var cmd ports.LeaveApprovalWorkflowCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = workflowID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLeaveApprovalWorkflow(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listLeaveApprovalWorkflowStepsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, workflowID uuid.UUID, operation string) {
	items, err := h.svc.ListLeaveApprovalWorkflowSteps(r.Context(), tenantID, workflowID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave approval workflow steps")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createLeaveApprovalWorkflowStepForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, workflowID uuid.UUID, operation string) {
	var cmd ports.LeaveApprovalWorkflowStepCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.WorkflowID = workflowID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLeaveApprovalWorkflowStep(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateLeaveApprovalWorkflowStepForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, stepID uuid.UUID, operation string) {
	var cmd ports.LeaveApprovalWorkflowStepCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = stepID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLeaveApprovalWorkflowStep(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listLeaveApprovalsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	approverID := h.actorIDFromRequest(r)
	if raw := r.URL.Query().Get("approver_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "parse approver id", err, "invalid approver_id")
			return
		}
		approverID = &parsed
	}
	if approverID == nil || *approverID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, nil, "approver_id is required")
		return
	}
	items, err := h.svc.ListPendingApprovalsByApprover(r.Context(), tenantID, *approverID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list leave approvals")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) approveLeaveForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, approvalID uuid.UUID, operation string) {
	var cmd ports.ApproveLeaveCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ApprovalID = approvalID
	if cmd.ApproverID == uuid.Nil {
		if actorID := h.actorIDFromRequest(r); actorID != nil {
			cmd.ApproverID = *actorID
		}
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ApproveLeave(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) rejectLeaveForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, approvalID uuid.UUID, operation string) {
	var cmd ports.RejectLeaveCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ApprovalID = approvalID
	if cmd.ApproverID == uuid.Nil {
		if actorID := h.actorIDFromRequest(r); actorID != nil {
			cmd.ApproverID = *actorID
		}
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.RejectLeave(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) cancelLeaveForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, leaveID uuid.UUID, operation string) {
	var cmd ports.CancelLeaveCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.LeaveID = leaveID
	if cmd.UserID == uuid.Nil {
		if actorID := h.actorIDFromRequest(r); actorID != nil {
			cmd.UserID = *actorID
		}
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CancelLeave(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) leaveApprovalWorkflowRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	workflowID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse workflow id", err, "invalid workflow id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, workflowID, true
}

func (h *Handler) superAdminLeaveApprovalWorkflowRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	workflowID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse workflow id", err, "invalid workflow id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, workflowID, true
}

func (h *Handler) leaveApprovalWorkflowStepRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	stepID, err := uuid.Parse(chi.URLParam(r, "stepID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse step id", err, "invalid step id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, stepID, true
}

func (h *Handler) superAdminLeaveApprovalWorkflowStepRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	stepID, err := uuid.Parse(chi.URLParam(r, "stepID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse step id", err, "invalid step id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, stepID, true
}

func (h *Handler) leaveApprovalRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	approvalID, err := uuid.Parse(chi.URLParam(r, "approvalID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse approval id", err, "invalid approval id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, approvalID, true
}

func (h *Handler) superAdminLeaveApprovalRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	approvalID, err := uuid.Parse(chi.URLParam(r, "approvalID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse approval id", err, "invalid approval id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, approvalID, true
}

func (h *Handler) leaveRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	leaveID, err := uuid.Parse(chi.URLParam(r, "leaveID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse leave id", err, "invalid leave id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, leaveID, true
}

func (h *Handler) superAdminLeaveRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	leaveID, err := uuid.Parse(chi.URLParam(r, "leaveID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse leave id", err, "invalid leave id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, leaveID, true
}
