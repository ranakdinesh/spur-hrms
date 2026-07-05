package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateSuccessionReviewCycle(ctx context.Context, item *domain.SuccessionReviewCycle, actorID *uuid.UUID) (*domain.SuccessionReviewCycle, error) {
	row, err := s.getQueries(ctx).CreateSuccessionReviewCycle(ctx, sqlc.CreateSuccessionReviewCycleParams{TenantID: item.TenantID, Code: item.Code, Name: item.Name, Status: item.Status, StartsOn: dateFromPtr(item.StartsOn), EndsOn: dateFromPtr(item.EndsOn), ConfidentialityLevel: item.ConfidentialityLevel, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create succession review cycle", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapSuccessionReviewCycle(row), nil
}

func (s *Store) UpdateSuccessionReviewCycle(ctx context.Context, item *domain.SuccessionReviewCycle, actorID *uuid.UUID) (*domain.SuccessionReviewCycle, error) {
	row, err := s.getQueries(ctx).UpdateSuccessionReviewCycle(ctx, sqlc.UpdateSuccessionReviewCycleParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Name: item.Name, StartsOn: dateFromPtr(item.StartsOn), EndsOn: dateFromPtr(item.EndsOn), ConfidentialityLevel: item.ConfidentialityLevel, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionReviewCycleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update succession review cycle", err, tenantIDField(item.TenantID), stringField("cycle_id", item.ID.String()))
	}
	return mapSuccessionReviewCycle(row), nil
}

func (s *Store) UpdateSuccessionReviewCycleStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.SuccessionReviewCycle, error) {
	row, err := s.getQueries(ctx).UpdateSuccessionReviewCycleStatus(ctx, sqlc.UpdateSuccessionReviewCycleStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionReviewCycleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update succession review cycle status", err, tenantIDField(tenantID), stringField("cycle_id", id.String()))
	}
	return mapSuccessionReviewCycle(row), nil
}

func (s *Store) GetSuccessionReviewCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SuccessionReviewCycle, error) {
	row, err := s.getQueries(ctx).GetSuccessionReviewCycle(ctx, sqlc.GetSuccessionReviewCycleParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionReviewCycleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get succession review cycle", err, tenantIDField(tenantID), stringField("cycle_id", id.String()))
	}
	return mapSuccessionReviewCycle(row), nil
}

func (s *Store) ListSuccessionReviewCycles(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionReviewCycle, error) {
	rows, err := s.getQueries(ctx).ListSuccessionReviewCycles(ctx, sqlc.ListSuccessionReviewCyclesParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list succession review cycles", err, tenantIDField(filter.TenantID))
	}
	return mapSuccessionReviewCycles(rows), nil
}

