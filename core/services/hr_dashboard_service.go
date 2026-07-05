package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetHRDashboard(ctx context.Context, query domain.HRDashboardQuery) (*domain.HRDashboard, error) {
	if query.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate hr dashboard tenant", err)
		return nil, err
	}
	now := time.Now().UTC()
	if query.Month <= 0 || query.Month > 12 {
		query.Month = int32(now.Month())
	}
	if query.Year <= 0 {
		query.Year = int32(now.Year())
	}
	windowStart := time.Date(int(query.Year), time.Month(query.Month), 1, 0, 0, 0, 0, time.UTC)
	windowEnd := windowStart.AddDate(0, 1, -1)
	if windowEnd.After(now) {
		windowEnd = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}

	employees, err := s.ListEmployees(ctx, query.TenantID)
	if err != nil {
		s.logError("load hr dashboard employees", err, serviceTenantIDField(query.TenantID))
		return nil, err
	}
	headcount := dashboardHeadcount(employees, windowStart, windowEnd)
	today := now.Format("2006-01-02")
	todayAttendance, err := s.GetAttendanceStatusSummary(ctx, ports.AttendanceStatusQuery{TenantID: query.TenantID, Date: today})
	if err != nil {
		s.logError("load hr dashboard today attendance", err, serviceTenantIDField(query.TenantID), serviceStringField("date", today))
		return nil, err
	}
	attendanceReport, err := s.GetAttendanceReport(ctx, ports.AttendanceReportQuery{TenantID: query.TenantID, StartDate: windowStart.Format("2006-01-02"), EndDate: windowEnd.Format("2006-01-02")})
	if err != nil {
		s.logError("load hr dashboard attendance report", err, serviceTenantIDField(query.TenantID))
		return nil, err
	}
	leave, err := s.dashboardHRLeave(ctx, query.TenantID, windowStart, windowEnd)
	if err != nil {
		return nil, err
	}
	payroll, err := s.dashboardHRPayroll(ctx, query.TenantID, query.Month, query.Year, headcount.ActiveEmployees)
	if err != nil {
		return nil, err
	}
	onboarding, err := s.dashboardHROnboarding(ctx, query.TenantID, employees)
	if err != nil {
		return nil, err
	}
	policies, err := s.dashboardHRPolicies(ctx, query.TenantID)
	if err != nil {
		return nil, err
	}
	celebrations, err := s.dashboardCelebrations(ctx, query.TenantID, uuid.Nil, now)
	if err != nil {
		return nil, err
	}
	return &domain.HRDashboard{
		GeneratedAt:       now,
		WindowStart:       windowStart,
		WindowEnd:         windowEnd,
		Headcount:         headcount,
		Attendance:        dashboardHRAttendance(todayAttendance, attendanceReport),
		Leave:             *leave,
		Payroll:           *payroll,
		Onboarding:        *onboarding,
		Policies:          *policies,
		Celebrations:      celebrations,
		UpcomingServices:  dashboardComingSoonServices(),
		OperationalAlerts: dashboardHRAlerts(todayAttendance, attendanceReport, leave, payroll, onboarding),
	}, nil
}

func dashboardHeadcount(employees []*domain.EmployeeListItem, windowStart time.Time, windowEnd time.Time) domain.HRDashboardHeadcount {
	departments := map[string]int32{}
	branches := map[string]int32{}
	designations := map[string]int32{}
	employmentTypes := map[string]int32{}
	result := domain.HRDashboardHeadcount{}
	for _, employee := range employees {
		if employee == nil {
			continue
		}
		result.TotalEmployees++
		if employee.Inactive {
			result.InactiveEmployees++
			continue
		}
		result.ActiveEmployees++
		if employee.JoiningDate != nil && !employee.JoiningDate.Before(windowStart) && !employee.JoiningDate.After(windowEnd) {
			result.NewJoinersThisMonth++
		}
		incrementDistribution(departments, valueOrUnassigned(employee.DepartmentName))
		incrementDistribution(branches, valueOrUnassigned(employee.BranchName))
		incrementDistribution(designations, valueOrUnassigned(employee.DesignationName))
		incrementDistribution(employmentTypes, valueOrUnassigned(employee.EmploymentTypeName))
	}
	result.Departments = topDistributions(departments, 8)
	result.Branches = topDistributions(branches, 8)
	result.Designations = topDistributions(designations, 8)
	result.EmploymentTypes = topDistributions(employmentTypes, 8)
	return result
}

