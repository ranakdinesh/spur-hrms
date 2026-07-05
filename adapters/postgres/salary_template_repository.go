package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateSalaryTemplate(ctx context.Context, item *domain.SalaryTemplate, actorID *uuid.UUID) (*domain.SalaryTemplate, error) {
	row, err := s.getQueries(ctx).CreateSalaryTemplate(ctx, sqlc.CreateSalaryTemplateParams{TenantID: item.TenantID, FyID: item.FYID, Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), TemplateType: item.TemplateType, AppliesTo: item.AppliesTo, CurrencyCode: item.CurrencyCode, EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), Notes: textFromPtr(item.Notes), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create salary template", err, tenantIDField(item.TenantID), stringField("salary_template_code", item.Code))
	}
	return mapSalaryTemplate(row), nil
}

func (s *Store) ListSalaryTemplates(ctx context.Context, tenantID uuid.UUID, fyID *uuid.UUID) ([]*domain.SalaryTemplate, error) {
	queries := s.getQueries(ctx)
	if fyID != nil && *fyID != uuid.Nil {
		rows, err := queries.ListSalaryTemplatesByFY(ctx, sqlc.ListSalaryTemplatesByFYParams{TenantID: tenantID, FyID: *fyID})
		if err != nil {
			return nil, s.logDBError(ctx, "list salary templates by fy", err, tenantIDField(tenantID), stringField("financial_year_id", fyID.String()))
		}
		return mapSalaryTemplates(rows), nil
	}
	rows, err := queries.ListSalaryTemplates(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list salary templates", err, tenantIDField(tenantID))
	}
	return mapSalaryTemplates(rows), nil
}

func (s *Store) GetSalaryTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SalaryTemplate, error) {
	row, err := s.getQueries(ctx).GetSalaryTemplate(ctx, sqlc.GetSalaryTemplateParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get salary template", err, tenantIDField(tenantID), stringField("salary_template_id", id.String()))
	}
	item := mapSalaryTemplate(row)
	children, err := s.ListSalaryTemplateItems(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	item.Items = children
	return item, nil
}

func (s *Store) GetActiveSalaryTemplate(ctx context.Context, tenantID uuid.UUID, fyID uuid.UUID) (*domain.SalaryTemplate, error) {
	row, err := s.getQueries(ctx).GetActiveSalaryTemplate(ctx, sqlc.GetActiveSalaryTemplateParams{TenantID: tenantID, FyID: fyID})
	if err != nil {
		return nil, s.logDBError(ctx, "get active salary template", err, tenantIDField(tenantID), stringField("financial_year_id", fyID.String()))
	}
	return mapSalaryTemplate(row), nil
}

func (s *Store) UpdateSalaryTemplate(ctx context.Context, item *domain.SalaryTemplate, actorID *uuid.UUID) (*domain.SalaryTemplate, error) {
	row, err := s.getQueries(ctx).UpdateSalaryTemplate(ctx, sqlc.UpdateSalaryTemplateParams{TenantID: item.TenantID, ID: item.ID, FyID: item.FYID, Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), TemplateType: item.TemplateType, AppliesTo: item.AppliesTo, CurrencyCode: item.CurrencyCode, EffectiveFrom: dateFromPtr(item.EffectiveFrom), EffectiveTo: dateFromPtr(item.EffectiveTo), Notes: textFromPtr(item.Notes), UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update salary template", err, tenantIDField(item.TenantID), stringField("salary_template_id", item.ID.String()))
	}
	return mapSalaryTemplate(row), nil
}

func (s *Store) DeleteSalaryTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteSalaryTemplate(ctx, sqlc.SoftDeleteSalaryTemplateParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete salary template", err, tenantIDField(tenantID), stringField("salary_template_id", id.String()))
	}
	return nil
}

