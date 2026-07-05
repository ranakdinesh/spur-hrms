package services

import (
	"context"
	"sort"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetAttendanceReport(ctx context.Context, query ports.AttendanceReportQuery) (*domain.AttendanceReport, error) {
	if query.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate attendance report tenant", err)
		return nil, err
	}
	start, end, err := parseDateRangeOrToday(query.StartDate, query.EndDate)
	if err != nil {
		s.logError("validate attendance report date range", err, serviceTenantIDField(query.TenantID), serviceStringField("start_date", query.StartDate), serviceStringField("end_date", query.EndDate))
		return nil, err
	}
	report := &domain.AttendanceReport{
		Filter:  domain.AttendanceReportFilter{TenantID: query.TenantID, UserID: query.UserID, DepartmentID: query.DepartmentID, BranchID: query.BranchID, StartDate: start, EndDate: end},
		Summary: domain.AttendanceReportSummary{StartDate: start, EndDate: end, ByStatus: map[string]int32{}, ByOutcome: map[string]int32{}},
		Rows:    []*domain.AttendanceReportRow{},
	}
	departmentAgg := map[string]*domain.AttendanceDepartmentReport{}
	dailyAgg := map[string]*domain.AttendanceDailyTrend{}
	workModeAgg := map[string]*domain.AttendanceWorkModeReport{}
	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {
		items, err := s.ListAttendanceDailyStatuses(ctx, ports.AttendanceStatusQuery{TenantID: query.TenantID, UserID: query.UserID, Date: day.Format("2006-01-02")})
		if err != nil {
			s.logError("list attendance report day statuses", err, serviceTenantIDField(query.TenantID), serviceStringField("date", day.Format("2006-01-02")))
			return nil, err
		}
		for _, item := range items {
			if item == nil || !attendanceReportRowMatches(item, query.DepartmentID, query.BranchID) {
				continue
			}
			if item.Status == domain.AttendanceStatusNotApplicable {
				continue
			}
			row := attendanceReportRowFromStatus(item)
			report.Rows = append(report.Rows, row)
			accumulateAttendanceReportSummary(&report.Summary, row)
			accumulateDepartmentAttendance(departmentAgg, row)
			accumulateDailyAttendance(dailyAgg, row)
			accumulateWorkModeAttendance(workModeAgg, row)
			if row.RuleOutcome == domain.AttendanceRuleOutcomeLate || row.LateMinutes > 0 {
				report.LateEmployees = append(report.LateEmployees, row)
			}
			if row.Status == domain.AttendanceStatusAbsent {
				report.AbsenceEmployees = append(report.AbsenceEmployees, row)
			}
			if row.Status == domain.AttendanceStatusIncomplete || row.RuleOutcome != "" && row.RuleOutcome != domain.AttendanceRuleOutcomeOnTime {
				report.ExceptionEmployees = append(report.ExceptionEmployees, row)
			}
		}
	}
	pending, err := s.attendanceRequests.ListAttendanceRequestsByStatus(ctx, query.TenantID, domain.LeaveStatusPending)
	if err == nil {
		report.Summary.PendingRequests = int32(len(pending))
	} else {
		s.log.Warn().Err(err).Str("tenant_id", query.TenantID.String()).Msg("hrms: attendance report pending request count skipped")
	}
	finalizeAttendanceReport(&report.Summary, departmentAgg, dailyAgg, workModeAgg, report)
	return report, nil
}

func attendanceReportRowMatches(item *domain.AttendanceDailyStatus, departmentID *uuid.UUID, branchID *uuid.UUID) bool {
	if departmentID != nil && *departmentID != uuid.Nil {
		return item.DepartmentID != nil && *item.DepartmentID == *departmentID
	}
	if branchID != nil && *branchID != uuid.Nil {
		return item.BranchID != nil && *item.BranchID == *branchID
	}
	return true
}

