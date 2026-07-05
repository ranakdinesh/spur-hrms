package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListWorkflowDefinitions(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.workflowTenantID(w, r, "list workflow definitions"); ok {
		h.listWorkflowDefinitionsForTenant(w, r, tenantID, "list workflow definitions")
	}
}

func (h *Handler) CreateWorkflowDefinition(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.workflowTenantID(w, r, "create workflow definition"); ok {
		h.createWorkflowDefinitionForTenant(w, r, tenantID, "create workflow definition")
	}
}

func (h *Handler) UpdateWorkflowDefinition(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowDefinitionID", "update workflow definition"); ok {
		h.updateWorkflowDefinitionForTenant(w, r, tenantID, id, "update workflow definition")
	}
}

func (h *Handler) ListWorkflowDefinitionSteps(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowDefinitionID", "list workflow definition steps"); ok {
		h.listWorkflowDefinitionStepsForTenant(w, r, tenantID, id, "list workflow definition steps")
	}
}

func (h *Handler) CreateWorkflowDefinitionStep(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowDefinitionID", "create workflow definition step"); ok {
		h.createWorkflowDefinitionStepForTenant(w, r, tenantID, id, "create workflow definition step")
	}
}

func (h *Handler) UpdateWorkflowDefinitionStep(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowStepID", "update workflow definition step"); ok {
		h.updateWorkflowDefinitionStepForTenant(w, r, tenantID, id, "update workflow definition step")
	}
}

func (h *Handler) ListOperationTemplates(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.workflowTenantID(w, r, "list operation templates"); ok {
		h.listOperationTemplatesForTenant(w, r, tenantID, "list operation templates")
	}
}

func (h *Handler) CreateOperationTemplate(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.workflowTenantID(w, r, "create operation template"); ok {
		h.createOperationTemplateForTenant(w, r, tenantID, "create operation template")
	}
}

func (h *Handler) UpdateOperationTemplate(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "operationTemplateID", "update operation template"); ok {
		h.updateOperationTemplateForTenant(w, r, tenantID, id, "update operation template")
	}
}

func (h *Handler) ListWorkflowTasks(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.workflowTenantID(w, r, "list workflow tasks"); ok {
		h.listWorkflowTasksForTenant(w, r, tenantID, "list workflow tasks")
	}
}

func (h *Handler) CreateWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.workflowTenantID(w, r, "create workflow task"); ok {
		h.createWorkflowTaskForTenant(w, r, tenantID, "create workflow task")
	}
}

func (h *Handler) GetWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowTaskID", "get workflow task"); ok {
		result, err := h.svc.GetWorkflowTaskWorkspace(r.Context(), tenantID, id)
		if err != nil {
			h.respondError(w, r, http.StatusNotFound, "get workflow task", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, result)
	}
}

func (h *Handler) UpdateWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowTaskID", "update workflow task"); ok {
		h.updateWorkflowTaskForTenant(w, r, tenantID, id, "update workflow task")
	}
}

func (h *Handler) ActWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowTaskID", "act workflow task"); ok {
		var cmd ports.WorkflowTaskActionCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode workflow action request", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.TaskID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		item, err := h.svc.ActWorkflowTask(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "act workflow task", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
	}
}

func (h *Handler) CreateWorkflowTaskComment(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowTaskID", "create workflow task comment"); ok {
		var cmd ports.WorkflowTaskCommentCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode workflow comment request", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.TaskID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		item, err := h.svc.CreateWorkflowTaskComment(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "create workflow task comment", err, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, item)
	}
}

func (h *Handler) CreateWorkflowTaskAttachment(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowTaskID", "create workflow task attachment"); ok {
		var cmd ports.WorkflowTaskAttachmentCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode workflow attachment request", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.TaskID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		item, err := h.svc.CreateWorkflowTaskAttachment(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "create workflow task attachment", err, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, item)
	}
}

func (h *Handler) WatchWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowTaskID", "watch workflow task"); ok {
		actorID := h.actorIDFromRequest(r)
		if actorID == nil {
			h.respondError(w, r, http.StatusUnauthorized, "watch workflow task", domain.ErrInvalidWorkflowTask, "user context is required")
			return
		}
		item, err := h.svc.WatchWorkflowTask(r.Context(), ports.WorkflowTaskWatchCommand{TenantID: tenantID, TaskID: id, WatcherUserID: *actorID, WatchReason: optionalStringQuery(r, "reason"), ActorID: actorID})
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "watch workflow task", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
	}
}

