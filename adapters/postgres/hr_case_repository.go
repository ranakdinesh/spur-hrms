package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranakdinesh/spur-hrms/adapters/postgres/sqlc"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

func (s *Store) CreateHRCaseCategory(ctx context.Context, item *domain.HRCaseCategory, actorID *uuid.UUID) (*domain.HRCaseCategory, error) {
	row, err := s.getQueries(ctx).CreateHRCaseCategory(ctx, sqlc.CreateHRCaseCategoryParams{TenantID: item.TenantID, Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), ConfidentialityDefault: item.ConfidentialityDefault, DefaultOwnerRole: textFromPtr(item.DefaultOwnerRole), IsActive: item.IsActive, CreatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create hr case category", err, tenantIDField(item.TenantID), stringField("code", item.Code))
	}
	return mapHRCaseCategory(row), nil
}

func (s *Store) UpdateHRCaseCategory(ctx context.Context, item *domain.HRCaseCategory, actorID *uuid.UUID) (*domain.HRCaseCategory, error) {
	row, err := s.getQueries(ctx).UpdateHRCaseCategory(ctx, sqlc.UpdateHRCaseCategoryParams{TenantID: item.TenantID, ID: item.ID, Code: item.Code, Name: item.Name, Description: textFromPtr(item.Description), ConfidentialityDefault: item.ConfidentialityDefault, DefaultOwnerRole: textFromPtr(item.DefaultOwnerRole), IsActive: item.IsActive, UpdatedBy: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update hr case category", err, tenantIDField(item.TenantID), stringField("hr_case_category_id", item.ID.String()))
	}
	return mapHRCaseCategory(row), nil
}

func (s *Store) ListHRCaseCategories(ctx context.Context, tenantID uuid.UUID, activeOnly *bool) ([]*domain.HRCaseCategory, error) {
	rows, err := s.getQueries(ctx).ListHRCaseCategories(ctx, sqlc.ListHRCaseCategoriesParams{TenantID: tenantID, IsActive: boolFromPtr(activeOnly)})
	if err != nil {
		return nil, s.logDBError(ctx, "list hr case categories", err, tenantIDField(tenantID))
	}
	return mapHRCaseCategories(rows), nil
}

func (s *Store) GetHRCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.HRCaseCategory, error) {
	row, err := s.getQueries(ctx).GetHRCaseCategory(ctx, sqlc.GetHRCaseCategoryParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrHRCaseCategoryNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get hr case category", err, tenantIDField(tenantID), stringField("hr_case_category_id", id.String()))
	}
	return mapHRCaseCategory(row), nil
}

func (s *Store) DeleteHRCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteHRCaseCategory(ctx, sqlc.SoftDeleteHRCaseCategoryParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete hr case category", err, tenantIDField(tenantID), stringField("hr_case_category_id", id.String()))
	}
	return nil
}