func (s *Store) ActivateSalaryTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.SalaryTemplate, error) {
	existing, err := s.GetSalaryTemplate(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	queries := s.getQueries(ctx)
	if err := queries.DeactivateSalaryTemplatesForFY(ctx, sqlc.DeactivateSalaryTemplatesForFYParams{TenantID: tenantID, FyID: existing.FYID, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return nil, s.logDBError(ctx, "deactivate salary templates for fy", err, tenantIDField(tenantID), stringField("financial_year_id", existing.FYID.String()))
	}
	row, err := queries.ActivateSalaryTemplate(ctx, sqlc.ActivateSalaryTemplateParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "activate salary template", err, tenantIDField(tenantID), stringField("salary_template_id", id.String()))
	}
	return mapSalaryTemplate(row), nil
}

func (s *Store) CreateSalaryTemplateItem(ctx context.Context, item *domain.SalaryTemplateItem, actorID *uuid.UUID) (*domain.SalaryTemplateItem, error) {
	row, err := s.getQueries(ctx).CreateSalaryTemplateItem(ctx, sqlc.CreateSalaryTemplateItemParams{TenantID: item.TenantID, TemplateID: item.TemplateID, ItemType: item.ItemType, Code: item.Code, Name: item.Name, Percentage: numericFromFloatPtr(item.Percentage), Amount: numericFromFloatPtr(item.Amount), CalculationMode: item.CalculationMode, CalculationBase: item.CalculationBase, Formula: textFromPtr(item.Formula), ContributionSide: item.ContributionSide, IsTaxExempt: item.IsTaxExempt, IsStatutory: item.IsStatutory, IsVariable: item.IsVariable, AffectsGross: item.AffectsGross, AffectsNet: item.AffectsNet, CapAmount: numericFromFloatPtr(item.CapAmount), MinAmount: numericFromFloatPtr(item.MinAmount), MaxAmount: numericFromFloatPtr(item.MaxAmount), SortOrder: item.SortOrder, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create salary template item", err, tenantIDField(item.TenantID), stringField("salary_template_id", item.TemplateID.String()))
	}
	return mapSalaryTemplateItem(row), nil
}

func (s *Store) ListSalaryTemplateItems(ctx context.Context, tenantID uuid.UUID, templateID uuid.UUID) ([]*domain.SalaryTemplateItem, error) {
	rows, err := s.getQueries(ctx).ListSalaryTemplateItems(ctx, sqlc.ListSalaryTemplateItemsParams{TenantID: tenantID, TemplateID: templateID})
	if err != nil {
		return nil, s.logDBError(ctx, "list salary template items", err, tenantIDField(tenantID), stringField("salary_template_id", templateID.String()))
	}
	return mapSalaryTemplateItems(rows), nil
}

func (s *Store) GetSalaryTemplateItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SalaryTemplateItem, error) {
	row, err := s.getQueries(ctx).GetSalaryTemplateItem(ctx, sqlc.GetSalaryTemplateItemParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get salary template item", err, tenantIDField(tenantID), stringField("salary_template_item_id", id.String()))
	}
	return mapSalaryTemplateItem(row), nil
}

func (s *Store) UpdateSalaryTemplateItem(ctx context.Context, item *domain.SalaryTemplateItem, actorID *uuid.UUID) (*domain.SalaryTemplateItem, error) {
	row, err := s.getQueries(ctx).UpdateSalaryTemplateItem(ctx, sqlc.UpdateSalaryTemplateItemParams{TenantID: item.TenantID, TemplateID: item.TemplateID, ID: item.ID, ItemType: item.ItemType, Code: item.Code, Name: item.Name, Percentage: numericFromFloatPtr(item.Percentage), Amount: numericFromFloatPtr(item.Amount), CalculationMode: item.CalculationMode, CalculationBase: item.CalculationBase, Formula: textFromPtr(item.Formula), ContributionSide: item.ContributionSide, IsTaxExempt: item.IsTaxExempt, IsStatutory: item.IsStatutory, IsVariable: item.IsVariable, AffectsGross: item.AffectsGross, AffectsNet: item.AffectsNet, CapAmount: numericFromFloatPtr(item.CapAmount), MinAmount: numericFromFloatPtr(item.MinAmount), MaxAmount: numericFromFloatPtr(item.MaxAmount), SortOrder: item.SortOrder, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update salary template item", err, tenantIDField(item.TenantID), stringField("salary_template_item_id", item.ID.String()))
	}
	return mapSalaryTemplateItem(row), nil
}

func (s *Store) DeleteSalaryTemplateItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteSalaryTemplateItem(ctx, sqlc.SoftDeleteSalaryTemplateItemParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete salary template item", err, tenantIDField(tenantID), stringField("salary_template_item_id", id.String()))
	}
	return nil
}
