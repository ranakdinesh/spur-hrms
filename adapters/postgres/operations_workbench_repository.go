package postgres

import (
	"context"

	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) ListOperationsWorkbenchCards(ctx context.Context, filter domain.OperationsWorkbenchFilter) ([]*domain.OperationsWorkbenchCard, error) {
	rows, err := s.getQueries(ctx).ListOperationsWorkbenchCards(ctx, sqlc.ListOperationsWorkbenchCardsParams{
		TenantID: filter.TenantID,
		Lane:     textFromPtr(filter.Lane),
		Category: textFromPtr(filter.Category),
		Severity: textFromPtr(filter.Severity),
		Search:   textFromPtr(filter.Search),
		Limit:    limitOrDefault(filter.Limit),
		Offset:   filter.Offset,
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list operations workbench cards", err, tenantIDField(filter.TenantID))
	}
	return mapOperationsWorkbenchCards(rows), nil
}
