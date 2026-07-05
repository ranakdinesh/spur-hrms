package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetPublicAgreement(w http.ResponseWriter, r *http.Request) {
	item, err := h.svc.GetAgreementBySignatureToken(r.Context(), chi.URLParam(r, "signatureToken"))
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get public agreement", err, "agreement not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) SignPublicAgreement(w http.ResponseWriter, r *http.Request) {
	var cmd ports.AgreementSignatureCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "sign public agreement", err, "invalid request body")
		return
	}
	cmd.Token = chi.URLParam(r, "signatureToken")
	ip := clientIP(r)
	ua := r.UserAgent()
	cmd.IPAddress = &ip
	cmd.UserAgent = &ua
	item, err := h.svc.SignAgreement(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "sign public agreement", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateAgreementTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create agreement template", err, "tenant context is required")
		return
	}
	h.createAgreementTemplate(w, r, tenantID, "create agreement template")
}

func (h *Handler) ListAgreementTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list agreement templates", err, "tenant context is required")
		return
	}
	h.listAgreementTemplates(w, r, tenantID, "list agreement templates")
}

func (h *Handler) UpdateAgreementTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.agreementTemplateRequestIDs(w, r, "update agreement template")
	if !ok {
		return
	}
	h.updateAgreementTemplate(w, r, tenantID, id, "update agreement template")
}

func (h *Handler) DeleteAgreementTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.agreementTemplateRequestIDs(w, r, "delete agreement template")
	if !ok {
		return
	}
	if err := h.svc.DeleteAgreementTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete agreement template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GenerateAgreement(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "generate agreement", err, "tenant context is required")
		return
	}
	h.generateAgreement(w, r, tenantID, "generate agreement")
}

func (h *Handler) ListAgreements(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list agreements", err, "tenant context is required")
		return
	}
	h.listAgreements(w, r, tenantID, "list agreements")
}

func (h *Handler) GetAgreement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.agreementRequestIDs(w, r, "get agreement")
	if !ok {
		return
	}
	item, err := h.svc.GetAgreement(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get agreement", err, "agreement not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateAgreementStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.agreementRequestIDs(w, r, "update agreement status")
	if !ok {
		return
	}
	h.updateAgreementStatus(w, r, tenantID, id, "update agreement status")
}

func (h *Handler) DownloadAgreementPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.agreementRequestIDs(w, r, "download agreement pdf")
	if !ok {
		return
	}
	h.downloadAgreementPDF(w, r, tenantID, id, "download agreement pdf")
}

func (h *Handler) ListAgreementEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.agreementRequestIDs(w, r, "list agreement events")
	if !ok {
		return
	}
	items, err := h.svc.ListAgreementEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list agreement events", err, "failed to list agreement events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteAgreement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.agreementRequestIDs(w, r, "delete agreement")
	if !ok {
		return
	}
	if err := h.svc.DeleteAgreement(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete agreement", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantAgreementTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant agreement template")
	if !ok {
		return
	}
	h.createAgreementTemplate(w, r, tenantID, "create tenant agreement template")
}

func (h *Handler) ListTenantAgreementTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant agreement templates")
	if !ok {
		return
	}
	h.listAgreementTemplates(w, r, tenantID, "list tenant agreement templates")
}

func (h *Handler) UpdateTenantAgreementTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "agreementTemplateID", "update tenant agreement template")
	if !ok {
		return
	}
	h.updateAgreementTemplate(w, r, tenantID, id, "update tenant agreement template")
}

func (h *Handler) DeleteTenantAgreementTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "agreementTemplateID", "delete tenant agreement template")
	if !ok {
		return
	}
	if err := h.svc.DeleteAgreementTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant agreement template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GenerateTenantAgreement(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "generate tenant agreement")
	if !ok {
		return
	}
	h.generateAgreement(w, r, tenantID, "generate tenant agreement")
}

func (h *Handler) ListTenantAgreements(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant agreements")
	if !ok {
		return
	}
	h.listAgreements(w, r, tenantID, "list tenant agreements")
}

func (h *Handler) GetTenantAgreement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "agreementID", "get tenant agreement")
	if !ok {
		return
	}
	item, err := h.svc.GetAgreement(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant agreement", err, "agreement not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantAgreementStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "agreementID", "update tenant agreement status")
	if !ok {
		return
	}
	h.updateAgreementStatus(w, r, tenantID, id, "update tenant agreement status")
}

func (h *Handler) DownloadTenantAgreementPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "agreementID", "download tenant agreement pdf")
	if !ok {
		return
	}
	h.downloadAgreementPDF(w, r, tenantID, id, "download tenant agreement pdf")
}

func (h *Handler) ListTenantAgreementEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "agreementID", "list tenant agreement events")
	if !ok {
		return
	}
	items, err := h.svc.ListAgreementEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant agreement events", err, "failed to list agreement events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteTenantAgreement(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "agreementID", "delete tenant agreement")
	if !ok {
		return
	}
	if err := h.svc.DeleteAgreement(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant agreement", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createAgreementTemplate(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AgreementTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateAgreementTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listAgreementTemplates(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAgreementTemplates(r.Context(), tenantID, optionalStringQuery(r, "agreement_type"))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list agreement templates")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateAgreementTemplate(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.AgreementTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateAgreementTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) generateAgreement(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.AgreementCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.GenerateAgreement(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listAgreements(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListAgreements(r.Context(), domain.AgreementFilter{
		TenantID:        tenantID,
		AgreementType:   optionalStringQuery(r, "agreement_type"),
		Status:          optionalStringQuery(r, "status"),
		WorkerProfileID: optionalUUIDQuery(r, "worker_profile_id"),
		EngagementID:    optionalUUIDQuery(r, "engagement_id"),
		ProjectID:       optionalUUIDQuery(r, "project_id"),
		Search:          optionalStringQuery(r, "search"),
	})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list agreements")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateAgreementStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.AgreementStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateAgreementStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) downloadAgreementPDF(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	data, name, err := h.svc.RenderAgreementPDF(r.Context(), tenantID, id, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondDownload(w, "application/pdf", name, data)
}

func (h *Handler) agreementTemplateRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "agreementTemplateID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid agreement template id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) agreementRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "agreementID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid agreement id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
