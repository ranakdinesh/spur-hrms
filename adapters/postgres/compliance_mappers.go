package postgres

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapComplianceRule(row sqlc.HrmsComplianceRule) *domain.ComplianceRule {
	return complianceRuleFromParts(
		row.ID, row.TenantID, row.Code, row.Title, ptrFromText(row.Description), row.Category, row.Scope, row.Severity,
		ptrFromText(row.ClassificationGroup), ptrFromUUID(row.WorkerTypeID), ptrFromText(row.EngagementType), ptrFromUUID(row.BranchID), ptrFromUUID(row.DepartmentID),
		row.CountryCode, ptrFromText(row.StateCode), row.TriggerEvent, row.DefaultDueDays, ptrFromInt4(row.RecurringDays), row.RequiresEvidence,
		ptrFromText(row.EvidenceLabel), ptrFromText(row.AutoDetectKey), row.BlocksPayroll, row.IsActive, ptrFromDate(row.EffectiveFrom), ptrFromDate(row.EffectiveTo),
		row.Metadata, row.Inactive, timeFromTimestamptz(row.CreatedAt), ptrFromUUID(row.CreatedBy), timeFromTimestamptz(row.UpdatedAt), ptrFromUUID(row.UpdatedBy),
		nil, nil, nil,
	)
}

func mapComplianceRules(rows []sqlc.HrmsComplianceRule) []*domain.ComplianceRule {
	items := make([]*domain.ComplianceRule, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapComplianceRule(row))
	}
	return items
}

func mapComplianceRuleListRow(row sqlc.ListComplianceRulesRow) *domain.ComplianceRule {
	return complianceRuleFromParts(
		row.ID, row.TenantID, row.Code, row.Title, ptrFromText(row.Description), row.Category, row.Scope, row.Severity,
		ptrFromText(row.ClassificationGroup), ptrFromUUID(row.WorkerTypeID), ptrFromText(row.EngagementType), ptrFromUUID(row.BranchID), ptrFromUUID(row.DepartmentID),
		row.CountryCode, ptrFromText(row.StateCode), row.TriggerEvent, row.DefaultDueDays, ptrFromInt4(row.RecurringDays), row.RequiresEvidence,
		ptrFromText(row.EvidenceLabel), ptrFromText(row.AutoDetectKey), row.BlocksPayroll, row.IsActive, ptrFromDate(row.EffectiveFrom), ptrFromDate(row.EffectiveTo),
		row.Metadata, row.Inactive, timeFromTimestamptz(row.CreatedAt), ptrFromUUID(row.CreatedBy), timeFromTimestamptz(row.UpdatedAt), ptrFromUUID(row.UpdatedBy),
		ptrFromText(row.WorkerTypeName), ptrFromText(row.BranchName), ptrFromText(row.DepartmentName),
	)
}

func mapComplianceRuleList(rows []sqlc.ListComplianceRulesRow) []*domain.ComplianceRule {
	items := make([]*domain.ComplianceRule, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapComplianceRuleListRow(row))
	}
	return items
}

func complianceRuleFromParts(id interface{}, tenantID interface{}, code string, title string, description *string, category string, scope string, severity string, classificationGroup *string, workerTypeID interface{}, engagementType *string, branchID interface{}, departmentID interface{}, countryCode string, stateCode *string, triggerEvent string, defaultDueDays int32, recurringDays *int32, requiresEvidence bool, evidenceLabel *string, autoDetectKey *string, blocksPayroll bool, isActive bool, effectiveFrom interface{}, effectiveTo interface{}, metadataBytes []byte, inactive bool, createdAt interface{}, createdBy interface{}, updatedAt interface{}, updatedBy interface{}, workerTypeName *string, branchName *string, departmentName *string) *domain.ComplianceRule {
	metadata := json.RawMessage(metadataBytes)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.ComplianceRule{
		ID:                  uuidFromAny(id),
		TenantID:            uuidFromAny(tenantID),
		Code:                code,
		Title:               title,
		Description:         description,
		Category:            category,
		Scope:               scope,
		Severity:            severity,
		ClassificationGroup: classificationGroup,
		WorkerTypeID:        uuidPtrFromAny(workerTypeID),
		EngagementType:      engagementType,
		BranchID:            uuidPtrFromAny(branchID),
		DepartmentID:        uuidPtrFromAny(departmentID),
		CountryCode:         countryCode,
		StateCode:           stateCode,
		TriggerEvent:        triggerEvent,
		DefaultDueDays:      defaultDueDays,
		RecurringDays:       recurringDays,
		RequiresEvidence:    requiresEvidence,
		EvidenceLabel:       evidenceLabel,
		AutoDetectKey:       autoDetectKey,
		BlocksPayroll:       blocksPayroll,
		IsActive:            isActive,
		EffectiveFrom:       timePtrFromAnyDate(effectiveFrom),
		EffectiveTo:         timePtrFromAnyDate(effectiveTo),
		Metadata:            metadata,
		Inactive:            inactive,
		CreatedAt:           timeFromAnyTimestamp(createdAt),
		CreatedBy:           uuidPtrFromAny(createdBy),
		UpdatedAt:           timeFromAnyTimestamp(updatedAt),
		UpdatedBy:           uuidPtrFromAny(updatedBy),
		WorkerTypeName:      workerTypeName,
		BranchName:          branchName,
		DepartmentName:      departmentName,
	}
}