func (s *Store) CreateHRCaseSLAPolicy(ctx context.Context, item *domain.HRCaseSLAPolicy, actorID *uuid.UUID) (*domain.HRCaseSLAPolicy, error) {
	row, err := s.getQueries(ctx).CreateHRCaseSLAPolicy(ctx, sqlc.CreateHRCaseSLAPolicyParams{TenantID: item.TenantID, CategoryID: uuidFromPtr(item.CategoryID), Priority: item.Priority, ResponseHours: item.ResponseHours, ResolutionHours: item.ResolutionHours, EscalationHours: item.EscalationHours, IsActive: item.IsActive, ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create hr case sla policy", err, tenantIDField(item.TenantID), stringField("priority", item.Priority))
	}
	return mapHRCaseSLAPolicy(row), nil
}

func (s *Store) UpdateHRCaseSLAPolicy(ctx context.Context, item *domain.HRCaseSLAPolicy, actorID *uuid.UUID) (*domain.HRCaseSLAPolicy, error) {
	row, err := s.getQueries(ctx).UpdateHRCaseSLAPolicy(ctx, sqlc.UpdateHRCaseSLAPolicyParams{TenantID: item.TenantID, ID: item.ID, CategoryID: uuidFromPtr(item.CategoryID), Priority: item.Priority, ResponseHours: item.ResponseHours, ResolutionHours: item.ResolutionHours, EscalationHours: item.EscalationHours, IsActive: item.IsActive, ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update hr case sla policy", err, tenantIDField(item.TenantID), stringField("hr_case_sla_policy_id", item.ID.String()))
	}
	return mapHRCaseSLAPolicy(row), nil
}

func (s *Store) ListHRCaseSLAPolicies(ctx context.Context, tenantID uuid.UUID, categoryID *uuid.UUID, priority *string) ([]*domain.HRCaseSLAPolicy, error) {
	rows, err := s.getQueries(ctx).ListHRCaseSLAPolicies(ctx, sqlc.ListHRCaseSLAPoliciesParams{TenantID: tenantID, CategoryID: uuidFromPtr(categoryID), Priority: textFromPtr(priority)})
	if err != nil {
		return nil, s.logDBError(ctx, "list hr case sla policies", err, tenantIDField(tenantID))
	}
	return mapHRCaseSLAPolicyList(rows), nil
}

func (s *Store) ResolveHRCaseSLAPolicy(ctx context.Context, tenantID uuid.UUID, categoryID *uuid.UUID, priority string) (*domain.HRCaseSLAPolicy, error) {
	row, err := s.getQueries(ctx).ResolveHRCaseSLAPolicy(ctx, sqlc.ResolveHRCaseSLAPolicyParams{TenantID: tenantID, CategoryID: uuidFromPtr(categoryID), Priority: priority})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrHRCaseSLANotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "resolve hr case sla policy", err, tenantIDField(tenantID), stringField("priority", priority))
	}
	return mapHRCaseSLAPolicy(row), nil
}

func (s *Store) DeleteHRCaseSLAPolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if err := s.getQueries(ctx).SoftDeleteHRCaseSLAPolicy(ctx, sqlc.SoftDeleteHRCaseSLAPolicyParams{TenantID: tenantID, ID: id, UpdatedBy: uuidFromPtr(actorID)}); err != nil {
		return s.logDBError(ctx, "delete hr case sla policy", err, tenantIDField(tenantID), stringField("hr_case_sla_policy_id", id.String()))
	}
	return nil
}

func (s *Store) CreateHRCase(ctx context.Context, item *domain.HRCase, actorID *uuid.UUID) (*domain.HRCase, error) {
	row, err := s.getQueries(ctx).CreateHRCase(ctx, sqlc.CreateHRCaseParams{TenantID: item.TenantID, CaseNumber: item.CaseNumber, CategoryID: uuidFromPtr(item.CategoryID), CaseType: item.CaseType, Title: item.Title, Description: item.Description, ConfidentialityLevel: item.ConfidentialityLevel, RequesterUserID: uuidFromPtr(item.RequesterUserID), SubjectEmployeeUserID: uuidFromPtr(item.SubjectEmployeeUserID), OwnerUserID: uuidFromPtr(item.OwnerUserID), OwnerRole: textFromPtr(item.OwnerRole), Status: item.Status, Priority: item.Priority, SourceChannel: item.SourceChannel, FirstResponseDueAt: timestamptzFromPtr(item.FirstResponseDueAt), DueAt: timestamptzFromPtr(item.DueAt), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create hr case", err, tenantIDField(item.TenantID), stringField("case_number", item.CaseNumber))
	}
	return mapHRCase(row), nil
}

func (s *Store) ListHRCases(ctx context.Context, filter domain.HRCaseFilter) ([]*domain.HRCase, error) {
	params := hrCaseFilterParams(filter)
	rows, err := s.getQueries(ctx).ListHRCases(ctx, sqlc.ListHRCasesParams(params))
	if err != nil {
		return nil, s.logDBError(ctx, "list hr cases", err, tenantIDField(filter.TenantID))
	}
	return mapHRCaseList(rows), nil
}

func (s *Store) CountHRCases(ctx context.Context, filter domain.HRCaseFilter) (int64, error) {
	params := hrCaseFilterParams(filter)
	total, err := s.getQueries(ctx).CountHRCases(ctx, sqlc.CountHRCasesParams{TenantID: params.TenantID, Status: params.Status, Priority: params.Priority, CategoryID: params.CategoryID, RequesterUserID: params.RequesterUserID, SubjectEmployeeUserID: params.SubjectEmployeeUserID, OwnerUserID: params.OwnerUserID, ConfidentialityLevel: params.ConfidentialityLevel, Search: params.Search})
	if err != nil {
		return 0, s.logDBError(ctx, "count hr cases", err, tenantIDField(filter.TenantID))
	}
	return total, nil
}

func (s *Store) GetHRCase(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.HRCase, error) {
	row, err := s.getQueries(ctx).GetHRCase(ctx, sqlc.GetHRCaseParams{TenantID: tenantID, ID: id})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrHRCaseNotFound
	}
	if err != nil {
		return nil, s.logDBError(ctx, "get hr case", err, tenantIDField(tenantID), stringField("hr_case_id", id.String()))
	}
	return mapHRCaseGetRow(row), nil
}

