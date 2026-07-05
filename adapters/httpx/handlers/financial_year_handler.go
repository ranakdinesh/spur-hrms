package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create financial year", err, "tenant context is required")
		return
	}
	var cmd ports.FinancialYearCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode financial year create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateFinancialYear(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create financial year", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListFinancialYears(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list financial years", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListFinancialYears(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list financial years", err, "failed to list financial years")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetActiveFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get active financial year", err, "tenant context is required")
		return
	}
	item, err := h.svc.GetActiveFinancialYear(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get active financial year", err, "active financial year not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) GetFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.financialYearRequestIDs(w, r, "get financial year")
	if !ok {
		return
	}
	item, err := h.svc.GetFinancialYear(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get financial year", err, "financial year not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.financialYearRequestIDs(w, r, "update financial year")
	if !ok {
		return
	}
	var cmd ports.FinancialYearCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode financial year update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateFinancialYear(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update financial year", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.financialYearRequestIDs(w, r, "delete financial year")
	if !ok {
		return
	}
	if err := h.svc.DeleteFinancialYear(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete financial year", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SetActiveFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.financialYearRequestIDs(w, r, "set active financial year")
	if !ok {
		return
	}
	item, err := h.svc.SetActiveFinancialYear(r.Context(), tenantID, id, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "set active financial year", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) SetFinancialYearLock(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.financialYearRequestIDs(w, r, "set financial year lock")
	if !ok {
		return
	}
	var cmd ports.FinancialYearLockCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode financial year lock request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.SetFinancialYearLock(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "set financial year lock", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateTenantFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant financial year")
	if !ok {
		return
	}
	var cmd ports.FinancialYearCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant financial year create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateFinancialYear(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant financial year", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListTenantFinancialYears(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant financial years")
	if !ok {
		return
	}
	items, err := h.svc.ListFinancialYears(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant financial years", err, "failed to list financial years")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetTenantActiveFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant active financial year")
	if !ok {
		return
	}
	item, err := h.svc.GetActiveFinancialYear(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant active financial year", err, "active financial year not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) GetTenantFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFinancialYearRequestIDs(w, r, "get tenant financial year")
	if !ok {
		return
	}
	item, err := h.svc.GetFinancialYear(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant financial year", err, "financial year not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFinancialYearRequestIDs(w, r, "update tenant financial year")
	if !ok {
		return
	}
	var cmd ports.FinancialYearCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant financial year update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateFinancialYear(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant financial year", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFinancialYearRequestIDs(w, r, "delete tenant financial year")
	if !ok {
		return
	}
	if err := h.svc.DeleteFinancialYear(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant financial year", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SetTenantActiveFinancialYear(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFinancialYearRequestIDs(w, r, "set tenant active financial year")
	if !ok {
		return
	}
	item, err := h.svc.SetActiveFinancialYear(r.Context(), tenantID, id, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "set tenant active financial year", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) SetTenantFinancialYearLock(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFinancialYearRequestIDs(w, r, "set tenant financial year lock")
	if !ok {
		return
	}
	var cmd ports.FinancialYearLockCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant financial year lock request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.SetFinancialYearLock(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "set tenant financial year lock", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) financialYearRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "financialYearID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid financial year id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminFinancialYearRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "financialYearID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid financial year id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
