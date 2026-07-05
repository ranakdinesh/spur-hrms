package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateWorkerType(ctx context.Context, item *domain.WorkerType, actorID *uuid.UUID) (*domain.WorkerType, error) {
	row, err := s.getQueries(ctx).CreateWorkerType(ctx, sqlc.CreateWorkerTypeParams{
		TenantID:            item.TenantID,
		Code:                item.Code,
		Name:                item.Name,
		ClassificationGroup: item.ClassificationGroup,
		Description:         textFromPtr(item.Description),
		AttendanceMode:      item.AttendanceMode,
		PayMode:             item.PayMode,
		TdsSection:          item.TDSSection,
		PfApplicable:        item.PFApplicable,
		EsicApplicable:      item.ESICApplicable,
		PtApplicable:        item.PTApplicable,
		LwfApplicable:       item.LWFApplicable,
		ClraApplicable:      item.CLRAApplicable,
		LeaveApplicable:     item.LeaveApplicable,
		OvertimeApplicable:  item.OvertimeApplicable,
		RequiresAgreement:   item.RequiresAgreement,
		RequiresInvoice:     item.RequiresInvoice,
		RequiresAttendance:  item.RequiresAttendance,
		StatutoryDefaults:   jsonBytesFromRaw(item.StatutoryDefaults),
		ComplianceNotes:     textFromPtr(item.ComplianceNotes),
		IsSystemDefault:     item.IsSystemDefault,
		SortOrder:           item.SortOrder,
		CreatedBy:           uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create worker type", err, tenantIDField(item.TenantID), stringField("worker_type_code", item.Code))
	}
	return mapWorkerType(row), nil
}

func (s *Store) ListWorkerTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.WorkerType, error) {
	rows, err := s.getQueries(ctx).ListWorkerTypes(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list worker types", err, tenantIDField(tenantID))
	}
	return mapWorkerTypes(rows), nil
}

func (s *Store) GetWorkerType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerType, error) {
	row, err := s.getQueries(ctx).GetWorkerType(ctx, sqlc.GetWorkerTypeParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get worker type", err, tenantIDField(tenantID), stringField("worker_type_id", id.String()))
	}
	return mapWorkerType(row), nil
}

func (s *Store) GetWorkerTypeByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.WorkerType, error) {
	row, err := s.getQueries(ctx).GetWorkerTypeByCode(ctx, sqlc.GetWorkerTypeByCodeParams{TenantID: tenantID, Code: code})
	if err != nil {
		return nil, s.logDBError(ctx, "get worker type by code", err, tenantIDField(tenantID), stringField("worker_type_code", code))
	}
	return mapWorkerType(row), nil
}

func (s *Store) UpdateWorkerType(ctx context.Context, item *domain.WorkerType, actorID *uuid.UUID) (*domain.WorkerType, error) {
	row, err := s.getQueries(ctx).UpdateWorkerType(ctx, sqlc.UpdateWorkerTypeParams{
		TenantID:            item.TenantID,
		ID:                  item.ID,
		Name:                item.Name,
		ClassificationGroup: item.ClassificationGroup,
		Description:         textFromPtr(item.Description),
		AttendanceMode:      item.AttendanceMode,
		PayMode:             item.PayMode,
		TdsSection:          item.TDSSection,
		PfApplicable:        item.PFApplicable,
		EsicApplicable:      item.ESICApplicable,
		PtApplicable:        item.PTApplicable,
		LwfApplicable:       item.LWFApplicable,
		ClraApplicable:      item.CLRAApplicable,
		LeaveApplicable:     item.LeaveApplicable,
		OvertimeApplicable:  item.OvertimeApplicable,
		RequiresAgreement:   item.RequiresAgreement,
		RequiresInvoice:     item.RequiresInvoice,
		RequiresAttendance:  item.RequiresAttendance,
		StatutoryDefaults:   jsonBytesFromRaw(item.StatutoryDefaults),
		ComplianceNotes:     textFromPtr(item.ComplianceNotes),
		SortOrder:           item.SortOrder,
		UpdatedBy:           uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update worker type", err, tenantIDField(item.TenantID), stringField("worker_type_id", item.ID.String()))
	}
	return mapWorkerType(row), nil
}

func (s *Store) DeleteWorkerType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteWorkerType(ctx, sqlc.SoftDeleteWorkerTypeParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete worker type", err, tenantIDField(tenantID), stringField("worker_type_id", id.String()))
	}
	return nil
}

