package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) ListAttendanceDailyStatuses(ctx context.Context, query ports.AttendanceStatusQuery) ([]*domain.AttendanceDailyStatus, error) {
	if query.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate attendance status tenant", err)
		return nil, err
	}
	day, err := parseAttendanceStatusDate(query.Date)
	if err != nil {
		s.logError("validate attendance status date", err, serviceTenantIDField(query.TenantID), serviceStringField("date", query.Date))
		return nil, err
	}
	employees, err := s.employees.ListEmployees(ctx, query.TenantID)
	if err != nil {
		s.logError("list attendance status employees", err, serviceTenantIDField(query.TenantID))
		return nil, err
	}
	if query.UserID != nil && *query.UserID != uuid.Nil {
		filtered := make([]*domain.EmployeeListItem, 0, 1)
		for _, employee := range employees {
			if employee != nil && employee.UserID == *query.UserID {
				filtered = append(filtered, employee)
				break
			}
		}
		employees = filtered
	}
	holidayItems, err := s.holidays.ListHolidaysByDateRange(ctx, query.TenantID, day, day)
	if err != nil {
		s.logError("list attendance status holidays", err, serviceTenantIDField(query.TenantID), serviceStringField("date", day.Format("2006-01-02")))
		return nil, err
	}
	items := make([]*domain.AttendanceDailyStatus, 0, len(employees))
	for _, employee := range employees {
		if employee == nil {
			continue
		}
		status, statusErr := s.resolveAttendanceDailyStatus(ctx, query.TenantID, employee, day, holidayItems)
		if statusErr != nil {
			s.logError("resolve attendance daily status", statusErr, serviceTenantIDField(query.TenantID), serviceStringField("user_id", employee.UserID.String()), serviceStringField("date", day.Format("2006-01-02")))
			return nil, statusErr
		}
		items = append(items, status)
	}
	return items, nil
}

func (s *TenantService) GetAttendanceStatusSummary(ctx context.Context, query ports.AttendanceStatusQuery) (*domain.AttendanceStatusSummary, error) {
	items, err := s.ListAttendanceDailyStatuses(ctx, query)
	if err != nil {
		return nil, err
	}
	day, _ := parseAttendanceStatusDate(query.Date)
	summary := &domain.AttendanceStatusSummary{Date: day, ByStatus: map[string]int32{}}
	for _, item := range items {
		if item == nil {
			continue
		}
		summary.TotalEmployees++
		summary.ByStatus[item.Status]++
		summary.TotalWorkedMinutes += item.WorkedMinutes
		switch item.Status {
		case domain.AttendanceStatusPresent:
			summary.Present++
		case domain.AttendanceStatusLeave:
			summary.Leave++
		case domain.AttendanceStatusAbsent:
			summary.Absent++
		case domain.AttendanceStatusHoliday:
			summary.Holiday++
		case domain.AttendanceStatusWeekoff:
			summary.Weekoff++
		case domain.AttendanceStatusIncomplete:
			summary.Incomplete++
		case domain.AttendanceStatusEmpty:
			summary.Empty++
		case domain.AttendanceStatusNotApplicable:
			summary.NotApplicable++
		}
	}
	return summary, nil
}

