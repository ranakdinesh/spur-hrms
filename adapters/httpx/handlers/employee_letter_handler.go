package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) GetPublicEmployeeLetter(w http.ResponseWriter, r *http.Request) {
	item, err := h.svc.GetEmployeeLetterBySignatureToken(r.Context(), chi.URLParam(r, "signatureToken"))
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get public employee letter", err, "employee letter not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) SignPublicEmployeeLetter(w http.ResponseWriter, r *http.Request) {
	var cmd ports.EmployeeLetterSignatureCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "sign public employee letter", err, "invalid request body")
		return
	}
	cmd.Token = chi.URLParam(r, "signatureToken")
	ip := clientIP(r)
	ua := r.UserAgent()
	cmd.IPAddress = &ip
	cmd.UserAgent = &ua
	item, err := h.svc.SignEmployeeLetter(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "sign public employee letter", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) CreateEmployeeLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create employee letter template", err, "tenant context is required")
		return
	}
	h.createEmployeeLetterTemplate(w, r, tenantID, "create employee letter template")
}

func (h *Handler) ListEmployeeLetterTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list employee letter templates", err, "tenant context is required")
		return
	}
	h.listEmployeeLetterTemplates(w, r, tenantID, "list employee letter templates")
}

func (h *Handler) UpdateEmployeeLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.employeeLetterTemplateRequestIDs(w, r, "update employee letter template")
	if !ok {
		return
	}
	h.updateEmployeeLetterTemplate(w, r, tenantID, id, "update employee letter template")
}

func (h *Handler) DeleteEmployeeLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.employeeLetterTemplateRequestIDs(w, r, "delete employee letter template")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeLetterTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete employee letter template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GenerateEmployeeLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "generate employee letter", err, "tenant context is required")
		return
	}
	h.generateEmployeeLetter(w, r, tenantID, "generate employee letter")
}

func (h *Handler) ListEmployeeLetters(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list employee letters", err, "tenant context is required")
		return
	}
	h.listEmployeeLetters(w, r, tenantID, "list employee letters")
}

func (h *Handler) GetEmployeeLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.employeeLetterRequestIDs(w, r, "get employee letter")
	if !ok {
		return
	}
	item, err := h.svc.GetEmployeeLetter(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get employee letter", err, "employee letter not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateEmployeeLetterStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.employeeLetterRequestIDs(w, r, "update employee letter status")
	if !ok {
		return
	}
	h.updateEmployeeLetterStatus(w, r, tenantID, id, "update employee letter status")
}

func (h *Handler) DownloadEmployeeLetterPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.employeeLetterRequestIDs(w, r, "download employee letter pdf")
	if !ok {
		return
	}
	h.downloadEmployeeLetterPDF(w, r, tenantID, id, "download employee letter pdf")
}

func (h *Handler) ListEmployeeLetterEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.employeeLetterRequestIDs(w, r, "list employee letter events")
	if !ok {
		return
	}
	items, err := h.svc.ListEmployeeLetterEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list employee letter events", err, "failed to list employee letter events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteEmployeeLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.employeeLetterRequestIDs(w, r, "delete employee letter")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeLetter(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete employee letter", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantEmployeeLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant employee letter template")
	if !ok {
		return
	}
	h.createEmployeeLetterTemplate(w, r, tenantID, "create tenant employee letter template")
}

func (h *Handler) ListTenantEmployeeLetterTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant employee letter templates")
	if !ok {
		return
	}
	h.listEmployeeLetterTemplates(w, r, tenantID, "list tenant employee letter templates")
}

func (h *Handler) UpdateTenantEmployeeLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employeeLetterTemplateID", "update tenant employee letter template")
	if !ok {
		return
	}
	h.updateEmployeeLetterTemplate(w, r, tenantID, id, "update tenant employee letter template")
}

func (h *Handler) DeleteTenantEmployeeLetterTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employeeLetterTemplateID", "delete tenant employee letter template")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeLetterTemplate(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant employee letter template", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GenerateTenantEmployeeLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "generate tenant employee letter")
	if !ok {
		return
	}
	h.generateEmployeeLetter(w, r, tenantID, "generate tenant employee letter")
}

func (h *Handler) ListTenantEmployeeLetters(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant employee letters")
	if !ok {
		return
	}
	h.listEmployeeLetters(w, r, tenantID, "list tenant employee letters")
}

func (h *Handler) GetTenantEmployeeLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employeeLetterID", "get tenant employee letter")
	if !ok {
		return
	}
	item, err := h.svc.GetEmployeeLetter(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant employee letter", err, "employee letter not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantEmployeeLetterStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employeeLetterID", "update tenant employee letter status")
	if !ok {
		return
	}
	h.updateEmployeeLetterStatus(w, r, tenantID, id, "update tenant employee letter status")
}

func (h *Handler) DownloadTenantEmployeeLetterPDF(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employeeLetterID", "download tenant employee letter pdf")
	if !ok {
		return
	}
	h.downloadEmployeeLetterPDF(w, r, tenantID, id, "download tenant employee letter pdf")
}

func (h *Handler) ListTenantEmployeeLetterEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employeeLetterID", "list tenant employee letter events")
	if !ok {
		return
	}
	items, err := h.svc.ListEmployeeLetterEvents(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant employee letter events", err, "failed to list employee letter events")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) DeleteTenantEmployeeLetter(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "employeeLetterID", "delete tenant employee letter")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeLetter(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant employee letter", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createEmployeeLetterTemplate(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EmployeeLetterTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateEmployeeLetterTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listEmployeeLetterTemplates(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListEmployeeLetterTemplates(r.Context(), tenantID, optionalStringQuery(r, "letter_type"))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list employee letter templates")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateEmployeeLetterTemplate(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.EmployeeLetterTemplateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEmployeeLetterTemplate(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) generateEmployeeLetter(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.EmployeeLetterCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.GenerateEmployeeLetter(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listEmployeeLetters(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	page, err := h.svc.ListEmployeeLetters(r.Context(), domain.EmployeeLetterFilter{TenantID: tenantID, EmployeeID: optionalUUIDQuery(r, "employee_id"), LetterType: optionalStringQuery(r, "letter_type"), Status: optionalStringQuery(r, "status"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 100), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list employee letters")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) updateEmployeeLetterStatus(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.EmployeeLetterStatusCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEmployeeLetterStatus(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) downloadEmployeeLetterPDF(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	data, name, err := h.svc.RenderEmployeeLetterPDF(r.Context(), tenantID, id, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondDownload(w, "application/pdf", name, data)
}

func (h *Handler) employeeLetterTemplateRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "employeeLetterTemplateID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid employee letter template id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) employeeLetterRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "employeeLetterID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid employee letter id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