func attendanceReportRowFromStatus(item *domain.AttendanceDailyStatus) *domain.AttendanceReportRow {
	row := &domain.AttendanceReportRow{TenantID: item.TenantID, UserID: item.UserID, EmployeeID: item.EmployeeID, EmployeeCode: item.EmployeeCode, Firstname: item.Firstname, Lastname: item.Lastname, DepartmentID: item.DepartmentID, DepartmentName: item.DepartmentName, BranchID: item.BranchID, BranchName: item.BranchName, Date: item.Date, Status: item.Status, Reason: item.Reason, RuleOutcome: item.RuleOutcome, FirstCheckIn: item.FirstCheckIn, LastCheckOut: item.LastCheckOut, WorkedMinutes: item.WorkedMinutes, LateMinutes: item.LateMinutes, EarlyExitMinutes: item.EarlyExitMinutes, PunchCount: int32(len(item.AttendanceRecords))}
	if item.Policy != nil {
		row.PolicyName = &item.Policy.Name
		row.ScheduleType = &item.Policy.ScheduleType
	}
	if item.Roster != nil {
		row.WorkMode = &item.Roster.WorkMode
	} else {
		for _, attendance := range item.AttendanceRecords {
			if attendance != nil && attendance.WorkMode != nil {
				row.WorkMode = attendance.WorkMode
				break
			}
		}
	}
	return row
}

func accumulateAttendanceReportSummary(summary *domain.AttendanceReportSummary, row *domain.AttendanceReportRow) {
	summary.EmployeeDays++
	summary.ByStatus[row.Status]++
	if row.RuleOutcome != "" {
		summary.ByOutcome[row.RuleOutcome]++
	}
	summary.TotalWorkedMinutes += row.WorkedMinutes
	switch row.Status {
	case domain.AttendanceStatusPresent:
		summary.PresentDays++
	case domain.AttendanceStatusAbsent:
		summary.AbsentDays++
	case domain.AttendanceStatusLeave:
		summary.LeaveDays++
	case domain.AttendanceStatusHoliday:
		summary.HolidayDays++
	case domain.AttendanceStatusWeekoff:
		summary.WeekoffDays++
	case domain.AttendanceStatusIncomplete:
		summary.IncompleteDays++
	case domain.AttendanceStatusEmpty:
		summary.EmptyDays++
	}
	if row.RuleOutcome == domain.AttendanceRuleOutcomeLate || row.LateMinutes > 0 {
		summary.LateDays++
	}
	if row.RuleOutcome == domain.AttendanceRuleOutcomeHalfDay || row.RuleOutcome == domain.AttendanceRuleOutcomeMissingCheckout {
		summary.HalfDays++
	}
	if row.RuleOutcome == domain.AttendanceRuleOutcomeEarlyExit || row.EarlyExitMinutes > 0 {
		summary.EarlyExitDays++
	}
}

func accumulateDepartmentAttendance(items map[string]*domain.AttendanceDepartmentReport, row *domain.AttendanceReportRow) {
	name := "Unassigned"
	if row.DepartmentName != nil && *row.DepartmentName != "" {
		name = *row.DepartmentName
	}
	item := items[name]
	if item == nil {
		item = &domain.AttendanceDepartmentReport{DepartmentName: name}
		items[name] = item
	}
	item.EmployeeDays++
	item.TotalWorkedMinutes += row.WorkedMinutes
	if row.Status == domain.AttendanceStatusPresent {
		item.PresentDays++
	}
	if row.Status == domain.AttendanceStatusAbsent {
		item.AbsentDays++
	}
	if row.Status == domain.AttendanceStatusIncomplete {
		item.IncompleteDays++
	}
	if row.RuleOutcome == domain.AttendanceRuleOutcomeLate || row.LateMinutes > 0 {
		item.LateDays++
	}
}

