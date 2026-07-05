package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateCandidateApplicantAccount(ctx context.Context, item *domain.CandidateApplicantAccount, actorID *uuid.UUID) (*domain.CandidateApplicantAccount, error) {
	row, err := s.getQueries(ctx).CreateCandidateApplicantAccount(ctx, sqlc.CreateCandidateApplicantAccountParams{
		TenantID:    item.TenantID,
		CandidateID: item.CandidateID,
		UserID:      item.UserID,
		Email:       item.Email,
		Status:      textFromString(item.Status),
		ConsentAt:   timestamptzFromTimePtr(item.ConsentAt),
		ConsentIp:   textFromPtr(item.ConsentIP),
		Metadata:    jsonBytes(item.Metadata),
		CreatedBy:   uuidFromPtr(actorID),
	})
	if err != nil {
		return nil, s.logDBError(ctx, "create candidate applicant account", fmt.Errorf("hrms: create candidate applicant account: %w", err), tenantIDField(item.TenantID), stringField("candidate_id", item.CandidateID.String()), stringField("user_id", item.UserID.String()))
	}
	return mapCandidateApplicantAccount(row), nil
}

func (s *Store) GetCandidateApplicantAccountByUser(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.CandidateApplicantAccount, error) {
	row, err := s.getQueries(ctx).GetCandidateApplicantAccountByUser(ctx, sqlc.GetCandidateApplicantAccountByUserParams{TenantID: tenantID, UserID: userID})
	if err != nil {
		return nil, s.logDBError(ctx, "get candidate applicant account by user", fmt.Errorf("hrms: get candidate applicant account by user: %w", err), tenantIDField(tenantID), stringField("user_id", userID.String()))
	}
	return mapCandidateApplicantAccount(row), nil
}

func (s *Store) GetCandidateApplicantAccountByCandidate(ctx context.Context, tenantID uuid.UUID, candidateID uuid.UUID) (*domain.CandidateApplicantAccount, error) {
	row, err := s.getQueries(ctx).GetCandidateApplicantAccountByCandidate(ctx, sqlc.GetCandidateApplicantAccountByCandidateParams{TenantID: tenantID, CandidateID: candidateID})
	if err != nil {
		return nil, s.logDBError(ctx, "get candidate applicant account by candidate", fmt.Errorf("hrms: get candidate applicant account by candidate: %w", err), tenantIDField(tenantID), stringField("candidate_id", candidateID.String()))
	}
	return mapCandidateApplicantAccount(row), nil
}

func (s *Store) CreateCandidateApplication(ctx context.Context, item *domain.CandidateApplication, actorID *uuid.UUID) (*domain.CandidateApplication, error) {
	row, err := s.getQueries(ctx).CreateCandidateApplication(ctx, candidateApplicationCreateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "create candidate application", fmt.Errorf("hrms: create candidate application: %w", err), tenantIDField(item.TenantID))
	}
	return mapCandidateApplication(row), nil
}

func (s *Store) ListCandidateApplications(ctx context.Context, filter domain.CandidateApplicationFilter) ([]*domain.CandidateApplication, error) {
	rows, err := s.getQueries(ctx).ListCandidateApplications(ctx, sqlc.ListCandidateApplicationsParams{TenantID: filter.TenantID, Status: textFromPtr(filter.Status), CandidateID: uuidFromPtr(filter.CandidateID), JobPostingID: uuidFromPtr(filter.JobPostingID), Search: textFromPtr(filter.Search), OffsetRows: filter.Offset, LimitRows: filter.Limit})
	if err != nil {
		return nil, s.logDBError(ctx, "list candidate applications", err, tenantIDField(filter.TenantID))
	}
	return mapCandidateApplications(rows), nil
}

func (s *Store) CountCandidateApplications(ctx context.Context, filter domain.CandidateApplicationFilter) (int64, error) {
	count, err := s.getQueries(ctx).CountCandidateApplications(ctx, sqlc.CountCandidateApplicationsParams{TenantID: filter.TenantID, Status: textFromPtr(filter.Status), CandidateID: uuidFromPtr(filter.CandidateID), JobPostingID: uuidFromPtr(filter.JobPostingID), Search: textFromPtr(filter.Search)})
	if err != nil {
		return 0, s.logDBError(ctx, "count candidate applications", err, tenantIDField(filter.TenantID))
	}
	return count, nil
}

func (s *Store) GetCandidateApplication(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CandidateApplication, error) {
	row, err := s.getQueries(ctx).GetCandidateApplication(ctx, sqlc.GetCandidateApplicationParams{TenantID: tenantID, ID: id})
	if err != nil {
		return nil, s.logDBError(ctx, "get candidate application", fmt.Errorf("hrms: get candidate application: %w", err), tenantIDField(tenantID), stringField("candidate_application_id", id.String()))
	}
	return mapCandidateApplication(row), nil
}

func (s *Store) UpdateCandidateApplication(ctx context.Context, item *domain.CandidateApplication, actorID *uuid.UUID) (*domain.CandidateApplication, error) {
	row, err := s.getQueries(ctx).UpdateCandidateApplication(ctx, candidateApplicationUpdateParams(item, actorID))
	if err != nil {
		return nil, s.logDBError(ctx, "update candidate application", fmt.Errorf("hrms: update candidate application: %w", err), tenantIDField(item.TenantID), stringField("candidate_application_id", item.ID.String()))
	}
	return mapCandidateApplication(row), nil
}

func (s *Store) MoveCandidateApplicationStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, comments *string, reason *string, actorID *uuid.UUID) (*domain.CandidateApplication, error) {
	row, err := s.getQueries(ctx).MoveCandidateApplicationStatus(ctx, sqlc.MoveCandidateApplicationStatusParams{TenantID: tenantID, ID: id, Status: status, Comments: textFromPtr(comments), StatusChangedBy: uuidFromPtr(actorID), RejectionReason: textFromPtr(reason)})
	if err != nil {
		return nil, s.logDBError(ctx, "move candidate application status", fmt.Errorf("hrms: move candidate application status: %w", err), tenantIDField(tenantID), stringField("candidate_application_id", id.String()), stringField("status", status))
	}
	return mapCandidateApplication(row), nil
}

func (s *Store) DeleteCandidateApplication(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteCandidateApplication(ctx, sqlc.SoftDeleteCandidateApplicationParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete candidate application", fmt.Errorf("hrms: delete candidate application: %w", err), tenantIDField(tenantID), stringField("candidate_application_id", id.String()))
	}
	return nil
}

func (s *Store) CreateCandidateApplicationEvent(ctx context.Context, event *domain.CandidateApplicationEvent, actorID *uuid.UUID) (*domain.CandidateApplicationEvent, error) {
	row, err := s.getQueries(ctx).CreateCandidateApplicationEvent(ctx, sqlc.CreateCandidateApplicationEventParams{TenantID: event.TenantID, ApplicationID: event.ApplicationID, FromStatus: textFromPtr(event.FromStatus), ToStatus: event.ToStatus, Action: event.Action, Reason: textFromPtr(event.Reason), Remarks: textFromPtr(event.Remarks), CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create candidate application event", fmt.Errorf("hrms: create candidate application event: %w", err), tenantIDField(event.TenantID), stringField("candidate_application_id", event.ApplicationID.String()))
	}
	return mapCandidateApplicationEvent(row), nil
}

func (s *Store) ListCandidateApplicationEvents(ctx context.Context, tenantID uuid.UUID, applicationID uuid.UUID) ([]*domain.CandidateApplicationEvent, error) {
	rows, err := s.getQueries(ctx).ListCandidateApplicationEvents(ctx, sqlc.ListCandidateApplicationEventsParams{TenantID: tenantID, ApplicationID: applicationID})
	if err != nil {
		return nil, s.logDBError(ctx, "list candidate application events", err, tenantIDField(tenantID), stringField("candidate_application_id", applicationID.String()))
	}
	return mapCandidateApplicationEvents(rows), nil
}

func candidateApplicationCreateParams(item *domain.CandidateApplication, actorID *uuid.UUID) sqlc.CreateCandidateApplicationParams {
	return sqlc.CreateCandidateApplicationParams{TenantID: item.TenantID, CandidateID: uuidFromPtr(item.CandidateID), JobPostingID: uuidFromPtr(item.JobPostingID), ResumeUrl: textFromPtr(item.ResumeURL), CoverLetter: textFromPtr(item.CoverLetter), CurrentCtc: numericFromFloatPtr(item.CurrentCTC), ExpectedCtc: numericFromFloatPtr(item.ExpectedCTC), NoticePeriod: int4FromPtr(item.NoticePeriod), ReferredBy: textFromPtr(item.ReferredBy), Source: textFromPtr(item.Source), Status: item.Status, Comments: textFromPtr(item.Comments), CreatedBy: uuidFromPtr(actorID), Column14: timestamptzFromTime(item.AppliedAt), RejectionReason: textFromPtr(item.RejectionReason), WithdrawalReason: textFromPtr(item.WithdrawalReason), SourceDetail: textFromPtr(item.SourceDetail), DuplicateOfApplicationID: uuidFromPtr(item.DuplicateOfApplicationID)}
}

func candidateApplicationUpdateParams(item *domain.CandidateApplication, actorID *uuid.UUID) sqlc.UpdateCandidateApplicationParams {
	return sqlc.UpdateCandidateApplicationParams{TenantID: item.TenantID, ID: item.ID, CandidateID: uuidFromPtr(item.CandidateID), JobPostingID: uuidFromPtr(item.JobPostingID), ResumeUrl: textFromPtr(item.ResumeURL), CoverLetter: textFromPtr(item.CoverLetter), CurrentCtc: numericFromFloatPtr(item.CurrentCTC), ExpectedCtc: numericFromFloatPtr(item.ExpectedCTC), NoticePeriod: int4FromPtr(item.NoticePeriod), ReferredBy: textFromPtr(item.ReferredBy), Source: textFromPtr(item.Source), Status: item.Status, Comments: textFromPtr(item.Comments), UpdatedBy: uuidFromPtr(actorID), AppliedAt: timestamptzFromTime(item.AppliedAt), RejectionReason: textFromPtr(item.RejectionReason), WithdrawalReason: textFromPtr(item.WithdrawalReason), SourceDetail: textFromPtr(item.SourceDetail), DuplicateOfApplicationID: uuidFromPtr(item.DuplicateOfApplicationID)}
}

func timestamptzFromTime(value time.Time) pgtype.Timestamptz {
	if value.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: value, Valid: true}
}

func timestamptzFromTimePtr(value *time.Time) pgtype.Timestamptz {
	if value == nil || value.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func jsonBytes(value map[string]any) []byte {
	if len(value) == 0 {
		return []byte(`{}`)
	}
	data, err := json.Marshal(value)
	if err != nil {
		return []byte(`{}`)
	}
	return data
}
