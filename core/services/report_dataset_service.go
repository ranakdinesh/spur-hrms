package services

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) BuildReportDataset(ctx context.Context, query ports.ReportDatasetQuery) (*domain.ReportDataset, error) {
	code := strings.TrimSpace(query.ReportCode)
	if code == "" && query.ReportID != uuid.Nil {
		item, err := s.reporting.GetReportCatalogItem(ctx, query.TenantID, query.ReportID)
		if err != nil {
			s.logError("get report catalog item for dataset", err, serviceTenantIDField(query.TenantID), serviceStringField("report_id", query.ReportID.String()))
			return nil, err
		}
		code = item.ReportCode
	}
	if code == "" {
		return nil, domain.ErrReportNotFound
	}
	var dataset *domain.ReportDataset
	var err error
	switch code {
	case "workforce.composition", "workforce.movement", "workforce.data_quality":
		dataset, err = s.buildWorkforceDataset(ctx, query, code)
	case "attendance.exceptions", "attendance.health", "time.attendance_health", "time.exception_watch":
		dataset, err = s.buildAttendanceDataset(ctx, query, code)
	case "leave.liability", "leave.utilization", "time.leave_utilization":
		dataset, err = s.buildLeaveDataset(ctx, query, code)
	case "payroll.readiness", "payroll.pay_group_readiness":
		dataset, err = s.buildPayrollReadinessDataset(ctx, query, code)
	case "payroll.consolidated_salary_sheet", "payroll.cost":
		dataset, err = s.buildSalarySheetDataset(ctx, query, code)
	case "payroll.reconciliation":
		dataset, err = s.buildPayrollReconciliationDataset(ctx, query, code)
	case "compliance.readiness", "payroll.compliance":
		dataset, err = s.buildComplianceDataset(ctx, query, code)
	case "recruitment.funnel", "recruitment.source_effectiveness":
		dataset, err = s.buildRecruitmentDataset(ctx, query, code)
	case "onboarding.completion":
		dataset, err = s.buildOnboardingLifecycleDataset(ctx, query, code)
	case "probation.due":
		dataset, err = s.buildProbationDueDataset(ctx, query, code)
	case "exit.pipeline", "exit.readiness":
		dataset, err = s.buildExitLifecycleDataset(ctx, query, code)
	default:
		return nil, domain.ErrReportNotFound
	}
	if err != nil {
		return nil, err
	}
	dataset.Branding = s.reportBranding(ctx, query.TenantID)
	return dataset, nil
}

func (s *TenantService) ExportReportDataset(ctx context.Context, query ports.ReportDatasetQuery, format string) (*ports.ReportDownload, error) {
	dataset, err := s.BuildReportDataset(ctx, query)
	if err != nil {
		return nil, err
	}
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case domain.ReportExportPDF:
		content, err := renderReportPDF(dataset)
		if err != nil {
			return nil, err
		}
		return &ports.ReportDownload{Content: content, FileName: reportFileName(dataset, "pdf"), ContentType: "application/pdf"}, nil
	case domain.ReportExportXLSX, "excel":
		content, err := renderReportXLSX(dataset)
		if err != nil {
			return nil, err
		}
		return &ports.ReportDownload{Content: content, FileName: reportFileName(dataset, "xlsx"), ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, nil
	default:
		return nil, domain.ErrReportExportJobInvalid
	}
}

func (s *TenantService) buildWorkforceDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	employees, err := s.employees.ListEmployees(ctx, query.TenantID)
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0, len(employees))
	active, inactive := 0, 0
	for _, employee := range employees {
		status := "active"
		if employee.Inactive {
			status = "inactive"
			inactive++
		} else {
			active++
		}
		rows = append(rows, []string{
			valueFromPtr(employee.EmployeeCode),
			employeeName(employee.Firstname, employee.Lastname),
			valueFromPtr(employee.Email),
			valueFromPtr(employee.BranchName),
			valueFromPtr(employee.DepartmentName),
			valueFromPtr(employee.DesignationName),
			valueFromPtr(employee.EmploymentTypeName),
			dateString(employee.JoiningDate),
			status,
		})
	}
	return &domain.ReportDataset{
		TenantID:    query.TenantID,
		ReportCode:  code,
		Title:       "Workforce Composition",
		Description: "Headcount, employee mix, reporting attributes, and data quality fields.",
		Columns:     reportColumns("employee_code", "Employee Code", "employee_name", "Employee", "email", "Email", "branch", "Branch", "department", "Department", "designation", "Designation", "employment_type", "Employment Type", "joining_date", "Joining Date", "status", "Status"),
		Rows:        rows,
		Summary:     []domain.ReportMetric{{Label: "Employees", Value: strconv.Itoa(len(employees))}, {Label: "Active", Value: strconv.Itoa(active)}, {Label: "Inactive", Value: strconv.Itoa(inactive)}},
	}, nil
}

