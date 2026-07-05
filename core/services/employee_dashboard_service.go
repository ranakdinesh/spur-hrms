package services

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetEmployeeDashboard(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.EmployeeDashboard, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employee dashboard tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidEmployeeUserID
		s.logError("validate employee dashboard user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	profile, err := s.GetEmployeeSelfProfile(ctx, tenantID, userID, &userID)
	if err != nil {
		s.logError("load employee dashboard profile", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	now := time.Now().UTC()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	today := now.Format("2006-01-02")

	todayStatuses, err := s.ListAttendanceDailyStatuses(ctx, ports.AttendanceStatusQuery{TenantID: tenantID, UserID: &userID, Date: today})
	if err != nil {
		s.logError("load employee dashboard today attendance", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	recentStatuses, err := s.dashboardRecentAttendance(ctx, tenantID, userID, now)
	if err != nil {
		return nil, err
	}
	report, err := s.GetAttendanceReport(ctx, ports.AttendanceReportQuery{TenantID: tenantID, UserID: &userID, StartDate: startOfMonth.Format("2006-01-02"), EndDate: today})
	if err != nil {
		s.logError("load employee dashboard attendance report", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	leave, err := s.dashboardLeave(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	payslips, err := s.dashboardPayslips(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	policies, err := s.dashboardPolicies(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	celebrations, err := s.dashboardCelebrations(ctx, tenantID, userID, now)
	if err != nil {
		return nil, err
	}

	var todayStatus *domain.AttendanceDailyStatus
	if len(todayStatuses) > 0 {
		todayStatus = todayStatuses[0]
	}
	return &domain.EmployeeDashboard{
		GeneratedAt:  now,
		Profile:      dashboardProfile(profile),
		Attendance:   dashboardAttendance(todayStatus, recentStatuses, report),
		Leave:        leave,
		Payslips:     payslips,
		Policies:     policies,
		Celebrations: celebrations,
		QuickTools:   dashboardQuickTools(),
		Onboarding:   profile.Onboarding,
	}, nil
}

func (s *TenantService) dashboardRecentAttendance(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, now time.Time) ([]*domain.AttendanceDailyStatus, error) {
	items := []*domain.AttendanceDailyStatus{}
	for day := now.AddDate(0, 0, -6); !day.After(now); day = day.AddDate(0, 0, 1) {
		statuses, err := s.ListAttendanceDailyStatuses(ctx, ports.AttendanceStatusQuery{TenantID: tenantID, UserID: &userID, Date: day.Format("2006-01-02")})
		if err != nil {
			s.logError("load employee dashboard recent attendance", err, serviceTenantIDField(tenantID), serviceStringField("date", day.Format("2006-01-02")))
			return nil, err
		}
		if len(statuses) > 0 {
			items = append(items, statuses[0])
		}
	}
	return items, nil
}

func (s *TenantService) dashboardLeave(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.EmployeeDashboardLeave, error) {
	balances, err := s.ListLeaveBalancesByUser(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	leaveTypes, err := s.ListLeaveTypes(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	typeNames := map[uuid.UUID]string{}
	for _, item := range leaveTypes {
		if item != nil {
			typeNames[item.ID] = item.Name
		}
	}
	leave := &domain.EmployeeDashboardLeave{Balances: []*domain.EmployeeDashboardLeaveBalance{}, RecentRequests: []*domain.Leave{}}
	for _, balance := range balances {
		if balance == nil {
			continue
		}
		name := typeNames[balance.LeaveTypeID]
		if name == "" {
			name = balance.LeaveTypeID.String()
		}
		leave.Balances = append(leave.Balances, &domain.EmployeeDashboardLeaveBalance{LeaveTypeID: balance.LeaveTypeID, LeaveTypeName: name, TotalDays: balance.TotalDays, UsedDays: balance.UsedDays, PendingDays: balance.PendingDays, BalanceDays: balance.BalanceDays})
		leave.AvailableDays += balance.BalanceDays
		leave.PendingDays += balance.PendingDays
		leave.UsedDays += balance.UsedDays
	}
	requests, err := s.ListLeavesByUser(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(requests, func(i, j int) bool { return requests[i].AppliedDate.After(requests[j].AppliedDate) })
	for _, request := range requests {
		if request == nil {
			continue
		}
		if request.Status == domain.LeaveStatusPending {
			leave.PendingRequests++
		}
		if len(leave.RecentRequests) < 5 {
			leave.RecentRequests = append(leave.RecentRequests, request)
		}
	}
	return leave, nil
}

func (s *TenantService) dashboardPayslips(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.EmployeeDashboardPayslip, error) {
	items, err := s.salarySlips.ListRecentSalarySlipsByUser(ctx, tenantID, userID, 6)
	if err != nil {
		s.logError("load employee dashboard payslips", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	result := make([]*domain.EmployeeDashboardPayslip, 0, len(items))
	for _, item := range items {
		if item != nil {
			result = append(result, &domain.EmployeeDashboardPayslip{ID: item.ID, Month: item.Month, Year: item.Year, NetSalary: item.NetSalary, PDFPath: item.PDFPath, CreatedAt: item.CreatedAt})
		}
	}
	return result, nil
}

func (s *TenantService) dashboardPolicies(ctx context.Context, tenantID uuid.UUID) ([]*domain.EmployeeDashboardPolicy, error) {
	items, err := s.ListCompanyPolicies(ctx, tenantID, nil)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.EmployeeDashboardPolicy, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		result = append(result, &domain.EmployeeDashboardPolicy{ID: item.ID, Title: item.Title, Description: item.Description, FilePath: item.FilePath, UpdatedAt: item.UpdatedAt})
		if len(result) == 5 {
			break
		}
	}
	return result, nil
}

func (s *TenantService) dashboardCelebrations(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, now time.Time) ([]*domain.EmployeeDashboardEvent, error) {
	types, err := s.celebrations.ListCelebrationTypes(ctx, tenantID)
	if err != nil {
		s.logError("load employee dashboard celebration types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	typeByID := map[uuid.UUID]*domain.CelebrationType{}
	for _, item := range types {
		if item != nil {
			typeByID[item.ID] = item
		}
	}
	items, err := s.celebrations.ListCelebrations(ctx, tenantID)
	if err != nil {
		s.logError("load employee dashboard celebrations", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	events := make([]*domain.EmployeeDashboardEvent, 0, len(items))
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	for _, item := range items {
		if item == nil || item.CelebrationDate == nil {
			continue
		}
		celebrationType := typeByID[item.CelebrationTypeID]
		typeName := "Celebration"
		isYearly := true
		if celebrationType != nil {
			typeName = celebrationType.Name
			isYearly = celebrationType.IsYearly
		}
		nextDate := nextDashboardCelebrationDate(*item.CelebrationDate, today, isYearly)
		if nextDate.Before(today) {
			continue
		}
		title := typeName
		if item.CustomTitle != nil && strings.TrimSpace(*item.CustomTitle) != "" {
			title = strings.TrimSpace(*item.CustomTitle)
		}
		events = append(events, &domain.EmployeeDashboardEvent{ID: item.ID, Title: title, TypeName: typeName, Date: nextDate, DaysUntil: int(nextDate.Sub(today).Hours() / 24), UserID: item.UserID, Description: item.Description, IsPersonalEvent: item.UserID != nil && *item.UserID == userID})
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Date.Equal(events[j].Date) {
			return events[i].Title < events[j].Title
		}
		return events[i].Date.Before(events[j].Date)
	})
	if len(events) > 6 {
		events = events[:6]
	}
	return events, nil
}

func dashboardProfile(profile *domain.EmployeeProfile) *domain.EmployeeDashboardProfile {
	if profile == nil || profile.Employee == nil {
		return nil
	}
	employee := profile.Employee
	return &domain.EmployeeDashboardProfile{
		EmployeeID:       employee.ID,
		UserID:           employee.UserID,
		EmployeeCode:     employee.EmployeeCode,
		Name:             employeeDisplayName(employee.Firstname, employee.MiddleName, employee.Lastname),
		Email:            employee.Email,
		Mobile:           employee.Mobile,
		DepartmentName:   employee.DepartmentName,
		BranchName:       employee.BranchName,
		DesignationName:  employee.DesignationName,
		EmploymentType:   employee.EmploymentTypeName,
		ProfilePhotoPath: employee.ProfilePhotoPath,
		JoiningDate:      employee.JoiningDate,
	}
}

func dashboardAttendance(today *domain.AttendanceDailyStatus, recent []*domain.AttendanceDailyStatus, report *domain.AttendanceReport) *domain.EmployeeDashboardAttendance {
	summary := &domain.AttendanceReportSummary{}
	if report != nil {
		summary = &report.Summary
	}
	work := domain.EmployeeDashboardWorkTime{MonthWorkedMinutes: summary.TotalWorkedMinutes, MonthWorkedHours: float64(summary.TotalWorkedMinutes) / 60, LateDays: summary.LateDays, EarlyExitDays: summary.EarlyExitDays}
	if today != nil {
		work.TodayWorkedMinutes = today.WorkedMinutes
	}
	return &domain.EmployeeDashboardAttendance{Today: today, MonthSummary: summary, RecentDays: recent, WorkTotals: work}
}

func dashboardQuickTools() []*domain.EmployeeDashboardTool {
	return []*domain.EmployeeDashboardTool{
		{Key: "attendance", Label: "Check in/out", Description: "Record today's attendance", Section: "attendance", Permission: "hrms.attendance.check_in"},
		{Key: "leave", Label: "Apply leave", Description: "Submit or track leave requests", Section: "leaves", Permission: "hrms.leaves.apply"},
		{Key: "payslips", Label: "Payslips", Description: "Download recent salary slips", Section: "payslips", Permission: "hrms.salary_slips.download"},
		{Key: "onboarding", Label: "Onboarding", Description: "Complete required documents", Section: "my-onboarding", Permission: "hrms.employees.documents.manage"},
		{Key: "policies", Label: "Policies", Description: "Open company policies", Section: "policies", Permission: "hrms.policies.list"},
	}
}

func employeeDisplayName(first string, middle *string, last *string) string {
	parts := []string{strings.TrimSpace(first)}
	if middle != nil && strings.TrimSpace(*middle) != "" {
		parts = append(parts, strings.TrimSpace(*middle))
	}
	if last != nil && strings.TrimSpace(*last) != "" {
		parts = append(parts, strings.TrimSpace(*last))
	}
	return strings.Join(parts, " ")
}

func nextDashboardCelebrationDate(source time.Time, today time.Time, yearly bool) time.Time {
	source = time.Date(source.Year(), source.Month(), source.Day(), 0, 0, 0, 0, time.UTC)
	if !yearly {
		return source
	}
	next := time.Date(today.Year(), source.Month(), source.Day(), 0, 0, 0, 0, time.UTC)
	if next.Before(today) {
		next = next.AddDate(1, 0, 0)
	}
	return next
}
