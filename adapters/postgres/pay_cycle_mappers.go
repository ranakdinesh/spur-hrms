package postgres

import (
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapPayCycle(row sqlc.HrmsPayCycle) *domain.PayCycle {
	return &domain.PayCycle{
		ID:                     row.ID,
		TenantID:               row.TenantID,
		Name:                   row.Name,
		CycleType:              row.CycleType,
		PayDay:                 ptrFromInt4(row.PayDay),
		StartDay:               ptrFromInt4(row.StartDay),
		EndDay:                 ptrFromInt4(row.EndDay),
		AttendanceSource:       row.AttendanceSource,
		AttendancePeriodType:   row.AttendancePeriodType,
		AttendanceCutoffDay:    row.AttendanceCutoffDay,
		PayoutTiming:           row.PayoutTiming,
		PayoutOffsetDays:       row.PayoutOffsetDays,
		IncludeWeeklyOffs:      row.IncludeWeeklyOffs,
		IncludeHolidays:        row.IncludeHolidays,
		ProrateJoiningExit:     row.ProrateJoiningExit,
		ProrationBasis:         row.ProrationBasis,
		AllowArrears:           row.AllowArrears,
		ArrearsMode:            row.ArrearsMode,
		AllowNegativeNetPay:    row.AllowNegativeNetPay,
		OvertimeComponentCode:  ptrFromText(row.OvertimeComponentCode),
		LWPComponentCode:       row.LwpComponentCode,
		RoundingMode:           row.RoundingMode,
		PaymentMode:            row.PaymentMode,
		PaymentFileFormat:      row.PaymentFileFormat,
		RequiresApproval:       row.RequiresApproval,
		AutoLockAfterApproval:  row.AutoLockAfterApproval,
		PayrollLockDay:         ptrFromInt4(row.PayrollLockDay),
		PFEnabled:              row.PfEnabled,
		PFEmployeeRate:         floatFromNumeric(row.PfEmployeeRate),
		PFEmployerRate:         floatFromNumeric(row.PfEmployerRate),
		PFWageCeiling:          floatFromNumeric(row.PfWageCeiling),
		PFApplyCeiling:         row.PfApplyCeiling,
		ESIEnabled:             row.EsiEnabled,
		ESIEmployeeRate:        floatFromNumeric(row.EsiEmployeeRate),
		ESIEmployerRate:        floatFromNumeric(row.EsiEmployerRate),
		ESIWageCeiling:         floatFromNumeric(row.EsiWageCeiling),
		ProfessionalTaxEnabled: row.ProfessionalTaxEnabled,
		TDSEnabled:             row.TdsEnabled,
		CountryCode:            row.CountryCode,
		StateCode:              ptrFromText(row.StateCode),
		Notes:                  ptrFromText(row.Notes),
		Inactive:               row.Inactive,
		CreatedAt:              timeFromTimestamptz(row.CreatedAt),
		CreatedBy:              ptrFromUUID(row.CreatedBy),
		UpdatedAt:              timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:              ptrFromUUID(row.UpdatedBy),
	}
}

func numericFromPayrollFloat(value float64) pgtype.Numeric {
	scaled := int64(math.Round(value * 10000))
	return pgtype.Numeric{Int: big.NewInt(scaled), Exp: -4, Valid: true}
}
