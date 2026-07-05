package postgres

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func mapAssetItem(row sqlc.HrmsAssetItem) *domain.AssetItem {
	return assetItemFromParts(row.ID, row.TenantID, row.AssetCode, row.AssetName, row.AssetType, row.Category, row.SerialNumber, row.Vendor, row.PurchaseDate, row.WarrantyUntil, row.OwnerUserID, row.CustodianWorkerProfileID, row.Status, row.CurrentAssignmentID, row.LocationLabel, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil)
}

func mapAssetItemList(rows []sqlc.ListAssetItemsRow) []*domain.AssetItem {
	items := make([]*domain.AssetItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, assetItemFromParts(row.ID, row.TenantID, row.AssetCode, row.AssetName, row.AssetType, row.Category, row.SerialNumber, row.Vendor, row.PurchaseDate, row.WarrantyUntil, row.OwnerUserID, row.CustodianWorkerProfileID, row.Status, row.CurrentAssignmentID, row.LocationLabel, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, ptrFromText(row.CustodianName), ptrFromText(row.CustodianCode)))
	}
	return items
}

func assetItemFromParts(id uuid.UUID, tenantID uuid.UUID, assetCode string, assetName string, assetType string, category string, serialNumber pgtype.Text, vendor pgtype.Text, purchaseDate pgtype.Date, warrantyUntil pgtype.Date, ownerUserID pgtype.UUID, custodianWorkerProfileID pgtype.UUID, status string, currentAssignmentID pgtype.UUID, locationLabel pgtype.Text, notes pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, custodianName *string, custodianCode *string) *domain.AssetItem {
	return &domain.AssetItem{ID: id, TenantID: tenantID, AssetCode: assetCode, AssetName: assetName, AssetType: assetType, Category: category, SerialNumber: ptrFromText(serialNumber), Vendor: ptrFromText(vendor), PurchaseDate: ptrFromDate(purchaseDate), WarrantyUntil: ptrFromDate(warrantyUntil), OwnerUserID: ptrFromUUID(ownerUserID), CustodianWorkerProfileID: ptrFromUUID(custodianWorkerProfileID), Status: status, CurrentAssignmentID: ptrFromUUID(currentAssignmentID), LocationLabel: ptrFromText(locationLabel), Notes: ptrFromText(notes), Metadata: assetAccessRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), CustodianName: custodianName, CustodianCode: custodianCode}
}

