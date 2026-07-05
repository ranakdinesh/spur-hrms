package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateSubscriptionPlan(w http.ResponseWriter, r *http.Request) {
	if !h.requireSuperAdmin(w, r, "create subscription plan") {
		return
	}
	var cmd ports.SubscriptionPlanCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode subscription plan request", err, "invalid request body")
		return
	}
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateSubscriptionPlan(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "create subscription plan", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	if !h.requireSuperAdmin(w, r, "list subscription plans") {
		return
	}
	items, err := h.svc.ListSubscriptionPlans(r.Context())
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list subscription plans", err, "failed to list subscription plans")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListActiveSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
	if !h.requireSuperAdmin(w, r, "list active subscription plans") {
		return
	}
	items, err := h.svc.ListActiveSubscriptionPlans(r.Context())
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list active subscription plans", err, "failed to list active subscription plans")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) GetSubscriptionPlan(w http.ResponseWriter, r *http.Request) {
	if !h.requireSuperAdmin(w, r, "get subscription plan") {
		return
	}
	id, ok := h.subscriptionPlanID(w, r, "get subscription plan")
	if !ok {
		return
	}
	item, err := h.svc.GetSubscriptionPlan(r.Context(), id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get subscription plan", err, "subscription plan not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateSubscriptionPlan(w http.ResponseWriter, r *http.Request) {
	if !h.requireSuperAdmin(w, r, "update subscription plan") {
		return
	}
	id, ok := h.subscriptionPlanID(w, r, "update subscription plan")
	if !ok {
		return
	}
	var cmd ports.SubscriptionPlanCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode subscription plan request", err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateSubscriptionPlan(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "update subscription plan", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) DeleteSubscriptionPlan(w http.ResponseWriter, r *http.Request) {
	if !h.requireSuperAdmin(w, r, "delete subscription plan") {
		return
	}
	id, ok := h.subscriptionPlanID(w, r, "delete subscription plan")
	if !ok {
		return
	}
	if err := h.svc.DeleteSubscriptionPlan(r.Context(), id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete subscription plan", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) requireSuperAdmin(w http.ResponseWriter, r *http.Request, operation string) bool {
	if !h.isSuperAdminRequest(r) {
		h.respondError(w, r, http.StatusForbidden, operation, nil, "super admin permission required")
		return false
	}
	return true
}

func (h *Handler) subscriptionPlanID(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "planID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid subscription plan id")
		return uuid.Nil, false
	}
	return id, true
}
