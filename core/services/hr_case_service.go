package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
	"github.com/ranakdinesh/spur-hrms/core/ports"
)

func (s *TenantService) CreateHRCaseCategory(ctx context.Context, cmd ports.HRCaseCategoryCommand) (*domain.HRCaseCategory, error) {
	item, err := hrCaseCategoryFromCommand(cmd)
	if err != nil {
		s.logError("validate hr case category", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.hrCases.CreateHRCaseCategory(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create hr case category", err, serviceTenantIDField(cmd.TenantID), serviceStringField("code", item.Code))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateHRCaseCategory(ctx context.Context, cmd ports.HRCaseCategoryCommand) (*domain.HRCaseCategory, error) {
	if cmd.ID == uuid.Nil {
		s.logError("validate hr case category update", domain.ErrInvalidHRCaseCategory, serviceTenantIDField(cmd.TenantID))
		return nil, domain.ErrInvalidHRCaseCategory
	}
	item, err := hrCaseCategoryFromCommand(cmd)
	if err != nil {
		s.logError("validate hr case category update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_category_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.hrCases.UpdateHRCaseCategory(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update hr case category", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_category_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListHRCaseCategories(ctx context.Context, tenantID uuid.UUID, activeOnly *bool) ([]*domain.HRCaseCategory, error) {
	items, err := s.hrCases.ListHRCaseCategories(ctx, tenantID, activeOnly)
	if err != nil {
		s.logError("list hr case categories", err, serviceTenantIDField(tenantID))
	}
	return items, err
}

func (s *TenantService) DeleteHRCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidHRCaseCategory
	}
	err := s.hrCases.DeleteHRCaseCategory(ctx, tenantID, id, actorID)
	if err != nil {
		s.logError("delete hr case category", err, serviceTenantIDField(tenantID), serviceStringField("hr_case_category_id", id.String()))
	}
	return err
}

func (s *TenantService) CreateHRCaseSLAPolicy(ctx context.Context, cmd ports.HRCaseSLAPolicyCommand) (*domain.HRCaseSLAPolicy, error) {
	item, err := hrCaseSLAFromCommand(cmd)
	if err != nil {
		s.logError("validate hr case sla policy", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	result, err := s.hrCases.CreateHRCaseSLAPolicy(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create hr case sla policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("priority", item.Priority))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) UpdateHRCaseSLAPolicy(ctx context.Context, cmd ports.HRCaseSLAPolicyCommand) (*domain.HRCaseSLAPolicy, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidHRCaseSLA
	}
	item, err := hrCaseSLAFromCommand(cmd)
	if err != nil {
		s.logError("validate hr case sla policy update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_sla_policy_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	result, err := s.hrCases.UpdateHRCaseSLAPolicy(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update hr case sla policy", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_sla_policy_id", cmd.ID.String()))
		return nil, err
	}
	return result, nil
}

func (s *TenantService) ListHRCaseSLAPolicies(ctx context.Context, tenantID uuid.UUID, categoryID *uuid.UUID, priority *string) ([]*domain.HRCaseSLAPolicy, error) {
	items, err := s.hrCases.ListHRCaseSLAPolicies(ctx, tenantID, categoryID, priority)
	if err != nil {
		s.logError("list hr case sla policies", err, serviceTenantIDField(tenantID))
	}
	return items, err
}

func (s *TenantService) DeleteHRCaseSLAPolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error {
	if tenantID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidHRCaseSLA
	}
	err := s.hrCases.DeleteHRCaseSLAPolicy(ctx, tenantID, id, actorID)
	if err != nil {
		s.logError("delete hr case sla policy", err, serviceTenantIDField(tenantID), serviceStringField("hr_case_sla_policy_id", id.String()))
	}
	return err
}

func (s *TenantService) CreateHRCase(ctx context.Context, cmd ports.HRCaseCommand) (*domain.HRCase, error) {
	item, err := hrCaseFromCommand(cmd)
	if err != nil {
		s.logError("validate hr case", err, serviceTenantIDField(cmd.TenantID))
		return nil, err
	}
	now := time.Now().UTC()
	item.CaseNumber = domain.NewHRCaseNumber(now)
	if item.RequesterUserID == nil && cmd.ActorID != nil {
		item.RequesterUserID = cmd.ActorID
	}
	if item.CategoryID != nil {
		if category, err := s.hrCases.GetHRCaseCategory(ctx, item.TenantID, *item.CategoryID); err == nil {
			if cmd.ConfidentialityLevel == "" {
				item.ConfidentialityLevel = category.ConfidentialityDefault
			}
			if item.OwnerRole == nil {
				item.OwnerRole = category.DefaultOwnerRole
			}
		} else if !errors.Is(err, domain.ErrHRCaseCategoryNotFound) {
			s.logError("get hr case category for case create", err, serviceTenantIDField(item.TenantID), serviceStringField("hr_case_category_id", item.CategoryID.String()))
			return nil, err
		}
	}
	if policy, err := s.hrCases.ResolveHRCaseSLAPolicy(ctx, item.TenantID, item.CategoryID, item.Priority); err == nil {
		first := now.Add(time.Duration(policy.ResponseHours) * time.Hour)
		due := now.Add(time.Duration(policy.ResolutionHours) * time.Hour)
		item.FirstResponseDueAt = &first
		item.DueAt = &due
	} else if !errors.Is(err, domain.ErrHRCaseSLANotFound) {
		s.logError("resolve hr case sla policy", err, serviceTenantIDField(item.TenantID), serviceStringField("priority", item.Priority))
		return nil, err
	}
	result, err := s.hrCases.CreateHRCase(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create hr case", err, serviceTenantIDField(item.TenantID), serviceStringField("title", item.Title))
		return nil, err
	}
	s.recordHRCaseEvent(ctx, result, "created", nil, &result.Status, nil, cmd.ActorID)
	return result, nil
}

func (s *TenantService) UpdateHRCase(ctx context.Context, cmd ports.HRCaseCommand) (*domain.HRCase, error) {
	if cmd.ID == uuid.Nil {
		return nil, domain.ErrInvalidHRCase
	}
	item, err := hrCaseFromCommand(cmd)
	if err != nil {
		s.logError("validate hr case update", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_id", cmd.ID.String()))
		return nil, err
	}
	item.ID = cmd.ID
	existing, _ := s.hrCases.GetHRCase(ctx, cmd.TenantID, cmd.ID)
	if existing != nil && item.DueAt == nil {
		item.DueAt = existing.DueAt
	}
	result, err := s.hrCases.UpdateHRCaseDetails(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("update hr case", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_id", cmd.ID.String()))
		return nil, err
	}
	s.recordHRCaseEvent(ctx, result, "updated", nil, &result.Status, nil, cmd.ActorID)
	return result, nil
}

func (s *TenantService) ListHRCases(ctx context.Context, filter domain.HRCaseFilter) (*domain.HRCasePage, error) {
	items, err := s.hrCases.ListHRCases(ctx, filter)
	if err != nil {
		s.logError("list hr cases", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	total, err := s.hrCases.CountHRCases(ctx, filter)
	if err != nil {
		s.logError("count hr cases", err, serviceTenantIDField(filter.TenantID))
		return nil, err
	}
	summary, _ := s.hrCases.GetHRCaseSummary(ctx, filter.TenantID)
	categories, _ := s.hrCases.ListHRCaseCategories(ctx, filter.TenantID, nil)
	return &domain.HRCasePage{Items: items, Total: total, Summary: summary, Categories: categories}, nil
}

func (s *TenantService) GetHRCaseWorkspace(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, includeInternal bool) (*domain.HRCaseWorkspace, error) {
	item, err := s.hrCases.GetHRCase(ctx, tenantID, id)
	if err != nil {
		s.logError("get hr case workspace case", err, serviceTenantIDField(tenantID), serviceStringField("hr_case_id", id.String()))
		return nil, err
	}
	comments, err := s.hrCases.ListHRCaseComments(ctx, tenantID, id, includeInternal)
	if err != nil {
		return nil, err
	}
	attachments, err := s.hrCases.ListHRCaseAttachments(ctx, tenantID, id, includeInternal)
	if err != nil {
		return nil, err
	}
	events, err := s.hrCases.ListHRCaseEvents(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	return &domain.HRCaseWorkspace{Case: item, Comments: comments, Attachments: attachments, Events: events}, nil
}

func (s *TenantService) UpdateHRCaseStatus(ctx context.Context, cmd ports.HRCaseStatusCommand) (*domain.HRCase, error) {
	status := cleanHRCaseStatus(cmd.Status)
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || status == "" {
		return nil, domain.ErrInvalidHRCase
	}
	before, _ := s.hrCases.GetHRCase(ctx, cmd.TenantID, cmd.ID)
	result, err := s.hrCases.UpdateHRCaseStatus(ctx, cmd.TenantID, cmd.ID, status, cmd.ResolutionSummary, cmd.ActorID)
	if err != nil {
		s.logError("update hr case status", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_id", cmd.ID.String()), serviceStringField("status", status))
		return nil, err
	}
	var from *string
	if before != nil {
		from = &before.Status
	}
	s.recordHRCaseEvent(ctx, result, "status_changed", from, &result.Status, cmd.Comment, cmd.ActorID)
	if cmd.Comment != nil && strings.TrimSpace(*cmd.Comment) != "" {
		_, _ = s.CreateHRCaseComment(ctx, ports.HRCaseCommentCommand{TenantID: cmd.TenantID, CaseID: cmd.ID, Visibility: domain.HRCaseVisibilityInternal, Body: *cmd.Comment, ActorID: cmd.ActorID})
	}
	return result, nil
}

func (s *TenantService) AssignHRCase(ctx context.Context, cmd ports.HRCaseAssignmentCommand) (*domain.HRCase, error) {
	if cmd.TenantID == uuid.Nil || cmd.ID == uuid.Nil || (cmd.OwnerUserID == nil && cmd.OwnerRole == nil) {
		return nil, domain.ErrInvalidHRCase
	}
	result, err := s.hrCases.UpdateHRCaseAssignment(ctx, cmd.TenantID, cmd.ID, cmd.OwnerUserID, cmd.OwnerRole, cmd.ActorID)
	if err != nil {
		s.logError("assign hr case", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_id", cmd.ID.String()))
		return nil, err
	}
	s.recordHRCaseEvent(ctx, result, "assigned", nil, &result.Status, cmd.Comment, cmd.ActorID)
	if cmd.Comment != nil && strings.TrimSpace(*cmd.Comment) != "" {
		_, _ = s.CreateHRCaseComment(ctx, ports.HRCaseCommentCommand{TenantID: cmd.TenantID, CaseID: cmd.ID, Visibility: domain.HRCaseVisibilityInternal, Body: *cmd.Comment, ActorID: cmd.ActorID})
	}
	return result, nil
}

func (s *TenantService) CreateHRCaseComment(ctx context.Context, cmd ports.HRCaseCommentCommand) (*domain.HRCaseComment, error) {
	item, err := domain.NewHRCaseComment(domain.HRCaseComment{TenantID: cmd.TenantID, CaseID: cmd.CaseID, AuthorUserID: cmd.ActorID, Visibility: cmd.Visibility, Body: cmd.Body})
	if err != nil {
		s.logError("validate hr case comment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_id", cmd.CaseID.String()))
		return nil, err
	}
	result, err := s.hrCases.CreateHRCaseComment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create hr case comment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_id", cmd.CaseID.String()))
		return nil, err
	}
	caseRef := &domain.HRCase{TenantID: cmd.TenantID, ID: cmd.CaseID}
	s.recordHRCaseEvent(ctx, caseRef, "comment_added", nil, nil, nil, cmd.ActorID)
	return result, nil
}

func (s *TenantService) CreateHRCaseAttachment(ctx context.Context, cmd ports.HRCaseAttachmentCommand) (*domain.HRCaseAttachment, error) {
	if s.hrCaseAttachmentStorage == nil {
		return nil, domain.ErrStorageProviderSettingsNotFound
	}
	content, err := base64.StdEncoding.DecodeString(cmd.FileContentBase64)
	if err != nil || len(content) == 0 {
		s.logError("decode hr case attachment", domain.ErrInvalidHRCaseAttachment, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_id", cmd.CaseID.String()))
		return nil, domain.ErrInvalidHRCaseAttachment
	}
	commentID := uuid.Nil
	if cmd.CommentID != nil {
		commentID = *cmd.CommentID
	}
	objectKey, err := s.hrCaseAttachmentStorage.StoreHRCaseAttachment(ctx, ports.StoreHRCaseAttachmentInput{TenantID: cmd.TenantID, CaseID: cmd.CaseID, CommentID: commentID, FileName: cmd.FileName, ContentType: cmd.ContentType, Content: content})
	if err != nil {
		s.logError("store hr case attachment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_id", cmd.CaseID.String()))
		return nil, err
	}
	item, err := domain.NewHRCaseAttachment(domain.HRCaseAttachment{TenantID: cmd.TenantID, CaseID: cmd.CaseID, CommentID: cmd.CommentID, FileName: cmd.FileName, ContentType: cmd.ContentType, ObjectKey: objectKey, Visibility: cmd.Visibility, UploadedBy: cmd.ActorID})
	if err != nil {
		return nil, err
	}
	result, err := s.hrCases.CreateHRCaseAttachment(ctx, item, cmd.ActorID)
	if err != nil {
		s.logError("create hr case attachment", err, serviceTenantIDField(cmd.TenantID), serviceStringField("hr_case_id", cmd.CaseID.String()))
		return nil, err
	}
	caseRef := &domain.HRCase{TenantID: cmd.TenantID, ID: cmd.CaseID}
	s.recordHRCaseEvent(ctx, caseRef, "attachment_added", nil, nil, nil, cmd.ActorID)
	return result, nil
}

func (s *TenantService) recordHRCaseEvent(ctx context.Context, item *domain.HRCase, eventType string, fromStatus *string, toStatus *string, comment *string, actorID *uuid.UUID) {
	if item == nil || item.TenantID == uuid.Nil || item.ID == uuid.Nil {
		return
	}
	_, _ = s.hrCases.CreateHRCaseEvent(ctx, &domain.HRCaseEvent{TenantID: item.TenantID, CaseID: item.ID, EventType: eventType, FromStatus: fromStatus, ToStatus: toStatus, ActorUserID: actorID, Comment: comment, Metadata: json.RawMessage(`{}`)}, actorID)
}

func hrCaseCategoryFromCommand(cmd ports.HRCaseCategoryCommand) (*domain.HRCaseCategory, error) {
	return domain.NewHRCaseCategory(domain.HRCaseCategory{TenantID: cmd.TenantID, Code: cmd.Code, Name: cmd.Name, Description: cmd.Description, ConfidentialityDefault: cmd.ConfidentialityDefault, DefaultOwnerRole: cmd.DefaultOwnerRole, IsActive: cmd.IsActive})
}

func hrCaseSLAFromCommand(cmd ports.HRCaseSLAPolicyCommand) (*domain.HRCaseSLAPolicy, error) {
	if cmd.ResponseHours == 0 {
		cmd.ResponseHours = 8
	}
	if cmd.ResolutionHours == 0 {
		cmd.ResolutionHours = 48
	}
	if cmd.EscalationHours == 0 {
		cmd.EscalationHours = 24
	}
	return domain.NewHRCaseSLAPolicy(domain.HRCaseSLAPolicy{TenantID: cmd.TenantID, CategoryID: cmd.CategoryID, Priority: cmd.Priority, ResponseHours: cmd.ResponseHours, ResolutionHours: cmd.ResolutionHours, EscalationHours: cmd.EscalationHours, IsActive: cmd.IsActive})
}

func hrCaseFromCommand(cmd ports.HRCaseCommand) (*domain.HRCase, error) {
	return domain.NewHRCase(domain.HRCase{TenantID: cmd.TenantID, CategoryID: cmd.CategoryID, CaseType: cmd.CaseType, Title: cmd.Title, Description: cmd.Description, ConfidentialityLevel: cmd.ConfidentialityLevel, RequesterUserID: cmd.RequesterUserID, SubjectEmployeeUserID: cmd.SubjectEmployeeUserID, OwnerUserID: cmd.OwnerUserID, OwnerRole: cmd.OwnerRole, Status: cmd.Status, Priority: cmd.Priority, SourceChannel: cmd.SourceChannel})
}

func cleanHRCaseStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case domain.HRCaseStatusNew, domain.HRCaseStatusOpen, domain.HRCaseStatusInProgress, domain.HRCaseStatusWaitingOnEmployee, domain.HRCaseStatusWaitingOnHR, domain.HRCaseStatusEscalated, domain.HRCaseStatusResolved, domain.HRCaseStatusClosed, domain.HRCaseStatusCancelled:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ""
	}
}