func mapAccessCatalogItem(row sqlc.HrmsAccessCatalogItem) *domain.AccessCatalogItem {
	return &domain.AccessCatalogItem{ID: row.ID, TenantID: row.TenantID, AccessCode: row.AccessCode, AccessName: row.AccessName, AccessType: row.AccessType, SystemName: ptrFromText(row.SystemName), OwnerUserID: ptrFromUUID(row.OwnerUserID), ProvisioningMethod: row.ProvisioningMethod, RequiresApproval: row.RequiresApproval, DefaultForOnboarding: row.DefaultForOnboarding, DefaultForExitRevocation: row.DefaultForExitRevocation, Status: row.Status, Notes: ptrFromText(row.Notes), Metadata: assetAccessRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAccessCatalogItems(rows []sqlc.HrmsAccessCatalogItem) []*domain.AccessCatalogItem {
	items := make([]*domain.AccessCatalogItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAccessCatalogItem(row))
	}
	return items
}

func mapAssetAssignment(row sqlc.HrmsAssetAssignment) *domain.AssetAssignment {
	return assetAssignmentFromParts(row.ID, row.TenantID, row.AssetID, row.WorkerProfileID, row.EmployeeID, row.CandidateOnboardingID, row.ExitRequestID, row.RequestedBy, row.ApprovedBy, row.IssuedBy, row.ReturnedBy, row.ApprovedAt, row.IssuedOn, row.ExpectedReturnOn, row.ReturnedOn, row.IssueCondition, row.ReturnCondition, row.DamageStatus, row.RecoveryAmount, row.Status, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil, nil)
}

func mapAssetAssignmentList(rows []sqlc.ListAssetAssignmentsRow) []*domain.AssetAssignment {
	items := make([]*domain.AssetAssignment, 0, len(rows))
	for _, row := range rows {
		items = append(items, assetAssignmentFromParts(row.ID, row.TenantID, row.AssetID, row.WorkerProfileID, row.EmployeeID, row.CandidateOnboardingID, row.ExitRequestID, row.RequestedBy, row.ApprovedBy, row.IssuedBy, row.ReturnedBy, row.ApprovedAt, row.IssuedOn, row.ExpectedReturnOn, row.ReturnedOn, row.IssueCondition, row.ReturnCondition, row.DamageStatus, row.RecoveryAmount, row.Status, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.AssetCode, &row.AssetName, &row.AssetType, &row.Category, &row.WorkerDisplayName, ptrFromText(row.WorkerCode)))
	}
	return items
}

func assetAssignmentFromParts(id uuid.UUID, tenantID uuid.UUID, assetID uuid.UUID, workerProfileID uuid.UUID, employeeID pgtype.UUID, candidateOnboardingID pgtype.UUID, exitRequestID pgtype.UUID, requestedBy pgtype.UUID, approvedBy pgtype.UUID, issuedBy pgtype.UUID, returnedBy pgtype.UUID, approvedAt pgtype.Timestamptz, issuedOn pgtype.Date, expectedReturnOn pgtype.Date, returnedOn pgtype.Date, issueCondition string, returnCondition pgtype.Text, damageStatus string, recoveryAmount pgtype.Numeric, status string, notes pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, assetCode *string, assetName *string, assetType *string, category *string, workerDisplayName *string, workerCode *string) *domain.AssetAssignment {
	return &domain.AssetAssignment{ID: id, TenantID: tenantID, AssetID: assetID, WorkerProfileID: workerProfileID, EmployeeID: ptrFromUUID(employeeID), CandidateOnboardingID: ptrFromUUID(candidateOnboardingID), ExitRequestID: ptrFromUUID(exitRequestID), RequestedBy: ptrFromUUID(requestedBy), ApprovedBy: ptrFromUUID(approvedBy), IssuedBy: ptrFromUUID(issuedBy), ReturnedBy: ptrFromUUID(returnedBy), ApprovedAt: ptrFromTimestamptz(approvedAt), IssuedOn: ptrFromDate(issuedOn), ExpectedReturnOn: ptrFromDate(expectedReturnOn), ReturnedOn: ptrFromDate(returnedOn), IssueCondition: issueCondition, ReturnCondition: ptrFromText(returnCondition), DamageStatus: damageStatus, RecoveryAmount: compFloat(recoveryAmount), Status: status, Notes: ptrFromText(notes), Metadata: assetAccessRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), AssetCode: assetCode, AssetName: assetName, AssetType: assetType, Category: category, WorkerDisplayName: workerDisplayName, WorkerCode: workerCode}
}

func mapAccessLifecycleTask(row sqlc.HrmsAccessLifecycleTask) *domain.AccessLifecycleTask {
	return accessLifecycleTaskFromParts(row.ID, row.TenantID, row.AccessItemID, row.WorkerProfileID, row.EmployeeID, row.CandidateOnboardingID, row.ExitRequestID, row.TaskType, row.RequestedBy, row.ApprovedBy, row.OwnerUserID, row.ApprovedAt, row.DueDate, row.CompletedAt, row.ExternalReference, row.Status, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, nil, nil, nil, nil, nil, nil)
}

func mapAccessLifecycleTaskList(rows []sqlc.ListAccessLifecycleTasksRow) []*domain.AccessLifecycleTask {
	items := make([]*domain.AccessLifecycleTask, 0, len(rows))
	for _, row := range rows {
		items = append(items, accessLifecycleTaskFromParts(row.ID, row.TenantID, row.AccessItemID, row.WorkerProfileID, row.EmployeeID, row.CandidateOnboardingID, row.ExitRequestID, row.TaskType, row.RequestedBy, row.ApprovedBy, row.OwnerUserID, row.ApprovedAt, row.DueDate, row.CompletedAt, row.ExternalReference, row.Status, row.Notes, row.Metadata, row.Inactive, row.CreatedAt, row.CreatedBy, row.UpdatedAt, row.UpdatedBy, &row.AccessCode, &row.AccessName, &row.AccessType, ptrFromText(row.SystemName), &row.WorkerDisplayName, ptrFromText(row.WorkerCode)))
	}
	return items
}

