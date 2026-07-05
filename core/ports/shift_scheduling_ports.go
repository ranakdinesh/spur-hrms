package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type ShiftSchedulingRepo interface {
	CreateShiftTemplate(ctx context.Context, item *domain.ShiftTemplate, actorID *uuid.UUID) (*domain.ShiftTemplate, error)
	UpdateShiftTemplate(ctx context.Context, item *domain.ShiftTemplate, actorID *uuid.UUID) (*domain.ShiftTemplate, error)
	ListShiftTemplates(ctx context.Context, tenantID uuid.UUID, activeOnly *bool, search *string, limit int32, offset int32) ([]*domain.ShiftTemplate, error)
	GetShiftTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ShiftTemplate, error)
	DeleteShiftTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateStaffingRequirement(ctx context.Context, item *domain.StaffingRequirement, actorID *uuid.UUID) (*domain.StaffingRequirement, error)
	UpdateStaffingRequirement(ctx context.Context, item *domain.StaffingRequirement, actorID *uuid.UUID) (*domain.StaffingRequirement, error)
	ListStaffingRequirements(ctx context.Context, filter domain.StaffingRequirementFilter) ([]*domain.StaffingRequirement, error)
	GetStaffingRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.StaffingRequirement, error)
	DeleteStaffingRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateShiftScheduleAssignment(ctx context.Context, item *domain.ShiftScheduleAssignment, actorID *uuid.UUID) (*domain.ShiftScheduleAssignment, error)
	UpdateShiftScheduleAssignment(ctx context.Context, item *domain.ShiftScheduleAssignment, actorID *uuid.UUID) (*domain.ShiftScheduleAssignment, error)
	UpdateShiftScheduleAssignmentStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, hasConflict bool, conflictReason *string, payrollBlocking bool, actorID *uuid.UUID) (*domain.ShiftScheduleAssignment, error)
	ListShiftScheduleAssignments(ctx context.Context, filter domain.ShiftScheduleFilter) ([]*domain.ShiftScheduleAssignment, error)
	GetShiftScheduleAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ShiftScheduleAssignment, error)
	ListShiftAssignmentsForWorkerDate(ctx context.Context, tenantID uuid.UUID, date string, workerProfileID *uuid.UUID, employeeUserID *uuid.UUID) ([]*domain.ShiftScheduleAssignment, error)
	DeleteShiftScheduleAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateShiftSwapRequest(ctx context.Context, item *domain.ShiftSwapRequest, actorID *uuid.UUID) (*domain.ShiftSwapRequest, error)
	UpdateShiftSwapRequestReview(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, reviewedBy uuid.UUID, remarks *string, payrollBlocking bool) (*domain.ShiftSwapRequest, error)
	ListShiftSwapRequests(ctx context.Context, filter domain.ShiftSwapFilter) ([]*domain.ShiftSwapRequest, error)
	GetShiftSwapRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ShiftSwapRequest, error)
	CreateShiftScheduleEvent(ctx context.Context, item *domain.ShiftScheduleEvent, actorID *uuid.UUID) (*domain.ShiftScheduleEvent, error)
	ListShiftScheduleEvents(ctx context.Context, filter domain.ShiftScheduleEventFilter) ([]*domain.ShiftScheduleEvent, error)
	GetShiftScheduleSummary(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.ShiftScheduleSummaryRow, error)
	ListShiftStaffingGaps(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.ShiftStaffingGapRow, error)
}

type ShiftTemplateCommand struct {
	ID                   uuid.UUID       `json:"id,omitempty"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	Code                 string          `json:"code"`
	Name                 string          `json:"name"`
	Description          *string         `json:"description,omitempty"`
	StartTime            string          `json:"start_time"`
	EndTime              string          `json:"end_time"`
	BreakMinutes         int32           `json:"break_minutes"`
	PaidMinutes          int32           `json:"paid_minutes"`
	WorkMode             string          `json:"work_mode"`
	LocationType         string          `json:"location_type"`
	AttendancePolicyID   *uuid.UUID      `json:"attendance_policy_id,omitempty"`
	AttendanceLocationID *uuid.UUID      `json:"attendance_location_id,omitempty"`
	AllowOvertime        bool            `json:"allow_overtime"`
	PayrollCode          *string         `json:"payroll_code,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	IsActive             bool            `json:"is_active"`
	ActorID              *uuid.UUID      `json:"-"`
}

