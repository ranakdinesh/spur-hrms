package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateHoliday(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create holiday", err, "tenant context is required")
		return
	}
	h.createHolidayForTenant(w, r, tenantID, "create holiday")
}

func (h *Handler) ListHolidays(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list holidays", err, "tenant context is required")
		return
	}
	h.listHolidaysForTenant(w, r, tenantID, "list holidays")
}

func (h *Handler) GetHoliday(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.holidayRequestIDs(w, r, "get holiday")
	if !ok {
		return
	}
	item, err := h.svc.GetHoliday(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get holiday", err, "holiday not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateHoliday(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.holidayRequestIDs(w, r, "update holiday")
	if !ok {
		return
	}
	h.updateHolidayForTenant(w, r, tenantID, id, "update holiday")
}

func (h *Handler) DeleteHoliday(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.holidayRequestIDs(w, r, "delete holiday")
	if !ok {
		return
	}
	if err := h.svc.DeleteHoliday(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete holiday", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantHoliday(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant holiday")
	if !ok {
		return
	}
	h.createHolidayForTenant(w, r, tenantID, "create tenant holiday")
}

func (h *Handler) ListTenantHolidays(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant holidays")
	if !ok {
		return
	}
	h.listHolidaysForTenant(w, r, tenantID, "list tenant holidays")
}

func (h *Handler) GetTenantHoliday(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminHolidayRequestIDs(w, r, "get tenant holiday")
	if !ok {
		return
	}
	item, err := h.svc.GetHoliday(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant holiday", err, "holiday not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantHoliday(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminHolidayRequestIDs(w, r, "update tenant holiday")
	if !ok {
		return
	}
	h.updateHolidayForTenant(w, r, tenantID, id, "update tenant holiday")
}

func (h *Handler) DeleteTenantHoliday(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminHolidayRequestIDs(w, r, "delete tenant holiday")
	if !ok {
		return
	}
	if err := h.svc.DeleteHoliday(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant holiday", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createHolidayForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.HolidayCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateHoliday(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateHolidayForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.HolidayCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateHoliday(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listHolidaysForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	cmd, ok := h.listHolidaysCommand(w, r, tenantID, operation)
	if !ok {
		return
	}
	items, err := h.svc.ListHolidays(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) listHolidaysCommand(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) (ports.ListHolidaysCommand, bool) {
	branchID, ok := h.optionalUUIDQuery(w, r, "branch_id", operation)
	if !ok {
		return ports.ListHolidaysCommand{}, false
	}
	fyID, ok := h.optionalUUIDQuery(w, r, "fy_id", operation)
	if !ok {
		return ports.ListHolidaysCommand{}, false
	}
	limit := int32(0)
	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid limit")
			return ports.ListHolidaysCommand{}, false
		}
		limit = int32(parsed)
	}
	return ports.ListHolidaysCommand{
		TenantID:  tenantID,
		BranchID:  branchID,
		FYID:      fyID,
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
		Upcoming:  r.URL.Query().Get("upcoming") == "true",
		Limit:     limit,
	}, true
}

func (h *Handler) holidayRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "holidayID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid holiday id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminHolidayRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "holidayID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid holiday id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
