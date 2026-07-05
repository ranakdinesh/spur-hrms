package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type EmployeeRelationsRepo interface {
	CreateERCaseCategory(ctx context.Context, item *domain.ERCaseCategory, actorID *uuid.UUID) (*domain.ERCaseCategory, error)
	UpdateERCaseCategory(ctx context.Context, item *domain.ERCaseCategory, actorID *uuid.UUID) (*domain.ERCaseCategory, error)
	ListERCaseCategories(ctx context.Context, tenantID uuid.UUID, activeOnly *bool) ([]*domain.ERCaseCategory, error)
	GetERCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ERCaseCategory, error)
	DeleteERCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateERCase(ctx context.Context, item *domain.ERCase, actorID *uuid.UUID) (*domain.ERCase, error)
	UpdateERCase(ctx context.Context, item *domain.ERCase, actorID *uuid.UUID) (*domain.ERCase, error)
	UpdateERCaseStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, resolutionSummary *string, actorID *uuid.UUID) (*domain.ERCase, error)
	UpdateERCaseLegalHold(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, enabled bool, reason *string, actorID *uuid.UUID) (*domain.ERCase, error)
	ListERCases(ctx context.Context, filter domain.ERCaseFilter) ([]*domain.ERCase, error)
	CountERCases(ctx context.Context, filter domain.ERCaseFilter) (int64, error)
	GetERCase(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ERCase, error)
	CreateERCaseParty(ctx context.Context, item *domain.ERCaseParty, actorID *uuid.UUID) (*domain.ERCaseParty, error)
	ListERCaseParties(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERCaseParty, error)
	CreateERAllegation(ctx context.Context, item *domain.ERAllegation, actorID *uuid.UUID) (*domain.ERAllegation, error)
	ListERAllegations(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERAllegation, error)
	CreateERInvestigationStep(ctx context.Context, item *domain.ERInvestigationStep, actorID *uuid.UUID) (*domain.ERInvestigationStep, error)
	UpdateERInvestigationStepStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, completedAt *time.Time, outcomeNotes *string, actorID *uuid.UUID) (*domain.ERInvestigationStep, error)
	ListERInvestigationSteps(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERInvestigationStep, error)
	CreateERWitnessNote(ctx context.Context, item *domain.ERWitnessNote, actorID *uuid.UUID) (*domain.ERWitnessNote, error)
	ListERWitnessNotes(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERWitnessNote, error)
	CreateEREvidenceAttachment(ctx context.Context, item *domain.EREvidenceAttachment, actorID *uuid.UUID) (*domain.EREvidenceAttachment, error)
	ListEREvidenceAttachments(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.EREvidenceAttachment, error)
	CreateERFinding(ctx context.Context, item *domain.ERFinding, actorID *uuid.UUID) (*domain.ERFinding, error)
	ListERFindings(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERFinding, error)
	CreateERActionPlan(ctx context.Context, item *domain.ERActionPlan, actorID *uuid.UUID) (*domain.ERActionPlan, error)
	UpdateERActionPlanStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, completedAt *time.Time, followUpNotes *string, actorID *uuid.UUID) (*domain.ERActionPlan, error)
	ListERActionPlans(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERActionPlan, error)
	CreateERCaseEvent(ctx context.Context, item *domain.ERCaseEvent, actorID *uuid.UUID) (*domain.ERCaseEvent, error)
	ListERCaseEvents(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.ERCaseEvent, error)
	GetERCaseSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.ERCaseSummaryRow, error)
}

