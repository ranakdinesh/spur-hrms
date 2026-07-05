package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	AssetStatusAvailable   = "available"
	AssetStatusReserved    = "reserved"
	AssetStatusIssued      = "issued"
	AssetStatusReturnDue   = "return_due"
	AssetStatusReturned    = "returned"
	AssetStatusMaintenance = "maintenance"
	AssetStatusDamaged     = "damaged"
	AssetStatusLost        = "lost"
	AssetStatusRetired     = "retired"

	AssetAssignmentRequested = "requested"
	AssetAssignmentApproved  = "approved"
	AssetAssignmentIssued    = "issued"
	AssetAssignmentReturnDue = "return_due"
	AssetAssignmentReturned  = "returned"
	AssetAssignmentDamaged   = "damaged"
	AssetAssignmentLost      = "lost"
	AssetAssignmentCancelled = "cancelled"

	AccessCatalogActive     = "active"
	AccessCatalogInactive   = "inactive"
	AccessCatalogDeprecated = "deprecated"

	AccessTaskProvision   = "provision"
	AccessTaskDeprovision = "deprovision"
	AccessTaskReview      = "review"
	AccessTaskChange      = "change"

	AccessTaskRequested   = "requested"
	AccessTaskApproved    = "approved"
	AccessTaskProvisioned = "provisioned"
	AccessTaskRevoked     = "revoked"
	AccessTaskReviewed    = "reviewed"
	AccessTaskRejected    = "rejected"
	AccessTaskCancelled   = "cancelled"
	AccessTaskBlocked     = "blocked"
)

var (
	ErrInvalidAssetItem        = errors.New("invalid asset item")
	ErrAssetItemNotFound       = errors.New("asset item not found")
	ErrInvalidAccessItem       = errors.New("invalid access catalog item")
	ErrAccessItemNotFound      = errors.New("access catalog item not found")
	ErrInvalidAssetAssignment  = errors.New("invalid asset assignment")
	ErrAssetAssignmentNotFound = errors.New("asset assignment not found")
	ErrInvalidAccessTask       = errors.New("invalid access lifecycle task")
	ErrAccessTaskNotFound      = errors.New("access lifecycle task not found")
)

type AssetItem struct {
	ID                       uuid.UUID       `json:"id"`
	TenantID                 uuid.UUID       `json:"tenant_id"`
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
	CurrentAssignmentID      *uuid.UUID      `json:"current_assignment_id,omitempty"`
	LocationLabel            *string         `json:"location_label,omitempty"`
	Notes                    *string         `json:"notes,omitempty"`
	Metadata                 json.RawMessage `json:"metadata,omitempty"`
	Inactive                 bool            `json:"inactive"`
	CreatedAt                time.Time       `json:"created_at"`
	CreatedBy                *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                time.Time       `json:"updated_at"`
	UpdatedBy                *uuid.UUID      `json:"updated_by,omitempty"`
	CustodianName            *string         `json:"custodian_name,omitempty"`
	CustodianCode            *string         `json:"custodian_code,omitempty"`
}

type AccessCatalogItem struct {
	ID                       uuid.UUID       `json:"id"`
	TenantID                 uuid.UUID       `json:"tenant_id"`
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
	Inactive                 bool            `json:"inactive"`
	CreatedAt                time.Time       `json:"created_at"`
	CreatedBy                *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                time.Time       `json:"updated_at"`
	UpdatedBy                *uuid.UUID      `json:"updated_by,omitempty"`
}

type AssetAssignment struct {
	ID                    uuid.UUID       `json:"id"`
	TenantID              uuid.UUID       `json:"tenant_id"`
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
	Inactive              bool            `json:"inactive"`
	CreatedAt             time.Time       `json:"created_at"`
	CreatedBy             *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt             time.Time       `json:"updated_at"`
	UpdatedBy             *uuid.UUID      `json:"updated_by,omitempty"`
	AssetCode             *string         `json:"asset_code,omitempty"`
	AssetName             *string         `json:"asset_name,omitempty"`
	AssetType             *string         `json:"asset_type,omitempty"`
	Category              *string         `json:"category,omitempty"`
	WorkerDisplayName     *string         `json:"worker_display_name,omitempty"`
	WorkerCode            *string         `json:"worker_code,omitempty"`
}

type AccessLifecycleTask struct {
	ID                    uuid.UUID       `json:"id"`
	TenantID              uuid.UUID       `json:"tenant_id"`
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
	Inactive              bool            `json:"inactive"`
	CreatedAt             time.Time       `json:"created_at"`
	CreatedBy             *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt             time.Time       `json:"updated_at"`
	UpdatedBy             *uuid.UUID      `json:"updated_by,omitempty"`
	AccessCode            *string         `json:"access_code,omitempty"`
	AccessName            *string         `json:"access_name,omitempty"`
	AccessType            *string         `json:"access_type,omitempty"`
	SystemName            *string         `json:"system_name,omitempty"`
	WorkerDisplayName     *string         `json:"worker_display_name,omitempty"`
	WorkerCode            *string         `json:"worker_code,omitempty"`
}

type AssetAccessEvent struct {
	ID         uuid.UUID       `json:"id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	SourceType string          `json:"source_type"`
	SourceID   *uuid.UUID      `json:"source_id,omitempty"`
	Action     string          `json:"action"`
	FromStatus *string         `json:"from_status,omitempty"`
	ToStatus   *string         `json:"to_status,omitempty"`
	Remarks    *string         `json:"remarks,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	Inactive   bool            `json:"inactive"`
	CreatedAt  time.Time       `json:"created_at"`
	CreatedBy  *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt  time.Time       `json:"updated_at"`
	UpdatedBy  *uuid.UUID      `json:"updated_by,omitempty"`
}