func accumulateDailyAttendance(items map[string]*domain.AttendanceDailyTrend, row *domain.AttendanceReportRow) {
	key := row.Date.Format("2006-01-02")
	item := items[key]
	if item == nil {
		item = &domain.AttendanceDailyTrend{Date: row.Date}
		items[key] = item
	}
	item.EmployeeDays++
	item.TotalWorkedMinutes += row.WorkedMinutes
	if row.Status == domain.AttendanceStatusPresent {
		item.PresentDays++
	}
	if row.Status == domain.AttendanceStatusAbsent {
		item.AbsentDays++
	}
	if row.RuleOutcome == domain.AttendanceRuleOutcomeLate || row.LateMinutes > 0 {
		item.LateDays++
	}
}

func accumulateWorkModeAttendance(items map[string]*domain.AttendanceWorkModeReport, row *domain.AttendanceReportRow) {
	mode := "unassigned"
	if row.WorkMode != nil && *row.WorkMode != "" {
		mode = *row.WorkMode
	}
	item := items[mode]
	if item == nil {
		item = &domain.AttendanceWorkModeReport{WorkMode: mode}
		items[mode] = item
	}
	item.Days++
	item.WorkedMinutes += row.WorkedMinutes
}

func finalizeAttendanceReport(summary *domain.AttendanceReportSummary, departmentAgg map[string]*domain.AttendanceDepartmentReport, dailyAgg map[string]*domain.AttendanceDailyTrend, workModeAgg map[string]*domain.AttendanceWorkModeReport, report *domain.AttendanceReport) {
	if summary.EmployeeDays > 0 {
		summary.AverageWorkedMinutes = summary.TotalWorkedMinutes / summary.EmployeeDays
		summary.AttendanceRate = percent(summary.PresentDays, summary.EmployeeDays)
		summary.AbsenteeismRate = percent(summary.AbsentDays, summary.EmployeeDays)
		summary.LateRate = percent(summary.LateDays, summary.EmployeeDays)
	}
	for _, item := range departmentAgg {
		if item.EmployeeDays > 0 {
			item.AverageWorkedMinutes = item.TotalWorkedMinutes / item.EmployeeDays
			item.AttendanceRate = percent(item.PresentDays, item.EmployeeDays)
		}
		report.Departments = append(report.Departments, item)
	}
	for _, item := range dailyAgg {
		if item.EmployeeDays > 0 {
			item.AverageWorkedMinutes = item.TotalWorkedMinutes / item.EmployeeDays
		}
		report.DailyTrends = append(report.DailyTrends, item)
	}
	for _, item := range workModeAgg {
		if summary.EmployeeDays > 0 {
			item.SharePercent = percent(item.Days, summary.EmployeeDays)
		}
		report.WorkModes = append(report.WorkModes, item)
	}
	sort.Slice(report.Departments, func(i, j int) bool {
		return report.Departments[i].DepartmentName < report.Departments[j].DepartmentName
	})
	sort.Slice(report.DailyTrends, func(i, j int) bool { return report.DailyTrends[i].Date.Before(report.DailyTrends[j].Date) })
	sort.Slice(report.WorkModes, func(i, j int) bool { return report.WorkModes[i].Days > report.WorkModes[j].Days })
	report.LateEmployees = trimReportRows(report.LateEmployees)
	report.AbsenceEmployees = trimReportRows(report.AbsenceEmployees)
	report.ExceptionEmployees = trimReportRows(report.ExceptionEmployees)
}

func trimReportRows(rows []*domain.AttendanceReportRow) []*domain.AttendanceReportRow {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Date.Equal(rows[j].Date) {
			return rows[i].Firstname < rows[j].Firstname
		}
		return rows[i].Date.After(rows[j].Date)
	})
	if len(rows) > 10 {
		return rows[:10]
	}
	return rows
}

func percent(part int32, total int32) float64 {
	if total <= 0 {
		return 0
	}
	return float64(part) * 100 / float64(total)
}
