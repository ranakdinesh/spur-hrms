package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateEmployeeExitRequest(ctx context.Context, item *domain.EmployeeExitRequest, actorID *uuid.UUID) (*domain.EmployeeExitRequest, error) {
	row, err := s.getQueries(ctx).CreateEmployeeExitRequest(ctx, sqlc.CreateEmployeeExitRequestParams{
		TenantID: item.TenantID, EmployeeID: item.EmployeeID, EmployeeUserID: item.EmployeeUserID, InitiatedBy: uuidFromPtr(item.InitiatedBy),
		Status: item.Status, ExitType: item.ExitType, Reason: textFromPtr(item.Reason), ResignationDate: dateFromPtr(item.ResignationDate),
		NoticeStartDate: dateFromPtr(item.NoticeStartDate), LastWorkingDate: dateFromTime(item.LastWorkingDate), RequestedRelievingDate: dateFromPtr(item.RequestedRelievingDate),
		FinalSettlementStatus: item.FinalSettlementStatus, AccessRevocationStatus: item.AccessRevocationStatus, AssetClearanceStatus: item.AssetClearanceStatus,
		HandoverStatus: item.HandoverStatus, ExitInterviewStatus: item.ExitInterviewStatus, Notes: textFromPtr(item.Notes), CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee exit request", err, tenantIDField(item.TenantID), stringField("employee_id", item.EmployeeID.String()))
	}
	return mapEmployeeExitRequest(row), nil
}

func (s *Store) ListEmployeeExitRequests(ctx context.Context, filter domain.EmployeeExitFilter) (*domain.EmployeeExitPage, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	params := sqlc.ListEmployeeExitRequestsParams{TenantID: filter.TenantID, Status: textFromPtr(filter.Status), EmployeeUserID: uuidFromPtr(filter.EmployeeUserID), Search: textFromPtr(filter.Search), Limit: limit, Offset: filter.Offset}
	rows, err := s.getQueries(ctx).ListEmployeeExitRequests(ctx, params)
	if err != nil {
		return nil, s.logDBError(ctx, "list employee exit requests", err, tenantIDField(filter.TenantID))
	}
	total, err := s.getQueries(ctx).CountEmployeeExitRequests(ctx, sqlc.CountEmployeeExitRequestsParams{TenantID: filter.TenantID, Status: textFromPtr(filter.Status), EmployeeUserID: uuidFromPtr(filter.EmployeeUserID), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "count employee exit requests", err, tenantIDField(filter.TenantID))
	}
	items := make([]*domain.EmployeeExitRequest, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapEmployeeExitListRow(row))
	}
	return &domain.EmployeeExitPage{Items: items, Total: total, Limit: limit, Offset: filter.Offset}, nil
}

func (s *Store) GetEmployeeExitRequest(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeExitRequest, error) {
	row, err := s.getQueries(ctx).GetEmployeeExitRequest(ctx, sqlc.GetEmployeeExitRequestParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmployeeExitNotFound
		}
		return nil, s.logDBError(ctx, "get employee exit request", err, tenantIDField(tenantID), stringField("exit_id", id.String()))
	}
	return mapEmployeeExitDetailRow(row), nil
}

func (s *Store) GetActiveEmployeeExitRequestByUserID(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.EmployeeExitRequest, error) {
	row, err := s.getQueries(ctx).GetActiveEmployeeExitRequestByUserID(ctx, sqlc.GetActiveEmployeeExitRequestByUserIDParams{TenantID: tenantID, EmployeeUserID: userID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmployeeExitNotFound
		}
		return nil, s.logDBError(ctx, "get active employee exit request", err, tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapEmployeeExitRequest(row), nil
}

func (s *Store) UpdateEmployeeExitRequestStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, approvedRelievingDate *time.Time, remarks *string, actorID *uuid.UUID) (*domain.EmployeeExitRequest, error) {
	row, err := s.getQueries(ctx).UpdateEmployeeExitRequestStatus(ctx, sqlc.UpdateEmployeeExitRequestStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID), ApprovedRelievingDate: dateFromPtr(approvedRelievingDate), RejectionReason: textFromPtr(remarks)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmployeeExitNotFound
		}
		return nil, s.logDBError(ctx, "update employee exit request status", err, tenantIDField(tenantID), stringField("exit_id", id.String()))
	}
	return mapEmployeeExitRequest(row), nil
}