type ERCaseCategoryCommand struct {
	ID               uuid.UUID  `json:"id,omitempty"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	Code             string     `json:"code"`
	Name             string     `json:"name"`
	CaseFamily       string     `json:"case_family"`
	Description      *string    `json:"description,omitempty"`
	DefaultSeverity  string     `json:"default_severity"`
	DefaultOwnerRole *string    `json:"default_owner_role,omitempty"`
	IsActive         bool       `json:"is_active"`
	ActorID          *uuid.UUID `json:"-"`
}

type ERCaseCommand struct {
	ID                      uuid.UUID  `json:"id,omitempty"`
	TenantID                uuid.UUID  `json:"tenant_id"`
	SourceHRCaseID          *uuid.UUID `json:"source_hr_case_id,omitempty"`
	CategoryID              *uuid.UUID `json:"category_id,omitempty"`
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
	DueAt                   *string    `json:"due_at,omitempty"`
	PrivacyNotes            *string    `json:"privacy_notes,omitempty"`
	ActorID                 *uuid.UUID `json:"-"`
}

type ERCaseStatusCommand struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	ID                uuid.UUID  `json:"id"`
	Status            string     `json:"status"`
	ResolutionSummary *string    `json:"resolution_summary,omitempty"`
	Comment           *string    `json:"comment,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}

type ERCaseLegalHoldCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Enabled  bool       `json:"enabled"`
	Reason   *string    `json:"reason,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type ERCasePartyCommand struct {
	TenantID            uuid.UUID  `json:"tenant_id"`
	ERCaseID            uuid.UUID  `json:"er_case_id"`
	PartyUserID         *uuid.UUID `json:"party_user_id,omitempty"`
	PartyName           *string    `json:"party_name,omitempty"`
	PartyRole           string     `json:"party_role"`
	RepresentationNotes *string    `json:"representation_notes,omitempty"`
	ContactNotes        *string    `json:"contact_notes,omitempty"`
	ActorID             *uuid.UUID `json:"-"`
}

type ERAllegationCommand struct {
	TenantID         uuid.UUID  `json:"tenant_id"`
	ERCaseID         uuid.UUID  `json:"er_case_id"`
	AllegationType   string     `json:"allegation_type"`
	IncidentDate     *string    `json:"incident_date,omitempty"`
	IncidentLocation *string    `json:"incident_location,omitempty"`
	Description      string     `json:"description"`
	PolicyReference  *string    `json:"policy_reference,omitempty"`
	Status           string     `json:"status"`
	ActorID          *uuid.UUID `json:"-"`
}

type ERInvestigationStepCommand struct {
	TenantID     uuid.UUID  `json:"tenant_id"`
	ID           uuid.UUID  `json:"id,omitempty"`
	ERCaseID     uuid.UUID  `json:"er_case_id"`
	StepType     string     `json:"step_type"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	OwnerUserID  *uuid.UUID `json:"owner_user_id,omitempty"`
	DueAt        *string    `json:"due_at,omitempty"`
	CompletedAt  *string    `json:"completed_at,omitempty"`
	Status       string     `json:"status"`
	OutcomeNotes *string    `json:"outcome_notes,omitempty"`
	ActorID      *uuid.UUID `json:"-"`
}

type ERWitnessNoteCommand struct {
	TenantID             uuid.UUID  `json:"tenant_id"`
	ERCaseID             uuid.UUID  `json:"er_case_id"`
	WitnessUserID        *uuid.UUID `json:"witness_user_id,omitempty"`
	WitnessName          *string    `json:"witness_name,omitempty"`
	InterviewAt          *string    `json:"interview_at,omitempty"`
	InterviewerUserID    *uuid.UUID `json:"interviewer_user_id,omitempty"`
	StatementSummary     string     `json:"statement_summary"`
	ConfidentialityLevel string     `json:"confidentiality_level"`
	ActorID              *uuid.UUID `json:"-"`
}

type EREvidenceAttachmentCommand struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	ERCaseID          uuid.UUID  `json:"er_case_id"`
	AllegationID      *uuid.UUID `json:"allegation_id,omitempty"`
	FileName          string     `json:"file_name"`
	ContentType       string     `json:"content_type"`
	FileContentBase64 string     `json:"file_content_base64"`
	EvidenceType      string     `json:"evidence_type"`
	Description       *string    `json:"description,omitempty"`
	LegalHold         bool       `json:"legal_hold"`
	ActorID           *uuid.UUID `json:"-"`
}

type ERFindingCommand struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	ERCaseID          uuid.UUID  `json:"er_case_id"`
	AllegationID      *uuid.UUID `json:"allegation_id,omitempty"`
	Finding           string     `json:"finding"`
	Rationale         string     `json:"rationale"`
	RecommendedAction *string    `json:"recommended_action,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}

type ERActionPlanCommand struct {
	TenantID         uuid.UUID  `json:"tenant_id"`
	ID               uuid.UUID  `json:"id,omitempty"`
	ERCaseID         uuid.UUID  `json:"er_case_id"`
	ActionType       string     `json:"action_type"`
	Description      string     `json:"description"`
	AssignedToUserID *uuid.UUID `json:"assigned_to_user_id,omitempty"`
	DueAt            *string    `json:"due_at,omitempty"`
	CompletedAt      *string    `json:"completed_at,omitempty"`
	Status           string     `json:"status"`
	FollowUpNotes    *string    `json:"follow_up_notes,omitempty"`
	ActorID          *uuid.UUID `json:"-"`
}
