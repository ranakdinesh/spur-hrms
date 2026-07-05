package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateOKRCycle(ctx context.Context, item *domain.OKRCycle, actorID *uuid.UUID) (*domain.OKRCycle, error) {
	row, err := s.getQueries(ctx).CreateOKRCycle(ctx, sqlc.CreateOKRCycleParams{TenantID: item.TenantID, Name: item.Name, CycleCode: item.CycleCode, Description: textFromPtr(item.Description), StartDate: dateFromPtr(&item.StartDate), EndDate: dateFromPtr(&item.EndDate), Status: item.Status, ReviewCadence: item.ReviewCadence, Column9: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create okr cycle", err, tenantIDField(item.TenantID), stringField("cycle_code", item.CycleCode))
	}
	return mapOKRCycle(row), nil
}

func (s *Store) UpdateOKRCycle(ctx context.Context, item *domain.OKRCycle, actorID *uuid.UUID) (*domain.OKRCycle, error) {
	row, err := s.getQueries(ctx).UpdateOKRCycle(ctx, sqlc.UpdateOKRCycleParams{TenantID: item.TenantID, ID: item.ID, Name: item.Name, CycleCode: item.CycleCode, Description: textFromPtr(item.Description), StartDate: dateFromPtr(&item.StartDate), EndDate: dateFromPtr(&item.EndDate), Status: item.Status, ReviewCadence: item.ReviewCadence, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOKRCycleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update okr cycle", err, tenantIDField(item.TenantID), stringField("cycle_id", item.ID.String()))
	}
	return mapOKRCycle(row), nil
}

func (s *Store) UpdateOKRCycleStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.OKRCycle, error) {
	row, err := s.getQueries(ctx).UpdateOKRCycleStatus(ctx, sqlc.UpdateOKRCycleStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOKRCycleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update okr cycle status", err, tenantIDField(tenantID), stringField("cycle_id", id.String()))
	}
	return mapOKRCycle(row), nil
}

func (s *Store) GetOKRCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.OKRCycle, error) {
	row, err := s.getQueries(ctx).GetOKRCycle(ctx, sqlc.GetOKRCycleParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOKRCycleNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get okr cycle", err, tenantIDField(tenantID), stringField("cycle_id", id.String()))
	}
	return mapOKRCycle(row), nil
}

func (s *Store) ListOKRCycles(ctx context.Context, filter domain.OKRCycleFilter) ([]*domain.OKRCycle, error) {
	rows, err := s.getQueries(ctx).ListOKRCycles(ctx, sqlc.ListOKRCyclesParams{TenantID: filter.TenantID, Column2: stringFromPtr(filter.Status), Column3: stringFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list okr cycles", err, tenantIDField(filter.TenantID))
	}
	return mapOKRCycles(rows), nil
}

func (s *Store) DeleteOKRCycle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteOKRCycle(ctx, sqlc.SoftDeleteOKRCycleParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete okr cycle", err, tenantIDField(tenantID), stringField("cycle_id", id.String()))
	}
	return nil
}

func (s *Store) CreateObjective(ctx context.Context, item *domain.Objective, actorID *uuid.UUID) (*domain.Objective, error) {
	row, err := s.getQueries(ctx).CreateObjective(ctx, okrCreateObjectiveParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create objective", err, tenantIDField(item.TenantID), stringField("title", item.Title))
	}
	return mapObjective(row), nil
}

func (s *Store) UpdateObjective(ctx context.Context, item *domain.Objective, actorID *uuid.UUID) (*domain.Objective, error) {
	row, err := s.getQueries(ctx).UpdateObjective(ctx, okrUpdateObjectiveParams(item, actorID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrObjectiveNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update objective", err, tenantIDField(item.TenantID), stringField("objective_id", item.ID.String()))
	}
	return mapObjective(row), nil
}

func (s *Store) UpdateObjectiveStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.Objective, error) {
	row, err := s.getQueries(ctx).UpdateObjectiveStatus(ctx, sqlc.UpdateObjectiveStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrObjectiveNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update objective status", err, tenantIDField(tenantID), stringField("objective_id", id.String()))
	}
	return mapObjective(row), nil
}

func (s *Store) RefreshObjectiveProgress(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) (*domain.Objective, error) {
	row, err := s.getQueries(ctx).RefreshObjectiveProgress(ctx, sqlc.RefreshObjectiveProgressParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrObjectiveNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "refresh objective progress", err, tenantIDField(tenantID), stringField("objective_id", id.String()))
	}
	return mapObjective(row), nil
}

