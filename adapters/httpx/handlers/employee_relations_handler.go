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

func (h *Handler) ListERCaseCategories(w http.ResponseWriter, r *http.Request) {
	h.listERCaseCategories(w, r, h.currentTenant(w, r, "list er case categories"), "list er case categories")
}

func (h *Handler) CreateERCaseCategory(w http.ResponseWriter, r *http.Request) {
	h.createERCaseCategory(w, r, h.currentTenant(w, r, "create er case category"), "create er case category")
}

func (h *Handler) UpdateERCaseCategory(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "update er case category")
	id, ok := erIDParam(w, r, h, "categoryID", "update er case category")
	if ok {
		h.updateERCaseCategory(w, r, tenantID, id, "update er case category")
	}
}

func (h *Handler) DeleteERCaseCategory(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "delete er case category")
	id, ok := erIDParam(w, r, h, "categoryID", "delete er case category")
	if !ok || tenantID == uuid.Nil {
		return
	}
	if err := h.svc.DeleteERCaseCategory(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondERError(w, r, "delete er case category", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListERCases(w http.ResponseWriter, r *http.Request) {
	h.listERCases(w, r, h.currentTenant(w, r, "list er cases"), "list er cases")
}

func (h *Handler) CreateERCase(w http.ResponseWriter, r *http.Request) {
	h.createERCase(w, r, h.currentTenant(w, r, "create er case"), "create er case")
}

func (h *Handler) GetERCase(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "get er case")
	id, ok := erIDParam(w, r, h, "caseID", "get er case")
	if ok {
		h.getERCase(w, r, tenantID, id, "get er case")
	}
}

func (h *Handler) UpdateERCase(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "update er case")
	id, ok := erIDParam(w, r, h, "caseID", "update er case")
	if ok {
		h.updateERCase(w, r, tenantID, id, "update er case")
	}
}

func (h *Handler) UpdateERCaseStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "update er case status")
	id, ok := erIDParam(w, r, h, "caseID", "update er case status")
	if ok {
		h.updateERCaseStatus(w, r, tenantID, id, "update er case status")
	}
}

func (h *Handler) UpdateERCaseLegalHold(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "update er case legal hold")
	id, ok := erIDParam(w, r, h, "caseID", "update er case legal hold")
	if ok {
		h.updateERCaseLegalHold(w, r, tenantID, id, "update er case legal hold")
	}
}

func (h *Handler) CreateERCaseParty(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "create er party")
	caseID, ok := erIDParam(w, r, h, "caseID", "create er party")
	if ok {
		h.createERParty(w, r, tenantID, caseID, "create er party")
	}
}

func (h *Handler) CreateERAllegation(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "create er allegation")
	caseID, ok := erIDParam(w, r, h, "caseID", "create er allegation")
	if ok {
		h.createERAllegation(w, r, tenantID, caseID, "create er allegation")
	}
}

func (h *Handler) CreateERInvestigationStep(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "create er step")
	caseID, ok := erIDParam(w, r, h, "caseID", "create er step")
	if ok {
		h.createERStep(w, r, tenantID, caseID, "create er step")
	}
}

func (h *Handler) UpdateERInvestigationStepStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "update er step status")
	stepID, ok := erIDParam(w, r, h, "stepID", "update er step status")
	if !ok || tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERInvestigationStepCommand
	if !h.decodeERBody(w, r, &cmd, "update er step status") {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = stepID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateERInvestigationStepStatus(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, "update er step status", err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateERWitnessNote(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "create er witness")
	caseID, ok := erIDParam(w, r, h, "caseID", "create er witness")
	if ok {
		h.createERWitness(w, r, tenantID, caseID, "create er witness")
	}
}

func (h *Handler) CreateEREvidenceAttachment(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "create er evidence")
	caseID, ok := erIDParam(w, r, h, "caseID", "create er evidence")
	if ok {
		h.createEREvidence(w, r, tenantID, caseID, "create er evidence")
	}
}

func (h *Handler) CreateERFinding(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "create er finding")
	caseID, ok := erIDParam(w, r, h, "caseID", "create er finding")
	if ok {
		h.createERFinding(w, r, tenantID, caseID, "create er finding")
	}
}

func (h *Handler) CreateERActionPlan(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "create er action")
	caseID, ok := erIDParam(w, r, h, "caseID", "create er action")
	if ok {
		h.createERAction(w, r, tenantID, caseID, "create er action")
	}
}

