package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateShiftTemplate(ctx context.Context, item *domain.ShiftTemplate, actorID *uuid.UUID) (*domain.ShiftTemplate, error) {
	row, err := s.getQueries(ctx).CreateShiftTemplate(ctx, sqlc.CreateShiftTemplateParams{TenantID: item.TenantID, Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), StartTime: timeFromClockString(item.StartTime), EndTime: timeFromClockString(item.EndTime), BreakMinutes: item.BreakMinutes, PaidMinutes: item.PaidMinutes, WorkMode: item.WorkMode, LocationType: item.LocationType, AttendancePolicyID: uuidFromPtr(item.AttendancePolicyID), AttendanceLocationID: uuidFromPtr(item.AttendanceLocationID), AllowOvertime: item.AllowOvertime, PayrollCode: textFromPtr(item.PayrollCode), Metadata: jsonBytesFromRaw(item.Metadata), IsActive: item.IsActive, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create shift template", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapShiftTemplate(row), nil
}

func (s *Store) UpdateShiftTemplate(ctx context.Context, item *domain.ShiftTemplate, actorID *uuid.UUID) (*domain.ShiftTemplate, error) {
	row, err := s.getQueries(ctx).UpdateShiftTemplate(ctx, sqlc.UpdateShiftTemplateParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), StartTime: timeFromClockString(item.StartTime), EndTime: timeFromClockString(item.EndTime), BreakMinutes: item.BreakMinutes, PaidMinutes: item.PaidMinutes, WorkMode: item.WorkMode, LocationType: item.LocationType, AttendancePolicyID: uuidFromPtr(item.AttendancePolicyID), AttendanceLocationID: uuidFromPtr(item.AttendanceLocationID), AllowOvertime: item.AllowOvertime, PayrollCode: textFromPtr(item.PayrollCode), Metadata: jsonBytesFromRaw(item.Metadata), IsActive: item.IsActive, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrShiftTemplateNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update shift template", err, tenantIDField(item.TenantID), stringField("shift_template_id", item.ID.String()))
	}
	return mapShiftTemplate(row), nil
}

func (s *Store) ListShiftTemplates(ctx context.Context, tenantID uuid.UUID, activeOnly *bool, search *string, limit int32, offset int32) ([]*domain.ShiftTemplate, error) {
	rows, err := s.getQueries(ctx).ListShiftTemplates(ctx, sqlc.ListShiftTemplatesParams{TenantID: tenantID, Limit: limitOrDefault(limit), Offset: offset, ActiveOnly: boolFromPtr(activeOnly), Search: textFromPtr(search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list shift templates", err, tenantIDField(tenantID))
	}
	return mapShiftTemplates(rows), nil
}

func (s *Store) GetShiftTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ShiftTemplate, error) {
	row, err := s.getQueries(ctx).GetShiftTemplate(ctx, sqlc.GetShiftTemplateParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrShiftTemplateNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get shift template", err, tenantIDField(tenantID), stringField("shift_template_id", id.String()))
	}
	return mapShiftTemplate(row), nil
}

func (s *Store) DeleteShiftTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteShiftTemplate(ctx, sqlc.SoftDeleteShiftTemplateParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete shift template", err, tenantIDField(tenantID), stringField("shift_template_id", id.String()))
	}
	return nil
}

func (s *Store) CreateStaffingRequirement(ctx context.Context, item *domain.StaffingRequirement, actorID *uuid.UUID) (*domain.StaffingRequirement, error) {
	row, err := s.getQueries(ctx).CreateStaffingRequirement(ctx, staffingRequirementCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create staffing requirement", err, tenantIDField(item.TenantID), stringField("name", item.Name))
	}
	return mapStaffingRequirement(row), nil
}

func (s *Store) UpdateStaffingRequirement(ctx context.Context, item *domain.StaffingRequirement, actorID *uuid.UUID) (*domain.StaffingRequirement, error) {
	params := staffingRequirementUpdateParams(item, actorID)
	row, err := s.getQueries(ctx).UpdateStaffingRequirement(ctx, params)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrStaffingRequirementNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update staffing requirement", err, tenantIDField(item.TenantID), stringField("staffing_requirement_id", item.ID.String()))
	}
	return mapStaffingRequirement(row), nil
}