func (s *TenantService) resolveAttendanceDailyStatus(ctx context.Context, tenantID uuid.UUID, employee *domain.EmployeeListItem, day time.Time, holidays []*domain.Holiday) (*domain.AttendanceDailyStatus, error) {
	dateKey := day.Format("2006-01-02")
	attendances, err := s.attendances.ListAttendancesByUserDate(ctx, tenantID, employee.UserID, dateKey)
	if err != nil {
		return nil, err
	}
	leaves, err := s.leaveRequests.ListOverlappingLeaves(ctx, tenantID, employee.UserID, dateKey, dateKey)
	if err != nil {
		return nil, err
	}
	workingHour, err := s.ResolveWorkingHour(ctx, ports.ResolveWorkingHourCommand{TenantID: tenantID, BranchID: employee.BranchID, UserID: &employee.UserID, DayOfWeek: day.Weekday().String()})
	if err != nil {
		return nil, err
	}
	policy, err := s.attendancePolicies.ResolveAttendancePolicy(ctx, tenantID, employee.UserID, employee.DepartmentID, employee.BranchID, dateKey)
	if err != nil && !errors.Is(err, domain.ErrAttendancePolicyNotFound) {
		return nil, err
	}
	roster, err := s.attendanceRosters.GetAttendanceRosterByUserDate(ctx, tenantID, employee.UserID, dateKey)
	if err != nil && !errors.Is(err, domain.ErrAttendanceRosterNotFound) {
		return nil, err
	}
	status := &domain.AttendanceDailyStatus{
		TenantID:          tenantID,
		UserID:            employee.UserID,
		EmployeeID:        employee.ID,
		EmployeeCode:      employee.EmployeeCode,
		Firstname:         employee.Firstname,
		Lastname:          employee.Lastname,
		DepartmentID:      employee.DepartmentID,
		DepartmentName:    employee.DepartmentName,
		BranchID:          employee.BranchID,
		BranchName:        employee.BranchName,
		Date:              day,
		Status:            domain.AttendanceStatusEmpty,
		Reason:            "No attendance activity yet.",
		WorkingHour:       workingHour,
		Policy:            policy,
		Roster:            roster,
		AttendanceRecords: attendances,
	}
	status.FirstCheckIn, status.LastCheckOut, status.WorkedMinutes = attendancePunchStats(attendances)
	status.LateMinutes, status.EarlyExitMinutes, status.RuleOutcome = attendanceRuleOutcome(status.FirstCheckIn, status.LastCheckOut, status.WorkedMinutes, workingHour, policy, roster)
	if employee.JoiningDate != nil && day.Before(normalizeDate(*employee.JoiningDate)) {
		status.Status = domain.AttendanceStatusNotApplicable
		status.Reason = "Employee had not joined on this date."
		return status, nil
	}
	if employee.ResignationDate != nil && day.After(normalizeDate(*employee.ResignationDate)) {
		status.Status = domain.AttendanceStatusNotApplicable
		status.Reason = "Employee was no longer active on this date."
		return status, nil
	}
	if !employee.AttendanceRequired {
		status.Status = domain.AttendanceStatusNotApplicable
		status.Reason = "Attendance is not required for this designation."
		return status, nil
	}
	if leave := approvedLeaveForDate(leaves, day); leave != nil {
		status.Status = domain.AttendanceStatusLeave
		status.Reason = "Approved leave covers this date."
		status.Leave = leave
		return status, nil
	}
	if explicit := explicitAttendanceStatus(attendances); explicit != "" {
		status.Status = explicit
		status.Reason = "Explicit attendance status recorded."
		return status, nil
	}
	if status.FirstCheckIn != nil {
		if status.LastCheckOut == nil && attendanceCheckoutRequired(day, status, time.Now().UTC()) {
			status.Status = domain.AttendanceStatusIncomplete
			status.RuleOutcome = domain.AttendanceRuleOutcomeMissingCheckout
			status.Reason = "Check-out is missing; attendance needs regularisation."
			return status, nil
		}
		if status.LastCheckOut == nil && status.RuleOutcome == domain.AttendanceRuleOutcomeMissingCheckout {
			status.RuleOutcome = ""
		}
		status.Status = domain.AttendanceStatusPresent
		if status.LastCheckOut == nil {
			status.Reason = "Check-in recorded; check-out is still pending."
		} else {
			status.Reason = "Check-in recorded."
		}
		if status.RuleOutcome != "" && status.RuleOutcome != domain.AttendanceRuleOutcomeOnTime {
			status.Reason = fmt.Sprintf("Check-in recorded; rule outcome: %s.", status.RuleOutcome)
		}
		return status, nil
	}
	if holiday := holidayForEmployee(holidays, employee.BranchID); holiday != nil {
		status.Status = domain.AttendanceStatusHoliday
		status.Reason = "Holiday calendar applies."
		status.Holiday = holiday
		return status, nil
	}
	if workingHour != nil && !workingHour.IsWorkingDay {
		status.Status = domain.AttendanceStatusWeekoff
		status.Reason = "Resolved working-hours mark this day as weekoff."
		return status, nil
	}
	if day.Before(normalizeDate(time.Now().UTC())) {
		status.Status = domain.AttendanceStatusAbsent
		status.Reason = "Past working day without attendance, leave, holiday, or weekoff."
		return status, nil
	}
	status.Status = domain.AttendanceStatusEmpty
	status.Reason = "Attendance is still open for this date."
	return status, nil
}

func attendanceRuleOutcome(firstIn *time.Time, lastOut *time.Time, workedMinutes int32, workingHour *domain.WorkingHour, policy *domain.AttendancePolicy, roster *domain.AttendanceRoster) (int32, int32, string) {
	if firstIn == nil {
		return 0, 0, ""
	}
	if lastOut == nil {
		return 0, 0, domain.AttendanceRuleOutcomeMissingCheckout
	}
	startMinutes, endMinutes, ok := scheduleMinutes(workingHour, roster)
	if !ok {
		return 0, 0, domain.AttendanceRuleOutcomeOnTime
	}
	firstMinutes := int32(firstIn.In(time.UTC).Hour()*60 + firstIn.In(time.UTC).Minute())
	late := firstMinutes - startMinutes
	graceLate := int32(0)
	graceEarly := int32(0)
	minHalf := int32(240)
	minFull := int32(420)
	if policy != nil {
		graceLate = policy.GraceLateMinutes
		graceEarly = policy.GraceEarlyMinutes
		minHalf = policy.MinHalfDayMinutes
		minFull = policy.MinFullDayMinutes
	}
	if late < 0 {
		late = 0
	}
	early := int32(0)
	if lastOut != nil {
		lastMinutes := int32(lastOut.In(time.UTC).Hour()*60 + lastOut.In(time.UTC).Minute())
		early = endMinutes - lastMinutes
		if early < 0 {
			early = 0
		}
	}
	outcome := domain.AttendanceRuleOutcomeOnTime
	if policy != nil && policy.AbsentLateAfterMinutes != nil && late > *policy.AbsentLateAfterMinutes {
		outcome = domain.AttendanceRuleOutcomeAbsent
	} else if policy != nil && policy.HalfDayLateAfterMinutes != nil && late > *policy.HalfDayLateAfterMinutes {
		outcome = domain.AttendanceRuleOutcomeHalfDay
	} else if workedMinutes > 0 && workedMinutes < minHalf {
		outcome = domain.AttendanceRuleOutcomeAbsent
	} else if workedMinutes > 0 && workedMinutes < minFull {
		outcome = domain.AttendanceRuleOutcomeHalfDay
	} else if late > graceLate {
		outcome = domain.AttendanceRuleOutcomeLate
	} else if early > graceEarly {
		outcome = domain.AttendanceRuleOutcomeEarlyExit
	}
	return late, early, outcome
}

