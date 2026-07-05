package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	TenantOperationCreateTenant         = "create_tenant"
	TenantOperationSuspendTenant        = "suspend_tenant"
	TenantOperationRestoreTenant        = "restore_tenant"
	TenantOperationScheduleDeleteTenant = "schedule_delete_tenant"
	TenantOperationCancelDeleteTenant   = "cancel_delete_tenant"
	TenantOperationModuleEnable         = "module_enable"
	TenantOperationModuleDisable        = "module_disable"
	TenantOperationStorageChange        = "storage_change"
	TenantOperationDomainBrandingChange = "domain_branding_change"
	TenantOperationAdminReassignment    = "admin_reassignment"
	TenantOperationDataExport           = "data_export"

	TenantOperationPendingValidation = "pending_validation"
	TenantOperationPendingApproval   = "pending_approval"
	TenantOperationApproved          = "approved"
	TenantOperationInProgress        = "in_progress"
	TenantOperationCompleted         = "completed"
	TenantOperationRejected          = "rejected"
	TenantOperationCancelled         = "cancelled"
	TenantOperationFailed            = "failed"

	TenantOperationActionCreated  = "created"
	TenantOperationActionValidate = "validate"
	TenantOperationActionApprove  = "approve"
	TenantOperationActionReject   = "reject"
	TenantOperationActionStart    = "start"
	TenantOperationActionComplete = "complete"
	TenantOperationActionFail     = "fail"
	TenantOperationActionCancel   = "cancel"
)

var (
	ErrInvalidTenantOperation       = errors.New("tenant operation request is invalid")
	ErrTenantOperationNotFound      = errors.New("tenant operation request not found")
	ErrInvalidTenantOperationAction = errors.New("tenant operation action is invalid")
)

type TenantOperationRequest struct {
	ID                uuid.UUID       `json:"id"`
	OperationNumber   string          `json:"operation_number"`
	OperationType     string          `json:"operation_type"`
	Title             string          `json:"title"`
	TargetTenantID    *uuid.UUID      `json:"target_tenant_id,omitempty"`
	TargetTenantName  *string         `json:"target_tenant_name,omitempty"`
	TargetTenantCode  *string         `json:"target_tenant_code,omitempty"`
	Status            string          `json:"status"`
	RiskLevel         string          `json:"risk_level"`
	Reason            string          `json:"reason"`
	RequestedBy       *uuid.UUID      `json:"requested_by,omitempty"`
	ApprovedBy        *uuid.UUID      `json:"approved_by,omitempty"`
	ApprovedAt        *time.Time      `json:"approved_at,omitempty"`
	CompletedBy       *uuid.UUID      `json:"completed_by,omitempty"`
	CompletedAt       *time.Time      `json:"completed_at,omitempty"`
	ApprovalRequired  bool            `json:"approval_required"`
	BackupRequired    bool            `json:"backup_required"`
	BackupConfirmed   bool            `json:"backup_confirmed"`
	RetentionUntil    *time.Time      `json:"retention_until,omitempty"`
	RequestPayload    json.RawMessage `json:"request_payload"`
	ValidationResults json.RawMessage `json:"validation_results"`
	RollbackMetadata  json.RawMessage `json:"rollback_metadata"`
	Metadata          json.RawMessage `json:"metadata"`
	Inactive          bool            `json:"inactive"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at"`
	UpdatedBy         *uuid.UUID      `json:"updated_by,omitempty"`
}

type TenantOperationEvent struct {
	ID          uuid.UUID       `json:"id"`
	RequestID   uuid.UUID       `json:"request_id"`
	Action      string          `json:"action"`
	FromStatus  *string         `json:"from_status,omitempty"`
	ToStatus    *string         `json:"to_status,omitempty"`
	ActorUserID *uuid.UUID      `json:"actor_user_id,omitempty"`
	Remarks     *string         `json:"remarks,omitempty"`
	Metadata    json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at"`
	UpdatedBy   *uuid.UUID      `json:"updated_by,omitempty"`
}

type TenantOperationFilter struct {
	Status         *string
	OperationType  *string
	RiskLevel      *string
	TargetTenantID *uuid.UUID
	Search         *string
	Limit          int32
	Offset         int32
}

type TenantOperationWorkspace struct {
	Requests []*TenantOperationRequest `json:"requests"`
	Summary  TenantOperationSummary    `json:"summary"`
}

