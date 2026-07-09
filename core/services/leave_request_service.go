package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) PreviewLeave(ctx context.Context, cmd ports.ApplyLeaveCommand) (*domain.LeavePreview, error) {
	startDate, endDate, err := parseLeaveDateRange(cmd.StartDate, cmd.EndDate)
	if err != nil {
		s.logError("parse preview leave dates", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.StartDayType == "" {
		cmd.StartDayType = domain.LeaveDayFullDay
	}
	if cmd.EndDayType == "" {
		cmd.EndDayType = domain.LeaveDayFullDay
	}
	baseDays, err := domain.CalculateLeaveDays(startDate, endDate, cmd.StartDayType, cmd.EndDayType)
	if err != nil {
		s.logError("calculate preview leave days", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	var fy *domain.FinancialYear
	if cmd.FYID == uuid.Nil {
		fy, err = s.financialYears.GetActiveFinancialYear(ctx, cmd.TenantID)
		if err != nil {
			s.logError("resolve preview leave active financial year", err, serviceTenantIDField(cmd.TenantID))
			return nil, err
		}
		cmd.FYID = fy.ID
	} else {
		fy, err = s.financialYears.GetFinancialYear(ctx, cmd.TenantID, cmd.FYID)
		if err != nil {
			s.logError("validate preview leave financial year", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
			return nil, err
		}
	}
	if startDate.Before(fy.StartDate) || endDate.After(fy.EndDate) {
		err := domain.ErrLeaveDatesOutsideFY
		s.logError("validate preview leave dates within fy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	employee, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID)
	if err != nil {
		s.logError("validate preview leave employee", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	leaveType, err := s.leaveTypes.GetLeaveType(ctx, cmd.TenantID, cmd.LeaveTypeID)
	if err != nil {
		s.logError("validate preview leave type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	preview := &domain.LeavePreview{
		TenantID:     cmd.TenantID,
		UserID:       cmd.UserID,
		LeaveTypeID:  cmd.LeaveTypeID,
		FYID:         cmd.FYID,
		StartDate:    startDate.Format("2006-01-02"),
		EndDate:      endDate.Format("2006-01-02"),
		StartDayType: cmd.StartDayType,
		EndDayType:   cmd.EndDayType,
		BaseDays:     baseDays,
		TotalDays:    baseDays,
		PaidLeave:    leaveType.IsPaid,
		Allowed:      true,
		Warnings:     []string{},
	}
	if !leaveType.IsEnabled {
		preview.Allowed = false
		preview.BlockingReasons = append(preview.BlockingReasons, domain.ErrLeaveTypeDisabled.Error())
	}
	templateRule, err := s.findApplicableLeaveRule(ctx, cmd.TenantID, employee, leaveType, cmd.FYID, startDate)
	if err != nil {
		s.logError("resolve preview leave policy rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	policy, err := s.leavePolicies.GetLeavePolicyByTypeAndFY(ctx, cmd.TenantID, cmd.LeaveTypeID, cmd.FYID)
	if err != nil && !errors.Is(err, domain.ErrLeavePolicyNotFound) {
		s.logError("load preview leave policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	if policy != nil && policy.IsSandwichApplicable {
		isSandwich, extraDays, err := s.calculateSandwichDays(ctx, cmd.TenantID, startDate, endDate, true, true)
		if err != nil {
			s.logError("calculate preview leave sandwich days", err, serviceTenantIDField(cmd.TenantID))
			return nil, err
		}
		preview.IsSandwich = isSandwich
		preview.SandwichDays = extraDays
		preview.TotalDays += extraDays
		if isSandwich && extraDays > 0 {
			preview.Warnings = append(preview.Warnings, fmt.Sprintf("Sandwich rule adds %.1f non-working day(s).", extraDays))
		}
	}
	effective, err := s.ResolveEmployeeLeavePolicySet(ctx, cmd.TenantID, cmd.UserID, startDate.Format("2006-01-02"), nil)
	if err == nil && effective != nil {
		preview.EffectivePolicy = effective.Policy
		for _, rule := range effective.LeaveRules {
			if rule.LeaveTypeID == cmd.LeaveTypeID {
				preview.EffectiveLeaveRule = rule
				preview.PayrollImpact = rule.PayrollImpact
				if !preview.IsSandwich && rule.SandwichEnabled {
					isSandwich, extraDays, err := s.calculateSandwichDays(ctx, cmd.TenantID, startDate, endDate, rule.SandwichIncludeWeeklyOff, rule.SandwichIncludePublicHoliday)
					if err != nil {
						s.logError("calculate preview leave policy-engine sandwich days", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
						return nil, err
					}
					preview.IsSandwich = isSandwich
					preview.SandwichDays = extraDays
					preview.TotalDays += extraDays
					if isSandwich && extraDays > 0 {
						preview.Warnings = append(preview.Warnings, fmt.Sprintf("Sandwich rule adds %.1f non-working day(s).", extraDays))
					}
				}
				if rule.AttachmentRequiredAfterDays != nil && preview.TotalDays > *rule.AttachmentRequiredAfterDays {
					preview.RequiresAttachment = true
					preview.AttachmentRequiredAfterDays = rule.AttachmentRequiredAfterDays
					preview.Warnings = append(preview.Warnings, "Attachment is required for this leave duration.")
				}
				if rule.NoticeRequiredAfterDays != nil && preview.TotalDays > *rule.NoticeRequiredAfterDays && rule.NoticeDays > 0 {
					preview.NoticeRequired = true
					preview.NoticeDays = rule.NoticeDays
					leadDays := int32(startDate.Sub(time.Now().UTC().Truncate(24*time.Hour)).Hours() / 24)
					if leadDays < rule.NoticeDays {
						preview.Warnings = append(preview.Warnings, fmt.Sprintf("This leave requires %d day(s) advance notice. Your manager can still review the request.", rule.NoticeDays))
					}
				}
				break
			}
		}
	}
	leave, err := domain.NewLeaveApplication(cmd.TenantID, cmd.UserID, cmd.LeaveTypeID, cmd.FYID, startDate, endDate, cmd.StartDayType, cmd.EndDayType, cmd.Reason, preview.TotalDays, preview.IsSandwich)
	if err != nil {
		s.logError("validate preview leave", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if err := validateLeaveRequestAgainstRule(leave, templateRule); err != nil {
		preview.Warnings = append(preview.Warnings, err.Error())
	}
	overlaps, err := s.leaveRequests.ListOverlappingLeaves(ctx, cmd.TenantID, cmd.UserID, cmd.StartDate, cmd.EndDate)
	if err != nil {
		s.logError("list preview leave overlaps", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	for _, existing := range overlaps {
		if cmd.ExcludeLeaveID != nil && existing.ID == *cmd.ExcludeLeaveID {
			continue
		}
		overlaps, err := domain.LeavesOverlap(leave, existing)
		if err != nil {
			s.logError("validate preview leave overlap mask", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_id", existing.ID.String()))
			return nil, err
		}
		if overlaps {
			preview.Allowed = false
			preview.BlockingReasons = append(preview.BlockingReasons, domain.ErrLeaveOverlap.Error())
			break
		}
	}
	if leaveType.IsPaid {
		balances, err := s.leaveBalances.ListLeaveBalancesByUser(ctx, cmd.TenantID, cmd.UserID)
		if err != nil {
			s.logError("list preview leave balances", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
			return nil, err
		}
		var balance *domain.LeaveBalance
		for _, item := range balances {
			if item.LeaveTypeID == cmd.LeaveTypeID && item.FYID == cmd.FYID {
				balance = item
				break
			}
		}
		if balance != nil {
			preview.BalanceBefore = balance.BalanceDays
			preview.PendingBefore = balance.PendingDays
			preview.UsedBefore = balance.UsedDays
		}
		preview.BalanceAfter = preview.BalanceBefore - preview.TotalDays
		preview.PendingAfter = preview.PendingBefore + preview.TotalDays
		negativeAllowance := 0.0
		if templateRule != nil && templateRule.NegativeBalanceAllowed {
			negativeAllowance = templateRule.MaxNegativeBalance
		}
		if preview.BalanceBefore+negativeAllowance < preview.TotalDays {
			preview.Warnings = append(preview.Warnings, domain.ErrLeaveBalanceInsufficient.Error()+". Approval may create leave without pay or negative balance.")
		}
	}
	return preview, nil
}

func (s *TenantService) ApplyLeave(ctx context.Context, cmd ports.ApplyLeaveCommand) (*domain.LeaveApplication, error) {
	startDate, endDate, err := parseLeaveDateRange(cmd.StartDate, cmd.EndDate)
	if err != nil {
		s.logError("parse apply leave dates", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if cmd.StartDayType == "" {
		cmd.StartDayType = domain.LeaveDayFullDay
	}
	if cmd.EndDayType == "" {
		cmd.EndDayType = domain.LeaveDayFullDay
	}
	days, err := domain.CalculateLeaveDays(startDate, endDate, cmd.StartDayType, cmd.EndDayType)
	if err != nil {
		s.logError("calculate apply leave days", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	var fy *domain.FinancialYear
	if cmd.FYID == uuid.Nil {
		fy, err = s.financialYears.GetActiveFinancialYear(ctx, cmd.TenantID)
		if err != nil {
			s.logError("resolve apply leave active financial year", err, serviceTenantIDField(cmd.TenantID))
			return nil, err
		}
		cmd.FYID = fy.ID
	} else {
		fy, err = s.financialYears.GetFinancialYear(ctx, cmd.TenantID, cmd.FYID)
		if err != nil {
			s.logError("validate apply leave financial year", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
			return nil, err
		}
	}
	if startDate.Before(fy.StartDate) || endDate.After(fy.EndDate) {
		err := domain.ErrLeaveDatesOutsideFY
		s.logError("validate apply leave dates within fy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", cmd.FYID.String()))
		return nil, err
	}
	employee, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID)
	if err != nil {
		s.logError("validate apply leave employee", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	leaveType, err := s.leaveTypes.GetLeaveType(ctx, cmd.TenantID, cmd.LeaveTypeID)
	if err != nil {
		s.logError("validate apply leave type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	if !leaveType.IsEnabled {
		err := domain.ErrLeaveTypeDisabled
		s.logError("validate apply leave type enabled", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	templateRule, err := s.findApplicableLeaveRule(ctx, cmd.TenantID, employee, leaveType, cmd.FYID, startDate)
	if err != nil {
		s.logError("resolve apply leave policy rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	policy, err := s.leavePolicies.GetLeavePolicyByTypeAndFY(ctx, cmd.TenantID, cmd.LeaveTypeID, cmd.FYID)
	if err != nil && !errors.Is(err, domain.ErrLeavePolicyNotFound) {
		s.logError("load apply leave policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
		return nil, err
	}
	sandwich := false
	if policy != nil && policy.IsSandwichApplicable {
		isSandwich, extraDays, err := s.calculateSandwichDays(ctx, cmd.TenantID, startDate, endDate, true, true)
		if err != nil {
			s.logError("calculate apply leave sandwich days", err, serviceTenantIDField(cmd.TenantID))
			return nil, err
		}
		sandwich = isSandwich
		days += extraDays
	}
	if !sandwich {
		effective, err := s.ResolveEmployeeLeavePolicySet(ctx, cmd.TenantID, cmd.UserID, startDate.Format("2006-01-02"), nil)
		if err == nil && effective != nil {
			for _, rule := range effective.LeaveRules {
				if rule.LeaveTypeID != cmd.LeaveTypeID || !rule.SandwichEnabled {
					continue
				}
				isSandwich, extraDays, err := s.calculateSandwichDays(ctx, cmd.TenantID, startDate, endDate, rule.SandwichIncludeWeeklyOff, rule.SandwichIncludePublicHoliday)
				if err != nil {
					s.logError("calculate apply leave policy-engine sandwich days", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", cmd.LeaveTypeID.String()))
					return nil, err
				}
				sandwich = isSandwich
				days += extraDays
				break
			}
		} else if err != nil {
			s.log.Warn().Err(err).Str("tenant_id", cmd.TenantID.String()).Str("user_id", cmd.UserID.String()).Msg("hrms: unable to resolve policy-engine leave rule during apply")
		}
	}
	leave, err := domain.NewLeaveApplication(cmd.TenantID, cmd.UserID, cmd.LeaveTypeID, cmd.FYID, startDate, endDate, cmd.StartDayType, cmd.EndDayType, cmd.Reason, days, sandwich)
	if err != nil {
		s.logError("validate apply leave", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if err := validateLeaveRequestAgainstRule(leave, templateRule); err != nil {
		s.log.Warn().Err(err).Str("tenant_id", cmd.TenantID.String()).Str("user_id", cmd.UserID.String()).Str("leave_type_id", cmd.LeaveTypeID.String()).Msg("hrms: discretionary leave submission includes policy warning")
	}
	overlaps, err := s.leaveRequests.ListOverlappingLeaves(ctx, cmd.TenantID, cmd.UserID, cmd.StartDate, cmd.EndDate)
	if err != nil {
		s.logError("list apply leave overlaps", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	for _, existing := range overlaps {
		overlaps, err := domain.LeavesOverlap(leave, existing)
		if err != nil {
			s.logError("validate apply leave overlap mask", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_id", existing.ID.String()))
			return nil, err
		}
		if overlaps {
			err := domain.ErrLeaveOverlap
			s.logError("validate apply leave overlap", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_id", existing.ID.String()))
			return nil, err
		}
	}
	fallbackApproverID := cmd.ApproverID
	if fallbackApproverID == nil || *fallbackApproverID == uuid.Nil {
		fallbackApproverID = employee.ReportingManagerID
	}
	if fallbackApproverID == nil || *fallbackApproverID == uuid.Nil {
		fallbackApproverID = cmd.ActorID
	}
	if fallbackApproverID == nil || *fallbackApproverID == uuid.Nil {
		fallbackApproverID = &cmd.UserID
		s.log.Warn().Str("tenant_id", cmd.TenantID.String()).Str("user_id", cmd.UserID.String()).Msg("hrms: leave approver missing, defaulting to applicant")
	}
	var application *domain.LeaveApplication
	err = s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		createdLeave, err := s.leaveRequests.CreateLeave(txCtx, leave, cmd.ActorID)
		if err != nil {
			return err
		}
		var updatedBalance *domain.LeaveBalance
		if leaveType.IsPaid {
			balance, err := s.ensureLeaveBalance(txCtx, cmd.TenantID, cmd.UserID, cmd.LeaveTypeID, cmd.FYID, cmd.ActorID)
			if err != nil {
				return err
			}
			if balance.BalanceDays < days {
				s.log.Warn().Str("tenant_id", cmd.TenantID.String()).Str("user_id", cmd.UserID.String()).Float64("balance_days", balance.BalanceDays).Float64("request_days", days).Msg("hrms: discretionary leave submission exceeds current balance")
			}
			before := *balance
			updatedBalance, err = s.leaveBalances.UpdateLeaveBalancePending(txCtx, cmd.TenantID, cmd.UserID, cmd.LeaveTypeID, cmd.FYID, days, cmd.ActorID)
			if err != nil {
				return err
			}
			ledgerRemarks := "leave application pending reservation"
			ledger := &domain.LeaveLedgerEntry{TenantID: cmd.TenantID, UserID: cmd.UserID, LeaveTypeID: cmd.LeaveTypeID, FYID: cmd.FYID, LeaveID: &createdLeave.ID, TransactionType: domain.LeaveLedgerDebit, Days: days, Remarks: &ledgerRemarks, SourceType: domain.LeaveLedgerSourceLeaveApply, SourceID: &createdLeave.ID, BalanceBefore: &before.BalanceDays, BalanceAfter: &updatedBalance.BalanceDays, PendingBefore: &before.PendingDays, PendingAfter: &updatedBalance.PendingDays, UsedBefore: &before.UsedDays, UsedAfter: &updatedBalance.UsedDays, Metadata: map[string]any{"leave_type_name": leaveType.Name}}
			if _, err := s.leaveBalances.CreateLeaveLedgerEntry(txCtx, ledger, cmd.ActorID); err != nil {
				return err
			}
		}
		createdApproval, err := s.createInitialLeaveApproval(txCtx, createdLeave, *fallbackApproverID, cmd.ActorID)
		if err != nil {
			return err
		}
		application = &domain.LeaveApplication{Leave: createdLeave, Approval: createdApproval, Balance: updatedBalance}
		return nil
	})
	if err != nil {
		s.logError("apply leave", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	if s.leaveNotifier != nil {
		if err := s.leaveNotifier.NotifyLeaveApplied(ctx, application); err != nil {
			s.logError("notify leave applied", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_id", application.Leave.ID.String()))
		}
	} else {
		s.log.Warn().Str("tenant_id", cmd.TenantID.String()).Str("leave_id", application.Leave.ID.String()).Msg("hrms: leave notification hook not configured")
	}
	s.notifyLeaveApplied(ctx, application, cmd.ActorID)
	return application, nil
}

func (s *TenantService) ListLeavesByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.Leave, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave list tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidLeaveUser
		s.logError("validate leave list user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.leaveRequests.ListLeavesByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("list leaves by user", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListLeavesByFY(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) ([]*domain.Leave, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate leave fy list tenant", err)
		return nil, err
	}
	if fyID == uuid.Nil {
		err := domain.ErrInvalidLeavePolicyFY
		s.logError("validate leave fy list fy", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.leaveRequests.ListLeavesByFY(ctx, tenantID, fyID)
	if err != nil {
		s.logError("list leaves by fy", err, serviceTenantIDField(tenantID), serviceStringField("financial_year_id", fyID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListLeaveReportRows(ctx context.Context, filter domain.LeaveReportFilter) ([]*domain.LeaveReportRow, error) {
	if err := validateLeaveReportFilter(filter); err != nil {
		s.logError("validate leave report rows filter", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	items, err := s.leaveRequests.ListLeaveReportRows(ctx, filter)
	if err != nil {
		s.logError("list leave report rows", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetLeaveReportSummary(ctx context.Context, filter domain.LeaveReportFilter) (*domain.LeaveReportSummary, error) {
	if err := validateLeaveReportFilter(filter); err != nil {
		s.logError("validate leave report summary filter", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	item, err := s.leaveRequests.GetLeaveReportSummary(ctx, filter)
	if err != nil {
		s.logError("get leave report summary", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return item, nil
}

func parseLeaveDateRange(startDate string, endDate string) (time.Time, time.Time, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrInvalidDateRange
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrInvalidDateRange
	}
	if end.Before(start) {
		return time.Time{}, time.Time{}, domain.ErrInvalidDateRange
	}
	return start, end, nil
}

func (s *TenantService) calculateSandwichDays(ctx context.Context, tenantID uuid.UUID, startDate time.Time, endDate time.Time, includeWeeklyOff bool, includePublicHoliday bool) (bool, float64, error) {
	holidayDates := []time.Time{}
	if includePublicHoliday {
		holidays, err := s.holidays.ListHolidaysByDateRange(ctx, tenantID, startDate, endDate)
		if err != nil {
			return false, 0, err
		}
		holidayDates = make([]time.Time, 0, len(holidays))
		for _, holiday := range holidays {
			if !holiday.IsOptional {
				holidayDates = append(holidayDates, holiday.Date)
			}
		}
	}
	weeklyOffs := []time.Weekday{}
	if includeWeeklyOff {
		weeklyOffs = []time.Weekday{time.Saturday, time.Sunday}
	}
	isSandwich, extraDays := domain.IsSandwich(startDate, endDate, holidayDates, weeklyOffs)
	return isSandwich, extraDays, nil
}

func validateLeaveReportFilter(filter domain.LeaveReportFilter) error {
	if filter.TenantID == uuid.Nil {
		return domain.ErrInvalidTenantID
	}
	if filter.StartDate != nil && filter.EndDate != nil && filter.EndDate.Before(*filter.StartDate) {
		return domain.ErrInvalidDateRange
	}
	if filter.Status != nil && *filter.Status != "" {
		switch *filter.Status {
		case domain.LeaveStatusPending, domain.LeaveStatusApproved, domain.LeaveStatusRejected, domain.LeaveStatusCanceled:
			return nil
		default:
			return domain.ErrInvalidEnumValue
		}
	}
	return nil
}

func (s *TenantService) findApplicableLeaveRule(ctx context.Context, tenantID uuid.UUID, employee *domain.Employee, leaveType *domain.LeaveType, fyID uuid.UUID, asOfDate time.Time) (*domain.LeavePolicyTemplateRule, error) {
	rules, err := s.leaveTemplates.ListLeavePolicyTemplateRulesByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	isProbation := employee.IsOnProbation(asOfDate)
	for _, rule := range rules {
		if leaveType == nil || rule.LeaveTypeID != leaveType.ID {
			continue
		}
		if rule.FYID != nil && *rule.FYID != fyID {
			continue
		}
		if !uuidPtrMatches(rule.EmploymentTypeID, employee.EmploymentTypeID) {
			continue
		}
		if !uuidPtrMatches(rule.DepartmentID, employee.DepartmentID) {
			continue
		}
		if !uuidPtrMatches(rule.DesignationID, employee.DesignationID) {
			continue
		}
		if rule.ProbationStatus != nil {
			if *rule.ProbationStatus == domain.LeaveProbationOnly && !isProbation {
				continue
			}
			if *rule.ProbationStatus == domain.LeaveProbationConfirmed && isProbation {
				continue
			}
		}
		return rule, nil
	}
	return nil, nil
}

func validateLeaveRequestAgainstRule(leave *domain.Leave, rule *domain.LeavePolicyTemplateRule) error {
	if leave == nil || rule == nil {
		return nil
	}
	if !rule.AllowHalfDay && (leave.StartDayType != domain.LeaveDayFullDay || leave.EndDayType != domain.LeaveDayFullDay) {
		return domain.ErrLeaveHalfDayNotAllowed
	}
	if rule.MinRequestDays > 0 && leave.Days < rule.MinRequestDays {
		return domain.ErrLeaveRequestBelowMinimum
	}
	if rule.MaxRequestDays != nil && *rule.MaxRequestDays > 0 && leave.Days > *rule.MaxRequestDays {
		return domain.ErrLeaveRequestAboveMaximum
	}
	return nil
}

func uuidPtrMatches(ruleValue *uuid.UUID, employeeValue *uuid.UUID) bool {
	if ruleValue == nil || *ruleValue == uuid.Nil {
		return true
	}
	if employeeValue == nil || *employeeValue == uuid.Nil {
		return false
	}
	return *ruleValue == *employeeValue
}
