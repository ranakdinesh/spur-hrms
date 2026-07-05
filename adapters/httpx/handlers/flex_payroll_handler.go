package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) CreateFlexPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create flex pay run", err, "tenant context is required")
		return
	}
	h.createFlexPayRunForTenant(w, r, tenantID, "create flex pay run")
}

func (h *Handler) ListFlexPayRuns(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list flex pay runs", err, "tenant context is required")
		return
	}
	h.listFlexPayRunsForTenant(w, r, tenantID, "list flex pay runs")
}

func (h *Handler) GetFlexPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.flexPayRunRequestIDs(w, r, "get flex pay run")
	if !ok {
		return
	}
	h.getFlexPayRunForTenant(w, r, tenantID, id, "get flex pay run")
}

func (h *Handler) DeleteFlexPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.flexPayRunRequestIDs(w, r, "delete flex pay run")
	if !ok {
		return
	}
	if err := h.svc.DeleteFlexPayRun(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete flex pay run", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GenerateFlexPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.flexPayRunRequestIDs(w, r, "generate flex pay run")
	if !ok {
		return
	}
	h.generateFlexPayRunForTenant(w, r, tenantID, id, "generate flex pay run")
}

func (h *Handler) SubmitFlexPayRun(w http.ResponseWriter, r *http.Request) {
	h.flexPayRunAction(w, r, "submit flex pay run", h.svc.SubmitFlexPayRun)
}

func (h *Handler) ApproveFlexPayRun(w http.ResponseWriter, r *http.Request) {
	h.flexPayRunAction(w, r, "approve flex pay run", h.svc.ApproveFlexPayRun)
}

func (h *Handler) RejectFlexPayRun(w http.ResponseWriter, r *http.Request) {
	h.flexPayRunAction(w, r, "reject flex pay run", h.svc.RejectFlexPayRun)
}

func (h *Handler) MarkFlexPayRunPaymentPending(w http.ResponseWriter, r *http.Request) {
	h.flexPayRunAction(w, r, "mark flex pay run payment pending", h.svc.MarkFlexPayRunPaymentPending)
}

func (h *Handler) MarkFlexPayRunPaid(w http.ResponseWriter, r *http.Request) {
	h.flexPayRunAction(w, r, "mark flex pay run paid", h.svc.MarkFlexPayRunPaid)
}

func (h *Handler) ExportFlexPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.flexPayRunRequestIDs(w, r, "export flex pay run")
	if !ok {
		return
	}
	h.exportFlexPayRunForTenant(w, r, tenantID, id, "export flex pay run")
}

func (h *Handler) CreateContractorInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create contractor invoice", err, "tenant context is required")
		return
	}
	h.createContractorInvoiceForTenant(w, r, tenantID, "create contractor invoice")
}

func (h *Handler) ListContractorInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list contractor invoices", err, "tenant context is required")
		return
	}
	h.listContractorInvoicesForTenant(w, r, tenantID, "list contractor invoices")
}

func (h *Handler) GetContractorInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.invoiceRequestIDs(w, r, "get contractor invoice")
	if !ok {
		return
	}
	item, err := h.svc.GetContractorInvoice(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get contractor invoice", err, "contractor invoice not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateContractorInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.invoiceRequestIDs(w, r, "update contractor invoice")
	if !ok {
		return
	}
	h.updateContractorInvoiceForTenant(w, r, tenantID, id, "update contractor invoice")
}

func (h *Handler) DeleteContractorInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.invoiceRequestIDs(w, r, "delete contractor invoice")
	if !ok {
		return
	}
	if err := h.svc.DeleteContractorInvoice(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete contractor invoice", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SubmitContractorInvoice(w http.ResponseWriter, r *http.Request) {
	h.contractorInvoiceAction(w, r, "submit contractor invoice", h.svc.SubmitContractorInvoice)
}

func (h *Handler) ApproveContractorInvoice(w http.ResponseWriter, r *http.Request) {
	h.contractorInvoiceAction(w, r, "approve contractor invoice", h.svc.ApproveContractorInvoice)
}

func (h *Handler) RejectContractorInvoice(w http.ResponseWriter, r *http.Request) {
	h.contractorInvoiceAction(w, r, "reject contractor invoice", h.svc.RejectContractorInvoice)
}

