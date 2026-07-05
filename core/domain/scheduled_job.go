package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	JobKeyCelebrationDaily      = "celebration_daily"
	JobKeyLeaveAccrualMonthly   = "leave_accrual_monthly"
	JobKeyHolidayReminderDaily  = "holiday_reminder_daily"
	JobKeySalaryReminderMonthly = "salary_reminder_monthly"

	JobStatusRunning   = "running"
	JobStatusSucceeded = "succeeded"
	JobStatusFailed    = "failed"
	JobStatusSkipped   = "skipped"
)

type ScheduledJobRun struct {
	ID             uuid.UUID      `json:"id"`
	TenantID       uuid.UUID      `json:"tenant_id"`
	JobKey         string         `json:"job_key"`
	RunDate        time.Time      `json:"run_date"`
	Status         string         `json:"status"`
	OwnerID        *string        `json:"owner_id,omitempty"`
	StartedAt      time.Time      `json:"started_at"`
	FinishedAt     *time.Time     `json:"finished_at,omitempty"`
	ProcessedCount int32          `json:"processed_count"`
	SuccessCount   int32          `json:"success_count"`
	FailedCount    int32          `json:"failed_count"`
	SkippedCount   int32          `json:"skipped_count"`
	ErrorMessage   *string        `json:"error_message,omitempty"`
	Metadata       map[string]any `json:"metadata"`
	Inactive       bool           `json:"inactive"`
	CreatedAt      time.Time      `json:"created_at"`
	CreatedBy      *uuid.UUID     `json:"created_by,omitempty"`
	UpdatedAt      time.Time      `json:"updated_at"`
	UpdatedBy      *uuid.UUID     `json:"updated_by,omitempty"`
}

type RegisteredScheduledJob struct {
	Key             string   `json:"key"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Schedule        string   `json:"schedule"`
	Timezone        string   `json:"timezone"`
	DefaultRunTime  string   `json:"default_run_time"`
	RecommendedMode string   `json:"recommended_mode"`
	IdempotencyKey  string   `json:"idempotency_key"`
	BackfillReady   bool     `json:"backfill_ready"`
	Channels        []string `json:"channels"`
}

type CelebrationNotification struct {
	TenantID          uuid.UUID `json:"tenant_id"`
	CelebrationID     uuid.UUID `json:"celebration_id"`
	CelebrationTypeID uuid.UUID `json:"celebration_type_id"`
	UserID            uuid.UUID `json:"user_id"`
	EmployeeName      string    `json:"employee_name"`
	EmployeeEmail     *string   `json:"employee_email,omitempty"`
	Title             string    `json:"title"`
	Message           string    `json:"message"`
	Channels          []string  `json:"channels"`
	ReferenceTable    string    `json:"reference_table"`
	ReferenceID       uuid.UUID `json:"reference_id"`
	IsGroup           bool      `json:"is_group"`
	RunDate           time.Time `json:"run_date"`
}

type CelebrationJobResult struct {
	JobKey        string                     `json:"job_key"`
	RunDate       time.Time                  `json:"run_date"`
	Runs          []*ScheduledJobRun         `json:"runs"`
	Notifications []*CelebrationNotification `json:"notifications"`
	Processed     int32                      `json:"processed"`
	Succeeded     int32                      `json:"succeeded"`
	Failed        int32                      `json:"failed"`
	Skipped       int32                      `json:"skipped"`
}

type ScheduledJobResult struct {
	JobKey    string             `json:"job_key"`
	RunDate   time.Time          `json:"run_date"`
	Runs      []*ScheduledJobRun `json:"runs"`
	Processed int32              `json:"processed"`
	Succeeded int32              `json:"succeeded"`
	Failed    int32              `json:"failed"`
	Skipped   int32              `json:"skipped"`
}
