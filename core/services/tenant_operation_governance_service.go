package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

const maxTenantOperationRequests = 500

func (s *TenantService) CreateTenantOperationRequest(ctx context.Context, cmd ports.TenantOperationCommand) (*domain.TenantOperationRequest, error) {
	retentionUntil, err := parseTenantOperationOptionalTime(cmd.RetentionUntil)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "parse tenant operation retention").Msg("invalid tenant operation retention")
		return nil, domain.ErrInvalidTenantOperation
	}
	item := &domain.TenantOperationRequest{
		ID:                cmd.ID,
		OperationNumber:   generateTenantOperationNumber(),
		OperationType:     cmd.OperationType,
		Title:             cmd.Title,
		TargetTenantID:    cmd.TargetTenantID,
		TargetTenantName:  cleanOptionalString(cmd.TargetTenantName),
		TargetTenantCode:  cleanOptionalString(cmd.TargetTenantCode),
		Status:            domain.TenantOperationPendingValidation,
		RiskLevel:         cmd.RiskLevel,
		Reason:            cmd.Reason,
		RequestedBy:       cmd.ActorID,
		ApprovalRequired:  true,
		BackupConfirmed:   cmd.BackupConfirmed,
		RetentionUntil:    retentionUntil,
		RequestPayload:    json.RawMessage(cmd.RequestPayload),
		ValidationResults: json.RawMessage(cmd.ValidationResults),
		RollbackMetadata:  json.RawMessage(cmd.RollbackMetadata),
		Metadata:          json.RawMessage(cmd.Metadata),
	}
	if err := domain.NormalizeTenantOperationRequest(item); err != nil {
		s.log.Warn().Err(err).Str("operation", "create tenant operation request").Str("operation_type", cmd.OperationType).Msg("invalid tenant operation request")
		return nil, err
	}
	created, err := s.tenantOperations.CreateTenantOperationRequest(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	toStatus := created.Status
	if _, err := s.tenantOperations.CreateTenantOperationEvent(ctx, &domain.TenantOperationEvent{RequestID: created.ID, Action: domain.TenantOperationActionCreated, ToStatus: &toStatus, ActorUserID: cmd.ActorID, Metadata: rawJSON(map[string]any{"research_pattern": "governed_request_validation_approval_execution"})}, cmd.ActorID); err != nil {
		return nil, err
	}
	return created, nil
}

