package permissions

import (
	"strings"

	identity "github.com/ranakdinesh/spur-identity"
)

const (
	ModuleCode        = "hrms"
	ModuleName        = "Spur HRMS"
	ModuleDescription = "Human resource management: tenant setup, employees, leave, attendance, payroll, notifications, celebrations, and onboarding."
)

const (
	TenantProfilesView   = "tenant_profiles.view"
	TenantProfilesUpdate = "tenant_profiles.update"
	TenantProvision      = "tenants.provision"
	TenantSettingsList   = "tenant_settings.list"
	TenantSettingsUpdate = "tenant_settings.update"

	BrandingView   = "branding.view"
	BrandingUpdate = "branding.update"

	BranchesList   = "branches.list"
	BranchesCreate = "branches.create"
	BranchesView   = "branches.view"
	BranchesUpdate = "branches.update"
	BranchesDelete = "branches.delete"

	DepartmentsList   = "departments.list"
	DepartmentsCreate = "departments.create"
	DepartmentsView   = "departments.view"
	DepartmentsUpdate = "departments.update"
	DepartmentsDelete = "departments.delete"

	DesignationsList                        = "designations.list"
	DesignationsCreate                      = "designations.create"
	DesignationsView                        = "designations.view"
	DesignationsUpdate                      = "designations.update"
	DesignationsDelete                      = "designations.delete"
	DesignationsMastersManage               = "designations.masters.manage"
	DesignationsAttendanceRequirementUpdate = "designations.attendance_requirement.update"

	WorkerTypesList                      = "worker_types.list"
	WorkerTypesCreate                    = "worker_types.create"
	WorkerTypesView                      = "worker_types.view"
	WorkerTypesUpdate                    = "worker_types.update"
	WorkerTypesDelete                    = "worker_types.delete"
	WorkerClassificationRulesList        = "worker_classification_rules.list"
	WorkerClassificationRulesManage      = "worker_classification_rules.manage"
	WorkersList                          = "workers.list"
	WorkersCreate                        = "workers.create"
	WorkersView                          = "workers.view"
	WorkersUpdate                        = "workers.update"
	WorkersDelete                        = "workers.delete"
	EngagementsList                      = "engagements.list"
	EngagementsCreate                    = "engagements.create"
	EngagementsView                      = "engagements.view"
	EngagementsUpdate                    = "engagements.update"
	EngagementsStatus                    = "engagements.status"
	EngagementsDelete                    = "engagements.delete"
	WorkLogsList                         = "work_logs.list"
	WorkLogsCreate                       = "work_logs.create"
	WorkLogsView                         = "work_logs.view"
	WorkLogsUpdate                       = "work_logs.update"
	WorkLogsSubmit                       = "work_logs.submit"
	WorkLogsApprove                      = "work_logs.approve"
	WorkLogsReject                       = "work_logs.reject"
	WorkLogsDelete                       = "work_logs.delete"
	WorkLogsReport                       = "work_logs.report"
	ProjectsList                         = "projects.list"
	ProjectsCreate                       = "projects.create"
	ProjectsView                         = "projects.view"
	ProjectsUpdate                       = "projects.update"
	ProjectsStatus                       = "projects.status"
	ProjectsDelete                       = "projects.delete"
	MilestonesList                       = "milestones.list"
	MilestonesCreate                     = "milestones.create"
	MilestonesView                       = "milestones.view"
	MilestonesUpdate                     = "milestones.update"
	MilestonesSubmit                     = "milestones.submit"
	MilestonesApprove                    = "milestones.approve"
	MilestonesReject                     = "milestones.reject"
	MilestonesDelete                     = "milestones.delete"
	MilestonesEventsView                 = "milestones.events.view"
	ComplianceRulesList                  = "compliance_rules.list"
	ComplianceRulesCreate                = "compliance_rules.create"
	ComplianceRulesView                  = "compliance_rules.view"
	ComplianceRulesUpdate                = "compliance_rules.update"
	ComplianceRulesDelete                = "compliance_rules.delete"
	ComplianceRulesSeed                  = "compliance_rules.seed"
	ComplianceChecklistList              = "compliance_checklist.list"
	ComplianceChecklistGenerate          = "compliance_checklist.generate"
	ComplianceChecklistReview            = "compliance_checklist.review"
	ComplianceChecklistEvidence          = "compliance_checklist.evidence"
	ComplianceChecklistWaive             = "compliance_checklist.waive"
	ComplianceChecklistDelete            = "compliance_checklist.delete"
	ComplianceEventsView                 = "compliance_events.view"
	SkillCategoriesList                  = "skill_categories.list"
	SkillCategoriesCreate                = "skill_categories.create"
	SkillCategoriesUpdate                = "skill_categories.update"
	SkillCategoriesDelete                = "skill_categories.delete"
	SkillsList                           = "skills.list"
	SkillsCreate                         = "skills.create"
	SkillsView                           = "skills.view"
	SkillsUpdate                         = "skills.update"
	SkillsDelete                         = "skills.delete"
	WorkerSkillsList                     = "worker_skills.list"
	WorkerSkillsCreate                   = "worker_skills.create"
	WorkerSkillsUpdate                   = "worker_skills.update"
	WorkerSkillsVerify                   = "worker_skills.verify"
	WorkerSkillsDelete                   = "worker_skills.delete"
	WorkerSkillAssessmentsList           = "worker_skill_assessments.list"
	WorkerSkillAssessmentsCreate         = "worker_skill_assessments.create"
	SkillsSummaryView                    = "skills.summary.view"
	ProjectSkillRequirementsList         = "project_skill_requirements.list"
	ProjectSkillRequirementsCreate       = "project_skill_requirements.create"
	ProjectSkillRequirementsView         = "project_skill_requirements.view"
	ProjectSkillRequirementsUpdate       = "project_skill_requirements.update"
	ProjectSkillRequirementsDelete       = "project_skill_requirements.delete"
	ProjectSkillGapsView                 = "project_skill_gaps.view"
	ProjectSkillDependenciesView         = "project_skill_dependencies.view"
	LearningCoursesList                  = "learning.courses.list"
	LearningCoursesManage                = "learning.courses.manage"
	LearningPathsList                    = "learning.paths.list"
	LearningPathsManage                  = "learning.paths.manage"
	LearningEnrollmentsList              = "learning.enrollments.list"
	LearningEnrollmentsAssign            = "learning.enrollments.assign"
	LearningEnrollmentsStatus            = "learning.enrollments.status"
	LearningCertificatesUpload           = "learning.certificates.upload"
	LearningRecommendationsView          = "learning.recommendations.view"
	LearningRecommendationsManage        = "learning.recommendations.manage"
	LearningSummaryView                  = "learning.summary.view"
	TalentMarketplaceOpportunitiesList   = "talent_marketplace.opportunities.list"
	TalentMarketplaceOpportunitiesCreate = "talent_marketplace.opportunities.create"
	TalentMarketplaceOpportunitiesView   = "talent_marketplace.opportunities.view"
	TalentMarketplaceOpportunitiesUpdate = "talent_marketplace.opportunities.update"
	TalentMarketplaceOpportunitiesDelete = "talent_marketplace.opportunities.delete"
	TalentMarketplaceApplicationsList    = "talent_marketplace.applications.list"
	TalentMarketplaceApplicationsCreate  = "talent_marketplace.applications.create"
	TalentMarketplaceApplicationsView    = "talent_marketplace.applications.view"
	TalentMarketplaceApplicationsUpdate  = "talent_marketplace.applications.update"
	TalentMarketplaceRecommendationsView = "talent_marketplace.recommendations.view"
	TalentMarketplaceEventsView          = "talent_marketplace.events.view"
	TalentMarketplaceFallbackManage      = "talent_marketplace.fallback.manage"
	OKRCyclesList                        = "okr.cycles.list"
	OKRCyclesCreate                      = "okr.cycles.create"
	OKRCyclesView                        = "okr.cycles.view"
	OKRCyclesUpdate                      = "okr.cycles.update"
	OKRCyclesDelete                      = "okr.cycles.delete"
	OKRCyclesStatus                      = "okr.cycles.status"
	OKRObjectivesList                    = "okr.objectives.list"
	OKRObjectivesCreate                  = "okr.objectives.create"
	OKRObjectivesView                    = "okr.objectives.view"
	OKRObjectivesUpdate                  = "okr.objectives.update"
	OKRObjectivesDelete                  = "okr.objectives.delete"
	OKRObjectivesStatus                  = "okr.objectives.status"
	OKRKeyResultsList                    = "okr.key_results.list"
	OKRKeyResultsCreate                  = "okr.key_results.create"
	OKRKeyResultsView                    = "okr.key_results.view"
	OKRKeyResultsUpdate                  = "okr.key_results.update"
	OKRKeyResultsDelete                  = "okr.key_results.delete"
	OKRCheckInsList                      = "okr.checkins.list"
	OKRCheckInsCreate                    = "okr.checkins.create"
	OKRSummaryView                       = "okr.summary.view"
	PerformanceCheckInsList              = "performance_checkins.list"
	PerformanceCheckInsCreate            = "performance_checkins.create"
	PerformanceCheckInsView              = "performance_checkins.view"
	PerformanceCheckInsUpdate            = "performance_checkins.update"
	PerformanceCheckInsSubmit            = "performance_checkins.submit"
	PerformanceCheckInsReview            = "performance_checkins.review"
	PerformanceCheckInsDelete            = "performance_checkins.delete"
	PerformanceCheckInsSummary           = "performance_checkins.summary"
	PerformanceCalibrationView           = "performance_calibration.view"
	PerformanceTimelineView              = "performance_timeline.view"
	FeedbackRequestsList                 = "feedback.requests.list"
	FeedbackRequestsCreate               = "feedback.requests.create"
	FeedbackRequestsView                 = "feedback.requests.view"
	FeedbackRequestsUpdate               = "feedback.requests.update"
	FeedbackRequestsStatus               = "feedback.requests.status"
	FeedbackResponsesList                = "feedback.responses.list"
	FeedbackResponsesCreate              = "feedback.responses.create"
	FeedbackResponsesView                = "feedback.responses.view"
	PulseSurveysList                     = "pulse_surveys.list"
	PulseSurveysCreate                   = "pulse_surveys.create"
	PulseSurveysView                     = "pulse_surveys.view"
	PulseSurveysUpdate                   = "pulse_surveys.update"
	PulseSurveysDelete                   = "pulse_surveys.delete"
	PulseSurveysStatus                   = "pulse_surveys.status"
	PulseQuestionsList                   = "pulse_questions.list"
	PulseQuestionsCreate                 = "pulse_questions.create"
	PulseQuestionsUpdate                 = "pulse_questions.update"
	PulseQuestionsDelete                 = "pulse_questions.delete"
	PulseResponsesList                   = "pulse_responses.list"
	PulseResponsesCreate                 = "pulse_responses.create"
	WellbeingScoresList                  = "wellbeing.scores.list"
	WellbeingScoresUpsert                = "wellbeing.scores.upsert"
	WellbeingAlertsList                  = "wellbeing.alerts.list"
	WellbeingAlertsReview                = "wellbeing.alerts.review"
	WellbeingAggregateView               = "wellbeing.view"

	LookupsList   = "lookups.list"
	LookupsCreate = "lookups.create"
	LookupsUpdate = "lookups.update"
	LookupsDelete = "lookups.delete"

	FinancialYearsList      = "financial_years.list"
	FinancialYearsCreate    = "financial_years.create"
	FinancialYearsView      = "financial_years.view"
	FinancialYearsUpdate    = "financial_years.update"
	FinancialYearsDelete    = "financial_years.delete"
	FinancialYearsSetActive = "financial_years.set_active"

	WorkingHoursList   = "working_hours.list"
	WorkingHoursCreate = "working_hours.create"
	WorkingHoursUpdate = "working_hours.update"
	WorkingHoursDelete = "working_hours.delete"
	WorkingHoursCopy   = "working_hours.copy_to_branch"

	HolidaysList   = "holidays.list"
	HolidaysCreate = "holidays.create"
	HolidaysView   = "holidays.view"
	HolidaysUpdate = "holidays.update"
	HolidaysDelete = "holidays.delete"

	PoliciesList    = "policies.list"
	PoliciesCreate  = "policies.create"
	PoliciesView    = "policies.view"
	PoliciesUpdate  = "policies.update"
	PoliciesDelete  = "policies.delete"
	PoliciesPublish = "policies.publish"

	TenantSubscriptionsList   = "tenant_subscriptions.list"
	TenantSubscriptionsCreate = "tenant_subscriptions.create"
	TenantSubscriptionsView   = "tenant_subscriptions.view"
	TenantSubscriptionsUpdate = "tenant_subscriptions.update"
	TenantSubscriptionsDelete = "tenant_subscriptions.delete"
	SubscriptionPlansList     = "subscription_plans.list"
	SubscriptionPlansCreate   = "subscription_plans.create"
	SubscriptionPlansView     = "subscription_plans.view"
	SubscriptionPlansUpdate   = "subscription_plans.update"
	SubscriptionPlansDelete   = "subscription_plans.delete"

	EmployeesList                 = "employees.list"
	EmployeesCreate               = "employees.create"
	EmployeesView                 = "employees.view"
	EmployeesUpdate               = "employees.update"
	EmployeesDelete               = "employees.delete"
	EmployeesDeactivate           = "employees.deactivate"
	EmployeesDocumentsManage      = "employees.documents.manage"
	EmployeesBankManage           = "employees.bank.manage"
	EmployeesStatutoryManage      = "employees.statutory.manage"
	EmployeesCredentialsManage    = "employees.credentials.manage"
	EmployeeLettersList           = "employee_letters.list"
	EmployeeLettersCreate         = "employee_letters.create"
	EmployeeLettersView           = "employee_letters.view"
	EmployeeLettersApprove        = "employee_letters.approve"
	EmployeeLettersSend           = "employee_letters.send"
	EmployeeLettersDownload       = "employee_letters.download"
	EmployeeLettersRevoke         = "employee_letters.revoke"
	EmployeeLetterTemplatesManage = "employee_letter_templates.manage"
	AgreementsList                = "agreements.list"
	AgreementsCreate              = "agreements.create"
	AgreementsView                = "agreements.view"
	AgreementsSend                = "agreements.send"
	AgreementsSign                = "agreements.sign"
	AgreementsRevoke              = "agreements.revoke"
	AgreementsDownload            = "agreements.download"
	AgreementsDelete              = "agreements.delete"
	AgreementsEventsView          = "agreements.events.view"
	AgreementTemplatesManage      = "agreement_templates.manage"

	EmployeeExitsList     = "employee_exits.list"
	EmployeeExitsCreate   = "employee_exits.create"
	EmployeeExitsView     = "employee_exits.view"
	EmployeeExitsApprove  = "employee_exits.approve"
	EmployeeExitsUpdate   = "employee_exits.update"
	EmployeeExitsComplete = "employee_exits.complete"
	EmployeeExitsCancel   = "employee_exits.cancel"

	LeaveTypesList   = "leave_types.list"
	LeaveTypesCreate = "leave_types.create"
	LeaveTypesView   = "leave_types.view"
	LeaveTypesUpdate = "leave_types.update"
	LeaveTypesDelete = "leave_types.delete"

	LeavePoliciesList   = "leave_policies.list"
	LeavePoliciesCreate = "leave_policies.create"
	LeavePoliciesView   = "leave_policies.view"
	LeavePoliciesUpdate = "leave_policies.update"
	LeavePoliciesDelete = "leave_policies.delete"

	LeaveTemplatesList           = "leave_templates.list"
	LeaveTemplatesCreate         = "leave_templates.create"
	LeaveTemplatesUpdate         = "leave_templates.update"
	LeaveTemplatesDelete         = "leave_templates.delete"
	LeaveTemplateRulesManage     = "leave_template_rules.manage"
	LeaveAssignmentsManage       = "leave_assignments.manage"
	LeaveBalancesList            = "leave_balances.list"
	LeaveBalancesUpdate          = "leave_balances.update"
	LeaveLedgerView              = "leave_ledger.view"
	LeaveAccrualRun              = "leave_accrual.run"
	LeaveApprovalWorkflowsList   = "leave_approval_workflows.list"
	LeaveApprovalWorkflowsManage = "leave_approval_workflows.manage"

	LeavesList    = "leaves.list"
	LeavesApply   = "leaves.apply"
	LeavesView    = "leaves.view"
	LeavesApprove = "leaves.approve"
	LeavesReject  = "leaves.reject"
	LeavesCancel  = "leaves.cancel"
	LeavesReport  = "leaves.report"

	AttendanceList                     = "attendance.list"
	AttendanceCheckIn                  = "attendance.check_in"
	AttendanceCheckOut                 = "attendance.check_out"
	AttendanceView                     = "attendance.view"
	AttendanceUpdate                   = "attendance.update"
	AttendanceRegularize               = "attendance.regularize"
	AttendanceReviewRequest            = "attendance.review_request"
	AttendanceReport                   = "attendance.report"
	AttendanceExceptionWorkflowsList   = "attendance_exception_workflows.list"
	AttendanceExceptionWorkflowsManage = "attendance_exception_workflows.manage"
	AttendancePayrollBlockersView      = "attendance_payroll_blockers.view"

	AttendanceLocationsList   = "attendance_locations.list"
	AttendanceLocationsCreate = "attendance_locations.create"
	AttendanceLocationsUpdate = "attendance_locations.update"
	AttendanceLocationsDelete = "attendance_locations.delete"

	AttendanceLocationAssignmentsList   = "attendance_location_assignments.list"
	AttendanceLocationAssignmentsCreate = "attendance_location_assignments.create"
	AttendanceLocationAssignmentsUpdate = "attendance_location_assignments.update"
	AttendanceLocationAssignmentsDelete = "attendance_location_assignments.delete"

	AttendanceDevicesList   = "attendance_devices.list"
	AttendanceDevicesCreate = "attendance_devices.create"
	AttendanceDevicesUpdate = "attendance_devices.update"
	AttendanceDevicesDelete = "attendance_devices.delete"

	EmployeeAttendanceDevicesList   = "employee_attendance_devices.list"
	EmployeeAttendanceDevicesCreate = "employee_attendance_devices.create"
	EmployeeAttendanceDevicesUpdate = "employee_attendance_devices.update"
	EmployeeAttendanceDevicesDelete = "employee_attendance_devices.delete"

	ShiftTemplatesList         = "shift_templates.list"
	ShiftTemplatesManage       = "shift_templates.manage"
	StaffingRequirementsList   = "staffing_requirements.list"
	StaffingRequirementsManage = "staffing_requirements.manage"
	ShiftAssignmentsList       = "shift_assignments.list"
	ShiftAssignmentsManage     = "shift_assignments.manage"
	ShiftAssignmentsPublish    = "shift_assignments.publish"
	ShiftAssignmentsLock       = "shift_assignments.lock"
	ShiftSwapsList             = "shift_swaps.list"
	ShiftSwapsCreate           = "shift_swaps.create"
	ShiftSwapsReview           = "shift_swaps.review"
	ShiftScheduleEventsView    = "shift_schedule_events.view"
	ShiftScheduleSummaryView   = "shift_schedule_summary.view"

	BenefitPlansList         = "benefit_plans.list"
	BenefitPlansManage       = "benefit_plans.manage"
	BenefitWindowsList       = "benefit_windows.list"
	BenefitWindowsManage     = "benefit_windows.manage"
	BenefitDependentsList    = "benefit_dependents.list"
	BenefitDependentsManage  = "benefit_dependents.manage"
	BenefitEnrollmentsList   = "benefit_enrollments.list"
	BenefitEnrollmentsManage = "benefit_enrollments.manage"
	BenefitEnrollmentsReview = "benefit_enrollments.review"
	BenefitClaimTypesList    = "benefit_claim_types.list"
	BenefitClaimTypesManage  = "benefit_claim_types.manage"
	BenefitClaimsList        = "benefit_claims.list"
	BenefitClaimsCreate      = "benefit_claims.create"
	BenefitClaimsReview      = "benefit_claims.review"
	BenefitClaimsPay         = "benefit_claims.pay"
	BenefitClaimsAttach      = "benefit_claims.attach"
	BenefitClaimsExport      = "benefit_claims.export"
	BenefitEventsView        = "benefit_events.view"
	BenefitSummaryView       = "benefits.summary.view"

	PayCyclesView   = "pay_cycles.view"
	PayCyclesUpdate = "pay_cycles.update"

	PayrollImportsList       = "payroll_imports.list"
	PayrollImportsCreate     = "payroll_imports.create"
	PayrollStatutoryRules    = "payroll_statutory_rules.manage"
	PayrollLocksManage       = "payroll_locks.manage"
	PayrollSalarySheetView   = "payroll_salary_sheet.view"
	PayrollSalarySheetExport = "payroll_salary_sheet.export"
	PayrollReconciliation    = "payroll_reconciliation.view"

	PayGroupsList   = "pay_groups.list"
	PayGroupsCreate = "pay_groups.create"
	PayGroupsView   = "pay_groups.view"
	PayGroupsUpdate = "pay_groups.update"
	PayGroupsDelete = "pay_groups.delete"
	PayRunsList     = "pay_runs.list"
	PayRunsCreate   = "pay_runs.create"
	PayRunsView     = "pay_runs.view"
	PayRunsAssess   = "pay_runs.assess"
	PayRunsFreeze   = "pay_runs.freeze"
	PayRunsGenerate = "pay_runs.generate"
	PayRunsLock     = "pay_runs.lock"

	FlexPayRunsList     = "flex_pay_runs.list"
	FlexPayRunsCreate   = "flex_pay_runs.create"
	FlexPayRunsView     = "flex_pay_runs.view"
	FlexPayRunsGenerate = "flex_pay_runs.generate"
	FlexPayRunsSubmit   = "flex_pay_runs.submit"
	FlexPayRunsApprove  = "flex_pay_runs.approve"
	FlexPayRunsReject   = "flex_pay_runs.reject"
	FlexPayRunsPay      = "flex_pay_runs.pay"
	FlexPayRunsExport   = "flex_pay_runs.export"
	FlexPayRunsDelete   = "flex_pay_runs.delete"

	ContractorInvoicesList    = "contractor_invoices.list"
	ContractorInvoicesCreate  = "contractor_invoices.create"
	ContractorInvoicesView    = "contractor_invoices.view"
	ContractorInvoicesUpdate  = "contractor_invoices.update"
	ContractorInvoicesSubmit  = "contractor_invoices.submit"
	ContractorInvoicesApprove = "contractor_invoices.approve"
	ContractorInvoicesReject  = "contractor_invoices.reject"
	ContractorInvoicesPay     = "contractor_invoices.pay"
	ContractorInvoicesDelete  = "contractor_invoices.delete"

	SalaryTemplatesList     = "salary_templates.list"
	SalaryTemplatesCreate   = "salary_templates.create"
	SalaryTemplatesView     = "salary_templates.view"
	SalaryTemplatesUpdate   = "salary_templates.update"
	SalaryTemplatesDelete   = "salary_templates.delete"
	SalaryTemplatesActivate = "salary_templates.activate"

	EmployeeSalariesList   = "employee_salaries.list"
	EmployeeSalariesCreate = "employee_salaries.create"
	EmployeeSalariesView   = "employee_salaries.view"
	EmployeeSalariesUpdate = "employee_salaries.update"
	EmployeeSalariesDelete = "employee_salaries.delete"

	SalarySlipsList       = "salary_slips.list"
	SalarySlipsGenerate   = "salary_slips.generate"
	SalarySlipsView       = "salary_slips.view"
	SalarySlipsRegenerate = "salary_slips.regenerate"
	SalarySlipsDownload   = "salary_slips.download"

	CompensationPayBandsList           = "compensation.pay_bands.list"
	CompensationPayBandsManage         = "compensation.pay_bands.manage"
	CompensationCyclesList             = "compensation.cycles.list"
	CompensationCyclesManage           = "compensation.cycles.manage"
	CompensationBudgetPoolsManage      = "compensation.budget_pools.manage"
	CompensationRecommendationsList    = "compensation.recommendations.list"
	CompensationRecommendationsManage  = "compensation.recommendations.manage"
	CompensationRecommendationsApprove = "compensation.recommendations.approve"
	CompensationPayrollHandoff         = "compensation.payroll_handoff"
	CompensationEquityView             = "compensation.equity.view"
	CompensationEquityManage           = "compensation.equity.manage"
	CompensationEventsView             = "compensation.events.view"
	CompensationSummaryView            = "compensation.summary.view"

	SuccessionCyclesList               = "succession.cycles.list"
	SuccessionCyclesManage             = "succession.cycles.manage"
	SuccessionCriticalRolesList        = "succession.critical_roles.list"
	SuccessionCriticalRolesManage      = "succession.critical_roles.manage"
	SuccessionSuccessorsList           = "succession.successors.list"
	SuccessionSuccessorsManage         = "succession.successors.manage"
	SuccessionDevelopmentActionsList   = "succession.development_actions.list"
	SuccessionDevelopmentActionsManage = "succession.development_actions.manage"
	SuccessionEventsView               = "succession.events.view"
	SuccessionSummaryView              = "succession.summary.view"
	SuccessionConfidentialView         = "succession.confidential.view"
	AssetItemsList                     = "asset_access.assets.list"
	AssetItemsManage                   = "asset_access.assets.manage"
	AccessCatalogList                  = "asset_access.catalog.list"
	AccessCatalogManage                = "asset_access.catalog.manage"
	AssetAssignmentsList               = "asset_access.assignments.list"
	AssetAssignmentsManage             = "asset_access.assignments.manage"
	AccessTasksList                    = "asset_access.tasks.list"
	AccessTasksManage                  = "asset_access.tasks.manage"
	AssetAccessEventsView              = "asset_access.events.view"
	AssetAccessSummaryView             = "asset_access.summary.view"

	DashboardEmployeeView       = "dashboard.employee.view"
	DashboardHRView             = "dashboard.hr.view"
	HRCommandCenterView         = "command_center.view"
	OperationCatalogView        = "operation_catalog.view"
	OperationsWorkbenchView     = "operations_workbench.view"
	TenantOperationsView        = "tenant_operations.view"
	TenantOperationsManage      = "tenant_operations.manage"
	WorkflowDefinitionsList     = "workflow.definitions.list"
	WorkflowDefinitionsManage   = "workflow.definitions.manage"
	OperationTemplatesList      = "workflow.operation_templates.list"
	OperationTemplatesManage    = "workflow.operation_templates.manage"
	WorkflowTasksList           = "workflow.tasks.list"
	WorkflowTasksCreate         = "workflow.tasks.create"
	WorkflowTasksView           = "workflow.tasks.view"
	WorkflowTasksUpdate         = "workflow.tasks.update"
	WorkflowTasksAct            = "workflow.tasks.act"
	WorkflowTasksComment        = "workflow.tasks.comment"
	WorkflowTasksAttach         = "workflow.tasks.attach"
	WorkflowTasksWatch          = "workflow.tasks.watch"
	WorkflowTasksRestrictedView = "workflow.tasks.restricted.view"
	WorkflowTasksSummaryView    = "workflow.tasks.summary.view"
	HRCasesList                 = "hr_cases.list"
	HRCasesCreate               = "hr_cases.create"
	HRCasesView                 = "hr_cases.view"
	HRCasesUpdate               = "hr_cases.update"
	HRCasesAssign               = "hr_cases.assign"
	HRCasesStatus               = "hr_cases.status"
	HRCasesComment              = "hr_cases.comment"
	HRCasesAttach               = "hr_cases.attach"
	HRCasesRestrictedView       = "hr_cases.restricted.view"
	HRCaseCategoriesManage      = "hr_case_categories.manage"
	HRCaseSLAManage             = "hr_case_sla.manage"
	ERCasesList                 = "employee_relations.cases.list"
	ERCasesCreate               = "employee_relations.cases.create"
	ERCasesView                 = "employee_relations.cases.view"
	ERCasesUpdate               = "employee_relations.cases.update"
	ERCasesStatus               = "employee_relations.cases.status"
	ERCasesLegalHold            = "employee_relations.cases.legal_hold"
	ERCaseCategoriesManage      = "employee_relations.categories.manage"
	ERCasePartiesManage         = "employee_relations.parties.manage"
	ERAllegationsManage         = "employee_relations.allegations.manage"
	ERStepsManage               = "employee_relations.steps.manage"
	ERWitnessNotesManage        = "employee_relations.witness_notes.manage"
	EREvidenceManage            = "employee_relations.evidence.manage"
	ERFindingsManage            = "employee_relations.findings.manage"
	ERActionPlansManage         = "employee_relations.action_plans.manage"
	EREventsView                = "employee_relations.events.view"
	ERRestrictedView            = "employee_relations.restricted.view"

	ReportsView     = "reports.view"
	ReportsManage   = "reports.manage"
	ReportsExport   = "reports.export"
	ReportsSchedule = "reports.schedule"

	InsightsView            = "insights.view"
	InsightsRefresh         = "insights.refresh"
	InsightsReview          = "insights.review"
	AIActionsView           = "ai_actions.view"
	AIActionsManage         = "ai_actions.manage"
	AIActionsOverride       = "ai_actions.override"
	AISignalsEmit           = "ai_signals.emit"
	AIAgentsView            = "ai_agents.view"
	AIAgentsRun             = "ai_agents.run"
	PeopleAnalyticsView     = "people_analytics.view"
	PrivacyConsentsManage   = "privacy.consents.manage"
	PrivacyErasureManage    = "privacy.erasure.manage"
	IntegrationHooksManage  = "integrations.hooks.manage"
	MobileConstraintsManage = "mobile.constraints.manage"

	CelebrationTypesManage = "celebration_types.manage"
	CelebrationsList       = "celebrations.list"
	CelebrationsCreate     = "celebrations.create"
	CelebrationsView       = "celebrations.view"
	CelebrationsUpdate     = "celebrations.update"
	CelebrationsDelete     = "celebrations.delete"
	CelebrationsSend       = "celebrations.send"

	NotificationsList               = "notifications.list"
	NotificationsSend               = "notifications.send"
	NotificationsRead               = "notifications.read"
	NotificationsPreferences        = "notifications.preferences"
	NotificationsMastersManage      = "notifications.masters.manage"
	NotificationsDeviceTokensManage = "notifications.device_tokens.manage"

	JobPositionsList   = "job_positions.list"
	JobPositionsCreate = "job_positions.create"
	JobPositionsView   = "job_positions.view"
	JobPositionsUpdate = "job_positions.update"
	JobPositionsDelete = "job_positions.delete"

	JobRequisitionsList    = "job_requisitions.list"
	JobRequisitionsCreate  = "job_requisitions.create"
	JobRequisitionsView    = "job_requisitions.view"
	JobRequisitionsUpdate  = "job_requisitions.update"
	JobRequisitionsApprove = "job_requisitions.approve"
	JobRequisitionsReject  = "job_requisitions.reject"
	JobRequisitionsClose   = "job_requisitions.close"

	JobPostingsList    = "job_postings.list"
	JobPostingsCreate  = "job_postings.create"
	JobPostingsView    = "job_postings.view"
	JobPostingsUpdate  = "job_postings.update"
	JobPostingsPublish = "job_postings.publish"
	JobPostingsClose   = "job_postings.close"

	CandidatesList   = "candidates.list"
	CandidatesCreate = "candidates.create"
	CandidatesView   = "candidates.view"
	CandidatesUpdate = "candidates.update"
	CandidatesDelete = "candidates.delete"

	CandidateApplicationsList   = "candidate_applications.list"
	CandidateApplicationsCreate = "candidate_applications.create"
	CandidateApplicationsView   = "candidate_applications.view"
	CandidateApplicationsUpdate = "candidate_applications.update"
	CandidateApplicationsMove   = "candidate_applications.move_stage"
	ApplicantPortalView         = "applicant.portal.view"

	InterviewRoundsList   = "interview_rounds.list"
	InterviewRoundsCreate = "interview_rounds.create"
	InterviewRoundsView   = "interview_rounds.view"
	InterviewRoundsUpdate = "interview_rounds.update"
	InterviewRoundsDelete = "interview_rounds.delete"

	OfferLettersList   = "offer_letters.list"
	OfferLettersCreate = "offer_letters.create"
	OfferLettersView   = "offer_letters.view"
	OfferLettersSend   = "offer_letters.send"
	OfferLettersUpdate = "offer_letters.update"
	OfferLettersRevoke = "offer_letters.revoke"

	OnboardingWorkflowsManage = "onboarding_workflows.manage"
	OnboardingList            = "onboarding.list"
	OnboardingStart           = "onboarding.start"
	OnboardingView            = "onboarding.view"
	OnboardingUpdate          = "onboarding.update"
	OnboardingCompleteTask    = "onboarding.complete_task"

	IntegrationsEmailManage    = "integrations.email.manage"
	IntegrationsSMSManage      = "integrations.sms.manage"
	IntegrationsWhatsAppManage = "integrations.whatsapp.manage"
	IntegrationsStorageManage  = "integrations.storage.manage"
	IntegrationsPushManage     = "integrations.push.manage"

	ScheduledJobsList = "scheduled_jobs.list"
	ScheduledJobsRun  = "scheduled_jobs.run"
)

