package domain

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceReportFilter struct {
	TenantID     uuid.UUID  `json:"tenant_id"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	BranchID     *uuid.UUID `json:"branch_id,omitempty"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      time.Time  `json:"end_date"`
}

type AttendanceReportRow struct {
	TenantID         uuid.UUID  `json:"tenant_id"`
	UserID           uuid.UUID  `json:"user_id"`
	EmployeeID       uuid.UUID  `json:"employee_id"`
	EmployeeCode     *string    `json:"employee_code,omitempty"`
	Firstname        string     `json:"firstname"`
	Lastname         *string    `json:"lastname,omitempty"`
	DepartmentID     *uuid.UUID `json:"department_id,omitempty"`
	DepartmentName   *string    `json:"department_name,omitempty"`
	BranchID         *uuid.UUID `json:"branch_id,omitempty"`
	BranchName       *string    `json:"branch_name,omitempty"`
	Date             time.Time  `json:"date"`
	Status           string     `json:"status"`
	Reason           string     `json:"reason"`
	RuleOutcome      string     `json:"rule_outcome,omitempty"`
	PolicyName       *string    `json:"policy_name,omitempty"`
	ScheduleType     *string    `json:"schedule_type,omitempty"`
	WorkMode         *string    `json:"work_mode,omitempty"`
	FirstCheckIn     *time.Time `json:"first_check_in,omitempty"`
	LastCheckOut     *time.Time `json:"last_check_out,omitempty"`
	WorkedMinutes    int32      `json:"worked_minutes"`
	LateMinutes      int32      `json:"late_minutes"`
	EarlyExitMinutes int32      `json:"early_exit_minutes"`
	PunchCount       int32      `json:"punch_count"`
}

type AttendanceReportSummary struct {
	StartDate            time.Time        `json:"start_date"`
	EndDate              time.Time        `json:"end_date"`
	EmployeeDays         int32            `json:"employee_days"`
	PresentDays          int32            `json:"present_days"`
	AbsentDays           int32            `json:"absent_days"`
	LeaveDays            int32            `json:"leave_days"`
	HolidayDays          int32            `json:"holiday_days"`
	WeekoffDays          int32            `json:"weekoff_days"`
	IncompleteDays       int32            `json:"incomplete_days"`
	EmptyDays            int32            `json:"empty_days"`
	LateDays             int32            `json:"late_days"`
	HalfDays             int32            `json:"half_days"`
	EarlyExitDays        int32            `json:"early_exit_days"`
	OvertimeDays         int32            `json:"overtime_days"`
	TotalWorkedMinutes   int32            `json:"total_worked_minutes"`
	AverageWorkedMinutes int32            `json:"average_worked_minutes"`
	AttendanceRate       float64          `json:"attendance_rate"`
	AbsenteeismRate      float64          `json:"absenteeism_rate"`
	LateRate             float64          `json:"late_rate"`
	PendingRequests      int32            `json:"pending_requests"`
	ByStatus             map[string]int32 `json:"by_status"`
	ByOutcome            map[string]int32 `json:"by_outcome"`
}

type AttendanceDepartmentReport struct {
	DepartmentName       string  `json:"department_name"`
	EmployeeDays         int32   `json:"employee_days"`
	PresentDays          int32   `json:"present_days"`
	AbsentDays           int32   `json:"absent_days"`
	IncompleteDays       int32   `json:"incomplete_days"`
	LateDays             int32   `json:"late_days"`
	TotalWorkedMinutes   int32   `json:"total_worked_minutes"`
	AverageWorkedMinutes int32   `json:"average_worked_minutes"`
	AttendanceRate       float64 `json:"attendance_rate"`
}

type AttendanceDailyTrend struct {
	Date                 time.Time `json:"date"`
	EmployeeDays         int32     `json:"employee_days"`
	PresentDays          int32     `json:"present_days"`
	AbsentDays           int32     `json:"absent_days"`
	LateDays             int32     `json:"late_days"`
	TotalWorkedMinutes   int32     `json:"total_worked_minutes"`
	AverageWorkedMinutes int32     `json:"average_worked_minutes"`
}

type AttendanceWorkModeReport struct {
	WorkMode      string  `json:"work_mode"`
	Days          int32   `json:"days"`
	WorkedMinutes int32   `json:"worked_minutes"`
	SharePercent  float64 `json:"share_percent"`
}

type AttendanceReport struct {
	Filter             AttendanceReportFilter        `json:"filter"`
	Summary            AttendanceReportSummary       `json:"summary"`
	Rows               []*AttendanceReportRow        `json:"rows"`
	Departments        []*AttendanceDepartmentReport `json:"departments"`
	DailyTrends        []*AttendanceDailyTrend       `json:"daily_trends"`
	WorkModes          []*AttendanceWorkModeReport   `json:"work_modes"`
	LateEmployees      []*AttendanceReportRow        `json:"late_employees"`
	AbsenceEmployees   []*AttendanceReportRow        `json:"absence_employees"`
	ExceptionEmployees []*AttendanceReportRow        `json:"exception_employees"`
}
