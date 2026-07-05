package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapAttendancePolicy(row sqlc.HrmsAttendancePolicy) *domain.AttendancePolicy {
	return &domain.AttendancePolicy{ID: row.ID, TenantID: row.TenantID, Name: row.Name, Code: row.Code, Description: ptrFromText(row.Description), BranchID: ptrFromUUID(row.BranchID), DepartmentID: ptrFromUUID(row.DepartmentID), UserID: ptrFromUUID(row.UserID), ScheduleType: row.ScheduleType, IsDefault: row.IsDefault, StandardWorkMinutes: row.StandardWorkMinutes, MinHalfDayMinutes: row.MinHalfDayMinutes, MinFullDayMinutes: row.MinFullDayMinutes, GraceLateMinutes: row.GraceLateMinutes, GraceEarlyMinutes: row.GraceEarlyMinutes, HalfDayLateAfterMinutes: ptrFromInt4(row.HalfDayLateAfterMinutes), AbsentLateAfterMinutes: ptrFromInt4(row.AbsentLateAfterMinutes), HalfDayEarlyBeforeMinutes: ptrFromInt4(row.HalfDayEarlyBeforeMinutes), AbsentEarlyBeforeMinutes: ptrFromInt4(row.AbsentEarlyBeforeMinutes), AllowFlexiHours: row.AllowFlexiHours, CoreStartTime: ptrFromClockTime(row.CoreStartTime), CoreEndTime: ptrFromClockTime(row.CoreEndTime), AllowWFH: row.AllowWfh, WFHDaysPerWeek: row.WfhDaysPerWeek, AllowPermanentRemote: row.AllowPermanentRemote, RequireGeo: row.RequireGeo, RequireDevice: row.RequireDevice, RegularizationWindowDays: row.RegularizationWindowDays, MaxRegularizationsPerMonth: row.MaxRegularizationsPerMonth, ApprovalMode: row.ApprovalMode, EffectiveFrom: ptrFromDate(row.EffectiveFrom), EffectiveTo: ptrFromDate(row.EffectiveTo), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAttendancePolicies(rows []sqlc.HrmsAttendancePolicy) []*domain.AttendancePolicy {
	items := make([]*domain.AttendancePolicy, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAttendancePolicy(row))
	}
	return items
}

func mapAttendanceRoster(row sqlc.HrmsAttendanceRoster) *domain.AttendanceRoster {
	return &domain.AttendanceRoster{ID: row.ID, TenantID: row.TenantID, UserID: row.UserID, PolicyID: ptrFromUUID(row.PolicyID), Date: timeFromDate(row.Date), StartTime: ptrFromClockTime(row.StartTime), EndTime: ptrFromClockTime(row.EndTime), BreakMinutes: row.BreakMinutes, WorkMode: row.WorkMode, LocationType: row.LocationType, Remarks: ptrFromText(row.Remarks), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAttendanceRosters(rows []sqlc.HrmsAttendanceRoster) []*domain.AttendanceRoster {
	items := make([]*domain.AttendanceRoster, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAttendanceRoster(row))
	}
	return items
}

func mapAttendanceRequest(row sqlc.HrmsAttendanceRequest) *domain.AttendanceRequest {
	return &domain.AttendanceRequest{ID: row.ID, TenantID: row.TenantID, UserID: row.UserID, Date: timeFromDate(row.Date), RequestedType: ptrFromText(row.RequestedType), RequestType: row.RequestType, RequestedCheckInAt: ptrFromTimestamptz(row.RequestedCheckinAt), RequestedCheckOutAt: ptrFromTimestamptz(row.RequestedCheckoutAt), RequestedWorkMode: ptrFromText(row.RequestedWorkMode), PolicyID: ptrFromUUID(row.PolicyID), RosterID: ptrFromUUID(row.RosterID), Reason: ptrFromText(row.Reason), Status: row.Status, ReviewedBy: ptrFromUUID(row.ReviewedBy), ReviewedAt: ptrFromTimestamptz(row.ReviewedAt), Remarks: ptrFromText(row.Remarks), WorkflowID: ptrFromUUID(row.WorkflowID), RouteMode: ptrFromText(row.RouteMode), EscalationDueAt: ptrFromTimestamptz(row.EscalationDueAt), EscalatedAt: ptrFromTimestamptz(row.EscalatedAt), PayrollBlocking: row.PayrollBlocking, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAttendanceRequests(rows []sqlc.HrmsAttendanceRequest) []*domain.AttendanceRequest {
	items := make([]*domain.AttendanceRequest, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAttendanceRequest(row))
	}
	return items
}

func int4FromPtr(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *value, Valid: true}
}
func ptrFromInt4(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	clean := value.Int32
	return &clean
}

func mapAttendanceExceptionWorkflow(row sqlc.HrmsAttendanceExceptionWorkflow) *domain.AttendanceExceptionWorkflow {
	return &domain.AttendanceExceptionWorkflow{ID: row.ID, TenantID: row.TenantID, Code: row.Code, Name: row.Name, Description: ptrFromText(row.Description), BranchID: ptrFromUUID(row.BranchID), DepartmentID: ptrFromUUID(row.DepartmentID), RequestType: row.RequestType, RouteMode: row.RouteMode, MaxRequestsPerMonth: row.MaxRequestsPerMonth, EscalationHours: row.EscalationHours, EscalationRouteMode: ptrFromText(row.EscalationRouteMode), BlockPayrollWhenPending: row.BlockPayrollWhenPending, AutoApprove: row.AutoApprove, IsActive: row.IsActive, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAttendanceExceptionWorkflows(rows []sqlc.HrmsAttendanceExceptionWorkflow) []*domain.AttendanceExceptionWorkflow {
	items := make([]*domain.AttendanceExceptionWorkflow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAttendanceExceptionWorkflow(row))
	}
	return items
}

func mapAttendanceExceptionEvent(row sqlc.HrmsAttendanceExceptionEvent) *domain.AttendanceExceptionEvent {
	return &domain.AttendanceExceptionEvent{ID: row.ID, TenantID: row.TenantID, AttendanceRequestID: row.AttendanceRequestID, WorkflowID: ptrFromUUID(row.WorkflowID), Action: row.Action, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), RoutedTo: ptrFromText(row.RoutedTo), Remarks: ptrFromText(row.Remarks), Metadata: row.Metadata, Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAttendanceExceptionEvents(rows []sqlc.HrmsAttendanceExceptionEvent) []*domain.AttendanceExceptionEvent {
	items := make([]*domain.AttendanceExceptionEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAttendanceExceptionEvent(row))
	}
	return items
}
func clockTimeFromPtr(value *string) pgtype.Time {
	if value == nil {
		return pgtype.Time{Valid: false}
	}
	return timeFromClockString(*value)
}
func ptrFromClockTime(value pgtype.Time) *string {
	if !value.Valid {
		return nil
	}
	clean := clockStringFromTime(value)
	return &clean
}
func parseDateKey(value string) (time.Time, error) { return time.Parse("2006-01-02", value) }
