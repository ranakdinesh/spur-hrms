package services

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) GetSalarySlipFormat(ctx context.Context, tenantID uuid.UUID) (*domain.SalarySlipFormat, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	format, err := s.salarySlips.GetSalarySlipFormat(ctx, tenantID)
	if err != nil {
		return defaultSalarySlipFormat(tenantID), nil
	}
	return format, nil
}

func (s *TenantService) UpsertSalarySlipFormat(ctx context.Context, cmd ports.SalarySlipFormatCommand) (*domain.SalarySlipFormat, error) {
	item, err := domain.NewSalarySlipFormat(domain.SalarySlipFormatInput{TenantID: cmd.TenantID, Title: cmd.Title, Subtitle: cmd.Subtitle, LogoPath: cmd.LogoPath, PrimaryColor: cmd.PrimaryColor, AccentColor: cmd.AccentColor, ShowLeaveBalance: cmd.ShowLeaveBalance, ShowYTDSummary: cmd.ShowYTDSummary, ShowEmployeeBank: cmd.ShowEmployeeBank, ShowEmployerContributions: cmd.ShowEmployerContributions, FooterText: cmd.FooterText, CustomFields: cmd.CustomFields})
	if err != nil {
		s.logError("validate salary slip format", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.salarySlips.UpsertSalarySlipFormat(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert salary slip format", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GenerateSalarySlip(ctx context.Context, cmd ports.GenerateSalarySlipCommand) (*domain.SalarySlip, error) {
	if cmd.TenantID == uuid.Nil || cmd.UserID == uuid.Nil || cmd.FYID == uuid.Nil || cmd.Month < 1 || cmd.Month > 12 || cmd.Year < 1 {
		err := domain.ErrInvalidSalarySlip
		s.logError("validate salary slip generation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if err := s.ensurePayrollPeriodOpen(ctx, cmd.TenantID, cmd.Month, cmd.Year); err != nil {
		s.logError("validate salary slip payroll lock", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	existing, existingErr := s.salarySlips.GetSalarySlipByPeriod(ctx, cmd.TenantID, cmd.UserID, cmd.Month, cmd.Year)
	if existingErr != nil && !isNoRows(existingErr) {
		return nil, existingErr
	}
	if existingErr == nil && !cmd.Regenerate {
		s.logError("validate salary slip regeneration flag", domain.ErrSalarySlipExists, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_slip_id", existing.ID.String()))
		return nil, domain.ErrSalarySlipExists
	}
	calculation, err := s.CalculateEmployeeSalary(ctx, ports.EmployeeSalaryCalculationCommand{TenantID: cmd.TenantID, UserID: cmd.UserID, FYID: cmd.FYID, Month: int(cmd.Month), Year: int(cmd.Year), PresentDays: cmd.PresentDays, AbsentDays: cmd.AbsentDays, TotalDays: cmd.TotalDays, IsSpecial: cmd.IsSpecial})
	if err != nil {
		return nil, err
	}
	salary, err := s.employeeSalaries.GetEmployeeSalary(ctx, cmd.TenantID, calculation.SalaryID)
	if err != nil {
		return nil, err
	}
	items := salaryResultItemsToSlipItems(cmd.TenantID, calculation.SalaryResult.Items)
	totalDeductions := calculation.SalaryResult.TotalDeductions
	netSalary := calculation.SalaryResult.NetSalary
	statItems, err := s.salarySlipStatutoryItems(ctx, cmd.TenantID, cmd.UserID, cmd.Month, cmd.Year, calculation.GrossSalary)
	if err != nil {
		return nil, err
	}
	for _, item := range statItems {
		items = append(items, item)
		totalDeductions += item.Amount
		netSalary -= item.Amount
	}
	slip := &domain.SalarySlip{TenantID: cmd.TenantID, UserID: cmd.UserID, FYID: cmd.FYID, TemplateID: salary.TemplateID, Month: cmd.Month, Year: cmd.Year, GrossSalary: calculation.GrossSalary, TotalEarnings: calculation.SalaryResult.TotalEarnings, TotalDeductions: totalDeductions, AbsentDeduction: calculation.SalaryResult.AbsentDeduction, NetSalary: netSalary, AbsentDays: int32(calculation.AbsentDays), PresentDays: int32(calculation.PresentDays), TotalDays: int32(calculation.TotalDays), LWPDays: float64(calculation.LWPDays), IsSpecial: calculation.IsSpecial, IsRegenerated: existingErr == nil && cmd.Regenerate}
	leaves, err := s.salarySlipLeavesSnapshot(ctx, cmd.TenantID, cmd.UserID, cmd.FYID)
	if err != nil {
		return nil, err
	}
	var result *domain.SalarySlip
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var saved *domain.SalarySlip
		var err error
		if existingErr == nil && cmd.Regenerate {
			slip.ID = existing.ID
			saved, err = s.salarySlips.UpdateSalarySlip(txCtx, slip, cmd.ActorID)
			if err != nil {
				return err
			}
			if err := s.salarySlips.DeleteSalarySlipItemsBySlip(txCtx, cmd.TenantID, saved.ID, cmd.ActorID); err != nil {
				return err
			}
			if err := s.salarySlips.DeleteSalarySlipLeavesBySlip(txCtx, cmd.TenantID, saved.ID, cmd.ActorID); err != nil {
				return err
			}
		} else {
			saved, err = s.salarySlips.CreateSalarySlip(txCtx, slip, cmd.ActorID)
			if err != nil {
				return err
			}
		}
		for _, item := range items {
			item.SlipID = saved.ID
			if _, err := s.salarySlips.CreateSalarySlipItem(txCtx, item, cmd.ActorID); err != nil {
				return err
			}
		}
		for _, leave := range leaves {
			leave.SlipID = saved.ID
			if _, err := s.salarySlips.CreateSalarySlipLeave(txCtx, leave, cmd.ActorID); err != nil {
				return err
			}
		}
		saved.Items = items
		saved.Leaves = leaves
		result = saved
		return nil
	})
	if err != nil {
		s.logError("generate salary slip", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if _, err := s.ensureSalarySlipPDF(ctx, result, cmd.ActorID); err != nil {
		s.logError("store generated salary slip pdf", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_slip_id", result.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) BulkGenerateSalarySlips(ctx context.Context, cmd ports.BulkGenerateSalarySlipsCommand) ([]*domain.SalarySlip, error) {
	employees, err := s.employees.ListEmployees(ctx, cmd.TenantID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.SalarySlip, 0, len(employees))
	for _, employee := range employees {
		if employee.Inactive {
			continue
		}
		slip, err := s.GenerateSalarySlip(ctx, ports.GenerateSalarySlipCommand{TenantID: cmd.TenantID, UserID: employee.UserID, FYID: cmd.FYID, Month: cmd.Month, Year: cmd.Year, Regenerate: cmd.Regenerate, ActorID: cmd.ActorID})
		if err != nil {
			s.log.Warn().Err(err).Str("tenant_id", cmd.TenantID.String()).Str("user_id", employee.UserID.String()).Msg("hrms: skipped salary slip generation")
			continue
		}
		result = append(result, slip)
	}
	return result, nil
}

func (s *TenantService) ListSalarySlipsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.SalarySlip, error) {
	items, err := s.salarySlips.ListSalarySlipsByUser(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	return s.hydrateSalarySlips(ctx, tenantID, items)
}

func (s *TenantService) ListSalarySlipsByTenantPeriod(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]*domain.SalarySlip, error) {
	items, err := s.salarySlips.ListSalarySlipsByTenantPeriod(ctx, tenantID, month, year)
	if err != nil {
		return nil, err
	}
	return s.hydrateSalarySlips(ctx, tenantID, items)
}

func (s *TenantService) GetSalarySlip(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SalarySlip, error) {
	item, err := s.salarySlips.GetSalarySlip(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if err := s.hydrateSalarySlip(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *TenantService) RenderSalarySlipPDF(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) ([]byte, string, error) {
	slip, err := s.GetSalarySlip(ctx, tenantID, id)
	if err != nil {
		return nil, "", err
	}
	data, err := s.ensureSalarySlipPDF(ctx, slip, nil)
	if err != nil {
		return nil, "", err
	}
	return data, salarySlipPDFName(slip), nil
}

func (s *TenantService) RenderRecentSalarySlipsZip(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, months int32) ([]byte, string, error) {
	if months <= 0 || months > 7 {
		months = 6
	}
	items, err := s.salarySlips.ListRecentSalarySlipsByUser(ctx, tenantID, userID, months)
	if err != nil {
		return nil, "", err
	}
	items, err = s.hydrateSalarySlips(ctx, tenantID, items)
	if err != nil {
		return nil, "", err
	}
	data, err := s.renderSlipsZip(ctx, items)
	if err != nil {
		return nil, "", err
	}
	return data, fmt.Sprintf("salary-slips-last-%d-months.zip", months), nil
}

func (s *TenantService) RenderTenantSalarySlipsZip(ctx context.Context, tenantID uuid.UUID, month int32, year int32) ([]byte, string, error) {
	items, err := s.ListSalarySlipsByTenantPeriod(ctx, tenantID, month, year)
	if err != nil {
		return nil, "", err
	}
	data, err := s.renderSlipsZip(ctx, items)
	if err != nil {
		return nil, "", err
	}
	return data, fmt.Sprintf("salary-slips-%04d-%02d.zip", year, month), nil
}

func (s *TenantService) hydrateSalarySlips(ctx context.Context, tenantID uuid.UUID, items []*domain.SalarySlip) ([]*domain.SalarySlip, error) {
	for _, item := range items {
		if err := s.hydrateSalarySlip(ctx, item); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *TenantService) hydrateSalarySlip(ctx context.Context, item *domain.SalarySlip) error {
	items, err := s.salarySlips.ListSalarySlipItems(ctx, item.TenantID, item.ID)
	if err != nil {
		return err
	}
	leaves, err := s.salarySlips.ListSalarySlipLeaves(ctx, item.TenantID, item.ID)
	if err != nil {
		return err
	}
	item.Items = items
	item.Leaves = leaves
	return nil
}

func (s *TenantService) renderSlipPDF(ctx context.Context, slip *domain.SalarySlip) ([]byte, error) {
	if s.salarySlipPDF == nil {
		return nil, domain.ErrSalarySlipFormatMissing
	}
	format, err := s.GetSalarySlipFormat(ctx, slip.TenantID)
	if err != nil {
		return nil, err
	}
	employee, _ := s.employees.GetEmployeeByUserID(ctx, slip.TenantID, slip.UserID)
	return s.salarySlipPDF.RenderSalarySlipPDF(ctx, ports.SalarySlipDocument{Format: format, Slip: slip, Employee: employee})
}

func (s *TenantService) ensureSalarySlipPDF(ctx context.Context, slip *domain.SalarySlip, actorID *uuid.UUID) ([]byte, error) {
	data, err := s.renderSlipPDF(ctx, slip)
	if err != nil {
		return nil, err
	}
	if s.salarySlipStorage == nil {
		return nil, domain.ErrSalarySlipFormatMissing
	}
	path, err := s.salarySlipStorage.StoreSalarySlipPDF(ctx, ports.StoreSalarySlipPDFInput{TenantID: slip.TenantID, UserID: slip.UserID, SlipID: slip.ID, Month: slip.Month, Year: slip.Year, FileName: salarySlipPDFName(slip), ContentType: "application/pdf", Content: data})
	if err != nil {
		return nil, err
	}
	if slip.PDFPath == nil || *slip.PDFPath != path {
		slip.PDFPath = &path
		updated, err := s.salarySlips.UpdateSalarySlip(ctx, slip, actorID)
		if err != nil {
			return nil, err
		}
		updated.Items = slip.Items
		updated.Leaves = slip.Leaves
		*slip = *updated
	}
	return data, nil
}

func (s *TenantService) renderSlipsZip(ctx context.Context, slips []*domain.SalarySlip) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, slip := range slips {
		data, err := s.ensureSalarySlipPDF(ctx, slip, nil)
		if err != nil {
			return nil, err
		}
		w, err := zw.Create(salarySlipPDFName(slip))
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(data); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func salaryResultItemsToSlipItems(tenantID uuid.UUID, items []domain.SalaryItem) []*domain.SalarySlipItem {
	result := make([]*domain.SalarySlipItem, 0, len(items))
	for _, item := range items {
		result = append(result, &domain.SalarySlipItem{TenantID: tenantID, ItemType: item.ItemType, Code: item.Code, Name: item.Name, Amount: item.Amount, SortOrder: int32(item.SortOrder)})
	}
	return result
}

func (s *TenantService) salarySlipStatutoryItems(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, month int32, year int32, grossSalary float64) ([]*domain.SalarySlipItem, error) {
	employee, err := s.employees.GetEmployeeByUserID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	profile, err := s.employees.GetEmployeeProfile(ctx, tenantID, employee.ID)
	if err != nil {
		return nil, err
	}
	if profile.Statutory == nil {
		return nil, nil
	}
	state := employee.State
	branchID := employee.BranchID
	effective := time.Date(int(year), time.Month(month), daysInMonth(int(year), time.Month(month)), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	items := make([]*domain.SalarySlipItem, 0, 2)
	if profile.Statutory.PTApplicable {
		if rule, err := s.payrollOperations.ResolvePayrollStatutoryRule(ctx, tenantID, domain.PayrollRulePT, state, branchID, effective, grossSalary, month); err == nil && rule.EmployeeAmount > 0 {
			items = append(items, &domain.SalarySlipItem{TenantID: tenantID, ItemType: domain.SalaryItemDeduction, Code: "PT", Name: "Professional Tax", Amount: rule.EmployeeAmount, SortOrder: 980})
		} else if err != nil && !errors.Is(err, domain.ErrPayrollStatutoryRuleNotFound) {
			return nil, err
		}
	}
	if profile.Statutory.LWFApplicable {
		if rule, err := s.payrollOperations.ResolvePayrollStatutoryRule(ctx, tenantID, domain.PayrollRuleLWF, state, branchID, effective, grossSalary, month); err == nil && rule.EmployeeAmount > 0 {
			items = append(items, &domain.SalarySlipItem{TenantID: tenantID, ItemType: domain.SalaryItemDeduction, Code: "LWF", Name: "Labour Welfare Fund", Amount: rule.EmployeeAmount, SortOrder: 981})
		} else if err != nil && !errors.Is(err, domain.ErrPayrollStatutoryRuleNotFound) {
			return nil, err
		}
	}
	return items, nil
}

func (s *TenantService) salarySlipLeavesSnapshot(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, fyID uuid.UUID) ([]*domain.SalarySlipLeave, error) {
	balances, err := s.leaveBalances.ListLeaveBalancesByUser(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	items := make([]*domain.SalarySlipLeave, 0, len(balances))
	for _, balance := range balances {
		if balance.FYID != fyID {
			continue
		}
		items = append(items, &domain.SalarySlipLeave{TenantID: tenantID, LeaveTypeID: balance.LeaveTypeID, TotalDays: balance.TotalDays, UsedDays: balance.UsedDays, BalanceDays: balance.BalanceDays})
	}
	return items, nil
}

func defaultSalarySlipFormat(tenantID uuid.UUID) *domain.SalarySlipFormat {
	item, _ := domain.NewSalarySlipFormat(domain.SalarySlipFormatInput{TenantID: tenantID, Title: "Salary Slip", PrimaryColor: "#111827", AccentColor: "#588368", ShowLeaveBalance: true, ShowEmployeeBank: true})
	return item
}

func salarySlipPDFName(slip *domain.SalarySlip) string {
	return fmt.Sprintf("salary-slip-%s-%04d-%02d.pdf", slip.UserID.String(), slip.Year, slip.Month)
}

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