func (s *TenantService) buildAttendanceDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	start, end, err := parseReportRange(query)
	if err != nil {
		return nil, err
	}
	report, err := s.GetAttendanceReport(ctx, ports.AttendanceReportQuery{TenantID: query.TenantID, StartDate: start, EndDate: end})
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0, len(report.Rows))
	for _, row := range report.Rows {
		if row.Status != domain.AttendanceStatusIncomplete && row.Status != domain.AttendanceStatusAbsent && row.RuleOutcome != "" && row.RuleOutcome != domain.AttendanceRuleOutcomeOnTime {
			rows = append(rows, []string{valueFromPtr(row.EmployeeCode), employeeName(row.Firstname, row.Lastname), valueFromPtr(row.DepartmentName), row.Date.Format("2006-01-02"), row.Status, row.RuleOutcome, minutesLabel(row.WorkedMinutes), minutesLabel(row.LateMinutes), minutesLabel(row.EarlyExitMinutes), valueFromPtr(row.WorkMode)})
		}
	}
	return &domain.ReportDataset{
		TenantID:    query.TenantID,
		ReportCode:  code,
		Title:       "Attendance Exceptions",
		Description: "Late, absent, incomplete, short-hours, and rule exception rows for payroll and HR review.",
		PeriodLabel: start + " to " + end,
		Columns:     reportColumns("employee_code", "Employee Code", "employee_name", "Employee", "department", "Department", "date", "Date", "status", "Status", "rule_outcome", "Rule Outcome", "worked", "Worked", "late", "Late", "early_exit", "Early Exit", "work_mode", "Work Mode"),
		Rows:        rows,
		Summary:     []domain.ReportMetric{{Label: "Employee Days", Value: fmt.Sprint(report.Summary.EmployeeDays)}, {Label: "Attendance Rate", Value: fmt.Sprintf("%.1f%%", report.Summary.AttendanceRate)}, {Label: "Late Days", Value: fmt.Sprint(report.Summary.LateDays)}, {Label: "Pending Requests", Value: fmt.Sprint(report.Summary.PendingRequests)}},
	}, nil
}

func (s *TenantService) buildLeaveDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	start, end, _ := parseOptionalReportRange(query)
	filter := domain.LeaveReportFilter{TenantID: query.TenantID, FYID: query.FYID, StartDate: start, EndDate: end}
	rowsData, err := s.ListLeaveReportRows(ctx, filter)
	if err != nil {
		return nil, err
	}
	summary, _ := s.GetLeaveReportSummary(ctx, filter)
	rows := make([][]string, 0, len(rowsData))
	for _, row := range rowsData {
		rows = append(rows, []string{valueFromPtr(row.EmployeeCode), employeeName(row.Firstname, row.Lastname), valueFromPtr(row.DepartmentName), valueFromPtr(row.LeaveTypeName), row.StartDate.Format("2006-01-02"), row.EndDate.Format("2006-01-02"), fmt.Sprintf("%.2f", row.Days), row.Status, valueFromPtr(row.Reason)})
	}
	metrics := []domain.ReportMetric{{Label: "Requests", Value: strconv.Itoa(len(rowsData))}}
	if summary != nil {
		metrics = []domain.ReportMetric{{Label: "Requests", Value: fmt.Sprint(summary.TotalRequests)}, {Label: "Days", Value: fmt.Sprintf("%.2f", summary.TotalDays)}, {Label: "Pending", Value: fmt.Sprint(summary.PendingCount)}, {Label: "Approved", Value: fmt.Sprint(summary.ApprovedCount)}}
	}
	return &domain.ReportDataset{TenantID: query.TenantID, ReportCode: code, Title: "Leave Liability and Utilization", Description: "Leave usage, pending approval exposure, and balances affecting payroll and staffing.", PeriodLabel: reportPeriodLabel(start, end), Columns: reportColumns("employee_code", "Employee Code", "employee_name", "Employee", "department", "Department", "leave_type", "Leave Type", "start_date", "Start", "end_date", "End", "days", "Days", "status", "Status", "reason", "Reason"), Rows: rows, Summary: metrics}, nil
}

func (s *TenantService) buildPayrollReadinessDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	month, year := reportMonthYear(query)
	reconcile, err := s.ListPayrollReconciliationRows(ctx, query.TenantID, month, year)
	if err != nil {
		return nil, err
	}
	start := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	blockers, _ := s.attendanceExceptionWorkflows.ListPayrollBlockingAttendanceRequests(ctx, query.TenantID, start.Format("2006-01-02"), start.AddDate(0, 1, -1).Format("2006-01-02"))
	rows := make([][]string, 0, len(reconcile)+len(blockers))
	exceptions := 0
	for _, row := range reconcile {
		if row.ReconciliationStatus != "ok" {
			exceptions++
		}
		rows = append(rows, []string{valueFromPtr(row.EmployeeCode), employeeName(row.Firstname, row.Lastname), valueFromPtr(row.DepartmentName), row.ReconciliationStatus, ptrInt32(row.PresentDays), ptrInt32(row.AbsentDays), ptrFloatString(row.LWPDays), ptrFloatString(row.NetSalary)})
	}
	for _, blocker := range blockers {
		rows = append(rows, []string{"", blocker.UserID.String(), "", "attendance_exception_pending", blocker.Date.Format("2006-01-02"), blocker.RequestType, blocker.Status, valueFromPtr(blocker.Reason)})
	}
	return &domain.ReportDataset{TenantID: query.TenantID, ReportCode: code, Title: "Payroll Readiness", Description: "Attendance, LOP, payslip and exception blockers before payroll finalization.", PeriodLabel: fmt.Sprintf("%02d/%04d", month, year), Columns: reportColumns("employee_code", "Employee Code", "employee_name", "Employee", "department", "Department", "status", "Readiness Status", "present_days", "Present", "absent_days", "Absent", "lwp_days", "LWP", "net_salary", "Net Salary"), Rows: rows, Summary: []domain.ReportMetric{{Label: "Employees Checked", Value: strconv.Itoa(len(reconcile))}, {Label: "Reconciliation Exceptions", Value: strconv.Itoa(exceptions)}, {Label: "Attendance Blockers", Value: strconv.Itoa(len(blockers))}}}, nil
}