type StaffingRequirementCommand struct {
	ID                   uuid.UUID       `json:"id,omitempty"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	Name                 string          `json:"name"`
	RequirementDate      string          `json:"requirement_date,omitempty"`
	StartDate            string          `json:"start_date,omitempty"`
	EndDate              string          `json:"end_date,omitempty"`
	DayOfWeek            *int32          `json:"day_of_week,omitempty"`
	BranchID             *uuid.UUID      `json:"branch_id,omitempty"`
	DepartmentID         *uuid.UUID      `json:"department_id,omitempty"`
	AttendanceLocationID *uuid.UUID      `json:"attendance_location_id,omitempty"`
	RoleLabel            *string         `json:"role_label,omitempty"`
	TeamLabel            *string         `json:"team_label,omitempty"`
	ShiftTemplateID      *uuid.UUID      `json:"shift_template_id,omitempty"`
	RequiredCount        int32           `json:"required_count"`
	MinCount             int32           `json:"min_count"`
	MaxCount             *int32          `json:"max_count,omitempty"`
	Priority             string          `json:"priority"`
	Status               string          `json:"status"`
	PayrollBlocking      bool            `json:"payroll_blocking"`
	Notes                *string         `json:"notes,omitempty"`
	Metadata             json.RawMessage `json:"metadata,omitempty"`
	ActorID              *uuid.UUID      `json:"-"`
}

type ShiftScheduleAssignmentCommand struct {
	ID                     uuid.UUID       `json:"id,omitempty"`
	TenantID               uuid.UUID       `json:"tenant_id"`
	ScheduleDate           string          `json:"schedule_date"`
	WorkerProfileID        *uuid.UUID      `json:"worker_profile_id,omitempty"`
	EmployeeUserID         *uuid.UUID      `json:"employee_user_id,omitempty"`
	ShiftTemplateID        *uuid.UUID      `json:"shift_template_id,omitempty"`
	AttendancePolicyID     *uuid.UUID      `json:"attendance_policy_id,omitempty"`
	AttendanceLocationID   *uuid.UUID      `json:"attendance_location_id,omitempty"`
	AttendanceRosterID     *uuid.UUID      `json:"attendance_roster_id,omitempty"`
	BranchID               *uuid.UUID      `json:"branch_id,omitempty"`
	DepartmentID           *uuid.UUID      `json:"department_id,omitempty"`
	StartTime              string          `json:"start_time"`
	EndTime                string          `json:"end_time"`
	BreakMinutes           int32           `json:"break_minutes"`
	WorkMode               string          `json:"work_mode"`
	LocationType           string          `json:"location_type"`
	RoleLabel              *string         `json:"role_label,omitempty"`
	TeamLabel              *string         `json:"team_label,omitempty"`
	Status                 string          `json:"status"`
	Source                 string          `json:"source"`
	OvertimePlannedMinutes int32           `json:"overtime_planned_minutes"`
	Notes                  *string         `json:"notes,omitempty"`
	Metadata               json.RawMessage `json:"metadata,omitempty"`
	ActorID                *uuid.UUID      `json:"-"`
}

type ShiftScheduleStatusCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ID       uuid.UUID  `json:"id,omitempty"`
	Status   string     `json:"status"`
	Remarks  *string    `json:"remarks,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type ShiftSwapRequestCommand struct {
	ID                       uuid.UUID       `json:"id,omitempty"`
	TenantID                 uuid.UUID       `json:"tenant_id"`
	RequesterAssignmentID    uuid.UUID       `json:"requester_assignment_id"`
	RequesterWorkerProfileID *uuid.UUID      `json:"requester_worker_profile_id,omitempty"`
	RequesterUserID          *uuid.UUID      `json:"requester_user_id,omitempty"`
	TargetWorkerProfileID    *uuid.UUID      `json:"target_worker_profile_id,omitempty"`
	TargetUserID             *uuid.UUID      `json:"target_user_id,omitempty"`
	OfferedAssignmentID      *uuid.UUID      `json:"offered_assignment_id,omitempty"`
	RequestedDate            string          `json:"requested_date,omitempty"`
	RequestedShiftTemplateID *uuid.UUID      `json:"requested_shift_template_id,omitempty"`
	Reason                   *string         `json:"reason,omitempty"`
	Metadata                 json.RawMessage `json:"metadata,omitempty"`
	ActorID                  *uuid.UUID      `json:"-"`
}

type ShiftSwapReviewCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	ID         uuid.UUID  `json:"id,omitempty"`
	Status     string     `json:"status"`
	Remarks    *string    `json:"remarks,omitempty"`
	ReviewerID uuid.UUID  `json:"-"`
	ActorID    *uuid.UUID `json:"-"`
}