type Permission struct {
	Key         string
	Description string
}

type RoleTemplate struct {
	Code        string
	Name        string
	Description string
	Permissions []string
}

var Catalog = []Permission{
	{TenantProfilesView, "View tenant HRMS profile."},
	{TenantProfilesUpdate, "Update tenant HRMS profile."},
	{TenantProvision, "Provision HRMS defaults after identity tenant creation."},
	{TenantSettingsList, "List tenant HRMS settings."},
	{TenantSettingsUpdate, "Update tenant HRMS settings."},
	{BrandingView, "View company branding."},
	{BrandingUpdate, "Update company branding."},
	{BranchesList, "List branches."},
	{BranchesCreate, "Create branches."},
	{BranchesView, "View branch details."},
	{BranchesUpdate, "Update branches."},
	{BranchesDelete, "Deactivate branches."},
	{DepartmentsList, "List departments."},
	{DepartmentsCreate, "Create departments."},
	{DepartmentsView, "View department details."},
	{DepartmentsUpdate, "Update departments."},
	{DepartmentsDelete, "Deactivate departments."},
	{DesignationsList, "List designations."},
	{DesignationsCreate, "Create designations."},
	{DesignationsView, "View designation details."},
	{DesignationsUpdate, "Update designations."},
	{DesignationsDelete, "Deactivate designations."},
	{DesignationsMastersManage, "Manage designation level and seniority master values."},
	{DesignationsAttendanceRequirementUpdate, "Update whether a designation requires attendance punches."},
	{WorkerTypesList, "List workforce type taxonomy records."},
	{WorkerTypesCreate, "Create workforce type taxonomy records."},
	{WorkerTypesView, "View workforce type taxonomy details."},
	{WorkerTypesUpdate, "Update workforce type taxonomy, attendance, pay, and compliance defaults."},
	{WorkerTypesDelete, "Deactivate workforce type taxonomy records."},
	{WorkerClassificationRulesList, "List workforce classification rules."},
	{WorkerClassificationRulesManage, "Create, update, and deactivate workforce classification rules."},
	{WorkersList, "List workforce hub worker profiles across employees and contingent workers."},
	{WorkersCreate, "Create worker profiles and link them to employees where applicable."},
	{WorkersView, "View worker profile details, classification, compliance, and payroll readiness."},
	{WorkersUpdate, "Update worker profiles, status, organization placement, and compliance readiness."},
	{WorkersDelete, "Deactivate worker profiles."},
	{EngagementsList, "List worker engagements, assignments, projects, retainers, and contract commitments."},
	{EngagementsCreate, "Create worker engagements with dates, cost centers, budgets, rates, and renewal tracking."},
	{EngagementsView, "View worker engagement details and lifecycle state."},
	{EngagementsUpdate, "Update worker engagements, organization placement, budget, rate, and renewal details."},
	{EngagementsStatus, "Pause, activate, complete, terminate, or cancel worker engagements."},
	{EngagementsDelete, "Deactivate worker engagements."},
	{WorkLogsList, "List hourly and contractor work logs."},
	{WorkLogsCreate, "Create hourly and contractor work logs."},
	{WorkLogsView, "View work log details."},
	{WorkLogsUpdate, "Update draft or rejected work logs."},
	{WorkLogsSubmit, "Submit work logs for approval."},
	{WorkLogsApprove, "Approve submitted work logs for payroll and analytics."},
	{WorkLogsReject, "Reject submitted work logs with review comments."},
	{WorkLogsDelete, "Deactivate draft or rejected work logs."},
	{WorkLogsReport, "View work log rollups and approved-hour reports."},
	{ProjectsList, "List projects with budget, owner, due-date, and milestone progress summaries."},
	{ProjectsCreate, "Create HRMS projects for employee and flexible workforce delivery."},
	{ProjectsView, "View project details, budget, department, manager, and status."},
	{ProjectsUpdate, "Update project metadata, ownership, budget, dates, and status fields."},
	{ProjectsStatus, "Activate, pause, complete, or cancel projects."},
	{ProjectsDelete, "Deactivate projects."},
	{MilestonesList, "List project milestones and acceptance queues."},
	{MilestonesCreate, "Create project milestones with acceptance criteria and payment trigger metadata."},
	{MilestonesView, "View project milestone details."},
	{MilestonesUpdate, "Update draft, open, or rejected project milestones."},
	{MilestonesSubmit, "Submit project milestones for acceptance."},
	{MilestonesApprove, "Accept submitted project milestones."},
	{MilestonesReject, "Reject submitted project milestones with review comments."},
	{MilestonesDelete, "Deactivate project milestones that are not accepted."},
	{MilestonesEventsView, "View project milestone audit events."},
	{ComplianceRulesList, "List compliance rule definitions."},
	{ComplianceRulesCreate, "Create compliance rules for workers and engagements."},
	{ComplianceRulesView, "View compliance rule details."},
	{ComplianceRulesUpdate, "Update compliance rule applicability, evidence, and payroll-blocking settings."},
	{ComplianceRulesDelete, "Deactivate compliance rules."},
	{ComplianceRulesSeed, "Seed default India workforce compliance rules."},
	{ComplianceChecklistList, "List worker and engagement compliance checklist items."},
	{ComplianceChecklistGenerate, "Generate applicable compliance checklist items."},
	{ComplianceChecklistReview, "Mark checklist items compliant, non-compliant, or not applicable."},
	{ComplianceChecklistEvidence, "Attach compliance evidence references."},
	{ComplianceChecklistWaive, "Waive compliance checklist items with audit reason."},
	{ComplianceChecklistDelete, "Deactivate compliance checklist items."},
	{ComplianceEventsView, "View compliance audit events."},
	{SkillCategoriesList, "List global and tenant skill categories."},
	{SkillCategoriesCreate, "Create tenant skill categories."},
	{SkillCategoriesUpdate, "Update tenant skill categories."},
	{SkillCategoriesDelete, "Deactivate tenant skill categories."},
	{SkillsList, "List global and tenant skill catalog records."},
	{SkillsCreate, "Create tenant skill catalog records."},
	{SkillsView, "View skill catalog details."},
	{SkillsUpdate, "Update tenant skills, requirements, and active status."},
	{SkillsDelete, "Deactivate tenant skills."},
	{WorkerSkillsList, "List worker skill profiles, proficiency, certificates, and verification status."},
	{WorkerSkillsCreate, "Add skills to worker profiles."},
	{WorkerSkillsUpdate, "Update worker skill proficiency, experience, certificates, and assessment data."},
	{WorkerSkillsVerify, "Endorse, verify, reject, or expire worker skills."},
	{WorkerSkillsDelete, "Deactivate worker skills."},
	{WorkerSkillAssessmentsList, "List worker skill assessments and evidence references."},
	{WorkerSkillAssessmentsCreate, "Create worker skill assessment records."},
	{SkillsSummaryView, "View skills inventory summary and expiring certificate counts."},
	{ProjectSkillRequirementsList, "List project and engagement skill requirements."},
	{ProjectSkillRequirementsCreate, "Create project and engagement skill requirements."},
	{ProjectSkillRequirementsView, "View project skill requirement details."},
	{ProjectSkillRequirementsUpdate, "Update project and engagement skill requirements."},
	{ProjectSkillRequirementsDelete, "Deactivate project and engagement skill requirements."},
	{ProjectSkillGapsView, "View project, engagement, and organisation skill gaps."},
	{ProjectSkillDependenciesView, "View single-person skill dependency risks."},
	{LearningCoursesList, "List learning catalog courses, compliance training, and AI-readiness courses."},
	{LearningCoursesManage, "Create, update, and deactivate learning catalog courses."},
	{LearningPathsList, "List learning paths connected to onboarding, compliance, skill gaps, and AI upskilling."},
	{LearningPathsManage, "Create, update, deactivate, and sequence learning paths and path courses."},
	{LearningEnrollmentsList, "List assigned, nominated, in-progress, overdue, and completed learning enrollments."},
	{LearningEnrollmentsAssign, "Assign or nominate workers for courses and learning paths."},
	{LearningEnrollmentsStatus, "Update learning enrollment lifecycle status and completion score."},
	{LearningCertificatesUpload, "Upload learning completion certificates to tenant storage."},
	{LearningRecommendationsView, "View learning recommendations from skill gaps, compliance, performance, and AI signals."},
	{LearningRecommendationsManage, "Create, generate, dismiss, accept, assign, and complete learning recommendations."},
	{LearningSummaryView, "View learning dashboard summary metrics."},
	{TalentMarketplaceOpportunitiesList, "List internal talent marketplace opportunities."},
	{TalentMarketplaceOpportunitiesCreate, "Create internal project, role, gig, mentorship, and backfill opportunities."},
	{TalentMarketplaceOpportunitiesView, "View internal talent marketplace opportunity details."},
	{TalentMarketplaceOpportunitiesUpdate, "Update opportunity status, visibility, staffing details, and fallback settings."},
	{TalentMarketplaceOpportunitiesDelete, "Deactivate internal talent marketplace opportunities."},
	{TalentMarketplaceApplicationsList, "List marketplace applications, recommendations, invitations, and worker decisions."},
	{TalentMarketplaceApplicationsCreate, "Create marketplace applications or invitations for workers."},
	{TalentMarketplaceApplicationsView, "View marketplace application details."},
	{TalentMarketplaceApplicationsUpdate, "Update marketplace application status, manager notes, and worker decisions."},
	{TalentMarketplaceRecommendationsView, "View skills-based worker recommendations for marketplace opportunities."},
	{TalentMarketplaceEventsView, "View marketplace opportunity and application event history."},
	{TalentMarketplaceFallbackManage, "Manage candidate fallback status for opportunities that need external hiring."},
	{OKRCyclesList, "List OKR cycles."},
	{OKRCyclesCreate, "Create OKR cycles."},
	{OKRCyclesView, "View OKR cycle details."},
	{OKRCyclesUpdate, "Update OKR cycle dates, cadence, metadata, and status."},
	{OKRCyclesDelete, "Deactivate OKR cycles."},
	{OKRCyclesStatus, "Activate, close, archive, or reopen OKR cycles."},
	{OKRObjectivesList, "List company, department, project, and worker objectives."},
	{OKRObjectivesCreate, "Create company, department, project, and worker objectives."},
	{OKRObjectivesView, "View objective details, owner, hierarchy, and progress."},
	{OKRObjectivesUpdate, "Update objective scope, owner, weight, dates, and progress."},
	{OKRObjectivesDelete, "Deactivate objectives."},
	{OKRObjectivesStatus, "Change objective lifecycle status."},
	{OKRKeyResultsList, "List measurable key results."},
	{OKRKeyResultsCreate, "Create measurable key results."},
	{OKRKeyResultsView, "View key result details and latest check-in."},
	{OKRKeyResultsUpdate, "Update key result targets, weights, confidence, and status."},
	{OKRKeyResultsDelete, "Deactivate key results."},
	{OKRCheckInsList, "List key result check-ins."},
	{OKRCheckInsCreate, "Create key result progress check-ins."},
	{OKRSummaryView, "View OKR summary rollups by owner type."},
	{PerformanceCheckInsList, "List worker performance check-ins."},
	{PerformanceCheckInsCreate, "Create weekly worker performance check-ins."},
	{PerformanceCheckInsView, "View performance check-in details."},
	{PerformanceCheckInsUpdate, "Update draft performance check-ins."},
	{PerformanceCheckInsSubmit, "Submit performance check-ins for manager review."},
	{PerformanceCheckInsReview, "Review, score, and calibrate performance check-ins."},
	{PerformanceCheckInsDelete, "Deactivate performance check-ins."},
	{PerformanceCheckInsSummary, "View performance check-in mood, status, and score summaries."},
	{PerformanceCalibrationView, "View calibration-ready performance extraction rows."},
	{PerformanceTimelineView, "View worker performance timeline events."},
	{FeedbackRequestsList, "List 360 and performance feedback requests."},
	{FeedbackRequestsCreate, "Create 360 and performance feedback requests."},
	{FeedbackRequestsView, "View feedback request details and response counts."},
	{FeedbackRequestsUpdate, "Update feedback request anonymity, visibility, due date, and prompt."},
	{FeedbackRequestsStatus, "Change feedback request status."},
	{FeedbackResponsesList, "List feedback responses with anonymity rules applied."},
	{FeedbackResponsesCreate, "Submit feedback responses."},
	{FeedbackResponsesView, "View feedback response details."},
	{PulseSurveysList, "List pulse, wellbeing, and engagement surveys."},
	{PulseSurveysCreate, "Create pulse and wellbeing surveys."},
	{PulseSurveysView, "View pulse survey details."},
	{PulseSurveysUpdate, "Update pulse survey audience, consent, privacy, and schedule settings."},
	{PulseSurveysDelete, "Deactivate pulse surveys."},
	{PulseSurveysStatus, "Activate, close, archive, or reopen pulse surveys."},
	{PulseQuestionsList, "List pulse survey questions."},
	{PulseQuestionsCreate, "Create pulse survey questions."},
	{PulseQuestionsUpdate, "Update pulse survey questions."},
	{PulseQuestionsDelete, "Deactivate pulse survey questions."},
	{PulseResponsesList, "List pulse responses with anonymity rules applied."},
	{PulseResponsesCreate, "Submit pulse survey responses."},
	{WellbeingScoresList, "List employee wellbeing scores."},
	{WellbeingScoresUpsert, "Create or update employee wellbeing scores."},
	{WellbeingAlertsList, "List HR-only wellbeing alerts."},
	{WellbeingAlertsReview, "Acknowledge, resolve, or dismiss wellbeing alerts."},
	{WellbeingAggregateView, "View aggregate wellbeing survey trends with anonymity thresholds."},
	{LookupsList, "List HRMS lookup values."},
	{LookupsCreate, "Create HRMS lookup values."},
	{LookupsUpdate, "Update HRMS lookup values."},
	{LookupsDelete, "Deactivate HRMS lookup values."},
	{FinancialYearsList, "List financial years."},
	{FinancialYearsCreate, "Create financial years."},
	{FinancialYearsView, "View financial year details."},
	{FinancialYearsUpdate, "Update financial years."},
	{FinancialYearsDelete, "Deactivate financial years."},
	{FinancialYearsSetActive, "Set active financial year."},
	{WorkingHoursList, "List working hours."},
	{WorkingHoursCreate, "Create working hours."},
	{WorkingHoursUpdate, "Update working hours."},
	{WorkingHoursDelete, "Deactivate working hours."},
	{WorkingHoursCopy, "Copy tenant working hours to a branch."},
	{HolidaysList, "List holidays."},
	{HolidaysCreate, "Create holidays."},
	{HolidaysView, "View holiday details."},
	{HolidaysUpdate, "Update holidays."},
	{HolidaysDelete, "Deactivate holidays."},
	{PoliciesList, "List company policies."},
	{PoliciesCreate, "Create company policies."},
	{PoliciesView, "View company policy details."},
	{PoliciesUpdate, "Update company policies."},
	{PoliciesDelete, "Deactivate company policies."},
	{PoliciesPublish, "Publish company policies."},
	{TenantSubscriptionsList, "List tenant subscriptions."},
	{TenantSubscriptionsCreate, "Create tenant subscriptions."},
	{TenantSubscriptionsView, "View tenant subscription details."},
	{TenantSubscriptionsUpdate, "Update tenant subscriptions."},
	{TenantSubscriptionsDelete, "Deactivate tenant subscriptions."},
	{SubscriptionPlansList, "List subscription plans."},
	{SubscriptionPlansCreate, "Create subscription plans."},
	{SubscriptionPlansView, "View subscription plan details."},
	{SubscriptionPlansUpdate, "Update subscription plans."},
	{SubscriptionPlansDelete, "Deactivate subscription plans."},
	{EmployeesList, "List employees."},
	{EmployeesCreate, "Create employees."},
	{EmployeesView, "View employee details."},
	{EmployeesUpdate, "Update employees."},
	{EmployeesDelete, "Deactivate employees."},
	{EmployeesDeactivate, "Deactivate employees and linked identity users."},
	{EmployeesDocumentsManage, "Manage employee documents."},
	{EmployeesBankManage, "Manage employee bank details."},
	{EmployeesStatutoryManage, "Manage employee statutory details."},
	{EmployeesCredentialsManage, "Resend employee credentials and reset temporary passwords."},
	{EmployeeLettersList, "List employee appointment, experience, and relieving letters."},
	{EmployeeLettersCreate, "Generate employee letters."},
	{EmployeeLettersView, "View employee letter details and rendered previews."},
	{EmployeeLettersApprove, "Approve employee letters before sending."},
	{EmployeeLettersSend, "Send employee letters and request signatures."},
	{EmployeeLettersDownload, "Download employee letter PDFs."},
	{EmployeeLettersRevoke, "Revoke employee letters."},
	{EmployeeLetterTemplatesManage, "Manage employee letter templates."},
	{AgreementsList, "List SOW, NDA, retainer, freelance, internship, and amendment agreements."},
	{AgreementsCreate, "Generate agreements from reusable templates with worker, engagement, and project links."},
	{AgreementsView, "View agreement details, rendered content, and signature state."},
	{AgreementsSend, "Send agreements for signer action."},
	{AgreementsSign, "Sign assigned agreements."},
	{AgreementsRevoke, "Revoke agreements before completion."},
	{AgreementsDownload, "Download agreement PDFs."},
	{AgreementsDelete, "Deactivate draft or obsolete agreements."},
	{AgreementsEventsView, "View agreement lifecycle audit events."},
	{AgreementTemplatesManage, "Manage reusable agreement templates."},
	{EmployeeExitsList, "List employee exit workflows."},
	{EmployeeExitsCreate, "Initiate employee exit workflows."},
	{EmployeeExitsView, "View employee exit workflow details."},
	{EmployeeExitsApprove, "Approve or reject employee exit workflows."},
	{EmployeeExitsUpdate, "Update employee exit checklist tasks."},
	{EmployeeExitsComplete, "Complete employee exits and deactivate access."},
	{EmployeeExitsCancel, "Cancel employee exit workflows."},
	{LeaveTypesList, "List leave types."},
	{LeaveTypesCreate, "Create leave types."},
	{LeaveTypesView, "View leave type details."},
	{LeaveTypesUpdate, "Update leave types."},
	{LeaveTypesDelete, "Deactivate leave types."},
	{LeavePoliciesList, "List leave policies."},
	{LeavePoliciesCreate, "Create leave policies."},
	{LeavePoliciesView, "View leave policy details."},
	{LeavePoliciesUpdate, "Update leave policies."},
	{LeavePoliciesDelete, "Deactivate leave policies."},
	{LeaveTemplatesList, "List leave policy templates."},
	{LeaveTemplatesCreate, "Create leave policy templates."},
	{LeaveTemplatesUpdate, "Update leave policy templates."},
	{LeaveTemplatesDelete, "Deactivate leave policy templates."},
	{LeaveTemplateRulesManage, "Manage leave policy template rules."},
	{LeaveAssignmentsManage, "Assign leave templates to employees."},
	{LeaveBalancesList, "List leave balances."},
	{LeaveBalancesUpdate, "Update leave balances."},
	{LeaveLedgerView, "View leave ledger entries."},
	{LeaveAccrualRun, "Run leave accrual jobs."},
	{LeaveApprovalWorkflowsList, "List leave approval workflows."},
	{LeaveApprovalWorkflowsManage, "Manage configurable leave approval workflows."},
	{LeavesList, "List leave requests."},
	{LeavesApply, "Apply for leave."},
	{LeavesView, "View leave request details."},
	{LeavesApprove, "Approve leave requests."},
	{LeavesReject, "Reject leave requests."},
	{LeavesCancel, "Cancel leave requests."},
	{LeavesReport, "View leave reports."},
	{AttendanceList, "List attendance records."},
	{AttendanceCheckIn, "Check in for attendance."},
	{AttendanceCheckOut, "Check out for attendance."},
	{AttendanceView, "View attendance details."},
	{AttendanceUpdate, "Update attendance records."},
	{AttendanceRegularize, "Create attendance regularisation requests."},
	{AttendanceReviewRequest, "Approve or reject attendance regularisation requests."},
	{AttendanceReport, "View attendance reports."},
	{AttendanceExceptionWorkflowsList, "List attendance exception workflow configurations."},
	{AttendanceExceptionWorkflowsManage, "Manage attendance exception workflows and payroll blocker rules."},
	{AttendancePayrollBlockersView, "View unresolved attendance exceptions that block payroll locking."},
	{AttendanceLocationsList, "List attendance locations."},
	{AttendanceLocationsCreate, "Create attendance locations."},
	{AttendanceLocationsUpdate, "Update attendance locations."},
	{AttendanceLocationsDelete, "Deactivate attendance locations."},
	{AttendanceLocationAssignmentsList, "List attendance location assignments."},
	{AttendanceLocationAssignmentsCreate, "Create attendance location assignments."},
	{AttendanceLocationAssignmentsUpdate, "Update attendance location assignments."},
	{AttendanceLocationAssignmentsDelete, "Deactivate attendance location assignments."},
	{AttendanceDevicesList, "List attendance devices."},
	{AttendanceDevicesCreate, "Create attendance devices."},
	{AttendanceDevicesUpdate, "Update attendance devices."},
	{AttendanceDevicesDelete, "Deactivate attendance devices."},
	{EmployeeAttendanceDevicesList, "List employee attendance device mappings."},
	{EmployeeAttendanceDevicesCreate, "Create employee attendance device mappings."},
	{EmployeeAttendanceDevicesUpdate, "Update employee attendance device mappings."},
	{EmployeeAttendanceDevicesDelete, "Deactivate employee attendance device mappings."},
	{ShiftTemplatesList, "List reusable shift templates."},
	{ShiftTemplatesManage, "Create, update, and deactivate shift templates."},
	{StaffingRequirementsList, "List staffing requirements and coverage rules."},
	{StaffingRequirementsManage, "Manage staffing requirements and payroll blocking rules."},
	{ShiftAssignmentsList, "List shift schedule assignments."},
	{ShiftAssignmentsManage, "Create and update shift schedule assignments."},
	{ShiftAssignmentsPublish, "Publish shift schedules to employees."},
	{ShiftAssignmentsLock, "Lock finalized shift schedules for attendance and payroll."},
	{ShiftSwapsList, "List shift swap requests."},
	{ShiftSwapsCreate, "Create shift swap requests."},
	{ShiftSwapsReview, "Approve or reject shift swap requests."},
	{ShiftScheduleEventsView, "View shift schedule audit events."},
	{ShiftScheduleSummaryView, "View shift schedule conflicts, gaps, and payroll blockers."},
	{BenefitPlansList, "List benefit plans, insurance policies, and reimbursement programs."},
	{BenefitPlansManage, "Create, update, and deactivate benefit plans and eligibility metadata."},
	{BenefitWindowsList, "List benefit enrollment windows."},
	{BenefitWindowsManage, "Create, update, close, and archive benefit enrollment windows."},
	{BenefitDependentsList, "List employee dependents and nominees for benefits."},
	{BenefitDependentsManage, "Create, update, and deactivate employee dependents and nominees."},
	{BenefitEnrollmentsList, "List benefit enrollments by employee, plan, and status."},
	{BenefitEnrollmentsManage, "Create and update benefit enrollments."},
	{BenefitEnrollmentsReview, "Approve or reject employee benefit enrollments."},
	{BenefitClaimTypesList, "List benefit reimbursement and claim types."},
	{BenefitClaimTypesManage, "Create, update, and deactivate benefit claim types, limits, and payroll components."},
	{BenefitClaimsList, "List employee benefit and reimbursement claims."},
	{BenefitClaimsCreate, "Submit benefit and reimbursement claims."},
	{BenefitClaimsReview, "Review, approve, reject, or cancel benefit claims."},
	{BenefitClaimsPay, "Mark approved benefit claims as paid."},
	{BenefitClaimsAttach, "Attach benefit claim evidence through tenant storage."},
	{BenefitClaimsExport, "Mark benefit claims ready or exported for payroll."},
	{BenefitEventsView, "View benefit and claim audit events."},
	{BenefitSummaryView, "View benefit plan, enrollment, claim, payable, and payroll-ready summary metrics."},
	{PayCyclesView, "View pay cycle configuration."},
	{PayCyclesUpdate, "Update pay cycle configuration."},
	{PayrollImportsList, "List payroll import batches and row results."},
	{PayrollImportsCreate, "Upload and apply payroll import data."},
	{PayrollStatutoryRules, "Manage payroll statutory PT and LWF rules."},
	{PayrollLocksManage, "Lock and unlock payroll periods with audit trail."},
	{PayrollSalarySheetView, "View consolidated salary sheets."},
	{PayrollSalarySheetExport, "Export consolidated salary sheets."},
	{PayrollReconciliation, "View payroll attendance and LOP reconciliation."},
	{PayGroupsList, "List payroll pay groups."},
	{PayGroupsCreate, "Create payroll pay groups."},
	{PayGroupsView, "View payroll pay group details and membership preview."},
	{PayGroupsUpdate, "Update payroll pay groups and membership overrides."},
	{PayGroupsDelete, "Deactivate payroll pay groups."},
	{PayRunsList, "List phased payroll pay runs."},
	{PayRunsCreate, "Create phased payroll pay runs for a pay group."},
	{PayRunsView, "View phased payroll pay run readiness and audit details."},
	{PayRunsAssess, "Assess payroll pay run readiness."},
	{PayRunsFreeze, "Freeze attendance, LOP, and adjustments for a payroll pay run."},
	{PayRunsGenerate, "Generate salary slips for a phased payroll pay run."},
	{PayRunsLock, "Lock or unlock phased payroll pay runs with audit trail."},
	{FlexPayRunsList, "List flexible payroll runs for non-employee and contingent workers."},
	{FlexPayRunsCreate, "Create flexible payroll runs for hourly, milestone, retainer, stipend, and invoice payments."},
	{FlexPayRunsView, "View flexible payroll run invoices, items, totals, and lifecycle events."},
	{FlexPayRunsGenerate, "Generate flexible payroll invoices from approved work logs and accepted milestones."},
	{FlexPayRunsSubmit, "Submit flexible payroll runs for approval."},
	{FlexPayRunsApprove, "Approve flexible payroll runs for payment processing."},
	{FlexPayRunsReject, "Reject flexible payroll runs with review comments."},
	{FlexPayRunsPay, "Mark flexible payroll runs as payment pending or paid."},
	{FlexPayRunsExport, "Export flexible payroll run payment data."},
	{FlexPayRunsDelete, "Deactivate draft, generated, or rejected flexible payroll runs."},
	{ContractorInvoicesList, "List contractor invoices."},
	{ContractorInvoicesCreate, "Create contractor invoices with TDS and GST details."},
	{ContractorInvoicesView, "View contractor invoice details."},
	{ContractorInvoicesUpdate, "Update draft or rejected contractor invoices."},
	{ContractorInvoicesSubmit, "Submit contractor invoices for approval."},
	{ContractorInvoicesApprove, "Approve contractor invoices for payment."},
	{ContractorInvoicesReject, "Reject contractor invoices with review comments."},
	{ContractorInvoicesPay, "Mark contractor invoices as paid."},
	{ContractorInvoicesDelete, "Deactivate draft or rejected contractor invoices."},
	{SalaryTemplatesList, "List salary templates."},
	{SalaryTemplatesCreate, "Create salary templates."},
	{SalaryTemplatesView, "View salary template details."},
	{SalaryTemplatesUpdate, "Update salary templates."},
	{SalaryTemplatesDelete, "Deactivate salary templates."},
	{SalaryTemplatesActivate, "Activate salary templates."},
	{EmployeeSalariesList, "List employee salary assignments."},
	{EmployeeSalariesCreate, "Create employee salary assignments."},
	{EmployeeSalariesView, "View employee salary assignments."},
	{EmployeeSalariesUpdate, "Update employee salary assignments."},
	{EmployeeSalariesDelete, "Deactivate employee salary assignments."},
	{SalarySlipsList, "List salary slips."},
	{SalarySlipsGenerate, "Generate salary slips."},
	{SalarySlipsView, "View salary slip details."},
	{SalarySlipsRegenerate, "Regenerate salary slips."},
	{SalarySlipsDownload, "Download salary slips."},
	{CompensationPayBandsList, "List compensation pay bands and salary ranges."},
	{CompensationPayBandsManage, "Create, update, and deactivate compensation pay bands."},
	{CompensationCyclesList, "List compensation review cycles."},
	{CompensationCyclesManage, "Create, update, submit, approve, finalize, or cancel compensation review cycles."},
	{CompensationBudgetPoolsManage, "Manage compensation review budget pools and allocations."},
	{CompensationRecommendationsList, "List employee compensation recommendations."},
	{CompensationRecommendationsManage, "Create and update compensation recommendations."},
	{CompensationRecommendationsApprove, "Approve, reject, or finalize compensation recommendations."},
	{CompensationPayrollHandoff, "Mark finalized compensation recommendations ready for payroll handoff."},
	{CompensationEquityView, "View compensation equity checks and band-position warnings."},
	{CompensationEquityManage, "Generate, acknowledge, resolve, or waive compensation equity checks."},
	{CompensationEventsView, "View compensation review audit events."},
	{CompensationSummaryView, "View compensation review summary metrics."},
	{SuccessionCyclesList, "List succession review cycles."},
	{SuccessionCyclesManage, "Create, update, close, archive, or cancel succession review cycles."},
	{SuccessionCriticalRolesList, "List critical roles, incumbents, emergency cover, and risk ratings."},
	{SuccessionCriticalRolesManage, "Create, update, deactivate, and change status for critical roles."},
	{SuccessionSuccessorsList, "List successor nominations, readiness levels, and readiness pipeline."},
	{SuccessionSuccessorsManage, "Nominate, update, approve, reject, or withdraw successors."},
	{SuccessionDevelopmentActionsList, "List successor development actions linked to learning and role readiness."},
	{SuccessionDevelopmentActionsManage, "Create, update, and complete successor development actions."},
	{SuccessionEventsView, "View succession planning audit events."},
	{SuccessionSummaryView, "View succession planning summary and coverage metrics."},
	{SuccessionConfidentialView, "Access confidential HR-only succession planning data."},
	{AssetItemsList, "List tenant assets such as laptops, ID cards, SIMs, and equipment."},
	{AssetItemsManage, "Create, update, deactivate, and change status for tenant assets."},
	{AccessCatalogList, "List software, system, facility, and role access catalog items."},
	{AccessCatalogManage, "Create, update, and deactivate access catalog items and defaults."},
	{AssetAssignmentsList, "List asset issue, return, damage, loss, and clearance records."},
	{AssetAssignmentsManage, "Issue, approve, return, mark damage/loss, and clear assigned assets."},
	{AccessTasksList, "List access provisioning, deprovisioning, review, and change tasks."},
	{AccessTasksManage, "Create, approve, complete, revoke, or block access lifecycle tasks."},
	{AssetAccessEventsView, "View asset and access lifecycle audit history."},
	{AssetAccessSummaryView, "View asset and access lifecycle summary metrics."},
	{DashboardEmployeeView, "View employee dashboard."},
	{DashboardHRView, "View HR dashboard."},
	{HRCommandCenterView, "View the HR command center landing workspace."},
	{OperationCatalogView, "View the role-scoped global operation catalog."},
	{OperationsWorkbenchView, "View the HR operations workbench and unified action queue."},
	{TenantOperationsView, "View governed platform tenant operation requests."},
	{TenantOperationsManage, "Create and act on governed tenant operation requests."},
	{WorkflowDefinitionsList, "List workflow definitions for task routing and approvals."},
	{WorkflowDefinitionsManage, "Create and update workflow definitions and workflow steps."},
	{OperationTemplatesList, "List reusable operation templates for request and approval workflows."},
	{OperationTemplatesManage, "Create and update operation templates and launch schemas."},
	{WorkflowTasksList, "List workflow task inbox, request, team, watch, completed, and delegated queues."},
	{WorkflowTasksCreate, "Create workflow tasks from operation templates or manual requests."},
	{WorkflowTasksView, "View workflow task details, timeline, comments, and attachments."},
	{WorkflowTasksUpdate, "Update workflow task metadata, assignment, priority, severity, and due dates."},
	{WorkflowTasksAct, "Approve, reject, request information, delegate, open, or complete workflow tasks."},
	{WorkflowTasksComment, "Add workflow task comments."},
	{WorkflowTasksAttach, "Attach files to workflow tasks through tenant storage."},
	{WorkflowTasksWatch, "Watch or unwatch workflow tasks."},
	{WorkflowTasksRestrictedView, "View restricted and confidential workflow tasks."},
	{WorkflowTasksSummaryView, "View workflow task summary counts and SLA-style metrics."},
	{HRCasesList, "List HR helpdesk cases."},
	{HRCasesCreate, "Create HR helpdesk cases."},
	{HRCasesView, "View HR helpdesk case details and timeline."},
	{HRCasesUpdate, "Update HR helpdesk case details."},
	{HRCasesAssign, "Assign HR helpdesk cases to owners or owner roles."},
	{HRCasesStatus, "Change HR helpdesk case status and resolution."},
	{HRCasesComment, "Add HR helpdesk case comments."},
	{HRCasesAttach, "Attach files to HR helpdesk cases."},
	{HRCasesRestrictedView, "View sensitive, restricted, and grievance HR helpdesk cases."},
	{HRCaseCategoriesManage, "Manage HR helpdesk request categories and default routing."},
	{HRCaseSLAManage, "Manage HR helpdesk SLA response, resolution, and escalation policies."},
	{ERCasesList, "List confidential employee relations cases."},
	{ERCasesCreate, "Create grievance, disciplinary, and employee relations cases."},
	{ERCasesView, "View employee relations case workspace details."},
	{ERCasesUpdate, "Update employee relations case intake, owner, and privacy details."},
	{ERCasesStatus, "Move employee relations cases through triage, investigation, findings, action, and closure."},
	{ERCasesLegalHold, "Enable or release legal hold on employee relations cases and evidence."},
	{ERCaseCategoriesManage, "Manage employee relations grievance and disciplinary categories."},
	{ERCasePartiesManage, "Manage employee relations complainants, respondents, witnesses, and investigators."},
	{ERAllegationsManage, "Manage employee relations allegations and incident records."},
	{ERStepsManage, "Manage employee relations investigation steps."},
	{ERWitnessNotesManage, "Manage restricted employee relations witness notes."},
	{EREvidenceManage, "Manage employee relations evidence attachments through tenant storage."},
	{ERFindingsManage, "Manage employee relations investigation findings."},
	{ERActionPlansManage, "Manage employee relations corrective and follow-up action plans."},
	{EREventsView, "View immutable employee relations audit events."},
	{ERRestrictedView, "View restricted, sensitive, and legal-hold employee relations records."},
	{ReportsView, "View the HRMS report catalog and saved report views."},
	{ReportsManage, "Create and manage report saved views and reporting metadata."},
	{ReportsExport, "Create and audit report export jobs."},
	{ReportsSchedule, "Create and manage scheduled report delivery."},
	{InsightsView, "View HRMS deterministic and AI-assisted insight queues."},
	{InsightsRefresh, "Refresh deterministic HRMS insight rules."},
	{InsightsReview, "Review, resolve, dismiss, or override HRMS insights."},
	{AIActionsView, "View AI signal, action, override, and event-bus audit layers."},
	{AIActionsManage, "Create and update AI action proposals and workflow signals."},
	{AIActionsOverride, "Record human overrides for AI action recommendations."},
	{AISignalsEmit, "Emit workflow events into the AI signal and event outbox layer."},
	{AIAgentsView, "View bounded HR AI agent definitions and guardrails."},
	{AIAgentsRun, "Run bounded HR AI agents with deterministic fallback recommendations."},
	{PeopleAnalyticsView, "View modern aggregate people analytics and risk heatmaps."},
	{PrivacyConsentsManage, "Manage HRMS privacy consents and consent evidence."},
	{PrivacyErasureManage, "Manage data subject erasure, export, correction, and restriction workflows."},
	{IntegrationHooksManage, "Manage HRMS ecosystem hooks for WhatsApp, Slack, email, Git, webhooks, and APIs."},
	{MobileConstraintsManage, "Manage mobile-first API workflow constraints for older and newer devices."},
	{CelebrationTypesManage, "Manage celebration types."},
	{CelebrationsList, "List celebrations."},
	{CelebrationsCreate, "Create celebrations."},
	{CelebrationsView, "View celebration details."},
	{CelebrationsUpdate, "Update celebrations."},
	{CelebrationsDelete, "Deactivate celebrations."},
	{CelebrationsSend, "Send celebration notifications."},
	{NotificationsList, "List notifications."},
	{NotificationsSend, "Send notifications."},
	{NotificationsRead, "Read notification inbox."},
	{NotificationsPreferences, "Manage notification preferences."},
	{NotificationsMastersManage, "Manage notification masters."},
	{NotificationsDeviceTokensManage, "Manage notification device tokens."},
	{JobPositionsList, "List job positions."},
	{JobPositionsCreate, "Create job positions."},
	{JobPositionsView, "View job position details."},
	{JobPositionsUpdate, "Update job positions."},
	{JobPositionsDelete, "Deactivate job positions."},
	{JobRequisitionsList, "List job requisitions."},
	{JobRequisitionsCreate, "Create job requisitions."},
	{JobRequisitionsView, "View job requisition details."},
	{JobRequisitionsUpdate, "Update job requisitions."},
	{JobRequisitionsApprove, "Approve job requisitions."},
	{JobRequisitionsReject, "Reject job requisitions."},
	{JobRequisitionsClose, "Close job requisitions."},
	{JobPostingsList, "List job postings."},
	{JobPostingsCreate, "Create job postings."},
	{JobPostingsView, "View job posting details."},
	{JobPostingsUpdate, "Update job postings."},
	{JobPostingsPublish, "Publish job postings."},
	{JobPostingsClose, "Close job postings."},
	{CandidatesList, "List candidates."},
	{CandidatesCreate, "Create candidates."},
	{CandidatesView, "View candidate details."},
	{CandidatesUpdate, "Update candidates."},
	{CandidatesDelete, "Deactivate candidates."},
	{CandidateApplicationsList, "List candidate applications."},
	{CandidateApplicationsCreate, "Create candidate applications."},
	{CandidateApplicationsView, "View candidate application details."},
	{CandidateApplicationsUpdate, "Update candidate applications."},
	{CandidateApplicationsMove, "Move candidate applications through pipeline stages."},
	{ApplicantPortalView, "View own applicant profile and application status."},
	{InterviewRoundsList, "List interview rounds."},
	{InterviewRoundsCreate, "Create interview rounds."},
	{InterviewRoundsView, "View interview round details."},
	{InterviewRoundsUpdate, "Update interview rounds."},
	{InterviewRoundsDelete, "Delete interview rounds."},
	{OfferLettersList, "List offer letters."},
	{OfferLettersCreate, "Create offer letters."},
	{OfferLettersView, "View offer letter details."},
	{OfferLettersSend, "Send offer letters."},
	{OfferLettersUpdate, "Update offer letters."},
	{OfferLettersRevoke, "Revoke offer letters."},
	{OnboardingWorkflowsManage, "Manage onboarding workflows."},
	{OnboardingList, "List candidate onboarding records."},
	{OnboardingStart, "Start candidate onboarding."},
	{OnboardingView, "View candidate onboarding details."},
	{OnboardingUpdate, "Update candidate onboarding records."},
	{OnboardingCompleteTask, "Complete candidate onboarding tasks."},
	{IntegrationsEmailManage, "Manage email integration settings."},
	{IntegrationsSMSManage, "Manage SMS integration settings."},
	{IntegrationsWhatsAppManage, "Manage WhatsApp integration settings."},
	{IntegrationsStorageManage, "Manage storage integration settings."},
	{IntegrationsPushManage, "Manage push notification integration settings."},
	{ScheduledJobsList, "List HRMS scheduled jobs."},
	{ScheduledJobsRun, "Run HRMS scheduled jobs."},
}