func (s *Store) CountWorkerClassificationRules(ctx context.Context, tenantID uuid.UUID, workerTypeID uuid.UUID) (int64, error) {
	count, err := s.getQueries(ctx).CountWorkerClassificationRules(ctx, sqlc.CountWorkerClassificationRulesParams{TenantID: tenantID, WorkerTypeID: workerTypeID})
	if err != nil {
		return 0, s.logDBError(ctx, "count worker classification rules", err, tenantIDField(tenantID), stringField("worker_type_id", workerTypeID.String()))
	}
	return count, nil
}

func (s *Store) CreateWorkerClassificationRule(ctx context.Context, item *domain.WorkerClassificationRule, actorID *uuid.UUID) (*domain.WorkerClassificationRule, error) {
	row, err := s.getQueries(ctx).CreateWorkerClassificationRule(ctx, sqlc.CreateWorkerClassificationRuleParams{
		TenantID:     item.TenantID,
		WorkerTypeID: item.WorkerTypeID,
		RuleName:     item.RuleName,
		RuleType:     item.RuleType,
		Priority:     item.Priority,
		Conditions:   jsonBytesFromRaw(item.Conditions),
		Outcome:      jsonBytesFromRaw(item.Outcome),
		Notes:        textFromPtr(item.Notes),
		CreatedBy:    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create worker classification rule", err, tenantIDField(item.TenantID), stringField("worker_type_id", item.WorkerTypeID.String()))
	}
	return mapWorkerClassificationRule(row), nil
}

func (s *Store) ListWorkerClassificationRules(ctx context.Context, tenantID uuid.UUID, workerTypeID *uuid.UUID) ([]*domain.WorkerClassificationRule, error) {
	rows, err := s.getQueries(ctx).ListWorkerClassificationRules(ctx, sqlc.ListWorkerClassificationRulesParams{TenantID: tenantID, WorkerTypeID: uuidFromPtr(workerTypeID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list worker classification rules", err, tenantIDField(tenantID))
	}
	return mapWorkerClassificationRules(rows), nil
}

func (s *Store) GetWorkerClassificationRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerClassificationRule, error) {
	row, err := s.getQueries(ctx).GetWorkerClassificationRule(ctx, sqlc.GetWorkerClassificationRuleParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get worker classification rule", err, tenantIDField(tenantID), stringField("rule_id", id.String()))
	}
	return mapWorkerClassificationRule(row), nil
}

func (s *Store) UpdateWorkerClassificationRule(ctx context.Context, item *domain.WorkerClassificationRule, actorID *uuid.UUID) (*domain.WorkerClassificationRule, error) {
	row, err := s.getQueries(ctx).UpdateWorkerClassificationRule(ctx, sqlc.UpdateWorkerClassificationRuleParams{
		TenantID:     item.TenantID,
		ID:           item.ID,
		WorkerTypeID: item.WorkerTypeID,
		RuleName:     item.RuleName,
		RuleType:     item.RuleType,
		Priority:     item.Priority,
		Conditions:   jsonBytesFromRaw(item.Conditions),
		Outcome:      jsonBytesFromRaw(item.Outcome),
		Notes:        textFromPtr(item.Notes),
		UpdatedBy:    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update worker classification rule", err, tenantIDField(item.TenantID), stringField("rule_id", item.ID.String()))
	}
	return mapWorkerClassificationRule(row), nil
}

func (s *Store) DeleteWorkerClassificationRule(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteWorkerClassificationRule(ctx, sqlc.SoftDeleteWorkerClassificationRuleParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete worker classification rule", err, tenantIDField(tenantID), stringField("rule_id", id.String()))
	}
	return nil
}
