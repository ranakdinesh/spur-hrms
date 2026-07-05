package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListTalentMarketplaceOpportunities(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list talent marketplace opportunities", err, "tenant context is required")
		return
	}
	h.listTalentMarketplaceOpportunitiesForTenant(w, r, tenantID, "list talent marketplace opportunities")
}

func (h *Handler) CreateTalentMarketplaceOpportunity(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create talent marketplace opportunity", err, "tenant context is required")
		return
	}
	h.createTalentMarketplaceOpportunityForTenant(w, r, tenantID, "create talent marketplace opportunity")
}

func (h *Handler) GetTalentMarketplaceOpportunity(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.talentMarketplaceOpportunityRequestIDs(w, r, "get talent marketplace opportunity")
	if !ok {
		return
	}
	item, err := h.svc.GetTalentMarketplaceOpportunity(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get talent marketplace opportunity", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTalentMarketplaceOpportunity(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.talentMarketplaceOpportunityRequestIDs(w, r, "update talent marketplace opportunity")
	if !ok {
		return
	}
	h.updateTalentMarketplaceOpportunityForTenant(w, r, tenantID, id, "update talent marketplace opportunity")
}

func (h *Handler) DeleteTalentMarketplaceOpportunity(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.talentMarketplaceOpportunityRequestIDs(w, r, "delete talent marketplace opportunity")
	if !ok {
		return
	}
	if err := h.svc.DeleteTalentMarketplaceOpportunity(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete talent marketplace opportunity", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTalentMarketplaceFallback(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.talentMarketplaceOpportunityRequestIDs(w, r, "update talent marketplace fallback")
	if !ok {
		return
	}
	h.updateTalentMarketplaceFallbackForTenant(w, r, tenantID, id, "update talent marketplace fallback")
}

func (h *Handler) ListTalentMarketplaceRecommendations(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.talentMarketplaceOpportunityRequestIDs(w, r, "list talent marketplace recommendations")
	if !ok {
		return
	}
	items, err := h.svc.ListTalentMarketplaceRecommendations(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list talent marketplace recommendations", err, "failed to list talent marketplace recommendations")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListTalentMarketplaceApplications(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list talent marketplace applications", err, "tenant context is required")
		return
	}
	h.listTalentMarketplaceApplicationsForTenant(w, r, tenantID, "list talent marketplace applications")
}

func (h *Handler) CreateTalentMarketplaceApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create talent marketplace application", err, "tenant context is required")
		return
	}
	h.createTalentMarketplaceApplicationForTenant(w, r, tenantID, "create talent marketplace application")
}

func (h *Handler) GetTalentMarketplaceApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.talentMarketplaceApplicationRequestIDs(w, r, "get talent marketplace application")
	if !ok {
		return
	}
	item, err := h.svc.GetTalentMarketplaceApplication(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get talent marketplace application", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTalentMarketplaceApplicationStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.talentMarketplaceApplicationRequestIDs(w, r, "update talent marketplace application status")
	if !ok {
		return
	}
	h.updateTalentMarketplaceApplicationStatusForTenant(w, r, tenantID, id, "update talent marketplace application status")
}

func (h *Handler) ListTalentMarketplaceEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list talent marketplace events", err, "tenant context is required")
		return
	}
	h.listTalentMarketplaceEventsForTenant(w, r, tenantID, "list talent marketplace events")
}

func (h *Handler) ListTenantTalentMarketplaceOpportunities(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant talent marketplace opportunities"); ok {
		h.listTalentMarketplaceOpportunitiesForTenant(w, r, tenantID, "list tenant talent marketplace opportunities")
	}
}

func (h *Handler) CreateTenantTalentMarketplaceOpportunity(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant talent marketplace opportunity"); ok {
		h.createTalentMarketplaceOpportunityForTenant(w, r, tenantID, "create tenant talent marketplace opportunity")
	}
}