func (s *TenantService) buildSalarySheetDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	month, year := reportMonthYear(query)
	data, err := s.ListConsolidatedSalarySheet(ctx, query.TenantID, month, year)
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0, len(data))
	net := 0.0
	for _, row := range data {
		net += row.NetSalary
		rows = append(rows, []string{valueFromPtr(row.EmployeeCode), employeeName(row.Firstname, row.Lastname), valueFromPtr(row.Email), valueFromPtr(row.BranchName), valueFromPtr(row.DepartmentName), moneyString(row.GrossSalary), moneyString(row.TotalEarnings), moneyString(row.TotalDeductions), moneyString(row.AbsentDeduction), moneyString(row.NetSalary), fmt.Sprint(row.PresentDays), fmt.Sprint(row.AbsentDays), moneyString(row.LWPDays)})
	}
	return &domain.ReportDataset{TenantID: query.TenantID, ReportCode: code, Title: "Consolidated Salary Sheet", Description: "Month-wise payroll output for Finance, HR and audit sharing.", PeriodLabel: fmt.Sprintf("%02d/%04d", month, year), Columns: reportColumns("employee_code", "Employee Code", "employee_name", "Employee", "email", "Email", "branch", "Branch", "department", "Department", "gross_salary", "Gross", "total_earnings", "Earnings", "total_deductions", "Deductions", "absent_deduction", "Absent Deduction", "net_salary", "Net", "present_days", "Present", "absent_days", "Absent", "lwp_days", "LWP"), Rows: rows, Summary: []domain.ReportMetric{{Label: "Rows", Value: strconv.Itoa(len(rows))}, {Label: "Net Payroll", Value: moneyString(net)}}}, nil
}

func (s *TenantService) buildPayrollReconciliationDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	month, year := reportMonthYear(query)
	data, err := s.ListPayrollReconciliationRows(ctx, query.TenantID, month, year)
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0, len(data))
	exceptions := 0
	for _, row := range data {
		if row.ReconciliationStatus != "ok" {
			exceptions++
		}
		rows = append(rows, []string{valueFromPtr(row.EmployeeCode), employeeName(row.Firstname, row.Lastname), valueFromPtr(row.Email), valueFromPtr(row.BranchName), valueFromPtr(row.DepartmentName), ptrInt32(row.PresentDays), ptrInt32(row.AbsentDays), ptrFloatString(row.LWPDays), ptrFloatString(row.NetSalary), row.ReconciliationStatus})
	}
	return &domain.ReportDataset{TenantID: query.TenantID, ReportCode: code, Title: "Payroll Attendance Reconciliation", Description: "Presents, absences, LOP and payslip readiness for payroll review.", PeriodLabel: fmt.Sprintf("%02d/%04d", month, year), Columns: reportColumns("employee_code", "Employee Code", "employee_name", "Employee", "email", "Email", "branch", "Branch", "department", "Department", "present_days", "Present", "absent_days", "Absent", "lwp_days", "LWP", "net_salary", "Net", "status", "Status"), Rows: rows, Summary: []domain.ReportMetric{{Label: "Rows", Value: strconv.Itoa(len(rows))}, {Label: "Exceptions", Value: strconv.Itoa(exceptions)}}}, nil
}

func (s *TenantService) buildComplianceDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	employees, err := s.employees.ListEmployees(ctx, query.TenantID)
	if err != nil {
		return nil, err
	}
	rows := [][]string{}
	for _, employee := range employees {
		if employee.Email == nil || employee.Mobile == nil || employee.DepartmentID == nil || employee.BranchID == nil || employee.ProbationStatus == domain.EmployeeProbationProbation {
			rows = append(rows, []string{valueFromPtr(employee.EmployeeCode), employeeName(employee.Firstname, employee.Lastname), valueFromPtr(employee.DepartmentName), valueFromPtr(employee.BranchName), complianceRisk(employee), employee.ProbationStatus})
		}
	}
	return &domain.ReportDataset{TenantID: query.TenantID, ReportCode: code, Title: "Compliance Readiness", Description: "Employee statutory, document, probation and payroll compliance exposure.", Columns: reportColumns("employee_code", "Employee Code", "employee_name", "Employee", "department", "Department", "branch", "Branch", "risk_area", "Risk Area", "probation_status", "Probation"), Rows: rows, Summary: []domain.ReportMetric{{Label: "Employees Reviewed", Value: strconv.Itoa(len(employees))}, {Label: "Risk Rows", Value: strconv.Itoa(len(rows))}}}, nil
}

