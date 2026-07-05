package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateCandidate(ctx context.Context, cmd ports.CandidateCommand) (*domain.Candidate, error) {
	item, err := domain.NewCandidate(candidateInput(cmd))
	if err != nil {
		s.logError("validate candidate create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.candidates.CreateCandidate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create candidate", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListCandidates(ctx context.Context, filter domain.CandidateFilter) (*domain.CandidatePage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate candidate list tenant", err)
		return nil, err
	}
	filter.Search = cleanStringPtr(filter.Search)
	filter.Source = cleanStringPtr(filter.Source)
	filter.Gender = cleanStringPtr(filter.Gender)
	limit, offset := normalizeListWindow(filter.Limit, filter.Offset)
	filter.Limit = limit
	filter.Offset = offset
	items, err := s.candidates.ListCandidates(ctx, filter)
	if err != nil {
		s.logError("list candidates", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.candidates.CountCandidates(ctx, filter)
	if err != nil {
		s.logError("count candidates", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	page := &domain.CandidatePage{Items: items, Total: total, Limit: limit, Offset: offset}
	if int64(offset)+int64(len(items)) < total {
		next := offset + limit
		page.NextOffset = &next
	}
	return page, nil
}

func (s *TenantService) GetCandidate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Candidate, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate candidate get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidCandidateID
		s.logError("validate candidate get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.candidates.GetCandidate(ctx, tenantID, id)
	if err != nil {
		s.logError("get candidate", err, serviceTenantIDField(tenantID), serviceStringField("candidate_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateCandidate(ctx context.Context, cmd ports.CandidateCommand) (*domain.Candidate, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidCandidateID
		s.logError("validate candidate update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetCandidate(ctx, cmd.TenantID, cmd.ID); err != nil {
		return nil, err
	}
	item, err := domain.NewCandidate(candidateInput(cmd))
	if err != nil {
		s.logError("validate candidate update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.candidates.UpdateCandidate(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update candidate", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) DeleteCandidate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetCandidate(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.candidates.DeleteCandidate(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete candidate", err, serviceTenantIDField(tenantID), serviceStringField("candidate_id", id.String()))
		return err
	}
	return nil
}

func candidateInput(cmd ports.CandidateCommand) domain.CandidateInput {
	return domain.CandidateInput{TenantID: cmd.TenantID, Firstname: cmd.Firstname, Lastname: cmd.Lastname, Email: cmd.Email, Phone: cmd.Phone, DOB: cmd.DOB, Gender: cmd.Gender, TotalExperience: cmd.TotalExperience, CurrentCompany: cmd.CurrentCompany, CurrentDesignation: cmd.CurrentDesignation, CurrentSalary: cmd.CurrentSalary, ExpectedSalary: cmd.ExpectedSalary, NoticePeriod: cmd.NoticePeriod, CurrentLocation: cmd.CurrentLocation, PreferredLocation: cmd.PreferredLocation, Source: cmd.Source, ResumeURL: cmd.ResumeURL}
}
