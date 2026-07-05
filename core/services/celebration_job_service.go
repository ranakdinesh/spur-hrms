package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

const celebrationJobLockTTLSeconds int32 = 900

func (s *TenantService) RunCelebrationDailyJob(ctx context.Context, cmd ports.RunCelebrationDailyJobCommand) (*domain.CelebrationJobResult, error) {
	runDate, err := parseJobRunDate(cmd.Date)
	if err != nil {
		s.logError("validate celebration daily job date", err, serviceStringField("date", cmd.Date))
		return nil, err
	}
	ownerID := strings.TrimSpace(cmd.OwnerID)
	if ownerID == "" {
		ownerID = defaultJobOwnerID()
	}
	result := &domain.CelebrationJobResult{JobKey: domain.JobKeyCelebrationDaily, RunDate: runDate}
	if cmd.TenantID != nil && *cmd.TenantID != uuid.Nil {
		tenantResult, err := s.runCelebrationDailyJobForTenant(ctx, *cmd.TenantID, runDate, ownerID, cmd.Force, cmd.ActorID)
		if err != nil {
			return nil, err
		}
		mergeCelebrationJobResult(result, tenantResult)
		return result, nil
	}
	if err := s.system.RunAsSystem(ctx, func(systemCtx context.Context) error {
		tenants, err := s.scheduledJobs.ListTenantProfiles(systemCtx)
		if err != nil {
			s.logError("list tenants for celebration daily job", err)
			return err
		}
		for _, tenant := range tenants {
			if tenant == nil || tenant.TenantID == uuid.Nil {
				continue
			}
			tenantResult, err := s.runCelebrationDailyJobForTenant(systemCtx, tenant.TenantID, runDate, ownerID, cmd.Force, cmd.ActorID)
			if err != nil {
				return err
			}
			mergeCelebrationJobResult(result, tenantResult)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListCelebrationDailyJobRuns(ctx context.Context, tenantID uuid.UUID, limit int32, offset int32) ([]*domain.ScheduledJobRun, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate celebration daily job runs tenant", err)
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	items, err := s.scheduledJobs.ListJobRuns(ctx, tenantID, domain.JobKeyCelebrationDaily, limit, offset)
	if err != nil {
		s.logError("list celebration daily job runs", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) runCelebrationDailyJobForTenant(ctx context.Context, tenantID uuid.UUID, runDate time.Time, ownerID string, force bool, actorID *uuid.UUID) (*domain.CelebrationJobResult, error) {
	result := &domain.CelebrationJobResult{JobKey: domain.JobKeyCelebrationDaily, RunDate: runDate}
	acquired, err := s.scheduledJobs.AcquireJobLock(ctx, tenantID, domain.JobKeyCelebrationDaily, ownerID, celebrationJobLockTTLSeconds)
	if err != nil {
		s.logError("acquire celebration daily job lock", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if !acquired {
		result.Skipped = 1
		return result, nil
	}
	defer func() {
		if err := s.scheduledJobs.ReleaseJobLock(ctx, tenantID, domain.JobKeyCelebrationDaily, ownerID); err != nil {
			s.logError("release celebration daily job lock", err, serviceTenantIDField(tenantID))
		}
	}()
	existing, err := s.scheduledJobs.GetJobRunByDate(ctx, tenantID, domain.JobKeyCelebrationDaily, runDate)
	if err != nil {
		s.logError("get celebration daily job run", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if existing != nil && existing.Status == domain.JobStatusSucceeded && !force {
		result.Runs = append(result.Runs, existing)
		result.Skipped = 1
		return result, nil
	}
	metadata := map[string]any{
		"date":                runDate.Format("2006-01-02"),
		"force":               force,
		"owner_id":            ownerID,
		"notifier_configured": s.celebrationNotifier != nil,
	}
	started, err := s.scheduledJobs.StartJobRun(ctx, tenantID, domain.JobKeyCelebrationDaily, runDate, ownerID, metadata, actorID)
	if err != nil {
		s.logError("start celebration daily job run", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result.Runs = append(result.Runs, started)

	notifications, processed, skipped, err := s.buildCelebrationNotifications(ctx, tenantID, runDate)
	if err != nil {
		failMessage := err.Error()
		failedRun, finishErr := s.scheduledJobs.FinishJobRun(ctx, tenantID, domain.JobKeyCelebrationDaily, runDate, domain.JobStatusFailed, processed, 0, 1, skipped, &failMessage, metadata, actorID)
		if finishErr != nil {
			s.logError("finish failed celebration daily job run", finishErr, serviceTenantIDField(tenantID))
		} else {
			result.Runs = append(result.Runs, failedRun)
		}
		return nil, err
	}

	var succeeded int32
	var failed int32
	for _, notification := range notifications {
		if notification == nil {
			continue
		}
		if s.celebrationNotifier == nil {
			skipped++
			continue
		}
		if err := s.celebrationNotifier.NotifyCelebration(ctx, *notification); err != nil {
			failed++
			s.logError("send celebration notification", err, serviceTenantIDField(tenantID), serviceStringField("user_id", notification.UserID.String()))
			continue
		}
		succeeded++
		result.Notifications = append(result.Notifications, notification)
	}
	status := domain.JobStatusSucceeded
	var errorMessage *string
	if failed > 0 {
		status = domain.JobStatusFailed
		msg := fmt.Sprintf("%d celebration notifications failed", failed)
		errorMessage = &msg
	} else if processed == 0 || succeeded == 0 && skipped > 0 {
		status = domain.JobStatusSkipped
	}
	metadata["due_notifications"] = processed
	metadata["success_count"] = succeeded
	metadata["failed_count"] = failed
	metadata["skipped_count"] = skipped
	finished, err := s.scheduledJobs.FinishJobRun(ctx, tenantID, domain.JobKeyCelebrationDaily, runDate, status, processed, succeeded, failed, skipped, errorMessage, metadata, actorID)
	if err != nil {
		s.logError("finish celebration daily job run", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result.Runs = append(result.Runs, finished)
	result.Processed = processed
	result.Succeeded = succeeded
	result.Failed = failed
	result.Skipped = skipped
	return result, nil
}

func (s *TenantService) buildCelebrationNotifications(ctx context.Context, tenantID uuid.UUID, runDate time.Time) ([]*domain.CelebrationNotification, int32, int32, error) {
	celebrations, err := s.ListCelebrations(ctx, tenantID)
	if err != nil {
		s.logError("load celebration daily items", err, serviceTenantIDField(tenantID))
		return nil, 0, 0, err
	}
	employees, err := s.ListEmployees(ctx, tenantID)
	if err != nil {
		s.logError("load celebration daily employees", err, serviceTenantIDField(tenantID))
		return nil, 0, 0, err
	}
	employeeByUser := map[uuid.UUID]*domain.EmployeeListItem{}
	activeEmployees := make([]*domain.EmployeeListItem, 0, len(employees))
	for _, employee := range employees {
		if !isActiveEmployeeForDate(employee, runDate) {
			continue
		}
		employeeByUser[employee.UserID] = employee
		activeEmployees = append(activeEmployees, employee)
	}
	var notifications []*domain.CelebrationNotification
	var processed int32
	var skipped int32
	for _, celebration := range celebrations {
		if celebration == nil || celebration.CelebrationDate == nil || !celebrationDueOn(*celebration.CelebrationDate, runDate, celebration.IsYearly) {
			continue
		}
		if celebration.UserID != nil && *celebration.UserID != uuid.Nil {
			processed++
			employee := employeeByUser[*celebration.UserID]
			if employee == nil {
				skipped++
				continue
			}
			notifications = append(notifications, celebrationNotification(tenantID, celebration, employee, runDate, false))
			continue
		}
		for _, employee := range activeEmployees {
			if employee == nil {
				continue
			}
			if celebration.BranchID != nil && (employee.BranchID == nil || *employee.BranchID != *celebration.BranchID) {
				continue
			}
			processed++
			notifications = append(notifications, celebrationNotification(tenantID, celebration, employee, runDate, true))
		}
	}
	return notifications, processed, skipped, nil
}

func celebrationNotification(tenantID uuid.UUID, celebration *domain.CelebrationListItem, employee *domain.EmployeeListItem, runDate time.Time, isGroup bool) *domain.CelebrationNotification {
	title := strings.TrimSpace(valueFromPtr(celebration.CustomTitle))
	if title == "" {
		title = celebration.CelebrationTypeName
	}
	name := employeeDisplayName(employee.Firstname, employee.MiddleName, employee.Lastname)
	message := strings.TrimSpace(valueFromPtr(celebration.Description))
	if message == "" {
		if !isGroup && celebration.UserID != nil && *celebration.UserID == employee.UserID {
			if strings.Contains(strings.ToLower(title), "birthday") {
				message = fmt.Sprintf("Happy Birthday, %s.", name)
			} else {
				message = fmt.Sprintf("Your %s is today.", strings.ToLower(title))
			}
		} else {
			message = fmt.Sprintf("%s: %s", title, name)
		}
	}
	channels := []string{domain.NotifChannelPush}
	if employee.Email != nil && strings.TrimSpace(*employee.Email) != "" {
		channels = append(channels, domain.NotifChannelEmail)
	}
	return &domain.CelebrationNotification{
		TenantID:          tenantID,
		CelebrationID:     celebration.ID,
		CelebrationTypeID: celebration.CelebrationTypeID,
		UserID:            employee.UserID,
		EmployeeName:      name,
		EmployeeEmail:     employee.Email,
		Title:             title,
		Message:           message,
		Channels:          channels,
		ReferenceTable:    domain.RefTableUserCelebration,
		ReferenceID:       celebration.ID,
		IsGroup:           isGroup,
		RunDate:           runDate,
	}
}

func celebrationDueOn(source time.Time, runDate time.Time, yearly bool) bool {
	source = normalizeDate(source)
	runDate = normalizeDate(runDate)
	if yearly {
		return source.Month() == runDate.Month() && source.Day() == runDate.Day()
	}
	return source.Equal(runDate)
}

func isActiveEmployeeForDate(employee *domain.EmployeeListItem, runDate time.Time) bool {
	if employee == nil || employee.Inactive {
		return false
	}
	if employee.ResignationDate != nil && normalizeDate(*employee.ResignationDate).Before(normalizeDate(runDate)) {
		return false
	}
	return true
}

func parseJobRunDate(value string) (time.Time, error) {
	parsed, err := parseOptionalDate(value)
	if err != nil {
		return time.Time{}, err
	}
	if parsed == nil {
		return normalizeDate(time.Now().UTC()), nil
	}
	return normalizeDate(*parsed), nil
}

func defaultJobOwnerID() string {
	hostname, err := os.Hostname()
	if err != nil || strings.TrimSpace(hostname) == "" {
		hostname = "hrms"
	}
	return fmt.Sprintf("%s-%d", hostname, time.Now().UTC().UnixNano())
}

func mergeCelebrationJobResult(r *domain.CelebrationJobResult, other *domain.CelebrationJobResult) {
	if r == nil || other == nil {
		return
	}
	r.Runs = append(r.Runs, other.Runs...)
	r.Notifications = append(r.Notifications, other.Notifications...)
	r.Processed += other.Processed
	r.Succeeded += other.Succeeded
	r.Failed += other.Failed
	r.Skipped += other.Skipped
}