func (h *Handler) UpdateERActionPlanStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := h.currentTenant(w, r, "update er action status")
	actionID, ok := erIDParam(w, r, h, "actionID", "update er action status")
	if !ok || tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERActionPlanCommand
	if !h.decodeERBody(w, r, &cmd, "update er action status") {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = actionID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateERActionPlanStatus(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, "update er action status", err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ListTenantERCaseCategories(w http.ResponseWriter, r *http.Request) {
	h.listERCaseCategories(w, r, h.superTenant(w, r, "list tenant er case categories"), "list tenant er case categories")
}

func (h *Handler) CreateTenantERCaseCategory(w http.ResponseWriter, r *http.Request) {
	h.createERCaseCategory(w, r, h.superTenant(w, r, "create tenant er case category"), "create tenant er case category")
}

func (h *Handler) UpdateTenantERCaseCategory(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "update tenant er case category")
	id, ok := erIDParam(w, r, h, "categoryID", "update tenant er case category")
	if ok {
		h.updateERCaseCategory(w, r, tenantID, id, "update tenant er case category")
	}
}

func (h *Handler) DeleteTenantERCaseCategory(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "delete tenant er case category")
	id, ok := erIDParam(w, r, h, "categoryID", "delete tenant er case category")
	if !ok || tenantID == uuid.Nil {
		return
	}
	if err := h.svc.DeleteERCaseCategory(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondERError(w, r, "delete tenant er case category", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantERCases(w http.ResponseWriter, r *http.Request) {
	h.listERCases(w, r, h.superTenant(w, r, "list tenant er cases"), "list tenant er cases")
}

func (h *Handler) CreateTenantERCase(w http.ResponseWriter, r *http.Request) {
	h.createERCase(w, r, h.superTenant(w, r, "create tenant er case"), "create tenant er case")
}

func (h *Handler) GetTenantERCase(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "get tenant er case")
	id, ok := erIDParam(w, r, h, "caseID", "get tenant er case")
	if ok {
		h.getERCase(w, r, tenantID, id, "get tenant er case")
	}
}

func (h *Handler) UpdateTenantERCase(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "update tenant er case")
	id, ok := erIDParam(w, r, h, "caseID", "update tenant er case")
	if ok {
		h.updateERCase(w, r, tenantID, id, "update tenant er case")
	}
}

func (h *Handler) UpdateTenantERCaseStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "update tenant er case status")
	id, ok := erIDParam(w, r, h, "caseID", "update tenant er case status")
	if ok {
		h.updateERCaseStatus(w, r, tenantID, id, "update tenant er case status")
	}
}

func (h *Handler) UpdateTenantERCaseLegalHold(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "update tenant er case legal hold")
	id, ok := erIDParam(w, r, h, "caseID", "update tenant er case legal hold")
	if ok {
		h.updateERCaseLegalHold(w, r, tenantID, id, "update tenant er case legal hold")
	}
}

func (h *Handler) CreateTenantERCaseParty(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "create tenant er party")
	caseID, ok := erIDParam(w, r, h, "caseID", "create tenant er party")
	if ok {
		h.createERParty(w, r, tenantID, caseID, "create tenant er party")
	}
}

func (h *Handler) CreateTenantERAllegation(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "create tenant er allegation")
	caseID, ok := erIDParam(w, r, h, "caseID", "create tenant er allegation")
	if ok {
		h.createERAllegation(w, r, tenantID, caseID, "create tenant er allegation")
	}
}

func (h *Handler) CreateTenantERInvestigationStep(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "create tenant er step")
	caseID, ok := erIDParam(w, r, h, "caseID", "create tenant er step")
	if ok {
		h.createERStep(w, r, tenantID, caseID, "create tenant er step")
	}
}

func (h *Handler) CreateTenantERWitnessNote(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "create tenant er witness")
	caseID, ok := erIDParam(w, r, h, "caseID", "create tenant er witness")
	if ok {
		h.createERWitness(w, r, tenantID, caseID, "create tenant er witness")
	}
}

func (h *Handler) CreateTenantEREvidenceAttachment(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "create tenant er evidence")
	caseID, ok := erIDParam(w, r, h, "caseID", "create tenant er evidence")
	if ok {
		h.createEREvidence(w, r, tenantID, caseID, "create tenant er evidence")
	}
}

