package ports

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type AttendanceRepo interface {
	CreateAttendance(ctx context.Context, item *domain.Attendance, actorID *uuid.UUID) (*domain.Attendance, error)
	UpdateAttendance(ctx context.Context, item *domain.Attendance, actorID *uuid.UUID) (*domain.Attendance, error)
	ListAttendancesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.Attendance, error)
	ListAttendancesByUserDate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, date string) ([]*domain.Attendance, error)
	ListAttendancesByDate(ctx context.Context, tenantID uuid.UUID, date string) ([]*domain.Attendance, error)
	GetAttendance(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Attendance, error)
	CreateDeviceLog(ctx context.Context, item *domain.DeviceLog, actorID *uuid.UUID) (*domain.DeviceLog, error)
	ListDeviceLogsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.DeviceLog, error)
}

type AttendancePunchCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	UserID     uuid.UUID  `json:"user_id"`
	Action     string     `json:"action"`
	Date       string     `json:"date"`
	Time       string     `json:"time"`
	Source     *string    `json:"source,omitempty"`
	Latitude   *float64   `json:"latitude,omitempty"`
	Longitude  *float64   `json:"longitude,omitempty"`
	WorkMode   *string    `json:"work_mode,omitempty"`
	Remarks    *string    `json:"remarks,omitempty"`
	DeviceID   *string    `json:"device_id,omitempty"`
	DeviceType *string    `json:"device_type,omitempty"`
	IPAddress  *string    `json:"ip_address,omitempty"`
	ActorID    *uuid.UUID `json:"-"`
}

type AttendanceStatusQuery struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	Date     string     `json:"date"`
}