func (h *Handler) GetTenantTalentMarketplaceOpportunity(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTalentMarketplaceOpportunityRequestIDs(w, r, "get tenant talent marketplace opportunity")
	if !ok {
		return
	}
	item, err := h.svc.GetTalentMarketplaceOpportunity(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant talent marketplace opportunity", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantTalentMarketplaceOpportunity(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTalentMarketplaceOpportunityRequestIDs(w, r, "update tenant talent marketplace opportunity")
	if !ok {
		return
	}
	h.updateTalentMarketplaceOpportunityForTenant(w, r, tenantID, id, "update tenant talent marketplace opportunity")
}

func (h *Handler) DeleteTenantTalentMarketplaceOpportunity(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTalentMarketplaceOpportunityRequestIDs(w, r, "delete tenant talent marketplace opportunity")
	if !ok {
		return
	}
	if err := h.svc.DeleteTalentMarketplaceOpportunity(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant talent marketplace opportunity", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTenantTalentMarketplaceFallback(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTalentMarketplaceOpportunityRequestIDs(w, r, "update tenant talent marketplace fallback")
	if !ok {
		return
	}
	h.updateTalentMarketplaceFallbackForTenant(w, r, tenantID, id, "update tenant talent marketplace fallback")
}

func (h *Handler) ListTenantTalentMarketplaceRecommendations(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTalentMarketplaceOpportunityRequestIDs(w, r, "list tenant talent marketplace recommendations")
	if !ok {
		return
	}
	items, err := h.svc.ListTalentMarketplaceRecommendations(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant talent marketplace recommendations", err, "failed to list tenant talent marketplace recommendations")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListTenantTalentMarketplaceApplications(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant talent marketplace applications"); ok {
		h.listTalentMarketplaceApplicationsForTenant(w, r, tenantID, "list tenant talent marketplace applications")
	}
}

func (h *Handler) CreateTenantTalentMarketplaceApplication(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant talent marketplace application"); ok {
		h.createTalentMarketplaceApplicationForTenant(w, r, tenantID, "create tenant talent marketplace application")
	}
}

func (h *Handler) GetTenantTalentMarketplaceApplication(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTalentMarketplaceApplicationRequestIDs(w, r, "get tenant talent marketplace application")
	if !ok {
		return
	}
	item, err := h.svc.GetTalentMarketplaceApplication(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant talent marketplace application", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantTalentMarketplaceApplicationStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminTalentMarketplaceApplicationRequestIDs(w, r, "update tenant talent marketplace application status")
	if !ok {
		return
	}
	h.updateTalentMarketplaceApplicationStatusForTenant(w, r, tenantID, id, "update tenant talent marketplace application status")
}

func (h *Handler) ListTenantTalentMarketplaceEvents(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant talent marketplace events"); ok {
		h.listTalentMarketplaceEventsForTenant(w, r, tenantID, "list tenant talent marketplace events")
	}
}

func (h *Handler) listTalentMarketplaceOpportunitiesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter, ok := h.talentMarketplaceOpportunityFilter(w, r, tenantID, operation)
	if !ok {
		return
	}
	items, err := h.svc.ListTalentMarketplaceOpportunities(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list talent marketplace opportunities")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createTalentMarketplaceOpportunityForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.TalentMarketplaceOpportunityCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateTalentMarketplaceOpportunity(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateTalentMarketplaceOpportunityForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.TalentMarketplaceOpportunityCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateTalentMarketplaceOpportunity(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateTalentMarketplaceFallbackForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.TalentMarketplaceFallbackCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ActivateTalentMarketplaceFallback(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listTalentMarketplaceApplicationsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter, ok := h.talentMarketplaceApplicationFilter(w, r, tenantID, operation)
	if !ok {
		return
	}
	items, err := h.svc.ListTalentMarketplaceApplications(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list talent marketplace applications")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createTalentMarketplaceApplicationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.TalentMarketplaceApplicationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateTalentMarketplaceApplication(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateTalentMarketplaceApplicationStatusForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.TalentMarketplaceApplicationStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateTalentMarketplaceApplicationStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listTalentMarketplaceEventsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	filter, ok := h.talentMarketplaceEventFilter(w, r, tenantID, operation)
	if !ok {
		return
	}
	items, err := h.svc.ListTalentMarketplaceEvents(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list talent marketplace events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) talentMarketplaceOpportunityFilter(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (domain.TalentMarketplaceOpportunityFilter, bool) {
	projectID, ok := h.optionalUUIDQuery(w, r, "project_id", operation)
	if !ok {
		return domain.TalentMarketplaceOpportunityFilter{}, false
	}
	engagementID, ok := h.optionalUUIDQuery(w, r, "engagement_id", operation)
	if !ok {
		return domain.TalentMarketplaceOpportunityFilter{}, false
	}
	return domain.TalentMarketplaceOpportunityFilter{
		TenantID:        tenantID,
		ProjectID:       projectID,
		EngagementID:    engagementID,
		Status:          optionalStringQuery(r, "status"),
		OpportunityType: optionalStringQuery(r, "opportunity_type"),
		Search:          optionalStringQuery(r, "search"),
	}, true
}

func (h *Handler) talentMarketplaceApplicationFilter(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (domain.TalentMarketplaceApplicationFilter, bool) {
	opportunityID, ok := h.optionalUUIDQuery(w, r, "opportunity_id", operation)
	if !ok {
		return domain.TalentMarketplaceApplicationFilter{}, false
	}
	workerProfileID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return domain.TalentMarketplaceApplicationFilter{}, false
	}
	return domain.TalentMarketplaceApplicationFilter{
		TenantID:        tenantID,
		OpportunityID:   opportunityID,
		WorkerProfileID: workerProfileID,
		Status:          optionalStringQuery(r, "status"),
	}, true
}

func (h *Handler) talentMarketplaceEventFilter(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (domain.TalentMarketplaceEventFilter, bool) {
	opportunityID, ok := h.optionalUUIDQuery(w, r, "opportunity_id", operation)
	if !ok {
		return domain.TalentMarketplaceEventFilter{}, false
	}
	applicationID, ok := h.optionalUUIDQuery(w, r, "application_id", operation)
	if !ok {
		return domain.TalentMarketplaceEventFilter{}, false
	}
	return domain.TalentMarketplaceEventFilter{TenantID: tenantID, OpportunityID: opportunityID, ApplicationID: applicationID}, true
}

func (h *Handler) talentMarketplaceOpportunityRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "opportunityID", operation, "invalid opportunity id")
	return tenantID, id, ok
}

func (h *Handler) talentMarketplaceApplicationRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "applicationID", operation, "invalid application id")
	return tenantID, id, ok
}

func (h *Handler) superAdminTalentMarketplaceOpportunityRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "opportunityID", operation, "invalid opportunity id")
	return tenantID, id, ok
}

func (h *Handler) superAdminTalentMarketplaceApplicationRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, ok := h.uuidURLParam(w, r, "applicationID", operation, "invalid application id")
	return tenantID, id, ok
}
