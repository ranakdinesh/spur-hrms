package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type WorkerProfileRepo interface {
	CreateWorkerProfile(ctx context.Context, item *domain.WorkerProfile, actorID *uuid.UUID) (*domain.WorkerProfile, error)
	UpdateWorkerProfile(ctx context.Context, item *domain.WorkerProfile, actorID *uuid.UUID) (*domain.WorkerProfile, error)
	GetWorkerProfile(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerProfile, error)
	GetWorkerProfileByEmployeeID(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID) (*domain.WorkerProfile, error)
	ListWorkerProfiles(ctx context.Context, filter domain.WorkerProfileFilter) ([]*domain.WorkerProfileListItem, error)
	DeleteWorkerProfile(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type WorkerProfileCommand struct {
	ID                 uuid.UUID       `json:"id,omitempty"`
	TenantID           uuid.UUID       `json:"tenant_id"`
	WorkerTypeID       uuid.UUID       `json:"worker_type_id"`
	EmployeeID         *uuid.UUID      `json:"employee_id,omitempty"`
	EmployeeUserID     *uuid.UUID      `json:"employee_user_id,omitempty"`
	WorkerCode         *string         `json:"worker_code,omitempty"`
	DisplayName        string          `json:"display_name"`
	LegalName          *string         `json:"legal_name,omitempty"`
	Email              *string         `json:"email,omitempty"`
	Mobile             *string         `json:"mobile,omitempty"`
	ProfileStatus      string          `json:"profile_status"`
	StartDate          string          `json:"start_date,omitempty"`
	EndDate            string          `json:"end_date,omitempty"`
	BranchID           *uuid.UUID      `json:"branch_id,omitempty"`
	DepartmentID       *uuid.UUID      `json:"department_id,omitempty"`
	ReportingManagerID *uuid.UUID      `json:"reporting_manager_id,omitempty"`
	WorkLocationLabel  *string         `json:"work_location_label,omitempty"`
	SourcePartner      *string         `json:"source_partner,omitempty"`
	ExternalReference  *string         `json:"external_reference,omitempty"`
	ComplianceStatus   string          `json:"compliance_status"`
	PayrollStatus      string          `json:"payroll_status"`
	Notes              *string         `json:"notes,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	ActorID            *uuid.UUID      `json:"-"`
}
