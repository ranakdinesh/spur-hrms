package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create branch", err, "tenant context is required")
		return
	}
	var cmd ports.BranchCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode branch create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	branch, err := h.svc.CreateBranch(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create branch", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, branch)
}

func (h *Handler) ListBranches(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list branches", err, "tenant context is required")
		return
	}
	branches, err := h.svc.ListBranches(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list branches", err, "failed to list branches")
		return
	}
	respondJSON(w, http.StatusOK, branches)
}

func (h *Handler) GetBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, branchID, ok := h.branchRequestIDs(w, r, "get branch")
	if !ok {
		return
	}
	branch, err := h.svc.GetBranch(r.Context(), tenantID, branchID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get branch", err, "branch not found")
		return
	}
	respondJSON(w, http.StatusOK, branch)
}

func (h *Handler) UpdateBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, branchID, ok := h.branchRequestIDs(w, r, "update branch")
	if !ok {
		return
	}
	var cmd ports.BranchCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode branch update request", err, "invalid request body")
		return
	}
	cmd.ID = branchID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	branch, err := h.svc.UpdateBranch(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update branch", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, branch)
}

func (h *Handler) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, branchID, ok := h.branchRequestIDs(w, r, "delete branch")
	if !ok {
		return
	}
	if err := h.svc.DeleteBranch(r.Context(), tenantID, branchID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete branch", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant branch")
	if !ok {
		return
	}
	var cmd ports.BranchCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant branch create request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	branch, err := h.svc.CreateBranch(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant branch", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, branch)
}

func (h *Handler) ListTenantBranches(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant branches")
	if !ok {
		return
	}
	branches, err := h.svc.ListBranches(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant branches", err, "failed to list branches")
		return
	}
	respondJSON(w, http.StatusOK, branches)
}

func (h *Handler) GetTenantBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, branchID, ok := h.superAdminBranchRequestIDs(w, r, "get tenant branch")
	if !ok {
		return
	}
	branch, err := h.svc.GetBranch(r.Context(), tenantID, branchID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant branch", err, "branch not found")
		return
	}
	respondJSON(w, http.StatusOK, branch)
}

func (h *Handler) UpdateTenantBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, branchID, ok := h.superAdminBranchRequestIDs(w, r, "update tenant branch")
	if !ok {
		return
	}
	var cmd ports.BranchCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant branch update request", err, "invalid request body")
		return
	}
	cmd.ID = branchID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	branch, err := h.svc.UpdateBranch(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update tenant branch", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, branch)
}

func (h *Handler) DeleteTenantBranch(w http.ResponseWriter, r *http.Request) {
	tenantID, branchID, ok := h.superAdminBranchRequestIDs(w, r, "delete tenant branch")
	if !ok {
		return
	}
	if err := h.svc.DeleteBranch(r.Context(), tenantID, branchID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant branch", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) branchRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	branchID, err := uuid.Parse(chi.URLParam(r, "branchID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid branch id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, branchID, true
}

func (h *Handler) superAdminTenantID(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	if !h.isSuperAdminRequest(r) {
		h.respondError(w, r, http.StatusForbidden, operation, nil, "super admin permission required")
		return uuid.Nil, false
	}
	tenantID, err := uuid.Parse(chi.URLParam(r, "tenantID"))
	if err != nil || tenantID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid tenant id")
		return uuid.Nil, false
	}
	return tenantID, true
}

func (h *Handler) superAdminBranchRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	branchID, err := uuid.Parse(chi.URLParam(r, "branchID"))
	if err != nil || branchID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid branch id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, branchID, true
}
