package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	BenefitPlanTypeInsurance     = "insurance"
	BenefitPlanTypeReimbursement = "reimbursement"
	BenefitPlanTypeAllowance     = "allowance"
	BenefitPlanTypeRetirement    = "retirement"
	BenefitPlanTypeWellness      = "wellness"
	BenefitPlanTypeOther         = "other"

	BenefitWindowStatusDraft    = "draft"
	BenefitWindowStatusOpen     = "open"
	BenefitWindowStatusClosed   = "closed"
	BenefitWindowStatusArchived = "archived"

	BenefitEnrollmentStatusDraft     = "draft"
	BenefitEnrollmentStatusSubmitted = "submitted"
	BenefitEnrollmentStatusApproved  = "approved"
	BenefitEnrollmentStatusRejected  = "rejected"
	BenefitEnrollmentStatusCancelled = "cancelled"
	BenefitEnrollmentStatusActive    = "active"
	BenefitEnrollmentStatusEnded     = "ended"

	BenefitClaimStatusDraft       = "draft"
	BenefitClaimStatusSubmitted   = "submitted"
	BenefitClaimStatusUnderReview = "under_review"
	BenefitClaimStatusApproved    = "approved"
	BenefitClaimStatusRejected    = "rejected"
	BenefitClaimStatusCancelled   = "cancelled"
	BenefitClaimStatusPaid        = "paid"

	BenefitPaymentStatusNotPayable = "not_payable"
	BenefitPaymentStatusPending    = "pending"
	BenefitPaymentStatusPaid       = "paid"
	BenefitPaymentStatusFailed     = "failed"

	BenefitPayrollExportStatusNotReady = "not_ready"
	BenefitPayrollExportStatusReady    = "ready"
	BenefitPayrollExportStatusExported = "exported"
	BenefitPayrollExportStatusBlocked  = "blocked"
)

var (
	ErrInvalidBenefitPlan              = errors.New("benefit plan is invalid")
	ErrBenefitPlanNotFound             = errors.New("benefit plan not found")
	ErrInvalidBenefitEnrollmentWindow  = errors.New("benefit enrollment window is invalid")
	ErrBenefitEnrollmentWindowNotFound = errors.New("benefit enrollment window not found")
	ErrInvalidBenefitDependent         = errors.New("benefit dependent is invalid")
	ErrBenefitDependentNotFound        = errors.New("benefit dependent not found")
	ErrInvalidBenefitEnrollment        = errors.New("benefit enrollment is invalid")
	ErrBenefitEnrollmentNotFound       = errors.New("benefit enrollment not found")
	ErrInvalidBenefitClaimType         = errors.New("benefit claim type is invalid")
	ErrBenefitClaimTypeNotFound        = errors.New("benefit claim type not found")
	ErrInvalidBenefitClaim             = errors.New("benefit claim is invalid")
	ErrBenefitClaimNotFound            = errors.New("benefit claim not found")
	ErrInvalidBenefitClaimAttachment   = errors.New("benefit claim attachment is invalid")
)

type BenefitPlan struct {
	ID                   uuid.UUID       `json:"id"`
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
	EligibilityRule      json.RawMessage `json:"eligibility_rule"`
	InsuranceMetadata    json.RawMessage `json:"insurance_metadata"`
	EffectiveFrom        *time.Time      `json:"effective_from,omitempty"`
	EffectiveTo          *time.Time      `json:"effective_to,omitempty"`
	IsActive             bool            `json:"is_active"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type BenefitEnrollmentWindow struct {
	ID            uuid.UUID       `json:"id"`
	TenantID      uuid.UUID       `json:"tenant_id"`
	PlanID        uuid.UUID       `json:"plan_id"`
	Name          string          `json:"name"`
	OpensOn       time.Time       `json:"opens_on"`
	ClosesOn      time.Time       `json:"closes_on"`
	EffectiveFrom *time.Time      `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time      `json:"effective_to,omitempty"`
	Status        string          `json:"status"`
	Metadata      json.RawMessage `json:"metadata"`
	Inactive      bool            `json:"inactive"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     *uuid.UUID      `json:"updated_by,omitempty"`
}

type BenefitDependent struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	EmployeeUserID    uuid.UUID       `json:"employee_user_id"`
	FullName          string          `json:"full_name"`
	Relationship      string          `json:"relationship"`
	DateOfBirth       *time.Time      `json:"date_of_birth,omitempty"`
	Gender            *string         `json:"gender,omitempty"`
	NomineePercentage *float64        `json:"nominee_percentage,omitempty"`
	IsNominee         bool            `json:"is_nominee"`
	Metadata          json.RawMessage `json:"metadata"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedBy         *uuid.UUID      `json:"updated_by,omitempty"`
}

