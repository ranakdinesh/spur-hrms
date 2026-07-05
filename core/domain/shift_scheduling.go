package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ShiftScheduleStatusDraft     = "draft"
	ShiftScheduleStatusPublished = "published"
	ShiftScheduleStatusLocked    = "locked"
	ShiftScheduleStatusCancelled = "cancelled"
	ShiftScheduleStatusCompleted = "completed"

	ShiftScheduleSourceManual   = "manual"
	ShiftScheduleSourceTemplate = "template"
	ShiftScheduleSourceImport   = "import"
	ShiftScheduleSourceSwap     = "swap"
	ShiftScheduleSourceSystem   = "system"

	ShiftSwapStatusPending   = "pending"
	ShiftSwapStatusApproved  = "approved"
	ShiftSwapStatusRejected  = "rejected"
	ShiftSwapStatusCancelled = "cancelled"

	StaffingRequirementStatusActive   = "active"
	StaffingRequirementStatusPaused   = "paused"
	StaffingRequirementStatusArchived = "archived"
)

var (
	ErrInvalidShiftTemplate        = errors.New("shift template is invalid")
	ErrShiftTemplateNotFound       = errors.New("shift template not found")
	ErrInvalidStaffingRequirement  = errors.New("staffing requirement is invalid")
	ErrStaffingRequirementNotFound = errors.New("staffing requirement not found")
	ErrInvalidShiftAssignment      = errors.New("shift schedule assignment is invalid")
	ErrShiftAssignmentNotFound     = errors.New("shift schedule assignment not found")
	ErrInvalidShiftSwapRequest     = errors.New("shift swap request is invalid")
	ErrShiftSwapRequestNotFound    = errors.New("shift swap request not found")
	ErrInvalidShiftScheduleEvent   = errors.New("shift schedule event is invalid")
)

type ShiftTemplate struct {
	ID                   uuid.UUID       `json:"id"`
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
	Metadata             json.RawMessage `json:"metadata"`
	IsActive             bool            `json:"is_active"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type StaffingRequirement struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	Name                 string          `json:"name"`
	RequirementDate      *time.Time      `json:"requirement_date,omitempty"`
	StartDate            *time.Time      `json:"start_date,omitempty"`
	EndDate              *time.Time      `json:"end_date,omitempty"`
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
	Metadata             json.RawMessage `json:"metadata"`
	Inactive             bool            `json:"inactive"`
	CreatedAt            time.Time       `json:"created_at"`
	CreatedBy            *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt            time.Time       `json:"updated_at"`
	UpdatedBy            *uuid.UUID      `json:"updated_by,omitempty"`
}

type ShiftScheduleAssignment struct {
	ID                     uuid.UUID       `json:"id"`
	TenantID               uuid.UUID       `json:"tenant_id"`
	ScheduleDate           time.Time       `json:"schedule_date"`
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
	HasConflict            bool            `json:"has_conflict"`
	ConflictReason         *string         `json:"conflict_reason,omitempty"`
	PayrollBlocking        bool            `json:"payroll_blocking"`
	Notes                  *string         `json:"notes,omitempty"`
	Metadata               json.RawMessage `json:"metadata"`
	Inactive               bool            `json:"inactive"`
	CreatedAt              time.Time       `json:"created_at"`
	CreatedBy              *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt              time.Time       `json:"updated_at"`
	UpdatedBy              *uuid.UUID      `json:"updated_by,omitempty"`
}

type ShiftSwapRequest struct {
	ID                       uuid.UUID       `json:"id"`
	TenantID                 uuid.UUID       `json:"tenant_id"`
	RequesterAssignmentID    uuid.UUID       `json:"requester_assignment_id"`
	RequesterWorkerProfileID *uuid.UUID      `json:"requester_worker_profile_id,omitempty"`
	RequesterUserID          *uuid.UUID      `json:"requester_user_id,omitempty"`
	TargetWorkerProfileID    *uuid.UUID      `json:"target_worker_profile_id,omitempty"`
	TargetUserID             *uuid.UUID      `json:"target_user_id,omitempty"`
	OfferedAssignmentID      *uuid.UUID      `json:"offered_assignment_id,omitempty"`
	RequestedDate            *time.Time      `json:"requested_date,omitempty"`
	RequestedShiftTemplateID *uuid.UUID      `json:"requested_shift_template_id,omitempty"`
	Reason                   *string         `json:"reason,omitempty"`
	Status                   string          `json:"status"`
	ReviewedBy               *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewedAt               *time.Time      `json:"reviewed_at,omitempty"`
	ReviewRemarks            *string         `json:"review_remarks,omitempty"`
	PayrollBlocking          bool            `json:"payroll_blocking"`
	Metadata                 json.RawMessage `json:"metadata"`
	Inactive                 bool            `json:"inactive"`
	CreatedAt                time.Time       `json:"created_at"`
	CreatedBy                *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt                time.Time       `json:"updated_at"`
	UpdatedBy                *uuid.UUID      `json:"updated_by,omitempty"`
}

type ShiftScheduleEvent struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	SourceType  string          `json:"source_type"`
	SourceID    uuid.UUID       `json:"source_id"`
	Action      string          `json:"action"`
	FromStatus  *string         `json:"from_status,omitempty"`
	ToStatus    *string         `json:"to_status,omitempty"`
	ActorUserID *uuid.UUID      `json:"actor_user_id,omitempty"`
	Remarks     *string         `json:"remarks,omitempty"`
	Metadata    json.RawMessage `json:"metadata"`
	Inactive    bool            `json:"inactive"`
	CreatedAt   time.Time       `json:"created_at"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at"`
	UpdatedBy   *uuid.UUID      `json:"updated_by,omitempty"`
}

