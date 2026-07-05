package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateWorkerType(ctx context.Context, cmd ports.WorkerTypeCommand) (*domain.WorkerType, error) {
	item, err := domain.NewWorkerType(workerTypeInputFromCommand(cmd, false))
	if err != nil {
		s.logError("validate worker type create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_type_code", cmd.Code))
		return nil, err
	}
	result, err := s.workerTypes.CreateWorkerType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create worker type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_type_code", item.Code))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("worker_type_id", result.ID.String()).Str("worker_type_code", result.Code).Msg("hrms: worker type created")
	return result, nil
}

func (s *TenantService) ListWorkerTypes(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.WorkerType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker type list tenant", err)
		return nil, err
	}
	result, err := s.workerTypes.ListWorkerTypes(ctx, tenantID)
	if err != nil {
		s.logError("list worker types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if len(result) == 0 {
		s.log.Warn().Str("tenant_id", tenantID.String()).Msg("hrms: worker type taxonomy missing, seeding defaults")
		result, err = s.seedDefaultWorkerTypes(ctx, tenantID, actorID)
		if err != nil {
			return nil, err
		}
	}
	if err := s.ensureWorkerClassificationRules(ctx, tenantID, result, actorID); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetWorkerType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker type get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkerTypeID
		s.logError("validate worker type get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.workerTypes.GetWorkerType(ctx, tenantID, id)
	if err != nil {
		s.logError("get worker type", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetWorkerTypeByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.WorkerType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker type by code tenant", err)
		return nil, err
	}
	clean := strings.ToLower(strings.TrimSpace(code))
	if clean == "" {
		err := domain.ErrInvalidWorkerTypeCode
		s.logError("validate worker type code", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.workerTypes.GetWorkerTypeByCode(ctx, tenantID, clean)
	if err != nil {
		s.logError("get worker type by code", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_code", clean))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateWorkerType(ctx context.Context, cmd ports.WorkerTypeCommand) (*domain.WorkerType, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidWorkerTypeID
		s.logError("validate worker type update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetWorkerType(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	cmd.Code = existing.Code
	item, err := domain.NewWorkerType(workerTypeInputFromCommand(cmd, existing.IsSystemDefault))
	if err != nil {
		s.logError("validate worker type update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_type_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.workerTypes.UpdateWorkerType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update worker type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_type_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("worker_type_id", result.ID.String()).Str("worker_type_code", result.Code).Msg("hrms: worker type updated")
	return result, nil
}

func (s *TenantService) DeleteWorkerType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker type delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkerTypeID
		s.logError("validate worker type delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.workerTypes.DeleteWorkerType(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete worker type", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_id", id.String()))
		return err
	}
	s.log.Info().Str("tenant_id", tenantID.String()).Str("worker_type_id", id.String()).Msg("hrms: worker type deactivated")
	return nil
}

func (s *TenantService) CreateWorkerClassificationRule(ctx context.Context, cmd ports.WorkerClassificationRuleCommand) (*domain.WorkerClassificationRule, error) {
	if _, err := s.GetWorkerType(ctx, cmd.TenantID, cmd.WorkerTypeID); err != nil {
		return nil, err
	}
	item, err := domain.NewWorkerClassificationRule(workerClassificationRuleInputFromCommand(cmd))
	if err != nil {
		s.logError("validate worker classification rule create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_type_id", cmd.WorkerTypeID.String()))
		return nil, err
	}
	result, err := s.workerTypes.CreateWorkerClassificationRule(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create worker classification rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("worker_type_id", cmd.WorkerTypeID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListWorkerClassificationRules(ctx context.Context, tenantID uuid.UUID, workerTypeID *uuid.UUID) ([]*domain.WorkerClassificationRule, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker classification rule list tenant", err)
		return nil, err
	}
	if workerTypeID != nil && *workerTypeID == uuid.Nil {
		err := domain.ErrInvalidWorkerTypeID
		s.logError("validate worker classification rule list worker type", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.workerTypes.ListWorkerClassificationRules(ctx, tenantID, workerTypeID)
	if err != nil {
		s.logError("list worker classification rules", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetWorkerClassificationRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerClassificationRule, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker classification rule get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkerClassificationRuleID
		s.logError("validate worker classification rule get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.workerTypes.GetWorkerClassificationRule(ctx, tenantID, id)
	if err != nil {
		s.logError("get worker classification rule", err, serviceTenantIDField(tenantID), serviceStringField("rule_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateWorkerClassificationRule(ctx context.Context, cmd ports.WorkerClassificationRuleCommand) (*domain.WorkerClassificationRule, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidWorkerClassificationRuleID
		s.logError("validate worker classification rule update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetWorkerType(ctx, cmd.TenantID, cmd.WorkerTypeID); err != nil {
		return nil, err
	}
	item, err := domain.NewWorkerClassificationRule(workerClassificationRuleInputFromCommand(cmd))
	if err != nil {
		s.logError("validate worker classification rule update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("rule_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.workerTypes.UpdateWorkerClassificationRule(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update worker classification rule", err, serviceTenantIDField(cmd.TenantID), serviceStringField("rule_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteWorkerClassificationRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate worker classification rule delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkerClassificationRuleID
		s.logError("validate worker classification rule delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.workerTypes.DeleteWorkerClassificationRule(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete worker classification rule", err, serviceTenantIDField(tenantID), serviceStringField("rule_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) seedDefaultWorkerTypes(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.WorkerType, error) {
	for _, input := range domain.DefaultWorkerTypeInputs(tenantID) {
		item, err := domain.NewWorkerType(input)
		if err != nil {
			s.logError("validate default worker type", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_code", input.Code))
			return nil, err
		}
		created, err := s.workerTypes.CreateWorkerType(ctx, item, actorID)
		if err != nil {
			s.logError("seed worker type", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_code", item.Code))
			return nil, err
		}
		rule, err := domain.NewWorkerClassificationRule(domain.DefaultWorkerClassificationRuleInput(tenantID, created))
		if err != nil {
			s.logError("validate default worker classification rule", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_code", item.Code))
			return nil, err
		}
		if _, err = s.workerTypes.CreateWorkerClassificationRule(ctx, rule, actorID); err != nil {
			s.logError("seed worker classification rule", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_id", created.ID.String()))
			return nil, err
		}
	}
	result, err := s.workerTypes.ListWorkerTypes(ctx, tenantID)
	if err != nil {
		s.logError("list seeded worker types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ensureWorkerClassificationRules(ctx context.Context, tenantID uuid.UUID, items []*domain.WorkerType, actorID *uuid.UUID) error {
	for _, item := range items {
		count, err := s.workerTypes.CountWorkerClassificationRules(ctx, tenantID, item.ID)
		if err != nil {
			s.logError("count worker classification rules", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_id", item.ID.String()))
			return err
		}
		if count > 0 {
			continue
		}
		rule, err := domain.NewWorkerClassificationRule(domain.DefaultWorkerClassificationRuleInput(tenantID, item))
		if err != nil {
			s.logError("validate missing worker classification rule", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_id", item.ID.String()))
			return err
		}
		if _, err = s.workerTypes.CreateWorkerClassificationRule(ctx, rule, actorID); err != nil {
			s.logError("create missing worker classification rule", err, serviceTenantIDField(tenantID), serviceStringField("worker_type_id", item.ID.String()))
			return err
		}
	}
	return nil
}

func workerTypeInputFromCommand(cmd ports.WorkerTypeCommand, isSystemDefault bool) domain.WorkerTypeInput {
	return domain.WorkerTypeInput{
		TenantID:            cmd.TenantID,
		Code:                cmd.Code,
		Name:                cmd.Name,
		ClassificationGroup: cmd.ClassificationGroup,
		Description:         cmd.Description,
		AttendanceMode:      cmd.AttendanceMode,
		PayMode:             cmd.PayMode,
		TDSSection:          cmd.TDSSection,
		PFApplicable:        cmd.PFApplicable,
		ESICApplicable:      cmd.ESICApplicable,
		PTApplicable:        cmd.PTApplicable,
		LWFApplicable:       cmd.LWFApplicable,
		CLRAApplicable:      cmd.CLRAApplicable,
		LeaveApplicable:     cmd.LeaveApplicable,
		OvertimeApplicable:  cmd.OvertimeApplicable,
		RequiresAgreement:   cmd.RequiresAgreement,
		RequiresInvoice:     cmd.RequiresInvoice,
		RequiresAttendance:  cmd.RequiresAttendance,
		StatutoryDefaults:   cmd.StatutoryDefaults,
		ComplianceNotes:     cmd.ComplianceNotes,
		IsSystemDefault:     isSystemDefault,
		SortOrder:           cmd.SortOrder,
	}
}

func workerClassificationRuleInputFromCommand(cmd ports.WorkerClassificationRuleCommand) domain.WorkerClassificationRuleInput {
	return domain.WorkerClassificationRuleInput{
		TenantID:     cmd.TenantID,
		WorkerTypeID: cmd.WorkerTypeID,
		RuleName:     cmd.RuleName,
		RuleType:     cmd.RuleType,
		Priority:     cmd.Priority,
		Conditions:   cmd.Conditions,
		Outcome:      cmd.Outcome,
		Notes:        cmd.Notes,
	}
}
