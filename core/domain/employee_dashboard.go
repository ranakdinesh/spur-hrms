package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCelebrationTypeID   = errors.New("celebration_type_id is required")
	ErrInvalidCelebrationTypeName = errors.New("celebration type name is required")
	ErrCelebrationTypeNotFound    = errors.New("celebration type not found")
	ErrInvalidCelebrationID       = errors.New("celebration_id is required")
	ErrInvalidCelebrationDate     = errors.New("celebration date is required")
	ErrInvalidCelebrationUser     = errors.New("user_id is required for this celebration type")
	ErrCelebrationNotFound        = errors.New("celebration not found")
	ErrCelebrationAlreadyExists   = errors.New("celebration already exists for this user and type")
)

type EmployeeDashboard struct {
	GeneratedAt  time.Time                    `json:"generated_at"`
	Profile      *EmployeeDashboardProfile    `json:"profile"`
	Attendance   *EmployeeDashboardAttendance `json:"attendance"`
	Leave        *EmployeeDashboardLeave      `json:"leave"`
	Payslips     []*EmployeeDashboardPayslip  `json:"payslips"`
	Policies     []*EmployeeDashboardPolicy   `json:"policies"`
	Celebrations []*EmployeeDashboardEvent    `json:"celebrations"`
	QuickTools   []*EmployeeDashboardTool     `json:"quick_tools"`
	Onboarding   EmployeeOnboardingStatus     `json:"onboarding"`
}

type EmployeeDashboardProfile struct {
	EmployeeID       uuid.UUID  `json:"employee_id"`
	UserID           uuid.UUID  `json:"user_id"`
	EmployeeCode     *string    `json:"employee_code,omitempty"`
	Name             string     `json:"name"`
	Email            *string    `json:"email,omitempty"`
	Mobile           *string    `json:"mobile,omitempty"`
	DepartmentName   *string    `json:"department_name,omitempty"`
	BranchName       *string    `json:"branch_name,omitempty"`
	DesignationName  *string    `json:"designation_name,omitempty"`
	EmploymentType   *string    `json:"employment_type,omitempty"`
	ProfilePhotoPath *string    `json:"profile_photo_path,omitempty"`
	JoiningDate      *time.Time `json:"joining_date,omitempty"`
}

type EmployeeDashboardAttendance struct {
	Today        *AttendanceDailyStatus    `json:"today,omitempty"`
	MonthSummary *AttendanceReportSummary  `json:"month_summary,omitempty"`
	RecentDays   []*AttendanceDailyStatus  `json:"recent_days"`
	WorkTotals   EmployeeDashboardWorkTime `json:"work_totals"`
}

type EmployeeDashboardWorkTime struct {
	TodayWorkedMinutes int32   `json:"today_worked_minutes"`
	MonthWorkedMinutes int32   `json:"month_worked_minutes"`
	MonthWorkedHours   float64 `json:"month_worked_hours"`
	LateDays           int32   `json:"late_days"`
	EarlyExitDays      int32   `json:"early_exit_days"`
}

type EmployeeDashboardLeave struct {
	Balances        []*EmployeeDashboardLeaveBalance `json:"balances"`
	RecentRequests  []*Leave                         `json:"recent_requests"`
	PendingRequests int32                            `json:"pending_requests"`
	AvailableDays   float64                          `json:"available_days"`
	PendingDays     float64                          `json:"pending_days"`
	UsedDays        float64                          `json:"used_days"`
}

type EmployeeDashboardLeaveBalance struct {
	LeaveTypeID   uuid.UUID `json:"leave_type_id"`
	LeaveTypeName string    `json:"leave_type_name"`
	TotalDays     float64   `json:"total_days"`
	UsedDays      float64   `json:"used_days"`
	PendingDays   float64   `json:"pending_days"`
	BalanceDays   float64   `json:"balance_days"`
}

type EmployeeDashboardPayslip struct {
	ID        uuid.UUID `json:"id"`
	Month     int32     `json:"month"`
	Year      int32     `json:"year"`
	NetSalary float64   `json:"net_salary"`
	PDFPath   *string   `json:"pdf_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type EmployeeDashboardPolicy struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	FilePath    *string   `json:"file_path,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CelebrationType struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	Name              string     `json:"name"`
	IsYearly          bool       `json:"is_yearly"`
	IsUserCelebration bool       `json:"is_user_celebration"`
	Inactive          bool       `json:"inactive"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty"`
}