func (h *Handler) MarkContractorInvoicePaid(w http.ResponseWriter, r *http.Request) {
	h.contractorInvoiceAction(w, r, "mark contractor invoice paid", h.svc.MarkContractorInvoicePaid)
}

func (h *Handler) CreateFlexPayRunItem(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.flexPayRunRequestIDs(w, r, "create flex pay run item")
	if !ok {
		return
	}
	h.createFlexPayRunItemForTenant(w, r, tenantID, id, "create flex pay run item")
}

func (h *Handler) ListFlexPayRunItems(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.flexPayRunRequestIDs(w, r, "list flex pay run items")
	if !ok {
		return
	}
	items, err := h.svc.ListFlexPayRunItems(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list flex pay run items", err, "failed to list flexible pay run items")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListFlexPayRunEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.flexPayRunRequestIDs(w, r, "list flex pay run events")
	if !ok {
		return
	}
	events, err := h.svc.ListFlexPayRunEvents(r.Context(), tenantID, &id, nil)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list flex pay run events", err, "failed to list flexible pay run events")
		return
	}
	respondJSON(w, http.StatusOK, events)
}

func (h *Handler) CreateTenantFlexPayRun(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant flex pay run"); ok {
		h.createFlexPayRunForTenant(w, r, tenantID, "create tenant flex pay run")
	}
}

func (h *Handler) ListTenantFlexPayRuns(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant flex pay runs"); ok {
		h.listFlexPayRunsForTenant(w, r, tenantID, "list tenant flex pay runs")
	}
}

func (h *Handler) GetTenantFlexPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFlexPayRunRequestIDs(w, r, "get tenant flex pay run")
	if ok {
		h.getFlexPayRunForTenant(w, r, tenantID, id, "get tenant flex pay run")
	}
}

func (h *Handler) DeleteTenantFlexPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFlexPayRunRequestIDs(w, r, "delete tenant flex pay run")
	if !ok {
		return
	}
	if err := h.svc.DeleteFlexPayRun(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant flex pay run", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GenerateTenantFlexPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFlexPayRunRequestIDs(w, r, "generate tenant flex pay run")
	if ok {
		h.generateFlexPayRunForTenant(w, r, tenantID, id, "generate tenant flex pay run")
	}
}

func (h *Handler) ExportTenantFlexPayRun(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFlexPayRunRequestIDs(w, r, "export tenant flex pay run")
	if ok {
		h.exportFlexPayRunForTenant(w, r, tenantID, id, "export tenant flex pay run")
	}
}

func (h *Handler) SubmitTenantFlexPayRun(w http.ResponseWriter, r *http.Request) {
	h.tenantFlexPayRunAction(w, r, "submit tenant flex pay run", h.svc.SubmitFlexPayRun)
}

func (h *Handler) ApproveTenantFlexPayRun(w http.ResponseWriter, r *http.Request) {
	h.tenantFlexPayRunAction(w, r, "approve tenant flex pay run", h.svc.ApproveFlexPayRun)
}

func (h *Handler) RejectTenantFlexPayRun(w http.ResponseWriter, r *http.Request) {
	h.tenantFlexPayRunAction(w, r, "reject tenant flex pay run", h.svc.RejectFlexPayRun)
}

func (h *Handler) MarkTenantFlexPayRunPaymentPending(w http.ResponseWriter, r *http.Request) {
	h.tenantFlexPayRunAction(w, r, "mark tenant flex pay run payment pending", h.svc.MarkFlexPayRunPaymentPending)
}

func (h *Handler) MarkTenantFlexPayRunPaid(w http.ResponseWriter, r *http.Request) {
	h.tenantFlexPayRunAction(w, r, "mark tenant flex pay run paid", h.svc.MarkFlexPayRunPaid)
}

func (h *Handler) CreateTenantContractorInvoice(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "create tenant contractor invoice"); ok {
		h.createContractorInvoiceForTenant(w, r, tenantID, "create tenant contractor invoice")
	}
}

func (h *Handler) ListTenantContractorInvoices(w http.ResponseWriter, r *http.Request) {
	if tenantID, ok := h.superAdminTenantID(w, r, "list tenant contractor invoices"); ok {
		h.listContractorInvoicesForTenant(w, r, tenantID, "list tenant contractor invoices")
	}
}

