package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

const compOffLeaveTypeShortcode = "CO"

func (s *TenantService) CreateCompOffRequest(ctx context.Context, cmd ports.CompOffRequestCommand) (*domain.CompOffRequest, error) {
	item, err := s.compOffRequestFromCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	created, err := s.compOffRequests.CreateCompOffRequest(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create comp-off request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return created, nil
}

func (s *TenantService) ListCompOffRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.CompOffRequest, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate comp-off list tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidLeaveBalanceUser
		s.logError("validate comp-off list user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.compOffRequests.ListCompOffRequestsByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("list comp-off requests by user", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListCompOffRequestsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.CompOffRequest, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate comp-off list tenant", err)
		return nil, err
	}
	normalized := domain.NormalizeCompOffStatus(status)
	if normalized == "" {
		err := domain.ErrInvalidCompOffStatus
		s.logError("validate comp-off list status", err, serviceTenantIDField(tenantID), serviceStringField("status", status))
		return nil, err
	}
	items, err := s.compOffRequests.ListCompOffRequestsByStatus(ctx, tenantID, normalized)
	if err != nil {
		s.logError("list comp-off requests by status", err, serviceTenantIDField(tenantID), serviceStringField("status", normalized))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ReviewCompOffRequest(ctx context.Context, cmd ports.CompOffReviewCommand) (*domain.CompOffRequest, error) {
	if cmd.TenantID == uuid.Nil || cmd.RequestID == uuid.Nil {
		err := domain.ErrInvalidCompOffRequest
		s.logError("validate comp-off review ids", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status := domain.NormalizeCompOffStatus(cmd.Status)
	if status != domain.CompOffStatusApproved && status != domain.CompOffStatusRejected {
		err := domain.ErrInvalidCompOffStatus
		s.logError("validate comp-off review status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("status", cmd.Status))
		return nil, err
	}
	var result *domain.CompOffRequest
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		existing, err := s.compOffRequests.GetCompOffRequest(txCtx, cmd.TenantID, cmd.RequestID)
		if err != nil {
			return err
		}
		if existing.Status != domain.CompOffStatusPending {
			return domain.ErrCompOffRequestNotPending
		}
		review := *existing
		review.Status = status
		review.ReviewRemarks = cleanCompOffString(cmd.Remarks)
		review.PayrollImpact = cmd.PayrollImpact
		review.Metadata = mergeCompOffMetadata(existing.Metadata, cmd.Metadata)
		if cmd.ExpiryDate != nil && strings.TrimSpace(*cmd.ExpiryDate) != "" {
			parsed, err := parseCompOffDate(*cmd.ExpiryDate)
			if err != nil {
				return err
			}
			review.ExpiryDate = &parsed
		}
		if status == domain.CompOffStatusApproved {
			approvedDays := existing.RequestedDays
			if cmd.ApprovedDays != nil {
				approvedDays = *cmd.ApprovedDays
			}
			if approvedDays <= 0 {
				return domain.ErrInvalidCompOffRequest
			}
			review.ApprovedDays = &approvedDays
			if _, err := s.leaveBalances.GetLeaveLedgerBySource(txCtx, existing.TenantID, existing.UserID, existing.LeaveTypeID, existing.FYID, domain.LeaveLedgerSourceCompOff, existing.ID); err == nil {
				return domain.ErrCompOffLedgerAlreadyExists
			} else if err != nil && !errors.Is(err, domain.ErrLeaveLedgerEntryNotFound) {
				return err
			}
			balance, err := s.ensureLeaveBalance(txCtx, existing.TenantID, existing.UserID, existing.LeaveTypeID, existing.FYID, cmd.ActorID)
			if err != nil {
				return err
			}
			before := *balance
			updated, err := s.leaveBalances.AddLeaveBalanceCredit(txCtx, existing.TenantID, existing.UserID, existing.LeaveTypeID, existing.FYID, approvedDays, cmd.ActorID)
			if err != nil {
				return err
			}
			remarks := "Comp-off credit approved"
			if review.ReviewRemarks != nil {
				remarks = *review.ReviewRemarks
			}
			sourceID := existing.ID
			entry := &domain.LeaveLedgerEntry{
				TenantID:        existing.TenantID,
				UserID:          existing.UserID,
				LeaveTypeID:     existing.LeaveTypeID,
				FYID:            existing.FYID,
				TransactionType: domain.LeaveLedgerCredit,
				Days:            approvedDays,
				Remarks:         &remarks,
				SourceType:      domain.LeaveLedgerSourceCompOff,
				SourceID:        &sourceID,
				BalanceBefore:   &before.BalanceDays,
				BalanceAfter:    &updated.BalanceDays,
				PendingBefore:   &before.PendingDays,
				PendingAfter:    &updated.PendingDays,
				UsedBefore:      &before.UsedDays,
				UsedAfter:       &updated.UsedDays,
				Metadata: map[string]any{
					"comp_off_request_id": existing.ID.String(),
					"work_date":           existing.WorkDate.Format("2006-01-02"),
					"worked_minutes":      existing.WorkedMinutes,
					"requested_days":      existing.RequestedDays,
					"approved_days":       approvedDays,
					"payroll_impact":      review.PayrollImpact,
				},
			}
			if review.ExpiryDate != nil {
				entry.Metadata["expiry_date"] = review.ExpiryDate.Format("2006-01-02")
			}
			if _, err := s.leaveBalances.CreateLeaveLedgerEntry(txCtx, entry, cmd.ActorID); err != nil {
				return err
			}
		}
		updated, err := s.compOffRequests.ReviewCompOffRequest(txCtx, &review, cmd.ActorID)
		if err != nil {
			return err
		}
		result = updated
		return nil
	})
	if err != nil {
		s.logError("review comp-off request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("comp_off_request_id", cmd.RequestID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) compOffRequestFromCommand(ctx context.Context, cmd ports.CompOffRequestCommand) (*domain.CompOffRequest, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate comp-off tenant", err)
		return nil, err
	}
	if cmd.UserID == uuid.Nil {
		err := domain.ErrInvalidLeaveBalanceUser
		s.logError("validate comp-off user", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	workDate, err := parseCompOffDate(cmd.WorkDate)
	if err != nil {
		s.logError("validate comp-off work date", err, serviceTenantIDField(cmd.TenantID), serviceStringField("work_date", cmd.WorkDate))
		return nil, err
	}
	fy, err := s.resolveCompOffFinancialYear(ctx, cmd.TenantID, cmd.FYID)
	if err != nil {
		return nil, err
	}
	if workDate.Before(fy.StartDate) || workDate.After(fy.EndDate) {
		err := domain.ErrInvalidCompOffRequest
		s.logError("validate comp-off work date within fy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("financial_year_id", fy.ID.String()))
		return nil, err
	}
	leaveType, err := s.resolveCompOffLeaveType(ctx, cmd.TenantID, cmd.LeaveTypeID)
	if err != nil {
		return nil, err
	}
	if !leaveType.IsEnabled {
		err := domain.ErrLeaveTypeDisabled
		s.logError("validate comp-off leave type enabled", err, serviceTenantIDField(cmd.TenantID), serviceStringField("leave_type_id", leaveType.ID.String()))
		return nil, err
	}
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID); err != nil {
		s.logError("validate comp-off employee exists", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	var expiryDate *time.Time
	if cmd.ExpiryDate != nil && strings.TrimSpace(*cmd.ExpiryDate) != "" {
		parsed, err := parseCompOffDate(*cmd.ExpiryDate)
		if err != nil {
			s.logError("validate comp-off expiry date", err, serviceTenantIDField(cmd.TenantID), serviceStringField("expiry_date", *cmd.ExpiryDate))
			return nil, err
		}
		expiryDate = &parsed
	}
	item := &domain.CompOffRequest{
		TenantID:           cmd.TenantID,
		UserID:             cmd.UserID,
		LeaveTypeID:        leaveType.ID,
		FYID:               fy.ID,
		WorkDate:           workDate,
		WorkedMinutes:      cmd.WorkedMinutes,
		RequestedDays:      cmd.RequestedDays,
		ExpiryDate:         expiryDate,
		Reason:             cleanCompOffString(cmd.Reason),
		Status:             domain.CompOffStatusPending,
		PayrollImpact:      cmd.PayrollImpact,
		SourceAttendanceID: cmd.SourceAttendanceID,
		SourceSegmentID:    cmd.SourceSegmentID,
		Metadata:           cmd.Metadata,
	}
	if err := domain.ValidateCompOffRequest(item); err != nil {
		s.logError("validate comp-off request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) resolveCompOffFinancialYear(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) (*domain.FinancialYear, error) {
	if fyID == uuid.Nil {
		fy, err := s.financialYears.GetActiveFinancialYear(ctx, tenantID)
		if err != nil {
			s.logError("resolve active financial year for comp-off", err, serviceTenantIDField(tenantID))
			return nil, err
		}
		return fy, nil
	}
	fy, err := s.financialYears.GetFinancialYear(ctx, tenantID, fyID)
	if err != nil {
		s.logError("resolve financial year for comp-off", err, serviceTenantIDField(tenantID), serviceStringField("financial_year_id", fyID.String()))
		return nil, err
	}
	return fy, nil
}

func (s *TenantService) resolveCompOffLeaveType(ctx context.Context, tenantID uuid.UUID, leaveTypeID uuid.UUID) (*domain.LeaveType, error) {
	if leaveTypeID != uuid.Nil {
		leaveType, err := s.leaveTypes.GetLeaveType(ctx, tenantID, leaveTypeID)
		if err != nil {
			s.logError("resolve comp-off leave type", err, serviceTenantIDField(tenantID), serviceStringField("leave_type_id", leaveTypeID.String()))
			return nil, err
		}
		return leaveType, nil
	}
	leaveType, err := s.leaveTypes.GetLeaveTypeByShortcode(ctx, tenantID, compOffLeaveTypeShortcode)
	if err != nil {
		s.logError("resolve default comp-off leave type", err, serviceTenantIDField(tenantID), serviceStringField("leave_type_shortcode", compOffLeaveTypeShortcode))
		return nil, err
	}
	return leaveType, nil
}

func parseCompOffDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func cleanCompOffString(value *string) *string {
	if value == nil {
		return nil
	}
	cleaned := strings.TrimSpace(*value)
	if cleaned == "" {
		return nil
	}
	return &cleaned
}

func mergeCompOffMetadata(existing map[string]any, updates map[string]any) map[string]any {
	merged := make(map[string]any, len(existing)+len(updates))
	for key, value := range existing {
		merged[key] = value
	}
	for key, value := range updates {
		merged[key] = value
	}
	return merged
}
