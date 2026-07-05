package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	AttendanceScheduleFixed         = "fixed"
	AttendanceScheduleFlexi         = "flexi"
	AttendanceScheduleDailyRoster   = "daily_roster"
	AttendanceScheduleWeeklyRoster  = "weekly_roster"
	AttendanceScheduleMonthlyRoster = "monthly_roster"

	AttendanceRequestRegularization      = "regularization"
	AttendanceRequestMissedPunch         = "missed_punch"
	AttendanceRequestLateExemption       = "late_exemption"
	AttendanceRequestEarlyExitExemption  = "early_exit_exemption"
	AttendanceRequestWFH                 = "wfh"
	AttendanceRequestRemoteWork          = "remote_work"
	AttendanceRequestHalfDay             = "halfday"
	AttendanceRequestAbsent              = "absent"
	AttendanceRequestOvertime            = "overtime"
	AttendancePolicyApprovalManager      = "manager"
	AttendancePolicyApprovalHR           = "hr"
	AttendancePolicyApprovalManagerHR    = "manager_hr"
	AttendancePolicyApprovalAuto         = "auto"
	AttendanceRuleOutcomeOnTime          = "on_time"
	AttendanceRuleOutcomeLate            = "late"
	AttendanceRuleOutcomeEarlyExit       = "early_exit"
	AttendanceRuleOutcomeShortHours      = "short_hours"
	AttendanceRuleOutcomeHalfDay         = "halfday"
	AttendanceRuleOutcomeAbsent          = "absent"
	AttendanceRuleOutcomeExemptionNeeded = "exemption_needed"
	AttendanceRuleOutcomeMissingCheckout = "missing_checkout"

	AttendanceLocationOffice      = "office"
	AttendanceLocationBranch      = "branch"
	AttendanceLocationWarehouse   = "warehouse"
	AttendanceLocationClientSite  = "client_site"
	AttendanceLocationFieldZone   = "field_zone"
	AttendanceLocationProjectSite = "project_site"
	AttendanceLocationRemote      = "remote"
	AttendanceLocationOther       = "other"

	AttendanceDeviceIntegrationPush       = "push"
	AttendanceDeviceIntegrationPoll       = "poll"
	AttendanceDeviceIntegrationFileImport = "file_import"
	AttendanceDeviceIntegrationAPI        = "api"
	AttendanceDeviceIntegrationEdgeAgent  = "edge_agent"

	AttendanceDeviceDirectionAuto         = "auto"
	AttendanceDeviceDirectionInOut        = "in_out"
	AttendanceDeviceDirectionEntryExit    = "entry_exit"
	AttendanceDeviceDirectionCheckinOnly  = "checkin_only"
	AttendanceDeviceDirectionCheckoutOnly = "checkout_only"

	AttendanceDeviceStatusActive      = "active"
	AttendanceDeviceStatusInactive    = "inactive"
	AttendanceDeviceStatusMaintenance = "maintenance"

	AttendanceCredentialBiometric   = "biometric"
	AttendanceCredentialFingerprint = "fingerprint"
	AttendanceCredentialFace        = "face"
	AttendanceCredentialCard        = "card"
	AttendanceCredentialPIN         = "pin"
	AttendanceCredentialMobile      = "mobile"
	AttendanceCredentialOther       = "other"
)