func (h *Handler) GetTenantContractorInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminInvoiceRequestIDs(w, r, "get tenant contractor invoice")
	if !ok {
		return
	}
	item, err := h.svc.GetContractorInvoice(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, "get tenant contractor invoice", err, "contractor invoice not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) UpdateTenantContractorInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminInvoiceRequestIDs(w, r, "update tenant contractor invoice")
	if ok {
		h.updateContractorInvoiceForTenant(w, r, tenantID, id, "update tenant contractor invoice")
	}
}

func (h *Handler) DeleteTenantContractorInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminInvoiceRequestIDs(w, r, "delete tenant contractor invoice")
	if !ok {
		return
	}
	if err := h.svc.DeleteContractorInvoice(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "delete tenant contractor invoice", err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SubmitTenantContractorInvoice(w http.ResponseWriter, r *http.Request) {
	h.tenantContractorInvoiceAction(w, r, "submit tenant contractor invoice", h.svc.SubmitContractorInvoice)
}

func (h *Handler) ApproveTenantContractorInvoice(w http.ResponseWriter, r *http.Request) {
	h.tenantContractorInvoiceAction(w, r, "approve tenant contractor invoice", h.svc.ApproveContractorInvoice)
}

func (h *Handler) RejectTenantContractorInvoice(w http.ResponseWriter, r *http.Request) {
	h.tenantContractorInvoiceAction(w, r, "reject tenant contractor invoice", h.svc.RejectContractorInvoice)
}

func (h *Handler) MarkTenantContractorInvoicePaid(w http.ResponseWriter, r *http.Request) {
	h.tenantContractorInvoiceAction(w, r, "mark tenant contractor invoice paid", h.svc.MarkContractorInvoicePaid)
}

func (h *Handler) CreateTenantFlexPayRunItem(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFlexPayRunRequestIDs(w, r, "create tenant flex pay run item")
	if ok {
		h.createFlexPayRunItemForTenant(w, r, tenantID, id, "create tenant flex pay run item")
	}
}

func (h *Handler) ListTenantFlexPayRunItems(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFlexPayRunRequestIDs(w, r, "list tenant flex pay run items")
	if !ok {
		return
	}
	items, err := h.svc.ListFlexPayRunItems(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant flex pay run items", err, "failed to list flexible pay run items")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) ListTenantFlexPayRunEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, id, ok := h.superAdminFlexPayRunRequestIDs(w, r, "list tenant flex pay run events")
	if !ok {
		return
	}
	events, err := h.svc.ListFlexPayRunEvents(r.Context(), tenantID, &id, nil)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "list tenant flex pay run events", err, "failed to list flexible pay run events")
		return
	}
	respondJSON(w, http.StatusOK, events)
}

func (h *Handler) createFlexPayRunForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.FlexPayRunCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateFlexPayRun(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listFlexPayRunsForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	dateFrom, ok := optionalFlexDateQuery(w, r, "date_from", operation, h)
	if !ok {
		return
	}
	dateTo, ok := optionalFlexDateQuery(w, r, "date_to", operation, h)
	if !ok {
		return
	}
	filter := domain.FlexPayRunFilter{TenantID: tenantID, Status: optionalStringQuery(r, "status"), RunType: optionalStringQuery(r, "run_type"), DateFrom: dateFrom, DateTo: dateTo, Search: optionalStringQuery(r, "search")}
	items, err := h.svc.ListFlexPayRuns(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list flexible pay runs")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getFlexPayRunForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	item, err := h.svc.GetFlexPayRun(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, "flexible pay run not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) generateFlexPayRunForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.FlexPayRunGenerateCommand
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&cmd)
	}
	cmd.TenantID = tenantID
	cmd.FlexPayRunID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.GenerateFlexPayRun(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) flexPayRunAction(w http.ResponseWriter, r *http.Request, operation string, fn func(context.Context, ports.FlexPayRunActionCommand) (*domain.FlexPayRun, error)) {
	tenantID, id, ok := h.flexPayRunRequestIDs(w, r, operation)
	if !ok {
		return
	}
	h.flexPayRunActionForTenant(w, r, tenantID, id, operation, fn)
}

