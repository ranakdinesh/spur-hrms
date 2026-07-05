package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type EmployeeDocumentRepo interface {
	CreateDocumentType(ctx context.Context, item *domain.DocumentType, actorID *uuid.UUID) (*domain.DocumentType, error)
	ListDocumentTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.DocumentType, error)
	GetDocumentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.DocumentType, error)
	UpdateDocumentType(ctx context.Context, item *domain.DocumentType, actorID *uuid.UUID) (*domain.DocumentType, error)
	DeleteDocumentType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateEmployeeDocument(ctx context.Context, item *domain.EmployeeDocument, actorID *uuid.UUID) (*domain.EmployeeDocument, error)
	GetEmployeeDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EmployeeDocument, error)
	ReviewEmployeeDocument(ctx context.Context, tenantID uuid.UUID, documentID uuid.UUID, status string, remarks *string, actorID *uuid.UUID) (*domain.EmployeeDocument, error)
	UpdateEmployeeDocument(ctx context.Context, item *domain.EmployeeDocument, actorID *uuid.UUID) (*domain.EmployeeDocument, error)
	DeleteEmployeeDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type EmployeeDocumentStorage interface {
	StoreEmployeeDocument(ctx context.Context, input StoreEmployeeDocumentInput) (string, error)
}

type StoreEmployeeDocumentInput struct {
	TenantID    uuid.UUID
	EmployeeID  uuid.UUID
	DocumentID  uuid.UUID
	FileName    string
	ContentType string
	Content     []byte
}

type DocumentTypeCommand struct {
	ID                  uuid.UUID  `json:"id,omitempty"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	Name                string     `json:"name"`
	Description         *string    `json:"description,omitempty"`
	IsRequired          *bool      `json:"is_required,omitempty"`
	Instructions        *string    `json:"instructions,omitempty"`
	AllowedContentTypes string     `json:"allowed_content_types,omitempty"`
	MaxFileSizeBytes    int64      `json:"max_file_size_bytes,omitempty"`
	DisplayOrder        int32      `json:"display_order"`
	ActorID             *uuid.UUID `json:"-"`
}

type EmployeeDocumentCommand struct {
	ID                uuid.UUID  `json:"id,omitempty"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	EmployeeID        uuid.UUID  `json:"employee_id"`
	DocumentTypeID    *uuid.UUID `json:"document_type_id,omitempty"`
	Title             *string    `json:"title,omitempty"`
	FilePath          *string    `json:"file_path,omitempty"`
	FileName          string     `json:"file_name,omitempty"`
	FileContentType   string     `json:"file_content_type,omitempty"`
	FileContentBase64 string     `json:"file_content_base64,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}

type EmployeeDocumentReviewCommand struct {
	TenantID   uuid.UUID  `json:"tenant_id"`
	EmployeeID uuid.UUID  `json:"employee_id"`
	DocumentID uuid.UUID  `json:"document_id"`
	Status     string     `json:"status"`
	Remarks    *string    `json:"remarks,omitempty"`
	ActorID    *uuid.UUID `json:"-"`
}