func (s *Store) CreateEmployeeExitTask(ctx context.Context, task *domain.EmployeeExitTask, actorID *uuid.UUID) (*domain.EmployeeExitTask, error) {
	row, err := s.getQueries(ctx).CreateEmployeeExitTask(ctx, sqlc.CreateEmployeeExitTaskParams{
		TenantID: task.TenantID, ExitRequestID: task.ExitRequestID, EmployeeUserID: task.EmployeeUserID, TaskKey: task.TaskKey, Title: task.Title,
		Description: textFromPtr(task.Description), OwnerRole: textFromPtr(task.OwnerRole), DueDate: dateFromPtr(task.DueDate), Status: task.Status,
		Remarks: textFromPtr(task.Remarks), SortOrder: task.SortOrder, CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee exit task", err, tenantIDField(task.TenantID), stringField("exit_id", task.ExitRequestID.String()), stringField("task_key", task.TaskKey))
	}
	return mapEmployeeExitTask(row), nil
}

func (s *Store) ListEmployeeExitTasks(ctx context.Context, tenantID uuid.UUID, exitRequestID uuid.UUID) ([]*domain.EmployeeExitTask, error) {
	rows, err := s.getQueries(ctx).ListEmployeeExitTasks(ctx, sqlc.ListEmployeeExitTasksParams{TenantID: tenantID, ExitRequestID: exitRequestID})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee exit tasks", err, tenantIDField(tenantID), stringField("exit_id", exitRequestID.String()))
	}
	return mapEmployeeExitTasks(rows), nil
}

func (s *Store) GetEmployeeExitTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeExitTask, error) {
	row, err := s.getQueries(ctx).GetEmployeeExitTask(ctx, sqlc.GetEmployeeExitTaskParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmployeeExitTaskNotFound
		}
		return nil, s.logDBError(ctx, "get employee exit task", err, tenantIDField(tenantID), stringField("task_id", id.String()))
	}
	return mapEmployeeExitTask(row), nil
}

func (s *Store) UpdateEmployeeExitTaskStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.EmployeeExitTask, error) {
	row, err := s.getQueries(ctx).UpdateEmployeeExitTaskStatus(ctx, sqlc.UpdateEmployeeExitTaskStatusParams{TenantID: tenantID, ID: id, Status: status, Remarks: textFromPtr(remarks), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEmployeeExitTaskNotFound
		}
		return nil, s.logDBError(ctx, "update employee exit task status", err, tenantIDField(tenantID), stringField("task_id", id.String()))
	}
	return mapEmployeeExitTask(row), nil
}

func (s *Store) CreateEmployeeExitEvent(ctx context.Context, event *domain.EmployeeExitEvent, actorID *uuid.UUID) (*domain.EmployeeExitEvent, error) {
	metadata := []byte(event.Metadata)
	if len(metadata) == 0 {
		metadata = []byte(`{}`)
	}
	row, err := s.getQueries(ctx).CreateEmployeeExitEvent(ctx, sqlc.CreateEmployeeExitEventParams{TenantID: event.TenantID, ExitRequestID: event.ExitRequestID, ExitTaskID: uuidFromPtr(event.ExitTaskID), Action: event.Action, FromStatus: textFromPtr(event.FromStatus), ToStatus: textFromPtr(event.ToStatus), Remarks: textFromPtr(event.Remarks), Metadata: metadata, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create employee exit event", err, tenantIDField(event.TenantID), stringField("exit_id", event.ExitRequestID.String()))
	}
	return mapEmployeeExitEvent(row), nil
}

func (s *Store) ListEmployeeExitEvents(ctx context.Context, tenantID uuid.UUID, exitRequestID uuid.UUID) ([]*domain.EmployeeExitEvent, error) {
	rows, err := s.getQueries(ctx).ListEmployeeExitEvents(ctx, sqlc.ListEmployeeExitEventsParams{TenantID: tenantID, ExitRequestID: exitRequestID})
	if err != nil {
		return nil, s.logDBError(ctx, "list employee exit events", err, tenantIDField(tenantID), stringField("exit_id", exitRequestID.String()))
	}
	return mapEmployeeExitEvents(rows), nil
}
