package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	OperationLaunchNavigate     = "navigate"
	OperationLaunchWorkflowTask = "workflow_task"
)

type OperationCatalogEntry struct {
	Key                    string          `json:"key"`
	Category               string          `json:"category"`
	Label                  string          `json:"label"`
	Description            string          `json:"description"`
	RequiredPermissions    []string        `json:"required_permissions"`
	TargetModule           string          `json:"target_module"`
	TargetSection          string          `json:"target_section"`
	WorkflowTemplateKey    string          `json:"workflow_template_key,omitempty"`
	WorkflowTemplateID     *uuid.UUID      `json:"workflow_template_id,omitempty"`
	LaunchMode             string          `json:"launch_mode"`
	RequiredApprovalPolicy string          `json:"required_approval_policy"`
	SourceServiceCommand   string          `json:"source_service_command"`
	SourceModule           string          `json:"source_module"`
	SourceType             string          `json:"source_type"`
	DefaultTitle           string          `json:"default_title"`
	DefaultPriority        int32           `json:"default_priority"`
	DefaultSeverity        string          `json:"default_severity"`
	AssigneeRole           string          `json:"assignee_role,omitempty"`
	AssigneeTeam           string          `json:"assignee_team,omitempty"`
	LaunchSchema           json.RawMessage `json:"launch_schema,omitempty"`
	Metadata               json.RawMessage `json:"metadata,omitempty"`
	SortOrder              int32           `json:"sort_order"`
	MobileEnabled          bool            `json:"mobile_enabled"`
}

type OperationCatalog struct {
	TenantID *uuid.UUID               `json:"tenant_id,omitempty"`
	Entries  []*OperationCatalogEntry `json:"entries"`
	Groups   []*OperationCatalogGroup `json:"groups"`
}

type OperationCatalogGroup struct {
	Category string `json:"category"`
	Count    int32  `json:"count"`
}
