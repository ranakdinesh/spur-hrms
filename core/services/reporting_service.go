package services

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) ListReportCatalog(ctx context.Context, query ports.ReportCatalogQuery) ([]*domain.ReportCatalogItem, error) {
	if query.TenantID == uuid.Nil {
		return nil, domain.ErrReportCatalogInvalid
	}
	if err := s.ensureDefaultReports(ctx, query.TenantID, query.ActorID); err != nil {
		s.logError("ensure default report catalog", err, serviceTenantIDField(query.TenantID))
		return nil, err
	}
	items, err := s.reporting.ListReportCatalog(ctx, query.TenantID, cleanCommandString(query.Module), cleanCommandString(query.Scope))
	if err != nil {
		s.logError("list report catalog", err, serviceTenantIDField(query.TenantID))
	}
	return items, err
}

func (s *TenantService) UpsertReportSavedView(ctx context.Context, cmd ports.ReportSavedViewCommand) (*domain.ReportSavedView, error) {
	ownerID := cmd.OwnerUserID
	if ownerID == nil {
		ownerID = cmd.ActorID
	}
	item, err := domain.NewReportSavedView(domain.ReportSavedView{ID: cmd.ID, TenantID: cmd.TenantID, ReportID: cmd.ReportID, Name: cmd.Name, Description: cleanCommandString(cmd.Description), Visibility: cmd.Visibility, Filters: cmd.Filters, Columns: cmd.Columns, IsFavorite: cmd.IsFavorite, OwnerUserID: ownerID})
	if err != nil {
		s.logError("validate report saved view", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.reporting.UpsertReportSavedView(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListReportSavedViews(ctx context.Context, query ports.ReportListQuery) ([]*domain.ReportSavedView, error) {
	return s.reporting.ListReportSavedViews(ctx, query.TenantID, query.ReportID)
}

func (s *TenantService) DeleteReportSavedView(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	return s.reporting.DeleteReportSavedView(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateReportExportJob(ctx context.Context, cmd ports.ReportExportJobCommand) (*domain.ReportExportJob, error) {
	item, err := domain.NewReportExportJob(domain.ReportExportJob{TenantID: cmd.TenantID, ReportID: cmd.ReportID, SavedViewID: cmd.SavedViewID, ExportFormat: cmd.ExportFormat, Status: domain.ReportExportQueued, Filters: cmd.Filters, RequestedBy: cmd.ActorID})
	if err != nil {
		s.logError("validate report export job", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.reporting.CreateReportExportJob(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListReportExportJobs(ctx context.Context, query ports.ReportListQuery) ([]*domain.ReportExportJob, error) {
	limit, offset := normalizeReportPaging(query.Limit, query.Offset)
	return s.reporting.ListReportExportJobs(ctx, query.TenantID, query.ReportID, cleanCommandString(query.Status), limit, offset)
}

func (s *TenantService) UpdateReportExportJobStatus(ctx context.Context, cmd ports.ReportExportJobStatusCommand) (*domain.ReportExportJob, error) {
	item, err := domain.NewReportExportJob(domain.ReportExportJob{ID: cmd.ID, TenantID: cmd.TenantID, ReportID: uuid.New(), Status: cmd.Status, ExportFormat: domain.ReportExportCSV, FileObjectKey: cleanCommandString(cmd.FileObjectKey), ErrorMessage: cleanCommandString(cmd.ErrorMessage)})
	if err != nil {
		s.logError("validate report export status", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item.ReportID = uuid.Nil
	return s.reporting.UpdateReportExportJobStatus(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpsertReportSchedule(ctx context.Context, cmd ports.ReportScheduleCommand) (*domain.ReportSchedule, error) {
	nextRunAt, err := parseOptionalRFC3339(cmd.NextRunAt)
	if err != nil {
		s.logError("parse report schedule next run", err, serviceTenantIDField(cmd.TenantID))
		return nil, domain.ErrReportScheduleInvalid
	}
	item, err := domain.NewReportSchedule(domain.ReportSchedule{ID: cmd.ID, TenantID: cmd.TenantID, ReportID: cmd.ReportID, SavedViewID: cmd.SavedViewID, Name: cmd.Name, Frequency: cmd.Frequency, Timezone: cmd.Timezone, DeliveryChannels: cmd.DeliveryChannels, RecipientUserIDs: cmd.RecipientUserIDs, RecipientEmails: cmd.RecipientEmails, NextRunAt: nextRunAt, IsActive: cmd.IsActive})
	if err != nil {
		s.logError("validate report schedule", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.reporting.UpsertReportSchedule(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListReportSchedules(ctx context.Context, query ports.ReportListQuery) ([]*domain.ReportSchedule, error) {
	return s.reporting.ListReportSchedules(ctx, query.TenantID, query.ReportID)
}

func (s *TenantService) DeleteReportSchedule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	return s.reporting.DeleteReportSchedule(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateReportSnapshot(ctx context.Context, cmd ports.ReportSnapshotCommand) (*domain.ReportSnapshot, error) {
	periodStart, err := parseOptionalDate(cmd.PeriodStart)
	if err != nil || periodStart == nil {
		return nil, domain.ErrReportSnapshotInvalid
	}
	periodEnd, err := parseOptionalDate(cmd.PeriodEnd)
	if err != nil || periodEnd == nil {
		return nil, domain.ErrReportSnapshotInvalid
	}
	item, err := domain.NewReportSnapshot(domain.ReportSnapshot{TenantID: cmd.TenantID, ReportID: cmd.ReportID, SavedViewID: cmd.SavedViewID, SnapshotKey: cmd.SnapshotKey, PeriodStart: *periodStart, PeriodEnd: *periodEnd, Filters: cmd.Filters, Summary: cmd.Summary, RowCount: cmd.RowCount, GeneratedBy: cmd.ActorID})
	if err != nil {
		s.logError("validate report snapshot", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.reporting.CreateReportSnapshot(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListReportSnapshots(ctx context.Context, query ports.ReportListQuery) ([]*domain.ReportSnapshot, error) {
	limit, offset := normalizeReportPaging(query.Limit, query.Offset)
	return s.reporting.ListReportSnapshots(ctx, query.TenantID, query.ReportID, limit, offset)
}

func (s *TenantService) ensureDefaultReports(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) error {
	for _, seed := range defaultReportCatalog(tenantID) {
		item, err := domain.NewReportCatalogItem(seed)
		if err != nil {
			return err
		}
		if _, err := s.reporting.UpsertReportCatalogItem(ctx, item, actorID); err != nil {
			return err
		}
	}
	return nil
}

func defaultReportCatalog(tenantID uuid.UUID) []domain.ReportCatalogItem {
	type reportSeed struct {
		Code        string
		Module      string
		Name        string
		Description string
		Category    string
		Filters     []string
		Columns     []string
		SortOrder   int32
	}
	seeds := []reportSeed{
		{Code: "workforce.composition", Module: "employees", Name: "Workforce Composition", Description: "Headcount, movement, location, department, designation, and employment-type mix.", Category: "Workforce", Filters: []string{"date", "financial_year", "branch", "department", "designation", "employment_type"}, Columns: []string{"employee_code", "employee_name", "branch", "department", "designation", "employment_type", "joining_date", "status"}, SortOrder: 10},
		{Code: "workforce.movement", Module: "employees", Name: "Workforce Movement", Description: "Joiners, exits, probation, confirmation, branch and department movement indicators.", Category: "Workforce", Filters: []string{"date_range", "branch", "department", "employment_type"}, Columns: []string{"employee_code", "employee_name", "joining_date", "resignation_date", "probation_status", "branch", "department"}, SortOrder: 11},
		{Code: "workforce.data_quality", Module: "employees", Name: "Employee Data Quality", Description: "Missing employee master data that affects payroll, reporting, compliance, and communication.", Category: "Workforce", Filters: []string{"branch", "department", "employment_type"}, Columns: []string{"employee_code", "employee_name", "email", "mobile", "branch", "department", "status"}, SortOrder: 12},
		{Code: "attendance.health", Module: "attendance", Name: "Attendance Health", Description: "Attendance rate, absenteeism, late trends, work-mode split, and incomplete day analysis.", Category: "Time & Leave", Filters: []string{"date_range", "branch", "department", "work_mode"}, Columns: []string{"employee_code", "employee_name", "date", "status", "rule_outcome", "worked_minutes", "late_minutes"}, SortOrder: 19},
		{Code: "attendance.exceptions", Module: "attendance", Name: "Attendance Exceptions", Description: "Late, missed punch, location mismatch, device punch, regularisation, and overtime exceptions.", Category: "Time & Leave", Filters: []string{"date_range", "pay_cycle", "branch", "department", "attendance_location", "device"}, Columns: []string{"employee_code", "employee_name", "date", "exception_type", "source", "location", "status"}, SortOrder: 20},
		{Code: "leave.liability", Module: "leave", Name: "Leave Liability", Description: "Leave balances, utilization, pending approvals, carry-forward, encashment, and policy exposure.", Category: "Time & Leave", Filters: []string{"financial_year", "leave_type", "branch", "department", "employment_type"}, Columns: []string{"employee_code", "employee_name", "leave_type", "opening", "earned", "used", "pending", "balance"}, SortOrder: 30},
		{Code: "leave.utilization", Module: "leave", Name: "Leave Utilization", Description: "Applied, approved, pending, rejected, canceled and sandwich leave usage by employee and department.", Category: "Time & Leave", Filters: []string{"financial_year", "date_range", "leave_type", "branch", "department", "status"}, Columns: []string{"employee_code", "employee_name", "leave_type", "start_date", "end_date", "days", "status"}, SortOrder: 31},
		{Code: "payroll.readiness", Module: "payroll", Name: "Payroll Readiness", Description: "Payroll blockers, attendance LOP status, missing salary structures, reimbursement adjustments, and statutory readiness.", Category: "Payroll", Filters: []string{"pay_cycle", "month", "year", "branch", "department", "pay_group"}, Columns: []string{"employee_code", "employee_name", "readiness_status", "blocker", "last_updated"}, SortOrder: 40},
		{Code: "payroll.cost", Module: "payroll", Name: "Payroll Cost", Description: "Gross, earnings, deductions, LOP deductions, net payroll and period cost movement.", Category: "Payroll", Filters: []string{"month", "year", "branch", "department", "pay_group"}, Columns: []string{"employee_code", "employee_name", "gross_salary", "total_earnings", "total_deductions", "net_salary"}, SortOrder: 41},
		{Code: "payroll.compliance", Module: "payroll", Name: "Payroll Compliance", Description: "Statutory readiness, employee master gaps, probation status, PT/LWF readiness, and payroll audit exposure.", Category: "Payroll", Filters: []string{"month", "year", "branch", "department", "state"}, Columns: []string{"employee_code", "employee_name", "risk_area", "status"}, SortOrder: 42},
		{Code: "payroll.reconciliation", Module: "payroll", Name: "Payroll Reconciliation", Description: "Payroll attendance and LOP reconciliation rows for salary processing.", Category: "Payroll", Filters: []string{"month", "year", "branch", "department"}, Columns: []string{"employee_code", "employee_name", "present_days", "absent_days", "lwp_days", "net_salary", "status"}, SortOrder: 43},
		{Code: "payroll.consolidated_salary_sheet", Module: "payroll", Name: "Consolidated Salary Sheet", Description: "Branded monthly salary sheet with gross, earnings, deductions, LOP and net payroll values.", Category: "Payroll", Filters: []string{"month", "year", "branch", "department"}, Columns: []string{"employee_code", "employee_name", "gross_salary", "total_earnings", "total_deductions", "net_salary"}, SortOrder: 44},
		{Code: "recruitment.funnel", Module: "recruitment", Name: "Recruitment Funnel", Description: "Open requisitions, source conversion, candidate-stage aging, interviews, offers, time-to-fill, and onboarding handoff.", Category: "Talent Lifecycle", Filters: []string{"date_range", "job_position", "department", "stage", "source"}, Columns: []string{"job_position", "candidate", "stage", "source", "age_days", "status"}, SortOrder: 50},
		{Code: "recruitment.source_effectiveness", Module: "recruitment", Name: "Source Effectiveness", Description: "Source/channel applications, interview throughput, offer acceptance proxy, joining conversion, rejection, withdrawal, and aging metrics.", Category: "Talent Lifecycle", Filters: []string{"date_range", "source", "job_position", "department"}, Columns: []string{"source", "applications", "interviews", "offered", "hired", "rejected", "joining_conversion"}, SortOrder: 51},
		{Code: "onboarding.completion", Module: "recruitment", Name: "Onboarding Completion", Description: "Candidate onboarding task progress, required task completion, missing onboarding documents/tasks, overdue work, and joining handoff status.", Category: "Talent Lifecycle", Filters: []string{"date_range", "status", "workflow"}, Columns: []string{"candidate", "workflow", "status", "progress", "completed_tasks", "overdue_tasks"}, SortOrder: 52},
		{Code: "probation.due", Module: "employees", Name: "Probation Confirmations Due", Description: "Employees whose probation confirmation, extension, or payroll-status review is overdue or due soon.", Category: "Talent Lifecycle", Filters: []string{"date", "branch", "department", "employment_type"}, Columns: []string{"employee_code", "employee_name", "probation_end", "status", "department", "branch"}, SortOrder: 53},
		{Code: "exit.pipeline", Module: "employees", Name: "Exit Pipeline", Description: "Resignation and exit pipeline, notice-period exposure, exit task progress, and reason analytics.", Category: "Talent Lifecycle", Filters: []string{"date_range", "status", "branch", "department", "exit_type"}, Columns: []string{"employee_code", "employee_name", "status", "exit_type", "last_working_date", "reason"}, SortOrder: 54},
		{Code: "exit.readiness", Module: "employees", Name: "Exit and F&F Readiness", Description: "Asset clearance, handover, access revocation, exit interview, final settlement readiness, blocked tasks, and completion exposure.", Category: "Talent Lifecycle", Filters: []string{"date_range", "status", "branch", "department", "readiness"}, Columns: []string{"employee_code", "employee_name", "fnf_status", "asset_status", "handover_status", "blocked_tasks", "readiness"}, SortOrder: 55},
		{Code: "compliance.readiness", Module: "compliance", Name: "Compliance Readiness", Description: "Missing documents, statutory details, policy acknowledgements, probation, exit, and audit follow-up exposure.", Category: "Compliance", Filters: []string{"date", "branch", "department", "document_type", "policy_type"}, Columns: []string{"employee_code", "employee_name", "risk_area", "required_item", "status", "due_date"}, SortOrder: 60},
	}
	items := make([]domain.ReportCatalogItem, 0, len(seeds))
	for _, seed := range seeds {
		description := seed.Description
		items = append(items, domain.ReportCatalogItem{
			TenantID:          tenantID,
			ReportCode:        seed.Code,
			Module:            seed.Module,
			Name:              seed.Name,
			Description:       &description,
			Category:          seed.Category,
			Scope:             domain.ReportScopeTenant,
			PermissionKey:     "reports.view",
			DefaultFilters:    rawJSON(map[string]any{}),
			SupportedFilters:  rawJSON(seed.Filters),
			OutputColumns:     rawJSON(seed.Columns),
			DrilldownContract: rawJSON(map[string]any{"route": "/reports/drilldown", "keys": []string{"tenant_id", "report_code", "filters"}}),
			IsSystem:          true,
			IsActive:          true,
			SortOrder:         seed.SortOrder,
		})
	}
	return items
}

func rawJSON(value any) json.RawMessage {
	data, _ := json.Marshal(value)
	return data
}

func normalizeReportPaging(limit int32, offset int32) (int32, int32) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
