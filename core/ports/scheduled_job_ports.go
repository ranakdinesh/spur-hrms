package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type ScheduledJobRepo interface {
	ListTenantProfiles(ctx context.Context) ([]*domain.TenantProfile, error)
	AcquireJobLock(ctx context.Context, tenantID uuid.UUID, jobKey string, ownerID string, ttlSeconds int32) (bool, error)
	ReleaseJobLock(ctx context.Context, tenantID uuid.UUID, jobKey string, ownerID string) error
	GetJobRunByDate(ctx context.Context, tenantID uuid.UUID, jobKey string, runDate time.Time) (*domain.ScheduledJobRun, error)
	StartJobRun(ctx context.Context, tenantID uuid.UUID, jobKey string, runDate time.Time, ownerID string, metadata map[string]any, actorID *uuid.UUID) (*domain.ScheduledJobRun, error)
	FinishJobRun(ctx context.Context, tenantID uuid.UUID, jobKey string, runDate time.Time, status string, processed int32, succeeded int32, failed int32, skipped int32, errorMessage *string, metadata map[string]any, actorID *uuid.UUID) (*domain.ScheduledJobRun, error)
	ListJobRuns(ctx context.Context, tenantID uuid.UUID, jobKey string, limit int32, offset int32) ([]*domain.ScheduledJobRun, error)
}

type CelebrationNotifier interface {
	NotifyCelebration(ctx context.Context, notification domain.CelebrationNotification) error
}

type RunCelebrationDailyJobCommand struct {
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
	Date     string     `json:"date,omitempty"`
	Force    bool       `json:"force"`
	OwnerID  string     `json:"owner_id,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type RunScheduledJobCommand struct {
	TenantID *uuid.UUID `json:"tenant_id,omitempty"`
	JobKey   string     `json:"job_key,omitempty"`
	Date     string     `json:"date,omitempty"`
	Force    bool       `json:"force"`
	OwnerID  string     `json:"owner_id,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}
