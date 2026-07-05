package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateComplianceRule(ctx context.Context, item *domain.ComplianceRule, actorID *uuid.UUID) (*domain.ComplianceRule, error) {
	row, err := s.getQueries(ctx).CreateComplianceRule(ctx, complianceRuleCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create compliance rule", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapComplianceRule(row), nil
}

func (s *Store) UpdateComplianceRule(ctx context.Context, item *domain.ComplianceRule, actorID *uuid.UUID) (*domain.ComplianceRule, error) {
	row, err := s.getQueries(ctx).UpdateComplianceRule(ctx, complianceRuleUpdateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "update compliance rule", err, tenantIDField(item.TenantID), stringField("compliance_rule_id", item.ID.String()))
	}
	return mapComplianceRule(row), nil
}

func (s *Store) GetComplianceRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ComplianceRule, error) {
	row, err := s.getQueries(ctx).GetComplianceRule(ctx, sqlc.GetComplianceRuleParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrComplianceRuleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get compliance rule", err, tenantIDField(tenantID), stringField("compliance_rule_id", id.String()))
	}
	return mapComplianceRule(row), nil
}

func (s *Store) GetComplianceRuleByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.ComplianceRule, error) {
	row, err := s.getQueries(ctx).GetComplianceRuleByCode(ctx, sqlc.GetComplianceRuleByCodeParams{TenantID: tenantID, Lower: code})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrComplianceRuleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get compliance rule by code", err, tenantIDField(tenantID), stringField("code", code))
	}
	return mapComplianceRule(row), nil
}

func (s *Store) ListComplianceRules(ctx context.Context, filter domain.ComplianceRuleFilter) ([]*domain.ComplianceRule, error) {
	rows, err := s.getQueries(ctx).ListComplianceRules(ctx, sqlc.ListComplianceRulesParams{TenantID: filter.TenantID, Category: textFromPtr(filter.Category), Scope: textFromPtr(filter.Scope), Severity: textFromPtr(filter.Severity), IsActive: complianceBoolFromPtr(filter.IsActive), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list compliance rules", err, tenantIDField(filter.TenantID))
	}
	return mapComplianceRuleList(rows), nil
}

func (s *Store) ListActiveComplianceRulesForWorker(ctx context.Context, tenantID uuid.UUID, workerProfileID uuid.UUID) ([]*domain.ComplianceRule, error) {
	rows, err := s.getQueries(ctx).ListActiveComplianceRulesForWorker(ctx, sqlc.ListActiveComplianceRulesForWorkerParams{TenantID: tenantID, ID: workerProfileID})
	if err != nil {
		return nil, s.logDBError(ctx, "list active compliance rules for worker", err, tenantIDField(tenantID), stringField("worker_profile_id", workerProfileID.String()))
	}
	return mapComplianceRules(rows), nil
}

func (s *Store) ListActiveComplianceRulesForEngagement(ctx context.Context, tenantID uuid.UUID, engagementID uuid.UUID) ([]*domain.ComplianceRule, error) {
	rows, err := s.getQueries(ctx).ListActiveComplianceRulesForEngagement(ctx, sqlc.ListActiveComplianceRulesForEngagementParams{TenantID: tenantID, ID: engagementID})
	if err != nil {
		return nil, s.logDBError(ctx, "list active compliance rules for engagement", err, tenantIDField(tenantID), stringField("engagement_id", engagementID.String()))
	}
	return mapComplianceRules(rows), nil
}

func (s *Store) DeleteComplianceRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteComplianceRule(ctx, sqlc.SoftDeleteComplianceRuleParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete compliance rule", err, tenantIDField(tenantID), stringField("compliance_rule_id", id.String()))
	}
	return nil
}

func (s *Store) CreateComplianceChecklistItem(ctx context.Context, item *domain.ComplianceChecklistItem, actorID *uuid.UUID) (*domain.ComplianceChecklistItem, error) {
	row, err := s.getQueries(ctx).CreateComplianceChecklistItem(ctx, sqlc.CreateComplianceChecklistItemParams{
		TenantID:            item.TenantID,
		RuleID:              item.RuleID,
		WorkerProfileID:     uuidFromPtr(item.WorkerProfileID),
		EngagementID:        uuidFromPtr(item.EngagementID),
		Status:              item.Status,
		DueDate:             dateFromPtr(item.DueDate),
		EvidencePath:        textFromPtr(item.EvidencePath),
		EvidenceFileName:    textFromPtr(item.EvidenceFileName),
		EvidenceContentType: textFromPtr(item.EvidenceContentType),
		DetectedValue:       textFromPtr(item.DetectedValue),
		Notes:               textFromPtr(item.Notes),
		Metadata:            jsonBytesFromRaw(item.Metadata),
		CreatedBy:           uuidFromPtr(actorID),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, s.logDBError(ctx, "create compliance checklist item", err, tenantIDField(item.TenantID), stringField("compliance_rule_id", item.RuleID.String()))
	}
	return mapComplianceChecklistItem(row), nil
}

func (s *Store) GetComplianceChecklistItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ComplianceChecklistItem, error) {
	row, err := s.getQueries(ctx).GetComplianceChecklistItem(ctx, sqlc.GetComplianceChecklistItemParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrComplianceChecklistItemNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get compliance checklist item", err, tenantIDField(tenantID), stringField("checklist_item_id", id.String()))
	}
	return mapComplianceChecklistItem(row), nil
}