var RoleTemplates = []RoleTemplate{
	{Code: "SUPER_ADMIN", Name: "HRMS Super Admin", Description: "Full HRMS access across all features.", Permissions: Keys()},
	{Code: "TENANT_ADMIN", Name: "Tenant Admin", Description: "Manage tenant HRMS setup, employees, leave, attendance, salary, notifications, and onboarding.", Permissions: tenantAdminPermissions()},
	{Code: "HR", Name: "HR", Description: "Manage employee lifecycle, leave, attendance, salary operations, celebrations, and onboarding.", Permissions: hrPermissions()},
	{Code: "MANAGER", Name: "Manager", Description: "Add team visibility, leave approvals, and attendance review permissions on top of the employee baseline.", Permissions: managerPermissions()},
	{Code: "EMPLOYEE", Name: "Employee", Description: "Baseline self-service role for every tenant user: own profile, attendance, leave, salary slips, notifications, policies, and dashboard.", Permissions: employeePermissions()},
	{Code: "APPLICANT", Name: "Applicant", Description: "External candidate access for own applicant profile and application status only.", Permissions: applicantPermissions()},
}

func Keys() []string {
	keys := make([]string, 0, len(Catalog))
	for _, permission := range Catalog {
		keys = append(keys, permission.Key)
	}
	return keys
}