func attendanceCheckoutRequired(day time.Time, status *domain.AttendanceDailyStatus, now time.Time) bool {
	if status == nil || status.FirstCheckIn == nil || status.LastCheckOut != nil {
		return false
	}
	today := normalizeDate(now)
	attendanceDay := normalizeDate(day)
	if attendanceDay.Before(today) {
		return true
	}
	if attendanceDay.After(today) {
		return false
	}
	_, endMinutes, ok := scheduleMinutes(status.WorkingHour, status.Roster)
	if !ok {
		return false
	}
	grace := int32(0)
	if status.Policy != nil {
		grace = status.Policy.GraceEarlyMinutes
	}
	nowMinutes := int32(now.In(time.UTC).Hour()*60 + now.In(time.UTC).Minute())
	return nowMinutes > endMinutes+grace
}

func scheduleMinutes(workingHour *domain.WorkingHour, roster *domain.AttendanceRoster) (int32, int32, bool) {
	if roster != nil && roster.StartTime != nil && roster.EndTime != nil {
		start, err1 := hourMinuteToInt(*roster.StartTime)
		end, err2 := hourMinuteToInt(*roster.EndTime)
		return start, end, err1 == nil && err2 == nil
	}
	if workingHour != nil && workingHour.IsWorkingDay {
		start, err1 := hourMinuteToInt(workingHour.StartTime)
		end, err2 := hourMinuteToInt(workingHour.EndTime)
		return start, end, err1 == nil && err2 == nil
	}
	return 0, 0, false
}

func hourMinuteToInt(value string) (int32, error) {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return 0, err
	}
	return int32(parsed.Hour()*60 + parsed.Minute()), nil
}

func parseAttendanceStatusDate(value string) (time.Time, error) {
	if value == "" {
		return normalizeDate(time.Now().UTC()), nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, domain.ErrInvalidAttendanceDate
	}
	return normalizeDate(parsed), nil
}

func normalizeDate(value time.Time) time.Time {
	return time.Date(value.UTC().Year(), value.UTC().Month(), value.UTC().Day(), 0, 0, 0, 0, time.UTC)
}

func approvedLeaveForDate(items []*domain.Leave, day time.Time) *domain.Leave {
	for _, item := range items {
		if item == nil || item.Status != domain.LeaveStatusApproved {
			continue
		}
		if !day.Before(normalizeDate(item.StartDate)) && !day.After(normalizeDate(item.EndDate)) {
			return item
		}
	}
	return nil
}

func explicitAttendanceStatus(items []*domain.Attendance) string {
	for _, item := range items {
		if item == nil || item.Status == nil {
			continue
		}
		switch *item.Status {
		case domain.AttendanceStatusAbsent, domain.AttendanceStatusHoliday, domain.AttendanceStatusLeave, domain.AttendanceStatusWeekoff:
			return *item.Status
		case domain.AttendanceStatusPresent:
			if item.Type == nil {
				return *item.Status
			}
		}
	}
	return ""
}

func attendancePunchStats(items []*domain.Attendance) (*time.Time, *time.Time, int32) {
	var firstCheckIn *time.Time
	var lastCheckOut *time.Time
	var openCheckIn *time.Time
	var workedMinutes int32
	for _, item := range items {
		if item == nil || item.Type == nil || item.Time == nil {
			continue
		}
		t := item.Time.UTC()
		switch *item.Type {
		case domain.AttendanceCheckin:
			if firstCheckIn == nil || t.Before(*firstCheckIn) {
				firstCheckIn = &t
			}
			openCheckIn = &t
		case domain.AttendanceCheckout:
			if lastCheckOut == nil || t.After(*lastCheckOut) {
				lastCheckOut = &t
			}
			if openCheckIn != nil && t.After(*openCheckIn) {
				workedMinutes += int32(t.Sub(*openCheckIn).Minutes())
				openCheckIn = nil
			}
		}
	}
	return firstCheckIn, lastCheckOut, workedMinutes
}

func holidayForEmployee(items []*domain.Holiday, branchID *uuid.UUID) *domain.Holiday {
	for _, item := range items {
		if item == nil || item.IsOptional {
			continue
		}
		if item.BranchID == nil || branchID == nil || *item.BranchID == *branchID {
			return item
		}
	}
	return nil
}
