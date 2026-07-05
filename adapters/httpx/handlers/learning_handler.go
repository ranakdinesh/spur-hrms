package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListLearningCourses(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listLearningCoursesForTenant(w, r, tenantID, "list learning courses")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list learning courses", err, "tenant context is required")
	}
}

func (h *Handler) CreateLearningCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createLearningCourseForTenant(w, r, tenantID, "create learning course")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create learning course", err, "tenant context is required")
	}
}

func (h *Handler) UpdateLearningCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.learningCourseRequestIDs(w, r, "update learning course"); ok {
		h.updateLearningCourseForTenant(w, r, tenantID, id, "update learning course")
	}
}

func (h *Handler) DeleteLearningCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.learningCourseRequestIDs(w, r, "delete learning course"); ok {
		if err := h.svc.DeleteLearningCourse(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete learning course", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListLearningPaths(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listLearningPathsForTenant(w, r, tenantID, "list learning paths")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list learning paths", err, "tenant context is required")
	}
}

func (h *Handler) CreateLearningPath(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createLearningPathForTenant(w, r, tenantID, "create learning path")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create learning path", err, "tenant context is required")
	}
}

func (h *Handler) UpdateLearningPath(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.learningPathRequestIDs(w, r, "update learning path"); ok {
		h.updateLearningPathForTenant(w, r, tenantID, id, "update learning path")
	}
}

func (h *Handler) DeleteLearningPath(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.learningPathRequestIDs(w, r, "delete learning path"); ok {
		if err := h.svc.DeleteLearningPath(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete learning path", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListLearningPathCourses(w http.ResponseWriter, r *http.Request) {
	if tenantID, pathID, ok := h.learningPathRequestIDs(w, r, "list learning path courses"); ok {
		items, err := h.svc.ListLearningPathCourses(r.Context(), tenantID, pathID)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, "list learning path courses", err, "failed to list learning path courses")
			return
		}
		respondJSON(w, http.StatusOK, items)
	}
}

func (h *Handler) UpsertLearningPathCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, pathID, ok := h.learningPathRequestIDs(w, r, "upsert learning path course"); ok {
		var cmd ports.LearningPathCourseCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode learning path course request", err, "invalid request body")
			return
		}
		cmd.TenantID = tenantID
		cmd.PathID = pathID
		cmd.ActorID = h.actorIDFromRequest(r)
		item, err := h.svc.UpsertLearningPathCourse(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "upsert learning path course", err, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, item)
	}
}

func (h *Handler) DeleteLearningPathCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.learningPathCourseRequestIDs(w, r, "delete learning path course"); ok {
		if err := h.svc.DeleteLearningPathCourse(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete learning path course", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListLearningEnrollments(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listLearningEnrollmentsForTenant(w, r, tenantID, "list learning enrollments")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list learning enrollments", err, "tenant context is required")
	}
}

func (h *Handler) CreateLearningEnrollment(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createLearningEnrollmentForTenant(w, r, tenantID, "create learning enrollment")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create learning enrollment", err, "tenant context is required")
	}
}

func (h *Handler) UpdateLearningEnrollmentStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.learningEnrollmentRequestIDs(w, r, "update learning enrollment status"); ok {
		h.updateLearningEnrollmentStatusForTenant(w, r, tenantID, id, "update learning enrollment status")
	}
}

func (h *Handler) UploadLearningCertificate(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.learningEnrollmentRequestIDs(w, r, "upload learning certificate"); ok {
		var cmd ports.LearningCertificateCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode learning certificate request", err, "invalid request body")
			return
		}
		cmd.TenantID = tenantID
		cmd.EnrollmentID = id
		cmd.ActorID = h.actorIDFromRequest(r)
		item, err := h.svc.UploadLearningCertificate(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "upload learning certificate", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
	}
}

func (h *Handler) DeleteLearningEnrollment(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.learningEnrollmentRequestIDs(w, r, "delete learning enrollment"); ok {
		if err := h.svc.DeleteLearningEnrollment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete learning enrollment", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListLearningRecommendations(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listLearningRecommendationsForTenant(w, r, tenantID, "list learning recommendations")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list learning recommendations", err, "tenant context is required")
	}
}

func (h *Handler) CreateLearningRecommendation(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createLearningRecommendationForTenant(w, r, tenantID, "create learning recommendation")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create learning recommendation", err, "tenant context is required")
	}
}

func (h *Handler) UpdateLearningRecommendationStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.learningRecommendationRequestIDs(w, r, "update learning recommendation status"); ok {
		h.updateLearningRecommendationStatusForTenant(w, r, tenantID, id, "update learning recommendation status")
	}
}