func (s *TenantService) buildRecruitmentDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	applications, err := s.candidates.ListCandidateApplications(ctx, domain.CandidateApplicationFilter{TenantID: query.TenantID, Limit: 10000})
	if err != nil {
		return nil, err
	}
	requisitions, _ := s.jobRequisitions.ListJobRequisitions(ctx, domain.JobRequisitionFilter{TenantID: query.TenantID, Limit: 10000})
	interviews, _ := s.candidates.ListInterviewRounds(ctx, domain.InterviewRoundFilter{TenantID: query.TenantID, Limit: 10000})
	if code == "recruitment.source_effectiveness" {
		return recruitmentSourceDataset(query, code, applications, interviews)
	}
	rows := make([][]string, 0, len(applications))
	stageCounts := map[string]int{}
	offerCount, hiredCount, rejectedCount, totalTimeToFillDays, filledCount := 0, 0, 0, 0, 0
	for _, app := range applications {
		stageCounts[app.Status]++
		switch app.Status {
		case domain.CandidateApplicationStatusOffered:
			offerCount++
		case domain.CandidateApplicationStatusHired:
			hiredCount++
			totalTimeToFillDays += daysBetween(app.AppliedAt, app.StatusChangedAt)
			filledCount++
		case domain.CandidateApplicationStatusRejected:
			rejectedCount++
		}
		rows = append(rows, []string{
			valueFromPtr(app.JobPostingTitle),
			candidateApplicationName(app),
			valueFromPtr(app.CandidateEmail),
			valueFromPtr(app.Source),
			app.Status,
			strconv.Itoa(app.DaysInStage),
			dateString(&app.AppliedAt),
			dateString(&app.StatusChangedAt),
			valueFromPtr(app.RejectionReason),
		})
	}
	openRequisitions := 0
	for _, req := range requisitions {
		if !req.Inactive && req.Status != "Closed" && req.Status != "Rejected" && req.Status != "Cancelled" {
			openRequisitions++
		}
	}
	return &domain.ReportDataset{
		TenantID:    query.TenantID,
		ReportCode:  code,
		Title:       "Recruitment Funnel",
		Description: "Requisition, candidate-stage aging, source, interview, offer and joining conversion drilldown.",
		Columns:     reportColumns("job_position", "Job Position", "candidate", "Candidate", "email", "Email", "source", "Source", "stage", "Stage", "age_days", "Age Days", "applied_at", "Applied", "stage_changed_at", "Stage Changed", "reason", "Reason"),
		Rows:        rows,
		Summary: []domain.ReportMetric{
			{Label: "Open Requisitions", Value: strconv.Itoa(openRequisitions)},
			{Label: "Applications", Value: strconv.Itoa(len(applications))},
			{Label: "Interviews", Value: strconv.Itoa(len(interviews))},
			{Label: "Offers", Value: strconv.Itoa(offerCount)},
			{Label: "Hired", Value: strconv.Itoa(hiredCount)},
			{Label: "Rejected", Value: strconv.Itoa(rejectedCount)},
			{Label: "Avg Time To Fill", Value: averageDaysLabel(totalTimeToFillDays, filledCount)},
			{Label: "Stage Mix", Value: stageMixLabel(stageCounts)},
		},
	}, nil
}

