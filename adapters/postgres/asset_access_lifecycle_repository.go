package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateAssetItem(ctx context.Context, item *domain.AssetItem, actorID *uuid.UUID) (*domain.AssetItem, error) {
	row, err := s.getQueries(ctx).CreateAssetItem(ctx, sqlc.CreateAssetItemParams{TenantID: item.TenantID, AssetCode: item.AssetCode, AssetName: item.AssetName, AssetType: item.AssetType, Category: item.Category, SerialNumber: textFromPtr(item.SerialNumber), Vendor: textFromPtr(item.Vendor), PurchaseDate: dateFromPtr(item.PurchaseDate), WarrantyUntil: dateFromPtr(item.WarrantyUntil), OwnerUserID: uuidFromPtr(item.OwnerUserID), CustodianWorkerProfileID: uuidFromPtr(item.CustodianWorkerProfileID), Status: item.Status, LocationLabel: textFromPtr(item.LocationLabel), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create asset item", err, tenantIDField(item.TenantID), stringField("asset_code", item.AssetCode))
	}
	return mapAssetItem(row), nil
}

func (s *Store) UpdateAssetItem(ctx context.Context, item *domain.AssetItem, actorID *uuid.UUID) (*domain.AssetItem, error) {
	row, err := s.getQueries(ctx).UpdateAssetItem(ctx, sqlc.UpdateAssetItemParams{TenantID: item.TenantID, ID: item.ID, AssetCode: item.AssetCode, AssetName: item.AssetName, AssetType: item.AssetType, Category: item.Category, SerialNumber: textFromPtr(item.SerialNumber), Vendor: textFromPtr(item.Vendor), PurchaseDate: dateFromPtr(item.PurchaseDate), WarrantyUntil: dateFromPtr(item.WarrantyUntil), OwnerUserID: uuidFromPtr(item.OwnerUserID), CustodianWorkerProfileID: uuidFromPtr(item.CustodianWorkerProfileID), Status: item.Status, LocationLabel: textFromPtr(item.LocationLabel), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAssetItemNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update asset item", err, tenantIDField(item.TenantID), stringField("asset_id", item.ID.String()))
	}
	return mapAssetItem(row), nil
}

func (s *Store) UpdateAssetItemStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.AssetItem, error) {
	row, err := s.getQueries(ctx).UpdateAssetItemStatus(ctx, sqlc.UpdateAssetItemStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAssetItemNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update asset item status", err, tenantIDField(tenantID), stringField("asset_id", id.String()))
	}
	return mapAssetItem(row), nil
}

func (s *Store) GetAssetItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AssetItem, error) {
	row, err := s.getQueries(ctx).GetAssetItem(ctx, sqlc.GetAssetItemParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAssetItemNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get asset item", err, tenantIDField(tenantID), stringField("asset_id", id.String()))
	}
	return mapAssetItem(row), nil
}

func (s *Store) ListAssetItems(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AssetItem, error) {
	rows, err := s.getQueries(ctx).ListAssetItems(ctx, sqlc.ListAssetItemsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, Status: textFromPtr(filter.Status), Category: textFromPtr(filter.Category), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list asset items", err, tenantIDField(filter.TenantID))
	}
	return mapAssetItemList(rows), nil
}