func (h *Handler) GenerateLearningRecommendations(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		items, err := h.svc.GenerateSkillGapLearningRecommendations(r.Context(), tenantID, h.actorIDFromRequest(r))
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "generate learning recommendations", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, items)
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "generate learning recommendations", err, "tenant context is required")
	}
}

func (h *Handler) GetLearningSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.getLearningSummaryForTenant(w, r, tenantID, "get learning summary")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "get learning summary", err, "tenant context is required")
	}
}

func (h *Handler) ListTenantLearningCourses(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant learning courses"); ok {
		h.listLearningCoursesForTenant(w, r, tenantID, "list tenant learning courses")
	}
}

func (h *Handler) CreateTenantLearningCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant learning course"); ok {
		h.createLearningCourseForTenant(w, r, tenantID, "create tenant learning course")
	}
}

func (h *Handler) UpdateTenantLearningCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminLearningRequestID(w, r, "courseID", "update tenant learning course", "invalid course id"); ok {
		h.updateLearningCourseForTenant(w, r, tenantID, id, "update tenant learning course")
	}
}

func (h *Handler) DeleteTenantLearningCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminLearningRequestID(w, r, "courseID", "delete tenant learning course", "invalid course id"); ok {
		if err := h.svc.DeleteLearningCourse(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete tenant learning course", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListTenantLearningPaths(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant learning paths"); ok {
		h.listLearningPathsForTenant(w, r, tenantID, "list tenant learning paths")
	}
}

func (h *Handler) CreateTenantLearningPath(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant learning path"); ok {
		h.createLearningPathForTenant(w, r, tenantID, "create tenant learning path")
	}
}

func (h *Handler) UpdateTenantLearningPath(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminLearningRequestID(w, r, "pathID", "update tenant learning path", "invalid path id"); ok {
		h.updateLearningPathForTenant(w, r, tenantID, id, "update tenant learning path")
	}
}

func (h *Handler) DeleteTenantLearningPath(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminLearningRequestID(w, r, "pathID", "delete tenant learning path", "invalid path id"); ok {
		if err := h.svc.DeleteLearningPath(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete tenant learning path", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListTenantLearningPathCourses(w http.ResponseWriter, r *http.Request) {
	if tenantID, pathID, ok := h.superAdminLearningRequestID(w, r, "pathID", "list tenant learning path courses", "invalid path id"); ok {
		items, err := h.svc.ListLearningPathCourses(r.Context(), tenantID, pathID)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, "list tenant learning path courses", err, "failed to list learning path courses")
			return
		}
		respondJSON(w, http.StatusOK, items)
	}
}

func (h *Handler) UpsertTenantLearningPathCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, pathID, ok := h.superAdminLearningRequestID(w, r, "pathID", "upsert tenant learning path course", "invalid path id"); ok {
		var cmd ports.LearningPathCourseCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode tenant learning path course request", err, "invalid request body")
			return
		}
		cmd.TenantID = tenantID
		cmd.PathID = pathID
		cmd.ActorID = h.actorIDFromRequest(r)
		item, err := h.svc.UpsertLearningPathCourse(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "upsert tenant learning path course", err, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, item)
	}
}

func (h *Handler) DeleteTenantLearningPathCourse(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminLearningRequestID(w, r, "pathCourseID", "delete tenant learning path course", "invalid path course id"); ok {
		if err := h.svc.DeleteLearningPathCourse(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete tenant learning path course", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListTenantLearningEnrollments(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant learning enrollments"); ok {
		h.listLearningEnrollmentsForTenant(w, r, tenantID, "list tenant learning enrollments")
	}
}

func (h *Handler) CreateTenantLearningEnrollment(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant learning enrollment"); ok {
		h.createLearningEnrollmentForTenant(w, r, tenantID, "create tenant learning enrollment")
	}
}

func (h *Handler) UpdateTenantLearningEnrollmentStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminLearningRequestID(w, r, "enrollmentID", "update tenant learning enrollment status", "invalid enrollment id"); ok {
		h.updateLearningEnrollmentStatusForTenant(w, r, tenantID, id, "update tenant learning enrollment status")
	}
}

func (h *Handler) UploadTenantLearningCertificate(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminLearningRequestID(w, r, "enrollmentID", "upload tenant learning certificate", "invalid enrollment id"); ok {
		var cmd ports.LearningCertificateCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode tenant learning certificate request", err, "invalid request body")
			return
		}
		cmd.TenantID = tenantID
		cmd.EnrollmentID = id
		cmd.ActorID = h.actorIDFromRequest(r)
		item, err := h.svc.UploadLearningCertificate(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "upload tenant learning certificate", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
	}
}

