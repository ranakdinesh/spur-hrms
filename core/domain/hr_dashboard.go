package domain

import (
	"time"

	"github.com/google/uuid"
)

type HRDashboard struct {
	GeneratedAt       time.Time                 `json:"generated_at"`
	WindowStart       time.Time                 `json:"window_start"`
	WindowEnd         time.Time                 `json:"window_end"`
	Headcount         HRDashboardHeadcount      `json:"headcount"`
	Attendance        HRDashboardAttendance     `json:"attendance"`
	Leave             HRDashboardLeave          `json:"leave"`
	Payroll           HRDashboardPayroll        `json:"payroll"`
	Onboarding        HRDashboardOnboarding     `json:"onboarding"`
	Policies          HRDashboardPolicies       `json:"policies"`
	Celebrations      []*EmployeeDashboardEvent `json:"celebrations"`
	UpcomingServices  []*HRDashboardComingSoon  `json:"upcoming_services"`
	OperationalAlerts []*HRDashboardAlert       `json:"operational_alerts"`
}

type HRDashboardHeadcount struct {
	TotalEmployees      int32                      `json:"total_employees"`
	ActiveEmployees     int32                      `json:"active_employees"`
	InactiveEmployees   int32                      `json:"inactive_employees"`
	NewJoinersThisMonth int32                      `json:"new_joiners_this_month"`
	Departments         []*HRDashboardDistribution `json:"departments"`
	Branches            []*HRDashboardDistribution `json:"branches"`
	Designations        []*HRDashboardDistribution `json:"designations"`
	EmploymentTypes     []*HRDashboardDistribution `json:"employment_types"`
}

type HRDashboardDistribution struct {
	Name  string `json:"name"`
	Count int32  `json:"count"`
}

type HRDashboardAttendance struct {
	TodaySummary       *AttendanceStatusSummary      `json:"today_summary,omitempty"`
	MonthSummary       *AttendanceReportSummary      `json:"month_summary,omitempty"`
	DailyTrends        []*AttendanceDailyTrend       `json:"daily_trends"`
	Departments        []*AttendanceDepartmentReport `json:"departments"`
	ExceptionEmployees []*AttendanceReportRow        `json:"exception_employees"`
}

type HRDashboardLeave struct {
	Summary         *LeaveReportSummary `json:"summary,omitempty"`
	PendingRequests int32               `json:"pending_requests"`
	RecentRequests  []*LeaveReportRow   `json:"recent_requests"`
}

type HRDashboardPayroll struct {
	Month            int32   `json:"month"`
	Year             int32   `json:"year"`
	GeneratedSlips   int32   `json:"generated_slips"`
	PendingSlips     int32   `json:"pending_slips"`
	TotalGrossSalary float64 `json:"total_gross_salary"`
	TotalNetSalary   float64 `json:"total_net_salary"`
	TotalDeductions  float64 `json:"total_deductions"`
}

type HRDashboardOnboarding struct {
	CompleteEmployees      int32 `json:"complete_employees"`
	IncompleteEmployees    int32 `json:"incomplete_employees"`
	PendingReviewDocuments int32 `json:"pending_review_documents"`
	RejectedDocuments      int32 `json:"rejected_documents"`
}

type HRDashboardPolicies struct {
	PublishedPolicies int32 `json:"published_policies"`
	RequiredDocuments int32 `json:"required_documents"`
}

type HRDashboardComingSoon struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
}

type HRDashboardAlert struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
	Detail   string `json:"detail"`
}

type HRDashboardQuery struct {
	TenantID uuid.UUID `json:"tenant_id"`
	Month    int32     `json:"month"`
	Year     int32     `json:"year"`
}
