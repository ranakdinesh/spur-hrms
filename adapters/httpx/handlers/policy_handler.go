package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreatePolicyType(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create policy type", err, "tenant context is required")
		return
	}
	h.createPolicyTypeForTenant(w, r, tenantID, "create policy type")
}

func (h *Handler) ListPolicyTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list policy types", err, "tenant context is required")
		return
	}
	h.listPolicyTypesForTenant(w, r, tenantID, "list policy types")
}

func (h *Handler) GetPolicyType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.policyTypeRequestIDs(w, r, "get policy type")
	if !ok {
		return
	}
	item, err := h.svc.GetPolicyType(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get policy type", err, "policy type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdatePolicyType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.policyTypeRequestIDs(w, r, "update policy type")
	if !ok {
		return
	}
	h.updatePolicyTypeForTenant(w, r, tenantID, id, "update policy type")
}

func (h *Handler) DeletePolicyType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.policyTypeRequestIDs(w, r, "delete policy type")
	if !ok {
		return
	}
	if err := h.svc.DeletePolicyType(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete policy type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateCompanyPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create company policy", err, "tenant context is required")
		return
	}
	h.createCompanyPolicyForTenant(w, r, tenantID, "create company policy")
}

func (h *Handler) ListCompanyPolicies(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list company policies", err, "tenant context is required")
		return
	}
	h.listCompanyPoliciesForTenant(w, r, tenantID, "list company policies")
}

func (h *Handler) GetCompanyPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.companyPolicyRequestIDs(w, r, "get company policy")
	if !ok {
		return
	}
	item, err := h.svc.GetCompanyPolicy(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get company policy", err, "company policy not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateCompanyPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.companyPolicyRequestIDs(w, r, "update company policy")
	if !ok {
		return
	}
	h.updateCompanyPolicyForTenant(w, r, tenantID, id, "update company policy")
}

func (h *Handler) DeleteCompanyPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.companyPolicyRequestIDs(w, r, "delete company policy")
	if !ok {
		return
	}
	if err := h.svc.DeleteCompanyPolicy(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete company policy", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantPolicyType(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant policy type")
	if !ok {
		return
	}
	h.createPolicyTypeForTenant(w, r, tenantID, "create tenant policy type")
}

func (h *Handler) ListTenantPolicyTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant policy types")
	if !ok {
		return
	}
	h.listPolicyTypesForTenant(w, r, tenantID, "list tenant policy types")
}

func (h *Handler) GetTenantPolicyType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPolicyTypeRequestIDs(w, r, "get tenant policy type")
	if !ok {
		return
	}
	item, err := h.svc.GetPolicyType(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant policy type", err, "policy type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantPolicyType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPolicyTypeRequestIDs(w, r, "update tenant policy type")
	if !ok {
		return
	}
	h.updatePolicyTypeForTenant(w, r, tenantID, id, "update tenant policy type")
}

func (h *Handler) DeleteTenantPolicyType(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminPolicyTypeRequestIDs(w, r, "delete tenant policy type")
	if !ok {
		return
	}
	if err := h.svc.DeletePolicyType(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant policy type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantCompanyPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant company policy")
	if !ok {
		return
	}
	h.createCompanyPolicyForTenant(w, r, tenantID, "create tenant company policy")
}

func (h *Handler) ListTenantCompanyPolicies(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant company policies")
	if !ok {
		return
	}
	h.listCompanyPoliciesForTenant(w, r, tenantID, "list tenant company policies")
}

func (h *Handler) GetTenantCompanyPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminCompanyPolicyRequestIDs(w, r, "get tenant company policy")
	if !ok {
		return
	}
	item, err := h.svc.GetCompanyPolicy(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant company policy", err, "company policy not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantCompanyPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminCompanyPolicyRequestIDs(w, r, "update tenant company policy")
	if !ok {
		return
	}
	h.updateCompanyPolicyForTenant(w, r, tenantID, id, "update tenant company policy")
}

func (h *Handler) DeleteTenantCompanyPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminCompanyPolicyRequestIDs(w, r, "delete tenant company policy")
	if !ok {
		return
	}
	if err := h.svc.DeleteCompanyPolicy(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant company policy", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createPolicyTypeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PolicyTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.IsSystem = false
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePolicyType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updatePolicyTypeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.PolicyTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePolicyType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listPolicyTypesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListPolicyTypes(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list policy types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createCompanyPolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CompanyPolicyCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateCompanyPolicy(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateCompanyPolicyForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.CompanyPolicyCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCompanyPolicy(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listCompanyPoliciesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	policyTypeID, ok := h.optionalUUIDQuery(w, r, "policy_type_id", operation)
	if !ok {
		return
	}
	items, err := h.svc.ListCompanyPolicies(r.Context(), tenantID, policyTypeID)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) policyTypeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "policyTypeID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid policy type id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) companyPolicyRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "policyID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid policy id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminPolicyTypeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "policyTypeID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid policy type id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminCompanyPolicyRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "policyID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid policy id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