func (s *Store) DeleteAssetItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAssetItem(ctx, sqlc.SoftDeleteAssetItemParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete asset item", err, tenantIDField(tenantID), stringField("asset_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAccessCatalogItem(ctx context.Context, item *domain.AccessCatalogItem, actorID *uuid.UUID) (*domain.AccessCatalogItem, error) {
	row, err := s.getQueries(ctx).CreateAccessCatalogItem(ctx, sqlc.CreateAccessCatalogItemParams{TenantID: item.TenantID, AccessCode: item.AccessCode, AccessName: item.AccessName, AccessType: item.AccessType, SystemName: textFromPtr(item.SystemName), OwnerUserID: uuidFromPtr(item.OwnerUserID), ProvisioningMethod: item.ProvisioningMethod, RequiresApproval: item.RequiresApproval, DefaultForOnboarding: item.DefaultForOnboarding, DefaultForExitRevocation: item.DefaultForExitRevocation, Status: item.Status, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create access catalog item", err, tenantIDField(item.TenantID), stringField("access_code", item.AccessCode))
	}
	return mapAccessCatalogItem(row), nil
}

func (s *Store) UpdateAccessCatalogItem(ctx context.Context, item *domain.AccessCatalogItem, actorID *uuid.UUID) (*domain.AccessCatalogItem, error) {
	row, err := s.getQueries(ctx).UpdateAccessCatalogItem(ctx, sqlc.UpdateAccessCatalogItemParams{TenantID: item.TenantID, ID: item.ID, AccessCode: item.AccessCode, AccessName: item.AccessName, AccessType: item.AccessType, SystemName: textFromPtr(item.SystemName), OwnerUserID: uuidFromPtr(item.OwnerUserID), ProvisioningMethod: item.ProvisioningMethod, RequiresApproval: item.RequiresApproval, DefaultForOnboarding: item.DefaultForOnboarding, DefaultForExitRevocation: item.DefaultForExitRevocation, Status: item.Status, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAccessItemNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update access catalog item", err, tenantIDField(item.TenantID), stringField("access_item_id", item.ID.String()))
	}
	return mapAccessCatalogItem(row), nil
}

func (s *Store) GetAccessCatalogItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AccessCatalogItem, error) {
	row, err := s.getQueries(ctx).GetAccessCatalogItem(ctx, sqlc.GetAccessCatalogItemParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAccessItemNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get access catalog item", err, tenantIDField(tenantID), stringField("access_item_id", id.String()))
	}
	return mapAccessCatalogItem(row), nil
}

func (s *Store) ListAccessCatalogItems(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AccessCatalogItem, error) {
	rows, err := s.getQueries(ctx).ListAccessCatalogItems(ctx, sqlc.ListAccessCatalogItemsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, Status: textFromPtr(filter.Status), AccessType: textFromPtr(filter.AccessType), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list access catalog items", err, tenantIDField(filter.TenantID))
	}
	return mapAccessCatalogItems(rows), nil
}

func (s *Store) DeleteAccessCatalogItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteAccessCatalogItem(ctx, sqlc.SoftDeleteAccessCatalogItemParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete access catalog item", err, tenantIDField(tenantID), stringField("access_item_id", id.String()))
	}
	return nil
}

func (s *Store) CreateAssetAssignment(ctx context.Context, item *domain.AssetAssignment, actorID *uuid.UUID) (*domain.AssetAssignment, error) {
	row, err := s.getQueries(ctx).CreateAssetAssignment(ctx, sqlc.CreateAssetAssignmentParams{TenantID: item.TenantID, AssetID: item.AssetID, WorkerProfileID: item.WorkerProfileID, EmployeeID: uuidFromPtr(item.EmployeeID), CandidateOnboardingID: uuidFromPtr(item.CandidateOnboardingID), ExitRequestID: uuidFromPtr(item.ExitRequestID), RequestedBy: uuidFromPtr(item.RequestedBy), ApprovedBy: uuidFromPtr(item.ApprovedBy), IssuedBy: uuidFromPtr(item.IssuedBy), ReturnedBy: uuidFromPtr(item.ReturnedBy), ApprovedAt: timestamptzFromPtr(item.ApprovedAt), IssuedOn: dateFromPtr(item.IssuedOn), ExpectedReturnOn: dateFromPtr(item.ExpectedReturnOn), ReturnedOn: dateFromPtr(item.ReturnedOn), IssueCondition: item.IssueCondition, ReturnCondition: textFromPtr(item.ReturnCondition), DamageStatus: item.DamageStatus, RecoveryAmount: compNumeric(item.RecoveryAmount), Status: item.Status, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create asset assignment", err, tenantIDField(item.TenantID), stringField("asset_id", item.AssetID.String()))
	}
	return mapAssetAssignment(row), nil
}

func (s *Store) UpdateAssetAssignment(ctx context.Context, item *domain.AssetAssignment, actorID *uuid.UUID) (*domain.AssetAssignment, error) {
	row, err := s.getQueries(ctx).UpdateAssetAssignment(ctx, sqlc.UpdateAssetAssignmentParams{TenantID: item.TenantID, ID: item.ID, WorkerProfileID: item.WorkerProfileID, EmployeeID: uuidFromPtr(item.EmployeeID), CandidateOnboardingID: uuidFromPtr(item.CandidateOnboardingID), ExitRequestID: uuidFromPtr(item.ExitRequestID), RequestedBy: uuidFromPtr(item.RequestedBy), ApprovedBy: uuidFromPtr(item.ApprovedBy), IssuedBy: uuidFromPtr(item.IssuedBy), ReturnedBy: uuidFromPtr(item.ReturnedBy), ApprovedAt: timestamptzFromPtr(item.ApprovedAt), IssuedOn: dateFromPtr(item.IssuedOn), ExpectedReturnOn: dateFromPtr(item.ExpectedReturnOn), ReturnedOn: dateFromPtr(item.ReturnedOn), IssueCondition: item.IssueCondition, ReturnCondition: textFromPtr(item.ReturnCondition), DamageStatus: item.DamageStatus, RecoveryAmount: compNumeric(item.RecoveryAmount), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAssetAssignmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update asset assignment", err, tenantIDField(item.TenantID), stringField("assignment_id", item.ID.String()))
	}
	return mapAssetAssignment(row), nil
}

func (s *Store) UpdateAssetAssignmentStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.AssetAssignment, error) {
	row, err := s.getQueries(ctx).UpdateAssetAssignmentStatus(ctx, sqlc.UpdateAssetAssignmentStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAssetAssignmentNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update asset assignment status", err, tenantIDField(tenantID), stringField("assignment_id", id.String()))
	}
	return mapAssetAssignment(row), nil
}

func (s *Store) ListAssetAssignments(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AssetAssignment, error) {
	rows, err := s.getQueries(ctx).ListAssetAssignments(ctx, sqlc.ListAssetAssignmentsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, WorkerProfileID: uuidFromPtr(filter.WorkerProfileID), AssetID: uuidFromPtr(filter.AssetID), ExitRequestID: uuidFromPtr(filter.ExitRequestID), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list asset assignments", err, tenantIDField(filter.TenantID))
	}
	return mapAssetAssignmentList(rows), nil
}

func (s *Store) CreateAccessLifecycleTask(ctx context.Context, item *domain.AccessLifecycleTask, actorID *uuid.UUID) (*domain.AccessLifecycleTask, error) {
	row, err := s.getQueries(ctx).CreateAccessLifecycleTask(ctx, sqlc.CreateAccessLifecycleTaskParams{TenantID: item.TenantID, AccessItemID: item.AccessItemID, WorkerProfileID: item.WorkerProfileID, EmployeeID: uuidFromPtr(item.EmployeeID), CandidateOnboardingID: uuidFromPtr(item.CandidateOnboardingID), ExitRequestID: uuidFromPtr(item.ExitRequestID), TaskType: item.TaskType, RequestedBy: uuidFromPtr(item.RequestedBy), ApprovedBy: uuidFromPtr(item.ApprovedBy), OwnerUserID: uuidFromPtr(item.OwnerUserID), ApprovedAt: timestamptzFromPtr(item.ApprovedAt), DueDate: dateFromPtr(item.DueDate), CompletedAt: timestamptzFromPtr(item.CompletedAt), ExternalReference: textFromPtr(item.ExternalReference), Status: item.Status, Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create access lifecycle task", err, tenantIDField(item.TenantID), stringField("access_item_id", item.AccessItemID.String()))
	}
	return mapAccessLifecycleTask(row), nil
}

func (s *Store) UpdateAccessLifecycleTask(ctx context.Context, item *domain.AccessLifecycleTask, actorID *uuid.UUID) (*domain.AccessLifecycleTask, error) {
	row, err := s.getQueries(ctx).UpdateAccessLifecycleTask(ctx, sqlc.UpdateAccessLifecycleTaskParams{TenantID: item.TenantID, ID: item.ID, WorkerProfileID: item.WorkerProfileID, EmployeeID: uuidFromPtr(item.EmployeeID), CandidateOnboardingID: uuidFromPtr(item.CandidateOnboardingID), ExitRequestID: uuidFromPtr(item.ExitRequestID), TaskType: item.TaskType, RequestedBy: uuidFromPtr(item.RequestedBy), ApprovedBy: uuidFromPtr(item.ApprovedBy), OwnerUserID: uuidFromPtr(item.OwnerUserID), ApprovedAt: timestamptzFromPtr(item.ApprovedAt), DueDate: dateFromPtr(item.DueDate), CompletedAt: timestamptzFromPtr(item.CompletedAt), ExternalReference: textFromPtr(item.ExternalReference), Notes: textFromPtr(item.Notes), Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAccessTaskNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update access lifecycle task", err, tenantIDField(item.TenantID), stringField("task_id", item.ID.String()))
	}
	return mapAccessLifecycleTask(row), nil
}

func (s *Store) UpdateAccessLifecycleTaskStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.AccessLifecycleTask, error) {
	row, err := s.getQueries(ctx).UpdateAccessLifecycleTaskStatus(ctx, sqlc.UpdateAccessLifecycleTaskStatusParams{TenantID: tenantID, ID: id, Status: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAccessTaskNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update access lifecycle task status", err, tenantIDField(tenantID), stringField("task_id", id.String()))
	}
	return mapAccessLifecycleTask(row), nil
}

func (s *Store) ListAccessLifecycleTasks(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AccessLifecycleTask, error) {
	rows, err := s.getQueries(ctx).ListAccessLifecycleTasks(ctx, sqlc.ListAccessLifecycleTasksParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, WorkerProfileID: uuidFromPtr(filter.WorkerProfileID), AccessItemID: uuidFromPtr(filter.AccessItemID), ExitRequestID: uuidFromPtr(filter.ExitRequestID), Status: textFromPtr(filter.Status), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list access lifecycle tasks", err, tenantIDField(filter.TenantID))
	}
	return mapAccessLifecycleTaskList(rows), nil
}

func (s *Store) CreateAssetAccessEvent(ctx context.Context, item *domain.AssetAccessEvent, actorID *uuid.UUID) (*domain.AssetAccessEvent, error) {
	row, err := s.getQueries(ctx).CreateAssetAccessEvent(ctx, sqlc.CreateAssetAccessEventParams{TenantID: item.TenantID, SourceType: item.SourceType, SourceID: uuidFromPtr(item.SourceID), Action: item.Action, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), Remarks: textFromPtr(item.Remarks), Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create asset access event", err, tenantIDField(item.TenantID), stringField("source_type", item.SourceType))
	}
	return mapAssetAccessEvent(row), nil
}

func (s *Store) ListAssetAccessEvents(ctx context.Context, filter domain.AssetAccessFilter, sourceType *string, sourceID *uuid.UUID) ([]*domain.AssetAccessEvent, error) {
	rows, err := s.getQueries(ctx).ListAssetAccessEvents(ctx, sqlc.ListAssetAccessEventsParams{TenantID: filter.TenantID, Limit: filter.Limit, Offset: filter.Offset, SourceType: textFromPtr(sourceType), SourceID: uuidFromPtr(sourceID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list asset access events", err, tenantIDField(filter.TenantID))
	}
	return mapAssetAccessEvents(rows), nil
}

func (s *Store) GetAssetAccessSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.AssetAccessSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetAssetAccessSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get asset access summary", err, tenantIDField(tenantID))
	}
	return mapAssetAccessSummary(rows), nil
}
