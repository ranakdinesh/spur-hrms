package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateWorkerProfile(ctx context.Context, item *domain.WorkerProfile, actorID *uuid.UUID) (*domain.WorkerProfile, error) {
	row, err := s.getQueries(ctx).CreateWorkerProfile(ctx, createWorkerProfileParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create worker profile", err, tenantIDField(item.TenantID), stringField("worker_type_id", item.WorkerTypeID.String()))
	}
	return mapWorkerProfile(row), nil
}

func (s *Store) UpdateWorkerProfile(ctx context.Context, item *domain.WorkerProfile, actorID *uuid.UUID) (*domain.WorkerProfile, error) {
	params := sqlc.UpdateWorkerProfileParams{
		TenantID:           item.TenantID,
		ID:                 item.ID,
		WorkerTypeID:       item.WorkerTypeID,
		EmployeeID:         uuidFromPtr(item.EmployeeID),
		EmployeeUserID:     uuidFromPtr(item.EmployeeUserID),
		WorkerCode:         textFromPtr(item.WorkerCode),
		DisplayName:        item.DisplayName,
		LegalName:          textFromPtr(item.LegalName),
		Email:              textFromPtr(item.Email),
		Mobile:             textFromPtr(item.Mobile),
		ProfileStatus:      item.ProfileStatus,
		StartDate:          dateFromPtr(item.StartDate),
		EndDate:            dateFromPtr(item.EndDate),
		BranchID:           uuidFromPtr(item.BranchID),
		DepartmentID:       uuidFromPtr(item.DepartmentID),
		ReportingManagerID: uuidFromPtr(item.ReportingManagerID),
		WorkLocationLabel:  textFromPtr(item.WorkLocationLabel),
		SourcePartner:      textFromPtr(item.SourcePartner),
		ExternalReference:  textFromPtr(item.ExternalReference),
		ComplianceStatus:   item.ComplianceStatus,
		PayrollStatus:      item.PayrollStatus,
		Notes:              textFromPtr(item.Notes),
		Metadata:           jsonBytesFromRaw(item.Metadata),
		UpdatedBy:          uuidFromPtr(actorID),
	}
	row, err := s.getQueries(ctx).UpdateWorkerProfile(ctx, params)
	if err != nil {
		return nil, s.logDBError(ctx, "update worker profile", err, tenantIDField(item.TenantID), stringField("worker_profile_id", item.ID.String()))
	}
	return mapWorkerProfile(row), nil
}

func (s *Store) GetWorkerProfile(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.WorkerProfile, error) {
	row, err := s.getQueries(ctx).GetWorkerProfile(ctx, sqlc.GetWorkerProfileParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get worker profile", err, tenantIDField(tenantID), stringField("worker_profile_id", id.String()))
	}
	return mapWorkerProfile(row), nil
}

func (s *Store) GetWorkerProfileByEmployeeID(ctx context.Context, tenantID uuid.UUID, employeeID uuid.UUID) (*domain.WorkerProfile, error) {
	row, err := s.getQueries(ctx).GetWorkerProfileByEmployeeID(ctx, sqlc.GetWorkerProfileByEmployeeIDParams{TenantID: tenantID, EmployeeID: uuidFromPtr(&employeeID)})
	if err != nil {
		return nil, s.logDBError(ctx, "get worker profile by employee", err, tenantIDField(tenantID), stringField("employee_id", employeeID.String()))
	}
	return mapWorkerProfile(row), nil
}

func (s *Store) ListWorkerProfiles(ctx context.Context, filter domain.WorkerProfileFilter) ([]*domain.WorkerProfileListItem, error) {
	rows, err := s.getQueries(ctx).ListWorkerProfiles(ctx, sqlc.ListWorkerProfilesParams{
		TenantID:            filter.TenantID,
		WorkerTypeID:        uuidFromPtr(filter.WorkerTypeID),
		ClassificationGroup: textFromPtr(filter.ClassificationGroup),
		ProfileStatus:       textFromPtr(filter.ProfileStatus),
		Search:              textFromPtr(filter.Search),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "list worker profiles", err, tenantIDField(filter.TenantID))
	}
	return mapWorkerProfileListItems(rows), nil
}

func (s *Store) DeleteWorkerProfile(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteWorkerProfile(ctx, sqlc.SoftDeleteWorkerProfileParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete worker profile", err, tenantIDField(tenantID), stringField("worker_profile_id", id.String()))
	}
	return nil
}

func createWorkerProfileParams(item *domain.WorkerProfile, actorID *uuid.UUID) sqlc.CreateWorkerProfileParams {
	return sqlc.CreateWorkerProfileParams{
		TenantID:           item.TenantID,
		WorkerTypeID:       item.WorkerTypeID,
		EmployeeID:         uuidFromPtr(item.EmployeeID),
		EmployeeUserID:     uuidFromPtr(item.EmployeeUserID),
		WorkerCode:         textFromPtr(item.WorkerCode),
		DisplayName:        item.DisplayName,
		LegalName:          textFromPtr(item.LegalName),
		Email:              textFromPtr(item.Email),
		Mobile:             textFromPtr(item.Mobile),
		ProfileStatus:      item.ProfileStatus,
		StartDate:          dateFromPtr(item.StartDate),
		EndDate:            dateFromPtr(item.EndDate),
		BranchID:           uuidFromPtr(item.BranchID),
		DepartmentID:       uuidFromPtr(item.DepartmentID),
		ReportingManagerID: uuidFromPtr(item.ReportingManagerID),
		WorkLocationLabel:  textFromPtr(item.WorkLocationLabel),
		SourcePartner:      textFromPtr(item.SourcePartner),
		ExternalReference:  textFromPtr(item.ExternalReference),
		ComplianceStatus:   item.ComplianceStatus,
		PayrollStatus:      item.PayrollStatus,
		Notes:              textFromPtr(item.Notes),
		Metadata:           jsonBytesFromRaw(item.Metadata),
		CreatedBy:          uuidFromPtr(actorID),
	}
}
