package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	AttendanceSourceWeb       = "web"
	AttendanceSourceMobile    = "mobile"
	AttendanceSourceKiosk     = "kiosk"
	AttendanceSourceBiometric = "biometric"
	AttendanceSourceAPI       = "api"

	AttendanceWorkModeOffice = "office"
	AttendanceWorkModeRemote = "remote"
	AttendanceWorkModeField  = "field"
	AttendanceWorkModeHybrid = "hybrid"
)

var (
	ErrInvalidAttendanceID         = errors.New("attendance id is required")
	ErrInvalidAttendanceDate       = errors.New("attendance date is required")
	ErrInvalidAttendanceAction     = errors.New("attendance action is invalid")
	ErrInvalidAttendanceSource     = errors.New("attendance source is invalid")
	ErrInvalidAttendanceWorkMode   = errors.New("attendance work mode is invalid")
	ErrInvalidAttendanceLocation   = errors.New("attendance location is invalid")
	ErrAttendanceLocationRequired  = errors.New("attendance location is required")
	ErrAttendanceAlreadyCheckedIn  = errors.New("employee is already checked in for this date")
	ErrAttendanceNotCheckedIn      = errors.New("employee must check in before checking out")
	ErrAttendanceAlreadyCheckedOut = errors.New("employee is already checked out for this date")
	ErrAttendanceNotRequired       = errors.New("attendance is not required for this employee designation")
	ErrAttendanceNotFound          = errors.New("attendance not found")
)

type Attendance struct {
	ID                         uuid.UUID  `json:"id"`
	TenantID                   uuid.UUID  `json:"tenant_id"`
	UserID                     uuid.UUID  `json:"user_id"`
	Date                       time.Time  `json:"date"`
	Time                       *time.Time `json:"time,omitempty"`
	Type                       *string    `json:"type,omitempty"`
	Status                     *string    `json:"status,omitempty"`
	Source                     *string    `json:"source,omitempty"`
	Latitude                   *float64   `json:"latitude,omitempty"`
	Longitude                  *float64   `json:"longitude,omitempty"`
	WorkMode                   *string    `json:"work_mode,omitempty"`
	Remarks                    *string    `json:"remarks,omitempty"`
	AttendanceLocationID       *uuid.UUID `json:"attendance_location_id,omitempty"`
	AttendanceDeviceID         *uuid.UUID `json:"attendance_device_id,omitempty"`
	RawAttendanceEventID       *uuid.UUID `json:"raw_attendance_event_id,omitempty"`
	LocationAccuracyMeters     *float64   `json:"location_accuracy_meters,omitempty"`
	LocationVerificationStatus string     `json:"location_verification_status"`
	Inactive                   bool       `json:"inactive"`
	CreatedAt                  time.Time  `json:"created_at"`
	CreatedBy                  *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt                  time.Time  `json:"updated_at"`
	UpdatedBy                  *uuid.UUID `json:"updated_by,omitempty"`
}

type DeviceLog struct {
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	UserID     uuid.UUID  `json:"user_id"`
	DeviceID   *string    `json:"device_id,omitempty"`
	DeviceType *string    `json:"device_type,omitempty"`
	IPAddress  *string    `json:"ip_address,omitempty"`
	Action     *string    `json:"action,omitempty"`
	Inactive   bool       `json:"inactive"`
	LoggedAt   time.Time  `json:"logged_at"`
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at"`
	UpdatedBy  *uuid.UUID `json:"updated_by,omitempty"`
}

type AttendancePunch struct {
	Attendance *Attendance `json:"attendance"`
	DeviceLog  *DeviceLog  `json:"device_log,omitempty"`
}

func NewAttendance(tenantID, userID uuid.UUID, date time.Time, punchTime time.Time, action string, source *string, latitude *float64, longitude *float64, workMode *string, remarks *string) (*Attendance, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidLeaveUser
	}
	if date.IsZero() {
		return nil, ErrInvalidAttendanceDate
	}
	normalizedAction, err := ValidateAttendanceType(action)
	if err != nil {
		return nil, ErrInvalidAttendanceAction
	}
	if source != nil {
		clean := strings.ToLower(strings.TrimSpace(*source))
		if !validAttendanceSource(clean) {
			return nil, ErrInvalidAttendanceSource
		}
		source = &clean
	}
	if workMode != nil {
		clean := strings.ToLower(strings.TrimSpace(*workMode))
		if !validAttendanceWorkMode(clean) {
			return nil, ErrInvalidAttendanceWorkMode
		}
		workMode = &clean
	}
	if latitude != nil && (*latitude < -90 || *latitude > 90) {
		return nil, ErrInvalidAttendanceLocation
	}
	if longitude != nil && (*longitude < -180 || *longitude > 180) {
		return nil, ErrInvalidAttendanceLocation
	}
	status := AttendanceStatusPresent
	locationStatus := "not_checked"
	now := time.Now().UTC()
	return &Attendance{TenantID: tenantID, UserID: userID, Date: dateOnly(date), Time: timePtr(punchTime.UTC()), Type: &normalizedAction, Status: &status, Source: cleanOptional(source), Latitude: latitude, Longitude: longitude, WorkMode: cleanOptional(workMode), Remarks: cleanOptional(remarks), LocationVerificationStatus: locationStatus, CreatedAt: now, UpdatedAt: now}, nil
}