type BenefitEnrollment struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	PlanID               uuid.UUID       `json:"plan_id"`
	WindowID             *uuid.UUID      `json:"window_id,omitempty"`
	EmployeeUserID       uuid.UUID       `json:"employee_user_id"`
	Status               string          `json:"status"`
	CoverageLevel        *string         `json:"coverage_level,omitempty"`
	SelectedAmount       *float64        `json:"selected_amount,omitempty"`
	EmployeeContribution float64         `json:"employee_contribution"`
	EmployerContribution float64         `json:"employer_contribution"`
	EffectiveFrom        *time.Time      `json:"effective_from,omitempty"`
	EffectiveTo          *time.Time      `json:"effective_to,omitempty"`
	SubmittedAt          *time.Time      `json:"submitted_at,omitempty"`
	ReviewedBy           *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewedAt           *time.Time      `json:"reviewed_at,omitempty"`
	ReviewRemarks        *string         `json:"review_remarks,omitempty"`
	Metadata             json.RawMessage `json:"metadata"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type BenefitClaimType struct {
	ID                   uuid.UUID       `json:"id"`
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
	EligibilityRule      json.RawMessage `json:"eligibility_rule"`
	IsActive             bool            `json:"is_active"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type BenefitClaim struct {
	ID                     uuid.UUID       `json:"id"`
	TenantID               uuid.UUID       `json:"tenant_id"`
	ClaimNumber            string          `json:"claim_number"`
	ClaimTypeID            uuid.UUID       `json:"claim_type_id"`
	PlanID                 *uuid.UUID      `json:"plan_id,omitempty"`
	EmployeeUserID         uuid.UUID       `json:"employee_user_id"`
	DependentID            *uuid.UUID      `json:"dependent_id,omitempty"`
	ExpenseDate            time.Time       `json:"expense_date"`
	SubmittedAt            *time.Time      `json:"submitted_at,omitempty"`
	ClaimAmount            float64         `json:"claim_amount"`
	ApprovedAmount         *float64        `json:"approved_amount,omitempty"`
	CurrencyCode           string          `json:"currency_code"`
	Status                 string          `json:"status"`
	PaymentStatus          string          `json:"payment_status"`
	PaymentReference       *string         `json:"payment_reference,omitempty"`
	PaidAt                 *time.Time      `json:"paid_at,omitempty"`
	ReviewedBy             *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewedAt             *time.Time      `json:"reviewed_at,omitempty"`
	ReviewRemarks          *string         `json:"review_remarks,omitempty"`
	PayrollExportStatus    string          `json:"payroll_export_status"`
	PayrollExportedAt      *time.Time      `json:"payroll_exported_at,omitempty"`
	PayrollExportReference *string         `json:"payroll_export_reference,omitempty"`
	Notes                  *string         `json:"notes,omitempty"`
	Metadata               json.RawMessage `json:"metadata"`
	Inactive               bool            `json:"inactive"`
	CreatedAt              time.Time       `json:"created_at"`
	CreatedBy              *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt              time.Time       `json:"updated_at"`
	UpdatedBy              *uuid.UUID      `json:"updated_by,omitempty"`
}

type BenefitClaimAttachment struct {
	ID             uuid.UUID       `json:"id"`
	TenantID       uuid.UUID       `json:"tenant_id"`
	ClaimID        uuid.UUID       `json:"claim_id"`
	FileName       string          `json:"file_name"`
	ContentType    string          `json:"content_type"`
	StoragePath    string          `json:"storage_path"`
	ChecksumSHA256 *string         `json:"checksum_sha256,omitempty"`
	SizeBytes      int64           `json:"size_bytes"`
	UploadedBy     *uuid.UUID      `json:"uploaded_by,omitempty"`
	Metadata       json.RawMessage `json:"metadata"`
	Inactive       bool            `json:"inactive"`
	CreatedAt      time.Time       `json:"created_at"`
	CreatedBy      *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt      time.Time       `json:"updated_at"`
	UpdatedBy      *uuid.UUID      `json:"updated_by,omitempty"`
}

