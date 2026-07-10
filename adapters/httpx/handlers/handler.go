package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/internal/logging"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
	"github.com/rs/zerolog"
)

type Handler struct {
	svc                    ports.TenantService
	tenantIDFromContext    func(context.Context) string
	userIDFromContext      func(context.Context) uuid.UUID
	isSuperAdmin           func(context.Context) bool
	rolesFromContext       func(context.Context) []string
	permissionsFromContext func(context.Context) []string
	log                    *zerolog.Logger
}

func New(svc ports.TenantService, tenantIDFromContext func(context.Context) string, userIDFromContext func(context.Context) uuid.UUID, isSuperAdmin func(context.Context) bool, rolesFromContext func(context.Context) []string, permissionsFromContext func(context.Context) []string, log ...*zerolog.Logger) *Handler {
	return &Handler{svc: svc, tenantIDFromContext: tenantIDFromContext, userIDFromContext: userIDFromContext, isSuperAdmin: isSuperAdmin, rolesFromContext: rolesFromContext, permissionsFromContext: permissionsFromContext, log: logging.Component(logging.First(log...), "http_handler")}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get tenant profile", err, "tenant context is required")
		return
	}
	profile, err := h.svc.GetTenantProfile(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant profile", err, "tenant profile not found")
		return
	}
	respondJSON(w, http.StatusOK, profile)
}

func (h *Handler) UpsertProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert tenant profile", err, "tenant context is required")
		return
	}
	var cmd ports.UpsertTenantProfileCmd
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant profile request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	profile, err := h.svc.UpsertTenantProfile(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert tenant profile", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, profile)
}

func (h *Handler) ResolveBySubdomain(w http.ResponseWriter, r *http.Request) {
	profile, err := h.svc.ResolveTenantBySubdomain(r.Context(), chi.URLParam(r, "subdomain"))
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "resolve tenant by subdomain", err, "tenant profile not found")
		return
	}
	respondJSON(w, http.StatusOK, publicTenantProfile(profile))
}

func (h *Handler) ResolveByActivationCode(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode activation-code request", err, "invalid request body")
		return
	}
	profile, err := h.svc.ResolveTenantByActivationCode(r.Context(), body.Code)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "resolve tenant by activation code", err, "tenant profile not found")
		return
	}
	respondJSON(w, http.StatusOK, publicTenantProfile(profile))
}

func (h *Handler) ListSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list tenant settings", err, "tenant context is required")
		return
	}
	settings, err := h.svc.ListTenantSettings(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant settings", err, "failed to list settings")
		return
	}
	respondJSON(w, http.StatusOK, settings)
}

func (h *Handler) GetSetting(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get tenant setting", err, "tenant context is required")
		return
	}
	setting, err := h.svc.GetTenantSetting(r.Context(), tenantID, chi.URLParam(r, "key"))
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant setting", err, "tenant setting not found")
		return
	}
	respondJSON(w, http.StatusOK, setting)
}

func (h *Handler) UpsertSetting(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert tenant setting", err, "tenant context is required")
		return
	}
	var cmd ports.UpsertTenantSettingCmd
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode tenant setting request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.Key = chi.URLParam(r, "key")
	setting, err := h.svc.UpsertTenantSetting(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "upsert tenant setting", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, setting)
}

func (h *Handler) DeleteSetting(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete tenant setting", err, "tenant context is required")
		return
	}
	if err := h.svc.DeleteTenantSetting(r.Context(), tenantID, chi.URLParam(r, "key")); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant setting", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	if r == nil {
		return uuid.Nil, errors.New("request is nil")
	}
	if h == nil || h.tenantIDFromContext == nil {
		return uuid.Nil, errors.New("tenant resolver is not configured")
	}
	return uuid.Parse(h.tenantIDFromContext(r.Context()))
}