type ShiftScheduleFilter struct {
	TenantID             uuid.UUID
	StartDate            string
	EndDate              string
	Status               *string
	WorkerProfileID      *uuid.UUID
	EmployeeUserID       *uuid.UUID
	BranchID             *uuid.UUID
	DepartmentID         *uuid.UUID
	AttendanceLocationID *uuid.UUID
	Limit                int32
	Offset               int32
}

type StaffingRequirementFilter struct {
	TenantID  uuid.UUID
	Status    *string
	StartDate string
	EndDate   string
	Limit     int32
	Offset    int32
}

type ShiftSwapFilter struct {
	TenantID        uuid.UUID
	Status          *string
	RequesterUserID *uuid.UUID
	TargetUserID    *uuid.UUID
	Limit           int32
	Offset          int32
}

type ShiftScheduleEventFilter struct {
	TenantID   uuid.UUID
	SourceType *string
	SourceID   *uuid.UUID
	Limit      int32
	Offset     int32
}

type ShiftScheduleSummaryRow struct {
	Metric      string `json:"metric"`
	MetricCount int32  `json:"metric_count"`
}

type ShiftStaffingGapRow struct {
	RequirementID        uuid.UUID  `json:"requirement_id"`
	RequirementName      string     `json:"requirement_name"`
	BranchID             *uuid.UUID `json:"branch_id,omitempty"`
	DepartmentID         *uuid.UUID `json:"department_id,omitempty"`
	AttendanceLocationID *uuid.UUID `json:"attendance_location_id,omitempty"`
	ShiftTemplateID      *uuid.UUID `json:"shift_template_id,omitempty"`
	RequiredCount        int32      `json:"required_count"`
	AssignedCount        int32      `json:"assigned_count"`
	GapCount             int32      `json:"gap_count"`
	Priority             string     `json:"priority"`
	PayrollBlocking      bool       `json:"payroll_blocking"`
}

