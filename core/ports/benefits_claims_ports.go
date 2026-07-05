package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type BenefitsClaimsRepo interface {
	CreateBenefitPlan(ctx context.Context, item *domain.BenefitPlan, actorID *uuid.UUID) (*domain.BenefitPlan, error)
	UpdateBenefitPlan(ctx context.Context, item *domain.BenefitPlan, actorID *uuid.UUID) (*domain.BenefitPlan, error)
	ListBenefitPlans(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitPlan, error)
	GetBenefitPlan(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitPlan, error)
	DeleteBenefitPlan(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateBenefitEnrollmentWindow(ctx context.Context, item *domain.BenefitEnrollmentWindow, actorID *uuid.UUID) (*domain.BenefitEnrollmentWindow, error)
	UpdateBenefitEnrollmentWindow(ctx context.Context, item *domain.BenefitEnrollmentWindow, actorID *uuid.UUID) (*domain.BenefitEnrollmentWindow, error)
	ListBenefitEnrollmentWindows(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitEnrollmentWindow, error)
	GetBenefitEnrollmentWindow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitEnrollmentWindow, error)
	DeleteBenefitEnrollmentWindow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateBenefitDependent(ctx context.Context, item *domain.BenefitDependent, actorID *uuid.UUID) (*domain.BenefitDependent, error)
	UpdateBenefitDependent(ctx context.Context, item *domain.BenefitDependent, actorID *uuid.UUID) (*domain.BenefitDependent, error)
	ListBenefitDependents(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitDependent, error)
	GetBenefitDependent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitDependent, error)
	DeleteBenefitDependent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateBenefitEnrollment(ctx context.Context, item *domain.BenefitEnrollment, actorID *uuid.UUID) (*domain.BenefitEnrollment, error)
	UpdateBenefitEnrollment(ctx context.Context, item *domain.BenefitEnrollment, actorID *uuid.UUID) (*domain.BenefitEnrollment, error)
	UpdateBenefitEnrollmentStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, reviewerID *uuid.UUID, remarks *string) (*domain.BenefitEnrollment, error)
	ListBenefitEnrollments(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitEnrollment, error)
	GetBenefitEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitEnrollment, error)
	DeleteBenefitEnrollment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateBenefitClaimType(ctx context.Context, item *domain.BenefitClaimType, actorID *uuid.UUID) (*domain.BenefitClaimType, error)
	UpdateBenefitClaimType(ctx context.Context, item *domain.BenefitClaimType, actorID *uuid.UUID) (*domain.BenefitClaimType, error)
	ListBenefitClaimTypes(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitClaimType, error)
	GetBenefitClaimType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitClaimType, error)
	DeleteBenefitClaimType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateBenefitClaim(ctx context.Context, item *domain.BenefitClaim, actorID *uuid.UUID) (*domain.BenefitClaim, error)
	UpdateBenefitClaim(ctx context.Context, item *domain.BenefitClaim, actorID *uuid.UUID) (*domain.BenefitClaim, error)
	UpdateBenefitClaimStatus(ctx context.Context, item *domain.BenefitClaim, actorID *uuid.UUID) (*domain.BenefitClaim, error)
	ListBenefitClaims(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitClaim, error)
	GetBenefitClaim(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.BenefitClaim, error)
	DeleteBenefitClaim(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	GetBenefitClaimLimitUsage(ctx context.Context, tenantID uuid.UUID, employeeUserID uuid.UUID, claimTypeID uuid.UUID, startDate string, endDate string) (float64, error)
	CreateBenefitClaimAttachment(ctx context.Context, item *domain.BenefitClaimAttachment, actorID *uuid.UUID) (*domain.BenefitClaimAttachment, error)
	ListBenefitClaimAttachments(ctx context.Context, tenantID uuid.UUID, claimID uuid.UUID) ([]*domain.BenefitClaimAttachment, error)
	CreateBenefitEvent(ctx context.Context, item *domain.BenefitEvent, actorID *uuid.UUID) (*domain.BenefitEvent, error)
	ListBenefitEvents(ctx context.Context, filter domain.BenefitFilter) ([]*domain.BenefitEvent, error)
	GetBenefitsSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.BenefitsSummaryRow, error)
}

