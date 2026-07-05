package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type EmployeeExitRepo interface {
	CreateEmployeeExitRequest(ctx context.Context, item *domain.EmployeeExitRequest, actorID *uuid.UUID) (*domain.EmployeeExitRequest, error)
	ListEmployeeExitRequests(ctx context.Context, filter domain.EmployeeExitFilter) (*domain.EmployeeExitPage, error)
	GetEmployeeExitRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeExitRequest, error)
	GetActiveEmployeeExitRequestByUserID(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.EmployeeExitRequest, error)
	UpdateEmployeeExitRequestStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, approvedRelievingDate *time.Time, remarks *string, actorID *uuid.UUID) (*domain.EmployeeExitRequest, error)
	CreateEmployeeExitTask(ctx context.Context, task *domain.EmployeeExitTask, actorID *uuid.UUID) (*domain.EmployeeExitTask, error)
	ListEmployeeExitTasks(ctx context.Context, tenantID uuid.UUID, exitRequestID uuid.UUID) ([]*domain.EmployeeExitTask, error)
	GetEmployeeExitTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeExitTask, error)
	UpdateEmployeeExitTaskStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.EmployeeExitTask, error)
	CreateEmployeeExitEvent(ctx context.Context, event *domain.EmployeeExitEvent, actorID *uuid.UUID) (*domain.EmployeeExitEvent, error)
	ListEmployeeExitEvents(ctx context.Context, tenantID uuid.UUID, exitRequestID uuid.UUID) ([]*domain.EmployeeExitEvent, error)
}

type CreateEmployeeExitCommand struct {
	TenantID               uuid.UUID  `json:"tenant_id"`
	EmployeeID             uuid.UUID  `json:"employee_id"`
	ExitType               string     `json:"exit_type"`
	Reason                 *string    `json:"reason,omitempty"`
	ResignationDate        string     `json:"resignation_date,omitempty"`
	NoticeStartDate        string     `json:"notice_start_date,omitempty"`
	LastWorkingDate        string     `json:"last_working_date"`
	RequestedRelievingDate string     `json:"requested_relieving_date,omitempty"`
	Notes                  *string    `json:"notes,omitempty"`
	ActorID                *uuid.UUID `json:"-"`
}

type EmployeeExitActionCommand struct {
	TenantID              uuid.UUID  `json:"tenant_id"`
	ExitID                uuid.UUID  `json:"exit_id"`
	ApprovedRelievingDate string     `json:"approved_relieving_date,omitempty"`
	Remarks               *string    `json:"remarks,omitempty"`
	ActorID               *uuid.UUID `json:"-"`
}

type EmployeeExitTaskStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	TaskID   uuid.UUID  `json:"task_id"`
	Status   string     `json:"status"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type EmployeeExitEventCommand struct {
	TenantID      uuid.UUID       `json:"tenant_id"`
	ExitRequestID uuid.UUID       `json:"exit_request_id"`
	ExitTaskID    *uuid.UUID      `json:"exit_task_id,omitempty"`
	Action        string          `json:"action"`
	FromStatus    *string         `json:"from_status,omitempty"`
	ToStatus      *string         `json:"to_status,omitempty"`
	Remarks       *string         `json:"remarks,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	ActorID       *uuid.UUID      `json:"-"`
}
