package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateShiftTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create shift template", err, "tenant context is required")
		return
	}
	h.createShiftTemplateForTenant(w, r, tenantID, "create shift template")
}

func (h *Handler) CreateTenantShiftTemplate(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant shift template"); ok {
		h.createShiftTemplateForTenant(w, r, tenantID, "create tenant shift template")
	}
}

func (h *Handler) ListShiftTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list shift templates", err, "tenant context is required")
		return
	}
	h.listShiftTemplatesForTenant(w, r, tenantID, "list shift templates")
}

func (h *Handler) ListTenantShiftTemplates(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant shift templates"); ok {
		h.listShiftTemplatesForTenant(w, r, tenantID, "list tenant shift templates")
	}
}

func (h *Handler) UpdateShiftTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "templateID", "update shift template")
	if !ok {
		return
	}
	h.updateShiftTemplateForTenant(w, r, tenantID, id, "update shift template")
}

func (h *Handler) UpdateTenantShiftTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "templateID", "update tenant shift template")
	if !ok {
		return
	}
	h.updateShiftTemplateForTenant(w, r, tenantID, id, "update tenant shift template")
}

func (h *Handler) DeleteShiftTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "templateID", "delete shift template")
	if !ok {
		return
	}
	if err := h.svc.DeleteShiftTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete shift template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteTenantShiftTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "templateID", "delete tenant shift template")
	if !ok {
		return
	}
	if err := h.svc.DeleteShiftTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant shift template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateStaffingRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create staffing requirement", err, "tenant context is required")
		return
	}
	h.createStaffingRequirementForTenant(w, r, tenantID, "create staffing requirement")
}

func (h *Handler) CreateTenantStaffingRequirement(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant staffing requirement"); ok {
		h.createStaffingRequirementForTenant(w, r, tenantID, "create tenant staffing requirement")
	}
}

func (h *Handler) ListStaffingRequirements(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list staffing requirements", err, "tenant context is required")
		return
	}
	h.listStaffingRequirementsForTenant(w, r, tenantID, "list staffing requirements")
}

func (h *Handler) ListTenantStaffingRequirements(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant staffing requirements"); ok {
		h.listStaffingRequirementsForTenant(w, r, tenantID, "list tenant staffing requirements")
	}
}

func (h *Handler) UpdateStaffingRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "requirementID", "update staffing requirement")
	if !ok {
		return
	}
	h.updateStaffingRequirementForTenant(w, r, tenantID, id, "update staffing requirement")
}

func (h *Handler) UpdateTenantStaffingRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "requirementID", "update tenant staffing requirement")
	if !ok {
		return
	}
	h.updateStaffingRequirementForTenant(w, r, tenantID, id, "update tenant staffing requirement")
}

func (h *Handler) DeleteStaffingRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "requirementID", "delete staffing requirement")
	if !ok {
		return
	}
	if err := h.svc.DeleteStaffingRequirement(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete staffing requirement", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteTenantStaffingRequirement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "requirementID", "delete tenant staffing requirement")
	if !ok {
		return
	}
	if err := h.svc.DeleteStaffingRequirement(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant staffing requirement", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateShiftScheduleAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create shift assignment", err, "tenant context is required")
		return
	}
	h.createShiftScheduleAssignmentForTenant(w, r, tenantID, "create shift assignment")
}

func (h *Handler) CreateTenantShiftScheduleAssignment(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant shift assignment"); ok {
		h.createShiftScheduleAssignmentForTenant(w, r, tenantID, "create tenant shift assignment")
	}
}

func (h *Handler) ListShiftScheduleAssignments(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list shift assignments", err, "tenant context is required")
		return
	}
	h.listShiftScheduleAssignmentsForTenant(w, r, tenantID, "list shift assignments")
}

func (h *Handler) ListTenantShiftScheduleAssignments(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant shift assignments"); ok {
		h.listShiftScheduleAssignmentsForTenant(w, r, tenantID, "list tenant shift assignments")
	}
}

func (h *Handler) UpdateShiftScheduleAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "assignmentID", "update shift assignment")
	if !ok {
		return
	}
	h.updateShiftScheduleAssignmentForTenant(w, r, tenantID, id, "update shift assignment")
}

func (h *Handler) UpdateTenantShiftScheduleAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "assignmentID", "update tenant shift assignment")
	if !ok {
		return
	}
	h.updateShiftScheduleAssignmentForTenant(w, r, tenantID, id, "update tenant shift assignment")
}

func (h *Handler) UpdateShiftScheduleAssignmentStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "assignmentID", "update shift assignment status")
	if !ok {
		return
	}
	h.updateShiftScheduleAssignmentStatusForTenant(w, r, tenantID, id, "update shift assignment status")
}

func (h *Handler) UpdateTenantShiftScheduleAssignmentStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "assignmentID", "update tenant shift assignment status")
	if !ok {
		return
	}
	h.updateShiftScheduleAssignmentStatusForTenant(w, r, tenantID, id, "update tenant shift assignment status")
}

func (h *Handler) DeleteShiftScheduleAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "assignmentID", "delete shift assignment")
	if !ok {
		return
	}
	if err := h.svc.DeleteShiftScheduleAssignment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete shift assignment", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteTenantShiftScheduleAssignment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "assignmentID", "delete tenant shift assignment")
	if !ok {
		return
	}
	if err := h.svc.DeleteShiftScheduleAssignment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant shift assignment", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateShiftSwapRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create shift swap request", err, "tenant context is required")
		return
	}
	h.createShiftSwapRequestForTenant(w, r, tenantID, "create shift swap request")
}

