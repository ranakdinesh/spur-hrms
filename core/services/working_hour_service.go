package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateWorkingHour(ctx context.Context, cmd ports.WorkingHourCommand) (*domain.WorkingHour, error) {
	item, err := domain.NewWorkingHour(domain.WorkingHourInput{
		TenantID:     cmd.TenantID,
		BranchID:     cmd.BranchID,
		UserID:       cmd.UserID,
		DayOfWeek:    cmd.DayOfWeek,
		IsWorkingDay: cmd.IsWorkingDay,
		StartTime:    cmd.StartTime,
		EndTime:      cmd.EndTime,
		BreakMinutes: cmd.BreakMinutes,
	})
	if err != nil {
		s.logError("validate working hour create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("day_of_week", cmd.DayOfWeek))
		return nil, err
	}
	result, err := s.workingHours.CreateWorkingHour(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create working hour", err, serviceTenantIDField(cmd.TenantID), serviceStringField("day_of_week", item.DayOfWeek))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListWorkingHours(ctx context.Context, tenantID uuid.UUID) ([]*domain.WorkingHour, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate working hour list tenant", err)
		return nil, err
	}
	result, err := s.workingHours.ListWorkingHours(ctx, tenantID)
	if err != nil {
		s.logError("list working hours", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	s.log.Warn().Str("tenant_id", tenantID.String()).Msg("hrms: working hours missing, seeding tenant defaults")
	for _, input := range domain.DefaultWorkingHourInputs(tenantID) {
		item, itemErr := domain.NewWorkingHour(input)
		if itemErr != nil {
			s.logError("validate default working hour", itemErr, serviceTenantIDField(tenantID), serviceStringField("day_of_week", input.DayOfWeek))
			return nil, itemErr
		}
		if _, itemErr = s.workingHours.CreateWorkingHour(ctx, item, nil); itemErr != nil {
			s.logError("seed working hour", itemErr, serviceTenantIDField(tenantID), serviceStringField("day_of_week", item.DayOfWeek))
			return nil, itemErr
		}
	}
	result, err = s.workingHours.ListWorkingHours(ctx, tenantID)
	if err != nil {
		s.logError("list seeded working hours", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetWorkingHour(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkingHour, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate working hour get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkingHourID
		s.logError("validate working hour get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.workingHours.GetWorkingHour(ctx, tenantID, id)
	if err != nil {
		s.logError("get working hour", err, serviceTenantIDField(tenantID), serviceStringField("working_hour_id", id.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ResolveWorkingHour(ctx context.Context, cmd ports.ResolveWorkingHourCommand) (*domain.WorkingHour, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate working hour resolve tenant", err)
		return nil, err
	}
	day, err := domain.NormalizeDayOfWeek(cmd.DayOfWeek)
	if err != nil {
		s.logError("validate working hour resolve day", err, serviceTenantIDField(cmd.TenantID), serviceStringField("day_of_week", cmd.DayOfWeek))
		return nil, err
	}
	result, err := s.workingHours.ResolveWorkingHour(ctx, cmd.TenantID, day, cmd.BranchID, cmd.UserID)
	if err == nil {
		result.Source = domain.WorkingHourScope(result)
		return result, nil
	}
	if errors.Is(err, domain.ErrWorkingHourNotFound) {
		fallback, fallbackErr := domain.DefaultWorkingHour(cmd.TenantID, day)
		if fallbackErr != nil {
			s.logError("build default working hour", fallbackErr, serviceTenantIDField(cmd.TenantID), serviceStringField("day_of_week", day))
			return nil, fallbackErr
		}
		return fallback, nil
	}
	s.logError("resolve working hour", err, serviceTenantIDField(cmd.TenantID), serviceStringField("day_of_week", day))
	return nil, err
}

func (s *TenantService) UpdateWorkingHour(ctx context.Context, cmd ports.WorkingHourCommand) (*domain.WorkingHour, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidWorkingHourID
		s.logError("validate working hour update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	item, err := domain.NewWorkingHour(domain.WorkingHourInput{
		TenantID:     cmd.TenantID,
		BranchID:     cmd.BranchID,
		UserID:       cmd.UserID,
		DayOfWeek:    cmd.DayOfWeek,
		IsWorkingDay: cmd.IsWorkingDay,
		StartTime:    cmd.StartTime,
		EndTime:      cmd.EndTime,
		BreakMinutes: cmd.BreakMinutes,
	})
	if err != nil {
		s.logError("validate working hour update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("working_hour_id", cmd.ID.String()), serviceStringField("day_of_week", cmd.DayOfWeek))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.workingHours.UpdateWorkingHour(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update working hour", err, serviceTenantIDField(cmd.TenantID), serviceStringField("working_hour_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteWorkingHour(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate working hour delete tenant", err)
		return err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidWorkingHourID
		s.logError("validate working hour delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.workingHours.DeleteWorkingHour(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete working hour", err, serviceTenantIDField(tenantID), serviceStringField("working_hour_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) CopyTenantWorkingHoursToBranch(ctx context.Context, cmd ports.CopyWorkingHoursCommand) ([]*domain.WorkingHour, error) {
	if cmd.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate working hour copy tenant", err)
		return nil, err
	}
	if cmd.BranchID == uuid.Nil {
		err := domain.ErrInvalidWorkingHourBranch
		s.logError("validate working hour copy branch", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.branches.GetBranch(ctx, cmd.TenantID, cmd.BranchID); err != nil {
		s.logError("validate working hour copy branch exists", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", cmd.BranchID.String()))
		return nil, err
	}
	var result []*domain.WorkingHour
	if err := s.system.RunAsSystem(ctx, func(txCtx context.Context) error {
		items, err := s.workingHours.ListWorkingHours(txCtx, cmd.TenantID)
		if err != nil {
			return err
		}
		if !hasTenantDefaultWorkingHours(items) {
			for _, input := range domain.DefaultWorkingHourInputs(cmd.TenantID) {
				item, itemErr := domain.NewWorkingHour(input)
				if itemErr != nil {
					return itemErr
				}
				if _, itemErr = s.workingHours.CreateWorkingHour(txCtx, item, cmd.ActorID); itemErr != nil {
					return itemErr
				}
			}
		}
		result, err = s.workingHours.CopyTenantWorkingHoursToBranch(txCtx, cmd.TenantID, cmd.BranchID, cmd.ActorID)
		return err
	}); err != nil {
		s.logError("copy tenant working hours to branch", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", cmd.BranchID.String()))
		return nil, err
	}
	if len(result) == 0 {
		err := domain.ErrNoTenantWorkingHoursToCopy
		s.logError("copy tenant working hours empty", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", cmd.BranchID.String()))
		return nil, err
	}
	return result, nil
}

func hasTenantDefaultWorkingHours(items []*domain.WorkingHour) bool {
	count := 0
	for _, item := range items {
		if item != nil && item.BranchID == nil && item.UserID == nil {
			count++
		}
	}
	return count >= 7
}
