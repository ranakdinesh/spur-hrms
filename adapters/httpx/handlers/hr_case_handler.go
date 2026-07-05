package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListHRCaseCategories(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list hr case categories", err, "tenant context is required")
		return
	}
	h.listHRCaseCategoriesForTenant(w, r, tenantID, "list hr case categories")
}

func (h *Handler) CreateHRCaseCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create hr case category", err, "tenant context is required")
		return
	}
	h.createHRCaseCategoryForTenant(w, r, tenantID, "create hr case category")
}

func (h *Handler) UpdateHRCaseCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "categoryID", "update hr case category")
	if !ok {
		return
	}
	h.updateHRCaseCategoryForTenant(w, r, tenantID, id, "update hr case category")
}

func (h *Handler) DeleteHRCaseCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "categoryID", "delete hr case category")
	if !ok {
		return
	}
	if err := h.svc.DeleteHRCaseCategory(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondHRCaseError(w, r, "delete hr case category", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListHRCaseSLAPolicies(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list hr case sla policies", err, "tenant context is required")
		return
	}
	h.listHRCaseSLAPoliciesForTenant(w, r, tenantID, "list hr case sla policies")
}

func (h *Handler) CreateHRCaseSLAPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create hr case sla policy", err, "tenant context is required")
		return
	}
	h.createHRCaseSLAPolicyForTenant(w, r, tenantID, "create hr case sla policy")
}

func (h *Handler) UpdateHRCaseSLAPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "slaPolicyID", "update hr case sla policy")
	if !ok {
		return
	}
	h.updateHRCaseSLAPolicyForTenant(w, r, tenantID, id, "update hr case sla policy")
}

func (h *Handler) DeleteHRCaseSLAPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "slaPolicyID", "delete hr case sla policy")
	if !ok {
		return
	}
	if err := h.svc.DeleteHRCaseSLAPolicy(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondHRCaseError(w, r, "delete hr case sla policy", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListHRCases(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list hr cases", err, "tenant context is required")
		return
	}
	h.listHRCasesForTenant(w, r, tenantID, "list hr cases")
}

func (h *Handler) CreateHRCase(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create hr case", err, "tenant context is required")
		return
	}
	h.createHRCaseForTenant(w, r, tenantID, "create hr case")
}

func (h *Handler) GetHRCase(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "caseID", "get hr case")
	if !ok {
		return
	}
	h.getHRCaseForTenant(w, r, tenantID, id, "get hr case")
}

func (h *Handler) UpdateHRCase(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "caseID", "update hr case")
	if !ok {
		return
	}
	h.updateHRCaseForTenant(w, r, tenantID, id, "update hr case")
}

func (h *Handler) AssignHRCase(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "caseID", "assign hr case")
	if !ok {
		return
	}
	h.assignHRCaseForTenant(w, r, tenantID, id, "assign hr case")
}

func (h *Handler) UpdateHRCaseStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "caseID", "update hr case status")
	if !ok {
		return
	}
	h.updateHRCaseStatusForTenant(w, r, tenantID, id, "update hr case status")
}

func (h *Handler) CreateHRCaseComment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "caseID", "create hr case comment")
	if !ok {
		return
	}
	h.createHRCaseCommentForTenant(w, r, tenantID, id, "create hr case comment")
}

func (h *Handler) CreateHRCaseAttachment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.hrCaseIDPair(w, r, "caseID", "create hr case attachment")
	if !ok {
		return
	}
	h.createHRCaseAttachmentForTenant(w, r, tenantID, id, "create hr case attachment")
}

func (h *Handler) ListTenantHRCaseCategories(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant hr case categories"); ok {
		h.listHRCaseCategoriesForTenant(w, r, tenantID, "list tenant hr case categories")
	}
}

func (h *Handler) CreateTenantHRCaseCategory(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant hr case category"); ok {
		h.createHRCaseCategoryForTenant(w, r, tenantID, "create tenant hr case category")
	}
}

func (h *Handler) UpdateTenantHRCaseCategory(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "categoryID", "update tenant hr case category"); ok {
		h.updateHRCaseCategoryForTenant(w, r, tenantID, id, "update tenant hr case category")
	}
}

func (h *Handler) DeleteTenantHRCaseCategory(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "categoryID", "delete tenant hr case category"); ok {
		if err := h.svc.DeleteHRCaseCategory(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondHRCaseError(w, r, "delete tenant hr case category", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListTenantHRCaseSLAPolicies(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant hr case sla policies"); ok {
		h.listHRCaseSLAPoliciesForTenant(w, r, tenantID, "list tenant hr case sla policies")
	}
}

func (h *Handler) CreateTenantHRCaseSLAPolicy(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant hr case sla policy"); ok {
		h.createHRCaseSLAPolicyForTenant(w, r, tenantID, "create tenant hr case sla policy")
	}
}

func (h *Handler) UpdateTenantHRCaseSLAPolicy(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "slaPolicyID", "update tenant hr case sla policy"); ok {
		h.updateHRCaseSLAPolicyForTenant(w, r, tenantID, id, "update tenant hr case sla policy")
	}
}

