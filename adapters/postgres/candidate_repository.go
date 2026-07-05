package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateCandidate(ctx context.Context, item *domain.Candidate, actorID *uuid.UUID) (*domain.Candidate, error) {
	row, err := s.getQueries(ctx).CreateCandidate(ctx, candidateCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create candidate", fmt.Errorf("hrms: create candidate: %w", err), tenantIDField(item.TenantID))
	}
	return mapCandidate(row), nil
}

func (s *Store) ListCandidates(ctx context.Context, filter domain.CandidateFilter) ([]*domain.Candidate, error) {
	rows, err := s.getQueries(ctx).ListCandidates(ctx, sqlc.ListCandidatesParams{TenantID: filter.TenantID, Search: textFromPtr(filter.Search), Source: textFromPtr(filter.Source), Gender: textFromPtr(filter.Gender), OffsetRows: filter.Offset, LimitRows: filter.Limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list candidates", err, tenantIDField(filter.TenantID))
	}
	return mapCandidates(rows), nil
}

func (s *Store) CountCandidates(ctx context.Context, filter domain.CandidateFilter) (int64, error) {
	count, err := s.getQueries(ctx).CountCandidates(ctx, sqlc.CountCandidatesParams{TenantID: filter.TenantID, Search: textFromPtr(filter.Search), Source: textFromPtr(filter.Source), Gender: textFromPtr(filter.Gender)})
	if err != nil {
		return 0, s.logDBError(ctx, "count candidates", err, tenantIDField(filter.TenantID))
	}
	return count, nil
}

func (s *Store) GetCandidate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Candidate, error) {
	row, err := s.getQueries(ctx).GetCandidate(ctx, sqlc.GetCandidateParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get candidate", fmt.Errorf("hrms: get candidate: %w", err), tenantIDField(tenantID), stringField("candidate_id", id.String()))
	}
	return mapCandidate(row), nil
}

func (s *Store) UpdateCandidate(ctx context.Context, item *domain.Candidate, actorID *uuid.UUID) (*domain.Candidate, error) {
	row, err := s.getQueries(ctx).UpdateCandidate(ctx, candidateUpdateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "update candidate", fmt.Errorf("hrms: update candidate: %w", err), tenantIDField(item.TenantID), stringField("candidate_id", item.ID.String()))
	}
	return mapCandidate(row), nil
}

func (s *Store) DeleteCandidate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCandidate(ctx, sqlc.SoftDeleteCandidateParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete candidate", fmt.Errorf("hrms: delete candidate: %w", err), tenantIDField(tenantID), stringField("candidate_id", id.String()))
	}
	return nil
}

func candidateCreateParams(item *domain.Candidate, actorID *uuid.UUID) sqlc.CreateCandidateParams {
	return sqlc.CreateCandidateParams{TenantID: item.TenantID, Firstname: textFromPtr(item.Firstname), Lastname: textFromPtr(item.Lastname), Email: textFromPtr(item.Email), Phone: textFromPtr(item.Phone), Dob: dateFromPtr(item.DOB), Gender: textFromPtr(item.Gender), TotalExperience: numericFromFloatPtr(item.TotalExperience), CurrentCompany: textFromPtr(item.CurrentCompany), CurrentDesignation: textFromPtr(item.CurrentDesignation), CurrentSalary: numericFromFloatPtr(item.CurrentSalary), ExpectedSalary: numericFromFloatPtr(item.ExpectedSalary), NoticePeriod: int4FromPtr(item.NoticePeriod), CurrentLocation: textFromPtr(item.CurrentLocation), PreferredLocation: textFromPtr(item.PreferredLocation), Source: textFromPtr(item.Source), ResumeUrl: textFromPtr(item.ResumeURL), CreatedBy: uuidFromPtr(actorID)}
}

func candidateUpdateParams(item *domain.Candidate, actorID *uuid.UUID) sqlc.UpdateCandidateParams {
	return sqlc.UpdateCandidateParams{TenantID: item.TenantID, ID: item.ID, Firstname: textFromPtr(item.Firstname), Lastname: textFromPtr(item.Lastname), Email: textFromPtr(item.Email), Phone: textFromPtr(item.Phone), Dob: dateFromPtr(item.DOB), Gender: textFromPtr(item.Gender), TotalExperience: numericFromFloatPtr(item.TotalExperience), CurrentCompany: textFromPtr(item.CurrentCompany), CurrentDesignation: textFromPtr(item.CurrentDesignation), CurrentSalary: numericFromFloatPtr(item.CurrentSalary), ExpectedSalary: numericFromFloatPtr(item.ExpectedSalary), NoticePeriod: int4FromPtr(item.NoticePeriod), CurrentLocation: textFromPtr(item.CurrentLocation), PreferredLocation: textFromPtr(item.PreferredLocation), Source: textFromPtr(item.Source), ResumeUrl: textFromPtr(item.ResumeURL), UpdatedBy: uuidFromPtr(actorID)}
}
