package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateFinancialYear(ctx context.Context, cmd ports.FinancialYearCommand) (*domain.FinancialYear, error) {
	startDate, endDate, err := parseFinancialYearPeriod(cmd.StartDate, cmd.EndDate)
	if err != nil {
		s.logError("validate financial year period", err, serviceTenantIDField(cmd.TenantID), serviceStringField("start_date", cmd.StartDate), serviceStringField("end_date", cmd.EndDate))
		return nil, err
	}
	ensureFinancialYearUsageFlags(&cmd)
	item, err := domain.NewFinancialYear(domain.FinancialYearInput{
		TenantID:      cmd.TenantID,
		Name:          cmd.Name,
		StartDate:     startDate,
		EndDate:       endDate,
		PayrollYear:   cmd.PayrollYear,
		LeaveYear:     cmd.LeaveYear,
		HolidayYear:   cmd.HolidayYear,
		ReportingYear: cmd.ReportingYear,
		CloseNote:     cmd.CloseNote,
	})
	if err != nil {
		s.logError("validate financial year create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_name", cmd.Name))
		return nil, err
	}
	result, err := s.financialYears.CreateFinancialYear(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create financial year", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_name", item.Name))
		return nil, err
	}
	if cmd.IsActive {
		result, err = s.SetActiveFinancialYear(ctx, cmd.TenantID, result.ID, cmd.ActorID)
		if err != nil {
			return nil, err
		}
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("financial_year_id", result.ID.String()).Msg("hrms: financial year created")
	return result, nil
}

func (s *TenantService) ListFinancialYears(ctx context.Context, tenantID uuid.UUID) ([]*domain.FinancialYear, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate financial year list tenant", err)
		return nil, err
	}
	result, err := s.financialYears.ListFinancialYears(ctx, tenantID)
	if err != nil {
		s.logError("list financial years", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetFinancialYear(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.FinancialYear, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate financial year get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidFinancialYearID
		s.logError("validate financial year get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.financialYears.GetFinancialYear(ctx, tenantID, id)
	if err != nil {
		s.logError("get financial year", err, serviceTenantIDField(tenantID), serviceStringField("financial_year_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetActiveFinancialYear(ctx context.Context, tenantID uuid.UUID) (*domain.FinancialYear, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate active financial year tenant", err)
		return nil, err
	}
	result, err := s.financialYears.GetActiveFinancialYear(ctx, tenantID)
	if err != nil {
		s.logError("get active financial year", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateFinancialYear(ctx context.Context, cmd ports.FinancialYearCommand) (*domain.FinancialYear, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidFinancialYearID
		s.logError("validate financial year update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetFinancialYear(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	if existing.IsLocked {
		err := domain.ErrFinancialYearLocked
		s.logError("validate financial year update lock", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.ID.String()))
		return nil, err
	}
	startDate, endDate, err := parseFinancialYearPeriod(cmd.StartDate, cmd.EndDate)
	if err != nil {
		s.logError("validate financial year update period", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.ID.String()))
		return nil, err
	}
	ensureFinancialYearUsageFlags(&cmd)
	item, err := domain.NewFinancialYear(domain.FinancialYearInput{
		TenantID:      cmd.TenantID,
		Name:          cmd.Name,
		StartDate:     startDate,
		EndDate:       endDate,
		PayrollYear:   cmd.PayrollYear,
		LeaveYear:     cmd.LeaveYear,
		HolidayYear:   cmd.HolidayYear,
		ReportingYear: cmd.ReportingYear,
		CloseNote:     cmd.CloseNote,
	})
	if err != nil {
		s.logError("validate financial year update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.financialYears.UpdateFinancialYear(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update financial year", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.ID.String()))
		return nil, err
	}
	if cmd.IsActive && !result.IsActive {
		result, err = s.SetActiveFinancialYear(ctx, cmd.TenantID, result.ID, cmd.ActorID)
		if err != nil {
			return nil, err
		}
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("financial_year_id", result.ID.String()).Msg("hrms: financial year updated")
	return result, nil
}

func (s *TenantService) DeleteFinancialYear(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate financial year delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidFinancialYearID
		s.logError("validate financial year delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	existing, err := s.GetFinancialYear(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.IsLocked {
		err := domain.ErrFinancialYearLocked
		s.logError("validate financial year delete lock", err, serviceTenantIDField(tenantID), serviceStringField("financial_year_id", id.String()))
		return err
	}
	if err := s.financialYears.DeleteFinancialYear(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete financial year", err, serviceTenantIDField(tenantID), serviceStringField("financial_year_id", id.String()))
		return err
	}
	s.log.Info().Str("tenant_id", tenantID.String()).Str("financial_year_id", id.String()).Msg("hrms: financial year deactivated")
	return nil
}

func (s *TenantService) SetActiveFinancialYear(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.FinancialYear, error) {
	existing, err := s.GetFinancialYear(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if existing.IsLocked {
		err := domain.ErrFinancialYearLocked
		s.logError("validate financial year active lock", err, serviceTenantIDField(tenantID), serviceStringField("financial_year_id", id.String()))
		return nil, err
	}
	result, err := s.financialYears.SetActiveFinancialYear(ctx, tenantID, id, actorID)
	if err != nil {
		s.logError("set active financial year", err, serviceTenantIDField(tenantID), serviceStringField("financial_year_id", id.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", tenantID.String()).Str("financial_year_id", id.String()).Msg("hrms: active financial year set")
	return result, nil
}

func (s *TenantService) SetFinancialYearLock(ctx context.Context, cmd ports.FinancialYearLockCommand) (*domain.FinancialYear, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate financial year lock tenant", err)
		return nil, err
	}
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidFinancialYearID
		s.logError("validate financial year lock id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.financialYears.SetFinancialYearLock(ctx, cmd.TenantID, cmd.ID, cmd.IsLocked, cmd.CloseNote, cmd.ActorID)
	if err != nil {
		s.logError("set financial year lock", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", cmd.TenantID.String()).Str("financial_year_id", cmd.ID.String()).Bool("is_locked", result.IsLocked).Msg("hrms: financial year lock updated")
	return result, nil
}

func parseFinancialYearPeriod(startDate string, endDate string) (time.Time, time.Time, error) {
	start, err := time.Parse("2006-01-02", strings.TrimSpace(startDate))
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrInvalidFinancialYearPeriod
	}
	end, err := time.Parse("2006-01-02", strings.TrimSpace(endDate))
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrInvalidFinancialYearPeriod
	}
	return start, end, nil
}

func ensureFinancialYearUsageFlags(cmd *ports.FinancialYearCommand) {
	if cmd == nil {
		return
	}
	if !cmd.PayrollYear && !cmd.LeaveYear && !cmd.HolidayYear && !cmd.ReportingYear {
		cmd.PayrollYear = true
		cmd.LeaveYear = true
		cmd.HolidayYear = true
		cmd.ReportingYear = true
	}
}
