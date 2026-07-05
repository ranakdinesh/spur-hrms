package services

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/pkg/permissions"
)

func (s *TenantService) ListOperationCatalog(ctx context.Context, query ports.OperationCatalogQuery) (*domain.OperationCatalog, error) {
	entries := defaultOperationCatalogEntries()
	if query.TenantID != uuid.Nil {
		active := true
		templates, err := s.ListOperationTemplates(ctx, query.TenantID, nil, nil, &active, nil, 200, 0)
		if err != nil {
			s.log.Error().Err(err).Str("operation", "list operation catalog").Str("tenant_id", query.TenantID.String()).Msg("operation template catalog merge failed")
			return nil, err
		}
		entries = append(entries, operationCatalogEntriesFromTemplates(templates)...)
	}
	if !query.IncludeAll {
		entries = filterOperationCatalogByPermissions(entries, query.Permissions)
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Category != entries[j].Category {
			return entries[i].Category < entries[j].Category
		}
		if entries[i].SortOrder != entries[j].SortOrder {
			return entries[i].SortOrder < entries[j].SortOrder
		}
		return entries[i].Label < entries[j].Label
	})
	var tenantID *uuid.UUID
	if query.TenantID != uuid.Nil {
		tenantID = &query.TenantID
	}
	return &domain.OperationCatalog{TenantID: tenantID, Entries: entries, Groups: buildOperationCatalogGroups(entries)}, nil
}

