package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapTenantSubscription(row sqlc.HrmsTenantSubscription) *domain.TenantSubscription {
	return &domain.TenantSubscription{
		ID:           row.ID,
		TenantID:     row.TenantID,
		PlanID:       ptrFromUUID(row.PlanID),
		StartDate:    ptrFromDate(row.StartDate),
		EndDate:      ptrFromDate(row.EndDate),
		Status:       row.Status,
		MaxEmployees: row.MaxEmployees,
		Inactive:     row.Inactive,
		CreatedAt:    timeFromTimestamptz(row.CreatedAt),
		CreatedBy:    ptrFromUUID(row.CreatedBy),
		UpdatedAt:    timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:    ptrFromUUID(row.UpdatedBy),
	}
}

func mapTenantSubscriptions(rows []sqlc.HrmsTenantSubscription) []*domain.TenantSubscription {
	items := make([]*domain.TenantSubscription, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapTenantSubscription(row))
	}
	return items
}