func accessLifecycleTaskFromParts(id uuid.UUID, tenantID uuid.UUID, accessItemID uuid.UUID, workerProfileID uuid.UUID, employeeID pgtype.UUID, candidateOnboardingID pgtype.UUID, exitRequestID pgtype.UUID, taskType string, requestedBy pgtype.UUID, approvedBy pgtype.UUID, ownerUserID pgtype.UUID, approvedAt pgtype.Timestamptz, dueDate pgtype.Date, completedAt pgtype.Timestamptz, externalReference pgtype.Text, status string, notes pgtype.Text, metadata []byte, inactive bool, createdAt pgtype.Timestamptz, createdBy pgtype.UUID, updatedAt pgtype.Timestamptz, updatedBy pgtype.UUID, accessCode *string, accessName *string, accessType *string, systemName *string, workerDisplayName *string, workerCode *string) *domain.AccessLifecycleTask {
	return &domain.AccessLifecycleTask{ID: id, TenantID: tenantID, AccessItemID: accessItemID, WorkerProfileID: workerProfileID, EmployeeID: ptrFromUUID(employeeID), CandidateOnboardingID: ptrFromUUID(candidateOnboardingID), ExitRequestID: ptrFromUUID(exitRequestID), TaskType: taskType, RequestedBy: ptrFromUUID(requestedBy), ApprovedBy: ptrFromUUID(approvedBy), OwnerUserID: ptrFromUUID(ownerUserID), ApprovedAt: ptrFromTimestamptz(approvedAt), DueDate: ptrFromDate(dueDate), CompletedAt: ptrFromTimestamptz(completedAt), ExternalReference: ptrFromText(externalReference), Status: status, Notes: ptrFromText(notes), Metadata: assetAccessRaw(metadata), Inactive: inactive, CreatedAt: timeFromTimestamptz(createdAt), CreatedBy: ptrFromUUID(createdBy), UpdatedAt: timeFromTimestamptz(updatedAt), UpdatedBy: ptrFromUUID(updatedBy), AccessCode: accessCode, AccessName: accessName, AccessType: accessType, SystemName: systemName, WorkerDisplayName: workerDisplayName, WorkerCode: workerCode}
}

func mapAssetAccessEvent(row sqlc.HrmsAssetAccessEvent) *domain.AssetAccessEvent {
	return &domain.AssetAccessEvent{ID: row.ID, TenantID: row.TenantID, SourceType: row.SourceType, SourceID: ptrFromUUID(row.SourceID), Action: row.Action, FromStatus: ptrFromText(row.FromStatus), ToStatus: ptrFromText(row.ToStatus), Remarks: ptrFromText(row.Remarks), Metadata: assetAccessRaw(row.Metadata), Inactive: row.Inactive, CreatedAt: timeFromTimestamptz(row.CreatedAt), CreatedBy: ptrFromUUID(row.CreatedBy), UpdatedAt: timeFromTimestamptz(row.UpdatedAt), UpdatedBy: ptrFromUUID(row.UpdatedBy)}
}

func mapAssetAccessEvents(rows []sqlc.HrmsAssetAccessEvent) []*domain.AssetAccessEvent {
	items := make([]*domain.AssetAccessEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapAssetAccessEvent(row))
	}
	return items
}

func mapAssetAccessSummary(rows []sqlc.GetAssetAccessSummaryRow) []*domain.AssetAccessSummaryRow {
	items := make([]*domain.AssetAccessSummaryRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, &domain.AssetAccessSummaryRow{Metric: row.Metric, MetricCount: row.MetricCount})
	}
	return items
}

func assetAccessRaw(value []byte) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(value)
}