func (s *TenantService) ListTenantOperationRequests(ctx context.Context, filter domain.TenantOperationFilter) (*domain.TenantOperationWorkspace, error) {
	filter.Limit = boundedTenantOperationLimit(filter.Limit)
	filter.Status = cleanOptionalString(filter.Status)
	filter.OperationType = cleanOptionalString(filter.OperationType)
	filter.RiskLevel = cleanOptionalString(filter.RiskLevel)
	filter.Search = cleanOptionalString(filter.Search)
	items, err := s.tenantOperations.ListTenantOperationRequests(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &domain.TenantOperationWorkspace{Requests: items, Summary: buildTenantOperationSummary(items)}, nil
}

func (s *TenantService) GetTenantOperationDetail(ctx context.Context, id uuid.UUID) (*domain.TenantOperationDetail, error) {
	if id == uuid.Nil {
		return nil, domain.ErrInvalidTenantOperation
	}
	item, err := s.tenantOperations.GetTenantOperationRequest(ctx, id)
	if err != nil {
		return nil, err
	}
	events, err := s.tenantOperations.ListTenantOperationEvents(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.TenantOperationDetail{Request: item, Events: events}, nil
}

func (s *TenantService) ActTenantOperationRequest(ctx context.Context, cmd ports.TenantOperationActionCommand) (*domain.TenantOperationRequest, error) {
	if cmd.ID == uuid.Nil || strings.TrimSpace(cmd.Action) == "" {
		return nil, domain.ErrInvalidTenantOperationAction
	}
	before, err := s.tenantOperations.GetTenantOperationRequest(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if domain.TenantOperationIsTerminal(before.Status) {
		return nil, domain.ErrInvalidTenantOperationAction
	}
	nextStatus, approvedBy, completedBy, backupConfirmed, validationResults, rollbackMetadata, metadata, err := tenantOperationActionUpdate(before, cmd)
	if err != nil {
		s.log.Warn().Err(err).Str("operation", "act tenant operation request").Str("request_id", cmd.ID.String()).Str("action", cmd.Action).Msg("invalid tenant operation action")
		return nil, err
	}
	updated, err := s.tenantOperations.UpdateTenantOperationRequestStatus(ctx, cmd.ID, nextStatus, approvedBy, completedBy, backupConfirmed, validationResults, rollbackMetadata, metadata, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	fromStatus, toStatus := before.Status, updated.Status
	if _, err := s.tenantOperations.CreateTenantOperationEvent(ctx, &domain.TenantOperationEvent{RequestID: cmd.ID, Action: strings.TrimSpace(cmd.Action), FromStatus: &fromStatus, ToStatus: &toStatus, ActorUserID: cmd.ActorID, Remarks: cleanOptionalString(cmd.Remarks), Metadata: rawJSON(map[string]any{"backup_confirmed": updated.BackupConfirmed, "risk_level": updated.RiskLevel})}, cmd.ActorID); err != nil {
		return nil, err
	}
	return updated, nil
}

func tenantOperationActionUpdate(before *domain.TenantOperationRequest, cmd ports.TenantOperationActionCommand) (string, *uuid.UUID, *uuid.UUID, *bool, json.RawMessage, json.RawMessage, json.RawMessage, error) {
	action := strings.TrimSpace(cmd.Action)
	metadata := json.RawMessage(cmd.Metadata)
	if len(metadata) == 0 {
		metadata = rawJSON(map[string]any{"action": action})
	}
	backupConfirmed := cmd.BackupConfirmed
	if before.BackupRequired && !before.BackupConfirmed && (action == domain.TenantOperationActionApprove || action == domain.TenantOperationActionStart || action == domain.TenantOperationActionComplete) {
		if backupConfirmed == nil || !*backupConfirmed {
			return "", nil, nil, nil, nil, nil, nil, domain.ErrInvalidTenantOperationAction
		}
	}
	switch action {
	case domain.TenantOperationActionValidate:
		if before.Status != domain.TenantOperationPendingValidation {
			return "", nil, nil, nil, nil, nil, nil, domain.ErrInvalidTenantOperationAction
		}
		validation := json.RawMessage(cmd.ValidationResults)
		if len(validation) == 0 {
			validation = rawJSON(map[string]any{"passed": true, "checks": []string{"request recorded", "required reason captured", "approval route pending"}})
		}
		return domain.TenantOperationPendingApproval, nil, nil, backupConfirmed, validation, nil, metadata, nil
	case domain.TenantOperationActionApprove:
		if before.Status != domain.TenantOperationPendingApproval {
			return "", nil, nil, nil, nil, nil, nil, domain.ErrInvalidTenantOperationAction
		}
		return domain.TenantOperationApproved, cmd.ActorID, nil, backupConfirmed, nil, json.RawMessage(cmd.RollbackMetadata), metadata, nil
	case domain.TenantOperationActionStart:
		if before.Status != domain.TenantOperationApproved {
			return "", nil, nil, nil, nil, nil, nil, domain.ErrInvalidTenantOperationAction
		}
		return domain.TenantOperationInProgress, nil, nil, backupConfirmed, nil, json.RawMessage(cmd.RollbackMetadata), metadata, nil
	case domain.TenantOperationActionComplete:
		if before.Status != domain.TenantOperationApproved && before.Status != domain.TenantOperationInProgress {
			return "", nil, nil, nil, nil, nil, nil, domain.ErrInvalidTenantOperationAction
		}
		return domain.TenantOperationCompleted, nil, cmd.ActorID, backupConfirmed, nil, json.RawMessage(cmd.RollbackMetadata), metadata, nil
	case domain.TenantOperationActionReject:
		if before.Status != domain.TenantOperationPendingValidation && before.Status != domain.TenantOperationPendingApproval && before.Status != domain.TenantOperationApproved {
			return "", nil, nil, nil, nil, nil, nil, domain.ErrInvalidTenantOperationAction
		}
		return domain.TenantOperationRejected, nil, cmd.ActorID, backupConfirmed, nil, json.RawMessage(cmd.RollbackMetadata), metadata, nil
	case domain.TenantOperationActionFail:
		if before.Status != domain.TenantOperationInProgress {
			return "", nil, nil, nil, nil, nil, nil, domain.ErrInvalidTenantOperationAction
		}
		return domain.TenantOperationFailed, nil, cmd.ActorID, backupConfirmed, json.RawMessage(cmd.ValidationResults), json.RawMessage(cmd.RollbackMetadata), metadata, nil
	case domain.TenantOperationActionCancel:
		return domain.TenantOperationCancelled, nil, cmd.ActorID, backupConfirmed, nil, json.RawMessage(cmd.RollbackMetadata), metadata, nil
	default:
		return "", nil, nil, nil, nil, nil, nil, domain.ErrInvalidTenantOperationAction
	}
}

func buildTenantOperationSummary(items []*domain.TenantOperationRequest) domain.TenantOperationSummary {
	summary := domain.TenantOperationSummary{ByStatus: map[string]int32{}, ByOperationType: map[string]int32{}}
	for _, item := range items {
		if item == nil {
			continue
		}
		summary.Total++
		summary.ByStatus[item.Status]++
		summary.ByOperationType[item.OperationType]++
		switch item.Status {
		case domain.TenantOperationPendingApproval:
			summary.PendingApproval++
		case domain.TenantOperationInProgress:
			summary.InProgress++
		case domain.TenantOperationCompleted:
			summary.Completed++
		}
		if item.RiskLevel == domain.WorkflowSeverityHigh || item.RiskLevel == domain.WorkflowSeverityCritical {
			summary.HighRisk++
		}
	}
	return summary
}

func parseTenantOperationOptionalTime(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil, nil
	}
	if parsed, err := time.Parse(time.RFC3339, clean); err == nil {
		return &parsed, nil
	}
	parsed, err := time.Parse("2006-01-02", clean)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func generateTenantOperationNumber() string {
	return fmt.Sprintf("TOP-%s-%s", time.Now().UTC().Format("20060102-150405"), strings.ToUpper(uuid.NewString()[:8]))
}

func boundedTenantOperationLimit(limit int32) int32 {
	if limit <= 0 {
		return 100
	}
	if limit > maxTenantOperationRequests {
		return maxTenantOperationRequests
	}
	return limit
}
