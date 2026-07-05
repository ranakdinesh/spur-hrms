package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreatePolicyType(ctx context.Context, item *domain.PolicyType, actorID *uuid.UUID) (*domain.PolicyType, error) {
	row, err := s.getQueries(ctx).CreatePolicyType(ctx, sqlc.CreatePolicyTypeParams{TenantID: uuidFromPtr(item.TenantID), Name: item.Name, IsSystem: item.IsSystem, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create policy type", err, optionalTenantIDField(item.TenantID), stringField("policy_type_name", item.Name))
	}
	return mapPolicyType(row), nil
}

func (s *Store) ListPolicyTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.PolicyType, error) {
	rows, err := s.getQueries(ctx).ListPolicyTypes(ctx, uuidFromPtr(&tenantID))
	if err != nil {
		return nil, s.logDBError(ctx, "list policy types", err, tenantIDField(tenantID))
	}
	return mapPolicyTypes(rows), nil
}

func (s *Store) GetPolicyType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PolicyType, error) {
	row, err := s.getQueries(ctx).GetPolicyType(ctx, sqlc.GetPolicyTypeParams{TenantID: uuidFromPtr(&tenantID), ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPolicyTypeNotFound
		}
		return nil, s.logDBError(ctx, "get policy type", err, tenantIDField(tenantID), stringField("policy_type_id", id.String()))
	}
	return mapPolicyType(row), nil
}

func (s *Store) UpdatePolicyType(ctx context.Context, item *domain.PolicyType, actorID *uuid.UUID) (*domain.PolicyType, error) {
	row, err := s.getQueries(ctx).UpdatePolicyType(ctx, sqlc.UpdatePolicyTypeParams{TenantID: uuidFromPtr(item.TenantID), ID: item.ID, Name: item.Name, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update policy type", err, optionalTenantIDField(item.TenantID), stringField("policy_type_id", item.ID.String()))
	}
	return mapPolicyType(row), nil
}

func (s *Store) DeletePolicyType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeletePolicyType(ctx, sqlc.SoftDeletePolicyTypeParams{TenantID: uuidFromPtr(&tenantID), ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete policy type", err, tenantIDField(tenantID), stringField("policy_type_id", id.String()))
	}
	return nil
}

func (s *Store) CreateCompanyPolicy(ctx context.Context, item *domain.CompanyPolicy, actorID *uuid.UUID) (*domain.CompanyPolicy, error) {
	row, err := s.getQueries(ctx).CreateCompanyPolicy(ctx, sqlc.CreateCompanyPolicyParams{
		TenantID:     item.TenantID,
		PolicyTypeID: uuidFromPtr(item.PolicyTypeID),
		Title:        item.Title,
		FilePath:     textFromPtr(item.FilePath),
		Description:  textFromPtr(item.Description),
		CreatedBy:    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create company policy", err, tenantIDField(item.TenantID), stringField("policy_title", item.Title))
	}
	return mapCompanyPolicy(row), nil
}

func (s *Store) ListCompanyPolicies(ctx context.Context, tenantID uuid.UUID, policyTypeID *uuid.UUID) ([]*domain.CompanyPolicy, error) {
	queries := s.getQueries(ctx)
	if policyTypeID != nil && *policyTypeID != uuid.Nil {
		rows, err := queries.ListCompanyPoliciesByType(ctx, sqlc.ListCompanyPoliciesByTypeParams{TenantID: tenantID, PolicyTypeID: uuidFromPtr(policyTypeID)})
		if err != nil {
			return nil, s.logDBError(ctx, "list company policies by type", err, tenantIDField(tenantID), stringField("policy_type_id", policyTypeID.String()))
		}
		return mapCompanyPolicies(rows), nil
	}
	rows, err := queries.ListCompanyPolicies(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "list company policies", err, tenantIDField(tenantID))
	}
	return mapCompanyPolicies(rows), nil
}

func (s *Store) GetCompanyPolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompanyPolicy, error) {
	row, err := s.getQueries(ctx).GetCompanyPolicy(ctx, sqlc.GetCompanyPolicyParams{TenantID: tenantID, ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCompanyPolicyNotFound
		}
		return nil, s.logDBError(ctx, "get company policy", err, tenantIDField(tenantID), stringField("policy_id", id.String()))
	}
	return mapCompanyPolicy(row), nil
}

func (s *Store) UpdateCompanyPolicy(ctx context.Context, item *domain.CompanyPolicy, actorID *uuid.UUID) (*domain.CompanyPolicy, error) {
	row, err := s.getQueries(ctx).UpdateCompanyPolicy(ctx, sqlc.UpdateCompanyPolicyParams{
		TenantID:     item.TenantID,
		ID:           item.ID,
		PolicyTypeID: uuidFromPtr(item.PolicyTypeID),
		Title:        item.Title,
		FilePath:     textFromPtr(item.FilePath),
		Description:  textFromPtr(item.Description),
		UpdatedBy:    uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "update company policy", err, tenantIDField(item.TenantID), stringField("policy_id", item.ID.String()))
	}
	return mapCompanyPolicy(row), nil
}

func (s *Store) DeleteCompanyPolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCompanyPolicy(ctx, sqlc.SoftDeleteCompanyPolicyParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete company policy", err, tenantIDField(tenantID), stringField("policy_id", id.String()))
	}
	return nil
}