func (h *Handler) UnwatchWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.workflowTenantIDAndParam(w, r, "workflowTaskID", "unwatch workflow task"); ok {
		actorID := h.actorIDFromRequest(r)
		if actorID == nil {
			h.respondError(w, r, http.StatusUnauthorized, "unwatch workflow task", domain.ErrInvalidWorkflowTask, "user context is required")
			return
		}
		if err := h.svc.UnwatchWorkflowTask(r.Context(), ports.WorkflowTaskWatchCommand{TenantID: tenantID, TaskID: id, WatcherUserID: *actorID, ActorID: actorID}); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "unwatch workflow task", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) GetWorkflowTaskSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.workflowTenantID(w, r, "get workflow task summary"); ok {
		items, err := h.svc.GetWorkflowTaskSummary(r.Context(), tenantID)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "get workflow task summary", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, items)
	}
}

func (h *Handler) ListTenantWorkflowDefinitions(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant workflow definitions"); ok {
		h.listWorkflowDefinitionsForTenant(w, r, tenantID, "list tenant workflow definitions")
	}
}
func (h *Handler) CreateTenantWorkflowDefinition(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant workflow definition"); ok {
		h.createWorkflowDefinitionForTenant(w, r, tenantID, "create tenant workflow definition")
	}
}
func (h *Handler) UpdateTenantWorkflowDefinition(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowDefinitionID", "update tenant workflow definition"); ok {
		h.updateWorkflowDefinitionForTenant(w, r, tenantID, id, "update tenant workflow definition")
	}
}
func (h *Handler) ListTenantWorkflowDefinitionSteps(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowDefinitionID", "list tenant workflow definition steps"); ok {
		h.listWorkflowDefinitionStepsForTenant(w, r, tenantID, id, "list tenant workflow definition steps")
	}
}
func (h *Handler) CreateTenantWorkflowDefinitionStep(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowDefinitionID", "create tenant workflow definition step"); ok {
		h.createWorkflowDefinitionStepForTenant(w, r, tenantID, id, "create tenant workflow definition step")
	}
}
func (h *Handler) UpdateTenantWorkflowDefinitionStep(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowStepID", "update tenant workflow definition step"); ok {
		h.updateWorkflowDefinitionStepForTenant(w, r, tenantID, id, "update tenant workflow definition step")
	}
}
func (h *Handler) ListTenantOperationTemplates(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant operation templates"); ok {
		h.listOperationTemplatesForTenant(w, r, tenantID, "list tenant operation templates")
	}
}
func (h *Handler) CreateTenantOperationTemplate(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant operation template"); ok {
		h.createOperationTemplateForTenant(w, r, tenantID, "create tenant operation template")
	}
}
func (h *Handler) UpdateTenantOperationTemplate(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "operationTemplateID", "update tenant operation template"); ok {
		h.updateOperationTemplateForTenant(w, r, tenantID, id, "update tenant operation template")
	}
}
func (h *Handler) ListTenantWorkflowTasks(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant workflow tasks"); ok {
		h.listWorkflowTasksForTenant(w, r, tenantID, "list tenant workflow tasks")
	}
}
func (h *Handler) CreateTenantWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant workflow task"); ok {
		h.createWorkflowTaskForTenant(w, r, tenantID, "create tenant workflow task")
	}
}
func (h *Handler) GetTenantWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowTaskID", "get tenant workflow task"); ok {
		result, err := h.svc.GetWorkflowTaskWorkspace(r.Context(), tenantID, id)
		if err != nil {
			h.respondError(w, r, http.StatusNotFound, "get tenant workflow task", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, result)
	}
}
func (h *Handler) UpdateTenantWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowTaskID", "update tenant workflow task"); ok {
		h.updateWorkflowTaskForTenant(w, r, tenantID, id, "update tenant workflow task")
	}
}
func (h *Handler) ActTenantWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowTaskID", "act tenant workflow task"); ok {
		var cmd ports.WorkflowTaskActionCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode tenant workflow action request", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.TaskID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		item, err := h.svc.ActWorkflowTask(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "act tenant workflow task", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
	}
}
func (h *Handler) CreateTenantWorkflowTaskComment(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowTaskID", "create tenant workflow task comment"); ok {
		var cmd ports.WorkflowTaskCommentCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode tenant workflow comment request", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.TaskID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		item, err := h.svc.CreateWorkflowTaskComment(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "create tenant workflow task comment", err, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, item)
	}
}
func (h *Handler) CreateTenantWorkflowTaskAttachment(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowTaskID", "create tenant workflow task attachment"); ok {
		var cmd ports.WorkflowTaskAttachmentCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode tenant workflow attachment request", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.TaskID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		item, err := h.svc.CreateWorkflowTaskAttachment(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "create tenant workflow task attachment", err, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, item)
	}
}
func (h *Handler) WatchTenantWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowTaskID", "watch tenant workflow task"); ok {
		actorID := h.actorIDFromRequest(r)
		if actorID == nil {
			h.respondError(w, r, http.StatusUnauthorized, "watch tenant workflow task", domain.ErrInvalidWorkflowTask, "user context is required")
			return
		}
		item, err := h.svc.WatchWorkflowTask(r.Context(), ports.WorkflowTaskWatchCommand{TenantID: tenantID, TaskID: id, WatcherUserID: *actorID, WatchReason: optionalStringQuery(r, "reason"), ActorID: actorID})
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "watch tenant workflow task", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
	}
}
func (h *Handler) UnwatchTenantWorkflowTask(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superWorkflowID(w, r, "workflowTaskID", "unwatch tenant workflow task"); ok {
		actorID := h.actorIDFromRequest(r)
		if actorID == nil {
			h.respondError(w, r, http.StatusUnauthorized, "unwatch tenant workflow task", domain.ErrInvalidWorkflowTask, "user context is required")
			return
		}
		if err := h.svc.UnwatchWorkflowTask(r.Context(), ports.WorkflowTaskWatchCommand{TenantID: tenantID, TaskID: id, WatcherUserID: *actorID, ActorID: actorID}); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "unwatch tenant workflow task", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
func (h *Handler) GetTenantWorkflowTaskSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant workflow task summary"); ok {
		items, err := h.svc.GetWorkflowTaskSummary(r.Context(), tenantID)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "get tenant workflow task summary", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, items)
	}
}

func (h *Handler) listWorkflowDefinitionsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListWorkflowDefinitions(r.Context(), tenantID, optionalStringQuery(r, "status"), optionalStringQuery(r, "module_key"), optionalStringQuery(r, "search"), queryInt32(r, "limit", 50), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createWorkflowDefinitionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.WorkflowDefinitionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkflowDefinition(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateWorkflowDefinitionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.WorkflowDefinitionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkflowDefinition(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listWorkflowDefinitionStepsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, definitionID uuid.UUID, operation string) {
	items, err := h.svc.ListWorkflowDefinitionSteps(r.Context(), tenantID, definitionID)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createWorkflowDefinitionStepForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, definitionID uuid.UUID, operation string) {
	var cmd ports.WorkflowDefinitionStepCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.WorkflowDefinitionID, cmd.ActorID = tenantID, definitionID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkflowDefinitionStep(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateWorkflowDefinitionStepForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.WorkflowDefinitionStepCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkflowDefinitionStep(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listOperationTemplatesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListOperationTemplates(r.Context(), tenantID, optionalStringQuery(r, "category"), optionalStringQuery(r, "source_module"), optionalBoolQuery(r, "active_only"), optionalStringQuery(r, "search"), queryInt32(r, "limit", 50), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createOperationTemplateForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.OperationTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateOperationTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateOperationTemplateForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.OperationTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateOperationTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listWorkflowTasksForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListWorkflowTasks(r.Context(), domain.WorkflowTaskFilter{TenantID: tenantID, ViewKey: optionalStringQuery(r, "view"), Status: optionalStringQuery(r, "status"), Severity: optionalStringQuery(r, "severity"), SourceModule: optionalStringQuery(r, "source_module"), Search: optionalStringQuery(r, "search"), ViewerUserID: h.actorIDFromRequest(r), ViewerRole: optionalStringQuery(r, "viewer_role"), ViewerTeam: optionalStringQuery(r, "viewer_team"), Limit: queryInt32(r, "limit", 50), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createWorkflowTaskForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.WorkflowTaskCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkflowTask(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateWorkflowTaskForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.WorkflowTaskCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkflowTask(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) workflowTenantID(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, false
	}
	return tenantID, true
}

func (h *Handler) workflowTenantIDAndParam(w http.ResponseWriter, r *http.Request, idParam string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.workflowTenantID(w, r, operation)
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

func (h *Handler) superWorkflowID(w http.ResponseWriter, r *http.Request, idParam string, operation string) (uuid.UUID, uuid.UUID, bool) {
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
