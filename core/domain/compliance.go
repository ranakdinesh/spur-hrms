package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ComplianceCategoryCLRA      = "clra"
	ComplianceCategoryFixedTerm = "fixed_term"
	ComplianceCategoryGigWorker = "gig_worker"
	ComplianceCategoryTDS       = "tds"
	ComplianceCategoryPF        = "pf"
	ComplianceCategoryESIC      = "esic"
	ComplianceCategoryPT        = "pt"
	ComplianceCategoryLWF       = "lwf"
	ComplianceCategoryDocument  = "document"
	ComplianceCategorySafety    = "safety"
	ComplianceCategoryContract  = "contract"
	ComplianceCategoryCustom    = "custom"

	ComplianceScopeWorker             = "worker"
	ComplianceScopeEngagement         = "engagement"
	ComplianceScopeWorkerOrEngagement = "worker_or_engagement"

	ComplianceSeverityLow      = "low"
	ComplianceSeverityMedium   = "medium"
	ComplianceSeverityHigh     = "high"
	ComplianceSeverityCritical = "critical"

	ComplianceStatusPending       = "pending"
	ComplianceStatusInReview      = "in_review"
	ComplianceStatusCompliant     = "compliant"
	ComplianceStatusNonCompliant  = "non_compliant"
	ComplianceStatusWaived        = "waived"
	ComplianceStatusExpired       = "expired"
	ComplianceStatusNotApplicable = "not_applicable"
)

var (
	ErrInvalidComplianceRule           = errors.New("compliance rule is invalid")
	ErrComplianceRuleNotFound          = errors.New("compliance rule not found")
	ErrInvalidComplianceChecklistItem  = errors.New("compliance checklist item is invalid")
	ErrComplianceChecklistItemNotFound = errors.New("compliance checklist item not found")
)

type ComplianceRule struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	Code                string          `json:"code"`
	Title               string          `json:"title"`
	Description         *string         `json:"description,omitempty"`
	Category            string          `json:"category"`
	Scope               string          `json:"scope"`
	Severity            string          `json:"severity"`
	ClassificationGroup *string         `json:"classification_group,omitempty"`
	WorkerTypeID        *uuid.UUID      `json:"worker_type_id,omitempty"`
	EngagementType      *string         `json:"engagement_type,omitempty"`
	BranchID            *uuid.UUID      `json:"branch_id,omitempty"`
	DepartmentID        *uuid.UUID      `json:"department_id,omitempty"`
	CountryCode         string          `json:"country_code"`
	StateCode           *string         `json:"state_code,omitempty"`
	TriggerEvent        string          `json:"trigger_event"`
	DefaultDueDays      int32           `json:"default_due_days"`
	RecurringDays       *int32          `json:"recurring_days,omitempty"`
	RequiresEvidence    bool            `json:"requires_evidence"`
	EvidenceLabel       *string         `json:"evidence_label,omitempty"`
	AutoDetectKey       *string         `json:"auto_detect_key,omitempty"`
	BlocksPayroll       bool            `json:"blocks_payroll"`
	IsActive            bool            `json:"is_active"`
	EffectiveFrom       *time.Time      `json:"effective_from,omitempty"`
	EffectiveTo         *time.Time      `json:"effective_to,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
	WorkerTypeName      *string         `json:"worker_type_name,omitempty"`
	BranchName          *string         `json:"branch_name,omitempty"`
	DepartmentName      *string         `json:"department_name,omitempty"`
}

type ComplianceChecklistItem struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	RuleID              uuid.UUID       `json:"rule_id"`
	WorkerProfileID     *uuid.UUID      `json:"worker_profile_id,omitempty"`
	EngagementID        *uuid.UUID      `json:"engagement_id,omitempty"`
	Status              string          `json:"status"`
	DueDate             *time.Time      `json:"due_date,omitempty"`
	CompletedAt         *time.Time      `json:"completed_at,omitempty"`
	ReviewedAt          *time.Time      `json:"reviewed_at,omitempty"`
	ReviewedBy          *uuid.UUID      `json:"reviewed_by,omitempty"`
	EvidencePath        *string         `json:"evidence_path,omitempty"`
	EvidenceFileName    *string         `json:"evidence_file_name,omitempty"`
	EvidenceContentType *string         `json:"evidence_content_type,omitempty"`
	EvidenceUploadedAt  *time.Time      `json:"evidence_uploaded_at,omitempty"`
	EvidenceUploadedBy  *uuid.UUID      `json:"evidence_uploaded_by,omitempty"`
	WaiverReason        *string         `json:"waiver_reason,omitempty"`
	WaiverUntil         *time.Time      `json:"waiver_until,omitempty"`
	WaivedAt            *time.Time      `json:"waived_at,omitempty"`
	WaivedBy            *uuid.UUID      `json:"waived_by,omitempty"`
	DetectedValue       *string         `json:"detected_value,omitempty"`
	Notes               *string         `json:"notes,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
	RuleCode            *string         `json:"rule_code,omitempty"`
	RuleTitle           *string         `json:"rule_title,omitempty"`
	RuleCategory        *string         `json:"rule_category,omitempty"`
	RuleScope           *string         `json:"rule_scope,omitempty"`
	RuleSeverity        *string         `json:"rule_severity,omitempty"`
	RequiresEvidence    bool            `json:"requires_evidence"`
	EvidenceLabel       *string         `json:"evidence_label,omitempty"`
	BlocksPayroll       bool            `json:"blocks_payroll"`
	WorkerDisplayName   *string         `json:"worker_display_name,omitempty"`
	WorkerCode          *string         `json:"worker_code,omitempty"`
	EngagementTitle     *string         `json:"engagement_title,omitempty"`
	EngagementCode      *string         `json:"engagement_code,omitempty"`
}