func Manifest() identity.Manifest {
	perms := make([]identity.ManifestPermission, 0, len(Catalog))
	for _, permission := range Catalog {
		perms = append(perms, identity.ManifestPermission{Slug: permission.Key, Description: permission.Description})
	}

	roleTemplates := make([]identity.ManifestRoleTemplate, 0, len(RoleTemplates))
	for _, template := range RoleTemplates {
		roleTemplates = append(roleTemplates, identity.ManifestRoleTemplate{
			Code:        template.Code,
			Name:        template.Name,
			Description: template.Description,
			Permissions: append([]string(nil), template.Permissions...),
		})
	}

	return identity.Manifest{
		Name:          ModuleName,
		Code:          ModuleCode,
		Description:   ModuleDescription,
		Permissions:   perms,
		RoleTemplates: roleTemplates,
	}
}

func ManifestRolePermissions(code string) []string {
	for _, template := range RoleTemplates {
		if strings.EqualFold(template.Code, code) {
			return append([]string(nil), template.Permissions...)
		}
	}
	return nil
}

func tenantAdminPermissions() []string {
	return without(Keys(), ScheduledJobsRun)
}

func hrPermissions() []string {
	return []string{
		BrandingView,
		BranchesList, BranchesView,
		DepartmentsList, DepartmentsCreate, DepartmentsView, DepartmentsUpdate,
		DesignationsList, DesignationsCreate, DesignationsView, DesignationsUpdate, DesignationsMastersManage,
		WorkerTypesList, WorkerTypesCreate, WorkerTypesView, WorkerTypesUpdate, WorkerClassificationRulesList, WorkerClassificationRulesManage,
		WorkersList, WorkersCreate, WorkersView, WorkersUpdate,
		EngagementsList, EngagementsCreate, EngagementsView, EngagementsUpdate, EngagementsStatus,
		WorkLogsList, WorkLogsCreate, WorkLogsView, WorkLogsUpdate, WorkLogsSubmit, WorkLogsApprove, WorkLogsReject, WorkLogsReport,
		ProjectsList, ProjectsCreate, ProjectsView, ProjectsUpdate, ProjectsStatus,
		MilestonesList, MilestonesCreate, MilestonesView, MilestonesUpdate, MilestonesSubmit, MilestonesApprove, MilestonesReject, MilestonesEventsView,
		ComplianceRulesList, ComplianceRulesCreate, ComplianceRulesView, ComplianceRulesUpdate, ComplianceRulesDelete, ComplianceRulesSeed,
		ComplianceChecklistList, ComplianceChecklistGenerate, ComplianceChecklistReview, ComplianceChecklistEvidence, ComplianceChecklistWaive, ComplianceChecklistDelete, ComplianceEventsView,
		SkillCategoriesList, SkillCategoriesCreate, SkillCategoriesUpdate, SkillCategoriesDelete,
		SkillsList, SkillsCreate, SkillsView, SkillsUpdate, SkillsDelete,
		WorkerSkillsList, WorkerSkillsCreate, WorkerSkillsUpdate, WorkerSkillsVerify, WorkerSkillsDelete,
		WorkerSkillAssessmentsList, WorkerSkillAssessmentsCreate, SkillsSummaryView,
		ProjectSkillRequirementsList, ProjectSkillRequirementsCreate, ProjectSkillRequirementsView, ProjectSkillRequirementsUpdate, ProjectSkillRequirementsDelete,
		ProjectSkillGapsView, ProjectSkillDependenciesView,
		LearningCoursesList, LearningCoursesManage, LearningPathsList, LearningPathsManage,
		LearningEnrollmentsList, LearningEnrollmentsAssign, LearningEnrollmentsStatus, LearningCertificatesUpload,
		LearningRecommendationsView, LearningRecommendationsManage, LearningSummaryView,
		TalentMarketplaceOpportunitiesList, TalentMarketplaceOpportunitiesCreate, TalentMarketplaceOpportunitiesView, TalentMarketplaceOpportunitiesUpdate, TalentMarketplaceOpportunitiesDelete,
		TalentMarketplaceApplicationsList, TalentMarketplaceApplicationsCreate, TalentMarketplaceApplicationsView, TalentMarketplaceApplicationsUpdate,
		TalentMarketplaceRecommendationsView, TalentMarketplaceEventsView, TalentMarketplaceFallbackManage,
		OKRCyclesList, OKRCyclesCreate, OKRCyclesView, OKRCyclesUpdate, OKRCyclesDelete, OKRCyclesStatus,
		OKRObjectivesList, OKRObjectivesCreate, OKRObjectivesView, OKRObjectivesUpdate, OKRObjectivesDelete, OKRObjectivesStatus,
		OKRKeyResultsList, OKRKeyResultsCreate, OKRKeyResultsView, OKRKeyResultsUpdate, OKRKeyResultsDelete,
		OKRCheckInsList, OKRCheckInsCreate, OKRSummaryView,
		PerformanceCheckInsList, PerformanceCheckInsCreate, PerformanceCheckInsView, PerformanceCheckInsUpdate, PerformanceCheckInsSubmit, PerformanceCheckInsReview, PerformanceCheckInsDelete, PerformanceCheckInsSummary,
		PerformanceCalibrationView, PerformanceTimelineView,
		FeedbackRequestsList, FeedbackRequestsCreate, FeedbackRequestsView, FeedbackRequestsUpdate, FeedbackRequestsStatus,
		FeedbackResponsesList, FeedbackResponsesCreate, FeedbackResponsesView,
		PulseSurveysList, PulseSurveysCreate, PulseSurveysView, PulseSurveysUpdate, PulseSurveysDelete, PulseSurveysStatus,
		PulseQuestionsList, PulseQuestionsCreate, PulseQuestionsUpdate, PulseQuestionsDelete,
		PulseResponsesList, PulseResponsesCreate,
		WellbeingScoresList, WellbeingScoresUpsert, WellbeingAlertsList, WellbeingAlertsReview, WellbeingAggregateView,
		LookupsList,
		FinancialYearsList, FinancialYearsView,
		WorkingHoursList, WorkingHoursCreate, WorkingHoursUpdate, WorkingHoursCopy,
		HolidaysList, HolidaysCreate, HolidaysView, HolidaysUpdate,
		PoliciesList, PoliciesCreate, PoliciesView, PoliciesUpdate, PoliciesPublish,
		EmployeesList, EmployeesCreate, EmployeesView, EmployeesUpdate, EmployeesDeactivate, EmployeesDocumentsManage, EmployeesBankManage, EmployeesStatutoryManage, EmployeesCredentialsManage,
		EmployeeLettersList, EmployeeLettersCreate, EmployeeLettersView, EmployeeLettersApprove, EmployeeLettersSend, EmployeeLettersDownload, EmployeeLettersRevoke, EmployeeLetterTemplatesManage,
		AgreementsList, AgreementsCreate, AgreementsView, AgreementsSend, AgreementsSign, AgreementsRevoke, AgreementsDownload, AgreementsDelete, AgreementsEventsView, AgreementTemplatesManage,
		EmployeeExitsList, EmployeeExitsCreate, EmployeeExitsView, EmployeeExitsApprove, EmployeeExitsUpdate, EmployeeExitsComplete, EmployeeExitsCancel,
		LeaveTypesList, LeaveTypesCreate, LeaveTypesView, LeaveTypesUpdate,
		LeavePoliciesList, LeavePoliciesCreate, LeavePoliciesView, LeavePoliciesUpdate,
		LeaveTemplatesList, LeaveTemplatesCreate, LeaveTemplatesUpdate, LeaveTemplateRulesManage, LeaveAssignmentsManage,
		LeaveBalancesList, LeaveBalancesUpdate, LeaveLedgerView, LeaveAccrualRun,
		LeaveApprovalWorkflowsList, LeaveApprovalWorkflowsManage,
		LeavesList, LeavesView, LeavesApprove, LeavesReject, LeavesReport,
		AttendanceList, AttendanceView, AttendanceUpdate, AttendanceReviewRequest, AttendanceReport,
		AttendanceExceptionWorkflowsList, AttendanceExceptionWorkflowsManage, AttendancePayrollBlockersView,
		AttendanceLocationsList, AttendanceLocationsCreate, AttendanceLocationsUpdate,
		AttendanceLocationAssignmentsList, AttendanceLocationAssignmentsCreate, AttendanceLocationAssignmentsUpdate,
		AttendanceDevicesList, AttendanceDevicesCreate, AttendanceDevicesUpdate,
		EmployeeAttendanceDevicesList, EmployeeAttendanceDevicesCreate, EmployeeAttendanceDevicesUpdate,
		ShiftTemplatesList, ShiftTemplatesManage,
		StaffingRequirementsList, StaffingRequirementsManage,
		ShiftAssignmentsList, ShiftAssignmentsManage, ShiftAssignmentsPublish, ShiftAssignmentsLock,
		ShiftSwapsList, ShiftSwapsCreate, ShiftSwapsReview,
		ShiftScheduleEventsView, ShiftScheduleSummaryView,
		BenefitPlansList, BenefitPlansManage,
		BenefitWindowsList, BenefitWindowsManage,
		BenefitDependentsList, BenefitDependentsManage,
		BenefitEnrollmentsList, BenefitEnrollmentsManage, BenefitEnrollmentsReview,
		BenefitClaimTypesList, BenefitClaimTypesManage,
		BenefitClaimsList, BenefitClaimsCreate, BenefitClaimsReview, BenefitClaimsPay, BenefitClaimsAttach, BenefitClaimsExport,
		BenefitEventsView, BenefitSummaryView,
		PayCyclesView, PayCyclesUpdate, PayrollImportsList, PayrollImportsCreate, PayrollStatutoryRules, PayrollLocksManage, PayrollSalarySheetView, PayrollSalarySheetExport, PayrollReconciliation,
		PayGroupsList, PayGroupsCreate, PayGroupsView, PayGroupsUpdate, PayGroupsDelete,
		PayRunsList, PayRunsCreate, PayRunsView, PayRunsAssess, PayRunsFreeze, PayRunsGenerate, PayRunsLock,
		FlexPayRunsList, FlexPayRunsCreate, FlexPayRunsView, FlexPayRunsGenerate, FlexPayRunsSubmit, FlexPayRunsApprove, FlexPayRunsReject, FlexPayRunsPay, FlexPayRunsExport, FlexPayRunsDelete,
		ContractorInvoicesList, ContractorInvoicesCreate, ContractorInvoicesView, ContractorInvoicesUpdate, ContractorInvoicesSubmit, ContractorInvoicesApprove, ContractorInvoicesReject, ContractorInvoicesPay, ContractorInvoicesDelete,
		SalaryTemplatesList, SalaryTemplatesCreate, SalaryTemplatesView, SalaryTemplatesUpdate, SalaryTemplatesActivate,
		EmployeeSalariesList, EmployeeSalariesCreate, EmployeeSalariesView, EmployeeSalariesUpdate,
		SalarySlipsList, SalarySlipsGenerate, SalarySlipsView, SalarySlipsRegenerate, SalarySlipsDownload,
		CompensationPayBandsList, CompensationPayBandsManage,
		CompensationCyclesList, CompensationCyclesManage, CompensationBudgetPoolsManage,
		CompensationRecommendationsList, CompensationRecommendationsManage, CompensationRecommendationsApprove, CompensationPayrollHandoff,
		CompensationEquityView, CompensationEquityManage, CompensationEventsView, CompensationSummaryView,
		SuccessionCyclesList, SuccessionCyclesManage,
		SuccessionCriticalRolesList, SuccessionCriticalRolesManage,
		SuccessionSuccessorsList, SuccessionSuccessorsManage,
		SuccessionDevelopmentActionsList, SuccessionDevelopmentActionsManage,
		SuccessionEventsView, SuccessionSummaryView, SuccessionConfidentialView,
		AssetItemsList, AssetItemsManage,
		AccessCatalogList, AccessCatalogManage,
		AssetAssignmentsList, AssetAssignmentsManage,
		AccessTasksList, AccessTasksManage,
		AssetAccessEventsView, AssetAccessSummaryView,
		DashboardHRView, DashboardEmployeeView, HRCommandCenterView, OperationCatalogView, OperationsWorkbenchView, TenantOperationsView, TenantOperationsManage,
		WorkflowDefinitionsList, WorkflowDefinitionsManage, OperationTemplatesList, OperationTemplatesManage,
		WorkflowTasksList, WorkflowTasksCreate, WorkflowTasksView, WorkflowTasksUpdate, WorkflowTasksAct, WorkflowTasksComment, WorkflowTasksAttach, WorkflowTasksWatch, WorkflowTasksRestrictedView, WorkflowTasksSummaryView,
		HRCasesList, HRCasesCreate, HRCasesView, HRCasesUpdate, HRCasesAssign, HRCasesStatus, HRCasesComment, HRCasesAttach, HRCasesRestrictedView, HRCaseCategoriesManage, HRCaseSLAManage,
		ERCasesList, ERCasesCreate, ERCasesView, ERCasesUpdate, ERCasesStatus, ERCasesLegalHold, ERCaseCategoriesManage, ERCasePartiesManage, ERAllegationsManage, ERStepsManage, ERWitnessNotesManage, EREvidenceManage, ERFindingsManage, ERActionPlansManage, EREventsView, ERRestrictedView,
		ReportsView, ReportsManage, ReportsExport, ReportsSchedule,
		InsightsView, InsightsRefresh, InsightsReview,
		AIActionsView, AIActionsManage, AIActionsOverride, AISignalsEmit, AIAgentsView, AIAgentsRun, PeopleAnalyticsView,
		PrivacyConsentsManage, PrivacyErasureManage, IntegrationHooksManage, MobileConstraintsManage,
		CelebrationTypesManage, CelebrationsList, CelebrationsCreate, CelebrationsView, CelebrationsUpdate, CelebrationsSend,
		NotificationsList, NotificationsSend, NotificationsRead, NotificationsPreferences,
		JobPositionsList, JobPositionsCreate, JobPositionsView, JobPositionsUpdate,
		JobRequisitionsList, JobRequisitionsCreate, JobRequisitionsView, JobRequisitionsUpdate,
		JobPostingsList, JobPostingsCreate, JobPostingsView, JobPostingsUpdate, JobPostingsPublish,
		CandidatesList, CandidatesCreate, CandidatesView, CandidatesUpdate,
		CandidateApplicationsList, CandidateApplicationsCreate, CandidateApplicationsView, CandidateApplicationsUpdate, CandidateApplicationsMove,
		InterviewRoundsList, InterviewRoundsCreate, InterviewRoundsView, InterviewRoundsUpdate,
		OfferLettersList, OfferLettersCreate, OfferLettersView, OfferLettersSend, OfferLettersUpdate, OfferLettersRevoke,
		OnboardingWorkflowsManage, OnboardingList, OnboardingStart, OnboardingView, OnboardingUpdate, OnboardingCompleteTask,
	}
}

