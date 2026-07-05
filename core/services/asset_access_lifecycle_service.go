package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateAssetItem(ctx context.Context, cmd ports.AssetItemCommand) (*domain.AssetItem, error) {
	if err := s.validateAssetWorkerRef(ctx, cmd.TenantID, cmd.CustodianWorkerProfileID); err != nil {
		return nil, err
	}
	item := assetItemFromCommand(cmd)
	if err := domain.ValidateAssetItem(item); err != nil {
		s.logError("validate asset item", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.assetAccess.CreateAssetItem(ctx, item, cmd.ActorID)
	if err == nil {
		_, _ = s.assetAccessEvent(ctx, result.TenantID, "asset_item", &result.ID, "created", nil, &result.Status, nil, cmd.ActorID)
	}
	return result, err
}

func (s *TenantService) UpdateAssetItem(ctx context.Context, cmd ports.AssetItemCommand) (*domain.AssetItem, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAssetItem
	}
	if err := s.validateAssetWorkerRef(ctx, cmd.TenantID, cmd.CustodianWorkerProfileID); err != nil {
		return nil, err
	}
	item := assetItemFromCommand(cmd)
	item.ID = cmd.ID
	if err := domain.ValidateAssetItem(item); err != nil {
		return nil, err
	}
	return s.assetAccess.UpdateAssetItem(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateAssetItemStatus(ctx context.Context, cmd ports.AssetAccessStatusCommand) (*domain.AssetItem, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAssetItem
	}
	status := domain.NormalizeAssetStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidAssetItem
	}
	before, _ := s.assetAccess.GetAssetItem(ctx, cmd.TenantID, cmd.ID)
	result, err := s.assetAccess.UpdateAssetItemStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	_, _ = s.assetAccessEvent(ctx, cmd.TenantID, "asset_item", &cmd.ID, "status_changed", from, &result.Status, cmd.Remarks, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListAssetItems(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AssetItem, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeAssetAccessPage(&filter.Limit, &filter.Offset)
	return s.assetAccess.ListAssetItems(ctx, filter)
}

func (s *TenantService) DeleteAssetItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidAssetItem
	}
	if err := s.assetAccess.DeleteAssetItem(ctx, tenantID, id, actorID); err != nil {
		return err
	}
	_, _ = s.assetAccessEvent(ctx, tenantID, "asset_item", &id, "deleted", nil, nil, nil, actorID)
	return nil
}

func (s *TenantService) CreateAccessCatalogItem(ctx context.Context, cmd ports.AccessCatalogItemCommand) (*domain.AccessCatalogItem, error) {
	item := accessCatalogItemFromCommand(cmd)
	if err := domain.ValidateAccessCatalogItem(item); err != nil {
		s.logError("validate access catalog item", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.assetAccess.CreateAccessCatalogItem(ctx, item, cmd.ActorID)
	if err == nil {
		_, _ = s.assetAccessEvent(ctx, result.TenantID, "access_catalog_item", &result.ID, "created", nil, &result.Status, nil, cmd.ActorID)
	}
	return result, err
}

func (s *TenantService) UpdateAccessCatalogItem(ctx context.Context, cmd ports.AccessCatalogItemCommand) (*domain.AccessCatalogItem, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAccessItem
	}
	item := accessCatalogItemFromCommand(cmd)
	item.ID = cmd.ID
	if err := domain.ValidateAccessCatalogItem(item); err != nil {
		return nil, err
	}
	return s.assetAccess.UpdateAccessCatalogItem(ctx, item, cmd.ActorID)
}

func (s *TenantService) ListAccessCatalogItems(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AccessCatalogItem, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeAssetAccessPage(&filter.Limit, &filter.Offset)
	return s.assetAccess.ListAccessCatalogItems(ctx, filter)
}

func (s *TenantService) DeleteAccessCatalogItem(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidAccessItem
	}
	if err := s.assetAccess.DeleteAccessCatalogItem(ctx, tenantID, id, actorID); err != nil {
		return err
	}
	_, _ = s.assetAccessEvent(ctx, tenantID, "access_catalog_item", &id, "deleted", nil, nil, nil, actorID)
	return nil
}

func (s *TenantService) CreateAssetAssignment(ctx context.Context, cmd ports.AssetAssignmentCommand) (*domain.AssetAssignment, error) {
	if _, err := s.assetAccess.GetAssetItem(ctx, cmd.TenantID, cmd.AssetID); err != nil {
		return nil, err
	}
	if err := s.validateAssetWorkerRef(ctx, cmd.TenantID, &cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	item := assetAssignmentFromCommand(cmd)
	if item.RequestedBy == nil {
		item.RequestedBy = cmd.ActorID
	}
	if err := domain.ValidateAssetAssignment(item); err != nil {
		s.logError("validate asset assignment", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.assetAccess.CreateAssetAssignment(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.assetAccessEvent(ctx, result.TenantID, "asset_assignment", &result.ID, "created", nil, &result.Status, nil, cmd.ActorID)
	_, _ = s.syncAssetStatusForAssignment(ctx, result, cmd.ActorID)
	return result, nil
}

func (s *TenantService) UpdateAssetAssignment(ctx context.Context, cmd ports.AssetAssignmentCommand) (*domain.AssetAssignment, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAssetAssignment
	}
	if err := s.validateAssetWorkerRef(ctx, cmd.TenantID, &cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	item := assetAssignmentFromCommand(cmd)
	item.ID = cmd.ID
	if err := domain.ValidateAssetAssignment(item); err != nil {
		return nil, err
	}
	result, err := s.assetAccess.UpdateAssetAssignment(ctx, item, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.syncAssetStatusForAssignment(ctx, result, cmd.ActorID)
	return result, nil
}

func (s *TenantService) UpdateAssetAssignmentStatus(ctx context.Context, cmd ports.AssetAccessStatusCommand) (*domain.AssetAssignment, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAssetAssignment
	}
	status := domain.NormalizeAssetAssignmentStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidAssetAssignment
	}
	result, err := s.assetAccess.UpdateAssetAssignmentStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.assetAccessEvent(ctx, cmd.TenantID, "asset_assignment", &cmd.ID, "status_changed", nil, &result.Status, cmd.Remarks, cmd.ActorID)
	_, _ = s.syncAssetStatusForAssignment(ctx, result, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListAssetAssignments(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AssetAssignment, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeAssetAccessPage(&filter.Limit, &filter.Offset)
	return s.assetAccess.ListAssetAssignments(ctx, filter)
}

func (s *TenantService) CreateAccessLifecycleTask(ctx context.Context, cmd ports.AccessLifecycleTaskCommand) (*domain.AccessLifecycleTask, error) {
	if _, err := s.assetAccess.GetAccessCatalogItem(ctx, cmd.TenantID, cmd.AccessItemID); err != nil {
		return nil, err
	}
	if err := s.validateAssetWorkerRef(ctx, cmd.TenantID, &cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	item := accessLifecycleTaskFromCommand(cmd)
	if item.RequestedBy == nil {
		item.RequestedBy = cmd.ActorID
	}
	if err := domain.ValidateAccessLifecycleTask(item); err != nil {
		s.logError("validate access lifecycle task", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.assetAccess.CreateAccessLifecycleTask(ctx, item, cmd.ActorID)
	if err == nil {
		_, _ = s.assetAccessEvent(ctx, result.TenantID, "access_lifecycle_task", &result.ID, "created", nil, &result.Status, nil, cmd.ActorID)
	}
	return result, err
}

func (s *TenantService) UpdateAccessLifecycleTask(ctx context.Context, cmd ports.AccessLifecycleTaskCommand) (*domain.AccessLifecycleTask, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAccessTask
	}
	if err := s.validateAssetWorkerRef(ctx, cmd.TenantID, &cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	item := accessLifecycleTaskFromCommand(cmd)
	item.ID = cmd.ID
	if err := domain.ValidateAccessLifecycleTask(item); err != nil {
		return nil, err
	}
	return s.assetAccess.UpdateAccessLifecycleTask(ctx, item, cmd.ActorID)
}

func (s *TenantService) UpdateAccessLifecycleTaskStatus(ctx context.Context, cmd ports.AssetAccessStatusCommand) (*domain.AccessLifecycleTask, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidAccessTask
	}
	status := domain.NormalizeAccessTaskStatus(cmd.Status)
	if status == "" {
		return nil, domain.ErrInvalidAccessTask
	}
	result, err := s.assetAccess.UpdateAccessLifecycleTaskStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	_, _ = s.assetAccessEvent(ctx, cmd.TenantID, "access_lifecycle_task", &cmd.ID, "status_changed", nil, &result.Status, cmd.Remarks, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListAccessLifecycleTasks(ctx context.Context, filter domain.AssetAccessFilter) ([]*domain.AccessLifecycleTask, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeAssetAccessPage(&filter.Limit, &filter.Offset)
	return s.assetAccess.ListAccessLifecycleTasks(ctx, filter)
}

func (s *TenantService) ListAssetAccessEvents(ctx context.Context, filter domain.AssetAccessFilter, sourceType *string, sourceID *uuid.UUID) ([]*domain.AssetAccessEvent, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	normalizeAssetAccessPage(&filter.Limit, &filter.Offset)
	return s.assetAccess.ListAssetAccessEvents(ctx, filter, sourceType, sourceID)
}

func (s *TenantService) GetAssetAccessSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.AssetAccessSummaryRow, error) {
	if tenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.assetAccess.GetAssetAccessSummary(ctx, tenantID)
}

func (s *TenantService) assetAccessEvent(ctx context.Context, tenantID uuid.UUID, sourceType string, sourceID *uuid.UUID, action string, fromStatus *string, toStatus *string, remarks *string, actorID *uuid.UUID) (*domain.AssetAccessEvent, error) {
	return s.assetAccess.CreateAssetAccessEvent(ctx, &domain.AssetAccessEvent{TenantID: tenantID, SourceType: sourceType, SourceID: sourceID, Action: action, FromStatus: fromStatus, ToStatus: toStatus, Remarks: remarks}, actorID)
}

func (s *TenantService) syncAssetStatusForAssignment(ctx context.Context, item *domain.AssetAssignment, actorID *uuid.UUID) (*domain.AssetItem, error) {
	status := domain.AssetStatusAvailable
	switch item.Status {
	case domain.AssetAssignmentApproved:
		status = domain.AssetStatusReserved
	case domain.AssetAssignmentIssued:
		status = domain.AssetStatusIssued
	case domain.AssetAssignmentReturnDue:
		status = domain.AssetStatusReturnDue
	case domain.AssetAssignmentDamaged:
		status = domain.AssetStatusDamaged
	case domain.AssetAssignmentLost:
		status = domain.AssetStatusLost
	case domain.AssetAssignmentReturned, domain.AssetAssignmentCancelled:
		status = domain.AssetStatusAvailable
	default:
		return nil, nil
	}
	return s.assetAccess.UpdateAssetItemStatus(ctx, item.TenantID, item.AssetID, status, actorID)
}

func (s *TenantService) validateAssetWorkerRef(ctx context.Context, tenantID uuid.UUID, workerProfileID *uuid.UUID) error {
	if workerProfileID == nil || *workerProfileID == uuid.Nil {
		return nil
	}
	_, err := s.workerProfiles.GetWorkerProfile(ctx, tenantID, *workerProfileID)
	return err
}

func assetItemFromCommand(cmd ports.AssetItemCommand) *domain.AssetItem {
	return &domain.AssetItem{ID: cmd.ID, TenantID: cmd.TenantID, AssetCode: cmd.AssetCode, AssetName: cmd.AssetName, AssetType: cmd.AssetType, Category: cmd.Category, SerialNumber: cmd.SerialNumber, Vendor: cmd.Vendor, PurchaseDate: cmd.PurchaseDate, WarrantyUntil: cmd.WarrantyUntil, OwnerUserID: cmd.OwnerUserID, CustodianWorkerProfileID: cmd.CustodianWorkerProfileID, Status: cmd.Status, LocationLabel: cmd.LocationLabel, Notes: cmd.Notes, Metadata: cmd.Metadata}
}

func accessCatalogItemFromCommand(cmd ports.AccessCatalogItemCommand) *domain.AccessCatalogItem {
	return &domain.AccessCatalogItem{ID: cmd.ID, TenantID: cmd.TenantID, AccessCode: cmd.AccessCode, AccessName: cmd.AccessName, AccessType: cmd.AccessType, SystemName: cmd.SystemName, OwnerUserID: cmd.OwnerUserID, ProvisioningMethod: cmd.ProvisioningMethod, RequiresApproval: cmd.RequiresApproval, DefaultForOnboarding: cmd.DefaultForOnboarding, DefaultForExitRevocation: cmd.DefaultForExitRevocation, Status: cmd.Status, Notes: cmd.Notes, Metadata: cmd.Metadata}
}

func assetAssignmentFromCommand(cmd ports.AssetAssignmentCommand) *domain.AssetAssignment {
	return &domain.AssetAssignment{ID: cmd.ID, TenantID: cmd.TenantID, AssetID: cmd.AssetID, WorkerProfileID: cmd.WorkerProfileID, EmployeeID: cmd.EmployeeID, CandidateOnboardingID: cmd.CandidateOnboardingID, ExitRequestID: cmd.ExitRequestID, RequestedBy: cmd.RequestedBy, ApprovedBy: cmd.ApprovedBy, IssuedBy: cmd.IssuedBy, ReturnedBy: cmd.ReturnedBy, ApprovedAt: cmd.ApprovedAt, IssuedOn: cmd.IssuedOn, ExpectedReturnOn: cmd.ExpectedReturnOn, ReturnedOn: cmd.ReturnedOn, IssueCondition: cmd.IssueCondition, ReturnCondition: cmd.ReturnCondition, DamageStatus: cmd.DamageStatus, RecoveryAmount: cmd.RecoveryAmount, Status: cmd.Status, Notes: cmd.Notes, Metadata: cmd.Metadata}
}

func accessLifecycleTaskFromCommand(cmd ports.AccessLifecycleTaskCommand) *domain.AccessLifecycleTask {
	return &domain.AccessLifecycleTask{ID: cmd.ID, TenantID: cmd.TenantID, AccessItemID: cmd.AccessItemID, WorkerProfileID: cmd.WorkerProfileID, EmployeeID: cmd.EmployeeID, CandidateOnboardingID: cmd.CandidateOnboardingID, ExitRequestID: cmd.ExitRequestID, TaskType: cmd.TaskType, RequestedBy: cmd.RequestedBy, ApprovedBy: cmd.ApprovedBy, OwnerUserID: cmd.OwnerUserID, ApprovedAt: cmd.ApprovedAt, DueDate: cmd.DueDate, CompletedAt: cmd.CompletedAt, ExternalReference: cmd.ExternalReference, Status: cmd.Status, Notes: cmd.Notes, Metadata: cmd.Metadata}
}

func normalizeAssetAccessPage(limit *int32, offset *int32) {
	if *limit <= 0 || *limit > 200 {
		*limit = 50
	}
	if *offset < 0 {
		*offset = 0
	}
}
