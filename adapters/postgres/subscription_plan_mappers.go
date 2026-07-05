package postgres

import (
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapSubscriptionPlan(row sqlc.HrmsSubscriptionPlan) *domain.SubscriptionPlan {
	return &domain.SubscriptionPlan{
		ID:                row.ID,
		Code:              row.Code,
		Name:              row.Name,
		Description:       ptrFromText(row.Description),
		PriceAmount:       floatFromNumeric(row.PriceAmount),
		PriceBasis:        row.PriceBasis,
		MinimumAmount:     floatFromNumeric(row.MinimumAmount),
		IncludedEmployees: row.IncludedEmployees,
		OverageAmount:     floatFromNumeric(row.OverageAmount),
		CurrencyCode:      row.CurrencyCode,
		BillingCycle:      row.BillingCycle,
		EmployeeLimit:     row.EmployeeLimit,
		TrialDays:         row.TrialDays,
		Visibility:        row.Visibility,
		IsActive:          row.IsActive,
		Inactive:          row.Inactive,
		CreatedAt:         timeFromTimestamptz(row.CreatedAt),
		CreatedBy:         ptrFromUUID(row.CreatedBy),
		UpdatedAt:         timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:         ptrFromUUID(row.UpdatedBy),
	}
}

func mapSubscriptionPlans(rows []sqlc.HrmsSubscriptionPlan) []*domain.SubscriptionPlan {
	items := make([]*domain.SubscriptionPlan, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapSubscriptionPlan(row))
	}
	return items
}

func numericFromMoney(value float64) pgtype.Numeric {
	scaled := int64(math.Round(value * 100))
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -2, Valid: true}
}
