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

func (h *Handler) ListTenantOperationRequests(w http.ResponseWriter, r *http.Request) {
	if !h.isSuperAdminRequest(r) {
		h.respondError(w, r, http.StatusForbidden, "list tenant operation requests", errors.New("super admin required"), "super admin access is required")
		return
	}
	result, err := h.svc.ListTenantOperationRequests(r.Context(), domain.TenantOperationFilter{
		Status:         optionalStringQuery(r, "status"),
		OperationType:  optionalStringQuery(r, "operation_type"),
		RiskLevel:      optionalStringQuery(r, "risk_level"),
		TargetTenantID: optionalUUIDQuery(r, "target_tenant_id"),
		Search:         optionalStringQuery(r, "search"),
		Limit:          queryInt32(r, "limit", 100),
		Offset:         queryInt32(r, "offset", 0),
	})
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "list tenant operation requests", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) CreateTenantOperationRequest(w http.ResponseWriter, r *http.Request) {
	if !h.isSuperAdminRequest(r) {
		h.respondError(w, r, http.StatusForbidden, "create tenant operation request", errors.New("super admin required"), "super admin access is required")
		return
	}
	var cmd ports.TenantOperationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant operation request", err, "invalid request body")
		return
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateTenantOperationRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create tenant operation request", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) GetTenantOperationRequest(w http.ResponseWriter, r *http.Request) {
	if !h.isSuperAdminRequest(r) {
		h.respondError(w, r, http.StatusForbidden, "get tenant operation request", errors.New("super admin required"), "super admin access is required")
		return
	}
	id, ok := h.tenantOperationRequestID(w, r, "get tenant operation request")
	if !ok {
		return
	}
	item, err := h.svc.GetTenantOperationDetail(r.Context(), id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant operation request", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ActTenantOperationRequest(w http.ResponseWriter, r *http.Request) {
	if !h.isSuperAdminRequest(r) {
		h.respondError(w, r, http.StatusForbidden, "act tenant operation request", errors.New("super admin required"), "super admin access is required")
		return
	}
	id, ok := h.tenantOperationRequestID(w, r, "act tenant operation request")
	if !ok {
		return
	}
	var cmd ports.TenantOperationActionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant operation action", err, "invalid request body")
		return
	}
	cmd.ID, cmd.ActorID = id, h.actorIDFromRequest(r)
	item, err := h.svc.ActTenantOperationRequest(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "act tenant operation request", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) tenantOperationRequestID(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "tenantOperationID"))
	if err != nil || id == uuid.Nil {
		if err == nil {
			err = domain.ErrInvalidTenantOperation
		}
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid tenant operation id")
		return uuid.Nil, false
	}
	return id, true
}
