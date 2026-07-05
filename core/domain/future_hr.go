package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPrivacyConsent      = errors.New("privacy consent is invalid")
	ErrInvalidDataErasureRequest  = errors.New("data erasure request is invalid")
	ErrInvalidIntegrationHook     = errors.New("ecosystem integration hook is invalid")
	ErrInvalidMobileAPIConstraint = errors.New("mobile api constraint is invalid")
	ErrInvalidBoundedAIAgentRun   = errors.New("bounded ai agent run is invalid")
)

type BoundedAIAgentDefinition struct {
	Key             string   `json:"key"`
	Name            string   `json:"name"`
	Workflow        string   `json:"workflow"`
	Severity        string   `json:"severity"`
	VisibilityScope string   `json:"visibility_scope"`
	Signals         []string `json:"signals"`
	Guardrails      []string `json:"guardrails"`
}

type BoundedAIAgentRunResult struct {
	Agents  []BoundedAIAgentDefinition `json:"agents"`
	Actions []*AIAgentActionLog        `json:"actions"`
}

type PeopleAnalyticsWorkspace struct {
	Workspace json.RawMessage `json:"workspace"`
}

type PrivacyConsent struct {
	ID              uuid.UUID       `json:"id"`
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
	GrantedAt       *time.Time      `json:"granted_at,omitempty"`
	RevokedAt       *time.Time      `json:"revoked_at,omitempty"`
	ExpiresAt       *time.Time      `json:"expires_at,omitempty"`
	Evidence        json.RawMessage `json:"evidence,omitempty"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
	Inactive        bool            `json:"inactive"`
	CreatedAt       time.Time       `json:"created_at"`
	CreatedBy       *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at"`
	UpdatedBy       *uuid.UUID      `json:"updated_by,omitempty"`
}

type DataErasureRequest struct {
	ID              uuid.UUID       `json:"id"`
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
	DueAt           *time.Time      `json:"due_at,omitempty"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty"`
	ReviewedBy      *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time      `json:"reviewed_at,omitempty"`
	AuditSummary    json.RawMessage `json:"audit_summary,omitempty"`
	Inactive        bool            `json:"inactive"`
	CreatedAt       time.Time       `json:"created_at"`
	CreatedBy       *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at"`
	UpdatedBy       *uuid.UUID      `json:"updated_by,omitempty"`
}

type EcosystemIntegrationHook struct {
	ID              uuid.UUID       `json:"id"`
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
	LastCheckedAt   *time.Time      `json:"last_checked_at,omitempty"`
	LastError       *string         `json:"last_error,omitempty"`
	Config          json.RawMessage `json:"config,omitempty"`
	Inactive        bool            `json:"inactive"`
	CreatedAt       time.Time       `json:"created_at"`
	CreatedBy       *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt       time.Time       `json:"updated_at"`
	UpdatedBy       *uuid.UUID      `json:"updated_by,omitempty"`
}

type MobileAPIConstraint struct {
	ID                    uuid.UUID       `json:"id"`
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
	Inactive              bool            `json:"inactive"`
	CreatedAt             time.Time       `json:"created_at"`
	CreatedBy             *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt             time.Time       `json:"updated_at"`
	UpdatedBy             *uuid.UUID      `json:"updated_by,omitempty"`
}

type PrivacyEcosystemFilter struct {
	TenantID    uuid.UUID
	Status      *string
	ConsentArea *string
	Priority    *string
	Channel     *string
	Workflow    *string
	Limit       int32
	Offset      int32
}

type PrivacyEcosystemWorkspace struct {
	Consents     []*PrivacyConsent           `json:"consents"`
	Erasure      []*DataErasureRequest       `json:"erasure_requests"`
	Integrations []*EcosystemIntegrationHook `json:"integrations"`
	Mobile       []*MobileAPIConstraint      `json:"mobile_constraints"`
	Summary      map[string]int32            `json:"summary"`
}

func NewPrivacyConsent(input PrivacyConsent) (*PrivacyConsent, error) {
	input.ConsentKey = strings.TrimSpace(input.ConsentKey)
	input.ConsentArea = strings.TrimSpace(input.ConsentArea)
	input.Purpose = strings.TrimSpace(input.Purpose)
	if input.TenantID == uuid.Nil || input.ConsentKey == "" || input.ConsentArea == "" || input.Purpose == "" {
		return nil, ErrInvalidPrivacyConsent
	}
	input.Status = defaultString(input.Status, "granted")
	input.LawfulBasis = defaultString(input.LawfulBasis, "consent")
	input.Channel = defaultString(input.Channel, "web")
	input.Source = defaultString(input.Source, "hrms")
	return &input, nil
}

func NewDataErasureRequest(input DataErasureRequest) (*DataErasureRequest, error) {
	input.RequestKey = strings.TrimSpace(input.RequestKey)
	input.Reason = strings.TrimSpace(input.Reason)
	if input.TenantID == uuid.Nil || input.RequestKey == "" || input.Reason == "" {
		return nil, ErrInvalidDataErasureRequest
	}
	input.RequestType = defaultString(input.RequestType, "erasure")
	input.Status = defaultString(input.Status, "intake")
	input.Priority = defaultString(input.Priority, "normal")
	return &input, nil
}

func NewEcosystemIntegrationHook(input EcosystemIntegrationHook) (*EcosystemIntegrationHook, error) {
	input.HookKey = strings.TrimSpace(input.HookKey)
	input.Provider = strings.TrimSpace(input.Provider)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	if input.TenantID == uuid.Nil || input.HookKey == "" || input.Provider == "" || input.DisplayName == "" {
		return nil, ErrInvalidIntegrationHook
	}
	input.Channel = defaultString(input.Channel, "webhook")
	input.Direction = defaultString(input.Direction, "outbound")
	input.Status = defaultString(input.Status, "draft")
	return &input, nil
}

func NewMobileAPIConstraint(input MobileAPIConstraint) (*MobileAPIConstraint, error) {
	input.ConstraintKey = strings.TrimSpace(input.ConstraintKey)
	input.Workflow = strings.TrimSpace(input.Workflow)
	if input.TenantID == uuid.Nil || input.ConstraintKey == "" || input.Workflow == "" {
		return nil, ErrInvalidMobileAPIConstraint
	}
	if input.MaxPayloadKB <= 0 {
		input.MaxPayloadKB = 256
	}
	input.Status = defaultString(input.Status, "active")
	return &input, nil
}

func defaultString(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
