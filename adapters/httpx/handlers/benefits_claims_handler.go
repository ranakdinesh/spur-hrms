package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateBenefitPlan(w http.ResponseWriter, r *http.Request) {
	h.createBenefitPlan(w, r, h.currentTenant(w, r, "create benefit plan"), "create benefit plan")
}
func (h *Handler) CreateTenantBenefitPlan(w http.ResponseWriter, r *http.Request) {
	h.createBenefitPlan(w, r, h.superTenant(w, r, "create tenant benefit plan"), "create tenant benefit plan")
}
func (h *Handler) ListBenefitPlans(w http.ResponseWriter, r *http.Request) {
	h.listBenefitPlans(w, r, h.currentTenant(w, r, "list benefit plans"), "list benefit plans")
}
func (h *Handler) ListTenantBenefitPlans(w http.ResponseWriter, r *http.Request) {
	h.listBenefitPlans(w, r, h.superTenant(w, r, "list tenant benefit plans"), "list tenant benefit plans")
}
func (h *Handler) UpdateBenefitPlan(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "planID", "update benefit plan")
	if ok {
		h.updateBenefitPlan(w, r, tenantID, id, "update benefit plan")
	}
}
func (h *Handler) UpdateTenantBenefitPlan(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "planID", "update tenant benefit plan")
	if ok {
		h.updateBenefitPlan(w, r, tenantID, id, "update tenant benefit plan")
	}
}
func (h *Handler) DeleteBenefitPlan(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "planID", "delete benefit plan")
	if ok {
		h.deleteBenefitPlan(w, r, tenantID, id, "delete benefit plan")
	}
}
func (h *Handler) DeleteTenantBenefitPlan(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "planID", "delete tenant benefit plan")
	if ok {
		h.deleteBenefitPlan(w, r, tenantID, id, "delete tenant benefit plan")
	}
}

func (h *Handler) CreateBenefitEnrollmentWindow(w http.ResponseWriter, r *http.Request) {
	h.createBenefitWindow(w, r, h.currentTenant(w, r, "create benefit window"), "create benefit window")
}
func (h *Handler) CreateTenantBenefitEnrollmentWindow(w http.ResponseWriter, r *http.Request) {
	h.createBenefitWindow(w, r, h.superTenant(w, r, "create tenant benefit window"), "create tenant benefit window")
}
func (h *Handler) ListBenefitEnrollmentWindows(w http.ResponseWriter, r *http.Request) {
	h.listBenefitWindows(w, r, h.currentTenant(w, r, "list benefit windows"), "list benefit windows")
}
func (h *Handler) ListTenantBenefitEnrollmentWindows(w http.ResponseWriter, r *http.Request) {
	h.listBenefitWindows(w, r, h.superTenant(w, r, "list tenant benefit windows"), "list tenant benefit windows")
}
func (h *Handler) UpdateBenefitEnrollmentWindow(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "windowID", "update benefit window")
	if ok {
		h.updateBenefitWindow(w, r, tenantID, id, "update benefit window")
	}
}
func (h *Handler) UpdateTenantBenefitEnrollmentWindow(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "windowID", "update tenant benefit window")
	if ok {
		h.updateBenefitWindow(w, r, tenantID, id, "update tenant benefit window")
	}
}
func (h *Handler) DeleteBenefitEnrollmentWindow(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "windowID", "delete benefit window")
	if ok {
		h.deleteBenefitWindow(w, r, tenantID, id, "delete benefit window")
	}
}
func (h *Handler) DeleteTenantBenefitEnrollmentWindow(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "windowID", "delete tenant benefit window")
	if ok {
		h.deleteBenefitWindow(w, r, tenantID, id, "delete tenant benefit window")
	}
}

func (h *Handler) CreateBenefitDependent(w http.ResponseWriter, r *http.Request) {
	h.createBenefitDependent(w, r, h.currentTenant(w, r, "create benefit dependent"), "create benefit dependent")
}
func (h *Handler) CreateTenantBenefitDependent(w http.ResponseWriter, r *http.Request) {
	h.createBenefitDependent(w, r, h.superTenant(w, r, "create tenant benefit dependent"), "create tenant benefit dependent")
}
func (h *Handler) ListBenefitDependents(w http.ResponseWriter, r *http.Request) {
	h.listBenefitDependents(w, r, h.currentTenant(w, r, "list benefit dependents"), "list benefit dependents")
}
func (h *Handler) ListTenantBenefitDependents(w http.ResponseWriter, r *http.Request) {
	h.listBenefitDependents(w, r, h.superTenant(w, r, "list tenant benefit dependents"), "list tenant benefit dependents")
}
func (h *Handler) UpdateBenefitDependent(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "dependentID", "update benefit dependent")
	if ok {
		h.updateBenefitDependent(w, r, tenantID, id, "update benefit dependent")
	}
}
func (h *Handler) UpdateTenantBenefitDependent(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "dependentID", "update tenant benefit dependent")
	if ok {
		h.updateBenefitDependent(w, r, tenantID, id, "update tenant benefit dependent")
	}
}
func (h *Handler) DeleteBenefitDependent(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "dependentID", "delete benefit dependent")
	if ok {
		h.deleteBenefitDependent(w, r, tenantID, id, "delete benefit dependent")
	}
}
func (h *Handler) DeleteTenantBenefitDependent(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "dependentID", "delete tenant benefit dependent")
	if ok {
		h.deleteBenefitDependent(w, r, tenantID, id, "delete tenant benefit dependent")
	}
}

