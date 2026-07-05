package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetTenantBranding(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "tenantID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant branding tenant id", err, "invalid tenant id")
		return
	}
	if !h.isSuperAdminRequest(r) {
		currentTenantID, err := h.tenantIDFromRequest(r)
		if err != nil {
			h.respondError(w, r, http.StatusUnauthorized, "get tenant branding", err, "tenant context is required")
			return
		}
		if currentTenantID != tenantID {
			h.respondError(w, r, http.StatusForbidden, "get tenant branding", nil, "tenant branding access denied")
			return
		}
	}
	branding, err := h.svc.GetTenantBranding(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant branding", err, "tenant branding not found")
		return
	}
	respondJSON(w, http.StatusOK, branding)
}

func (h *Handler) GetCurrentTenantBranding(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get current tenant branding", err, "tenant context is required")
		return
	}
	branding, err := h.svc.GetTenantBranding(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get current tenant branding", err, "tenant branding not found")
		return
	}
	respondJSON(w, http.StatusOK, branding)
}

func (h *Handler) UpsertTenantBranding(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "tenantID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant branding tenant id", err, "invalid tenant id")
		return
	}
	if !h.isSuperAdminRequest(r) {
		currentTenantID, err := h.tenantIDFromRequest(r)
		if err != nil {
			h.respondError(w, r, http.StatusUnauthorized, "upsert tenant branding", err, "tenant context is required")
			return
		}
		if currentTenantID != tenantID {
			h.respondError(w, r, http.StatusForbidden, "upsert tenant branding", nil, "tenant branding access denied")
			return
		}
	}
	var cmd ports.UpsertTenantBrandingCmd
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant branding request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	if !h.isSuperAdminRequest(r) {
		profile, err := h.svc.GetTenantProfile(r.Context(), tenantID)
		if err != nil {
			h.respondError(w, r, http.StatusNotFound, "get tenant profile for branding", err, "tenant profile not found")
			return
		}
		cmd.Subdomain = profile.Subdomain
		if cmd.DisplayName == nil {
			cmd.DisplayName = profile.DisplayName
		}
	}
	branding, err := h.svc.UpsertTenantBranding(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert tenant branding", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, branding)
}

func (h *Handler) UpsertCurrentTenantBranding(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert current tenant branding", err, "tenant context is required")
		return
	}
	profile, err := h.svc.GetTenantProfile(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get current tenant profile for branding", err, "tenant profile not found")
		return
	}
	var cmd ports.UpsertTenantBrandingCmd
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode current tenant branding request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.Subdomain = profile.Subdomain
	cmd.ActorID = h.actorIDFromRequest(r)
	if cmd.DisplayName == nil {
		cmd.DisplayName = profile.DisplayName
	}
	branding, err := h.svc.UpsertTenantBranding(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert current tenant branding", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, branding)
}

func (h *Handler) ResolveBrandingByHost(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		host = r.Host
	}
	branding, err := h.svc.ResolveTenantBrandingByHost(r.Context(), host)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "resolve tenant branding by host", err, "tenant branding not found")
		return
	}
	respondJSON(w, http.StatusOK, branding)
}