type TenantOperationSummary struct {
	Total           int32            `json:"total"`
	PendingApproval int32            `json:"pending_approval"`
	InProgress      int32            `json:"in_progress"`
	HighRisk        int32            `json:"high_risk"`
	Completed       int32            `json:"completed"`
	ByStatus        map[string]int32 `json:"by_status"`
	ByOperationType map[string]int32 `json:"by_operation_type"`
}

type TenantOperationDetail struct {
	Request *TenantOperationRequest `json:"request"`
	Events  []*TenantOperationEvent `json:"events"`
}

func NormalizeTenantOperationRequest(item *TenantOperationRequest) error {
	if item == nil || strings.TrimSpace(item.OperationType) == "" || strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Reason) == "" {
		return ErrInvalidTenantOperation
	}
	item.OperationType = strings.TrimSpace(item.OperationType)
	item.Title = strings.TrimSpace(item.Title)
	item.Reason = strings.TrimSpace(item.Reason)
	item.Status = normalizeTenantOperationStatus(item.Status)
	if !isTenantOperationType(item.OperationType) {
		return ErrInvalidTenantOperation
	}
	if strings.TrimSpace(item.RiskLevel) == "" {
		item.RiskLevel = TenantOperationDefaultRisk(item.OperationType)
	} else {
		item.RiskLevel = normalizeTenantOperationRisk(item.RiskLevel)
	}
	item.ApprovalRequired = true
	item.BackupRequired = item.BackupRequired || TenantOperationRequiresBackup(item.OperationType)
	item.RequestPayload = ensureJSONObject(item.RequestPayload)
	item.ValidationResults = ensureJSONObject(item.ValidationResults)
	item.RollbackMetadata = ensureJSONObject(item.RollbackMetadata)
	item.Metadata = ensureJSONObject(item.Metadata)
	if tenantOperationNeedsTarget(item.OperationType) && (item.TargetTenantID == nil || *item.TargetTenantID == uuid.Nil) {
		return ErrInvalidTenantOperation
	}
	return nil
}

func TenantOperationIsTerminal(status string) bool {
	switch strings.TrimSpace(status) {
	case TenantOperationCompleted, TenantOperationRejected, TenantOperationCancelled, TenantOperationFailed:
		return true
	default:
		return false
	}
}

func isTenantOperationType(operationType string) bool {
	switch operationType {
	case TenantOperationCreateTenant, TenantOperationSuspendTenant, TenantOperationRestoreTenant, TenantOperationScheduleDeleteTenant, TenantOperationCancelDeleteTenant, TenantOperationModuleEnable, TenantOperationModuleDisable, TenantOperationStorageChange, TenantOperationDomainBrandingChange, TenantOperationAdminReassignment, TenantOperationDataExport:
		return true
	default:
		return false
	}
}

func TenantOperationDefaultRisk(operationType string) string {
	switch operationType {
	case TenantOperationScheduleDeleteTenant, TenantOperationStorageChange, TenantOperationAdminReassignment:
		return WorkflowSeverityCritical
	case TenantOperationSuspendTenant, TenantOperationModuleDisable, TenantOperationDomainBrandingChange:
		return WorkflowSeverityHigh
	case TenantOperationCreateTenant, TenantOperationDataExport:
		return WorkflowSeverityMedium
	default:
		return WorkflowSeverityLow
	}
}

func TenantOperationRequiresBackup(operationType string) bool {
	switch operationType {
	case TenantOperationScheduleDeleteTenant, TenantOperationSuspendTenant, TenantOperationStorageChange, TenantOperationDataExport:
		return true
	default:
		return false
	}
}

func tenantOperationNeedsTarget(operationType string) bool {
	return operationType != TenantOperationCreateTenant
}

func normalizeTenantOperationStatus(status string) string {
	switch strings.TrimSpace(status) {
	case TenantOperationPendingValidation, TenantOperationPendingApproval, TenantOperationApproved, TenantOperationInProgress, TenantOperationCompleted, TenantOperationRejected, TenantOperationCancelled, TenantOperationFailed:
		return strings.TrimSpace(status)
	default:
		return TenantOperationPendingValidation
	}
}

func normalizeTenantOperationRisk(risk string) string {
	switch strings.TrimSpace(risk) {
	case WorkflowSeverityLow, WorkflowSeverityMedium, WorkflowSeverityHigh, WorkflowSeverityCritical:
		return strings.TrimSpace(risk)
	default:
		return WorkflowSeverityMedium
	}
}