func (h *Handler) CreateBenefitEnrollment(w http.ResponseWriter, r *http.Request) {
	h.createBenefitEnrollment(w, r, h.currentTenant(w, r, "create benefit enrollment"), "create benefit enrollment")
}
func (h *Handler) CreateTenantBenefitEnrollment(w http.ResponseWriter, r *http.Request) {
	h.createBenefitEnrollment(w, r, h.superTenant(w, r, "create tenant benefit enrollment"), "create tenant benefit enrollment")
}
func (h *Handler) ListBenefitEnrollments(w http.ResponseWriter, r *http.Request) {
	h.listBenefitEnrollments(w, r, h.currentTenant(w, r, "list benefit enrollments"), "list benefit enrollments")
}
func (h *Handler) ListTenantBenefitEnrollments(w http.ResponseWriter, r *http.Request) {
	h.listBenefitEnrollments(w, r, h.superTenant(w, r, "list tenant benefit enrollments"), "list tenant benefit enrollments")
}
func (h *Handler) UpdateBenefitEnrollment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "enrollmentID", "update benefit enrollment")
	if ok {
		h.updateBenefitEnrollment(w, r, tenantID, id, "update benefit enrollment")
	}
}
func (h *Handler) UpdateTenantBenefitEnrollment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "enrollmentID", "update tenant benefit enrollment")
	if ok {
		h.updateBenefitEnrollment(w, r, tenantID, id, "update tenant benefit enrollment")
	}
}
func (h *Handler) UpdateBenefitEnrollmentStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "enrollmentID", "update benefit enrollment status")
	if ok {
		h.updateBenefitEnrollmentStatus(w, r, tenantID, id, "update benefit enrollment status")
	}
}
func (h *Handler) UpdateTenantBenefitEnrollmentStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "enrollmentID", "update tenant benefit enrollment status")
	if ok {
		h.updateBenefitEnrollmentStatus(w, r, tenantID, id, "update tenant benefit enrollment status")
	}
}
func (h *Handler) DeleteBenefitEnrollment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "enrollmentID", "delete benefit enrollment")
	if ok {
		h.deleteBenefitEnrollment(w, r, tenantID, id, "delete benefit enrollment")
	}
}
func (h *Handler) DeleteTenantBenefitEnrollment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "enrollmentID", "delete tenant benefit enrollment")
	if ok {
		h.deleteBenefitEnrollment(w, r, tenantID, id, "delete tenant benefit enrollment")
	}
}

func (h *Handler) CreateBenefitClaimType(w http.ResponseWriter, r *http.Request) {
	h.createBenefitClaimType(w, r, h.currentTenant(w, r, "create benefit claim type"), "create benefit claim type")
}
func (h *Handler) CreateTenantBenefitClaimType(w http.ResponseWriter, r *http.Request) {
	h.createBenefitClaimType(w, r, h.superTenant(w, r, "create tenant benefit claim type"), "create tenant benefit claim type")
}
func (h *Handler) ListBenefitClaimTypes(w http.ResponseWriter, r *http.Request) {
	h.listBenefitClaimTypes(w, r, h.currentTenant(w, r, "list benefit claim types"), "list benefit claim types")
}
func (h *Handler) ListTenantBenefitClaimTypes(w http.ResponseWriter, r *http.Request) {
	h.listBenefitClaimTypes(w, r, h.superTenant(w, r, "list tenant benefit claim types"), "list tenant benefit claim types")
}
func (h *Handler) UpdateBenefitClaimType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "claimTypeID", "update benefit claim type")
	if ok {
		h.updateBenefitClaimType(w, r, tenantID, id, "update benefit claim type")
	}
}
func (h *Handler) UpdateTenantBenefitClaimType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "claimTypeID", "update tenant benefit claim type")
	if ok {
		h.updateBenefitClaimType(w, r, tenantID, id, "update tenant benefit claim type")
	}
}
func (h *Handler) DeleteBenefitClaimType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "claimTypeID", "delete benefit claim type")
	if ok {
		h.deleteBenefitClaimType(w, r, tenantID, id, "delete benefit claim type")
	}
}
func (h *Handler) DeleteTenantBenefitClaimType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "claimTypeID", "delete tenant benefit claim type")
	if ok {
		h.deleteBenefitClaimType(w, r, tenantID, id, "delete tenant benefit claim type")
	}
}