func managerPermissions() []string {
	return []string{
		EmployeesList, EmployeesView, WorkersList, WorkersView, EngagementsList, EngagementsView,
		WorkLogsList, WorkLogsView, WorkLogsApprove, WorkLogsReject, WorkLogsReport,
		ProjectsList, ProjectsView, MilestonesList, MilestonesView, MilestonesApprove, MilestonesReject, MilestonesEventsView,
		ComplianceRulesList, ComplianceRulesView, ComplianceChecklistList, ComplianceEventsView,
		SkillCategoriesList, SkillsList, SkillsView, WorkerSkillsList, WorkerSkillsVerify, WorkerSkillAssessmentsList, SkillsSummaryView,
		ProjectSkillRequirementsList, ProjectSkillRequirementsView, ProjectSkillGapsView, ProjectSkillDependenciesView,
		LearningCoursesList, LearningPathsList, LearningEnrollmentsList, LearningEnrollmentsAssign, LearningEnrollmentsStatus,
		LearningRecommendationsView, LearningRecommendationsManage, LearningSummaryView,
		TalentMarketplaceOpportunitiesList, TalentMarketplaceOpportunitiesCreate, TalentMarketplaceOpportunitiesView, TalentMarketplaceOpportunitiesUpdate,
		TalentMarketplaceApplicationsList, TalentMarketplaceApplicationsCreate, TalentMarketplaceApplicationsView, TalentMarketplaceApplicationsUpdate,
		TalentMarketplaceRecommendationsView, TalentMarketplaceEventsView, TalentMarketplaceFallbackManage,
		OKRCyclesList, OKRCyclesView,
		OKRObjectivesList, OKRObjectivesCreate, OKRObjectivesView, OKRObjectivesUpdate, OKRObjectivesStatus,
		OKRKeyResultsList, OKRKeyResultsCreate, OKRKeyResultsView, OKRKeyResultsUpdate,
		OKRCheckInsList, OKRCheckInsCreate, OKRSummaryView,
		PerformanceCheckInsList, PerformanceCheckInsCreate, PerformanceCheckInsView, PerformanceCheckInsUpdate, PerformanceCheckInsSubmit, PerformanceCheckInsReview, PerformanceCheckInsSummary,
		PerformanceCalibrationView, PerformanceTimelineView,
		FeedbackRequestsList, FeedbackRequestsCreate, FeedbackRequestsView, FeedbackRequestsUpdate, FeedbackRequestsStatus,
		FeedbackResponsesList, FeedbackResponsesCreate, FeedbackResponsesView,
		PulseSurveysList, PulseSurveysView, PulseQuestionsList, PulseResponsesCreate, WellbeingScoresList, WellbeingAggregateView,
		AgreementsList, AgreementsView, AgreementsDownload, AgreementsEventsView,
		LeavesList, LeavesView, LeavesApprove, LeavesReject, LeavesReport,
		AttendanceList, AttendanceView, AttendanceReviewRequest, AttendanceReport,
		ShiftAssignmentsList, ShiftSwapsList, ShiftSwapsReview, ShiftScheduleSummaryView,
		BenefitPlansList, BenefitWindowsList, BenefitEnrollmentsList, BenefitClaimsList, BenefitClaimsReview, BenefitSummaryView,
		HRCasesList, HRCasesCreate, HRCasesView, HRCasesStatus, HRCasesComment, HRCasesAttach,
		ERCasesList, ERCasesCreate, ERCasesView, ERCasesStatus, ERStepsManage, EREventsView,
		CompensationCyclesList, CompensationRecommendationsList, CompensationRecommendationsManage, CompensationRecommendationsApprove, CompensationEquityView, CompensationSummaryView,
		ReportsView, ReportsExport, InsightsView, AIActionsView, AIAgentsView, PeopleAnalyticsView, HRCommandCenterView, OperationCatalogView, OperationsWorkbenchView, TenantOperationsView,
		WorkflowDefinitionsList, OperationTemplatesList,
		WorkflowTasksList, WorkflowTasksCreate, WorkflowTasksView, WorkflowTasksUpdate, WorkflowTasksAct, WorkflowTasksComment, WorkflowTasksAttach, WorkflowTasksWatch, WorkflowTasksSummaryView,
	}
}