func NewShiftTemplate(item ShiftTemplate) (*ShiftTemplate, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Code) == "" || strings.TrimSpace(item.Name) == "" {
		return nil, ErrInvalidShiftTemplate
	}
	start, end, ok := validateShiftClockRange(item.StartTime, item.EndTime)
	if !ok || item.BreakMinutes < 0 || item.PaidMinutes < 0 {
		return nil, ErrInvalidShiftTemplate
	}
	mode, err := normalizeShiftWorkMode(item.WorkMode)
	if err != nil {
		return nil, err
	}
	locationType, err := normalizeShiftLocationType(item.LocationType)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	item.Code = strings.ToUpper(strings.TrimSpace(item.Code))
	item.Name = strings.TrimSpace(item.Name)
	item.Description = cleanOptional(item.Description)
	item.StartTime = start
	item.EndTime = end
	item.WorkMode = mode
	item.LocationType = locationType
	item.AttendancePolicyID = cleanUUIDOptional(item.AttendancePolicyID)
	item.AttendanceLocationID = cleanUUIDOptional(item.AttendanceLocationID)
	item.PayrollCode = cleanOptional(item.PayrollCode)
	item.Metadata = normalizeShiftJSONRaw(item.Metadata)
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewStaffingRequirement(item StaffingRequirement) (*StaffingRequirement, error) {
	if item.TenantID == uuid.Nil || strings.TrimSpace(item.Name) == "" || item.RequiredCount < 0 || item.MinCount < 0 {
		return nil, ErrInvalidStaffingRequirement
	}
	if item.MaxCount != nil && *item.MaxCount < item.MinCount {
		return nil, ErrInvalidStaffingRequirement
	}
	if item.DayOfWeek != nil && (*item.DayOfWeek < 0 || *item.DayOfWeek > 6) {
		return nil, ErrInvalidStaffingRequirement
	}
	if item.RequirementDate == nil && item.StartDate == nil {
		return nil, ErrInvalidStaffingRequirement
	}
	if item.StartDate != nil && item.EndDate != nil && item.EndDate.Before(*item.StartDate) {
		return nil, ErrInvalidStaffingRequirement
	}
	priority, err := normalizeShiftPriority(item.Priority)
	if err != nil {
		return nil, ErrInvalidStaffingRequirement
	}
	status, err := NormalizeStaffingRequirementStatus(item.Status)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	item.Name = strings.TrimSpace(item.Name)
	item.RequirementDate = datePtrOnly(item.RequirementDate)
	item.StartDate = datePtrOnly(item.StartDate)
	item.EndDate = datePtrOnly(item.EndDate)
	item.BranchID = cleanUUIDOptional(item.BranchID)
	item.DepartmentID = cleanUUIDOptional(item.DepartmentID)
	item.AttendanceLocationID = cleanUUIDOptional(item.AttendanceLocationID)
	item.RoleLabel = cleanOptional(item.RoleLabel)
	item.TeamLabel = cleanOptional(item.TeamLabel)
	item.ShiftTemplateID = cleanUUIDOptional(item.ShiftTemplateID)
	item.Priority = priority
	item.Status = status
	item.Notes = cleanOptional(item.Notes)
	item.Metadata = normalizeShiftJSONRaw(item.Metadata)
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewShiftScheduleAssignment(item ShiftScheduleAssignment) (*ShiftScheduleAssignment, error) {
	if item.TenantID == uuid.Nil || item.ScheduleDate.IsZero() {
		return nil, ErrInvalidShiftAssignment
	}
	item.WorkerProfileID = cleanUUIDOptional(item.WorkerProfileID)
	item.EmployeeUserID = cleanUUIDOptional(item.EmployeeUserID)
	if item.WorkerProfileID == nil && item.EmployeeUserID == nil {
		return nil, ErrInvalidShiftAssignment
	}
	start, end, ok := validateShiftClockRange(item.StartTime, item.EndTime)
	if !ok || item.BreakMinutes < 0 || item.OvertimePlannedMinutes < 0 {
		return nil, ErrInvalidShiftAssignment
	}
	status, err := NormalizeShiftAssignmentStatus(item.Status)
	if err != nil {
		return nil, err
	}
	source, err := NormalizeShiftAssignmentSource(item.Source)
	if err != nil {
		return nil, err
	}
	mode, err := normalizeShiftWorkMode(item.WorkMode)
	if err != nil {
		return nil, err
	}
	locationType, err := normalizeShiftLocationType(item.LocationType)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	item.ScheduleDate = dateOnly(item.ScheduleDate)
	item.ShiftTemplateID = cleanUUIDOptional(item.ShiftTemplateID)
	item.AttendancePolicyID = cleanUUIDOptional(item.AttendancePolicyID)
	item.AttendanceLocationID = cleanUUIDOptional(item.AttendanceLocationID)
	item.AttendanceRosterID = cleanUUIDOptional(item.AttendanceRosterID)
	item.BranchID = cleanUUIDOptional(item.BranchID)
	item.DepartmentID = cleanUUIDOptional(item.DepartmentID)
	item.StartTime = start
	item.EndTime = end
	item.WorkMode = mode
	item.LocationType = locationType
	item.RoleLabel = cleanOptional(item.RoleLabel)
	item.TeamLabel = cleanOptional(item.TeamLabel)
	item.Status = status
	item.Source = source
	item.ConflictReason = cleanOptional(item.ConflictReason)
	item.Notes = cleanOptional(item.Notes)
	item.Metadata = normalizeShiftJSONRaw(item.Metadata)
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewShiftSwapRequest(item ShiftSwapRequest) (*ShiftSwapRequest, error) {
	if item.TenantID == uuid.Nil || item.RequesterAssignmentID == uuid.Nil {
		return nil, ErrInvalidShiftSwapRequest
	}
	item.RequesterWorkerProfileID = cleanUUIDOptional(item.RequesterWorkerProfileID)
	item.RequesterUserID = cleanUUIDOptional(item.RequesterUserID)
	item.TargetWorkerProfileID = cleanUUIDOptional(item.TargetWorkerProfileID)
	item.TargetUserID = cleanUUIDOptional(item.TargetUserID)
	item.OfferedAssignmentID = cleanUUIDOptional(item.OfferedAssignmentID)
	item.RequestedShiftTemplateID = cleanUUIDOptional(item.RequestedShiftTemplateID)
	if item.TargetWorkerProfileID == nil && item.TargetUserID == nil && item.OfferedAssignmentID == nil && item.RequestedShiftTemplateID == nil {
		return nil, ErrInvalidShiftSwapRequest
	}
	status, err := NormalizeShiftSwapStatus(item.Status)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	item.RequestedDate = datePtrOnly(item.RequestedDate)
	item.Reason = cleanOptional(item.Reason)
	item.Status = status
	item.ReviewRemarks = cleanOptional(item.ReviewRemarks)
	item.Metadata = normalizeShiftJSONRaw(item.Metadata)
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NewShiftScheduleEvent(item ShiftScheduleEvent) (*ShiftScheduleEvent, error) {
	if item.TenantID == uuid.Nil || item.SourceID == uuid.Nil || strings.TrimSpace(item.SourceType) == "" || strings.TrimSpace(item.Action) == "" {
		return nil, ErrInvalidShiftScheduleEvent
	}
	sourceType := strings.ToLower(strings.TrimSpace(item.SourceType))
	switch sourceType {
	case "template", "requirement", "assignment", "swap_request":
	default:
		return nil, ErrInvalidShiftScheduleEvent
	}
	now := time.Now().UTC()
	item.SourceType = sourceType
	item.Action = strings.ToLower(strings.TrimSpace(item.Action))
	item.FromStatus = cleanOptional(item.FromStatus)
	item.ToStatus = cleanOptional(item.ToStatus)
	item.ActorUserID = cleanUUIDOptional(item.ActorUserID)
	item.Remarks = cleanOptional(item.Remarks)
	item.Metadata = normalizeShiftJSONRaw(item.Metadata)
	item.CreatedAt = now
	item.UpdatedAt = now
	return &item, nil
}

func NormalizeShiftAssignmentStatus(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = ShiftScheduleStatusDraft
	}
	switch clean {
	case ShiftScheduleStatusDraft, ShiftScheduleStatusPublished, ShiftScheduleStatusLocked, ShiftScheduleStatusCancelled, ShiftScheduleStatusCompleted:
		return clean, nil
	}
	return "", ErrInvalidShiftAssignment
}

func NormalizeShiftAssignmentSource(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = ShiftScheduleSourceManual
	}
	switch clean {
	case ShiftScheduleSourceManual, ShiftScheduleSourceTemplate, ShiftScheduleSourceImport, ShiftScheduleSourceSwap, ShiftScheduleSourceSystem:
		return clean, nil
	}
	return "", ErrInvalidShiftAssignment
}

func NormalizeShiftSwapStatus(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = ShiftSwapStatusPending
	}
	switch clean {
	case ShiftSwapStatusPending, ShiftSwapStatusApproved, ShiftSwapStatusRejected, ShiftSwapStatusCancelled:
		return clean, nil
	}
	return "", ErrInvalidShiftSwapRequest
}

