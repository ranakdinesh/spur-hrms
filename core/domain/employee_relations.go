package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ERCaseFamilyGrievance         = "grievance"
	ERCaseFamilyDisciplinary      = "disciplinary"
	ERCaseFamilyHarassment        = "harassment"
	ERCaseFamilyEthics            = "ethics"
	ERCaseFamilyWorkplaceConflict = "workplace_conflict"
	ERCaseFamilyPolicyViolation   = "policy_violation"
	ERCaseFamilyOther             = "other"

	ERSeverityLow      = "low"
	ERSeverityMedium   = "medium"
	ERSeverityHigh     = "high"
	ERSeverityCritical = "critical"

	ERCaseStatusIntake        = "intake"
	ERCaseStatusTriage        = "triage"
	ERCaseStatusInvestigation = "investigation"
	ERCaseStatusFindings      = "findings"
	ERCaseStatusActionPlan    = "action_plan"
	ERCaseStatusMonitoring    = "monitoring"
	ERCaseStatusClosed        = "closed"
	ERCaseStatusCancelled     = "cancelled"

	ERConfidentialityRestricted = "restricted"
	ERConfidentialitySensitive  = "sensitive"
	ERConfidentialityLegalHold  = "legal_hold"
)

var (
	ErrInvalidERCase           = errors.New("employee relations case is invalid")
	ErrERCaseNotFound          = errors.New("employee relations case not found")
	ErrInvalidERCaseCategory   = errors.New("employee relations category is invalid")
	ErrERCaseCategoryNotFound  = errors.New("employee relations category not found")
	ErrInvalidERCaseParty      = errors.New("employee relations case party is invalid")
	ErrInvalidERAllegation     = errors.New("employee relations allegation is invalid")
	ErrInvalidERStep           = errors.New("employee relations investigation step is invalid")
	ErrInvalidERWitnessNote    = errors.New("employee relations witness note is invalid")
	ErrInvalidEREvidence       = errors.New("employee relations evidence is invalid")
	ErrInvalidERFinding        = errors.New("employee relations finding is invalid")
	ErrInvalidERActionPlan     = errors.New("employee relations action plan is invalid")
	ErrInvalidERCaseTransition = errors.New("employee relations status transition is invalid")
)

type ERCaseCategory struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	Code             string     `json:"code"`
	Name             string     `json:"name"`
	CaseFamily       string     `json:"case_family"`
	Description      *string    `json:"description,omitempty"`
	DefaultSeverity  string     `json:"default_severity"`
	DefaultOwnerRole *string    `json:"default_owner_role,omitempty"`
	IsActive         bool       `json:"is_active"`
	Inactive         bool       `json:"inactive"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
}

type ERCase struct {
	ID                      uuid.UUID  `json:"id"`
	TenantID                uuid.UUID  `json:"tenant_id"`
	CaseNumber              string     `json:"case_number"`
	SourceHRCaseID          *uuid.UUID `json:"source_hr_case_id,omitempty"`
	CategoryID              *uuid.UUID `json:"category_id,omitempty"`
	CategoryName            *string    `json:"category_name,omitempty"`
	CategoryCode            *string    `json:"category_code,omitempty"`
	Title                   string     `json:"title"`
	IntakeSummary           string     `json:"intake_summary"`
	CaseFamily              string     `json:"case_family"`
	Severity                string     `json:"severity"`
	Status                  string     `json:"status"`
	ConfidentialityLevel    string     `json:"confidentiality_level"`
	ComplainantUserID       *uuid.UUID `json:"complainant_user_id,omitempty"`
	SubjectEmployeeUserID   *uuid.UUID `json:"subject_employee_user_id,omitempty"`
	OwnerUserID             *uuid.UUID `json:"owner_user_id,omitempty"`
	OwnerRole               *string    `json:"owner_role,omitempty"`
	InvestigationLeadUserID *uuid.UUID `json:"investigation_lead_user_id,omitempty"`
	LegalHold               bool       `json:"legal_hold"`
	LegalHoldReason         *string    `json:"legal_hold_reason,omitempty"`
	LegalHoldAt             *time.Time `json:"legal_hold_at,omitempty"`
	LegalHoldBy             *uuid.UUID `json:"legal_hold_by,omitempty"`
	DueAt                   *time.Time `json:"due_at,omitempty"`
	ClosedAt                *time.Time `json:"closed_at,omitempty"`
	ResolutionSummary       *string    `json:"resolution_summary,omitempty"`
	PrivacyNotes            *string    `json:"privacy_notes,omitempty"`
	Inactive                bool       `json:"inactive"`
	CreatedAt               time.Time  `json:"created_at"`
	CreatedBy               *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt               time.Time  `json:"updated_at"`
	UpdatedBy               *uuid.UUID `json:"updated_by,omitempty"`
	AllegationCount         int32      `json:"allegation_count"`
	EvidenceCount           int32      `json:"evidence_count"`
	OpenActionCount         int32      `json:"open_action_count"`
}

type ERCaseParty struct {
	ID                  uuid.UUID  `json:"id"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	ERCaseID            uuid.UUID  `json:"er_case_id"`
	PartyUserID         *uuid.UUID `json:"party_user_id,omitempty"`
	PartyName           *string    `json:"party_name,omitempty"`
	PartyRole           string     `json:"party_role"`
	RepresentationNotes *string    `json:"representation_notes,omitempty"`
	ContactNotes        *string    `json:"contact_notes,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt           time.Time  `json:"updated_at"`
	UpdatedBy           *uuid.UUID `json:"updated_by,omitempty"`
}

type ERAllegation struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	ERCaseID         uuid.UUID  `json:"er_case_id"`
	AllegationType   string     `json:"allegation_type"`
	IncidentDate     *time.Time `json:"incident_date,omitempty"`
	IncidentLocation *string    `json:"incident_location,omitempty"`
	Description      string     `json:"description"`
	PolicyReference  *string    `json:"policy_reference,omitempty"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
}