type BenefitEvent struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	SourceType  string          `json:"source_type"`
	SourceID    uuid.UUID       `json:"source_id"`
	Action      string          `json:"action"`
	FromStatus  *string         `json:"from_status,omitempty"`
	ToStatus    *string         `json:"to_status,omitempty"`
	ActorUserID *uuid.UUID      `json:"actor_user_id,omitempty"`
	Remarks     *string         `json:"remarks,omitempty"`
	Metadata    json.RawMessage `json:"metadata"`
	Inactive    bool            `json:"inactive"`
	CreatedAt   time.Time       `json:"created_at"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at"`
	UpdatedBy   *uuid.UUID      `json:"updated_by,omitempty"`
}

type BenefitFilter struct {
	TenantID            uuid.UUID
	ActiveOnly          *bool
	Search              *string
	PlanType            *string
	PlanID              *uuid.UUID
	Status              *string
	EmployeeUserID      *uuid.UUID
	ClaimTypeID         *uuid.UUID
	PaymentStatus       *string
	PayrollExportStatus *string
	NomineesOnly        *bool
	SourceType          *string
	SourceID            *uuid.UUID
	Limit               int32
	Offset              int32
}

type BenefitsSummaryRow struct {
	Metric      string  `json:"metric"`
	MetricCount int32   `json:"metric_count"`
	Amount      float64 `json:"amount"`
}

func NewBenefitPlan(item BenefitPlan) (*BenefitPlan, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidBenefitPlan
	}
	item.Code = strings.ToUpper(strings.TrimSpace(item.Code))
	item.Name = strings.TrimSpace(item.Name)
	item.PlanType = normalizeBenefitValue(item.PlanType, BenefitPlanTypeOther, []string{BenefitPlanTypeInsurance, BenefitPlanTypeReimbursement, BenefitPlanTypeAllowance, BenefitPlanTypeRetirement, BenefitPlanTypeWellness, BenefitPlanTypeOther})
	item.CurrencyCode = normalizeCurrency(item.CurrencyCode)
	item.Description = cleanOptional(item.Description)
	item.ProviderName = cleanOptional(item.ProviderName)
	item.PolicyNumber = cleanOptional(item.PolicyNumber)
	item.EligibilityRule = normalizeBenefitJSON(item.EligibilityRule)
	item.InsuranceMetadata = normalizeBenefitJSON(item.InsuranceMetadata)
	if item.EmployeeContribution < 0 || item.EmployerContribution < 0 || negativeFloatPtr(item.CoverageAmount) {
		return nil, ErrInvalidBenefitPlan
	}
	return &item, nil
}

func NewBenefitEnrollmentWindow(item BenefitEnrollmentWindow) (*BenefitEnrollmentWindow, error) {
	if item.TenantID == uuid.Nil || item.PlanID == uuid.Nil || strings.TrimSpace(item.Name) == "" || item.OpensOn.IsZero() || item.ClosesOn.IsZero() || item.ClosesOn.Before(item.OpensOn) {
		return nil, ErrInvalidBenefitEnrollmentWindow
	}
	item.Name = strings.TrimSpace(item.Name)
	item.Status = normalizeBenefitValue(item.Status, BenefitWindowStatusOpen, []string{BenefitWindowStatusDraft, BenefitWindowStatusOpen, BenefitWindowStatusClosed, BenefitWindowStatusArchived})
	item.Metadata = normalizeBenefitJSON(item.Metadata)
	return &item, nil
}

func NewBenefitDependent(item BenefitDependent) (*BenefitDependent, error) {
	if item.TenantID == uuid.Nil || item.EmployeeUserID == uuid.Nil || strings.TrimSpace(item.FullName) == "" {
		return nil, ErrInvalidBenefitDependent
	}
	item.FullName = strings.TrimSpace(item.FullName)
	item.Relationship = normalizeBenefitValue(item.Relationship, "other", []string{"spouse", "child", "parent", "sibling", "nominee", "other"})
	item.Gender = cleanOptional(item.Gender)
	item.Metadata = normalizeBenefitJSON(item.Metadata)
	if negativeFloatPtr(item.NomineePercentage) || (item.NomineePercentage != nil && *item.NomineePercentage > 100) {
		return nil, ErrInvalidBenefitDependent
	}
	return &item, nil
}