var (
	ErrInvalidAttendancePolicyID           = errors.New("attendance_policy_id is required")
	ErrInvalidAttendancePolicyName         = errors.New("attendance policy name is required")
	ErrInvalidAttendancePolicyCode         = errors.New("attendance policy code is required")
	ErrInvalidAttendanceSchedule           = errors.New("attendance schedule_type is invalid")
	ErrInvalidAttendanceApproval           = errors.New("attendance approval mode is invalid")
	ErrInvalidAttendancePolicyScope        = errors.New("attendance policy can target tenant, branch, department, or user, not multiple scopes")
	ErrAttendancePolicyNotFound            = errors.New("attendance policy not found")
	ErrInvalidAttendanceRosterID           = errors.New("attendance_roster_id is required")
	ErrInvalidAttendanceRosterDate         = errors.New("attendance roster date is required")
	ErrInvalidAttendanceRosterTime         = errors.New("attendance roster time is invalid")
	ErrAttendanceRosterNotFound            = errors.New("attendance roster not found")
	ErrInvalidAttendanceRequestID          = errors.New("attendance_request_id is required")
	ErrInvalidAttendanceRequestType        = errors.New("attendance request_type is invalid")
	ErrInvalidAttendanceReview             = errors.New("attendance review status is invalid")
	ErrAttendanceRequestNotFound           = errors.New("attendance request not found")
	ErrInvalidAttendanceExceptionWorkflow  = errors.New("attendance exception workflow is invalid")
	ErrAttendanceExceptionWorkflowNotFound = errors.New("attendance exception workflow not found")
	ErrInvalidAttendanceLocationID         = errors.New("attendance_location_id is required")
	ErrInvalidAttendanceLocationConfig     = errors.New("attendance location configuration is invalid")
	ErrAttendanceLocationNotFound          = errors.New("attendance location not found")
	ErrInvalidAttendanceDeviceID           = errors.New("attendance_device_id is required")
	ErrInvalidAttendanceDevice             = errors.New("attendance device is invalid")
	ErrAttendanceDeviceNotFound            = errors.New("attendance device not found")
	ErrInvalidDeviceUserID                 = errors.New("device_user_id is required")
	ErrEmployeeAttendanceDeviceNotFound    = errors.New("employee attendance device mapping not found")
)

type AttendancePolicy struct {
	ID                         uuid.UUID  `json:"id"`
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
	EffectiveFrom              *time.Time `json:"effective_from,omitempty"`
	EffectiveTo                *time.Time `json:"effective_to,omitempty"`
	Inactive                   bool       `json:"inactive"`
	CreatedAt                  time.Time  `json:"created_at"`
	CreatedBy                  *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt                  time.Time  `json:"updated_at"`
	UpdatedBy                  *uuid.UUID `json:"updated_by,omitempty"`
}

type AttendanceRoster struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	UserID       uuid.UUID  `json:"user_id"`
	PolicyID     *uuid.UUID `json:"policy_id,omitempty"`
	Date         time.Time  `json:"date"`
	StartTime    *string    `json:"start_time,omitempty"`
	EndTime      *string    `json:"end_time,omitempty"`
	BreakMinutes int32      `json:"break_minutes"`
	WorkMode     string     `json:"work_mode"`
	LocationType string     `json:"location_type"`
	Remarks      *string    `json:"remarks,omitempty"`
	Inactive     bool       `json:"inactive"`
	CreatedAt    time.Time  `json:"created_at"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty"`
}

type AttendanceRequest struct {
	ID                  uuid.UUID  `json:"id"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	UserID              uuid.UUID  `json:"user_id"`
	Date                time.Time  `json:"date"`
	RequestedType       *string    `json:"requested_type,omitempty"`
	RequestType         string     `json:"request_type"`
	RequestedCheckInAt  *time.Time `json:"requested_checkin_at,omitempty"`
	RequestedCheckOutAt *time.Time `json:"requested_checkout_at,omitempty"`
	RequestedWorkMode   *string    `json:"requested_work_mode,omitempty"`
	PolicyID            *uuid.UUID `json:"policy_id,omitempty"`
	RosterID            *uuid.UUID `json:"roster_id,omitempty"`
	Reason              *string    `json:"reason,omitempty"`
	Status              string     `json:"status"`
	ReviewedBy          *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt          *time.Time `json:"reviewed_at,omitempty"`
	Remarks             *string    `json:"remarks,omitempty"`
	WorkflowID          *uuid.UUID `json:"workflow_id,omitempty"`
	RouteMode           *string    `json:"route_mode,omitempty"`
	EscalationDueAt     *time.Time `json:"escalation_due_at,omitempty"`
	EscalatedAt         *time.Time `json:"escalated_at,omitempty"`
	PayrollBlocking     bool       `json:"payroll_blocking"`
	Inactive            bool       `json:"inactive"`
	CreatedAt           time.Time  `json:"created_at"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt           time.Time  `json:"updated_at"`
	UpdatedBy           *uuid.UUID `json:"updated_by,omitempty"`
}