func (s *Store) GetObjective(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Objective, error) {
	row, err := s.getQueries(ctx).GetObjective(ctx, sqlc.GetObjectiveParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrObjectiveNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get objective", err, tenantIDField(tenantID), stringField("objective_id", id.String()))
	}
	return mapObjective(row), nil
}

func (s *Store) ListObjectives(ctx context.Context, filter domain.ObjectiveFilter) ([]*domain.Objective, error) {
	rows, err := s.getQueries(ctx).ListObjectives(ctx, sqlc.ListObjectivesParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.CycleID), Column3: stringFromPtr(filter.OwnerType), Column4: stringFromPtr(filter.Status), Column5: stringFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list objectives", err, tenantIDField(filter.TenantID))
	}
	return mapObjectiveRows(rows), nil
}

func (s *Store) DeleteObjective(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteObjective(ctx, sqlc.SoftDeleteObjectiveParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete objective", err, tenantIDField(tenantID), stringField("objective_id", id.String()))
	}
	return nil
}

func (s *Store) CreateKeyResult(ctx context.Context, item *domain.KeyResult, actorID *uuid.UUID) (*domain.KeyResult, error) {
	row, err := s.getQueries(ctx).CreateKeyResult(ctx, sqlc.CreateKeyResultParams{TenantID: item.TenantID, ObjectiveID: item.ObjectiveID, Title: item.Title, Description: textFromPtr(item.Description), MetricType: item.MetricType, StartValue: numericFromFloat(item.StartValue), TargetValue: numericFromFloat(item.TargetValue), CurrentValue: numericFromFloat(item.CurrentValue), Column9: numericFromFloat(item.ProgressPercent), Confidence: item.Confidence, Status: item.Status, Column12: numericFromFloat(item.Weight), UnitLabel: textFromPtr(item.UnitLabel), DueDate: dateFromPtr(item.DueDate), Column15: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create key result", err, tenantIDField(item.TenantID), stringField("objective_id", item.ObjectiveID.String()))
	}
	return mapKeyResult(row), nil
}

func (s *Store) UpdateKeyResult(ctx context.Context, item *domain.KeyResult, actorID *uuid.UUID) (*domain.KeyResult, error) {
	row, err := s.getQueries(ctx).UpdateKeyResult(ctx, sqlc.UpdateKeyResultParams{TenantID: item.TenantID, ID: item.ID, ObjectiveID: item.ObjectiveID, Title: item.Title, Description: textFromPtr(item.Description), MetricType: item.MetricType, StartValue: numericFromFloat(item.StartValue), TargetValue: numericFromFloat(item.TargetValue), CurrentValue: numericFromFloat(item.CurrentValue), ProgressPercent: numericFromFloat(item.ProgressPercent), Confidence: item.Confidence, Status: item.Status, Weight: numericFromFloat(item.Weight), UnitLabel: textFromPtr(item.UnitLabel), DueDate: dateFromPtr(item.DueDate), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrKeyResultNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update key result", err, tenantIDField(item.TenantID), stringField("key_result_id", item.ID.String()))
	}
	return mapKeyResult(row), nil
}

func (s *Store) UpdateKeyResultProgress(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, value float64, progress float64, confidence string, status string, actorID *uuid.UUID) (*domain.KeyResult, error) {
	row, err := s.getQueries(ctx).UpdateKeyResultProgress(ctx, sqlc.UpdateKeyResultProgressParams{TenantID: tenantID, ID: id, CurrentValue: numericFromFloat(value), ProgressPercent: numericFromFloat(progress), Confidence: confidence, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrKeyResultNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update key result progress", err, tenantIDField(tenantID), stringField("key_result_id", id.String()))
	}
	return mapKeyResult(row), nil
}

func (s *Store) GetKeyResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.KeyResult, error) {
	row, err := s.getQueries(ctx).GetKeyResult(ctx, sqlc.GetKeyResultParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrKeyResultNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get key result", err, tenantIDField(tenantID), stringField("key_result_id", id.String()))
	}
	return mapKeyResult(row), nil
}

func (s *Store) ListKeyResults(ctx context.Context, filter domain.KeyResultFilter) ([]*domain.KeyResult, error) {
	rows, err := s.getQueries(ctx).ListKeyResults(ctx, sqlc.ListKeyResultsParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.ObjectiveID), Column3: stringFromPtr(filter.Status)})
	if err != nil {
		return nil, s.logDBError(ctx, "list key results", err, tenantIDField(filter.TenantID))
	}
	return mapKeyResultRows(rows), nil
}

