package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListCompensationPayBands(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listCompensationPayBandsForTenant(w, r, tenantID, "list compensation pay bands")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list compensation pay bands", err, "tenant context is required")
	}
}

func (h *Handler) CreateCompensationPayBand(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createCompensationPayBandForTenant(w, r, tenantID, "create compensation pay band")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create compensation pay band", err, "tenant context is required")
	}
}

func (h *Handler) UpdateCompensationPayBand(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.compTenantAndID(w, r, "payBandID", "update compensation pay band"); ok {
		h.updateCompensationPayBandForTenant(w, r, tenantID, id, "update compensation pay band")
	}
}

func (h *Handler) DeleteCompensationPayBand(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.compTenantAndID(w, r, "payBandID", "delete compensation pay band"); ok {
		if err := h.svc.DeleteCompensationPayBand(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete compensation pay band", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) ListCompensationCycles(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listCompensationCyclesForTenant(w, r, tenantID, "list compensation cycles")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list compensation cycles", err, "tenant context is required")
	}
}

func (h *Handler) CreateCompensationCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createCompensationCycleForTenant(w, r, tenantID, "create compensation cycle")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create compensation cycle", err, "tenant context is required")
	}
}

func (h *Handler) UpdateCompensationCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.compTenantAndID(w, r, "cycleID", "update compensation cycle"); ok {
		h.updateCompensationCycleForTenant(w, r, tenantID, id, "update compensation cycle")
	}
}

func (h *Handler) UpdateCompensationCycleStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.compTenantAndID(w, r, "cycleID", "update compensation cycle status"); ok {
		var cmd ports.CompensationStatusCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode compensation cycle status", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		item, err := h.svc.UpdateCompensationCycleStatus(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "update compensation cycle status", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
	}
}

func (h *Handler) ListCompensationBudgetPools(w http.ResponseWriter, r *http.Request) {
	if tenantID, cycleID, ok := h.compTenantAndID(w, r, "cycleID", "list compensation budget pools"); ok {
		items, err := h.svc.ListCompensationBudgetPools(r.Context(), tenantID, cycleID)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, "list compensation budget pools", err, "failed to list budget pools")
			return
		}
		respondJSON(w, http.StatusOK, items)
	}
}

func (h *Handler) CreateCompensationBudgetPool(w http.ResponseWriter, r *http.Request) {
	if tenantID, cycleID, ok := h.compTenantAndID(w, r, "cycleID", "create compensation budget pool"); ok {
		var cmd ports.CompensationBudgetPoolCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode compensation budget pool", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.CycleID, cmd.ActorID = tenantID, cycleID, h.actorIDFromRequest(r)
		item, err := h.svc.CreateCompensationBudgetPool(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "create compensation budget pool", err, err.Error())
			return
		}
		respondJSON(w, http.StatusCreated, item)
	}
}

func (h *Handler) UpdateCompensationBudgetPool(w http.ResponseWriter, r *http.Request) {
	tenantID, cycleID, ok := h.compTenantAndID(w, r, "cycleID", "update compensation budget pool")
	if !ok {
		return
	}
	poolID, err := uuid.Parse(chi.URLParam(r, "poolID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update compensation budget pool", err, "invalid budget pool id")
		return
	}
	var cmd ports.CompensationBudgetPoolCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode compensation budget pool", err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.CycleID, cmd.ID, cmd.ActorID = tenantID, cycleID, poolID, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCompensationBudgetPool(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update compensation budget pool", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteCompensationBudgetPool(w http.ResponseWriter, r *http.Request) {
	tenantID, cycleID, ok := h.compTenantAndID(w, r, "cycleID", "delete compensation budget pool")
	if !ok {
		return
	}
	poolID, err := uuid.Parse(chi.URLParam(r, "poolID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete compensation budget pool", err, "invalid budget pool id")
		return
	}
	if err := h.svc.DeleteCompensationBudgetPool(r.Context(), tenantID, cycleID, poolID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete compensation budget pool", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListCompensationRecommendations(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listCompensationRecommendationsForTenant(w, r, tenantID, "list compensation recommendations")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list compensation recommendations", err, "tenant context is required")
	}
}

func (h *Handler) CreateCompensationRecommendation(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.createCompensationRecommendationForTenant(w, r, tenantID, "create compensation recommendation")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "create compensation recommendation", err, "tenant context is required")
	}
}

func (h *Handler) UpdateCompensationRecommendation(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.compTenantAndID(w, r, "recommendationID", "update compensation recommendation"); ok {
		h.updateCompensationRecommendationForTenant(w, r, tenantID, id, "update compensation recommendation")
	}
}

func (h *Handler) UpdateCompensationRecommendationStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.compTenantAndID(w, r, "recommendationID", "update compensation recommendation status"); ok {
		var cmd ports.CompensationStatusCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode compensation recommendation status", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		item, err := h.svc.UpdateCompensationRecommendationStatus(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "update compensation recommendation status", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
	}
}

func (h *Handler) DeleteCompensationRecommendation(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.compTenantAndID(w, r, "recommendationID", "delete compensation recommendation"); ok {
		if err := h.svc.DeleteCompensationRecommendation(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "delete compensation recommendation", err, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) GenerateCompensationEquityChecks(w http.ResponseWriter, r *http.Request) {
	if tenantID, cycleID, ok := h.compTenantAndID(w, r, "cycleID", "generate compensation equity checks"); ok {
		items, err := h.svc.GenerateCompensationEquityChecks(r.Context(), tenantID, cycleID, h.actorIDFromRequest(r))
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "generate compensation equity checks", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, items)
	}
}

func (h *Handler) ListCompensationEquityChecks(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listCompensationEquityChecksForTenant(w, r, tenantID, "list compensation equity checks")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list compensation equity checks", err, "tenant context is required")
	}
}

func (h *Handler) UpdateCompensationEquityCheckStatus(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.compTenantAndID(w, r, "equityCheckID", "update compensation equity check status"); ok {
		var cmd ports.CompensationStatusCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			h.respondError(w, r, http.StatusBadRequest, "decode compensation equity check status", err, "invalid request body")
			return
		}
		cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
		item, err := h.svc.UpdateCompensationEquityCheckStatus(r.Context(), cmd)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, "update compensation equity check status", err, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, item)
	}
}

func (h *Handler) ListCompensationEvents(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.listCompensationEventsForTenant(w, r, tenantID, "list compensation events")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "list compensation events", err, "tenant context is required")
	}
}

func (h *Handler) GetCompensationSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, err := h.tenantIDFromRequest(r); err == nil {
		h.getCompensationSummaryForTenant(w, r, tenantID, "get compensation summary")
	} else {
		h.respondError(w, r, http.StatusUnauthorized, "get compensation summary", err, "tenant context is required")
	}
}

func (h *Handler) ListTenantCompensationPayBands(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant compensation pay bands"); ok {
		h.listCompensationPayBandsForTenant(w, r, tenantID, "list tenant compensation pay bands")
	}
}

func (h *Handler) CreateTenantCompensationPayBand(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant compensation pay band"); ok {
		h.createCompensationPayBandForTenant(w, r, tenantID, "create tenant compensation pay band")
	}
}

func (h *Handler) UpdateTenantCompensationPayBand(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.tenantCompID(w, r, "payBandID", "update tenant compensation pay band"); ok {
		h.updateCompensationPayBandForTenant(w, r, tenantID, id, "update tenant compensation pay band")
	}
}