func (s *Store) ListComplianceChecklistItems(ctx context.Context, filter domain.ComplianceChecklistFilter) ([]*domain.ComplianceChecklistItem, error) {
	rows, err := s.getQueries(ctx).ListComplianceChecklistItems(ctx, sqlc.ListComplianceChecklistItemsParams{TenantID: filter.TenantID, WorkerProfileID: uuidFromPtr(filter.WorkerProfileID), EngagementID: uuidFromPtr(filter.EngagementID), RuleID: uuidFromPtr(filter.RuleID), Status: textFromPtr(filter.Status), Category: textFromPtr(filter.Category), DueBefore: dateFromPtr(filter.DueBefore), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list compliance checklist items", err, tenantIDField(filter.TenantID))
	}
	return mapComplianceChecklistItemList(rows), nil
}

func (s *Store) UpdateComplianceChecklistStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, notes *string, actorID *uuid.UUID) (*domain.ComplianceChecklistItem, error) {
	row, err := s.getQueries(ctx).UpdateComplianceChecklistStatus(ctx, sqlc.UpdateComplianceChecklistStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID), Notes: textFromPtr(notes)})
	if err != nil {
		return nil, s.logDBError(ctx, "update compliance checklist status", err, tenantIDField(tenantID), stringField("checklist_item_id", id.String()), stringField("status", status))
	}
	return mapComplianceChecklistItem(row), nil
}

func (s *Store) UpdateComplianceChecklistEvidence(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, evidencePath *string, evidenceFileName *string, evidenceContentType *string, notes *string, actorID *uuid.UUID) (*domain.ComplianceChecklistItem, error) {
	row, err := s.getQueries(ctx).UpdateComplianceChecklistEvidence(ctx, sqlc.UpdateComplianceChecklistEvidenceParams{TenantID: tenantID, ID: id, EvidencePath: textFromPtr(evidencePath), EvidenceFileName: textFromPtr(evidenceFileName), EvidenceContentType: textFromPtr(evidenceContentType), EvidenceUploadedBy: uuidFromPtr(actorID), Notes: textFromPtr(notes)})
	if err != nil {
		return nil, s.logDBError(ctx, "update compliance checklist evidence", err, tenantIDField(tenantID), stringField("checklist_item_id", id.String()))
	}
	return mapComplianceChecklistItem(row), nil
}

func (s *Store) WaiveComplianceChecklistItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, waiverReason string, waiverUntil *time.Time, notes *string, actorID *uuid.UUID) (*domain.ComplianceChecklistItem, error) {
	row, err := s.getQueries(ctx).WaiveComplianceChecklistItem(ctx, sqlc.WaiveComplianceChecklistItemParams{TenantID: tenantID, ID: id, WaiverReason: textFromString(waiverReason), WaiverUntil: dateFromPtr(waiverUntil), WaivedBy: uuidFromPtr(actorID), Notes: textFromPtr(notes)})
	if err != nil {
		return nil, s.logDBError(ctx, "waive compliance checklist item", err, tenantIDField(tenantID), stringField("checklist_item_id", id.String()))
	}
	return mapComplianceChecklistItem(row), nil
}

func (s *Store) RefreshComplianceChecklistDueStatus(ctx context.Context, tenantID uuid.UUID) error {
	if err := s.getQueries(ctx).RefreshComplianceChecklistDueStatus(ctx, tenantID); err != nil {
		return s.logDBError(ctx, "refresh compliance checklist due status", err, tenantIDField(tenantID))
	}
	return nil
}

func (s *Store) DeleteComplianceChecklistItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteComplianceChecklistItem(ctx, sqlc.SoftDeleteComplianceChecklistItemParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete compliance checklist item", err, tenantIDField(tenantID), stringField("checklist_item_id", id.String()))
	}
	return nil
}

func (s *Store) GetComplianceSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.ComplianceSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetComplianceSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get compliance summary", err, tenantIDField(tenantID))
	}
	return mapComplianceSummary(rows), nil
}