func (h *Handler) DeleteTenantHRCaseSLAPolicy(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "slaPolicyID", "delete tenant hr case sla policy"); ok {
		if err := h.svc.DeleteHRCaseSLAPolicy(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondHRCaseError(w, r, "delete tenant hr case sla policy", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListTenantHRCases(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant hr cases"); ok {
		h.listHRCasesForTenant(w, r, tenantID, "list tenant hr cases")
	}
}

func (h *Handler) CreateTenantHRCase(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant hr case"); ok {
		h.createHRCaseForTenant(w, r, tenantID, "create tenant hr case")
	}
}

func (h *Handler) GetTenantHRCase(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "caseID", "get tenant hr case"); ok {
		h.getHRCaseForTenant(w, r, tenantID, id, "get tenant hr case")
	}
}

func (h *Handler) UpdateTenantHRCase(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "caseID", "update tenant hr case"); ok {
		h.updateHRCaseForTenant(w, r, tenantID, id, "update tenant hr case")
	}
}

func (h *Handler) AssignTenantHRCase(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "caseID", "assign tenant hr case"); ok {
		h.assignHRCaseForTenant(w, r, tenantID, id, "assign tenant hr case")
	}
}

func (h *Handler) UpdateTenantHRCaseStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "caseID", "update tenant hr case status"); ok {
		h.updateHRCaseStatusForTenant(w, r, tenantID, id, "update tenant hr case status")
	}
}

func (h *Handler) CreateTenantHRCaseComment(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "caseID", "create tenant hr case comment"); ok {
		h.createHRCaseCommentForTenant(w, r, tenantID, id, "create tenant hr case comment")
	}
}

func (h *Handler) CreateTenantHRCaseAttachment(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.superAdminHRCaseIDPair(w, r, "caseID", "create tenant hr case attachment"); ok {
		h.createHRCaseAttachmentForTenant(w, r, tenantID, id, "create tenant hr case attachment")
	}
}

func (h *Handler) listHRCaseCategoriesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListHRCaseCategories(r.Context(), tenantID, optionalBoolQuery(r, "active"))
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createHRCaseCategoryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.HRCaseCategoryCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateHRCaseCategory(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateHRCaseCategoryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.HRCaseCategoryCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateHRCaseCategory(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listHRCaseSLAPoliciesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListHRCaseSLAPolicies(r.Context(), tenantID, optionalUUIDQuery(r, "category_id"), optionalStringQuery(r, "priority"))
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createHRCaseSLAPolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.HRCaseSLAPolicyCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateHRCaseSLAPolicy(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateHRCaseSLAPolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.HRCaseSLAPolicyCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateHRCaseSLAPolicy(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listHRCasesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter := domain.HRCaseFilter{
		TenantID:              tenantID,
		Status:                optionalStringQuery(r, "status"),
		Priority:              optionalStringQuery(r, "priority"),
		CategoryID:            optionalUUIDQuery(r, "category_id"),
		RequesterUserID:       optionalUUIDQuery(r, "requester_user_id"),
		SubjectEmployeeUserID: optionalUUIDQuery(r, "subject_employee_user_id"),
		OwnerUserID:           optionalUUIDQuery(r, "owner_user_id"),
		ConfidentialityLevel:  optionalStringQuery(r, "confidentiality_level"),
		Search:                optionalStringQuery(r, "search"),
		Limit:                 queryInt32(r, "limit", 50),
		Offset:                queryInt32(r, "offset", 0),
	}
	page, err := h.svc.ListHRCases(r.Context(), filter)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) createHRCaseForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.HRCaseCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateHRCase(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getHRCaseForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	includeInternal := optionalBoolQuery(r, "include_internal")
	include := includeInternal != nil && *includeInternal
	item, err := h.svc.GetHRCaseWorkspace(r.Context(), tenantID, id, include)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateHRCaseForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.HRCaseCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateHRCase(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) assignHRCaseForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.HRCaseAssignmentCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.AssignHRCase(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateHRCaseStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.HRCaseStatusCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateHRCaseStatus(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createHRCaseCommentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.HRCaseCommentCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.CaseID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateHRCaseComment(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) createHRCaseAttachmentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.HRCaseAttachmentCommand
	if !h.decodeHRCaseJSON(w, r, operation, &cmd) {
		return
	}
	cmd.TenantID = tenantID
	cmd.CaseID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateHRCaseAttachment(r.Context(), cmd)
	if err != nil {
		h.respondHRCaseError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) hrCaseIDPair(w http.ResponseWriter, r *http.Request, idParam string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, idParam))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid record id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminHRCaseIDPair(w http.ResponseWriter, r *http.Request, idParam string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, idParam))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid record id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) decodeHRCaseJSON(w http.ResponseWriter, r *http.Request, operation string, dest any) bool {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return false
	}
	return true
}

func (h *Handler) respondHRCaseError(w http.ResponseWriter, r *http.Request, operation string, err error) {
	status := http.StatusBadRequest
	message := err.Error()
	if errors.Is(err, domain.ErrHRCaseNotFound) || errors.Is(err, domain.ErrHRCaseCategoryNotFound) || errors.Is(err, domain.ErrHRCaseSLANotFound) {
		status = http.StatusNotFound
		message = "record not found"
	} else if errors.Is(err, domain.ErrStorageProviderSettingsNotFound) {
		status = http.StatusConflict
		message = "storage provider settings are required"
	}
	h.respondError(w, r, status, operation, err, message)
}
