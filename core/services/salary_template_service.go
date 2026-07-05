package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateSalaryTemplate(ctx context.Context, cmd ports.SalaryTemplateCommand) (*domain.SalaryTemplate, error) {
	item, err := salaryTemplateFromCommand(cmd)
	if err != nil {
		s.logError("validate salary template create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.salaryTemplates.CreateSalaryTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create salary template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_template_code", item.Code))
		return nil, err
	}
	if cmd.IsActive {
		result, err = s.ActivateSalaryTemplate(ctx, cmd.TenantID, result.ID, cmd.ActorID)
		if err != nil {
			return nil, err
		}
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("salary_template_id", result.ID.String()).Msg("hrms: salary template created")
	return result, nil
}

func (s *TenantService) ListSalaryTemplates(ctx context.Context, tenantID uuid.UUID, fyID *uuid.UUID) ([]*domain.SalaryTemplate, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate salary template list tenant", err)
		return nil, err
	}
	items, err := s.salaryTemplates.ListSalaryTemplates(ctx, tenantID, fyID)
	if err != nil {
		s.logError("list salary templates", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) GetSalaryTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SalaryTemplate, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidSalaryTemplateID
		s.logError("validate salary template get", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.salaryTemplates.GetSalaryTemplate(ctx, tenantID, id)
	if err != nil {
		s.logError("get salary template", err, serviceTenantIDField(tenantID), serviceStringField("salary_template_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) GetActiveSalaryTemplate(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) (*domain.SalaryTemplate, error) {
	if tenantID == uuid.Nil || fyID == uuid.Nil {
		err := domain.ErrInvalidSalaryTemplate
		s.logError("validate active salary template", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.salaryTemplates.GetActiveSalaryTemplate(ctx, tenantID, fyID)
	if err != nil {
		s.logError("get active salary template", err, serviceTenantIDField(tenantID), serviceStringField("financial_year_id", fyID.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateSalaryTemplate(ctx context.Context, cmd ports.SalaryTemplateCommand) (*domain.SalaryTemplate, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidSalaryTemplateID
		s.logError("validate salary template update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := salaryTemplateFromCommand(cmd)
	if err != nil {
		s.logError("validate salary template update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_template_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.salaryTemplates.UpdateSalaryTemplate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update salary template", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_template_id", cmd.ID.String()))
		return nil, err
	}
	if cmd.IsActive && !result.IsActive {
		result, err = s.ActivateSalaryTemplate(ctx, cmd.TenantID, result.ID, cmd.ActorID)
		if err != nil {
			return nil, err
		}
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("salary_template_id", result.ID.String()).Msg("hrms: salary template updated")
	return result, nil
}

func (s *TenantService) DeleteSalaryTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidSalaryTemplateID
		s.logError("validate salary template delete", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.salaryTemplates.DeleteSalaryTemplate(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete salary template", err, serviceTenantIDField(tenantID), serviceStringField("salary_template_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) ActivateSalaryTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.SalaryTemplate, error) {
	existing, err := s.GetSalaryTemplate(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	result, err := s.salaryTemplates.ActivateSalaryTemplate(ctx, tenantID, existing.ID, actorID)
	if err != nil {
		s.logError("activate salary template", err, serviceTenantIDField(tenantID), serviceStringField("salary_template_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) CreateSalaryTemplateItem(ctx context.Context, cmd ports.SalaryTemplateItemCommand) (*domain.SalaryTemplateItem, error) {
	item, err := salaryTemplateItemFromCommand(cmd)
	if err != nil {
		s.logError("validate salary template item create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.salaryTemplates.CreateSalaryTemplateItem(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create salary template item", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_template_id", cmd.TemplateID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListSalaryTemplateItems(ctx context.Context, tenantID uuid.UUID, templateID uuid.UUID) ([]*domain.SalaryTemplateItem, error) {
	if tenantID == uuid.Nil || templateID == uuid.Nil {
		err := domain.ErrInvalidSalaryTemplateID
		s.logError("validate salary template item list", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	items, err := s.salaryTemplates.ListSalaryTemplateItems(ctx, tenantID, templateID)
	if err != nil {
		s.logError("list salary template items", err, serviceTenantIDField(tenantID), serviceStringField("salary_template_id", templateID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) UpdateSalaryTemplateItem(ctx context.Context, cmd ports.SalaryTemplateItemCommand) (*domain.SalaryTemplateItem, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidSalaryItemID
		s.logError("validate salary template item update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := salaryTemplateItemFromCommand(cmd)
	if err != nil {
		s.logError("validate salary template item update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_template_item_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.salaryTemplates.UpdateSalaryTemplateItem(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update salary template item", err, serviceTenantIDField(cmd.TenantID), serviceStringField("salary_template_item_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteSalaryTemplateItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		err := domain.ErrInvalidSalaryItemID
		s.logError("validate salary template item delete", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.salaryTemplates.DeleteSalaryTemplateItem(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete salary template item", err, serviceTenantIDField(tenantID), serviceStringField("salary_template_item_id", id.String()))
		return err
	}
	return nil
}

func salaryTemplateFromCommand(cmd ports.SalaryTemplateCommand) (*domain.SalaryTemplate, error) {
	effectiveFrom, err := parseOptionalSalaryTemplateDate(cmd.EffectiveFrom)
	if err != nil {
		return nil, err
	}
	effectiveTo, err := parseOptionalSalaryTemplateDate(cmd.EffectiveTo)
	if err != nil {
		return nil, err
	}
	return domain.NewSalaryTemplate(domain.SalaryTemplateInput{TenantID: cmd.TenantID, FYID: cmd.FYID, Code: cmd.Code, Name: cmd.Name, Description: cmd.Description, TemplateType: cmd.TemplateType, AppliesTo: cmd.AppliesTo, CurrencyCode: cmd.CurrencyCode, EffectiveFrom: effectiveFrom, EffectiveTo: effectiveTo, Notes: cmd.Notes})
}

func salaryTemplateItemFromCommand(cmd ports.SalaryTemplateItemCommand) (*domain.SalaryTemplateItem, error) {
	return domain.NewSalaryTemplateItem(domain.SalaryTemplateItemInput{TenantID: cmd.TenantID, TemplateID: cmd.TemplateID, ItemType: cmd.ItemType, Code: cmd.Code, Name: cmd.Name, Percentage: cmd.Percentage, Amount: cmd.Amount, CalculationMode: cmd.CalculationMode, CalculationBase: cmd.CalculationBase, Formula: cmd.Formula, ContributionSide: cmd.ContributionSide, IsTaxExempt: cmd.IsTaxExempt, IsStatutory: cmd.IsStatutory, IsVariable: cmd.IsVariable, AffectsGross: cmd.AffectsGross, AffectsNet: cmd.AffectsNet, CapAmount: cmd.CapAmount, MinAmount: cmd.MinAmount, MaxAmount: cmd.MaxAmount, SortOrder: cmd.SortOrder})
}

func parseOptionalSalaryTemplateDate(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil, domain.ErrInvalidSalaryTemplate
	}
	return &parsed, nil
}