func (s *Store) ListStaffingRequirements(ctx context.Context, filter domain.StaffingRequirementFilter) ([]*domain.StaffingRequirement, error) {
	rows, err := s.getQueries(ctx).ListStaffingRequirements(ctx, sqlc.ListStaffingRequirementsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, Status: textFromPtr(filter.Status), StartDate: dateFromString(filter.StartDate), EndDate: dateFromString(filter.EndDate)})
	if err != nil {
		return nil, s.logDBError(ctx, "list staffing requirements", err, tenantIDField(filter.TenantID))
	}
	return mapStaffingRequirements(rows), nil
}

func (s *Store) GetStaffingRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.StaffingRequirement, error) {
	row, err := s.getQueries(ctx).GetStaffingRequirement(ctx, sqlc.GetStaffingRequirementParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrStaffingRequirementNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get staffing requirement", err, tenantIDField(tenantID), stringField("staffing_requirement_id", id.String()))
	}
	return mapStaffingRequirement(row), nil
}

func (s *Store) DeleteStaffingRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteStaffingRequirement(ctx, sqlc.SoftDeleteStaffingRequirementParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete staffing requirement", err, tenantIDField(tenantID), stringField("staffing_requirement_id", id.String()))
	}
	return nil
}

func (s *Store) CreateShiftScheduleAssignment(ctx context.Context, item *domain.ShiftScheduleAssignment, actorID *uuid.UUID) (*domain.ShiftScheduleAssignment, error) {
	row, err := s.getQueries(ctx).CreateShiftScheduleAssignment(ctx, shiftAssignmentCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create shift schedule assignment", err, tenantIDField(item.TenantID))
	}
	return mapShiftScheduleAssignment(row), nil
}