func (h *Handler) CreateTenantShiftSwapRequest(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant shift swap request"); ok {
		h.createShiftSwapRequestForTenant(w, r, tenantID, "create tenant shift swap request")
	}
}

func (h *Handler) ListShiftSwapRequests(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list shift swap requests", err, "tenant context is required")
		return
	}
	h.listShiftSwapRequestsForTenant(w, r, tenantID, "list shift swap requests")
}

func (h *Handler) ListTenantShiftSwapRequests(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant shift swap requests"); ok {
		h.listShiftSwapRequestsForTenant(w, r, tenantID, "list tenant shift swap requests")
	}
}

func (h *Handler) ReviewShiftSwapRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "swapID", "review shift swap request")
	if !ok {
		return
	}
	h.reviewShiftSwapRequestForTenant(w, r, tenantID, id, "review shift swap request")
}

func (h *Handler) ReviewTenantShiftSwapRequest(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "swapID", "review tenant shift swap request")
	if !ok {
		return
	}
	h.reviewShiftSwapRequestForTenant(w, r, tenantID, id, "review tenant shift swap request")
}

func (h *Handler) ListShiftScheduleEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list shift schedule events", err, "tenant context is required")
		return
	}
	h.listShiftScheduleEventsForTenant(w, r, tenantID, "list shift schedule events")
}

func (h *Handler) ListTenantShiftScheduleEvents(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant shift schedule events"); ok {
		h.listShiftScheduleEventsForTenant(w, r, tenantID, "list tenant shift schedule events")
	}
}

func (h *Handler) GetShiftScheduleSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get shift schedule summary", err, "tenant context is required")
		return
	}
	h.getShiftScheduleSummaryForTenant(w, r, tenantID, "get shift schedule summary")
}

func (h *Handler) GetTenantShiftScheduleSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant shift schedule summary"); ok {
		h.getShiftScheduleSummaryForTenant(w, r, tenantID, "get tenant shift schedule summary")
	}
}

func (h *Handler) ListShiftStaffingGaps(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list shift staffing gaps", err, "tenant context is required")
		return
	}
	h.listShiftStaffingGapsForTenant(w, r, tenantID, "list shift staffing gaps")
}

func (h *Handler) ListTenantShiftStaffingGaps(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant shift staffing gaps"); ok {
		h.listShiftStaffingGapsForTenant(w, r, tenantID, "list tenant shift staffing gaps")
	}
}

func (h *Handler) createShiftTemplateForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ShiftTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateShiftTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateShiftTemplateForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ShiftTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateShiftTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listShiftTemplatesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListShiftTemplates(r.Context(), tenantID, optionalBoolQuery(r, "active_only"), optionalStringQuery(r, "search"), queryInt32(r, "limit", 100), queryInt32(r, "offset", 0))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createStaffingRequirementForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.StaffingRequirementCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateStaffingRequirement(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateStaffingRequirementForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.StaffingRequirementCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateStaffingRequirement(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listStaffingRequirementsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter := domain.StaffingRequirementFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), StartDate: r.URL.Query().Get("start_date"), EndDate: r.URL.Query().Get("end_date"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)}
	items, err := h.svc.ListStaffingRequirements(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createShiftScheduleAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ShiftScheduleAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateShiftScheduleAssignment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateShiftScheduleAssignmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ShiftScheduleAssignmentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateShiftScheduleAssignment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateShiftScheduleAssignmentStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ShiftScheduleStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateShiftScheduleAssignmentStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listShiftScheduleAssignmentsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter := domain.ShiftScheduleFilter{TenantID: tenantID, StartDate: r.URL.Query().Get("start_date"), EndDate: r.URL.Query().Get("end_date"), Status: optionalStringQuery(r, "status"), WorkerProfileID: optionalUUIDQuery(r, "worker_profile_id"), EmployeeUserID: optionalUUIDQuery(r, "employee_user_id"), BranchID: optionalUUIDQuery(r, "branch_id"), DepartmentID: optionalUUIDQuery(r, "department_id"), AttendanceLocationID: optionalUUIDQuery(r, "attendance_location_id"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)}
	items, err := h.svc.ListShiftScheduleAssignments(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createShiftSwapRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ShiftSwapRequestCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateShiftSwapRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listShiftSwapRequestsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter := domain.ShiftSwapFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), RequesterUserID: optionalUUIDQuery(r, "requester_user_id"), TargetUserID: optionalUUIDQuery(r, "target_user_id"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)}
	items, err := h.svc.ListShiftSwapRequests(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) reviewShiftSwapRequestForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ShiftSwapReviewCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	actorID := h.actorIDFromRequest(r)
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, actorID
	if actorID != nil {
		cmd.ReviewerID = *actorID
	}
	item, err := h.svc.ReviewShiftSwapRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listShiftScheduleEventsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter := domain.ShiftScheduleEventFilter{TenantID: tenantID, SourceType: optionalStringQuery(r, "source_type"), SourceID: optionalUUIDQuery(r, "source_id"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)}
	items, err := h.svc.ListShiftScheduleEvents(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getShiftScheduleSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.GetShiftScheduleSummary(r.Context(), tenantID, r.URL.Query().Get("start_date"), r.URL.Query().Get("end_date"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listShiftStaffingGapsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListShiftStaffingGaps(r.Context(), tenantID, r.URL.Query().Get("start_date"), r.URL.Query().Get("end_date"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) superAdminTenantAndURLUUID(w http.ResponseWriter, r *http.Request, param string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid identifier")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