func recruitmentSourceDataset(query ports.ReportDatasetQuery, code string, applications []*domain.CandidateApplication, interviews []*domain.InterviewRound) (*domain.ReportDataset, error) {
	type sourceStats struct {
		applications    int
		screening       int
		interview       int
		offered         int
		hired           int
		rejected        int
		withdrawn       int
		stageAgeDays    int
		timeToFillDays  int
		filledCount     int
		interviewRounds int
	}
	stats := map[string]*sourceStats{}
	applicationSource := map[uuid.UUID]string{}
	for _, app := range applications {
		source := firstNonEmptyString(valueFromPtr(app.Source), "Unknown")
		row := stats[source]
		if row == nil {
			row = &sourceStats{}
			stats[source] = row
		}
		applicationSource[app.ID] = source
		row.applications++
		row.stageAgeDays += app.DaysInStage
		switch app.Status {
		case domain.CandidateApplicationStatusScreening:
			row.screening++
		case domain.CandidateApplicationStatusInterview:
			row.interview++
		case domain.CandidateApplicationStatusOffered:
			row.offered++
		case domain.CandidateApplicationStatusHired:
			row.hired++
			row.timeToFillDays += daysBetween(app.AppliedAt, app.StatusChangedAt)
			row.filledCount++
		case domain.CandidateApplicationStatusRejected:
			row.rejected++
		case domain.CandidateApplicationStatusWithdrawn:
			row.withdrawn++
		}
	}
	for _, round := range interviews {
		source := firstNonEmptyString(applicationSource[round.ApplicationID], "Unknown")
		row := stats[source]
		if row == nil {
			row = &sourceStats{}
			stats[source] = row
		}
		row.interviewRounds++
	}
	rows := make([][]string, 0, len(stats))
	for source, row := range stats {
		rows = append(rows, []string{
			source,
			strconv.Itoa(row.applications),
			strconv.Itoa(row.screening),
			strconv.Itoa(row.interview),
			strconv.Itoa(row.interviewRounds),
			strconv.Itoa(row.offered),
			strconv.Itoa(row.hired),
			strconv.Itoa(row.rejected),
			strconv.Itoa(row.withdrawn),
			percentString(row.hired, row.applications),
			percentString(row.hired, row.offered+row.hired),
			averageDaysLabel(row.stageAgeDays, row.applications),
			averageDaysLabel(row.timeToFillDays, row.filledCount),
		})
	}
	return &domain.ReportDataset{
		TenantID:    query.TenantID,
		ReportCode:  code,
		Title:       "Recruitment Source Effectiveness",
		Description: "Source/channel conversion, interview throughput, offer acceptance proxy, joining conversion, and stage aging.",
		Columns:     reportColumns("source", "Source", "applications", "Applications", "screening", "Screening", "interview_stage", "Interview Stage", "interview_rounds", "Interview Rounds", "offered", "Offered", "hired", "Hired", "rejected", "Rejected", "withdrawn", "Withdrawn", "joining_conversion", "Joining Conversion", "offer_acceptance", "Offer Acceptance", "avg_stage_age", "Avg Stage Age", "avg_time_to_fill", "Avg Time To Fill"),
		Rows:        rows,
		Summary:     []domain.ReportMetric{{Label: "Sources", Value: strconv.Itoa(len(stats))}, {Label: "Applications", Value: strconv.Itoa(len(applications))}, {Label: "Interviews", Value: strconv.Itoa(len(interviews))}},
	}, nil
}

func (s *TenantService) buildOnboardingLifecycleDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	page, err := s.candidateOnboardings.ListCandidateOnboardings(ctx, domain.CandidateOnboardingFilter{TenantID: query.TenantID, Limit: 10000})
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0, len(page.Items))
	completed, overdue, missingRequired := 0, 0, 0
	for _, item := range page.Items {
		if item.OnboardingStatus == domain.OnboardStatusCompleted {
			completed++
		}
		if item.OverdueTasks > 0 {
			overdue++
		}
		missing := item.RequiredTasks - item.CompletedRequiredTasks
		if missing > 0 {
			missingRequired += int(missing)
		}
		rows = append(rows, []string{
			candidateOnboardingName(item),
			valueFromPtr(item.CandidateEmail),
			valueFromPtr(item.WorkflowName),
			item.OnboardingStatus,
			fmt.Sprint(item.ProgressPercentage) + "%",
			fmt.Sprint(item.CompletedTasks),
			fmt.Sprint(item.TotalTasks),
			fmt.Sprint(item.CompletedRequiredTasks),
			fmt.Sprint(item.RequiredTasks),
			fmt.Sprint(item.OverdueTasks),
			dateString(item.StartedAt),
			dateString(item.CompletedAt),
		})
	}
	return &domain.ReportDataset{
		TenantID:    query.TenantID,
		ReportCode:  code,
		Title:       "Onboarding Completion",
		Description: "Joining conversion handoff, onboarding task completion, missing required documents/tasks, and overdue onboarding work.",
		Columns:     reportColumns("candidate", "Candidate", "email", "Email", "workflow", "Workflow", "status", "Status", "progress", "Progress", "completed_tasks", "Completed Tasks", "total_tasks", "Total Tasks", "required_completed", "Required Done", "required_tasks", "Required Tasks", "overdue_tasks", "Overdue Tasks", "started_at", "Started", "completed_at", "Completed"),
		Rows:        rows,
		Summary:     []domain.ReportMetric{{Label: "Onboardings", Value: strconv.Itoa(len(page.Items))}, {Label: "Completed", Value: strconv.Itoa(completed)}, {Label: "Overdue", Value: strconv.Itoa(overdue)}, {Label: "Missing Required", Value: strconv.Itoa(missingRequired)}},
	}, nil
}