func (s *Store) UpdateShiftScheduleAssignment(ctx context.Context, item *domain.ShiftScheduleAssignment, actorID *uuid.UUID) (*domain.ShiftScheduleAssignment, error) {
	row, err := s.getQueries(ctx).UpdateShiftScheduleAssignment(ctx, shiftAssignmentUpdateParams(item, actorID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrShiftAssignmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update shift schedule assignment", err, tenantIDField(item.TenantID), stringField("shift_assignment_id", item.ID.String()))
	}
	return mapShiftScheduleAssignment(row), nil
}

func (s *Store) UpdateShiftScheduleAssignmentStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, hasConflict bool, conflictReason *string, payrollBlocking bool, actorID *uuid.UUID) (*domain.ShiftScheduleAssignment, error) {
	row, err := s.getQueries(ctx).UpdateShiftScheduleAssignmentStatus(ctx, sqlc.UpdateShiftScheduleAssignmentStatusParams{TenantID: tenantID, ID: id, Status: status, HasConflict: hasConflict, ConflictReason: textFromPtr(conflictReason), PayrollBlocking: payrollBlocking, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrShiftAssignmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update shift schedule assignment status", err, tenantIDField(tenantID), stringField("shift_assignment_id", id.String()), stringField("status", status))
	}
	return mapShiftScheduleAssignment(row), nil
}

func (s *Store) ListShiftScheduleAssignments(ctx context.Context, filter domain.ShiftScheduleFilter) ([]*domain.ShiftScheduleAssignment, error) {
	rows, err := s.getQueries(ctx).ListShiftScheduleAssignments(ctx, sqlc.ListShiftScheduleAssignmentsParams{TenantID: filter.TenantID, ScheduleDate: dateFromString(filter.StartDate), ScheduleDate_2: dateFromString(filter.EndDate), Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, Status: textFromPtr(filter.Status), WorkerProfileID: uuidFromPtr(filter.WorkerProfileID), EmployeeUserID: uuidFromPtr(filter.EmployeeUserID), BranchID: uuidFromPtr(filter.BranchID), DepartmentID: uuidFromPtr(filter.DepartmentID), AttendanceLocationID: uuidFromPtr(filter.AttendanceLocationID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list shift schedule assignments", err, tenantIDField(filter.TenantID))
	}
	return mapShiftScheduleAssignments(rows), nil
}

func (s *Store) GetShiftScheduleAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ShiftScheduleAssignment, error) {
	row, err := s.getQueries(ctx).GetShiftScheduleAssignment(ctx, sqlc.GetShiftScheduleAssignmentParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrShiftAssignmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get shift schedule assignment", err, tenantIDField(tenantID), stringField("shift_assignment_id", id.String()))
	}
	return mapShiftScheduleAssignment(row), nil
}

func (s *Store) ListShiftAssignmentsForWorkerDate(ctx context.Context, tenantID uuid.UUID, date string, workerProfileID *uuid.UUID, employeeUserID *uuid.UUID) ([]*domain.ShiftScheduleAssignment, error) {
	rows, err := s.getQueries(ctx).ListShiftAssignmentsForWorkerDate(ctx, sqlc.ListShiftAssignmentsForWorkerDateParams{TenantID: tenantID, ScheduleDate: dateFromString(date), WorkerProfileID: uuidFromPtr(workerProfileID), EmployeeUserID: uuidFromPtr(employeeUserID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list shift assignments for worker date", err, tenantIDField(tenantID))
	}
	return mapShiftScheduleAssignments(rows), nil
}

func (s *Store) DeleteShiftScheduleAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteShiftScheduleAssignment(ctx, sqlc.SoftDeleteShiftScheduleAssignmentParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete shift schedule assignment", err, tenantIDField(tenantID), stringField("shift_assignment_id", id.String()))
	}
	return nil
}

func (s *Store) CreateShiftSwapRequest(ctx context.Context, item *domain.ShiftSwapRequest, actorID *uuid.UUID) (*domain.ShiftSwapRequest, error) {
	row, err := s.getQueries(ctx).CreateShiftSwapRequest(ctx, sqlc.CreateShiftSwapRequestParams{TenantID: item.TenantID, RequesterAssignmentID: item.RequesterAssignmentID, RequesterWorkerProfileID: uuidFromPtr(item.RequesterWorkerProfileID), RequesterUserID: uuidFromPtr(item.RequesterUserID), TargetWorkerProfileID: uuidFromPtr(item.TargetWorkerProfileID), TargetUserID: uuidFromPtr(item.TargetUserID), OfferedAssignmentID: uuidFromPtr(item.OfferedAssignmentID), RequestedDate: dateFromPtr(item.RequestedDate), RequestedShiftTemplateID: uuidFromPtr(item.RequestedShiftTemplateID), Reason: textFromPtr(item.Reason), Status: item.Status, PayrollBlocking: item.PayrollBlocking, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create shift swap request", err, tenantIDField(item.TenantID), stringField("assignment_id", item.RequesterAssignmentID.String()))
	}
	return mapShiftSwapRequest(row), nil
}

func (s *Store) UpdateShiftSwapRequestReview(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, reviewedBy uuid.UUID, remarks *string, payrollBlocking bool) (*domain.ShiftSwapRequest, error) {
	row, err := s.getQueries(ctx).UpdateShiftSwapRequestReview(ctx, sqlc.UpdateShiftSwapRequestReviewParams{TenantID: tenantID, ID: id, Status: status, ReviewedBy: uuidFromPtr(&reviewedBy), ReviewRemarks: textFromPtr(remarks), PayrollBlocking: payrollBlocking})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrShiftSwapRequestNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "review shift swap request", err, tenantIDField(tenantID), stringField("shift_swap_id", id.String()))
	}
	return mapShiftSwapRequest(row), nil
}

func (s *Store) ListShiftSwapRequests(ctx context.Context, filter domain.ShiftSwapFilter) ([]*domain.ShiftSwapRequest, error) {
	rows, err := s.getQueries(ctx).ListShiftSwapRequests(ctx, sqlc.ListShiftSwapRequestsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, Status: textFromPtr(filter.Status), RequesterUserID: uuidFromPtr(filter.RequesterUserID), TargetUserID: uuidFromPtr(filter.TargetUserID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list shift swap requests", err, tenantIDField(filter.TenantID))
	}
	return mapShiftSwapRequests(rows), nil
}

func (s *Store) GetShiftSwapRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ShiftSwapRequest, error) {
	row, err := s.getQueries(ctx).GetShiftSwapRequest(ctx, sqlc.GetShiftSwapRequestParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrShiftSwapRequestNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get shift swap request", err, tenantIDField(tenantID), stringField("shift_swap_id", id.String()))
	}
	return mapShiftSwapRequest(row), nil
}

func (s *Store) CreateShiftScheduleEvent(ctx context.Context, item *domain.ShiftScheduleEvent, actorID *uuid.UUID) (*domain.ShiftScheduleEvent, error) {
	row, err := s.getQueries(ctx).CreateShiftScheduleEvent(ctx, sqlc.CreateShiftScheduleEventParams{TenantID: item.TenantID, SourceType: item.SourceType, SourceID: item.SourceID, Action: item.Action, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), ActorUserID: uuidFromPtr(item.ActorUserID), Remarks: textFromPtr(item.Remarks), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create shift schedule event", err, tenantIDField(item.TenantID), stringField("source_id", item.SourceID.String()))
	}
	return mapShiftScheduleEvent(row), nil
}

func (s *Store) ListShiftScheduleEvents(ctx context.Context, filter domain.ShiftScheduleEventFilter) ([]*domain.ShiftScheduleEvent, error) {
	rows, err := s.getQueries(ctx).ListShiftScheduleEvents(ctx, sqlc.ListShiftScheduleEventsParams{TenantID: filter.TenantID, Limit: limitOrDefault(filter.Limit), Offset: filter.Offset, SourceType: textFromPtr(filter.SourceType), SourceID: uuidFromPtr(filter.SourceID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list shift schedule events", err, tenantIDField(filter.TenantID))
	}
	return mapShiftScheduleEvents(rows), nil
}

func (s *Store) GetShiftScheduleSummary(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.ShiftScheduleSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetShiftScheduleSummary(ctx, sqlc.GetShiftScheduleSummaryParams{TenantID: tenantID, ScheduleDate: dateFromString(startDate), ScheduleDate_2: dateFromString(endDate)})
	if err != nil {
		return nil, s.logDBError(ctx, "get shift schedule summary", err, tenantIDField(tenantID))
	}
	return mapShiftScheduleSummaryRows(rows), nil
}

func (s *Store) ListShiftStaffingGaps(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.ShiftStaffingGapRow, error) {
	rows, err := s.getQueries(ctx).ListShiftStaffingGaps(ctx, sqlc.ListShiftStaffingGapsParams{TenantID: tenantID, ScheduleDate: dateFromString(startDate), ScheduleDate_2: dateFromString(endDate)})
	if err != nil {
		return nil, s.logDBError(ctx, "list shift staffing gaps", err, tenantIDField(tenantID))
	}
	return mapShiftStaffingGapRows(rows), nil
}

func staffingRequirementCreateParams(item *domain.StaffingRequirement, actorID *uuid.UUID) sqlc.CreateStaffingRequirementParams {
	return sqlc.CreateStaffingRequirementParams{TenantID: item.TenantID, Name: item.Name, RequirementDate: dateFromPtr(item.RequirementDate), StartDate: dateFromPtr(item.StartDate), EndDate: dateFromPtr(item.EndDate), DayOfWeek: int4FromPtr(item.DayOfWeek), BranchID: uuidFromPtr(item.BranchID), DepartmentID: uuidFromPtr(item.DepartmentID), AttendanceLocationID: uuidFromPtr(item.AttendanceLocationID), RoleLabel: textFromPtr(item.RoleLabel), TeamLabel: textFromPtr(item.TeamLabel), ShiftTemplateID: uuidFromPtr(item.ShiftTemplateID), RequiredCount: item.RequiredCount, MinCount: item.MinCount, MaxCount: int4FromPtr(item.MaxCount), Priority: item.Priority, Status: item.Status, PayrollBlocking: item.PayrollBlocking, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func staffingRequirementUpdateParams(item *domain.StaffingRequirement, actorID *uuid.UUID) sqlc.UpdateStaffingRequirementParams {
	return sqlc.UpdateStaffingRequirementParams{TenantID: item.TenantID, ID: item.ID, Name: item.Name, RequirementDate: dateFromPtr(item.RequirementDate), StartDate: dateFromPtr(item.StartDate), EndDate: dateFromPtr(item.EndDate), DayOfWeek: int4FromPtr(item.DayOfWeek), BranchID: uuidFromPtr(item.BranchID), DepartmentID: uuidFromPtr(item.DepartmentID), AttendanceLocationID: uuidFromPtr(item.AttendanceLocationID), RoleLabel: textFromPtr(item.RoleLabel), TeamLabel: textFromPtr(item.TeamLabel), ShiftTemplateID: uuidFromPtr(item.ShiftTemplateID), RequiredCount: item.RequiredCount, MinCount: item.MinCount, MaxCount: int4FromPtr(item.MaxCount), Priority: item.Priority, Status: item.Status, PayrollBlocking: item.PayrollBlocking, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}

func shiftAssignmentCreateParams(item *domain.ShiftScheduleAssignment, actorID *uuid.UUID) sqlc.CreateShiftScheduleAssignmentParams {
	return sqlc.CreateShiftScheduleAssignmentParams{TenantID: item.TenantID, ScheduleDate: dateFromTime(item.ScheduleDate), WorkerProfileID: uuidFromPtr(item.WorkerProfileID), EmployeeUserID: uuidFromPtr(item.EmployeeUserID), ShiftTemplateID: uuidFromPtr(item.ShiftTemplateID), AttendancePolicyID: uuidFromPtr(item.AttendancePolicyID), AttendanceLocationID: uuidFromPtr(item.AttendanceLocationID), AttendanceRosterID: uuidFromPtr(item.AttendanceRosterID), BranchID: uuidFromPtr(item.BranchID), DepartmentID: uuidFromPtr(item.DepartmentID), StartTime: timeFromClockString(item.StartTime), EndTime: timeFromClockString(item.EndTime), BreakMinutes: item.BreakMinutes, WorkMode: item.WorkMode, LocationType: item.LocationType, RoleLabel: textFromPtr(item.RoleLabel), TeamLabel: textFromPtr(item.TeamLabel), Status: item.Status, Source: item.Source, OvertimePlannedMinutes: item.OvertimePlannedMinutes, HasConflict: item.HasConflict, ConflictReason: textFromPtr(item.ConflictReason), PayrollBlocking: item.PayrollBlocking, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func shiftAssignmentUpdateParams(item *domain.ShiftScheduleAssignment, actorID *uuid.UUID) sqlc.UpdateShiftScheduleAssignmentParams {
	return sqlc.UpdateShiftScheduleAssignmentParams{TenantID: item.TenantID, ID: item.ID, ScheduleDate: dateFromTime(item.ScheduleDate), WorkerProfileID: uuidFromPtr(item.WorkerProfileID), EmployeeUserID: uuidFromPtr(item.EmployeeUserID), ShiftTemplateID: uuidFromPtr(item.ShiftTemplateID), AttendancePolicyID: uuidFromPtr(item.AttendancePolicyID), AttendanceLocationID: uuidFromPtr(item.AttendanceLocationID), AttendanceRosterID: uuidFromPtr(item.AttendanceRosterID), BranchID: uuidFromPtr(item.BranchID), DepartmentID: uuidFromPtr(item.DepartmentID), StartTime: timeFromClockString(item.StartTime), EndTime: timeFromClockString(item.EndTime), BreakMinutes: item.BreakMinutes, WorkMode: item.WorkMode, LocationType: item.LocationType, RoleLabel: textFromPtr(item.RoleLabel), TeamLabel: textFromPtr(item.TeamLabel), Status: item.Status, Source: item.Source, OvertimePlannedMinutes: item.OvertimePlannedMinutes, HasConflict: item.HasConflict, ConflictReason: textFromPtr(item.ConflictReason), PayrollBlocking: item.PayrollBlocking, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}
