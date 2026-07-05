package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateDocumentType(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create document type", err, "tenant context is required")
		return
	}
	h.createDocumentTypeForTenant(w, r, tenantID, "create document type")
}

func (h *Handler) ListDocumentTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list document types", err, "tenant context is required")
		return
	}
	h.listDocumentTypesForTenant(w, r, tenantID, "list document types")
}

func (h *Handler) GetDocumentType(w http.ResponseWriter, r *http.Request) {
	tenantID, documentTypeID, ok := h.documentTypeRequestIDs(w, r, "get document type")
	if !ok {
		return
	}
	item, err := h.svc.GetDocumentType(r.Context(), tenantID, documentTypeID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get document type", err, "document type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateDocumentType(w http.ResponseWriter, r *http.Request) {
	tenantID, documentTypeID, ok := h.documentTypeRequestIDs(w, r, "update document type")
	if !ok {
		return
	}
	h.updateDocumentTypeForTenant(w, r, tenantID, documentTypeID, "update document type")
}

func (h *Handler) DeleteDocumentType(w http.ResponseWriter, r *http.Request) {
	tenantID, documentTypeID, ok := h.documentTypeRequestIDs(w, r, "delete document type")
	if !ok {
		return
	}
	if err := h.svc.DeleteDocumentType(r.Context(), tenantID, documentTypeID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete document type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListEmployeeDocuments(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.employeeDocumentRequestIDs(w, r, "list employee documents")
	if !ok {
		return
	}
	profile, err := h.svc.GetEmployeeProfile(r.Context(), tenantID, employeeID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "list employee documents", err, "employee profile not found")
		return
	}
	respondJSON(w, http.StatusOK, profile.Documents)
}

func (h *Handler) CreateEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.employeeDocumentRequestIDs(w, r, "create employee document")
	if !ok {
		return
	}
	h.createEmployeeDocumentForTenant(w, r, tenantID, employeeID, "create employee document")
}

func (h *Handler) CreateMyEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, profile, ok := h.myEmployeeProfileRequest(w, r, "create my employee document")
	if !ok {
		return
	}
	h.createEmployeeDocumentForTenant(w, r, tenantID, profile.Employee.ID, "create my employee document")
}

func (h *Handler) UpdateEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, documentID, ok := h.employeeDocumentItemRequestIDs(w, r, "update employee document")
	if !ok {
		return
	}
	h.updateEmployeeDocumentForTenant(w, r, tenantID, employeeID, documentID, "update employee document")
}

func (h *Handler) UpdateMyEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, profile, ok := h.myEmployeeProfileRequest(w, r, "update my employee document")
	if !ok {
		return
	}
	documentID, err := uuid.Parse(chi.URLParam(r, "documentID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse my employee document id", err, "invalid document id")
		return
	}
	h.updateEmployeeDocumentForTenant(w, r, tenantID, profile.Employee.ID, documentID, "update my employee document")
}

func (h *Handler) ReviewEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, documentID, ok := h.employeeDocumentItemRequestIDs(w, r, "review employee document")
	if !ok {
		return
	}
	h.reviewEmployeeDocumentForTenant(w, r, tenantID, employeeID, documentID, "review employee document")
}

func (h *Handler) DeleteEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, documentID, ok := h.employeeDocumentItemRequestIDs(w, r, "delete employee document")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeDocument(r.Context(), tenantID, employeeID, documentID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete employee document", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantDocumentType(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant document type")
	if !ok {
		return
	}
	h.createDocumentTypeForTenant(w, r, tenantID, "create tenant document type")
}

func (h *Handler) ListTenantDocumentTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant document types")
	if !ok {
		return
	}
	h.listDocumentTypesForTenant(w, r, tenantID, "list tenant document types")
}