func NewBenefitEnrollment(item BenefitEnrollment) (*BenefitEnrollment, error) {
	if item.TenantID == uuid.Nil || item.PlanID == uuid.Nil || item.EmployeeUserID == uuid.Nil {
		return nil, ErrInvalidBenefitEnrollment
	}
	item.Status = normalizeBenefitValue(item.Status, BenefitEnrollmentStatusSubmitted, []string{BenefitEnrollmentStatusDraft, BenefitEnrollmentStatusSubmitted, BenefitEnrollmentStatusApproved, BenefitEnrollmentStatusRejected, BenefitEnrollmentStatusCancelled, BenefitEnrollmentStatusActive, BenefitEnrollmentStatusEnded})
	item.CoverageLevel = cleanOptional(item.CoverageLevel)
	item.ReviewRemarks = cleanOptional(item.ReviewRemarks)
	item.Metadata = normalizeBenefitJSON(item.Metadata)
	if item.EmployeeContribution < 0 || item.EmployerContribution < 0 || negativeFloatPtr(item.SelectedAmount) {
		return nil, ErrInvalidBenefitEnrollment
	}
	return &item, nil
}

func NewBenefitClaimType(item BenefitClaimType) (*BenefitClaimType, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidBenefitClaimType
	}
	item.Code = strings.ToUpper(strings.TrimSpace(item.Code))
	item.Name = strings.TrimSpace(item.Name)
	item.Description = cleanOptional(item.Description)
	item.PayrollComponentCode = cleanOptional(item.PayrollComponentCode)
	item.EligibilityRule = normalizeBenefitJSON(item.EligibilityRule)
	if negativeFloatPtr(item.AnnualLimit) || negativeFloatPtr(item.PerClaimLimit) {
		return nil, ErrInvalidBenefitClaimType
	}
	return &item, nil
}

func NewBenefitClaim(item BenefitClaim) (*BenefitClaim, error) {
	if item.TenantID == uuid.Nil || item.ClaimTypeID == uuid.Nil || item.EmployeeUserID == uuid.Nil || item.ExpenseDate.IsZero() || item.ClaimAmount <= 0 {
		return nil, ErrInvalidBenefitClaim
	}
	item.ClaimNumber = strings.TrimSpace(item.ClaimNumber)
	item.CurrencyCode = normalizeCurrency(item.CurrencyCode)
	item.Status = normalizeBenefitValue(item.Status, BenefitClaimStatusDraft, []string{BenefitClaimStatusDraft, BenefitClaimStatusSubmitted, BenefitClaimStatusUnderReview, BenefitClaimStatusApproved, BenefitClaimStatusRejected, BenefitClaimStatusCancelled, BenefitClaimStatusPaid})
	item.PaymentStatus = normalizeBenefitValue(item.PaymentStatus, BenefitPaymentStatusNotPayable, []string{BenefitPaymentStatusNotPayable, BenefitPaymentStatusPending, BenefitPaymentStatusPaid, BenefitPaymentStatusFailed})
	item.PayrollExportStatus = normalizeBenefitValue(item.PayrollExportStatus, BenefitPayrollExportStatusNotReady, []string{BenefitPayrollExportStatusNotReady, BenefitPayrollExportStatusReady, BenefitPayrollExportStatusExported, BenefitPayrollExportStatusBlocked})
	item.ReviewRemarks = cleanOptional(item.ReviewRemarks)
	item.PaymentReference = cleanOptional(item.PaymentReference)
	item.PayrollExportReference = cleanOptional(item.PayrollExportReference)
	item.Notes = cleanOptional(item.Notes)
	item.Metadata = normalizeBenefitJSON(item.Metadata)
	if negativeFloatPtr(item.ApprovedAmount) {
		return nil, ErrInvalidBenefitClaim
	}
	return &item, nil
}

func NewBenefitClaimAttachment(item BenefitClaimAttachment) (*BenefitClaimAttachment, error) {
	if item.TenantID == uuid.Nil || item.ClaimID == uuid.Nil || strings.TrimSpace(item.FileName) == "" || strings.TrimSpace(item.StoragePath) == "" {
		return nil, ErrInvalidBenefitClaimAttachment
	}
	item.FileName = strings.TrimSpace(item.FileName)
	item.ContentType = strings.TrimSpace(item.ContentType)
	if item.ContentType == "" {
		item.ContentType = "application/octet-stream"
	}
	item.Metadata = normalizeBenefitJSON(item.Metadata)
	return &item, nil
}

func normalizeBenefitJSON(value json.RawMessage) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(`{}`)
	}
	return value
}

func normalizeBenefitValue(value string, fallback string, allowed []string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		return fallback
	}
	for _, candidate := range allowed {
		if clean == candidate {
			return clean
		}
	}
	return fallback
}

func normalizeCurrency(value string) string {
	clean := strings.ToUpper(strings.TrimSpace(value))
	if clean == "" {
		return "INR"
	}
	return clean
}

func negativeFloatPtr(value *float64) bool {
	return value != nil && *value < 0
}