type AttendanceExceptionWorkflow struct {
	ID                      uuid.UUID  `json:"id"`
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
	Inactive                bool       `json:"inactive"`
	CreatedAt               time.Time  `json:"created_at"`
	CreatedBy               *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt               time.Time  `json:"updated_at"`
	UpdatedBy               *uuid.UUID `json:"updated_by,omitempty"`
}

type AttendanceExceptionEvent struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	AttendanceRequestID uuid.UUID       `json:"attendance_request_id"`
	WorkflowID          *uuid.UUID      `json:"workflow_id,omitempty"`
	Action              string          `json:"action"`
	FromStatus          *string         `json:"from_status,omitempty"`
	ToStatus            *string         `json:"to_status,omitempty"`
	RoutedTo            *string         `json:"routed_to,omitempty"`
	Remarks             *string         `json:"remarks,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
}

type AttendanceLocation struct {
	ID            uuid.UUID  `json:"id"`
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
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type AttendanceLocationAssignment struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	LocationID    uuid.UUID  `json:"location_id"`
	BranchID      *uuid.UUID `json:"branch_id,omitempty"`
	DepartmentID  *uuid.UUID `json:"department_id,omitempty"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	Inactive      bool       `json:"inactive"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     *uuid.UUID `json:"updated_by,omitempty"`
}

type AttendanceDevice struct {
	ID                   uuid.UUID       `json:"id"`
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
	LastSeenAt           *time.Time      `json:"last_seen_at,omitempty"`
	Config               json.RawMessage `json:"config,omitempty"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type EmployeeAttendanceDevice struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	UserID         uuid.UUID  `json:"user_id"`
	DeviceID       uuid.UUID  `json:"device_id"`
	DeviceUserID   string     `json:"device_user_id"`
	CredentialType string     `json:"credential_type"`
	CardNumber     *string    `json:"card_number,omitempty"`
	EffectiveFrom  *time.Time `json:"effective_from,omitempty"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
	Inactive       bool       `json:"inactive"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
}

type RawAttendanceEvent struct {
	ID                      uuid.UUID       `json:"id"`
	TenantID                uuid.UUID       `json:"tenant_id"`
	DeviceID                uuid.UUID       `json:"device_id"`
	EmployeeDeviceMappingID *uuid.UUID      `json:"employee_device_mapping_id,omitempty"`
	AttendanceID            *uuid.UUID      `json:"attendance_id,omitempty"`
	ExternalEventID         *string         `json:"external_event_id,omitempty"`
	DeviceUserID            *string         `json:"device_user_id,omitempty"`
	EventTime               time.Time       `json:"event_time"`
	EventType               *string         `json:"event_type,omitempty"`
	AttendanceType          *string         `json:"attendance_type,omitempty"`
	ImportBatchID           *string         `json:"import_batch_id,omitempty"`
	ProcessingStatus        string          `json:"processing_status"`
	ProcessingError         *string         `json:"processing_error,omitempty"`
	RawPayload              json.RawMessage `json:"raw_payload,omitempty"`
	ProcessedAt             *time.Time      `json:"processed_at,omitempty"`
	Inactive                bool            `json:"inactive"`
	CreatedAt               time.Time       `json:"created_at"`
	CreatedBy               *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt               time.Time       `json:"updated_at"`
	UpdatedBy               *uuid.UUID      `json:"updated_by,omitempty"`
}