func (h *Handler) tenantFlexPayRunAction(w http.ResponseWriter, r *http.Request, operation string, fn func(context.Context, ports.FlexPayRunActionCommand) (*domain.FlexPayRun, error)) {
	tenantID, id, ok := h.superAdminFlexPayRunRequestIDs(w, r, operation)
	if !ok {
		return
	}
	h.flexPayRunActionForTenant(w, r, tenantID, id, operation, fn)
}

func (h *Handler) flexPayRunActionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string, fn func(context.Context, ports.FlexPayRunActionCommand) (*domain.FlexPayRun, error)) {
	var cmd ports.FlexPayRunActionCommand
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&cmd)
	}
	cmd.TenantID = tenantID
	cmd.FlexPayRunID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := fn(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) exportFlexPayRunForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	data, name, err := h.svc.ExportFlexPayRunCSV(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	_, _ = w.Write(data)
}

func (h *Handler) createContractorInvoiceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.ContractorInvoiceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateContractorInvoice(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateContractorInvoiceForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string) {
	var cmd ports.ContractorInvoiceCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.ID = id
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdateContractorInvoice(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listContractorInvoicesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	flexPayRunID, ok := h.optionalUUIDQuery(w, r, "flex_pay_run_id", operation)
	if !ok {
		return
	}
	workerProfileID, ok := h.optionalUUIDQuery(w, r, "worker_profile_id", operation)
	if !ok {
		return
	}
	filter := domain.ContractorInvoiceFilter{TenantID: tenantID, FlexPayRunID: flexPayRunID, WorkerProfileID: workerProfileID, Status: optionalStringQuery(r, "status"), Search: optionalStringQuery(r, "search")}
	items, err := h.svc.ListContractorInvoices(r.Context(), filter)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list contractor invoices")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) contractorInvoiceAction(w http.ResponseWriter, r *http.Request, operation string, fn func(context.Context, ports.ContractorInvoiceActionCommand) (*domain.ContractorInvoice, error)) {
	tenantID, id, ok := h.invoiceRequestIDs(w, r, operation)
	if !ok {
		return
	}
	h.contractorInvoiceActionForTenant(w, r, tenantID, id, operation, fn)
}

func (h *Handler) tenantContractorInvoiceAction(w http.ResponseWriter, r *http.Request, operation string, fn func(context.Context, ports.ContractorInvoiceActionCommand) (*domain.ContractorInvoice, error)) {
	tenantID, id, ok := h.superAdminInvoiceRequestIDs(w, r, operation)
	if !ok {
		return
	}
	h.contractorInvoiceActionForTenant(w, r, tenantID, id, operation, fn)
}

func (h *Handler) contractorInvoiceActionForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, id uuid.UUID, operation string, fn func(context.Context, ports.ContractorInvoiceActionCommand) (*domain.ContractorInvoice, error)) {
	var cmd ports.ContractorInvoiceActionCommand
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&cmd)
	}
	cmd.TenantID = tenantID
	cmd.InvoiceID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := fn(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) createFlexPayRunItemForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, flexPayRunID uuid.UUID, operation string) {
	var cmd ports.FlexPayRunItemCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.FlexPayRunID = flexPayRunID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreateFlexPayRunItem(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) flexPayRunRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "flexPayRunID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid flexible pay run id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminFlexPayRunRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "flexPayRunID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid flexible pay run id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) invoiceRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, operation, err, "tenant context is required")
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "invoiceID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid contractor invoice id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func (h *Handler) superAdminInvoiceRequestIDs(w http.ResponseWriter, r *http.Request, operation string) (uuid.UUID, uuid.UUID, bool) {
	tenantID, ok := h.superAdminTenantID(w, r, operation)
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	id, err := uuid.Parse(chi.URLParam(r, "invoiceID"))
	if err != nil || id == uuid.Nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid contractor invoice id")
		return uuid.Nil, uuid.Nil, false
	}
	return tenantID, id, true
}

func optionalFlexDateQuery(w http.ResponseWriter, r *http.Request, key string, operation string, h *Handler) (*time.Time, bool) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, true
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, "invalid "+key)
		return nil, false
	}
	return &parsed, true
}
