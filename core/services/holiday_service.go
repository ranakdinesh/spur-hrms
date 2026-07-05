package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateHoliday(ctx context.Context, cmd ports.HolidayCommand) (*domain.Holiday, error) {
	item, err := s.buildHoliday(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := s.holidays.CreateHoliday(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create holiday", err, serviceTenantIDField(cmd.TenantID), serviceStringField("holiday_name", item.Name))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("holiday_id", result.ID.String()).Msg("hrms: holiday created")
	return result, nil
}

func (s *TenantService) ListHolidays(ctx context.Context, cmd ports.ListHolidaysCommand) ([]*domain.Holiday, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate holiday list tenant", err)
		return nil, err
	}
	if cmd.Upcoming {
		limit := cmd.Limit
		if limit <= 0 || limit > 100 {
			limit = 10
		}
		result, err := s.holidays.ListUpcomingHolidays(ctx, cmd.TenantID, limit)
		if err != nil {
			s.logError("list upcoming holidays", err, serviceTenantIDField(cmd.TenantID))
			return nil, err
		}
		return result, nil
	}
	if cmd.FYID != nil && *cmd.FYID != uuid.Nil {
		result, err := s.holidays.ListHolidaysByFinancialYear(ctx, cmd.TenantID, *cmd.FYID)
		if err != nil {
			s.logError("list holidays by financial year", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
			return nil, err
		}
		return result, nil
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		result, err := s.holidays.ListHolidaysByBranch(ctx, cmd.TenantID, *cmd.BranchID)
		if err != nil {
			s.logError("list holidays by branch", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", cmd.BranchID.String()))
			return nil, err
		}
		return result, nil
	}
	if strings.TrimSpace(cmd.StartDate) != "" || strings.TrimSpace(cmd.EndDate) != "" {
		startDate, endDate, err := parseDateRange(cmd.StartDate, cmd.EndDate)
		if err != nil {
			s.logError("validate holiday date range", err, serviceTenantIDField(cmd.TenantID), serviceStringField("start_date", cmd.StartDate), serviceStringField("end_date", cmd.EndDate))
			return nil, err
		}
		result, err := s.holidays.ListHolidaysByDateRange(ctx, cmd.TenantID, startDate, endDate)
		if err != nil {
			s.logError("list holidays by date range", err, serviceTenantIDField(cmd.TenantID))
			return nil, err
		}
		return result, nil
	}
	result, err := s.holidays.ListHolidays(ctx, cmd.TenantID)
	if err != nil {
		s.logError("list holidays", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetHoliday(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Holiday, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate holiday get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidHolidayID
		s.logError("validate holiday get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.holidays.GetHoliday(ctx, tenantID, id)
	if err != nil {
		s.logError("get holiday", err, serviceTenantIDField(tenantID), serviceStringField("holiday_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateHoliday(ctx context.Context, cmd ports.HolidayCommand) (*domain.Holiday, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidHolidayID
		s.logError("validate holiday update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetHoliday(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := s.buildHoliday(ctx, cmd)
	if err != nil {
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.holidays.UpdateHoliday(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update holiday", err, serviceTenantIDField(cmd.TenantID), serviceStringField("holiday_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("holiday_id", result.ID.String()).Msg("hrms: holiday updated")
	return result, nil
}

func (s *TenantService) DeleteHoliday(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate holiday delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidHolidayID
		s.logError("validate holiday delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if _, err := s.GetHoliday(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.holidays.DeleteHoliday(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete holiday", err, serviceTenantIDField(tenantID), serviceStringField("holiday_id", id.String()))
		return err
	}
	s.log.Info().Str("tenant_id", tenantID.String()).Str("holiday_id", id.String()).Msg("hrms: holiday deactivated")
	return nil
}

func (s *TenantService) buildHoliday(ctx context.Context, cmd ports.HolidayCommand) (*domain.Holiday, error) {
	date, err := parseHolidayDate(cmd.Date)
	if err != nil {
		s.logError("validate holiday date", err, serviceTenantIDField(cmd.TenantID), serviceStringField("holiday_date", cmd.Date))
		return nil, err
	}
	if cmd.BranchID != nil && *cmd.BranchID != uuid.Nil {
		if _, err := s.branches.GetBranch(ctx, cmd.TenantID, *cmd.BranchID); err != nil {
			s.logError("validate holiday branch", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", cmd.BranchID.String()))
			return nil, err
		}
	}
	if cmd.FYID != nil && *cmd.FYID != uuid.Nil {
		fy, err := s.financialYears.GetFinancialYear(ctx, cmd.TenantID, *cmd.FYID)
		if err != nil {
			s.logError("validate holiday financial year", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
			return nil, err
		}
		if date.Before(fy.StartDate) || date.After(fy.EndDate) {
			err := domain.ErrInvalidHolidayDate
			s.logError("validate holiday date in financial year", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", fy.ID.String()))
			return nil, err
		}
	}
	item, err := domain.NewHoliday(domain.HolidayInput{TenantID: cmd.TenantID, BranchID: cmd.BranchID, FYID: cmd.FYID, Name: cmd.Name, Date: date, IsOptional: cmd.IsOptional})
	if err != nil {
		s.logError("validate holiday", err, serviceTenantIDField(cmd.TenantID), serviceStringField("holiday_name", cmd.Name))
		return nil, err
	}
	return item, nil
}

func parseHolidayDate(value string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, domain.ErrInvalidHolidayDate
	}
	return date, nil
}

func parseDateRange(startDate string, endDate string) (time.Time, time.Time, error) {
	start, err := parseHolidayDate(startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := parseHolidayDate(endDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, domain.ErrInvalidHolidayDate
	}
	return start, end, nil
}
