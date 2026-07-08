package services

import (
	"context"

	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
	"github.com/ranakdinesh/spur-hrms/internal/logging"
	"github.com/rs/zerolog"
)

type TenantService struct {
	profiles                     ports.TenantProfileRepo
	settings                     ports.TenantSettingsRepo
	branding                     ports.TenantBrandingRepo
	branches                     ports.BranchRepo
	departments                  ports.DepartmentRepo
	designations                 ports.DesignationRepo
	designationMasters           ports.DesignationMasterRepo
	workerTypes                  ports.WorkerTypeRepo
	workerProfiles               ports.WorkerProfileRepo
	engagements                  ports.EngagementRepo
	workLogs                     ports.WorkLogRepo
	projects                     ports.ProjectRepo
	compliance                   ports.ComplianceRepo
	skills                       ports.SkillsRepo
	skillGaps                    ports.SkillGapRepo
	learning                     ports.LearningRepo
	compensationReview           ports.CompensationReviewRepo
	successionPlanning           ports.SuccessionPlanningRepo
	assetAccess                  ports.AssetAccessLifecycleRepo
	talentMarketplace            ports.TalentMarketplaceRepo
	okrs                         ports.OKRRepo
	performance                  ports.PerformanceRepo
	wellbeing                    ports.WellbeingRepo
	agreements                   ports.AgreementRepo
	workingHours                 ports.WorkingHourRepo
	financialYears               ports.FinancialYearRepo
	holidays                     ports.HolidayRepo
	policyEngine                 ports.PolicyEngineRepo
	leaveTypes                   ports.LeaveTypeRepo
	leavePolicies                ports.LeavePolicyRepo
	leaveTemplates               ports.LeaveTemplateRepo
	leaveBalances                ports.LeaveBalanceRepo
	leaveRequests                ports.LeaveRequestRepo
	compOffRequests              ports.CompOffRequestRepo
	approvalWorkflows            ports.LeaveApprovalWorkflowRepo
	lookups                      ports.EmploymentLookupRepo
	policies                     ports.PolicyRepo
	subscriptionPlans            ports.SubscriptionPlanRepo
	subscriptions                ports.TenantSubscriptionRepo
	employees                    ports.EmployeeRepo
	employeeCredentialEvents     ports.EmployeeCredentialEventRepo
	employeeExits                ports.EmployeeExitRepo
	employeeDocuments            ports.EmployeeDocumentRepo
	attendances                  ports.AttendanceRepo
	attendancePolicies           ports.AttendancePolicyRepo
	attendanceRosters            ports.AttendanceRosterRepo
	attendanceRequests           ports.AttendanceRequestRepo
	overtimeRequests             ports.OvertimeRequestRepo
	attendanceExceptionWorkflows ports.AttendanceExceptionWorkflowRepo
	attendanceLocations          ports.AttendanceLocationRepo
	attendanceDevices            ports.AttendanceDeviceRepo
	shiftScheduling              ports.ShiftSchedulingRepo
	payCycles                    ports.PayCycleRepo
	salaryTemplates              ports.SalaryTemplateRepo
	employeeSalaries             ports.EmployeeSalaryRepo
	salarySlips                  ports.SalarySlipRepo
	payrollOperations            ports.PayrollOperationsRepo
	payGroups                    ports.PayGroupRepo
	flexPayroll                  ports.FlexPayrollRepo
	reporting                    ports.ReportingRepo
	insights                     ports.InsightRepo
	aiActions                    ports.AIActionLayerRepo
	peopleAnalytics              ports.PeopleAnalyticsRepo
	privacyEcosystem             ports.PrivacyEcosystemRepo
	operationsWorkbench          ports.OperationsWorkbenchRepo
	workflowTasks                ports.WorkflowTaskEngineRepo
	tenantOperations             ports.TenantOperationGovernanceRepo
	benefitsClaims               ports.BenefitsClaimsRepo
	employeeRelations            ports.EmployeeRelationsRepo
	hrCases                      ports.HRCaseRepo
	celebrations                 ports.CelebrationRepo
	scheduledJobs                ports.ScheduledJobRepo
	notifications                ports.NotificationRepo
	emailProviders               ports.EmailProviderRepo
	communicationProviders       ports.CommunicationProviderRepo
	storageProviders             ports.StorageProviderRepo
	pushProviders                ports.PushProviderRepo
	jobPositions                 ports.JobPositionRepo
	jobRequisitions              ports.JobRequisitionRepo
	jobPostings                  ports.JobPostingRepo
	candidates                   ports.CandidateRepo
	offerLetters                 ports.OfferLetterRepo
	employeeLetters              ports.EmployeeLetterRepo
	onboardingWorkflows          ports.OnboardingWorkflowRepo
	candidateOnboardings         ports.CandidateOnboardingRepo
	salarySlipPDF                ports.SalarySlipPDFRenderer
	salarySlipStorage            ports.SalarySlipStorage
	employeeLetterPDF            ports.EmployeeLetterPDFRenderer
	employeeLetterStorage        ports.EmployeeLetterStorage
	agreementPDF                 ports.AgreementPDFRenderer
	agreementStorage             ports.AgreementStorage
	employeeIdentity             ports.EmployeeIdentityPort
	insightScorer                ports.InsightScoringPort
	aiEventPublisher             ports.AIEventPublisherPort
	aiActionSidecar              ports.AIActionSidecarPort
	legacyPasswordMigration      ports.LegacyPasswordMigrationPort
	policyStorage                ports.PolicyFileStorage
	documentStorage              ports.EmployeeDocumentStorage
	hrCaseAttachmentStorage      ports.HRCaseAttachmentStorage
	learningCertificateStorage   ports.LearningCertificateStorage
	objectStorage                ports.ObjectStorage
	defaultStorageProvider       *domain.StorageProviderSettings
	leaveNotifier                ports.LeaveNotifier
	policyNotifier               ports.PolicyNotifier
	celebrationNotifier          ports.CelebrationNotifier
	emailDelivery                ports.EmailDeliverySender
	defaultEmailProvider         *domain.EmailProviderSettings
	globalEmailProviderOnly      bool
	communicationDelivery        ports.CommunicationDeliverySender
	defaultCommunicationProvider *domain.CommunicationProviderSettings
	pushDelivery                 ports.PushDeliverySender
	defaultPushProvider          *domain.PushProviderSettings
	registrationEmail            ports.TenantRegistrationEmailSender
	system                       ports.SystemRunner
	log                          *zerolog.Logger
}

