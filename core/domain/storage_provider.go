package domain

import (
	"encoding/json"
	"errors"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	StorageProviderS3    = "s3"
	StorageProviderMinIO = "minio"
)

var (
	ErrStorageProviderSettingsNotFound = errors.New("storage provider settings not found")
	ErrInvalidStorageProvider          = errors.New("storage provider is invalid")
	ErrInvalidStorageBucket            = errors.New("storage bucket is required")
	ErrInvalidStorageEndpoint          = errors.New("storage endpoint is required")
	ErrInvalidStorageCredentials       = errors.New("storage credentials are required")
	ErrStorageProviderDisabled         = errors.New("storage provider is disabled")
	ErrStorageFileTooLarge             = errors.New("storage file exceeds configured size limit")
	ErrStorageContentTypeNotAllowed    = errors.New("storage content type is not allowed")
)

var bucketPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)

type StorageProviderSettings struct {
	ID                  uuid.UUID       `json:"id"`
	TenantID            uuid.UUID       `json:"tenant_id"`
	Provider            string          `json:"provider"`
	IsEnabled           bool            `json:"is_enabled"`
	Bucket              string          `json:"bucket"`
	Region              *string         `json:"region,omitempty"`
	Endpoint            *string         `json:"endpoint,omitempty"`
	AccessKeyID         *string         `json:"-"`
	SecretAccessKey     *string         `json:"-"`
	UseSSL              bool            `json:"use_ssl"`
	ForcePathStyle      bool            `json:"force_path_style"`
	ObjectPrefix        *string         `json:"object_prefix,omitempty"`
	PublicBaseURL       *string         `json:"public_base_url,omitempty"`
	MaxFileSizeBytes    int64           `json:"max_file_size_bytes"`
	AllowedContentTypes *string         `json:"allowed_content_types,omitempty"`
	LastTestAt          *time.Time      `json:"last_test_at,omitempty"`
	LastTestStatus      *string         `json:"last_test_status,omitempty"`
	LastTestMessage     *string         `json:"last_test_message,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`
	HasAccessKeyID      bool            `json:"has_access_key_id"`
	HasSecretAccessKey  bool            `json:"has_secret_access_key"`
	ReadinessHints      []string        `json:"readiness_hints,omitempty"`
	Inactive            bool            `json:"inactive"`
	CreatedAt           time.Time       `json:"created_at"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at"`
	UpdatedBy           *uuid.UUID      `json:"updated_by,omitempty"`
}

type StorageProviderSettingsInput struct {
	TenantID            uuid.UUID
	Provider            string
	IsEnabled           bool
	Bucket              string
	Region              *string
	Endpoint            *string
	AccessKeyID         *string
	SecretAccessKey     *string
	UseSSL              bool
	ForcePathStyle      bool
	ObjectPrefix        *string
	PublicBaseURL       *string
	MaxFileSizeBytes    int64
	AllowedContentTypes *string
	Metadata            json.RawMessage
}

func NewStorageProviderSettings(input StorageProviderSettingsInput) (*StorageProviderSettings, error) {
	if input.TenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	provider, err := ValidateStorageProvider(input.Provider)
	if err != nil {
		return nil, err
	}
	bucket := strings.ToLower(strings.TrimSpace(input.Bucket))
	if !bucketPattern.MatchString(bucket) {
		return nil, ErrInvalidStorageBucket
	}
	endpoint := cleanOptional(input.Endpoint)
	if provider == StorageProviderMinIO && input.IsEnabled && endpoint == nil {
		return nil, ErrInvalidStorageEndpoint
	}
	accessKey := cleanOptional(input.AccessKeyID)
	secretKey := cleanOptional(input.SecretAccessKey)
	if input.IsEnabled && (accessKey == nil || secretKey == nil) {
		return nil, ErrInvalidStorageCredentials
	}
	maxSize := input.MaxFileSizeBytes
	if maxSize <= 0 {
		maxSize = 10 * 1024 * 1024
	}
	metadata := input.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	now := time.Now().UTC()
	return &StorageProviderSettings{
		TenantID: input.TenantID, Provider: provider, IsEnabled: input.IsEnabled, Bucket: bucket, Region: cleanOptional(input.Region), Endpoint: endpoint, AccessKeyID: accessKey, SecretAccessKey: secretKey,
		UseSSL: input.UseSSL, ForcePathStyle: input.ForcePathStyle, ObjectPrefix: cleanStoragePrefix(input.ObjectPrefix), PublicBaseURL: cleanOptional(input.PublicBaseURL), MaxFileSizeBytes: maxSize, AllowedContentTypes: cleanOptional(input.AllowedContentTypes), Metadata: metadata,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

func ValidateStorageProvider(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", StorageProviderMinIO:
		return StorageProviderMinIO, nil
	case StorageProviderS3:
		return StorageProviderS3, nil
	default:
		return "", ErrInvalidStorageProvider
	}
}

func (s *StorageProviderSettings) RedactSecrets() *StorageProviderSettings {
	if s == nil {
		return nil
	}
	s.HasAccessKeyID = s.AccessKeyID != nil && strings.TrimSpace(*s.AccessKeyID) != ""
	s.HasSecretAccessKey = s.SecretAccessKey != nil && strings.TrimSpace(*s.SecretAccessKey) != ""
	s.AccessKeyID = nil
	s.SecretAccessKey = nil
	s.ReadinessHints = []string{"Use MinIO for local or self-hosted tenant storage.", "Use S3 for managed cloud object storage.", "Keep buckets private and serve downloads through signed URLs or backend-controlled paths.", "Use tenant prefixes to avoid object key collisions across modules."}
	return s
}

func StorageObjectKey(prefix *string, tenantID uuid.UUID, category string, ownerID uuid.UUID, entityID uuid.UUID, fileName string) string {
	parts := []string{}
	if prefix != nil && strings.TrimSpace(*prefix) != "" {
		parts = append(parts, strings.Trim(strings.TrimSpace(*prefix), "/"))
	}
	parts = append(parts, "tenants", tenantID.String(), strings.Trim(category, "/"))
	if ownerID != uuid.Nil {
		parts = append(parts, ownerID.String())
	}
	if entityID != uuid.Nil {
		parts = append(parts, entityID.String())
	}
	parts = append(parts, safeStorageFileName(fileName))
	return path.Join(parts...)
}

func cleanStoragePrefix(value *string) *string {
	clean := cleanOptional(value)
	if clean == nil {
		return nil
	}
	next := strings.Trim(strings.TrimSpace(*clean), "/")
	if next == "" {
		return nil
	}
	return &next
}

func safeStorageFileName(value string) string {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return "file"
	}
	clean = path.Base(clean)
	clean = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' || r == ':' || r == 0 {
			return '-'
		}
		return r
	}, clean)
	return clean
}
