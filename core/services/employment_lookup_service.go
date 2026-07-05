package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateEmploymentType(ctx context.Context, cmd ports.EmploymentTypeCommand) (*domain.EmploymentType, error) {
	item, err := domain.NewEmploymentType(domain.EmploymentTypeInput{TenantID: cmd.TenantID, Name: cmd.Name})
	if err != nil {
		s.logError("validate employment type create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employment_type_name", cmd.Name))
		return nil, err
	}
	result, err := s.lookups.CreateEmploymentType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create employment type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employment_type_name", item.Name))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListEmploymentTypes(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.EmploymentType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employment type list tenant", err)
		return nil, err
	}
	result, err := s.lookups.ListEmploymentTypes(ctx, tenantID)
	if err != nil {
		s.logError("list employment types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	s.log.Warn().Str("tenant_id", tenantID.String()).Msg("hrms: employment type lookups missing, seeding defaults")
	for _, input := range domain.DefaultEmploymentTypeInputs(tenantID) {
		item, itemErr := domain.NewEmploymentType(input)
		if itemErr != nil {
			s.logError("validate default employment type", itemErr, serviceTenantIDField(tenantID), serviceStringField("employment_type_name", input.Name))
			return nil, itemErr
		}
		if _, itemErr = s.lookups.CreateEmploymentType(ctx, item, actorID); itemErr != nil {
			s.logError("seed employment type", itemErr, serviceTenantIDField(tenantID), serviceStringField("employment_type_name", item.Name))
			return nil, itemErr
		}
	}
	result, err = s.lookups.ListEmploymentTypes(ctx, tenantID)
	if err != nil {
		s.logError("list seeded employment types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetEmploymentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmploymentType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employment type get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidEmploymentTypeID
		s.logError("validate employment type get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.lookups.GetEmploymentType(ctx, tenantID, id)
	if err != nil {
		s.logError("get employment type", err, serviceTenantIDField(tenantID), serviceStringField("employment_type_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateEmploymentType(ctx context.Context, cmd ports.EmploymentTypeCommand) (*domain.EmploymentType, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidEmploymentTypeID
		s.logError("validate employment type update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewEmploymentType(domain.EmploymentTypeInput{TenantID: cmd.TenantID, Name: cmd.Name})
	if err != nil {
		s.logError("validate employment type update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employment_type_id", cmd.ID.String()), serviceStringField("employment_type_name", cmd.Name))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.lookups.UpdateEmploymentType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update employment type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("employment_type_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteEmploymentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate employment type delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidEmploymentTypeID
		s.logError("validate employment type delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.lookups.DeleteEmploymentType(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete employment type", err, serviceTenantIDField(tenantID), serviceStringField("employment_type_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateMaritalStatus(ctx context.Context, cmd ports.MaritalStatusCommand) (*domain.MaritalStatus, error) {
	item, err := domain.NewMaritalStatus(domain.MaritalStatusInput{TenantID: cmd.TenantID, Name: cmd.Name})
	if err != nil {
		s.logError("validate marital status create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("marital_status_name", cmd.Name))
		return nil, err
	}
	result, err := s.lookups.CreateMaritalStatus(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create marital status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("marital_status_name", item.Name))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListMaritalStatuses(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.MaritalStatus, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate marital status list tenant", err)
		return nil, err
	}
	result, err := s.lookups.ListMaritalStatuses(ctx, tenantID)
	if err != nil {
		s.logError("list marital statuses", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	s.log.Warn().Str("tenant_id", tenantID.String()).Msg("hrms: marital status lookups missing, seeding defaults")
	for _, input := range domain.DefaultMaritalStatusInputs(tenantID) {
		item, itemErr := domain.NewMaritalStatus(input)
		if itemErr != nil {
			s.logError("validate default marital status", itemErr, serviceTenantIDField(tenantID), serviceStringField("marital_status_name", input.Name))
			return nil, itemErr
		}
		if _, itemErr = s.lookups.CreateMaritalStatus(ctx, item, actorID); itemErr != nil {
			s.logError("seed marital status", itemErr, serviceTenantIDField(tenantID), serviceStringField("marital_status_name", item.Name))
			return nil, itemErr
		}
	}
	result, err = s.lookups.ListMaritalStatuses(ctx, tenantID)
	if err != nil {
		s.logError("list seeded marital statuses", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetMaritalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.MaritalStatus, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate marital status get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidMaritalStatusID
		s.logError("validate marital status get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.lookups.GetMaritalStatus(ctx, tenantID, id)
	if err != nil {
		s.logError("get marital status", err, serviceTenantIDField(tenantID), serviceStringField("marital_status_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateMaritalStatus(ctx context.Context, cmd ports.MaritalStatusCommand) (*domain.MaritalStatus, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidMaritalStatusID
		s.logError("validate marital status update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewMaritalStatus(domain.MaritalStatusInput{TenantID: cmd.TenantID, Name: cmd.Name})
	if err != nil {
		s.logError("validate marital status update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("marital_status_id", cmd.ID.String()), serviceStringField("marital_status_name", cmd.Name))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.lookups.UpdateMaritalStatus(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update marital status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("marital_status_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteMaritalStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate marital status delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidMaritalStatusID
		s.logError("validate marital status delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.lookups.DeleteMaritalStatus(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete marital status", err, serviceTenantIDField(tenantID), serviceStringField("marital_status_id", id.String()))
		return err
	}
	return nil
}
