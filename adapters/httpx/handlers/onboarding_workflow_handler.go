package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListOnboardingWorkflows(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list onboarding workflows", err, "tenant context is required")
		return
	}
	h.listOnboardingWorkflows(w, r, tenantID, "list onboarding workflows")
}

func (h *Handler) CreateOnboardingWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create onboarding workflow", err, "tenant context is required")
		return
	}
	h.createOnboardingWorkflow(w, r, tenantID, "create onboarding workflow")
}

func (h *Handler) UpdateOnboardingWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.onboardingWorkflowRequestIDs(w, r, "update onboarding workflow")
	if !ok {
		return
	}
	h.updateOnboardingWorkflow(w, r, tenantID, id, "update onboarding workflow")
}

func (h *Handler) DeleteOnboardingWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.onboardingWorkflowRequestIDs(w, r, "delete onboarding workflow")
	if !ok {
		return
	}
	if err := h.svc.DeleteOnboardingWorkflow(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete onboarding workflow", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateOnboardingTask(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.onboardingWorkflowRequestIDs(w, r, "create onboarding task")
	if !ok {
		return
	}
	h.createOnboardingTask(w, r, tenantID, workflowID, "create onboarding task")
}

func (h *Handler) UpdateOnboardingTask(w http.ResponseWriter, r *http.Request) {
	tenantID, taskID, ok := h.onboardingTaskRequestIDs(w, r, "update onboarding task")
	if !ok {
		return
	}
	h.updateOnboardingTask(w, r, tenantID, taskID, "update onboarding task")
}

func (h *Handler) DeleteOnboardingTask(w http.ResponseWriter, r *http.Request) {
	tenantID, taskID, ok := h.onboardingTaskRequestIDs(w, r, "delete onboarding task")
	if !ok {
		return
	}
	if err := h.svc.DeleteOnboardingTask(r.Context(), tenantID, taskID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete onboarding task", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListOnboardingWorkflowAssignments(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list onboarding assignments", err, "tenant context is required")
		return
	}
	h.listOnboardingAssignments(w, r, tenantID, "list onboarding assignments")
}

func (h *Handler) CreateOnboardingWorkflowAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create onboarding assignment", err, "tenant context is required")
		return
	}
	h.createOnboardingAssignment(w, r, tenantID, "create onboarding assignment")
}

func (h *Handler) UpdateOnboardingWorkflowAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.onboardingAssignmentRequestIDs(w, r, "update onboarding assignment")
	if !ok {
		return
	}
	h.updateOnboardingAssignment(w, r, tenantID, id, "update onboarding assignment")
}

func (h *Handler) DeleteOnboardingWorkflowAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.onboardingAssignmentRequestIDs(w, r, "delete onboarding assignment")
	if !ok {
		return
	}
	if err := h.svc.DeleteOnboardingWorkflowAssignment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete onboarding assignment", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantOnboardingWorkflows(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant onboarding workflows")
	if !ok {
		return
	}
	h.listOnboardingWorkflows(w, r, tenantID, "list tenant onboarding workflows")
}

func (h *Handler) CreateTenantOnboardingWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant onboarding workflow")
	if !ok {
		return
	}
	h.createOnboardingWorkflow(w, r, tenantID, "create tenant onboarding workflow")
}

func (h *Handler) UpdateTenantOnboardingWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "onboardingWorkflowID", "update tenant onboarding workflow")
	if !ok {
		return
	}
	h.updateOnboardingWorkflow(w, r, tenantID, id, "update tenant onboarding workflow")
}

func (h *Handler) DeleteTenantOnboardingWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "onboardingWorkflowID", "delete tenant onboarding workflow")
	if !ok {
		return
	}
	if err := h.svc.DeleteOnboardingWorkflow(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant onboarding workflow", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantOnboardingTask(w http.ResponseWriter, r *http.Request) {
	tenantID, workflowID, ok := h.superAdminLookupRequestIDs(w, r, "onboardingWorkflowID", "create tenant onboarding task")
	if !ok {
		return
	}
	h.createOnboardingTask(w, r, tenantID, workflowID, "create tenant onboarding task")
}

func (h *Handler) UpdateTenantOnboardingTask(w http.ResponseWriter, r *http.Request) {
	tenantID, taskID, ok := h.superAdminLookupRequestIDs(w, r, "onboardingTaskID", "update tenant onboarding task")
	if !ok {
		return
	}
	h.updateOnboardingTask(w, r, tenantID, taskID, "update tenant onboarding task")
}

func (h *Handler) DeleteTenantOnboardingTask(w http.ResponseWriter, r *http.Request) {
	tenantID, taskID, ok := h.superAdminLookupRequestIDs(w, r, "onboardingTaskID", "delete tenant onboarding task")
	if !ok {
		return
	}
	if err := h.svc.DeleteOnboardingTask(r.Context(), tenantID, taskID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant onboarding task", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantOnboardingWorkflowAssignments(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant onboarding assignments")
	if !ok {
		return
	}
	h.listOnboardingAssignments(w, r, tenantID, "list tenant onboarding assignments")
}

func (h *Handler) CreateTenantOnboardingWorkflowAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant onboarding assignment")
	if !ok {
		return
	}
	h.createOnboardingAssignment(w, r, tenantID, "create tenant onboarding assignment")
}

func (h *Handler) UpdateTenantOnboardingWorkflowAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "onboardingAssignmentID", "update tenant onboarding assignment")
	if !ok {
		return
	}
	h.updateOnboardingAssignment(w, r, tenantID, id, "update tenant onboarding assignment")
}

func (h *Handler) DeleteTenantOnboardingWorkflowAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "onboardingAssignmentID", "delete tenant onboarding assignment")
	if !ok {
		return
	}
	if err := h.svc.DeleteOnboardingWorkflowAssignment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant onboarding assignment", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createOnboardingWorkflow(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.OnboardingWorkflowCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateOnboardingWorkflow(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listOnboardingWorkflows(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListOnboardingWorkflows(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list workflows")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateOnboardingWorkflow(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.OnboardingWorkflowCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateOnboardingWorkflow(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createOnboardingTask(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, workflowID uuid.UUID, operation string) {
	var cmd ports.OnboardingTaskCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.WorkflowID = workflowID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateOnboardingTask(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateOnboardingTask(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.OnboardingTaskCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateOnboardingTask(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listOnboardingAssignments(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListOnboardingWorkflowAssignments(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list assignments")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createOnboardingAssignment(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.OnboardingAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateOnboardingWorkflowAssignment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateOnboardingAssignment(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.OnboardingAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateOnboardingWorkflowAssignment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) onboardingWorkflowRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "onboardingWorkflowID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid onboarding workflow id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) onboardingTaskRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "onboardingTaskID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid onboarding task id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) onboardingAssignmentRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "onboardingAssignmentID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid onboarding assignment id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
