package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListProjectSkillRequirements(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list project skill requirements", err, "tenant context is required")
		return
	}
	h.listProjectSkillRequirementsForTenant(w, r, tenantID, "list project skill requirements")
}

func (h *Handler) CreateProjectSkillRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create project skill requirement", err, "tenant context is required")
		return
	}
	h.createProjectSkillRequirementForTenant(w, r, tenantID, "create project skill requirement")
}

func (h *Handler) UpdateProjectSkillRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.projectSkillRequirementRequestIDs(w, r, "update project skill requirement")
	if !ok {
		return
	}
	h.updateProjectSkillRequirementForTenant(w, r, tenantID, id, "update project skill requirement")
}

func (h *Handler) DeleteProjectSkillRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.projectSkillRequirementRequestIDs(w, r, "delete project skill requirement")
	if !ok {
		return
	}
	if err := h.svc.DeleteProjectSkillRequirement(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete project skill requirement", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListProjectSkillGapRows(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list project skill gaps", err, "tenant context is required")
		return
	}
	h.listProjectSkillGapRowsForTenant(w, r, tenantID, "list project skill gaps")
}

func (h *Handler) ListSkillGapSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list skill gap summary", err, "tenant context is required")
		return
	}
	h.listSkillGapSummaryForTenant(w, r, tenantID, "list skill gap summary")
}

func (h *Handler) ListSinglePersonSkillDependencies(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list single person skill dependencies", err, "tenant context is required")
		return
	}
	h.listSinglePersonSkillDependenciesForTenant(w, r, tenantID, "list single person skill dependencies")
}

func (h *Handler) ListTenantProjectSkillRequirements(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant project skill requirements"); ok {
		h.listProjectSkillRequirementsForTenant(w, r, tenantID, "list tenant project skill requirements")
	}
}

func (h *Handler) CreateTenantProjectSkillRequirement(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant project skill requirement"); ok {
		h.createProjectSkillRequirementForTenant(w, r, tenantID, "create tenant project skill requirement")
	}
}

func (h *Handler) UpdateTenantProjectSkillRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminProjectSkillRequirementRequestIDs(w, r, "update tenant project skill requirement")
	if !ok {
		return
	}
	h.updateProjectSkillRequirementForTenant(w, r, tenantID, id, "update tenant project skill requirement")
}

func (h *Handler) DeleteTenantProjectSkillRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminProjectSkillRequirementRequestIDs(w, r, "delete tenant project skill requirement")
	if !ok {
		return
	}
	if err := h.svc.DeleteProjectSkillRequirement(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant project skill requirement", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantProjectSkillGapRows(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant project skill gaps"); ok {
		h.listProjectSkillGapRowsForTenant(w, r, tenantID, "list tenant project skill gaps")
	}
}

func (h *Handler) ListTenantSkillGapSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant skill gap summary"); ok {
		h.listSkillGapSummaryForTenant(w, r, tenantID, "list tenant skill gap summary")
	}
}

func (h *Handler) ListTenantSinglePersonSkillDependencies(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant single person skill dependencies"); ok {
		h.listSinglePersonSkillDependenciesForTenant(w, r, tenantID, "list tenant single person skill dependencies")
	}
}

func (h *Handler) listProjectSkillRequirementsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter, ok := h.projectSkillRequirementFilter(w, r, tenantID, operation)
	if !ok {
		return
	}
	items, err := h.svc.ListProjectSkillRequirements(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list project skill requirements")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createProjectSkillRequirementForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ProjectSkillRequirementCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateProjectSkillRequirement(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateProjectSkillRequirementForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ProjectSkillRequirementCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateProjectSkillRequirement(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listProjectSkillGapRowsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter, ok := h.projectSkillRequirementFilter(w, r, tenantID, operation)
	if !ok {
		return
	}
	items, err := h.svc.ListProjectSkillGapRows(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list project skill gaps")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listSkillGapSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	projectID, ok := h.optionalUUIDQuery(w, r, "project_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListSkillGapSummary(r.Context(), tenantID, projectID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list skill gap summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listSinglePersonSkillDependenciesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	projectID, ok := h.optionalUUIDQuery(w, r, "project_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListSinglePersonSkillDependencies(r.Context(), tenantID, projectID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list single person skill dependencies")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) projectSkillRequirementFilter(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (domain.ProjectSkillRequirementFilter, bool) {
	projectID, ok := h.optionalUUIDQuery(w, r, "project_id", operation)
	if !ok {
		return domain.ProjectSkillRequirementFilter{}, false
	}
	engagementID, ok := h.optionalUUIDQuery(w, r, "engagement_id", operation)
	if !ok {
		return domain.ProjectSkillRequirementFilter{}, false
	}
	skillID, ok := h.optionalUUIDQuery(w, r, "skill_id", operation)
	if !ok {
		return domain.ProjectSkillRequirementFilter{}, false
	}
	return domain.ProjectSkillRequirementFilter{TenantID: tenantID, ProjectID: projectID, EngagementID: engagementID, SkillID: skillID, Importance: optionalStringQuery(r, "importance"), Search: optionalStringQuery(r, "search")}, true
}

func (h *Handler) projectSkillRequirementRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "requirementID", operation, "invalid requirement id")
	return tenantID, id, ok
}

func (h *Handler) superAdminProjectSkillRequirementRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "requirementID", operation, "invalid requirement id")
	return tenantID, id, ok
}