func (h *Handler) CreateBenefitClaim(w http.ResponseWriter, r *http.Request) {
	h.createBenefitClaim(w, r, h.currentTenant(w, r, "create benefit claim"), "create benefit claim")
}
func (h *Handler) CreateTenantBenefitClaim(w http.ResponseWriter, r *http.Request) {
	h.createBenefitClaim(w, r, h.superTenant(w, r, "create tenant benefit claim"), "create tenant benefit claim")
}
func (h *Handler) ListBenefitClaims(w http.ResponseWriter, r *http.Request) {
	h.listBenefitClaims(w, r, h.currentTenant(w, r, "list benefit claims"), "list benefit claims")
}
func (h *Handler) ListTenantBenefitClaims(w http.ResponseWriter, r *http.Request) {
	h.listBenefitClaims(w, r, h.superTenant(w, r, "list tenant benefit claims"), "list tenant benefit claims")
}
func (h *Handler) UpdateBenefitClaim(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "claimID", "update benefit claim")
	if ok {
		h.updateBenefitClaim(w, r, tenantID, id, "update benefit claim")
	}
}
func (h *Handler) UpdateTenantBenefitClaim(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "claimID", "update tenant benefit claim")
	if ok {
		h.updateBenefitClaim(w, r, tenantID, id, "update tenant benefit claim")
	}
}
func (h *Handler) UpdateBenefitClaimStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "claimID", "update benefit claim status")
	if ok {
		h.updateBenefitClaimStatus(w, r, tenantID, id, "update benefit claim status")
	}
}
func (h *Handler) UpdateTenantBenefitClaimStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "claimID", "update tenant benefit claim status")
	if ok {
		h.updateBenefitClaimStatus(w, r, tenantID, id, "update tenant benefit claim status")
	}
}
func (h *Handler) DeleteBenefitClaim(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "claimID", "delete benefit claim")
	if ok {
		h.deleteBenefitClaim(w, r, tenantID, id, "delete benefit claim")
	}
}
func (h *Handler) DeleteTenantBenefitClaim(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "claimID", "delete tenant benefit claim")
	if ok {
		h.deleteBenefitClaim(w, r, tenantID, id, "delete tenant benefit claim")
	}
}
func (h *Handler) CreateBenefitClaimAttachment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "claimID", "create benefit claim attachment")
	if ok {
		h.createBenefitClaimAttachment(w, r, tenantID, id, "create benefit claim attachment")
	}
}
func (h *Handler) CreateTenantBenefitClaimAttachment(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "claimID", "create tenant benefit claim attachment")
	if ok {
		h.createBenefitClaimAttachment(w, r, tenantID, id, "create tenant benefit claim attachment")
	}
}
func (h *Handler) ListBenefitClaimAttachments(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.tenantAndURLUUID(w, r, "claimID", "list benefit claim attachments")
	if ok {
		h.listBenefitClaimAttachments(w, r, tenantID, id, "list benefit claim attachments")
	}
}
func (h *Handler) ListTenantBenefitClaimAttachments(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTenantAndURLUUID(w, r, "claimID", "list tenant benefit claim attachments")
	if ok {
		h.listBenefitClaimAttachments(w, r, tenantID, id, "list tenant benefit claim attachments")
	}
}
func (h *Handler) ListBenefitEvents(w http.ResponseWriter, r *http.Request) {
	h.listBenefitEvents(w, r, h.currentTenant(w, r, "list benefit events"), "list benefit events")
}
func (h *Handler) ListTenantBenefitEvents(w http.ResponseWriter, r *http.Request) {
	h.listBenefitEvents(w, r, h.superTenant(w, r, "list tenant benefit events"), "list tenant benefit events")
}
func (h *Handler) GetBenefitsSummary(w http.ResponseWriter, r *http.Request) {
	h.getBenefitsSummary(w, r, h.currentTenant(w, r, "get benefits summary"), "get benefits summary")
}
func (h *Handler) GetTenantBenefitsSummary(w http.ResponseWriter, r *http.Request) {
	h.getBenefitsSummary(w, r, h.superTenant(w, r, "get tenant benefits summary"), "get tenant benefits summary")
}