type CelebrationTypeInput struct {
	TenantID          uuid.UUID
	Name              string
	IsYearly          bool
	IsUserCelebration bool
}

func NewCelebrationType(input CelebrationTypeInput) (*CelebrationType, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidCelebrationTypeName
	}
	now := time.Now().UTC()
	return &CelebrationType{TenantID: input.TenantID, Name: name, IsYearly: input.IsYearly, IsUserCelebration: input.IsUserCelebration, CreatedAt: now, UpdatedAt: now}, nil
}

func DefaultCelebrationTypeInputs(tenantID uuid.UUID) []CelebrationTypeInput {
	return []CelebrationTypeInput{
		{TenantID: tenantID, Name: "Birthday", IsYearly: true, IsUserCelebration: true},
		{TenantID: tenantID, Name: "Work Anniversary", IsYearly: true, IsUserCelebration: true},
		{TenantID: tenantID, Name: "Company Foundation Day", IsYearly: true, IsUserCelebration: false},
		{TenantID: tenantID, Name: "Festival", IsYearly: true, IsUserCelebration: false},
		{TenantID: tenantID, Name: "Team Event", IsYearly: false, IsUserCelebration: false},
	}
}

type Celebration struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	BranchID          *uuid.UUID `json:"branch_id,omitempty"`
	UserID            *uuid.UUID `json:"user_id,omitempty"`
	CelebrationTypeID uuid.UUID  `json:"celebration_type_id"`
	CelebrationDate   *time.Time `json:"celebration_date,omitempty"`
	CustomTitle       *string    `json:"custom_title,omitempty"`
	Description       *string    `json:"description,omitempty"`
	Inactive          bool       `json:"inactive"`
	CreatedAt         time.Time  `json:"created_at"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty"`
}

type CelebrationInput struct {
	TenantID          uuid.UUID
	BranchID          *uuid.UUID
	UserID            *uuid.UUID
	CelebrationTypeID uuid.UUID
	CelebrationDate   *time.Time
	CustomTitle       *string
	Description       *string
}

func NewCelebration(input CelebrationInput) (*Celebration, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if input.CelebrationTypeID == uuid.Nil {
		return nil, ErrInvalidCelebrationTypeID
	}
	if input.CelebrationDate == nil || input.CelebrationDate.IsZero() {
		return nil, ErrInvalidCelebrationDate
	}
	date := dateOnly(*input.CelebrationDate)
	now := time.Now().UTC()
	return &Celebration{TenantID: input.TenantID, BranchID: cleanUUIDOptional(input.BranchID), UserID: cleanUUIDOptional(input.UserID), CelebrationTypeID: input.CelebrationTypeID, CelebrationDate: &date, CustomTitle: cleanOptional(input.CustomTitle), Description: cleanOptional(input.Description), CreatedAt: now, UpdatedAt: now}, nil
}

type CelebrationListItem struct {
	Celebration
	CelebrationTypeName     string     `json:"celebration_type_name"`
	IsYearly                bool       `json:"is_yearly"`
	IsUserCelebration       bool       `json:"is_user_celebration"`
	EmployeeName            *string    `json:"employee_name,omitempty"`
	EmployeeCode            *string    `json:"employee_code,omitempty"`
	BranchName              *string    `json:"branch_name,omitempty"`
	NextOccurrenceDate      *time.Time `json:"next_occurrence_date,omitempty"`
	DaysUntilNextOccurrence *int       `json:"days_until_next_occurrence,omitempty"`
}

type EmployeeDashboardEvent struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	TypeName        string     `json:"type_name"`
	Date            time.Time  `json:"date"`
	DaysUntil       int        `json:"days_until"`
	UserID          *uuid.UUID `json:"user_id,omitempty"`
	Description     *string    `json:"description,omitempty"`
	IsPersonalEvent bool       `json:"is_personal_event"`
}

type EmployeeDashboardTool struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Section     string `json:"section"`
	Permission  string `json:"permission"`
}