type BenefitPlanCommand struct {
	ID                   uuid.UUID       `json:"id,omitempty"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	Code                 string          `json:"code"`
	Name                 string          `json:"name"`
	PlanType             string          `json:"plan_type"`
	Description          *string         `json:"description,omitempty"`
	ProviderName         *string         `json:"provider_name,omitempty"`
	PolicyNumber         *string         `json:"policy_number,omitempty"`
	CoverageAmount       *float64        `json:"coverage_amount,omitempty"`
	EmployerContribution float64         `json:"employer_contribution"`
	EmployeeContribution float64         `json:"employee_contribution"`
	CurrencyCode         string          `json:"currency_code"`
	EligibilityRule      json.RawMessage `json:"eligibility_rule,omitempty"`
	InsuranceMetadata    json.RawMessage `json:"insurance_metadata,omitempty"`
	EffectiveFrom        string          `json:"effective_from,omitempty"`
	EffectiveTo          string          `json:"effective_to,omitempty"`
	IsActive             bool            `json:"is_active"`
	ActorID              *uuid.UUID      `json:"-"`
}

type BenefitEnrollmentWindowCommand struct {
	ID            uuid.UUID       `json:"id,omitempty"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	PlanID        uuid.UUID       `json:"plan_id"`
	Name          string          `json:"name"`
	OpensOn       string          `json:"opens_on"`
	ClosesOn      string          `json:"closes_on"`
	EffectiveFrom string          `json:"effective_from,omitempty"`
	EffectiveTo   string          `json:"effective_to,omitempty"`
	Status        string          `json:"status"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	ActorID       *uuid.UUID      `json:"-"`
}

type BenefitDependentCommand struct {
	ID                uuid.UUID       `json:"id,omitempty"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	EmployeeUserID    uuid.UUID       `json:"employee_user_id"`
	FullName          string          `json:"full_name"`
	Relationship      string          `json:"relationship"`
	DateOfBirth       string          `json:"date_of_birth,omitempty"`
	Gender            *string         `json:"gender,omitempty"`
	NomineePercentage *float64        `json:"nominee_percentage,omitempty"`
	IsNominee         bool            `json:"is_nominee"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	ActorID           *uuid.UUID      `json:"-"`
}

type BenefitEnrollmentCommand struct {
	ID                   uuid.UUID       `json:"id,omitempty"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	PlanID               uuid.UUID       `json:"plan_id"`
	WindowID             *uuid.UUID      `json:"window_id,omitempty"`
	EmployeeUserID       uuid.UUID       `json:"employee_user_id"`
	Status               string          `json:"status"`
	CoverageLevel        *string         `json:"coverage_level,omitempty"`
	SelectedAmount       *float64        `json:"selected_amount,omitempty"`
	EmployeeContribution float64         `json:"employee_contribution"`
	EmployerContribution float64         `json:"employer_contribution"`
	EffectiveFrom        string          `json:"effective_from,omitempty"`
	EffectiveTo          string          `json:"effective_to,omitempty"`
	SubmittedAt          string          `json:"submitted_at,omitempty"`
	ReviewRemarks        *string         `json:"review_remarks,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	ActorID              *uuid.UUID      `json:"-"`
}

type BenefitEnrollmentStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id,omitempty"`
	Status   string     `json:"status"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type BenefitClaimTypeCommand struct {
	ID                   uuid.UUID       `json:"id,omitempty"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	PlanID               *uuid.UUID      `json:"plan_id,omitempty"`
	Code                 string          `json:"code"`
	Name                 string          `json:"name"`
	Description          *string         `json:"description,omitempty"`
	AnnualLimit          *float64        `json:"annual_limit,omitempty"`
	PerClaimLimit        *float64        `json:"per_claim_limit,omitempty"`
	RequiresAttachment   bool            `json:"requires_attachment"`
	Taxable              bool            `json:"taxable"`
	PayrollComponentCode *string         `json:"payroll_component_code,omitempty"`
	EligibilityRule      json.RawMessage `json:"eligibility_rule,omitempty"`
	IsActive             bool            `json:"is_active"`
	ActorID              *uuid.UUID      `json:"-"`
}

type BenefitClaimCommand struct {
	ID                     uuid.UUID       `json:"id,omitempty"`
	TenantID               uuid.UUID       `json:"tenant_id"`
	ClaimNumber            string          `json:"claim_number,omitempty"`
	ClaimTypeID            uuid.UUID       `json:"claim_type_id"`
	PlanID                 *uuid.UUID      `json:"plan_id,omitempty"`
	EmployeeUserID         uuid.UUID       `json:"employee_user_id"`
	DependentID            *uuid.UUID      `json:"dependent_id,omitempty"`
	ExpenseDate            string          `json:"expense_date"`
	SubmittedAt            string          `json:"submitted_at,omitempty"`
	ClaimAmount            float64         `json:"claim_amount"`
	ApprovedAmount         *float64        `json:"approved_amount,omitempty"`
	CurrencyCode           string          `json:"currency_code"`
	Status                 string          `json:"status"`
	PaymentStatus          string          `json:"payment_status"`
	PaymentReference       *string         `json:"payment_reference,omitempty"`
	PaidAt                 string          `json:"paid_at,omitempty"`
	ReviewRemarks          *string         `json:"review_remarks,omitempty"`
	PayrollExportStatus    string          `json:"payroll_export_status"`
	PayrollExportedAt      string          `json:"payroll_exported_at,omitempty"`
	PayrollExportReference *string         `json:"payroll_export_reference,omitempty"`
	Notes                  *string         `json:"notes,omitempty"`
	Metadata               json.RawMessage `json:"metadata,omitempty"`
	ActorID                *uuid.UUID      `json:"-"`
}

type BenefitClaimStatusCommand struct {
	TenantID               uuid.UUID  `json:"tenant_id"`
	ID                     uuid.UUID  `json:"id,omitempty"`
	Status                 string     `json:"status"`
	ApprovedAmount         *float64   `json:"approved_amount,omitempty"`
	PaymentReference       *string    `json:"payment_reference,omitempty"`
	PayrollExportReference *string    `json:"payroll_export_reference,omitempty"`
	Remarks                *string    `json:"remarks,omitempty"`
	ActorID                *uuid.UUID `json:"-"`
}

type BenefitClaimAttachmentCommand struct {
	TenantID      uuid.UUID       `json:"tenant_id"`
	ClaimID       uuid.UUID       `json:"claim_id"`
	FileName      string          `json:"file_name"`
	ContentType   string          `json:"content_type"`
	ContentBase64 string          `json:"content_base64"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	ActorID       *uuid.UUID      `json:"-"`
}
