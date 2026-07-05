package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) RefreshInsights(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (*domain.InsightWorkspace, error) {
	candidates, err := s.buildDeterministicInsights(ctx, tenantID)
	if err != nil {
		s.logError("build deterministic insights", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if s.insightScorer != nil {
		scored, err := s.insightScorer.ScoreInsights(ctx, tenantID, candidates)
		if err != nil {
			s.logError("score insights", err, serviceTenantIDField(tenantID))
			return nil, err
		}
		candidates = scored
	}
	for _, candidate := range candidates {
		item, err := s.insights.UpsertInsight(ctx, candidate, actorID)
		if err != nil {
			s.logError("upsert insight", err, serviceTenantIDField(tenantID), serviceStringField("insight_key", candidate.InsightKey))
			return nil, err
		}
		toStatus := item.Status
		if _, err := s.insights.CreateInsightEvent(ctx, &domain.InsightEvent{TenantID: tenantID, InsightID: item.ID, Action: "refreshed", ToStatus: &toStatus, Metadata: rawJSON(map[string]any{"source": item.Source, "score": item.Score})}, actorID); err != nil {
			s.logError("record insight refresh event", err, serviceTenantIDField(tenantID), serviceStringField("insight_id", item.ID.String()))
			return nil, err
		}
		s.recordInsightAIAction(ctx, item, actorID)
	}
	return s.ListInsightWorkspace(ctx, domain.InsightFilter{TenantID: tenantID, Limit: 200})
}

func (s *TenantService) ListInsightWorkspace(ctx context.Context, filter domain.InsightFilter) (*domain.InsightWorkspace, error) {
	items, err := s.insights.ListInsights(ctx, filter)
	if err != nil {
		s.logError("list insights", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return &domain.InsightWorkspace{Items: items, Summary: insightSummary(items)}, nil
}

func (s *TenantService) UpdateInsightStatus(ctx context.Context, cmd ports.InsightStatusCommand) (*domain.Insight, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		err := domain.ErrInvalidInsightID
		s.logError("validate insight status update", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status, err := domain.ValidateInsightStatus(cmd.Status)
	if err != nil {
		s.logError("validate insight status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("status", cmd.Status))
		return nil, err
	}
	before, err := s.insights.GetInsight(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		s.logError("get insight before status update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("insight_id", cmd.ID.String()))
		return nil, err
	}
	cmd.Status = status
	if cmd.ResolutionNote == nil {
		cmd.ResolutionNote = cmd.Remarks
	}
	updated, err := s.insights.UpdateInsightStatus(ctx, cmd)
	if err != nil {
		s.logError("update insight status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("insight_id", cmd.ID.String()))
		return nil, err
	}
	fromStatus, toStatus := before.Status, updated.Status
	if _, err := s.insights.CreateInsightEvent(ctx, &domain.InsightEvent{TenantID: cmd.TenantID, InsightID: cmd.ID, Action: "status_changed", FromStatus: &fromStatus, ToStatus: &toStatus, Remarks: cmd.Remarks, Metadata: rawJSON(map[string]any{"assigned_to": cmd.AssignedTo})}, cmd.ActorID); err != nil {
		s.logError("record insight status event", err, serviceTenantIDField(cmd.TenantID), serviceStringField("insight_id", cmd.ID.String()))
		return nil, err
	}
	if status == domain.InsightStatusOverridden || status == domain.InsightStatusDismissed {
		reason := "Human review changed insight outcome."
		if cmd.Remarks != nil && *cmd.Remarks != "" {
			reason = *cmd.Remarks
		}
		if _, err := s.OverrideAIAction(ctx, ports.AIOverrideCommand{TenantID: cmd.TenantID, InsightID: &cmd.ID, OverrideType: "insight_status", OriginalStatus: &fromStatus, OverrideStatus: status, Reason: reason, Decision: domain.AIOverrideDecisionManualAction, Metadata: rawJSON(map[string]any{"source": "insight_status_update"}), ActorID: cmd.ActorID}); err != nil {
			s.log.Warn().Err(err).Str("tenant_id", cmd.TenantID.String()).Str("insight_id", cmd.ID.String()).Msg("hrms: insight status updated but ai override tracking failed")
		}
	}
	return updated, nil
}

func (s *TenantService) ListInsightEvents(ctx context.Context, tenantID uuid.UUID, insightID uuid.UUID) ([]*domain.InsightEvent, error) {
	if tenantID == uuid.Nil || insightID == uuid.Nil {
		return nil, domain.ErrInvalidInsightID
	}
	events, err := s.insights.ListInsightEvents(ctx, tenantID, insightID)
	if err != nil {
		s.logError("list insight events", err, serviceTenantIDField(tenantID), serviceStringField("insight_id", insightID.String()))
		return nil, err
	}
	return events, nil
}

func (s *TenantService) buildDeterministicInsights(ctx context.Context, tenantID uuid.UUID) ([]*domain.Insight, error) {
	now := time.Now().UTC()
	start := now.AddDate(0, 0, -30).Format("2006-01-02")
	end := now.Format("2006-01-02")
	items := []*domain.Insight{}

	employees, err := s.employees.ListEmployees(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	employeeByUser := map[uuid.UUID]*domain.EmployeeListItem{}
	for _, employee := range employees {
		if employee.UserID != uuid.Nil {
			employeeByUser[employee.UserID] = employee
		}
		if complianceRisk(employee) != "" && !employee.Inactive {
			items = append(items, newInsight(tenantID, "compliance", "compliance_risk", "compliance-master-"+employee.UserID.String(), domain.InsightSeverityMedium, "Compliance data gap", fmt.Sprintf("%s has missing master data: %s.", employeeName(employee.Firstname, employee.Lastname), complianceRisk(employee)), 72, 68, "employee", &employee.UserID, []string{complianceRisk(employee)}, []string{"Complete employee master data before payroll and compliance reporting."}, map[string]any{"employee_code": valueFromPtr(employee.EmployeeCode), "department": valueFromPtr(employee.DepartmentName)}))
		}
		if employee.ProbationStatus == domain.EmployeeProbationProbation && employee.ProbationEndDate != nil && employee.ProbationEndDate.Before(now.AddDate(0, 0, 15)) && !employee.Inactive {
			severity := domain.InsightSeverityMedium
			if employee.ProbationEndDate.Before(now) {
				severity = domain.InsightSeverityHigh
			}
			items = append(items, newInsight(tenantID, "attrition", "probation_action_due", "probation-"+employee.UserID.String(), severity, "Probation action due", fmt.Sprintf("%s needs probation confirmation, extension, or manager review.", employeeName(employee.Firstname, employee.Lastname)), 74, 70, "employee", &employee.UserID, []string{"Probation end date is " + employee.ProbationEndDate.Format("2006-01-02")}, []string{"Ask manager for confirmation recommendation and update employee probation status."}, map[string]any{"probation_end": employee.ProbationEndDate.Format("2006-01-02")}))
		}
	}

	attendanceReport, err := s.GetAttendanceReport(ctx, ports.AttendanceReportQuery{TenantID: tenantID, StartDate: start, EndDate: end})
	if err == nil && attendanceReport != nil {
		items = append(items, attendanceInsights(tenantID, attendanceReport, employeeByUser)...)
	} else if err != nil {
		s.log.Warn().Err(err).Str("tenant_id", tenantID.String()).Msg("hrms: skipped attendance insights")
	}

	leaveStart, leaveEnd := now.AddDate(0, 0, -90), now
	leaveRows, err := s.ListLeaveReportRows(ctx, domain.LeaveReportFilter{TenantID: tenantID, StartDate: &leaveStart, EndDate: &leaveEnd})
	if err == nil {
		items = append(items, leaveInsights(tenantID, leaveRows)...)
	} else {
		s.log.Warn().Err(err).Str("tenant_id", tenantID.String()).Msg("hrms: skipped leave insights")
	}

	month, year := int32(now.Month()), int32(now.Year())
	reconcile, err := s.ListPayrollReconciliationRows(ctx, tenantID, month, year)
	if err == nil {
		items = append(items, payrollInsights(tenantID, month, year, reconcile)...)
	} else {
		s.log.Warn().Err(err).Str("tenant_id", tenantID.String()).Msg("hrms: skipped payroll insights")
	}

	onboardings, err := s.candidateOnboardings.ListCandidateOnboardings(ctx, domain.CandidateOnboardingFilter{TenantID: tenantID, Limit: 10000})
	if err == nil {
		items = append(items, onboardingInsights(tenantID, now, onboardings.Items)...)
	} else {
		s.log.Warn().Err(err).Str("tenant_id", tenantID.String()).Msg("hrms: skipped onboarding insights")
	}

	exits, err := s.employeeExits.ListEmployeeExitRequests(ctx, domain.EmployeeExitFilter{TenantID: tenantID, Limit: 10000})
	if err == nil {
		items = append(items, exitInsights(tenantID, now, exits.Items)...)
	} else {
		s.log.Warn().Err(err).Str("tenant_id", tenantID.String()).Msg("hrms: skipped exit insights")
	}

	return items, nil
}

func attendanceInsights(tenantID uuid.UUID, report *domain.AttendanceReport, employees map[uuid.UUID]*domain.EmployeeListItem) []*domain.Insight {
	type stats struct{ late, absent, incomplete int }
	byUser := map[uuid.UUID]stats{}
	for _, row := range report.Rows {
		current := byUser[row.UserID]
		if row.LateMinutes > 0 {
			current.late++
		}
		if row.Status == domain.AttendanceStatusAbsent {
			current.absent++
		}
		if row.Status == domain.AttendanceStatusIncomplete {
			current.incomplete++
		}
		byUser[row.UserID] = current
	}
	items := []*domain.Insight{}
	for userID, stat := range byUser {
		if stat.late < 4 && stat.absent < 3 && stat.incomplete < 2 {
			continue
		}
		name := userID.String()
		if employee := employees[userID]; employee != nil {
			name = employeeName(employee.Firstname, employee.Lastname)
		}
		severity := domain.InsightSeverityMedium
		if stat.absent >= 5 || stat.incomplete >= 4 {
			severity = domain.InsightSeverityHigh
		}
		items = append(items, newInsight(tenantID, "attendance", "attendance_anomaly", "attendance-"+userID.String(), severity, "Attendance anomaly", fmt.Sprintf("%s has repeated attendance exceptions in the last 30 days.", name), 82, float64(stat.late*8+stat.absent*15+stat.incomplete*12), "employee", &userID, []string{fmt.Sprintf("%d late, %d absent, %d incomplete days", stat.late, stat.absent, stat.incomplete)}, []string{"Review shift fit, manager context, biometric/device issues, and regularisation queue."}, map[string]any{"late_days": stat.late, "absent_days": stat.absent, "incomplete_days": stat.incomplete}))
	}
	if report.Summary.EmployeeDays > 0 && report.Summary.AttendanceRate < 85 {
		items = append(items, newInsight(tenantID, "engagement", "engagement_health", "tenant-attendance-rate", domain.InsightSeverityHigh, "Low attendance health", fmt.Sprintf("Tenant attendance rate is %.1f%% for the current 30-day window.", report.Summary.AttendanceRate), 78, 85-report.Summary.AttendanceRate, "tenant", nil, []string{"Attendance rate below 85% threshold"}, []string{"Review department-level attendance, manager staffing notes, location/device issues, and leave overlap."}, map[string]any{"attendance_rate": report.Summary.AttendanceRate}))
	}
	return items
}

func leaveInsights(tenantID uuid.UUID, rows []*domain.LeaveReportRow) []*domain.Insight {
	type stats struct {
		requests int
		days     float64
	}
	byUser := map[uuid.UUID]stats{}
	for _, row := range rows {
		if row.Status != "approved" && row.Status != "Approved" && row.Status != "pending" && row.Status != "Pending" {
			continue
		}
		current := byUser[row.UserID]
		current.requests++
		current.days += row.Days
		byUser[row.UserID] = current
	}
	items := []*domain.Insight{}
	for userID, stat := range byUser {
		if stat.requests < 4 && stat.days < 10 {
			continue
		}
		items = append(items, newInsight(tenantID, "leave", "leave_abuse_signal", "leave-"+userID.String(), domain.InsightSeverityMedium, "Leave usage signal", "Employee leave usage is above the review threshold for the last 90 days.", 70, stat.days, "employee", &userID, []string{fmt.Sprintf("%d requests covering %.1f days", stat.requests, stat.days)}, []string{"Review leave type mix, medical/document requirements, manager notes, and policy fit before action."}, map[string]any{"requests": stat.requests, "days": stat.days}))
	}
	return items
}

func payrollInsights(tenantID uuid.UUID, month int32, year int32, rows []*domain.PayrollReconciliationRow) []*domain.Insight {
	items := []*domain.Insight{}
	for _, row := range rows {
		if row.ReconciliationStatus == "ok" {
			continue
		}
		items = append(items, newInsight(tenantID, "payroll", "payroll_anomaly", fmt.Sprintf("payroll-%d-%d-%s", year, month, row.UserID), domain.InsightSeverityHigh, "Payroll reconciliation anomaly", fmt.Sprintf("%s has payroll reconciliation status %s.", employeeName(row.Firstname, row.Lastname), row.ReconciliationStatus), 86, 90, "employee", &row.UserID, []string{"Payroll reconciliation status: " + row.ReconciliationStatus}, []string{"Resolve attendance/LOP/salary-slip blockers before payroll lock."}, map[string]any{"month": month, "year": year, "net_salary": row.NetSalary}))
	}
	return items
}

func onboardingInsights(tenantID uuid.UUID, now time.Time, rows []*domain.CandidateOnboarding) []*domain.Insight {
	items := []*domain.Insight{}
	for _, row := range rows {
		if row.OnboardingStatus == domain.OnboardStatusCompleted || (row.OverdueTasks == 0 && row.ProgressPercentage >= 70) {
			continue
		}
		key := "onboarding-" + row.ID.String()
		items = append(items, newInsight(tenantID, "onboarding", "onboarding_delay", key, domain.InsightSeverityMedium, "Onboarding delay risk", fmt.Sprintf("%s onboarding is delayed or missing required tasks.", candidateOnboardingName(row)), 76, float64(100-row.ProgressPercentage), "candidate_onboarding", &row.ID, []string{fmt.Sprintf("%d overdue tasks, %d%% progress", row.OverdueTasks, row.ProgressPercentage)}, []string{"Escalate overdue tasks and required document collection before joining day."}, map[string]any{"candidate_id": row.CandidateID, "overdue_tasks": row.OverdueTasks, "progress": row.ProgressPercentage}))
	}
	return items
}

func exitInsights(tenantID uuid.UUID, now time.Time, rows []*domain.EmployeeExitRequest) []*domain.Insight {
	items := []*domain.Insight{}
	for _, row := range rows {
		if row.Status == domain.EmployeeExitStatusCompleted || row.Status == domain.EmployeeExitStatusCanceled || row.Status == domain.EmployeeExitStatusRejected {
			continue
		}
		if row.BlockedTasks == 0 && row.CompletedTasks >= row.TotalTasks && !row.LastWorkingDate.Before(now) {
			continue
		}
		severity := domain.InsightSeverityMedium
		if row.BlockedTasks > 0 || row.LastWorkingDate.Before(now) {
			severity = domain.InsightSeverityHigh
		}
		items = append(items, newInsight(tenantID, "exit", "exit_readiness_risk", "exit-"+row.ID.String(), severity, "Exit readiness risk", fmt.Sprintf("%s exit workflow needs clearance/F&F attention.", employeeName(valueFromPtr(row.EmployeeFirstname), row.EmployeeLastname)), 80, float64(row.BlockedTasks*25+row.TotalTasks-row.CompletedTasks), "employee_exit", &row.ID, []string{fmt.Sprintf("%d/%d tasks complete, %d blocked", row.CompletedTasks, row.TotalTasks, row.BlockedTasks)}, []string{"Review asset, handover, access revocation, exit interview, and final settlement readiness."}, map[string]any{"last_working_date": row.LastWorkingDate.Format("2006-01-02"), "fnf_status": row.FinalSettlementStatus}))
	}
	return items
}

func newInsight(tenantID uuid.UUID, category string, insightType string, key string, severity string, title string, summary string, confidence float64, score float64, entityType string, entityID *uuid.UUID, reasons []string, recommendations []string, context map[string]any) *domain.Insight {
	now := time.Now().UTC()
	explainability := map[string]any{"engine": "deterministic_rules_v1", "human_review_required": true, "ai_assisted": false}
	return &domain.Insight{TenantID: tenantID, InsightKey: key, InsightType: insightType, Category: category, Severity: severity, Status: domain.InsightStatusOpen, Title: title, Summary: summary, ConfidenceScore: confidence, Score: score, Source: domain.InsightSourceDeterministic, EntityType: &entityType, EntityID: entityID, EmployeeUserID: entityEmployeeID(entityType, entityID), Reasons: rawJSON(reasons), Recommendations: rawJSON(recommendations), Context: rawJSON(context), Explainability: rawJSON(explainability), DetectedAt: now}
}

func entityEmployeeID(entityType string, entityID *uuid.UUID) *uuid.UUID {
	if entityType == "employee" {
		return entityID
	}
	return nil
}

func insightSummary(items []*domain.Insight) domain.InsightSummary {
	summary := domain.InsightSummary{ByCategory: map[string]int32{}}
	for _, item := range items {
		summary.Total++
		summary.ByCategory[item.Category]++
		switch item.Status {
		case domain.InsightStatusOpen:
			summary.Open++
		case domain.InsightStatusReviewing:
			summary.Reviewing++
		case domain.InsightStatusResolved:
			summary.Resolved++
		case domain.InsightStatusDismissed:
			summary.Dismissed++
		case domain.InsightStatusOverridden:
			summary.Overridden++
		}
		switch item.Severity {
		case domain.InsightSeverityCritical:
			summary.Critical++
		case domain.InsightSeverityHigh:
			summary.High++
		case domain.InsightSeverityMedium:
			summary.Medium++
		case domain.InsightSeverityLow:
			summary.Low++
		}
	}
	return summary
}

func intLabel(value int32) string {
	return strconv.Itoa(int(value))
}
