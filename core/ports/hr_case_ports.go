package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type HRCaseRepo interface {
	CreateHRCaseCategory(ctx context.Context, item *domain.HRCaseCategory, actorID *uuid.UUID) (*domain.HRCaseCategory, error)
	UpdateHRCaseCategory(ctx context.Context, item *domain.HRCaseCategory, actorID *uuid.UUID) (*domain.HRCaseCategory, error)
	ListHRCaseCategories(ctx context.Context, tenantID uuid.UUID, activeOnly *bool) ([]*domain.HRCaseCategory, error)
	GetHRCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.HRCaseCategory, error)
	DeleteHRCaseCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateHRCaseSLAPolicy(ctx context.Context, item *domain.HRCaseSLAPolicy, actorID *uuid.UUID) (*domain.HRCaseSLAPolicy, error)
	UpdateHRCaseSLAPolicy(ctx context.Context, item *domain.HRCaseSLAPolicy, actorID *uuid.UUID) (*domain.HRCaseSLAPolicy, error)
	ListHRCaseSLAPolicies(ctx context.Context, tenantID uuid.UUID, categoryID *uuid.UUID, priority *string) ([]*domain.HRCaseSLAPolicy, error)
	ResolveHRCaseSLAPolicy(ctx context.Context, tenantID uuid.UUID, categoryID *uuid.UUID, priority string) (*domain.HRCaseSLAPolicy, error)
	DeleteHRCaseSLAPolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateHRCase(ctx context.Context, item *domain.HRCase, actorID *uuid.UUID) (*domain.HRCase, error)
	ListHRCases(ctx context.Context, filter domain.HRCaseFilter) ([]*domain.HRCase, error)
	CountHRCases(ctx context.Context, filter domain.HRCaseFilter) (int64, error)
	GetHRCase(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.HRCase, error)
	UpdateHRCaseDetails(ctx context.Context, item *domain.HRCase, actorID *uuid.UUID) (*domain.HRCase, error)
	UpdateHRCaseAssignment(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, ownerUserID *uuid.UUID, ownerRole *string, actorID *uuid.UUID) (*domain.HRCase, error)
	UpdateHRCaseStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, resolutionSummary *string, actorID *uuid.UUID) (*domain.HRCase, error)
	CreateHRCaseComment(ctx context.Context, item *domain.HRCaseComment, actorID *uuid.UUID) (*domain.HRCaseComment, error)
	ListHRCaseComments(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID, includeInternal bool) ([]*domain.HRCaseComment, error)
	CreateHRCaseAttachment(ctx context.Context, item *domain.HRCaseAttachment, actorID *uuid.UUID) (*domain.HRCaseAttachment, error)
	ListHRCaseAttachments(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID, includeInternal bool) ([]*domain.HRCaseAttachment, error)
	CreateHRCaseEvent(ctx context.Context, item *domain.HRCaseEvent, actorID *uuid.UUID) (*domain.HRCaseEvent, error)
	ListHRCaseEvents(ctx context.Context, tenantID uuid.UUID, caseID uuid.UUID) ([]*domain.HRCaseEvent, error)
	GetHRCaseSummary(ctx context.Context, tenantID uuid.UUID) ([]*domain.HRCaseSummaryRow, error)
}

type HRCaseAttachmentStorage interface {
	StoreHRCaseAttachment(ctx context.Context, input StoreHRCaseAttachmentInput) (string, error)
}

type StoreHRCaseAttachmentInput struct {
	TenantID    uuid.UUID
	CaseID      uuid.UUID
	CommentID   uuid.UUID
	FileName    string
	ContentType string
	Content     []byte
}

type HRCaseCategoryCommand struct {
	ID                     uuid.UUID  `json:"id,omitempty"`
	TenantID               uuid.UUID  `json:"tenant_id"`
	Code                   string     `json:"code"`
	Name                   string     `json:"name"`
	Description            *string    `json:"description,omitempty"`
	ConfidentialityDefault string     `json:"confidentiality_default"`
	DefaultOwnerRole       *string    `json:"default_owner_role,omitempty"`
	IsActive               bool       `json:"is_active"`
	ActorID                *uuid.UUID `json:"-"`
}

type HRCaseSLAPolicyCommand struct {
	ID              uuid.UUID  `json:"id,omitempty"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty"`
	Priority        string     `json:"priority"`
	ResponseHours   int32      `json:"response_hours"`
	ResolutionHours int32      `json:"resolution_hours"`
	EscalationHours int32      `json:"escalation_hours"`
	IsActive        bool       `json:"is_active"`
	ActorID         *uuid.UUID `json:"-"`
}

type HRCaseCommand struct {
	ID                    uuid.UUID  `json:"id,omitempty"`
	TenantID              uuid.UUID  `json:"tenant_id"`
	CategoryID            *uuid.UUID `json:"category_id,omitempty"`
	CaseType              string     `json:"case_type"`
	Title                 string     `json:"title"`
	Description           string     `json:"description"`
	ConfidentialityLevel  string     `json:"confidentiality_level"`
	RequesterUserID       *uuid.UUID `json:"requester_user_id,omitempty"`
	SubjectEmployeeUserID *uuid.UUID `json:"subject_employee_user_id,omitempty"`
	OwnerUserID           *uuid.UUID `json:"owner_user_id,omitempty"`
	OwnerRole             *string    `json:"owner_role,omitempty"`
	Status                string     `json:"status"`
	Priority              string     `json:"priority"`
	SourceChannel         string     `json:"source_channel"`
	ActorID               *uuid.UUID `json:"-"`
}

type HRCaseStatusCommand struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	ID                uuid.UUID  `json:"id"`
	Status            string     `json:"status"`
	ResolutionSummary *string    `json:"resolution_summary,omitempty"`
	Comment           *string    `json:"comment,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}

type HRCaseAssignmentCommand struct {
	TenantID    uuid.UUID  `json:"tenant_id"`
	ID          uuid.UUID  `json:"id"`
	OwnerUserID *uuid.UUID `json:"owner_user_id,omitempty"`
	OwnerRole   *string    `json:"owner_role,omitempty"`
	Comment     *string    `json:"comment,omitempty"`
	ActorID     *uuid.UUID `json:"-"`
}

type HRCaseCommentCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	CaseID     uuid.UUID  `json:"case_id"`
	Visibility string     `json:"visibility"`
	Body       string     `json:"body"`
	ActorID    *uuid.UUID `json:"-"`
}

type HRCaseAttachmentCommand struct {
	TenantID          uuid.UUID  `json:"tenant_id"`
	CaseID            uuid.UUID  `json:"case_id"`
	CommentID         *uuid.UUID `json:"comment_id,omitempty"`
	FileName          string     `json:"file_name"`
	ContentType       string     `json:"content_type"`
	FileContentBase64 string     `json:"file_content_base64"`
	Visibility        string     `json:"visibility"`
	ActorID           *uuid.UUID `json:"-"`
}