func dashboardHRAttendance(today *domain.AttendanceStatusSummary, report *domain.AttendanceReport) domain.HRDashboardAttendance {
	result := domain.HRDashboardAttendance{TodaySummary: today}
	if report == nil {
		return result
	}
	result.MonthSummary = &report.Summary
	result.DailyTrends = report.DailyTrends
	result.Departments = report.Departments
	result.ExceptionEmployees = report.ExceptionEmployees
	if len(result.ExceptionEmployees) > 10 {
		result.ExceptionEmployees = result.ExceptionEmployees[:10]
	}
	return result
}

func (s *TenantService) dashboardHRLeave(ctx context.Context, tenantID uuid.UUID, windowStart time.Time, windowEnd time.Time) (*domain.HRDashboardLeave, error) {
	filter := domain.LeaveReportFilter{TenantID: tenantID, StartDate: &windowStart, EndDate: &windowEnd}
	summary, err := s.GetLeaveReportSummary(ctx, filter)
	if err != nil {
		s.logError("load hr dashboard leave summary", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	rows, err := s.ListLeaveReportRows(ctx, filter)
	if err != nil {
		s.logError("load hr dashboard leave rows", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].AppliedDate.After(rows[j].AppliedDate) })
	result := &domain.HRDashboardLeave{Summary: summary, PendingRequests: summary.PendingCount}
	for _, row := range rows {
		if row != nil && row.Status == domain.LeaveStatusPending && len(result.RecentRequests) < 8 {
			result.RecentRequests = append(result.RecentRequests, row)
		}
	}
	return result, nil
}

func (s *TenantService) dashboardHRPayroll(ctx context.Context, tenantID uuid.UUID, month int32, year int32, activeEmployees int32) (*domain.HRDashboardPayroll, error) {
	slips, err := s.ListSalarySlipsByTenantPeriod(ctx, tenantID, month, year)
	if err != nil {
		s.logError("load hr dashboard salary slips", err, serviceTenantIDField(tenantID), serviceStringField("period", fmt.Sprintf("%02d-%d", month, year)))
		return nil, err
	}
	result := &domain.HRDashboardPayroll{Month: month, Year: year, GeneratedSlips: int32(len(slips))}
	for _, slip := range slips {
		if slip == nil {
			continue
		}
		result.TotalGrossSalary += slip.GrossSalary
		result.TotalNetSalary += slip.NetSalary
		result.TotalDeductions += slip.TotalDeductions
	}
	if activeEmployees > result.GeneratedSlips {
		result.PendingSlips = activeEmployees - result.GeneratedSlips
	}
	return result, nil
}

func (s *TenantService) dashboardHROnboarding(ctx context.Context, tenantID uuid.UUID, employees []*domain.EmployeeListItem) (*domain.HRDashboardOnboarding, error) {
	result := &domain.HRDashboardOnboarding{}
	for _, employee := range employees {
		if employee == nil || employee.Inactive {
			continue
		}
		profile, err := s.GetEmployeeProfile(ctx, tenantID, employee.ID, nil)
		if err != nil {
			s.log.Warn().Err(err).Str("tenant_id", tenantID.String()).Str("employee_id", employee.ID.String()).Msg("hrms: skipped employee onboarding dashboard row")
			continue
		}
		if profile.Onboarding.IsComplete {
			result.CompleteEmployees++
		} else {
			result.IncompleteEmployees++
		}
		result.PendingReviewDocuments += int32(profile.Onboarding.PendingReviewDocuments)
		result.RejectedDocuments += int32(profile.Onboarding.RejectedDocuments)
	}
	return result, nil
}

func (s *TenantService) dashboardHRPolicies(ctx context.Context, tenantID uuid.UUID) (*domain.HRDashboardPolicies, error) {
	policies, err := s.ListCompanyPolicies(ctx, tenantID, nil)
	if err != nil {
		s.logError("load hr dashboard policies", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	documentTypes, err := s.ListDocumentTypes(ctx, tenantID)
	if err != nil {
		s.logError("load hr dashboard document types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result := &domain.HRDashboardPolicies{PublishedPolicies: int32(len(policies))}
	for _, documentType := range documentTypes {
		if documentType != nil && documentType.IsRequired {
			result.RequiredDocuments++
		}
	}
	return result, nil
}

func dashboardHRAlerts(today *domain.AttendanceStatusSummary, report *domain.AttendanceReport, leave *domain.HRDashboardLeave, payroll *domain.HRDashboardPayroll, onboarding *domain.HRDashboardOnboarding) []*domain.HRDashboardAlert {
	alerts := []*domain.HRDashboardAlert{}
	if today != nil && today.Absent > 0 {
		alerts = append(alerts, &domain.HRDashboardAlert{Key: "today_absent", Title: "Absences today", Severity: "warning", Detail: fmt.Sprintf("%d employees are absent today.", today.Absent)})
	}
	if today != nil && today.Incomplete > 0 {
		alerts = append(alerts, &domain.HRDashboardAlert{Key: "today_incomplete_attendance", Title: "Incomplete attendance", Severity: "warning", Detail: fmt.Sprintf("%d employees have check-in without check-out today.", today.Incomplete)})
	}
	if report != nil && report.Summary.IncompleteDays > 0 {
		alerts = append(alerts, &domain.HRDashboardAlert{Key: "monthly_incomplete_attendance", Title: "Attendance regularisation pending", Severity: "warning", Detail: fmt.Sprintf("%d attendance days need checkout regularisation this month.", report.Summary.IncompleteDays)})
	}
	if report != nil && report.Summary.LateDays > 0 {
		alerts = append(alerts, &domain.HRDashboardAlert{Key: "late_arrivals", Title: "Late arrivals", Severity: "info", Detail: fmt.Sprintf("%d late attendance days this month.", report.Summary.LateDays)})
	}
	if leave != nil && leave.PendingRequests > 0 {
		alerts = append(alerts, &domain.HRDashboardAlert{Key: "pending_leave", Title: "Leave approvals pending", Severity: "warning", Detail: fmt.Sprintf("%d leave requests are pending.", leave.PendingRequests)})
	}
	if payroll != nil && payroll.PendingSlips > 0 {
		alerts = append(alerts, &domain.HRDashboardAlert{Key: "pending_payroll", Title: "Payslips pending", Severity: "warning", Detail: fmt.Sprintf("%d active employees do not have payslips for the selected month.", payroll.PendingSlips)})
	}
	if onboarding != nil && onboarding.IncompleteEmployees > 0 {
		alerts = append(alerts, &domain.HRDashboardAlert{Key: "onboarding_incomplete", Title: "Onboarding incomplete", Severity: "info", Detail: fmt.Sprintf("%d active employees have incomplete onboarding.", onboarding.IncompleteEmployees)})
	}
	return alerts
}

func dashboardComingSoonServices() []*domain.HRDashboardComingSoon {
	return []*domain.HRDashboardComingSoon{
		{Key: "recruitment", Title: "Recruitment pipeline", Description: "Open positions, applicants, interview stages, time-to-fill, and offer acceptance.", Reason: "Track hiring activity from requisition to offer."},
		{Key: "performance", Title: "Performance and goals", Description: "Review cycles, goal completion, manager ratings, and performance risk signals.", Reason: "Review employee goals and performance checkpoints."},
		{Key: "training", Title: "Training and compliance", Description: "Mandatory training, certification expiry, course completion, and learning hours.", Reason: "Monitor learning and compliance readiness."},
		{Key: "engagement", Title: "Engagement pulse", Description: "Survey scores, sentiment trends, recognition, and retention risk.", Reason: "Watch team sentiment and retention signals."},
	}
}

func incrementDistribution(items map[string]int32, name string) {
	items[name]++
}

func topDistributions(items map[string]int32, limit int) []*domain.HRDashboardDistribution {
	result := make([]*domain.HRDashboardDistribution, 0, len(items))
	for name, count := range items {
		result = append(result, &domain.HRDashboardDistribution{Name: name, Count: count})
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Name < result[j].Name
		}
		return result[i].Count > result[j].Count
	})
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result
}

func valueOrUnassigned(value *string) string {
	if value == nil || *value == "" {
		return "Unassigned"
	}
	return *value
}
