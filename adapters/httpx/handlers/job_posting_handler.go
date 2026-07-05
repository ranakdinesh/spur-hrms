package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create job posting", err, "tenant context is required")
		return
	}
	h.createJobPosting(w, r, tenantID, "create job posting")
}

func (h *Handler) ListJobPostings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list job postings", err, "tenant context is required")
		return
	}
	h.listJobPostings(w, r, tenantID, "list job postings")
}

func (h *Handler) ListPublishedJobPostings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list published job postings", err, "tenant context is required")
		return
	}
	h.listPublishedJobPostings(w, r, tenantID, "list published job postings")
}

func (h *Handler) GetJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobPostingRequestIDs(w, r, "get job posting")
	if !ok {
		return
	}
	item, err := h.svc.GetJobPosting(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get job posting", err, "job posting not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobPostingRequestIDs(w, r, "update job posting")
	if !ok {
		return
	}
	h.updateJobPosting(w, r, tenantID, id, "update job posting")
}

func (h *Handler) PublishJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobPostingRequestIDs(w, r, "publish job posting")
	if !ok {
		return
	}
	h.publishJobPosting(w, r, tenantID, id, "publish job posting")
}

func (h *Handler) CloseJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobPostingRequestIDs(w, r, "close job posting")
	if !ok {
		return
	}
	item, err := h.svc.CloseJobPosting(r.Context(), tenantID, id, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "close job posting", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ExpireJobPostings(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "expire job postings", err, "tenant context is required")
		return
	}
	h.expireJobPostings(w, r, tenantID, "expire job postings")
}

func (h *Handler) DeleteJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.jobPostingRequestIDs(w, r, "delete job posting")
	if !ok {
		return
	}
	if err := h.svc.DeleteJobPosting(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete job posting", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTenantJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant job posting")
	if !ok {
		return
	}
	h.createJobPosting(w, r, tenantID, "create tenant job posting")
}

func (h *Handler) ListTenantJobPostings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant job postings")
	if !ok {
		return
	}
	h.listJobPostings(w, r, tenantID, "list tenant job postings")
}

func (h *Handler) ListTenantPublishedJobPostings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant published job postings")
	if !ok {
		return
	}
	h.listPublishedJobPostings(w, r, tenantID, "list tenant published job postings")
}

func (h *Handler) GetTenantJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobPostingID", "get tenant job posting")
	if !ok {
		return
	}
	item, err := h.svc.GetJobPosting(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant job posting", err, "job posting not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobPostingID", "update tenant job posting")
	if !ok {
		return
	}
	h.updateJobPosting(w, r, tenantID, id, "update tenant job posting")
}

func (h *Handler) PublishTenantJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobPostingID", "publish tenant job posting")
	if !ok {
		return
	}
	h.publishJobPosting(w, r, tenantID, id, "publish tenant job posting")
}

func (h *Handler) CloseTenantJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobPostingID", "close tenant job posting")
	if !ok {
		return
	}
	item, err := h.svc.CloseJobPosting(r.Context(), tenantID, id, h.actorIDFromRequest(r))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "close tenant job posting", err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) ExpireTenantJobPostings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "expire tenant job postings")
	if !ok {
		return
	}
	h.expireJobPostings(w, r, tenantID, "expire tenant job postings")
}

func (h *Handler) DeleteTenantJobPosting(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminLookupRequestIDs(w, r, "jobPostingID", "delete tenant job posting")
	if !ok {
		return
	}
	if err := h.svc.DeleteJobPosting(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant job posting", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) createJobPosting(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.JobPostingCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateJobPosting(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateJobPosting(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.JobPostingCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateJobPosting(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listJobPostings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	page, err := h.svc.ListJobPostings(r.Context(), domain.JobPostingFilter{TenantID: tenantID, JobStatus: optionalStringQuery(r, "job_status"), IsPublished: optionalBoolQuery(r, "is_published"), DepartmentID: optionalUUIDQuery(r, "department_id"), Search: optionalStringQuery(r, "search"), Limit: queryInt32(r, "limit", 25), Offset: queryInt32(r, "offset", 0)})
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list job postings")
		return
	}
	respondJSON(w, http.StatusOK, page)
}

func (h *Handler) listPublishedJobPostings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListPublishedJobPostings(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list published job postings")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) publishJobPosting(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.JobPostingPublishCommand
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, http.ErrBodyNotAllowed) {
			h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
			return
		}
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.PublishJobPosting(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) expireJobPostings(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ExpireJobPostings(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to expire job postings")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) jobPostingRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "jobPostingID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid job posting id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}