func (s *TenantService) buildProbationDueDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	employees, err := s.employees.ListEmployees(ctx, query.TenantID)
	if err != nil {
		return nil, err
	}
	today := time.Now().UTC()
	dueUntil := today.AddDate(0, 0, 45)
	rows := [][]string{}
	overdue, dueSoon := 0, 0
	for _, employee := range employees {
		if employee.Inactive || employee.ProbationStatus != domain.EmployeeProbationProbation || employee.ProbationEndDate == nil {
			continue
		}
		if employee.ProbationEndDate.After(dueUntil) {
			continue
		}
		status := "due_soon"
		if employee.ProbationEndDate.Before(today) {
			status = "overdue"
			overdue++
		} else {
			dueSoon++
		}
		rows = append(rows, []string{
			valueFromPtr(employee.EmployeeCode),
			employeeName(employee.Firstname, employee.Lastname),
			valueFromPtr(employee.DepartmentName),
			valueFromPtr(employee.BranchName),
			dateString(employee.JoiningDate),
			dateString(employee.ProbationStartDate),
			dateString(employee.ProbationEndDate),
			fmt.Sprint(employee.ProbationDurationDays),
			status,
		})
	}
	return &domain.ReportDataset{
		TenantID:    query.TenantID,
		ReportCode:  code,
		Title:       "Probation Confirmations Due",
		Description: "Employees requiring probation confirmation, extension, or payroll-status review.",
		Columns:     reportColumns("employee_code", "Employee Code", "employee_name", "Employee", "department", "Department", "branch", "Branch", "joining_date", "Joining Date", "probation_start", "Probation Start", "probation_end", "Probation End", "duration_days", "Duration", "status", "Status"),
		Rows:        rows,
		Summary:     []domain.ReportMetric{{Label: "Due Rows", Value: strconv.Itoa(len(rows))}, {Label: "Overdue", Value: strconv.Itoa(overdue)}, {Label: "Due Soon", Value: strconv.Itoa(dueSoon)}},
	}, nil
}

func (s *TenantService) buildExitLifecycleDataset(ctx context.Context, query ports.ReportDatasetQuery, code string) (*domain.ReportDataset, error) {
	page, err := s.employeeExits.ListEmployeeExitRequests(ctx, domain.EmployeeExitFilter{TenantID: query.TenantID, Limit: 10000})
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0, len(page.Items))
	pipeline, ready, blocked, overdueNotice := 0, 0, 0, 0
	today := time.Now().UTC()
	reasonCounts := map[string]int{}
	for _, item := range page.Items {
		if item.Status != domain.EmployeeExitStatusCompleted && item.Status != domain.EmployeeExitStatusCanceled && item.Status != domain.EmployeeExitStatusRejected {
			pipeline++
		}
		if item.BlockedTasks > 0 {
			blocked++
		}
		if exitReady(item) {
			ready++
		}
		if item.LastWorkingDate.Before(today) && item.Status != domain.EmployeeExitStatusCompleted && item.Status != domain.EmployeeExitStatusCanceled {
			overdueNotice++
		}
		reasonCounts[firstNonEmptyString(valueFromPtr(item.Reason), item.ExitType)]++
		rows = append(rows, []string{
			valueFromPtr(item.EmployeeCode),
			employeeName(valueFromPtr(item.EmployeeFirstname), item.EmployeeLastname),
			valueFromPtr(item.DepartmentName),
			valueFromPtr(item.BranchName),
			item.Status,
			item.ExitType,
			valueFromPtr(item.Reason),
			dateString(item.ResignationDate),
			dateString(item.NoticeStartDate),
			item.LastWorkingDate.Format("2006-01-02"),
			item.FinalSettlementStatus,
			item.AssetClearanceStatus,
			item.HandoverStatus,
			item.AccessRevocationStatus,
			item.ExitInterviewStatus,
			fmt.Sprint(item.CompletedTasks),
			fmt.Sprint(item.TotalTasks),
			fmt.Sprint(item.BlockedTasks),
			exitReadinessLabel(item),
		})
	}
	title := "Exit Pipeline"
	if code == "exit.readiness" {
		title = "Exit and F&F Readiness"
	}
	return &domain.ReportDataset{
		TenantID:    query.TenantID,
		ReportCode:  code,
		Title:       title,
		Description: "Resignation pipeline, notice-period status, exit task aging, asset/handover/access/F&F readiness, and exit-reason analytics.",
		Columns:     reportColumns("employee_code", "Employee Code", "employee_name", "Employee", "department", "Department", "branch", "Branch", "status", "Status", "exit_type", "Exit Type", "reason", "Reason", "resignation_date", "Resignation", "notice_start", "Notice Start", "last_working_date", "Last Working", "fnf_status", "F&F", "asset_status", "Assets", "handover_status", "Handover", "access_status", "Access", "interview_status", "Interview", "completed_tasks", "Done", "total_tasks", "Tasks", "blocked_tasks", "Blocked", "readiness", "Readiness"),
		Rows:        rows,
		Summary:     []domain.ReportMetric{{Label: "Exit Requests", Value: strconv.Itoa(len(page.Items))}, {Label: "Active Pipeline", Value: strconv.Itoa(pipeline)}, {Label: "Ready", Value: strconv.Itoa(ready)}, {Label: "Blocked", Value: strconv.Itoa(blocked)}, {Label: "Notice Overdue", Value: strconv.Itoa(overdueNotice)}, {Label: "Reason Mix", Value: stageMixLabel(reasonCounts)}},
	}, nil
}

func (s *TenantService) reportBranding(ctx context.Context, tenantID uuid.UUID) *domain.ReportBranding {
	branding, err := s.branding.GetTenantBranding(ctx, tenantID)
	if err != nil || branding == nil {
		return &domain.ReportBranding{DisplayName: "Setika HRMS", PrimaryColor: "#588368", SecondaryColor: "#2f6f7d"}
	}
	return &domain.ReportBranding{DisplayName: firstNonEmptyString(valueFromPtr(branding.DisplayName), "Setika HRMS"), PrimaryColor: firstNonEmptyString(branding.PrimaryColor, "#588368"), SecondaryColor: firstNonEmptyString(branding.SecondaryColor, "#2f6f7d"), LogoPath: valueFromPtr(branding.LogoPath)}
}

