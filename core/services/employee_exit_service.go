package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateEmployeeExit(ctx context.Context, cmd ports.CreateEmployeeExitCommand) (*domain.EmployeeExitRequest, error) {
	if cmd.TenantID == uuid.Nil || cmd.EmployeeID == uuid.Nil {
		err := domain.ErrInvalidEmployeeExitID
		s.logError("validate employee exit create", err)
		return nil, err
	}
	exitType, err := domain.ValidateEmployeeExitType(cmd.ExitType)
	if err != nil {
		s.logError("validate employee exit type", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	lastWorkingDate, err := parseRequiredDate(cmd.LastWorkingDate)
	if err != nil {
		s.logError("validate employee exit last working date", err, serviceTenantIDField(cmd.TenantID))
		return nil, domain.ErrInvalidEmployeeExitDate
	}
	resignationDate, err := parseOptionalDate(cmd.ResignationDate)
	if err != nil {
		return nil, domain.ErrInvalidEmployeeExitDate
	}
	noticeStartDate, err := parseOptionalDate(cmd.NoticeStartDate)
	if err != nil {
		return nil, domain.ErrInvalidEmployeeExitDate
	}
	requestedRelievingDate, err := parseOptionalDate(cmd.RequestedRelievingDate)
	if err != nil {
		return nil, domain.ErrInvalidEmployeeExitDate
	}
	if noticeStartDate != nil && lastWorkingDate.Before(*noticeStartDate) {
		return nil, domain.ErrInvalidEmployeeExitDate
	}
	profile, err := s.employees.GetEmployeeProfile(ctx, cmd.TenantID, cmd.EmployeeID)
	if err != nil {
		s.logError("get employee before exit create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	if existing, err := s.employeeExits.GetActiveEmployeeExitRequestByUserID(ctx, cmd.TenantID, profile.Employee.UserID); err == nil {
		return s.enrichEmployeeExit(ctx, existing)
	} else if !errors.Is(err, domain.ErrEmployeeExitNotFound) {
		s.logError("lookup active employee exit", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", profile.Employee.UserID.String()))
		return nil, err
	}
	item := &domain.EmployeeExitRequest{
		TenantID:               cmd.TenantID,
		EmployeeID:             cmd.EmployeeID,
		EmployeeUserID:         profile.Employee.UserID,
		InitiatedBy:            cmd.ActorID,
		Status:                 domain.EmployeeExitStatusSubmitted,
		ExitType:               exitType,
		Reason:                 cleanCommandString(cmd.Reason),
		ResignationDate:        resignationDate,
		NoticeStartDate:        noticeStartDate,
		LastWorkingDate:        *lastWorkingDate,
		RequestedRelievingDate: requestedRelievingDate,
		FinalSettlementStatus:  "pending",
		AccessRevocationStatus: "pending",
		AssetClearanceStatus:   "pending",
		HandoverStatus:         "pending",
		ExitInterviewStatus:    "pending",
		Notes:                  cleanCommandString(cmd.Notes),
	}
	var created *domain.EmployeeExitRequest
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var createErr error
		created, createErr = s.employeeExits.CreateEmployeeExitRequest(txCtx, item, cmd.ActorID)
		if createErr != nil {
			return createErr
		}
		for _, template := range domain.DefaultEmployeeExitTaskTemplates(*lastWorkingDate) {
			task := &domain.EmployeeExitTask{
				TenantID:       cmd.TenantID,
				ExitRequestID:  created.ID,
				EmployeeUserID: profile.Employee.UserID,
				TaskKey:        template.Key,
				Title:          template.Title,
				Description:    stringPtr(template.Description),
				OwnerRole:      stringPtr(template.OwnerRole),
				DueDate:        template.DueDate,
				Status:         domain.EmployeeExitTaskPending,
				SortOrder:      template.SortOrder,
			}
			if _, err := s.employeeExits.CreateEmployeeExitTask(txCtx, task, cmd.ActorID); err != nil {
				return err
			}
		}
		_, eventErr := s.createEmployeeExitEvent(txCtx, ports.EmployeeExitEventCommand{TenantID: cmd.TenantID, ExitRequestID: created.ID, Action: "submitted", ToStatus: &created.Status, ActorID: cmd.ActorID})
		return eventErr
	}); err != nil {
		s.logError("create employee exit transaction", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employee_id", cmd.EmployeeID.String()))
		return nil, err
	}
	return s.GetEmployeeExit(ctx, cmd.TenantID, created.ID)
}

func (s *TenantService) ListEmployeeExits(ctx context.Context, filter domain.EmployeeExitFilter) (*domain.EmployeeExitPage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate list employee exits tenant", err)
		return nil, err
	}
	page, err := s.employeeExits.ListEmployeeExitRequests(ctx, filter)
	if err != nil {
		s.logError("list employee exits", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return page, nil
}

func (s *TenantService) GetEmployeeExit(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeExitRequest, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidEmployeeExitID
		s.logError("validate get employee exit", err)
		return nil, err
	}
	item, err := s.employeeExits.GetEmployeeExitRequest(ctx, tenantID, id)
	if err != nil {
		s.logError("get employee exit", err, serviceTenantIDField(tenantID), serviceStringField("exit_id", id.String()))
		return nil, err
	}
	return s.enrichEmployeeExit(ctx, item)
}

func (s *TenantService) ApproveEmployeeExit(ctx context.Context, cmd ports.EmployeeExitActionCommand) (*domain.EmployeeExitRequest, error) {
	approvedRelievingDate, err := parseOptionalDate(cmd.ApprovedRelievingDate)
	if err != nil {
		return nil, domain.ErrInvalidEmployeeExitDate
	}
	return s.transitionEmployeeExit(ctx, cmd, domain.EmployeeExitStatusApproved, approvedRelievingDate, "approved")
}

func (s *TenantService) RejectEmployeeExit(ctx context.Context, cmd ports.EmployeeExitActionCommand) (*domain.EmployeeExitRequest, error) {
	return s.transitionEmployeeExit(ctx, cmd, domain.EmployeeExitStatusRejected, nil, "rejected")
}

func (s *TenantService) CancelEmployeeExit(ctx context.Context, cmd ports.EmployeeExitActionCommand) (*domain.EmployeeExitRequest, error) {
	return s.transitionEmployeeExit(ctx, cmd, domain.EmployeeExitStatusCanceled, nil, "canceled")
}

func (s *TenantService) CompleteEmployeeExit(ctx context.Context, cmd ports.EmployeeExitActionCommand) (*domain.EmployeeExitRequest, error) {
	if cmd.TenantID == uuid.Nil || cmd.ExitID == uuid.Nil {
		return nil, domain.ErrInvalidEmployeeExitID
	}
	item, err := s.employeeExits.GetEmployeeExitRequest(ctx, cmd.TenantID, cmd.ExitID)
	if err != nil {
		return nil, err
	}
	if item.Status != domain.EmployeeExitStatusApproved {
		return nil, domain.ErrInvalidEmployeeExitStatus
	}
	tasks, err := s.employeeExits.ListEmployeeExitTasks(ctx, cmd.TenantID, cmd.ExitID)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		if task.Status != domain.EmployeeExitTaskCompleted && task.Status != domain.EmployeeExitTaskWaived {
			return nil, domain.ErrEmployeeExitCompletionBlock
		}
	}
	if s.employeeIdentity == nil {
		err := domain.ErrEmployeeIdentityPortMissing
		s.logError("complete employee exit identity port missing", err, serviceTenantIDField(cmd.TenantID), serviceStringField("exit_id", cmd.ExitID.String()))
		return nil, err
	}
	var updated *domain.EmployeeExitRequest
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var updateErr error
		updated, updateErr = s.employeeExits.UpdateEmployeeExitRequestStatus(txCtx, cmd.TenantID, cmd.ExitID, domain.EmployeeExitStatusCompleted, nil, cmd.Remarks, cmd.ActorID)
		if updateErr != nil {
			return updateErr
		}
		if err := s.employees.DeactivateEmployee(txCtx, cmd.TenantID, item.EmployeeID, cmd.ActorID); err != nil {
			return err
		}
		if err := s.employeeIdentity.DeactivateEmployeeIdentity(txCtx, ports.DeactivateEmployeeIdentityCommand{TenantID: cmd.TenantID, UserID: item.EmployeeUserID, ActorID: cmd.ActorID}); err != nil {
			return err
		}
		_, eventErr := s.createEmployeeExitEvent(txCtx, ports.EmployeeExitEventCommand{TenantID: cmd.TenantID, ExitRequestID: cmd.ExitID, Action: "completed", FromStatus: &item.Status, ToStatus: &updated.Status, Remarks: cmd.Remarks, ActorID: cmd.ActorID})
		return eventErr
	}); err != nil {
		s.logError("complete employee exit transaction", err, serviceTenantIDField(cmd.TenantID), serviceStringField("exit_id", cmd.ExitID.String()))
		return nil, err
	}
	return s.GetEmployeeExit(ctx, cmd.TenantID, updated.ID)
}

