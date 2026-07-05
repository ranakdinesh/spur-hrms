package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListDesignationLevelCodes(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list designation level codes", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListDesignationLevelCodes(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list designation level codes", err, "failed to list designation level codes")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateDesignationLevelCode(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create designation level code", err, "tenant context is required")
		return
	}
	var cmd ports.DesignationLevelCodeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode designation level code create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateDesignationLevelCode(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create designation level code", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateDesignationLevelCode(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.designationMasterRequestIDs(w, r, "levelCodeID", "update designation level code")
	if !ok {
		return
	}
	var cmd ports.DesignationLevelCodeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode designation level code update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateDesignationLevelCode(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update designation level code", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteDesignationLevelCode(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.designationMasterRequestIDs(w, r, "levelCodeID", "delete designation level code")
	if !ok {
		return
	}
	if err := h.svc.DeleteDesignationLevelCode(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete designation level code", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantDesignationLevelCodes(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant designation level codes")
	if !ok {
		return
	}
	items, err := h.svc.ListDesignationLevelCodes(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant designation level codes", err, "failed to list designation level codes")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateTenantDesignationLevelCode(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant designation level code")
	if !ok {
		return
	}
	var cmd ports.DesignationLevelCodeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant designation level code create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateDesignationLevelCode(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant designation level code", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateTenantDesignationLevelCode(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminDesignationMasterRequestIDs(w, r, "levelCodeID", "update tenant designation level code")
	if !ok {
		return
	}
	var cmd ports.DesignationLevelCodeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant designation level code update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateDesignationLevelCode(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant designation level code", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantDesignationLevelCode(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminDesignationMasterRequestIDs(w, r, "levelCodeID", "delete tenant designation level code")
	if !ok {
		return
	}
	if err := h.svc.DeleteDesignationLevelCode(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant designation level code", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListDesignationSeniorityRanks(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list designation seniority ranks", err, "tenant context is required")
		return
	}
	items, err := h.svc.ListDesignationSeniorityRanks(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list designation seniority ranks", err, "failed to list designation seniority ranks")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateDesignationSeniorityRank(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create designation seniority rank", err, "tenant context is required")
		return
	}
	var cmd ports.DesignationSeniorityRankCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode designation seniority rank create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateDesignationSeniorityRank(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create designation seniority rank", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateDesignationSeniorityRank(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.designationMasterRequestIDs(w, r, "rankID", "update designation seniority rank")
	if !ok {
		return
	}
	var cmd ports.DesignationSeniorityRankCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode designation seniority rank update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateDesignationSeniorityRank(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update designation seniority rank", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteDesignationSeniorityRank(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.designationMasterRequestIDs(w, r, "rankID", "delete designation seniority rank")
	if !ok {
		return
	}
	if err := h.svc.DeleteDesignationSeniorityRank(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete designation seniority rank", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantDesignationSeniorityRanks(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant designation seniority ranks")
	if !ok {
		return
	}
	items, err := h.svc.ListDesignationSeniorityRanks(r.Context(), tenantID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant designation seniority ranks", err, "failed to list designation seniority ranks")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateTenantDesignationSeniorityRank(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant designation seniority rank")
	if !ok {
		return
	}
	var cmd ports.DesignationSeniorityRankCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant designation seniority rank create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateDesignationSeniorityRank(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant designation seniority rank", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) UpdateTenantDesignationSeniorityRank(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminDesignationMasterRequestIDs(w, r, "rankID", "update tenant designation seniority rank")
	if !ok {
		return
	}
	var cmd ports.DesignationSeniorityRankCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant designation seniority rank update request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateDesignationSeniorityRank(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant designation seniority rank", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteTenantDesignationSeniorityRank(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminDesignationMasterRequestIDs(w, r, "rankID", "delete tenant designation seniority rank")
	if !ok {
		return
	}
	if err := h.svc.DeleteDesignationSeniorityRank(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant designation seniority rank", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) designationMasterRequestIDs(w http.ResponseWriter, r *http.Request, param string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid designation master id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminDesignationMasterRequestIDs(w http.ResponseWriter, r *http.Request, param string, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid designation master id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