func (s *Store) UpdateHRCaseDetails(ctx context.Context, item *domain.HRCase, actorID *uuid.UUID) (*domain.HRCase, error) {
	row, err := s.getQueries(ctx).UpdateHRCaseDetails(ctx, sqlc.UpdateHRCaseDetailsParams{TenantID: item.TenantID, ID: item.ID, CaseType: item.CaseType, Title: item.Title, Description: item.Description, ConfidentialityLevel: item.ConfidentialityLevel, Priority: item.Priority, CategoryID: uuidFromPtr(item.CategoryID), SubjectEmployeeUserID: uuidFromPtr(item.SubjectEmployeeUserID), DueAt: timestamptzFromPtr(item.DueAt), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update hr case details", err, tenantIDField(item.TenantID), stringField("hr_case_id", item.ID.String()))
	}
	return mapHRCase(row), nil
}

func (s *Store) UpdateHRCaseAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, ownerUserID *uuid.UUID, ownerRole *string, actorID *uuid.UUID) (*domain.HRCase, error) {
	row, err := s.getQueries(ctx).UpdateHRCaseAssignment(ctx, sqlc.UpdateHRCaseAssignmentParams{TenantID: tenantID, ID: id, OwnerUserID: uuidFromPtr(ownerUserID), OwnerRole: textFromPtr(ownerRole), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update hr case assignment", err, tenantIDField(tenantID), stringField("hr_case_id", id.String()))
	}
	return mapHRCase(row), nil
}

func (s *Store) UpdateHRCaseStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, resolutionSummary *string, actorID *uuid.UUID) (*domain.HRCase, error) {
	row, err := s.getQueries(ctx).UpdateHRCaseStatus(ctx, sqlc.UpdateHRCaseStatusParams{TenantID: tenantID, ID: id, Status: status, ResolutionSummary: textFromPtr(resolutionSummary), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "update hr case status", err, tenantIDField(tenantID), stringField("hr_case_id", id.String()), stringField("status", status))
	}
	return mapHRCase(row), nil
}