type AttendancePolicyInput struct {
	TenantID                   uuid.UUID
	Name                       string
	Code                       string
	Description                *string
	BranchID                   *uuid.UUID
	DepartmentID               *uuid.UUID
	UserID                     *uuid.UUID
	ScheduleType               string
	IsDefault                  bool
	StandardWorkMinutes        int32
	MinHalfDayMinutes          int32
	MinFullDayMinutes          int32
	GraceLateMinutes           int32
	GraceEarlyMinutes          int32
	HalfDayLateAfterMinutes    *int32
	AbsentLateAfterMinutes     *int32
	HalfDayEarlyBeforeMinutes  *int32
	AbsentEarlyBeforeMinutes   *int32
	AllowFlexiHours            bool
	CoreStartTime              *string
	CoreEndTime                *string
	AllowWFH                   bool
	WFHDaysPerWeek             int32
	AllowPermanentRemote       bool
	RequireGeo                 bool
	RequireDevice              bool
	RegularizationWindowDays   int32
	MaxRegularizationsPerMonth int32
	ApprovalMode               string
	EffectiveFrom              *time.Time
	EffectiveTo                *time.Time
}

type AttendanceLocationInput struct {
	TenantID      uuid.UUID
	BranchID      *uuid.UUID
	Code          string
	Name          string
	LocationType  string
	Latitude      *float64
	Longitude     *float64
	RadiusMeters  int32
	Address       *string
	City          *string
	State         *string
	Country       *string
	Pincode       *string
	EffectiveFrom *time.Time
	EffectiveTo   *time.Time
}

type AttendanceLocationAssignmentInput struct {
	TenantID      uuid.UUID
	LocationID    uuid.UUID
	BranchID      *uuid.UUID
	DepartmentID  *uuid.UUID
	UserID        *uuid.UUID
	EffectiveFrom *time.Time
	EffectiveTo   *time.Time
}

type AttendanceDeviceInput struct {
	TenantID             uuid.UUID
	AttendanceLocationID *uuid.UUID
	BranchID             *uuid.UUID
	Code                 string
	Name                 string
	Vendor               *string
	Model                *string
	SerialNumber         *string
	DeviceIdentifier     *string
	IntegrationType      string
	DirectionMode        string
	Timezone             string
	Status               string
	Config               json.RawMessage
}

type AttendanceExceptionWorkflowInput struct {
	TenantID                uuid.UUID
	Code                    string
	Name                    string
	Description             *string
	BranchID                *uuid.UUID
	DepartmentID            *uuid.UUID
	RequestType             string
	RouteMode               string
	MaxRequestsPerMonth     int32
	EscalationHours         int32
	EscalationRouteMode     *string
	BlockPayrollWhenPending bool
	AutoApprove             bool
	IsActive                bool
}

type EmployeeAttendanceDeviceInput struct {
	TenantID       uuid.UUID
	UserID         uuid.UUID
	DeviceID       uuid.UUID
	DeviceUserID   string
	CredentialType string
	CardNumber     *string
	EffectiveFrom  *time.Time
	EffectiveTo    *time.Time
}

