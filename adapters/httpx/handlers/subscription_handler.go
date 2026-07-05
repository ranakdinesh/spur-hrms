package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateTenantSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant subscription")
	if !ok {
		return
	}
	h.createSubscriptionForTenant(w, r, tenantID, "create tenant subscription")
}

func (h *Handler) ListTenantSubscriptions(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant subscriptions")
	if !ok {
		return
	}
	h.listSubscriptionsForTenant(w, r, tenantID, "list tenant subscriptions")
}

func (h *Handler) GetTenantCurrentSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant current subscription")
	if !ok {
		return
	}
	item, err := h.svc.GetCurrentTenantSubscription(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant current subscription", err, "current subscription not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) GetTenantSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminSubscriptionRequestIDs(w, r, "get tenant subscription")
	if !ok {
		return
	}
	item, err := h.svc.GetTenantSubscription(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant subscription", err, "subscription not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminSubscriptionRequestIDs(w, r, "update tenant subscription")
	if !ok {
		return
	}
	h.updateSubscriptionForTenant(w, r, tenantID, id, "update tenant subscription")
}

func (h *Handler) DeleteTenantSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminSubscriptionRequestIDs(w, r, "delete tenant subscription")
	if !ok {
		return
	}
	if err := h.svc.DeleteTenantSubscription(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant subscription", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list subscriptions", err, "tenant context is required")
		return
	}
	h.listSubscriptionsForTenant(w, r, tenantID, "list subscriptions")
}

func (h *Handler) GetCurrentSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get current subscription", err, "tenant context is required")
		return
	}
	item, err := h.svc.GetCurrentTenantSubscription(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get current subscription", err, "current subscription not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.subscriptionRequestIDs(w, r, "get subscription")
	if !ok {
		return
	}
	item, err := h.svc.GetTenantSubscription(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get subscription", err, "subscription not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createSubscriptionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.TenantSubscriptionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateTenantSubscription(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listSubscriptionsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListTenantSubscriptions(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list subscriptions")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateSubscriptionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.TenantSubscriptionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateTenantSubscription(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) subscriptionRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "subscriptionID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid subscription id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminSubscriptionRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "subscriptionID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid subscription id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
