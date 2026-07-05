package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetPublicOfferLetter(w http.ResponseWriter, r *http.Request) {
	item, err := h.svc.GetOfferLetterBySignatureToken(r.Context(), chi.URLParam(r, "signatureToken"))
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get public offer", err, "offer letter not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) SignPublicOfferLetter(w http.ResponseWriter, r *http.Request) {
	var cmd ports.OfferLetterSignatureCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "sign public offer", err, "invalid request body")
		return
	}
	cmd.Token = chi.URLParam(r, "signatureToken")
	ip := clientIP(r)
	ua := r.UserAgent()
	cmd.IPAddress = &ip
	cmd.UserAgent = &ua
	item, err := h.svc.SignOfferLetter(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "sign public offer", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateOfferLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create offer template", err, "tenant context is required")
		return
	}
	h.createOfferLetterTemplate(w, r, tenantID, "create offer template")
}

func (h *Handler) ListOfferLetterTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list offer templates", err, "tenant context is required")
		return
	}
	h.listOfferLetterTemplates(w, r, tenantID, "list offer templates")
}

func (h *Handler) UpdateOfferLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.offerTemplateRequestIDs(w, r, "update offer template")
	if !ok {
		return
	}
	h.updateOfferLetterTemplate(w, r, tenantID, id, "update offer template")
}

func (h *Handler) DeleteOfferLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.offerTemplateRequestIDs(w, r, "delete offer template")
	if !ok {
		return
	}
	if err := h.svc.DeleteOfferLetterTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete offer template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateOfferLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create offer", err, "tenant context is required")
		return
	}
	h.createOfferLetter(w, r, tenantID, "create offer")
}

func (h *Handler) ListOfferLetters(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list offers", err, "tenant context is required")
		return
	}
	h.listOfferLetters(w, r, tenantID, "list offers")
}

func (h *Handler) GetOfferLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.offerLetterRequestIDs(w, r, "get offer")
	if !ok {
		return
	}
	item, err := h.svc.GetOfferLetter(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get offer", err, "offer letter not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateOfferLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.offerLetterRequestIDs(w, r, "update offer")
	if !ok {
		return
	}
	h.updateOfferLetter(w, r, tenantID, id, "update offer")
}

func (h *Handler) UpdateOfferLetterStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.offerLetterRequestIDs(w, r, "update offer status")
	if !ok {
		return
	}
	h.updateOfferLetterStatus(w, r, tenantID, id, "update offer status")
}

func (h *Handler) ListOfferLetterEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.offerLetterRequestIDs(w, r, "list offer events")
	if !ok {
		return
	}
	items, err := h.svc.ListOfferLetterEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list offer events", err, "failed to list offer events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteOfferLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.offerLetterRequestIDs(w, r, "delete offer")
	if !ok {
		return
	}
	if err := h.svc.DeleteOfferLetter(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete offer", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantOfferLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant offer template")
	if !ok {
		return
	}
	h.createOfferLetterTemplate(w, r, tenantID, "create tenant offer template")
}

func (h *Handler) ListTenantOfferLetterTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant offer templates")
	if !ok {
		return
	}
	h.listOfferLetterTemplates(w, r, tenantID, "list tenant offer templates")
}

func (h *Handler) UpdateTenantOfferLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "offerTemplateID", "update tenant offer template")
	if !ok {
		return
	}
	h.updateOfferLetterTemplate(w, r, tenantID, id, "update tenant offer template")
}

func (h *Handler) DeleteTenantOfferLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "offerTemplateID", "delete tenant offer template")
	if !ok {
		return
	}
	if err := h.svc.DeleteOfferLetterTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant offer template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantOfferLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant offer")
	if !ok {
		return
	}
	h.createOfferLetter(w, r, tenantID, "create tenant offer")
}

func (h *Handler) ListTenantOfferLetters(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant offers")
	if !ok {
		return
	}
	h.listOfferLetters(w, r, tenantID, "list tenant offers")
}

func (h *Handler) GetTenantOfferLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "offerLetterID", "get tenant offer")
	if !ok {
		return
	}
	item, err := h.svc.GetOfferLetter(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant offer", err, "offer letter not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantOfferLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "offerLetterID", "update tenant offer")
	if !ok {
		return
	}
	h.updateOfferLetter(w, r, tenantID, id, "update tenant offer")
}

func (h *Handler) UpdateTenantOfferLetterStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "offerLetterID", "update tenant offer status")
	if !ok {
		return
	}
	h.updateOfferLetterStatus(w, r, tenantID, id, "update tenant offer status")
}

func (h *Handler) ListTenantOfferLetterEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "offerLetterID", "list tenant offer events")
	if !ok {
		return
	}
	items, err := h.svc.ListOfferLetterEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant offer events", err, "failed to list offer events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteTenantOfferLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "offerLetterID", "delete tenant offer")
	if !ok {
		return
	}
	if err := h.svc.DeleteOfferLetter(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant offer", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createOfferLetterTemplate(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.OfferLetterTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateOfferLetterTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listOfferLetterTemplates(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListOfferLetterTemplates(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list offer templates")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateOfferLetterTemplate(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.OfferLetterTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateOfferLetterTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createOfferLetter(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.OfferLetterCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateOfferLetter(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listOfferLetters(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	page, err := h.svc.ListOfferLetters(r.Context(), domain.OfferLetterFilter{TenantID: tenantID, ApplicationID: optionalUUIDQuery(r, "application_id"), Status: optionalStringQuery(r, "status"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list offers")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) updateOfferLetter(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.OfferLetterCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateOfferLetter(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) updateOfferLetterStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.OfferLetterStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateOfferLetterStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) offerTemplateRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "offerTemplateID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid offer template id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) offerLetterRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "offerLetterID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid offer letter id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
