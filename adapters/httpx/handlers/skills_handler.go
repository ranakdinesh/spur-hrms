package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListSkillCategories(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list skill categories", err, "tenant context is required")
		return
	}
	h.listSkillCategoriesForTenant(w, r, tenantID, "list skill categories")
}

func (h *Handler) CreateSkillCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create skill category", err, "tenant context is required")
		return
	}
	h.createSkillCategoryForTenant(w, r, tenantID, "create skill category")
}

func (h *Handler) UpdateSkillCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.skillCategoryRequestIDs(w, r, "update skill category")
	if !ok {
		return
	}
	h.updateSkillCategoryForTenant(w, r, tenantID, id, "update skill category")
}

func (h *Handler) DeleteSkillCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.skillCategoryRequestIDs(w, r, "delete skill category")
	if !ok {
		return
	}
	if err := h.svc.DeleteSkillCategory(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete skill category", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListSkills(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list skills", err, "tenant context is required")
		return
	}
	h.listSkillsForTenant(w, r, tenantID, "list skills")
}

func (h *Handler) CreateSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create skill", err, "tenant context is required")
		return
	}
	h.createSkillForTenant(w, r, tenantID, "create skill")
}

func (h *Handler) UpdateSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.skillRequestIDs(w, r, "update skill")
	if !ok {
		return
	}
	h.updateSkillForTenant(w, r, tenantID, id, "update skill")
}

func (h *Handler) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.skillRequestIDs(w, r, "delete skill")
	if !ok {
		return
	}
	if err := h.svc.DeleteSkill(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete skill", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListWorkerSkills(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list worker skills", err, "tenant context is required")
		return
	}
	h.listWorkerSkillsForTenant(w, r, tenantID, "list worker skills")
}

func (h *Handler) CreateWorkerSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create worker skill", err, "tenant context is required")
		return
	}
	h.createWorkerSkillForTenant(w, r, tenantID, "create worker skill")
}

func (h *Handler) UpdateWorkerSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.workerSkillRequestIDs(w, r, "update worker skill")
	if !ok {
		return
	}
	h.updateWorkerSkillForTenant(w, r, tenantID, id, "update worker skill")
}

func (h *Handler) VerifyWorkerSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.workerSkillRequestIDs(w, r, "verify worker skill")
	if !ok {
		return
	}
	h.verifyWorkerSkillForTenant(w, r, tenantID, id, "verify worker skill")
}

func (h *Handler) DeleteWorkerSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.workerSkillRequestIDs(w, r, "delete worker skill")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkerSkill(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete worker skill", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListWorkerSkillAssessments(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list worker skill assessments", err, "tenant context is required")
		return
	}
	h.listWorkerSkillAssessmentsForTenant(w, r, tenantID, "list worker skill assessments")
}

func (h *Handler) CreateWorkerSkillAssessment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create worker skill assessment", err, "tenant context is required")
		return
	}
	h.createWorkerSkillAssessmentForTenant(w, r, tenantID, "create worker skill assessment")
}

func (h *Handler) GetSkillsSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get skills summary", err, "tenant context is required")
		return
	}
	h.getSkillsSummaryForTenant(w, r, tenantID, "get skills summary")
}

func (h *Handler) ListTenantSkillCategories(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant skill categories"); ok {
		h.listSkillCategoriesForTenant(w, r, tenantID, "list tenant skill categories")
	}
}

func (h *Handler) CreateTenantSkillCategory(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant skill category"); ok {
		h.createSkillCategoryForTenant(w, r, tenantID, "create tenant skill category")
	}
}

func (h *Handler) UpdateTenantSkillCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminSkillCategoryRequestIDs(w, r, "update tenant skill category")
	if ok {
		h.updateSkillCategoryForTenant(w, r, tenantID, id, "update tenant skill category")
	}
}

func (h *Handler) DeleteTenantSkillCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminSkillCategoryRequestIDs(w, r, "delete tenant skill category")
	if !ok {
		return
	}
	if err := h.svc.DeleteSkillCategory(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant skill category", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantSkills(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant skills"); ok {
		h.listSkillsForTenant(w, r, tenantID, "list tenant skills")
	}
}

func (h *Handler) CreateTenantSkill(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant skill"); ok {
		h.createSkillForTenant(w, r, tenantID, "create tenant skill")
	}
}

func (h *Handler) UpdateTenantSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminSkillRequestIDs(w, r, "update tenant skill")
	if ok {
		h.updateSkillForTenant(w, r, tenantID, id, "update tenant skill")
	}
}

func (h *Handler) DeleteTenantSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminSkillRequestIDs(w, r, "delete tenant skill")
	if !ok {
		return
	}
	if err := h.svc.DeleteSkill(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant skill", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantWorkerSkills(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant worker skills"); ok {
		h.listWorkerSkillsForTenant(w, r, tenantID, "list tenant worker skills")
	}
}

func (h *Handler) CreateTenantWorkerSkill(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant worker skill"); ok {
		h.createWorkerSkillForTenant(w, r, tenantID, "create tenant worker skill")
	}
}

func (h *Handler) UpdateTenantWorkerSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminWorkerSkillRequestIDs(w, r, "update tenant worker skill")
	if ok {
		h.updateWorkerSkillForTenant(w, r, tenantID, id, "update tenant worker skill")
	}
}

func (h *Handler) VerifyTenantWorkerSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminWorkerSkillRequestIDs(w, r, "verify tenant worker skill")
	if ok {
		h.verifyWorkerSkillForTenant(w, r, tenantID, id, "verify tenant worker skill")
	}
}

func (h *Handler) DeleteTenantWorkerSkill(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminWorkerSkillRequestIDs(w, r, "delete tenant worker skill")
	if !ok {
		return
	}
	if err := h.svc.DeleteWorkerSkill(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant worker skill", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantWorkerSkillAssessments(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant worker skill assessments"); ok {
		h.listWorkerSkillAssessmentsForTenant(w, r, tenantID, "list tenant worker skill assessments")
	}
}

func (h *Handler) CreateTenantWorkerSkillAssessment(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant worker skill assessment"); ok {
		h.createWorkerSkillAssessmentForTenant(w, r, tenantID, "create tenant worker skill assessment")
	}
}

func (h *Handler) GetTenantSkillsSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant skills summary"); ok {
		h.getSkillsSummaryForTenant(w, r, tenantID, "get tenant skills summary")
	}
}

func (h *Handler) listSkillCategoriesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	parentID, ok := h.optionalUUIDQuery(w, r, "parent_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListSkillCategories(r.Context(), domain.SkillCategoryFilter{TenantID: tenantID, SourceScope: optionalStringQuery(r, "source_scope"), ParentID: parentID, Search: optionalStringQuery(r, "search")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list skill categories")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createSkillCategoryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.SkillCategoryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateSkillCategory(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateSkillCategoryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.SkillCategoryCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateSkillCategory(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listSkillsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	categoryID, ok := h.optionalUUIDQuery(w, r, "category_id", operation)
	if !ok {
		return
	}
	active := optionalBoolQuery(r, "is_active")
	items, err := h.svc.ListSkills(r.Context(), domain.SkillFilter{TenantID: tenantID, CategoryID: categoryID, SkillType: optionalStringQuery(r, "skill_type"), SourceScope: optionalStringQuery(r, "source_scope"), IsActive: active, Search: optionalStringQuery(r, "search")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list skills")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createSkillForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.SkillCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateSkill(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateSkillForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.SkillCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateSkill(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listWorkerSkillsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerProfileID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	skillID, ok := h.optionalUUIDQuery(w, r, "skill_id", operation)
	if !ok {
		return
	}
	categoryID, ok := h.optionalUUIDQuery(w, r, "category_id", operation)
	if !ok {
		return
	}
	expiringBefore, ok := optionalDateQuery(w, r, "certificate_expiring_before", operation, h)
	if !ok {
		return
	}
	items, err := h.svc.ListWorkerSkills(r.Context(), domain.WorkerSkillFilter{TenantID: tenantID, WorkerProfileID: workerProfileID, SkillID: skillID, CategoryID: categoryID, Proficiency: optionalStringQuery(r, "proficiency"), VerificationStatus: optionalStringQuery(r, "verification_status"), CertificateExpiringBefore: expiringBefore, Search: optionalStringQuery(r, "search")})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list worker skills")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createWorkerSkillForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.WorkerSkillCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkerSkill(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateWorkerSkillForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.WorkerSkillCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkerSkill(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) verifyWorkerSkillForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.WorkerSkillVerificationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateWorkerSkillVerification(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listWorkerSkillAssessmentsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerSkillID, ok := h.optionalUUIDQuery(w, r, "worker_skill_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListWorkerSkillAssessments(r.Context(), tenantID, workerSkillID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list worker skill assessments")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createWorkerSkillAssessmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.WorkerSkillAssessmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateWorkerSkillAssessment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getSkillsSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.GetSkillsSummary(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to get skills summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) skillCategoryRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "skillCategoryID", operation, "invalid skill category id")
	return tenantID, id, ok
}

func (h *Handler) skillRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "skillID", operation, "invalid skill id")
	return tenantID, id, ok
}

func (h *Handler) workerSkillRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "workerSkillID", operation, "invalid worker skill id")
	return tenantID, id, ok
}

func (h *Handler) superAdminSkillCategoryRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "skillCategoryID", operation, "invalid skill category id")
	return tenantID, id, ok
}

func (h *Handler) superAdminSkillRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "skillID", operation, "invalid skill id")
	return tenantID, id, ok
}

func (h *Handler) superAdminWorkerSkillRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "workerSkillID", operation, "invalid worker skill id")
	return tenantID, id, ok
}

func (h *Handler) uuidURLParam(w http.ResponseWriter, r *http.Request, key string, operation string, message string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, key))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, message)
		return uuid.Nil, false
	}
	return id, true
}