func (h *Handler) createBenefitPlan(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	var cmd ports.BenefitPlanCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateBenefitPlan(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateBenefitPlan(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.BenefitPlanCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateBenefitPlan(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listBenefitPlans(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	items, err := h.svc.ListBenefitPlans(r.Context(), h.benefitFilter(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list benefit plans")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) deleteBenefitPlan(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if err := h.svc.DeleteBenefitPlan(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createBenefitWindow(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.BenefitEnrollmentWindowCommand
	if tenantID == uuid.Nil || !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateBenefitEnrollmentWindow(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) updateBenefitWindow(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.BenefitEnrollmentWindowCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateBenefitEnrollmentWindow(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) listBenefitWindows(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	items, err := h.svc.ListBenefitEnrollmentWindows(r.Context(), h.benefitFilter(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list benefit windows")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) deleteBenefitWindow(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if err := h.svc.DeleteBenefitEnrollmentWindow(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createBenefitDependent(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.BenefitDependentCommand
	if tenantID == uuid.Nil || !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateBenefitDependent(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) updateBenefitDependent(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.BenefitDependentCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateBenefitDependent(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) listBenefitDependents(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	items, err := h.svc.ListBenefitDependents(r.Context(), h.benefitFilter(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list benefit dependents")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) deleteBenefitDependent(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if err := h.svc.DeleteBenefitDependent(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createBenefitEnrollment(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.BenefitEnrollmentCommand
	if tenantID == uuid.Nil || !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateBenefitEnrollment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) updateBenefitEnrollment(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.BenefitEnrollmentCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateBenefitEnrollment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) listBenefitEnrollments(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	items, err := h.svc.ListBenefitEnrollments(r.Context(), h.benefitFilter(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list benefit enrollments")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) updateBenefitEnrollmentStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.BenefitEnrollmentStatusCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateBenefitEnrollmentStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) deleteBenefitEnrollment(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if err := h.svc.DeleteBenefitEnrollment(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createBenefitClaimType(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.BenefitClaimTypeCommand
	if tenantID == uuid.Nil || !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateBenefitClaimType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) updateBenefitClaimType(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.BenefitClaimTypeCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateBenefitClaimType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) listBenefitClaimTypes(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	items, err := h.svc.ListBenefitClaimTypes(r.Context(), h.benefitFilter(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list benefit claim types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) deleteBenefitClaimType(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if err := h.svc.DeleteBenefitClaimType(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createBenefitClaim(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.BenefitClaimCommand
	if tenantID == uuid.Nil || !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateBenefitClaim(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) updateBenefitClaim(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.BenefitClaimCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateBenefitClaim(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) listBenefitClaims(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	items, err := h.svc.ListBenefitClaims(r.Context(), h.benefitFilter(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list benefit claims")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) updateBenefitClaimStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.BenefitClaimStatusCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateBenefitClaimStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}
func (h *Handler) deleteBenefitClaim(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	if err := h.svc.DeleteBenefitClaim(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) createBenefitClaimAttachment(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, claimID uuid.UUID, operation string) {
	var cmd ports.BenefitClaimAttachmentCommand
	if !h.decodeBenefitBody(w, r, &cmd, operation) {
		return
	}
	cmd.TenantID, cmd.ClaimID, cmd.ActorID = tenantID, claimID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateBenefitClaimAttachment(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}
func (h *Handler) listBenefitClaimAttachments(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, claimID uuid.UUID, operation string) {
	items, err := h.svc.ListBenefitClaimAttachments(r.Context(), tenantID, claimID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list benefit claim attachments")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) listBenefitEvents(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	items, err := h.svc.ListBenefitEvents(r.Context(), h.benefitFilter(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list benefit events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}
func (h *Handler) getBenefitsSummary(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	if tenantID == uuid.Nil {
		return
	}
	items, err := h.svc.GetBenefitsSummary(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to get benefits summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) currentTenant(w http.ResponseWriter, r *http.Request, operation string) uuid.UUID {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil
	}
	return tenantID
}

func (h *Handler) superTenant(w http.ResponseWriter, r *http.Request, operation string) uuid.UUID {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil
	}
	return tenantID
}

func (h *Handler) decodeBenefitBody(w http.ResponseWriter, r *http.Request, v any, operation string) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return false
	}
	return true
}

func (h *Handler) benefitFilter(r *http.Request, tenantID uuid.UUID) domain.BenefitFilter {
	return domain.BenefitFilter{
		TenantID: tenantID, ActiveOnly: optionalBoolQuery(r, "active_only"), Search: optionalStringQuery(r, "search"),
		PlanType: optionalStringQuery(r, "plan_type"), PlanID: optionalUUIDQuery(r, "plan_id"), Status: optionalStringQuery(r, "status"),
		EmployeeUserID: optionalUUIDQuery(r, "employee_user_id"), ClaimTypeID: optionalUUIDQuery(r, "claim_type_id"),
		PaymentStatus: optionalStringQuery(r, "payment_status"), PayrollExportStatus: optionalStringQuery(r, "payroll_export_status"),
		NomineesOnly: optionalBoolQuery(r, "nominees_only"), SourceType: optionalStringQuery(r, "source_type"), SourceID: optionalUUIDQuery(r, "source_id"),
		Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0),
	}
}