type AssetAccessSummaryRow struct {
	Metric      string `json:"metric"`
	MetricCount int64  `json:"metric_count"`
}

type AssetAccessFilter struct {
	TenantID        uuid.UUID
	WorkerProfileID *uuid.UUID
	AssetID         *uuid.UUID
	AccessItemID    *uuid.UUID
	ExitRequestID   *uuid.UUID
	Status          *string
	Category        *string
	AccessType      *string
	Search          *string
	Limit           int32
	Offset          int32
}

func NormalizeAssetStatus(value string) string {
	return normalizeAssetAllowed(value, []string{AssetStatusAvailable, AssetStatusReserved, AssetStatusIssued, AssetStatusReturnDue, AssetStatusReturned, AssetStatusMaintenance, AssetStatusDamaged, AssetStatusLost, AssetStatusRetired})
}

func NormalizeAssetAssignmentStatus(value string) string {
	return normalizeAssetAllowed(value, []string{AssetAssignmentRequested, AssetAssignmentApproved, AssetAssignmentIssued, AssetAssignmentReturnDue, AssetAssignmentReturned, AssetAssignmentDamaged, AssetAssignmentLost, AssetAssignmentCancelled})
}

func NormalizeAccessCatalogStatus(value string) string {
	return normalizeAssetAllowed(value, []string{AccessCatalogActive, AccessCatalogInactive, AccessCatalogDeprecated})
}

func NormalizeAccessTaskType(value string) string {
	return normalizeAssetAllowed(value, []string{AccessTaskProvision, AccessTaskDeprovision, AccessTaskReview, AccessTaskChange})
}

func NormalizeAccessTaskStatus(value string) string {
	return normalizeAssetAllowed(value, []string{AccessTaskRequested, AccessTaskApproved, AccessTaskProvisioned, AccessTaskRevoked, AccessTaskReviewed, AccessTaskRejected, AccessTaskCancelled, AccessTaskBlocked})
}

func ValidateAssetItem(item *AssetItem) error {
	if item == nil || item.TenantID == uuid.Nil || strings.TrimSpace(item.AssetCode) == "" || strings.TrimSpace(item.AssetName) == "" {
		return ErrInvalidAssetItem
	}
	item.AssetCode = strings.TrimSpace(item.AssetCode)
	item.AssetName = strings.TrimSpace(item.AssetName)
	item.AssetType = defaultAssetString(item.AssetType, "hardware")
	item.Category = defaultAssetString(item.Category, "general")
	item.Status = defaultAssetAllowed(item.Status, AssetStatusAvailable, NormalizeAssetStatus)
	if item.Status == "" {
		return ErrInvalidAssetItem
	}
	return nil
}

func ValidateAccessCatalogItem(item *AccessCatalogItem) error {
	if item == nil || item.TenantID == uuid.Nil || strings.TrimSpace(item.AccessCode) == "" || strings.TrimSpace(item.AccessName) == "" {
		return ErrInvalidAccessItem
	}
	item.AccessCode = strings.TrimSpace(item.AccessCode)
	item.AccessName = strings.TrimSpace(item.AccessName)
	item.AccessType = defaultAssetString(item.AccessType, "software")
	item.ProvisioningMethod = defaultAssetString(item.ProvisioningMethod, "manual")
	item.Status = defaultAssetAllowed(item.Status, AccessCatalogActive, NormalizeAccessCatalogStatus)
	if item.Status == "" {
		return ErrInvalidAccessItem
	}
	return nil
}

func ValidateAssetAssignment(item *AssetAssignment) error {
	if item == nil || item.TenantID == uuid.Nil || item.AssetID == uuid.Nil || item.WorkerProfileID == uuid.Nil {
		return ErrInvalidAssetAssignment
	}
	item.IssueCondition = defaultAssetString(item.IssueCondition, "good")
	item.DamageStatus = defaultAssetAllowed(item.DamageStatus, "none", func(value string) string {
		return normalizeAssetAllowed(value, []string{"none", "minor", "major", "lost", "recovered"})
	})
	item.Status = defaultAssetAllowed(item.Status, AssetAssignmentRequested, NormalizeAssetAssignmentStatus)
	if item.Status == "" || item.DamageStatus == "" || item.RecoveryAmount < 0 {
		return ErrInvalidAssetAssignment
	}
	return nil
}

func ValidateAccessLifecycleTask(item *AccessLifecycleTask) error {
	if item == nil || item.TenantID == uuid.Nil || item.AccessItemID == uuid.Nil || item.WorkerProfileID == uuid.Nil {
		return ErrInvalidAccessTask
	}
	item.TaskType = defaultAssetAllowed(item.TaskType, AccessTaskProvision, NormalizeAccessTaskType)
	item.Status = defaultAssetAllowed(item.Status, AccessTaskRequested, NormalizeAccessTaskStatus)
	if item.TaskType == "" || item.Status == "" {
		return ErrInvalidAccessTask
	}
	return nil
}

func normalizeAssetAllowed(value string, allowed []string) string {
	normalized := normalizeWorkerProfileEnum(value, "")
	for _, item := range allowed {
		if normalized == item {
			return normalized
		}
	}
	return ""
}

func defaultAssetAllowed(value string, fallback string, normalizer func(string) string) string {
	normalized := normalizer(value)
	if normalized != "" {
		return normalized
	}
	return fallback
}

func defaultAssetString(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return normalizeWorkerProfileEnum(value, fallback)
}