func NewAttendancePolicy(input AttendancePolicyInput) (*AttendancePolicy, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidAttendancePolicyName
	}
	code := strings.ToLower(strings.TrimSpace(input.Code))
	if code == "" {
		return nil, ErrInvalidAttendancePolicyCode
	}
	if scopeCount(input.BranchID, input.DepartmentID, input.UserID) > 1 {
		return nil, ErrInvalidAttendancePolicyScope
	}
	scheduleType, err := NormalizeAttendanceScheduleType(input.ScheduleType)
	if err != nil {
		return nil, err
	}
	approvalMode, err := NormalizeAttendanceApprovalMode(input.ApprovalMode)
	if err != nil {
		return nil, err
	}
	if input.StandardWorkMinutes <= 0 {
		input.StandardWorkMinutes = 480
	}
	if input.MinHalfDayMinutes <= 0 {
		input.MinHalfDayMinutes = 240
	}
	if input.MinFullDayMinutes <= 0 {
		input.MinFullDayMinutes = 420
	}
	if input.RegularizationWindowDays < 0 {
		input.RegularizationWindowDays = 0
	}
	if input.MaxRegularizationsPerMonth < 0 {
		input.MaxRegularizationsPerMonth = 0
	}
	now := time.Now().UTC()
	return &AttendancePolicy{TenantID: input.TenantID, Name: name, Code: code, Description: cleanOptional(input.Description), BranchID: cleanUUIDOptional(input.BranchID), DepartmentID: cleanUUIDOptional(input.DepartmentID), UserID: cleanUUIDOptional(input.UserID), ScheduleType: scheduleType, IsDefault: input.IsDefault, StandardWorkMinutes: input.StandardWorkMinutes, MinHalfDayMinutes: input.MinHalfDayMinutes, MinFullDayMinutes: input.MinFullDayMinutes, GraceLateMinutes: input.GraceLateMinutes, GraceEarlyMinutes: input.GraceEarlyMinutes, HalfDayLateAfterMinutes: input.HalfDayLateAfterMinutes, AbsentLateAfterMinutes: input.AbsentLateAfterMinutes, HalfDayEarlyBeforeMinutes: input.HalfDayEarlyBeforeMinutes, AbsentEarlyBeforeMinutes: input.AbsentEarlyBeforeMinutes, AllowFlexiHours: input.AllowFlexiHours, CoreStartTime: cleanOptional(input.CoreStartTime), CoreEndTime: cleanOptional(input.CoreEndTime), AllowWFH: input.AllowWFH, WFHDaysPerWeek: input.WFHDaysPerWeek, AllowPermanentRemote: input.AllowPermanentRemote, RequireGeo: input.RequireGeo, RequireDevice: input.RequireDevice, RegularizationWindowDays: input.RegularizationWindowDays, MaxRegularizationsPerMonth: input.MaxRegularizationsPerMonth, ApprovalMode: approvalMode, EffectiveFrom: datePtr(input.EffectiveFrom), EffectiveTo: datePtr(input.EffectiveTo), CreatedAt: now, UpdatedAt: now}, nil
}

func NewAttendanceLocation(input AttendanceLocationInput) (*AttendanceLocation, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code := strings.ToLower(strings.TrimSpace(input.Code))
	name := strings.TrimSpace(input.Name)
	if code == "" || name == "" {
		return nil, ErrInvalidAttendanceLocationConfig
	}
	locationType := normalizeAttendanceLocationType(input.LocationType)
	if locationType == "" {
		return nil, ErrInvalidAttendanceLocationConfig
	}
	if input.RadiusMeters < 0 {
		return nil, ErrInvalidAttendanceLocationConfig
	}
	if input.Latitude != nil && (*input.Latitude < -90 || *input.Latitude > 90) {
		return nil, ErrInvalidAttendanceLocationConfig
	}
	if input.Longitude != nil && (*input.Longitude < -180 || *input.Longitude > 180) {
		return nil, ErrInvalidAttendanceLocationConfig
	}
	if input.RadiusMeters == 0 {
		input.RadiusMeters = 100
	}
	now := time.Now().UTC()
	return &AttendanceLocation{TenantID: input.TenantID, BranchID: cleanUUIDOptional(input.BranchID), Code: code, Name: name, LocationType: locationType, Latitude: input.Latitude, Longitude: input.Longitude, RadiusMeters: input.RadiusMeters, Address: cleanOptional(input.Address), City: cleanOptional(input.City), State: cleanOptional(input.State), Country: cleanOptional(input.Country), Pincode: cleanOptional(input.Pincode), EffectiveFrom: datePtr(input.EffectiveFrom), EffectiveTo: datePtr(input.EffectiveTo), CreatedAt: now, UpdatedAt: now}, nil
}

