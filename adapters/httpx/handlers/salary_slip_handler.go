package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetSalarySlipFormat(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get salary slip format", err, "tenant context is required")
		return
	}
	h.getSalarySlipFormatForTenant(w, r, tenantID, "get salary slip format")
}

func (h *Handler) UpsertSalarySlipFormat(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "upsert salary slip format", err, "tenant context is required")
		return
	}
	h.upsertSalarySlipFormatForTenant(w, r, tenantID, "upsert salary slip format")
}

func (h *Handler) GenerateSalarySlip(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "generate salary slip", err, "tenant context is required")
		return
	}
	h.generateSalarySlipForTenant(w, r, tenantID, "generate salary slip")
}

func (h *Handler) ListSalarySlips(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list salary slips", err, "tenant context is required")
		return
	}
	h.listSalarySlipsForTenant(w, r, tenantID, "list salary slips")
}

func (h *Handler) DownloadSalarySlipPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, slipID, ok := h.salarySlipRequestIDs(w, r, "download salary slip pdf")
	if !ok {
		return
	}
	h.downloadSalarySlipPDFForTenant(w, r, tenantID, slipID, "download salary slip pdf")
}

func (h *Handler) DownloadMyRecentSalarySlips(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "download recent salary slips", err, "tenant context is required")
		return
	}
	userID := h.userIDFromContext(r.Context())
	months := int32(queryInt(r, "months", 6))
	data, name, err := h.svc.RenderRecentSalarySlipsZip(r.Context(), tenantID, userID, months)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "download recent salary slips", err, err.Error())
		return
	}
	respondDownload(w, "application/zip", name, data)
}

func (h *Handler) GetTenantSalarySlipFormat(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant salary slip format")
	if !ok {
		return
	}
	h.getSalarySlipFormatForTenant(w, r, tenantID, "get tenant salary slip format")
}

func (h *Handler) UpsertTenantSalarySlipFormat(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "upsert tenant salary slip format")
	if !ok {
		return
	}
	h.upsertSalarySlipFormatForTenant(w, r, tenantID, "upsert tenant salary slip format")
}

func (h *Handler) GenerateTenantSalarySlip(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "generate tenant salary slip")
	if !ok {
		return
	}
	h.generateSalarySlipForTenant(w, r, tenantID, "generate tenant salary slip")
}

func (h *Handler) BulkGenerateTenantSalarySlips(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "bulk generate tenant salary slips")
	if !ok {
		return
	}
	var cmd ports.BulkGenerateSalarySlipsCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode bulk salary slip request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	items, err := h.svc.BulkGenerateSalarySlips(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "bulk generate tenant salary slips", err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, items)
}

func (h *Handler) ListTenantSalarySlips(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant salary slips")
	if !ok {
		return
	}
	h.listSalarySlipsForTenant(w, r, tenantID, "list tenant salary slips")
}

func (h *Handler) DownloadTenantSalarySlipPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, slipID, ok := h.tenantSalarySlipRequestIDs(w, r, "download tenant salary slip pdf")
	if !ok {
		return
	}
	h.downloadSalarySlipPDFForTenant(w, r, tenantID, slipID, "download tenant salary slip pdf")
}

func (h *Handler) DownloadTenantSalarySlipsBulkPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "download tenant salary slips bulk pdf")
	if !ok {
		return
	}
	month := int32(queryInt(r, "month", 0))
	year := int32(queryInt(r, "year", 0))
	data, name, err := h.svc.RenderTenantSalarySlipsZip(r.Context(), tenantID, month, year)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "download tenant salary slips bulk pdf", err, err.Error())
		return
	}
	respondDownload(w, "application/zip", name, data)
}

func (h *Handler) getSalarySlipFormatForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	item, err := h.svc.GetSalarySlipFormat(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to get salary slip format")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) upsertSalarySlipFormatForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.SalarySlipFormatCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode salary slip format request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertSalarySlipFormat(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) generateSalarySlipForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.GenerateSalarySlipCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode salary slip request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.GenerateSalarySlip(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listSalarySlipsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	userIDRaw := r.URL.Query().Get("user_id")
	if userIDRaw != "" {
		userID, err := uuid.Parse(userIDRaw)
		if err != nil {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid user id")
			return
		}
		items, err := h.svc.ListSalarySlipsByUser(r.Context(), tenantID, userID)
		if err != nil {
			h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list salary slips")
			return
		}
		respondJSON(w, http.StatusOK, items)
		return
	}
	items, err := h.svc.ListSalarySlipsByTenantPeriod(r.Context(), tenantID, int32(queryInt(r, "month", 0)), int32(queryInt(r, "year", 0)))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list salary slips")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) downloadSalarySlipPDFForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, slipID uuid.UUID, operation string) {
	data, name, err := h.svc.RenderSalarySlipPDF(r.Context(), tenantID, slipID)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondDownload(w, "application/pdf", name, data)
}

func (h *Handler) salarySlipRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	slipID, err := uuid.Parse(chi.URLParam(r, "slipID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid salary slip id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, slipID, true
}

func (h *Handler) tenantSalarySlipRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	slipID, err := uuid.Parse(chi.URLParam(r, "slipID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid salary slip id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, slipID, true
}

func queryInt(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func respondDownload(w http.ResponseWriter, contentType string, filename string, data []byte) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