type ERInvestigationStep struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	ERCaseID     uuid.UUID  `json:"er_case_id"`
	StepType     string     `json:"step_type"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	OwnerUserID  *uuid.UUID `json:"owner_user_id,omitempty"`
	DueAt        *time.Time `json:"due_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Status       string     `json:"status"`
	OutcomeNotes *string    `json:"outcome_notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
}

type ERWitnessNote struct {
	ID                   uuid.UUID  `json:"id"`
	TenantID             uuid.UUID  `json:"tenant_id"`
	ERCaseID             uuid.UUID  `json:"er_case_id"`
	WitnessUserID        *uuid.UUID `json:"witness_user_id,omitempty"`
	WitnessName          *string    `json:"witness_name,omitempty"`
	InterviewAt          *time.Time `json:"interview_at,omitempty"`
	InterviewerUserID    *uuid.UUID `json:"interviewer_user_id,omitempty"`
	StatementSummary     string     `json:"statement_summary"`
	ConfidentialityLevel string     `json:"confidentiality_level"`
	CreatedAt            time.Time  `json:"created_at"`
	CreatedBy            *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by,omitempty"`
}

type EREvidenceAttachment struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	ERCaseID     uuid.UUID  `json:"er_case_id"`
	AllegationID *uuid.UUID `json:"allegation_id,omitempty"`
	FileName     string     `json:"file_name"`
	ContentType  string     `json:"content_type"`
	StoragePath  string     `json:"storage_path"`
	ChecksumSHA  *string    `json:"checksum_sha256,omitempty"`
	SizeBytes    int64      `json:"size_bytes"`
	EvidenceType string     `json:"evidence_type"`
	Description  *string    `json:"description,omitempty"`
	UploadedBy   *uuid.UUID `json:"uploaded_by,omitempty"`
	LegalHold    bool       `json:"legal_hold"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
}

type ERFinding struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	ERCaseID          uuid.UUID  `json:"er_case_id"`
	AllegationID      *uuid.UUID `json:"allegation_id,omitempty"`
	Finding           string     `json:"finding"`
	Rationale         string     `json:"rationale"`
	RecommendedAction *string    `json:"recommended_action,omitempty"`
	DecidedBy         *uuid.UUID `json:"decided_by,omitempty"`
	DecidedAt         *time.Time `json:"decided_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty"`
}

type ERActionPlan struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	ERCaseID         uuid.UUID  `json:"er_case_id"`
	ActionType       string     `json:"action_type"`
	Description      string     `json:"description"`
	AssignedToUserID *uuid.UUID `json:"assigned_to_user_id,omitempty"`
	DueAt            *time.Time `json:"due_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	Status           string     `json:"status"`
	FollowUpNotes    *string    `json:"follow_up_notes,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
}

type ERCaseEvent struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	ERCaseID    uuid.UUID       `json:"er_case_id"`
	EventType   string          `json:"event_type"`
	FromStatus  *string         `json:"from_status,omitempty"`
	ToStatus    *string         `json:"to_status,omitempty"`
	ActorUserID *uuid.UUID      `json:"actor_user_id,omitempty"`
	Comment     *string         `json:"comment,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty"`
}

type ERCaseFilter struct {
	TenantID              uuid.UUID
	Status                *string
	Severity              *string
	CaseFamily            *string
	CategoryID            *uuid.UUID
	OwnerUserID           *uuid.UUID
	SubjectEmployeeUserID *uuid.UUID
	ComplainantUserID     *uuid.UUID
	LegalHold             *bool
	Search                *string
	Limit                 int32
	Offset                int32
}

type ERCaseSummaryRow struct {
	Status         string `json:"status"`
	Severity       string `json:"severity"`
	CaseCount      int32  `json:"case_count"`
	LegalHoldCount int32  `json:"legal_hold_count"`
	OverdueCount   int32  `json:"overdue_count"`
}

type ERCasePage struct {
	Items      []*ERCase           `json:"items"`
	Total      int64               `json:"total"`
	Summary    []*ERCaseSummaryRow `json:"summary"`
	Categories []*ERCaseCategory   `json:"categories,omitempty"`
}

type ERCaseWorkspace struct {
	Case        *ERCase                 `json:"case"`
	Parties     []*ERCaseParty          `json:"parties"`
	Allegations []*ERAllegation         `json:"allegations"`
	Steps       []*ERInvestigationStep  `json:"steps"`
	Witnesses   []*ERWitnessNote        `json:"witnesses"`
	Evidence    []*EREvidenceAttachment `json:"evidence"`
	Findings    []*ERFinding            `json:"findings"`
	Actions     []*ERActionPlan         `json:"actions"`
	Events      []*ERCaseEvent          `json:"events"`
}

func NewERCaseCategory(item ERCaseCategory) (*ERCaseCategory, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidERCaseCategory
	}
	item.Code = strings.ToUpper(strings.TrimSpace(item.Code))
	item.Name = strings.TrimSpace(item.Name)
	item.CaseFamily = normalizeERFamily(item.CaseFamily)
	item.DefaultSeverity = normalizeERSeverity(item.DefaultSeverity)
	item.Description = cleanOptional(item.Description)
	item.DefaultOwnerRole = cleanOptional(item.DefaultOwnerRole)
	return &item, nil
}

func NewERCase(item ERCase) (*ERCase, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.IntakeSummary) == "" {
		return nil, ErrInvalidERCase
	}
	item.Title = strings.TrimSpace(item.Title)
	item.IntakeSummary = strings.TrimSpace(item.IntakeSummary)
	item.CaseFamily = normalizeERFamily(item.CaseFamily)
	item.Severity = normalizeERSeverity(item.Severity)
	item.Status = normalizeERStatus(item.Status)
	item.ConfidentialityLevel = normalizeERConfidentiality(item.ConfidentialityLevel)
	item.OwnerRole = cleanOptional(item.OwnerRole)
	item.LegalHoldReason = cleanOptional(item.LegalHoldReason)
	item.ResolutionSummary = cleanOptional(item.ResolutionSummary)
	item.PrivacyNotes = cleanOptional(item.PrivacyNotes)
	return &item, nil
}

func NewERCaseParty(item ERCaseParty) (*ERCaseParty, error) {
	if item.TenantID == uuid.Nil || item.ERCaseID == uuid.Nil || (item.PartyUserID == nil && strings.TrimSpace(valueFromPtr(item.PartyName)) == "") {
		return nil, ErrInvalidERCaseParty
	}
	item.PartyRole = normalizeERPartyRole(item.PartyRole)
	item.PartyName = cleanOptional(item.PartyName)
	item.RepresentationNotes = cleanOptional(item.RepresentationNotes)
	item.ContactNotes = cleanOptional(item.ContactNotes)
	return &item, nil
}

func NewERAllegation(item ERAllegation) (*ERAllegation, error) {
	if item.TenantID == uuid.Nil || item.ERCaseID == uuid.Nil || strings.TrimSpace(item.Description) == "" {
		return nil, ErrInvalidERAllegation
	}
	item.AllegationType = defaultString(strings.TrimSpace(item.AllegationType), "general")
	item.Description = strings.TrimSpace(item.Description)
	item.Status = normalizeERAllegationStatus(item.Status)
	item.IncidentLocation = cleanOptional(item.IncidentLocation)
	item.PolicyReference = cleanOptional(item.PolicyReference)
	return &item, nil
}

func NewERInvestigationStep(item ERInvestigationStep) (*ERInvestigationStep, error) {
	if item.TenantID == uuid.Nil || item.ERCaseID == uuid.Nil || strings.TrimSpace(item.Title) == "" {
		return nil, ErrInvalidERStep
	}
	item.StepType = defaultString(strings.TrimSpace(item.StepType), "investigation")
	item.Title = strings.TrimSpace(item.Title)
	item.Status = normalizeERStepStatus(item.Status)
	item.Description = cleanOptional(item.Description)
	item.OutcomeNotes = cleanOptional(item.OutcomeNotes)
	return &item, nil
}

func NewERWitnessNote(item ERWitnessNote) (*ERWitnessNote, error) {
	if item.TenantID == uuid.Nil || item.ERCaseID == uuid.Nil || strings.TrimSpace(item.StatementSummary) == "" || (item.WitnessUserID == nil && strings.TrimSpace(valueFromPtr(item.WitnessName)) == "") {
		return nil, ErrInvalidERWitnessNote
	}
	item.WitnessName = cleanOptional(item.WitnessName)
	item.StatementSummary = strings.TrimSpace(item.StatementSummary)
	item.ConfidentialityLevel = normalizeERConfidentiality(item.ConfidentialityLevel)
	return &item, nil
}

func NewEREvidenceAttachment(item EREvidenceAttachment) (*EREvidenceAttachment, error) {
	if item.TenantID == uuid.Nil || item.ERCaseID == uuid.Nil || strings.TrimSpace(item.FileName) == "" || strings.TrimSpace(item.StoragePath) == "" {
		return nil, ErrInvalidEREvidence
	}
	item.FileName = strings.TrimSpace(item.FileName)
	item.ContentType = defaultString(strings.TrimSpace(item.ContentType), "application/octet-stream")
	item.StoragePath = strings.TrimSpace(item.StoragePath)
	item.EvidenceType = defaultString(strings.TrimSpace(item.EvidenceType), "document")
	item.Description = cleanOptional(item.Description)
	item.ChecksumSHA = cleanOptional(item.ChecksumSHA)
	return &item, nil
}

func NewERFinding(item ERFinding) (*ERFinding, error) {
	if item.TenantID == uuid.Nil || item.ERCaseID == uuid.Nil || strings.TrimSpace(item.Rationale) == "" {
		return nil, ErrInvalidERFinding
	}
	item.Finding = normalizeERFinding(item.Finding)
	item.Rationale = strings.TrimSpace(item.Rationale)
	item.RecommendedAction = cleanOptional(item.RecommendedAction)
	return &item, nil
}

func NewERActionPlan(item ERActionPlan) (*ERActionPlan, error) {
	if item.TenantID == uuid.Nil || item.ERCaseID == uuid.Nil || strings.TrimSpace(item.Description) == "" {
		return nil, ErrInvalidERActionPlan
	}
	item.ActionType = defaultString(strings.TrimSpace(item.ActionType), "corrective_action")
	item.Description = strings.TrimSpace(item.Description)
	item.Status = normalizeERActionStatus(item.Status)
	item.FollowUpNotes = cleanOptional(item.FollowUpNotes)
	return &item, nil
}

func NewERCaseNumber(now time.Time) string {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return fmt.Sprintf("ER-%s-%06d", now.Format("20060102"), now.UnixNano()%1000000)
}

func normalizeERFamily(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case ERCaseFamilyDisciplinary, ERCaseFamilyHarassment, ERCaseFamilyEthics, ERCaseFamilyWorkplaceConflict, ERCaseFamilyPolicyViolation, ERCaseFamilyOther:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ERCaseFamilyGrievance
	}
}

func normalizeERSeverity(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case ERSeverityLow, ERSeverityHigh, ERSeverityCritical:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ERSeverityMedium
	}
}

func normalizeERStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case ERCaseStatusTriage, ERCaseStatusInvestigation, ERCaseStatusFindings, ERCaseStatusActionPlan, ERCaseStatusMonitoring, ERCaseStatusClosed, ERCaseStatusCancelled:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ERCaseStatusIntake
	}
}

func normalizeERConfidentiality(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case ERConfidentialitySensitive, ERConfidentialityLegalHold:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ERConfidentialityRestricted
	}
}

func normalizeERPartyRole(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "complainant", "respondent", "witness", "investigator", "hr_partner", "legal", "manager":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "other"
	}
}

func normalizeERAllegationStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "substantiated", "unsubstantiated", "inconclusive", "withdrawn":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "open"
	}
}

func normalizeERStepStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "in_progress", "completed", "skipped", "blocked":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "pending"
	}
}

func normalizeERFinding(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "unsubstantiated", "inconclusive", "partially_substantiated", "withdrawn":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "substantiated"
	}
}

func normalizeERActionStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "in_progress", "completed", "cancelled", "overdue":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "pending"
	}
}

func valueFromPtr(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