func (h *Handler) DeleteTenantCompensationPayBand(w http.ResponseWriter, r *http.Request) {
	h.DeleteCompensationPayBand(w, r)
}
func (h *Handler) ListTenantCompensationCycles(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant compensation cycles"); ok {
		h.listCompensationCyclesForTenant(w, r, tenantID, "list tenant compensation cycles")
	}
}
func (h *Handler) CreateTenantCompensationCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant compensation cycle"); ok {
		h.createCompensationCycleForTenant(w, r, tenantID, "create tenant compensation cycle")
	}
}
func (h *Handler) UpdateTenantCompensationCycle(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.tenantCompID(w, r, "cycleID", "update tenant compensation cycle"); ok {
		h.updateCompensationCycleForTenant(w, r, tenantID, id, "update tenant compensation cycle")
	}
}
func (h *Handler) UpdateTenantCompensationCycleStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateCompensationCycleStatus(w, r)
}
func (h *Handler) ListTenantCompensationBudgetPools(w http.ResponseWriter, r *http.Request) {
	h.ListCompensationBudgetPools(w, r)
}
func (h *Handler) CreateTenantCompensationBudgetPool(w http.ResponseWriter, r *http.Request) {
	h.CreateCompensationBudgetPool(w, r)
}
func (h *Handler) UpdateTenantCompensationBudgetPool(w http.ResponseWriter, r *http.Request) {
	h.UpdateCompensationBudgetPool(w, r)
}
func (h *Handler) DeleteTenantCompensationBudgetPool(w http.ResponseWriter, r *http.Request) {
	h.DeleteCompensationBudgetPool(w, r)
}
func (h *Handler) ListTenantCompensationRecommendations(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant compensation recommendations"); ok {
		h.listCompensationRecommendationsForTenant(w, r, tenantID, "list tenant compensation recommendations")
	}
}
func (h *Handler) CreateTenantCompensationRecommendation(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant compensation recommendation"); ok {
		h.createCompensationRecommendationForTenant(w, r, tenantID, "create tenant compensation recommendation")
	}
}
func (h *Handler) UpdateTenantCompensationRecommendation(w http.ResponseWriter, r *http.Request) {
	if tenantID, id, ok := h.tenantCompID(w, r, "recommendationID", "update tenant compensation recommendation"); ok {
		h.updateCompensationRecommendationForTenant(w, r, tenantID, id, "update tenant compensation recommendation")
	}
}
func (h *Handler) UpdateTenantCompensationRecommendationStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateCompensationRecommendationStatus(w, r)
}
func (h *Handler) DeleteTenantCompensationRecommendation(w http.ResponseWriter, r *http.Request) {
	h.DeleteCompensationRecommendation(w, r)
}
func (h *Handler) GenerateTenantCompensationEquityChecks(w http.ResponseWriter, r *http.Request) {
	h.GenerateCompensationEquityChecks(w, r)
}
func (h *Handler) ListTenantCompensationEquityChecks(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant compensation equity checks"); ok {
		h.listCompensationEquityChecksForTenant(w, r, tenantID, "list tenant compensation equity checks")
	}
}
func (h *Handler) UpdateTenantCompensationEquityCheckStatus(w http.ResponseWriter, r *http.Request) {
	h.UpdateCompensationEquityCheckStatus(w, r)
}
func (h *Handler) ListTenantCompensationEvents(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant compensation events"); ok {
		h.listCompensationEventsForTenant(w, r, tenantID, "list tenant compensation events")
	}
}
func (h *Handler) GetTenantCompensationSummary(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "get tenant compensation summary"); ok {
		h.getCompensationSummaryForTenant(w, r, tenantID, "get tenant compensation summary")
	}
}

func (h *Handler) listCompensationPayBandsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListCompensationPayBands(r.Context(), domain.CompensationPayBandFilter{TenantID: tenantID, IsActive: optionalBoolQuery(r, "is_active"), CurrencyCode: optionalStringQuery(r, "currency_code"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 50), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list pay bands")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createCompensationPayBandForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CompensationPayBandCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateCompensationPayBand(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateCompensationPayBandForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.CompensationPayBandCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCompensationPayBand(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listCompensationCyclesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListCompensationCycles(r.Context(), compFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list cycles")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createCompensationCycleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CompensationCycleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateCompensationCycle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateCompensationCycleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.CompensationCycleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCompensationCycle(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listCompensationRecommendationsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListCompensationRecommendations(r.Context(), compFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list recommendations")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createCompensationRecommendationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.CompensationRecommendationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ActorID = tenantID, h.actorIDFromRequest(r)
	item, err := h.svc.CreateCompensationRecommendation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateCompensationRecommendationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.CompensationRecommendationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID, cmd.ID, cmd.ActorID = tenantID, id, h.actorIDFromRequest(r)
	item, err := h.svc.UpdateCompensationRecommendation(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listCompensationEquityChecksForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListCompensationEquityChecks(r.Context(), compFilterFromRequest(r, tenantID))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list equity checks")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listCompensationEventsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListCompensationEvents(r.Context(), compFilterFromRequest(r, tenantID), optionalStringQuery(r, "source_type"), optionalUUIDQuery(r, "source_id"))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list compensation events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getCompensationSummaryForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.GetCompensationSummary(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to get compensation summary")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func compFilterFromRequest(r *http.Request, tenantID uuid.UUID) domain.CompensationFilter {
	return domain.CompensationFilter{TenantID: tenantID, CycleID: optionalUUIDQuery(r, "cycle_id"), Status: optionalStringQuery(r, "status"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 50), Offset: queryInt32(r, "offset", 0)}
}

func (h *Handler) compTenantAndID(w http.ResponseWriter, r *http.Request, idParam string, operation string) (uuid.UUID, uuid.UUID, bool) {
	if tenantIDValue := chi.URLParam(r, "tenantID"); tenantIDValue != "" {
		return h.tenantCompID(w, r, idParam, operation)
	}
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, idParam))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) tenantCompID(w http.ResponseWriter, r *http.Request, idParam string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, idParam))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