func NormalizeStaffingRequirementStatus(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = StaffingRequirementStatusActive
	}
	switch clean {
	case StaffingRequirementStatusActive, StaffingRequirementStatusPaused, StaffingRequirementStatusArchived:
		return clean, nil
	}
	return "", ErrInvalidStaffingRequirement
}

func validateShiftClockRange(startTime, endTime string) (string, string, bool) {
	start := strings.TrimSpace(startTime)
	end := strings.TrimSpace(endTime)
	startParsed, startErr := time.Parse("15:04", start)
	endParsed, endErr := time.Parse("15:04", end)
	if startErr != nil || endErr != nil || !endParsed.After(startParsed) {
		return "", "", false
	}
	return start, end, true
}

func normalizeShiftWorkMode(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = AttendanceWorkModeOffice
	}
	if !validAttendanceWorkMode(clean) {
		return "", ErrInvalidShiftAssignment
	}
	return clean, nil
}

func normalizeShiftLocationType(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = "office"
	}
	switch clean {
	case "office", "remote", "field", "client_site", "hybrid", "branch", "warehouse", "project_site", "other":
		return clean, nil
	}
	return "", ErrInvalidShiftAssignment
}

func normalizeShiftPriority(value string) (string, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "" {
		clean = WorkflowSeverityMedium
	}
	switch clean {
	case WorkflowSeverityLow, WorkflowSeverityMedium, WorkflowSeverityHigh, WorkflowSeverityCritical:
		return clean, nil
	}
	return "", ErrInvalidStaffingRequirement
}

func normalizeShiftJSONRaw(value json.RawMessage) json.RawMessage {
	if len(value) == 0 || !json.Valid(value) {
		return json.RawMessage(`{}`)
	}
	return value
}

func datePtrOnly(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	clean := dateOnly(*value)
	return &clean
}
