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

func (s *TenantService) CreateComplianceRule(ctx context.Context, cmd ports.ComplianceRuleCommand) (*domain.ComplianceRule, error) {
	item, err := complianceRuleFromCommand(cmd)
	if err != nil {
		s.logError("validate compliance rule", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.compliance.CreateComplianceRule(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create compliance rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	_, _ = s.compliance.CreateComplianceEvent(ctx, &domain.ComplianceEvent{TenantID: result.TenantID, RuleID: &result.ID, EventType: "rule_created", ToStatus: &result.Severity, ActorID: cmd.ActorID, Metadata: json.RawMessage(`{}`)})
	return result, nil
}

func (s *TenantService) UpdateComplianceRule(ctx context.Context, cmd ports.ComplianceRuleCommand) (*domain.ComplianceRule, error) {
	item, err := complianceRuleFromCommand(cmd)
	if err != nil || cmd.ID == uuid.Nil {
		if err == nil {
			err = domain.ErrInvalidComplianceRule
		}
		s.logError("validate compliance rule update", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.compliance.UpdateComplianceRule(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update compliance rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("compliance_rule_id", cmd.ID.String()))
		return nil, err
	}
	_, _ = s.compliance.CreateComplianceEvent(ctx, &domain.ComplianceEvent{TenantID: result.TenantID, RuleID: &result.ID, EventType: "rule_updated", ToStatus: &result.Severity, ActorID: cmd.ActorID, Metadata: json.RawMessage(`{}`)})
	return result, nil
}

func (s *TenantService) GetComplianceRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.ComplianceRule, error) {
	return s.compliance.GetComplianceRule(ctx, tenantID, id)
}

func (s *TenantService) ListComplianceRules(ctx context.Context, filter domain.ComplianceRuleFilter) ([]*domain.ComplianceRule, error) {
	return s.compliance.ListComplianceRules(ctx, filter)
}

func (s *TenantService) DeleteComplianceRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	err := s.compliance.DeleteComplianceRule(ctx, tenantID, id, actorID)
	if err != nil {
		s.logError("delete compliance rule", err, serviceTenantIDField(tenantID), serviceStringField("compliance_rule_id", id.String()))
		return err
	}
	ruleID := id
	_, _ = s.compliance.CreateComplianceEvent(ctx, &domain.ComplianceEvent{TenantID: tenantID, RuleID: &ruleID, EventType: "rule_deleted", ActorID: actorID, Metadata: json.RawMessage(`{}`)})
	return nil
}

func (s *TenantService) SeedDefaultComplianceRules(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.ComplianceRule, error) {
	created := make([]*domain.ComplianceRule, 0)
	for _, preset := range domain.ComplianceDefaultRules() {
		if _, err := s.compliance.GetComplianceRuleByCode(ctx, tenantID, preset.Code); err == nil {
			continue
		} else if !errors.Is(err, domain.ErrComplianceRuleNotFound) {
			s.logError("get existing default compliance rule", err, serviceTenantIDField(tenantID), serviceStringField("code", preset.Code))
			return nil, err
		}
		rule, err := domain.NewComplianceRule(domain.ComplianceRule{
			TenantID:            tenantID,
			Code:                preset.Code,
			Title:               preset.Title,
			Description:         complianceStringPtr(preset.Description),
			Category:            preset.Category,
			Scope:               preset.Scope,
			Severity:            preset.Severity,
			ClassificationGroup: complianceStringPtr(preset.ClassificationGroup),
			EngagementType:      complianceStringPtr(preset.EngagementType),
			TriggerEvent:        preset.TriggerEvent,
			DefaultDueDays:      preset.DefaultDueDays,
			RecurringDays:       preset.RecurringDays,
			RequiresEvidence:    preset.RequiresEvidence,
			EvidenceLabel:       complianceStringPtr(preset.EvidenceLabel),
			AutoDetectKey:       complianceStringPtr(preset.AutoDetectKey),
			BlocksPayroll:       preset.BlocksPayroll,
			IsActive:            true,
			CountryCode:         "IN",
			Metadata:            json.RawMessage(`{"preset":"india_default"}`),
		})
		if err != nil {
			s.logError("validate default compliance rule", err, serviceTenantIDField(tenantID), serviceStringField("code", preset.Code))
			return nil, err
		}
		result, err := s.compliance.CreateComplianceRule(ctx, rule, actorID)
		if err != nil {
			s.logError("seed default compliance rule", err, serviceTenantIDField(tenantID), serviceStringField("code", preset.Code))
			return nil, err
		}
		created = append(created, result)
		_, _ = s.compliance.CreateComplianceEvent(ctx, &domain.ComplianceEvent{TenantID: tenantID, RuleID: &result.ID, EventType: "rule_seeded", ActorID: actorID, Metadata: json.RawMessage(`{"preset":"india_default"}`)})
	}
	return created, nil
}

func (s *TenantService) GenerateComplianceChecklist(ctx context.Context, cmd ports.ComplianceChecklistGenerateCommand) ([]*domain.ComplianceChecklistItem, error) {
	if cmd.TenantID == uuid.Nil || ((cmd.WorkerProfileID == nil || *cmd.WorkerProfileID == uuid.Nil) && (cmd.EngagementID == nil || *cmd.EngagementID == uuid.Nil)) {
		s.logError("validate compliance checklist generation", domain.ErrInvalidComplianceChecklistItem, serviceTenantIDField(cmd.TenantID))
		return nil, domain.ErrInvalidComplianceChecklistItem
	}
	now := time.Now().UTC()
	var rules []*domain.ComplianceRule
	var err error
	var workerID *uuid.UUID
	var engagementID *uuid.UUID
	if cmd.EngagementID != nil && *cmd.EngagementID != uuid.Nil {
		engagementID = cmd.EngagementID
		rules, err = s.compliance.ListActiveComplianceRulesForEngagement(ctx, cmd.TenantID, *cmd.EngagementID)
	} else {
		workerID = cmd.WorkerProfileID
		rules, err = s.compliance.ListActiveComplianceRulesForWorker(ctx, cmd.TenantID, *cmd.WorkerProfileID)
	}
	if err != nil {
		s.logError("list applicable compliance rules", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	created := make([]*domain.ComplianceChecklistItem, 0, len(rules))
	for _, rule := range rules {
		item, err := domain.NewComplianceChecklistItem(domain.ComplianceChecklistItem{
			TenantID:        cmd.TenantID,
			RuleID:          rule.ID,
			WorkerProfileID: workerID,
			EngagementID:    engagementID,
			Status:          domain.ComplianceStatusPending,
			DueDate:         domain.ComplianceDueDateFrom(rule, now),
			Metadata:        json.RawMessage(`{"generated":true}`),
		})
		if err != nil {
			s.logError("validate generated compliance checklist item", err, serviceTenantIDField(cmd.TenantID), serviceStringField("compliance_rule_id", rule.ID.String()))
			return nil, err
		}
		result, err := s.compliance.CreateComplianceChecklistItem(ctx, item, cmd.ActorID)
		if err != nil {
			s.logError("create generated compliance checklist item", err, serviceTenantIDField(cmd.TenantID), serviceStringField("compliance_rule_id", rule.ID.String()))
			return nil, err
		}
		if result == nil {
			continue
		}
		created = append(created, result)
		_, _ = s.compliance.CreateComplianceEvent(ctx, &domain.ComplianceEvent{TenantID: cmd.TenantID, ChecklistItemID: &result.ID, RuleID: &rule.ID, EventType: "checklist_generated", ToStatus: &result.Status, ActorID: cmd.ActorID, Metadata: json.RawMessage(`{}`)})
	}
	return created, nil
}

func (s *TenantService) ListComplianceChecklistItems(ctx context.Context, filter domain.ComplianceChecklistFilter) ([]*domain.ComplianceChecklistItem, error) {
	if err := s.compliance.RefreshComplianceChecklistDueStatus(ctx, filter.TenantID); err != nil {
		s.logError("refresh compliance checklist due statuses", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	return s.compliance.ListComplianceChecklistItems(ctx, filter)
}

func (s *TenantService) UpdateComplianceChecklistStatus(ctx context.Context, cmd ports.ComplianceChecklistStatusCommand) (*domain.ComplianceChecklistItem, error) {
	status := complianceStatus(cmd.Status)
	if status == "" {
		s.logError("validate compliance checklist status", domain.ErrInvalidComplianceChecklistItem, serviceTenantIDField(cmd.TenantID), serviceStringField("status", cmd.Status))
		return nil, domain.ErrInvalidComplianceChecklistItem
	}
	before, _ := s.compliance.GetComplianceChecklistItem(ctx, cmd.TenantID, cmd.ID)
	result, err := s.compliance.UpdateComplianceChecklistStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.Notes, cmd.ActorID)
	if err != nil {
		s.logError("update compliance checklist status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("checklist_item_id", cmd.ID.String()))
		return nil, err
	}
	from := ""
	if before != nil {
		from = before.Status
	}
	_, _ = s.compliance.CreateComplianceEvent(ctx, &domain.ComplianceEvent{TenantID: cmd.TenantID, ChecklistItemID: &result.ID, RuleID: &result.RuleID, EventType: "status_changed", FromStatus: complianceStringPtr(from), ToStatus: &result.Status, Comment: cmd.Notes, ActorID: cmd.ActorID, Metadata: json.RawMessage(`{}`)})
	return result, nil
}

func (s *TenantService) UpdateComplianceChecklistEvidence(ctx context.Context, cmd ports.ComplianceEvidenceCommand) (*domain.ComplianceChecklistItem, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || (cmd.EvidencePath == nil && cmd.EvidenceFileName == nil) {
		s.logError("validate compliance evidence", domain.ErrInvalidComplianceChecklistItem, serviceTenantIDField(cmd.TenantID))
		return nil, domain.ErrInvalidComplianceChecklistItem
	}
	before, _ := s.compliance.GetComplianceChecklistItem(ctx, cmd.TenantID, cmd.ID)
	result, err := s.compliance.UpdateComplianceChecklistEvidence(ctx, cmd.TenantID, cmd.ID, cmd.EvidencePath, cmd.EvidenceFileName, cmd.EvidenceContentType, cmd.Notes, cmd.ActorID)
	if err != nil {
		s.logError("update compliance evidence", err, serviceTenantIDField(cmd.TenantID), serviceStringField("checklist_item_id", cmd.ID.String()))
		return nil, err
	}
	from := ""
	if before != nil {
		from = before.Status
	}
	_, _ = s.compliance.CreateComplianceEvent(ctx, &domain.ComplianceEvent{TenantID: cmd.TenantID, ChecklistItemID: &result.ID, RuleID: &result.RuleID, EventType: "evidence_uploaded", FromStatus: complianceStringPtr(from), ToStatus: &result.Status, Comment: cmd.Notes, ActorID: cmd.ActorID, Metadata: json.RawMessage(`{}`)})
	return result, nil
}

func (s *TenantService) WaiveComplianceChecklistItem(ctx context.Context, cmd ports.ComplianceWaiverCommand) (*domain.ComplianceChecklistItem, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || cmd.WaiverReason == "" {
		s.logError("validate compliance waiver", domain.ErrInvalidComplianceChecklistItem, serviceTenantIDField(cmd.TenantID))
		return nil, domain.ErrInvalidComplianceChecklistItem
	}
	before, _ := s.compliance.GetComplianceChecklistItem(ctx, cmd.TenantID, cmd.ID)
	result, err := s.compliance.WaiveComplianceChecklistItem(ctx, cmd.TenantID, cmd.ID, cmd.WaiverReason, cmd.WaiverUntil, cmd.Notes, cmd.ActorID)
	if err != nil {
		s.logError("waive compliance checklist item", err, serviceTenantIDField(cmd.TenantID), serviceStringField("checklist_item_id", cmd.ID.String()))
		return nil, err
	}
	from := ""
	if before != nil {
		from = before.Status
	}
	_, _ = s.compliance.CreateComplianceEvent(ctx, &domain.ComplianceEvent{TenantID: cmd.TenantID, ChecklistItemID: &result.ID, RuleID: &result.RuleID, EventType: "waived", FromStatus: complianceStringPtr(from), ToStatus: &result.Status, Comment: &cmd.WaiverReason, ActorID: cmd.ActorID, Metadata: json.RawMessage(`{}`)})
	return result, nil
}

func (s *TenantService) DeleteComplianceChecklistItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	err := s.compliance.DeleteComplianceChecklistItem(ctx, tenantID, id, actorID)
	if err != nil {
		s.logError("delete compliance checklist item", err, serviceTenantIDField(tenantID), serviceStringField("checklist_item_id", id.String()))
		return err
	}
	itemID := id
	_, _ = s.compliance.CreateComplianceEvent(ctx, &domain.ComplianceEvent{TenantID: tenantID, ChecklistItemID: &itemID, EventType: "checklist_deleted", ActorID: actorID, Metadata: json.RawMessage(`{}`)})
	return nil
}

func (s *TenantService) GetComplianceSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.ComplianceSummaryRow, error) {
	if err := s.compliance.RefreshComplianceChecklistDueStatus(ctx, tenantID); err != nil {
		s.logError("refresh compliance summary due statuses", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return s.compliance.GetComplianceSummary(ctx, tenantID)
}

func (s *TenantService) ListComplianceEvents(ctx context.Context, tenantID uuid.UUID, checklistItemID *uuid.UUID, ruleID *uuid.UUID) ([]*domain.ComplianceEvent, error) {
	return s.compliance.ListComplianceEvents(ctx, tenantID, checklistItemID, ruleID)
}

func complianceRuleFromCommand(cmd ports.ComplianceRuleCommand) (*domain.ComplianceRule, error) {
	return domain.NewComplianceRule(domain.ComplianceRule{TenantID: cmd.TenantID, Code: cmd.Code, Title: cmd.Title, Description: cmd.Description, Category: cmd.Category, Scope: cmd.Scope, Severity: cmd.Severity, ClassificationGroup: cmd.ClassificationGroup, WorkerTypeID: cmd.WorkerTypeID, EngagementType: cmd.EngagementType, BranchID: cmd.BranchID, DepartmentID: cmd.DepartmentID, CountryCode: cmd.CountryCode, StateCode: cmd.StateCode, TriggerEvent: cmd.TriggerEvent, DefaultDueDays: cmd.DefaultDueDays, RecurringDays: cmd.RecurringDays, RequiresEvidence: cmd.RequiresEvidence, EvidenceLabel: cmd.EvidenceLabel, AutoDetectKey: cmd.AutoDetectKey, BlocksPayroll: cmd.BlocksPayroll, IsActive: cmd.IsActive, EffectiveFrom: cmd.EffectiveFrom, EffectiveTo: cmd.EffectiveTo, Metadata: cmd.Metadata})
}

func complianceStatus(value string) string {
	status := domain.ComplianceStatusPending
	if value != "" {
		status = value
	}
	switch status {
	case domain.ComplianceStatusPending, domain.ComplianceStatusInReview, domain.ComplianceStatusCompliant, domain.ComplianceStatusNonCompliant, domain.ComplianceStatusWaived, domain.ComplianceStatusExpired, domain.ComplianceStatusNotApplicable:
		return status
	default:
		return ""
	}
}

func complianceStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