func NewTenantService(repo interface {
	ports.TenantProfileRepo
	ports.TenantSettingsRepo
	ports.TenantBrandingRepo
	ports.BranchRepo
	ports.DepartmentRepo
	ports.DesignationRepo
	ports.DesignationMasterRepo
	ports.WorkerTypeRepo
	ports.WorkerProfileRepo
	ports.EngagementRepo
	ports.WorkLogRepo
	ports.ProjectRepo
	ports.ComplianceRepo
	ports.SkillsRepo
	ports.SkillGapRepo
	ports.LearningRepo
	ports.CompensationReviewRepo
	ports.SuccessionPlanningRepo
	ports.AssetAccessLifecycleRepo
	ports.TalentMarketplaceRepo
	ports.OKRRepo
	ports.PerformanceRepo
	ports.WellbeingRepo
	ports.AgreementRepo
	ports.WorkingHourRepo
	ports.FinancialYearRepo
	ports.HolidayRepo
	ports.PolicyEngineRepo
	ports.LeaveTypeRepo
	ports.LeavePolicyRepo
	ports.LeaveTemplateRepo
	ports.LeaveBalanceRepo
	ports.LeaveRequestRepo
	ports.CompOffRequestRepo
	ports.LeaveApprovalWorkflowRepo
	ports.EmploymentLookupRepo
	ports.PolicyRepo
	ports.SubscriptionPlanRepo
	ports.TenantSubscriptionRepo
	ports.EmployeeRepo
	ports.EmployeeCredentialEventRepo
	ports.EmployeeExitRepo
	ports.EmployeeDocumentRepo
	ports.AttendanceRepo
	ports.AttendancePolicyRepo
	ports.AttendanceRosterRepo
	ports.AttendanceRequestRepo
	ports.OvertimeRequestRepo
	ports.AttendanceExceptionWorkflowRepo
	ports.AttendanceLocationRepo
	ports.AttendanceDeviceRepo
	ports.ShiftSchedulingRepo
	ports.PayCycleRepo
	ports.SalaryTemplateRepo
	ports.EmployeeSalaryRepo
	ports.SalarySlipRepo
	ports.PayrollOperationsRepo
	ports.PayGroupRepo
	ports.FlexPayrollRepo
	ports.ReportingRepo
	ports.InsightRepo
	ports.AIActionLayerRepo
	ports.PeopleAnalyticsRepo
	ports.PrivacyEcosystemRepo
	ports.OperationsWorkbenchRepo
	ports.WorkflowTaskEngineRepo
	ports.TenantOperationGovernanceRepo
	ports.BenefitsClaimsRepo
	ports.EmployeeRelationsRepo
	ports.HRCaseRepo
	ports.CelebrationRepo
	ports.ScheduledJobRepo
	ports.NotificationRepo
	ports.EmailProviderRepo
	ports.CommunicationProviderRepo
	ports.StorageProviderRepo
	ports.PushProviderRepo
	ports.JobPositionRepo
	ports.JobRequisitionRepo
	ports.JobPostingRepo
	ports.CandidateRepo
	ports.OfferLetterRepo
	ports.EmployeeLetterRepo
	ports.OnboardingWorkflowRepo
	ports.CandidateOnboardingRepo
	ports.SystemRunner
}, log *zerolog.Logger, opts ...TenantServiceOption) *TenantService {
	svc := &TenantService{profiles: repo, settings: repo, branding: repo, branches: repo, departments: repo, designations: repo, designationMasters: repo, workerTypes: repo, workerProfiles: repo, engagements: repo, workLogs: repo, projects: repo, compliance: repo, skills: repo, skillGaps: repo, learning: repo, compensationReview: repo, successionPlanning: repo, assetAccess: repo, talentMarketplace: repo, okrs: repo, performance: repo, wellbeing: repo, agreements: repo, workingHours: repo, financialYears: repo, holidays: repo, policyEngine: repo, leaveTypes: repo, leavePolicies: repo, leaveTemplates: repo, leaveBalances: repo, leaveRequests: repo, compOffRequests: repo, approvalWorkflows: repo, lookups: repo, policies: repo, subscriptionPlans: repo, subscriptions: repo, employees: repo, employeeCredentialEvents: repo, employeeExits: repo, employeeDocuments: repo, attendances: repo, attendancePolicies: repo, attendanceRosters: repo, attendanceRequests: repo, overtimeRequests: repo, attendanceExceptionWorkflows: repo, attendanceLocations: repo, attendanceDevices: repo, shiftScheduling: repo, payCycles: repo, salaryTemplates: repo, employeeSalaries: repo, salarySlips: repo, payrollOperations: repo, payGroups: repo, flexPayroll: repo, reporting: repo, insights: repo, aiActions: repo, peopleAnalytics: repo, privacyEcosystem: repo, operationsWorkbench: repo, workflowTasks: repo, tenantOperations: repo, benefitsClaims: repo, employeeRelations: repo, hrCases: repo, celebrations: repo, scheduledJobs: repo, notifications: repo, emailProviders: repo, communicationProviders: repo, storageProviders: repo, pushProviders: repo, jobPositions: repo, jobRequisitions: repo, jobPostings: repo, candidates: repo, offerLetters: repo, employeeLetters: repo, onboardingWorkflows: repo, candidateOnboardings: repo, system: repo, log: logging.Component(log, "service")}
	for _, opt := range opts {
		if opt != nil {
			opt(svc)
		}
	}
	return svc
}

