package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateCandidateOnboarding(ctx context.Context, item *domain.CandidateOnboarding, actorID *uuid.UUID) (*domain.CandidateOnboarding, error) {
	row, err := s.getQueries(ctx).CreateCandidateOnboarding(ctx, sqlc.CreateCandidateOnboardingParams{TenantID: item.TenantID, CandidateID: item.CandidateID, WorkflowID: item.WorkflowID, OnboardingStatus: item.OnboardingStatus, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create candidate onboarding", err, tenantIDField(item.TenantID), stringField("candidate_id", item.CandidateID.String()), stringField("workflow_id", item.WorkflowID.String()))
	}
	return mapCandidateOnboardingBase(row), nil
}

func (s *Store) ListCandidateOnboardings(ctx context.Context, filter domain.CandidateOnboardingFilter) (*domain.CandidateOnboardingPage, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	q := s.getQueries(ctx)
	params := sqlc.ListCandidateOnboardingsParams{TenantID: filter.TenantID, Limit: limit, Offset: filter.Offset, Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)}
	rows, err := q.ListCandidateOnboardings(ctx, params)
	if err != nil {
		return nil, s.logDBError(ctx, "list candidate onboardings", err, tenantIDField(filter.TenantID))
	}
	total, err := q.CountCandidateOnboardings(ctx, sqlc.CountCandidateOnboardingsParams{TenantID: filter.TenantID, Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "count candidate onboardings", err, tenantIDField(filter.TenantID))
	}
	items := make([]*domain.CandidateOnboarding, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCandidateOnboardingListRow(row))
	}
	return &domain.CandidateOnboardingPage{Items: items, Total: total, Limit: limit, Offset: filter.Offset}, nil
}

func (s *Store) GetCandidateOnboarding(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CandidateOnboarding, error) {
	row, err := s.getQueries(ctx).GetCandidateOnboarding(ctx, sqlc.GetCandidateOnboardingParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCandidateOnboardingNotFound
		}
		return nil, s.logDBError(ctx, "get candidate onboarding", err, tenantIDField(tenantID), stringField("candidate_onboarding_id", id.String()))
	}
	return mapCandidateOnboardingRow(row), nil
}

func (s *Store) GetCandidateOnboardingByCandidate(ctx context.Context, tenantID uuid.UUID, candidateID uuid.UUID) (*domain.CandidateOnboarding, error) {
	row, err := s.getQueries(ctx).GetCandidateOnboardingByCandidate(ctx, sqlc.GetCandidateOnboardingByCandidateParams{TenantID: tenantID, CandidateID: candidateID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCandidateOnboardingNotFound
		}
		return nil, s.logDBError(ctx, "get candidate onboarding by candidate", err, tenantIDField(tenantID), stringField("candidate_id", candidateID.String()))
	}
	return mapCandidateOnboardingByCandidateRow(row), nil
}

func (s *Store) GetDefaultOnboardingWorkflow(ctx context.Context, tenantID uuid.UUID) (*domain.OnboardingWorkflow, error) {
	row, err := s.getQueries(ctx).GetDefaultOnboardingWorkflow(ctx, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOnboardingWorkflowNotFound
		}
		return nil, s.logDBError(ctx, "get default onboarding workflow", err, tenantIDField(tenantID))
	}
	return mapOnboardingWorkflow(row), nil
}

func (s *Store) ResolveOnboardingWorkflowForCandidate(ctx context.Context, tenantID uuid.UUID, candidateID uuid.UUID) (*domain.OnboardingWorkflow, error) {
	row, err := s.getQueries(ctx).ResolveOnboardingWorkflowForCandidate(ctx, sqlc.ResolveOnboardingWorkflowForCandidateParams{TenantID: tenantID, CandidateID: uuidFromPtr(&candidateID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOnboardingWorkflowNotFound
		}
		return nil, s.logDBError(ctx, "resolve onboarding workflow for candidate", err, tenantIDField(tenantID), stringField("candidate_id", candidateID.String()))
	}
	return mapOnboardingWorkflow(row), nil
}

func (s *Store) CreateCandidateOnboardingTasksFromWorkflow(ctx context.Context, tenantID uuid.UUID, candidateOnboardingID uuid.UUID, workflowID uuid.UUID, actorID *uuid.UUID) ([]*domain.CandidateOnboardingTask, error) {
	rows, err := s.getQueries(ctx).CreateCandidateOnboardingTasksFromWorkflow(ctx, sqlc.CreateCandidateOnboardingTasksFromWorkflowParams{TenantID: tenantID, CandidateOnboardingID: candidateOnboardingID, WorkflowID: workflowID, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create candidate onboarding task snapshots", err, tenantIDField(tenantID), stringField("candidate_onboarding_id", candidateOnboardingID.String()), stringField("workflow_id", workflowID.String()))
	}
	items := make([]*domain.CandidateOnboardingTask, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCandidateOnboardingTaskBase(row))
	}
	return items, nil
}

func (s *Store) RecalculateCandidateOnboardingProgress(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.CandidateOnboarding, error) {
	row, err := s.getQueries(ctx).RecalculateCandidateOnboardingProgress(ctx, sqlc.RecalculateCandidateOnboardingProgressParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCandidateOnboardingNotFound
		}
		return nil, s.logDBError(ctx, "recalculate candidate onboarding progress", err, tenantIDField(tenantID), stringField("candidate_onboarding_id", id.String()))
	}
	return mapCandidateOnboardingBase(row), nil
}

