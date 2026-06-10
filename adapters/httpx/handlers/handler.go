package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-platform/httpserver"
	"y/core/ports"
)

// Handler handles all HTTP requests for the Hrms module.
type Handler struct {
	svc ports.HrmsService
}

func New(svc ports.HrmsService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := uuid.MustParse(httpserver.GetTenantID(r.Context()))

	items, err := h.svc.List(r.Context(), tenantID)
	if err != nil {
		respondError(w, 500, "failed to list")
		return
	}
	respondJSON(w, 200, items)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := uuid.MustParse(httpserver.GetTenantID(r.Context()))
	userID   := uuid.MustParse(httpserver.GetUserID(r.Context()))

	cmd := ports.CreateHrmsCmd{
		TenantID:  tenantID,
		CreatedBy: userID,
	}
	// TODO: decode additional fields from r.Body into cmd

	item, err := h.svc.Create(r.Context(), cmd)
	if err != nil {
		respondError(w, 500, "failed to create")
		return
	}
	respondJSON(w, 201, item)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := uuid.MustParse(httpserver.GetTenantID(r.Context()))
	id       := uuid.MustParse(chi.URLParam(r, "id"))

	item, err := h.svc.Get(r.Context(), id, tenantID)
	if err != nil || item == nil {
		respondError(w, 404, "not found")
		return
	}
	respondJSON(w, 200, item)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := uuid.MustParse(httpserver.GetTenantID(r.Context()))
	id       := uuid.MustParse(chi.URLParam(r, "id"))

	if err := h.svc.Delete(r.Context(), id, tenantID); err != nil {
		respondError(w, 500, "failed to delete")
		return
	}
	w.WriteHeader(204)
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