func (s *Store) CreateHRCaseComment(ctx context.Context, item *domain.HRCaseComment, actorID *uuid.UUID) (*domain.HRCaseComment, error) {
	row, err := s.getQueries(ctx).CreateHRCaseComment(ctx, sqlc.CreateHRCaseCommentParams{TenantID: item.TenantID, CaseID: item.CaseID, AuthorUserID: uuidFromPtr(item.AuthorUserID), Visibility: item.Visibility, Body: item.Body, ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create hr case comment", err, tenantIDField(item.TenantID), stringField("hr_case_id", item.CaseID.String()))
	}
	return mapHRCaseComment(row), nil
}

func (s *Store) ListHRCaseComments(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID, includeInternal bool) ([]*domain.HRCaseComment, error) {
	rows, err := s.getQueries(ctx).ListHRCaseComments(ctx, sqlc.ListHRCaseCommentsParams{TenantID: tenantID, CaseID: caseID, IncludeInternal: boolFromPtr(&includeInternal)})
	if err != nil {
		return nil, s.logDBError(ctx, "list hr case comments", err, tenantIDField(tenantID), stringField("hr_case_id", caseID.String()))
	}
	return mapHRCaseComments(rows), nil
}

func (s *Store) CreateHRCaseAttachment(ctx context.Context, item *domain.HRCaseAttachment, actorID *uuid.UUID) (*domain.HRCaseAttachment, error) {
	row, err := s.getQueries(ctx).CreateHRCaseAttachment(ctx, sqlc.CreateHRCaseAttachmentParams{TenantID: item.TenantID, CaseID: item.CaseID, CommentID: uuidFromPtr(item.CommentID), FileName: item.FileName, ContentType: item.ContentType, ObjectKey: item.ObjectKey, Visibility: item.Visibility, UploadedBy: uuidFromPtr(item.UploadedBy), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create hr case attachment", err, tenantIDField(item.TenantID), stringField("hr_case_id", item.CaseID.String()))
	}
	return mapHRCaseAttachment(row), nil
}

func (s *Store) ListHRCaseAttachments(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID, includeInternal bool) ([]*domain.HRCaseAttachment, error) {
	rows, err := s.getQueries(ctx).ListHRCaseAttachments(ctx, sqlc.ListHRCaseAttachmentsParams{TenantID: tenantID, CaseID: caseID, IncludeInternal: boolFromPtr(&includeInternal)})
	if err != nil {
		return nil, s.logDBError(ctx, "list hr case attachments", err, tenantIDField(tenantID), stringField("hr_case_id", caseID.String()))
	}
	return mapHRCaseAttachments(rows), nil
}

func (s *Store) CreateHRCaseEvent(ctx context.Context, item *domain.HRCaseEvent, actorID *uuid.UUID) (*domain.HRCaseEvent, error) {
	row, err := s.getQueries(ctx).CreateHRCaseEvent(ctx, sqlc.CreateHRCaseEventParams{TenantID: item.TenantID, CaseID: item.CaseID, EventType: item.EventType, FromStatus: textFromPtr(item.FromStatus), ToStatus: textFromPtr(item.ToStatus), ActorUserID: uuidFromPtr(item.ActorUserID), Comment: textFromPtr(item.Comment), Metadata: jsonBytesFromRaw(item.Metadata), ActorID: uuidFromPtr(actorID)})
	if err != nil {
		return nil, s.logDBError(ctx, "create hr case event", err, tenantIDField(item.TenantID), stringField("hr_case_id", item.CaseID.String()), stringField("event_type", item.EventType))
	}
	return mapHRCaseEvent(row), nil
}

func (s *Store) ListHRCaseEvents(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.HRCaseEvent, error) {
	rows, err := s.getQueries(ctx).ListHRCaseEvents(ctx, sqlc.ListHRCaseEventsParams{TenantID: tenantID, CaseID: caseID})
	if err != nil {
		return nil, s.logDBError(ctx, "list hr case events", err, tenantIDField(tenantID), stringField("hr_case_id", caseID.String()))
	}
	return mapHRCaseEvents(rows), nil
}

func (s *Store) GetHRCaseSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.HRCaseSummaryRow, error) {
	rows, err := s.getQueries(ctx).GetHRCaseSummary(ctx, tenantID)
	if err != nil {
		return nil, s.logDBError(ctx, "get hr case summary", err, tenantIDField(tenantID))
	}
	return mapHRCaseSummary(rows), nil
}

func hrCaseFilterParams(filter domain.HRCaseFilter) sqlc.ListHRCasesParams {
	limit := filter.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	return sqlc.ListHRCasesParams{TenantID: filter.TenantID, Status: textFromPtr(filter.Status), Priority: textFromPtr(filter.Priority), CategoryID: uuidFromPtr(filter.CategoryID), RequesterUserID: uuidFromPtr(filter.RequesterUserID), SubjectEmployeeUserID: uuidFromPtr(filter.SubjectEmployeeUserID), OwnerUserID: uuidFromPtr(filter.OwnerUserID), ConfidentialityLevel: textFromPtr(filter.ConfidentialityLevel), Search: textFromPtr(filter.Search), Limit: limit, Offset: offset}
}
