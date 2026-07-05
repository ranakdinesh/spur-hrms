package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateSuccessionReviewCycle(ctx context.Context, cmd ports.SuccessionReviewCycleCommand) (*domain.SuccessionReviewCycle, error) {
	item := &domain.SuccessionReviewCycle{TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, Status: cmd.Status, StartsOn: cmd.StartsOn, EndsOn: cmd.EndsOn, ConfidentialityLevel: cmd.ConfidentialityLevel, Notes: cmd.Notes, Metadata: cmd.Metadata}
	if err := domain.ValidateSuccessionReviewCycle(item); err != nil {
		s.logError("validate succession review cycle", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.successionPlanning.CreateSuccessionReviewCycle(ctx, item, cmd.ActorID)
	if err == nil {
		_, _ = s.successionEvent(ctx, result.TenantID, "cycle", &result.ID, "created", nil, &result.Status, nil, cmd.ActorID)
	}
	return result, err
}

func (s *TenantService) UpdateSuccessionReviewCycle(ctx context.Context, cmd ports.SuccessionReviewCycleCommand) (*domain.SuccessionReviewCycle, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSuccessionReviewCycle
	}
	existing, err := s.successionPlanning.GetSuccessionReviewCycle(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	item := &domain.SuccessionReviewCycle{ID: cmd.ID, TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, Status: existing.Status, StartsOn: cmd.StartsOn, EndsOn: cmd.EndsOn, ConfidentialityLevel: cmd.ConfidentialityLevel, Notes: cmd.Notes, Metadata: cmd.Metadata}
	if err := domain.ValidateSuccessionReviewCycle(item); err != nil {
		return nil, err
	}
	return s.successionPlanning.UpdateSuccessionReviewCycle(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateSuccessionReviewCycleStatus(ctx context.Context, cmd ports.SuccessionStatusCommand) (*domain.SuccessionReviewCycle, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSuccessionReviewCycle
	}
	status := domain.NormalizeSuccessionCycleStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidSuccessionReviewCycle
	}
	before, _ := s.successionPlanning.GetSuccessionReviewCycle(ctx, cmd.TenantID, cmd.ID)
	result, err := s.successionPlanning.UpdateSuccessionReviewCycleStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	_, _ = s.successionEvent(ctx, cmd.TenantID, "cycle", &cmd.ID, "status_changed", from, &result.Status, cmd.Remarks, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListSuccessionReviewCycles(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionReviewCycle, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeSuccessionPage(&filter.Limit, &filter.Offset)
	return s.successionPlanning.ListSuccessionReviewCycles(ctx, filter)
}

func (s *TenantService) CreateSuccessionCriticalRole(ctx context.Context, cmd ports.SuccessionCriticalRoleCommand) (*domain.SuccessionCriticalRole, error) {
	if cmd.CycleID != nil {
		if _, err := s.successionPlanning.GetSuccessionReviewCycle(ctx, cmd.TenantID, *cmd.CycleID); err != nil {
			return nil, err
		}
	}
	if err := s.validateSuccessionWorkerRefs(ctx, cmd.TenantID, cmd.IncumbentWorkerProfileID, cmd.EmergencyCoverWorkerProfileID); err != nil {
		return nil, err
	}
	item := successionCriticalRoleFromCommand(cmd)
	if err := domain.ValidateSuccessionCriticalRole(item); err != nil {
		s.logError("validate succession critical role", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.successionPlanning.CreateSuccessionCriticalRole(ctx, item, cmd.ActorID)
	if err == nil {
		_, _ = s.successionEvent(ctx, result.TenantID, "critical_role", &result.ID, "created", nil, &result.Status, nil, cmd.ActorID)
	}
	return result, err
}

func (s *TenantService) UpdateSuccessionCriticalRole(ctx context.Context, cmd ports.SuccessionCriticalRoleCommand) (*domain.SuccessionCriticalRole, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSuccessionCriticalRole
	}
	if cmd.CycleID != nil {
		if _, err := s.successionPlanning.GetSuccessionReviewCycle(ctx, cmd.TenantID, *cmd.CycleID); err != nil {
			return nil, err
		}
	}
	if err := s.validateSuccessionWorkerRefs(ctx, cmd.TenantID, cmd.IncumbentWorkerProfileID, cmd.EmergencyCoverWorkerProfileID); err != nil {
		return nil, err
	}
	item := successionCriticalRoleFromCommand(cmd)
	item.ID = cmd.ID
	if err := domain.ValidateSuccessionCriticalRole(item); err != nil {
		return nil, err
	}
	return s.successionPlanning.UpdateSuccessionCriticalRole(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateSuccessionCriticalRoleStatus(ctx context.Context, cmd ports.SuccessionStatusCommand) (*domain.SuccessionCriticalRole, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSuccessionCriticalRole
	}
	status := domain.NormalizeSuccessionRoleStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidSuccessionCriticalRole
	}
	before, _ := s.successionPlanning.GetSuccessionCriticalRole(ctx, cmd.TenantID, cmd.ID)
	result, err := s.successionPlanning.UpdateSuccessionCriticalRoleStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	_, _ = s.successionEvent(ctx, cmd.TenantID, "critical_role", &cmd.ID, "status_changed", from, &result.Status, cmd.Remarks, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListSuccessionCriticalRoles(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionCriticalRole, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeSuccessionPage(&filter.Limit, &filter.Offset)
	return s.successionPlanning.ListSuccessionCriticalRoles(ctx, filter)
}

func (s *TenantService) DeleteSuccessionCriticalRole(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidSuccessionCriticalRole
	}
	return s.successionPlanning.DeleteSuccessionCriticalRole(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateSuccessionSuccessorNomination(ctx context.Context, cmd ports.SuccessionSuccessorNominationCommand) (*domain.SuccessionSuccessorNomination, error) {
	if _, err := s.successionPlanning.GetSuccessionCriticalRole(ctx, cmd.TenantID, cmd.CriticalRoleID); err != nil {
		return nil, err
	}
	if _, err := s.workerProfiles.GetWorkerProfile(ctx, cmd.TenantID, cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	item := successionNominationFromCommand(cmd)
	if item.NominatedBy == nil {
		item.NominatedBy = cmd.ActorID
	}
	if err := domain.ValidateSuccessionSuccessorNomination(item); err != nil {
		return nil, err
	}
	result, err := s.successionPlanning.CreateSuccessionSuccessorNomination(ctx, item, cmd.ActorID)
	if err == nil {
		_, _ = s.successionEvent(ctx, result.TenantID, "successor_nomination", &result.ID, "created", nil, &result.NominationStatus, nil, cmd.ActorID)
	}
	return result, err
}

func (s *TenantService) UpdateSuccessionSuccessorNomination(ctx context.Context, cmd ports.SuccessionSuccessorNominationCommand) (*domain.SuccessionSuccessorNomination, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSuccessionNomination
	}
	item := successionNominationFromCommand(cmd)
	item.ID = cmd.ID
	if err := domain.ValidateSuccessionSuccessorNomination(item); err != nil {
		return nil, err
	}
	return s.successionPlanning.UpdateSuccessionSuccessorNomination(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateSuccessionSuccessorNominationStatus(ctx context.Context, cmd ports.SuccessionStatusCommand) (*domain.SuccessionSuccessorNomination, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSuccessionNomination
	}
	status := domain.NormalizeSuccessionNominationStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidSuccessionNomination
	}
	result, err := s.successionPlanning.UpdateSuccessionSuccessorNominationStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.successionEvent(ctx, result.TenantID, "successor_nomination", &result.ID, "status_changed", nil, &result.NominationStatus, cmd.Remarks, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListSuccessionSuccessorNominations(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionSuccessorNomination, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeSuccessionPage(&filter.Limit, &filter.Offset)
	return s.successionPlanning.ListSuccessionSuccessorNominations(ctx, filter)
}

func (s *TenantService) CreateSuccessionDevelopmentAction(ctx context.Context, cmd ports.SuccessionDevelopmentActionCommand) (*domain.SuccessionDevelopmentAction, error) {
	if _, err := s.workerProfiles.GetWorkerProfile(ctx, cmd.TenantID, cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	item := successionDevelopmentActionFromCommand(cmd)
	if err := domain.ValidateSuccessionDevelopmentAction(item); err != nil {
		return nil, err
	}
	result, err := s.successionPlanning.CreateSuccessionDevelopmentAction(ctx, item, cmd.ActorID)
	if err == nil {
		_, _ = s.successionEvent(ctx, result.TenantID, "development_action", &result.ID, "created", nil, &result.Status, nil, cmd.ActorID)
	}
	return result, err
}

func (s *TenantService) UpdateSuccessionDevelopmentAction(ctx context.Context, cmd ports.SuccessionDevelopmentActionCommand) (*domain.SuccessionDevelopmentAction, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSuccessionDevelopmentAction
	}
	item := successionDevelopmentActionFromCommand(cmd)
	item.ID = cmd.ID
	if err := domain.ValidateSuccessionDevelopmentAction(item); err != nil {
		return nil, err
	}
	return s.successionPlanning.UpdateSuccessionDevelopmentAction(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateSuccessionDevelopmentActionStatus(ctx context.Context, cmd ports.SuccessionStatusCommand) (*domain.SuccessionDevelopmentAction, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidSuccessionDevelopmentAction
	}
	status := domain.NormalizeSuccessionActionStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidSuccessionDevelopmentAction
	}
	result, err := s.successionPlanning.UpdateSuccessionDevelopmentActionStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.successionEvent(ctx, result.TenantID, "development_action", &result.ID, "status_changed", nil, &result.Status, cmd.Remarks, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListSuccessionDevelopmentActions(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionDevelopmentAction, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeSuccessionPage(&filter.Limit, &filter.Offset)
	return s.successionPlanning.ListSuccessionDevelopmentActions(ctx, filter)
}

func (s *TenantService) ListSuccessionEvents(ctx context.Context, filter domain.SuccessionFilter, sourceType *string, sourceID *uuid.UUID) ([]*domain.SuccessionEvent, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeSuccessionPage(&filter.Limit, &filter.Offset)
	return s.successionPlanning.ListSuccessionEvents(ctx, filter, sourceType, sourceID)
}

func (s *TenantService) GetSuccessionSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.SuccessionSummaryRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.successionPlanning.GetSuccessionSummary(ctx, tenantID)
}

func successionCriticalRoleFromCommand(cmd ports.SuccessionCriticalRoleCommand) *domain.SuccessionCriticalRole {
	return &domain.SuccessionCriticalRole{TenantID: cmd.TenantID, CycleID: cmd.CycleID, Code: cmd.Code, Title: cmd.Title, DepartmentID: cmd.DepartmentID, DesignationID: cmd.DesignationID, IncumbentWorkerProfileID: cmd.IncumbentWorkerProfileID, EmergencyCoverWorkerProfileID: cmd.EmergencyCoverWorkerProfileID, Criticality: cmd.Criticality, ImpactLevel: cmd.ImpactLevel, VacancyRisk: cmd.VacancyRisk, AttritionRisk: cmd.AttritionRisk, ReadinessTarget: cmd.ReadinessTarget, SuccessorRequiredCount: cmd.SuccessorRequiredCount, RoleSummary: cmd.RoleSummary, Status: cmd.Status, Metadata: cmd.Metadata}
}

func successionNominationFromCommand(cmd ports.SuccessionSuccessorNominationCommand) *domain.SuccessionSuccessorNomination {
	return &domain.SuccessionSuccessorNomination{TenantID: cmd.TenantID, CriticalRoleID: cmd.CriticalRoleID, WorkerProfileID: cmd.WorkerProfileID, NominatedBy: cmd.NominatedBy, ReadinessLevel: cmd.ReadinessLevel, ReadinessMonths: cmd.ReadinessMonths, PotentialRating: cmd.PotentialRating, PerformanceRating: cmd.PerformanceRating, RetentionRisk: cmd.RetentionRisk, MobilityPreference: cmd.MobilityPreference, NominationStatus: cmd.NominationStatus, DevelopmentNotes: cmd.DevelopmentNotes, Metadata: cmd.Metadata}
}

func successionDevelopmentActionFromCommand(cmd ports.SuccessionDevelopmentActionCommand) *domain.SuccessionDevelopmentAction {
	return &domain.SuccessionDevelopmentAction{TenantID: cmd.TenantID, NominationID: cmd.NominationID, CriticalRoleID: cmd.CriticalRoleID, WorkerProfileID: cmd.WorkerProfileID, ActionType: cmd.ActionType, Title: cmd.Title, LearningCourseID: cmd.LearningCourseID, LearningPathID: cmd.LearningPathID, OwnerUserID: cmd.OwnerUserID, DueDate: cmd.DueDate, Status: cmd.Status, Notes: cmd.Notes, Metadata: cmd.Metadata}
}

func (s *TenantService) successionEvent(ctx context.Context, tenantID uuid.UUID, sourceType string, sourceID *uuid.UUID, action string, fromStatus *string, toStatus *string, remarks *string, actorID *uuid.UUID) (*domain.SuccessionEvent, error) {
	return s.successionPlanning.CreateSuccessionEvent(ctx, &domain.SuccessionEvent{TenantID: tenantID, SourceType: sourceType, SourceID: sourceID, Action: action, FromStatus: fromStatus, ToStatus: toStatus, Remarks: remarks}, actorID)
}

func (s *TenantService) validateSuccessionWorkerRefs(ctx context.Context, tenantID uuid.UUID, workerIDs ...*uuid.UUID) error {
	for _, id := range workerIDs {
		if id == nil || *id == uuid.Nil {
			continue
		}
		if _, err := s.workerProfiles.GetWorkerProfile(ctx, tenantID, *id); err != nil {
			return err
		}
	}
	return nil
}

func normalizeSuccessionPage(limit *int32, offset *int32) {
	if *limit <= 0 || *limit > 200 {
		*limit = 50
	}
	if *offset < 0 {
		*offset = 0
	}
}
