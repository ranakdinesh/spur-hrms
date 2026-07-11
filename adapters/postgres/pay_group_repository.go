package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreatePayGroup(ctx context.Context, item *domain.PayGroup, actorID *uuid.UUID) (*domain.PayGroup, error) {
	row, err := s.getQueries(ctx).CreatePayGroup(ctx, sqlc.CreatePayGroupParams{
		TenantID:         item.TenantID,
		Code:             item.Code,
		Name:             item.Name,
		Description:      textFromPtr(item.Description),
		GroupingType:     item.GroupingType,
		BranchID:         uuidFromPtr(item.BranchID),
		DepartmentID:     uuidFromPtr(item.DepartmentID),
		EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID),
		ReportingTag:     textFromPtr(item.ReportingTag),
		Rules:            []byte(item.Rules),
		IsActive:         item.IsActive,
		CreatedBy:        uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create pay group", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapPayGroup(row), nil
}

func (s *Store) UpdatePayGroup(ctx context.Context, item *domain.PayGroup, actorID *uuid.UUID) (*domain.PayGroup, error) {
	row, err := s.getQueries(ctx).UpdatePayGroup(ctx, sqlc.UpdatePayGroupParams{
		TenantID:         item.TenantID,
		ID:               item.ID,
		Code:             item.Code,
		Name:             item.Name,
		Description:      textFromPtr(item.Description),
		GroupingType:     item.GroupingType,
		BranchID:         uuidFromPtr(item.BranchID),
		DepartmentID:     uuidFromPtr(item.DepartmentID),
		EmploymentTypeID: uuidFromPtr(item.EmploymentTypeID),
		ReportingTag:     textFromPtr(item.ReportingTag),
		Rules:            []byte(item.Rules),
		IsActive:         item.IsActive,
		UpdatedBy:        uuidFromPtr(actorID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPayGroupNotFound
		}
		return nil, s.logDBError(ctx, "update pay group", err, tenantIDField(item.TenantID), stringField("pay_group_id", item.ID.String()))
	}
	return mapPayGroup(row), nil
}

func (s *Store) GetPayGroup(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayGroup, error) {
	row, err := s.getQueries(ctx).GetPayGroup(ctx, sqlc.GetPayGroupParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPayGroupNotFound
		}
		return nil, s.logDBError(ctx, "get pay group", err, tenantIDField(tenantID), stringField("pay_group_id", id.String()))
	}
	return mapPayGroup(row), nil
}

func (s *Store) ListPayGroups(ctx context.Context, tenantID uuid.UUID) ([]*domain.PayGroup, error) {
	rows, err := s.getQueries(ctx).ListPayGroups(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list pay groups", err, tenantIDField(tenantID))
	}
	return mapPayGroups(rows), nil
}

func (s *Store) DeletePayGroup(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePayGroup(ctx, sqlc.SoftDeletePayGroupParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete pay group", err, tenantIDField(tenantID), stringField("pay_group_id", id.String()))
	}
	return nil
}

func (s *Store) UpsertPayGroupMember(ctx context.Context, item *domain.PayGroupMember, actorID *uuid.UUID) (*domain.PayGroupMember, error) {
	row, err := s.getQueries(ctx).UpsertPayGroupMember(ctx, sqlc.UpsertPayGroupMemberParams{
		TenantID:       item.TenantID,
		PayGroupID:     item.PayGroupID,
		UserID:         item.UserID,
		MembershipType: item.MembershipType,
		EffectiveFrom:  dateFromPtr(item.EffectiveFrom),
		EffectiveTo:    dateFromPtr(item.EffectiveTo),
		CreatedBy:      uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert pay group member", err, tenantIDField(item.TenantID), stringField("pay_group_id", item.PayGroupID.String()), stringField("user_id", item.UserID.String()))
	}
	return mapPayGroupMember(row), nil
}

func (s *Store) ListPayGroupMembers(ctx context.Context, tenantID uuid.UUID, payGroupID uuid.UUID) ([]*domain.PayGroupMember, error) {
	rows, err := s.getQueries(ctx).ListPayGroupMembers(ctx, sqlc.ListPayGroupMembersParams{TenantID: tenantID, PayGroupID: payGroupID})
	if err != nil {
		return nil, s.logDBError(ctx, "list pay group members", err, tenantIDField(tenantID), stringField("pay_group_id", payGroupID.String()))
	}
	return mapPayGroupMembers(rows), nil
}

func (s *Store) DeletePayGroupMember(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePayGroupMember(ctx, sqlc.SoftDeletePayGroupMemberParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete pay group member", err, tenantIDField(tenantID), stringField("member_id", id.String()))
	}
	return nil
}

func (s *Store) ListPayGroupEmployees(ctx context.Context, tenantID uuid.UUID, payGroupID uuid.UUID) ([]*domain.PayGroupEmployee, error) {
	rows, err := s.getQueries(ctx).ListPayGroupEmployees(ctx, sqlc.ListPayGroupEmployeesParams{TenantID: tenantID, ID: payGroupID})
	if err != nil {
		return nil, s.logDBError(ctx, "list pay group employees", err, tenantIDField(tenantID), stringField("pay_group_id", payGroupID.String()))
	}
	return mapPayGroupEmployees(rows), nil
}

func (s *Store) CreatePayRun(ctx context.Context, item *domain.PayRun, actorID *uuid.UUID) (*domain.PayRun, error) {
	row, err := s.getQueries(ctx).CreatePayRun(ctx, sqlc.CreatePayRunParams{
		TenantID:       item.TenantID,
		PayGroupID:     item.PayGroupID,
		FyID:           item.FYID,
		Month:          item.Month,
		Year:           item.Year,
		Status:         item.Status,
		EmployeeCount:  item.EmployeeCount,
		ReadyCount:     item.ReadyCount,
		BlockedCount:   item.BlockedCount,
		GeneratedCount: item.GeneratedCount,
		Readiness:      []byte(item.Readiness),
		Notes:          textFromPtr(item.Notes),
		CreatedBy:      uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create pay run", err, tenantIDField(item.TenantID), stringField("pay_group_id", item.PayGroupID.String()))
	}
	return mapPayRun(row), nil
}

func (s *Store) GetPayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PayRun, error) {
	row, err := s.getQueries(ctx).GetPayRun(ctx, sqlc.GetPayRunParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPayRunNotFound
		}
		return nil, s.logDBError(ctx, "get pay run", err, tenantIDField(tenantID), stringField("pay_run_id", id.String()))
	}
	return mapPayRun(row), nil
}

func (s *Store) ListPayRuns(ctx context.Context, tenantID uuid.UUID, payGroupID *uuid.UUID, month *int32, year *int32) ([]*domain.PayRun, error) {
	rows, err := s.getQueries(ctx).ListPayRuns(ctx, sqlc.ListPayRunsParams{
		TenantID:   tenantID,
		PayGroupID: uuidFromPtr(payGroupID),
		Month:      int4FromPtr(month),
		Year:       int4FromPtr(year),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list pay runs", err, tenantIDField(tenantID))
	}
	return mapPayRuns(rows), nil
}

func (s *Store) UpdatePayRunStatus(ctx context.Context, item *domain.PayRun, actorID *uuid.UUID) (*domain.PayRun, error) {
	row, err := s.getQueries(ctx).UpdatePayRunStatus(ctx, sqlc.UpdatePayRunStatusParams{
		TenantID:       item.TenantID,
		ID:             item.ID,
		Status:         item.Status,
		EmployeeCount:  item.EmployeeCount,
		ReadyCount:     item.ReadyCount,
		BlockedCount:   item.BlockedCount,
		GeneratedCount: item.GeneratedCount,
		Readiness:      []byte(item.Readiness),
		Notes:          textFromPtr(item.Notes),
		UpdatedBy:      uuidFromPtr(actorID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPayRunNotFound
		}
		return nil, s.logDBError(ctx, "update pay run status", err, tenantIDField(item.TenantID), stringField("pay_run_id", item.ID.String()))
	}
	return mapPayRun(row), nil
}

func (s *Store) DeletePayRun(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePayRun(ctx, sqlc.SoftDeletePayRunParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete pay run", err, tenantIDField(tenantID), stringField("pay_run_id", id.String()))
	}
	return nil
}

func (s *Store) UpsertPayRunEmployee(ctx context.Context, item *domain.PayRunEmployee, actorID *uuid.UUID) (*domain.PayRunEmployee, error) {
	row, err := s.getQueries(ctx).UpsertPayRunEmployee(ctx, sqlc.UpsertPayRunEmployeeParams{
		TenantID:        item.TenantID,
		PayRunID:        item.PayRunID,
		UserID:          item.UserID,
		ReadinessStatus: item.ReadinessStatus,
		BlockerReason:   textFromPtr(item.BlockerReason),
		SalarySlipID:    uuidFromPtr(item.SalarySlipID),
		GeneratedAt:     timestamptzFromPtr(item.GeneratedAt),
		CreatedBy:       uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "upsert pay run employee", err, tenantIDField(item.TenantID), stringField("pay_run_id", item.PayRunID.String()), stringField("user_id", item.UserID.String()))
	}
	return mapPayRunEmployeeRecord(row), nil
}

func (s *Store) ListPayRunEmployees(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) ([]*domain.PayRunEmployee, error) {
	rows, err := s.getQueries(ctx).ListPayRunEmployees(ctx, sqlc.ListPayRunEmployeesParams{TenantID: tenantID, PayRunID: payRunID})
	if err != nil {
		return nil, s.logDBError(ctx, "list pay run employees", err, tenantIDField(tenantID), stringField("pay_run_id", payRunID.String()))
	}
	return mapPayRunEmployees(rows), nil
}

func (s *Store) DeletePayRunLedger(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID, actorID *uuid.UUID) error {
	params := sqlc.DeletePayRunLedgerParams{TenantID: tenantID, PayRunID: payRunID, UpdatedBy: uuidFromPtr(actorID)}
	if err := s.getQueries(ctx).DeletePayRunLedger(ctx, params); err != nil {
		return s.logDBError(ctx, "delete pay run input ledger", err, tenantIDField(tenantID), stringField("pay_run_id", payRunID.String()))
	}
	componentParams := sqlc.DeletePayRunComponentLedgerParams{TenantID: tenantID, PayRunID: payRunID, UpdatedBy: uuidFromPtr(actorID)}
	if err := s.getQueries(ctx).DeletePayRunComponentLedger(ctx, componentParams); err != nil {
		return s.logDBError(ctx, "delete pay run component ledger", err, tenantIDField(tenantID), stringField("pay_run_id", payRunID.String()))
	}
	return nil
}

func (s *Store) CreatePayRunInput(ctx context.Context, item *domain.PayRunInput, actorID *uuid.UUID) (*domain.PayRunInput, error) {
	row, err := s.getQueries(ctx).CreatePayRunInput(ctx, sqlc.CreatePayRunInputParams{
		TenantID:    item.TenantID,
		PayRunID:    item.PayRunID,
		UserID:      item.UserID,
		InputType:   item.InputType,
		SourceType:  item.SourceType,
		SourceID:    uuidFromPtr(item.SourceID),
		Description: item.Description,
		Quantity:    numericFromFloatPtr(item.Quantity),
		Amount:      numericFromFloatPtr(item.Amount),
		Metadata:    []byte(item.Metadata),
		CreatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create pay run input", err, tenantIDField(item.TenantID), stringField("pay_run_id", item.PayRunID.String()), stringField("user_id", item.UserID.String()))
	}
	return mapPayRunInputRecord(row), nil
}

func (s *Store) CreatePayRunComponent(ctx context.Context, item *domain.PayRunComponent, actorID *uuid.UUID) (*domain.PayRunComponent, error) {
	row, err := s.getQueries(ctx).CreatePayRunComponent(ctx, sqlc.CreatePayRunComponentParams{
		TenantID:         item.TenantID,
		PayRunID:         item.PayRunID,
		UserID:           item.UserID,
		ComponentType:    item.ComponentType,
		Code:             item.Code,
		Name:             item.Name,
		Amount:           numericFromFloat(item.Amount),
		SourceInputID:    uuidFromPtr(item.SourceInputID),
		SalaryTemplateID: uuidFromPtr(item.SalaryTemplateID),
		Taxable:          item.Taxable,
		Statutory:        item.Statutory,
		EmployerCost:     item.EmployerCost,
		SortOrder:        item.SortOrder,
		Metadata:         []byte(item.Metadata),
		CreatedBy:        uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create pay run component", err, tenantIDField(item.TenantID), stringField("pay_run_id", item.PayRunID.String()), stringField("user_id", item.UserID.String()), stringField("code", item.Code))
	}
	return mapPayRunComponentRecord(row), nil
}

func (s *Store) ListPayRunInputs(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) ([]*domain.PayRunInput, error) {
	rows, err := s.getQueries(ctx).ListPayRunInputs(ctx, sqlc.ListPayRunInputsParams{TenantID: tenantID, PayRunID: payRunID})
	if err != nil {
		return nil, s.logDBError(ctx, "list pay run inputs", err, tenantIDField(tenantID), stringField("pay_run_id", payRunID.String()))
	}
	return mapPayRunInputs(rows), nil
}

func (s *Store) ListPayRunComponents(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) ([]*domain.PayRunComponent, error) {
	rows, err := s.getQueries(ctx).ListPayRunComponents(ctx, sqlc.ListPayRunComponentsParams{TenantID: tenantID, PayRunID: payRunID})
	if err != nil {
		return nil, s.logDBError(ctx, "list pay run components", err, tenantIDField(tenantID), stringField("pay_run_id", payRunID.String()))
	}
	return mapPayRunComponents(rows), nil
}

func (s *Store) GetPayRunLedgerSummary(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) (*domain.PayRunLedgerSummary, error) {
	row, err := s.getQueries(ctx).GetPayRunLedgerSummary(ctx, sqlc.GetPayRunLedgerSummaryParams{TenantID: tenantID, ID: payRunID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.PayRunLedgerSummary{PayRunID: payRunID}, nil
		}
		return nil, s.logDBError(ctx, "get pay run ledger summary", err, tenantIDField(tenantID), stringField("pay_run_id", payRunID.String()))
	}
	return mapPayRunLedgerSummary(row), nil
}

func (s *Store) CreatePayRunEvent(ctx context.Context, item *domain.PayRunEvent, actorID *uuid.UUID) (*domain.PayRunEvent, error) {
	row, err := s.getQueries(ctx).CreatePayRunEvent(ctx, sqlc.CreatePayRunEventParams{
		TenantID:   item.TenantID,
		PayRunID:   item.PayRunID,
		Action:     item.Action,
		FromStatus: textFromPtr(item.FromStatus),
		ToStatus:   textFromPtr(item.ToStatus),
		Remarks:    textFromPtr(item.Remarks),
		Metadata:   []byte(item.Metadata),
		CreatedBy:  uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create pay run event", err, tenantIDField(item.TenantID), stringField("pay_run_id", item.PayRunID.String()), stringField("action", item.Action))
	}
	return mapPayRunEvent(row), nil
}

func (s *Store) ListPayRunEvents(ctx context.Context, tenantID uuid.UUID, payRunID uuid.UUID) ([]*domain.PayRunEvent, error) {
	rows, err := s.getQueries(ctx).ListPayRunEvents(ctx, sqlc.ListPayRunEventsParams{TenantID: tenantID, PayRunID: payRunID})
	if err != nil {
		return nil, s.logDBError(ctx, "list pay run events", err, tenantIDField(tenantID), stringField("pay_run_id", payRunID.String()))
	}
	return mapPayRunEvents(rows), nil
}