func (h *Handler) actorIDFromRequest(r *http.Request) *uuid.UUID {
	if h == nil || h.userIDFromContext == nil || r == nil {
		return nil
	}
	userID := h.userIDFromContext(r.Context())
	if userID == uuid.Nil {
		return nil
	}
	return &userID
}

func (h *Handler) isSuperAdminRequest(r *http.Request) bool {
	if h == nil || h.isSuperAdmin == nil || r == nil {
		return false
	}
	return h.isSuperAdmin(r.Context())
}

func (h *Handler) requirePermission(w http.ResponseWriter, r *http.Request, operation string, permission string) bool {
	if h.isSuperAdminRequest(r) {
		return true
	}
	if h.hasPermission(r, permission) {
		return true
	}
	h.respondError(w, r, http.StatusForbidden, operation, nil, "permission required")
	return false
}

func (h *Handler) hasPermission(r *http.Request, permission string) bool {
	if r == nil || strings.TrimSpace(permission) == "" {
		return false
	}
	required := strings.TrimSpace(permission)
	fullRequired := required
	if !strings.HasPrefix(required, permissions.ModuleCode+".") {
		fullRequired = permissions.ModuleCode + "." + required
	}
	localRequired := strings.TrimPrefix(required, permissions.ModuleCode+".")

	if h == nil || h.permissionsFromContext == nil {
		return false
	}
	for _, granted := range h.permissionsFromContext(r.Context()) {
		granted = strings.TrimSpace(granted)
		if granted == "" {
			continue
		}
		if granted == fullRequired || granted == localRequired {
			return true
		}
	}
	return false
}

func (h *Handler) hasAnyPermission(r *http.Request, permissions ...string) bool {
	for _, permission := range permissions {
		if h.hasPermission(r, permission) {
			return true
		}
	}
	return false
}

func (h *Handler) requireOwnUserOrPermission(w http.ResponseWriter, r *http.Request, operation string, targetUserID uuid.UUID, selfPermissions []string, operationsPermissions []string) bool {
	if h.isSuperAdminRequest(r) {
		return true
	}
	if h.hasAnyPermission(r, operationsPermissions...) {
		return true
	}
	actorID := h.actorIDFromRequest(r)
	if actorID != nil && *actorID == targetUserID && h.hasAnyPermission(r, selfPermissions...) {
		return true
	}
	h.respondError(w, r, http.StatusForbidden, operation, nil, "permission required")
	return false
}

func (h *Handler) respondError(w http.ResponseWriter, r *http.Request, status int, operation string, err error, msg string) {
	h.logRequestFailure(r, status, operation, err)
	respondError(w, status, msg)
}

func (h *Handler) logRequestFailure(r *http.Request, status int, operation string, err error) {
	if h == nil || h.log == nil {
		return
	}
	event := h.log.Warn()
	if status >= http.StatusInternalServerError {
		event = h.log.Error()
	}
	if err != nil {
		event.Err(err)
	}
	event.Int("status", status).Str("operation", operation)
	if r != nil {
		event.Str("method", r.Method).Str("path", r.URL.Path)
		if routeContext := chi.RouteContext(r.Context()); routeContext != nil {
			if route := routeContext.RoutePattern(); route != "" {
				event.Str("route", route)
			}
		}
		if h.tenantIDFromContext != nil {
			if tenantRaw := h.tenantIDFromContext(r.Context()); tenantRaw != "" {
				event.Str("tenant_id", tenantRaw)
			}
		}
	}
	event.Msg("hrms request failed")
}

func publicTenantProfile(profile *domain.TenantProfile) any {
	return struct {
		TenantID      uuid.UUID `json:"tenant_id"`
		Subdomain     string    `json:"subdomain"`
		DisplayName   *string   `json:"display_name,omitempty"`
		LogoObjectKey *string   `json:"logo_object_key,omitempty"`
	}{
		TenantID:      profile.TenantID,
		Subdomain:     profile.Subdomain,
		DisplayName:   profile.DisplayName,
		LogoObjectKey: profile.LogoObjectKey,
	}
}

func respondJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