func (h *Handler) GetTenantDocumentType(w http.ResponseWriter, r *http.Request) {
	tenantID, documentTypeID, ok := h.superAdminDocumentTypeRequestIDs(w, r, "get tenant document type")
	if !ok {
		return
	}
	item, err := h.svc.GetDocumentType(r.Context(), tenantID, documentTypeID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant document type", err, "document type not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantDocumentType(w http.ResponseWriter, r *http.Request) {
	tenantID, documentTypeID, ok := h.superAdminDocumentTypeRequestIDs(w, r, "update tenant document type")
	if !ok {
		return
	}
	h.updateDocumentTypeForTenant(w, r, tenantID, documentTypeID, "update tenant document type")
}

func (h *Handler) DeleteTenantDocumentType(w http.ResponseWriter, r *http.Request) {
	tenantID, documentTypeID, ok := h.superAdminDocumentTypeRequestIDs(w, r, "delete tenant document type")
	if !ok {
		return
	}
	if err := h.svc.DeleteDocumentType(r.Context(), tenantID, documentTypeID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant document type", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTenantEmployeeDocuments(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.superAdminEmployeeDocumentRequestIDs(w, r, "list tenant employee documents")
	if !ok {
		return
	}
	profile, err := h.svc.GetEmployeeProfile(r.Context(), tenantID, employeeID, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "list tenant employee documents", err, "employee profile not found")
		return
	}
	respondJSON(w, http.StatusOK, profile.Documents)
}

func (h *Handler) CreateTenantEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, ok := h.superAdminEmployeeDocumentRequestIDs(w, r, "create tenant employee document")
	if !ok {
		return
	}
	h.createEmployeeDocumentForTenant(w, r, tenantID, employeeID, "create tenant employee document")
}

func (h *Handler) UpdateTenantEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, documentID, ok := h.superAdminEmployeeDocumentItemRequestIDs(w, r, "update tenant employee document")
	if !ok {
		return
	}
	h.updateEmployeeDocumentForTenant(w, r, tenantID, employeeID, documentID, "update tenant employee document")
}

func (h *Handler) ReviewTenantEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, documentID, ok := h.superAdminEmployeeDocumentItemRequestIDs(w, r, "review tenant employee document")
	if !ok {
		return
	}
	h.reviewEmployeeDocumentForTenant(w, r, tenantID, employeeID, documentID, "review tenant employee document")
}

func (h *Handler) DeleteTenantEmployeeDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, employeeID, documentID, ok := h.superAdminEmployeeDocumentItemRequestIDs(w, r, "delete tenant employee document")
	if !ok {
		return
	}
	if err := h.svc.DeleteEmployeeDocument(r.Context(), tenantID, employeeID, documentID, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant employee document", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createDocumentTypeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.DocumentTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateDocumentType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listDocumentTypesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListDocumentTypes(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list document types")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) updateDocumentTypeForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, documentTypeID uuid.UUID, operation string) {
	var cmd ports.DocumentTypeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = documentTypeID
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateDocumentType(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createEmployeeDocumentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, employeeID uuid.UUID, operation string) {
	var cmd ports.EmployeeDocumentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.EmployeeID = employeeID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateEmployeeDocument(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateEmployeeDocumentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, employeeID uuid.UUID, documentID uuid.UUID, operation string) {
	var cmd ports.EmployeeDocumentCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.ID = documentID
	cmd.TenantID = tenantID
	cmd.EmployeeID = employeeID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateEmployeeDocument(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) reviewEmployeeDocumentForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, employeeID uuid.UUID, documentID uuid.UUID, operation string) {
	var cmd ports.EmployeeDocumentReviewCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation+" request", err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.EmployeeID = employeeID
	cmd.DocumentID = documentID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ReviewEmployeeDocument(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) documentTypeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	documentTypeID, err := uuid.Parse(chi.URLParam(r, "documentTypeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse document type id", err, "invalid document type id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, documentTypeID, true
}

func (h *Handler) superAdminDocumentTypeRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	documentTypeID, err := uuid.Parse(chi.URLParam(r, "documentTypeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant document type id", err, "invalid document type id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, documentTypeID, true
}

func (h *Handler) employeeDocumentRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse employee id", err, "invalid employee id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, employeeID, true
}

func (h *Handler) myEmployeeProfileRequest(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, *domain.EmployeeProfile, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, nil, false
	}
	actorID := h.actorIDFromRequest(r)
	if actorID == nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, domain.ErrInvalidEmployeeUserID, "user context is required")
		return uuid.Nil, nil, false
	}
	profile, err := h.svc.GetEmployeeSelfProfile(r.Context(), tenantID, *actorID, actorID)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, "employee profile not found")
		return uuid.Nil, nil, false
	}
	return tenantID, profile, true
}

func (h *Handler) employeeDocumentItemRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, uuid.UUID, bool) {
	tenantID, employeeID, ok := h.employeeDocumentRequestIDs(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	documentID, err := uuid.Parse(chi.URLParam(r, "documentID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse employee document id", err, "invalid employee document id")
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	return tenantID, employeeID, documentID, true
}

func (h *Handler) superAdminEmployeeDocumentRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	employeeID, err := uuid.Parse(chi.URLParam(r, "employeeID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant employee id", err, "invalid employee id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, employeeID, true
}

func (h *Handler) superAdminEmployeeDocumentItemRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, uuid.UUID, bool) {
	tenantID, employeeID, ok := h.superAdminEmployeeDocumentRequestIDs(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	documentID, err := uuid.Parse(chi.URLParam(r, "documentID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse tenant employee document id", err, "invalid employee document id")
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	return tenantID, employeeID, documentID, true
}