func (s *Store) ListCandidateOnboardingTasks(ctx context.Context, tenantID uuid.UUID, candidateOnboardingID uuid.UUID) ([]*domain.CandidateOnboardingTask, error) {
	rows, err := s.getQueries(ctx).ListCandidateOnboardingTasks(ctx, sqlc.ListCandidateOnboardingTasksParams{TenantID: tenantID, CandidateOnboardingID: candidateOnboardingID})
	if err != nil {
		return nil, s.logDBError(ctx, "list candidate onboarding tasks", err, tenantIDField(tenantID), stringField("candidate_onboarding_id", candidateOnboardingID.String()))
	}
	return mapCandidateOnboardingTasks(rows), nil
}

func (s *Store) GetCandidateOnboardingTask(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CandidateOnboardingTask, error) {
	row, err := s.getQueries(ctx).GetCandidateOnboardingTask(ctx, sqlc.GetCandidateOnboardingTaskParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCandidateOnboardingTaskNotFound
		}
		return nil, s.logDBError(ctx, "get candidate onboarding task", err, tenantIDField(tenantID), stringField("candidate_onboarding_task_id", id.String()))
	}
	return mapCandidateOnboardingTaskGetRow(row), nil
}

func (s *Store) UpdateCandidateOnboardingTaskStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.CandidateOnboardingTask, error) {
	row, err := s.getQueries(ctx).UpdateCandidateOnboardingTaskStatus(ctx, sqlc.UpdateCandidateOnboardingTaskStatusParams{TenantID: tenantID, ID: id, Status: status, Remarks: textFromPtr(remarks), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCandidateOnboardingTaskNotFound
		}
		return nil, s.logDBError(ctx, "update candidate onboarding task status", err, tenantIDField(tenantID), stringField("candidate_onboarding_task_id", id.String()))
	}
	return mapCandidateOnboardingTaskBase(row), nil
}

func (s *Store) CreateCandidateOnboardingEvent(ctx context.Context, event *domain.CandidateOnboardingEvent, actorID *uuid.UUID) (*domain.CandidateOnboardingEvent, error) {
	metadata := []byte(event.Metadata)
	if len(metadata) == 0 {
		metadata = []byte(`{}`)
	}
	row, err := s.getQueries(ctx).CreateCandidateOnboardingEvent(ctx, sqlc.CreateCandidateOnboardingEventParams{TenantID: event.TenantID, CandidateOnboardingID: event.CandidateOnboardingID, CandidateOnboardingTaskID: uuidFromPtr(event.CandidateOnboardingTaskID), Action: event.Action, FromStatus: textFromPtr(event.FromStatus), ToStatus: textFromPtr(event.ToStatus), Remarks: textFromPtr(event.Remarks), Column8: metadata, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create candidate onboarding event", err, tenantIDField(event.TenantID), stringField("candidate_onboarding_id", event.CandidateOnboardingID.String()))
	}
	return mapCandidateOnboardingEvent(row), nil
}

func (s *Store) ListCandidateOnboardingEvents(ctx context.Context, tenantID uuid.UUID, candidateOnboardingID uuid.UUID) ([]*domain.CandidateOnboardingEvent, error) {
	rows, err := s.getQueries(ctx).ListCandidateOnboardingEvents(ctx, sqlc.ListCandidateOnboardingEventsParams{TenantID: tenantID, CandidateOnboardingID: candidateOnboardingID})
	if err != nil {
		return nil, s.logDBError(ctx, "list candidate onboarding events", err, tenantIDField(tenantID), stringField("candidate_onboarding_id", candidateOnboardingID.String()))
	}
	return mapCandidateOnboardingEvents(rows), nil
}

func (s *Store) DeleteCandidateOnboarding(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCandidateOnboarding(ctx, sqlc.SoftDeleteCandidateOnboardingParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete candidate onboarding", err, tenantIDField(tenantID), stringField("candidate_onboarding_id", id.String()))
	}
	return nil
}
