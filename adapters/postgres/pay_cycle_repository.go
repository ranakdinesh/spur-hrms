package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) GetPayCycle(ctx context.Context, tenantID uuid.UUID) (*domain.PayCycle, error) {
	row, err := s.getQueries(ctx).GetPayCycle(ctx, tenantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPayCycleNotFound
		}
		return nil, s.logDBError(ctx, "get pay cycle", err, tenantIDField(tenantID))
	}
	return mapPayCycle(row), nil
}

func (s *Store) UpsertPayCycle(ctx context.Context, item *domain.PayCycle, actorID *uuid.UUID) (*domain.PayCycle, error) {
	row, err := s.getQueries(ctx).UpsertPayCycle(ctx, sqlc.UpsertPayCycleParams{
		TenantID:               item.TenantID,
		Name:                   item.Name,
		CycleType:              item.CycleType,
		PayDay:                 int4FromPtr(item.PayDay),
		StartDay:               int4FromPtr(item.StartDay),
		EndDay:                 int4FromPtr(item.EndDay),
		AttendanceSource:       item.AttendanceSource,
		AttendancePeriodType:   item.AttendancePeriodType,
		AttendanceCutoffDay:    item.AttendanceCutoffDay,
		PayoutTiming:           item.PayoutTiming,
		PayoutOffsetDays:       item.PayoutOffsetDays,
		IncludeWeeklyOffs:      item.IncludeWeeklyOffs,
		IncludeHolidays:        item.IncludeHolidays,
		ProrateJoiningExit:     item.ProrateJoiningExit,
		ProrationBasis:         item.ProrationBasis,
		AllowArrears:           item.AllowArrears,
		ArrearsMode:            item.ArrearsMode,
		AllowNegativeNetPay:    item.AllowNegativeNetPay,
		OvertimeComponentCode:  textFromPtr(item.OvertimeComponentCode),
		LwpComponentCode:       item.LWPComponentCode,
		RoundingMode:           item.RoundingMode,
		PaymentMode:            item.PaymentMode,
		PaymentFileFormat:      item.PaymentFileFormat,
		RequiresApproval:       item.RequiresApproval,
		AutoLockAfterApproval:  item.AutoLockAfterApproval,
		PayrollLockDay:         int4FromPtr(item.PayrollLockDay),
		PfEnabled:              item.PFEnabled,
		PfEmployeeRate:         numericFromPayrollFloat(item.PFEmployeeRate),
		PfEmployerRate:         numericFromPayrollFloat(item.PFEmployerRate),
		PfWageCeiling:          numericFromPayrollFloat(item.PFWageCeiling),
		PfApplyCeiling:         item.PFApplyCeiling,
		EsiEnabled:             item.ESIEnabled,
		EsiEmployeeRate:        numericFromPayrollFloat(item.ESIEmployeeRate),
		EsiEmployerRate:        numericFromPayrollFloat(item.ESIEmployerRate),
		EsiWageCeiling:         numericFromPayrollFloat(item.ESIWageCeiling),
		ProfessionalTaxEnabled: item.ProfessionalTaxEnabled,
		TdsEnabled:             item.TDSEnabled,
		CountryCode:            item.CountryCode,
		StateCode:              textFromPtr(item.StateCode),
		Notes:                  textFromPtr(item.Notes),
		CreatedBy:              uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert pay cycle", err, tenantIDField(item.TenantID))
	}
	return mapPayCycle(row), nil
}

func (s *Store) DeletePayCycle(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePayCycle(ctx, sqlc.SoftDeletePayCycleParams{TenantID: tenantID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete pay cycle", err, tenantIDField(tenantID))
	}
	return nil
}