func NewDeviceLog(tenantID, userID uuid.UUID, deviceID *string, deviceType *string, ipAddress *string, action string) (*DeviceLog, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidLeaveUser
	}
	normalizedAction, err := ValidateAttendanceType(action)
	if err != nil {
		return nil, ErrInvalidAttendanceAction
	}
	now := time.Now().UTC()
	return &DeviceLog{TenantID: tenantID, UserID: userID, DeviceID: cleanOptional(deviceID), DeviceType: cleanOptional(deviceType), IPAddress: cleanOptional(ipAddress), Action: &normalizedAction, LoggedAt: now, CreatedAt: now, UpdatedAt: now}, nil
}

func validAttendanceSource(value string) bool {
	switch value {
	case AttendanceSourceWeb, AttendanceSourceMobile, AttendanceSourceKiosk, AttendanceSourceBiometric, AttendanceSourceAPI:
		return true
	default:
		return false
	}
}

func validAttendanceWorkMode(value string) bool {
	switch value {
	case AttendanceWorkModeOffice, AttendanceWorkModeRemote, AttendanceWorkModeField, AttendanceWorkModeHybrid:
		return true
	default:
		return false
	}
}

func timePtr(value time.Time) *time.Time {
	return &value
}

const (
	AttendanceStatusNotApplicable = "not_applicable"
	AttendanceStatusEmpty         = "empty"
)

type AttendanceDailyStatus struct {
	TenantID          uuid.UUID         `json:"tenant_id"`
	UserID            uuid.UUID         `json:"user_id"`
	EmployeeID        uuid.UUID         `json:"employee_id"`
	EmployeeCode      *string           `json:"employee_code,omitempty"`
	Firstname         string            `json:"firstname"`
	Lastname          *string           `json:"lastname,omitempty"`
	DepartmentID      *uuid.UUID        `json:"department_id,omitempty"`
	DepartmentName    *string           `json:"department_name,omitempty"`
	BranchID          *uuid.UUID        `json:"branch_id,omitempty"`
	BranchName        *string           `json:"branch_name,omitempty"`
	Date              time.Time         `json:"date"`
	Status            string            `json:"status"`
	Reason            string            `json:"reason"`
	WorkingHour       *WorkingHour      `json:"working_hour,omitempty"`
	Policy            *AttendancePolicy `json:"policy,omitempty"`
	Roster            *AttendanceRoster `json:"roster,omitempty"`
	Holiday           *Holiday          `json:"holiday,omitempty"`
	Leave             *Leave            `json:"leave,omitempty"`
	FirstCheckIn      *time.Time        `json:"first_check_in,omitempty"`
	LastCheckOut      *time.Time        `json:"last_check_out,omitempty"`
	WorkedMinutes     int32             `json:"worked_minutes"`
	LateMinutes       int32             `json:"late_minutes"`
	EarlyExitMinutes  int32             `json:"early_exit_minutes"`
	RuleOutcome       string            `json:"rule_outcome,omitempty"`
	AttendanceRecords []*Attendance     `json:"attendance_records"`
}

type AttendanceStatusSummary struct {
	Date               time.Time        `json:"date"`
	TotalEmployees     int32            `json:"total_employees"`
	Present            int32            `json:"present"`
	Leave              int32            `json:"leave"`
	Absent             int32            `json:"absent"`
	Holiday            int32            `json:"holiday"`
	Weekoff            int32            `json:"weekoff"`
	Incomplete         int32            `json:"incomplete"`
	Empty              int32            `json:"empty"`
	NotApplicable      int32            `json:"not_applicable"`
	ByStatus           map[string]int32 `json:"by_status"`
	TotalWorkedMinutes int32            `json:"total_worked_minutes"`
}
