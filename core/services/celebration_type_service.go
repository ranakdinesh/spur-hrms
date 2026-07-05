package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateCelebrationType(ctx context.Context, cmd ports.CelebrationTypeCommand) (*domain.CelebrationType, error) {
	item, err := domain.NewCelebrationType(domain.CelebrationTypeInput{TenantID: cmd.TenantID, Name: cmd.Name, IsYearly: cmd.IsYearly, IsUserCelebration: cmd.IsUserCelebration})
	if err != nil {
		s.logError("validate celebration type create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("celebration_type_name", cmd.Name))
		return nil, err
	}
	result, err := s.celebrations.CreateCelebrationType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create celebration type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("celebration_type_name", item.Name))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListCelebrationTypes(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.CelebrationType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate celebration type list tenant", err)
		return nil, err
	}
	result, err := s.celebrations.ListCelebrationTypes(ctx, tenantID)
	if err != nil {
		s.logError("list celebration types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	s.log.Warn().Str("tenant_id", tenantID.String()).Msg("hrms: celebration types missing, seeding defaults")
	if _, err := s.seedCelebrationTypes(ctx, tenantID, actorID); err != nil {
		return nil, err
	}
	result, err = s.celebrations.ListCelebrationTypes(ctx, tenantID)
	if err != nil {
		s.logError("list seeded celebration types", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetCelebrationType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CelebrationType, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate celebration type get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidCelebrationTypeID
		s.logError("validate celebration type get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.celebrations.GetCelebrationType(ctx, tenantID, id)
	if err != nil {
		s.logError("get celebration type", err, serviceTenantIDField(tenantID), serviceStringField("celebration_type_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateCelebrationType(ctx context.Context, cmd ports.CelebrationTypeCommand) (*domain.CelebrationType, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidCelebrationTypeID
		s.logError("validate celebration type update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewCelebrationType(domain.CelebrationTypeInput{TenantID: cmd.TenantID, Name: cmd.Name, IsYearly: cmd.IsYearly, IsUserCelebration: cmd.IsUserCelebration})
	if err != nil {
		s.logError("validate celebration type update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("celebration_type_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.celebrations.UpdateCelebrationType(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update celebration type", err, serviceTenantIDField(cmd.TenantID), serviceStringField("celebration_type_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteCelebrationType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate celebration type delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidCelebrationTypeID
		s.logError("validate celebration type delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.celebrations.DeleteCelebrationType(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete celebration type", err, serviceTenantIDField(tenantID), serviceStringField("celebration_type_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) seedCelebrationTypes(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) (int, error) {
	items, err := s.celebrations.ListCelebrationTypes(ctx, tenantID)
	if err != nil || len(items) > 0 {
		return 0, err
	}
	count := 0
	for _, input := range domain.DefaultCelebrationTypeInputs(tenantID) {
		item, err := domain.NewCelebrationType(input)
		if err != nil {
			s.logError("validate default celebration type", err, serviceTenantIDField(tenantID), serviceStringField("celebration_type_name", input.Name))
			return count, err
		}
		if _, err = s.celebrations.CreateCelebrationType(ctx, item, actorID); err != nil {
			s.logError("seed celebration type", err, serviceTenantIDField(tenantID), serviceStringField("celebration_type_name", item.Name))
			return count, err
		}
		count++
	}
	return count, nil
}
