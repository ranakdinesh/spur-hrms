package postgres

import (
	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapOperationsWorkbenchCard(row sqlc.ListOperationsWorkbenchCardsRow) *domain.OperationsWorkbenchCard {
	var employeeUserID *uuid.UUID
	if row.EmployeeUserID != uuid.Nil {
		employeeUserID = &row.EmployeeUserID
	}
	return &domain.OperationsWorkbenchCard{
		TenantID:       row.TenantID,
		CardKey:        row.CardKey,
		Lane:           row.Lane,
		Category:       row.Category,
		SourceModule:   row.SourceModule,
		SourceType:     row.SourceType,
		SourceID:       row.SourceID,
		EmployeeUserID: employeeUserID,
		Title:          row.Title,
		Summary:        row.Summary,
		Status:         row.Status,
		Severity:       row.Severity,
		Priority:       row.Priority,
		DueAt:          ptrFromTimestamptz(row.DueAt),
		DetectedAt:     timeFromTimestamptz(row.DetectedAt),
		ActionLabel:    row.ActionLabel,
		RouteSection:   row.RouteSection,
		RouteRecordID:  &row.RouteRecordID,
		Metadata:       jsonRawDefault(row.Metadata, `{}`),
	}
}

func mapOperationsWorkbenchCards(rows []sqlc.ListOperationsWorkbenchCardsRow) []*domain.OperationsWorkbenchCard {
	items := make([]*domain.OperationsWorkbenchCard, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapOperationsWorkbenchCard(row))
	}
	return items
}
