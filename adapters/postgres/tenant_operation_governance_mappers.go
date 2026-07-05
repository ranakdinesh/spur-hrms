package postgres

import (
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapTenantOperationRequest(row sqlc.HrmsTenantOperationRequest) *domain.TenantOperationRequest {
	return &domain.TenantOperationRequest{
		ID:                row.ID,
		OperationNumber:   row.OperationNumber,
		OperationType:     row.OperationType,
		Title:             row.Title,
		TargetTenantID:    ptrFromUUID(row.TargetTenantID),
		TargetTenantName:  ptrFromText(row.TargetTenantName),
		TargetTenantCode:  ptrFromText(row.TargetTenantCode),
		Status:            row.Status,
		RiskLevel:         row.RiskLevel,
		Reason:            row.Reason,
		RequestedBy:       ptrFromUUID(row.RequestedBy),
		ApprovedBy:        ptrFromUUID(row.ApprovedBy),
		ApprovedAt:        ptrFromTimestamptz(row.ApprovedAt),
		CompletedBy:       ptrFromUUID(row.CompletedBy),
		CompletedAt:       ptrFromTimestamptz(row.CompletedAt),
		ApprovalRequired:  row.ApprovalRequired,
		BackupRequired:    row.BackupRequired,
		BackupConfirmed:   row.BackupConfirmed,
		RetentionUntil:    ptrFromTimestamptz(row.RetentionUntil),
		RequestPayload:    jsonRawDefault(row.RequestPayload, `{}`),
		ValidationResults: jsonRawDefault(row.ValidationResults, `{}`),
		RollbackMetadata:  jsonRawDefault(row.RollbackMetadata, `{}`),
		Metadata:          jsonRawDefault(row.Metadata, `{}`),
		Inactive:          row.Inactive,
		CreatedAt:         timeFromTimestamptz(row.CreatedAt),
		CreatedBy:         ptrFromUUID(row.CreatedBy),
		UpdatedAt:         timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:         ptrFromUUID(row.UpdatedBy),
	}
}

func mapTenantOperationRequests(rows []sqlc.HrmsTenantOperationRequest) []*domain.TenantOperationRequest {
	items := make([]*domain.TenantOperationRequest, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapTenantOperationRequest(row))
	}
	return items
}

func mapTenantOperationEvent(row sqlc.HrmsTenantOperationEvent) *domain.TenantOperationEvent {
	return &domain.TenantOperationEvent{
		ID:          row.ID,
		RequestID:   row.RequestID,
		Action:      row.Action,
		FromStatus:  ptrFromText(row.FromStatus),
		ToStatus:    ptrFromText(row.ToStatus),
		ActorUserID: ptrFromUUID(row.ActorUserID),
		Remarks:     ptrFromText(row.Remarks),
		Metadata:    jsonRawDefault(row.Metadata, `{}`),
		CreatedAt:   timeFromTimestamptz(row.CreatedAt),
		CreatedBy:   ptrFromUUID(row.CreatedBy),
		UpdatedAt:   timeFromTimestamptz(row.UpdatedAt),
		UpdatedBy:   ptrFromUUID(row.UpdatedBy),
	}
}

func mapTenantOperationEvents(rows []sqlc.HrmsTenantOperationEvent) []*domain.TenantOperationEvent {
	items := make([]*domain.TenantOperationEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapTenantOperationEvent(row))
	}
	return items
}