func NewAttendanceLocationAssignment(input AttendanceLocationAssignmentInput) (*AttendanceLocationAssignment, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.LocationID == uuid.Nil {
		return nil, ErrInvalidAttendanceLocationID
	}
	if scopeCount(input.BranchID, input.DepartmentID, input.UserID) > 1 {
		return nil, ErrInvalidAttendancePolicyScope
	}
	now := time.Now().UTC()
	return &AttendanceLocationAssignment{TenantID: input.TenantID, LocationID: input.LocationID, BranchID: cleanUUIDOptional(input.BranchID), DepartmentID: cleanUUIDOptional(input.DepartmentID), UserID: cleanUUIDOptional(input.UserID), EffectiveFrom: datePtr(input.EffectiveFrom), EffectiveTo: datePtr(input.EffectiveTo), CreatedAt: now, UpdatedAt: now}, nil
}

func NewAttendanceDevice(input AttendanceDeviceInput) (*AttendanceDevice, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code := strings.ToLower(strings.TrimSpace(input.Code))
	name := strings.TrimSpace(input.Name)
	if code == "" || name == "" {
		return nil, ErrInvalidAttendanceDevice
	}
	integrationType := normalizeAttendanceDeviceIntegration(input.IntegrationType)
	directionMode := normalizeAttendanceDeviceDirection(input.DirectionMode)
	status := normalizeAttendanceDeviceStatus(input.Status)
	if integrationType == "" || directionMode == "" || status == "" {
		return nil, ErrInvalidAttendanceDevice
	}
	timezone := strings.TrimSpace(input.Timezone)
	if timezone == "" {
		timezone = "UTC"
	}
	config := input.Config
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}
	now := time.Now().UTC()
	return &AttendanceDevice{TenantID: input.TenantID, AttendanceLocationID: cleanUUIDOptional(input.AttendanceLocationID), BranchID: cleanUUIDOptional(input.BranchID), Code: code, Name: name, Vendor: cleanOptional(input.Vendor), Model: cleanOptional(input.Model), SerialNumber: cleanOptional(input.SerialNumber), DeviceIdentifier: cleanOptional(input.DeviceIdentifier), IntegrationType: integrationType, DirectionMode: directionMode, Timezone: timezone, Status: status, Config: config, CreatedAt: now, UpdatedAt: now}, nil
}

func NewEmployeeAttendanceDevice(input EmployeeAttendanceDeviceInput) (*EmployeeAttendanceDevice, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.UserID == uuid.Nil {
		return nil, ErrInvalidEmployeeUserID
	}
	if input.DeviceID == uuid.Nil {
		return nil, ErrInvalidAttendanceDeviceID
	}
	deviceUserID := strings.TrimSpace(input.DeviceUserID)
	if deviceUserID == "" {
		return nil, ErrInvalidDeviceUserID
	}
	credentialType := normalizeAttendanceCredentialType(input.CredentialType)
	if credentialType == "" {
		return nil, ErrInvalidAttendanceDevice
	}
	now := time.Now().UTC()
	return &EmployeeAttendanceDevice{TenantID: input.TenantID, UserID: input.UserID, DeviceID: input.DeviceID, DeviceUserID: deviceUserID, CredentialType: credentialType, CardNumber: cleanOptional(input.CardNumber), EffectiveFrom: datePtr(input.EffectiveFrom), EffectiveTo: datePtr(input.EffectiveTo), CreatedAt: now, UpdatedAt: now}, nil
}

