package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateDesignation(ctx context.Context, cmd ports.DesignationCommand) (*domain.Designation, error) {
	designation, err := domain.NewDesignation(domain.DesignationInput{
		TenantID:           cmd.TenantID,
		Name:               cmd.Name,
		LevelCode:          cmd.LevelCode,
		SeniorityRank:      cmd.SeniorityRank,
		Description:        cmd.Description,
		AttendanceRequired: cmd.AttendanceRequired,
	})
	if err != nil {
		s.logError("validate designation create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("designation_name", cmd.Name), serviceStringField("designation_level_code", cmd.LevelCode))
		return nil, err
	}
	result, err := s.designations.CreateDesignation(ctx, designation, cmd.ActorID)
	if err != nil {
		s.logError("create designation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("designation_name", designation.Name))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("designation_id", result.ID.String()).Msg("hrms: designation created")
	return result, nil
}

func (s *TenantService) ListDesignations(ctx context.Context, tenantID uuid.UUID) ([]*domain.Designation, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate designation list tenant", err)
		return nil, err
	}
	result, err := s.designations.ListDesignations(ctx, tenantID)
	if err != nil {
		s.logError("list designations", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetDesignation(ctx context.Context, tenantID uuid.UUID, designationID uuid.UUID) (*domain.Designation, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate designation get tenant", err)
		return nil, err
	}
	if designationID == uuid.Nil {
		err := domain.ErrInvalidDesignationID
		s.logError("validate designation get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.designations.GetDesignation(ctx, tenantID, designationID)
	if err != nil {
		s.logError("get designation", err, serviceTenantIDField(tenantID), serviceStringField("designation_id", designationID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateDesignation(ctx context.Context, cmd ports.DesignationCommand) (*domain.Designation, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidDesignationID
		s.logError("validate designation update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	designation, err := domain.NewDesignation(domain.DesignationInput{
		TenantID:           cmd.TenantID,
		Name:               cmd.Name,
		LevelCode:          cmd.LevelCode,
		SeniorityRank:      cmd.SeniorityRank,
		Description:        cmd.Description,
		AttendanceRequired: cmd.AttendanceRequired,
	})
	if err != nil {
		s.logError("validate designation update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("designation_id", cmd.ID.String()), serviceStringField("designation_name", cmd.Name), serviceStringField("designation_level_code", cmd.LevelCode))
		return nil, err
	}
	designation.ID = cmd.ID
	if cmd.AttendanceRequired == nil {
		existing, existingErr := s.designations.GetDesignation(ctx, cmd.TenantID, cmd.ID)
		if existingErr != nil {
			s.logError("get designation before update", existingErr, serviceTenantIDField(cmd.TenantID), serviceStringField("designation_id", cmd.ID.String()))
			return nil, existingErr
		}
		designation.AttendanceRequired = existing.AttendanceRequired
	}
	result, err := s.designations.UpdateDesignation(ctx, designation, cmd.ActorID)
	if err != nil {
		s.logError("update designation", err, serviceTenantIDField(cmd.TenantID), serviceStringField("designation_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("designation_id", result.ID.String()).Msg("hrms: designation updated")
	return result, nil
}

func (s *TenantService) DeleteDesignation(ctx context.Context, tenantID uuid.UUID, designationID uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate designation delete tenant", err)
		return err
	}
	if designationID == uuid.Nil {
		err := domain.ErrInvalidDesignationID
		s.logError("validate designation delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.designations.DeleteDesignation(ctx, tenantID, designationID, actorID); err != nil {
		s.logError("delete designation", err, serviceTenantIDField(tenantID), serviceStringField("designation_id", designationID.String()))
		return err
	}
	s.log.Info().Str("tenant_id", tenantID.String()).Str("designation_id", designationID.String()).Msg("hrms: designation deactivated")
	return nil
}

func (s *TenantService) CreateDesignationLevelCode(ctx context.Context, cmd ports.DesignationLevelCodeCommand) (*domain.DesignationLevelCode, error) {
	item, err := domain.NewDesignationLevelCode(domain.DesignationLevelCodeInput{
		TenantID:    cmd.TenantID,
		Code:        cmd.Code,
		Label:       cmd.Label,
		Description: cmd.Description,
		SortOrder:   cmd.SortOrder,
	})
	if err != nil {
		s.logError("validate designation level code create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("level_code", cmd.Code))
		return nil, err
	}
	result, err := s.designationMasters.CreateDesignationLevelCode(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create designation level code", err, serviceTenantIDField(cmd.TenantID), serviceStringField("level_code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListDesignationLevelCodes(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.DesignationLevelCode, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate designation level code list tenant", err)
		return nil, err
	}
	result, err := s.designationMasters.ListDesignationLevelCodes(ctx, tenantID)
	if err != nil {
		s.logError("list designation level codes", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	s.log.Warn().Str("tenant_id", tenantID.String()).Msg("hrms: designation level code masters missing, seeding defaults")
	for _, input := range domain.DefaultDesignationLevelCodeInputs(tenantID) {
		item, itemErr := domain.NewDesignationLevelCode(input)
		if itemErr != nil {
			s.logError("validate default designation level code", itemErr, serviceTenantIDField(tenantID), serviceStringField("level_code", input.Code))
			return nil, itemErr
		}
		if _, itemErr = s.designationMasters.CreateDesignationLevelCode(ctx, item, actorID); itemErr != nil {
			s.logError("seed designation level code", itemErr, serviceTenantIDField(tenantID), serviceStringField("level_code", item.Code))
			return nil, itemErr
		}
	}
	result, err = s.designationMasters.ListDesignationLevelCodes(ctx, tenantID)
	if err != nil {
		s.logError("list seeded designation level codes", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateDesignationLevelCode(ctx context.Context, cmd ports.DesignationLevelCodeCommand) (*domain.DesignationLevelCode, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidDesignationLevelCodeID
		s.logError("validate designation level code update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewDesignationLevelCode(domain.DesignationLevelCodeInput{
		TenantID:    cmd.TenantID,
		Code:        cmd.Code,
		Label:       cmd.Label,
		Description: cmd.Description,
		SortOrder:   cmd.SortOrder,
	})
	if err != nil {
		s.logError("validate designation level code update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("level_code_id", cmd.ID.String()), serviceStringField("level_code", cmd.Code))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.designationMasters.UpdateDesignationLevelCode(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update designation level code", err, serviceTenantIDField(cmd.TenantID), serviceStringField("level_code_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteDesignationLevelCode(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate designation level code delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidDesignationLevelCodeID
		s.logError("validate designation level code delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.designationMasters.DeleteDesignationLevelCode(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete designation level code", err, serviceTenantIDField(tenantID), serviceStringField("level_code_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CreateDesignationSeniorityRank(ctx context.Context, cmd ports.DesignationSeniorityRankCommand) (*domain.DesignationSeniorityRank, error) {
	item, err := domain.NewDesignationSeniorityRank(domain.DesignationSeniorityRankInput{
		TenantID:    cmd.TenantID,
		RankValue:   cmd.RankValue,
		Label:       cmd.Label,
		Description: cmd.Description,
		SortOrder:   cmd.SortOrder,
	})
	if err != nil {
		s.logError("validate designation seniority rank create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.designationMasters.CreateDesignationSeniorityRank(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create designation seniority rank", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListDesignationSeniorityRanks(ctx context.Context, tenantID uuid.UUID, actorID *uuid.UUID) ([]*domain.DesignationSeniorityRank, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate designation seniority rank list tenant", err)
		return nil, err
	}
	result, err := s.designationMasters.ListDesignationSeniorityRanks(ctx, tenantID)
	if err != nil {
		s.logError("list designation seniority ranks", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	s.log.Warn().Str("tenant_id", tenantID.String()).Msg("hrms: designation seniority rank masters missing, seeding defaults")
	for _, input := range domain.DefaultDesignationSeniorityRankInputs(tenantID) {
		item, itemErr := domain.NewDesignationSeniorityRank(input)
		if itemErr != nil {
			s.logError("validate default designation seniority rank", itemErr, serviceTenantIDField(tenantID))
			return nil, itemErr
		}
		if _, itemErr = s.designationMasters.CreateDesignationSeniorityRank(ctx, item, actorID); itemErr != nil {
			s.logError("seed designation seniority rank", itemErr, serviceTenantIDField(tenantID))
			return nil, itemErr
		}
	}
	result, err = s.designationMasters.ListDesignationSeniorityRanks(ctx, tenantID)
	if err != nil {
		s.logError("list seeded designation seniority ranks", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateDesignationSeniorityRank(ctx context.Context, cmd ports.DesignationSeniorityRankCommand) (*domain.DesignationSeniorityRank, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidDesignationSeniorityRankID
		s.logError("validate designation seniority rank update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewDesignationSeniorityRank(domain.DesignationSeniorityRankInput{
		TenantID:    cmd.TenantID,
		RankValue:   cmd.RankValue,
		Label:       cmd.Label,
		Description: cmd.Description,
		SortOrder:   cmd.SortOrder,
	})
	if err != nil {
		s.logError("validate designation seniority rank update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("rank_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.designationMasters.UpdateDesignationSeniorityRank(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update designation seniority rank", err, serviceTenantIDField(cmd.TenantID), serviceStringField("rank_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteDesignationSeniorityRank(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate designation seniority rank delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidDesignationSeniorityRankID
		s.logError("validate designation seniority rank delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.designationMasters.DeleteDesignationSeniorityRank(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete designation seniority rank", err, serviceTenantIDField(tenantID), serviceStringField("rank_id", id.String()))
		return err
	}
	return nil
}
