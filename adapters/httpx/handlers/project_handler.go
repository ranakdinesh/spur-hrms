package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create project", err, "tenant context is required")
		return
	}
	h.createProjectForTenant(w, r, tenantID, "create project")
}

func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list projects", err, "tenant context is required")
		return
	}
	h.listProjectsForTenant(w, r, tenantID, "list projects")
}

func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	tenantID, projectID, ok := h.projectRequestIDs(w, r, "get project")
	if !ok {
		return
	}
	item, err := h.svc.GetProject(r.Context(), tenantID, projectID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get project", err, "project not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	tenantID, projectID, ok := h.projectRequestIDs(w, r, "update project")
	if !ok {
		return
	}
	h.updateProjectForTenant(w, r, tenantID, projectID, "update project")
}

func (h *Handler) UpdateProjectStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, projectID, ok := h.projectRequestIDs(w, r, "update project status")
	if !ok {
		return
	}
	h.updateProjectStatusForTenant(w, r, tenantID, projectID, "update project status")
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	tenantID, projectID, ok := h.projectRequestIDs(w, r, "delete project")
	if !ok {
		return
	}
	if err := h.svc.DeleteProject(r.Context(), tenantID, projectID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete project", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create project milestone", err, "tenant context is required")
		return
	}
	h.createProjectMilestoneForTenant(w, r, tenantID, "create project milestone")
}

func (h *Handler) ListProjectMilestones(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list project milestones", err, "tenant context is required")
		return
	}
	h.listProjectMilestonesForTenant(w, r, tenantID, "list project milestones")
}

func (h *Handler) GetProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.projectMilestoneRequestIDs(w, r, "get project milestone")
	if !ok {
		return
	}
	item, err := h.svc.GetProjectMilestone(r.Context(), tenantID, milestoneID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get project milestone", err, "project milestone not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.projectMilestoneRequestIDs(w, r, "update project milestone")
	if !ok {
		return
	}
	h.updateProjectMilestoneForTenant(w, r, tenantID, milestoneID, "update project milestone")
}

func (h *Handler) SubmitProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.projectMilestoneRequestIDs(w, r, "submit project milestone")
	if !ok {
		return
	}
	item, err := h.svc.SubmitProjectMilestone(r.Context(), tenantID, milestoneID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "submit project milestone", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ReviewProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.projectMilestoneRequestIDs(w, r, "review project milestone")
	if !ok {
		return
	}
	h.reviewProjectMilestoneForTenant(w, r, tenantID, milestoneID, "review project milestone")
}