func defaultOperationCatalogEntries() []*domain.OperationCatalogEntry {
	return []*domain.OperationCatalogEntry{
		catalogEntry("employee.apply_leave", "Employee Request", "Apply Leave", "Submit a leave request through the leave workflow.", []string{permissions.LeavesApply}, "hrms", "leaves", "", domain.OperationLaunchNavigate, "leave_policy", "leave.apply", "leave", "leave_request", "Leave request", 50, "medium", "", 10, true),
		catalogEntry("employee.regularize_attendance", "Employee Request", "Regularize Attendance", "Request correction for a missed or incorrect attendance punch.", []string{permissions.AttendanceReviewRequest}, "hrms", "attendance", "", domain.OperationLaunchNavigate, "attendance_exception_policy", "attendance.regularization.request", "attendance", "regularization_request", "Attendance regularization", 45, "medium", "HR", 20, true),
		catalogEntry("employee.submit_claim", "Employee Request", "Submit Claim", "Create a reimbursement or benefit claim with evidence.", []string{permissions.BenefitClaimsCreate}, "hrms", "benefits-claims", "", domain.OperationLaunchNavigate, "benefit_claim_policy", "benefits.claim.create", "benefits", "benefit_claim", "Benefit claim", 55, "medium", "HR", 30, true),
		catalogEntry("employee.request_document", "Employee Request", "Request Document", "Ask HR for a certificate, letter, or employee document.", []string{permissions.HRCasesCreate}, "hrms", "hr-helpdesk", "", domain.OperationLaunchNavigate, "hr_case_sla", "hr_case.document_request", "hr_cases", "document_request", "Document request", 60, "low", "HR", 40, true),
		catalogEntry("employee.raise_case", "Employee Request", "Raise HR Case", "Open a helpdesk case for HR support or policy questions.", []string{permissions.HRCasesCreate}, "hrms", "hr-helpdesk", "", domain.OperationLaunchNavigate, "hr_case_sla", "hr_case.create", "hr_cases", "employee_case", "HR case", 55, "medium", "HR", 50, true),
		catalogEntry("hr.add_employee", "HR Operation", "Add Employee", "Create an employee record and start profile completion.", []string{permissions.EmployeesCreate}, "hrms", "employees", "", domain.OperationLaunchNavigate, "employee_master_control", "employees.create", "employees", "employee", "New employee", 35, "medium", "HR", 110, false),
		catalogEntry("hr.start_onboarding", "HR Operation", "Start Onboarding", "Start candidate or employee onboarding tasks.", []string{permissions.OnboardingStart}, "hrms", "candidate-onboarding", "", domain.OperationLaunchNavigate, "onboarding_workflow", "onboarding.start", "onboarding", "onboarding", "Onboarding request", 35, "medium", "HR", 120, false),
		catalogEntry("hr.create_task", "HR Operation", "Create Workflow Task", "Create a routed workflow task with comments and attachments.", []string{permissions.WorkflowTasksCreate}, "hrms", "workflow-inbox", "", domain.OperationLaunchWorkflowTask, "workflow_definition", "workflow.task.create", "workflow", "manual_task", "Workflow task", 50, "medium", "HR", 130, true),
		catalogEntry("tenant.add_department", "Tenant Operation", "Add Department", "Create a tenant department for employee classification and reporting.", []string{permissions.DepartmentsCreate}, "hrms", "departments", "", domain.OperationLaunchNavigate, "tenant_admin", "departments.create", "departments", "department", "Department setup", 60, "low", "HR", 210, false),
		catalogEntry("tenant.configure_attendance", "Tenant Operation", "Configure Attendance", "Open attendance locations, devices, and working-hour setup.", []string{permissions.AttendanceLocationsCreate, permissions.WorkingHoursCreate}, "hrms", "attendance", "", domain.OperationLaunchNavigate, "tenant_admin", "attendance.configure", "attendance", "attendance_setup", "Attendance setup", 45, "medium", "HR", 220, false),
		catalogEntry("tenant.storage_provider", "Tenant Operation", "Configure Storage", "Set tenant document storage for MinIO, S3, or a tenant-owned provider.", []string{permissions.IntegrationsStorageManage}, "hrms", "storage-providers", "", domain.OperationLaunchNavigate, "tenant_admin", "storage.configure", "storage", "storage_provider", "Storage provider setup", 35, "high", "HR", 230, false),
		catalogEntry("payroll.create_pay_group", "Payroll Operation", "Create Pay Group", "Group employees for phased payroll processing.", []string{permissions.PayGroupsCreate}, "hrms", "flexible-payroll", "", domain.OperationLaunchNavigate, "payroll_admin", "pay_groups.create", "payroll", "pay_group", "Pay group setup", 40, "medium", "Payroll", 310, false),
		catalogEntry("payroll.upload_template", "Payroll Operation", "Upload Payroll Data", "Import payroll inputs, adjustments, or attendance-linked payroll data.", []string{permissions.PayrollImportsCreate}, "hrms", "payroll-operations", "", domain.OperationLaunchNavigate, "payroll_lock_policy", "payroll.import", "payroll", "payroll_import", "Payroll import", 30, "high", "Payroll", 320, false),
		catalogEntry("document.generate_letter", "Document Workflow", "Generate HR Letter", "Generate appointment, experience, relieving, or other HR letters.", []string{permissions.EmployeeLettersCreate}, "hrms", "employee-letters", "", domain.OperationLaunchNavigate, "document_approval", "employee_letters.create", "documents", "employee_letter", "Employee letter", 45, "medium", "HR", 410, false),
		catalogEntry("document.send_agreement", "Document Workflow", "Send Agreement", "Create and send a worker agreement or contract for signature.", []string{permissions.AgreementsCreate}, "hrms", "agreements", "", domain.OperationLaunchNavigate, "document_signing", "agreements.create", "documents", "agreement", "Agreement request", 40, "medium", "HR", 420, false),
		catalogEntry("asset.request_access", "Asset/Access Request", "Request Asset or Access", "Request equipment, software, facility, or system access.", []string{permissions.AccessTasksManage}, "hrms", "asset-access", "", domain.OperationLaunchNavigate, "asset_access_policy", "asset_access.task.create", "asset_access", "access_task", "Asset/access request", 45, "medium", "IT", 510, true),
		catalogEntry("case.raise_grievance", "Case/Grievance", "Raise Grievance", "Open a restricted case for grievance, disciplinary, or employee-relations review.", []string{permissions.ERCasesCreate}, "hrms", "employee-relations", "", domain.OperationLaunchNavigate, "er_confidential_workflow", "employee_relations.case.create", "employee_relations", "er_case", "Employee relations case", 20, "high", "HR", 610, false),
		catalogEntry("config.change_request", "Configuration Change", "Configuration Change Request", "Create a controlled task for policy, setup, or module configuration changes.", []string{permissions.WorkflowTasksCreate}, "hrms", "workflow-inbox", "configuration_change", domain.OperationLaunchWorkflowTask, "admin_approval", "workflow.config_change", "configuration", "change_request", "Configuration change request", 25, "high", "HR", 710, false),
		catalogEntry("superadmin.tenant_operation", "Tenant Operation", "Tenant Operation", "Start tenant, module, domain, storage, or admin-governance operation.", []string{permissions.TenantOperationsManage, "superadmin.tenants"}, "platform", "tenant-operations", "tenant_operation", domain.OperationLaunchNavigate, "super_admin_approval", "tenant.operation.request", "tenant", "tenant_operation", "Tenant operation", 20, "high", "Super Admin", 810, false),
	}
}

func catalogEntry(key, category, label, description string, required []string, targetModule, targetSection, workflowTemplateKey, launchMode, approvalPolicy, sourceCommand, sourceModule, sourceType, defaultTitle string, priority int32, severity, assigneeRole string, sortOrder int32, mobile bool) *domain.OperationCatalogEntry {
	return &domain.OperationCatalogEntry{
		Key:                    key,
		Category:               category,
		Label:                  label,
		Description:            description,
		RequiredPermissions:    required,
		TargetModule:           targetModule,
		TargetSection:          targetSection,
		WorkflowTemplateKey:    workflowTemplateKey,
		LaunchMode:             launchMode,
		RequiredApprovalPolicy: approvalPolicy,
		SourceServiceCommand:   sourceCommand,
		SourceModule:           sourceModule,
		SourceType:             sourceType,
		DefaultTitle:           defaultTitle,
		DefaultPriority:        priority,
		DefaultSeverity:        severity,
		AssigneeRole:           assigneeRole,
		LaunchSchema:           json.RawMessage(`{}`),
		Metadata:               json.RawMessage(`{}`),
		SortOrder:              sortOrder,
		MobileEnabled:          mobile,
	}
}