func mapComplianceChecklistItem(row sqlc.HrmsComplianceChecklistItem) *domain.ComplianceChecklistItem {
	return complianceChecklistItemFromParts(row.ID, row.TenantID, row.RuleID, row.WorkerProfileID, row.EngagementID, row.Status, row.DueDate, row.CompletedAt, row.ReviewedAt, row.ReviewedBy, row.EvidencePath, row.EvidenceFileName, row.EvidenceContentType, row.EvidenceUploadedAt, row.EvidenceUploadedBy, row.WaiverReason, row.WaiverUntil, row.WaivedAt, row.WaivedBy, row.DetectedValue, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil, false, nil, false, nil, nil, nil, nil)
}

func mapComplianceChecklistItemListRow(row sqlc.ListComplianceChecklistItemsRow) *domain.ComplianceChecklistItem {
	ruleCode := row.RuleCode
	ruleTitle := row.RuleTitle
	ruleCategory := row.RuleCategory
	ruleScope := row.RuleScope
	ruleSeverity := row.RuleSeverity
	return complianceChecklistItemFromParts(row.ID, row.TenantID, row.RuleID, row.WorkerProfileID, row.EngagementID, row.Status, row.DueDate, row.CompletedAt, row.ReviewedAt, row.ReviewedBy, row.EvidencePath, row.EvidenceFileName, row.EvidenceContentType, row.EvidenceUploadedAt, row.EvidenceUploadedBy, row.WaiverReason, row.WaiverUntil, row.WaivedAt, row.WaivedBy, row.DetectedValue, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &ruleCode, &ruleTitle, &ruleCategory, &ruleScope, &ruleSeverity, row.RequiresEvidence, ptrFromText(row.EvidenceLabel), row.BlocksPayroll, ptrFromText(row.WorkerDisplayName), ptrFromText(row.WorkerCode), ptrFromText(row.EngagementTitle), ptrFromText(row.EngagementCode))
}

func mapComplianceChecklistItemList(rows []sqlc.ListComplianceChecklistItemsRow) []*domain.ComplianceChecklistItem {
	items := make([]*domain.ComplianceChecklistItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapComplianceChecklistItemListRow(row))
	}
	return items
}

func complianceChecklistItemFromParts(id interface{}, tenantID interface{}, ruleID interface{}, workerProfileID interface{}, engagementID interface{}, status string, dueDate interface{}, completedAt interface{}, reviewedAt interface{}, reviewedBy interface{}, evidencePath interface{}, evidenceFileName interface{}, evidenceContentType interface{}, evidenceUploadedAt interface{}, evidenceUploadedBy interface{}, waiverReason interface{}, waiverUntil interface{}, waivedAt interface{}, waivedBy interface{}, detectedValue interface{}, notes interface{}, metadataBytes []byte, inactive bool, createdAt interface{}, createdBy interface{}, updatedAt interface{}, updatedBy interface{}, ruleCode *string, ruleTitle *string, ruleCategory *string, ruleScope *string, ruleSeverity *string, requiresEvidence bool, evidenceLabel *string, blocksPayroll bool, workerDisplayName *string, workerCode *string, engagementTitle *string, engagementCode *string) *domain.ComplianceChecklistItem {
	metadata := json.RawMessage(metadataBytes)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.ComplianceChecklistItem{
		ID:                  uuidFromAny(id),
		TenantID:            uuidFromAny(tenantID),
		RuleID:              uuidFromAny(ruleID),
		WorkerProfileID:     uuidPtrFromAny(workerProfileID),
		EngagementID:        uuidPtrFromAny(engagementID),
		Status:              status,
		DueDate:             timePtrFromAnyDate(dueDate),
		CompletedAt:         timePtrFromAnyTimestamp(completedAt),
		ReviewedAt:          timePtrFromAnyTimestamp(reviewedAt),
		ReviewedBy:          uuidPtrFromAny(reviewedBy),
		EvidencePath:        textPtrFromAny(evidencePath),
		EvidenceFileName:    textPtrFromAny(evidenceFileName),
		EvidenceContentType: textPtrFromAny(evidenceContentType),
		EvidenceUploadedAt:  timePtrFromAnyTimestamp(evidenceUploadedAt),
		EvidenceUploadedBy:  uuidPtrFromAny(evidenceUploadedBy),
		WaiverReason:        textPtrFromAny(waiverReason),
		WaiverUntil:         timePtrFromAnyDate(waiverUntil),
		WaivedAt:            timePtrFromAnyTimestamp(waivedAt),
		WaivedBy:            uuidPtrFromAny(waivedBy),
		DetectedValue:       textPtrFromAny(detectedValue),
		Notes:               textPtrFromAny(notes),
		Metadata:            metadata,
		Inactive:            inactive,
		CreatedAt:           timeFromAnyTimestamp(createdAt),
		CreatedBy:           uuidPtrFromAny(createdBy),
		UpdatedAt:           timeFromAnyTimestamp(updatedAt),
		UpdatedBy:           uuidPtrFromAny(updatedBy),
		RuleCode:            ruleCode,
		RuleTitle:           ruleTitle,
		RuleCategory:        ruleCategory,
		RuleScope:           ruleScope,
		RuleSeverity:        ruleSeverity,
		RequiresEvidence:    requiresEvidence,
		EvidenceLabel:       evidenceLabel,
		BlocksPayroll:       blocksPayroll,
		WorkerDisplayName:   workerDisplayName,
		WorkerCode:          workerCode,
		EngagementTitle:     engagementTitle,
		EngagementCode:      engagementCode,
	}
}

