package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateTalentMarketplaceOpportunity(ctx context.Context, cmd ports.TalentMarketplaceOpportunityCommand) (*domain.TalentMarketplaceOpportunity, error) {
	item, err := s.prepareTalentOpportunity(ctx, cmd)
	if err != nil {
		s.logError("validate talent marketplace opportunity", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.talentMarketplace.CreateTalentMarketplaceOpportunity(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create talent marketplace opportunity", err, serviceTenantIDField(cmd.TenantID), serviceStringField("title", cmd.Title))
		return nil, err
	}
	_, _ = s.talentMarketplace.CreateTalentMarketplaceEvent(ctx, talentMarketplaceEvent(cmd.TenantID, &result.ID, nil, cmd.ActorID, "opportunity_created", nil, &result.Status, nil))
	return result, nil
}

func (s *TenantService) UpdateTalentMarketplaceOpportunity(ctx context.Context, cmd ports.TalentMarketplaceOpportunityCommand) (*domain.TalentMarketplaceOpportunity, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidTalentOpportunity
	}
	existing, err := s.talentMarketplace.GetTalentMarketplaceOpportunity(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	item, err := s.prepareTalentOpportunity(ctx, cmd)
	if err != nil {
		s.logError("validate talent marketplace opportunity update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("opportunity_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.talentMarketplace.UpdateTalentMarketplaceOpportunity(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update talent marketplace opportunity", err, serviceTenantIDField(cmd.TenantID), serviceStringField("opportunity_id", cmd.ID.String()))
		return nil, err
	}
	if existing.Status != result.Status {
		_, _ = s.talentMarketplace.CreateTalentMarketplaceEvent(ctx, talentMarketplaceEvent(cmd.TenantID, &result.ID, nil, cmd.ActorID, "opportunity_status_changed", &existing.Status, &result.Status, nil))
	}
	return result, nil
}

func (s *TenantService) ActivateTalentMarketplaceFallback(ctx context.Context, cmd ports.TalentMarketplaceFallbackCommand) (*domain.TalentMarketplaceOpportunity, error) {
	status, err := domain.ValidateTalentFallbackStatus(cmd.Status)
	if err != nil {
		return nil, err
	}
	existing, err := s.talentMarketplace.GetTalentMarketplaceOpportunity(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	result, err := s.talentMarketplace.UpdateTalentMarketplaceOpportunityFallback(ctx, cmd.TenantID, cmd.ID, status, cmd.ActorID)
	if err != nil {
		s.logError("activate candidate fallback", err, serviceTenantIDField(cmd.TenantID), serviceStringField("opportunity_id", cmd.ID.String()))
		return nil, err
	}
	_, _ = s.talentMarketplace.CreateTalentMarketplaceEvent(ctx, talentMarketplaceEvent(cmd.TenantID, &result.ID, nil, cmd.ActorID, "candidate_fallback_updated", &existing.CandidateFallbackStatus, &result.CandidateFallbackStatus, cmd.Notes))
	return result, nil
}

func (s *TenantService) GetTalentMarketplaceOpportunity(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.TalentMarketplaceOpportunity, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidTalentOpportunity
	}
	return s.talentMarketplace.GetTalentMarketplaceOpportunity(ctx, tenantID, id)
}

func (s *TenantService) ListTalentMarketplaceOpportunities(ctx context.Context, filter domain.TalentMarketplaceOpportunityFilter) ([]*domain.TalentMarketplaceOpportunity, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	filter.Search = domain.NormalizeTalentMarketplaceSearch(filter.Search)
	return s.talentMarketplace.ListTalentMarketplaceOpportunities(ctx, filter)
}

func (s *TenantService) DeleteTalentMarketplaceOpportunity(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidTalentOpportunity
	}
	return s.talentMarketplace.DeleteTalentMarketplaceOpportunity(ctx, tenantID, id, actorID)
}

func (s *TenantService) CreateTalentMarketplaceApplication(ctx context.Context, cmd ports.TalentMarketplaceApplicationCommand) (*domain.TalentMarketplaceApplication, error) {
	if _, err := s.GetTalentMarketplaceOpportunity(ctx, cmd.TenantID, cmd.OpportunityID); err != nil {
		return nil, err
	}
	if _, err := s.GetWorkerProfile(ctx, cmd.TenantID, cmd.WorkerProfileID); err != nil {
		return nil, err
	}
	item, err := domain.NewTalentMarketplaceApplication(domain.TalentMarketplaceApplicationInput{TenantID: cmd.TenantID, OpportunityID: cmd.OpportunityID, WorkerProfileID: cmd.WorkerProfileID, Status: cmd.Status, MatchScore: cmd.MatchScore, MatchReasons: cmd.MatchReasons, WorkerNote: cmd.WorkerNote, ManagerNote: cmd.ManagerNote})
	if err != nil {
		s.logError("validate talent marketplace application", err, serviceTenantIDField(cmd.TenantID), serviceStringField("opportunity_id", cmd.OpportunityID.String()))
		return nil, err
	}
	result, err := s.talentMarketplace.CreateTalentMarketplaceApplication(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create talent marketplace application", err, serviceTenantIDField(cmd.TenantID), serviceStringField("opportunity_id", cmd.OpportunityID.String()), serviceStringField("worker_profile_id", cmd.WorkerProfileID.String()))
		return nil, err
	}
	_, _ = s.talentMarketplace.CreateTalentMarketplaceEvent(ctx, talentMarketplaceEvent(cmd.TenantID, &result.OpportunityID, &result.ID, cmd.ActorID, "application_created", nil, &result.Status, nil))
	return result, nil
}

func (s *TenantService) UpdateTalentMarketplaceApplicationStatus(ctx context.Context, cmd ports.TalentMarketplaceApplicationStatusCommand) (*domain.TalentMarketplaceApplication, error) {
	status, err := domain.ValidateTalentApplicationStatus(cmd.Status)
	if err != nil {
		return nil, err
	}
	existing, err := s.talentMarketplace.GetTalentMarketplaceApplication(ctx, cmd.TenantID, cmd.ID)
	if err != nil {
		return nil, err
	}
	result, err := s.talentMarketplace.UpdateTalentMarketplaceApplicationStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.WorkerNote, cmd.ManagerNote, cmd.ActorID)
	if err != nil {
		s.logError("update talent marketplace application status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("application_id", cmd.ID.String()))
		return nil, err
	}
	_, _ = s.talentMarketplace.CreateTalentMarketplaceEvent(ctx, talentMarketplaceEvent(cmd.TenantID, &result.OpportunityID, &result.ID, cmd.ActorID, "application_status_changed", &existing.Status, &result.Status, firstNonNilString(cmd.ManagerNote, cmd.WorkerNote)))
	return result, nil
}

func (s *TenantService) GetTalentMarketplaceApplication(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.TalentMarketplaceApplication, error) {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return nil, domain.ErrInvalidTalentApplication
	}
	return s.talentMarketplace.GetTalentMarketplaceApplication(ctx, tenantID, id)
}

func (s *TenantService) ListTalentMarketplaceApplications(ctx context.Context, filter domain.TalentMarketplaceApplicationFilter) ([]*domain.TalentMarketplaceApplication, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.talentMarketplace.ListTalentMarketplaceApplications(ctx, filter)
}

func (s *TenantService) ListTalentMarketplaceRecommendations(ctx context.Context, tenantID uuid.UUID, opportunityID uuid.UUID) ([]*domain.TalentMarketplaceRecommendation, error) {
	if tenantID == uuid.Nil || opportunityID == uuid.Nil {
		return nil, domain.ErrInvalidTalentOpportunity
	}
	return s.talentMarketplace.ListTalentMarketplaceRecommendations(ctx, tenantID, opportunityID)
}

func (s *TenantService) ListTalentMarketplaceEvents(ctx context.Context, filter domain.TalentMarketplaceEventFilter) ([]*domain.TalentMarketplaceEvent, error) {
	if filter.TenantID == uuid.Nil {
		return nil, domain.ErrInvalidTenantID
	}
	return s.talentMarketplace.ListTalentMarketplaceEvents(ctx, filter)
}

func (s *TenantService) prepareTalentOpportunity(ctx context.Context, cmd ports.TalentMarketplaceOpportunityCommand) (*domain.TalentMarketplaceOpportunity, error) {
	if cmd.ProjectID != nil {
		if _, err := s.GetProject(ctx, cmd.TenantID, *cmd.ProjectID); err != nil {
			return nil, err
		}
	}
	if cmd.EngagementID != nil {
		if _, err := s.GetEngagement(ctx, cmd.TenantID, *cmd.EngagementID); err != nil {
			return nil, err
		}
	}
	if cmd.SourceRequirementID != nil {
		if _, err := s.GetProjectSkillRequirement(ctx, cmd.TenantID, *cmd.SourceRequirementID); err != nil {
			return nil, err
		}
	}
	if cmd.JobPostingID != nil {
		if _, err := s.GetJobPosting(ctx, cmd.TenantID, *cmd.JobPostingID); err != nil {
			return nil, err
		}
	}
	startDate, err := parseWorkerProfileDate(cmd.StartDate)
	if err != nil {
		return nil, err
	}
	dueDate, err := parseWorkerProfileDate(cmd.DueDate)
	if err != nil {
		return nil, err
	}
	return domain.NewTalentMarketplaceOpportunity(domain.TalentMarketplaceOpportunityInput{TenantID: cmd.TenantID, ProjectID: cmd.ProjectID, EngagementID: cmd.EngagementID, SourceRequirementID: cmd.SourceRequirementID, JobPostingID: cmd.JobPostingID, Title: cmd.Title, Description: cmd.Description, OpportunityType: cmd.OpportunityType, Status: cmd.Status, Visibility: cmd.Visibility, Priority: cmd.Priority, Seats: cmd.Seats, LocationMode: cmd.LocationMode, MinAllocationPercent: cmd.MinAllocationPercent, DurationLabel: cmd.DurationLabel, StartDate: startDate, DueDate: dueDate, CandidateFallbackEnabled: cmd.CandidateFallbackEnabled, CandidateFallbackStatus: cmd.CandidateFallbackStatus, Metadata: cmd.Metadata})
}

func talentMarketplaceEvent(tenantID uuid.UUID, opportunityID *uuid.UUID, applicationID *uuid.UUID, actorID *uuid.UUID, action string, fromStatus *string, toStatus *string, notes *string) *domain.TalentMarketplaceEvent {
	return &domain.TalentMarketplaceEvent{TenantID: tenantID, OpportunityID: opportunityID, ApplicationID: applicationID, ActorUserID: actorID, Action: action, FromStatus: fromStatus, ToStatus: toStatus, Notes: notes, Metadata: []byte(`{}`)}
}

func firstNonNilString(values ...*string) *string {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}