func employeePermissions() []string {
	return []string{
		TenantProfilesView,
		EmployeesView, EmployeesDocumentsManage, EmployeesBankManage, EmployeesStatutoryManage, EmployeeLettersView, EmployeeLettersDownload, EmployeeExitsView,
		AgreementsList, AgreementsView, AgreementsSign, AgreementsDownload,
		WorkersView, EngagementsView, WorkLogsList, WorkLogsCreate, WorkLogsView, WorkLogsUpdate, WorkLogsSubmit,
		ProjectsList, ProjectsView, MilestonesList, MilestonesView, MilestonesSubmit,
		SkillsList, SkillsView, WorkerSkillsList, WorkerSkillsCreate, WorkerSkillsUpdate, WorkerSkillAssessmentsList, WorkerSkillAssessmentsCreate,
		LearningCoursesList, LearningPathsList, LearningEnrollmentsList, LearningEnrollmentsStatus, LearningCertificatesUpload, LearningRecommendationsView,
		TalentMarketplaceOpportunitiesList, TalentMarketplaceOpportunitiesView,
		TalentMarketplaceApplicationsList, TalentMarketplaceApplicationsCreate, TalentMarketplaceApplicationsView, TalentMarketplaceApplicationsUpdate,
		OKRCyclesList, OKRCyclesView,
		OKRObjectivesList, OKRObjectivesView,
		OKRKeyResultsList, OKRKeyResultsView, OKRKeyResultsUpdate,
		OKRCheckInsList, OKRCheckInsCreate, OKRSummaryView,
		PerformanceCheckInsList, PerformanceCheckInsCreate, PerformanceCheckInsView, PerformanceCheckInsUpdate, PerformanceCheckInsSubmit,
		PerformanceTimelineView,
		FeedbackRequestsList, FeedbackRequestsView,
		FeedbackResponsesList, FeedbackResponsesCreate, FeedbackResponsesView,
		PulseSurveysList, PulseSurveysView, PulseQuestionsList, PulseResponsesCreate, WellbeingScoresList,
		LeaveTypesList, LeaveBalancesList,
		LeavesList, LeavesApply, LeavesView, LeavesCancel,
		AttendanceList, AttendanceCheckIn, AttendanceCheckOut, AttendanceView, AttendanceRegularize,
		ShiftAssignmentsList, ShiftSwapsList, ShiftSwapsCreate,
		BenefitPlansList, BenefitWindowsList, BenefitDependentsList, BenefitDependentsManage, BenefitEnrollmentsList, BenefitEnrollmentsManage, BenefitClaimTypesList, BenefitClaimsList, BenefitClaimsCreate, BenefitClaimsAttach, BenefitSummaryView,
		SalarySlipsList, SalarySlipsView, SalarySlipsDownload,
		DashboardEmployeeView,
		WorkflowTasksList, WorkflowTasksCreate, WorkflowTasksView, WorkflowTasksComment, WorkflowTasksAttach, WorkflowTasksWatch, OperationCatalogView,
		HRCasesCreate,
		ERCasesCreate,
		CelebrationsList, CelebrationsView,
		NotificationsRead, NotificationsPreferences, NotificationsDeviceTokensManage,
		PoliciesList, PoliciesView,
	}
}

func applicantPermissions() []string {
	return []string{ApplicantPortalView}
}

func without(values []string, skipped ...string) []string {
	skip := make(map[string]struct{}, len(skipped))
	for _, value := range skipped {
		skip[value] = struct{}{}
	}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := skip[value]; !ok {
			out = append(out, value)
		}
	}
	return out
}
