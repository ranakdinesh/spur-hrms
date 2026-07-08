package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateOvertimeRequest(ctx context.Context, cmd ports.OvertimeRequestCommand) (*domain.OvertimeRequest, error) {
	item, err := s.overtimeRequestFromCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}
	created, err := s.overtimeRequests.CreateOvertimeRequest(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create overtime request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return created, nil
}

func (s *TenantService) ListOvertimeRequestsByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) ([]*domain.OvertimeRequest, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate overtime list tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidEmployeeUserID
		s.logError("validate overtime list user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.overtimeRequests.ListOvertimeRequestsByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("list overtime requests by user", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListOvertimeRequestsByStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.OvertimeRequest, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate overtime list tenant", err)
		return nil, err
	}
	normalized := domain.NormalizeOvertimeStatus(status)
	if normalized == "" {
		err := domain.ErrInvalidOvertimeStatus
		s.logError("validate overtime list status", err, serviceTenantIDField(tenantID), serviceStringField("status", status))
		return nil, err
	}
	items, err := s.overtimeRequests.ListOvertimeRequestsByStatus(ctx, tenantID, normalized)
	if err != nil {
		s.logError("list overtime requests by status", err, serviceTenantIDField(tenantID), serviceStringField("status", normalized))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ListOvertimeRequestsByPayrollExportStatus(ctx context.Context, tenantID uuid.UUID, status string) ([]*domain.OvertimeRequest, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate overtime export list tenant", err)
		return nil, err
	}
	normalized := domain.NormalizeOvertimePayrollExportStatus(status)
	if normalized == "" {
		err := domain.ErrInvalidOvertimeRequest
		s.logError("validate overtime export list status", err, serviceTenantIDField(tenantID), serviceStringField("payroll_export_status", status))
		return nil, err
	}
	items, err := s.overtimeRequests.ListOvertimeRequestsByPayrollExportStatus(ctx, tenantID, normalized)
	if err != nil {
		s.logError("list overtime requests by payroll export status", err, serviceTenantIDField(tenantID), serviceStringField("payroll_export_status", normalized))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) ReviewOvertimeRequest(ctx context.Context, cmd ports.OvertimeReviewCommand) (*domain.OvertimeRequest, error) {
	if cmd.TenantID == uuid.Nil || cmd.RequestID == uuid.Nil {
		err := domain.ErrInvalidOvertimeRequest
		s.logError("validate overtime review ids", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	status := domain.NormalizeOvertimeStatus(cmd.Status)
	if status != domain.OvertimeStatusApproved && status != domain.OvertimeStatusRejected {
		err := domain.ErrInvalidOvertimeStatus
		s.logError("validate overtime review status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("status", cmd.Status))
		return nil, err
	}
	var result *domain.OvertimeRequest
	err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		existing, err := s.overtimeRequests.GetOvertimeRequest(txCtx, cmd.TenantID, cmd.RequestID)
		if err != nil {
			return err
		}
		if existing.Status != domain.OvertimeStatusPending {
			return domain.ErrOvertimeRequestNotPending
		}
		review := *existing
		review.Status = status
		review.ReviewRemarks = cleanOvertimeString(cmd.Remarks)
		review.Metadata = mergeOvertimeMetadata(existing.Metadata, cmd.Metadata)
		review.CalculationType = firstOvertimeCalculation(cmd.CalculationType, existing.CalculationType)
		review.RateMultiplier = firstPositiveFloat(cmd.RateMultiplier, existing.RateMultiplier, 1)
		if cmd.PayrollComponentCode != nil {
			review.PayrollComponentCode = cleanOvertimeString(cmd.PayrollComponentCode)
		}
		if status == domain.OvertimeStatusApproved {
			approvedMinutes := existing.RequestedMinutes
			if cmd.ApprovedMinutes != nil {
				approvedMinutes = *cmd.ApprovedMinutes
			}
			if approvedMinutes <= 0 {
				return domain.ErrInvalidOvertimeRequest
			}
			review.ApprovedMinutes = &approvedMinutes
			review.PayrollExportStatus = domain.OvertimePayrollExportReady
		} else {
			review.ApprovedMinutes = nil
			review.PayrollExportStatus = domain.OvertimePayrollExportNotApplicable
		}
		if err := domain.ValidateOvertimeRequest(&review); err != nil {
			return err
		}
		updated, err := s.overtimeRequests.ReviewOvertimeRequest(txCtx, &review, cmd.ActorID)
		if err != nil {
			return err
		}
		result = updated
		return nil
	})
	if err != nil {
		s.logError("review overtime request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("overtime_request_id", cmd.RequestID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) overtimeRequestFromCommand(ctx context.Context, cmd ports.OvertimeRequestCommand) (*domain.OvertimeRequest, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate overtime tenant", err)
		return nil, err
	}
	if cmd.UserID == uuid.Nil {
		err := domain.ErrInvalidEmployeeUserID
		s.logError("validate overtime user", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	workDate, err := parseOvertimeDate(cmd.WorkDate)
	if err != nil {
		s.logError("validate overtime work date", err, serviceTenantIDField(cmd.TenantID), serviceStringField("work_date", cmd.WorkDate))
		return nil, err
	}
	if _, err := s.employees.GetEmployeeByUserID(ctx, cmd.TenantID, cmd.UserID); err != nil {
		s.logError("validate overtime employee exists", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	item := &domain.OvertimeRequest{
		TenantID:             cmd.TenantID,
		UserID:               cmd.UserID,
		WorkDate:             workDate,
		RequestedMinutes:     cmd.RequestedMinutes,
		Reason:               cleanOvertimeString(cmd.Reason),
		Status:               domain.OvertimeStatusPending,
		CalculationType:      firstOvertimeCalculation(cmd.CalculationType, domain.OvertimeCalculationMultiplier),
		RateMultiplier:       firstPositiveFloat(cmd.RateMultiplier, 1),
		PayrollComponentCode: cleanOvertimeString(cmd.PayrollComponentCode),
		PayrollExportStatus:  domain.OvertimePayrollExportNotReady,
		SourceAttendanceID:   cmd.SourceAttendanceID,
		SourceSegmentID:      cmd.SourceSegmentID,
		Metadata:             cmd.Metadata,
	}
	if err := domain.ValidateOvertimeRequest(item); err != nil {
		s.logError("validate overtime request", err, serviceTenantIDField(cmd.TenantID), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return item, nil
}

func parseOvertimeDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func cleanOvertimeString(value *string) *string {
	if value == nil {
		return nil
	}
	cleaned := strings.TrimSpace(*value)
	if cleaned == "" {
		return nil
	}
	return &cleaned
}

func firstOvertimeCalculation(values ...string) string {
	for _, value := range values {
		if normalized := domain.NormalizeOvertimeCalculationType(value); normalized != "" {
			return normalized
		}
	}
	return domain.OvertimeCalculationMultiplier
}

func firstPositiveFloat(values ...float64) float64 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}

func mergeOvertimeMetadata(existing map[string]any, updates map[string]any) map[string]any {
	merged := make(map[string]any, len(existing)+len(updates))
	for key, value := range existing {
		merged[key] = value
	}
	for key, value := range updates {
		merged[key] = value
	}
	return merged
}
