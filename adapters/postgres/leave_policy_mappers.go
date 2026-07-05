package postgres

import (
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapLeavePolicy(row sqlc.HrmsLeavePolicy) *domain.LeavePolicy {
	return &domain.LeavePolicy{
		ID:                   row.ID,
		TenantID:             row.TenantID,
		LeaveTypeID:          row.LeaveTypeID,
		FYID:                 row.FyID,
		TotalDays:            floatFromNumeric(row.TotalDays),
		AllocationType:       row.AllocationType,
		Jan:                  row.Jan,
		Feb:                  row.Feb,
		Mar:                  row.Mar,
		Apr:                  row.Apr,
		May:                  row.May,
		Jun:                  row.Jun,
		Jul:                  row.Jul,
		Aug:                  row.Aug,
		Sep:                  row.Sep,
		Oct:                  row.Oct,
		Nov:                  row.Nov,
		Dec:                  row.Dec,
		IsSandwichApplicable: row.IsSandwichApplicable,
		Inactive:             row.Inactive,
		CreatedAt:            timeFromTimestamptz(row.CreatedAt),
		CreatedBy:            ptrFromUUID(row.CreatedBy),
		UpdatedAt:            timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:            ptrFromUUID(row.UpdatedBy),
	}
}

func mapLeavePolicies(rows []sqlc.HrmsLeavePolicy) []*domain.LeavePolicy {
	items := make([]*domain.LeavePolicy, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeavePolicy(row))
	}
	return items
}

func numericFromFloat(value float64) pgtype.Numeric {
	scaled := int64(math.Round(value * 10))
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -1, Valid: true}
}

func floatFromNumeric(value pgtype.Numeric) float64 {
	if !value.Valid {
		return 0
	}
	floatValue, err := value.Float64Value()
	if err != nil || !floatValue.Valid {
		return 0
	}
	return floatValue.Float64
}
