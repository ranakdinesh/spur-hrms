package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create salary template", err, "tenant context is required")
		return
	}
	h.createSalaryTemplateForTenant(w, r, tenantID, "create salary template")
}

func (h *Handler) ListSalaryTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list salary templates", err, "tenant context is required")
		return
	}
	h.listSalaryTemplatesForTenant(w, r, tenantID, "list salary templates")
}

func (h *Handler) GetSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.salaryTemplateRequestIDs(w, r, "get salary template")
	if !ok {
		return
	}
	item, err := h.svc.GetSalaryTemplate(r.Context(), tenantID, templateID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get salary template", err, "salary template not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.salaryTemplateRequestIDs(w, r, "update salary template")
	if !ok {
		return
	}
	h.updateSalaryTemplateForTenant(w, r, tenantID, templateID, "update salary template")
}

func (h *Handler) DeleteSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.salaryTemplateRequestIDs(w, r, "delete salary template")
	if !ok {
		return
	}
	if err := h.svc.DeleteSalaryTemplate(r.Context(), tenantID, templateID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete salary template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ActivateSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.salaryTemplateRequestIDs(w, r, "activate salary template")
	if !ok {
		return
	}
	item, err := h.svc.ActivateSalaryTemplate(r.Context(), tenantID, templateID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "activate salary template", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateSalaryTemplateItem(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.salaryTemplateRequestIDs(w, r, "create salary template item")
	if !ok {
		return
	}
	h.createSalaryTemplateItemForTenant(w, r, tenantID, templateID, "create salary template item")
}

func (h *Handler) ListSalaryTemplateItems(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.salaryTemplateRequestIDs(w, r, "list salary template items")
	if !ok {
		return
	}
	items, err := h.svc.ListSalaryTemplateItems(r.Context(), tenantID, templateID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list salary template items", err, "failed to list salary template items")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpdateSalaryTemplateItem(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, itemID, ok := h.salaryTemplateItemRequestIDs(w, r, "update salary template item")
	if !ok {
		return
	}
	h.updateSalaryTemplateItemForTenant(w, r, tenantID, templateID, itemID, "update salary template item")
}

func (h *Handler) DeleteSalaryTemplateItem(w http.ResponseWriter, r *http.Request) {
	tenantID, _, itemID, ok := h.salaryTemplateItemRequestIDs(w, r, "delete salary template item")
	if !ok {
		return
	}
	if err := h.svc.DeleteSalaryTemplateItem(r.Context(), tenantID, itemID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete salary template item", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant salary template")
	if !ok {
		return
	}
	h.createSalaryTemplateForTenant(w, r, tenantID, "create tenant salary template")
}

func (h *Handler) ListTenantSalaryTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant salary templates")
	if !ok {
		return
	}
	h.listSalaryTemplatesForTenant(w, r, tenantID, "list tenant salary templates")
}

func (h *Handler) GetTenantSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.tenantSalaryTemplateRequestIDs(w, r, "get tenant salary template")
	if !ok {
		return
	}
	item, err := h.svc.GetSalaryTemplate(r.Context(), tenantID, templateID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant salary template", err, "salary template not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.tenantSalaryTemplateRequestIDs(w, r, "update tenant salary template")
	if !ok {
		return
	}
	h.updateSalaryTemplateForTenant(w, r, tenantID, templateID, "update tenant salary template")
}

func (h *Handler) DeleteTenantSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.tenantSalaryTemplateRequestIDs(w, r, "delete tenant salary template")
	if !ok {
		return
	}
	if err := h.svc.DeleteSalaryTemplate(r.Context(), tenantID, templateID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant salary template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ActivateTenantSalaryTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.tenantSalaryTemplateRequestIDs(w, r, "activate tenant salary template")
	if !ok {
		return
	}
	item, err := h.svc.ActivateSalaryTemplate(r.Context(), tenantID, templateID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "activate tenant salary template", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateTenantSalaryTemplateItem(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.tenantSalaryTemplateRequestIDs(w, r, "create tenant salary template item")
	if !ok {
		return
	}
	h.createSalaryTemplateItemForTenant(w, r, tenantID, templateID, "create tenant salary template item")
}

func (h *Handler) ListTenantSalaryTemplateItems(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, ok := h.tenantSalaryTemplateRequestIDs(w, r, "list tenant salary template items")
	if !ok {
		return
	}
	items, err := h.svc.ListSalaryTemplateItems(r.Context(), tenantID, templateID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant salary template items", err, "failed to list salary template items")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) UpdateTenantSalaryTemplateItem(w http.ResponseWriter, r *http.Request) {
	tenantID, templateID, itemID, ok := h.tenantSalaryTemplateItemRequestIDs(w, r, "update tenant salary template item")
	if !ok {
		return
	}
	h.updateSalaryTemplateItemForTenant(w, r, tenantID, templateID, itemID, "update tenant salary template item")
}

func (h *Handler) DeleteTenantSalaryTemplateItem(w http.ResponseWriter, r *http.Request) {
	tenantID, _, itemID, ok := h.tenantSalaryTemplateItemRequestIDs(w, r, "delete tenant salary template item")
	if !ok {
		return
	}
	if err := h.svc.DeleteSalaryTemplateItem(r.Context(), tenantID, itemID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant salary template item", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createSalaryTemplateForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.SalaryTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode salary template request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateSalaryTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listSalaryTemplatesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var fyID *uuid.UUID
	if raw := r.URL.Query().Get("fy_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid financial year id")
			return
		}
		fyID = &parsed
	}
	items, err := h.svc.ListSalaryTemplates(r.Context(), tenantID, fyID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list salary templates")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateSalaryTemplateForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, templateID uuid.UUID, operation string) {
	var cmd ports.SalaryTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode salary template request", err, "invalid request body")
		return
	}
	cmd.ID = templateID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateSalaryTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createSalaryTemplateItemForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, templateID uuid.UUID, operation string) {
	var cmd ports.SalaryTemplateItemCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode salary template item request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.TemplateID = templateID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateSalaryTemplateItem(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateSalaryTemplateItemForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, templateID uuid.UUID, itemID uuid.UUID, operation string) {
	var cmd ports.SalaryTemplateItemCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode salary template item request", err, "invalid request body")
		return
	}
	cmd.ID = itemID
	cmd.TenantID = tenantID
	cmd.TemplateID = templateID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateSalaryTemplateItem(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) salaryTemplateRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	templateID, err := uuid.Parse(chi.URLParam(r, "templateID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid salary template id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, templateID, true
}

func (h *Handler) salaryTemplateItemRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, uuid.UUID, bool) {
	tenantID, templateID, ok := h.salaryTemplateRequestIDs(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid salary template item id")
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	return tenantID, templateID, itemID, true
}

func (h *Handler) tenantSalaryTemplateRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	templateID, err := uuid.Parse(chi.URLParam(r, "templateID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid salary template id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, templateID, true
}

func (h *Handler) tenantSalaryTemplateItemRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, uuid.UUID, bool) {
	tenantID, templateID, ok := h.tenantSalaryTemplateRequestIDs(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	itemID, err := uuid.Parse(chi.URLParam(r, "itemID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid salary template item id")
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	return tenantID, templateID, itemID, true
}
