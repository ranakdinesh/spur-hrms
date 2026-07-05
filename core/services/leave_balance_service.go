package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) UpsertLeaveBalance(ctx context.Context, cmd ports.LeaveBalanceCommand) (*domain.LeaveBalance, error) {
	item := &domain.LeaveBalance{TenantID: cmd.TenantID, UserID: cmd.UserID, LeaveTypeID: cmd.LeaveTypeID, FYID: cmd.FYID, TotalDays: cmd.TotalDays, UsedDays: cmd.UsedDays, PendingDays: cmd.PendingDays}
	if err := s.validateLeaveBalanceRefs(ctx, item); err != nil {
		return nil, err
	}
	result, err := s.leaveBalances.UpsertLeaveBalance(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("upsert leave balance", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListLeaveBalancesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.LeaveBalance, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave balance list tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidLeaveBalanceUser
		s.logError("validate leave balance list user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.leaveBalances.ListLeaveBalancesByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("list leave balances by user", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListLeaveBalancesByTenantFY(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) ([]*domain.LeaveBalance, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave balance report tenant", err)
		return nil, err
	}
	if fyID == uuid.Nil {
		err := domain.ErrInvalidLeavePolicyFY
		s.logError("validate leave balance report fy", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.leaveBalances.ListLeaveBalancesByTenantFY(ctx, tenantID, fyID)
	if err != nil {
		s.logError("list leave balances by tenant fy", err, serviceTenantIDField(tenantID), serviceStringField("financial_year_id", fyID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) AdjustLeaveBalance(ctx context.Context, cmd ports.LeaveBalanceAdjustmentCommand) (*domain.LeaveBalance, error) {
	if cmd.SourceType == "" {
		cmd.SourceType = domain.LeaveLedgerSourceManualAdjustment
	}
	if cmd.TransactionType == "" {
		cmd.TransactionType = domain.LeaveLedgerCredit
	}
	entry := &domain.LeaveLedgerEntry{TenantID: cmd.TenantID, UserID: cmd.UserID, LeaveTypeID: cmd.LeaveTypeID, FYID: cmd.FYID, TransactionType: cmd.TransactionType, Days: cmd.Days, Remarks: cmd.Remarks, SourceType: cmd.SourceType, SourceID: cmd.SourceID, LeaveID: cmd.LeaveID, Metadata: cmd.Metadata}
	if err := domain.ValidateLeaveLedgerEntry(entry); err != nil {
		s.logError("validate leave balance adjustment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	var result *domain.LeaveBalance
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		balance, err := s.ensureLeaveBalance(txCtx, cmd.TenantID, cmd.UserID, cmd.LeaveTypeID, cmd.FYID, cmd.ActorID)
		if err != nil {
			return err
		}
		before := *balance
		delta := cmd.Days
		if cmd.TransactionType == domain.LeaveLedgerDebit {
			delta = -delta
		}
		updated, err := s.leaveBalances.AddLeaveBalanceCredit(txCtx, cmd.TenantID, cmd.UserID, cmd.LeaveTypeID, cmd.FYID, delta, cmd.ActorID)
		if err != nil {
			return err
		}
		entry.BalanceBefore = &before.BalanceDays
		entry.BalanceAfter = &updated.BalanceDays
		entry.PendingBefore = &before.PendingDays
		entry.PendingAfter = &updated.PendingDays
		entry.UsedBefore = &before.UsedDays
		entry.UsedAfter = &updated.UsedDays
		if _, err := s.leaveBalances.CreateLeaveLedgerEntry(txCtx, entry, cmd.ActorID); err != nil {
			return err
		}
		result = updated
		return nil
	})
	if err != nil {
		s.logError("adjust leave balance", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListLeaveLedgerByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.LeaveLedgerEntry, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave ledger list tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidLeaveBalanceUser
		s.logError("validate leave ledger list user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.leaveBalances.ListLeaveLedgerByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("list leave ledger by user", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) RunLeaveAccrual(ctx context.Context, cmd ports.RunLeaveAccrualCommand) ([]*domain.LeaveLedgerEntry, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave accrual tenant", err)
		return nil, err
	}
	if cmd.FYID == uuid.Nil {
		err := domain.ErrInvalidLeavePolicyFY
		s.logError("validate leave accrual fy", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.Month < 1 || cmd.Month > 12 {
		err := domain.ErrInvalidLeaveAllocationMonth
		s.logError("validate leave accrual month", err, serviceTenantIDField(cmd.TenantID), serviceStringField("month", fmt.Sprint(cmd.Month)))
		return nil, err
	}
	financialYear, err := s.financialYears.GetFinancialYear(ctx, cmd.TenantID, cmd.FYID)
	if err != nil {
		s.logError("validate leave accrual fy exists", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	var ledgerEntries []*domain.LeaveLedgerEntry
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		rules, err := s.leaveTemplates.ListLeavePolicyTemplateRulesByTenant(txCtx, cmd.TenantID)
		if err != nil {
			return err
		}
		allEmployees, err := s.employees.ListEmployees(txCtx, cmd.TenantID)
		if err != nil {
			return err
		}
		employeeByUserID := make(map[uuid.UUID]*domain.Employee, len(allEmployees))
		for _, item := range allEmployees {
			if item == nil {
				continue
			}
			employee := item.Employee
			employeeByUserID[employee.UserID] = &employee
		}
		asOfDate := leaveAccrualAsOfDate(cmd.Month, financialYear)
		for _, rule := range rules {
			if rule.FYID != nil && *rule.FYID != cmd.FYID {
				continue
			}
			if rule.AccrualFrequency != domain.LeaveAccrualFrequencyMonthly && rule.AccrualFrequency != domain.LeaveAccrualFrequencyYearly {
				continue
			}
			employees, err := s.resolveLeaveAccrualEmployees(txCtx, cmd.TenantID, cmd.FYID, rule, allEmployees, employeeByUserID, asOfDate)
			if err != nil {
				return err
			}
			for _, employee := range employees {
				isProbation := employee.IsOnProbation(asOfDate)
				if isProbation {
					leaveType, err := s.leaveTypes.GetLeaveType(txCtx, cmd.TenantID, rule.LeaveTypeID)
					if err != nil {
						return err
					}
					if isEarnedLeaveType(leaveType) && (rule.ProbationStatus == nil || *rule.ProbationStatus != domain.LeaveProbationOnly) {
						continue
					}
				}
				joiningDate := asOfDate
				if employee.JoiningDate != nil {
					joiningDate = *employee.JoiningDate
				}
				result, err := domain.CalculateLeaveAccrual(rule, domain.LeaveAccrualInput{JoiningDate: joiningDate, ConfirmationDate: employee.ProbationConfirmedAt, AsOfDate: asOfDate, PeriodStart: financialYear.StartDate, PeriodEnd: financialYear.EndDate, IsProbation: isProbation})
				if err != nil {
					return err
				}
				if result.Days <= 0 {
					continue
				}
				sourceID := accrualSourceID(cmd.TenantID, employee.UserID, cmd.FYID, rule.ID, cmd.Month)
				if _, err := s.leaveBalances.GetLeaveLedgerBySource(txCtx, cmd.TenantID, employee.UserID, rule.LeaveTypeID, cmd.FYID, result.SourceType, sourceID); err == nil {
					continue
				} else if err != nil && !errors.Is(err, domain.ErrLeaveLedgerEntryNotFound) {
					return err
				}
				balance, err := s.ensureLeaveBalance(txCtx, cmd.TenantID, employee.UserID, rule.LeaveTypeID, cmd.FYID, cmd.ActorID)
				if err != nil {
					return err
				}
				before := *balance
				creditDays := cappedAccrualDays(result.Days, rule, balance)
				if creditDays <= 0 {
					continue
				}
				updated, err := s.leaveBalances.AddLeaveBalanceCredit(txCtx, cmd.TenantID, employee.UserID, rule.LeaveTypeID, cmd.FYID, creditDays, cmd.ActorID)
				if err != nil {
					return err
				}
				remarks := result.Remarks
				entry := &domain.LeaveLedgerEntry{TenantID: cmd.TenantID, UserID: employee.UserID, LeaveTypeID: rule.LeaveTypeID, FYID: cmd.FYID, TransactionType: domain.LeaveLedgerCredit, Days: creditDays, Remarks: &remarks, SourceType: result.SourceType, SourceID: &sourceID, BalanceBefore: &before.BalanceDays, BalanceAfter: &updated.BalanceDays, PendingBefore: &before.PendingDays, PendingAfter: &updated.PendingDays, UsedBefore: &before.UsedDays, UsedAfter: &updated.UsedDays, Metadata: map[string]any{"month": cmd.Month, "rule_id": rule.ID.String(), "template_id": rule.TemplateID.String(), "calculated_days": result.Days}}
				created, err := s.leaveBalances.CreateLeaveLedgerEntry(txCtx, entry, cmd.ActorID)
				if err != nil {
					return err
				}
				ledgerEntries = append(ledgerEntries, created)
			}
		}
		return nil
	})
	if err != nil {
		s.logError("run leave accrual", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	return ledgerEntries, nil
}

func leaveAccrualAsOfDate(month int32, financialYear *domain.FinancialYear) time.Time {
	now := time.Now().UTC()
	if month < 1 || month > 12 {
		return now
	}
	if financialYear != nil {
		start := financialYear.StartDate
		end := financialYear.EndDate
		for year := start.Year(); year <= end.Year(); year++ {
			firstOfNextMonth := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC)
			candidate := firstOfNextMonth.AddDate(0, 0, -1)
			if !candidate.Before(start) && !candidate.After(end) {
				return candidate
			}
		}
	}
	firstOfNextMonth := time.Date(now.Year(), time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC)
	return firstOfNextMonth.AddDate(0, 0, -1)
}

func (s *TenantService) resolveLeaveAccrualEmployees(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID, rule *domain.LeavePolicyTemplateRule, allEmployees []*domain.EmployeeListItem, employeeByUserID map[uuid.UUID]*domain.Employee, asOfDate time.Time) ([]*domain.Employee, error) {
	candidates := make(map[uuid.UUID]*domain.Employee)
	assignments, err := s.leaveTemplates.ListLeavePolicyAssignmentsByTemplate(ctx, tenantID, rule.TemplateID)
	if err != nil {
		return nil, err
	}
	for _, assignment := range assignments {
		if assignment.FYID != nil && *assignment.FYID != fyID {
			continue
		}
		employee := employeeByUserID[assignment.UserID]
		if employee == nil {
			employee, err = s.employees.GetEmployeeByUserID(ctx, tenantID, assignment.UserID)
			if err != nil {
				return nil, err
			}
		}
		if leaveRuleMatchesEmployee(rule, employee, asOfDate) {
			candidates[employee.UserID] = employee
		}
	}
	if boolFromLeaveRuleConfig(rule.CalculationConfig, "auto_apply_by_scope", false) {
		for _, item := range allEmployees {
			if item == nil {
				continue
			}
			employee := item.Employee
			if leaveRuleMatchesEmployee(rule, &employee, asOfDate) {
				candidates[employee.UserID] = &employee
			}
		}
	}
	items := make([]*domain.Employee, 0, len(candidates))
	for _, employee := range candidates {
		items = append(items, employee)
	}
	return items, nil
}

func leaveRuleMatchesEmployee(rule *domain.LeavePolicyTemplateRule, employee *domain.Employee, asOfDate time.Time) bool {
	if rule == nil || employee == nil || employee.Inactive {
		return false
	}
	if !uuidPtrMatches(rule.EmploymentTypeID, employee.EmploymentTypeID) {
		return false
	}
	if !uuidPtrMatches(rule.DepartmentID, employee.DepartmentID) {
		return false
	}
	if !uuidPtrMatches(rule.DesignationID, employee.DesignationID) {
		return false
	}
	if rule.ProbationStatus != nil {
		isProbation := employee.IsOnProbation(asOfDate)
		if *rule.ProbationStatus == domain.LeaveProbationOnly && !isProbation {
			return false
		}
		if *rule.ProbationStatus == domain.LeaveProbationConfirmed && isProbation {
			return false
		}
	}
	return true
}

func boolFromLeaveRuleConfig(config map[string]any, key string, fallback bool) bool {
	if config == nil {
		return fallback
	}
	switch value := config[key].(type) {
	case bool:
		return value
	case string:
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "true", "yes", "1", "on":
			return true
		case "false", "no", "0", "off":
			return false
		default:
			return fallback
		}
	default:
		return fallback
	}
}

func isEarnedLeaveType(leaveType *domain.LeaveType) bool {
	if leaveType == nil {
		return false
	}
	shortcode := strings.ToLower(strings.TrimSpace(valueFromPtr(leaveType.Shortcode)))
	name := strings.ToLower(strings.TrimSpace(leaveType.Name))
	return shortcode == "el" || shortcode == domain.LeaveTypeShortEarnLeave || strings.Contains(name, "earned leave") || strings.Contains(name, "earn leave")
}

func cappedAccrualDays(days float64, rule *domain.LeavePolicyTemplateRule, balance *domain.LeaveBalance) float64 {
	creditDays := days
	if rule.AnnualEntitlement > 0 {
		remaining := rule.AnnualEntitlement - balance.TotalDays
		if remaining < creditDays {
			creditDays = remaining
		}
	}
	if rule.MaxBalance != nil && *rule.MaxBalance >= 0 {
		remaining := *rule.MaxBalance - balance.TotalDays
		if remaining < creditDays {
			creditDays = remaining
		}
	}
	if creditDays <= 0 {
		return 0
	}
	return math.Round(creditDays*100) / 100
}

func (s *TenantService) validateLeaveBalanceRefs(ctx context.Context, item *domain.LeaveBalance) error {
	if err := domain.ValidateLeaveBalance(item); err != nil {
		s.logError("validate leave balance", err, serviceTenantIDField(item.TenantID))
		return err
	}
	if _, err := s.employees.GetEmployeeByUserID(ctx, item.TenantID, item.UserID); err != nil {
		s.logError("validate leave balance employee", err, serviceTenantIDField(item.TenantID), serviceStringField("user_id", item.UserID.String()))
		return err
	}
	if _, err := s.leaveTypes.GetLeaveType(ctx, item.TenantID, item.LeaveTypeID); err != nil {
		s.logError("validate leave balance leave type", err, serviceTenantIDField(item.TenantID), serviceStringField("leave_type_id", item.LeaveTypeID.String()))
		return err
	}
	if _, err := s.financialYears.GetFinancialYear(ctx, item.TenantID, item.FYID); err != nil {
		s.logError("validate leave balance financial year", err, serviceTenantIDField(item.TenantID), serviceStringField("financial_year_id", item.FYID.String()))
		return err
	}
	return nil
}

func (s *TenantService) ensureLeaveBalance(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, leaveTypeID uuid.UUID, fyID uuid.UUID, actorID *uuid.UUID) (*domain.LeaveBalance, error) {
	balance, err := s.leaveBalances.GetLeaveBalance(ctx, tenantID, userID, leaveTypeID, fyID)
	if err == nil {
		return balance, nil
	}
	if !errors.Is(err, domain.ErrLeaveBalanceNotFound) {
		return nil, err
	}
	item := &domain.LeaveBalance{TenantID: tenantID, UserID: userID, LeaveTypeID: leaveTypeID, FYID: fyID}
	if err := s.validateLeaveBalanceRefs(ctx, item); err != nil {
		return nil, err
	}
	return s.leaveBalances.UpsertLeaveBalance(ctx, item, actorID)
}

func accrualSourceID(tenantID uuid.UUID, userID uuid.UUID, fyID uuid.UUID, ruleID uuid.UUID, month int32) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("hrms-leave-accrual:%s:%s:%s:%s:%d", tenantID, userID, fyID, ruleID, month)))
}