type TenantServiceOption func(*TenantService)

func WithPolicyFileStorage(storage ports.PolicyFileStorage) TenantServiceOption {
	return func(s *TenantService) {
		s.policyStorage = storage
	}
}

func WithEmployeeDocumentStorage(storage ports.EmployeeDocumentStorage) TenantServiceOption {
	return func(s *TenantService) {
		s.documentStorage = storage
	}
}

func WithHRCaseAttachmentStorage(storage ports.HRCaseAttachmentStorage) TenantServiceOption {
	return func(s *TenantService) {
		s.hrCaseAttachmentStorage = storage
	}
}

func WithLearningCertificateStorage(storage ports.LearningCertificateStorage) TenantServiceOption {
	return func(s *TenantService) {
		s.learningCertificateStorage = storage
	}
}

func WithObjectStorage(storage ports.ObjectStorage) TenantServiceOption {
	return func(s *TenantService) {
		s.objectStorage = storage
	}
}

func WithDefaultStorageProvider(settings *domain.StorageProviderSettings) TenantServiceOption {
	return func(s *TenantService) {
		s.defaultStorageProvider = settings
	}
}

func WithLeaveNotifier(notifier ports.LeaveNotifier) TenantServiceOption {
	return func(s *TenantService) {
		s.leaveNotifier = notifier
	}
}

