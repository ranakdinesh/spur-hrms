package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

const registeredJobLockTTLSeconds int32 = 1800

var errUnknownScheduledJob = errors.New("scheduled job is not registered")

type scheduledJobExecutor func(context.Context, uuid.UUID, time.Time, *uuid.UUID) (scheduledJobCounts, map[string]any, error)

type scheduledJobCounts struct {
	processed int32
	succeeded int32
	failed    int32
	skipped   int32
}

func (s *TenantService) ListRegisteredScheduledJobs(ctx context.Context) ([]domain.RegisteredScheduledJob, error) {
	return registeredScheduledJobs(), nil
}

func (s *TenantService) RunScheduledJob(ctx context.Context, cmd ports.RunScheduledJobCommand) (*domain.ScheduledJobResult, error) {
	jobKey := strings.TrimSpace(cmd.JobKey)
	executor, err := s.scheduledJobExecutor(jobKey)
	if err != nil {
		s.logError("validate scheduled job key", err, serviceStringField("job_key", jobKey))
		return nil, err
	}
	if jobKey == domain.JobKeyCelebrationDaily {
		result, err := s.RunCelebrationDailyJob(ctx, ports.RunCelebrationDailyJobCommand{TenantID: cmd.TenantID, Date: cmd.Date, Force: cmd.Force, OwnerID: cmd.OwnerID, ActorID: cmd.ActorID})
		if err != nil {
			return nil, err
		}
		return &domain.ScheduledJobResult{JobKey: result.JobKey, RunDate: result.RunDate, Runs: result.Runs, Processed: result.Processed, Succeeded: result.Succeeded, Failed: result.Failed, Skipped: result.Skipped}, nil
	}
	runDate, err := parseJobRunDate(cmd.Date)
	if err != nil {
		s.logError("validate scheduled job date", err, serviceStringField("job_key", jobKey), serviceStringField("date", cmd.Date))
		return nil, err
	}
	ownerID := strings.TrimSpace(cmd.OwnerID)
	if ownerID == "" {
		ownerID = defaultJobOwnerID()
	}
	result := &domain.ScheduledJobResult{JobKey: jobKey, RunDate: runDate}
	if cmd.TenantID != nil && *cmd.TenantID != uuid.Nil {
		tenantResult, err := s.runRegisteredScheduledJobForTenant(ctx, *cmd.TenantID, jobKey, runDate, ownerID, cmd.Force, cmd.ActorID, executor)
		if err != nil {
			return nil, err
		}
		mergeScheduledJobResult(result, tenantResult)
		return result, nil
	}
	if err := s.system.RunAsSystem(ctx, func(systemCtx context.Context) error {
		tenants, err := s.scheduledJobs.ListTenantProfiles(systemCtx)
		if err != nil {
			s.logError("list tenants for registered scheduled job", err, serviceStringField("job_key", jobKey))
			return err
		}
		for _, tenant := range tenants {
			if tenant == nil || tenant.TenantID == uuid.Nil {
				continue
			}
			tenantResult, err := s.runRegisteredScheduledJobForTenant(systemCtx, tenant.TenantID, jobKey, runDate, ownerID, cmd.Force, cmd.ActorID, executor)
			if err != nil {
				return err
			}
			mergeScheduledJobResult(result, tenantResult)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListScheduledJobRuns(ctx context.Context, tenantID uuid.UUID, jobKey string, limit int32, offset int32) ([]*domain.ScheduledJobRun, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate scheduled job runs tenant", err, serviceStringField("job_key", jobKey))
		return nil, err
	}
	if _, err := s.scheduledJobExecutor(jobKey); err != nil {
		s.logError("validate scheduled job runs key", err, serviceTenantIDField(tenantID), serviceStringField("job_key", jobKey))
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	items, err := s.scheduledJobs.ListJobRuns(ctx, tenantID, strings.TrimSpace(jobKey), limit, offset)
	if err != nil {
		s.logError("list scheduled job runs", err, serviceTenantIDField(tenantID), serviceStringField("job_key", jobKey))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) scheduledJobExecutor(jobKey string) (scheduledJobExecutor, error) {
	switch strings.TrimSpace(jobKey) {
	case domain.JobKeyCelebrationDaily:
		return nil, nil
	case domain.JobKeyLeaveAccrualMonthly:
		return s.executeLeaveAccrualMonthlyJob, nil
	case domain.JobKeyHolidayReminderDaily:
		return s.executeHolidayReminderDailyJob, nil
	case domain.JobKeySalaryReminderMonthly:
		return s.executeSalaryReminderMonthlyJob, nil
	default:
		return nil, errUnknownScheduledJob
	}
}

func (s *TenantService) runRegisteredScheduledJobForTenant(ctx context.Context, tenantID uuid.UUID, jobKey string, runDate time.Time, ownerID string, force bool, actorID *uuid.UUID, executor scheduledJobExecutor) (*domain.ScheduledJobResult, error) {
	result := &domain.ScheduledJobResult{JobKey: jobKey, RunDate: runDate}
	acquired, err := s.scheduledJobs.AcquireJobLock(ctx, tenantID, jobKey, ownerID, registeredJobLockTTLSeconds)
	if err != nil {
		s.logError("acquire registered scheduled job lock", err, serviceTenantIDField(tenantID), serviceStringField("job_key", jobKey))
		return nil, err
	}
	if !acquired {
		result.Skipped = 1
		return result, nil
	}
	defer func() {
		if err := s.scheduledJobs.ReleaseJobLock(ctx, tenantID, jobKey, ownerID); err != nil {
			s.logError("release registered scheduled job lock", err, serviceTenantIDField(tenantID), serviceStringField("job_key", jobKey))
		}
	}()
	existing, err := s.scheduledJobs.GetJobRunByDate(ctx, tenantID, jobKey, runDate)
	if err != nil {
		s.logError("get registered scheduled job run", err, serviceTenantIDField(tenantID), serviceStringField("job_key", jobKey))
		return nil, err
	}
	if existing != nil && existing.Status == domain.JobStatusSucceeded && !force {
		result.Runs = append(result.Runs, existing)
		result.Skipped = 1
		return result, nil
	}
	metadata := map[string]any{"date": runDate.Format("2006-01-02"), "force": force, "owner_id": ownerID, "job_key": jobKey}
	started, err := s.scheduledJobs.StartJobRun(ctx, tenantID, jobKey, runDate, ownerID, metadata, actorID)
	if err != nil {
		s.logError("start registered scheduled job run", err, serviceTenantIDField(tenantID), serviceStringField("job_key", jobKey))
		return nil, err
	}
	result.Runs = append(result.Runs, started)

	counts, extraMetadata, err := executor(ctx, tenantID, runDate, actorID)
	for key, value := range extraMetadata {
		metadata[key] = value
	}
	status := domain.JobStatusSucceeded
	var errorMessage *string
	if err != nil {
		status = domain.JobStatusFailed
		message := err.Error()
		errorMessage = &message
		counts.failed++
	} else if counts.processed == 0 || counts.succeeded == 0 && counts.skipped > 0 {
		status = domain.JobStatusSkipped
	}
	finished, finishErr := s.scheduledJobs.FinishJobRun(ctx, tenantID, jobKey, runDate, status, counts.processed, counts.succeeded, counts.failed, counts.skipped, errorMessage, metadata, actorID)
	if finishErr != nil {
		s.logError("finish registered scheduled job run", finishErr, serviceTenantIDField(tenantID), serviceStringField("job_key", jobKey))
		return nil, finishErr
	}
	result.Runs = append(result.Runs, finished)
	result.Processed = counts.processed
	result.Succeeded = counts.succeeded
	result.Failed = counts.failed
	result.Skipped = counts.skipped
	if err != nil {
		s.logError("execute registered scheduled job", err, serviceTenantIDField(tenantID), serviceStringField("job_key", jobKey))
		return result, err
	}
	return result, nil
}

func (s *TenantService) executeLeaveAccrualMonthlyJob(ctx context.Context, tenantID uuid.UUID, runDate time.Time, actorID *uuid.UUID) (scheduledJobCounts, map[string]any, error) {
	fy, err := s.financialYears.GetActiveFinancialYear(ctx, tenantID)
	if err != nil {
		return scheduledJobCounts{}, nil, err
	}
	entries, err := s.RunLeaveAccrual(ctx, ports.RunLeaveAccrualCommand{TenantID: tenantID, FYID: fy.ID, Month: int32(runDate.Month()), ActorID: actorID})
	if err != nil {
		return scheduledJobCounts{}, nil, err
	}
	count := int32(len(entries))
	return scheduledJobCounts{processed: count, succeeded: count}, map[string]any{"financial_year_id": fy.ID.String(), "month": int(runDate.Month()), "ledger_entries": count}, nil
}

func (s *TenantService) executeHolidayReminderDailyJob(ctx context.Context, tenantID uuid.UUID, runDate time.Time, actorID *uuid.UUID) (scheduledJobCounts, map[string]any, error) {
	reminderDate := normalizeDate(runDate).AddDate(0, 0, 1)
	holidays, err := s.holidays.ListHolidaysByDateRange(ctx, tenantID, reminderDate, reminderDate)
	if err != nil {
		return scheduledJobCounts{}, nil, err
	}
	counts := scheduledJobCounts{}
	for _, holiday := range holidays {
		if holiday == nil || holiday.Inactive {
			continue
		}
		counts.processed++
		employees, err := s.employees.ListEmployees(ctx, tenantID)
		if err != nil {
			return counts, nil, err
		}
		recipients := activeEmployeeUserIDsForHoliday(employees, holiday.BranchID, runDate)
		if len(recipients) == 0 {
			counts.skipped++
			continue
		}
		title := "Upcoming holiday"
		message := fmt.Sprintf("%s is on %s.", holiday.Name, reminderDate.Format("02 Jan 2006"))
		referenceTable := domain.RefTableHoliday
		_, err = s.sendNotification(ctx, domain.NotificationSendInput{TenantID: tenantID, NotificationCode: domain.NotifGeneralNotif, UserIDs: recipients, Title: title, Message: message, ReferenceTable: &referenceTable, ReferenceID: &holiday.ID, Channels: []string{domain.NotifChannelPush, domain.NotifChannelEmail}, ActorID: actorID})
		if err != nil {
			counts.failed++
			s.logError("send holiday reminder notification", err, serviceTenantIDField(tenantID), serviceStringField("holiday_id", holiday.ID.String()))
			continue
		}
		counts.succeeded += int32(len(recipients))
	}
	return counts, map[string]any{"reminder_date": reminderDate.Format("2006-01-02"), "holiday_count": len(holidays)}, nil
}

func (s *TenantService) executeSalaryReminderMonthlyJob(ctx context.Context, tenantID uuid.UUID, runDate time.Time, actorID *uuid.UUID) (scheduledJobCounts, map[string]any, error) {
	employees, err := s.employees.ListEmployees(ctx, tenantID)
	if err != nil {
		return scheduledJobCounts{}, nil, err
	}
	activeUsers := activeEmployeeUserIDsForHoliday(employees, nil, runDate)
	month := int32(runDate.Month())
	year := int32(runDate.Year())
	slips, err := s.salarySlips.ListSalarySlipsByTenantPeriod(ctx, tenantID, month, year)
	if err != nil {
		return scheduledJobCounts{}, nil, err
	}
	pending := int32(len(activeUsers) - len(slips))
	if pending < 0 {
		pending = 0
	}
	counts := scheduledJobCounts{processed: int32(len(activeUsers)), succeeded: int32(len(slips)), skipped: pending}
	return counts, map[string]any{"month": month, "year": year, "active_employee_count": len(activeUsers), "salary_slip_count": len(slips), "pending_salary_slip_count": pending, "recommendation": "Generate or review monthly payslips before payout."}, nil
}

func registeredScheduledJobs() []domain.RegisteredScheduledJob {
	return []domain.RegisteredScheduledJob{
		{Key: domain.JobKeyCelebrationDaily, Name: "Celebration Notifications", Description: "Sends birthday, anniversary, festival, and tenant celebration notifications.", Schedule: "Daily", Timezone: "tenant/local", DefaultRunTime: "09:00", RecommendedMode: "forbid_overlap", IdempotencyKey: "tenant_id + job_key + run_date", BackfillReady: true, Channels: []string{domain.NotifChannelPush, domain.NotifChannelEmail}},
		{Key: domain.JobKeyLeaveAccrualMonthly, Name: "Monthly Leave Accrual", Description: "Credits monthly leave accruals from active leave templates and policies.", Schedule: "Monthly", Timezone: "tenant/local", DefaultRunTime: "02:00 on day 1", RecommendedMode: "forbid_overlap", IdempotencyKey: "tenant_id + user_id + leave_type_id + fy_id + month", BackfillReady: true, Channels: []string{"ledger"}},
		{Key: domain.JobKeyHolidayReminderDaily, Name: "Holiday Reminders", Description: "Notifies active employees about holidays due the next day.", Schedule: "Daily", Timezone: "tenant/local", DefaultRunTime: "10:00", RecommendedMode: "forbid_overlap", IdempotencyKey: "tenant_id + job_key + run_date", BackfillReady: true, Channels: []string{domain.NotifChannelPush, domain.NotifChannelEmail}},
		{Key: domain.JobKeySalaryReminderMonthly, Name: "Salary Processing Reminder", Description: "Checks payroll readiness and records pending payslip counts before payout.", Schedule: "Monthly", Timezone: "tenant/local", DefaultRunTime: "09:00 on payroll cutoff", RecommendedMode: "forbid_overlap", IdempotencyKey: "tenant_id + job_key + run_date", BackfillReady: true, Channels: []string{"run_history", "metadata"}},
	}
}

func activeEmployeeUserIDsForHoliday(employees []*domain.EmployeeListItem, branchID *uuid.UUID, runDate time.Time) []uuid.UUID {
	userIDs := make([]uuid.UUID, 0, len(employees))
	for _, employee := range employees {
		if !isActiveEmployeeForDate(employee, runDate) {
			continue
		}
		if branchID != nil && (employee.BranchID == nil || *employee.BranchID != *branchID) {
			continue
		}
		userIDs = append(userIDs, employee.UserID)
	}
	return userIDs
}

func mergeScheduledJobResult(r *domain.ScheduledJobResult, other *domain.ScheduledJobResult) {
	if r == nil || other == nil {
		return
	}
	r.Runs = append(r.Runs, other.Runs...)
	r.Processed += other.Processed
	r.Succeeded += other.Succeeded
	r.Failed += other.Failed
	r.Skipped += other.Skipped
}