func (s *TenantService) UpdateEmployeeExitTaskStatus(ctx context.Context, cmd ports.EmployeeExitTaskStatusCommand) (*domain.EmployeeExitRequest, error) {
	if cmd.TenantID == uuid.Nil || cmd.TaskID == uuid.Nil {
		return nil, domain.ErrInvalidEmployeeExitID
	}
	status, err := domain.ValidateEmployeeExitTaskStatus(cmd.Status)
	if err != nil {
		return nil, err
	}
	before, err := s.employeeExits.GetEmployeeExitTask(ctx, cmd.TenantID, cmd.TaskID)
	if err != nil {
		return nil, err
	}
	updated, err := s.employeeExits.UpdateEmployeeExitTaskStatus(ctx, cmd.TenantID, cmd.TaskID, status, cmd.Remarks, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	if _, err := s.createEmployeeExitEvent(ctx, ports.EmployeeExitEventCommand{TenantID: cmd.TenantID, ExitRequestID: updated.ExitRequestID, ExitTaskID: &updated.ID, Action: "task_status_changed", FromStatus: &before.Status, ToStatus: &updated.Status, Remarks: cmd.Remarks, ActorID: cmd.ActorID}); err != nil {
		s.logError("record employee exit task event", err, serviceTenantIDField(cmd.TenantID), serviceStringField("task_id", cmd.TaskID.String()))
	}
	return s.GetEmployeeExit(ctx, cmd.TenantID, updated.ExitRequestID)
}

func (s *TenantService) transitionEmployeeExit(ctx context.Context, cmd ports.EmployeeExitActionCommand, status string, approvedRelievingDate *time.Time, action string) (*domain.EmployeeExitRequest, error) {
	if cmd.TenantID == uuid.Nil || cmd.ExitID == uuid.Nil {
		return nil, domain.ErrInvalidEmployeeExitID
	}
	before, err := s.employeeExits.GetEmployeeExitRequest(ctx, cmd.TenantID, cmd.ExitID)
	if err != nil {
		return nil, err
	}
	updated, err := s.employeeExits.UpdateEmployeeExitRequestStatus(ctx, cmd.TenantID, cmd.ExitID, status, approvedRelievingDate, cmd.Remarks, cmd.ActorID)
	if err != nil {
		s.logError("transition employee exit", err, serviceTenantIDField(cmd.TenantID), serviceStringField("exit_id", cmd.ExitID.String()), serviceStringField("status", status))
		return nil, err
	}
	if _, err := s.createEmployeeExitEvent(ctx, ports.EmployeeExitEventCommand{TenantID: cmd.TenantID, ExitRequestID: cmd.ExitID, Action: action, FromStatus: &before.Status, ToStatus: &updated.Status, Remarks: cmd.Remarks, ActorID: cmd.ActorID}); err != nil {
		s.logError("record employee exit event", err, serviceTenantIDField(cmd.TenantID), serviceStringField("exit_id", cmd.ExitID.String()))
	}
	return s.GetEmployeeExit(ctx, cmd.TenantID, updated.ID)
}

func (s *TenantService) enrichEmployeeExit(ctx context.Context, item *domain.EmployeeExitRequest) (*domain.EmployeeExitRequest, error) {
	if item == nil {
		return nil, domain.ErrEmployeeExitNotFound
	}
	tasks, err := s.employeeExits.ListEmployeeExitTasks(ctx, item.TenantID, item.ID)
	if err != nil {
		return nil, err
	}
	events, err := s.employeeExits.ListEmployeeExitEvents(ctx, item.TenantID, item.ID)
	if err != nil {
		return nil, err
	}
	item.Tasks = tasks
	item.Events = events
	item.TotalTasks = int32(len(tasks))
	item.CompletedTasks = 0
	item.BlockedTasks = 0
	for _, task := range tasks {
		if task.Status == domain.EmployeeExitTaskCompleted || task.Status == domain.EmployeeExitTaskWaived {
			item.CompletedTasks++
		}
		if task.Status == domain.EmployeeExitTaskBlocked {
			item.BlockedTasks++
		}
	}
	return item, nil
}

func (s *TenantService) createEmployeeExitEvent(ctx context.Context, cmd ports.EmployeeExitEventCommand) (*domain.EmployeeExitEvent, error) {
	event := &domain.EmployeeExitEvent{TenantID: cmd.TenantID, ExitRequestID: cmd.ExitRequestID, ExitTaskID: cmd.ExitTaskID, Action: cmd.Action, FromStatus: cmd.FromStatus, ToStatus: cmd.ToStatus, Remarks: cmd.Remarks, Metadata: cmd.Metadata}
	return s.employeeExits.CreateEmployeeExitEvent(ctx, event, cmd.ActorID)
}

func parseRequiredDate(value string) (*time.Time, error) {
	parsed, err := parseOptionalDate(value)
	if err != nil {
		return nil, err
	}
	if parsed == nil {
		return nil, domain.ErrInvalidEmployeeExitDate
	}
	return parsed, nil
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
