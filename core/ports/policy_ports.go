package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

type PolicyRepo interface {
	CreatePolicyType(ctx context.Context, item *domain.PolicyType, actorID *uuid.UUID) (*domain.PolicyType, error)
	ListPolicyTypes(ctx context.Context, tenantID uuid.UUID) ([]*domain.PolicyType, error)
	GetPolicyType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.PolicyType, error)
	UpdatePolicyType(ctx context.Context, item *domain.PolicyType, actorID *uuid.UUID) (*domain.PolicyType, error)
	DeletePolicyType(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
	CreateCompanyPolicy(ctx context.Context, item *domain.CompanyPolicy, actorID *uuid.UUID) (*domain.CompanyPolicy, error)
	ListCompanyPolicies(ctx context.Context, tenantID uuid.UUID, policyTypeID *uuid.UUID) ([]*domain.CompanyPolicy, error)
	GetCompanyPolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.CompanyPolicy, error)
	UpdateCompanyPolicy(ctx context.Context, item *domain.CompanyPolicy, actorID *uuid.UUID) (*domain.CompanyPolicy, error)
	DeleteCompanyPolicy(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type PolicyFileStorage interface {
	StorePolicyFile(ctx context.Context, input StorePolicyFileInput) (string, error)
}

type PolicyNotifier interface {
	CompanyPolicyChanged(ctx context.Context, event CompanyPolicyChangedEvent) error
}

type StorePolicyFileInput struct {
	TenantID    uuid.UUID
	PolicyID    uuid.UUID
	FileName    string
	ContentType string
	Content     []byte
}

type CompanyPolicyChangedEvent struct {
	TenantID uuid.UUID
	PolicyID uuid.UUID
	Title    string
	Action   string
	ActorID  *uuid.UUID
}

type PolicyTypeCommand struct {
	ID       uuid.UUID  `json:"id,omitempty"`
	TenantID uuid.UUID  `json:"tenant_id"`
	Name     string     `json:"name"`
	IsSystem bool       `json:"is_system,omitempty"`
	ActorID  *uuid.UUID `json:"-"`
}

type CompanyPolicyCommand struct {
	ID                uuid.UUID  `json:"id,omitempty"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	PolicyTypeID      *uuid.UUID `json:"policy_type_id,omitempty"`
	Title             string     `json:"title"`
	FilePath          *string    `json:"file_path,omitempty"`
	Description       *string    `json:"description,omitempty"`
	FileName          string     `json:"file_name,omitempty"`
	FileContentType   string     `json:"file_content_type,omitempty"`
	FileContentBase64 string     `json:"file_content_base64,omitempty"`
	ActorID           *uuid.UUID `json:"-"`
}
