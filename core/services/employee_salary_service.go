package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) AssignEmployeeSalary(ctx context.Context, cmd ports.EmployeeSalaryCommand) (*domain.EmployeeSalary, error) {
	item, err := employeeSalaryFromCommand(cmd)
	if err != nil {
		s.logError("validate employee salary assignment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	template, err := s.salaryTemplates.GetSalaryTemplate(ctx, cmd.TenantID, cmd.TemplateID)
	if err != nil {
		s.logError("get salary template for employee salary", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_template_id", cmd.TemplateID.String()))
		return nil, err
	}
	if template.FYID != cmd.FYID {
		err := domain.ErrInvalidSalaryTemplate
		s.logError("validate salary template fy for employee salary", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_template_id", cmd.TemplateID.String()), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	structures, err := calculateEmployeeSalaryStructures(item, template.Items)
	if err != nil {
		s.logError("calculate employee salary structure", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_template_id", cmd.TemplateID.String()))
		return nil, err
	}
	var result *domain.EmployeeSalary
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		if err := s.employeeSalaries.DeleteEmployeeSalariesByUserFY(txCtx, item.TenantID, item.UserID, item.FYID, cmd.ActorID); err != nil {
			return err
		}
		if err := s.employeeSalaries.DeleteEmployeeSalaryStructuresByUserFY(txCtx, item.TenantID, item.UserID, item.FYID, cmd.ActorID); err != nil {
			return err
		}
		created, err := s.employeeSalaries.CreateEmployeeSalary(txCtx, item, cmd.ActorID)
		if err != nil {
			return err
		}
		for _, structure := range structures {
			if _, err := s.employeeSalaries.CreateEmployeeSalaryStructure(txCtx, structure, cmd.ActorID); err != nil {
				return err
			}
		}
		created.Template = template
		created.Structures = structures
		result = created
		return nil
	})
	if err != nil {
		s.logError("assign employee salary", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("user_id", result.UserID.String()).Str("employee_salary_id", result.ID.String()).Msg("hrms: employee salary assigned")
	return result, nil
}

func (s *TenantService) ListEmployeeSalariesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.EmployeeSalary, error) {
	if tenantID == uuid.Nil || userID == uuid.Nil {
		err := domain.ErrInvalidEmployeeUserID
		s.logError("validate employee salary list", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.employeeSalaries.ListEmployeeSalariesByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("list employee salaries", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	for _, item := range items {
		structures, err := s.employeeSalaries.ListEmployeeSalaryStructures(ctx, tenantID, item.UserID, item.FYID)
		if err != nil {
			return nil, err
		}
		item.Structures = structures
	}
	return items, nil
}

func (s *TenantService) GetEmployeeSalary(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeSalary, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidEmployeeSalaryID
		s.logError("validate employee salary get", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.employeeSalaries.GetEmployeeSalary(ctx, tenantID, id)
	if err != nil {
		s.logError("get employee salary", err, serviceTenantIDField(tenantID), serviceStringField("employee_salary_id", id.String()))
		return nil, err
	}
	structures, err := s.employeeSalaries.ListEmployeeSalaryStructures(ctx, tenantID, item.UserID, item.FYID)
	if err != nil {
		return nil, err
	}
	item.Structures = structures
	return item, nil
}

func (s *TenantService) DeleteEmployeeSalary(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidEmployeeSalaryID
		s.logError("validate employee salary delete", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.employeeSalaries.DeleteEmployeeSalary(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete employee salary", err, serviceTenantIDField(tenantID), serviceStringField("employee_salary_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CalculateEmployeeSalary(ctx context.Context, cmd ports.EmployeeSalaryCalculationCommand) (*ports.EmployeeSalaryCalculation, error) {
	if cmd.TenantID == uuid.Nil || cmd.UserID == uuid.Nil || cmd.FYID == uuid.Nil || cmd.Month < 1 || cmd.Month > 12 || cmd.Year < 1 {
		err := domain.ErrInvalidEmployeeSalary
		s.logError("validate employee salary calculation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	salary, err := s.employeeSalaryForCalculation(ctx, cmd)
	if err != nil {
		s.logError("load employee salary for calculation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	structures, err := s.employeeSalaries.ListEmployeeSalaryStructures(ctx, cmd.TenantID, salary.UserID, salary.FYID)
	if err != nil {
		s.logError("list employee salary structures for calculation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", salary.UserID.String()), serviceStringField("financial_year_id", salary.FYID.String()))
		return nil, err
	}
	totalDays := daysInMonth(cmd.Year, time.Month(cmd.Month))
	if cmd.TotalDays != nil && *cmd.TotalDays > 0 {
		totalDays = *cmd.TotalDays
	}
	absentDays := 0
	if cmd.AbsentDays != nil && *cmd.AbsentDays >= 0 {
		absentDays = *cmd.AbsentDays
	} else {
		absentDays, err = s.countAbsentDays(ctx, cmd.TenantID, cmd.UserID, cmd.Year, cmd.Month)
		if err != nil {
			return nil, err
		}
	}
	presentDays := totalDays - absentDays
	if cmd.PresentDays != nil && *cmd.PresentDays >= 0 {
		presentDays = *cmd.PresentDays
	}
	if presentDays < 0 {
		presentDays = 0
	}
	if absentDays < 0 {
		absentDays = 0
	}
	items := make([]domain.SalaryItem, 0, len(structures))
	for _, item := range structures {
		items = append(items, domain.SalaryItem{ItemType: item.ItemType, Code: item.Code, Name: item.Name, Amount: item.Amount, SortOrder: int(item.SortOrder)})
	}
	result := domain.CalculateUserSalary(salary.GrossSalary, items, presentDays, absentDays, totalDays, cmd.IsSpecial)
	return &ports.EmployeeSalaryCalculation{TenantID: cmd.TenantID, UserID: salary.UserID, FYID: salary.FYID, SalaryID: salary.ID, Month: cmd.Month, Year: cmd.Year, PresentDays: presentDays, AbsentDays: absentDays, TotalDays: totalDays, LWPDays: absentDays, IsSpecial: cmd.IsSpecial, GrossSalary: salary.GrossSalary, SalaryResult: result}, nil
}

func employeeSalaryFromCommand(cmd ports.EmployeeSalaryCommand) (*domain.EmployeeSalary, error) {
	effectiveFrom, err := parseOptionalEmployeeSalaryDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	return domain.NewEmployeeSalary(domain.EmployeeSalaryInput{TenantID: cmd.TenantID, UserID: cmd.UserID, FYID: cmd.FYID, TemplateID: cmd.TemplateID, GrossSalary: cmd.GrossSalary, EffectiveFrom: effectiveFrom})
}

func (s *TenantService) employeeSalaryForCalculation(ctx context.Context, cmd ports.EmployeeSalaryCalculationCommand) (*domain.EmployeeSalary, error) {
	if cmd.SalaryID != nil && *cmd.SalaryID != uuid.Nil {
		return s.employeeSalaries.GetEmployeeSalary(ctx, cmd.TenantID, *cmd.SalaryID)
	}
	items, err := s.employeeSalaries.ListEmployeeSalariesByUser(ctx, cmd.TenantID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.FYID == cmd.FYID {
			return item, nil
		}
	}
	return nil, domain.ErrInvalidEmployeeSalaryID
}

func (s *TenantService) countAbsentDays(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, year int, month int) (int, error) {
	totalDays := daysInMonth(year, time.Month(month))
	absentDays := 0
	for day := 1; day <= totalDays; day++ {
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		items, err := s.ListAttendanceDailyStatuses(ctx, ports.AttendanceStatusQuery{TenantID: tenantID, UserID: &userID, Date: date})
		if err != nil {
			s.logError("calculate salary attendance status", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()), serviceStringField("date", date))
			return 0, err
		}
		for _, item := range items {
			if item.UserID == userID && item.Status == domain.AttendanceStatusAbsent {
				absentDays++
				break
			}
		}
	}
	return absentDays, nil
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func parseOptionalEmployeeSalaryDate(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*value))
	if err != nil {
		return nil, domain.ErrInvalidEmployeeSalary
	}
	return &parsed, nil
}

func calculateEmployeeSalaryStructures(salary *domain.EmployeeSalary, items []*domain.SalaryTemplateItem) ([]*domain.EmployeeSalaryStructure, error) {
	structures := make([]*domain.EmployeeSalaryStructure, 0, len(items))
	amountsByCode := map[string]float64{"ctc": salary.GrossSalary, "gross": salary.GrossSalary}
	for _, item := range items {
		if item == nil || item.ItemType == domain.SalaryItemEmployerContribution || item.ItemType == domain.SalaryItemReimbursement {
			continue
		}
		amount := calculateTemplateItemAmount(item, salary.GrossSalary, amountsByCode)
		structure, err := domain.NewEmployeeSalaryStructure(salary.TenantID, salary.UserID, salary.FYID, salary.TemplateID, item, amount)
		if err != nil {
			return nil, err
		}
		amountsByCode[item.Code] = structure.Amount
		if item.Code == domain.SalaryCodeBasic {
			amountsByCode["basic"] = structure.Amount
		}
		structures = append(structures, structure)
	}
	return structures, nil
}

func calculateTemplateItemAmount(item *domain.SalaryTemplateItem, grossSalary float64, amountsByCode map[string]float64) float64 {
	var amount float64
	switch item.CalculationMode {
	case domain.SalaryCalculationFixed:
		if item.Amount != nil {
			amount = *item.Amount
		}
	case domain.SalaryCalculationPercentage:
		base := grossSalary
		if value, ok := amountsByCode[item.CalculationBase]; ok {
			base = value
		}
		if item.Percentage != nil {
			amount = base * (*item.Percentage) / 100
		}
	case domain.SalaryCalculationManual:
		if item.Amount != nil {
			amount = *item.Amount
		}
	}
	if item.CapAmount != nil && amount > *item.CapAmount {
		amount = *item.CapAmount
	}
	if item.MinAmount != nil && amount < *item.MinAmount {
		amount = *item.MinAmount
	}
	if item.MaxAmount != nil && amount > *item.MaxAmount {
		amount = *item.MaxAmount
	}
	if amount < 0 {
		return 0
	}
	return amount
}
