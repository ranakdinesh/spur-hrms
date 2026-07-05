package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type AssetAccessLifecycleRepo interface {
	CreateAssetItem(ctx context.Context, item *domain.AssetItem, actorID *uuid.UUID) (*domain.AssetItem, error)
	UpdateAssetItem(ctx context.Context, item *domain.AssetItem, actorID *uuid.UUID) (*domain.AssetItem, error)
	UpdateAssetItemStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.AssetItem, error)
	GetAssetItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AssetItem, error)
	ListAssetItems(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AssetItem, error)
	DeleteAssetItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateAccessCatalogItem(ctx context.Context, item *domain.AccessCatalogItem, actorID *uuid.UUID) (*domain.AccessCatalogItem, error)
	UpdateAccessCatalogItem(ctx context.Context, item *domain.AccessCatalogItem, actorID *uuid.UUID) (*domain.AccessCatalogItem, error)
	GetAccessCatalogItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AccessCatalogItem, error)
	ListAccessCatalogItems(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AccessCatalogItem, error)
	DeleteAccessCatalogItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateAssetAssignment(ctx context.Context, item *domain.AssetAssignment, actorID *uuid.UUID) (*domain.AssetAssignment, error)
	UpdateAssetAssignment(ctx context.Context, item *domain.AssetAssignment, actorID *uuid.UUID) (*domain.AssetAssignment, error)
	UpdateAssetAssignmentStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.AssetAssignment, error)
	ListAssetAssignments(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AssetAssignment, error)
	CreateAccessLifecycleTask(ctx context.Context, item *domain.AccessLifecycleTask, actorID *uuid.UUID) (*domain.AccessLifecycleTask, error)
	UpdateAccessLifecycleTask(ctx context.Context, item *domain.AccessLifecycleTask, actorID *uuid.UUID) (*domain.AccessLifecycleTask, error)
	UpdateAccessLifecycleTaskStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.AccessLifecycleTask, error)
	ListAccessLifecycleTasks(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AccessLifecycleTask, error)
	CreateAssetAccessEvent(ctx context.Context, item *domain.AssetAccessEvent, actorID *uuid.UUID) (*domain.AssetAccessEvent, error)
	ListAssetAccessEvents(ctx context.Context, filter domain.AssetAccessFilter, sourceType *string, sourceID *uuid.UUID) ([]*domain.AssetAccessEvent, error)
	GetAssetAccessSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.AssetAccessSummaryRow, error)
}

type AssetItemCommand struct {
	TenantID                 uuid.UUID       `json:"tenant_id"`
	ID                       uuid.UUID       `json:"id,omitempty"`
	AssetCode                string          `json:"asset_code"`
	AssetName                string          `json:"asset_name"`
	AssetType                string          `json:"asset_type"`
	Category                 string          `json:"category"`
	SerialNumber             *string         `json:"serial_number,omitempty"`
	Vendor                   *string         `json:"vendor,omitempty"`
	PurchaseDate             *time.Time      `json:"purchase_date,omitempty"`
	WarrantyUntil            *time.Time      `json:"warranty_until,omitempty"`
	OwnerUserID              *uuid.UUID      `json:"owner_user_id,omitempty"`
	CustodianWorkerProfileID *uuid.UUID      `json:"custodian_worker_profile_id,omitempty"`
	Status                   string          `json:"status"`
	LocationLabel            *string         `json:"location_label,omitempty"`
	Notes                    *string         `json:"notes,omitempty"`
	Metadata                 json.RawMessage `json:"metadata,omitempty"`
	ActorID                  *uuid.UUID      `json:"-"`
}

type AccessCatalogItemCommand struct {
	TenantID                 uuid.UUID       `json:"tenant_id"`
	ID                       uuid.UUID       `json:"id,omitempty"`
	AccessCode               string          `json:"access_code"`
	AccessName               string          `json:"access_name"`
	AccessType               string          `json:"access_type"`
	SystemName               *string         `json:"system_name,omitempty"`
	OwnerUserID              *uuid.UUID      `json:"owner_user_id,omitempty"`
	ProvisioningMethod       string          `json:"provisioning_method"`
	RequiresApproval         bool            `json:"requires_approval"`
	DefaultForOnboarding     bool            `json:"default_for_onboarding"`
	DefaultForExitRevocation bool            `json:"default_for_exit_revocation"`
	Status                   string          `json:"status"`
	Notes                    *string         `json:"notes,omitempty"`
	Metadata                 json.RawMessage `json:"metadata,omitempty"`
	ActorID                  *uuid.UUID      `json:"-"`
}

type AssetAssignmentCommand struct {
	TenantID              uuid.UUID       `json:"tenant_id"`
	ID                    uuid.UUID       `json:"id,omitempty"`
	AssetID               uuid.UUID       `json:"asset_id"`
	WorkerProfileID       uuid.UUID       `json:"worker_profile_id"`
	EmployeeID            *uuid.UUID      `json:"employee_id,omitempty"`
	CandidateOnboardingID *uuid.UUID      `json:"candidate_onboarding_id,omitempty"`
	ExitRequestID         *uuid.UUID      `json:"exit_request_id,omitempty"`
	RequestedBy           *uuid.UUID      `json:"requested_by,omitempty"`
	ApprovedBy            *uuid.UUID      `json:"approved_by,omitempty"`
	IssuedBy              *uuid.UUID      `json:"issued_by,omitempty"`
	ReturnedBy            *uuid.UUID      `json:"returned_by,omitempty"`
	ApprovedAt            *time.Time      `json:"approved_at,omitempty"`
	IssuedOn              *time.Time      `json:"issued_on,omitempty"`
	ExpectedReturnOn      *time.Time      `json:"expected_return_on,omitempty"`
	ReturnedOn            *time.Time      `json:"returned_on,omitempty"`
	IssueCondition        string          `json:"issue_condition"`
	ReturnCondition       *string         `json:"return_condition,omitempty"`
	DamageStatus          string          `json:"damage_status"`
	RecoveryAmount        float64         `json:"recovery_amount"`
	Status                string          `json:"status"`
	Notes                 *string         `json:"notes,omitempty"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`
	ActorID               *uuid.UUID      `json:"-"`
}

type AccessLifecycleTaskCommand struct {
	TenantID              uuid.UUID       `json:"tenant_id"`
	ID                    uuid.UUID       `json:"id,omitempty"`
	AccessItemID          uuid.UUID       `json:"access_item_id"`
	WorkerProfileID       uuid.UUID       `json:"worker_profile_id"`
	EmployeeID            *uuid.UUID      `json:"employee_id,omitempty"`
	CandidateOnboardingID *uuid.UUID      `json:"candidate_onboarding_id,omitempty"`
	ExitRequestID         *uuid.UUID      `json:"exit_request_id,omitempty"`
	TaskType              string          `json:"task_type"`
	RequestedBy           *uuid.UUID      `json:"requested_by,omitempty"`
	ApprovedBy            *uuid.UUID      `json:"approved_by,omitempty"`
	OwnerUserID           *uuid.UUID      `json:"owner_user_id,omitempty"`
	ApprovedAt            *time.Time      `json:"approved_at,omitempty"`
	DueDate               *time.Time      `json:"due_date,omitempty"`
	CompletedAt           *time.Time      `json:"completed_at,omitempty"`
	ExternalReference     *string         `json:"external_reference,omitempty"`
	Status                string          `json:"status"`
	Notes                 *string         `json:"notes,omitempty"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`
	ActorID               *uuid.UUID      `json:"-"`
}

type AssetAccessStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id"`
	Status   string     `json:"status"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}