func NewAttendanceExceptionWorkflow(input AttendanceExceptionWorkflowInput) (*AttendanceExceptionWorkflow, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	code := strings.ToLower(strings.TrimSpace(input.Code))
	name := strings.TrimSpace(input.Name)
	if code == "" || name == "" {
		return nil, ErrInvalidAttendanceExceptionWorkflow
	}
	if scopeCount(input.BranchID, input.DepartmentID) > 1 {
		return nil, ErrInvalidAttendancePolicyScope
	}
	requestType, err := NormalizeAttendanceRequestType(input.RequestType)
	if err != nil {
		return nil, err
	}
	routeMode, err := NormalizeAttendanceApprovalMode(input.RouteMode)
	if err != nil {
		return nil, err
	}
	if input.EscalationRouteMode != nil {
		escalationRoute, err := NormalizeAttendanceApprovalMode(*input.EscalationRouteMode)
		if err != nil {
			return nil, err
		}
		input.EscalationRouteMode = &escalationRoute
	}
	if input.MaxRequestsPerMonth < 0 || input.EscalationHours < 0 {
		return nil, ErrInvalidAttendanceExceptionWorkflow
	}
	now := time.Now().UTC()
	return &AttendanceExceptionWorkflow{TenantID: input.TenantID, Code: code, Name: name, Description: cleanOptional(input.Description), BranchID: cleanUUIDOptional(input.BranchID), DepartmentID: cleanUUIDOptional(input.DepartmentID), RequestType: requestType, RouteMode: routeMode, MaxRequestsPerMonth: input.MaxRequestsPerMonth, EscalationHours: input.EscalationHours, EscalationRouteMode: cleanOptional(input.EscalationRouteMode), BlockPayrollWhenPending: input.BlockPayrollWhenPending, AutoApprove: input.AutoApprove, IsActive: input.IsActive, CreatedAt: now, UpdatedAt: now}, nil
}

func normalizeAttendanceLocationType(value string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = AttendanceLocationOffice
	}
	switch clean {
	case AttendanceLocationOffice, AttendanceLocationBranch, AttendanceLocationWarehouse, AttendanceLocationClientSite, AttendanceLocationFieldZone, AttendanceLocationProjectSite, AttendanceLocationRemote, AttendanceLocationOther:
		return clean
	}
	return ""
}

func normalizeAttendanceDeviceIntegration(value string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = AttendanceDeviceIntegrationEdgeAgent
	}
	switch clean {
	case AttendanceDeviceIntegrationPush, AttendanceDeviceIntegrationPoll, AttendanceDeviceIntegrationFileImport, AttendanceDeviceIntegrationAPI, AttendanceDeviceIntegrationEdgeAgent:
		return clean
	}
	return ""
}

func normalizeAttendanceDeviceDirection(value string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = AttendanceDeviceDirectionAuto
	}
	switch clean {
	case AttendanceDeviceDirectionAuto, AttendanceDeviceDirectionInOut, AttendanceDeviceDirectionEntryExit, AttendanceDeviceDirectionCheckinOnly, AttendanceDeviceDirectionCheckoutOnly:
		return clean
	}
	return ""
}

func normalizeAttendanceDeviceStatus(value string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = AttendanceDeviceStatusActive
	}
	switch clean {
	case AttendanceDeviceStatusActive, AttendanceDeviceStatusInactive, AttendanceDeviceStatusMaintenance:
		return clean
	}
	return ""
}

func normalizeAttendanceCredentialType(value string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = AttendanceCredentialBiometric
	}
	switch clean {
	case AttendanceCredentialBiometric, AttendanceCredentialFingerprint, AttendanceCredentialFace, AttendanceCredentialCard, AttendanceCredentialPIN, AttendanceCredentialMobile, AttendanceCredentialOther:
		return clean
	}
	return ""
}

func NormalizeAttendanceScheduleType(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = AttendanceScheduleFixed
	}
	switch clean {
	case AttendanceScheduleFixed, AttendanceScheduleFlexi, AttendanceScheduleDailyRoster, AttendanceScheduleWeeklyRoster, AttendanceScheduleMonthlyRoster:
		return clean, nil
	}
	return "", ErrInvalidAttendanceSchedule
}

