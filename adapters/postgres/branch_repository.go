package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateBranch(ctx context.Context, branch *domain.Branch, actorID *uuid.UUID) (*domain.Branch, error) {
	row, err := s.getQueries(ctx).CreateBranch(ctx, sqlc.CreateBranchParams{
		TenantID:            branch.TenantID,
		BranchName:          branch.Name,
		Address:             textFromPtr(branch.Address),
		City:                textFromPtr(branch.City),
		State:               textFromPtr(branch.State),
		Country:             textFromPtr(branch.Country),
		Pincode:             textFromPtr(branch.Pincode),
		Phone:               textFromPtr(branch.Phone),
		BranchManagerUserID: uuidFromPtr(branch.BranchManagerUserID),
		HrUserID:            uuidFromPtr(branch.HRUserID),
		AccountsUserID:      uuidFromPtr(branch.AccountsUserID),
		CreatedBy:           uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create branch", err, tenantIDField(branch.TenantID), stringField("branch_name", branch.Name))
	}
	return mapBranch(row), nil
}

func (s *Store) ListBranches(ctx context.Context, tenantID uuid.UUID) ([]*domain.Branch, error) {
	rows, err := s.getQueries(ctx).ListBranches(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list branches", err, tenantIDField(tenantID))
	}
	return mapBranches(rows), nil
}

func (s *Store) GetBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID) (*domain.Branch, error) {
	row, err := s.getQueries(ctx).GetBranch(ctx, sqlc.GetBranchParams{
		TenantID: tenantID,
		ID:       branchID,
	})
	if err != nil {
		return nil, s.logDBError(ctx, "get branch", err, tenantIDField(tenantID), stringField("branch_id", branchID.String()))
	}
	return mapBranch(row), nil
}

func (s *Store) UpdateBranch(ctx context.Context, branch *domain.Branch, actorID *uuid.UUID) (*domain.Branch, error) {
	row, err := s.getQueries(ctx).UpdateBranch(ctx, sqlc.UpdateBranchParams{
		TenantID:            branch.TenantID,
		ID:                  branch.ID,
		BranchName:          branch.Name,
		Address:             textFromPtr(branch.Address),
		City:                textFromPtr(branch.City),
		State:               textFromPtr(branch.State),
		Country:             textFromPtr(branch.Country),
		Pincode:             textFromPtr(branch.Pincode),
		Phone:               textFromPtr(branch.Phone),
		BranchManagerUserID: uuidFromPtr(branch.BranchManagerUserID),
		HrUserID:            uuidFromPtr(branch.HRUserID),
		AccountsUserID:      uuidFromPtr(branch.AccountsUserID),
		UpdatedBy:           uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update branch", err, tenantIDField(branch.TenantID), stringField("branch_id", branch.ID.String()))
	}
	return mapBranch(row), nil
}

func (s *Store) DeleteBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteBranch(ctx, sqlc.SoftDeleteBranchParams{
		TenantID:  tenantID,
		ID:        branchID,
		UpdatedBy: uuidFromPtr(actorID),
	}); err != nil {
		return s.logDBError(ctx, "delete branch", err, tenantIDField(tenantID), stringField("branch_id", branchID.String()))
	}
	return nil
}
