package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapScheduledJobRun(row sqlc.HrmsJobRun) (*domain.ScheduledJobRun, error) {
	item := &domain.ScheduledJobRun{
		ID:             row.ID,
		TenantID:       row.TenantID,
		JobKey:         row.JobKey,
		RunDate:        timeFromDate(row.RunDate),
		Status:         row.Status,
		OwnerID:        ptrFromText(row.OwnerID),
		StartedAt:      timeFromTimestamptz(row.StartedAt),
		FinishedAt:     ptrFromTimestamptz(row.FinishedAt),
		ProcessedCount: row.ProcessedCount,
		SuccessCount:   row.SuccessCount,
		FailedCount:    row.FailedCount,
		SkippedCount:   row.SkippedCount,
		ErrorMessage:   ptrFromText(row.ErrorMessage),
		Metadata:       map[string]any{},
		Inactive:       row.Inactive,
		CreatedAt:      timeFromTimestamptz(row.CreatedAt),
		CreatedBy:      ptrFromUUID(row.CreatedBy),
		UpdatedAt:      timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:      ptrFromUUID(row.UpdatedBy),
	}
	if len(row.Metadata) > 0 {
		if err := json.Unmarshal(row.Metadata, &item.Metadata); err != nil {
			return nil, fmt.Errorf("hrms: unmarshal scheduled job metadata: %w", err)
		}
	}
	return item, nil
}

func mapScheduledJobRuns(rows []sqlc.HrmsJobRun) ([]*domain.ScheduledJobRun, error) {
	items := make([]*domain.ScheduledJobRun, 0, len(rows))
	for _, row := range rows {
		item, err := mapScheduledJobRun(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