func NormalizeAttendanceApprovalMode(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = AttendancePolicyApprovalManager
	}
	switch clean {
	case AttendancePolicyApprovalManager, AttendancePolicyApprovalHR, AttendancePolicyApprovalManagerHR, AttendancePolicyApprovalAuto:
		return clean, nil
	}
	return "", ErrInvalidAttendanceApproval
}

func NormalizeAttendanceRequestType(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = AttendanceRequestRegularization
	}
	switch clean {
	case AttendanceRequestRegularization, AttendanceRequestMissedPunch, AttendanceRequestLateExemption, AttendanceRequestEarlyExitExemption, AttendanceRequestWFH, AttendanceRequestRemoteWork, AttendanceRequestHalfDay, AttendanceRequestAbsent, AttendanceRequestOvertime:
		return clean, nil
	}
	return "", ErrInvalidAttendanceRequestType
}

func NewAttendanceRoster(tenantID, userID uuid.UUID, policyID *uuid.UUID, date time.Time, startTime, endTime *string, breakMinutes int32, workMode, locationType string, remarks *string) (*AttendanceRoster, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidLeaveUser
	}
	if date.IsZero() {
		return nil, ErrInvalidAttendanceRosterDate
	}
	if breakMinutes < 0 {
		return nil, ErrInvalidAttendanceRosterTime
	}
	mode := strings.ToLower(strings.TrimSpace(workMode))
	if mode == "" {
		mode = AttendanceWorkModeOffice
	}
	if !validAttendanceWorkMode(mode) {
		return nil, ErrInvalidAttendanceWorkMode
	}
	loc := strings.ToLower(strings.TrimSpace(locationType))
	if loc == "" {
		loc = "office"
	}
	switch loc {
	case "office", "remote", "field", "client_site", "hybrid":
	default:
		return nil, ErrInvalidAttendanceWorkMode
	}
	now := time.Now().UTC()
	return &AttendanceRoster{TenantID: tenantID, UserID: userID, PolicyID: cleanUUIDOptional(policyID), Date: dateOnly(date), StartTime: cleanOptional(startTime), EndTime: cleanOptional(endTime), BreakMinutes: breakMinutes, WorkMode: mode, LocationType: loc, Remarks: cleanOptional(remarks), CreatedAt: now, UpdatedAt: now}, nil
}

func NewAttendanceRequest(tenantID, userID uuid.UUID, date time.Time, requestType string, requestedType *string, checkIn, checkOut *time.Time, workMode *string, policyID, rosterID *uuid.UUID, reason *string) (*AttendanceRequest, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidLeaveUser
	}
	if date.IsZero() {
		return nil, ErrInvalidAttendanceDate
	}
	requestType, err := NormalizeAttendanceRequestType(requestType)
	if err != nil {
		return nil, err
	}
	if workMode != nil {
		mode := strings.ToLower(strings.TrimSpace(*workMode))
		if !validAttendanceWorkMode(mode) {
			return nil, ErrInvalidAttendanceWorkMode
		}
		workMode = &mode
	}
	now := time.Now().UTC()
	return &AttendanceRequest{TenantID: tenantID, UserID: userID, Date: dateOnly(date), RequestedType: cleanOptional(requestedType), RequestType: requestType, RequestedCheckInAt: timeOptionalUTC(checkIn), RequestedCheckOutAt: timeOptionalUTC(checkOut), RequestedWorkMode: cleanOptional(workMode), PolicyID: cleanUUIDOptional(policyID), RosterID: cleanUUIDOptional(rosterID), Reason: cleanOptional(reason), Status: LeaveStatusPending, CreatedAt: now, UpdatedAt: now}, nil
}

func timeOptionalUTC(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	clean := value.UTC()
	return &clean
}

func scopeCount(values ...*uuid.UUID) int {
	count := 0
	for _, value := range values {
		if cleanUUIDOptional(value) != nil {
			count++
		}
	}
	return count
}