func WithPolicyNotifier(notifier ports.PolicyNotifier) TenantServiceOption {
	return func(s *TenantService) {
		s.policyNotifier = notifier
	}
}

func WithCelebrationNotifier(notifier ports.CelebrationNotifier) TenantServiceOption {
	return func(s *TenantService) {
		s.celebrationNotifier = notifier
	}
}

func WithEmailDeliverySender(sender ports.EmailDeliverySender) TenantServiceOption {
	return func(s *TenantService) {
		s.emailDelivery = sender
	}
}

func WithDefaultEmailProvider(settings *domain.EmailProviderSettings) TenantServiceOption {
	return func(s *TenantService) {
		s.defaultEmailProvider = settings
	}
}

func WithGlobalEmailProviderOnly(enabled bool) TenantServiceOption {
	return func(s *TenantService) {
		s.globalEmailProviderOnly = enabled
	}
}

func WithCommunicationDeliverySender(sender ports.CommunicationDeliverySender) TenantServiceOption {
	return func(s *TenantService) {
		s.communicationDelivery = sender
	}
}

func WithDefaultCommunicationProvider(settings *domain.CommunicationProviderSettings) TenantServiceOption {
	return func(s *TenantService) {
		s.defaultCommunicationProvider = settings
	}
}

func WithPushDeliverySender(sender ports.PushDeliverySender) TenantServiceOption {
	return func(s *TenantService) {
		s.pushDelivery = sender
	}
}

func WithDefaultPushProvider(settings *domain.PushProviderSettings) TenantServiceOption {
	return func(s *TenantService) {
		s.defaultPushProvider = settings
	}
}

func WithTenantRegistrationEmailSender(sender ports.TenantRegistrationEmailSender) TenantServiceOption {
	return func(s *TenantService) {
		s.registrationEmail = sender
	}
}

func WithEmployeeIdentityPort(identity ports.EmployeeIdentityPort) TenantServiceOption {
	return func(s *TenantService) {
		s.employeeIdentity = identity
	}
}

func WithLegacyPasswordMigrationPort(migration ports.LegacyPasswordMigrationPort) TenantServiceOption {
	return func(s *TenantService) {
		s.legacyPasswordMigration = migration
	}
}

func WithAIEventPublisher(publisher ports.AIEventPublisherPort) TenantServiceOption {
	return func(s *TenantService) {
		s.aiEventPublisher = publisher
	}
}

func WithAIActionSidecar(sidecar ports.AIActionSidecarPort) TenantServiceOption {
	return func(s *TenantService) {
		s.aiActionSidecar = sidecar
	}
}

func WithSalarySlipPDFRenderer(renderer ports.SalarySlipPDFRenderer) TenantServiceOption {
	return func(s *TenantService) {
		s.salarySlipPDF = renderer
	}
}

func WithSalarySlipStorage(storage ports.SalarySlipStorage) TenantServiceOption {
	return func(s *TenantService) {
		s.salarySlipStorage = storage
	}
}

func WithEmployeeLetterPDFRenderer(renderer ports.EmployeeLetterPDFRenderer) TenantServiceOption {
	return func(s *TenantService) {
		s.employeeLetterPDF = renderer
	}
}

func WithEmployeeLetterStorage(storage ports.EmployeeLetterStorage) TenantServiceOption {
	return func(s *TenantService) {
		s.employeeLetterStorage = storage
	}
}

func WithAgreementPDFRenderer(renderer ports.AgreementPDFRenderer) TenantServiceOption {
	return func(s *TenantService) {
		s.agreementPDF = renderer
	}
}

func WithAgreementStorage(storage ports.AgreementStorage) TenantServiceOption {
	return func(s *TenantService) {
		s.agreementStorage = storage
	}
}

func (s *TenantService) RunAsSystem(ctx context.Context, fn func(context.Context) error) error {
	if err := s.system.RunAsSystem(ctx, fn); err != nil {
		s.logError("run as system", err)
		return err
	}
	return nil
}
