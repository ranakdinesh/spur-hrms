package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateCandidateApplication(ctx context.Context, cmd ports.CandidateApplicationCommand) (*domain.CandidateApplication, error) {
	item, err := domain.NewCandidateApplication(candidateApplicationInput(cmd))
	if err != nil {
		s.logError("validate candidate application create", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if err := s.validateCandidateApplicationReferences(ctx, item); err != nil {
		return nil, err
	}
	result, err := s.candidates.CreateCandidateApplication(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create candidate application", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.candidates.CreateCandidateApplicationEvent(ctx, &domain.CandidateApplicationEvent{TenantID: cmd.TenantID, ApplicationID: result.ID, ToStatus: result.Status, Action: "created", Remarks: result.Comments}, cmd.ActorID); err != nil {
		s.logError("create candidate application event", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_application_id", result.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) LinkCandidateApplicantAccount(ctx context.Context, cmd ports.CandidateApplicantAccountCommand) (*domain.CandidateApplicantAccount, error) {
	item, err := domain.NewCandidateApplicantAccount(candidateApplicantAccountInput(cmd))
	if err != nil {
		s.logError("validate candidate applicant account", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	if _, err := s.GetCandidate(ctx, item.TenantID, item.CandidateID); err != nil {
		return nil, err
	}
	result, err := s.candidates.CreateCandidateApplicantAccount(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create candidate applicant account", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_id", cmd.CandidateID.String()), serviceStringField("user_id", cmd.UserID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) GetApplicantPortal(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.ApplicantPortal, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate applicant portal tenant", err)
		return nil, err
	}
	if userID == uuid.Nil {
		err := domain.ErrInvalidApplicantUserID
		s.logError("validate applicant portal user", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	account, err := s.candidates.GetCandidateApplicantAccountByUser(ctx, tenantID, userID)
	if err != nil {
		s.logError("get candidate applicant account", err, serviceTenantIDField(tenantID), serviceStringField("user_id", userID.String()))
		return nil, err
	}
	candidate, err := s.GetCandidate(ctx, tenantID, account.CandidateID)
	if err != nil {
		return nil, err
	}
	applications, err := s.candidates.ListCandidateApplications(ctx, domain.CandidateApplicationFilter{
		TenantID:    tenantID,
		CandidateID: &account.CandidateID,
		Limit:       100,
		Offset:      0,
	})
	if err != nil {
		s.logError("list applicant portal applications", err, serviceTenantIDField(tenantID), serviceStringField("candidate_id", account.CandidateID.String()))
		return nil, err
	}
	return &domain.ApplicantPortal{Account: account, Candidate: candidate, Applications: applications}, nil
}

func (s *TenantService) ListCandidateApplications(ctx context.Context, filter domain.CandidateApplicationFilter) (*domain.CandidateApplicationPage, error) {
	if filter.TenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate candidate application list tenant", err)
		return nil, err
	}
	filter.Search = cleanStringPtr(filter.Search)
	filter.Status = cleanStringPtr(filter.Status)
	if filter.Status != nil {
		if _, err := domain.ValidateCandidateApplicationStatus(filter.Status); err != nil {
			s.logError("validate candidate application list status", err, serviceTenantIDField(filter.TenantID), serviceStringField("status", *filter.Status))
			return nil, err
		}
	}
	limit, offset := normalizeListWindow(filter.Limit, filter.Offset)
	filter.Limit = limit
	filter.Offset = offset
	items, err := s.candidates.ListCandidateApplications(ctx, filter)
	if err != nil {
		s.logError("list candidate applications", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.candidates.CountCandidateApplications(ctx, filter)
	if err != nil {
		s.logError("count candidate applications", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	page := &domain.CandidateApplicationPage{Items: items, Total: total, Limit: limit, Offset: offset}
	if int64(offset)+int64(len(items)) < total {
		next := offset + limit
		page.NextOffset = &next
	}
	return page, nil
}

func (s *TenantService) GetCandidateApplication(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CandidateApplication, error) {
	if tenantID == uuid.Nil {
		err := domain.ErrInvalidTenantID
		s.logError("validate candidate application get tenant", err)
		return nil, err
	}
	if id == uuid.Nil {
		err := domain.ErrInvalidCandidateApplicationID
		s.logError("validate candidate application get id", err, serviceTenantIDField(tenantID))
		return nil, err
	}
	item, err := s.candidates.GetCandidateApplication(ctx, tenantID, id)
	if err != nil {
		s.logError("get candidate application", err, serviceTenantIDField(tenantID), serviceStringField("candidate_application_id", id.String()))
		return nil, err
	}
	return item, nil
}

func (s *TenantService) UpdateCandidateApplication(ctx context.Context, cmd ports.CandidateApplicationCommand) (*domain.CandidateApplication, error) {
	if cmd.ID == uuid.Nil {
		err := domain.ErrInvalidCandidateApplicationID
		s.logError("validate candidate application update id", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	existing, err := s.GetCandidateApplication(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	item, err := domain.NewCandidateApplication(candidateApplicationInput(cmd))
	if err != nil {
		s.logError("validate candidate application update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_application_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	if err := s.validateCandidateApplicationReferences(ctx, item); err != nil {
		return nil, err
	}
	result, err := s.candidates.UpdateCandidateApplication(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update candidate application", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_application_id", cmd.ID.String()))
		return nil, err
	}
	if existing.Status != result.Status {
		if _, err := s.candidates.CreateCandidateApplicationEvent(ctx, &domain.CandidateApplicationEvent{TenantID: cmd.TenantID, ApplicationID: result.ID, FromStatus: &existing.Status, ToStatus: result.Status, Action: "updated", Reason: stageReason(result), Remarks: result.Comments}, cmd.ActorID); err != nil {
			s.logError("create candidate application update event", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_application_id", result.ID.String()))
			return nil, err
		}
	}
	return result, nil
}

func (s *TenantService) MoveCandidateApplicationStatus(ctx context.Context, cmd ports.CandidateApplicationMoveCommand) (*domain.CandidateApplication, error) {
	status, err := domain.ValidateCandidateApplicationStatus(&cmd.Status)
	if err != nil {
		s.logError("validate candidate application move status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("status", cmd.Status))
		return nil, err
	}
	existing, err := s.GetCandidateApplication(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	result, err := s.candidates.MoveCandidateApplicationStatus(ctx, cmd.TenantID, cmd.ID, *status, cleanStringPtr(cmd.Comments), cleanStringPtr(cmd.Reason), cmd.ActorID)
	if err != nil {
		s.logError("move candidate application status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_application_id", cmd.ID.String()), serviceStringField("status", *status))
		return nil, err
	}
	if _, err := s.candidates.CreateCandidateApplicationEvent(ctx, &domain.CandidateApplicationEvent{TenantID: cmd.TenantID, ApplicationID: result.ID, FromStatus: &existing.Status, ToStatus: result.Status, Action: "moved", Reason: cleanStringPtr(cmd.Reason), Remarks: cleanStringPtr(cmd.Comments)}, cmd.ActorID); err != nil {
		s.logError("create candidate application move event", err, serviceTenantIDField(cmd.TenantID), serviceStringField("candidate_application_id", result.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListCandidateApplicationEvents(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID) ([]*domain.CandidateApplicationEvent, error) {
	if _, err := s.GetCandidateApplication(ctx, tenantID, applicationID); err != nil {
		return nil, err
	}
	items, err := s.candidates.ListCandidateApplicationEvents(ctx, tenantID, applicationID)
	if err != nil {
		s.logError("list candidate application events", err, serviceTenantIDField(tenantID), serviceStringField("candidate_application_id", applicationID.String()))
		return nil, err
	}
	return items, nil
}

func (s *TenantService) DeleteCandidateApplication(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if _, err := s.GetCandidateApplication(ctx, tenantID, id); err != nil {
		return err
	}
	if err := s.candidates.DeleteCandidateApplication(ctx, tenantID, id, actorID); err != nil {
		s.logError("delete candidate application", err, serviceTenantIDField(tenantID), serviceStringField("candidate_application_id", id.String()))
		return err
	}
	return nil
}

func (s *TenantService) validateCandidateApplicationReferences(ctx context.Context, item *domain.CandidateApplication) error {
	if item.CandidateID != nil {
		if _, err := s.GetCandidate(ctx, item.TenantID, *item.CandidateID); err != nil {
			return err
		}
	}
	if item.JobPostingID != nil {
		if _, err := s.GetJobPosting(ctx, item.TenantID, *item.JobPostingID); err != nil {
			return err
		}
	}
	return nil
}

func candidateApplicationInput(cmd ports.CandidateApplicationCommand) domain.CandidateApplicationInput {
	return domain.CandidateApplicationInput{TenantID: cmd.TenantID, CandidateID: cmd.CandidateID, JobPostingID: cmd.JobPostingID, ResumeURL: cmd.ResumeURL, CoverLetter: cmd.CoverLetter, CurrentCTC: cmd.CurrentCTC, ExpectedCTC: cmd.ExpectedCTC, NoticePeriod: cmd.NoticePeriod, ReferredBy: cmd.ReferredBy, Source: cmd.Source, SourceDetail: cmd.SourceDetail, Status: cmd.Status, Comments: cmd.Comments, AppliedAt: cmd.AppliedAt, RejectionReason: cmd.RejectionReason, WithdrawalReason: cmd.WithdrawalReason, DuplicateOfApplicationID: cmd.DuplicateOfApplicationID}
}

func candidateApplicantAccountInput(cmd ports.CandidateApplicantAccountCommand) domain.CandidateApplicantAccountInput {
	return domain.CandidateApplicantAccountInput{TenantID: cmd.TenantID, CandidateID: cmd.CandidateID, UserID: cmd.UserID, Email: cmd.Email, Status: cmd.Status, ConsentAt: cmd.ConsentAt, ConsentIP: cmd.ConsentIP, Metadata: cmd.Metadata}
}

func stageReason(item *domain.CandidateApplication) *string {
	if item == nil {
		return nil
	}
	if item.Status == domain.CandidateApplicationStatusRejected {
		return item.RejectionReason
	}
	if item.Status == domain.CandidateApplicationStatusWithdrawn {
		return item.WithdrawalReason
	}
	return nil
}
