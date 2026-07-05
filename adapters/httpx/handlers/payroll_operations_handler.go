package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (h *Handler) ListPayrollPeriodLocks(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list payroll locks", err, "tenant context is required")
		return
	}
	h.listPayrollPeriodLocksForTenant(w, r, tenantID, "list payroll locks")
}

func (h *Handler) UpsertPayrollPeriodLock(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "set payroll lock", err, "tenant context is required")
		return
	}
	h.upsertPayrollPeriodLockForTenant(w, r, tenantID, "set payroll lock")
}

func (h *Handler) ListTenantPayrollPeriodLocks(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant payroll locks")
	if !ok {
		return
	}
	h.listPayrollPeriodLocksForTenant(w, r, tenantID, "list tenant payroll locks")
}

func (h *Handler) UpsertTenantPayrollPeriodLock(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "set tenant payroll lock")
	if !ok {
		return
	}
	h.upsertPayrollPeriodLockForTenant(w, r, tenantID, "set tenant payroll lock")
}

func (h *Handler) ListPayrollStatutoryRules(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list payroll statutory rules", err, "tenant context is required")
		return
	}
	h.listPayrollStatutoryRulesForTenant(w, r, tenantID, "list payroll statutory rules")
}

func (h *Handler) CreatePayrollStatutoryRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "create payroll statutory rule", err, "tenant context is required")
		return
	}
	h.createPayrollStatutoryRuleForTenant(w, r, tenantID, "create payroll statutory rule")
}

func (h *Handler) UpdatePayrollStatutoryRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "update payroll statutory rule", err, "tenant context is required")
		return
	}
	h.updatePayrollStatutoryRuleForTenant(w, r, tenantID, "update payroll statutory rule")
}

func (h *Handler) DeletePayrollStatutoryRule(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "delete payroll statutory rule", err, "tenant context is required")
		return
	}
	h.deletePayrollStatutoryRuleForTenant(w, r, tenantID, "delete payroll statutory rule")
}

func (h *Handler) ListTenantPayrollStatutoryRules(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant payroll statutory rules")
	if !ok {
		return
	}
	h.listPayrollStatutoryRulesForTenant(w, r, tenantID, "list tenant payroll statutory rules")
}

func (h *Handler) CreateTenantPayrollStatutoryRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "create tenant payroll statutory rule")
	if !ok {
		return
	}
	h.createPayrollStatutoryRuleForTenant(w, r, tenantID, "create tenant payroll statutory rule")
}

func (h *Handler) UpdateTenantPayrollStatutoryRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "update tenant payroll statutory rule")
	if !ok {
		return
	}
	h.updatePayrollStatutoryRuleForTenant(w, r, tenantID, "update tenant payroll statutory rule")
}

func (h *Handler) DeleteTenantPayrollStatutoryRule(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "delete tenant payroll statutory rule")
	if !ok {
		return
	}
	h.deletePayrollStatutoryRuleForTenant(w, r, tenantID, "delete tenant payroll statutory rule")
}

func (h *Handler) ImportPayrollData(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "import payroll data", err, "tenant context is required")
		return
	}
	h.importPayrollDataForTenant(w, r, tenantID, "import payroll data")
}

func (h *Handler) ListPayrollImportBatches(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list payroll imports", err, "tenant context is required")
		return
	}
	h.listPayrollImportBatchesForTenant(w, r, tenantID, "list payroll imports")
}

func (h *Handler) GetPayrollImportBatch(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "get payroll import", err, "tenant context is required")
		return
	}
	h.getPayrollImportBatchForTenant(w, r, tenantID, "get payroll import")
}

func (h *Handler) ImportTenantPayrollData(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "import tenant payroll data")
	if ok {
		h.importPayrollDataForTenant(w, r, tenantID, "import tenant payroll data")
	}
}
func (h *Handler) ListTenantPayrollImportBatches(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant payroll imports")
	if ok {
		h.listPayrollImportBatchesForTenant(w, r, tenantID, "list tenant payroll imports")
	}
}
func (h *Handler) GetTenantPayrollImportBatch(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "get tenant payroll import")
	if ok {
		h.getPayrollImportBatchForTenant(w, r, tenantID, "get tenant payroll import")
	}
}

