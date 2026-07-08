package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
)

func (h *Handler) ListPayGroups(w http.ResponseWriter, r *http.Request) {
	if !h.requirePermission(w, r, "list pay groups", permissions.PayGroupsList) {
		return
	}
	tenantID, ok := h.currentTenantID(w, r, "list pay groups")
	if ok {
		h.listPayGroupsForTenant(w, r, tenantID, "list pay groups")
	}
}

func (h *Handler) CreatePayGroup(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.currentTenantID(w, r, "create pay group")
	if ok {
		h.createPayGroupForTenant(w, r, tenantID, "create pay group")
	}
}

func (h *Handler) GetPayGroup(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.payGroupRequestIDs(w, r, "get pay group")
	if ok {
		h.getPayGroupForTenant(w, r, tenantID, payGroupID, "get pay group")
	}
}

func (h *Handler) UpdatePayGroup(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.payGroupRequestIDs(w, r, "update pay group")
	if ok {
		h.updatePayGroupForTenant(w, r, tenantID, payGroupID, "update pay group")
	}
}

func (h *Handler) DeletePayGroup(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.payGroupRequestIDs(w, r, "delete pay group")
	if !ok {
		return
	}
	if err := h.svc.DeletePayGroup(r.Context(), tenantID, payGroupID, h.actorIDFromRequest(r)); err != nil {
		h.respondPayGroupError(w, r, "delete pay group", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListPayGroupEmployees(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.payGroupRequestIDs(w, r, "list pay group employees")
	if ok {
		h.listPayGroupEmployeesForTenant(w, r, tenantID, payGroupID, "list pay group employees")
	}
}

func (h *Handler) UpsertPayGroupMember(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.payGroupRequestIDs(w, r, "upsert pay group member")
	if ok {
		h.upsertPayGroupMemberForTenant(w, r, tenantID, payGroupID, "upsert pay group member")
	}
}

func (h *Handler) DeletePayGroupMember(w http.ResponseWriter, r *http.Request) {
	tenantID, memberID, ok := h.payGroupMemberRequestIDs(w, r, "delete pay group member")
	if !ok {
		return
	}
	if err := h.svc.DeletePayGroupMember(r.Context(), tenantID, memberID, h.actorIDFromRequest(r)); err != nil {
		h.respondPayGroupError(w, r, "delete pay group member", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListPayRuns(w http.ResponseWriter, r *http.Request) {
	if !h.requirePermission(w, r, "list pay runs", permissions.PayRunsList) {
		return
	}
	tenantID, ok := h.currentTenantID(w, r, "list pay runs")
	if ok {
		h.listPayRunsForTenant(w, r, tenantID, "list pay runs")
	}
}

func (h *Handler) CreatePayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.currentTenantID(w, r, "create pay run")
	if ok {
		h.createPayRunForTenant(w, r, tenantID, "create pay run")
	}
}

func (h *Handler) GetPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, payRunID, ok := h.payRunRequestIDs(w, r, "get pay run")
	if ok {
		h.getPayRunForTenant(w, r, tenantID, payRunID, "get pay run")
	}
}

func (h *Handler) AssessPayRunReadiness(w http.ResponseWriter, r *http.Request) {
	tenantID, payRunID, ok := h.payRunRequestIDs(w, r, "assess pay run")
	if !ok {
		return
	}
	item, err := h.svc.AssessPayRunReadiness(r.Context(), tenantID, payRunID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondPayGroupError(w, r, "assess pay run", err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) FreezePayRun(w http.ResponseWriter, r *http.Request) {
	h.payRunAction(w, r, false, h.svc.FreezePayRun, "freeze pay run")
}

func (h *Handler) GeneratePayRun(w http.ResponseWriter, r *http.Request) {
	h.payRunAction(w, r, false, h.svc.GeneratePayRun, "generate pay run")
}

func (h *Handler) LockPayRun(w http.ResponseWriter, r *http.Request) {
	h.payRunAction(w, r, false, h.svc.LockPayRun, "lock pay run")
}

func (h *Handler) UnlockPayRun(w http.ResponseWriter, r *http.Request) {
	h.payRunAction(w, r, false, h.svc.UnlockPayRun, "unlock pay run")
}

func (h *Handler) DeletePayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, payRunID, ok := h.payRunRequestIDs(w, r, "delete pay run")
	if !ok {
		return
	}
	if err := h.svc.DeletePayRun(r.Context(), tenantID, payRunID, h.actorIDFromRequest(r)); err != nil {
		h.respondPayGroupError(w, r, "delete pay run", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantPayGroups(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant pay groups")
	if ok {
		h.listPayGroupsForTenant(w, r, tenantID, "list tenant pay groups")
	}
}

func (h *Handler) CreateTenantPayGroup(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant pay group")
	if ok {
		h.createPayGroupForTenant(w, r, tenantID, "create tenant pay group")
	}
}

func (h *Handler) GetTenantPayGroup(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.tenantPayGroupRequestIDs(w, r, "get tenant pay group")
	if ok {
		h.getPayGroupForTenant(w, r, tenantID, payGroupID, "get tenant pay group")
	}
}

func (h *Handler) UpdateTenantPayGroup(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.tenantPayGroupRequestIDs(w, r, "update tenant pay group")
	if ok {
		h.updatePayGroupForTenant(w, r, tenantID, payGroupID, "update tenant pay group")
	}
}

func (h *Handler) DeleteTenantPayGroup(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.tenantPayGroupRequestIDs(w, r, "delete tenant pay group")
	if !ok {
		return
	}
	if err := h.svc.DeletePayGroup(r.Context(), tenantID, payGroupID, h.actorIDFromRequest(r)); err != nil {
		h.respondPayGroupError(w, r, "delete tenant pay group", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantPayGroupEmployees(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.tenantPayGroupRequestIDs(w, r, "list tenant pay group employees")
	if ok {
		h.listPayGroupEmployeesForTenant(w, r, tenantID, payGroupID, "list tenant pay group employees")
	}
}

func (h *Handler) UpsertTenantPayGroupMember(w http.ResponseWriter, r *http.Request) {
	tenantID, payGroupID, ok := h.tenantPayGroupRequestIDs(w, r, "upsert tenant pay group member")
	if ok {
		h.upsertPayGroupMemberForTenant(w, r, tenantID, payGroupID, "upsert tenant pay group member")
	}
}

func (h *Handler) DeleteTenantPayGroupMember(w http.ResponseWriter, r *http.Request) {
	tenantID, memberID, ok := h.tenantPayGroupMemberRequestIDs(w, r, "delete tenant pay group member")
	if !ok {
		return
	}
	if err := h.svc.DeletePayGroupMember(r.Context(), tenantID, memberID, h.actorIDFromRequest(r)); err != nil {
		h.respondPayGroupError(w, r, "delete tenant pay group member", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantPayRuns(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant pay runs")
	if ok {
		h.listPayRunsForTenant(w, r, tenantID, "list tenant pay runs")
	}
}

func (h *Handler) CreateTenantPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant pay run")
	if ok {
		h.createPayRunForTenant(w, r, tenantID, "create tenant pay run")
	}
}

func (h *Handler) GetTenantPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, payRunID, ok := h.tenantPayRunRequestIDs(w, r, "get tenant pay run")
	if ok {
		h.getPayRunForTenant(w, r, tenantID, payRunID, "get tenant pay run")
	}
}

func (h *Handler) AssessTenantPayRunReadiness(w http.ResponseWriter, r *http.Request) {
	tenantID, payRunID, ok := h.tenantPayRunRequestIDs(w, r, "assess tenant pay run")
	if !ok {
		return
	}
	item, err := h.svc.AssessPayRunReadiness(r.Context(), tenantID, payRunID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondPayGroupError(w, r, "assess tenant pay run", err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) FreezeTenantPayRun(w http.ResponseWriter, r *http.Request) {
	h.payRunAction(w, r, true, h.svc.FreezePayRun, "freeze tenant pay run")
}

func (h *Handler) GenerateTenantPayRun(w http.ResponseWriter, r *http.Request) {
	h.payRunAction(w, r, true, h.svc.GeneratePayRun, "generate tenant pay run")
}

func (h *Handler) LockTenantPayRun(w http.ResponseWriter, r *http.Request) {
	h.payRunAction(w, r, true, h.svc.LockPayRun, "lock tenant pay run")
}

func (h *Handler) UnlockTenantPayRun(w http.ResponseWriter, r *http.Request) {
	h.payRunAction(w, r, true, h.svc.UnlockPayRun, "unlock tenant pay run")
}

func (h *Handler) DeleteTenantPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, payRunID, ok := h.tenantPayRunRequestIDs(w, r, "delete tenant pay run")
	if !ok {
		return
	}
	if err := h.svc.DeletePayRun(r.Context(), tenantID, payRunID, h.actorIDFromRequest(r)); err != nil {
		h.respondPayGroupError(w, r, "delete tenant pay run", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listPayGroupsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListPayGroups(r.Context(), tenantID)
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPayGroupForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PayGroupCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePayGroup(r.Context(), cmd)
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getPayGroupForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, payGroupID uuid.UUID, operation string) {
	item, err := h.svc.GetPayGroup(r.Context(), tenantID, payGroupID)
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updatePayGroupForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, payGroupID uuid.UUID, operation string) {
	var cmd ports.PayGroupCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = payGroupID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePayGroup(r.Context(), cmd)
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listPayGroupEmployeesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, payGroupID uuid.UUID, operation string) {
	items, err := h.svc.ListPayGroupEmployees(r.Context(), tenantID, payGroupID)
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) upsertPayGroupMemberForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, payGroupID uuid.UUID, operation string) {
	var cmd ports.PayGroupMemberCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.PayGroupID = payGroupID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertPayGroupMember(r.Context(), cmd)
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listPayRunsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListPayRuns(r.Context(), ports.PayRunListQuery{TenantID: tenantID, PayGroupID: optionalUUIDQuery(r, "pay_group_id"), Month: optionalInt32Query(r, "month"), Year: optionalInt32Query(r, "year")})
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPayRunForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PayRunCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePayRun(r.Context(), cmd)
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) getPayRunForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, payRunID uuid.UUID, operation string) {
	item, err := h.svc.GetPayRun(r.Context(), tenantID, payRunID)
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) payRunAction(w http.ResponseWriter, r *http.Request, tenantScoped bool, fn func(context.Context, ports.PayRunActionCommand) (*domain.PayRun, error), operation string) {
	var tenantID uuid.UUID
	var payRunID uuid.UUID
	var ok bool
	if tenantScoped {
		tenantID, payRunID, ok = h.tenantPayRunRequestIDs(w, r, operation)
	} else {
		tenantID, payRunID, ok = h.payRunRequestIDs(w, r, operation)
	}
	if !ok {
		return
	}
	cmd, ok := h.decodePayRunAction(r, operation, w)
	if !ok {
		return
	}
	cmd.TenantID = tenantID
	cmd.PayRunID = payRunID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := fn(r.Context(), cmd)
	if err != nil {
		h.respondPayGroupError(w, r, operation, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) decodePayRunAction(r *http.Request, operation string, w http.ResponseWriter) (ports.PayRunActionCommand, bool) {
	var cmd ports.PayRunActionCommand
	if r.Body == nil {
		return cmd, true
	}
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, http.ErrBodyNotAllowed) {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return cmd, false
	}
	return cmd, true
}

func (h *Handler) currentTenantID(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, false
	}
	return tenantID, true
}

func (h *Handler) payGroupRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.currentTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	payGroupID, err := uuid.Parse(chi.URLParam(r, "payGroupID"))
	if err != nil || payGroupID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid pay group id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, payGroupID, true
}

func (h *Handler) tenantPayGroupRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	payGroupID, err := uuid.Parse(chi.URLParam(r, "payGroupID"))
	if err != nil || payGroupID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid pay group id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, payGroupID, true
}

func (h *Handler) payGroupMemberRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.currentTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	memberID, err := uuid.Parse(chi.URLParam(r, "memberID"))
	if err != nil || memberID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid pay group member id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, memberID, true
}

func (h *Handler) tenantPayGroupMemberRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	memberID, err := uuid.Parse(chi.URLParam(r, "memberID"))
	if err != nil || memberID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid pay group member id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, memberID, true
}

func (h *Handler) payRunRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.currentTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	payRunID, err := uuid.Parse(chi.URLParam(r, "payRunID"))
	if err != nil || payRunID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid pay run id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, payRunID, true
}

func (h *Handler) tenantPayRunRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	payRunID, err := uuid.Parse(chi.URLParam(r, "payRunID"))
	if err != nil || payRunID == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid pay run id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, payRunID, true
}

func (h *Handler) respondPayGroupError(w http.ResponseWriter, r *http.Request, operation string, err error) {
	status := http.StatusBadRequest
	msg := err.Error()
	if errors.Is(err, domain.ErrPayGroupNotFound) || errors.Is(err, domain.ErrPayRunNotFound) {
		status = http.StatusNotFound
		msg = "record not found"
	} else if errors.Is(err, domain.ErrPayRunLocked) || errors.Is(err, domain.ErrPayRunBlocked) {
		status = http.StatusConflict
	}
	h.respondError(w, r, status, operation, err, msg)
}

func optionalInt32Query(r *http.Request, key string) *int32 {
	if r == nil {
		return nil
	}
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return nil
	}
	result := int32(parsed)
	return &result
}
