package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapCelebrationType(row sqlc.HrmsCelebrationType) *domain.CelebrationType {
	return &domain.CelebrationType{
		ID:                row.ID,
		TenantID:          row.TenantID,
		Name:              row.Name,
		IsYearly:          row.IsYearly,
		IsUserCelebration: row.IsUserCelebration,
		Inactive:          row.Inactive,
		CreatedAt:         timeFromTimestamptz(row.CreatedAt),
		CreatedBy:         ptrFromUUID(row.CreatedBy),
		UpdatedAt:         timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:         ptrFromUUID(row.UpdatedBy),
	}
}

func mapCelebrationTypes(rows []sqlc.HrmsCelebrationType) []*domain.CelebrationType {
	items := make([]*domain.CelebrationType, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCelebrationType(row))
	}
	return items
}

func mapCelebration(row sqlc.HrmsCelebration) *domain.Celebration {
	return &domain.Celebration{
		ID:                row.ID,
		TenantID:          row.TenantID,
		BranchID:          ptrFromUUID(row.BranchID),
		UserID:            ptrFromUUID(row.UserID),
		CelebrationTypeID: row.CelebrationTypeID,
		CelebrationDate:   ptrFromDate(row.CelebrationDate),
		CustomTitle:       ptrFromText(row.CustomTitle),
		Description:       ptrFromText(row.Description),
		Inactive:          row.Inactive,
		CreatedAt:         timeFromTimestamptz(row.CreatedAt),
		CreatedBy:         ptrFromUUID(row.CreatedBy),
		UpdatedAt:         timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:         ptrFromUUID(row.UpdatedBy),
	}
}

func mapCelebrations(rows []sqlc.HrmsCelebration) []*domain.Celebration {
	items := make([]*domain.Celebration, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapCelebration(row))
	}
	return items
}
