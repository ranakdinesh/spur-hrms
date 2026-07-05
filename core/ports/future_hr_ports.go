package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type PeopleAnalyticsRepo interface {
	GetPeopleAnalyticsWorkspace(ctx context.Context, tenantID uuid.UUID) (*domain.PeopleAnalyticsWorkspace, error)
}

type PrivacyEcosystemRepo interface {
	UpsertPrivacyConsent(ctx context.Context, item *domain.PrivacyConsent, actorID *uuid.UUID) (*domain.PrivacyConsent, error)
	ListPrivacyConsents(ctx context.Context, filter domain.PrivacyEcosystemFilter) ([]*domain.PrivacyConsent, error)
	CreateDataErasureRequest(ctx context.Context, item *domain.DataErasureRequest, actorID *uuid.UUID) (*domain.DataErasureRequest, error)
	ListDataErasureRequests(ctx context.Context, filter domain.PrivacyEcosystemFilter) ([]*domain.DataErasureRequest, error)
	UpdateDataErasureRequestStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, retainedReason *string, auditSummary json.RawMessage, actorID *uuid.UUID) (*domain.DataErasureRequest, error)
	UpsertEcosystemIntegrationHook(ctx context.Context, item *domain.EcosystemIntegrationHook, actorID *uuid.UUID) (*domain.EcosystemIntegrationHook, error)
	ListEcosystemIntegrationHooks(ctx context.Context, filter domain.PrivacyEcosystemFilter) ([]*domain.EcosystemIntegrationHook, error)
	UpsertMobileAPIConstraint(ctx context.Context, item *domain.MobileAPIConstraint, actorID *uuid.UUID) (*domain.MobileAPIConstraint, error)
	ListMobileAPIConstraints(ctx context.Context, filter domain.PrivacyEcosystemFilter) ([]*domain.MobileAPIConstraint, error)
}

type BoundedAIAgentRunCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	Agents   []string   `json:"agents"`
	ActorID  *uuid.UUID `json:"-"`
}

type PrivacyConsentCommand struct {
	TenantID        uuid.UUID       `json:"tenant_id"`
	EmployeeUserID  *uuid.UUID      `json:"employee_user_id,omitempty"`
	WorkerProfileID *uuid.UUID      `json:"worker_profile_id,omitempty"`
	ConsentKey      string          `json:"consent_key"`
	ConsentArea     string          `json:"consent_area"`
	Status          string          `json:"status"`
	LawfulBasis     string          `json:"lawful_basis"`
	Channel         string          `json:"channel"`
	Source          string          `json:"source"`
	Purpose         string          `json:"purpose"`
	GrantedAt       string          `json:"granted_at,omitempty"`
	RevokedAt       string          `json:"revoked_at,omitempty"`
	ExpiresAt       string          `json:"expires_at,omitempty"`
	Evidence        json.RawMessage `json:"evidence,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}

type DataErasureRequestCommand struct {
	TenantID        uuid.UUID       `json:"tenant_id"`
	RequestKey      string          `json:"request_key"`
	SubjectUserID   *uuid.UUID      `json:"subject_user_id,omitempty"`
	WorkerProfileID *uuid.UUID      `json:"worker_profile_id,omitempty"`
	RequestType     string          `json:"request_type"`
	Status          string          `json:"status"`
	Priority        string          `json:"priority"`
	RequestedBy     *uuid.UUID      `json:"requested_by,omitempty"`
	Reason          string          `json:"reason"`
	Scope           json.RawMessage `json:"scope,omitempty"`
	RetainedReason  *string         `json:"retained_reason,omitempty"`
	DueAt           string          `json:"due_at,omitempty"`
	AuditSummary    json.RawMessage `json:"audit_summary,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}

type DataErasureStatusCommand struct {
	TenantID       uuid.UUID       `json:"tenant_id"`
	ID             uuid.UUID       `json:"id"`
	Status         string          `json:"status"`
	RetainedReason *string         `json:"retained_reason,omitempty"`
	AuditSummary   json.RawMessage `json:"audit_summary,omitempty"`
	ActorID        *uuid.UUID      `json:"-"`
}

type EcosystemIntegrationHookCommand struct {
	TenantID        uuid.UUID       `json:"tenant_id"`
	HookKey         string          `json:"hook_key"`
	Provider        string          `json:"provider"`
	Channel         string          `json:"channel"`
	Direction       string          `json:"direction"`
	Status          string          `json:"status"`
	DisplayName     string          `json:"display_name"`
	EndpointURL     *string         `json:"endpoint_url,omitempty"`
	EventTypes      []string        `json:"event_types"`
	SecretRef       *string         `json:"secret_ref,omitempty"`
	ConsentRequired bool            `json:"consent_required"`
	MobileSafe      bool            `json:"mobile_safe"`
	Config          json.RawMessage `json:"config,omitempty"`
	ActorID         *uuid.UUID      `json:"-"`
}

type MobileAPIConstraintCommand struct {
	TenantID              uuid.UUID       `json:"tenant_id"`
	ConstraintKey         string          `json:"constraint_key"`
	Workflow              string          `json:"workflow"`
	MinAndroidVersion     *string         `json:"min_android_version,omitempty"`
	MinIOSVersion         *string         `json:"min_ios_version,omitempty"`
	OfflineSupported      bool            `json:"offline_supported"`
	LowBandwidthMode      bool            `json:"low_bandwidth_mode"`
	RequiresLocation      bool            `json:"requires_location"`
	RequiresDeviceBinding bool            `json:"requires_device_binding"`
	MaxPayloadKB          int32           `json:"max_payload_kb"`
	Status                string          `json:"status"`
	Notes                 *string         `json:"notes,omitempty"`
	Config                json.RawMessage `json:"config,omitempty"`
	ActorID               *uuid.UUID      `json:"-"`
}