func (h *Handler) CreateTenantERFinding(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "create tenant er finding")
	caseID, ok := erIDParam(w, r, h, "caseID", "create tenant er finding")
	if ok {
		h.createERFinding(w, r, tenantID, caseID, "create tenant er finding")
	}
}

func (h *Handler) CreateTenantERActionPlan(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "create tenant er action")
	caseID, ok := erIDParam(w, r, h, "caseID", "create tenant er action")
	if ok {
		h.createERAction(w, r, tenantID, caseID, "create tenant er action")
	}
}

func (h *Handler) UpdateTenantERInvestigationStepStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "update tenant er step status")
	stepID, ok := erIDParam(w, r, h, "stepID", "update tenant er step status")
	if !ok || tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERInvestigationStepCommand
	if !h.decodeERBody(w, r, &cmd, "update tenant er step status") {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = stepID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateERInvestigationStepStatus(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, "update tenant er step status", err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantERActionPlanStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := h.superTenant(w, r, "update tenant er action status")
	actionID, ok := erIDParam(w, r, h, "actionID", "update tenant er action status")
	if !ok || tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERActionPlanCommand
	if !h.decodeERBody(w, r, &cmd, "update tenant er action status") {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = actionID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateERActionPlanStatus(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, "update tenant er action status", err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listERCaseCategories(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	items, err := h.svc.ListERCaseCategories(r.Context(), tenantID, optionalBoolQuery(r, "active_only"))
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createERCaseCategory(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERCaseCategoryCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateERCaseCategory(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateERCaseCategory(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERCaseCategoryCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateERCaseCategory(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listERCases(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	page, err := h.svc.ListERCases(r.Context(), erFilter(r, tenantID))
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) createERCase(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERCaseCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateERCase(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getERCase(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	item, err := h.svc.GetERCaseWorkspace(r.Context(), tenantID, id)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateERCase(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERCaseCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateERCase(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateERCaseStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERCaseStatusCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateERCaseStatus(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateERCaseLegalHold(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERCaseLegalHoldCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateERCaseLegalHold(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createERParty(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, caseID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERCasePartyCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ERCaseID = caseID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateERCaseParty(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) createERAllegation(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, caseID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERAllegationCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ERCaseID = caseID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateERAllegation(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) createERStep(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, caseID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERInvestigationStepCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ERCaseID = caseID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateERInvestigationStep(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) createERWitness(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, caseID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERWitnessNoteCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ERCaseID = caseID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateERWitnessNote(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) createEREvidence(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, caseID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.EREvidenceAttachmentCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ERCaseID = caseID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateEREvidenceAttachment(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) createERFinding(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, caseID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERFindingCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ERCaseID = caseID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateERFinding(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) createERAction(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, caseID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.ERActionPlanCommand
	if !h.decodeERBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID = tenantID
	cmd.ERCaseID = caseID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateERActionPlan(r.Context(), cmd)
	if err != nil {
		h.respondERError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) decodeERBody(w http.ResponseWriter, r *http.Request, v any, operation string) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return false
	}
	return true
}

func erFilter(r *http.Request, tenantID uuid.UUID) domain.ERCaseFilter {
	return domain.ERCaseFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), Severity: optionalStringQuery(r, "severity"), CaseFamily: optionalStringQuery(r, "case_family"), CategoryID: optionalUUIDQuery(r, "category_id"), OwnerUserID: optionalUUIDQuery(r, "owner_user_id"), SubjectEmployeeUserID: optionalUUIDQuery(r, "subject_employee_user_id"), ComplainantUserID: optionalUUIDQuery(r, "complainant_user_id"), LegalHold: optionalBoolQuery(r, "legal_hold"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)}
}

func erIDParam(w http.ResponseWriter, r *http.Request, h *Handler, key string, operation string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, key))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid identifier")
		return uuid.Nil, false
	}
	return id, true
}

func (h *Handler) respondERError(w http.ResponseWriter, r *http.Request, operation string, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, domain.ErrERCaseNotFound) || errors.Is(err, domain.ErrERCaseCategoryNotFound) {
		status = http.StatusNotFound
	}
	if errors.Is(err, domain.ErrStorageProviderSettingsNotFound) {
		status = http.StatusConflict
	}
	h.respondError(w, r, status, operation, err, err.Error())
}