type AttendanceReportQuery struct {
	TenantID     uuid.UUID  `json:"tenant_id"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	BranchID     *uuid.UUID `json:"branch_id,omitempty"`
	StartDate    string     `json:"start_date"`
	EndDate      string     `json:"end_date"`
}

type AttendancePolicyRepo interface {
	CreateAttendancePolicy(ctx context.Context, item *domain.AttendancePolicy, actorID *uuid.UUID) (*domain.AttendancePolicy, error)
	ListAttendancePolicies(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendancePolicy, error)
	GetAttendancePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendancePolicy, error)
	ResolveAttendancePolicy(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, departmentID *uuid.UUID, branchID *uuid.UUID, date string) (*domain.AttendancePolicy, error)
	UpdateAttendancePolicy(ctx context.Context, item *domain.AttendancePolicy, actorID *uuid.UUID) (*domain.AttendancePolicy, error)
	DeleteAttendancePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type AttendanceRosterRepo interface {
	CreateAttendanceRoster(ctx context.Context, item *domain.AttendanceRoster, actorID *uuid.UUID) (*domain.AttendanceRoster, error)
	ListAttendanceRostersByDateRange(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceRoster, error)
	ListAttendanceRostersByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceRoster, error)
	GetAttendanceRosterByUserDate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, date string) (*domain.AttendanceRoster, error)
	UpdateAttendanceRoster(ctx context.Context, item *domain.AttendanceRoster, actorID *uuid.UUID) (*domain.AttendanceRoster, error)
	DeleteAttendanceRoster(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type AttendanceRequestRepo interface {
	CreateAttendanceRequest(ctx context.Context, item *domain.AttendanceRequest, actorID *uuid.UUID) (*domain.AttendanceRequest, error)
	ListAttendanceRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.AttendanceRequest, error)
	ListAttendanceRequestsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.AttendanceRequest, error)
	GetAttendanceRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendanceRequest, error)
	SetAttendanceRequestWorkflow(ctx context.Context, tenantID uuid.UUID, requestID uuid.UUID, workflowID *uuid.UUID, routeMode *string, escalationDueAt *string, payrollBlocking bool, actorID *uuid.UUID) (*domain.AttendanceRequest, error)
	UpdateAttendanceRequestReview(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, reviewerID uuid.UUID, remarks *string) (*domain.AttendanceRequest, error)
	DeleteAttendanceRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type AttendanceExceptionWorkflowRepo interface {
	CreateAttendanceExceptionWorkflow(ctx context.Context, item *domain.AttendanceExceptionWorkflow, actorID *uuid.UUID) (*domain.AttendanceExceptionWorkflow, error)
	ListAttendanceExceptionWorkflows(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendanceExceptionWorkflow, error)
	GetAttendanceExceptionWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendanceExceptionWorkflow, error)
	ResolveAttendanceExceptionWorkflow(ctx context.Context, tenantID uuid.UUID, requestType string, departmentID *uuid.UUID, branchID *uuid.UUID) (*domain.AttendanceExceptionWorkflow, error)
	UpdateAttendanceExceptionWorkflow(ctx context.Context, item *domain.AttendanceExceptionWorkflow, actorID *uuid.UUID) (*domain.AttendanceExceptionWorkflow, error)
	DeleteAttendanceExceptionWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateAttendanceExceptionEvent(ctx context.Context, item *domain.AttendanceExceptionEvent, actorID *uuid.UUID) (*domain.AttendanceExceptionEvent, error)
	ListAttendanceExceptionEvents(ctx context.Context, tenantID uuid.UUID, attendanceRequestID uuid.UUID) ([]*domain.AttendanceExceptionEvent, error)
	ListPayrollBlockingAttendanceRequests(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceRequest, error)
}

type AttendanceLocationRepo interface {
	CreateAttendanceLocation(ctx context.Context, item *domain.AttendanceLocation, actorID *uuid.UUID) (*domain.AttendanceLocation, error)
	ListAttendanceLocations(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendanceLocation, error)
	GetAttendanceLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendanceLocation, error)
	UpdateAttendanceLocation(ctx context.Context, item *domain.AttendanceLocation, actorID *uuid.UUID) (*domain.AttendanceLocation, error)
	DeleteAttendanceLocation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateAttendanceLocationAssignment(ctx context.Context, item *domain.AttendanceLocationAssignment, actorID *uuid.UUID) (*domain.AttendanceLocationAssignment, error)
	ListAttendanceLocationAssignments(ctx context.Context, tenantID uuid.UUID, locationID *uuid.UUID) ([]*domain.AttendanceLocationAssignment, error)
	UpdateAttendanceLocationAssignment(ctx context.Context, item *domain.AttendanceLocationAssignment, actorID *uuid.UUID) (*domain.AttendanceLocationAssignment, error)
	DeleteAttendanceLocationAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type AttendanceDeviceRepo interface {
	CreateAttendanceDevice(ctx context.Context, item *domain.AttendanceDevice, actorID *uuid.UUID) (*domain.AttendanceDevice, error)
	ListAttendanceDevices(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendanceDevice, error)
	GetAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendanceDevice, error)
	UpdateAttendanceDevice(ctx context.Context, item *domain.AttendanceDevice, actorID *uuid.UUID) (*domain.AttendanceDevice, error)
	DeleteAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateEmployeeAttendanceDevice(ctx context.Context, item *domain.EmployeeAttendanceDevice, actorID *uuid.UUID) (*domain.EmployeeAttendanceDevice, error)
	ListEmployeeAttendanceDevices(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID) ([]*domain.EmployeeAttendanceDevice, error)
	GetEmployeeAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeAttendanceDevice, error)
	UpdateEmployeeAttendanceDevice(ctx context.Context, item *domain.EmployeeAttendanceDevice, actorID *uuid.UUID) (*domain.EmployeeAttendanceDevice, error)
	DeleteEmployeeAttendanceDevice(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateRawAttendanceEvent(ctx context.Context, item *domain.RawAttendanceEvent, actorID *uuid.UUID) (*domain.RawAttendanceEvent, error)
	ListRawAttendanceEvents(ctx context.Context, tenantID uuid.UUID, limit int32) ([]*domain.RawAttendanceEvent, error)
}

type AttendancePolicyCommand struct {
	ID                         uuid.UUID  `json:"id,omitempty"`
	TenantID                   uuid.UUID  `json:"tenant_id"`
	Name                       string     `json:"name"`
	Code                       string     `json:"code"`
	Description                *string    `json:"description,omitempty"`
	BranchID                   *uuid.UUID `json:"branch_id,omitempty"`
	DepartmentID               *uuid.UUID `json:"department_id,omitempty"`
	UserID                     *uuid.UUID `json:"user_id,omitempty"`
	ScheduleType               string     `json:"schedule_type"`
	IsDefault                  bool       `json:"is_default"`
	StandardWorkMinutes        int32      `json:"standard_work_minutes"`
	MinHalfDayMinutes          int32      `json:"min_half_day_minutes"`
	MinFullDayMinutes          int32      `json:"min_full_day_minutes"`
	GraceLateMinutes           int32      `json:"grace_late_minutes"`
	GraceEarlyMinutes          int32      `json:"grace_early_minutes"`
	HalfDayLateAfterMinutes    *int32     `json:"half_day_late_after_minutes,omitempty"`
	AbsentLateAfterMinutes     *int32     `json:"absent_late_after_minutes,omitempty"`
	HalfDayEarlyBeforeMinutes  *int32     `json:"half_day_early_before_minutes,omitempty"`
	AbsentEarlyBeforeMinutes   *int32     `json:"absent_early_before_minutes,omitempty"`
	AllowFlexiHours            bool       `json:"allow_flexi_hours"`
	CoreStartTime              *string    `json:"core_start_time,omitempty"`
	CoreEndTime                *string    `json:"core_end_time,omitempty"`
	AllowWFH                   bool       `json:"allow_wfh"`
	WFHDaysPerWeek             int32      `json:"wfh_days_per_week"`
	AllowPermanentRemote       bool       `json:"allow_permanent_remote"`
	RequireGeo                 bool       `json:"require_geo"`
	RequireDevice              bool       `json:"require_device"`
	RegularizationWindowDays   int32      `json:"regularization_window_days"`
	MaxRegularizationsPerMonth int32      `json:"max_regularizations_per_month"`
	ApprovalMode               string     `json:"approval_mode"`
	EffectiveFrom              string     `json:"effective_from,omitempty"`
	EffectiveTo                string     `json:"effective_to,omitempty"`
	ActorID                    *uuid.UUID `json:"-"`
}

type AttendanceRosterCommand struct {
	ID           uuid.UUID  `json:"id,omitempty"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	UserID       uuid.UUID  `json:"user_id"`
	PolicyID     *uuid.UUID `json:"policy_id,omitempty"`
	Date         string     `json:"date"`
	StartTime    *string    `json:"start_time,omitempty"`
	EndTime      *string    `json:"end_time,omitempty"`
	BreakMinutes int32      `json:"break_minutes"`
	WorkMode     string     `json:"work_mode"`
	LocationType string     `json:"location_type"`
	Remarks      *string    `json:"remarks,omitempty"`
	ActorID      *uuid.UUID `json:"-"`
}