func operationCatalogEntriesFromTemplates(templates []*domain.OperationTemplate) []*domain.OperationCatalogEntry {
	entries := make([]*domain.OperationCatalogEntry, 0, len(templates))
	for _, template := range templates {
		if template == nil || !template.IsActive {
			continue
		}
		metadata := map[string]any{}
		_ = json.Unmarshal(template.Metadata, &metadata)
		required := stringSliceFromAny(metadata["required_permissions"])
		targetSection := stringFromAny(metadata["target_section"])
		if targetSection == "" {
			targetSection = "workflow-inbox"
		}
		launchMode := stringFromAny(metadata["launch_mode"])
		if launchMode == "" {
			launchMode = domain.OperationLaunchWorkflowTask
		}
		approvalPolicy := stringFromAny(metadata["approval_policy"])
		if approvalPolicy == "" {
			approvalPolicy = "workflow_definition"
		}
		sourceCommand := stringFromAny(metadata["source_service_command"])
		if sourceCommand == "" {
			sourceCommand = template.SourceModule + "." + template.SourceType
		}
		entries = append(entries, &domain.OperationCatalogEntry{
			Key:                    "template." + template.TemplateKey,
			Category:               titleFromSnake(template.Category),
			Label:                  template.Name,
			Description:            stringFromAny(metadata["description"]),
			RequiredPermissions:    required,
			TargetModule:           template.SourceModule,
			TargetSection:          targetSection,
			WorkflowTemplateKey:    template.TemplateKey,
			WorkflowTemplateID:     &template.ID,
			LaunchMode:             launchMode,
			RequiredApprovalPolicy: approvalPolicy,
			SourceServiceCommand:   sourceCommand,
			SourceModule:           template.SourceModule,
			SourceType:             template.SourceType,
			DefaultTitle:           template.Name,
			DefaultPriority:        template.DefaultPriority,
			DefaultSeverity:        template.DefaultSeverity,
			LaunchSchema:           template.LaunchSchema,
			Metadata:               template.Metadata,
			SortOrder:              900,
			MobileEnabled:          boolFromAny(metadata["mobile_enabled"]),
		})
	}
	return entries
}

func filterOperationCatalogByPermissions(entries []*domain.OperationCatalogEntry, rawPermissions []string) []*domain.OperationCatalogEntry {
	if len(rawPermissions) == 0 {
		return entries
	}
	allowed := normalizedPermissionSet(rawPermissions)
	filtered := make([]*domain.OperationCatalogEntry, 0, len(entries))
	for _, entry := range entries {
		if entry == nil || hasAnyOperationCatalogPermission(allowed, entry.RequiredPermissions) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func normalizedPermissionSet(raw []string) map[string]struct{} {
	set := make(map[string]struct{}, len(raw)*2)
	for _, item := range raw {
		clean := strings.TrimSpace(item)
		if clean == "" {
			continue
		}
		set[clean] = struct{}{}
		if strings.HasPrefix(clean, "hrms.") {
			set[strings.TrimPrefix(clean, "hrms.")] = struct{}{}
		} else {
			set["hrms."+clean] = struct{}{}
		}
	}
	return set
}

func hasAnyOperationCatalogPermission(set map[string]struct{}, required []string) bool {
	if len(required) == 0 {
		return true
	}
	for _, permission := range required {
		if _, ok := set[permission]; ok {
			return true
		}
		if _, ok := set["hrms."+permission]; ok {
			return true
		}
	}
	return false
}

func buildOperationCatalogGroups(entries []*domain.OperationCatalogEntry) []*domain.OperationCatalogGroup {
	counts := map[string]int32{}
	for _, entry := range entries {
		if entry != nil {
			counts[entry.Category]++
		}
	}
	groups := make([]*domain.OperationCatalogGroup, 0, len(counts))
	for category, count := range counts {
		groups = append(groups, &domain.OperationCatalogGroup{Category: category, Count: count})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Category < groups[j].Category })
	return groups
}

func stringSliceFromAny(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if text, ok := item.(string); ok && strings.TrimSpace(text) != "" {
			out = append(out, strings.TrimSpace(text))
		}
	}
	return out
}

func stringFromAny(value any) string {
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text)
	}
	return ""
}

func boolFromAny(value any) bool {
	flag, _ := value.(bool)
	return flag
}

func titleFromSnake(value string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(value, "_", " "))
	if clean == "" {
		return "Workflow Template"
	}
	parts := strings.Fields(clean)
	for index, part := range parts {
		parts[index] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, " ")
}
