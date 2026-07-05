package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) ListTenantProfiles(ctx context.Context) ([]*domain.TenantProfile, error) {
	rows, err := s.getQueries(ctx).ListTenantProfiles(ctx)
	if err != nil {
		return nil, s.logDBError(ctx, "list tenant profiles", err)
	}
	items := make([]*domain.TenantProfile, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapTenantProfile(row))
	}
	return items, nil
}

func (s *Store) AcquireJobLock(ctx context.Context, tenantID uuid.UUID, jobKey string, ownerID string, ttlSeconds int32) (bool, error) {
	if ttlSeconds <= 0 {
		ttlSeconds = 900
	}
	_, err := s.getQueries(ctx).AcquireJobLock(ctx, sqlc.AcquireJobLockParams{
		TenantID: tenantID,
		JobKey:   jobKey,
		Column3:  strconv.Itoa(int(ttlSeconds)),
		OwnerID:  ownerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, s.logDBError(ctx, "acquire job lock", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	return true, nil
}

func (s *Store) ReleaseJobLock(ctx context.Context, tenantID uuid.UUID, jobKey string, ownerID string) error {
	if err := s.getQueries(ctx).ReleaseJobLock(ctx, sqlc.ReleaseJobLockParams{TenantID: tenantID, JobKey: jobKey, OwnerID: ownerID}); err != nil {
		return s.logDBError(ctx, "release job lock", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	return nil
}

func (s *Store) GetJobRunByDate(ctx context.Context, tenantID uuid.UUID, jobKey string, runDate time.Time) (*domain.ScheduledJobRun, error) {
	row, err := s.getQueries(ctx).GetJobRunByDate(ctx, sqlc.GetJobRunByDateParams{TenantID: tenantID, JobKey: jobKey, RunDate: dateFromTime(runDate)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, s.logDBError(ctx, "get job run by date", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	item, err := mapScheduledJobRun(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map job run by date", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	return item, nil
}

func (s *Store) StartJobRun(ctx context.Context, tenantID uuid.UUID, jobKey string, runDate time.Time, ownerID string, metadata map[string]any, actorID *uuid.UUID) (*domain.ScheduledJobRun, error) {
	payload, err := marshalJobMetadata(metadata)
	if err != nil {
		return nil, s.logDBError(ctx, "marshal job run metadata", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	row, err := s.getQueries(ctx).UpsertJobRunStarted(ctx, sqlc.UpsertJobRunStartedParams{
		TenantID:  tenantID,
		JobKey:    jobKey,
		RunDate:   dateFromTime(runDate),
		OwnerID:   textFromString(ownerID),
		Metadata:  payload,
		CreatedBy: uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "start job run", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	item, err := mapScheduledJobRun(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map started job run", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	return item, nil
}

func (s *Store) FinishJobRun(ctx context.Context, tenantID uuid.UUID, jobKey string, runDate time.Time, status string, processed int32, succeeded int32, failed int32, skipped int32, errorMessage *string, metadata map[string]any, actorID *uuid.UUID) (*domain.ScheduledJobRun, error) {
	payload, err := marshalJobMetadata(metadata)
	if err != nil {
		return nil, s.logDBError(ctx, "marshal finished job run metadata", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	row, err := s.getQueries(ctx).FinishJobRun(ctx, sqlc.FinishJobRunParams{
		TenantID:       tenantID,
		JobKey:         jobKey,
		RunDate:        dateFromTime(runDate),
		Status:         status,
		ProcessedCount: processed,
		SuccessCount:   succeeded,
		FailedCount:    failed,
		SkippedCount:   skipped,
		ErrorMessage:   textFromPtr(errorMessage),
		Metadata:       payload,
		UpdatedBy:      uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "finish job run", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	item, err := mapScheduledJobRun(row)
	if err != nil {
		return nil, s.logDBError(ctx, "map finished job run", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	return item, nil
}

func (s *Store) ListJobRuns(ctx context.Context, tenantID uuid.UUID, jobKey string, limit int32, offset int32) ([]*domain.ScheduledJobRun, error) {
	rows, err := s.getQueries(ctx).ListJobRuns(ctx, sqlc.ListJobRunsParams{TenantID: tenantID, JobKey: jobKey, Limit: limit, Offset: offset})
	if err != nil {
		return nil, s.logDBError(ctx, "list job runs", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	items, err := mapScheduledJobRuns(rows)
	if err != nil {
		return nil, s.logDBError(ctx, "map job runs", err, tenantIDField(tenantID), stringField("job_key", jobKey))
	}
	return items, nil
}

func marshalJobMetadata(metadata map[string]any) ([]byte, error) {
	if metadata == nil {
		metadata = map[string]any{}
	}
	payload, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("hrms: marshal scheduled job metadata: %w", err)
	}
	return payload, nil
}