func (h *Handler) ListConsolidatedSalarySheet(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list salary sheet", err, "tenant context is required")
		return
	}
	h.listConsolidatedSalarySheetForTenant(w, r, tenantID, "list salary sheet")
}
func (h *Handler) ExportConsolidatedSalarySheet(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "export salary sheet", err, "tenant context is required")
		return
	}
	h.exportConsolidatedSalarySheetForTenant(w, r, tenantID, "export salary sheet")
}
func (h *Handler) ListPayrollReconciliation(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.tenantIDFromRequest(r)
	if err != nil {
		h.respondError(w, r, http.StatusUnauthorized, "list payroll reconciliation", err, "tenant context is required")
		return
	}
	h.listPayrollReconciliationForTenant(w, r, tenantID, "list payroll reconciliation")
}
func (h *Handler) ListTenantConsolidatedSalarySheet(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant salary sheet")
	if ok {
		h.listConsolidatedSalarySheetForTenant(w, r, tenantID, "list tenant salary sheet")
	}
}
func (h *Handler) ExportTenantConsolidatedSalarySheet(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "export tenant salary sheet")
	if ok {
		h.exportConsolidatedSalarySheetForTenant(w, r, tenantID, "export tenant salary sheet")
	}
}
func (h *Handler) ListTenantPayrollReconciliation(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := h.superAdminTenantID(w, r, "list tenant payroll reconciliation")
	if ok {
		h.listPayrollReconciliationForTenant(w, r, tenantID, "list tenant payroll reconciliation")
	}
}

func (h *Handler) listPayrollPeriodLocksForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	items, err := h.svc.ListPayrollPeriodLocks(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list payroll locks")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) upsertPayrollPeriodLockForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PayrollPeriodLockCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpsertPayrollPeriodLock(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) listPayrollStatutoryRulesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var ruleType *string
	if value := r.URL.Query().Get("rule_type"); value != "" {
		ruleType = &value
	}
	items, err := h.svc.ListPayrollStatutoryRules(r.Context(), tenantID, ruleType)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list statutory rules")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) createPayrollStatutoryRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PayrollStatutoryRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.CreatePayrollStatutoryRule(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) updatePayrollStatutoryRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "ruleID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse rule id", err, "invalid rule id")
		return
	}
	var cmd ports.PayrollStatutoryRuleCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ID = id
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.UpdatePayrollStatutoryRule(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *Handler) deletePayrollStatutoryRuleForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "ruleID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse rule id", err, "invalid rule id")
		return
	}
	if err := h.svc.DeletePayrollStatutoryRule(r.Context(), tenantID, id, h.actorIDFromRequest(r)); err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) importPayrollDataForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	var cmd ports.PayrollImportCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "decode "+operation, err, "invalid request body")
		return
	}
	cmd.TenantID = tenantID
	cmd.ActorID = h.actorIDFromRequest(r)
	item, err := h.svc.ImportPayrollData(r.Context(), cmd)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, item)
}

func (h *Handler) listPayrollImportBatchesForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	items, err := h.svc.ListPayrollImportBatches(r.Context(), tenantID, int32(limit), int32(offset))
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, operation, err, "failed to list payroll imports")
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getPayrollImportBatchForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	id, err := uuid.Parse(chi.URLParam(r, "batchID"))
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse batch id", err, "invalid batch id")
		return
	}
	item, err := h.svc.GetPayrollImportBatch(r.Context(), tenantID, id)
	if err != nil {
		h.respondError(w, r, http.StatusNotFound, operation, err, "payroll import not found")
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func payrollPeriodParams(r *http.Request) (int32, int32, error) {
	month, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil {
		return 0, 0, err
	}
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		return 0, 0, err
	}
	return int32(month), int32(year), nil
}

func (h *Handler) listConsolidatedSalarySheetForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	month, year, err := payrollPeriodParams(r)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse payroll period", err, "month and year are required")
		return
	}
	items, err := h.svc.ListConsolidatedSalarySheet(r.Context(), tenantID, month, year)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) exportConsolidatedSalarySheetForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	month, year, err := payrollPeriodParams(r)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse payroll period", err, "month and year are required")
		return
	}
	format := r.URL.Query().Get("format")
	if format == "pdf" || format == "xlsx" || format == "excel" {
		h.downloadReportForTenant(w, r, ports.ReportDatasetQuery{TenantID: tenantID, ReportCode: "payroll.consolidated_salary_sheet", Month: month, Year: year}, operation)
		return
	}
	data, name, err := h.svc.ExportConsolidatedSalarySheetCSV(r.Context(), tenantID, month, year)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	_, _ = w.Write(data)
}

func (h *Handler) listPayrollReconciliationForTenant(w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, operation string) {
	month, year, err := payrollPeriodParams(r)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, "parse payroll period", err, "month and year are required")
		return
	}
	items, err := h.svc.ListPayrollReconciliationRows(r.Context(), tenantID, month, year)
	if err != nil {
		h.respondError(w, r, http.StatusBadRequest, operation, err, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, items)
}
