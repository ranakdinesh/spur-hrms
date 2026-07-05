package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateTalentMarketplaceOpportunity(ctx context.Context, item *domain.TalentMarketplaceOpportunity, actorID *uuid.UUID) (*domain.TalentMarketplaceOpportunity, error) {
	row, err := s.getQueries(ctx).CreateTalentMarketplaceOpportunity(ctx, talentOpportunityCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create talent marketplace opportunity", err, tenantIDField(item.TenantID), stringField("title", item.Title))
	}
	return mapTalentMarketplaceOpportunity(row), nil
}

func (s *Store) UpdateTalentMarketplaceOpportunity(ctx context.Context, item *domain.TalentMarketplaceOpportunity, actorID *uuid.UUID) (*domain.TalentMarketplaceOpportunity, error) {
	row, err := s.getQueries(ctx).UpdateTalentMarketplaceOpportunity(ctx, talentOpportunityUpdateParams(item, actorID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTalentOpportunityNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update talent marketplace opportunity", err, tenantIDField(item.TenantID), stringField("opportunity_id", item.ID.String()))
	}
	return mapTalentMarketplaceOpportunity(row), nil
}

func (s *Store) UpdateTalentMarketplaceOpportunityFallback(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, actorID *uuid.UUID) (*domain.TalentMarketplaceOpportunity, error) {
	row, err := s.getQueries(ctx).UpdateTalentMarketplaceOpportunityFallback(ctx, sqlc.UpdateTalentMarketplaceOpportunityFallbackParams{TenantID: tenantID, ID: id, CandidateFallbackStatus: status, UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTalentOpportunityNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update talent marketplace fallback", err, tenantIDField(tenantID), stringField("opportunity_id", id.String()))
	}
	return mapTalentMarketplaceOpportunity(row), nil
}

func (s *Store) GetTalentMarketplaceOpportunity(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.TalentMarketplaceOpportunity, error) {
	row, err := s.getQueries(ctx).GetTalentMarketplaceOpportunity(ctx, sqlc.GetTalentMarketplaceOpportunityParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTalentOpportunityNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get talent marketplace opportunity", err, tenantIDField(tenantID), stringField("opportunity_id", id.String()))
	}
	return mapTalentMarketplaceOpportunity(row), nil
}

func (s *Store) ListTalentMarketplaceOpportunities(ctx context.Context, filter domain.TalentMarketplaceOpportunityFilter) ([]*domain.TalentMarketplaceOpportunity, error) {
	rows, err := s.getQueries(ctx).ListTalentMarketplaceOpportunities(ctx, sqlc.ListTalentMarketplaceOpportunitiesParams{TenantID: filter.TenantID, ProjectID: uuidFromPtr(filter.ProjectID), EngagementID: uuidFromPtr(filter.EngagementID), Status: textFromPtr(filter.Status), OpportunityType: textFromPtr(filter.OpportunityType), Search: textFromPtr(filter.Search)})
	if err != nil {
		return nil, s.logDBError(ctx, "list talent marketplace opportunities", err, tenantIDField(filter.TenantID))
	}
	return mapTalentMarketplaceOpportunityRows(rows), nil
}

func (s *Store) DeleteTalentMarketplaceOpportunity(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteTalentMarketplaceOpportunity(ctx, sqlc.SoftDeleteTalentMarketplaceOpportunityParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete talent marketplace opportunity", err, tenantIDField(tenantID), stringField("opportunity_id", id.String()))
	}
	return nil
}

func (s *Store) CreateTalentMarketplaceApplication(ctx context.Context, item *domain.TalentMarketplaceApplication, actorID *uuid.UUID) (*domain.TalentMarketplaceApplication, error) {
	row, err := s.getQueries(ctx).CreateTalentMarketplaceApplication(ctx, sqlc.CreateTalentMarketplaceApplicationParams{TenantID: item.TenantID, OpportunityID: item.OpportunityID, WorkerProfileID: item.WorkerProfileID, Status: item.Status, MatchScore: numericFromFloatPtr(item.MatchScore), MatchReasons: jsonBytesFromRaw(item.MatchReasons), WorkerNote: textFromPtr(item.WorkerNote), ManagerNote: textFromPtr(item.ManagerNote), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create talent marketplace application", err, tenantIDField(item.TenantID), stringField("opportunity_id", item.OpportunityID.String()), stringField("worker_profile_id", item.WorkerProfileID.String()))
	}
	return mapTalentMarketplaceApplication(row), nil
}

func (s *Store) UpdateTalentMarketplaceApplicationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, workerNote *string, managerNote *string, actorID *uuid.UUID) (*domain.TalentMarketplaceApplication, error) {
	row, err := s.getQueries(ctx).UpdateTalentMarketplaceApplicationStatus(ctx, sqlc.UpdateTalentMarketplaceApplicationStatusParams{TenantID: tenantID, ID: id, Status: status, WorkerNote: textFromPtr(workerNote), ManagerNote: textFromPtr(managerNote), UpdatedBy: uuidFromPtr(actorID)})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTalentApplicationNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "update talent marketplace application status", err, tenantIDField(tenantID), stringField("application_id", id.String()), stringField("status", status))
	}
	return mapTalentMarketplaceApplication(row), nil
}

func (s *Store) GetTalentMarketplaceApplication(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.TalentMarketplaceApplication, error) {
	row, err := s.getQueries(ctx).GetTalentMarketplaceApplication(ctx, sqlc.GetTalentMarketplaceApplicationParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTalentApplicationNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get talent marketplace application", err, tenantIDField(tenantID), stringField("application_id", id.String()))
	}
	return mapTalentMarketplaceApplication(row), nil
}

