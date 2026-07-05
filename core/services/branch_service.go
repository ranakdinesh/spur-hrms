package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateBranch(ctx context.Context, cmd ports.BranchCommand) (*domain.Branch, error) {
	branch, err := domain.NewBranch(domain.BranchInput{
		TenantID:            cmd.TenantID,
		Name:                cmd.Name,
		Address:             cmd.Address,
		City:                cmd.City,
		State:               cmd.State,
		Country:             cmd.Country,
		Pincode:             cmd.Pincode,
		Phone:               cmd.Phone,
		BranchManagerUserID: cmd.BranchManagerUserID,
		HRUserID:            cmd.HRUserID,
		AccountsUserID:      cmd.AccountsUserID,
	})
	if err != nil {
		s.logError("validate branch create", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_name", cmd.Name))
		return nil, err
	}
	result, err := s.branches.CreateBranch(ctx, branch, cmd.ActorID)
	if err != nil {
		s.logError("create branch", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_name", branch.Name))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("branch_id", result.ID.String()).Msg("hrms: branch created")
	return result, nil
}

func (s *TenantService) ListBranches(ctx context.Context, tenantID uuid.UUID) ([]*domain.Branch, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate branch list tenant", err)
		return nil, err
	}
	result, err := s.branches.ListBranches(ctx, tenantID)
	if err != nil {
		s.logError("list branches", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID) (*domain.Branch, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate branch get tenant", err)
		return nil, err
	}
	if branchID == uuid.Nil {
		err := domain.ErrInvalidBranchID
		s.logError("validate branch get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	result, err := s.branches.GetBranch(ctx, tenantID, branchID)
	if err != nil {
		s.logError("get branch", err, serviceTenantIDField(tenantID), serviceStringField("branch_id", branchID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateBranch(ctx context.Context, cmd ports.BranchCommand) (*domain.Branch, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidBranchID
		s.logError("validate branch update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	branch, err := domain.NewBranch(domain.BranchInput{
		TenantID:            cmd.TenantID,
		Name:                cmd.Name,
		Address:             cmd.Address,
		City:                cmd.City,
		State:               cmd.State,
		Country:             cmd.Country,
		Pincode:             cmd.Pincode,
		Phone:               cmd.Phone,
		BranchManagerUserID: cmd.BranchManagerUserID,
		HRUserID:            cmd.HRUserID,
		AccountsUserID:      cmd.AccountsUserID,
	})
	if err != nil {
		s.logError("validate branch update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", cmd.ID.String()), serviceStringField("branch_name", cmd.Name))
		return nil, err
	}
	branch.ID = cmd.ID
	result, err := s.branches.UpdateBranch(ctx, branch, cmd.ActorID)
	if err != nil {
		s.logError("update branch", err, serviceTenantIDField(cmd.TenantID), serviceStringField("branch_id", cmd.ID.String()))
		return nil, err
	}
	s.log.Info().Str("tenant_id", result.TenantID.String()).Str("branch_id", result.ID.String()).Msg("hrms: branch updated")
	return result, nil
}

func (s *TenantService) DeleteBranch(ctx context.Context, tenantID uuid.UUID, branchID uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate branch delete tenant", err)
		return err
	}
	if branchID == uuid.Nil {
		err := domain.ErrInvalidBranchID
		s.logError("validate branch delete id", err, serviceTenantIDField(tenantID))
		return err
	}
	if err := s.branches.DeleteBranch(ctx, tenantID, branchID, actorID); err != nil {
		s.logError("delete branch", err, serviceTenantIDField(tenantID), serviceStringField("branch_id", branchID.String()))
		return err
	}
	s.log.Info().Str("tenant_id", tenantID.String()).Str("branch_id", branchID.String()).Msg("hrms: branch deactivated")
	return nil
}