type AttendanceRequestCommand struct {
	TenantID            uuid.UUID  `json:"tenant_id"`
	UserID              uuid.UUID  `json:"user_id"`
	Date                string     `json:"date"`
	RequestedType       *string    `json:"requested_type,omitempty"`
	RequestType         string     `json:"request_type"`
	RequestedCheckInAt  *string    `json:"requested_checkin_at,omitempty"`
	RequestedCheckOutAt *string    `json:"requested_checkout_at,omitempty"`
	RequestedWorkMode   *string    `json:"requested_work_mode,omitempty"`
	PolicyID            *uuid.UUID `json:"policy_id,omitempty"`
	RosterID            *uuid.UUID `json:"roster_id,omitempty"`
	Reason              *string    `json:"reason,omitempty"`
	ActorID             *uuid.UUID `json:"-"`
}

type AttendanceExceptionWorkflowCommand struct {
	ID                      uuid.UUID  `json:"id,omitempty"`
	TenantID                uuid.UUID  `json:"tenant_id"`
	Code                    string     `json:"code"`
	Name                    string     `json:"name"`
	Description             *string    `json:"description,omitempty"`
	BranchID                *uuid.UUID `json:"branch_id,omitempty"`
	DepartmentID            *uuid.UUID `json:"department_id,omitempty"`
	RequestType             string     `json:"request_type"`
	RouteMode               string     `json:"route_mode"`
	MaxRequestsPerMonth     int32      `json:"max_requests_per_month"`
	EscalationHours         int32      `json:"escalation_hours"`
	EscalationRouteMode     *string    `json:"escalation_route_mode,omitempty"`
	BlockPayrollWhenPending bool       `json:"block_payroll_when_pending"`
	AutoApprove             bool       `json:"auto_approve"`
	IsActive                bool       `json:"is_active"`
	ActorID                 *uuid.UUID `json:"-"`
}

