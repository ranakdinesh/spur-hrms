package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateOKRCycle(ctx context.Context, cmd ports.OKRCycleCommand) (*domain.OKRCycle, error) {
	item, err := s.prepareOKRCycle(cmd)
	if err != nil {
		s.logError("validate okr cycle", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.okrs.CreateOKRCycle(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateOKRCycle(ctx context.Context, cmd ports.OKRCycleCommand) (*domain.OKRCycle, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidOKRCycle
	}
	item, err := s.prepareOKRCycle(cmd)
	if err != nil {
		s.logError("validate okr cycle update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("cycle_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	return s.okrs.UpdateOKRCycle(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateOKRCycleStatus(ctx context.Context, cmd ports.OKRStatusCommand) (*domain.OKRCycle, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || !stringIn(cmd.Status, []string{domain.OKRCycleStatusDraft, domain.OKRCycleStatusActive, domain.OKRCycleStatusClosed, domain.OKRCycleStatusArchived}) {
		return nil, domain.ErrInvalidOKRCycle
	}
	return s.okrs.UpdateOKRCycleStatus(ctx, cmd.TenantID, cmd.ID, cmd.Status, cmd.ActorID)
}

func (s *TenantService) GetOKRCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OKRCycle, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidOKRCycle
	}
	return s.okrs.GetOKRCycle(ctx, tenantID, id)
}

func (s *TenantService) ListOKRCycles(ctx context.Context, filter domain.OKRCycleFilter) ([]*domain.OKRCycle, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.Search = domain.NormalizeOKRSearch(filter.Search)
	return s.okrs.ListOKRCycles(ctx, filter)
}

func (s *TenantService) DeleteOKRCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidOKRCycle
	}
	return s.okrs.DeleteOKRCycle(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateObjective(ctx context.Context, cmd ports.ObjectiveCommand) (*domain.Objective, error) {
	item, err := s.prepareObjective(ctx, cmd)
	if err != nil {
		s.logError("validate objective", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return s.okrs.CreateObjective(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateObjective(ctx context.Context, cmd ports.ObjectiveCommand) (*domain.Objective, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidObjective
	}
	item, err := s.prepareObjective(ctx, cmd)
	if err != nil {
		s.logError("validate objective update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("objective_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	return s.okrs.UpdateObjective(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateObjectiveStatus(ctx context.Context, cmd ports.OKRStatusCommand) (*domain.Objective, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || !stringIn(cmd.Status, []string{domain.ObjectiveStatusDraft, domain.ObjectiveStatusActive, domain.ObjectiveStatusAtRisk, domain.ObjectiveStatusCompleted, domain.ObjectiveStatusClosed, domain.ObjectiveStatusCancelled}) {
		return nil, domain.ErrInvalidObjective
	}
	return s.okrs.UpdateObjectiveStatus(ctx, cmd.TenantID, cmd.ID, cmd.Status, cmd.ActorID)
}

func (s *TenantService) GetObjective(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Objective, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidObjective
	}
	return s.okrs.GetObjective(ctx, tenantID, id)
}

func (s *TenantService) ListObjectives(ctx context.Context, filter domain.ObjectiveFilter) ([]*domain.Objective, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.Search = domain.NormalizeOKRSearch(filter.Search)
	return s.okrs.ListObjectives(ctx, filter)
}

func (s *TenantService) DeleteObjective(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidObjective
	}
	return s.okrs.DeleteObjective(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateKeyResult(ctx context.Context, cmd ports.KeyResultCommand) (*domain.KeyResult, error) {
	item, err := s.prepareKeyResult(ctx, cmd)
	if err != nil {
		s.logError("validate key result", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.okrs.CreateKeyResult(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.okrs.RefreshObjectiveProgress(ctx, result.TenantID, result.ObjectiveID, cmd.ActorID)
	return result, nil
}

func (s *TenantService) UpdateKeyResult(ctx context.Context, cmd ports.KeyResultCommand) (*domain.KeyResult, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidKeyResult
	}
	item, err := s.prepareKeyResult(ctx, cmd)
	if err != nil {
		s.logError("validate key result update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("key_result_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.okrs.UpdateKeyResult(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.okrs.RefreshObjectiveProgress(ctx, result.TenantID, result.ObjectiveID, cmd.ActorID)
	return result, nil
}

func (s *TenantService) GetKeyResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.KeyResult, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidKeyResult
	}
	return s.okrs.GetKeyResult(ctx, tenantID, id)
}

func (s *TenantService) ListKeyResults(ctx context.Context, filter domain.KeyResultFilter) ([]*domain.KeyResult, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.okrs.ListKeyResults(ctx, filter)
}

func (s *TenantService) DeleteKeyResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidKeyResult
	}
	existing, err := s.okrs.GetKeyResult(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if err := s.okrs.DeleteKeyResult(ctx, tenantID, id, actorID); err != nil {
		return err
	}
	_, _ = s.okrs.RefreshObjectiveProgress(ctx, tenantID, existing.ObjectiveID, actorID)
	return nil
}

func (s *TenantService) CreateKeyResultCheckIn(ctx context.Context, cmd ports.KeyResultCheckInCommand) (*domain.KeyResultCheckIn, error) {
	existing, err := s.GetKeyResult(ctx, cmd.TenantID, cmd.KeyResultID)
	if err != nil {
		return nil, err
	}
	checkInDate, err := parseWorkerProfileDate(cmd.CheckInDate)
	if err != nil {
		return nil, err
	}
	progress := cmd.ProgressPercent
	if progress == nil {
		calculated := domain.CalculateKeyResultProgress(existing.MetricType, existing.StartValue, existing.TargetValue, cmd.Value)
		progress = &calculated
	}
	item, err := domain.NewKeyResultCheckIn(domain.KeyResultCheckInInput{TenantID: cmd.TenantID, KeyResultID: cmd.KeyResultID, CheckInDate: checkInDate, Value: cmd.Value, ProgressPercent: progress, Confidence: cmd.Confidence, Status: cmd.Status, Note: cmd.Note, Metadata: cmd.Metadata})
	if err != nil {
		s.logError("validate key result check-in", err, serviceTenantIDField(cmd.TenantID), serviceStringField("key_result_id", cmd.KeyResultID.String()))
		return nil, err
	}
	result, err := s.okrs.CreateKeyResultCheckIn(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.okrs.UpdateKeyResultProgress(ctx, cmd.TenantID, cmd.KeyResultID, result.Value, result.ProgressPercent, result.Confidence, result.Status, cmd.ActorID)
	_, _ = s.okrs.RefreshObjectiveProgress(ctx, cmd.TenantID, existing.ObjectiveID, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListKeyResultCheckIns(ctx context.Context, filter domain.KeyResultCheckInFilter) ([]*domain.KeyResultCheckIn, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.okrs.ListKeyResultCheckIns(ctx, filter)
}

func (s *TenantService) GetOKRSummary(ctx context.Context, tenantID uuid.UUID, cycleID *uuid.UUID) ([]*domain.OKRSummaryRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.okrs.GetOKRSummary(ctx, tenantID, cycleID)
}

func (s *TenantService) prepareOKRCycle(cmd ports.OKRCycleCommand) (*domain.OKRCycle, error) {
	startDate, err := parseWorkerProfileDate(cmd.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := parseWorkerProfileDate(cmd.EndDate)
	if err != nil {
		return nil, err
	}
	return domain.NewOKRCycle(domain.OKRCycleInput{TenantID: cmd.TenantID, Name: cmd.Name, CycleCode: cmd.CycleCode, Description: cmd.Description, StartDate: startDate, EndDate: endDate, Status: cmd.Status, ReviewCadence: cmd.ReviewCadence, Metadata: cmd.Metadata})
}

func (s *TenantService) prepareObjective(ctx context.Context, cmd ports.ObjectiveCommand) (*domain.Objective, error) {
	if _, err := s.GetOKRCycle(ctx, cmd.TenantID, cmd.CycleID); err != nil {
		return nil, err
	}
	if cmd.ParentObjectiveID != nil {
		if _, err := s.GetObjective(ctx, cmd.TenantID, *cmd.ParentObjectiveID); err != nil {
			return nil, err
		}
	}
	if cmd.OwnerWorkerProfileID != nil {
		if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, *cmd.OwnerWorkerProfileID); err != nil {
			return nil, err
		}
	}
	if cmd.OwnerDepartmentID != nil {
		if _, err := s.GetDepartment(ctx, cmd.TenantID, *cmd.OwnerDepartmentID); err != nil {
			return nil, err
		}
	}
	if cmd.OwnerProjectID != nil {
		if _, err := s.GetProject(ctx, cmd.TenantID, *cmd.OwnerProjectID); err != nil {
			return nil, err
		}
	}
	startDate, err := parseWorkerProfileDate(cmd.StartDate)
	if err != nil {
		return nil, err
	}
	dueDate, err := parseWorkerProfileDate(cmd.DueDate)
	if err != nil {
		return nil, err
	}
	return domain.NewObjective(domain.ObjectiveInput{TenantID: cmd.TenantID, CycleID: cmd.CycleID, ParentObjectiveID: cmd.ParentObjectiveID, OwnerType: cmd.OwnerType, OwnerWorkerProfileID: cmd.OwnerWorkerProfileID, OwnerDepartmentID: cmd.OwnerDepartmentID, OwnerProjectID: cmd.OwnerProjectID, Title: cmd.Title, Description: cmd.Description, Status: cmd.Status, Priority: cmd.Priority, ProgressPercent: cmd.ProgressPercent, Weight: cmd.Weight, StartDate: startDate, DueDate: dueDate, Metadata: cmd.Metadata})
}

func (s *TenantService) prepareKeyResult(ctx context.Context, cmd ports.KeyResultCommand) (*domain.KeyResult, error) {
	if _, err := s.GetObjective(ctx, cmd.TenantID, cmd.ObjectiveID); err != nil {
		return nil, err
	}
	dueDate, err := parseWorkerProfileDate(cmd.DueDate)
	if err != nil {
		return nil, err
	}
	return domain.NewKeyResult(domain.KeyResultInput{TenantID: cmd.TenantID, ObjectiveID: cmd.ObjectiveID, Title: cmd.Title, Description: cmd.Description, MetricType: cmd.MetricType, StartValue: cmd.StartValue, TargetValue: cmd.TargetValue, CurrentValue: cmd.CurrentValue, ProgressPercent: cmd.ProgressPercent, Confidence: cmd.Confidence, Status: cmd.Status, Weight: cmd.Weight, UnitLabel: cmd.UnitLabel, DueDate: dueDate, Metadata: cmd.Metadata})
}

func stringIn(value string, values []string) bool {
	for _, candidate := range values {
		if value == candidate {
			return true
		}
	}
	return false
}
