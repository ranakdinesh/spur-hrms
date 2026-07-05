package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetPayCycle(ctx context.Context, tenantID uuid.UUID) (*domain.PayCycle, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate pay cycle tenant", err)
		return nil, err
	}
	item, err := s.payCycles.GetPayCycle(ctx, tenantID)
	if err != nil {
		s.logError("get pay cycle", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpsertPayCycle(ctx context.Context, cmd ports.PayCycleCommand) (*domain.PayCycle, error) {
	item, err := domain.NewPayCycle(domain.PayCycleInput{
		TenantID:               cmd.TenantID,
		Name:                   cmd.Name,
		CycleType:              cmd.CycleType,
		PayDay:                 cmd.PayDay,
		StartDay:               cmd.StartDay,
		EndDay:                 cmd.EndDay,
		AttendanceSource:       cmd.AttendanceSource,
		AttendancePeriodType:   cmd.AttendancePeriodType,
		AttendanceCutoffDay:    cmd.AttendanceCutoffDay,
		PayoutTiming:           cmd.PayoutTiming,
		PayoutOffsetDays:       cmd.PayoutOffsetDays,
		IncludeWeeklyOffs:      cmd.IncludeWeeklyOffs,
		IncludeHolidays:        cmd.IncludeHolidays,
		ProrateJoiningExit:     cmd.ProrateJoiningExit,
		ProrationBasis:         cmd.ProrationBasis,
		AllowArrears:           cmd.AllowArrears,
		ArrearsMode:            cmd.ArrearsMode,
		AllowNegativeNetPay:    cmd.AllowNegativeNetPay,
		OvertimeComponentCode:  cmd.OvertimeComponentCode,
		LWPComponentCode:       cmd.LWPComponentCode,
		RoundingMode:           cmd.RoundingMode,
		PaymentMode:            cmd.PaymentMode,
		PaymentFileFormat:      cmd.PaymentFileFormat,
		RequiresApproval:       cmd.RequiresApproval,
		AutoLockAfterApproval:  cmd.AutoLockAfterApproval,
		PayrollLockDay:         cmd.PayrollLockDay,
		PFEnabled:              cmd.PFEnabled,
		PFEmployeeRate:         cmd.PFEmployeeRate,
		PFEmployerRate:         cmd.PFEmployerRate,
		PFWageCeiling:          cmd.PFWageCeiling,
		PFApplyCeiling:         cmd.PFApplyCeiling,
		ESIEnabled:             cmd.ESIEnabled,
		ESIEmployeeRate:        cmd.ESIEmployeeRate,
		ESIEmployerRate:        cmd.ESIEmployerRate,
		ESIWageCeiling:         cmd.ESIWageCeiling,
		ProfessionalTaxEnabled: cmd.ProfessionalTaxEnabled,
		TDSEnabled:             cmd.TDSEnabled,
		CountryCode:            cmd.CountryCode,
		StateCode:              cmd.StateCode,
		Notes:                  cmd.Notes,
	})
	if err != nil {
		s.logError("validate pay cycle", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.payCycles.UpsertPayCycle(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert pay cycle", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeletePayCycle(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate pay cycle delete tenant", err)
		return err
	}
	if err := s.payCycles.DeletePayCycle(ctx, tenantID, actorID); err != nil {
		s.logError("delete pay cycle", err, serviceTenantIDField(tenantID))
		return err
	}
	return nil
}

func (s *TenantService) ResolvePayCyclePeriod(ctx context.Context, query ports.PayCyclePeriodQuery) (*domain.PayCyclePeriod, error) {
	if query.TenantID == uuid.Nil || query.Month < 1 || query.Month > 12 || query.Year < 1900 {
		err := domain.ErrInvalidPayCycleConfig
		s.logError("validate pay cycle period", err, serviceTenantIDField(query.TenantID))
		return nil, err
	}
	cycle, err := s.GetPayCycle(ctx, query.TenantID)
	if err != nil {
		return nil, err
	}
	return resolvePayCyclePeriod(cycle, query.Month, query.Year), nil
}

func resolvePayCyclePeriod(cycle *domain.PayCycle, month int, year int) *domain.PayCyclePeriod {
	periodStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, -1)
	attendanceStart := periodStart
	attendanceEnd := periodEnd
	if cycle.AttendancePeriodType == domain.PayrollAttendancePeriodPreviousMonth {
		attendanceStart = periodStart.AddDate(0, -1, 0)
		attendanceEnd = periodStart.AddDate(0, 0, -1)
	} else if cycle.AttendancePeriodType == domain.PayrollAttendancePeriodCustomDays && cycle.StartDay != nil && cycle.EndDay != nil {
		attendanceStart = dateInMonth(year, month, *cycle.StartDay)
		attendanceEnd = dateInMonth(year, month, *cycle.EndDay)
		if attendanceEnd.Before(attendanceStart) {
			attendanceStart = attendanceStart.AddDate(0, -1, 0)
		}
	}
	cutoff := dateInMonth(attendanceEnd.Year(), int(attendanceEnd.Month()), cycle.AttendanceCutoffDay)
	payoutBase := periodStart
	if cycle.PayoutTiming == domain.PayrollPayoutNextMonth {
		payoutBase = periodStart.AddDate(0, 1, 0)
	}
	payDay := int32(30)
	if cycle.PayDay != nil {
		payDay = *cycle.PayDay
	}
	payoutDate := dateInMonth(payoutBase.Year(), int(payoutBase.Month()), payDay).AddDate(0, 0, int(cycle.PayoutOffsetDays))
	var lockDate *time.Time
	if cycle.PayrollLockDay != nil {
		lock := dateInMonth(periodStart.Year(), int(periodStart.Month()), *cycle.PayrollLockDay)
		lockDate = &lock
	}
	return &domain.PayCyclePeriod{PeriodStart: periodStart, PeriodEnd: periodEnd, AttendanceStart: attendanceStart, AttendanceEnd: attendanceEnd, AttendanceCutoff: cutoff, PayoutDate: payoutDate, LockDate: lockDate, Month: month, Year: year}
}

func dateInMonth(year int, month int, day int32) time.Time {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	last := start.AddDate(0, 1, -1).Day()
	cleanDay := int(day)
	if cleanDay < 1 {
		cleanDay = 1
	}
	if cleanDay > last {
		cleanDay = last
	}
	return time.Date(year, time.Month(month), cleanDay, 0, 0, 0, 0, time.UTC)
}