func (h *Handler) DeleteTenantLearningEnrollment(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminLearningRequestID(w, r, "enrollmentID", "delete tenant learning enrollment", "invalid enrollment id"); ok {
		if err := h.svc.DeleteLearningEnrollment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete tenant learning enrollment", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListTenantLearningRecommendations(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant learning recommendations"); ok {
		h.listLearningRecommendationsForTenant(w, r, tenantID, "list tenant learning recommendations")
	}
}

func (h *Handler) CreateTenantLearningRecommendation(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant learning recommendation"); ok {
		h.createLearningRecommendationForTenant(w, r, tenantID, "create tenant learning recommendation")
	}
}

func (h *Handler) UpdateTenantLearningRecommendationStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminLearningRequestID(w, r, "recommendationID", "update tenant learning recommendation status", "invalid recommendation id"); ok {
		h.updateLearningRecommendationStatusForTenant(w, r, tenantID, id, "update tenant learning recommendation status")
	}
}

func (h *Handler) GenerateTenantLearningRecommendations(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "generate tenant learning recommendations"); ok {
		items, err := h.svc.GenerateSkillGapLearningRecommendations(r.Context(), tenantID, h.actorIDFromRequest(r))
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "generate tenant learning recommendations", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, items)
	}
}

func (h *Handler) GetTenantLearningSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant learning summary"); ok {
		h.getLearningSummaryForTenant(w, r, tenantID, "get tenant learning summary")
	}
}

func (h *Handler) listLearningCoursesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	skillID, ok := h.optionalUUIDQuery(w, r, "skill_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListLearningCourses(r.Context(), domain.LearningCourseFilter{TenantID: tenantID, CourseType: optionalStringQuery(r, "course_type"), SkillID: skillID, Mandatory: optionalBoolQuery(r, "mandatory"), AIReadiness: optionalBoolQuery(r, "ai_readiness"), IsActive: optionalBoolQuery(r, "is_active"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list learning courses")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createLearningCourseForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LearningCourseCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLearningCourse(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateLearningCourseForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.LearningCourseCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLearningCourse(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listLearningPathsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	skillID, ok := h.optionalUUIDQuery(w, r, "skill_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListLearningPaths(r.Context(), domain.LearningPathFilter{TenantID: tenantID, PathType: optionalStringQuery(r, "path_type"), SkillID: skillID, AIReadiness: optionalBoolQuery(r, "ai_readiness"), IsActive: optionalBoolQuery(r, "is_active"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list learning paths")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createLearningPathForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LearningPathCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLearningPath(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateLearningPathForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.LearningPathCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLearningPath(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listLearningEnrollmentsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	courseID, ok := h.optionalUUIDQuery(w, r, "course_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListLearningEnrollments(r.Context(), domain.LearningEnrollmentFilter{TenantID: tenantID, WorkerProfileID: workerID, CourseID: courseID, Status: optionalStringQuery(r, "status"), AssignmentSource: optionalStringQuery(r, "assignment_source"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list learning enrollments")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createLearningEnrollmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LearningEnrollmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLearningEnrollment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateLearningEnrollmentStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.LearningEnrollmentStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLearningEnrollmentStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listLearningRecommendationsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	workerID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	skillID, ok := h.optionalUUIDQuery(w, r, "skill_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListLearningRecommendations(r.Context(), domain.LearningRecommendationFilter{TenantID: tenantID, WorkerProfileID: workerID, SkillID: skillID, SourceType: optionalStringQuery(r, "source_type"), Status: optionalStringQuery(r, "status"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list learning recommendations")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createLearningRecommendationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.LearningRecommendationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateLearningRecommendation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateLearningRecommendationStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.LearningRecommendationStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateLearningRecommendationStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) getLearningSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.GetLearningSummary(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to get learning summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) learningCourseRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "courseID", operation, "invalid course id")
	return tenantID, id, ok
}

func (h *Handler) learningPathRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "pathID", operation, "invalid path id")
	return tenantID, id, ok
}

func (h *Handler) learningPathCourseRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "pathCourseID", operation, "invalid path course id")
	return tenantID, id, ok
}

func (h *Handler) learningEnrollmentRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "enrollmentID", operation, "invalid enrollment id")
	return tenantID, id, ok
}

func (h *Handler) learningRecommendationRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "recommendationID", operation, "invalid recommendation id")
	return tenantID, id, ok
}

func (h *Handler) superAdminLearningRequestID(w http.ResponseWriter, r *http.Request, param string, operation string, message string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, param, operation, message)
	return tenantID, id, ok
}