func renderReportPDF(dataset *domain.ReportDataset) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 12)
	pdf.AddPage()
	r, g, b := reportHexColor("#588368")
	brandName := "Setika HRMS"
	if dataset.Branding != nil {
		brandName = dataset.Branding.DisplayName
		r, g, b = reportHexColor(dataset.Branding.PrimaryColor)
	}
	pdf.SetFillColor(r, g, b)
	pdf.Rect(10, 10, 277, 12, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(277, 12, brandName+" - "+dataset.Title, "", 1, "C", false, 0, "")
	pdf.SetTextColor(17, 24, 39)
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "", 9)
	if dataset.PeriodLabel != "" {
		pdf.CellFormat(277, 5, "Period: "+dataset.PeriodLabel, "", 1, "L", false, 0, "")
	}
	for _, metric := range dataset.Summary {
		pdf.CellFormat(55, 6, metric.Label+": "+metric.Value, "1", 0, "L", false, 0, "")
	}
	pdf.Ln(9)
	width := 277.0 / float64(maxInt(1, len(dataset.Columns)))
	pdf.SetFont("Helvetica", "B", 7)
	pdf.SetFillColor(238, 244, 241)
	for _, col := range dataset.Columns {
		pdf.CellFormat(width, 6, truncate(col.Label, 24), "1", 0, "L", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Helvetica", "", 7)
	for _, row := range dataset.Rows {
		for i := range dataset.Columns {
			value := ""
			if i < len(row) {
				value = row[i]
			}
			pdf.CellFormat(width, 5, truncate(value, 32), "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func renderReportXLSX(dataset *domain.ReportDataset) ([]byte, error) {
	var buf bytes.Buffer
	zipper := zip.NewWriter(&buf)
	files := map[string]string{
		"[Content_Types].xml":        xlsxContentTypes,
		"_rels/.rels":                xlsxRels,
		"xl/workbook.xml":            xlsxWorkbook,
		"xl/_rels/workbook.xml.rels": xlsxWorkbookRels,
		"xl/styles.xml":              xlsxStyles,
		"xl/worksheets/sheet1.xml":   xlsxSheet(dataset),
		"docProps/core.xml":          xlsxCore(dataset),
		"docProps/app.xml":           xlsxApp,
	}
	for name, content := range files {
		writer, err := zipper.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			return nil, err
		}
	}
	if err := zipper.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func xlsxSheet(dataset *domain.ReportDataset) string {
	rows := [][]string{{dataset.Title}}
	if dataset.Branding != nil {
		rows = append(rows, []string{"Tenant", dataset.Branding.DisplayName})
	}
	if dataset.PeriodLabel != "" {
		rows = append(rows, []string{"Period", dataset.PeriodLabel})
	}
	for _, metric := range dataset.Summary {
		rows = append(rows, []string{metric.Label, metric.Value})
	}
	header := make([]string, 0, len(dataset.Columns))
	for _, col := range dataset.Columns {
		header = append(header, col.Label)
	}
	rows = append(rows, header)
	rows = append(rows, dataset.Rows...)
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for rIdx, row := range rows {
		b.WriteString(`<row r="` + strconv.Itoa(rIdx+1) + `">`)
		for cIdx, cell := range row {
			ref := columnName(cIdx+1) + strconv.Itoa(rIdx+1)
			b.WriteString(`<c r="` + ref + `" t="inlineStr"><is><t>` + xmlEscape(cell) + `</t></is></c>`)
		}
		b.WriteString(`</row>`)
	}
	b.WriteString(`</sheetData></worksheet>`)
	return b.String()
}

func reportColumns(values ...string) []domain.ReportColumn {
	cols := make([]domain.ReportColumn, 0, len(values)/2)
	for i := 0; i+1 < len(values); i += 2 {
		cols = append(cols, domain.ReportColumn{Key: values[i], Label: values[i+1]})
	}
	return cols
}

func parseReportRange(query ports.ReportDatasetQuery) (string, string, error) {
	start, end, err := parseDateRangeOrToday(query.StartDate, query.EndDate)
	if err != nil {
		return "", "", err
	}
	return start.Format("2006-01-02"), end.Format("2006-01-02"), nil
}

func parseOptionalReportRange(query ports.ReportDatasetQuery) (*time.Time, *time.Time, error) {
	if strings.TrimSpace(query.StartDate) == "" || strings.TrimSpace(query.EndDate) == "" {
		return nil, nil, nil
	}
	start, end, err := parseDateRangeOrToday(query.StartDate, query.EndDate)
	if err != nil {
		return nil, nil, err
	}
	return &start, &end, nil
}

func reportMonthYear(query ports.ReportDatasetQuery) (int32, int32) {
	now := time.Now()
	month, year := query.Month, query.Year
	if month == 0 {
		month = int32(now.Month())
	}
	if year == 0 {
		year = int32(now.Year())
	}
	return month, year
}

func reportPeriodLabel(start *time.Time, end *time.Time) string {
	if start == nil || end == nil {
		return ""
	}
	return start.Format("2006-01-02") + " to " + end.Format("2006-01-02")
}

func employeeName(first string, last *string) string {
	return strings.TrimSpace(first + " " + valueFromPtr(last))
}

func dateString(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02")
}

func minutesLabel(value int32) string {
	if value <= 0 {
		return "0m"
	}
	return fmt.Sprintf("%dh %02dm", value/60, value%60)
}

func ptrInt32(value *int32) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(*value)
}

func ptrFloatString(value *float64) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%.2f", *value)
}

func complianceRisk(employee *domain.EmployeeListItem) string {
	risks := []string{}
	if employee.Email == nil {
		risks = append(risks, "email missing")
	}
	if employee.Mobile == nil {
		risks = append(risks, "mobile missing")
	}
	if employee.DepartmentID == nil {
		risks = append(risks, "department missing")
	}
	if employee.BranchID == nil {
		risks = append(risks, "branch missing")
	}
	if employee.ProbationStatus == domain.EmployeeProbationProbation {
		risks = append(risks, "probation pending")
	}
	return strings.Join(risks, ", ")
}

func candidateApplicationName(app *domain.CandidateApplication) string {
	return strings.TrimSpace(valueFromPtr(app.CandidateFirstname) + " " + valueFromPtr(app.CandidateLastname))
}

func candidateOnboardingName(item *domain.CandidateOnboarding) string {
	return strings.TrimSpace(valueFromPtr(item.CandidateFirstname) + " " + valueFromPtr(item.CandidateLastname))
}

func percentString(numerator int, denominator int) string {
	if denominator <= 0 {
		return "0.0%"
	}
	return fmt.Sprintf("%.1f%%", float64(numerator)*100/float64(denominator))
}

func averageDaysLabel(totalDays int, count int) string {
	if count <= 0 {
		return "0d"
	}
	return fmt.Sprintf("%.1fd", float64(totalDays)/float64(count))
}

func daysBetween(start time.Time, end time.Time) int {
	if end.Before(start) {
		return 0
	}
	return int(end.Sub(start).Hours() / 24)
}

func stageMixLabel(values map[string]int) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(values))
	for key, value := range values {
		parts = append(parts, fmt.Sprintf("%s:%d", key, value))
	}
	return strings.Join(parts, ", ")
}

