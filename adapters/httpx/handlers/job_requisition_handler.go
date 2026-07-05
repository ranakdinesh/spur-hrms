package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create job requisition", err, "tenant context is required")
		return
	}
	h.createJobRequisition(w, r, tenantID, "create job requisition")
}

func (h *Handler) ListJobRequisitions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list job requisitions", err, "tenant context is required")
		return
	}
	h.listJobRequisitions(w, r, tenantID, "list job requisitions")
}

func (h *Handler) GetJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobRequisitionRequestIDs(w, r, "get job requisition")
	if !ok {
		return
	}
	item, err := h.svc.GetJobRequisition(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get job requisition", err, "job requisition not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobRequisitionRequestIDs(w, r, "update job requisition")
	if !ok {
		return
	}
	h.updateJobRequisition(w, r, tenantID, id, "update job requisition")
}

func (h *Handler) DeleteJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobRequisitionRequestIDs(w, r, "delete job requisition")
	if !ok {
		return
	}
	if err := h.svc.DeleteJobRequisition(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete job requisition", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SubmitJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobRequisitionRequestIDs(w, r, "submit job requisition")
	if !ok {
		return
	}
	h.requisitionAction(w, r, ports.JobRequisitionActionCommand{TenantID: tenantID, ID: id, ActorID: h.actorIDFromRequest(r)}, h.svc.SubmitJobRequisition, "submit job requisition")
}

func (h *Handler) ApproveJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobRequisitionRequestIDs(w, r, "approve job requisition")
	if !ok {
		return
	}
	h.requisitionAction(w, r, ports.JobRequisitionActionCommand{TenantID: tenantID, ID: id, ActorID: h.actorIDFromRequest(r)}, h.svc.ApproveJobRequisition, "approve job requisition")
}

func (h *Handler) RejectJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobRequisitionRequestIDs(w, r, "reject job requisition")
	if !ok {
		return
	}
	h.requisitionAction(w, r, ports.JobRequisitionActionCommand{TenantID: tenantID, ID: id, ActorID: h.actorIDFromRequest(r)}, h.svc.RejectJobRequisition, "reject job requisition")
}

func (h *Handler) CloseJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobRequisitionRequestIDs(w, r, "close job requisition")
	if !ok {
		return
	}
	h.requisitionAction(w, r, ports.JobRequisitionActionCommand{TenantID: tenantID, ID: id, ActorID: h.actorIDFromRequest(r)}, h.svc.CloseJobRequisition, "close job requisition")
}

func (h *Handler) ListJobRequisitionLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobRequisitionRequestIDs(w, r, "list job requisition logs")
	if !ok {
		return
	}
	items, err := h.svc.ListJobRequisitionLogs(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list job requisition logs", err, "failed to list requisition logs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) CreateTenantJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant job requisition")
	if !ok {
		return
	}
	h.createJobRequisition(w, r, tenantID, "create tenant job requisition")
}

func (h *Handler) ListTenantJobRequisitions(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant job requisitions")
	if !ok {
		return
	}
	h.listJobRequisitions(w, r, tenantID, "list tenant job requisitions")
}

func (h *Handler) GetTenantJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobRequisitionID", "get tenant job requisition")
	if !ok {
		return
	}
	item, err := h.svc.GetJobRequisition(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant job requisition", err, "job requisition not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobRequisitionID", "update tenant job requisition")
	if !ok {
		return
	}
	h.updateJobRequisition(w, r, tenantID, id, "update tenant job requisition")
}

func (h *Handler) DeleteTenantJobRequisition(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobRequisitionID", "delete tenant job requisition")
	if !ok {
		return
	}
	if err := h.svc.DeleteJobRequisition(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant job requisition", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SubmitTenantJobRequisition(w http.ResponseWriter, r *http.Request) {
	h.tenantRequisitionAction(w, r, h.svc.SubmitJobRequisition, "submit tenant job requisition")
}

func (h *Handler) ApproveTenantJobRequisition(w http.ResponseWriter, r *http.Request) {
	h.tenantRequisitionAction(w, r, h.svc.ApproveJobRequisition, "approve tenant job requisition")
}

func (h *Handler) RejectTenantJobRequisition(w http.ResponseWriter, r *http.Request) {
	h.tenantRequisitionAction(w, r, h.svc.RejectJobRequisition, "reject tenant job requisition")
}

func (h *Handler) CloseTenantJobRequisition(w http.ResponseWriter, r *http.Request) {
	h.tenantRequisitionAction(w, r, h.svc.CloseJobRequisition, "close tenant job requisition")
}

func (h *Handler) ListTenantJobRequisitionLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobRequisitionID", "list tenant job requisition logs")
	if !ok {
		return
	}
	items, err := h.svc.ListJobRequisitionLogs(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant job requisition logs", err, "failed to list requisition logs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createJobRequisition(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.JobRequisitionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateJobRequisition(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateJobRequisition(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.JobRequisitionCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateJobRequisition(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listJobRequisitions(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	page, err := h.svc.ListJobRequisitions(r.Context(), domain.JobRequisitionFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), JobPositionID: optionalUUIDQuery(r, "job_position_id"), DepartmentID: optionalUUIDQuery(r, "department_id"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 25), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list job requisitions")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) requisitionAction(w http.ResponseWriter, r *http.Request, cmd ports.JobRequisitionActionCommand, fn func(context.Context, ports.JobRequisitionActionCommand) (*domain.JobRequisition, error), operation string) {
	var body struct {
		Remarks *string `json:"remarks,omitempty"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body)
	}
	cmd.Remarks = body.Remarks
	item, err := fn(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) tenantRequisitionAction(w http.ResponseWriter, r *http.Request, fn func(context.Context, ports.JobRequisitionActionCommand) (*domain.JobRequisition, error), operation string) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobRequisitionID", operation)
	if !ok {
		return
	}
	h.requisitionAction(w, r, ports.JobRequisitionActionCommand{TenantID: tenantID, ID: id, ActorID: h.actorIDFromRequest(r)}, fn, operation)
}

func (h *Handler) jobRequisitionRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "jobRequisitionID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid job requisition id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