type ComplianceEvent struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	ChecklistItemID *uuid.UUID      `json:"checklist_item_id,omitempty"`
	RuleID          *uuid.UUID      `json:"rule_id,omitempty"`
	EventType       string          `json:"event_type"`
	FromStatus      *string         `json:"from_status,omitempty"`
	ToStatus        *string         `json:"to_status,omitempty"`
	Comment         *string         `json:"comment,omitempty"`
	ActorID         *uuid.UUID      `json:"actor_id,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

type ComplianceRuleFilter struct {
	TenantID uuid.UUID
	Category *string
	Scope    *string
	Severity *string
	IsActive *bool
	Search   *string
}

type ComplianceChecklistFilter struct {
	TenantID        uuid.UUID
	WorkerProfileID *uuid.UUID
	EngagementID    *uuid.UUID
	RuleID          *uuid.UUID
	Status          *string
	Category        *string
	DueBefore       *time.Time
	Search          *string
}

type ComplianceSummaryRow struct {
	Category            string `json:"category"`
	Status              string `json:"status"`
	ItemCount           int32  `json:"item_count"`
	PayrollBlockerCount int32  `json:"payroll_blocker_count"`
	DueSoonCount        int32  `json:"due_soon_count"`
}

type ComplianceDefaultRule struct {
	Code                string
	Title               string
	Description         string
	Category            string
	Scope               string
	Severity            string
	ClassificationGroup string
	EngagementType      string
	TriggerEvent        string
	DefaultDueDays      int32
	RecurringDays       *int32
	RequiresEvidence    bool
	EvidenceLabel       string
	AutoDetectKey       string
	BlocksPayroll       bool
}

func NewComplianceRule(item ComplianceRule) (*ComplianceRule, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Title) == "" {
		return nil, ErrInvalidComplianceRule
	}
	category := normalizeWorkerProfileEnum(item.Category, ComplianceCategoryCustom)
	if !containsString(complianceCategories(), category) {
		return nil, ErrInvalidComplianceRule
	}
	scope := normalizeWorkerProfileEnum(item.Scope, ComplianceScopeWorker)
	if !containsString([]string{ComplianceScopeWorker, ComplianceScopeEngagement, ComplianceScopeWorkerOrEngagement}, scope) {
		return nil, ErrInvalidComplianceRule
	}
	severity := normalizeWorkerProfileEnum(item.Severity, ComplianceSeverityMedium)
	if !containsString([]string{ComplianceSeverityLow, ComplianceSeverityMedium, ComplianceSeverityHigh, ComplianceSeverityCritical}, severity) {
		return nil, ErrInvalidComplianceRule
	}
	if item.DefaultDueDays < 0 || (item.RecurringDays != nil && *item.RecurringDays <= 0) {
		return nil, ErrInvalidComplianceRule
	}
	if item.EffectiveFrom != nil && item.EffectiveTo != nil && item.EffectiveTo.Before(*item.EffectiveFrom) {
		return nil, ErrInvalidComplianceRule
	}
	metadata := normalizeWorkerJSONObject(item.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidComplianceRule
	}
	item.Code = strings.ToUpper(strings.TrimSpace(item.Code))
	item.Title = strings.TrimSpace(item.Title)
	item.Description = cleanOptional(item.Description)
	item.Category = category
	item.Scope = scope
	item.Severity = severity
	item.ClassificationGroup = cleanOptionalLower(item.ClassificationGroup)
	item.WorkerTypeID = cleanUUIDOptional(item.WorkerTypeID)
	item.EngagementType = cleanOptionalLower(item.EngagementType)
	item.BranchID = cleanUUIDOptional(item.BranchID)
	item.DepartmentID = cleanUUIDOptional(item.DepartmentID)
	item.CountryCode = strings.ToUpper(strings.TrimSpace(item.CountryCode))
	if item.CountryCode == "" {
		item.CountryCode = "IN"
	}
	item.StateCode = cleanOptionalUpper(item.StateCode)
	item.TriggerEvent = normalizeWorkerProfileEnum(item.TriggerEvent, "onboarding")
	item.EvidenceLabel = cleanOptional(item.EvidenceLabel)
	item.AutoDetectKey = cleanOptionalLower(item.AutoDetectKey)
	item.EffectiveFrom = datePtrUTC(item.EffectiveFrom)
	item.EffectiveTo = datePtrUTC(item.EffectiveTo)
	item.Metadata = metadata
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now
	return &item, nil
}

func NewComplianceChecklistItem(item ComplianceChecklistItem) (*ComplianceChecklistItem, error) {
	if item.TenantID == uuid.Nil || item.RuleID == uuid.Nil || (item.WorkerProfileID == nil && item.EngagementID == nil) {
		return nil, ErrInvalidComplianceChecklistItem
	}
	status := normalizeWorkerProfileEnum(item.Status, ComplianceStatusPending)
	if !containsString(complianceStatuses(), status) {
		return nil, ErrInvalidComplianceChecklistItem
	}
	if status == ComplianceStatusWaived && strings.TrimSpace(stringPtrValue(item.WaiverReason)) == "" {
		return nil, ErrInvalidComplianceChecklistItem
	}
	metadata := normalizeWorkerJSONObject(item.Metadata, "{}")
	if !json.Valid(metadata) || !jsonObject(metadata) {
		return nil, ErrInvalidComplianceChecklistItem
	}
	item.WorkerProfileID = cleanUUIDOptional(item.WorkerProfileID)
	item.EngagementID = cleanUUIDOptional(item.EngagementID)
	item.Status = status
	item.DueDate = datePtrUTC(item.DueDate)
	item.EvidencePath = cleanOptional(item.EvidencePath)
	item.EvidenceFileName = cleanOptional(item.EvidenceFileName)
	item.EvidenceContentType = cleanOptional(item.EvidenceContentType)
	item.WaiverReason = cleanOptional(item.WaiverReason)
	item.WaiverUntil = datePtrUTC(item.WaiverUntil)
	item.DetectedValue = cleanOptional(item.DetectedValue)
	item.Notes = cleanOptional(item.Notes)
	item.Metadata = metadata
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now
	return &item, nil
}

func ComplianceDefaultRules() []ComplianceDefaultRule {
	ninety := int32(90)
	yearly := int32(365)
	return []ComplianceDefaultRule{
		{Code: "DOC-ID", Title: "Identity and address evidence", Description: "Capture government identity/address evidence for the worker record.", Category: ComplianceCategoryDocument, Scope: ComplianceScopeWorker, Severity: ComplianceSeverityHigh, TriggerEvent: "onboarding", DefaultDueDays: 0, RequiresEvidence: true, EvidenceLabel: "Identity/address document", AutoDetectKey: "identity_document"},
		{Code: "CONTRACT-SIGNED", Title: "Agreement or SOW signed", Description: "Ensure the engagement has a signed agreement, offer, SOW, or contractor terms before work/payment.", Category: ComplianceCategoryContract, Scope: ComplianceScopeEngagement, Severity: ComplianceSeverityCritical, TriggerEvent: "engagement_start", DefaultDueDays: 0, RequiresEvidence: true, EvidenceLabel: "Signed agreement", AutoDetectKey: "agreement_signed", BlocksPayroll: true},
		{Code: "TDS-PAN", Title: "PAN and TDS section verified", Description: "Verify PAN and the applicable TDS section/rate before contractor or professional payments.", Category: ComplianceCategoryTDS, Scope: ComplianceScopeWorkerOrEngagement, Severity: ComplianceSeverityHigh, TriggerEvent: "payroll_ready", DefaultDueDays: 0, RequiresEvidence: true, EvidenceLabel: "PAN/TDS verification", AutoDetectKey: "tds_ready", BlocksPayroll: true},
		{Code: "CLRA-CONTRACTOR", Title: "CLRA contractor applicability reviewed", Description: "For contractor/agency workforce, record principal-employer and contractor compliance evidence where applicable.", Category: ComplianceCategoryCLRA, Scope: ComplianceScopeEngagement, Severity: ComplianceSeverityCritical, ClassificationGroup: "contractor", TriggerEvent: "engagement_start", DefaultDueDays: 0, RecurringDays: &yearly, RequiresEvidence: true, EvidenceLabel: "CLRA registration/license evidence", AutoDetectKey: "clra_applicable", BlocksPayroll: true},
		{Code: "PF-APPLICABILITY", Title: "PF applicability reviewed", Description: "Confirm PF applicability and evidence for eligible employee, fixed-term, contractor, or agency workers.", Category: ComplianceCategoryPF, Scope: ComplianceScopeWorkerOrEngagement, Severity: ComplianceSeverityHigh, TriggerEvent: "onboarding", DefaultDueDays: 7, RequiresEvidence: true, EvidenceLabel: "PF declaration/evidence", AutoDetectKey: "pf_applicable", BlocksPayroll: true},
		{Code: "ESIC-APPLICABILITY", Title: "ESIC applicability reviewed", Description: "Confirm ESIC coverage evidence or non-applicability for eligible workers.", Category: ComplianceCategoryESIC, Scope: ComplianceScopeWorkerOrEngagement, Severity: ComplianceSeverityHigh, TriggerEvent: "onboarding", DefaultDueDays: 7, RequiresEvidence: true, EvidenceLabel: "ESIC declaration/evidence", AutoDetectKey: "esic_applicable", BlocksPayroll: true},
		{Code: "PT-LWF-STATE", Title: "State PT/LWF applicability reviewed", Description: "Review Professional Tax and Labour Welfare Fund applicability based on the worker location/state.", Category: ComplianceCategoryPT, Scope: ComplianceScopeWorker, Severity: ComplianceSeverityMedium, TriggerEvent: "onboarding", DefaultDueDays: 7, RequiresEvidence: false, AutoDetectKey: "state_statutory_ready"},
		{Code: "FIXED-TERM-BENEFITS", Title: "Fixed-term benefits parity reviewed", Description: "For fixed-term workers, verify benefits and tenure terms before activation.", Category: ComplianceCategoryFixedTerm, Scope: ComplianceScopeEngagement, Severity: ComplianceSeverityHigh, EngagementType: EngagementTypeFixedTerm, TriggerEvent: "engagement_start", DefaultDueDays: 7, RequiresEvidence: true, EvidenceLabel: "Fixed-term review note", AutoDetectKey: "fixed_term_review"},
		{Code: "GIG-WORKER-TERMS", Title: "Gig worker terms and consent captured", Description: "Capture gig/platform worker terms, consent, and safety/contact evidence where applicable.", Category: ComplianceCategoryGigWorker, Scope: ComplianceScopeWorker, Severity: ComplianceSeverityMedium, ClassificationGroup: "gig", TriggerEvent: "onboarding", DefaultDueDays: 7, RecurringDays: &ninety, RequiresEvidence: true, EvidenceLabel: "Gig worker consent/evidence", AutoDetectKey: "gig_terms"},
	}
}

func ComplianceDueDateFrom(rule *ComplianceRule, base time.Time) *time.Time {
	if rule == nil {
		return nil
	}
	due := datePtrUTC(&base)
	if due == nil {
		now := time.Now().UTC()
		due = datePtrUTC(&now)
	}
	result := due.AddDate(0, 0, int(rule.DefaultDueDays))
	return &result
}

func complianceCategories() []string {
	return []string{ComplianceCategoryCLRA, ComplianceCategoryFixedTerm, ComplianceCategoryGigWorker, ComplianceCategoryTDS, ComplianceCategoryPF, ComplianceCategoryESIC, ComplianceCategoryPT, ComplianceCategoryLWF, ComplianceCategoryDocument, ComplianceCategorySafety, ComplianceCategoryContract, ComplianceCategoryCustom}
}

func complianceStatuses() []string {
	return []string{ComplianceStatusPending, ComplianceStatusInReview, ComplianceStatusCompliant, ComplianceStatusNonCompliant, ComplianceStatusWaived, ComplianceStatusExpired, ComplianceStatusNotApplicable}
}

func cleanOptionalLower(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.ToLower(strings.TrimSpace(*value))
	if clean == "" {
		return nil
	}
	return &clean
}

func cleanOptionalUpper(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.ToUpper(strings.TrimSpace(*value))
	if clean == "" {
		return nil
	}
	return &clean
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
