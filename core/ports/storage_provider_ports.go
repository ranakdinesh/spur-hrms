package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranakdinesh/spur-hrms/core/domain"
)

const (
	StorageCategoryPolicyFile         = "policies"
	StorageCategoryEmployeeDoc        = "employee-documents"
	StorageCategoryEmployeeLetter     = "employee-letters"
	StorageCategoryAgreement          = "agreements"
	StorageCategoryHRCase             = "hr-cases"
	StorageCategoryLearningCert       = "learning-certificates"
	StorageCategorySalarySlip         = "salary-slips"
	StorageCategoryCandidateResume    = "candidate-resumes"
	StorageCategoryProfilePhoto       = "profile-photos"
	StorageCategoryWorkflowAttachment = "workflow-task-attachments"
	StorageCategoryBenefitClaim       = "benefit-claim-attachments"
	StorageCategoryEREvidence         = "employee-relations-evidence"
)

type StorageProviderRepo interface {
	GetStorageProviderSettings(ctx context.Context, tenantID uuid.UUID) (*domain.StorageProviderSettings, error)
	UpsertStorageProviderSettings(ctx context.Context, item *domain.StorageProviderSettings, actorID *uuid.UUID) (*domain.StorageProviderSettings, error)
	UpdateStorageProviderTestResult(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, message *string, actorID *uuid.UUID) (*domain.StorageProviderSettings, error)
	DeleteStorageProviderSettings(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, actorID *uuid.UUID) error
}

type ObjectStorage interface {
	PutObject(ctx context.Context, settings *domain.StorageProviderSettings, input StoreObjectInput) (string, error)
	TestStorage(ctx context.Context, settings *domain.StorageProviderSettings) error
}

type StoreObjectInput struct {
	TenantID    uuid.UUID
	Category    string
	OwnerID     uuid.UUID
	EntityID    uuid.UUID
	FileName    string
	ContentType string
	Content     []byte
}

type StorageProviderSettingsCommand struct {
	TenantID            uuid.UUID  `json:"tenant_id"`
	Provider            string     `json:"provider"`
	IsEnabled           bool       `json:"is_enabled"`
	Bucket              string     `json:"bucket"`
	Region              *string    `json:"region,omitempty"`
	Endpoint            *string    `json:"endpoint,omitempty"`
	AccessKeyID         *string    `json:"access_key_id,omitempty"`
	SecretAccessKey     *string    `json:"secret_access_key,omitempty"`
	UseSSL              bool       `json:"use_ssl"`
	ForcePathStyle      bool       `json:"force_path_style"`
	ObjectPrefix        *string    `json:"object_prefix,omitempty"`
	PublicBaseURL       *string    `json:"public_base_url,omitempty"`
	MaxFileSizeBytes    int64      `json:"max_file_size_bytes"`
	AllowedContentTypes *string    `json:"allowed_content_types,omitempty"`
	ActorID             *uuid.UUID `json:"-"`
}

type StorageProviderTestCommand struct {
	TenantID uuid.UUID  `json:"tenant_id"`
	ActorID  *uuid.UUID `json:"-"`
}