func (s *Store) CreateComplianceEvent(ctx context.Context, item *domain.ComplianceEvent) (*domain.ComplianceEvent, error) {
	row, err := s.getQueries(ctx).CreateComplianceEvent(ctx, sqlc.CreateComplianceEventParams{TenantID: item.TenantID, ChecklistItemID: uuidFromPtr(item.ChecklistItemID), RuleID: uuidFromPtr(item.RuleID), EventType: item.EventType, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), Comment: textFromPtr(item.Comment), ActorID: uuidFromPtr(item.ActorID), Metadata: jsonBytesFromRaw(item.Metadata)})
	if err != nil {
		return nil, s.logDBError(ctx, "create compliance event", err, tenantIDField(item.TenantID), stringField("event_type", item.EventType))
	}
	return mapComplianceEvent(row), nil
}

func (s *Store) ListComplianceEvents(ctx context.Context, tenantID uuid.UUID, checklistItemID *uuid.UUID, ruleID *uuid.UUID) ([]*domain.ComplianceEvent, error) {
	rows, err := s.getQueries(ctx).ListComplianceEvents(ctx, sqlc.ListComplianceEventsParams{TenantID: tenantID, ChecklistItemID: uuidFromPtr(checklistItemID), RuleID: uuidFromPtr(ruleID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list compliance events", err, tenantIDField(tenantID))
	}
	return mapComplianceEvents(rows), nil
}

func complianceRuleCreateParams(item *domain.ComplianceRule, actorID *uuid.UUID) sqlc.CreateComplianceRuleParams {
	return sqlc.CreateComplianceRuleParams{
		TenantID:            item.TenantID,
		Code:                item.Code,
		Title:               item.Title,
		Description:         textFromPtr(item.Description),
		Category:            item.Category,
		Scope:               item.Scope,
		Severity:            item.Severity,
		ClassificationGroup: textFromPtr(item.ClassificationGroup),
		WorkerTypeID:        uuidFromPtr(item.WorkerTypeID),
		EngagementType:      textFromPtr(item.EngagementType),
		BranchID:            uuidFromPtr(item.BranchID),
		DepartmentID:        uuidFromPtr(item.DepartmentID),
		CountryCode:         item.CountryCode,
		StateCode:           textFromPtr(item.StateCode),
		TriggerEvent:        item.TriggerEvent,
		DefaultDueDays:      item.DefaultDueDays,
		RecurringDays:       int4FromPtr(item.RecurringDays),
		RequiresEvidence:    item.RequiresEvidence,
		EvidenceLabel:       textFromPtr(item.EvidenceLabel),
		AutoDetectKey:       textFromPtr(item.AutoDetectKey),
		BlocksPayroll:       item.BlocksPayroll,
		IsActive:            item.IsActive,
		EffectiveFrom:       dateFromPtr(item.EffectiveFrom),
		EffectiveTo:         dateFromPtr(item.EffectiveTo),
		Metadata:            jsonBytesFromRaw(item.Metadata),
		CreatedBy:           uuidFromPtr(actorID),
	}
}

func complianceRuleUpdateParams(item *domain.ComplianceRule, actorID *uuid.UUID) sqlc.UpdateComplianceRuleParams {
	return sqlc.UpdateComplianceRuleParams{
		TenantID:            item.TenantID,
		ID:                  item.ID,
		Code:                item.Code,
		Title:               item.Title,
		Description:         textFromPtr(item.Description),
		Category:            item.Category,
		Scope:               item.Scope,
		Severity:            item.Severity,
		ClassificationGroup: textFromPtr(item.ClassificationGroup),
		WorkerTypeID:        uuidFromPtr(item.WorkerTypeID),
		EngagementType:      textFromPtr(item.EngagementType),
		BranchID:            uuidFromPtr(item.BranchID),
		DepartmentID:        uuidFromPtr(item.DepartmentID),
		CountryCode:         item.CountryCode,
		StateCode:           textFromPtr(item.StateCode),
		TriggerEvent:        item.TriggerEvent,
		DefaultDueDays:      item.DefaultDueDays,
		RecurringDays:       int4FromPtr(item.RecurringDays),
		RequiresEvidence:    item.RequiresEvidence,
		EvidenceLabel:       textFromPtr(item.EvidenceLabel),
		AutoDetectKey:       textFromPtr(item.AutoDetectKey),
		BlocksPayroll:       item.BlocksPayroll,
		IsActive:            item.IsActive,
		EffectiveFrom:       dateFromPtr(item.EffectiveFrom),
		EffectiveTo:         dateFromPtr(item.EffectiveTo),
		Metadata:            jsonBytesFromRaw(item.Metadata),
		UpdatedBy:           uuidFromPtr(actorID),
	}
}

func complianceBoolFromPtr(value *bool) pgtype.Bool {
	if value == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *value, Valid: true}
}