func mapComplianceSummary(rows []sqlc.GetComplianceSummaryRow) []*domain.ComplianceSummaryRow {
	items := make([]*domain.ComplianceSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.ComplianceSummaryRow{Category: row.Category, Status: row.Status, ItemCount: row.ItemCount, PayrollBlockerCount: row.PayrollBlockerCount, DueSoonCount: row.DueSoonCount})
	}
	return items
}

func mapComplianceEvent(row sqlc.HrmsComplianceEvent) *domain.ComplianceEvent {
	metadata := json.RawMessage(row.Metadata)
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return &domain.ComplianceEvent{ID: row.ID, TenantID: row.TenantID, ChecklistItemID: ptrFromUUID(row.ChecklistItemID), RuleID: ptrFromUUID(row.RuleID), EventType: row.EventType, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), Comment: ptrFromText(row.Comment), ActorID: ptrFromUUID(row.ActorID), Metadata: metadata, CreatedAt: timeFromTimestamptz(row.CreatedAt)}
}

func mapComplianceEvents(rows []sqlc.HrmsComplianceEvent) []*domain.ComplianceEvent {
	items := make([]*domain.ComplianceEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapComplianceEvent(row))
	}
	return items
}

func uuidFromAny(value interface{}) uuid.UUID {
	switch typed := value.(type) {
	case uuid.UUID:
		return typed
	case *uuid.UUID:
		if typed == nil {
			return uuid.Nil
		}
		return *typed
	case pgtype.UUID:
		if !typed.Valid {
			return uuid.Nil
		}
		return uuid.UUID(typed.Bytes)
	default:
		return uuid.Nil
	}
}

func uuidPtrFromAny(value interface{}) *uuid.UUID {
	switch typed := value.(type) {
	case uuid.UUID:
		if typed == uuid.Nil {
			return nil
		}
		clean := typed
		return &clean
	case *uuid.UUID:
		if typed == nil || *typed == uuid.Nil {
			return nil
		}
		return typed
	case pgtype.UUID:
		return ptrFromUUID(typed)
	default:
		return nil
	}
}

func textPtrFromAny(value interface{}) *string {
	switch typed := value.(type) {
	case *string:
		return typed
	case string:
		if typed == "" {
			return nil
		}
		clean := typed
		return &clean
	case pgtype.Text:
		return ptrFromText(typed)
	default:
		return nil
	}
}

func timePtrFromAnyDate(value interface{}) *time.Time {
	switch typed := value.(type) {
	case *time.Time:
		return typed
	case time.Time:
		if typed.IsZero() {
			return nil
		}
		clean := typed
		return &clean
	case pgtype.Date:
		return ptrFromDate(typed)
	default:
		return nil
	}
}

func timePtrFromAnyTimestamp(value interface{}) *time.Time {
	switch typed := value.(type) {
	case *time.Time:
		return typed
	case time.Time:
		if typed.IsZero() {
			return nil
		}
		clean := typed
		return &clean
	case pgtype.Timestamptz:
		return ptrFromTimestamptz(typed)
	default:
		return nil
	}
}

func timeFromAnyTimestamp(value interface{}) time.Time {
	if parsed := timePtrFromAnyTimestamp(value); parsed != nil {
		return *parsed
	}
	return time.Time{}
}
