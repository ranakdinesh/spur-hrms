package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type ComplianceRepo interface {
	CreateComplianceRule(ctx context.Context, item *domain.ComplianceRule, actorID *uuid.UUID) (*domain.ComplianceRule, error)
	UpdateComplianceRule(ctx context.Context, item *domain.ComplianceRule, actorID *uuid.UUID) (*domain.ComplianceRule, error)
	GetComplianceRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ComplianceRule, error)
	GetComplianceRuleByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.ComplianceRule, error)
	ListComplianceRules(ctx context.Context, filter domain.ComplianceRuleFilter) ([]*domain.ComplianceRule, error)
	ListActiveComplianceRulesForWorker(ctx context.Context, tenantID uuid.UUID, workerProfileID uuid.UUID) ([]*domain.ComplianceRule, error)
	ListActiveComplianceRulesForEngagement(ctx context.Context, tenantID uuid.UUID, engagementID uuid.UUID) ([]*domain.ComplianceRule, error)
	DeleteComplianceRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateComplianceChecklistItem(ctx context.Context, item *domain.ComplianceChecklistItem, actorID *uuid.UUID) (*domain.ComplianceChecklistItem, error)
	GetComplianceChecklistItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ComplianceChecklistItem, error)
	ListComplianceChecklistItems(ctx context.Context, filter domain.ComplianceChecklistFilter) ([]*domain.ComplianceChecklistItem, error)
	UpdateComplianceChecklistStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, notes *string, actorID *uuid.UUID) (*domain.ComplianceChecklistItem, error)
	UpdateComplianceChecklistEvidence(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, evidencePath *string, evidenceFileName *string, evidenceContentType *string, notes *string, actorID *uuid.UUID) (*domain.ComplianceChecklistItem, error)
	WaiveComplianceChecklistItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, waiverReason string, waiverUntil *time.Time, notes *string, actorID *uuid.UUID) (*domain.ComplianceChecklistItem, error)
	RefreshComplianceChecklistDueStatus(ctx context.Context, tenantID uuid.UUID) error
	DeleteComplianceChecklistItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	GetComplianceSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.ComplianceSummaryRow, error)
	CreateComplianceEvent(ctx context.Context, item *domain.ComplianceEvent) (*domain.ComplianceEvent, error)
	ListComplianceEvents(ctx context.Context, tenantID uuid.UUID, checklistItemID *uuid.UUID, ruleID *uuid.UUID) ([]*domain.ComplianceEvent, error)
}

type ComplianceRuleCommand struct {
	TenantID            uuid.UUID       `json:"tenant_id"`
	ID                  uuid.UUID       `json:"id,omitempty"`
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
	ActorID             *uuid.UUID      `json:"actor_id,omitempty"`
}

type ComplianceChecklistGenerateCommand struct {
	TenantID        uuid.UUID  `json:"tenant_id"`
	WorkerProfileID *uuid.UUID `json:"worker_profile_id,omitempty"`
	EngagementID    *uuid.UUID `json:"engagement_id,omitempty"`
	ActorID         *uuid.UUID `json:"actor_id,omitempty"`
}

type ComplianceChecklistStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Notes    *string    `json:"notes,omitempty"`
	ActorID  *uuid.UUID `json:"actor_id,omitempty"`
}

type ComplianceEvidenceCommand struct {
	TenantID            uuid.UUID  `json:"tenant_id"`
	ID                  uuid.UUID  `json:"id"`
	EvidencePath        *string    `json:"evidence_path,omitempty"`
	EvidenceFileName    *string    `json:"evidence_file_name,omitempty"`
	EvidenceContentType *string    `json:"evidence_content_type,omitempty"`
	Notes               *string    `json:"notes,omitempty"`
	ActorID             *uuid.UUID `json:"actor_id,omitempty"`
}

type ComplianceWaiverCommand struct {
	TenantID     uuid.UUID  `json:"tenant_id"`
	ID           uuid.UUID  `json:"id"`
	WaiverReason string     `json:"waiver_reason"`
	WaiverUntil  *time.Time `json:"waiver_until,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	ActorID      *uuid.UUID `json:"actor_id,omitempty"`
}