func (s *Store) DeleteKeyResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteKeyResult(ctx, sqlc.SoftDeleteKeyResultParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete key result", err, tenantIDField(tenantID), stringField("key_result_id", id.String()))
	}
	return nil
}

func (s *Store) CreateKeyResultCheckIn(ctx context.Context, item *domain.KeyResultCheckIn, actorID *uuid.UUID) (*domain.KeyResultCheckIn, error) {
	row, err := s.getQueries(ctx).CreateKeyResultCheckIn(ctx, sqlc.CreateKeyResultCheckInParams{TenantID: item.TenantID, KeyResultID: item.KeyResultID, CheckinDate: dateFromPtr(&item.CheckInDate), Value: numericFromFloat(item.Value), ProgressPercent: numericFromFloat(item.ProgressPercent), Confidence: item.Confidence, Status: item.Status, Note: textFromPtr(item.Note), Column9: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create key result check-in", err, tenantIDField(item.TenantID), stringField("key_result_id", item.KeyResultID.String()))
	}
	return mapKeyResultCheckIn(row), nil
}

func (s *Store) ListKeyResultCheckIns(ctx context.Context, filter domain.KeyResultCheckInFilter) ([]*domain.KeyResultCheckIn, error) {
	rows, err := s.getQueries(ctx).ListKeyResultCheckIns(ctx, sqlc.ListKeyResultCheckInsParams{TenantID: filter.TenantID, Column2: uuidValueFromPtr(filter.KeyResultID), Column3: uuidValueFromPtr(filter.ObjectiveID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list key result check-ins", err, tenantIDField(filter.TenantID))
	}
	return mapKeyResultCheckInRows(rows), nil
}

func (s *Store) GetOKRSummary(ctx context.Context, tenantID uuid.UUID, cycleID *uuid.UUID) ([]*domain.OKRSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetOKRSummary(ctx, sqlc.GetOKRSummaryParams{TenantID: tenantID, Column2: uuidValueFromPtr(cycleID)})
	if err != nil {
		return nil, s.logDBError(ctx, "get okr summary", err, tenantIDField(tenantID))
	}
	return mapOKRSummaryRows(rows), nil
}

func okrCreateObjectiveParams(item *domain.Objective, actorID *uuid.UUID) sqlc.CreateObjectiveParams {
	return sqlc.CreateObjectiveParams{TenantID: item.TenantID, CycleID: item.CycleID, ParentObjectiveID: uuidFromPtr(item.ParentObjectiveID), OwnerType: item.OwnerType, OwnerWorkerProfileID: uuidFromPtr(item.OwnerWorkerProfileID), OwnerDepartmentID: uuidFromPtr(item.OwnerDepartmentID), OwnerProjectID: uuidFromPtr(item.OwnerProjectID), Title: item.Title, Description: textFromPtr(item.Description), Status: item.Status, Priority: item.Priority, Column12: numericFromFloat(item.ProgressPercent), Column13: numericFromFloat(item.Weight), StartDate: dateFromPtr(item.StartDate), DueDate: dateFromPtr(item.DueDate), Column16: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func okrUpdateObjectiveParams(item *domain.Objective, actorID *uuid.UUID) sqlc.UpdateObjectiveParams {
	return sqlc.UpdateObjectiveParams{TenantID: item.TenantID, ID: item.ID, CycleID: item.CycleID, ParentObjectiveID: uuidFromPtr(item.ParentObjectiveID), OwnerType: item.OwnerType, OwnerWorkerProfileID: uuidFromPtr(item.OwnerWorkerProfileID), OwnerDepartmentID: uuidFromPtr(item.OwnerDepartmentID), OwnerProjectID: uuidFromPtr(item.OwnerProjectID), Title: item.Title, Description: textFromPtr(item.Description), Status: item.Status, Priority: item.Priority, ProgressPercent: numericFromFloat(item.ProgressPercent), Weight: numericFromFloat(item.Weight), StartDate: dateFromPtr(item.StartDate), DueDate: dateFromPtr(item.DueDate), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}

func stringFromPtr(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func uuidValueFromPtr(value *uuid.UUID) uuid.UUID {
	if value == nil {
		return uuid.Nil
	}
	return *value
}
