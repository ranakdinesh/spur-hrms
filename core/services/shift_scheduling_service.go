package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateShiftTemplate(ctx context.Context, cmd ports.ShiftTemplateCommand) (*domain.ShiftTemplate, error) {
	item, err := s.buildShiftTemplate(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.shiftScheduling.CreateShiftTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create shift template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	to := boolState(result.IsActive)
	_ = s.createShiftScheduleEvent(ctx, cmd.TenantID, "template", result.ID, "created", nil, &to, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateShiftTemplate(ctx context.Context, cmd ports.ShiftTemplateCommand) (*domain.ShiftTemplate, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidShiftTemplate
	}
	item, err := s.buildShiftTemplate(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	before, _ := s.shiftScheduling.GetShiftTemplate(ctx, cmd.TenantID, cmd.ID)
	result, err := s.shiftScheduling.UpdateShiftTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update shift template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("template_id", cmd.ID.String()))
		return nil, err
	}
	var from *string
	if before != nil {
		state := boolState(before.IsActive)
		from = &state
	}
	to := boolState(result.IsActive)
	_ = s.createShiftScheduleEvent(ctx, cmd.TenantID, "template", result.ID, "updated", from, &to, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) ListShiftTemplates(ctx context.Context, tenantID uuid.UUID, activeOnly *bool, search *string, limit int32, offset int32) ([]*domain.ShiftTemplate, error) {
	items, err := s.shiftScheduling.ListShiftTemplates(ctx, tenantID, activeOnly, search, limit, offset)
	if err != nil {
		s.logError("list shift templates", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) DeleteShiftTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidShiftTemplate
	}
	if err := s.shiftScheduling.DeleteShiftTemplate(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete shift template", err, serviceTenantIDField(tenantID), serviceStringField("template_id", id.String()))
		return err
	}
	_ = s.createShiftScheduleEvent(ctx, tenantID, "template", id, "deleted", nil, nil, actorID, nil)
	return nil
}

func (s *TenantService) CreateStaffingRequirement(ctx context.Context, cmd ports.StaffingRequirementCommand) (*domain.StaffingRequirement, error) {
	item, err := s.buildStaffingRequirement(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.shiftScheduling.CreateStaffingRequirement(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create staffing requirement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("name", item.Name))
		return nil, err
	}
	_ = s.createShiftScheduleEvent(ctx, cmd.TenantID, "requirement", result.ID, "created", nil, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) UpdateStaffingRequirement(ctx context.Context, cmd ports.StaffingRequirementCommand) (*domain.StaffingRequirement, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidStaffingRequirement
	}
	item, err := s.buildStaffingRequirement(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	before, _ := s.shiftScheduling.GetStaffingRequirement(ctx, cmd.TenantID, cmd.ID)
	result, err := s.shiftScheduling.UpdateStaffingRequirement(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update staffing requirement", err, serviceTenantIDField(cmd.TenantID), serviceStringField("requirement_id", cmd.ID.String()))
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	_ = s.createShiftScheduleEvent(ctx, cmd.TenantID, "requirement", result.ID, "updated", from, &result.Status, cmd.ActorID, nil)
	return result, nil
}

func (s *TenantService) ListStaffingRequirements(ctx context.Context, filter domain.StaffingRequirementFilter) ([]*domain.StaffingRequirement, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	items, err := s.shiftScheduling.ListStaffingRequirements(ctx, filter)
	if err != nil {
		s.logError("list staffing requirements", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) DeleteStaffingRequirement(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidStaffingRequirement
	}
	if err := s.shiftScheduling.DeleteStaffingRequirement(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete staffing requirement", err, serviceTenantIDField(tenantID), serviceStringField("requirement_id", id.String()))
		return err
	}
	_ = s.createShiftScheduleEvent(ctx, tenantID, "requirement", id, "deleted", nil, nil, actorID, nil)
	return nil
}

func (s *TenantService) CreateShiftScheduleAssignment(ctx context.Context, cmd ports.ShiftScheduleAssignmentCommand) (*domain.ShiftScheduleAssignment, error) {
	item, err := s.buildShiftAssignment(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if err := s.applyShiftConflict(ctx, item, uuid.Nil); err != nil {
		return nil, err
	}
	result, err := s.shiftScheduling.CreateShiftScheduleAssignment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create shift assignment", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	_ = s.createShiftScheduleEvent(ctx, cmd.TenantID, "assignment", result.ID, "created", nil, &result.Status, cmd.ActorID, result.ConflictReason)
	return result, nil
}

func (s *TenantService) UpdateShiftScheduleAssignment(ctx context.Context, cmd ports.ShiftScheduleAssignmentCommand) (*domain.ShiftScheduleAssignment, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidShiftAssignment
	}
	item, err := s.buildShiftAssignment(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	if err := s.applyShiftConflict(ctx, item, cmd.ID); err != nil {
		return nil, err
	}
	before, _ := s.shiftScheduling.GetShiftScheduleAssignment(ctx, cmd.TenantID, cmd.ID)
	result, err := s.shiftScheduling.UpdateShiftScheduleAssignment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update shift assignment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("assignment_id", cmd.ID.String()))
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	_ = s.createShiftScheduleEvent(ctx, cmd.TenantID, "assignment", result.ID, "updated", from, &result.Status, cmd.ActorID, result.ConflictReason)
	return result, nil
}

func (s *TenantService) UpdateShiftScheduleAssignmentStatus(ctx context.Context, cmd ports.ShiftScheduleStatusCommand) (*domain.ShiftScheduleAssignment, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidShiftAssignment
	}
	status, err := domain.NormalizeShiftAssignmentStatus(cmd.Status)
	if err != nil {
		return nil, err
	}
	before, err := s.shiftScheduling.GetShiftScheduleAssignment(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		s.logError("get shift assignment for status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("assignment_id", cmd.ID.String()))
		return nil, err
	}
	check := *before
	check.Status = status
	if err := s.applyShiftConflict(ctx, &check, cmd.ID); err != nil {
		return nil, err
	}
	result, err := s.shiftScheduling.UpdateShiftScheduleAssignmentStatus(ctx, cmd.TenantID, cmd.ID, status, check.HasConflict, check.ConflictReason, check.PayrollBlocking, cmd.ActorID)
	if err != nil {
		s.logError("update shift assignment status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("assignment_id", cmd.ID.String()))
		return nil, err
	}
	_ = s.createShiftScheduleEvent(ctx, cmd.TenantID, "assignment", result.ID, "status_updated", &before.Status, &result.Status, cmd.ActorID, cmd.Remarks)
	return result, nil
}

func (s *TenantService) ListShiftScheduleAssignments(ctx context.Context, filter domain.ShiftScheduleFilter) ([]*domain.ShiftScheduleAssignment, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	start, end, err := parseDateRangeOrToday(filter.StartDate, filter.EndDate)
	if err != nil {
		return nil, err
	}
	filter.StartDate = start.Format("2006-01-02")
	filter.EndDate = end.Format("2006-01-02")
	items, err := s.shiftScheduling.ListShiftScheduleAssignments(ctx, filter)
	if err != nil {
		s.logError("list shift assignments", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) DeleteShiftScheduleAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidShiftAssignment
	}
	if err := s.shiftScheduling.DeleteShiftScheduleAssignment(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete shift assignment", err, serviceTenantIDField(tenantID), serviceStringField("assignment_id", id.String()))
		return err
	}
	_ = s.createShiftScheduleEvent(ctx, tenantID, "assignment", id, "deleted", nil, nil, actorID, nil)
	return nil
}

func (s *TenantService) CreateShiftSwapRequest(ctx context.Context, cmd ports.ShiftSwapRequestCommand) (*domain.ShiftSwapRequest, error) {
	requestedDate, err := parseOptionalDate(cmd.RequestedDate)
	if err != nil {
		return nil, err
	}
	assignment, err := s.shiftScheduling.GetShiftScheduleAssignment(ctx, cmd.TenantID, cmd.RequesterAssignmentID)
	if err != nil {
		s.logError("validate shift swap assignment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("assignment_id", cmd.RequesterAssignmentID.String()))
		return nil, err
	}
	item, err := domain.NewShiftSwapRequest(domain.ShiftSwapRequest{
		TenantID:                 cmd.TenantID,
		RequesterAssignmentID:    cmd.RequesterAssignmentID,
		RequesterWorkerProfileID: valueOrFallbackUUIDPtr(cmd.RequesterWorkerProfileID, assignment.WorkerProfileID),
		RequesterUserID:          valueOrFallbackUUIDPtr(cmd.RequesterUserID, assignment.EmployeeUserID),
		TargetWorkerProfileID:    cmd.TargetWorkerProfileID,
		TargetUserID:             cmd.TargetUserID,
		OfferedAssignmentID:      cmd.OfferedAssignmentID,
		RequestedDate:            requestedDate,
		RequestedShiftTemplateID: cmd.RequestedShiftTemplateID,
		Reason:                   cmd.Reason,
		Status:                   domain.ShiftSwapStatusPending,
		PayrollBlocking:          true,
		Metadata:                 cmd.Metadata,
	})
	if err != nil {
		return nil, err
	}
	result, err := s.shiftScheduling.CreateShiftSwapRequest(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create shift swap request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("assignment_id", cmd.RequesterAssignmentID.String()))
		return nil, err
	}
	_ = s.createShiftScheduleEvent(ctx, cmd.TenantID, "swap_request", result.ID, "created", nil, &result.Status, cmd.ActorID, result.Reason)
	return result, nil
}

func (s *TenantService) ReviewShiftSwapRequest(ctx context.Context, cmd ports.ShiftSwapReviewCommand) (*domain.ShiftSwapRequest, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || cmd.ReviewerID == uuid.Nil {
		return nil, domain.ErrInvalidShiftSwapRequest
	}
	status, err := domain.NormalizeShiftSwapStatus(cmd.Status)
	if err != nil || status == domain.ShiftSwapStatusPending {
		return nil, domain.ErrInvalidShiftSwapRequest
	}
	before, err := s.shiftScheduling.GetShiftSwapRequest(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		s.logError("get shift swap for review", err, serviceTenantIDField(cmd.TenantID), serviceStringField("swap_id", cmd.ID.String()))
		return nil, err
	}
	payrollBlocking := status == domain.ShiftSwapStatusPending
	result, err := s.shiftScheduling.UpdateShiftSwapRequestReview(ctx, cmd.TenantID, cmd.ID, status, cmd.ReviewerID, cmd.Remarks, payrollBlocking)
	if err != nil {
		s.logError("review shift swap request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("swap_id", cmd.ID.String()))
		return nil, err
	}
	_ = s.createShiftScheduleEvent(ctx, cmd.TenantID, "swap_request", result.ID, "reviewed", &before.Status, &result.Status, &cmd.ReviewerID, cmd.Remarks)
	return result, nil
}

func (s *TenantService) ListShiftSwapRequests(ctx context.Context, filter domain.ShiftSwapFilter) ([]*domain.ShiftSwapRequest, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	items, err := s.shiftScheduling.ListShiftSwapRequests(ctx, filter)
	if err != nil {
		s.logError("list shift swap requests", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListShiftScheduleEvents(ctx context.Context, filter domain.ShiftScheduleEventFilter) ([]*domain.ShiftScheduleEvent, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	items, err := s.shiftScheduling.ListShiftScheduleEvents(ctx, filter)
	if err != nil {
		s.logError("list shift schedule events", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetShiftScheduleSummary(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.ShiftScheduleSummaryRow, error) {
	start, end, err := parseDateRangeOrToday(startDate, endDate)
	if err != nil {
		return nil, err
	}
	items, err := s.shiftScheduling.GetShiftScheduleSummary(ctx, tenantID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		s.logError("get shift schedule summary", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListShiftStaffingGaps(ctx context.Context, tenantID uuid.UUID, startDate string, endDate string) ([]*domain.ShiftStaffingGapRow, error) {
	start, end, err := parseDateRangeOrToday(startDate, endDate)
	if err != nil {
		return nil, err
	}
	items, err := s.shiftScheduling.ListShiftStaffingGaps(ctx, tenantID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		s.logError("list shift staffing gaps", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) buildShiftTemplate(ctx context.Context, cmd ports.ShiftTemplateCommand) (*domain.ShiftTemplate, error) {
	if cmd.AttendancePolicyID != nil && *cmd.AttendancePolicyID != uuid.Nil {
		if _, err := s.attendancePolicies.GetAttendancePolicy(ctx, cmd.TenantID, *cmd.AttendancePolicyID); err != nil {
			return nil, err
		}
	}
	if cmd.AttendanceLocationID != nil && *cmd.AttendanceLocationID != uuid.Nil {
		if _, err := s.attendanceLocations.GetAttendanceLocation(ctx, cmd.TenantID, *cmd.AttendanceLocationID); err != nil {
			return nil, err
		}
	}
	return domain.NewShiftTemplate(domain.ShiftTemplate{
		TenantID:             cmd.TenantID,
		Code:                 cmd.Code,
		Name:                 cmd.Name,
		Description:          cmd.Description,
		StartTime:            cmd.StartTime,
		EndTime:              cmd.EndTime,
		BreakMinutes:         cmd.BreakMinutes,
		PaidMinutes:          cmd.PaidMinutes,
		WorkMode:             cmd.WorkMode,
		LocationType:         cmd.LocationType,
		AttendancePolicyID:   cmd.AttendancePolicyID,
		AttendanceLocationID: cmd.AttendanceLocationID,
		AllowOvertime:        cmd.AllowOvertime,
		PayrollCode:          cmd.PayrollCode,
		Metadata:             cmd.Metadata,
		IsActive:             cmd.IsActive,
	})
}

func (s *TenantService) buildStaffingRequirement(ctx context.Context, cmd ports.StaffingRequirementCommand) (*domain.StaffingRequirement, error) {
	requirementDate, err := parseOptionalDate(cmd.RequirementDate)
	if err != nil {
		return nil, err
	}
	startDate, err := parseOptionalDate(cmd.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := parseOptionalDate(cmd.EndDate)
	if err != nil {
		return nil, err
	}
	if cmd.ShiftTemplateID != nil && *cmd.ShiftTemplateID != uuid.Nil {
		if _, err := s.shiftScheduling.GetShiftTemplate(ctx, cmd.TenantID, *cmd.ShiftTemplateID); err != nil {
			return nil, err
		}
	}
	return domain.NewStaffingRequirement(domain.StaffingRequirement{
		TenantID:             cmd.TenantID,
		Name:                 cmd.Name,
		RequirementDate:      requirementDate,
		StartDate:            startDate,
		EndDate:              endDate,
		DayOfWeek:            cmd.DayOfWeek,
		BranchID:             cmd.BranchID,
		DepartmentID:         cmd.DepartmentID,
		AttendanceLocationID: cmd.AttendanceLocationID,
		RoleLabel:            cmd.RoleLabel,
		TeamLabel:            cmd.TeamLabel,
		ShiftTemplateID:      cmd.ShiftTemplateID,
		RequiredCount:        cmd.RequiredCount,
		MinCount:             cmd.MinCount,
		MaxCount:             cmd.MaxCount,
		Priority:             cmd.Priority,
		Status:               cmd.Status,
		PayrollBlocking:      cmd.PayrollBlocking,
		Notes:                cmd.Notes,
		Metadata:             cmd.Metadata,
	})
}

func (s *TenantService) buildShiftAssignment(ctx context.Context, cmd ports.ShiftScheduleAssignmentCommand) (*domain.ShiftScheduleAssignment, error) {
	date, err := parseAttendanceDate(cmd.ScheduleDate, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	if cmd.ShiftTemplateID != nil && *cmd.ShiftTemplateID != uuid.Nil {
		template, err := s.shiftScheduling.GetShiftTemplate(ctx, cmd.TenantID, *cmd.ShiftTemplateID)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(cmd.StartTime) == "" {
			cmd.StartTime = template.StartTime
		}
		if strings.TrimSpace(cmd.EndTime) == "" {
			cmd.EndTime = template.EndTime
		}
		if cmd.BreakMinutes == 0 {
			cmd.BreakMinutes = template.BreakMinutes
		}
		if strings.TrimSpace(cmd.WorkMode) == "" {
			cmd.WorkMode = template.WorkMode
		}
		if strings.TrimSpace(cmd.LocationType) == "" {
			cmd.LocationType = template.LocationType
		}
		if cmd.AttendancePolicyID == nil {
			cmd.AttendancePolicyID = template.AttendancePolicyID
		}
		if cmd.AttendanceLocationID == nil {
			cmd.AttendanceLocationID = template.AttendanceLocationID
		}
	}
	if cmd.EmployeeUserID != nil && *cmd.EmployeeUserID != uuid.Nil {
		if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, *cmd.EmployeeUserID); err != nil {
			return nil, err
		}
	}
	if cmd.AttendancePolicyID != nil && *cmd.AttendancePolicyID != uuid.Nil {
		if _, err := s.attendancePolicies.GetAttendancePolicy(ctx, cmd.TenantID, *cmd.AttendancePolicyID); err != nil {
			return nil, err
		}
	}
	if cmd.AttendanceLocationID != nil && *cmd.AttendanceLocationID != uuid.Nil {
		if _, err := s.attendanceLocations.GetAttendanceLocation(ctx, cmd.TenantID, *cmd.AttendanceLocationID); err != nil {
			return nil, err
		}
	}
	return domain.NewShiftScheduleAssignment(domain.ShiftScheduleAssignment{
		TenantID:               cmd.TenantID,
		ScheduleDate:           date,
		WorkerProfileID:        cmd.WorkerProfileID,
		EmployeeUserID:         cmd.EmployeeUserID,
		ShiftTemplateID:        cmd.ShiftTemplateID,
		AttendancePolicyID:     cmd.AttendancePolicyID,
		AttendanceLocationID:   cmd.AttendanceLocationID,
		AttendanceRosterID:     cmd.AttendanceRosterID,
		BranchID:               cmd.BranchID,
		DepartmentID:           cmd.DepartmentID,
		StartTime:              cmd.StartTime,
		EndTime:                cmd.EndTime,
		BreakMinutes:           cmd.BreakMinutes,
		WorkMode:               cmd.WorkMode,
		LocationType:           cmd.LocationType,
		RoleLabel:              cmd.RoleLabel,
		TeamLabel:              cmd.TeamLabel,
		Status:                 cmd.Status,
		Source:                 cmd.Source,
		OvertimePlannedMinutes: cmd.OvertimePlannedMinutes,
		Notes:                  cmd.Notes,
		Metadata:               cmd.Metadata,
	})
}

func (s *TenantService) applyShiftConflict(ctx context.Context, item *domain.ShiftScheduleAssignment, excludeID uuid.UUID) error {
	if item == nil {
		return nil
	}
	existing, err := s.shiftScheduling.ListShiftAssignmentsForWorkerDate(ctx, item.TenantID, item.ScheduleDate.Format("2006-01-02"), item.WorkerProfileID, item.EmployeeUserID)
	if err != nil {
		s.logError("list shift assignments for conflict", err, serviceTenantIDField(item.TenantID))
		return err
	}
	item.HasConflict = false
	item.ConflictReason = nil
	item.PayrollBlocking = false
	for _, other := range existing {
		if other == nil || other.ID == excludeID || other.Status == domain.ShiftScheduleStatusCancelled {
			continue
		}
		if shiftsOverlap(item.StartTime, item.EndTime, other.StartTime, other.EndTime) {
			reason := "Overlaps another assigned shift on the same date"
			item.HasConflict = true
			item.ConflictReason = &reason
			item.PayrollBlocking = true
			return nil
		}
	}
	return nil
}

func (s *TenantService) createShiftScheduleEvent(ctx context.Context, tenantID uuid.UUID, sourceType string, sourceID uuid.UUID, action string, fromStatus *string, toStatus *string, actorID *uuid.UUID, remarks *string) error {
	actorUserID := actorID
	item, err := domain.NewShiftScheduleEvent(domain.ShiftScheduleEvent{TenantID: tenantID, SourceType: sourceType, SourceID: sourceID, Action: action, FromStatus: fromStatus, ToStatus: toStatus, ActorUserID: actorUserID, Remarks: remarks})
	if err != nil {
		return err
	}
	_, err = s.shiftScheduling.CreateShiftScheduleEvent(ctx, item, actorID)
	return err
}

func shiftsOverlap(startA string, endA string, startB string, endB string) bool {
	aStart, aErr := clockMinutes(startA)
	aEnd, bErr := clockMinutes(endA)
	bStart, cErr := clockMinutes(startB)
	bEnd, dErr := clockMinutes(endB)
	if aErr != nil || bErr != nil || cErr != nil || dErr != nil {
		return false
	}
	return aStart < bEnd && bStart < aEnd
}

func clockMinutes(value string) (int, error) {
	parsed, err := time.Parse("15:04", strings.TrimSpace(value))
	if err != nil {
		return 0, err
	}
	return parsed.Hour()*60 + parsed.Minute(), nil
}

func valueOrFallbackUUIDPtr(value *uuid.UUID, fallback *uuid.UUID) *uuid.UUID {
	if value != nil && *value != uuid.Nil {
		return value
	}
	return fallback
}

func boolState(value bool) string {
	if value {
		return "active"
	}
	return "inactive"
}
