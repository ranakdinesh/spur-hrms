package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateAttendancePolicy(ctx context.Context, item *domain.AttendancePolicy, actorID *uuid.UUID) (*domain.AttendancePolicy, error) {
	row, err := s.getQueries(ctx).CreateAttendancePolicy(ctx, attendancePolicyCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance policy", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapAttendancePolicy(row), nil
}

func (s *Store) ListAttendancePolicies(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendancePolicy, error) {
	rows, err := s.getQueries(ctx).ListAttendancePolicies(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance policies", err, tenantIDField(tenantID))
	}
	return mapAttendancePolicies(rows), nil
}

func (s *Store) GetAttendancePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendancePolicy, error) {
	row, err := s.getQueries(ctx).GetAttendancePolicy(ctx, sqlc.GetAttendancePolicyParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendancePolicyNotFound
		}
		return nil, s.logDBError(ctx, "get attendance policy", err, tenantIDField(tenantID), stringField("policy_id", id.String()))
	}
	return mapAttendancePolicy(row), nil
}

func (s *Store) ResolveAttendancePolicy(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, departmentID *uuid.UUID, branchID *uuid.UUID, date string) (*domain.AttendancePolicy, error) {
	parsed, err := parseDateKey(date)
	if err != nil {
		return nil, err
	}
	row, err := s.getQueries(ctx).ResolveAttendancePolicy(ctx, sqlc.ResolveAttendancePolicyParams{TenantID: tenantID, UserID: uuidFromPtr(&userID), DepartmentID: uuidFromPtr(departmentID), BranchID: uuidFromPtr(branchID), EffectiveFrom: dateFromTime(parsed)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendancePolicyNotFound
		}
		return nil, s.logDBError(ctx, "resolve attendance policy", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapAttendancePolicy(row), nil
}

func (s *Store) UpdateAttendancePolicy(ctx context.Context, item *domain.AttendancePolicy, actorID *uuid.UUID) (*domain.AttendancePolicy, error) {
	params := attendancePolicyUpdateParams(item, actorID)
	row, err := s.getQueries(ctx).UpdateAttendancePolicy(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendancePolicyNotFound
		}
		return nil, s.logDBError(ctx, "update attendance policy", err, tenantIDField(item.TenantID), stringField("policy_id", item.ID.String()))
	}
	return mapAttendancePolicy(row), nil
}

func (s *Store) DeleteAttendancePolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAttendancePolicy(ctx, sqlc.SoftDeleteAttendancePolicyParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete attendance policy", err, tenantIDField(tenantID), stringField("policy_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAttendanceRoster(ctx context.Context, item *domain.AttendanceRoster, actorID *uuid.UUID) (*domain.AttendanceRoster, error) {
	row, err := s.getQueries(ctx).CreateAttendanceRoster(ctx, sqlc.CreateAttendanceRosterParams{TenantID: item.TenantID, UserID: item.UserID, PolicyID: uuidFromPtr(item.PolicyID), Date: dateFromTime(item.Date), StartTime: clockTimeFromPtr(item.StartTime), EndTime: clockTimeFromPtr(item.EndTime), BreakMinutes: item.BreakMinutes, WorkMode: item.WorkMode, LocationType: item.LocationType, Remarks: textFromPtr(item.Remarks), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance roster", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapAttendanceRoster(row), nil
}

func (s *Store) ListAttendanceRostersByDateRange(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceRoster, error) {
	start, err := parseDateKey(startDate)
	if err != nil {
		return nil, err
	}
	end, err := parseDateKey(endDate)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListAttendanceRostersByDateRange(ctx, sqlc.ListAttendanceRostersByDateRangeParams{TenantID: tenantID, Date: dateFromTime(start), Date_2: dateFromTime(end)})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance rosters", err, tenantIDField(tenantID))
	}
	return mapAttendanceRosters(rows), nil
}

func (s *Store) ListAttendanceRostersByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceRoster, error) {
	start, err := parseDateKey(startDate)
	if err != nil {
		return nil, err
	}
	end, err := parseDateKey(endDate)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListAttendanceRostersByUser(ctx, sqlc.ListAttendanceRostersByUserParams{TenantID: tenantID, UserID: userID, Date: dateFromTime(start), Date_2: dateFromTime(end)})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance rosters by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapAttendanceRosters(rows), nil
}

func (s *Store) GetAttendanceRosterByUserDate(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, date string) (*domain.AttendanceRoster, error) {
	parsed, err := parseDateKey(date)
	if err != nil {
		return nil, err
	}
	row, err := s.getQueries(ctx).GetAttendanceRosterByUserDate(ctx, sqlc.GetAttendanceRosterByUserDateParams{TenantID: tenantID, UserID: userID, Date: dateFromTime(parsed)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceRosterNotFound
		}
		return nil, s.logDBError(ctx, "get attendance roster by user date", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapAttendanceRoster(row), nil
}

func (s *Store) UpdateAttendanceRoster(ctx context.Context, item *domain.AttendanceRoster, actorID *uuid.UUID) (*domain.AttendanceRoster, error) {
	row, err := s.getQueries(ctx).UpdateAttendanceRoster(ctx, sqlc.UpdateAttendanceRosterParams{TenantID: item.TenantID, ID: item.ID, PolicyID: uuidFromPtr(item.PolicyID), Date: dateFromTime(item.Date), StartTime: clockTimeFromPtr(item.StartTime), EndTime: clockTimeFromPtr(item.EndTime), BreakMinutes: item.BreakMinutes, WorkMode: item.WorkMode, LocationType: item.LocationType, Remarks: textFromPtr(item.Remarks), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceRosterNotFound
		}
		return nil, s.logDBError(ctx, "update attendance roster", err, tenantIDField(item.TenantID), stringField("roster_id", item.ID.String()))
	}
	return mapAttendanceRoster(row), nil
}

func (s *Store) DeleteAttendanceRoster(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAttendanceRoster(ctx, sqlc.SoftDeleteAttendanceRosterParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete attendance roster", err, tenantIDField(tenantID), stringField("roster_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAttendanceRequest(ctx context.Context, item *domain.AttendanceRequest, actorID *uuid.UUID) (*domain.AttendanceRequest, error) {
	row, err := s.getQueries(ctx).CreateAttendanceRequest(ctx, sqlc.CreateAttendanceRequestParams{TenantID: item.TenantID, UserID: item.UserID, Date: dateFromTime(item.Date), RequestedType: textFromPtr(item.RequestedType), RequestType: item.RequestType, RequestedCheckinAt: timestamptzFromPtr(item.RequestedCheckInAt), RequestedCheckoutAt: timestamptzFromPtr(item.RequestedCheckOutAt), RequestedWorkMode: textFromPtr(item.RequestedWorkMode), PolicyID: uuidFromPtr(item.PolicyID), RosterID: uuidFromPtr(item.RosterID), Reason: textFromPtr(item.Reason), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance request", err, tenantIDField(item.TenantID), stringField("user_id", item.UserID.String()))
	}
	return mapAttendanceRequest(row), nil
}

func (s *Store) ListAttendanceRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.AttendanceRequest, error) {
	rows, err := s.getQueries(ctx).ListAttendanceRequestsByUser(ctx, sqlc.ListAttendanceRequestsByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance requests by user", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapAttendanceRequests(rows), nil
}

func (s *Store) ListAttendanceRequestsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.AttendanceRequest, error) {
	rows, err := s.getQueries(ctx).ListAttendanceRequestsByStatus(ctx, sqlc.ListAttendanceRequestsByStatusParams{TenantID: tenantID, Status: status})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance requests by status", err, tenantIDField(tenantID), stringField("status", status))
	}
	return mapAttendanceRequests(rows), nil
}

func (s *Store) GetAttendanceRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendanceRequest, error) {
	row, err := s.getQueries(ctx).GetAttendanceRequest(ctx, sqlc.GetAttendanceRequestParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceRequestNotFound
		}
		return nil, s.logDBError(ctx, "get attendance request", err, tenantIDField(tenantID), stringField("request_id", id.String()))
	}
	return mapAttendanceRequest(row), nil
}

func (s *Store) UpdateAttendanceRequestReview(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, reviewerID uuid.UUID, remarks *string) (*domain.AttendanceRequest, error) {
	row, err := s.getQueries(ctx).UpdateAttendanceRequestReview(ctx, sqlc.UpdateAttendanceRequestReviewParams{TenantID: tenantID, ID: id, Status: status, ReviewedBy: uuidFromPtr(&reviewerID), Remarks: textFromPtr(remarks)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceRequestNotFound
		}
		return nil, s.logDBError(ctx, "review attendance request", err, tenantIDField(tenantID), stringField("request_id", id.String()))
	}
	return mapAttendanceRequest(row), nil
}

func (s *Store) SetAttendanceRequestWorkflow(ctx context.Context, tenantID uuid.UUID, requestID uuid.UUID, workflowID *uuid.UUID, routeMode *string, escalationDueAt *string, payrollBlocking bool, actorID *uuid.UUID) (*domain.AttendanceRequest, error) {
	var dueAt pgtype.Timestamptz
	if escalationDueAt != nil && *escalationDueAt != "" {
		parsed, err := time.Parse(time.RFC3339, *escalationDueAt)
		if err != nil {
			return nil, err
		}
		dueAt = pgtype.Timestamptz{Time: parsed.UTC(), Valid: true}
	}
	row, err := s.getQueries(ctx).SetAttendanceRequestWorkflow(ctx, sqlc.SetAttendanceRequestWorkflowParams{TenantID: tenantID, ID: requestID, WorkflowID: uuidFromPtr(workflowID), RouteMode: textFromPtr(routeMode), EscalationDueAt: dueAt, PayrollBlocking: payrollBlocking, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceRequestNotFound
		}
		return nil, s.logDBError(ctx, "set attendance request workflow", err, tenantIDField(tenantID), stringField("request_id", requestID.String()))
	}
	return mapAttendanceRequest(row), nil
}

func (s *Store) DeleteAttendanceRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAttendanceRequest(ctx, sqlc.SoftDeleteAttendanceRequestParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete attendance request", err, tenantIDField(tenantID), stringField("request_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAttendanceExceptionWorkflow(ctx context.Context, item *domain.AttendanceExceptionWorkflow, actorID *uuid.UUID) (*domain.AttendanceExceptionWorkflow, error) {
	row, err := s.getQueries(ctx).CreateAttendanceExceptionWorkflow(ctx, sqlc.CreateAttendanceExceptionWorkflowParams{TenantID: item.TenantID, Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), BranchID: uuidFromPtr(item.BranchID), DepartmentID: uuidFromPtr(item.DepartmentID), RequestType: item.RequestType, RouteMode: item.RouteMode, MaxRequestsPerMonth: item.MaxRequestsPerMonth, EscalationHours: item.EscalationHours, EscalationRouteMode: textFromPtr(item.EscalationRouteMode), BlockPayrollWhenPending: item.BlockPayrollWhenPending, AutoApprove: item.AutoApprove, IsActive: item.IsActive, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance exception workflow", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapAttendanceExceptionWorkflow(row), nil
}

func (s *Store) ListAttendanceExceptionWorkflows(ctx context.Context, tenantID uuid.UUID) ([]*domain.AttendanceExceptionWorkflow, error) {
	rows, err := s.getQueries(ctx).ListAttendanceExceptionWorkflows(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance exception workflows", err, tenantIDField(tenantID))
	}
	return mapAttendanceExceptionWorkflows(rows), nil
}

func (s *Store) GetAttendanceExceptionWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AttendanceExceptionWorkflow, error) {
	row, err := s.getQueries(ctx).GetAttendanceExceptionWorkflow(ctx, sqlc.GetAttendanceExceptionWorkflowParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceExceptionWorkflowNotFound
		}
		return nil, s.logDBError(ctx, "get attendance exception workflow", err, tenantIDField(tenantID), stringField("workflow_id", id.String()))
	}
	return mapAttendanceExceptionWorkflow(row), nil
}

func (s *Store) ResolveAttendanceExceptionWorkflow(ctx context.Context, tenantID uuid.UUID, requestType string, departmentID *uuid.UUID, branchID *uuid.UUID) (*domain.AttendanceExceptionWorkflow, error) {
	row, err := s.getQueries(ctx).ResolveAttendanceExceptionWorkflow(ctx, sqlc.ResolveAttendanceExceptionWorkflowParams{TenantID: tenantID, RequestType: requestType, DepartmentID: uuidFromPtr(departmentID), BranchID: uuidFromPtr(branchID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceExceptionWorkflowNotFound
		}
		return nil, s.logDBError(ctx, "resolve attendance exception workflow", err, tenantIDField(tenantID), stringField("request_type", requestType))
	}
	return mapAttendanceExceptionWorkflow(row), nil
}

func (s *Store) UpdateAttendanceExceptionWorkflow(ctx context.Context, item *domain.AttendanceExceptionWorkflow, actorID *uuid.UUID) (*domain.AttendanceExceptionWorkflow, error) {
	row, err := s.getQueries(ctx).UpdateAttendanceExceptionWorkflow(ctx, sqlc.UpdateAttendanceExceptionWorkflowParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), BranchID: uuidFromPtr(item.BranchID), DepartmentID: uuidFromPtr(item.DepartmentID), RequestType: item.RequestType, RouteMode: item.RouteMode, MaxRequestsPerMonth: item.MaxRequestsPerMonth, EscalationHours: item.EscalationHours, EscalationRouteMode: textFromPtr(item.EscalationRouteMode), BlockPayrollWhenPending: item.BlockPayrollWhenPending, AutoApprove: item.AutoApprove, IsActive: item.IsActive, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttendanceExceptionWorkflowNotFound
		}
		return nil, s.logDBError(ctx, "update attendance exception workflow", err, tenantIDField(item.TenantID), stringField("workflow_id", item.ID.String()))
	}
	return mapAttendanceExceptionWorkflow(row), nil
}

func (s *Store) DeleteAttendanceExceptionWorkflow(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAttendanceExceptionWorkflow(ctx, sqlc.SoftDeleteAttendanceExceptionWorkflowParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete attendance exception workflow", err, tenantIDField(tenantID), stringField("workflow_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAttendanceExceptionEvent(ctx context.Context, item *domain.AttendanceExceptionEvent, actorID *uuid.UUID) (*domain.AttendanceExceptionEvent, error) {
	metadata := item.Metadata
	if len(metadata) == 0 {
		metadata = []byte(`{}`)
	}
	row, err := s.getQueries(ctx).CreateAttendanceExceptionEvent(ctx, sqlc.CreateAttendanceExceptionEventParams{TenantID: item.TenantID, AttendanceRequestID: item.AttendanceRequestID, WorkflowID: uuidFromPtr(item.WorkflowID), Action: item.Action, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), RoutedTo: textFromPtr(item.RoutedTo), Remarks: textFromPtr(item.Remarks), Metadata: metadata, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create attendance exception event", err, tenantIDField(item.TenantID), stringField("request_id", item.AttendanceRequestID.String()))
	}
	return mapAttendanceExceptionEvent(row), nil
}

func (s *Store) ListAttendanceExceptionEvents(ctx context.Context, tenantID uuid.UUID, attendanceRequestID uuid.UUID) ([]*domain.AttendanceExceptionEvent, error) {
	rows, err := s.getQueries(ctx).ListAttendanceExceptionEvents(ctx, sqlc.ListAttendanceExceptionEventsParams{TenantID: tenantID, AttendanceRequestID: attendanceRequestID})
	if err != nil {
		return nil, s.logDBError(ctx, "list attendance exception events", err, tenantIDField(tenantID), stringField("request_id", attendanceRequestID.String()))
	}
	return mapAttendanceExceptionEvents(rows), nil
}

func (s *Store) ListPayrollBlockingAttendanceRequests(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.AttendanceRequest, error) {
	start, err := parseDateKey(startDate)
	if err != nil {
		return nil, err
	}
	end, err := parseDateKey(endDate)
	if err != nil {
		return nil, err
	}
	rows, err := s.getQueries(ctx).ListPayrollBlockingAttendanceRequests(ctx, sqlc.ListPayrollBlockingAttendanceRequestsParams{TenantID: tenantID, Date: dateFromTime(start), Date_2: dateFromTime(end)})
	if err != nil {
		return nil, s.logDBError(ctx, "list payroll blocking attendance requests", err, tenantIDField(tenantID))
	}
	return mapAttendanceRequests(rows), nil
}

func attendancePolicyCreateParams(item *domain.AttendancePolicy, actorID *uuid.UUID) sqlc.CreateAttendancePolicyParams {
	return sqlc.CreateAttendancePolicyParams{TenantID: item.TenantID, Name: item.Name, Code: item.Code, Description: textFromPtr(item.Description), BranchID: uuidFromPtr(item.BranchID), DepartmentID: uuidFromPtr(item.DepartmentID), UserID: uuidFromPtr(item.UserID), ScheduleType: item.ScheduleType, IsDefault: item.IsDefault, StandardWorkMinutes: item.StandardWorkMinutes, MinHalfDayMinutes: item.MinHalfDayMinutes, MinFullDayMinutes: item.MinFullDayMinutes, GraceLateMinutes: item.GraceLateMinutes, GraceEarlyMinutes: item.GraceEarlyMinutes, HalfDayLateAfterMinutes: int4FromPtr(item.HalfDayLateAfterMinutes), AbsentLateAfterMinutes: int4FromPtr(item.AbsentLateAfterMinutes), HalfDayEarlyBeforeMinutes: int4FromPtr(item.HalfDayEarlyBeforeMinutes), AbsentEarlyBeforeMinutes: int4FromPtr(item.AbsentEarlyBeforeMinutes), AllowFlexiHours: item.AllowFlexiHours, CoreStartTime: clockTimeFromPtr(item.CoreStartTime), CoreEndTime: clockTimeFromPtr(item.CoreEndTime), AllowWfh: item.AllowWFH, WfhDaysPerWeek: item.WFHDaysPerWeek, AllowPermanentRemote: item.AllowPermanentRemote, RequireGeo: item.RequireGeo, RequireDevice: item.RequireDevice, RegularizationWindowDays: item.RegularizationWindowDays, MaxRegularizationsPerMonth: item.MaxRegularizationsPerMonth, ApprovalMode: item.ApprovalMode, EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), CreatedBy: uuidFromPtr(actorID)}
}

func attendancePolicyUpdateParams(item *domain.AttendancePolicy, actorID *uuid.UUID) sqlc.UpdateAttendancePolicyParams {
	p := attendancePolicyCreateParams(item, actorID)
	return sqlc.UpdateAttendancePolicyParams{TenantID: p.TenantID, ID: item.ID, Name: p.Name, Code: p.Code, Description: p.Description, BranchID: p.BranchID, DepartmentID: p.DepartmentID, UserID: p.UserID, ScheduleType: p.ScheduleType, IsDefault: p.IsDefault, StandardWorkMinutes: p.StandardWorkMinutes, MinHalfDayMinutes: p.MinHalfDayMinutes, MinFullDayMinutes: p.MinFullDayMinutes, GraceLateMinutes: p.GraceLateMinutes, GraceEarlyMinutes: p.GraceEarlyMinutes, HalfDayLateAfterMinutes: p.HalfDayLateAfterMinutes, AbsentLateAfterMinutes: p.AbsentLateAfterMinutes, HalfDayEarlyBeforeMinutes: p.HalfDayEarlyBeforeMinutes, AbsentEarlyBeforeMinutes: p.AbsentEarlyBeforeMinutes, AllowFlexiHours: p.AllowFlexiHours, CoreStartTime: p.CoreStartTime, CoreEndTime: p.CoreEndTime, AllowWfh: p.AllowWfh, WfhDaysPerWeek: p.WfhDaysPerWeek, AllowPermanentRemote: p.AllowPermanentRemote, RequireGeo: p.RequireGeo, RequireDevice: p.RequireDevice, RegularizationWindowDays: p.RegularizationWindowDays, MaxRegularizationsPerMonth: p.MaxRegularizationsPerMonth, ApprovalMode: p.ApprovalMode, EffectiveFrom: p.EffectiveFrom, EffectiveTo: p.EffectiveTo, UpdatedBy: uuidFromPtr(actorID)}
}

var _ = time.UTC