type AttendanceReviewCommand struct {
	TenantID   uuid.UUID `json:"tenant_id"`
	RequestID  uuid.UUID `json:"request_id"`
	Status     string    `json:"status"`
	Remarks    *string   `json:"remarks,omitempty"`
	ReviewerID uuid.UUID `json:"-"`
}

type AttendanceLocationCommand struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	BranchID      *uuid.UUID `json:"branch_id,omitempty"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	LocationType  string     `json:"location_type"`
	Latitude      *float64   `json:"latitude,omitempty"`
	Longitude     *float64   `json:"longitude,omitempty"`
	RadiusMeters  int32      `json:"radius_meters"`
	Address       *string    `json:"address,omitempty"`
	City          *string    `json:"city,omitempty"`
	State         *string    `json:"state,omitempty"`
	Country       *string    `json:"country,omitempty"`
	Pincode       *string    `json:"pincode,omitempty"`
	EffectiveFrom string     `json:"effective_from,omitempty"`
	EffectiveTo   string     `json:"effective_to,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}

type AttendanceLocationAssignmentCommand struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	LocationID    uuid.UUID  `json:"location_id"`
	BranchID      *uuid.UUID `json:"branch_id,omitempty"`
	DepartmentID  *uuid.UUID `json:"department_id,omitempty"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	EffectiveFrom string     `json:"effective_from,omitempty"`
	EffectiveTo   string     `json:"effective_to,omitempty"`
	ActorID       *uuid.UUID `json:"-"`
}

type AttendanceDeviceCommand struct {
	ID                   uuid.UUID       `json:"id,omitempty"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	AttendanceLocationID *uuid.UUID      `json:"attendance_location_id,omitempty"`
	BranchID             *uuid.UUID      `json:"branch_id,omitempty"`
	Code                 string          `json:"code"`
	Name                 string          `json:"name"`
	Vendor               *string         `json:"vendor,omitempty"`
	Model                *string         `json:"model,omitempty"`
	SerialNumber         *string         `json:"serial_number,omitempty"`
	DeviceIdentifier     *string         `json:"device_identifier,omitempty"`
	IntegrationType      string          `json:"integration_type"`
	DirectionMode        string          `json:"direction_mode"`
	Timezone             string          `json:"timezone"`
	Status               string          `json:"status"`
	Config               json.RawMessage `json:"config,omitempty"`
	ActorID              *uuid.UUID      `json:"-"`
}

type EmployeeAttendanceDeviceCommand struct {
	ID             uuid.UUID  `json:"id,omitempty"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	UserID         uuid.UUID  `json:"user_id"`
	DeviceID       uuid.UUID  `json:"device_id"`
	DeviceUserID   string     `json:"device_user_id"`
	CredentialType string     `json:"credential_type"`
	CardNumber     *string    `json:"card_number,omitempty"`
	EffectiveFrom  string     `json:"effective_from,omitempty"`
	EffectiveTo    string     `json:"effective_to,omitempty"`
	ActorID        *uuid.UUID `json:"-"`
}