func (h *Handler) DeleteProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.projectMilestoneRequestIDs(w, r, "delete project milestone")
	if !ok {
		return
	}
	if err := h.svc.DeleteProjectMilestone(r.Context(), tenantID, milestoneID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete project milestone", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListProjectMilestoneEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.projectMilestoneRequestIDs(w, r, "list project milestone events")
	if !ok {
		return
	}
	items, err := h.svc.ListProjectMilestoneEvents(r.Context(), tenantID, milestoneID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list project milestone events", err, "failed to list project milestone events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateTenantProject(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant project")
	if !ok {
		return
	}
	h.createProjectForTenant(w, r, tenantID, "create tenant project")
}

func (h *Handler) ListTenantProjects(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant projects")
	if !ok {
		return
	}
	h.listProjectsForTenant(w, r, tenantID, "list tenant projects")
}

func (h *Handler) GetTenantProject(w http.ResponseWriter, r *http.Request) {
	tenantID, projectID, ok := h.superAdminProjectRequestIDs(w, r, "get tenant project")
	if !ok {
		return
	}
	item, err := h.svc.GetProject(r.Context(), tenantID, projectID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant project", err, "project not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantProject(w http.ResponseWriter, r *http.Request) {
	tenantID, projectID, ok := h.superAdminProjectRequestIDs(w, r, "update tenant project")
	if !ok {
		return
	}
	h.updateProjectForTenant(w, r, tenantID, projectID, "update tenant project")
}

func (h *Handler) UpdateTenantProjectStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, projectID, ok := h.superAdminProjectRequestIDs(w, r, "update tenant project status")
	if !ok {
		return
	}
	h.updateProjectStatusForTenant(w, r, tenantID, projectID, "update tenant project status")
}

func (h *Handler) DeleteTenantProject(w http.ResponseWriter, r *http.Request) {
	tenantID, projectID, ok := h.superAdminProjectRequestIDs(w, r, "delete tenant project")
	if !ok {
		return
	}
	if err := h.svc.DeleteProject(r.Context(), tenantID, projectID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant project", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant project milestone")
	if !ok {
		return
	}
	h.createProjectMilestoneForTenant(w, r, tenantID, "create tenant project milestone")
}

func (h *Handler) ListTenantProjectMilestones(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant project milestones")
	if !ok {
		return
	}
	h.listProjectMilestonesForTenant(w, r, tenantID, "list tenant project milestones")
}

func (h *Handler) GetTenantProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.superAdminProjectMilestoneRequestIDs(w, r, "get tenant project milestone")
	if !ok {
		return
	}
	item, err := h.svc.GetProjectMilestone(r.Context(), tenantID, milestoneID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant project milestone", err, "project milestone not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.superAdminProjectMilestoneRequestIDs(w, r, "update tenant project milestone")
	if !ok {
		return
	}
	h.updateProjectMilestoneForTenant(w, r, tenantID, milestoneID, "update tenant project milestone")
}

func (h *Handler) SubmitTenantProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.superAdminProjectMilestoneRequestIDs(w, r, "submit tenant project milestone")
	if !ok {
		return
	}
	item, err := h.svc.SubmitProjectMilestone(r.Context(), tenantID, milestoneID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "submit tenant project milestone", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ReviewTenantProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.superAdminProjectMilestoneRequestIDs(w, r, "review tenant project milestone")
	if !ok {
		return
	}
	h.reviewProjectMilestoneForTenant(w, r, tenantID, milestoneID, "review tenant project milestone")
}

func (h *Handler) DeleteTenantProjectMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.superAdminProjectMilestoneRequestIDs(w, r, "delete tenant project milestone")
	if !ok {
		return
	}
	if err := h.svc.DeleteProjectMilestone(r.Context(), tenantID, milestoneID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant project milestone", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantProjectMilestoneEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, milestoneID, ok := h.superAdminProjectMilestoneRequestIDs(w, r, "list tenant project milestone events")
	if !ok {
		return
	}
	items, err := h.svc.ListProjectMilestoneEvents(r.Context(), tenantID, milestoneID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant project milestone events", err, "failed to list project milestone events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createProjectForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ProjectCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateProject(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateProjectForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, projectID uuid.UUID, operation string) {
	var cmd ports.ProjectCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = projectID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateProject(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateProjectStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, projectID uuid.UUID, operation string) {
	var cmd ports.ProjectStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ProjectID = projectID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateProjectStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createProjectMilestoneForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ProjectMilestoneCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateProjectMilestone(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateProjectMilestoneForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, milestoneID uuid.UUID, operation string) {
	var cmd ports.ProjectMilestoneCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = milestoneID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateProjectMilestone(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) reviewProjectMilestoneForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, milestoneID uuid.UUID, operation string) {
	var cmd ports.ProjectMilestoneStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.MilestoneID = milestoneID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ReviewProjectMilestone(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listProjectsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	departmentID, ok := h.optionalUUIDQuery(w, r, "department_id", operation)
	if !ok {
		return
	}
	branchID, ok := h.optionalUUIDQuery(w, r, "branch_id", operation)
	if !ok {
		return
	}
	projectManagerID, ok := h.optionalUUIDQuery(w, r, "project_manager_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListProjects(r.Context(), domain.ProjectFilter{
		TenantID:         tenantID,
		Status:           optionalStringQuery(r, "status"),
		DepartmentID:     departmentID,
		BranchID:         branchID,
		ProjectManagerID: projectManagerID,
		Search:           optionalStringQuery(r, "search"),
	})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list projects")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listProjectMilestonesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	projectID, ok := h.optionalUUIDQuery(w, r, "project_id", operation)
	if !ok {
		return
	}
	engagementID, ok := h.optionalUUIDQuery(w, r, "engagement_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListProjectMilestones(r.Context(), domain.ProjectMilestoneFilter{
		TenantID:     tenantID,
		ProjectID:    projectID,
		EngagementID: engagementID,
		Status:       optionalStringQuery(r, "status"),
		Search:       optionalStringQuery(r, "search"),
	})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list project milestones")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) projectRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	projectID, ok := h.projectIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, projectID, true
}

func (h *Handler) superAdminProjectRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	projectID, ok := h.projectIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, projectID, true
}

func (h *Handler) projectMilestoneRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	milestoneID, ok := h.milestoneIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, milestoneID, true
}

func (h *Handler) superAdminProjectMilestoneRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	milestoneID, ok := h.milestoneIDFromURL(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, milestoneID, true
}

func (h *Handler) projectIDFromURL(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	value := chi.URLParam(r, "projectID")
	id, err := uuid.Parse(value)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid project id")
		return uuid.Nil, false
	}
	return id, true
}

func (h *Handler) milestoneIDFromURL(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	value := chi.URLParam(r, "milestoneID")
	id, err := uuid.Parse(value)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid milestone id")
		return uuid.Nil, false
	}
	return id, true
}