func (s *Store) ListTalentMarketplaceApplications(ctx context.Context, filter domain.TalentMarketplaceApplicationFilter) ([]*domain.TalentMarketplaceApplication, error) {
	rows, err := s.getQueries(ctx).ListTalentMarketplaceApplications(ctx, sqlc.ListTalentMarketplaceApplicationsParams{TenantID: filter.TenantID, OpportunityID: uuidFromPtr(filter.OpportunityID), WorkerProfileID: uuidFromPtr(filter.WorkerProfileID), Status: textFromPtr(filter.Status)})
	if err != nil {
		return nil, s.logDBError(ctx, "list talent marketplace applications", err, tenantIDField(filter.TenantID))
	}
	return mapTalentMarketplaceApplicationRows(rows), nil
}

func (s *Store) ListTalentMarketplaceRecommendations(ctx context.Context, tenantID uuid.UUID, opportunityID uuid.UUID) ([]*domain.TalentMarketplaceRecommendation, error) {
	rows, err := s.getQueries(ctx).ListTalentMarketplaceRecommendations(ctx, sqlc.ListTalentMarketplaceRecommendationsParams{TenantID: tenantID, OpportunityID: opportunityID})
	if err != nil {
		return nil, s.logDBError(ctx, "list talent marketplace recommendations", err, tenantIDField(tenantID), stringField("opportunity_id", opportunityID.String()))
	}
	return mapTalentMarketplaceRecommendations(rows), nil
}

func (s *Store) CreateTalentMarketplaceEvent(ctx context.Context, event *domain.TalentMarketplaceEvent) (*domain.TalentMarketplaceEvent, error) {
	row, err := s.getQueries(ctx).CreateTalentMarketplaceEvent(ctx, sqlc.CreateTalentMarketplaceEventParams{TenantID: event.TenantID, OpportunityID: uuidFromPtr(event.OpportunityID), ApplicationID: uuidFromPtr(event.ApplicationID), ActorUserID: uuidFromPtr(event.ActorUserID), Action: event.Action, FromStatus: textFromPtr(event.FromStatus), ToStatus: textFromPtr(event.ToStatus), Notes: textFromPtr(event.Notes), Metadata: jsonBytesFromRaw(event.Metadata)})
	if err != nil {
		return nil, s.logDBError(ctx, "create talent marketplace event", err, tenantIDField(event.TenantID), stringField("action", event.Action))
	}
	return mapTalentMarketplaceEvent(row), nil
}

func (s *Store) ListTalentMarketplaceEvents(ctx context.Context, filter domain.TalentMarketplaceEventFilter) ([]*domain.TalentMarketplaceEvent, error) {
	rows, err := s.getQueries(ctx).ListTalentMarketplaceEvents(ctx, sqlc.ListTalentMarketplaceEventsParams{TenantID: filter.TenantID, OpportunityID: uuidFromPtr(filter.OpportunityID), ApplicationID: uuidFromPtr(filter.ApplicationID)})
	if err != nil {
		return nil, s.logDBError(ctx, "list talent marketplace events", err, tenantIDField(filter.TenantID))
	}
	return mapTalentMarketplaceEventRows(rows), nil
}

func talentOpportunityCreateParams(item *domain.TalentMarketplaceOpportunity, actorID *uuid.UUID) sqlc.CreateTalentMarketplaceOpportunityParams {
	return sqlc.CreateTalentMarketplaceOpportunityParams{TenantID: item.TenantID, ProjectID: uuidFromPtr(item.ProjectID), EngagementID: uuidFromPtr(item.EngagementID), SourceRequirementID: uuidFromPtr(item.SourceRequirementID), JobPostingID: uuidFromPtr(item.JobPostingID), Title: item.Title, Description: textFromPtr(item.Description), OpportunityType: item.OpportunityType, Status: item.Status, Visibility: item.Visibility, Priority: item.Priority, Seats: item.Seats, LocationMode: item.LocationMode, MinAllocationPercent: int4FromPtr(item.MinAllocationPercent), DurationLabel: textFromPtr(item.DurationLabel), StartDate: dateFromPtr(item.StartDate), DueDate: dateFromPtr(item.DueDate), CandidateFallbackEnabled: item.CandidateFallbackEnabled, CandidateFallbackStatus: item.CandidateFallbackStatus, Metadata: jsonBytesFromRaw(item.Metadata), CreatedBy: uuidFromPtr(actorID)}
}

func talentOpportunityUpdateParams(item *domain.TalentMarketplaceOpportunity, actorID *uuid.UUID) sqlc.UpdateTalentMarketplaceOpportunityParams {
	return sqlc.UpdateTalentMarketplaceOpportunityParams{TenantID: item.TenantID, ID: item.ID, ProjectID: uuidFromPtr(item.ProjectID), EngagementID: uuidFromPtr(item.EngagementID), SourceRequirementID: uuidFromPtr(item.SourceRequirementID), JobPostingID: uuidFromPtr(item.JobPostingID), Title: item.Title, Description: textFromPtr(item.Description), OpportunityType: item.OpportunityType, Status: item.Status, Visibility: item.Visibility, Priority: item.Priority, Seats: item.Seats, LocationMode: item.LocationMode, MinAllocationPercent: int4FromPtr(item.MinAllocationPercent), DurationLabel: textFromPtr(item.DurationLabel), StartDate: dateFromPtr(item.StartDate), DueDate: dateFromPtr(item.DueDate), CandidateFallbackEnabled: item.CandidateFallbackEnabled, CandidateFallbackStatus: item.CandidateFallbackStatus, Metadata: jsonBytesFromRaw(item.Metadata), UpdatedBy: uuidFromPtr(actorID)}
}