func (s *Store) CreateSuccessionCriticalRole(ctx context.Context, item *domain.SuccessionCriticalRole, actorID *uuid.UUID) (*domain.SuccessionCriticalRole, error) {
	row, err := s.getQueries(ctx).CreateSuccessionCriticalRole(ctx, sqlc.CreateSuccessionCriticalRoleParams{TenantID: item.TenantID, CycleID: uuidFromPtr(item.CycleID), Code: item.Code, Title: item.Title, DepartmentID: uuidFromPtr(item.DepartmentID), DesignationID: uuidFromPtr(item.DesignationID), IncumbentWorkerProfileID: uuidFromPtr(item.IncumbentWorkerProfileID), EmergencyCoverWorkerProfileID: uuidFromPtr(item.EmergencyCoverWorkerProfileID), Criticality: item.Criticality, ImpactLevel: item.ImpactLevel, VacancyRisk: item.VacancyRisk, AttritionRisk: item.AttritionRisk, ReadinessTarget: item.ReadinessTarget, SuccessorRequiredCount: item.SuccessorRequiredCount, RoleSummary: textFromPtr(item.RoleSummary), Status: item.Status, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create succession critical role", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapSuccessionCriticalRole(row), nil
}

func (s *Store) UpdateSuccessionCriticalRole(ctx context.Context, item *domain.SuccessionCriticalRole, actorID *uuid.UUID) (*domain.SuccessionCriticalRole, error) {
	row, err := s.getQueries(ctx).UpdateSuccessionCriticalRole(ctx, sqlc.UpdateSuccessionCriticalRoleParams{TenantID: item.TenantID, ID: item.ID, CycleID: uuidFromPtr(item.CycleID), Code: item.Code, Title: item.Title, DepartmentID: uuidFromPtr(item.DepartmentID), DesignationID: uuidFromPtr(item.DesignationID), IncumbentWorkerProfileID: uuidFromPtr(item.IncumbentWorkerProfileID), EmergencyCoverWorkerProfileID: uuidFromPtr(item.EmergencyCoverWorkerProfileID), Criticality: item.Criticality, ImpactLevel: item.ImpactLevel, VacancyRisk: item.VacancyRisk, AttritionRisk: item.AttritionRisk, ReadinessTarget: item.ReadinessTarget, SuccessorRequiredCount: item.SuccessorRequiredCount, RoleSummary: textFromPtr(item.RoleSummary), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionCriticalRoleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update succession critical role", err, tenantIDField(item.TenantID), stringField("role_id", item.ID.String()))
	}
	return mapSuccessionCriticalRole(row), nil
}

func (s *Store) UpdateSuccessionCriticalRoleStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.SuccessionCriticalRole, error) {
	row, err := s.getQueries(ctx).UpdateSuccessionCriticalRoleStatus(ctx, sqlc.UpdateSuccessionCriticalRoleStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionCriticalRoleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update succession critical role status", err, tenantIDField(tenantID), stringField("role_id", id.String()))
	}
	return mapSuccessionCriticalRole(row), nil
}

func (s *Store) GetSuccessionCriticalRole(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.SuccessionCriticalRole, error) {
	row, err := s.getQueries(ctx).GetSuccessionCriticalRole(ctx, sqlc.GetSuccessionCriticalRoleParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionCriticalRoleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get succession critical role", err, tenantIDField(tenantID), stringField("role_id", id.String()))
	}
	return mapSuccessionCriticalRoleDetail(row), nil
}

func (s *Store) ListSuccessionCriticalRoles(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionCriticalRole, error) {
	rows, err := s.getQueries(ctx).ListSuccessionCriticalRoles(ctx, sqlc.ListSuccessionCriticalRolesParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, CycleID: uuidFromPtr(filter.CycleID), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list succession critical roles", err, tenantIDField(filter.TenantID))
	}
	return mapSuccessionCriticalRoleList(rows), nil
}

func (s *Store) DeleteSuccessionCriticalRole(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteSuccessionCriticalRole(ctx, sqlc.SoftDeleteSuccessionCriticalRoleParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete succession critical role", err, tenantIDField(tenantID), stringField("role_id", id.String()))
	}
	return nil
}

func (s *Store) CreateSuccessionSuccessorNomination(ctx context.Context, item *domain.SuccessionSuccessorNomination, actorID *uuid.UUID) (*domain.SuccessionSuccessorNomination, error) {
	row, err := s.getQueries(ctx).CreateSuccessionSuccessorNomination(ctx, sqlc.CreateSuccessionSuccessorNominationParams{TenantID: item.TenantID, CriticalRoleID: item.CriticalRoleID, WorkerProfileID: item.WorkerProfileID, NominatedBy: uuidFromPtr(item.NominatedBy), ReadinessLevel: item.ReadinessLevel, ReadinessMonths: item.ReadinessMonths, PotentialRating: textFromPtr(item.PotentialRating), PerformanceRating: textFromPtr(item.PerformanceRating), RetentionRisk: item.RetentionRisk, MobilityPreference: textFromPtr(item.MobilityPreference), NominationStatus: item.NominationStatus, DevelopmentNotes: textFromPtr(item.DevelopmentNotes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create succession nomination", err, tenantIDField(item.TenantID), stringField("role_id", item.CriticalRoleID.String()))
	}
	return mapSuccessionNomination(row), nil
}

func (s *Store) UpdateSuccessionSuccessorNomination(ctx context.Context, item *domain.SuccessionSuccessorNomination, actorID *uuid.UUID) (*domain.SuccessionSuccessorNomination, error) {
	row, err := s.getQueries(ctx).UpdateSuccessionSuccessorNomination(ctx, sqlc.UpdateSuccessionSuccessorNominationParams{TenantID: item.TenantID, ID: item.ID, CriticalRoleID: item.CriticalRoleID, ReadinessLevel: item.ReadinessLevel, ReadinessMonths: item.ReadinessMonths, PotentialRating: textFromPtr(item.PotentialRating), PerformanceRating: textFromPtr(item.PerformanceRating), RetentionRisk: item.RetentionRisk, MobilityPreference: textFromPtr(item.MobilityPreference), DevelopmentNotes: textFromPtr(item.DevelopmentNotes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionNominationNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update succession nomination", err, tenantIDField(item.TenantID), stringField("nomination_id", item.ID.String()))
	}
	return mapSuccessionNomination(row), nil
}

func (s *Store) UpdateSuccessionSuccessorNominationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.SuccessionSuccessorNomination, error) {
	row, err := s.getQueries(ctx).UpdateSuccessionSuccessorNominationStatus(ctx, sqlc.UpdateSuccessionSuccessorNominationStatusParams{TenantID: tenantID, ID: id, NominationStatus: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionNominationNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update succession nomination status", err, tenantIDField(tenantID), stringField("nomination_id", id.String()))
	}
	return mapSuccessionNomination(row), nil
}

func (s *Store) ListSuccessionSuccessorNominations(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionSuccessorNomination, error) {
	rows, err := s.getQueries(ctx).ListSuccessionSuccessorNominations(ctx, sqlc.ListSuccessionSuccessorNominationsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, CriticalRoleID: uuidFromPtr(filter.CriticalRoleID), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list succession nominations", err, tenantIDField(filter.TenantID))
	}
	return mapSuccessionNominationList(rows), nil
}

func (s *Store) CreateSuccessionDevelopmentAction(ctx context.Context, item *domain.SuccessionDevelopmentAction, actorID *uuid.UUID) (*domain.SuccessionDevelopmentAction, error) {
	row, err := s.getQueries(ctx).CreateSuccessionDevelopmentAction(ctx, sqlc.CreateSuccessionDevelopmentActionParams{TenantID: item.TenantID, NominationID: uuidFromPtr(item.NominationID), CriticalRoleID: uuidFromPtr(item.CriticalRoleID), WorkerProfileID: item.WorkerProfileID, ActionType: item.ActionType, Title: item.Title, LearningCourseID: uuidFromPtr(item.LearningCourseID), LearningPathID: uuidFromPtr(item.LearningPathID), OwnerUserID: uuidFromPtr(item.OwnerUserID), DueDate: dateFromPtr(item.DueDate), Status: item.Status, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create succession development action", err, tenantIDField(item.TenantID), stringField("worker_profile_id", item.WorkerProfileID.String()))
	}
	return mapSuccessionDevelopmentAction(row), nil
}

func (s *Store) UpdateSuccessionDevelopmentAction(ctx context.Context, item *domain.SuccessionDevelopmentAction, actorID *uuid.UUID) (*domain.SuccessionDevelopmentAction, error) {
	row, err := s.getQueries(ctx).UpdateSuccessionDevelopmentAction(ctx, sqlc.UpdateSuccessionDevelopmentActionParams{TenantID: item.TenantID, ID: item.ID, NominationID: uuidFromPtr(item.NominationID), CriticalRoleID: uuidFromPtr(item.CriticalRoleID), WorkerProfileID: item.WorkerProfileID, ActionType: item.ActionType, Title: item.Title, LearningCourseID: uuidFromPtr(item.LearningCourseID), LearningPathID: uuidFromPtr(item.LearningPathID), OwnerUserID: uuidFromPtr(item.OwnerUserID), DueDate: dateFromPtr(item.DueDate), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionDevelopmentActionNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update succession development action", err, tenantIDField(item.TenantID), stringField("action_id", item.ID.String()))
	}
	return mapSuccessionDevelopmentAction(row), nil
}

func (s *Store) UpdateSuccessionDevelopmentActionStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.SuccessionDevelopmentAction, error) {
	row, err := s.getQueries(ctx).UpdateSuccessionDevelopmentActionStatus(ctx, sqlc.UpdateSuccessionDevelopmentActionStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrSuccessionDevelopmentActionNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update succession development action status", err, tenantIDField(tenantID), stringField("action_id", id.String()))
	}
	return mapSuccessionDevelopmentAction(row), nil
}

func (s *Store) ListSuccessionDevelopmentActions(ctx context.Context, filter domain.SuccessionFilter) ([]*domain.SuccessionDevelopmentAction, error) {
	rows, err := s.getQueries(ctx).ListSuccessionDevelopmentActions(ctx, sqlc.ListSuccessionDevelopmentActionsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, CriticalRoleID: uuidFromPtr(filter.CriticalRoleID), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list succession development actions", err, tenantIDField(filter.TenantID))
	}
	return mapSuccessionDevelopmentActionList(rows), nil
}

func (s *Store) CreateSuccessionEvent(ctx context.Context, item *domain.SuccessionEvent, actorID *uuid.UUID) (*domain.SuccessionEvent, error) {
	row, err := s.getQueries(ctx).CreateSuccessionEvent(ctx, sqlc.CreateSuccessionEventParams{TenantID: item.TenantID, SourceType: item.SourceType, SourceID: uuidFromPtr(item.SourceID), Action: item.Action, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), Remarks: textFromPtr(item.Remarks), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create succession event", err, tenantIDField(item.TenantID), stringField("source_type", item.SourceType))
	}
	return mapSuccessionEvent(row), nil
}

func (s *Store) ListSuccessionEvents(ctx context.Context, filter domain.SuccessionFilter, sourceType *string, sourceID *uuid.UUID) ([]*domain.SuccessionEvent, error) {
	rows, err := s.getQueries(ctx).ListSuccessionEvents(ctx, sqlc.ListSuccessionEventsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, SourceType: textFromPtr(sourceType), SourceID: uuidFromPtr(sourceID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list succession events", err, tenantIDField(filter.TenantID))
	}
	return mapSuccessionEvents(rows), nil
}

func (s *Store) GetSuccessionSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.SuccessionSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetSuccessionSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get succession summary", err, tenantIDField(tenantID))
	}
	return mapSuccessionSummary(rows), nil
}