func exitReady(item *domain.EmployeeExitRequest) bool {
	return item.FinalSettlementStatus == domain.EmployeeExitTaskCompleted &&
		item.AssetClearanceStatus == domain.EmployeeExitTaskCompleted &&
		item.HandoverStatus == domain.EmployeeExitTaskCompleted &&
		item.AccessRevocationStatus == domain.EmployeeExitTaskCompleted &&
		(item.TotalTasks == 0 || item.CompletedTasks >= item.TotalTasks) &&
		item.BlockedTasks == 0
}

func exitReadinessLabel(item *domain.EmployeeExitRequest) string {
	if exitReady(item) {
		return "ready"
	}
	if item.BlockedTasks > 0 {
		return "blocked"
	}
	return "pending"
}

func reportFileName(dataset *domain.ReportDataset, ext string) string {
	name := strings.ToLower(strings.TrimSpace(dataset.ReportCode))
	if name == "" {
		name = dataset.Title
	}
	name = strings.NewReplacer(" ", "-", "/", "-", "\\", "-", "_", "-").Replace(name)
	return path.Base(path.Clean(name)) + "." + ext
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func reportHexColor(value string) (int, int, int) {
	value = strings.TrimPrefix(strings.TrimSpace(value), "#")
	if len(value) != 6 {
		return 88, 131, 104
	}
	var r, g, b int
	if _, err := fmt.Sscanf(value, "%02x%02x%02x", &r, &g, &b); err != nil {
		return 88, 131, 104
	}
	return r, g, b
}

func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max-1] + "."
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func columnName(n int) string {
	name := ""
	for n > 0 {
		n--
		name = string(rune('A'+n%26)) + name
		n /= 26
	}
	return name
}

func xmlEscape(value string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(value))
	return buf.String()
}

func xlsxCore(dataset *domain.ReportDataset) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:title>` + xmlEscape(dataset.Title) + `</dc:title><dc:creator>Setika HRMS</dc:creator></cp:coreProperties>`
}

const xlsxContentTypes = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/><Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/><Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/><Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/></Types>`
const xlsxRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/><Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/></Relationships>`
const xlsxWorkbook = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="Report" sheetId="1" r:id="rId1"/></sheets></workbook>`
const xlsxWorkbookRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/></Relationships>`
const xlsxStyles = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><fonts count="1"><font><sz val="11"/><name val="Calibri"/></font></fonts><fills count="1"><fill><patternFill patternType="none"/></fill></fills><borders count="1"><border/></borders><cellStyleXfs count="1"><xf/></cellStyleXfs><cellXfs count="1"><xf/></cellXfs></styleSheet>`
const xlsxApp = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties"><Application>Setika HRMS</Application></Properties>`
