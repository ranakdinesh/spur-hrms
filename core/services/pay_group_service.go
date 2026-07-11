package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreatePayGroup(ctx context.Context, cmd ports.PayGroupCommand) (*domain.PayGroup, error) {
	if cmd.ID == uuid.Nil && !cmd.IsActive {
		cmd.IsActive = true
	}
	item, err := domain.NewPayGroup(payGroupFromCommand(cmd))
	if err != nil {
		s.logError("validate pay group", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.payGroups.CreatePayGroup(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create pay group", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", cmd.Code))
		return nil, err
	}
	return s.hydratePayGroup(ctx, result)
}

func (s *TenantService) UpdatePayGroup(ctx context.Context, cmd ports.PayGroupCommand) (*domain.PayGroup, error) {
	item, err := domain.NewPayGroup(payGroupFromCommand(cmd))
	if err != nil || cmd.ID == uuid.Nil {
		if err == nil {
			err = domain.ErrInvalidPayGroup
		}
		s.logError("validate pay group update", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.payGroups.UpdatePayGroup(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update pay group", err, serviceTenantIDField(cmd.TenantID), serviceStringField("pay_group_id", cmd.ID.String()))
		return nil, err
	}
	return s.hydratePayGroup(ctx, result)
}

func (s *TenantService) ListPayGroups(ctx context.Context, tenantID uuid.UUID) ([]*domain.PayGroup, error) {
	items, err := s.payGroups.ListPayGroups(ctx, tenantID)
	if err != nil {
		s.logError("list pay groups", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	for _, item := range items {
		if _, err := s.hydratePayGroup(ctx, item); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *TenantService) GetPayGroup(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayGroup, error) {
	item, err := s.payGroups.GetPayGroup(ctx, tenantID, id)
	if err != nil {
		s.logError("get pay group", err, serviceTenantIDField(tenantID), serviceStringField("pay_group_id", id.String()))
		return nil, err
	}
	return s.hydratePayGroup(ctx, item)
}

func (s *TenantService) DeletePayGroup(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.payGroups.DeletePayGroup(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete pay group", err, serviceTenantIDField(tenantID), serviceStringField("pay_group_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) UpsertPayGroupMember(ctx context.Context, cmd ports.PayGroupMemberCommand) (*domain.PayGroupMember, error) {
	from, err := parseOptionalDate(cmd.EffectiveFrom)
	if err != nil {
		s.logError("validate pay group member effective_from", err, serviceTenantIDField(cmd.TenantID))
		return nil, domain.ErrInvalidPayGroupMember
	}
	to, err := parseOptionalDate(cmd.EffectiveTo)
	if err != nil {
		s.logError("validate pay group member effective_to", err, serviceTenantIDField(cmd.TenantID))
		return nil, domain.ErrInvalidPayGroupMember
	}
	item, err := domain.NewPayGroupMember(domain.PayGroupMember{
		TenantID:       cmd.TenantID,
		PayGroupID:     cmd.PayGroupID,
		UserID:         cmd.UserID,
		MembershipType: cmd.MembershipType,
		EffectiveFrom:  from,
		EffectiveTo:    to,
	})
	if err != nil {
		s.logError("validate pay group member", err, serviceTenantIDField(cmd.TenantID), serviceStringField("pay_group_id", cmd.PayGroupID.String()))
		return nil, err
	}
	result, err := s.payGroups.UpsertPayGroupMember(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert pay group member", err, serviceTenantIDField(cmd.TenantID), serviceStringField("pay_group_id", cmd.PayGroupID.String()), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeletePayGroupMember(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.payGroups.DeletePayGroupMember(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete pay group member", err, serviceTenantIDField(tenantID), serviceStringField("member_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ListPayGroupEmployees(ctx context.Context, tenantID uuid.UUID, payGroupID uuid.UUID) ([]*domain.PayGroupEmployee, error) {
	items, err := s.payGroups.ListPayGroupEmployees(ctx, tenantID, payGroupID)
	if err != nil {
		s.logError("list pay group employees", err, serviceTenantIDField(tenantID), serviceStringField("pay_group_id", payGroupID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) CreatePayRun(ctx context.Context, cmd ports.PayRunCommand) (*domain.PayRun, error) {
	if err := s.ensurePayrollPeriodOpen(ctx, cmd.TenantID, cmd.Month, cmd.Year); err != nil {
		s.logError("validate pay run payroll lock", err, serviceTenantIDField(cmd.TenantID), serviceStringField("pay_group_id", cmd.PayGroupID.String()))
		return nil, err
	}
	if _, err := s.payGroups.GetPayGroup(ctx, cmd.TenantID, cmd.PayGroupID); err != nil {
		s.logError("validate pay run pay group", err, serviceTenantIDField(cmd.TenantID), serviceStringField("pay_group_id", cmd.PayGroupID.String()))
		return nil, err
	}
	employees, err := s.payGroups.ListPayGroupEmployees(ctx, cmd.TenantID, cmd.PayGroupID)
	if err != nil {
		return nil, err
	}
	if len(employees) == 0 {
		s.logError("validate pay run employees", domain.ErrInvalidPayRun, serviceTenantIDField(cmd.TenantID), serviceStringField("pay_group_id", cmd.PayGroupID.String()))
		return nil, domain.ErrInvalidPayRun
	}
	item, err := domain.NewPayRun(domain.PayRun{
		TenantID:       cmd.TenantID,
		PayGroupID:     cmd.PayGroupID,
		FYID:           cmd.FYID,
		Month:          cmd.Month,
		Year:           cmd.Year,
		Status:         domain.PayRunDraft,
		EmployeeCount:  int32(len(employees)),
		ReadyCount:     0,
		BlockedCount:   0,
		GeneratedCount: 0,
		Readiness:      payRunReadinessJSON(int32(len(employees)), 0, 0, 0),
		Notes:          cleanCommandString(cmd.Notes),
	})
	if err != nil {
		s.logError("validate pay run", err, serviceTenantIDField(cmd.TenantID), serviceStringField("pay_group_id", cmd.PayGroupID.String()))
		return nil, err
	}
	var result *domain.PayRun
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		result, err = s.payGroups.CreatePayRun(txCtx, item, cmd.ActorID)
		if err != nil {
			return err
		}
		for _, employee := range employees {
			entry, err := domain.NewPayRunEmployee(domain.PayRunEmployee{TenantID: cmd.TenantID, PayRunID: result.ID, UserID: employee.UserID, ReadinessStatus: domain.PayRunEmployeePending})
			if err != nil {
				return err
			}
			if _, err := s.payGroups.UpsertPayRunEmployee(txCtx, entry, cmd.ActorID); err != nil {
				return err
			}
		}
		return s.createPayRunEvent(txCtx, result, "created", nil, &result.Status, cmd.Notes, cmd.ActorID, map[string]any{"employee_count": len(employees)})
	})
	if err != nil {
		s.logError("create pay run", err, serviceTenantIDField(cmd.TenantID), serviceStringField("pay_group_id", cmd.PayGroupID.String()))
		return nil, err
	}
	return s.AssessPayRunReadiness(ctx, cmd.TenantID, result.ID, cmd.ActorID)
}

func (s *TenantService) ListPayRuns(ctx context.Context, query ports.PayRunListQuery) ([]*domain.PayRun, error) {
	items, err := s.payGroups.ListPayRuns(ctx, query.TenantID, query.PayGroupID, query.Month, query.Year)
	if err != nil {
		s.logError("list pay runs", err, serviceTenantIDField(query.TenantID))
		return nil, err
	}
	for _, item := range items {
		if err := s.hydratePayRun(ctx, item); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *TenantService) GetPayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayRun, error) {
	item, err := s.payGroups.GetPayRun(ctx, tenantID, id)
	if err != nil {
		s.logError("get pay run", err, serviceTenantIDField(tenantID), serviceStringField("pay_run_id", id.String()))
		return nil, err
	}
	return item, s.hydratePayRun(ctx, item)
}

func (s *TenantService) GetPayRunCommandCenter(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayRunCommandCenter, error) {
	run, err := s.GetPayRun(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	summary, err := s.payGroups.GetPayRunLedgerSummary(ctx, tenantID, id)
	if err != nil {
		s.logError("get pay run command center summary", err, serviceTenantIDField(tenantID), serviceStringField("pay_run_id", id.String()))
		return nil, err
	}
	inputs, err := s.payGroups.ListPayRunInputs(ctx, tenantID, id)
	if err != nil {
		s.logError("get pay run command center inputs", err, serviceTenantIDField(tenantID), serviceStringField("pay_run_id", id.String()))
		return nil, err
	}
	components, err := s.payGroups.ListPayRunComponents(ctx, tenantID, id)
	if err != nil {
		s.logError("get pay run command center components", err, serviceTenantIDField(tenantID), serviceStringField("pay_run_id", id.String()))
		return nil, err
	}
	return &domain.PayRunCommandCenter{Run: run, Summary: summary, Inputs: inputs, Components: components}, nil
}

func (s *TenantService) AssessPayRunReadiness(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.PayRun, error) {
	run, err := s.payGroups.GetPayRun(ctx, tenantID, id)
	if err != nil {
		s.logError("assess pay run load", err, serviceTenantIDField(tenantID), serviceStringField("pay_run_id", id.String()))
		return nil, err
	}
	if run.Status == domain.PayRunLocked {
		return nil, domain.ErrPayRunLocked
	}
	employees, err := s.payGroups.ListPayRunEmployees(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	reconciliationByUser := map[uuid.UUID]*domain.PayrollReconciliationRow{}
	reconciliation, err := s.payrollOperations.ListPayrollReconciliationRows(ctx, tenantID, run.Month, run.Year)
	if err == nil {
		for _, row := range reconciliation {
			reconciliationByUser[row.UserID] = row
		}
	}
	for _, employee := range employees {
		s.assessPayRunEmployee(ctx, run, employee, reconciliationByUser[employee.UserID])
	}
	run.EmployeeCount, run.ReadyCount, run.BlockedCount, run.GeneratedCount = payRunEmployeeCounts(employees)
	if run.BlockedCount > 0 {
		run.Status = domain.PayRunBlocked
	} else {
		run.Status = domain.PayRunReadinessReady
	}
	run.Readiness = payRunReadinessJSON(run.EmployeeCount, run.ReadyCount, run.BlockedCount, run.GeneratedCount)
	var result *domain.PayRun
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		for _, employee := range employees {
			if _, err := s.payGroups.UpsertPayRunEmployee(txCtx, employee, actorID); err != nil {
				return err
			}
		}
		var err error
		result, err = s.payGroups.UpdatePayRunStatus(txCtx, run, actorID)
		if err != nil {
			return err
		}
		if err := s.rebuildPayRunDraftLedger(txCtx, result, employees, actorID); err != nil {
			return err
		}
		return s.createPayRunEvent(txCtx, result, "readiness_assessed", nil, &result.Status, nil, actorID, map[string]any{"ready_count": result.ReadyCount, "blocked_count": result.BlockedCount, "generated_count": result.GeneratedCount})
	})
	if err != nil {
		s.logError("assess pay run readiness", err, serviceTenantIDField(tenantID), serviceStringField("pay_run_id", id.String()))
		return nil, err
	}
	return result, s.hydratePayRun(ctx, result)
}

func (s *TenantService) rebuildPayRunDraftLedger(ctx context.Context, run *domain.PayRun, employees []*domain.PayRunEmployee, actorID *uuid.UUID) error {
	if run == nil {
		return domain.ErrInvalidPayRun
	}
	if err := s.payGroups.DeletePayRunLedger(ctx, run.TenantID, run.ID, actorID); err != nil {
		return err
	}
	for _, employee := range employees {
		if employee.ReadinessStatus != domain.PayRunEmployeeReady && employee.ReadinessStatus != domain.PayRunEmployeeGenerated {
			continue
		}
		calculation, err := s.CalculateEmployeeSalary(ctx, ports.EmployeeSalaryCalculationCommand{TenantID: run.TenantID, UserID: employee.UserID, FYID: run.FYID, Month: int(run.Month), Year: int(run.Year)})
		if err != nil {
			return err
		}
		salary, err := s.employeeSalaries.GetEmployeeSalary(ctx, run.TenantID, calculation.SalaryID)
		if err != nil {
			return err
		}
		attendanceInput, err := s.createPayRunInput(ctx, run, employee.UserID, "attendance", "salary_calculation", "Attendance, leave and LOP snapshot", floatPtr(float64(calculation.PresentDays)), floatPtr(calculation.SalaryResult.AbsentDeduction), actorID, map[string]any{
			"present_days": calculation.PresentDays,
			"absent_days":  calculation.AbsentDays,
			"total_days":   calculation.TotalDays,
			"lwp_days":     calculation.LWPDays,
		})
		if err != nil {
			return err
		}
		salaryInput, err := s.createPayRunInput(ctx, run, employee.UserID, "salary", "employee_salary", "Employee salary structure snapshot", nil, floatPtr(calculation.GrossSalary), actorID, map[string]any{
			"salary_id":   calculation.SalaryID.String(),
			"template_id": salary.TemplateID.String(),
		})
		if err != nil {
			return err
		}
		for _, item := range calculation.SalaryResult.Items {
			inputID := &salaryInput.ID
			if item.Code == domain.SalaryCodeLWP {
				inputID = &attendanceInput.ID
			}
			if _, err := s.payGroups.CreatePayRunComponent(ctx, &domain.PayRunComponent{
				TenantID:         run.TenantID,
				PayRunID:         run.ID,
				UserID:           employee.UserID,
				ComponentType:    item.ItemType,
				Code:             item.Code,
				Name:             item.Name,
				Amount:           item.Amount,
				SourceInputID:    inputID,
				SalaryTemplateID: &salary.TemplateID,
				Taxable:          item.ItemType == domain.SalaryItemEarning,
				SortOrder:        int32(item.SortOrder),
				Metadata:         jsonMap(map[string]any{"source": "salary_calculation"}),
			}, actorID); err != nil {
				return err
			}
		}
		statutoryItems, err := s.salarySlipStatutoryItems(ctx, run.TenantID, employee.UserID, run.Month, run.Year, calculation.GrossSalary)
		if err != nil {
			return err
		}
		for _, item := range statutoryItems {
			if _, err := s.payGroups.CreatePayRunComponent(ctx, &domain.PayRunComponent{
				TenantID:         run.TenantID,
				PayRunID:         run.ID,
				UserID:           employee.UserID,
				ComponentType:    item.ItemType,
				Code:             item.Code,
				Name:             item.Name,
				Amount:           item.Amount,
				SourceInputID:    &salaryInput.ID,
				SalaryTemplateID: &salary.TemplateID,
				Statutory:        true,
				SortOrder:        item.SortOrder,
				Metadata:         jsonMap(map[string]any{"source": "payroll_statutory_rule"}),
			}, actorID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *TenantService) createPayRunInput(ctx context.Context, run *domain.PayRun, userID uuid.UUID, inputType string, sourceType string, description string, quantity *float64, amount *float64, actorID *uuid.UUID, metadata map[string]any) (*domain.PayRunInput, error) {
	return s.payGroups.CreatePayRunInput(ctx, &domain.PayRunInput{
		TenantID:    run.TenantID,
		PayRunID:    run.ID,
		UserID:      userID,
		InputType:   inputType,
		SourceType:  sourceType,
		Description: description,
		Quantity:    quantity,
		Amount:      amount,
		Metadata:    jsonMap(metadata),
	}, actorID)
}

func (s *TenantService) FreezePayRun(ctx context.Context, cmd ports.PayRunActionCommand) (*domain.PayRun, error) {
	assessed, err := s.AssessPayRunReadiness(ctx, cmd.TenantID, cmd.PayRunID, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	if assessed.BlockedCount > 0 {
		s.logError("freeze blocked pay run", domain.ErrPayRunBlocked, serviceTenantIDField(cmd.TenantID), serviceStringField("pay_run_id", cmd.PayRunID.String()))
		return nil, domain.ErrPayRunBlocked
	}
	return s.movePayRunStatus(ctx, assessed, domain.PayRunFrozen, "frozen", cmd.Remarks, cmd.ActorID, nil)
}

func (s *TenantService) GeneratePayRun(ctx context.Context, cmd ports.PayRunActionCommand) (*domain.PayRun, error) {
	run, err := s.payGroups.GetPayRun(ctx, cmd.TenantID, cmd.PayRunID)
	if err != nil {
		return nil, err
	}
	if err := s.ensurePayrollPeriodOpen(ctx, run.TenantID, run.Month, run.Year); err != nil {
		return nil, err
	}
	if run.Status == domain.PayRunLocked {
		return nil, domain.ErrPayRunLocked
	}
	if run.Status != domain.PayRunFrozen && run.Status != domain.PayRunReadinessReady && run.Status != domain.PayRunGenerated && run.Status != domain.PayRunUnlocked {
		run, err = s.AssessPayRunReadiness(ctx, cmd.TenantID, cmd.PayRunID, cmd.ActorID)
		if err != nil {
			return nil, err
		}
		if run.BlockedCount > 0 {
			return nil, domain.ErrPayRunBlocked
		}
	}
	run, err = s.movePayRunStatus(ctx, run, domain.PayRunProcessing, "generation_started", cmd.Remarks, cmd.ActorID, nil)
	if err != nil {
		return nil, err
	}
	employees, err := s.payGroups.ListPayRunEmployees(ctx, run.TenantID, run.ID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	for _, employee := range employees {
		if employee.ReadinessStatus == domain.PayRunEmployeeBlocked || employee.ReadinessStatus == domain.PayRunEmployeeGenerated {
			continue
		}
		slip, err := s.GenerateSalarySlip(ctx, ports.GenerateSalarySlipCommand{TenantID: run.TenantID, UserID: employee.UserID, FYID: run.FYID, Month: run.Month, Year: run.Year, Regenerate: cmd.Regenerate, ActorID: cmd.ActorID})
		if err != nil {
			message := err.Error()
			employee.ReadinessStatus = domain.PayRunEmployeeFailed
			employee.BlockerReason = &message
			if errors.Is(err, domain.ErrSalarySlipExists) {
				employee.ReadinessStatus = domain.PayRunEmployeeBlocked
			}
			_, _ = s.payGroups.UpsertPayRunEmployee(ctx, employee, cmd.ActorID)
			continue
		}
		employee.ReadinessStatus = domain.PayRunEmployeeGenerated
		employee.BlockerReason = nil
		employee.SalarySlipID = &slip.ID
		employee.GeneratedAt = &now
		_, _ = s.payGroups.UpsertPayRunEmployee(ctx, employee, cmd.ActorID)
	}
	employees, err = s.payGroups.ListPayRunEmployees(ctx, run.TenantID, run.ID)
	if err != nil {
		return nil, err
	}
	run.EmployeeCount, run.ReadyCount, run.BlockedCount, run.GeneratedCount = payRunEmployeeCounts(employees)
	run.Status = domain.PayRunGenerated
	if run.BlockedCount > 0 || run.ReadyCount > 0 {
		run.Status = domain.PayRunFailed
	}
	if run.GeneratedCount > 0 && run.BlockedCount == 0 && run.ReadyCount == 0 {
		run.Status = domain.PayRunGenerated
	}
	run.Readiness = payRunReadinessJSON(run.EmployeeCount, run.ReadyCount, run.BlockedCount, run.GeneratedCount)
	result, err := s.payGroups.UpdatePayRunStatus(ctx, run, cmd.ActorID)
	if err != nil {
		s.logError("update generated pay run", err, serviceTenantIDField(run.TenantID), serviceStringField("pay_run_id", run.ID.String()))
		return nil, err
	}
	if err := s.createPayRunEvent(ctx, result, "generated", nil, &result.Status, cmd.Remarks, cmd.ActorID, map[string]any{"generated_count": result.GeneratedCount, "blocked_count": result.BlockedCount}); err != nil {
		return nil, err
	}
	return result, s.hydratePayRun(ctx, result)
}

func (s *TenantService) LockPayRun(ctx context.Context, cmd ports.PayRunActionCommand) (*domain.PayRun, error) {
	run, err := s.payGroups.GetPayRun(ctx, cmd.TenantID, cmd.PayRunID)
	if err != nil {
		return nil, err
	}
	if run.Status != domain.PayRunGenerated && run.Status != domain.PayRunUnlocked {
		return nil, domain.ErrInvalidPayRun
	}
	return s.movePayRunStatus(ctx, run, domain.PayRunLocked, "locked", cmd.Remarks, cmd.ActorID, nil)
}

func (s *TenantService) UnlockPayRun(ctx context.Context, cmd ports.PayRunActionCommand) (*domain.PayRun, error) {
	run, err := s.payGroups.GetPayRun(ctx, cmd.TenantID, cmd.PayRunID)
	if err != nil {
		return nil, err
	}
	if run.Status != domain.PayRunLocked {
		return nil, domain.ErrInvalidPayRun
	}
	return s.movePayRunStatus(ctx, run, domain.PayRunUnlocked, "unlocked", cmd.Remarks, cmd.ActorID, nil)
}

func (s *TenantService) DeletePayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	run, err := s.payGroups.GetPayRun(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if run.Status == domain.PayRunLocked {
		return domain.ErrPayRunLocked
	}
	if err := s.payGroups.DeletePayRun(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete pay run", err, serviceTenantIDField(tenantID), serviceStringField("pay_run_id", id.String()))
		return err
	}
	return nil
}

func payGroupFromCommand(cmd ports.PayGroupCommand) domain.PayGroup {
	return domain.PayGroup{
		ID:               cmd.ID,
		TenantID:         cmd.TenantID,
		Code:             cmd.Code,
		Name:             cmd.Name,
		Description:      cleanCommandString(cmd.Description),
		GroupingType:     cmd.GroupingType,
		BranchID:         cmd.BranchID,
		DepartmentID:     cmd.DepartmentID,
		EmploymentTypeID: cmd.EmploymentTypeID,
		ReportingTag:     cleanCommandString(cmd.ReportingTag),
		Rules:            cmd.Rules,
		IsActive:         cmd.IsActive,
	}
}

func (s *TenantService) hydratePayGroup(ctx context.Context, item *domain.PayGroup) (*domain.PayGroup, error) {
	members, err := s.payGroups.ListPayGroupMembers(ctx, item.TenantID, item.ID)
	if err != nil {
		return nil, err
	}
	employees, err := s.payGroups.ListPayGroupEmployees(ctx, item.TenantID, item.ID)
	if err != nil {
		return nil, err
	}
	item.Members = members
	item.EmployeeCount = int32(len(employees))
	return item, nil
}

func (s *TenantService) hydratePayRun(ctx context.Context, item *domain.PayRun) error {
	employees, err := s.payGroups.ListPayRunEmployees(ctx, item.TenantID, item.ID)
	if err != nil {
		return err
	}
	events, err := s.payGroups.ListPayRunEvents(ctx, item.TenantID, item.ID)
	if err != nil {
		return err
	}
	item.Employees = employees
	item.Events = events
	return nil
}

func (s *TenantService) assessPayRunEmployee(ctx context.Context, run *domain.PayRun, employee *domain.PayRunEmployee, reconciliation *domain.PayrollReconciliationRow) {
	employee.ReadinessStatus = domain.PayRunEmployeeReady
	employee.BlockerReason = nil
	employee.SalarySlipID = nil
	employee.GeneratedAt = nil
	salaries, err := s.employeeSalaries.ListEmployeeSalariesByUser(ctx, run.TenantID, employee.UserID)
	if err != nil {
		message := "salary assignment could not be checked"
		employee.ReadinessStatus = domain.PayRunEmployeeBlocked
		employee.BlockerReason = &message
		return
	}
	hasFYSalary := false
	for _, salary := range salaries {
		if salary.FYID == run.FYID && !salary.Inactive {
			hasFYSalary = true
			break
		}
	}
	if !hasFYSalary {
		message := "salary assignment missing for financial year"
		employee.ReadinessStatus = domain.PayRunEmployeeBlocked
		employee.BlockerReason = &message
		return
	}
	if reconciliation != nil && reconciliation.ReconciliationStatus == "lop_without_deduction" {
		message := "existing salary slip has LOP without deduction"
		employee.ReadinessStatus = domain.PayRunEmployeeBlocked
		employee.BlockerReason = &message
		return
	}
	slip, err := s.salarySlips.GetSalarySlipByPeriod(ctx, run.TenantID, employee.UserID, run.Month, run.Year)
	if err == nil {
		employee.ReadinessStatus = domain.PayRunEmployeeGenerated
		employee.SalarySlipID = &slip.ID
		employee.GeneratedAt = &slip.CreatedAt
		return
	}
	if err != nil && !isNoRows(err) {
		message := "salary slip duplicate check failed"
		employee.ReadinessStatus = domain.PayRunEmployeeBlocked
		employee.BlockerReason = &message
	}
}

func (s *TenantService) movePayRunStatus(ctx context.Context, run *domain.PayRun, status string, action string, remarks *string, actorID *uuid.UUID, metadata map[string]any) (*domain.PayRun, error) {
	before := run.Status
	run.Status = status
	run.Readiness = payRunReadinessJSON(run.EmployeeCount, run.ReadyCount, run.BlockedCount, run.GeneratedCount)
	var result *domain.PayRun
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		var err error
		result, err = s.payGroups.UpdatePayRunStatus(txCtx, run, actorID)
		if err != nil {
			return err
		}
		return s.createPayRunEvent(txCtx, result, action, &before, &result.Status, remarks, actorID, metadata)
	})
	if err != nil {
		s.logError("move pay run status", err, serviceTenantIDField(run.TenantID), serviceStringField("pay_run_id", run.ID.String()), serviceStringField("status", status))
		return nil, err
	}
	return result, s.hydratePayRun(ctx, result)
}

func (s *TenantService) createPayRunEvent(ctx context.Context, run *domain.PayRun, action string, fromStatus *string, toStatus *string, remarks *string, actorID *uuid.UUID, metadata map[string]any) error {
	raw, _ := json.Marshal(metadata)
	if len(raw) == 0 || string(raw) == "null" {
		raw = []byte(`{}`)
	}
	event, err := domain.NewPayRunEvent(domain.PayRunEvent{TenantID: run.TenantID, PayRunID: run.ID, Action: action, FromStatus: fromStatus, ToStatus: toStatus, Remarks: cleanCommandString(remarks), Metadata: raw})
	if err != nil {
		return err
	}
	_, err = s.payGroups.CreatePayRunEvent(ctx, event, actorID)
	return err
}

func payRunEmployeeCounts(employees []*domain.PayRunEmployee) (employeeCount int32, readyCount int32, blockedCount int32, generatedCount int32) {
	employeeCount = int32(len(employees))
	for _, employee := range employees {
		switch employee.ReadinessStatus {
		case domain.PayRunEmployeeReady:
			readyCount++
		case domain.PayRunEmployeeBlocked, domain.PayRunEmployeeFailed:
			blockedCount++
		case domain.PayRunEmployeeGenerated:
			generatedCount++
		}
	}
	return employeeCount, readyCount, blockedCount, generatedCount
}

func payRunReadinessJSON(employeeCount int32, readyCount int32, blockedCount int32, generatedCount int32) json.RawMessage {
	raw, _ := json.Marshal(map[string]any{
		"employee_count":  employeeCount,
		"ready_count":     readyCount,
		"blocked_count":   blockedCount,
		"generated_count": generatedCount,
		"assessed_at":     time.Now().UTC().Format(time.RFC3339),
	})
	return raw
}

func jsonMap(value map[string]any) json.RawMessage {
	raw, _ := json.Marshal(value)
	if len(raw) == 0 || string(raw) == "null" {
		return json.RawMessage(`{}`)
	}
	return raw
}

func floatPtr(value float64) *float64 {
	return &value
}
