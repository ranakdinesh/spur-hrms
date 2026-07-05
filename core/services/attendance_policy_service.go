package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateAttendancePolicy(ctx context.Context, cmd ports.AttendancePolicyCommand) (*domain.AttendancePolicy, error) {
	item, err := s.buildAttendancePolicy(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.attendancePolicies.CreateAttendancePolicy(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create attendance policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAttendancePolicies(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendancePolicy, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate attendance policy list tenant", err)
		return nil, err
	}
	items, err := s.attendancePolicies.ListAttendancePolicies(ctx, tenantID)
	if err != nil {
		s.logError("list attendance policies", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateAttendancePolicy(ctx context.Context, cmd ports.AttendancePolicyCommand) (*domain.AttendancePolicy, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidAttendancePolicyID
		s.logError("validate attendance policy update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := s.buildAttendancePolicy(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.attendancePolicies.UpdateAttendancePolicy(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update attendance policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("policy_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteAttendancePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if id == uuid.Nil {
		return domain.ErrInvalidAttendancePolicyID
	}
	if err := s.attendancePolicies.DeleteAttendancePolicy(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete attendance policy", err, serviceTenantIDField(tenantID), serviceStringField("policy_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateAttendanceRoster(ctx context.Context, cmd ports.AttendanceRosterCommand) (*domain.AttendanceRoster, error) {
	item, err := s.buildAttendanceRoster(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.attendanceRosters.CreateAttendanceRoster(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create attendance roster", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAttendanceRostersByDateRange(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceRoster, error) {
	start, end, err := parseDateRangeOrToday(startDate, endDate)
	if err != nil {
		return nil, err
	}
	items, err := s.attendanceRosters.ListAttendanceRostersByDateRange(ctx, tenantID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		s.logError("list attendance rosters", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListAttendanceRostersByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceRoster, error) {
	start, end, err := parseDateRangeOrToday(startDate, endDate)
	if err != nil {
		return nil, err
	}
	items, err := s.attendanceRosters.ListAttendanceRostersByUser(ctx, tenantID, userID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		s.logError("list attendance rosters by user", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateAttendanceRoster(ctx context.Context, cmd ports.AttendanceRosterCommand) (*domain.AttendanceRoster, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAttendanceRosterID
	}
	item, err := s.buildAttendanceRoster(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.attendanceRosters.UpdateAttendanceRoster(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update attendance roster", err, serviceTenantIDField(cmd.TenantID), serviceStringField("roster_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteAttendanceRoster(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if id == uuid.Nil {
		return domain.ErrInvalidAttendanceRosterID
	}
	if err := s.attendanceRosters.DeleteAttendanceRoster(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete attendance roster", err, serviceTenantIDField(tenantID), serviceStringField("roster_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateAttendanceRequest(ctx context.Context, cmd ports.AttendanceRequestCommand) (*domain.AttendanceRequest, error) {
	date, err := parseAttendanceDate(cmd.Date, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	checkIn, err := parseOptionalRFC3339(cmd.RequestedCheckInAt)
	if err != nil {
		return nil, err
	}
	checkOut, err := parseOptionalRFC3339(cmd.RequestedCheckOutAt)
	if err != nil {
		return nil, err
	}
	employee, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID)
	if err != nil {
		s.logError("validate attendance request employee", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	item, err := domain.NewAttendanceRequest(cmd.TenantID, cmd.UserID, date, cmd.RequestType, cmd.RequestedType, checkIn, checkOut, cmd.RequestedWorkMode, cmd.PolicyID, cmd.RosterID, cmd.Reason)
	if err != nil {
		s.logError("validate attendance request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	result, err := s.attendanceRequests.CreateAttendanceRequest(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create attendance request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	workflow, err := s.resolveExceptionWorkflow(ctx, cmd.TenantID, item.RequestType, employee.DepartmentID, employee.BranchID)
	if err != nil {
		return nil, err
	}
	if workflow != nil {
		result, err = s.applyWorkflowToAttendanceRequest(ctx, result, workflow, cmd.ActorID)
		if err != nil {
			return nil, err
		}
		if workflow.AutoApprove || workflow.RouteMode == domain.AttendancePolicyApprovalAuto {
			reviewerID := cmd.UserID
			if cmd.ActorID != nil && *cmd.ActorID != uuid.Nil {
				reviewerID = *cmd.ActorID
			}
			return s.ReviewAttendanceRequest(ctx, ports.AttendanceReviewCommand{TenantID: cmd.TenantID, RequestID: result.ID, Status: domain.LeaveStatusApproved, Remarks: stringPtr("Auto-approved by attendance exception workflow"), ReviewerID: reviewerID})
		}
	}
	return result, nil
}

func (s *TenantService) ListAttendanceRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.AttendanceRequest, error) {
	items, err := s.attendanceRequests.ListAttendanceRequestsByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("list attendance requests by user", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListAttendanceRequestsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.AttendanceRequest, error) {
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "" {
		status = domain.LeaveStatusPending
	}
	items, err := s.attendanceRequests.ListAttendanceRequestsByStatus(ctx, tenantID, status)
	if err != nil {
		s.logError("list attendance requests by status", err, serviceTenantIDField(tenantID), serviceStringField("status", status))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ReviewAttendanceRequest(ctx context.Context, cmd ports.AttendanceReviewCommand) (*domain.AttendanceRequest, error) {
	if cmd.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	if cmd.RequestID == uuid.Nil {
		return nil, domain.ErrInvalidAttendanceRequestID
	}
	status, err := domain.ValidateLeaveStatus(cmd.Status)
	if err != nil || status == domain.LeaveStatusPending {
		return nil, domain.ErrInvalidAttendanceReview
	}
	if cmd.ReviewerID == uuid.Nil {
		return nil, domain.ErrInvalidEmployeeUserID
	}
	var reviewed *domain.AttendanceRequest
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		before, err := s.attendanceRequests.GetAttendanceRequest(txCtx, cmd.TenantID, cmd.RequestID)
		if err != nil {
			return err
		}
		item, err := s.attendanceRequests.UpdateAttendanceRequestReview(txCtx, cmd.TenantID, cmd.RequestID, status, cmd.ReviewerID, cmd.Remarks)
		if err != nil {
			return err
		}
		reviewed = item
		if err := s.createAttendanceExceptionEvent(txCtx, item, "reviewed", &before.Status, &item.Status, item.RouteMode, cmd.Remarks, &cmd.ReviewerID); err != nil {
			return err
		}
		if status == domain.LeaveStatusApproved {
			return s.applyApprovedAttendanceRequest(txCtx, item, &cmd.ReviewerID)
		}
		return nil
	})
	if err != nil {
		s.logError("review attendance request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("request_id", cmd.RequestID.String()))
		return nil, err
	}
	return reviewed, nil
}

func (s *TenantService) CreateAttendanceExceptionWorkflow(ctx context.Context, cmd ports.AttendanceExceptionWorkflowCommand) (*domain.AttendanceExceptionWorkflow, error) {
	item, err := s.buildAttendanceExceptionWorkflow(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.attendanceExceptionWorkflows.CreateAttendanceExceptionWorkflow(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create attendance exception workflow", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListAttendanceExceptionWorkflows(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendanceExceptionWorkflow, error) {
	items, err := s.attendanceExceptionWorkflows.ListAttendanceExceptionWorkflows(ctx, tenantID)
	if err != nil {
		s.logError("list attendance exception workflows", err, serviceTenantIDField(tenantID))
	}
	return items, err
}

func (s *TenantService) UpdateAttendanceExceptionWorkflow(ctx context.Context, cmd ports.AttendanceExceptionWorkflowCommand) (*domain.AttendanceExceptionWorkflow, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAttendanceExceptionWorkflow
	}
	item, err := s.buildAttendanceExceptionWorkflow(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.attendanceExceptionWorkflows.UpdateAttendanceExceptionWorkflow(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update attendance exception workflow", err, serviceTenantIDField(cmd.TenantID), serviceStringField("workflow_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteAttendanceExceptionWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if id == uuid.Nil {
		return domain.ErrInvalidAttendanceExceptionWorkflow
	}
	if err := s.attendanceExceptionWorkflows.DeleteAttendanceExceptionWorkflow(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete attendance exception workflow", err, serviceTenantIDField(tenantID), serviceStringField("workflow_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListAttendanceExceptionEvents(ctx context.Context, tenantID uuid.UUID, attendanceRequestID uuid.UUID) ([]*domain.AttendanceExceptionEvent, error) {
	if attendanceRequestID == uuid.Nil {
		return nil, domain.ErrInvalidAttendanceRequestID
	}
	items, err := s.attendanceExceptionWorkflows.ListAttendanceExceptionEvents(ctx, tenantID, attendanceRequestID)
	if err != nil {
		s.logError("list attendance exception events", err, serviceTenantIDField(tenantID), serviceStringField("request_id", attendanceRequestID.String()))
	}
	return items, err
}

func (s *TenantService) ListPayrollBlockingAttendanceRequests(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceRequest, error) {
	start, end, err := parseDateRangeOrToday(startDate, endDate)
	if err != nil {
		return nil, err
	}
	items, err := s.attendanceExceptionWorkflows.ListPayrollBlockingAttendanceRequests(ctx, tenantID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		s.logError("list payroll blocking attendance requests", err, serviceTenantIDField(tenantID))
	}
	return items, err
}

func (s *TenantService) resolveExceptionWorkflow(ctx context.Context, tenantID uuid.UUID, requestType string, departmentID *uuid.UUID, branchID *uuid.UUID) (*domain.AttendanceExceptionWorkflow, error) {
	workflow, err := s.attendanceExceptionWorkflows.ResolveAttendanceExceptionWorkflow(ctx, tenantID, requestType, departmentID, branchID)
	if err != nil {
		if errors.Is(err, domain.ErrAttendanceExceptionWorkflowNotFound) {
			return nil, nil
		}
		s.logError("resolve attendance exception workflow", err, serviceTenantIDField(tenantID), serviceStringField("request_type", requestType))
		return nil, err
	}
	return workflow, nil
}

func (s *TenantService) applyWorkflowToAttendanceRequest(ctx context.Context, request *domain.AttendanceRequest, workflow *domain.AttendanceExceptionWorkflow, actorID *uuid.UUID) (*domain.AttendanceRequest, error) {
	if request == nil || workflow == nil {
		return request, nil
	}
	var dueAt *string
	if workflow.EscalationHours > 0 {
		value := time.Now().UTC().Add(time.Duration(workflow.EscalationHours) * time.Hour).Format(time.RFC3339)
		dueAt = &value
	}
	workflowID := workflow.ID
	routeMode := workflow.RouteMode
	updated, err := s.attendanceRequests.SetAttendanceRequestWorkflow(ctx, request.TenantID, request.ID, &workflowID, &routeMode, dueAt, workflow.BlockPayrollWhenPending, actorID)
	if err != nil {
		s.logError("set attendance request workflow", err, serviceTenantIDField(request.TenantID), serviceStringField("request_id", request.ID.String()))
		return nil, err
	}
	if err := s.createAttendanceExceptionEvent(ctx, updated, "submitted", nil, &updated.Status, &routeMode, updated.Reason, actorID); err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *TenantService) createAttendanceExceptionEvent(ctx context.Context, request *domain.AttendanceRequest, action string, fromStatus *string, toStatus *string, routedTo *string, remarks *string, actorID *uuid.UUID) error {
	if request == nil {
		return nil
	}
	event := &domain.AttendanceExceptionEvent{TenantID: request.TenantID, AttendanceRequestID: request.ID, WorkflowID: request.WorkflowID, Action: action, FromStatus: fromStatus, ToStatus: toStatus, RoutedTo: routedTo, Remarks: remarks}
	_, err := s.attendanceExceptionWorkflows.CreateAttendanceExceptionEvent(ctx, event, actorID)
	if err != nil {
		s.logError("create attendance exception event", err, serviceTenantIDField(request.TenantID), serviceStringField("request_id", request.ID.String()))
	}
	return err
}

func (s *TenantService) buildAttendanceExceptionWorkflow(ctx context.Context, cmd ports.AttendanceExceptionWorkflowCommand) (*domain.AttendanceExceptionWorkflow, error) {
	if cmd.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.branches.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			return nil, err
		}
	}
	if cmd.DepartmentID != nil && *cmd.DepartmentID != uuid.Nil {
		if _, err := s.departments.GetDepartment(ctx, cmd.TenantID, *cmd.DepartmentID); err != nil {
			return nil, err
		}
	}
	if cmd.ID == uuid.Nil && !cmd.IsActive {
		cmd.IsActive = true
	}
	return domain.NewAttendanceExceptionWorkflow(domain.AttendanceExceptionWorkflowInput{TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, Description: cmd.Description, BranchID: cmd.BranchID, DepartmentID: cmd.DepartmentID, RequestType: cmd.RequestType, RouteMode: cmd.RouteMode, MaxRequestsPerMonth: cmd.MaxRequestsPerMonth, EscalationHours: cmd.EscalationHours, EscalationRouteMode: cmd.EscalationRouteMode, BlockPayrollWhenPending: cmd.BlockPayrollWhenPending, AutoApprove: cmd.AutoApprove, IsActive: cmd.IsActive})
}

func (s *TenantService) buildAttendancePolicy(ctx context.Context, cmd ports.AttendancePolicyCommand) (*domain.AttendancePolicy, error) {
	effectiveFrom, err := parseAttendancePolicyOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	effectiveTo, err := parseAttendancePolicyOptionalDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.branches.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			return nil, err
		}
	}
	if cmd.DepartmentID != nil && *cmd.DepartmentID != uuid.Nil {
		if _, err := s.departments.GetDepartment(ctx, cmd.TenantID, *cmd.DepartmentID); err != nil {
			return nil, err
		}
	}
	if cmd.UserID != nil && *cmd.UserID != uuid.Nil {
		if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, *cmd.UserID); err != nil {
			return nil, err
		}
	}
	return domain.NewAttendancePolicy(domain.AttendancePolicyInput{TenantID: cmd.TenantID, Name: cmd.Name, Code: cmd.Code, Description: cmd.Description, BranchID: cmd.BranchID, DepartmentID: cmd.DepartmentID, UserID: cmd.UserID, ScheduleType: cmd.ScheduleType, IsDefault: cmd.IsDefault, StandardWorkMinutes: cmd.StandardWorkMinutes, MinHalfDayMinutes: cmd.MinHalfDayMinutes, MinFullDayMinutes: cmd.MinFullDayMinutes, GraceLateMinutes: cmd.GraceLateMinutes, GraceEarlyMinutes: cmd.GraceEarlyMinutes, HalfDayLateAfterMinutes: cmd.HalfDayLateAfterMinutes, AbsentLateAfterMinutes: cmd.AbsentLateAfterMinutes, HalfDayEarlyBeforeMinutes: cmd.HalfDayEarlyBeforeMinutes, AbsentEarlyBeforeMinutes: cmd.AbsentEarlyBeforeMinutes, AllowFlexiHours: cmd.AllowFlexiHours, CoreStartTime: cmd.CoreStartTime, CoreEndTime: cmd.CoreEndTime, AllowWFH: cmd.AllowWFH, WFHDaysPerWeek: cmd.WFHDaysPerWeek, AllowPermanentRemote: cmd.AllowPermanentRemote, RequireGeo: cmd.RequireGeo, RequireDevice: cmd.RequireDevice, RegularizationWindowDays: cmd.RegularizationWindowDays, MaxRegularizationsPerMonth: cmd.MaxRegularizationsPerMonth, ApprovalMode: cmd.ApprovalMode, EffectiveFrom: effectiveFrom, EffectiveTo: effectiveTo})
}

func (s *TenantService) buildAttendanceRoster(ctx context.Context, cmd ports.AttendanceRosterCommand) (*domain.AttendanceRoster, error) {
	date, err := parseAttendanceDate(cmd.Date, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID); err != nil {
		return nil, err
	}
	if cmd.PolicyID != nil && *cmd.PolicyID != uuid.Nil {
		if _, err := s.attendancePolicies.GetAttendancePolicy(ctx, cmd.TenantID, *cmd.PolicyID); err != nil {
			return nil, err
		}
	}
	return domain.NewAttendanceRoster(cmd.TenantID, cmd.UserID, cmd.PolicyID, date, cmd.StartTime, cmd.EndTime, cmd.BreakMinutes, cmd.WorkMode, cmd.LocationType, cmd.Remarks)
}

func (s *TenantService) applyApprovedAttendanceRequest(ctx context.Context, item *domain.AttendanceRequest, actorID *uuid.UUID) error {
	if item == nil {
		return nil
	}
	status := domain.AttendanceStatusPresent
	if item.RequestType == domain.AttendanceRequestAbsent {
		status = domain.AttendanceStatusAbsent
	}
	workMode := item.RequestedWorkMode
	if workMode == nil && (item.RequestType == domain.AttendanceRequestWFH || item.RequestType == domain.AttendanceRequestRemoteWork) {
		mode := domain.AttendanceWorkModeRemote
		workMode = &mode
	}
	if item.RequestedCheckInAt != nil {
		if _, err := s.attendances.CreateAttendance(ctx, &domain.Attendance{TenantID: item.TenantID, UserID: item.UserID, Date: item.Date, Time: item.RequestedCheckInAt, Type: strPtr(domain.AttendanceCheckin), Status: &status, WorkMode: workMode, Remarks: item.Reason}, actorID); err != nil {
			return err
		}
	}
	if item.RequestedCheckOutAt != nil {
		if _, err := s.attendances.CreateAttendance(ctx, &domain.Attendance{TenantID: item.TenantID, UserID: item.UserID, Date: item.Date, Time: item.RequestedCheckOutAt, Type: strPtr(domain.AttendanceCheckout), Status: &status, WorkMode: workMode, Remarks: item.Reason}, actorID); err != nil {
			return err
		}
	}
	if item.RequestedCheckInAt == nil && item.RequestedCheckOutAt == nil {
		now := time.Now().UTC()
		_, err := s.attendances.CreateAttendance(ctx, &domain.Attendance{TenantID: item.TenantID, UserID: item.UserID, Date: item.Date, Time: &now, Status: &status, WorkMode: workMode, Remarks: item.Reason}, actorID)
		return err
	}
	return nil
}

func parseAttendancePolicyOptionalDate(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return nil, domain.ErrInvalidAttendanceDate
	}
	return &parsed, nil
}
func parseOptionalRFC3339(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*value))
	if err != nil {
		return nil, domain.ErrInvalidAttendanceDate
	}
	parsed = parsed.UTC()
	return &parsed, nil
}
